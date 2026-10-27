package grouptestwrite

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

const sqlServerTestTableOwner = "gateway"

type sqlServerTestPlan struct {
	proof        *dbtarget.SQLServerManagedRowProof
	schemaDigest string
	layout       *dbtarget.LongRowLayout
	members      []sqlServerTestMember
}

type sqlServerTestMember struct {
	member    workspace.WriteGroupMember
	key       string
	kind      measurement.ExactType
	value     measurement.ExactValue
	reduction string
}

type sqlServerTestWriteLayoutReader interface {
	SQLServerTestWriteLayout(ctx context.Context, groupID, expectedGroupRevision string) (*workspace.WriteGroupRuntimeLayout, error)
}

// validSQLServerTestWriteAttempt checks attempt identity and the saved scope
// without re-running the frozen initial-layout digest validator. The workspace
// adapter owns initial-proof and current-catalog/ledger reconciliation.
func validSQLServerTestWriteAttempt(group *workspace.WriteGroup, layout *workspace.WriteGroupRuntimeLayout) bool {
	if group == nil || layout == nil || layout.Dialect != dbtarget.SQLDialectSQLServer ||
		layout.SchemaDigest == "" || layout.SchemaDigest != group.Destination.SchemaDigest ||
		!layout.ReceiptTableReady || !workspace.LongRecordingShapeSupported(group, layout.TagTypes) {
		return false
	}
	table, ownerColumn, ok := workspace.LongStorageIdentity(group)
	if !ok || table != group.Destination.TableName {
		return false
	}
	proof := layout.SQLServer
	if proof == nil || proof.ConnectorID != group.Destination.ConnectorID || proof.ConnectorRevision != group.Destination.ConnectorRevision ||
		proof.Group == nil || proof.Group.Inspection == nil || proof.Group.Evidence == nil ||
		proof.Receipt == nil || proof.Receipt.Inspection == nil || proof.Receipt.Evidence == nil {
		return false
	}
	groupInspection, receiptInspection := proof.Group.Inspection, proof.Receipt.Inspection
	groupEvidence, receiptEvidence := proof.Group.Evidence, proof.Receipt.Evidence
	if groupInspection.Status != dbtarget.TableInspectionExists || groupInspection.Table != table ||
		groupInspection.Schema != group.Destination.TableSchema || receiptInspection.Status != dbtarget.TableInspectionExists ||
		receiptInspection.Table != dbtarget.EffectReceiptTable || receiptInspection.Schema != group.Destination.TableSchema ||
		proof.Group.Indexes == nil || proof.Receipt.Indexes == nil ||
		!slices.EqualFunc(proof.Group.Indexes, groupInspection.Indexes, sameSQLServerTestIndex) ||
		!slices.EqualFunc(proof.Receipt.Indexes, receiptInspection.Indexes, sameSQLServerTestIndex) ||
		!slices.Equal(layout.Columns, groupInspection.Columns) ||
		!groupEvidence.MetadataVerified || !receiptEvidence.MetadataVerified || groupEvidence.TableType != "USER_TABLE" ||
		receiptEvidence.TableType != "USER_TABLE" || groupEvidence.Database != group.Destination.Database ||
		receiptEvidence.Database != group.Destination.Database || groupEvidence.Schema != group.Destination.TableSchema ||
		receiptEvidence.Schema != group.Destination.TableSchema || groupEvidence.CapabilityRevision == "" ||
		receiptEvidence.CapabilityRevision != groupEvidence.CapabilityRevision || groupEvidence.ServerVersion == "" ||
		receiptEvidence.ServerVersion != groupEvidence.ServerVersion || groupEvidence.ServerEdition == "" ||
		receiptEvidence.ServerEdition != groupEvidence.ServerEdition || groupEvidence.EngineEdition <= 0 ||
		receiptEvidence.EngineEdition != groupEvidence.EngineEdition || (groupEvidence.CompatibilityLevel != 160 && groupEvidence.CompatibilityLevel != 170) ||
		receiptEvidence.CompatibilityLevel != groupEvidence.CompatibilityLevel || !receiptEvidence.Plain ||
		receiptEvidence.NativeMonthly != nil {
		return false
	}
	ownerDefault, ok := groupEvidence.ColumnDefaults[ownerColumn]
	if !ok || ownerDefault == nil || *ownerDefault != sqlServerTestTableOwner {
		return false
	}
	if group.StorageMode == workspace.WriteGroupStorageModeNativeMonthly {
		partition := groupEvidence.NativeMonthly
		return !groupEvidence.Plain && groupEvidence.PartitionInspected && partition != nil &&
			partition.FunctionExists && partition.SchemeExists && partition.Exclusive && partition.IndexesAligned &&
			partition.FunctionName != "" && partition.SchemeName != "" && len(partition.Partitions) >= 6
	}
	return group.StorageMode == workspace.WriteGroupStorageModePlain && groupEvidence.Plain && groupEvidence.NativeMonthly == nil
}

func sameSQLServerTestIndex(left, right dbtarget.IndexInfo) bool {
	return left.Name == right.Name && left.Unique == right.Unique && slices.Equal(left.Columns, right.Columns)
}

func (s *Service) buildSQLServerTestPlan(ctx context.Context, group *workspace.WriteGroup, connector *schema.DatabaseConnector) (*plan, error) {
	if group == nil || connector == nil || group.RecordingProfile == nil || len(group.Members) == 0 ||
		len(group.Members) > dbtarget.MaxGroupLongTestRows || group.LayoutVersion != workspace.WriteGroupLayoutManagedLongV1 ||
		(group.StorageMode != workspace.WriteGroupStorageModePlain && group.StorageMode != workspace.WriteGroupStorageModeNativeMonthly) ||
		group.Destination.StorageStrategy != workspace.WriteGroupStorageStrategyManaged ||
		!strings.EqualFold(strings.TrimSpace(group.WritePolicy.DedupeCapability), groupdelivery.DedupeReceipt) {
		return nil, unsupported(ReasonGroupNotTestable)
	}
	reader, ok := s.deps.Groups.(sqlServerTestWriteLayoutReader)
	if !ok {
		return nil, unsupported(ReasonLayoutBlocked)
	}
	runtimeLayout, err := reader.SQLServerTestWriteLayout(ctx, group.ID, group.Revision)
	if err != nil {
		if errors.Is(err, workspace.ErrWriteGroupRevisionConflict) {
			return nil, fmt.Errorf("%w: group revision changed", ErrDestinationChanged)
		}
		return nil, unsupported(ReasonLayoutBlocked)
	}
	if runtimeLayout == nil || runtimeLayout.Dialect != dbtarget.SQLDialectSQLServer || runtimeLayout.SQLServer == nil ||
		runtimeLayout.SchemaDigest == "" || !runtimeLayout.ReceiptTableReady {
		return nil, unsupported(ReasonLayoutBlocked)
	}
	if !validSQLServerTestWriteAttempt(group, runtimeLayout) {
		return nil, unsupported(ReasonLayoutBlocked)
	}
	ownerColumn := ""
	for _, column := range runtimeLayout.Columns {
		if strings.HasPrefix(column.Name, "_gw_owner_") {
			ownerColumn = column.Name
			break
		}
	}
	if ownerColumn == "" {
		return nil, unsupported(ReasonOwnershipUnsafe)
	}

	reduction := group.RecordingProfile.EffectiveSampleReduction()
	layoutMembers := make([]dbtarget.LongRowMember, 0, len(group.Members))
	fixtures := make([]sqlServerTestMember, 0, len(group.Members))
	previewMembers := make([]memberFixture, 0, len(group.Members))
	for _, member := range group.Members {
		tagType, ok := runtimeLayout.TagTypes[member.TagID]
		if !ok {
			tag, tagErr := s.deps.Tags.GetByID(ctx, member.TagID)
			if tagErr != nil || tag == nil {
				return nil, unsupported(ReasonTagUnavailable)
			}
			tagType = tag.DataType
		}
		kind, ok := workspace.LongMemberStoredTypeForMember(reduction, member, tagType)
		if !ok {
			return nil, unsupported(ReasonLayoutBlocked)
		}
		value, err := fixtureValue(kind)
		if err != nil {
			return nil, unsupported(ReasonTagTypeUnsupported)
		}
		key := memberKey(member.DeviceID, member.PointID, member.TagID)
		measurementID := ""
		if member.MeasurementID != nil {
			measurementID = strings.TrimSpace(*member.MeasurementID)
		}
		reductionLabel := string(reduction)
		if reduction == workspace.WriteGroupSampleReductionMean && member.ValueRole != workspace.WriteGroupValueRoleGauge {
			reductionLabel = ""
		}
		layoutMembers = append(layoutMembers, dbtarget.LongRowMember{
			MemberKey: key, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
			CircuitID: member.CircuitID, MeasurementID: measurementID, Type: kind, Reduction: reductionLabel,
			ValueRole: string(member.ValueRole), SourceValueType: member.SourceValueType,
		})
		fixtures = append(fixtures, sqlServerTestMember{member: member, key: key, kind: kind, value: value, reduction: reductionLabel})
		previewMembers = append(previewMembers, memberFixture{column: sqlServerValueSlot(kind), kind: kind, value: value, key: key})
	}
	layout, issues := dbtarget.NewLongRowLayout(dbtarget.LongRowSpec{
		Dialect: dbtarget.SQLDialectSQLServer, Columns: runtimeLayout.Columns,
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision,
		ProfileKind: string(group.RecordingProfile.Kind), ProfileInterval: group.RecordingProfile.Interval,
		ProfileUnit: string(group.RecordingProfile.Unit), Members: layoutMembers,
	})
	if len(issues) > 0 {
		return nil, unsupported(ReasonLayoutBlocked, issues...)
	}
	scope, err := snapshot.DestinationScope(group.Destination.ConnectorID, group.Destination.ConnectorRevision,
		group.Destination.Database, group.Destination.TableSchema, group.Destination.TableName)
	if err != nil {
		return nil, unsupported(ReasonGroupNotTestable)
	}
	p := &plan{
		group: group, connector: connector, dialect: dbtarget.SQLDialectSQLServer,
		columns: runtimeLayout.Columns, members: previewMembers, ownerColumn: ownerColumn,
		strategy: dbtarget.GroupEffectReceipt, scope: scope,
		sqlserver: &sqlServerTestPlan{proof: runtimeLayout.SQLServer, schemaDigest: runtimeLayout.SchemaDigest, layout: layout, members: fixtures},
	}
	if p.digest, err = p.contentDigest(); err != nil {
		return nil, err
	}
	if _, err := p.sqlserverRows(ownerValue("op-preview"), time.Unix(0, 0).UTC()); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *plan) sqlserverRows(owner string, at time.Time) ([]dbtarget.EncodedRow, error) {
	if p == nil || p.group == nil || p.group.RecordingProfile == nil || p.group.RecordingProfile.IntervalSeconds() <= 0 ||
		p.sqlserver == nil || p.sqlserver.layout == nil || len(p.sqlserver.members) == 0 {
		return nil, unsupported(ReasonGroupNotTestable)
	}
	at = at.UTC().Truncate(time.Microsecond)
	interval := time.Duration(p.group.RecordingProfile.IntervalSeconds()) * time.Second
	end := at.Add(interval).UTC().Truncate(time.Microsecond)
	rows := make([]dbtarget.EncodedRow, 0, len(p.sqlserver.members))
	for index, member := range p.sqlserver.members {
		recordID, err := snapshot.RecordID(p.group.WorkspaceID, p.group.ID, p.group.Revision, owner+"|"+member.key, at)
		if err != nil {
			return nil, unsupported(ReasonGroupNotTestable)
		}
		effectKey, err := snapshot.EffectKey(p.scope, recordID)
		if err != nil {
			return nil, unsupported(ReasonGroupNotTestable)
		}
		result := snapshot.MemberResult{
			MemberKey: member.key, Status: snapshot.MemberOK,
			Sample: &snapshot.Sample{
				SampleID: fmt.Sprintf("%s-%d", owner, index), MemberKey: member.key,
				ObservedAt: at, ReceivedAt: at, Quality: schema.QualityGood, Value: member.value,
			},
		}
		if member.reduction == string(workspace.WriteGroupSampleReductionMean) && member.member.ValueRole == workspace.WriteGroupValueRoleGauge {
			mean, err := measurement.Mean([]measurement.ExactValue{member.value})
			if err != nil {
				return nil, unsupported(ReasonSQLValueBlocked)
			}
			result.Reduction = &mean
		}
		outcome := snapshot.Outcome{
			Kind: snapshot.OutcomeRow, EntityKey: owner, RecordID: recordID, EffectKey: effectKey,
			BucketStart: at, BucketEnd: end, Members: []snapshot.MemberResult{result},
		}
		row, err := p.sqlserver.layout.EncodeRow(outcome)
		if err != nil {
			return nil, unsupported(ReasonSQLValueBlocked)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func sqlServerValueSlot(kind measurement.ExactType) string {
	switch kind {
	case measurement.ExactInt64:
		return "value_int"
	case measurement.ExactUint64, measurement.ExactDecimal:
		return "value_decimal"
	case measurement.ExactFloat64:
		return "value_float"
	case measurement.ExactBool:
		return "value_bool"
	case measurement.ExactText:
		return "value_text"
	}
	return ""
}

func sqlServerPreview(p *plan, owner string, at time.Time) (values []PreviewValue, recordIDs []string, cleanup string, err error) {
	rows, err := p.sqlserverRows(owner, at)
	if err != nil {
		return nil, nil, "", err
	}
	if len(rows) == 0 {
		return nil, nil, "", unsupported(ReasonGroupNotTestable)
	}
	recordIDs = make([]string, len(rows))
	for index, row := range rows {
		recordIDs[index] = row.RecordID
	}
	types := make(map[string]string, len(p.columns))
	for _, column := range p.columns {
		types[strings.ToLower(column.Name)] = column.DataType
	}
	values = make([]PreviewValue, 0, len(rows[0].Cells))
	for _, cell := range rows[0].Cells {
		values = append(values, PreviewValue{Column: cell.Column, Type: previewType(types[strings.ToLower(cell.Column)]), Value: previewValue(cell.Value)})
	}
	cleanup = "remove only the previewed record IDs and their matching receipts; retain all other rows"
	return values, recordIDs, cleanup, nil
}
