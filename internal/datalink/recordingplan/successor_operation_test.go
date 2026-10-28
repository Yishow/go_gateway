package recordingplan

import (
	"database/sql"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func successorReserveFixture() *SuccessorReserveRecord {
	return &SuccessorReserveRecord{
		Scope:               SuccessorOperationScope{WorkspaceID: "W1", RequestID: "O2", LineageID: "L1", PredecessorGroupID: "G1", SuccessorGroupID: "G2"},
		IntentDigest:        strings.Repeat("a", 64),
		Intent:              json.RawMessage(`{"version":"reserve/v1","operation":"O2","group":"G2","expected_revision":"R1"}`),
		Result:              json.RawMessage(`{"workspace_revision":"R2","group":{"id":"G2","generation":2}}`),
		SchemaSourceBinding: json.RawMessage(`{"source_digest":"source-1","group":"G2"}`),
		SourceEvidence:      json.RawMessage(`{"no_history":true,"source_family":"sqlite"}`),
	}
}

func beginSuccessorTx(t *testing.T, db *sql.DB) *sql.Tx {
	t.Helper()
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = tx.Rollback() })
	return tx
}

func TestSuccessorLedgerReserveSurvivesCommitAndExactReplay(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	want := successorReserveFixture()
	got, err := RecordSuccessorReserveInTx(t.Context(), tx, want)
	require.NoError(t, err, "O2 must atomically preserve original G2 result")
	require.JSONEq(t, string(want.Result), string(got.Result))
	require.NoError(t, tx.Commit())
	tx = beginSuccessorTx(t, db)
	read, err := ReadSuccessorReserveInTx(t.Context(), tx, "W1", "O2")
	require.NoError(t, err)
	require.NotNil(t, read)
	require.JSONEq(t, string(want.Intent), string(read.Intent))
	replay, err := RecordSuccessorReserveInTx(t.Context(), tx, want)
	require.NoError(t, err)
	require.JSONEq(t, string(want.Result), string(replay.Result))
	var count int
	require.NoError(t, tx.QueryRow(`SELECT COUNT(*) FROM managed_schema_operations`).Scan(&count))
	require.Equal(t, 1, count, "same O2 never creates G3 or another root")
	encoded, err := json.Marshal(read)
	require.NoError(t, err)
	require.JSONEq(t, `{}`, string(encoded), "private ledger must not become a raw DTO")
}

func TestSuccessorLedgerReserveRejectsChangedFullIntentAndResult(t *testing.T) {
	cases := map[string]func(*SuccessorReserveRecord){
		"same-digest-other-intent": func(r *SuccessorReserveRecord) {
			r.Intent = json.RawMessage(`{"version":"reserve/v1","operation":"O2","group":"G3","expected_revision":"R1"}`)
		},
		"result": func(r *SuccessorReserveRecord) {
			r.Result = json.RawMessage(`{"workspace_revision":"R9","group":{"id":"G3","generation":3}}`)
		},
		"schema-binding": func(r *SuccessorReserveRecord) {
			r.SchemaSourceBinding = json.RawMessage(`{"source_digest":"source-9","group":"G2"}`)
		},
		"source-evidence": func(r *SuccessorReserveRecord) {
			r.SourceEvidence = json.RawMessage(`{"no_history":true,"source_family":"mysql"}`)
		},
		"digest":      func(r *SuccessorReserveRecord) { r.IntentDigest = strings.Repeat("b", 64) },
		"lineage":     func(r *SuccessorReserveRecord) { r.Scope.LineageID = "L9" },
		"predecessor": func(r *SuccessorReserveRecord) { r.Scope.PredecessorGroupID = "G9" },
		"successor":   func(r *SuccessorReserveRecord) { r.Scope.SuccessorGroupID = "G3" },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			db := recordingRepositoryDB(t)
			tx := beginSuccessorTx(t, db)
			original, err := RecordSuccessorReserveInTx(t.Context(), tx, successorReserveFixture())
			require.NoError(t, err)
			changed := successorReserveFixture()
			change(changed)
			_, err = RecordSuccessorReserveInTx(t.Context(), tx, changed)
			require.ErrorIs(t, err, ErrSchemaOperationMismatch)
			retained, err := ReadSuccessorReserveInTx(t.Context(), tx, "W1", "O2")
			require.NoError(t, err)
			require.Equal(t, original, retained)
		})
	}
}

func TestSuccessorLedgerReserveRollbackLeavesNoAuthority(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	_, err := RecordSuccessorReserveInTx(t.Context(), tx, successorReserveFixture())
	require.NoError(t, err)
	// A later configuration checkpoint failure aborts the entire caller commit.
	_, err = tx.Exec(`INSERT INTO missing_configuration_checkpoint VALUES (1)`)
	require.Error(t, err)
	require.NoError(t, tx.Rollback())
	tx = beginSuccessorTx(t, db)
	record, err := ReadSuccessorReserveInTx(t.Context(), tx, "W1", "O2")
	require.NoError(t, err)
	require.Nil(t, record)
	first, err := ReadSuccessorFirstStartInTx(t.Context(), tx, "W1", SuccessorReserveOperationID("W1", "O2"))
	require.NoError(t, err)
	require.Nil(t, first)
}

func TestSuccessorLedgerCanonicalJSONRejectsAmbiguityWithoutLeakingPayload(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	original := successorReserveFixture()
	_, err := RecordSuccessorReserveInTx(t.Context(), tx, original)
	require.NoError(t, err)
	original.Intent = json.RawMessage(`{ "expected_revision":"R1", "group":"G2", "operation":"O2", "version":"reserve/v1" }`)
	_, err = RecordSuccessorReserveInTx(t.Context(), tx, original)
	require.NoError(t, err, "JSON object ordering is not a distinct normalized intent")
	for _, payload := range []string{`{"secret":"do-not-echo","secret":"other"}`, `[]`, `null`, `{"unterminated":"do-not-echo"`, `{} {}`} {
		original.Intent = json.RawMessage(payload)
		_, err = RecordSuccessorReserveInTx(t.Context(), tx, original)
		require.ErrorIs(t, err, ErrSchemaOperationResult)
		require.NotContains(t, err.Error(), "do-not-echo")
	}
}
