package recordingplan

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func successorFirstStartFixture() *SuccessorFirstStartRecord {
	root := successorReserveFixture()
	return &SuccessorFirstStartRecord{Scope: root.Scope,
		ReserveOperationID: SuccessorReserveOperationID("W1", "O2"), IntentDigest: strings.Repeat("b", 64),
		Intent: json.RawMessage(`{"kind":"first_start","boundary":"2026-10-29T00:00:00Z","schema":"S2"}`),
		Result: json.RawMessage(`{"group_id":"G2","applied_revision":"A2","effective_at":"2026-10-29T00:00:00Z"}`),
	}
}

func TestSuccessorLedgerFirstStartReadIncludesAllAndRejectsUncertainEvidence(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	root := successorReserveFixture()
	_, err := RecordSuccessorReserveInTx(t.Context(), tx, root)
	require.NoError(t, err)
	first := successorFirstStartFixture()
	_, err = RecordSuccessorFirstStartInTx(t.Context(), tx, first)
	require.NoError(t, err)
	all, err := ReadSuccessorFirstStartsInTx(t.Context(), tx, "W1")
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.JSONEq(t, string(first.Result), string(all[0].Result))
	foreign, err := ReadSuccessorFirstStartsInTx(t.Context(), tx, "foreign")
	require.NoError(t, err)
	require.Empty(t, foreign)
	_, err = tx.Exec(`UPDATE managed_schema_operations SET status='unknown' WHERE action='successor_first_start'`)
	require.NoError(t, err)
	_, err = ReadSuccessorFirstStartsInTx(t.Context(), tx, "W1")
	require.ErrorIs(t, err, ErrSchemaOperationResult, "unknown activation cannot disappear from no-history proof")
}

func TestSuccessorLedgerFirstStartIsOneImmutableBoundaryPerRoot(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	root := successorReserveFixture()
	_, err := RecordSuccessorReserveInTx(t.Context(), tx, root)
	require.NoError(t, err)
	first := successorFirstStartFixture()
	got, err := RecordSuccessorFirstStartInTx(t.Context(), tx, first)
	require.NoError(t, err)
	require.JSONEq(t, string(first.Result), string(got.Result))
	require.NoError(t, tx.Commit())
	tx = beginSuccessorTx(t, db)
	replay, err := RecordSuccessorFirstStartInTx(t.Context(), tx, first)
	require.NoError(t, err)
	require.Equal(t, got, replay)
	// Keep the claimed digest unchanged: full frozen intent is still authority.
	first.Intent = json.RawMessage(`{"kind":"first_start","boundary":"2026-10-30T00:00:00Z","schema":"S2"}`)
	_, err = RecordSuccessorFirstStartInTx(t.Context(), tx, first)
	require.ErrorIs(t, err, ErrSchemaOperationMismatch)
	original, err := ReadSuccessorFirstStartInTx(t.Context(), tx, "W1", first.ReserveOperationID)
	require.NoError(t, err)
	require.Equal(t, got, original)
	var count int
	require.NoError(t, tx.QueryRow(`SELECT COUNT(*) FROM managed_schema_operations WHERE action='successor_first_start'`).Scan(&count))
	require.Equal(t, 1, count)
	encoded, err := json.Marshal(original)
	require.NoError(t, err)
	require.JSONEq(t, `{}`, string(encoded))
}

func TestSuccessorLedgerFirstStartRequiresExactExistingReserveScope(t *testing.T) {
	for _, mode := range []string{"absent-root", "foreign-root", "different-generation", "different-lineage", "different-predecessor"} {
		t.Run(mode, func(t *testing.T) {
			db := recordingRepositoryDB(t)
			tx := beginSuccessorTx(t, db)
			first := successorFirstStartFixture()
			if mode != "absent-root" {
				_, err := RecordSuccessorReserveInTx(t.Context(), tx, successorReserveFixture())
				require.NoError(t, err)
			}
			switch mode {
			case "foreign-root":
				first.ReserveOperationID = SuccessorReserveOperationID("W9", "O2")
			case "different-generation":
				first.Scope.SuccessorGroupID = "G3"
			case "different-lineage":
				first.Scope.LineageID = "L9"
			case "different-predecessor":
				first.Scope.PredecessorGroupID = "G9"
			}
			_, err := RecordSuccessorFirstStartInTx(t.Context(), tx, first)
			require.ErrorIs(t, err, ErrSchemaOperationMismatch)
		})
	}
}
