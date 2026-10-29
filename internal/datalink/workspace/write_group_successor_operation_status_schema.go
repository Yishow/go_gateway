package workspace

import (
	"context"
	"database/sql"
	"errors"
	"reflect"

	"go-gateway/internal/datalink/recordingplan"
)

func (s *WriteGroupService) successorSchemaOperationsInTx(ctx context.Context, tx *sql.Tx, reserved *WriteGroupSuccessorPreparation, binding *successorReserveSchemaBinding) ([]WriteGroupSuccessorSchemaStatus, error) {
	rows, err := tx.QueryContext(ctx, s.repo.query(`SELECT token FROM managed_schema_preview_tokens WHERE workspace_id=$1 AND plan_id=$2 ORDER BY created_at,operation_id`), binding.WorkspaceID, binding.GroupID)
	if err != nil {
		return nil, err
	}
	var candidates []string
	for rows.Next() {
		var candidate string
		if err := rows.Scan(&candidate); err != nil {
			_ = rows.Close()
			return nil, err
		}
		candidates = append(candidates, candidate)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	result := make([]WriteGroupSuccessorSchemaStatus, 0, len(candidates))
	for _, candidate := range candidates {
		token, err := recordingplan.StoredSchemaPreviewTokenInTx(ctx, tx, binding.WorkspaceID, candidate)
		if err != nil {
			return nil, ErrWriteGroupLineageCorrupt
		}
		// A token for different immutable sources is not a result of this root.
		if token.SourceDigest != binding.SourceDigest {
			continue
		}
		if err := successorStatusTokenScope(token, reserved, binding); err != nil {
			return nil, err
		}
		if err := successorStatusPriorSchemaInTx(ctx, tx, token, reserved, binding); err != nil {
			return nil, err
		}
		item, err := s.successorSchemaStatusInTx(ctx, tx, token)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

// A later preview may legitimately carry the schema proof saved by an earlier
// preview. Authenticate that exact chain, never replace it with a latest proof.
func successorStatusPriorSchemaInTx(ctx context.Context, tx *sql.Tx, token *recordingplan.SchemaPreviewToken, reserved *WriteGroupSuccessorPreparation, binding *successorReserveSchemaBinding) error {
	seen := map[string]bool{token.OperationID: true}
	for token.SchemaRevision != "" {
		if seen[token.SchemaRevision] {
			return ErrWriteGroupLineageCorrupt
		}
		seen[token.SchemaRevision] = true
		previous, err := recordingplan.SucceededSchemaTokenInTx(ctx, tx, binding.WorkspaceID, token.SchemaRevision, token.SchemaDigest)
		if err != nil {
			return ErrWriteGroupLineageCorrupt
		}
		if err := successorStatusTokenScope(previous, reserved, binding); err != nil {
			return err
		}
		token = previous
	}
	return nil
}

func (s *WriteGroupService) successorSchemaStatusInTx(ctx context.Context, tx *sql.Tx, token *recordingplan.SchemaPreviewToken) (WriteGroupSuccessorSchemaStatus, error) {
	result := WriteGroupSuccessorSchemaStatus{OperationID: token.OperationID, Status: "preview_not_submitted"}
	var tokenID, workspaceID, scope, action, status, verified string
	err := tx.QueryRowContext(ctx, s.repo.query(`SELECT token,workspace_id,scope_key,action,status,verified_digest FROM managed_schema_operations WHERE operation_id=$1`), token.OperationID).Scan(&tokenID, &workspaceID, &scope, &action, &status, &verified)
	if errors.Is(err, sql.ErrNoRows) {
		return result, nil
	}
	if err != nil {
		return result, err
	}
	if tokenID != token.Token || workspaceID != token.WorkspaceID || scope != token.ScopeKey() || action != recordingplan.SchemaApplyAction {
		return result, ErrWriteGroupLineageCorrupt
	}
	result.Status = status
	// Never copy raw ledger reasons/actions: old rows can contain target errors.
	switch recordingplan.SchemaOperationStatus(status) {
	case recordingplan.SchemaOperationPending, recordingplan.SchemaOperationRunning:
		result.NextAction = "read_operation_status"
	case recordingplan.SchemaOperationSucceeded:
		if !successorStatusDigestValid(verified) {
			return result, ErrWriteGroupLineageCorrupt
		}
	case recordingplan.SchemaOperationFailed:
		result.ReasonCode, result.NextAction = "schema_failed", "review_schema_operation"
	case recordingplan.SchemaOperationPartial:
		result.ReasonCode, result.NextAction = "schema_partial", "reconcile_schema_operation"
	case recordingplan.SchemaOperationUnknown:
		result.ReasonCode, result.NextAction = "schema_unknown", "reconcile_schema_operation"
	default:
		return result, ErrWriteGroupLineageCorrupt
	}
	return result, nil
}

func successorStatusTokenScope(token *recordingplan.SchemaPreviewToken, reserved *WriteGroupSuccessorPreparation, binding *successorReserveSchemaBinding) error {
	g := reserved.Group
	if token.WorkspaceID != binding.WorkspaceID || token.PlanID != g.ID || token.PlanRevision != g.Revision ||
		token.ConnectorID != g.Destination.ConnectorID || token.ConnectorRevision != g.Destination.ConnectorRevision || token.Dialect != binding.Dialect ||
		token.Database != g.Destination.Database || token.Schema != g.Destination.TableSchema || token.TablePrefix != "" ||
		token.Action != recordingplan.SchemaApplyAction || token.SourceDigest != binding.SourceDigest || token.GroupLayout == nil ||
		(token.SchemaRevision == "") != (token.SchemaDigest == "") {
		return ErrWriteGroupLineageCorrupt
	}
	var policy *bool
	if managedLongGroup(g) {
		policy = g.WritePolicy.RecreateMissing
	}
	if !reflect.DeepEqual(policy, token.RecreateMissing) {
		return ErrWriteGroupLineageCorrupt
	}
	layout, err := successorStatusExpectedLayout(g, binding, token.GroupLayout)
	if err != nil {
		return ErrWriteGroupLineageCorrupt
	}
	validated, err := recordingplan.ValidateGroupSchemaLayout(*layout)
	if err != nil || !reflect.DeepEqual(&validated, token.GroupLayout) {
		return ErrWriteGroupLineageCorrupt
	}
	return nil
}

func successorStatusExpectedLayout(group *WriteGroup, binding *successorReserveSchemaBinding, recorded *recordingplan.GroupSchemaLayout) (*recordingplan.GroupSchemaLayout, error) {
	switch binding.Dialect {
	case writeGroupConnectorKindMySQL:
		if recorded.MySQL == nil {
			return nil, ErrWriteGroupLineageCorrupt
		}
		return successorMySQLProofLayout(group, binding.SourceFamily, recorded)
	case writeGroupConnectorKindSQLServer:
		if recorded.SQLServer == nil || !successorActivationEngineSupported(group, binding.Dialect) ||
			(group.StorageMode == WriteGroupStorageModeNativeMonthly) != (recorded.SQLServer.NativeMonthly != nil) {
			return nil, ErrWriteGroupLineageCorrupt
		}
		table, owner, ok := LongStorageIdentity(group)
		if !ok || table != group.Destination.TableName {
			return nil, ErrWriteGroupLineageCorrupt
		}
		options := *recorded.SQLServer
		options.Database, options.Schema = group.Destination.Database, group.Destination.TableSchema
		layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, options)
		return &layout, err
	default:
		layout, err := managedSchemaLayout(group, binding.TagTypes, binding.Dialect)
		if err != nil {
			return nil, err
		}
		if layout.Partition != nil {
			// Storage shape is independently derived. Only the extension's
			// installation namespace comes from the already authenticated token,
			// matching the existing schema-confirmation partition contract.
			if recorded.Partition == nil || recorded.Partition.ExtensionSchema == "" {
				return nil, ErrWriteGroupLineageCorrupt
			}
			layout.Partition.ExtensionSchema = recorded.Partition.ExtensionSchema
		}
		return layout, nil
	}
}
