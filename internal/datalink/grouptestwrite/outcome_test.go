package grouptestwrite

import (
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
)

// reject installs a statement filter on the destination connection.
func (h *harness) reject(filter func(query string) error) {
	h.faults.Reject.Store(&filter)
}

func isInsertInto(query string) bool {
	return strings.HasPrefix(strings.TrimSpace(strings.ToUpper(query)), "INSERT INTO \"READINGS\"") ||
		strings.HasPrefix(strings.TrimSpace(query), `INSERT INTO "readings"`)
}

func (h *harness) mustExec(statement string) {
	h.t.Helper()
	if _, err := h.probe.ExecContext(h.t.Context(), statement); err != nil {
		h.t.Fatal(err)
	}
}

func TestTypedReadbackEvidenceRequiresTheWholePayloadToMatch(t *testing.T) {
	t.Run("same identity different value", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		// The row is found by its owner, but its stored value is not what was written.
		h.mustExec(`CREATE TRIGGER skew AFTER INSERT ON readings WHEN NEW.entity LIKE 'gw-test-%'
			BEGIN UPDATE readings SET temperature = 99.25 WHERE entity = NEW.entity; END`)
		outcome, err := h.confirm(h.preview())
		if err != nil {
			t.Fatal(err)
		}
		op := outcome.Operation
		if op.WriteOutcome != WriteUnverified || op.Reason != reasonReadbackMismatch || op.Status != recordingplan.SchemaOperationPartial {
			t.Fatalf("a value mismatch is written_unverified with a reason, never verified: %+v", op)
		}
		if op.CleanupStatus != CleanupCleaned || h.count(`entity LIKE 'gw-test-%'`) != 0 {
			t.Fatalf("an unverified row is still removed because it is operation-owned: %+v", op)
		}
	})
	t.Run("value truncated by the column", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.mustExec(`CREATE TRIGGER lossy AFTER INSERT ON readings WHEN NEW.entity LIKE 'gw-test-%'
			BEGIN UPDATE readings SET big = big - 1 WHERE entity = NEW.entity; END`)
		outcome, err := h.confirm(h.preview())
		if err != nil || outcome.Operation.WriteOutcome != WriteUnverified || outcome.Operation.Reason != reasonReadbackMismatch {
			t.Fatalf("an off-by-one integer must not be verified: %+v err=%v", outcome, err)
		}
	})
	t.Run("provenance differs", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.mustExec(`CREATE TRIGGER prov AFTER INSERT ON readings WHEN NEW.entity LIKE 'gw-test-%'
			BEGIN UPDATE readings SET prov = '[]' WHERE entity = NEW.entity; END`)
		outcome, err := h.confirm(h.preview())
		if err != nil || outcome.Operation.WriteOutcome != WriteUnverified || outcome.Operation.Reason != reasonReadbackMismatch {
			t.Fatalf("different provenance must not be verified: %+v err=%v", outcome, err)
		}
	})
	t.Run("readback not permitted", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.reject(func(query string) error {
			if strings.HasPrefix(strings.TrimSpace(query), "SELECT") && strings.Contains(query, `"readings"`) {
				return errors.New("select denied")
			}
			return nil
		})
		outcome, err := h.confirm(h.preview())
		if err != nil {
			t.Fatal(err)
		}
		op := outcome.Operation
		if op.WriteOutcome != WriteUnverified || op.Status != recordingplan.SchemaOperationPartial ||
			(op.Reason != reasonReadbackFailed && op.Reason != reasonReadbackDenied) {
			t.Fatalf("a committed write that cannot be read is written_unverified: %+v", op)
		}
		if op.CleanupStatus != CleanupUnknown || op.CleanupReason != reasonReadbackFailed {
			t.Fatalf("cleanup readback failure must remain unknown: %+v", op)
		}
		if h.count(`entity LIKE 'gw-test-%'`) != 0 {
			t.Fatalf("cleanup can remove the owned row but cannot claim it was verified: %+v", op)
		}
	})
}

func TestOwnedCleanupStatusIsIndependentOfWriteVerification(t *testing.T) {
	t.Run("delete not permitted", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.reject(func(query string) error {
			if strings.HasPrefix(strings.TrimSpace(query), "DELETE") {
				return errors.New("delete denied")
			}
			return nil
		})
		preview := h.preview()
		outcome, err := h.confirm(preview)
		if err != nil {
			t.Fatal(err)
		}
		op := outcome.Operation
		if op.WriteOutcome != WriteVerified || op.CleanupStatus != CleanupFailed || op.CleanupReason == "" {
			t.Fatalf("verified write and failed cleanup are reported separately: %+v", op)
		}
		if h.count(`entity LIKE 'gw-test-%'`) != 1 {
			t.Fatal("the row that could not be removed must still be reported as present, never as cleaned")
		}
		// A retry returns the saved result; it neither writes nor deletes again.
		h.faults.Reject.Store(nil)
		again, err := h.confirm(preview)
		if err != nil || again.Claim != recordingplan.ClaimCompleted || again.Operation.CleanupStatus != CleanupFailed {
			t.Fatalf("a retry after a failed cleanup returns the retained result: %+v err=%v", again, err)
		}
		if h.count(`entity LIKE 'gw-test-%'`) != 1 || h.count(`1 = 1`) != 2 {
			t.Fatal("a retry must not insert a second row")
		}
	})
	t.Run("cleanup commit response lost", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.faults.Lose.Store(true)
		var inserts atomic.Int32
		h.reject(func(query string) error {
			if isInsertInto(query) {
				inserts.Add(1)
			}
			return nil
		})
		outcome, err := h.confirm(h.preview())
		if err != nil {
			t.Fatal(err)
		}
		op := outcome.Operation
		if op.WriteOutcome != WriteVerified || op.CleanupStatus != CleanupUnknown {
			t.Fatalf("a lost cleanup response is unknown, not cleaned and not failed: %+v", op)
		}
		if inserts.Load() != 1 {
			t.Fatalf("the ambiguous write commit must not be inserted twice, got %d inserts", inserts.Load())
		}
	})
}

func TestWriteOutcomesWhenTheTargetRefusesOrIsAmbiguous(t *testing.T) {
	t.Run("destination refuses the row", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.mustExec(`CREATE TRIGGER refuse BEFORE INSERT ON readings WHEN NEW.entity LIKE 'gw-test-%'
			BEGIN SELECT RAISE(ABORT, 'refused'); END`)
		outcome, err := h.confirm(h.preview())
		if err != nil {
			t.Fatal(err)
		}
		op := outcome.Operation
		if op.WriteOutcome != WriteFailed || op.Status != recordingplan.SchemaOperationFailed || op.CleanupStatus != CleanupNotAttempted || op.Reason != reasonRowRejected {
			t.Fatalf("a row the destination rejects is failed with nothing to clean: %+v", op)
		}
	})
	t.Run("destination fails before commit", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		failure := errors.New("connection reset")
		h.faults.ExecError.Store(&failure)
		outcome, err := h.confirm(h.preview())
		if err != nil {
			t.Fatal(err)
		}
		if op := outcome.Operation; op.WriteOutcome != WriteFailed || op.CleanupStatus != CleanupNotAttempted || op.Reason != reasonDestinationDown {
			t.Fatalf("a pre-commit failure is failed, never unknown: %+v", op)
		}
		if h.count(`entity LIKE 'gw-test-%'`) != 0 {
			t.Fatal("nothing may have been written")
		}
	})
	t.Run("destination offline", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.dest.offline.Store(true)
		outcome, err := h.confirm(h.preview())
		if err != nil {
			t.Fatal(err)
		}
		if op := outcome.Operation; op.WriteOutcome != WriteFailed || op.CleanupStatus != CleanupNotAttempted || op.Reason != reasonDestinationDown {
			t.Fatalf("an unreachable destination fails the test write before anything is sent: %+v", op)
		}
	})
	t.Run("connector edited after preview", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		preview := h.preview()
		h.dest.revision = "identity-2"
		h.groups.group.Destination.ConnectorRevision = "identity-2"
		if _, err := h.confirm(preview); !errors.Is(err, recordingplan.ErrPreviewTokenStale) {
			t.Fatalf("an edited connector makes the preview stale and writes nothing: %v", err)
		}
		if h.dest.opened.Load() != 0 {
			t.Fatal("no connection may be opened for a stale preview")
		}
	})
	t.Run("lost commit response is resolved from the target", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.faults.Lose.Store(true)
		var inserts atomic.Int32
		h.reject(func(query string) error {
			if isInsertInto(query) {
				inserts.Add(1)
			}
			if strings.HasPrefix(strings.TrimSpace(query), "DELETE") {
				h.faults.Lose.Store(false) // only the write response is lost
			}
			return nil
		})
		outcome, err := h.confirm(h.preview())
		if err != nil {
			t.Fatal(err)
		}
		op := outcome.Operation
		if op.WriteOutcome != WriteVerified || op.CleanupStatus != CleanupCleaned || inserts.Load() != 1 {
			t.Fatalf("a committed row whose response was lost is found by readback and written once: %+v inserts=%d", op, inserts.Load())
		}
	})
	t.Run("lost commit with no row is unknown", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.faults.Lose.Store(true)
		h.reject(func(query string) error {
			if strings.HasPrefix(strings.TrimSpace(query), "SELECT") && strings.Contains(query, `"readings"`) {
				h.mustExec(`DELETE FROM readings WHERE entity LIKE 'gw-test-%'`)
			}
			return nil
		})
		outcome, err := h.confirm(h.preview())
		if err != nil {
			t.Fatal(err)
		}
		if op := outcome.Operation; op.WriteOutcome != WriteUnknown || op.Status != recordingplan.SchemaOperationUnknown || op.Reason != reasonCommitAmbiguous {
			t.Fatalf("an ambiguous commit the target cannot settle is unknown: %+v", op)
		}
	})
}
