package workspace

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"reflect"

	"go-gateway/internal/datalink/recordingplan"
)

// Succeeded token facts select capability evidence, never group ownership or
// mode. The caller must subsequently authenticate the complete current scope
// with VerifiedSchemaProofInTx; this helper grants no independent authority.
func successorSQLServerProofLayoutInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup, recorded *recordingplan.GroupSchemaLayout) (*recordingplan.GroupSchemaLayout, error) {
	if recorded == nil || recorded.SQLServer == nil || !successorActivationEngineSupported(group, writeGroupConnectorKindSQLServer) ||
		(group.StorageMode == WriteGroupStorageModeNativeMonthly) != (recorded.SQLServer.NativeMonthly != nil) {
		return nil, ErrWriteGroupValidation
	}
	table, owner, ok := LongStorageIdentity(group)
	if !ok || table != group.Destination.TableName {
		return nil, ErrWriteGroupValidation
	}
	options := *recorded.SQLServer
	options.Database, options.Schema = group.Destination.Database, group.Destination.TableSchema
	if group.StorageMode == WriteGroupStorageModeNativeMonthly {
		if err := successorSQLServerPartitionProofInTx(ctx, tx, group, table, owner, options.NativeMonthly); err != nil {
			return nil, err
		}
	}
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, options)
	if err != nil {
		return nil, err
	}
	return &layout, nil
}

// Match the existing native confirmation contract inside the same local
// transaction. A shared scope may have been confirmed by another group; its
// immutable partition basis, not the last confirming operation ID, must match.
func successorSQLServerPartitionProofInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup, table, owner string, native *recordingplan.SQLServerPartitionEvidence) error {
	lineage, err := recordingplan.ReadStorageLineage(ctx, tx, StorageScopeOf(group))
	if err != nil {
		return err
	}
	if lineage.TableName != table || lineage.OwnerColumn != owner {
		return ErrWriteGroupValidation
	}
	if lineage.Layout == nil {
		if group.Destination.SchemaRevision != "" || lineage.ConfirmedOperationID != "" || lineage.LayoutDigest != recordingplan.UnconfirmedSQLServerStorageDigest(StorageScopeOf(group)) {
			return ErrWriteGroupValidation
		}
		function, scheme, err := recordingplan.SQLServerNativeMonthlyObjectNames(table)
		if err != nil || native == nil || native.FunctionName != function || native.SchemeName != scheme {
			return ErrWriteGroupValidation
		}
		return nil
	}
	if !lineage.Confirmed() || lineage.Layout.SQLServer == nil || lineage.Layout.SQLServer.NativeMonthly == nil ||
		lineage.Layout.TableName != table || lineage.Layout.OwnerColumn != owner ||
		lineage.Layout.SQLServer.Database != group.Destination.Database || lineage.Layout.SQLServer.Schema != group.Destination.TableSchema ||
		!reflect.DeepEqual(lineage.Layout.SQLServer.NativeMonthly, native) {
		return ErrWriteGroupValidation
	}
	encoded, err := json.Marshal(lineage.Layout)
	if err != nil {
		return err
	}
	hash := sha256.Sum256(encoded)
	if lineage.LayoutDigest != hex.EncodeToString(hash[:]) {
		return ErrWriteGroupValidation
	}
	return nil
}
