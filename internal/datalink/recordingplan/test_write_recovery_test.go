package recordingplan

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func seedRecovery(t *testing.T, svc *Service, scope TestWriteScope, status SchemaOperationStatus, write, cleanup string) (*SchemaPreviewToken, *SchemaOperation) {
	t.Helper()
	token, err := svc.PrepareTestWritePreview(t.Context(), scope, "content")
	require.NoError(t, err)
	op, claim, err := svc.ClaimTestWrite(t.Context(), token)
	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, claim)
	if status != SchemaOperationRunning {
		op, err = svc.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, SchemaOperationResult{Status: status, PayloadDigest: "content", WriteOutcome: write, CleanupStatus: cleanup, Reason: "secret-endpoint", NextAction: "password=private", Detail: "private-proof"})
		require.NoError(t, err)
	}
	return token, op
}

func TestTestWriteRecoveryResolvedHistoryDoesNotHideUnresolved(t *testing.T) {
	for _, sqlRepo := range []bool{false, true} {
		t.Run(fmt.Sprint(sqlRepo), func(t *testing.T) {
			var repo Repository = NewMemoryRepository()
			if sqlRepo {
				db := openLedgerDB(t, filepath.Join(t.TempDir(), "ledger.db"))
				applyRecordingDDL(t, db)
				repo = NewSQLRepository(db)
			}
			svc := NewService(repo)
			scope := testWriteScope()
			for i := 0; i < 1001; i++ {
				token, _ := seedRecovery(t, svc, scope, SchemaOperationSucceeded, "written_verified", "cleaned")
				require.NoError(t, repo.DeletePreviewToken(t.Context(), token.Token))
			}
			_, pending := seedRecovery(t, svc, scope, SchemaOperationRunning, "", "")
			result, err := svc.ListUnresolvedTestWrites(t.Context(), scope.WorkspaceID)
			require.NoError(t, err)
			require.True(t, result.Complete)
			require.Len(t, result.Operations, 1)
			require.Equal(t, pending.OperationID, result.Operations[0].OperationID)
		})
	}
}

func TestTestWriteRecoveryDurableHistoricalBindingsAndNoMutation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ledger.db")
	db := openLedgerDB(t, path)
	applyRecordingDDL(t, db)
	svc := NewService(NewSQLRepository(db))
	scope := testWriteScope()
	firstToken, first := seedRecovery(t, svc, scope, SchemaOperationSucceeded, "written_verified", "failed")
	scope.GroupID = "group-2"
	scope.GroupRevision = "rev-2"
	scope.ConnectorRevision = "identity-2"
	scope.Table = "other"
	secondToken, second := seedRecovery(t, svc, scope, SchemaOperationUnknown, "unknown", "unknown")
	// Expiry is not deletion and cannot discard an existing operation.
	_, err := db.Exec(`UPDATE managed_schema_preview_tokens SET expires_at=?`, time.Now().Add(-time.Hour))
	require.NoError(t, err)
	require.NoError(t, db.Close())
	reopened := openLedgerDB(t, path)
	repo := NewSQLRepository(reopened)
	svc = NewService(repo)
	var before string
	require.NoError(t, reopened.QueryRow(`SELECT group_concat(operation_id || status || owner || detail) FROM managed_schema_operations`).Scan(&before))
	result, err := svc.ListUnresolvedTestWrites(t.Context(), scope.WorkspaceID)
	require.NoError(t, err)
	require.True(t, result.Complete)
	require.Len(t, result.Operations, 2)
	byID := map[string]TestWriteRecoveryOperation{}
	for _, item := range result.Operations {
		byID[item.OperationID] = item
	}
	require.Equal(t, "group-1", byID[first.OperationID].GroupID)
	require.Equal(t, "rev-1", byID[first.OperationID].GroupRevision)
	require.Equal(t, "identity-1", byID[first.OperationID].ConnectorRevision)
	require.Equal(t, "group-2", byID[second.OperationID].GroupID)
	require.NotEqual(t, byID[first.OperationID].TargetFingerprint, byID[second.OperationID].TargetFingerprint)
	encoded, err := json.Marshal(result)
	require.NoError(t, err)
	for _, secret := range []string{firstToken.Token, secondToken.Token, "secret-endpoint", "password=private", "private-proof", "/tmp/line.db", "scope_key", "payload_digest"} {
		require.NotContains(t, string(encoded), secret)
	}
	var after string
	require.NoError(t, reopened.QueryRow(`SELECT group_concat(operation_id || status || owner || detail) FROM managed_schema_operations`).Scan(&after))
	require.Equal(t, before, after)
	var count int
	require.NoError(t, reopened.QueryRow(`SELECT count(*) FROM managed_schema_preview_tokens`).Scan(&count))
	require.Equal(t, 2, count)
	foreign, err := svc.ListUnresolvedTestWrites(t.Context(), "foreign")
	require.NoError(t, err)
	require.Empty(t, foreign.Operations)
}

func TestTestWriteRecoveryClassificationAndIsolation(t *testing.T) {
	cases := []struct {
		status         SchemaOperationStatus
		write, cleanup string
		unresolved     bool
	}{
		{SchemaOperationRunning, "", "", true}, {SchemaOperationPending, "", "", true},
		{SchemaOperationSucceeded, "written_verified", "failed", true}, {SchemaOperationSucceeded, "written_verified", "unknown", true},
		{SchemaOperationSucceeded, "written_verified", "cleaned", false}, {SchemaOperationFailed, "failed", "not_attempted", false},
		{SchemaOperationUnknown, "unknown", "cleaned", true},
		{SchemaOperationPartial, "written_unverified", "cleaned", false},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			repo := NewMemoryRepository()
			svc := NewService(repo)
			_, op := seedRecovery(t, svc, testWriteScope(), SchemaOperationRunning, "", "")
			stored := repo.operations[op.OperationID]
			stored.Status = tc.status
			stored.WriteOutcome = tc.write
			stored.CleanupStatus = tc.cleanup
			if !stored.Status.IsActive() {
				now := time.Now().UTC()
				stored.CompletedAt = &now
			}
			repo.operations[op.OperationID] = stored
			// Unrelated records with absent bindings must never enter this workspace query.
			repo.operations["schema"] = SchemaOperation{OperationID: "schema", WorkspaceID: "ws-1", Action: SchemaApplyAction}
			repo.operations["foreign"] = SchemaOperation{OperationID: "foreign", WorkspaceID: "foreign", Action: TestWriteAction}
			got, err := svc.ListUnresolvedTestWrites(t.Context(), "ws-1")
			require.NoError(t, err)
			require.Len(t, got.Operations, map[bool]int{true: 1, false: 0}[tc.unresolved])
			if !tc.unresolved {
				b, err := json.Marshal(got)
				require.NoError(t, err)
				require.Contains(t, string(b), `"operations":[]`)
			}
		})
	}
}

func TestTestWriteRecoveryCorruptionFailsClosed(t *testing.T) {
	for _, corruption := range []string{"missing", "foreign", "group", "digest", "scope", "payload", "status", "write", "cleanup", "operation", "kind", "json"} {
		t.Run(corruption, func(t *testing.T) {
			db := openLedgerDB(t, filepath.Join(t.TempDir(), "ledger.db"))
			applyRecordingDDL(t, db)
			svc := NewService(NewSQLRepository(db))
			token, _ := seedRecovery(t, svc, testWriteScope(), SchemaOperationRunning, "", "")
			queries := map[string]string{
				"missing": `DELETE FROM managed_schema_preview_tokens`, "foreign": `UPDATE managed_schema_preview_tokens SET workspace_id='foreign'`,
				"group": `UPDATE managed_schema_preview_tokens SET plan_id='altered'`, "digest": `UPDATE managed_schema_preview_tokens SET digest='bad'`,
				"scope": `UPDATE managed_schema_operations SET scope_key='bad'`, "payload": `UPDATE managed_schema_operations SET payload_digest='bad'`,
				"status": `UPDATE managed_schema_operations SET status='private-endpoint'`, "write": `UPDATE managed_schema_operations SET write_outcome='private-endpoint'`,
				"cleanup": `UPDATE managed_schema_operations SET cleanup_status='private-endpoint'`, "operation": `UPDATE managed_schema_preview_tokens SET operation_id='wrong'`,
				"kind": `UPDATE managed_schema_preview_tokens SET action='schema_apply'`, "json": `UPDATE managed_schema_preview_tokens SET statements='{'`,
			}
			_, err := db.Exec(queries[corruption])
			require.NoError(t, err)
			got, err := svc.ListUnresolvedTestWrites(t.Context(), "ws-1")
			require.ErrorIs(t, err, ErrTestWriteRecoveryUnavailable)
			require.Nil(t, got)
			require.False(t, strings.Contains(err.Error(), token.Token))
			require.NotContains(t, err.Error(), "private-endpoint")
		})
	}
}

func TestTestWriteRecoveryBoundIsHonest(t *testing.T) {
	for _, sqlRepo := range []bool{false, true} {
		t.Run(fmt.Sprint(sqlRepo), func(t *testing.T) {
			var repo Repository = NewMemoryRepository()
			if sqlRepo {
				db := openLedgerDB(t, filepath.Join(t.TempDir(), "ledger.db"))
				applyRecordingDDL(t, db)
				repo = NewSQLRepository(db)
			}
			svc := NewService(repo)
			scope := testWriteScope()
			for i := 0; i < 1000; i++ {
				scope.Table = fmt.Sprint("table_", i)
				seedRecovery(t, svc, scope, SchemaOperationRunning, "", "")
			}
			result, err := svc.ListUnresolvedTestWrites(t.Context(), scope.WorkspaceID)
			require.NoError(t, err)
			require.Len(t, result.Operations, 1000)
			scope.Table = "over_limit"
			seedRecovery(t, svc, scope, SchemaOperationRunning, "", "")
			result, err = svc.ListUnresolvedTestWrites(t.Context(), scope.WorkspaceID)
			require.ErrorIs(t, err, ErrTestWriteRecoveryUnavailable)
			require.Nil(t, result)
		})
	}
}

func TestTestWriteRecoveryOptionalFacts(t *testing.T) {
	svc := NewService(NewMemoryRepository())
	_, op := seedRecovery(t, svc, testWriteScope(), SchemaOperationRunning, "", "")
	result, err := svc.ListUnresolvedTestWrites(t.Context(), "ws-1")
	require.NoError(t, err)
	bytes, err := json.Marshal(result.Operations[0])
	require.NoError(t, err)
	require.NotContains(t, string(bytes), "write_outcome")
	require.NotContains(t, string(bytes), "cleanup_status")
	require.NotContains(t, string(bytes), "completed_at")
	_, err = svc.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded, PayloadDigest: "content", WriteOutcome: "written_verified", CleanupStatus: "failed"})
	require.NoError(t, err)
	result, err = svc.ListUnresolvedTestWrites(t.Context(), "ws-1")
	require.NoError(t, err)
	require.NotNil(t, result.Operations[0].CompletedAt)
}

func TestTestWriteRecoveryCorruptTerminalCannotDisappear(t *testing.T) {
	for _, sqlRepo := range []bool{false, true} {
		for _, corruption := range []string{"missing-completion", "failed-cleaned", "wrong-status", "partial-missing-completion"} {
			t.Run(fmt.Sprint(sqlRepo, "/", corruption), func(t *testing.T) {
				memory := NewMemoryRepository()
				var repo Repository = memory
				var db *sql.DB
				if sqlRepo {
					db = openLedgerDB(t, filepath.Join(t.TempDir(), "ledger.db"))
					applyRecordingDDL(t, db)
					repo = NewSQLRepository(db)
				}
				svc := NewService(repo)
				_, op := seedRecovery(t, svc, testWriteScope(), SchemaOperationSucceeded, "written_verified", "cleaned")
				if sqlRepo {
					queries := map[string]string{"partial-missing-completion": `UPDATE managed_schema_operations SET status='partial',write_outcome='written_unverified',completed_at=NULL`, "missing-completion": `UPDATE managed_schema_operations SET completed_at=NULL`, "failed-cleaned": `UPDATE managed_schema_operations SET status='failed',write_outcome='failed'`, "wrong-status": `UPDATE managed_schema_operations SET status='failed'`}
					_, err := db.Exec(queries[corruption])
					require.NoError(t, err)
				} else {
					value := memory.operations[op.OperationID]
					switch corruption {
					case "partial-missing-completion":
						value.Status = SchemaOperationPartial
						value.WriteOutcome = "written_unverified"
						value.CompletedAt = nil
					case "missing-completion":
						value.CompletedAt = nil
					case "failed-cleaned":
						value.Status = SchemaOperationFailed
						value.WriteOutcome = "failed"
					case "wrong-status":
						value.Status = SchemaOperationFailed
					}
					memory.operations[op.OperationID] = value
				}
				got, err := svc.ListUnresolvedTestWrites(t.Context(), "ws-1")
				require.ErrorIs(t, err, ErrTestWriteRecoveryUnavailable)
				require.Nil(t, got)
			})
		}
	}
}
