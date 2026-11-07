package main

import (
	"context"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// Explicitly selected measurement: one MC 3E peer, 32 points, one-second
// polling, one mapped tag and a locally unreachable PostgreSQL connector.
func TestProductionSnapshotWorker_MCUnreachableThroughput(t *testing.T) {
	if os.Getenv("GATEWAY_TASK13_THROUGHPUT") != "1" {
		t.Skip("explicit local throughput measurement")
	}
	connection, err := (&net.Dialer{Timeout: 200 * time.Millisecond}).DialContext(t.Context(), "tcp", "127.0.0.1:1")
	if connection != nil {
		_ = connection.Close()
	}
	require.Error(t, err, "the local PostgreSQL test destination must be unreachable")
	f, load, _ := newAcceptanceScheduledLoad(t, 1, 32, time.Millisecond)
	destination, path, connector, _ := productionSnapshotTarget(t, f)
	require.NoError(t, f.services.runtime.Start(t.Context()))
	require.Eventually(t, func() bool { return snapshotRowCount(t.Context(), destination, "first IS NOT NULL") > 0 }, 8*time.Second, 20*time.Millisecond)

	postgres := schema.DatabaseConnectorKindPostgres
	badConfig := dbtarget.ConnectionConfig{"host": "127.0.0.1", "port": "1", "user": "fixture", "database": "fixture", "sslmode": "disable"}
	connector, err = f.services.dbTarget.Update(t.Context(), connector.ID, dbtarget.UpdateConnectorRequest{
		Kind: &postgres, ConnectionConfig: &badConfig, ExpectedIdentityRevision: &connector.IdentityRevision,
	})
	require.NoError(t, err)
	var failedHistory []dbtarget.WriteHistoryRecord
	require.Eventually(t, func() bool {
		failedHistory, err = f.services.dbTarget.ListWriteHistory(t.Context(), connector.ID, 100)
		if err != nil {
			return false
		}
		for _, entry := range failedHistory {
			if entry.Status == "failed" && strings.Contains(entry.ErrorSummary, "127.0.0.1:1") {
				return true
			}
		}
		return false
	}, 10*time.Second, 20*time.Millisecond, "production worker must actually attempt the unreachable connector")
	warm := f.services.runtime.Snapshot().CollectedTotal
	require.Eventually(t, func() bool { return f.services.runtime.Snapshot().CollectedTotal >= warm+32 }, 5*time.Second, 10*time.Millisecond)

	start := time.Now()
	before := f.services.runtime.Snapshot()
	<-time.After(22 * time.Second)
	after := f.services.runtime.Snapshot()
	elapsed := time.Since(start)
	delta := after.CollectedTotal - before.CollectedTotal
	rate := float64(delta) / elapsed.Seconds()
	wire := load[0].metrics.snapshot()
	history, historyErr := f.services.dbTarget.ListWriteHistory(t.Context(), connector.ID, 100)
	failures := 0
	for _, entry := range history {
		if entry.Status == "failed" && strings.Contains(entry.ErrorSummary, "127.0.0.1:1") {
			failures++
		}
	}
	t.Logf("MC_3E points=32 cadence=1s mapped_tags=1 postgres=127.0.0.1:1 collected_before=%d collected_after=%d elapsed=%.3fs delta=%d rate=%.3f/s target=30.4/s postgres_failures=%d mc_requests=%d mc_responses=%d runtime_before=%+v runtime_after=%+v", before.CollectedTotal, after.CollectedTotal, elapsed.Seconds(), delta, rate, failures, wire.Requests, wire.Responses, before, after)

	// Recover the disposable target and prove that a live identity change is
	// observed before shutdown, so cleanup does not bypass the active backoff.
	rowsBeforeRecovery := snapshotRowCount(t.Context(), destination, "1=1")
	sqlite := schema.DatabaseConnectorKindSQLite
	goodConfig := dbtarget.ConnectionConfig{"dsn": path}
	connector, err = f.services.dbTarget.Update(t.Context(), connector.ID, dbtarget.UpdateConnectorRequest{
		Kind: &sqlite, ConnectionConfig: &goodConfig, ExpectedIdentityRevision: &connector.IdentityRevision,
	})
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		current, getErr := f.services.dbTarget.GetByID(t.Context(), connector.ID)
		return getErr == nil && current.LastWriteStatus == "success" && snapshotRowCount(t.Context(), destination, "1=1") > rowsBeforeRecovery
	}, 35*time.Second, 20*time.Millisecond, "updated connector must resume at the disposable target")
	require.NoError(t, f.services.runtime.Stop(context.Background()))
	require.NoError(t, historyErr)
	require.GreaterOrEqual(t, failures, 2, "retry attempts must be recorded during the measurement")
	require.Greater(t, wire.Responses, uint64(0), "MC requests must receive actual TCP responses")
	require.Empty(t, f.capture.acceptanceInterest.snapshot().bad)
	require.GreaterOrEqual(t, rate, 32.0*0.95)
}
