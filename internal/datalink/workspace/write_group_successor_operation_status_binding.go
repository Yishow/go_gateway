package workspace

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
)

const successorReserveSchemaBindingVersion = "successor-reserve-schema-binding/v1"

// The original schema source facts are private, immutable ledger evidence. They
// are not refreshed from current tags/connectors when reading historical status.
type successorReserveSchemaBinding struct {
	Version             string                     `json:"version"`
	WorkspaceID         string                     `json:"workspace_id"`
	LineageID           string                     `json:"lineage_id"`
	PredecessorGroupID  string                     `json:"predecessor_group_id"`
	GroupID             string                     `json:"group_id"`
	OperationID         string                     `json:"operation_id"`
	IntentDigest        string                     `json:"intent_digest"`
	ManagedSourceDigest string                     `json:"managed_source_digest"`
	SourceDigest        string                     `json:"source_digest"`
	Dialect             string                     `json:"dialect"`
	SourceFamily        string                     `json:"source_family"`
	TagTypes            map[string]schema.DataType `json:"tag_types"`
}

func (s *WriteGroupService) makeSuccessorReserveSchemaBinding(ctx context.Context, tx *sql.Tx, result *WriteGroupSuccessorPreparation, evidence *successorReserveSourceEvidence) (json.RawMessage, error) {
	if result == nil || result.Group == nil {
		return nil, ErrWriteGroupLineageCorrupt
	}
	types, err := loadWriteGroupTagTypes(ctx, tx, s.repo.query, result.Group)
	if err != nil {
		return nil, err
	}
	dialect, err := s.successorDriverInTx(ctx, tx, result.Group.Destination.ConnectorID)
	if err != nil {
		return nil, err
	}
	if dialect == "mariadb" {
		dialect = writeGroupConnectorKindMySQL
	}
	family, err := s.successorReserveBindingFamily(ctx, tx, result, dialect, evidence)
	if err != nil {
		return nil, err
	}
	managed, err := managedSourceDigest(result.Group, types)
	if err != nil {
		return nil, err
	}
	b := successorReserveSchemaBinding{Version: successorReserveSchemaBindingVersion, WorkspaceID: result.Group.WorkspaceID,
		LineageID: result.LineageID, PredecessorGroupID: result.PredecessorGroupID, GroupID: result.Group.ID,
		OperationID: result.OperationID, IntentDigest: result.IntentDigest, ManagedSourceDigest: managed, Dialect: dialect, SourceFamily: family, TagTypes: types}
	b.SourceDigest, err = b.sourceDigest()
	if err != nil {
		return nil, err
	}
	return json.Marshal(b)
}

func (b successorReserveSchemaBinding) sourceDigest() (string, error) {
	encoded, err := json.Marshal([]string{"write-group-successor-schema/v1", b.LineageID, b.PredecessorGroupID, b.GroupID, b.OperationID, b.IntentDigest, b.ManagedSourceDigest})
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}

func decodeSuccessorReserveSchemaBinding(raw json.RawMessage, result *WriteGroupSuccessorPreparation) (*successorReserveSchemaBinding, error) {
	var b successorReserveSchemaBinding
	if result == nil || result.Group == nil || decodeSuccessorJSON(raw, &b) != nil || b.Version != successorReserveSchemaBindingVersion ||
		b.WorkspaceID != result.Group.WorkspaceID || b.LineageID != result.LineageID || b.PredecessorGroupID != result.PredecessorGroupID ||
		b.GroupID != result.Group.ID || b.OperationID != result.OperationID || b.IntentDigest != result.IntentDigest || len(b.TagTypes) == 0 {
		return nil, ErrWriteGroupLineageCorrupt
	}
	switch b.Dialect {
	case writeGroupConnectorKindSQLite, writeGroupConnectorKindPostgres, writeGroupConnectorKindMySQL, writeGroupConnectorKindSQLServer:
	default:
		return nil, ErrWriteGroupLineageCorrupt
	}
	if b.Dialect == writeGroupConnectorKindMySQL {
		if b.SourceFamily != "mysql" && b.SourceFamily != "mariadb" {
			return nil, ErrWriteGroupLineageCorrupt
		}
	} else if b.SourceFamily != b.Dialect {
		return nil, ErrWriteGroupLineageCorrupt
	}
	managed, err := managedSourceDigest(result.Group, b.TagTypes)
	if err != nil || managed != b.ManagedSourceDigest {
		return nil, ErrWriteGroupLineageCorrupt
	}
	digest, err := b.sourceDigest()
	if err != nil || digest != b.SourceDigest {
		return nil, ErrWriteGroupLineageCorrupt
	}
	return &b, nil
}

func (s *WriteGroupService) successorReserveBindingFamily(ctx context.Context, tx *sql.Tx, result *WriteGroupSuccessorPreparation, dialect string, evidence *successorReserveSourceEvidence) (string, error) {
	if dialect != writeGroupConnectorKindMySQL {
		return dialect, nil
	}
	if evidence == nil || evidence.Source == nil || evidence.Source.ID != result.PredecessorGroupID || evidence.Source.WorkspaceID != result.Group.WorkspaceID {
		return "", ErrWriteGroupLineageCorrupt
	}
	var family string
	if evidence.Source.AppliedRevision != "" {
		var err error
		family, err = s.repo.recreateSourceEngineInTx(ctx, tx, evidence.Source)
		if err != nil {
			return "", err
		}
	} else {
		if !successorBindingCapability(evidence.SourceMySQL, evidence.Source) {
			return "", ErrWriteGroupLineageCorrupt
		}
		family = evidence.SourceMySQL.ServerFamily
	}
	if !successorBindingCapability(evidence.TargetMySQL, result.Group) || evidence.TargetMySQL.ServerFamily != family {
		return "", ErrWriteGroupLineageCorrupt
	}
	return family, nil
}

func successorBindingCapability(cap *dbtarget.MySQLRecordingCapability, group *WriteGroup) bool {
	return cap != nil && cap.Connected && cap.MetadataStatus == "verified" && cap.TransactionalEngine && cap.SupportsManagedLong &&
		cap.CapabilityRevision != "" && dbtarget.MySQLRecordingVersionVerified(cap.ServerFamily, cap.ServerVersion) &&
		cap.ConnectorID == group.Destination.ConnectorID && cap.ConnectorRevision == group.Destination.ConnectorRevision && cap.Database == group.Destination.Database &&
		(group.StorageMode != WriteGroupStorageModeNativeMonthly || cap.SupportsNativePartition)
}
