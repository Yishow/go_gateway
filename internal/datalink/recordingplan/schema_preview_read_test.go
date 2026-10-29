package recordingplan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStoredSchemaPreviewReadAuthenticatesExpiredAndUnsubmittedEvidence(t *testing.T) {
	db := recordingRepositoryDB(t)
	svc := NewService(NewSQLRepository(db))
	token := preparedPreview(t, svc, previewScope())
	// Expiry is an execution guard, not grounds for hiding historical evidence.
	_, err := db.Exec(`UPDATE managed_schema_preview_tokens SET expires_at=? WHERE token=?`, time.Now().Add(-time.Hour), token.Token)
	require.NoError(t, err)
	tx := beginSuccessorTx(t, db)
	got, err := StoredSchemaPreviewTokenInTx(t.Context(), tx, token.WorkspaceID, token.Token)
	require.NoError(t, err)
	require.Equal(t, token.OperationID, got.OperationID)
	require.True(t, got.ExpiresAt.Before(time.Now()))
	_, err = StoredSchemaPreviewTokenInTx(t.Context(), tx, "foreign", token.Token)
	require.ErrorIs(t, err, ErrSchemaOperationMismatch)
	_, err = tx.Exec(`UPDATE managed_schema_preview_tokens SET database_name='tampered' WHERE token=?`, token.Token)
	require.NoError(t, err)
	_, err = StoredSchemaPreviewTokenInTx(t.Context(), tx, token.WorkspaceID, token.Token)
	require.ErrorIs(t, err, ErrPreviewTokenStale)
}

func TestStoredSchemaPreviewReadRejectsMalformedAndLegacyEvidence(t *testing.T) {
	for _, mode := range []string{"nil-transaction", "legacy", "invalid-tables", "invalid-statements", "unknown-token"} {
		t.Run(mode, func(t *testing.T) {
			db := recordingRepositoryDB(t)
			token := preparedPreview(t, NewService(NewSQLRepository(db)), previewScope())
			tx := beginSuccessorTx(t, db)
			var err error
			switch mode {
			case "nil-transaction":
				tx = nil
			case "legacy":
				_, err = tx.Exec(`UPDATE managed_schema_preview_tokens SET digest='' WHERE token=?`, token.Token)
			case "invalid-tables":
				_, err = tx.Exec(`UPDATE managed_schema_preview_tokens SET tables='{' WHERE token=?`, token.Token)
			case "invalid-statements":
				_, err = tx.Exec(`UPDATE managed_schema_preview_tokens SET statements='{' WHERE token=?`, token.Token)
			case "unknown-token":
				token.Token = "missing-private-preview-token"
			}
			require.NoError(t, err)
			_, err = StoredSchemaPreviewTokenInTx(t.Context(), tx, token.WorkspaceID, token.Token)
			require.Error(t, err)
			if mode == "unknown-token" {
				require.NotContains(t, err.Error(), token.Token)
			}
		})
	}
}
