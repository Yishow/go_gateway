package recordingplan

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestTestWriteRecoveryRetainsSupportedCleanupClaim(t *testing.T) {
	for _, sqlRepo := range []bool{false, true} {
		for _, cleanup := range []string{"failed", "unknown"} {
			t.Run(fmt.Sprint(sqlRepo, "/", cleanup), func(t *testing.T) {
				var repo Repository = NewMemoryRepository()
				if sqlRepo {
					db := openLedgerDB(t, filepath.Join(t.TempDir(), "ledger.db"))
					applyRecordingDDL(t, db)
					repo = NewSQLRepository(db)
				}
				svc := NewService(repo)
				scope := testWriteScope()
				scope.Dialect = "sqlserver"
				token, finished := seedRecovery(t, svc, scope, SchemaOperationSucceeded, "written_verified", cleanup)
				claimed, outcome, err := svc.ClaimTestWriteCleanup(t.Context(), scope.WorkspaceID, token.Token, token.OperationID, scope, "content")
				require.NoError(t, err)
				require.Equal(t, ClaimAcquired, outcome)
				got, err := svc.ListUnresolvedTestWrites(t.Context(), scope.WorkspaceID)
				require.NoError(t, err)
				require.Len(t, got.Operations, 1)
				item := got.Operations[0]
				require.Equal(t, claimed.OperationID, item.OperationID)
				require.Equal(t, SchemaOperationRunning, item.Status)
				require.Equal(t, "written_verified", item.WriteOutcome)
				require.Equal(t, cleanup, item.CleanupStatus)
				require.Equal(t, finished.CompletedAt, item.CompletedAt)
				after, err := repo.GetSchemaOperation(t.Context(), scope.WorkspaceID, claimed.OperationID)
				require.NoError(t, err)
				require.Equal(t, claimed, after)
			})
		}
	}
}

func TestTestWriteRecoveryRejectsOtherActiveResultShapes(t *testing.T) {
	for _, change := range []string{"pending", "unverified", "cleaned", "no-proof", "no-completion", "other-dialect"} {
		t.Run(change, func(t *testing.T) {
			repo := NewMemoryRepository()
			svc := NewService(repo)
			scope := testWriteScope()
			scope.Dialect = "sqlserver"
			if change == "other-dialect" {
				scope.Dialect = "mysql"
			}
			token, op := seedRecovery(t, svc, scope, SchemaOperationSucceeded, "written_verified", "failed")
			_, _, err := svc.ClaimTestWriteCleanup(t.Context(), scope.WorkspaceID, token.Token, op.OperationID, scope, "content")
			require.NoError(t, err)
			value := repo.operations[op.OperationID]
			switch change {
			case "pending":
				value.Status = SchemaOperationPending
			case "unverified":
				value.WriteOutcome = "written_unverified"
			case "cleaned":
				value.CleanupStatus = "cleaned"
			case "no-proof":
				value.Detail = ""
			case "no-completion":
				value.CompletedAt = nil
			}
			repo.operations[op.OperationID] = value
			got, err := svc.ListUnresolvedTestWrites(t.Context(), scope.WorkspaceID)
			require.ErrorIs(t, err, ErrTestWriteRecoveryUnavailable)
			require.Nil(t, got)
		})
	}
}

func TestTestWriteRecoverySQLMalformedResolvedRowsCannotDisappear(t *testing.T) {
	for _, corruption := range []string{"invalid-all-times", "invalid-completion", "null-operation-id"} {
		t.Run(corruption, func(t *testing.T) {
			db := openLedgerDB(t, filepath.Join(t.TempDir(), "ledger.db"))
			applyRecordingDDL(t, db)
			svc := NewService(NewSQLRepository(db))
			_, _ = seedRecovery(t, svc, testWriteScope(), SchemaOperationSucceeded, "written_verified", "cleaned")
			queries := map[string]string{
				"invalid-all-times":  `UPDATE managed_schema_operations SET created_at='invalid',updated_at='invalid',completed_at='invalid'`,
				"invalid-completion": `UPDATE managed_schema_operations SET completed_at='invalid'`,
				"null-operation-id":  `UPDATE managed_schema_operations SET operation_id=NULL`,
			}
			_, err := db.Exec(queries[corruption])
			require.NoError(t, err)
			result, err := svc.ListUnresolvedTestWrites(t.Context(), "ws-1")
			require.ErrorIs(t, err, ErrTestWriteRecoveryUnavailable)
			require.Nil(t, result)
		})
	}
}

func TestTestWriteRecoveryCommittedUnverifiedCleanupParity(t *testing.T) {
	for _, sqlRepo := range []bool{false, true} {
		for _, cleanup := range []string{"cleaned", "failed", "unknown"} {
			t.Run(fmt.Sprint(sqlRepo, "/", cleanup), func(t *testing.T) {
				var repo Repository = NewMemoryRepository()
				if sqlRepo {
					db := openLedgerDB(t, filepath.Join(t.TempDir(), "ledger.db"))
					applyRecordingDDL(t, db)
					repo = NewSQLRepository(db)
				}
				svc := NewService(repo)
				token, original := seedRecovery(t, svc, testWriteScope(), SchemaOperationPartial, "written_unverified", cleanup)
				got, err := svc.ListUnresolvedTestWrites(t.Context(), "ws-1")
				require.NoError(t, err)
				if cleanup == "cleaned" {
					require.Empty(t, got.Operations)
				} else {
					require.Len(t, got.Operations, 1)
					require.Equal(t, original.OperationID, got.Operations[0].OperationID)
				}
				// Recovery classification must never rewrite the historical readback facts.
				persisted, err := repo.GetSchemaOperation(t.Context(), "ws-1", original.OperationID)
				require.NoError(t, err)
				require.Equal(t, original, persisted)
				ordinary, err := svc.GetSchemaOperation(t.Context(), "ws-1", original.OperationID)
				require.NoError(t, err)
				require.Equal(t, SchemaOperationPartial, ordinary.Status)
				require.Equal(t, "written_unverified", ordinary.WriteOutcome)
				require.Equal(t, cleanup, ordinary.CleanupStatus)
				if cleanup == "cleaned" {
					require.NoError(t, repo.DeletePreviewToken(t.Context(), token.Token))
					got, err = svc.ListUnresolvedTestWrites(t.Context(), "ws-1")
					require.NoError(t, err)
					require.Empty(t, got.Operations)
				}
			})
		}
	}
}
