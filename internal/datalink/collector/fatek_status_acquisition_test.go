package collector

import (
	"testing"

	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestFatekMalformedStatusPayloadAcrossEntrypoints(t *testing.T) {
	cases := []struct {
		name    string
		payload string
	}{
		{name: "empty", payload: ""},
		{name: "non-binary", payload: "X"},
		{name: "long", payload: "01"},
	}
	for _, tc := range cases {
		for _, entry := range []string{"periodic", "manual", "read-now"} {
			t.Run(tc.name+"/"+entry, func(t *testing.T) {
				config, requests := newFatekScalarLoopback(t, tc.payload)
				manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
				t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })

				configForTest := DefaultSchedulerConfig()
				configForTest.DefaultRetryCount = 0
				scheduler := NewScheduler(configForTest, manager)
				deviceRecord := &schema.Device{
					ID:               "fatek-status-acquisition",
					Protocol:         schema.ProtocolFatekFBs,
					ConnectionConfig: config,
				}
				point := &schema.Point{
					ID:       "fatek-status-point",
					DeviceID: deviceRecord.ID,
					Address:  "X0",
					DataType: schema.DataTypeBool,
				}
				scheduler.AddDevice(deviceRecord)
				scheduler.AddPoint(point)

				switch entry {
				case "periodic":
					conn, err := manager.GetOrCreate(t.Context(), deviceRecord.ID, deviceRecord.Protocol, config)
					require.NoError(t, err)
					err = scheduler.pollPointWithAcquisition(t.Context(), conn, scheduler.pointInfos[point.ID], "fatek-status-acquisition")
					require.Error(t, err)
					requireBadCollectedValue(t, <-scheduler.ValueChannel())
				case "manual":
					values := scheduler.PollNowContext(t.Context(), []string{point.ID})
					require.Len(t, values, 1)
					requireBadCollectedValue(t, values[0])
				case "read-now":
					repo := device.NewMemoryRepository()
					require.NoError(t, repo.Create(t.Context(), deviceRecord))
					result, err := device.NewService(repo, manager).ReadPointNow(t.Context(), point)
					require.Error(t, err)
					require.Equal(t, schema.QualityBad, result.Quality)
					require.Nil(t, result.Value)
					require.NotEmpty(t, result.Error)
				default:
					t.Fatalf("unknown entrypoint %q", entry)
				}

				request := <-requests
				require.Equal(t, "44", request.Command)
				require.Equal(t, 1, request.Count)
			})
		}
	}
}

func requireBadCollectedValue(t *testing.T, value CollectedValue) {
	t.Helper()
	require.Equal(t, schema.QualityBad, value.Quality)
	require.Nil(t, value.Value)
	require.NotEmpty(t, value.Error)
}
