package workspace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// mysqlNativeMonthlyLayout retains confirmed UTC boundaries across clock
// rollovers. Only a new, unconfirmed scope starts from the current UTC month.
func (s *WriteGroupService) mysqlNativeMonthlyLayout(ctx context.Context, group *WriteGroup) (*recordingplan.MySQLPartitionEvidence, error) {
	lineage, err := recordingplan.ReadStorageLineage(ctx, s.repo.db, StorageScopeOf(group))
	if err != nil {
		return nil, fmt.Errorf("%w: read native monthly confirmation: %w", ErrWriteGroupValidation, err)
	}
	table, owner, ok := LongStorageIdentity(group)
	if !ok || lineage.TableName != table || lineage.OwnerColumn != owner {
		return nil, ErrWriteGroupValidation
	}
	if lineage.Layout == nil {
		if lineage.ConfirmedOperationID != "" || lineage.LayoutDigest != recordingplan.UnconfirmedMySQLStorageDigest(StorageScopeOf(group)) {
			return nil, ErrWriteGroupValidation
		}
		return recordingplan.NewMySQLNativeMonthlyPartition(s.now())
	}
	if !lineage.Confirmed() || lineage.Layout.MySQL == nil || lineage.Layout.MySQL.NativeMonthly == nil ||
		lineage.Layout.TableName != table || lineage.Layout.OwnerColumn != owner {
		return nil, ErrWriteGroupValidation
	}
	encoded, err := json.Marshal(lineage.Layout)
	if err != nil {
		return nil, fmt.Errorf("encode native monthly confirmation: %w", err)
	}
	hash := sha256.Sum256(encoded)
	if lineage.LayoutDigest != hex.EncodeToString(hash[:]) {
		return nil, ErrWriteGroupValidation
	}
	return lineage.Layout.MySQL.NativeMonthly, nil
}

func recordingPlanMySQLPartitions(evidence *dbtarget.MySQLPartitionEvidence) *recordingplan.MySQLPartitionEvidence {
	if evidence == nil {
		return nil
	}
	result := &recordingplan.MySQLPartitionEvidence{Method: evidence.Method, Expression: evidence.Expression,
		Partitions: make([]recordingplan.MySQLPartitionBoundary, len(evidence.Partitions))}
	for i, partition := range evidence.Partitions {
		result.Partitions[i] = recordingplan.MySQLPartitionBoundary{Name: partition.Name, LessThan: partition.LessThan, Ordinal: partition.Ordinal}
	}
	return result
}
