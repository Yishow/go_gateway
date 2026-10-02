package grouptestwrite

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
)

func requireUnsupported(t *testing.T, err error, reason string) {
	t.Helper()
	var unsupported *UnsupportedError
	if !errors.As(err, &unsupported) || unsupported.Reason != reason {
		t.Fatalf("want unsupported %q, got %v", reason, err)
	}
}

func TestWritePreviewIsReadOnlyAndScoped(t *testing.T) {
	h := newHarness(t, harnessOptions{})
	preview := h.preview()
	if preview.Action != recordingplan.TestWriteAction || preview.OperationID == "" || preview.Token == "" ||
		preview.OwnerColumn != "entity" || preview.Target.Table != "readings" || len(preview.Values) != 5 {
		t.Fatalf("preview must describe the exact row and its target: %+v", preview)
	}
	if h.count(`1 = 1`) != 1 {
		t.Fatal("a preview must not insert anything into the target")
	}
	if _, err := h.ledger.GetSchemaOperation(t.Context(), testWorkspaceID, preview.OperationID); !errors.Is(err, recordingplan.ErrSchemaOperationNotFound) {
		t.Fatalf("a preview must not claim an operation: %v", err)
	}
	if h.dest.opened.Load() != 0 {
		t.Fatal("a preview reads real table metadata through the inspector, never through a write connection")
	}
}

func TestWritePreviewRejectsUnsupportedOwnershipAndTargets(t *testing.T) {
	t.Run("no entity key column", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.groups.group.RowPolicy.EntityKeyColumn = ""
		_, err := h.svc.Preview(t.Context(), h.ws, testGroupID)
		requireUnsupported(t, err, ReasonOwnershipUnsafe)
		if h.count(`1 = 1`) != 1 {
			t.Fatal("an unsupported preview mutates nothing")
		}
	})
	t.Run("table missing", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.groups.group.Destination.TableName = "nope"
		_, err := h.svc.Preview(t.Context(), h.ws, testGroupID)
		requireUnsupported(t, err, ReasonTableMissing)
	})
	t.Run("column missing", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.groups.group.Members[0].TargetColumn = "ghost"
		_, err := h.svc.Preview(t.Context(), h.ws, testGroupID)
		requireUnsupported(t, err, ReasonLayoutBlocked)
	})
	t.Run("receipt table missing", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.groups.group.WritePolicy.DedupeCapability = "receipt"
		_, err := h.svc.Preview(t.Context(), h.ws, testGroupID)
		requireUnsupported(t, err, ReasonReceiptTableMissing)
	})
	t.Run("connector revision changed", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		h.dest.revision = "identity-2"
		if _, err := h.svc.Preview(t.Context(), h.ws, testGroupID); !errors.Is(err, ErrDestinationChanged) {
			t.Fatalf("a moved connector revision is a stale destination, got %v", err)
		}
	})
	t.Run("unknown group", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		if _, err := h.svc.Preview(t.Context(), h.ws, "ghost"); !errors.Is(err, ErrGroupNotFound) {
			t.Fatalf("got %v", err)
		}
	})
}

func TestWriteConfirmWritesReadsBackAndCleansOnlyItsOwnRow(t *testing.T) {
	for name, receipt := range map[string]bool{"no dedupe": false, "receipt dedupe": true} {
		t.Run(name, func(t *testing.T) {
			h := newHarness(t, harnessOptions{receipt: receipt})
			preview := h.preview()
			outcome, err := h.confirm(preview)
			if err != nil {
				t.Fatalf("confirm: %v", err)
			}
			op := outcome.Operation
			if outcome.Claim != recordingplan.ClaimAcquired || op.WriteOutcome != WriteVerified || op.CleanupStatus != CleanupCleaned ||
				op.Status != recordingplan.SchemaOperationSucceeded {
				t.Fatalf("a full-chain test write is verified and cleaned: %+v", op)
			}
			if h.count(`entity LIKE 'gw-test-%'`) != 0 {
				t.Fatal("the test row must be gone after cleanup")
			}
			if note, temperature := h.productionRow(); note != "production" || temperature != 20.5 || h.count(`1 = 1`) != 1 {
				t.Fatalf("the production neighbor must be unchanged: note=%q temperature=%v", note, temperature)
			}
			if receipt {
				var receipts int
				if err := h.probe.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM gw_effect_receipts`).Scan(&receipts); err != nil || receipts != 0 {
					t.Fatalf("the test's destination receipt is cleaned with its row: count=%d err=%v", receipts, err)
				}
			}
		})
	}
}

func TestWriteResubmissionReturnsRetainedResultWithoutWritingAgain(t *testing.T) {
	h := newHarness(t, harnessOptions{})
	preview := h.preview()
	first, err := h.confirm(preview)
	if err != nil || first.Operation.WriteOutcome != WriteVerified {
		t.Fatalf("first confirm: %+v err=%v", first, err)
	}
	opened := h.dest.opened.Load()

	// Even after the token itself has expired, the saved result is returned.
	h.expireToken(preview.Token)
	again, err := h.confirm(preview)
	if err != nil || again.Claim != recordingplan.ClaimCompleted || again.Operation.WriteOutcome != WriteVerified {
		t.Fatalf("a retained result is returned (200) even for an expired token: %+v err=%v", again, err)
	}
	if h.dest.opened.Load() != opened || h.count(`1 = 1`) != 1 {
		t.Fatal("a retried confirmation must not touch the target again")
	}
}

func (h *harness) expireToken(token string) {
	h.t.Helper()
	stored, err := h.ledger.PreviewTokenForWorkspace(h.t.Context(), testWorkspaceID, token)
	if err != nil {
		h.t.Fatal(err)
	}
	stored.ExpiresAt = time.Now().UTC().Add(-time.Minute)
	if err := h.repo.SavePreviewToken(h.t.Context(), stored); err != nil {
		h.t.Fatal(err)
	}
}

func TestWriteSafeRejectionsHappenBeforeAnyTargetMutation(t *testing.T) {
	t.Run("schema token", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		schemaToken, err := h.ledger.PrepareSchemaPreview(t.Context(), recordingplan.SchemaPreviewScope{
			WorkspaceID: testWorkspaceID, WorkspaceRevision: testRevision, PlanID: "plan-1", PlanRevision: "r1",
			ConnectorID: testConnector, ConnectorRevision: testConnRev, Dialect: "sqlite", Database: h.path, Schema: "main", TablePrefix: "gw_record_",
		}, func(_ context.Context, _ string) (recordingplan.TargetTableInspection, error) {
			return recordingplan.TargetTableInspection{Status: "missing"}, nil
		})
		if err != nil {
			t.Fatal(err)
		}
		_, err = h.svc.Confirm(t.Context(), h.ws, Confirmation{Token: schemaToken.Token, OperationID: schemaToken.OperationID})
		if !errors.Is(err, recordingplan.ErrPreviewTokenKind) {
			t.Fatalf("a schema token must not authorize a test write, got %v", err)
		}
		if h.dest.opened.Load() != 0 || h.count(`1 = 1`) != 1 {
			t.Fatal("nothing may reach the target")
		}
	})
	t.Run("stale group revision", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		preview := h.preview()
		h.groups.group.Revision = "rev-2"
		_, err := h.confirm(preview)
		if !errors.Is(err, recordingplan.ErrPreviewTokenStale) {
			t.Fatalf("an edited group makes the preview stale, got %v", err)
		}
		if h.dest.opened.Load() != 0 {
			t.Fatal("a stale confirmation must not open the target")
		}
	})
	t.Run("changed test content", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		preview := h.preview()
		h.groups.group.Members[0].TargetColumn = "note"
		_, err := h.confirm(preview)
		if err == nil {
			t.Fatal("changed test content must not be written under the old confirmation")
		}
		if h.dest.opened.Load() != 0 || h.count(`entity LIKE 'gw-test-%'`) != 0 {
			t.Fatal("nothing may reach the target")
		}
	})
	t.Run("workspace revision moved", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		preview := h.preview()
		h.ws.Revision = "setup-2"
		if _, err := h.confirm(preview); !errors.Is(err, recordingplan.ErrPreviewTokenStale) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("expired and unclaimed", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		preview := h.preview()
		h.expireToken(preview.Token)
		if _, err := h.confirm(preview); !errors.Is(err, recordingplan.ErrPreviewTokenExpired) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("foreign workspace and unknown token", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		preview := h.preview()
		if _, err := h.svc.Confirm(t.Context(), Workspace{ID: "ws-other", Revision: testRevision}, Confirmation{Token: preview.Token, OperationID: preview.OperationID}); !errors.Is(err, recordingplan.ErrPreviewTokenNotFound) {
			t.Fatalf("a foreign workspace sees an unknown token, got %v", err)
		}
		if _, err := h.svc.Confirm(t.Context(), h.ws, Confirmation{Token: "tok-unknown", OperationID: "op-x"}); !errors.Is(err, recordingplan.ErrPreviewTokenNotFound) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("operation of another token", func(t *testing.T) {
		h := newHarness(t, harnessOptions{})
		preview := h.preview()
		if _, err := h.svc.Confirm(t.Context(), h.ws, Confirmation{Token: preview.Token, OperationID: "op-someone-else"}); !errors.Is(err, recordingplan.ErrSchemaOperationMismatch) {
			t.Fatalf("got %v", err)
		}
	})
}

func TestWriteRunningDuplicateAndBusyScope(t *testing.T) {
	h := newHarness(t, harnessOptions{})
	preview := h.preview()
	other := h.preview()
	h.dest.gate = make(chan struct{})

	done := make(chan *Outcome, 1)
	go func() {
		outcome, err := h.confirm(preview)
		if err != nil {
			t.Errorf("first confirm: %v", err)
		}
		done <- outcome
	}()
	waitForClaim(t, h, preview.OperationID)

	dup, err := h.confirm(preview)
	if err != nil || dup.Claim != recordingplan.ClaimInProgress || dup.Operation.OperationID != preview.OperationID {
		t.Fatalf("a duplicate while running is in progress (202): %+v err=%v", dup, err)
	}
	busy, err := h.confirm(other)
	if err != nil || busy.Claim != recordingplan.ClaimScopeBusy || busy.Operation.OperationID != preview.OperationID {
		t.Fatalf("a different operation on the same table is busy (409) and names the holder: %+v err=%v", busy, err)
	}
	close(h.dest.gate)
	if first := <-done; first == nil || first.Operation.WriteOutcome != WriteVerified {
		t.Fatalf("the first confirmation completes normally: %+v", first)
	}
	if h.count(`1 = 1`) != 1 {
		t.Fatal("exactly one test row was written and cleaned")
	}
	if retry, err := h.confirm(other); err != nil || retry.Operation.WriteOutcome != WriteVerified {
		t.Fatalf("the other operation can run once the table is free: %+v err=%v", retry, err)
	}
}

func waitForClaim(t *testing.T, h *harness, operationID string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if op, err := h.ledger.GetSchemaOperation(t.Context(), testWorkspaceID, operationID); err == nil && op.Status.IsActive() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("the first confirmation never claimed its operation")
}

func TestWriteConfirmationMustAddressTheGroupOfItsPreview(t *testing.T) {
	h := newHarness(t, harnessOptions{})
	preview := h.preview()
	_, err := h.svc.Confirm(t.Context(), h.ws, Confirmation{Token: preview.Token, OperationID: preview.OperationID, GroupID: "group-other"})
	if !errors.Is(err, ErrGroupMismatch) {
		t.Fatalf("a token cannot be used through another group's address, got %v", err)
	}
	if h.dest.opened.Load() != 0 || h.count(`1 = 1`) != 1 {
		t.Fatal("a mismatched address must not touch the target")
	}
	outcome, err := h.svc.Confirm(t.Context(), h.ws, Confirmation{Token: preview.Token, OperationID: preview.OperationID, GroupID: testGroupID})
	if err != nil || outcome.Operation.WriteOutcome != WriteVerified {
		t.Fatalf("the matching group address works: %+v err=%v", outcome, err)
	}
	// A retained result is not served through another group's address either.
	if _, err := h.svc.Confirm(t.Context(), h.ws, Confirmation{Token: preview.Token, OperationID: preview.OperationID, GroupID: "group-other"}); !errors.Is(err, ErrGroupMismatch) {
		t.Fatalf("got %v", err)
	}
}
