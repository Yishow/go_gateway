package workspace

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/recordingplan"
)

// sqlServerNativeMonthlyLayout retains confirmed UTC boundaries across clock
// rollovers. Only a new, unconfirmed scope starts from the current UTC month.
func (s *WriteGroupService) sqlServerNativeMonthlyLayout(ctx context.Context, group *WriteGroup) (*recordingplan.SQLServerPartitionEvidence, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("%w: sqlserver storage repository is unavailable", ErrWriteGroupValidation)
	}
	lineage, err := recordingplan.ReadStorageLineage(ctx, s.repo.db, StorageScopeOf(group))
	if err != nil {
		if errors.Is(err, recordingplan.ErrStorageLineageNotFound) {
			return nil, fmt.Errorf("%w: sqlserver storage scope is not saved", ErrWriteGroupValidation)
		}
		return nil, fmt.Errorf("%w: read sqlserver native monthly confirmation: %w", ErrWriteGroupValidation, err)
	}
	table, owner, ok := LongStorageIdentity(group)
	if !ok || lineage.TableName != table || lineage.OwnerColumn != owner {
		return nil, ErrWriteGroupValidation
	}
	if lineage.Layout == nil {
		if lineage.ConfirmedOperationID != "" || lineage.LayoutDigest != recordingplan.UnconfirmedSQLServerStorageDigest(StorageScopeOf(group)) {
			return nil, ErrWriteGroupValidation
		}
		functionName, schemeName, err := recordingplan.SQLServerNativeMonthlyObjectNames(table)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrWriteGroupValidation, err)
		}
		return recordingplan.NewSQLServerNativeMonthlyPartition(s.now(), functionName, schemeName)
	}
	if !lineage.Confirmed() || lineage.Layout.SQLServer == nil || lineage.Layout.SQLServer.NativeMonthly == nil ||
		lineage.Layout.TableName != table || lineage.Layout.OwnerColumn != owner {
		return nil, ErrWriteGroupValidation
	}
	encoded, err := json.Marshal(lineage.Layout)
	if err != nil {
		return nil, fmt.Errorf("encode sqlserver native monthly confirmation: %w", err)
	}
	hash := sha256.Sum256(encoded)
	if lineage.LayoutDigest != hex.EncodeToString(hash[:]) {
		return nil, ErrWriteGroupValidation
	}
	return recordingplan.CloneSQLServerPartitionEvidence(lineage.Layout.SQLServer.NativeMonthly), nil
}
