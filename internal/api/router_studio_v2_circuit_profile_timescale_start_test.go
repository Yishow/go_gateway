package api

import (
	"testing"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// Run the actual recording-start path with explicitly selected Timescale groups,
// rather than calling a plain table on a PostgreSQL server that has the extension.
func TestRecordingStartAcceptsBothTimescaleProfileGroups(t *testing.T) {
	f := newManagedDestinationFixture(t, "postgres", "uint64")
	if !timescaleInstalled(t, f) {
		t.Skip("requires a disposable TimescaleDB fixture; this case is NOT RUN without the extension")
	}
	insertCounterMeasurement(t, f, "counter-ts", "counter", "m3")
	pair := profilePair{
		f:          f,
		instant:    ensureProfileGroup(t, f, "instantaneous", 15, "second", "", "timescale"),
		cumulative: ensureProfileGroup(t, f, "cumulative", 10, "minute", "counter-ts", "timescale"),
	}
	for _, group := range []map[string]any{pair.instant, pair.cumulative} {
		id := group["id"].(string)
		require.Equal(t, "timescale", group["storage_mode"])
		require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, previewLong(t, f, id)).Status)
	}
	instantDestination := pair.instant["destination"].(map[string]any)
	cumulativeDestination := pair.cumulative["destination"].(map[string]any)
	require.Equal(t, instantDestination["table_name"], cumulativeDestination["table_name"])
	var hypertables int
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM timescaledb_information.hypertables
		WHERE hypertable_schema=$1 AND hypertable_name=$2`, f.namespace, instantDestination["table_name"]).Scan(&hypertables))
	require.Equal(t, 1, hypertables)
	svc := newRecordingStartServiceFixture(t, managedGroupFixture{writeGroupRouterFixture: f.writeGroupRouterFixture, groups: f.groups}, recordingStartSuccess())
	request := pair.request(t, "timescale-dual-basic-start")
	operation, err := svc.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, operation.Status)
	require.Len(t, operation.Groups, 2)
	for _, group := range operation.Groups {
		require.True(t, group.Saved && group.Ready && group.Applied)
	}
	require.Equal(t, 1, svc.activation.Calls())
	again, err := svc.start.Start(t.Context(), request)
	require.NoError(t, err)
	require.Equal(t, operation.OperationID, again.OperationID)
	require.Equal(t, 1, svc.activation.Calls())
}
