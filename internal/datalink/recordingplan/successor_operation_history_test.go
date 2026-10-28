package recordingplan

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func insertForeignSuccessorHistory(t *testing.T, tx *sql.Tx, id, workspaceID, raw string) {
	t.Helper()
	_, err := tx.Exec(`INSERT INTO managed_schema_operations
  (operation_id,token,workspace_id,scope_key,owner,action,status,created_at,updated_at,detail)
  VALUES(?,?,?,?,?,'successor_first_start','unknown',CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,?)`,
		id, "token-"+id, workspaceID, "scope-"+id, "owner-"+id, raw)
	require.NoError(t, err)
}

func TestSuccessorLedgerLineageHistoryRejectsForeignAttributedEvidence(t *testing.T) {
	cases := map[string]string{
		"lineage-only":              `{"scope":{"lineage_id":"L1","predecessor_group_id":"F1","successor_group_id":"F2"}}`,
		"case-insensitive-scope":    `{"Scope":{"LINEAGE_ID":"L1"}}`,
		"predecessor":               `{"scope":{"lineage_id":"foreign","predecessor_group_id":"G1","successor_group_id":"F2"}}`,
		"successor":                 `{"scope":{"lineage_id":"foreign","predecessor_group_id":"F1","successor_group_id":"G2"}}`,
		"old-generation":            `{"scope":{"successor_group_id":"G0"}}`,
		"malformed-tail":            `{"scope":{"predecessor_group_id":"G1"},"intent":`,
		"large-number-before-scope": `{"unused":1e99999,"scope":{"lineage_id":"L1"}}`,
		"duplicate-scope":           `{"scope":{"predecessor_group_id":"G1"},"scope":{"predecessor_group_id":"F1"}}`,
		"duplicate-identity":        `{"scope":{"predecessor_group_id":"G1","predecessor_group_id":"F1"}}`,
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			db := recordingRepositoryDB(t)
			tx := beginSuccessorTx(t, db)
			insertForeignSuccessorHistory(t, tx, "foreign-start", "foreign", raw)
			got, err := ReadSuccessorFirstStartsForLineageInTx(t.Context(), tx, "W1", "L1", []string{"G0", "G1", "G2"})
			require.ErrorIs(t, err, ErrSchemaOperationMismatch)
			require.Nil(t, got)
		})
	}
}

func TestSuccessorLedgerLineageHistoryRejectsForeignExactLocalRoot(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	root := successorReserveFixture()
	_, err := RecordSuccessorReserveInTx(t.Context(), tx, root)
	require.NoError(t, err)
	raw := fmt.Sprintf(`{"reserve_operation_id":%q,"scope":{"lineage_id":"foreign","request_id":"unrelated","predecessor_group_id":"F1","successor_group_id":"F2"}}`, SuccessorReserveOperationID("W1", "O2"))
	insertForeignSuccessorHistory(t, tx, "foreign-start", "foreign", raw)
	_, err = ReadSuccessorFirstStartsForLineageInTx(t.Context(), tx, "W1", "L1", []string{"G1", "G2"})
	require.ErrorIs(t, err, ErrSchemaOperationMismatch)
}

func TestSuccessorLedgerLineageHistoryIgnoresUnrelatedForeignEvidence(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	unrelated := successorReserveFixture()
	unrelated.Scope.RequestID = "other"
	unrelated.Scope.LineageID = "L9"
	unrelated.Scope.PredecessorGroupID = "G8"
	unrelated.Scope.SuccessorGroupID = "G9"
	_, err := RecordSuccessorReserveInTx(t.Context(), tx, unrelated)
	require.NoError(t, err)
	raws := []string{
		`{`,
		`{"scope":{"lineage_id":"L9","predecessor_group_id":"G11","successor_group_id":"G22","request_id":"O2"}}`,
		`{"notes":{"lineage_id":"L1","predecessor_group_id":"G1"},"intent":{"scope":{"successor_group_id":"G2"}}}`,
		`{"scope":{"lineage_id":"prefix-L1","predecessor_group_id":"G1-suffix"},"note":"G1 G2 L1"}`,
		fmt.Sprintf(`{"reserve_operation_id":%q}`, SuccessorReserveOperationID("W1", "other")),
	}
	for i, raw := range raws {
		insertForeignSuccessorHistory(t, tx, fmt.Sprintf("foreign-%d", i), "foreign", raw)
	}
	got, err := ReadSuccessorFirstStartsForLineageInTx(t.Context(), tx, "W1", "L1", []string{"G1", "G2"})
	require.NoError(t, err)
	require.Empty(t, got)
}

func TestSuccessorLedgerLineageHistoryStillValidatesEveryLocalRow(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	_, err := RecordSuccessorReserveInTx(t.Context(), tx, successorReserveFixture())
	require.NoError(t, err)
	_, err = RecordSuccessorFirstStartInTx(t.Context(), tx, successorFirstStartFixture())
	require.NoError(t, err)
	got, err := ReadSuccessorFirstStartsForLineageInTx(t.Context(), tx, "W1", "L1", []string{"G1", "G2"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	insertForeignSuccessorHistory(t, tx, "local-corrupt", "W1", `{`)
	_, err = ReadSuccessorFirstStartsForLineageInTx(t.Context(), tx, "W1", "L1", []string{"G1", "G2"})
	require.ErrorIs(t, err, ErrSchemaOperationResult)
	encoded, err := json.Marshal(got)
	require.NoError(t, err)
	require.JSONEq(t, `[{}]`, string(encoded))
}

func TestSuccessorLedgerLineageHistoryRejectsUnboundedInput(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	for _, tc := range []struct {
		lineage string
		groups  []string
	}{{"", []string{"G1"}}, {"L1", nil}, {"L1", []string{" G1"}}} {
		_, err := ReadSuccessorFirstStartsForLineageInTx(t.Context(), tx, "W1", tc.lineage, tc.groups)
		require.ErrorIs(t, err, ErrSchemaOperationResult)
	}
}
