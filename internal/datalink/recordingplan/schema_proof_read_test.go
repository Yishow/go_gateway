package recordingplan

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVerifiedSchemaProofUsesCallerSnapshotAndExactReplay(t *testing.T) {
	db := recordingRepositoryDB(t)
	db.SetMaxOpenConns(1)
	repo := NewSQLRepository(db)
	svc := NewService(repo)
	scope := previewScope()
	token := preparedPreview(t, svc, scope)
	op, _, err := svc.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	digest := strings.Repeat("a", 64)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()
	verified, err := VerifiedSchemaProofInTx(ctx, tx, scope, op.OperationID, digest)
	require.NoError(t, err)
	require.Equal(t, token.Token, verified.Token)
	changed := scope
	changed.WorkspaceRevision = "new-current-cas"
	_, err = VerifiedSchemaProofInTx(ctx, tx, changed, op.OperationID, digest)
	require.ErrorIs(t, err, ErrPreviewTokenStale)
	changed.SchemaRevision, changed.SchemaDigest = op.OperationID, digest
	_, err = VerifiedSchemaProofInTx(ctx, tx, changed, op.OperationID, digest)
	require.NoError(t, err, "exact already-recorded proof replays read-only")
	changed.ConnectorRevision = "foreign"
	_, err = VerifiedSchemaProofInTx(ctx, tx, changed, op.OperationID, digest)
	require.ErrorIs(t, err, ErrPreviewTokenStale)
	_, err = tx.ExecContext(ctx, `UPDATE managed_schema_preview_tokens SET database_name='foreign-uncommitted' WHERE token=?`, token.Token)
	require.NoError(t, err)
	_, err = VerifiedSchemaProofInTx(ctx, tx, scope, op.OperationID, digest)
	require.ErrorIs(t, err, ErrSchemaOperationMismatch, "must see caller's uncommitted destination and its mismatched ledger scope")
}

func TestVerifiedSchemaProofRejectsNilTransactionAndForeignOperation(t *testing.T) {
	_, err := VerifiedSchemaProofInTx(t.Context(), nil, previewScope(), "op", strings.Repeat("a", 64))
	require.ErrorIs(t, err, ErrSchemaOperationResult)
	for _, mode := range []string{"foreign-workspace", "foreign-op", "foreign-scope", "wrong-status", "wrong-digest"} {
		t.Run(mode, func(t *testing.T) {
			db := recordingRepositoryDB(t)
			repo := NewSQLRepository(db)
			svc := NewService(repo)
			scope := previewScope()
			token := preparedPreview(t, svc, scope)
			op, _, err := svc.ClaimSchemaApply(t.Context(), token)
			require.NoError(t, err)
			digest := strings.Repeat("a", 64)
			_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded, VerifiedDigest: digest})
			require.NoError(t, err)
			tx, err := db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			defer tx.Rollback()
			switch mode {
			case "foreign-workspace":
				scope.WorkspaceID = "foreign"
			case "foreign-op":
				_, err = tx.Exec(`UPDATE managed_schema_preview_tokens SET operation_id='foreign' WHERE token=?`, token.Token)
			case "foreign-scope":
				_, err = tx.Exec(`UPDATE managed_schema_operations SET scope_key='foreign' WHERE operation_id=?`, op.OperationID)
			case "wrong-status":
				_, err = tx.Exec(`UPDATE managed_schema_operations SET status='unknown' WHERE operation_id=?`, op.OperationID)
			case "wrong-digest":
				digest = strings.Repeat("b", 64)
			}
			require.NoError(t, err)
			_, err = VerifiedSchemaProofInTx(t.Context(), tx, scope, op.OperationID, digest)
			require.True(t, errors.Is(err, ErrSchemaOperationNotFound) || errors.Is(err, ErrSchemaOperationMismatch) || errors.Is(err, ErrSchemaOperationResult))
		})
	}
}
