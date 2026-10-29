package recordingplan

import (
	"encoding/json"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessorLedgerDomainIDsAreStableAndUnambiguous(t *testing.T) {
	// Independently derived SHA-256 vectors for the documented JSON tuple domains.
	reserve := "successor-reserve-9c45ae6624cf7c9bc7312ae1d4642815cb727d9ea079114a899539896635da06"
	require.Equal(t, reserve, SuccessorReserveOperationID("W1", "O2"))
	require.Equal(t, "successor-first-start-ebb4e1731c61cde9d53e0eda6be74bb06ec7556f4a09a57173f6b23351d7a898", SuccessorFirstStartOperationID("W1", reserve))
	require.NotEqual(t, SuccessorReserveOperationID("a:b", "c"), SuccessorReserveOperationID("a", "b:c"))
	require.NotEqual(t, reserve, SuccessorReserveOperationID("W9", "O2"))
	require.NotEqual(t, reserve, SuccessorFirstStartOperationID("W1", reserve))
}

func TestSuccessorLedgerConcurrentConnectionsKeepOneResult(t *testing.T) {
	path := filepath.Join(t.TempDir(), "successor-ledger.db")
	first := openLedgerDB(t, path)
	applyRecordingDDL(t, first)
	second := openLedgerDB(t, path)
	const attempts = 8
	errors := make(chan error, attempts)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := range attempts {
		wg.Go(func() {
			<-start
			db := first
			if i%2 != 0 {
				db = second
			}
			tx, err := db.BeginTx(t.Context(), nil)
			if err != nil {
				errors <- err
				return
			}
			defer tx.Rollback()
			_, err = RecordSuccessorReserveInTx(t.Context(), tx, successorReserveFixture())
			if err == nil {
				err = tx.Commit()
			}
			errors <- err
		})
	}
	close(start)
	wg.Wait()
	close(errors)
	for err := range errors {
		require.NoError(t, err)
	}
	var count int
	require.NoError(t, first.QueryRow(`SELECT COUNT(*) FROM managed_schema_operations`).Scan(&count))
	require.Equal(t, 1, count)
	require.NoError(t, first.Close())
	require.NoError(t, second.Close())
	// Lost acknowledgement plus a new process connection sees the same G2 result.
	restarted := openLedgerDB(t, path)
	tx := beginSuccessorTx(t, restarted)
	result, err := ReadSuccessorReserveInTx(t.Context(), tx, "W1", "O2")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.JSONEq(t, string(successorReserveFixture().Result), string(result.Result))
	_, err = RecordSuccessorReserveInTx(t.Context(), tx, successorReserveFixture())
	require.NoError(t, err)
}

func TestSuccessorLedgerFirstStartRollbackPreservesReserveAndAllowsNewReview(t *testing.T) {
	db := recordingRepositoryDB(t)
	tx := beginSuccessorTx(t, db)
	_, err := RecordSuccessorReserveInTx(t.Context(), tx, successorReserveFixture())
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
	tx = beginSuccessorTx(t, db)
	first := successorFirstStartFixture()
	_, err = RecordSuccessorFirstStartInTx(t.Context(), tx, first)
	require.NoError(t, err)
	require.NoError(t, tx.Rollback())
	tx = beginSuccessorTx(t, db)
	read, err := ReadSuccessorFirstStartInTx(t.Context(), tx, "W1", first.ReserveOperationID)
	require.NoError(t, err)
	require.Nil(t, read)
	root, err := ReadSuccessorReserveInTx(t.Context(), tx, "W1", "O2")
	require.NoError(t, err)
	require.NotNil(t, root)
	first.Intent = json.RawMessage(`{"kind":"first_start","boundary":"2026-10-30T00:00:00Z","schema":"S2"}`)
	_, err = RecordSuccessorFirstStartInTx(t.Context(), tx, first)
	require.NoError(t, err, "only a fully rolled back attempt permits a different frozen review")
}

func TestSuccessorLedgerFirstStartListRejectsOrphanAndMalformedEvidence(t *testing.T) {
	for _, mode := range []string{"malformed", "orphan", "wrong-id", "foreign-scope", "nil-transaction"} {
		t.Run(mode, func(t *testing.T) {
			db := recordingRepositoryDB(t)
			tx := beginSuccessorTx(t, db)
			_, err := RecordSuccessorReserveInTx(t.Context(), tx, successorReserveFixture())
			require.NoError(t, err)
			_, err = RecordSuccessorFirstStartInTx(t.Context(), tx, successorFirstStartFixture())
			require.NoError(t, err)
			switch mode {
			case "malformed":
				_, err = tx.Exec(`UPDATE managed_schema_operations SET detail='{' WHERE action='successor_first_start'`)
			case "orphan":
				_, err = tx.Exec(`DELETE FROM managed_schema_operations WHERE action='successor_reserve'`)
			case "wrong-id":
				_, err = tx.Exec(`UPDATE managed_schema_operations SET operation_id='foreign' WHERE action='successor_first_start'`)
			case "foreign-scope":
				_, err = tx.Exec(`UPDATE managed_schema_operations SET scope_key='foreign' WHERE action='successor_first_start'`)
			case "nil-transaction":
				tx = nil
			}
			require.NoError(t, err)
			got, err := ReadSuccessorFirstStartsInTx(t.Context(), tx, "W1")
			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}
