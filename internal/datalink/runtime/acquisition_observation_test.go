package runtime

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

type acquisitionInstallFixture struct {
	service      *Service
	reader       *sourceObservationReader
	projection   *workspace.RuntimeProjection
	observations chan collector.AcquisitionObservation
}

func newAcquisitionInstallFixture(t *testing.T, observe bool) acquisitionInstallFixture {
	t.Helper()
	server := virtualmodbus.NewServer(memory.NewMemoryBank(64))
	require.NoError(t, server.Start(0))
	t.Cleanup(func() { require.NoError(t, server.Stop()) })
	config, err := json.Marshal(schema.ConnectionConfigModbusTCP{Host: "127.0.0.1", Port: server.Port(), SlaveID: 1, Timeout: 1})
	require.NoError(t, err)
	p := &workspace.RuntimeProjection{WorkspaceID: "W", Version: "actual-input-v1", DeviceIDs: []string{"D1", "D2"}}
	for _, id := range []string{"1", "2"} {
		p.Devices = append(p.Devices, &schema.Device{ID: "D" + id, Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: string(config)})
		p.Points = append(p.Points, &schema.Point{ID: "P" + id, DeviceID: "D" + id, Address: "4000" + id, Function: "03", DataType: schema.DataTypeInt16, Mode: schema.PointModeReadOnly, Enabled: true})
		p.Rules = append(p.Rules, &schema.SourceRule{ID: "R" + id, DeviceID: "D" + id, RevisionID: "original-rule-" + id, Enabled: true})
		p.RuleLinks = append(p.RuleLinks, &schema.SourceRuleLink{RuleID: "R" + id, PointID: "P" + id})
	}
	r := &sourceObservationReader{record: &workspace.Record{ID: "W", DatabaseSetupRevision: "observed-W-revision", OrderedDeviceIDs: []string{"D1", "D2"}}, projection: p}
	observations := make(chan collector.AcquisitionObservation, 32)
	cfg := collector.DefaultSchedulerConfig()
	cfg.DefaultRetryCount = 0
	if observe {
		cfg.AcquisitionObserver = func(v collector.AcquisitionObservation) { observations <- v }
	}
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	s := &Service{scheduler: collector.NewScheduler(cfg, manager), workspace: r, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{}}
	return acquisitionInstallFixture{s, r, p, observations}
}
func (f acquisitionInstallFixture) poll(t *testing.T, id string) collector.AcquisitionObservation {
	t.Helper()
	values := f.service.scheduler.PollNowContext(t.Context(), []string{id})
	require.Len(t, values, 1)
	require.Equal(t, schema.QualityGood, values[0].Quality)
	select {
	case observed := <-f.observations:
		return observed
	default:
		t.Fatal("actual final acquisition observation missing")
		return collector.AcquisitionObservation{}
	}
}

func TestAcquisitionInstallOriginalRuleRevisionNotFingerprintOrTagSourceRevision(t *testing.T) {
	f := newAcquisitionInstallFixture(t, true)
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	observed := f.poll(t, "P1")
	require.Equal(t, "original-rule-1", observed.Binding.RuleRevision)
	require.Equal(t, "R1", observed.Binding.RuleID)
	require.Equal(t, "W", observed.Binding.WorkspaceID)
	require.Equal(t, "observed-W-revision", observed.Binding.WorkspaceRevision)
	fingerprint, err := measurement.AcquisitionConfigFingerprint(*f.projection.Devices[0], *f.projection.Points[0])
	require.NoError(t, err)
	require.Equal(t, fingerprint, observed.Binding.ConfigFingerprint)
	require.NotEqual(t, fingerprint, observed.Binding.RuleRevision)
	require.Equal(t, "D1", observed.DeviceID)
	require.Equal(t, "P1", observed.PointID)
	require.NotEmpty(t, observed.AcquisitionID)
	require.Equal(t, 2, f.reader.reads)
	require.Zero(t, f.reader.projections, "no bootstrap-capable projection read for logging")
}

func TestAcquisitionInstallNilObserverAddsNoMetadataReads(t *testing.T) {
	f := newAcquisitionInstallFixture(t, false)
	f.reader.panicRead = true
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	require.Zero(t, f.reader.reads)
}

func TestAcquisitionInstallMissingForeignAmbiguousOrChangedProofUnavailable(t *testing.T) {
	cases := map[string]func(acquisitionInstallFixture){
		"manual point":               func(f acquisitionInstallFixture) { f.projection.RuleLinks = nil },
		"missing rule":               func(f acquisitionInstallFixture) { f.projection.Rules = nil },
		"missing revision":           func(f acquisitionInstallFixture) { f.projection.Rules[0].RevisionID = "" },
		"missing projection version": func(f acquisitionInstallFixture) { f.projection.Version = "" },
		"changed installed version": func(f acquisitionInstallFixture) {
			f.reader.beforeRead = func(n int) {
				if n == 2 {
					f.service.markDeviceProjectionAligned("D1", &workspace.RuntimeProjection{Version: "different-installed-input"})
				}
			}
		},
		"disabled source": func(f acquisitionInstallFixture) { f.projection.Rules[0].Enabled = false },
		"duplicate points": func(f acquisitionInstallFixture) {
			f.projection.Points = append(f.projection.Points, f.projection.Points[0])
		},
		"duplicate devices": func(f acquisitionInstallFixture) {
			f.projection.Devices = append(f.projection.Devices, f.projection.Devices[0])
		},
		"foreign rule": func(f acquisitionInstallFixture) { f.projection.Rules[0].DeviceID = "foreign" },
		"duplicate links": func(f acquisitionInstallFixture) {
			f.projection.RuleLinks = append(f.projection.RuleLinks, f.projection.RuleLinks[0])
		},
		"duplicate rules": func(f acquisitionInstallFixture) {
			f.projection.Rules = append(f.projection.Rules, f.projection.Rules[0])
		},
		"empty workspace revision": func(f acquisitionInstallFixture) { f.reader.record.DatabaseSetupRevision = "" },
		"foreign workspace":        func(f acquisitionInstallFixture) { f.reader.record.ID = "other" },
		"foreign member":           func(f acquisitionInstallFixture) { f.reader.record.OrderedDeviceIDs = []string{"D2"} },
		"metadata error":           func(f acquisitionInstallFixture) { f.reader.failure = errors.New("password_fixture") },
		"metadata panic":           func(f acquisitionInstallFixture) { f.reader.panicRead = true },
		"metadata deleted": func(f acquisitionInstallFixture) {
			f.reader.beforeRead = func(n int) {
				if n == 2 {
					f.reader.record = nil
				}
			}
		},
		"metadata changed": func(f acquisitionInstallFixture) {
			f.reader.beforeRead = func(n int) {
				if n == 2 {
					f.reader.record.DatabaseSetupRevision = "new"
				}
			}
		},
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			f := newAcquisitionInstallFixture(t, true)
			change(f)
			require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
			require.Empty(t, f.poll(t, "P1").Binding.RuleRevision)
			require.Zero(t, f.reader.projections)
		})
	}
}

func TestAcquisitionInstallTransitionalPollHasNoUnprovenBinding(t *testing.T) {
	f := newAcquisitionInstallFixture(t, true)
	f.reader.beforeRead = func(n int) {
		if n == 2 {
			require.Empty(t, f.poll(t, "P1").Binding.RuleRevision)
		}
	}
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	require.Equal(t, "original-rule-1", f.poll(t, "P1").Binding.RuleRevision)
}

func TestAcquisitionInstallLateOldBindingCannotOverwriteRenameOnlyRevision(t *testing.T) {
	f := newAcquisitionInstallFixture(t, true)
	f.reader.beforeRead = func(n int) {
		if n == 2 {
			newer := *f.projection
			newer.Rules = append([]*schema.SourceRule(nil), f.projection.Rules...)
			rule := *newer.Rules[0]
			rule.NamingPrefix = "Renamed"
			rule.RevisionID = "original-rule-renamed"
			newer.Rules[0] = &rule
			// Even with an unchanged external version hint, the actual install lease must win.
			require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), &newer))
		}
	}
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	require.Equal(t, "original-rule-renamed", f.poll(t, "P1").Binding.RuleRevision)
}

func TestAcquisitionInstallRemoveReaddCannotReuseOldLease(t *testing.T) {
	f := newAcquisitionInstallFixture(t, true)
	f.reader.beforeRead = func(n int) {
		if n == 2 {
			f.service.scheduler.RemovePoint("P1")
			f.service.scheduler.AddPoint(f.projection.Points[0])
		}
	}
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	require.Empty(t, f.poll(t, "P1").Binding.RuleRevision)
}

func TestAcquisitionInstallScopedPartialFailureClearsOnlyTouchedProof(t *testing.T) {
	f := newAcquisitionInstallFixture(t, true)
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	next := *f.projection
	next.Version = "partial-v2"
	next.Mappings = []*schema.Mapping{{ID: "broken", PointID: "P1", TagID: "missing", Enabled: true, TransformPipeline: "[]"}}
	require.Error(t, f.service.ApplyWorkspaceProjectionForDevices(t.Context(), &next, []string{"D1"}))
	require.Empty(t, f.poll(t, "P1").Binding.RuleRevision)
	require.Equal(t, "original-rule-2", f.poll(t, "P2").Binding.RuleRevision)
}

func TestAcquisitionInstallScopedSuccessRetainsUnselectedDeviceProof(t *testing.T) {
	f := newAcquisitionInstallFixture(t, true)
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	next := *f.projection
	next.Version = "selected-v2"
	next.Rules = append([]*schema.SourceRule(nil), f.projection.Rules...)
	rule := *next.Rules[0]
	rule.RevisionID = "selected-rule-v2"
	next.Rules[0] = &rule
	require.NoError(t, f.service.ApplyWorkspaceProjectionForDevices(t.Context(), &next, []string{"D1"}))
	require.Equal(t, "selected-rule-v2", f.poll(t, "P1").Binding.RuleRevision)
	require.Equal(t, "original-rule-2", f.poll(t, "P2").Binding.RuleRevision)
}

func TestAcquisitionInstallPointScopedSuccessRetainsSameDeviceUnselectedProof(t *testing.T) {
	f := newAcquisitionInstallFixture(t, true)
	f.projection.Points[1].DeviceID = "D1"
	f.projection.Rules[1].DeviceID = "D1"
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
	next := *f.projection
	next.Version = "one-point-v2"
	next.Rules = append([]*schema.SourceRule(nil), f.projection.Rules...)
	rule := *next.Rules[0]
	rule.RevisionID = "selected-point-rule-v2"
	next.Rules[0] = &rule
	require.NoError(t, f.service.ApplyWorkspaceProjectionForDevicesAndPoints(t.Context(), &next, []string{"D1"}, []string{"P1"}))
	require.Equal(t, "selected-point-rule-v2", f.poll(t, "P1").Binding.RuleRevision)
	require.Equal(t, "original-rule-2", f.poll(t, "P2").Binding.RuleRevision)
}
