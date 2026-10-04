package grouptestwrite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

// Safe reasons a group cannot be test-written. They never carry values.
const (
	ReasonDestinationMissing  = "destination-missing"
	ReasonDestinationRevision = "destination-revision-changed"
	ReasonDestinationKind     = "destination-unsupported"
	ReasonTableUnavailable    = "table-unavailable"
	ReasonTableMissing        = "table-missing"
	ReasonLayoutBlocked       = "layout-blocked"
	ReasonOwnershipUnsafe     = "test-ownership-unsupported"
	ReasonReceiptTableMissing = "receipt-table-missing"
	ReasonTagUnavailable      = "tag-unavailable"
	ReasonTagTypeUnsupported  = "tag-type-unsupported"
	ReasonSQLValueBlocked     = "test-value-blocked"
	ReasonGroupNotTestable    = "group-not-testable"
)

// UnsupportedError says why a group cannot be test-written, without mutating
// anything. Issues carry only layout issue codes and column names.
type UnsupportedError struct {
	Reason string
	Issues []dbtarget.LayoutIssue
}

func (e *UnsupportedError) Error() string { return "test write unsupported: " + e.Reason }

// ErrGroupNotFound means the group does not exist in the workspace.
var ErrGroupNotFound = errors.New("write group not found")

// ErrDestinationChanged means the saved destination no longer matches the revision
// the preview or the group was bound to.
var ErrDestinationChanged = errors.New("write group destination changed")

type memberFixture struct {
	column string
	kind   measurement.ExactType
	value  measurement.ExactValue
	key    string
}

// plan is everything a test write needs, derived only from the saved group,
// the saved connector and the real inspected table. It is rebuilt at
// confirmation and compared to the preview by digest.
type plan struct {
	group       *workspace.WriteGroup
	connector   *schema.DatabaseConnector
	dialect     dbtarget.SQLDialect
	columns     []dbtarget.ColumnInfo
	layout      *dbtarget.GroupRowLayout
	members     []memberFixture
	ownerColumn string
	strategy    string
	digest      string
	scope       string
}

// fixtureValue is a fixed, recognizable value per exact type. It deliberately
// includes a number that a float64 cannot hold, so a lossy path is detected.
func fixtureValue(kind measurement.ExactType) (measurement.ExactValue, error) {
	switch kind {
	case measurement.ExactBool:
		return measurement.NewBool(true), nil
	case measurement.ExactText:
		return measurement.NewText("gw-test-value"), nil
	case measurement.ExactInt64:
		return measurement.NewInt64(-4242), nil
	case measurement.ExactUint64:
		return measurement.NewUint64(9007199254740993), nil
	case measurement.ExactFloat64:
		return measurement.NewFloat64(1.5)
	}
	return measurement.ExactValue{}, fmt.Errorf("%w: unsupported type", measurement.ErrExactValueInvalid)
}

func memberKey(deviceID, pointID, tagID string) string {
	encoded, err := json.Marshal([]string{deviceID, pointID, tagID})
	if err != nil {
		return deviceID + "\x1f" + pointID + "\x1f" + tagID
	}
	return string(encoded)
}

func groupDeviceID(group *workspace.WriteGroup) string {
	if len(group.Members) == 0 {
		return ""
	}
	deviceID := group.Members[0].DeviceID
	for _, member := range group.Members[1:] {
		if member.DeviceID != deviceID {
			return ""
		}
	}
	return deviceID
}

func unsupported(reason string, issues ...dbtarget.LayoutIssue) error {
	return &UnsupportedError{Reason: reason, Issues: issues}
}

// buildPlan resolves the saved group, destination and real table into a plan.
func (s *Service) buildPlan(ctx context.Context, groupID string) (*plan, error) {
	saved, err := s.deps.Groups.Get(ctx, groupID)
	if err != nil {
		if errors.Is(err, workspace.ErrWriteGroupNotFound) {
			return nil, ErrGroupNotFound
		}
		return nil, fmt.Errorf("read write group: %w", err)
	}
	group := saved.Group
	if group == nil || group.Status == workspace.WriteGroupStatusDeleted {
		return nil, ErrGroupNotFound
	}
	connector, err := s.deps.Destinations.GetByID(ctx, group.Destination.ConnectorID)
	switch {
	case errors.Is(err, dbtarget.ErrConnectorNotFound):
		return nil, unsupported(ReasonDestinationMissing)
	case err != nil:
		return nil, fmt.Errorf("read destination connector: %w", err)
	case connector.IdentityRevision != group.Destination.ConnectorRevision:
		return nil, fmt.Errorf("%w: connector identity changed", ErrDestinationChanged)
	}
	dialect, ok := dialectFor(connector.Kind)
	if !ok {
		return nil, unsupported(ReasonDestinationKind)
	}
	inspection, err := s.deps.Inspector.InspectTable(ctx, group.Destination.ConnectorID, group.Destination.TableSchema, group.Destination.TableName)
	switch {
	case err != nil || inspection == nil:
		return nil, unsupported(ReasonTableUnavailable)
	case inspection.Status == dbtarget.TableInspectionMissing:
		return nil, unsupported(ReasonTableMissing)
	case inspection.Status != dbtarget.TableInspectionExists:
		return nil, unsupported(ReasonTableUnavailable)
	}
	// Ownership must be provable by a value in the row itself. Without an
	// entity key column no column can carry an operation-owned marker, and the
	// only cleanup left would be a guess.
	ownerColumn := strings.TrimSpace(group.RowPolicy.EntityKeyColumn)
	if ownerColumn == "" {
		return nil, unsupported(ReasonOwnershipUnsafe)
	}

	members := make([]memberFixture, 0, len(group.Members))
	layoutMembers := make([]dbtarget.GroupRowMember, 0, len(group.Members))
	for _, member := range group.Members {
		tag, err := s.deps.Tags.GetByID(ctx, member.TagID)
		if err != nil || tag == nil {
			return nil, unsupported(ReasonTagUnavailable)
		}
		kind, ok := measurement.ExactTypeForTag(tag.DataType)
		if !ok {
			return nil, unsupported(ReasonTagTypeUnsupported)
		}
		value, err := fixtureValue(kind)
		if err != nil {
			return nil, unsupported(ReasonTagTypeUnsupported)
		}
		key := memberKey(member.DeviceID, member.PointID, member.TagID)
		members = append(members, memberFixture{column: member.TargetColumn, kind: kind, value: value, key: key})
		layoutMembers = append(layoutMembers, dbtarget.GroupRowMember{MemberKey: key, Column: member.TargetColumn, Type: kind, Required: true})
	}
	deviceID := groupDeviceID(group)
	layout, issues := dbtarget.NewGroupRowLayout(dbtarget.GroupRowSpec{
		Dialect: dialect, Columns: inspection.Columns, Members: layoutMembers,
		EntityKeyed: true, EntityKeyColumn: ownerColumn,
		RecordKeyColumn: group.RowPolicy.RecordKeyColumn, BucketStartColumn: group.RowPolicy.BucketStartColumn,
		ProvenanceColumn: group.RowPolicy.ProvenanceColumn,
		GroupIDColumn:    group.RowPolicy.GroupIDColumn, GroupID: group.ID,
		DeviceIDColumn: group.RowPolicy.DeviceIDColumn, DeviceID: deviceID,
	})
	if len(issues) > 0 {
		return nil, unsupported(ReasonLayoutBlocked, issues...)
	}
	strategy := dbtarget.GroupEffectNone
	if strings.EqualFold(strings.TrimSpace(group.WritePolicy.DedupeCapability), groupdelivery.DedupeReceipt) {
		receipts, err := s.deps.Inspector.InspectTable(ctx, group.Destination.ConnectorID, group.Destination.TableSchema, dbtarget.EffectReceiptTable)
		if err != nil || receipts == nil || receipts.Status != dbtarget.TableInspectionExists {
			return nil, unsupported(ReasonReceiptTableMissing)
		}
		strategy = dbtarget.GroupEffectReceipt
	}
	scope, err := snapshot.DestinationScope(group.Destination.ConnectorID, group.Destination.ConnectorRevision,
		group.Destination.Database, group.Destination.TableSchema, group.Destination.TableName)
	if err != nil {
		return nil, unsupported(ReasonGroupNotTestable)
	}
	p := &plan{
		group: group, connector: connector, dialect: dialect, columns: inspection.Columns, layout: layout,
		members: members, ownerColumn: ownerColumn, strategy: strategy, scope: scope,
	}
	if p.digest, err = p.contentDigest(); err != nil {
		return nil, err
	}
	// A value the real column cannot hold is found now, at preview, not at write.
	if _, err := p.row("gw-test-00000000-0000-0000-0000-000000000000", time.Unix(0, 0).UTC()); err != nil {
		return nil, err
	}
	return p, nil
}

func dialectFor(kind schema.DatabaseConnectorKind) (dbtarget.SQLDialect, bool) {
	switch kind {
	case schema.DatabaseConnectorKindSQLite:
		return dbtarget.SQLDialectSQLite, true
	case schema.DatabaseConnectorKindPostgres:
		return dbtarget.SQLDialectPostgres, true
	}
	return "", false
}

// contentDigest binds what a confirmation will write apart from the owner
// value and time, which are derived from the operation: target table, columns,
// types, fixture values, ownership column and dedupe strategy.
func (p *plan) contentDigest() (string, error) {
	type member struct {
		Column string `json:"column"`
		Type   string `json:"type"`
		Value  any    `json:"value"`
	}
	type rowIdentity struct {
		EntityKeyColumn   string `json:"entity_key_column"`
		RecordKeyColumn   string `json:"record_key_column"`
		BucketStartColumn string `json:"bucket_start_column"`
		ProvenanceColumn  string `json:"provenance_column"`
		GroupIDColumn     string `json:"group_id_column"`
		GroupID           string `json:"group_id"`
		DeviceIDColumn    string `json:"device_id_column"`
		DeviceID          string `json:"device_id"`
	}
	list := make([]member, 0, len(p.members))
	for _, m := range p.members {
		encoded, err := json.Marshal(m.value)
		if err != nil {
			return "", fmt.Errorf("encode test value: %w", err)
		}
		list = append(list, member{Column: m.column, Type: string(m.kind), Value: json.RawMessage(encoded)})
	}
	parts := []any{
		"test-write-content-v1", p.group.Destination.ConnectorID, p.group.Destination.TableSchema, p.group.Destination.TableName,
		p.ownerColumn, p.group.RowPolicy.ProvenanceColumn, p.strategy,
	}
	identity := rowIdentity{
		EntityKeyColumn: p.ownerColumn, RecordKeyColumn: p.group.RowPolicy.RecordKeyColumn,
		BucketStartColumn: p.group.RowPolicy.BucketStartColumn, ProvenanceColumn: p.group.RowPolicy.ProvenanceColumn,
		GroupIDColumn: p.group.RowPolicy.GroupIDColumn, GroupID: p.group.ID,
		DeviceIDColumn: p.group.RowPolicy.DeviceIDColumn, DeviceID: groupDeviceID(p.group),
	}
	if identity.RecordKeyColumn != "" || identity.BucketStartColumn != "" || identity.GroupIDColumn != "" || identity.DeviceIDColumn != "" {
		parts = append(parts, identity)
	}
	parts = append(parts, list)
	payload, err := json.Marshal(parts)
	if err != nil {
		return "", fmt.Errorf("encode test content: %w", err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

// ownerValue is the operation-owned marker written to the entity key column.
func ownerValue(operationID string) string {
	return dbtarget.TestRowOwnerPrefix + strings.TrimPrefix(operationID, "op-")
}

// row encodes the test row for one operation through the production layout.
func (p *plan) row(owner string, at time.Time) (dbtarget.EncodedRow, error) {
	at = at.UTC()
	recordID, err := snapshot.RecordID(p.group.WorkspaceID, p.group.ID, p.group.Revision, owner, at)
	if err != nil {
		return dbtarget.EncodedRow{}, unsupported(ReasonGroupNotTestable)
	}
	effectKey, err := snapshot.EffectKey(p.scope, recordID)
	if err != nil {
		return dbtarget.EncodedRow{}, unsupported(ReasonGroupNotTestable)
	}
	results := make([]snapshot.MemberResult, 0, len(p.members))
	for i, m := range p.members {
		results = append(results, snapshot.MemberResult{
			MemberKey: m.key, Status: snapshot.MemberOK,
			Sample: &snapshot.Sample{
				SampleID: fmt.Sprintf("%s-%d", owner, i), MemberKey: m.key, ObservedAt: at,
				Quality: schema.QualityGood, Value: m.value,
			},
		})
	}
	outcome := snapshot.Outcome{
		Kind: snapshot.OutcomeRow, EntityKey: owner, RecordID: recordID, EffectKey: effectKey,
		BucketStart: at, BucketEnd: at, Members: results,
	}
	row, err := p.layout.EncodeRow(outcome)
	if err != nil {
		var blocked *dbtarget.GroupRowError
		if errors.As(err, &blocked) {
			return dbtarget.EncodedRow{}, unsupported(ReasonSQLValueBlocked,
				dbtarget.LayoutIssue{Code: blocked.Code, Column: blocked.Column})
		}
		return dbtarget.EncodedRow{}, unsupported(ReasonSQLValueBlocked)
	}
	return row, nil
}

func (p *plan) ownedRef(owner string) dbtarget.OwnedRowRef {
	return dbtarget.OwnedRowRef{
		Kind: p.connector.Kind, SchemaName: p.group.Destination.TableSchema, TableName: p.group.Destination.TableName,
		OwnerColumn: p.ownerColumn, OwnerValue: owner,
	}
}
