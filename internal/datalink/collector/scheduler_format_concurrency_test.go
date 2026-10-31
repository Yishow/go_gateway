package collector

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"sync"
	"testing"
)

func TestModernTwoPointConcurrentReadsRetainFormatFingerprint(t *testing.T) {
	config, requests := mcWordFixture(t, []uint16{0x1122, 0x3344})
	connector.Register(schema.ProtocolMC3E, adapters.NewMC3EConnector)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { manager.CloseAll() })
	s := NewScheduler(DefaultSchedulerConfig(), manager)
	d := &schema.Device{ID: "D", Protocol: schema.ProtocolMC3E, ConnectionConfig: config}
	s.AddDevice(d)
	fingerprints := map[string]string{}
	wants := map[string]uint32{"A": 0x11223344, "B": 0x33441122}
	for _, id := range []string{"A", "B"} {
		format := "ABCD"
		if id == "B" {
			format = "CDAB"
		}
		p := &schema.Point{ID: id, DeviceID: d.ID, RevisionID: "same", Address: "D0", DataType: schema.DataTypeUint32, FormatMode: schema.PointFormatOverride, FormatOverride: format, FormatVersion: 1, Enabled: true}
		profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionRevision: connector.ReadConnectionRevision(config)})
		require.NoError(t, err)
		require.NotZero(t, s.AddPointWithObservationLease(p, profile))
		fingerprints[id], err = measurement.AcquisitionConfigFingerprint(*d, *p, profile)
		require.NoError(t, err)
	}
	require.NotEqual(t, fingerprints["A"], fingerprints["B"])
	results := make(chan []CollectedValue, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); results <- s.PollNowContext(t.Context(), []string{"A", "B"}) }()
	}
	wg.Wait()
	close(results)
	for values := range results {
		require.Len(t, values, 2)
		for _, v := range values {
			require.Equal(t, schema.QualityGood, v.Quality)
			require.Equal(t, wants[v.PointID], v.Value)
			require.Equal(t, fingerprints[v.PointID], v.ConfigFingerprint)
			require.True(t, s.AcceptsInstallationSample(v.PointID, v.InstallationLease, v.InstallationScope))
		}
	}
	for i := 0; i < 4; i++ {
		require.Equal(t, 2, <-requests)
	}
}
