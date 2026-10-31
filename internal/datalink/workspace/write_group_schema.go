package workspace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

const managedSQLiteSchemaName = "main"

// managedTextColumnType is the exact-text column type: text itself, a SQLite
// decimal and a SQLite unsigned 64-bit integer, none of which pass through a float.
const managedTextColumnType = "TEXT"

// ManagedSchemaScope resolves a preview entirely from the canonical saved
// group and live source revisions. It performs no destination inspection.
func (s *WriteGroupService) ManagedSchemaScope(ctx context.Context, id string, mutation WriteGroupMutation) (recordingplan.SchemaPreviewScope, error) {
	if err := s.validate(); err != nil {
		return recordingplan.SchemaPreviewScope{}, err
	}
	snapshot, err := s.readinessSnapshot(ctx, id, nil)
	if err != nil {
		return recordingplan.SchemaPreviewScope{}, err
	}
	if snapshot.validateErr != nil {
		return recordingplan.SchemaPreviewScope{}, snapshot.validateErr
	}
	group := snapshot.validated
	if snapshot.record.ID != mutation.WorkspaceID || group.Status == WriteGroupStatusDeleted {
		return recordingplan.SchemaPreviewScope{}, writeGroupNotFound("prepare group schema")
	}
	if snapshot.record.DatabaseSetupRevision != mutation.ExpectedWorkspaceRevision || group.Revision != mutation.ExpectedGroupRevision || group.Destination.ConnectorRevision != mutation.ExpectedConnectorRevision {
		return recordingplan.SchemaPreviewScope{}, writeGroupRevisionConflict("managed schema scope is stale")
	}
	if group.Destination.StorageStrategy != WriteGroupStorageStrategyManaged || group.RowPolicy.RecordKeyColumn == "" {
		return recordingplan.SchemaPreviewScope{}, fmt.Errorf("managed schema requires a saved managed layout: %w", ErrWriteGroupValidation)
	}
	if snapshot.connectorKind == writeGroupConnectorKindSQLite && group.Destination.TableSchema != managedSQLiteSchemaName {
		return recordingplan.SchemaPreviewScope{}, fmt.Errorf("managed SQLite schema must be main: %w", ErrWriteGroupValidation)
	}
	// Mixed intent uses only the existing supported result-layout authority.
	// Check before legacy dialect branches, which may probe the destination.
	if isMixedWriteGroup(group) {
		if _, err := managedSchemaLayout(group, snapshot.tagTypes, snapshot.connectorKind); err != nil {
			return recordingplan.SchemaPreviewScope{}, err
		}
	}
	var layout *recordingplan.GroupSchemaLayout
	switch snapshot.connectorKind {
	case writeGroupConnectorKindSQLServer:
		var err error
		layout, err = s.managedSQLServerSchemaLayout(ctx, group)
		if err != nil {
			return recordingplan.SchemaPreviewScope{}, err
		}
	case writeGroupConnectorKindMySQL:
		if !group.isLongLayout() || group.StorageMode != WriteGroupStorageModePlain && group.StorageMode != WriteGroupStorageModeNativeMonthly {
			return recordingplan.SchemaPreviewScope{}, fmt.Errorf("%w: mysql requires plain or native monthly managed long schema", ErrWriteGroupValidation)
		}
		var err error
		layout, err = s.managedMySQLSchemaLayout(ctx, group)
		if err != nil {
			return recordingplan.SchemaPreviewScope{}, err
		}
	default:
		var err error
		layout, err = managedSchemaLayout(group, snapshot.tagTypes, snapshot.connectorKind)
		if err != nil {
			return recordingplan.SchemaPreviewScope{}, err
		}
	}
	sourceDigest, err := managedSourceDigest(group, snapshot.tagTypes, snapshot.counterBindings)
	if err != nil {
		return recordingplan.SchemaPreviewScope{}, err
	}
	scope := recordingplan.SchemaPreviewScope{
		WorkspaceID: snapshot.record.ID, WorkspaceRevision: snapshot.record.DatabaseSetupRevision,
		PlanID: group.ID, PlanRevision: group.Revision, ConnectorID: group.Destination.ConnectorID,
		ConnectorRevision: group.Destination.ConnectorRevision, Dialect: snapshot.connectorKind,
		Database: group.Destination.Database, Schema: group.Destination.TableSchema, TablePrefix: group.Destination.TableName,
		GroupLayout: layout, SourceDigest: sourceDigest,
		SchemaRevision: group.Destination.SchemaRevision, SchemaDigest: group.Destination.SchemaDigest,
	}
	if managedLongGroup(group) && group.WritePolicy.RecreateMissing != nil {
		// The confirmation shows, and binds, the saved recovery policy. A group
		// that never confirmed one shows none and so grants none.
		policy := *group.WritePolicy.RecreateMissing
		scope.RecreateMissing = &policy
	}
	return scope, nil
}

// managedSourceDigest binds the sources a schema confirmation was made for. A
// long group also binds its profile, layout and mode, which are part of the
// confirmed scope; a wide group keeps the original two-part digest. Recording
// start records the same digest at each checkpoint, so both must come from here.
func managedSourceDigest(group *WriteGroup, tagTypes map[string]schema.DataType, bindings ...CounterBindingProofs) (string, error) {
	sources := []any{group.Members, tagTypes}
	if isMixedWriteGroup(group) {
		var supplied CounterBindingProofs
		if len(bindings) > 0 {
			supplied = bindings[0]
		}
		proofs, err := mixedSchemaCounterProofs(group, supplied)
		if err != nil {
			return "", fmt.Errorf("%w: %w", ErrWriteGroupValidation, err)
		}
		sources = append(sources, proofs)
		sources = append(sources, group.WorkspaceID, group.ID, group.RecordID, group.RowPolicy, group.RecordingProfile, group.LayoutVersion, group.StorageMode)
	}
	switch group.EffectiveLayoutVersion() {
	case WriteGroupLayoutManagedLongV1:
		sources = append(sources, canonicalDigestProfile(group.RecordingProfile), group.LayoutVersion, group.StorageMode)
	case WriteGroupLayoutManagedWideProfileV1:
		// The new wide profile binds its profile and layout beside the member
		// metadata; the legacy wide digest never carried them.
		sources = append(sources, canonicalDigestProfile(group.RecordingProfile), group.LayoutVersion)
	}
	encoded, err := json.Marshal(sources)
	if err != nil {
		return "", fmt.Errorf("encode managed source scope: %w", err)
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

// canonicalDigestProfile pins the digest to the effective reduction: an
// omitted reduction and an explicitly pinned latest are the same confirmed
// scope, so both marshal without the key and a no-op edit never invalidates a
// confirmation. A mean keeps the key and forces a fresh confirmation.
func canonicalDigestProfile(profile *WriteGroupRecordingProfile) *WriteGroupRecordingProfile {
	if profile == nil {
		return nil
	}
	canonical := *profile
	if canonical.EffectiveSampleReduction() == WriteGroupSampleReductionLatest {
		canonical.SampleReduction = ""
	}
	return &canonical
}

func managedSchemaLayout(group *WriteGroup, types map[string]schema.DataType, dialect string) (*recordingplan.GroupSchemaLayout, error) {
	if isMixedWriteGroup(group) && dialect != writeGroupConnectorKindSQLite && dialect != writeGroupConnectorKindPostgres {
		return nil, fmt.Errorf("%w: mixed schema dialect unsupported", ErrWriteGroupValidation)
	}
	if isMixedWriteGroup(group) && !mixedSchemaShapeSupported(group, types) {
		return nil, fmt.Errorf("%w: mixed schema requires complete frozen calculation facts and explicit dispatch mapping", ErrWriteGroupValidation)
	}
	if group.isWideProfile() && !isMixedWriteGroup(group) && !WideProfileShapeSupported(group, types) {
		// The wide profile layout is prepared only from frozen result types that
		// are complete and verified; anything else is never given the legacy wide
		// DDL as a stand-in.
		return nil, fmt.Errorf("%w: the %s layout needs complete, verified recording facts",
			ErrWriteGroupValidation, WriteGroupLayoutManagedWideProfileV1)
	}
	if group.isLongLayout() {
		for _, member := range group.Members {
			if _, ok := measurement.ExactTypeForTag(types[member.TagID]); !ok {
				return nil, fmt.Errorf("managed exact type is unsupported: %w", ErrWriteGroupValidation)
			}
		}
		return longStorageScopeOf(group).layout(dialect)
	}
	owner := sha256.Sum256([]byte(group.WorkspaceID + "\x00" + group.ID))
	layout := &recordingplan.GroupSchemaLayout{TableName: group.Destination.TableName, OwnerColumn: "_gw_owner_" + hex.EncodeToString(owner[:12])}
	add := func(name, sqlType string, nullable, primary bool) {
		if name != "" {
			layout.Columns = append(layout.Columns, recordingplan.GroupSchemaColumn{Name: name, SQLType: sqlType, Nullable: nullable, PrimaryKey: primary})
		}
	}
	add(group.RowPolicy.RecordKeyColumn, "TEXT", false, true)
	add(group.RowPolicy.GroupIDColumn, "TEXT", false, false)
	add(group.RowPolicy.DeviceIDColumn, "TEXT", false, false)
	add(group.RowPolicy.BucketStartColumn, "TEXT", false, false)
	add(group.RowPolicy.DueAtColumn, "TEXT", false, false)
	add(group.RowPolicy.ProvenanceColumn, "TEXT", false, false)
	add(group.RowPolicy.EntityKeyColumn, "TEXT", false, false)
	add(layout.OwnerColumn, "TEXT", false, false)
	for _, member := range group.Members {
		exact, ok := memberExactType(group, member, types)
		if !ok {
			return nil, fmt.Errorf("managed exact type is unsupported: %w", ErrWriteGroupValidation)
		}
		sqlType, ok := managedSQLType(exact, dialect)
		if !ok {
			return nil, fmt.Errorf("managed exact type is unsupported: %w", ErrWriteGroupValidation)
		}
		add(member.TargetColumn, sqlType, true, false)
	}
	return layout, nil
}

// managedSQLType maps the exact type a column holds to its managed SQL type. A
// decimal is exact text on SQLite and an unbounded NUMERIC on PostgreSQL, so a
// mean never passes through a float column.
func managedSQLType(exact measurement.ExactType, dialect string) (string, bool) {
	postgres := dialect == writeGroupConnectorKindPostgres
	switch exact {
	case measurement.ExactDecimal:
		if postgres {
			return "NUMERIC", true
		}
		return managedTextColumnType, true
	case measurement.ExactText:
		return managedTextColumnType, true
	case measurement.ExactBool:
		if postgres {
			return "BOOLEAN", true
		}
		return "INTEGER", true
	case measurement.ExactInt64:
		if postgres {
			return "BIGINT", true
		}
		return "INTEGER", true
	case measurement.ExactUint64:
		if postgres {
			return "NUMERIC(20,0)", true
		}
		return managedTextColumnType, true
	case measurement.ExactFloat64:
		if postgres {
			return "DOUBLE PRECISION", true
		}
		return "REAL", true
	}
	return "", false
}
