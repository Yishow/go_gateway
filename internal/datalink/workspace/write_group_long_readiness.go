package workspace

import (
	"context"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

// ManagedWriteGroupIndexInspector reads a managed table's explicit indexes at
// the connector revision a readiness snapshot captured. Only the long layout,
// which declares indexes, needs it; wide groups never call it.
type ManagedWriteGroupIndexInspector interface {
	InspectTableIndexesAtRevision(ctx context.Context, connectorID, expectedRevision, schemaName, tableName string) ([]dbtarget.IndexInfo, error)
}

// attachManagedIndexes adds real catalog index evidence to a long table's
// inspection. Without an inspector that can read indexes the layout cannot be
// confirmed, so readiness fails closed instead of assuming the index exists.
func (s *WriteGroupService) attachManagedIndexes(ctx context.Context, group *WriteGroup, inspection *dbtarget.TableInspection) error {
	reader, ok := s.managedTableInspector.(ManagedWriteGroupIndexInspector)
	if !ok {
		return ErrWriteGroupReadinessUnavailable
	}
	indexes, err := reader.InspectTableIndexesAtRevision(ctx,
		group.Destination.ConnectorID, group.Destination.ConnectorRevision,
		group.Destination.TableSchema, group.Destination.TableName)
	if err != nil {
		return err
	}
	inspection.Indexes, inspection.IndexesInspected = indexes, true
	return nil
}

// attachManagedHypertable adds the real catalog hypertable evidence to a
// timescale group's table inspection. Without an inspector that can read it the
// layout cannot be confirmed, so readiness fails closed rather than assuming the
// table is a hypertable.
func (s *WriteGroupService) attachManagedHypertable(ctx context.Context, group *WriteGroup, inspection *dbtarget.TableInspection) error {
	reader, ok := s.managedTableInspector.(ManagedWriteGroupHypertableInspector)
	if !ok {
		return ErrWriteGroupReadinessUnavailable
	}
	hypertable, err := reader.InspectHypertableAtRevision(ctx,
		group.Destination.ConnectorID, group.Destination.ConnectorRevision,
		group.Destination.TableSchema, group.Destination.TableName)
	if err != nil {
		return err
	}
	inspection.Hypertable, inspection.HypertableInspected = hypertable, true
	return nil
}

// ManagedWriteGroupHypertableInspector reads a table's hypertable evidence at
// the connector revision a readiness snapshot captured. Only a timescale group
// needs it.
type ManagedWriteGroupHypertableInspector interface {
	InspectHypertableAtRevision(ctx context.Context, connectorID, expectedRevision, schemaName, tableName string) (*dbtarget.HypertableInfo, error)
}

// The production inspector must keep reading hypertable evidence: readiness
// asserts the interface at run time and fails closed without it, so a silent
// loss of the method would block every timescale group.
var _ ManagedWriteGroupHypertableInspector = (*dbtarget.ManagedTableInspector)(nil)

// evaluateLongGroupInspection checks the shared long table. Sources are rows,
// not columns, so there is no per-member column check; the exact fixed layout,
// owner marker, composite key and index must all match the real catalog.
func evaluateLongGroupInspection(group *WriteGroup, tagTypes map[string]schema.DataType, inspection *dbtarget.TableInspection, dialect string) (bool, []ReadinessIssue) {
	layout, err := managedSchemaLayout(group, tagTypes, dialect)
	if err != nil {
		return false, []ReadinessIssue{writeGroupReadinessIssue(
			"tag-type-unverified", "a source tag type could not be verified", group.ID,
		)}
	}
	if err := recordingplan.ValidateManagedGroupTableInspection(*layout, recordingPlanTableInspection(inspection)); err != nil {
		return false, []ReadinessIssue{writeGroupReadinessIssue(
			"managed-schema-unverified", "managed destination ownership and canonical schema could not be verified", group.ID,
		)}
	}
	var issues []ReadinessIssue
	if inspection.Schema != group.Destination.TableSchema {
		issues = append(issues, writeGroupReadinessIssue(
			"schema-scope-mismatch", "destination schema does not match the saved binding", group.ID,
		))
	}
	if inspection.Table != group.Destination.TableName {
		issues = append(issues, writeGroupReadinessIssue(
			"schema-scope-mismatch", "destination table does not match the saved binding", group.ID,
		))
	}
	return len(issues) == 0, issues
}

// longRuntimeSupported reports whether the runtime really writes this long
// group: plain or explicitly selected timescale storage (the same typed table
// and writer; a timescale table is only a hypertable) with an instantaneous or
// cumulative profile, through the durable closure that also carries cumulative
// baselines. Any other storage mode or profile stays blocked, so an unsupported
// group never looks writable just because its schema was confirmed.
func longRuntimeSupported(group *WriteGroup) bool {
	if !group.isLongLayout() || group.RecordingProfile == nil ||
		group.StorageMode != WriteGroupStorageModePlain && group.StorageMode != WriteGroupStorageModeTimescale {
		return false
	}
	kind := group.RecordingProfile.Kind
	return kind == WriteGroupProfileInstantaneous || kind == WriteGroupProfileCumulative
}

// cumulativeWithoutMeasurements reports a cumulative group with a source that
// has no saved measurement. Counter semantics are never inferred from a tag name
// or a unit; the source must name the measurement that confirms them.
func cumulativeWithoutMeasurements(group *WriteGroup) bool {
	if group.RecordingProfile == nil || group.RecordingProfile.Kind != WriteGroupProfileCumulative {
		return false
	}
	for _, member := range group.Members {
		if member.MeasurementID == nil || *member.MeasurementID == "" {
			return true
		}
	}
	return false
}

// longRuntimeUnavailableIssue keeps a confirmed long schema from being read as
// a writer that is ready when the runtime does not write its storage mode or
// profile.
func longRuntimeUnavailableIssue(groupID string) ReadinessIssue {
	return writeGroupReadinessIssue(
		"long-runtime-unavailable",
		"long-table recording is not available at runtime for this storage mode or profile; the confirmed schema does not make this group writable",
		groupID,
	)
}

// counterSemanticsRequiredIssue blocks a cumulative group whose source has no
// saved measurement to confirm that it is a counter.
func counterSemanticsRequiredIssue(groupID string) ReadinessIssue {
	return writeGroupReadinessIssue(
		"counter-semantics-required",
		"every cumulative source needs a saved measurement that confirms counter or signed_counter semantics",
		groupID,
	)
}
