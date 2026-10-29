package workspace

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

func TestSuccessorSchemaCannotReplaceNewerProofWithOldLedgerSuccess(t *testing.T) {
	h, id, req, api := successorSchemaFixture(t)
	scope, err := api.SuccessorSchemaScope(t.Context(), id, req)
	require.NoError(t, err)
	repo := recordingplan.NewSQLRepository(h.db)
	ledger := recordingplan.NewService(repo)
	token, err := ledger.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
		return recordingplan.TargetTableInspection{Status: "missing"}, nil
	})
	require.NoError(t, err)
	op, _, err := ledger.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	// Seed retained older successful ledger evidence, not a claim that this test
	// performed external DDL. The actual SQL path is covered by API tests.
	digest := strings.Repeat("a", 64)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_groups SET destination_schema_revision='newer-schema-operation',destination_schema_digest=? WHERE id=?`, strings.Repeat("b", 64), req.Current.GroupID)
	require.NoError(t, err)
	before := recreateRepositorySnapshot(t, h)
	_, err = api.ConfirmSuccessorSchema(t.Context(), id, req, op.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	t.Logf("proof refusal: %v", err)
	require.Error(t, err, "an older valid ledger result must not overwrite a newer schema proof")
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
}

func TestSuccessorSchemaRejectsMismatchedLedgerEvidence(t *testing.T) {
	for _, mode := range []string{"failed", "foreign-plan", "foreign-connector", "source", "layout", "created", "legacy-prefix", "malformed", "database", "schema", "dialect", "token-digest", "statements"} {
		t.Run(mode, func(t *testing.T) {
			h, id, req, api := successorSchemaFixture(t)
			scope, err := api.SuccessorSchemaScope(t.Context(), id, req)
			require.NoError(t, err)
			repo := recordingplan.NewSQLRepository(h.db)
			ledger := recordingplan.NewService(repo)
			token, err := ledger.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
				return recordingplan.TargetTableInspection{Status: "missing"}, nil
			})
			require.NoError(t, err)
			op, _, err := ledger.ClaimSchemaApply(t.Context(), token)
			require.NoError(t, err)
			digest := strings.Repeat("a", 64)
			_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
			require.NoError(t, err)
			evidence := &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true}
			switch mode {
			case "failed":
				_, err = h.db.Exec(`UPDATE managed_schema_operations SET status='failed' WHERE operation_id=?`, op.OperationID)
			case "foreign-plan":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET plan_id='foreign' WHERE token=?`, token.Token)
			case "foreign-connector":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET connector_revision='foreign' WHERE token=?`, token.Token)
			case "source":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET tables=json_set(tables,'$.source_digest','foreign') WHERE token=?`, token.Token)
			case "layout":
				evidence.Layout.TableName = "foreign"
			case "created":
				evidence.CreatedTable = false
			case "legacy-prefix":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET table_prefix='legacy_prefix' WHERE token=?`, token.Token)
			case "database":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET database_name='foreign' WHERE token=?`, token.Token)
			case "schema":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET schema_name='foreign' WHERE token=?`, token.Token)
			case "dialect":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET dialect='postgres' WHERE token=?`, token.Token)
			case "token-digest":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET digest='foreign' WHERE token=?`, token.Token)
			case "statements":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET statements='[]' WHERE token=?`, token.Token)
			case "malformed":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET tables='broken-json' WHERE token=?`, token.Token)
			}
			require.NoError(t, err)
			before := recreateRepositorySnapshot(t, h)
			_, err = api.ConfirmSuccessorSchema(t.Context(), id, req, op.OperationID, digest, evidence)
			require.Error(t, err)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func TestSuccessorSchemaOldTokenCannotAuthorizeNewCurrentCAS(t *testing.T) {
	h, id, req, api := successorSchemaFixture(t)
	scope, err := api.SuccessorSchemaScope(t.Context(), id, req)
	require.NoError(t, err)
	repo := recordingplan.NewSQLRepository(h.db)
	ledger := recordingplan.NewService(repo)
	token, err := ledger.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
		return recordingplan.TargetTableInspection{Status: "missing"}, nil
	})
	require.NoError(t, err)
	op, _, err := ledger.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	digest := strings.Repeat("a", 64)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	updated, err := h.workspace.UpdateDatabaseSetup(t.Context(), req.Current.WorkspaceRevision, func(context.Context, *sql.Tx, *Record) error { return nil })
	require.NoError(t, err)
	req.Current.WorkspaceRevision = updated.DatabaseSetupRevision
	before := recreateRepositorySnapshot(t, h)
	_, err = api.ConfirmSuccessorSchema(t.Context(), id, req, op.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	require.Error(t, err, "a new current CAS cannot silently re-authorize an old token")
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
}
