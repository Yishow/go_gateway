package grouptestwrite

import (
	"context"
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// The kill test runs the real service in a child process and SIGKILLs it in the
// middle of cleanup; the parent then restarts against the same ledger and
// destination files. Only the process boundary belongs to the test.

const (
	killRoleEnv    = "GW_TESTWRITE_ROLE"
	killDestEnv    = "GW_TESTWRITE_DEST"
	killLedgerEnv  = "GW_TESTWRITE_LEDGER"
	killTokenEnv   = "GW_TESTWRITE_TOKEN"
	killOpEnv      = "GW_TESTWRITE_OP"
	killReadyEnv   = "GW_TESTWRITE_READY"
	killReceiptEnv = "GW_TESTWRITE_RECEIPT"
)

// TestTestWriteChildProcess is the child; it does nothing unless spawned.
func TestTestWriteChildProcess(t *testing.T) {
	if os.Getenv(killRoleEnv) == "" {
		t.Skip("helper process; only runs when spawned by the kill test")
	}
	probe, err := sql.Open("sqlite", os.Getenv(killDestEnv))
	require.NoError(t, err)
	probe.SetMaxOpenConns(1)
	ledgerDB := openLedger(t, os.Getenv(killLedgerEnv))
	h := attachHarness(t, os.Getenv(killDestEnv), probe, harnessOptions{
		receipt: os.Getenv(killReceiptEnv) == "1", repo: recordingplan.NewSQLRepository(ledgerDB),
	})
	// Hang at the cleanup statement: the row is committed and verified, the
	// cleanup intent is saved, and the process is about to be killed.
	h.reject(func(query string) error {
		if strings.HasPrefix(strings.TrimSpace(query), "DELETE") {
			if err := os.WriteFile(os.Getenv(killReadyEnv), []byte("ready"), 0o600); err != nil {
				return err
			}
			time.Sleep(time.Hour)
		}
		return nil
	})
	_, _ = h.svc.Confirm(context.Background(), h.ws, Confirmation{Token: os.Getenv(killTokenEnv), OperationID: os.Getenv(killOpEnv)})
}

func TestOwnedCleanupAndRestartSurvivesARealKill(t *testing.T) {
	for name, receipt := range map[string]bool{"no dedupe": false, "receipt dedupe": true} {
		t.Run(name, func(t *testing.T) {
			h, ledgerDB := newDurableHarness(t, receipt)
			ledgerPath := ledgerFilePath(t, ledgerDB)
			preview := h.preview()
			ready := filepath.Join(t.TempDir(), "ready")

			cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestTestWriteChildProcess$", "-test.v")
			receiptFlag := "0"
			if receipt {
				receiptFlag = "1"
			}
			cmd.Env = append(os.Environ(), killRoleEnv+"=confirm", killDestEnv+"="+h.path, killLedgerEnv+"="+ledgerPath,
				killTokenEnv+"="+preview.Token, killOpEnv+"="+preview.OperationID, killReadyEnv+"="+ready, killReceiptEnv+"="+receiptFlag)
			require.NoError(t, cmd.Start())
			exited := make(chan struct{})
			go func() { _ = cmd.Wait(); close(exited) }()
			t.Cleanup(func() { _ = cmd.Process.Kill(); <-exited })

			require.Eventually(t, func() bool { _, err := os.Stat(ready); return err == nil }, 30*time.Second, 20*time.Millisecond,
				"the child must reach the cleanup step")
			require.Equal(t, 2, h.count(`1 = 1`), "at the kill point the test row is committed next to the production row")
			require.NoError(t, cmd.Process.Signal(syscall.SIGKILL))
			<-exited

			op, err := h.ledger.GetSchemaOperation(t.Context(), testWorkspaceID, preview.OperationID)
			require.NoError(t, err)
			require.True(t, op.Status.IsActive(), "the killed process never recorded a result: %+v", op)
			_, err = ledgerDB.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`,
				time.Now().UTC().Add(-2*recordingplan.TestWriteLease), preview.OperationID)
			require.NoError(t, err)

			inserts := h.countInserts()
			outcome, err := h.confirm(preview)
			require.NoError(t, err)
			result := outcome.Operation
			require.Equal(t, WriteVerified, result.WriteOutcome, "%+v", result)
			require.Equal(t, CleanupCleaned, result.CleanupStatus, "%+v", result)
			require.Zero(t, inserts.Load(), "the restart reconciles; it must not write again")
			require.Equal(t, 1, h.count(`1 = 1`), "only the production row remains")
			require.Equal(t, 1, h.count(`entity = 'line-a' AND note = 'production'`))

			again, err := h.confirm(preview)
			require.NoError(t, err)
			require.Equal(t, recordingplan.ClaimCompleted, again.Claim, "a further retry returns the retained result")
			require.Zero(t, inserts.Load())
		})
	}
}

// ledgerFilePath returns the file behind a ledger handle.
func ledgerFilePath(t *testing.T, db *sql.DB) string {
	t.Helper()
	var seq int
	var name, file string
	require.NoError(t, db.QueryRowContext(t.Context(), `PRAGMA database_list`).Scan(&seq, &name, &file))
	return file
}
