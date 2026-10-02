package grouptestwrite

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

var ledgerMigrations = []string{
	"019_recording_plans_sqlite.up.sql",
	"022_schema_preview_scope_sqlite.up.sql",
	"023_schema_operations_sqlite.up.sql",
	"026_operation_test_write_sqlite.up.sql",
}

// openLedger opens the gateway's operation ledger as a SQLite file, the way
// separate gateway processes share it.
func openLedger(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func applyLedgerSchema(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, name := range ledgerMigrations {
		ddl, err := os.ReadFile(filepath.Join("..", "schema", "migrations", name))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.ExecContext(t.Context(), string(ddl)); err != nil {
			t.Fatalf("apply %s: %v", name, err)
		}
	}
}

func newDurableHarness(t *testing.T, receipt bool) (*harness, *sql.DB) {
	t.Helper()
	ledgerDB := openLedger(t, filepath.Join(t.TempDir(), "ledger.db"))
	applyLedgerSchema(t, ledgerDB)
	return newHarness(t, harnessOptions{receipt: receipt, repo: recordingplan.NewSQLRepository(ledgerDB)}), ledgerDB
}

// strand claims the preview's operation as a gateway that then died: it is
// running, owned by someone else, and its lease has run out.
func (h *harness) strand(ledgerDB *sql.DB, preview *Preview, state progress) {
	h.t.Helper()
	token, err := h.ledger.PreviewTokenForWorkspace(h.t.Context(), testWorkspaceID, preview.Token)
	if err != nil {
		h.t.Fatal(err)
	}
	op, outcome, err := h.ledger.ClaimTestWrite(h.t.Context(), token)
	if err != nil || outcome != recordingplan.ClaimAcquired {
		h.t.Fatalf("claim: outcome=%q err=%v", outcome, err)
	}
	state.OwnerValue, state.Table, state.Schema = ownerValue(preview.OperationID), "readings", "main"
	state.OwnerColumn, state.Strategy = "entity", "none"
	if h.receipts {
		state.Strategy = dbtarget.GroupEffectReceipt
	}
	if state.Phase != "" {
		if err := h.ledger.SaveTestWriteProgress(h.t.Context(), op.OperationID, op.Owner, state.encode()); err != nil {
			h.t.Fatal(err)
		}
	}
	old := time.Now().UTC().Add(-2 * recordingplan.TestWriteLease)
	if _, err := ledgerDB.ExecContext(h.t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`, old, op.OperationID); err != nil {
		h.t.Fatal(err)
	}
}

// commitOwnedRow writes the operation's test row exactly as the service would,
// standing in for a write that committed before the gateway died.
func (h *harness) commitOwnedRow(preview *Preview) {
	h.t.Helper()
	p, err := h.svc.buildPlan(h.t.Context(), testGroupID)
	if err != nil {
		h.t.Fatal(err)
	}
	token, err := h.ledger.PreviewTokenForWorkspace(h.t.Context(), testWorkspaceID, preview.Token)
	if err != nil {
		h.t.Fatal(err)
	}
	row, err := p.row(ownerValue(preview.OperationID), token.CreatedAt)
	if err != nil {
		h.t.Fatal(err)
	}
	digest, err := groupdelivery.RowPayloadDigest(row)
	if err != nil {
		h.t.Fatal(err)
	}
	if _, err := dbtarget.InsertGroupRow(h.t.Context(), h.probe, dbtarget.GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindSQLite, TableName: "readings", Row: row, PayloadDigest: digest,
		Strategy: p.strategy, CommittedAt: time.Now().UTC().Format(time.RFC3339Nano),
	}); err != nil {
		h.t.Fatal(err)
	}
}

// countInserts counts INSERTs into the test table that reach the destination.
func (h *harness) countInserts() *atomic.Int32 {
	var inserts atomic.Int32
	h.reject(func(query string) error {
		if isInsertInto(query) {
			inserts.Add(1)
		}
		return nil
	})
	return &inserts
}

func TestOwnedCleanupAndRestartResumeFromRecordedProgress(t *testing.T) {
	t.Run("died after the commit before the result was saved", func(t *testing.T) {
		h, ledgerDB := newDurableHarness(t, false)
		preview := h.preview()
		h.strand(ledgerDB, preview, progress{Phase: phaseWriting})
		h.commitOwnedRow(preview)
		inserts := h.countInserts()

		outcome, err := h.confirm(preview)
		if err != nil {
			t.Fatal(err)
		}
		op := outcome.Operation
		if op.WriteOutcome != WriteVerified || op.CleanupStatus != CleanupCleaned || inserts.Load() != 0 {
			t.Fatalf("a restart reconciles from the target and never inserts again: %+v inserts=%d", op, inserts.Load())
		}
		if h.count(`1 = 1`) != 1 {
			t.Fatal("only the production row may remain")
		}
	})
	t.Run("died while cleaning with the row still there", func(t *testing.T) {
		h, ledgerDB := newDurableHarness(t, false)
		preview := h.preview()
		h.strand(ledgerDB, preview, progress{Phase: phaseCleaning, WriteOutcome: WriteVerified})
		h.commitOwnedRow(preview)
		inserts := h.countInserts()

		outcome, err := h.confirm(preview)
		if err != nil {
			t.Fatal(err)
		}
		if op := outcome.Operation; op.WriteOutcome != WriteVerified || op.CleanupStatus != CleanupCleaned || inserts.Load() != 0 || h.count(`1 = 1`) != 1 {
			t.Fatalf("the recorded verified write is kept and cleanup finishes: %+v inserts=%d", op, inserts.Load())
		}
	})
	t.Run("died after cleanup", func(t *testing.T) {
		h, ledgerDB := newDurableHarness(t, false)
		preview := h.preview()
		h.strand(ledgerDB, preview, progress{Phase: phaseCleaning, WriteOutcome: WriteVerified})
		inserts := h.countInserts()

		outcome, err := h.confirm(preview)
		if err != nil {
			t.Fatal(err)
		}
		if op := outcome.Operation; op.WriteOutcome != WriteVerified || op.CleanupStatus != CleanupCleaned || inserts.Load() != 0 || h.count(`1 = 1`) != 1 {
			t.Fatalf("cleanup that already happened is never turned into a new write: %+v inserts=%d", op, inserts.Load())
		}
	})
	t.Run("died before writing", func(t *testing.T) {
		h, ledgerDB := newDurableHarness(t, false)
		preview := h.preview()
		h.strand(ledgerDB, preview, progress{})
		outcome, err := h.confirm(preview)
		if err != nil {
			t.Fatal(err)
		}
		if op := outcome.Operation; op.WriteOutcome != WriteVerified || op.CleanupStatus != CleanupCleaned || h.count(`1 = 1`) != 1 {
			t.Fatalf("an operation that never wrote runs from the start once: %+v", op)
		}
	})
	t.Run("write phase with no row and no dedupe stays unknown", func(t *testing.T) {
		h, ledgerDB := newDurableHarness(t, false)
		preview := h.preview()
		h.strand(ledgerDB, preview, progress{Phase: phaseWriting})
		inserts := h.countInserts()

		outcome, err := h.confirm(preview)
		if err != nil {
			t.Fatal(err)
		}
		op := outcome.Operation
		if op.WriteOutcome != WriteUnknown || op.Status != recordingplan.SchemaOperationUnknown || inserts.Load() != 0 {
			t.Fatalf("without dedupe a possibly committed write is not repeated: %+v inserts=%d", op, inserts.Load())
		}
	})
	t.Run("write phase with no row and receipt dedupe retries the same effect", func(t *testing.T) {
		h, ledgerDB := newDurableHarness(t, true)
		preview := h.preview()
		h.strand(ledgerDB, preview, progress{Phase: phaseWriting})
		inserts := h.countInserts()

		outcome, err := h.confirm(preview)
		if err != nil {
			t.Fatal(err)
		}
		if op := outcome.Operation; op.WriteOutcome != WriteVerified || op.CleanupStatus != CleanupCleaned || inserts.Load() != 1 {
			t.Fatalf("the receipt makes the retry safe, so it writes once: %+v inserts=%d", op, inserts.Load())
		}
	})
	t.Run("receipt exists for a committed effect", func(t *testing.T) {
		h, ledgerDB := newDurableHarness(t, true)
		preview := h.preview()
		h.strand(ledgerDB, preview, progress{Phase: phaseWriting})
		h.commitOwnedRow(preview)
		inserts := h.countInserts()

		outcome, err := h.confirm(preview)
		if err != nil {
			t.Fatal(err)
		}
		if op := outcome.Operation; op.WriteOutcome != WriteVerified || op.CleanupStatus != CleanupCleaned || inserts.Load() != 0 {
			t.Fatalf("a committed effect with its receipt is verified and cleaned without a second insert: %+v inserts=%d", op, inserts.Load())
		}
	})
	t.Run("group edited after the write", func(t *testing.T) {
		h, ledgerDB := newDurableHarness(t, false)
		preview := h.preview()
		h.strand(ledgerDB, preview, progress{Phase: phaseWriting})
		h.commitOwnedRow(preview)
		h.groups.group.Revision = "rev-2"

		outcome, err := h.confirm(preview)
		if err != nil {
			t.Fatal(err)
		}
		op := outcome.Operation
		if op.WriteOutcome != WriteUnverified || op.Reason != reasonGroupChanged || op.CleanupStatus != CleanupCleaned {
			t.Fatalf("an edited group cannot reproduce the expected row; the owned row is reported unverified and still cleaned: %+v", op)
		}
		if h.count(`entity LIKE 'gw-test-%'`) != 0 {
			t.Fatal("the owned row must be removed")
		}
	})
	t.Run("a live lease is not adopted", func(t *testing.T) {
		h, _ := newDurableHarness(t, false)
		preview := h.preview()
		token, err := h.ledger.PreviewTokenForWorkspace(t.Context(), testWorkspaceID, preview.Token)
		if err != nil {
			t.Fatal(err)
		}
		if _, _, err := h.ledger.ClaimTestWrite(t.Context(), token); err != nil {
			t.Fatal(err)
		}
		inserts := h.countInserts()
		outcome, err := h.confirm(preview)
		if err != nil || outcome.Claim != recordingplan.ClaimInProgress || inserts.Load() != 0 {
			t.Fatalf("an operation still inside its lease is reported running and left alone: %+v err=%v", outcome, err)
		}
	})
}

func TestOwnedCleanupNeverReachesProductionNeighborsAfterRestart(t *testing.T) {
	h, ledgerDB := newDurableHarness(t, false)
	h.mustExec(`INSERT INTO readings (entity, temperature, running, batch, big, prov, note) VALUES ('gw-test-someone-else', 1, 1, 'x', 1, '[]', 'foreign test')`)
	preview := h.preview()
	h.strand(ledgerDB, preview, progress{Phase: phaseCleaning, WriteOutcome: WriteVerified})
	h.commitOwnedRow(preview)

	if outcome, err := h.confirm(preview); err != nil || outcome.Operation.CleanupStatus != CleanupCleaned {
		t.Fatalf("cleanup: %+v err=%v", outcome, err)
	}
	if h.count(`entity = 'line-a'`) != 1 || h.count(`entity = 'gw-test-someone-else'`) != 1 || strings.Contains(preview.OwnerValue, "someone-else") {
		t.Fatal("neither the production row nor another operation's test row may be touched")
	}
}

func TestUnreadableProgressIsTreatedAsAStartedWrite(t *testing.T) {
	h, ledgerDB := newDurableHarness(t, false)
	preview := h.preview()
	h.strand(ledgerDB, preview, progress{})
	token, err := h.ledger.PreviewTokenForWorkspace(t.Context(), testWorkspaceID, preview.Token)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ledgerDB.ExecContext(t.Context(), `UPDATE managed_schema_operations SET detail = ? WHERE operation_id = ?`, "{not json", token.OperationID); err != nil {
		t.Fatal(err)
	}
	inserts := h.countInserts()
	outcome, err := h.confirm(preview)
	if err != nil {
		t.Fatal(err)
	}
	if op := outcome.Operation; op.WriteOutcome != WriteUnknown || inserts.Load() != 0 {
		t.Fatalf("corrupted progress without dedupe must not write again: %+v inserts=%d", op, inserts.Load())
	}
}
