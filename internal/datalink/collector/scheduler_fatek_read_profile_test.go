package collector

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"testing"
)

func TestFatekResolvedProfileAcrossProductionEntrypoints(t *testing.T) {
	for _, dtype := range []schema.DataType{schema.DataTypeUint32, schema.DataTypeUint64} {
		for _, address := range []string{"D3", "DR3"} {
			for _, mode := range []schema.PointFormatMode{schema.PointFormatOverride, schema.PointFormatInherit} {
				for _, entry := range []string{"scheduled", "manual", "direct"} {
					t.Run(address+"/"+string(dtype)+"/"+string(mode)+"/"+entry, func(t *testing.T) {
						words := []uint16{0x1122, 0x3344}
						var want any = uint32(0x11223344)
						if dtype == schema.DataTypeUint64 {
							words = append(words, 0x5566, 0x7788)
							want = uint64(0x1122334455667788)
						}
						payload := "33441122"
						if dtype == schema.DataTypeUint64 {
							payload = "7788556633441122"
						}
						config, requests := newFatekScalarLoopback(t, payload)
						connector.Register(schema.ProtocolFatekFBs, adapters.NewFatekConnector)
						manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
						t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
						d := &schema.Device{ID: "F", Protocol: schema.ProtocolFatekFBs, ConnectionConfig: config}
						p := &schema.Point{ID: "P", DeviceID: d.ID, Address: address, DataType: dtype, FormatMode: mode, FormatVersion: 1, RevisionID: "point-1", Enabled: true}
						if mode == schema.PointFormatOverride {
							p.FormatOverride = "CDAB"
						}
						profile, err := connector.ResolveReadProfile(*p, decoding.FormatContext{Protocol: d.Protocol, ConnectionFormat: "ABCD", ConnectionRevision: connector.ReadConnectionRevision(config), RuleOwned: true, RuleFormat: "CDAB", RuleRevision: "rule-1"})
						require.NoError(t, err)
						cfg := DefaultSchedulerConfig()
						cfg.DefaultRetryCount = 0
						s := NewScheduler(cfg, manager)
						s.AddDevice(d)
						require.NotZero(t, s.AddPointWithObservationLease(p, profile))
						var value any
						switch entry {
						case "scheduled":
							conn, err := manager.GetOrCreate(t.Context(), d.ID, d.Protocol, config)
							require.NoError(t, err)
							require.NoError(t, s.pollPointWithAcquisition(t.Context(), conn, s.pointInfos[p.ID], "scheduled"))
							value = (<-s.ValueChannel()).Value
						case "manual":
							values := s.PollNowContext(t.Context(), []string{p.ID})
							require.Len(t, values, 1)
							require.Equal(t, schema.QualityGood, values[0].Quality)
							value = values[0].Value
						case "direct":
							repo := device.NewMemoryRepository()
							require.NoError(t, repo.Create(t.Context(), d))
							result, err := device.NewService(repo, manager).ReadPointNow(t.Context(), p, profile)
							require.NoError(t, err)
							value = result.Value
						}
						require.Equal(t, want, value)
						request := <-requests
						count := len(words)
						if address == "DR3" {
							count /= 2
						}
						require.Equal(t, count, request.Count)
						require.Equal(t, "46", request.Command)
					})
				}
			}
		}
	}

}
