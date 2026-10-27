package grouptestwrite

import (
	"context"
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

type mysqlTestPlan struct {
	proof        *dbtarget.MySQLManagedRowProof
	schemaDigest string
	layout       *dbtarget.LongRowLayout
	members      []mysqlTestMember
	columns      []string
}

type mysqlTestMember struct {
	member    workspace.WriteGroupMember
	key       string
	kind      measurement.ExactType
	value     measurement.ExactValue
	reduction string
}

type mysqlTestWriteLayoutReader interface {
	MySQLTestWriteLayout(ctx context.Context, groupID, expectedGroupRevision string) (*workspace.WriteGroupRuntimeLayout, error)
}

func (s *Service) buildMySQLTestPlan(ctx context.Context, group *workspace.WriteGroup, connector *schema.DatabaseConnector) (*plan, error) {
	if group == nil || connector == nil || group.RecordingProfile == nil ||
		len(group.Members) == 0 || len(group.Members) > dbtarget.MaxGroupLongTestRows ||
		group.LayoutVersion != workspace.WriteGroupLayoutManagedLongV1 ||
		(group.StorageMode != workspace.WriteGroupStorageModePlain && group.StorageMode != workspace.WriteGroupStorageModeNativeMonthly) ||
		group.Destination.StorageStrategy != workspace.WriteGroupStorageStrategyManaged ||
		!strings.EqualFold(strings.TrimSpace(group.WritePolicy.DedupeCapability), groupdelivery.DedupeReceipt) {
		return nil, unsupported(ReasonGroupNotTestable)
	}
	reader, ok := s.deps.Groups.(mysqlTestWriteLayoutReader)
	if !ok {
		return nil, unsupported(ReasonLayoutBlocked)
	}
	runtimeLayout, err := reader.MySQLTestWriteLayout(ctx, group.ID, group.Revision)
	if err != nil {
		if errors.Is(err, workspace.ErrWriteGroupRevisionConflict) {
			return nil, fmt.Errorf("%w: group revision changed", ErrDestinationChanged)
		}
		return nil, unsupported(ReasonLayoutBlocked)
	}
	if err := validateMySQLTestWriteLayout(group, runtimeLayout); err != nil {
		return nil, err
	}
	proof := runtimeLayout.MySQL
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
	fixtures := make([]mysqlTestMember, 0, len(group.Members))
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
		layoutMembers = append(layoutMembers, dbtarget.LongRowMember{
			MemberKey: key, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID,
			CircuitID: member.CircuitID, MeasurementID: measurementID, Type: kind, Reduction: reductionLabel,
			ValueRole: string(member.ValueRole), SourceValueType: member.SourceValueType,
		})
		fixtures = append(fixtures, mysqlTestMember{member: member, key: key, kind: kind, value: value, reduction: reductionLabel})
		previewMembers = append(previewMembers, memberFixture{column: mysqlValueSlot(kind), kind: kind, value: value, key: key})
	}
	layout, issues := dbtarget.NewLongRowLayout(dbtarget.LongRowSpec{
		Dialect: dbtarget.SQLDialectMySQL, Columns: runtimeLayout.Columns,
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
	columns := make([]string, len(runtimeLayout.Columns))
	for i, column := range runtimeLayout.Columns {
		columns[i] = column.Name
	}
	p := &plan{
		group: group, connector: connector, dialect: dbtarget.SQLDialectMySQL,
		columns: runtimeLayout.Columns, members: previewMembers, ownerColumn: ownerColumn,
		strategy: dbtarget.GroupEffectReceipt, scope: scope,
		mysql: &mysqlTestPlan{proof: proof, schemaDigest: runtimeLayout.SchemaDigest, layout: layout, members: fixtures, columns: columns},
	}
	if p.digest, err = p.contentDigest(); err != nil {
		return nil, err
	}
	if _, err := p.mysqlRows(ownerValue("op-preview"), time.Unix(0, 0).UTC()); err != nil {
		return nil, err
	}
	return p, nil
}

func validateMySQLTestWriteLayout(group *workspace.WriteGroup, runtimeLayout *workspace.WriteGroupRuntimeLayout) error {
	if runtimeLayout == nil || runtimeLayout.Dialect != dbtarget.SQLDialectMySQL || runtimeLayout.MySQL == nil ||
		!runtimeLayout.ReceiptTableReady || strings.TrimSpace(runtimeLayout.SchemaDigest) == "" {
		return unsupported(ReasonLayoutBlocked)
	}
	proof := runtimeLayout.MySQL
	if proof.ConnectorID != group.Destination.ConnectorID || proof.ConnectorRevision != group.Destination.ConnectorRevision ||
		proof.Group == nil || proof.Receipt == nil || proof.Group.Inspection == nil || proof.Receipt.Inspection == nil ||
		proof.Group.Evidence == nil || proof.Receipt.Evidence == nil ||
		proof.Group.Inspection.Table != group.Destination.TableName ||
		proof.Receipt.Inspection.Table != dbtarget.EffectReceiptTable ||
		proof.Group.Inspection.Status != dbtarget.TableInspectionExists ||
		proof.Receipt.Inspection.Status != dbtarget.TableInspectionExists ||
		proof.Group.Evidence.Database != group.Destination.Database ||
		proof.Receipt.Evidence.Database != group.Destination.Database ||
		proof.Group.Evidence.ServerFamily != proof.Receipt.Evidence.ServerFamily ||
		proof.Group.Evidence.ServerVersion != proof.Receipt.Evidence.ServerVersion ||
		proof.Group.Evidence.CapabilityRevision != proof.Receipt.Evidence.CapabilityRevision {
		return unsupported(ReasonLayoutBlocked)
	}
	native := group.StorageMode == workspace.WriteGroupStorageModeNativeMonthly
	if native != (proof.Group.Evidence.Partitions != nil) {
		return unsupported(ReasonLayoutBlocked)
	}
	return nil
}

func (p *plan) mysqlRows(owner string, at time.Time) ([]dbtarget.EncodedRow, error) {
	if p.mysql == nil {
		return nil, unsupported(ReasonGroupNotTestable)
	}
	at = at.UTC()
	rows := make([]dbtarget.EncodedRow, 0, len(p.mysql.members))
	for index, fixture := range p.mysql.members {
		recordID, err := snapshot.RecordID(p.group.WorkspaceID, p.group.ID, p.group.Revision, owner+"|"+fixture.key, at)
		if err != nil {
			return nil, unsupported(ReasonGroupNotTestable)
		}
		effectKey, err := snapshot.EffectKey(p.scope, recordID)
		if err != nil {
			return nil, unsupported(ReasonGroupNotTestable)
		}
		result := snapshot.MemberResult{
			MemberKey: fixture.key, Status: snapshot.MemberOK,
			Sample: &snapshot.Sample{
				SampleID: fmt.Sprintf("%s-%d", owner, index), MemberKey: fixture.key,
				ObservedAt: at, ReceivedAt: at, Quality: schema.QualityGood, Value: fixture.value,
			},
		}
		if fixture.reduction == string(workspace.WriteGroupSampleReductionMean) {
			mean, err := measurement.Mean([]measurement.ExactValue{fixture.value})
			if err != nil {
				return nil, unsupported(ReasonSQLValueBlocked)
			}
			result.Reduction = &mean
		}
		outcome := snapshot.Outcome{
			Kind: snapshot.OutcomeRow, EntityKey: owner, RecordID: recordID, EffectKey: effectKey,
			BucketStart: at, BucketEnd: at, Members: []snapshot.MemberResult{result},
		}
		row, err := p.mysql.layout.EncodeRow(outcome)
		if err != nil {
			return nil, unsupported(ReasonSQLValueBlocked)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func mysqlValueSlot(kind measurement.ExactType) string {
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
