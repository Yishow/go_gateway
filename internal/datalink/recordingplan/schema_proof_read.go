package recordingplan

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type schemaProofReader interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// readSchemaPreviewToken is shared by ordinary reads and transaction-bound proof.
func readSchemaPreviewToken(ctx context.Context, q schemaProofReader, token string) (*SchemaPreviewToken, error) {
	query := `
		SELECT token, workspace_id, plan_id, plan_revision, connector_id, table_prefix,
			statements, expires_at, created_at, operation_id, action, workspace_revision,
			connector_revision, dialect, database_name, schema_name, no_change_reason, digest, tables
		FROM managed_schema_preview_tokens
		WHERE token = $1
	`
	query = adaptPlaceholders(query)

	var t SchemaPreviewToken
	var stmtsStr, tablesStr string
	err := q.QueryRowContext(ctx, query, token).Scan(
		&t.Token, &t.WorkspaceID, &t.PlanID, &t.PlanRevision,
		&t.ConnectorID, &t.TablePrefix, &stmtsStr, &t.ExpiresAt, &t.CreatedAt,
		&t.OperationID, &t.Action, &t.WorkspaceRevision, &t.ConnectorRevision, &t.Dialect,
		&t.Database, &t.Schema, &t.NoChangeReason, &t.Digest, &tablesStr,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", ErrPreviewTokenNotFound, token)
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(stmtsStr), &t.Statements); err != nil {
		return nil, fmt.Errorf("decode preview statements: %w", err)
	}
	if err := decodePreviewTables(tablesStr, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// VerifiedSchemaProofInTx reads an already succeeded schema operation and its
// complete token inside the caller's transaction. It grants no execution rights.
// Current must be resolved by the owning service, not copied from request input.
func VerifiedSchemaProofInTx(ctx context.Context, tx *sql.Tx, current SchemaPreviewScope, operationID, verifiedDigest string) (*SchemaPreviewToken, error) {
	if tx == nil || operationID == "" || len(verifiedDigest) != 64 {
		return nil, ErrSchemaOperationResult
	}
	current, err := normalizePreviewScope(current)
	if err != nil {
		return nil, err
	}
	token, err := SucceededSchemaTokenInTx(ctx, tx, current.WorkspaceID, operationID, verifiedDigest)
	if err != nil {
		return nil, err
	}
	// An exact proof already saved by this operation is read-only replay. Only
	// its own schema/workspace advancement is reconciled; every other scope fact
	// still goes through the ordinary complete-token mismatch validator.
	if current.SchemaRevision == operationID && current.SchemaDigest == verifiedDigest {
		current.WorkspaceRevision = token.WorkspaceRevision
		current.SchemaRevision, current.SchemaDigest = token.SchemaRevision, token.SchemaDigest
	}
	if field := previewScopeMismatch(token, current); field != "" {
		return nil, fmt.Errorf("%w: %s changed", ErrPreviewTokenStale, field)
	}
	return token, nil
}

// SucceededSchemaTokenInTx reads and authenticates the existing operation/token
// ledger binding and token digest. It does NOT verify the caller's current scope
// or grant mutation authority; callers must still use VerifiedSchemaProofInTx.
// It exists only to derive backend-specific expected layouts from durable proof.
func SucceededSchemaTokenInTx(ctx context.Context, tx *sql.Tx, workspaceID, operationID, verifiedDigest string) (*SchemaPreviewToken, error) {
	if tx == nil || workspaceID == "" || operationID == "" || len(verifiedDigest) != 64 {
		return nil, ErrSchemaOperationResult
	}
	var tokenID, scopeKey, action, status, digest string
	err := tx.QueryRowContext(ctx, adaptPlaceholders(`SELECT token,scope_key,action,status,verified_digest FROM managed_schema_operations WHERE workspace_id=$1 AND operation_id=$2`), workspaceID, operationID).Scan(&tokenID, &scopeKey, &action, &status, &digest)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrSchemaOperationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("read verified schema operation: %w", err)
	}
	if status != string(SchemaOperationSucceeded) || digest != verifiedDigest || action != SchemaApplyAction {
		return nil, ErrSchemaOperationResult
	}
	token, err := readSchemaPreviewToken(ctx, tx, tokenID)
	if err != nil {
		return nil, err
	}
	if token.WorkspaceID != workspaceID || token.OperationID != operationID || token.Action != action || schemaOperationScopeKey(token) != scopeKey {
		return nil, ErrSchemaOperationMismatch
	}
	if isLegacyPreviewToken(token) {
		return nil, ErrPreviewTokenLegacy
	}
	actualDigest, err := previewDigest(token)
	if err != nil {
		return nil, err
	}
	if subtle.ConstantTimeCompare([]byte(actualDigest), []byte(token.Digest)) != 1 {
		return nil, ErrPreviewTokenStale
	}
	return token, nil
}
