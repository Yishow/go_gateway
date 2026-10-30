package collector

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
	"reflect"
	"testing"

	"go-gateway/internal/datalink/schema"
)

type pointObservationInstaller interface {
	HasAcquisitionObserver() bool
	AddPointWithObservationLease(*schema.Point) uint64
}

func TestSchedulerObservationClosedObserverContract(t *testing.T) {
	field, ok := reflect.TypeFor[SchedulerConfig]().FieldByName("AcquisitionObserver")
	if !ok {
		t.Fatal("scheduler has no final-result acquisition observer")
	}
	if field.Type.Kind() != reflect.Func || field.Type.NumIn() != 1 || field.Type.NumOut() != 0 {
		t.Fatal("acquisition observer must accept one closed value and return no domain outcome")
	}
	observation := field.Type.In(0)
	want := []string{"Binding", "DeviceID", "PointID", "AcquisitionID", "ConfigFingerprint", "Quality", "QualityReason", "ReceivedAt"}
	if observation.Kind() != reflect.Struct || observation.NumField() != len(want) {
		t.Fatal("observer must exclude arbitrary payload, values, raw bytes and errors")
	}
	for _, name := range want {
		if _, ok := observation.FieldByName(name); !ok {
			t.Fatalf("missing closed observation field %s", name)
		}
	}
}

func TestSchedulerObservationInstallHasMonotonicLease(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	installer, ok := any(s).(pointObservationInstaller)
	if !ok {
		t.Fatal("scheduler has no exact-install observation lease")
	}
	if installer.HasAcquisitionObserver() {
		t.Fatal("ordinary scheduler unexpectedly enables observations")
	}
	p := &schema.Point{ID: "P", DeviceID: "D", Address: "40001", DataType: schema.DataTypeUint16}
	first := installer.AddPointWithObservationLease(p)
	s.RemovePoint(p.ID)
	second := installer.AddPointWithObservationLease(p)
	if first == 0 || second <= first {
		t.Fatal("remove/re-add reused or lost the exact installation lease")
	}
}

func observationFixture(t *testing.T, s *Scheduler, protocol schema.ProtocolType, id string) (*schema.Device, *schema.Point, PointObservationBinding) {
	t.Helper()
	d := &schema.Device{ID: "D-" + id, Protocol: protocol, ConnectionConfig: `{}`}
	p := &schema.Point{ID: "P-" + id, DeviceID: d.ID, Address: "40001", DataType: schema.DataTypeUint64, Mode: schema.PointModeReadOnly}
	s.AddDevice(d)
	fingerprint, err := measurement.AcquisitionConfigFingerprint(*d, *p)
	require.NoError(t, err)
	b := PointObservationBinding{WorkspaceID: "W", WorkspaceRevision: "W1", DeviceID: d.ID, PointID: p.ID, RuleID: "S-" + id, RuleRevision: "R1", ConfigFingerprint: fingerprint}
	lease := s.AddPointWithObservationLease(p)
	require.True(t, s.BindPointObservation(p.ID, lease, b))
	return d, p, b
}

func observedResult(t *testing.T, observations <-chan AcquisitionObservation) AcquisitionObservation {
	t.Helper()
	select {
	case o := <-observations:
		return o
	default:
		t.Fatal("final normalized acquisition result was not observed")
		return AcquisitionObservation{}
	}
}

func TestSchedulerObservationBindingExactLeaseAndConfiguration(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	d, p, binding := observationFixture(t, s, "observation-binding", "A")
	first := s.pointInfos[p.ID].observationLease
	for _, mutate := range []func(*PointObservationBinding){
		func(b *PointObservationBinding) { b.PointID = "foreign" },
		func(b *PointObservationBinding) { b.DeviceID = "foreign" },
		func(b *PointObservationBinding) { b.ConfigFingerprint = "wrong" },
		func(b *PointObservationBinding) { b.WorkspaceRevision = "" },
		func(b *PointObservationBinding) { b.RuleRevision = "" },
	} {
		bad := binding
		mutate(&bad)
		require.False(t, s.BindPointObservation(p.ID, first, bad))
		require.Equal(t, binding, s.pointInfos[p.ID].observationBinding)
	}
	require.False(t, s.BindPointObservation("missing", first, binding))
	s.RemovePoint(p.ID)
	require.False(t, s.BindPointObservation(p.ID, first, binding))
	second := s.AddPointWithObservationLease(p)
	require.Greater(t, second, first)
	require.False(t, s.BindPointObservation(p.ID, first, binding))
	require.True(t, s.BindPointObservation(p.ID, second, binding))
	s.AddPoint(p)
	require.Empty(t, s.pointInfos[p.ID].observationBinding)
	require.False(t, s.BindPointObservation(p.ID, second, binding))
	third := s.AddPointWithObservationLease(p)
	d.ConnectionConfig = `{"host":"changed"}`
	s.AddDevice(d)
	require.False(t, s.BindPointObservation(p.ID, third, binding))
	require.Empty(t, s.pointInfos[p.ID].observationBinding)
}

func TestSchedulerObservationLeaseOverflowDoesNotChangeNativeInstall(t *testing.T) {
	s := NewScheduler(DefaultSchedulerConfig(), nil)
	_, p, binding := observationFixture(t, s, "observation-overflow", "A")
	s.pointObservationLease = math.MaxUint64 - 1
	require.Equal(t, uint64(math.MaxUint64), s.AddPointWithObservationLease(p))
	require.True(t, s.BindPointObservation(p.ID, math.MaxUint64, binding))
	for range 2 {
		p.Address = "40002"
		require.Zero(t, s.AddPointWithObservationLease(p))
		require.Equal(t, p.Address, s.pointInfos[p.ID].Address)
		require.Empty(t, s.pointInfos[p.ID].observationBinding)
		require.False(t, s.BindPointObservation(p.ID, 0, binding))
	}
}

func TestSchedulerObservationFinalNormalizedManualResultHasNoRawPayload(t *testing.T) {
	for _, tc := range []struct {
		name     string
		protocol acquisitionFactsProtocol
		quality  schema.QualityFlag
		reason   string
	}{
		{"good", acquisitionFactsProtocol{}, schema.QualityGood, ""},
		{"error despite good", acquisitionFactsProtocol{resultErr: "password_fixture"}, schema.QualityBad, "read-failed"},
		{"unknown quality", acquisitionFactsProtocol{quality: "raw_private_quality", resultErr: "password_fixture"}, schema.QualityBad, "quality-unknown"},
		{"read error", acquisitionFactsProtocol{readErr: errors.New("password_fixture")}, schema.QualityBad, "read-failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			protocol := schema.ProtocolType("observation-final")
			connector.Register(protocol, func() connector.Protocol { return tc.protocol })
			defer connector.Unregister(protocol)
			observations := make(chan AcquisitionObservation, 2)
			cfg := DefaultSchedulerConfig()
			cfg.AcquisitionObserver = func(o AcquisitionObservation) { observations <- o }
			cfg.Clock = func() time.Time { return time.Date(2026, 10, 29, 10, 0, 0, 0, time.UTC) }
			cfg.AcquisitionIDFactory = func() string { return "actual-acquisition" }
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			defer manager.CloseAll()
			s := NewScheduler(cfg, manager)
			_, p, binding := observationFixture(t, s, protocol, "A")
			require.True(t, s.HasAcquisitionObserver())
			values := s.PollNowContext(t.Context(), []string{p.ID})
			require.Len(t, values, 1)
			o := observedResult(t, observations)
			require.Equal(t, binding, o.Binding)
			require.Equal(t, tc.quality, o.Quality)
			require.Equal(t, tc.reason, o.QualityReason)
			require.Equal(t, values[0].ReceivedAt, o.ReceivedAt)
			require.Equal(t, values[0].AcquisitionID, o.AcquisitionID)
			require.Equal(t, values[0].ConfigFingerprint, o.ConfigFingerprint)
			require.Equal(t, uint64(9007199254740993), values[0].Value)
			encoded, err := json.Marshal(o)
			require.NoError(t, err)
			for _, secret := range []string{"password_fixture", "9007199254740993", "raw_private_quality", "RawBytes", "Error", "ConnectionConfig", "Value"} {
				require.NotContains(t, string(encoded), secret)
			}
			require.Empty(t, observations)
		})
	}
}

func TestSchedulerObservationPeriodicRetriesPublishOnlyFinalResult(t *testing.T) {
	protocol := schema.ProtocolType("observation-retry")
	adapter := &retryFactsProtocol{}
	connector.Register(protocol, func() connector.Protocol { return adapter })
	defer connector.Unregister(protocol)
	observations := make(chan AcquisitionObservation, 4)
	cfg := DefaultSchedulerConfig()
	cfg.DefaultRetryCount, cfg.DefaultRetryDelay = 1, 0
	cfg.AcquisitionObserver = func(o AcquisitionObservation) { observations <- o }
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	defer manager.CloseAll()
	s := NewScheduler(cfg, manager)
	d, p, binding := observationFixture(t, s, protocol, "A")
	s.pollDevicePoints(d.ID, []pointInfo{s.pointInfos[p.ID]})
	value := <-s.ValueChannel()
	require.Equal(t, int32(2), adapter.calls.Load())
	require.Equal(t, schema.QualityGood, value.Quality)
	o := observedResult(t, observations)
	require.Equal(t, binding, o.Binding)
	require.Equal(t, value.AcquisitionID, o.AcquisitionID)
	require.Equal(t, schema.QualityGood, o.Quality)
	require.Empty(t, observations)
}

type observationPausedProtocol struct {
	acquisitionFactsProtocol
	entered, release chan struct{}
	calls            atomic.Int32
}

func (p *observationPausedProtocol) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	if p.calls.Add(1) == 1 {
		close(p.entered)
		select {
		case <-p.release:
		case <-ctx.Done():
			return connector.ReadResult{}, ctx.Err()
		}
	}
	return p.acquisitionFactsProtocol.Read(ctx, req)
}

func TestSchedulerObservationInFlightRenameKeepsOriginalRevision(t *testing.T) {
	protocol := schema.ProtocolType("observation-paused")
	adapter := &observationPausedProtocol{entered: make(chan struct{}), release: make(chan struct{})}
	connector.Register(protocol, func() connector.Protocol { return adapter })
	defer connector.Unregister(protocol)
	observations := make(chan AcquisitionObservation, 4)
	cfg := DefaultSchedulerConfig()
	cfg.AcquisitionObserver = func(o AcquisitionObservation) { observations <- o }
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	defer manager.CloseAll()
	s := NewScheduler(cfg, manager)
	_, p, original := observationFixture(t, s, protocol, "A")
	completed := make(chan []CollectedValue, 1)
	go func() { completed <- s.PollNowContext(t.Context(), []string{p.ID}) }()
	select {
	case <-adapter.entered:
	case <-time.After(time.Second):
		t.Fatal("deferred poll did not begin")
	}
	renamed := original
	renamed.RuleRevision = "R2"
	lease := s.AddPointWithObservationLease(p) // same point/fingerprint, new source revision
	require.True(t, s.BindPointObservation(p.ID, lease, renamed))
	close(adapter.release)
	require.Len(t, <-completed, 1)
	require.Equal(t, original, observedResult(t, observations).Binding)
	require.Len(t, s.PollNowContext(t.Context(), []string{p.ID}), 1)
	require.Equal(t, renamed, observedResult(t, observations).Binding)
}

func TestSchedulerObservationConnectionAndCircuitFailureAreFinal(t *testing.T) {
	for _, circuit := range []bool{false, true} {
		t.Run(map[bool]string{false: "connection", true: "circuit"}[circuit], func(t *testing.T) {
			observations := make(chan AcquisitionObservation, 2)
			cfg := DefaultSchedulerConfig()
			cfg.AcquisitionObserver = func(o AcquisitionObservation) { observations <- o }
			cfg.BreakerConfig = health.BreakerConfig{ErrorThreshold: 0.2, UnstableThreshold: 0.5, WindowSize: 1, CooldownPeriod: time.Hour}
			s := NewScheduler(cfg, connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
			d, p, binding := observationFixture(t, s, "unregistered-observation-protocol", "A")
			if circuit {
				s.deviceBreakers[d.ID].ReportResult(errors.New("private error"))
				s.deviceBreakers[d.ID].ReportResult(errors.New("private error"))
			}
			values := s.PollNowContext(t.Context(), []string{p.ID})
			require.Len(t, values, 1)
			o := observedResult(t, observations)
			require.Equal(t, binding, o.Binding)
			require.Equal(t, map[bool]string{false: "connection-failed", true: "circuit-open"}[circuit], o.QualityReason)
			require.Equal(t, values[0].Quality, o.Quality)
			require.Empty(t, observations)
		})
	}
}

func TestSchedulerObservationNilPanicAndFullCapturePreserveNativeResult(t *testing.T) {
	for _, mode := range []string{"nil", "panic", "full"} {
		t.Run(mode, func(t *testing.T) {
			protocol := schema.ProtocolType("observation-failure")
			connector.Register(protocol, func() connector.Protocol { return acquisitionFactsProtocol{} })
			defer connector.Unregister(protocol)
			cfg := DefaultSchedulerConfig()
			var dropped atomic.Int32
			if mode == "panic" {
				cfg.AcquisitionObserver = func(AcquisitionObservation) { panic("private callback failure") }
			}
			if mode == "full" {
				full := make(chan AcquisitionObservation)
				cfg.AcquisitionObserver = func(o AcquisitionObservation) {
					select {
					case full <- o:
					default:
						dropped.Add(1)
					}
				}
			}
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			defer manager.CloseAll()
			s := NewScheduler(cfg, manager)
			_, p, _ := observationFixture(t, s, protocol, "A")
			values := s.PollNowContext(t.Context(), []string{p.ID})
			require.Len(t, values, 1)
			require.Equal(t, schema.QualityGood, values[0].Quality)
			require.Equal(t, uint64(9007199254740993), values[0].Value)
			require.Empty(t, values[0].Error)
			if mode == "full" {
				require.Equal(t, int32(1), dropped.Load())
			}
		})
	}
}

func TestSchedulerObservationTwoDevicesKeepIndependentBindings(t *testing.T) {
	protocol := schema.ProtocolType("observation-devices")
	connector.Register(protocol, func() connector.Protocol { return acquisitionFactsProtocol{readErr: errors.New("secret failure")} })
	defer connector.Unregister(protocol)
	observations := make(chan AcquisitionObservation, 4)
	cfg := DefaultSchedulerConfig()
	cfg.AcquisitionObserver = func(o AcquisitionObservation) { observations <- o }
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	defer manager.CloseAll()
	s := NewScheduler(cfg, manager)
	_, p1, b1 := observationFixture(t, s, protocol, "A")
	_, p2, b2 := observationFixture(t, s, protocol, "B")
	require.Len(t, s.PollNowContext(t.Context(), []string{p1.ID, p2.ID}), 2)
	want := map[string]PointObservationBinding{p1.ID: b1, p2.ID: b2}
	for range 2 {
		o := observedResult(t, observations)
		require.Equal(t, want[o.PointID], o.Binding)
		require.Equal(t, o.Binding.DeviceID, o.DeviceID)
		require.Equal(t, schema.QualityBad, o.Quality)
		require.Equal(t, "read-failed", o.QualityReason)
		require.False(t, strings.Contains(o.QualityReason, "secret"))
		delete(want, o.PointID)
	}
	require.Empty(t, want)
}

type observationFailProtocol struct {
	acquisitionFactsProtocol
	calls atomic.Int32
}

func (p *observationFailProtocol) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	p.calls.Add(1)
	return p.acquisitionFactsProtocol.Read(ctx, req)
}

func TestSchedulerObservationPeriodicFailurePreservesNativeError(t *testing.T) {
	protocol := schema.ProtocolType("observation-final-failure")
	failure := errors.New("private native read failure")
	adapter := &observationFailProtocol{acquisitionFactsProtocol: acquisitionFactsProtocol{readErr: failure}}
	connector.Register(protocol, func() connector.Protocol { return adapter })
	defer connector.Unregister(protocol)
	observations := make(chan AcquisitionObservation, 2)
	cfg := DefaultSchedulerConfig()
	cfg.DefaultRetryCount, cfg.DefaultRetryDelay = 1, 0
	cfg.AcquisitionObserver = func(o AcquisitionObservation) { observations <- o; panic("logging failed") }
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	defer manager.CloseAll()
	s := NewScheduler(cfg, manager)
	d, p, binding := observationFixture(t, s, protocol, "A")
	conn, err := manager.GetOrCreate(t.Context(), d.ID, d.Protocol, d.ConnectionConfig)
	require.NoError(t, err)
	err = s.pollPointWithAcquisition(t.Context(), conn, s.pointInfos[p.ID], "actual-final-failure")
	require.ErrorIs(t, err, failure)
	require.Equal(t, int32(2), adapter.calls.Load())
	value := <-s.ValueChannel()
	require.Equal(t, failure.Error(), value.Error)
	require.Equal(t, schema.QualityBad, value.Quality)
	o := observedResult(t, observations)
	require.Equal(t, binding, o.Binding)
	require.Equal(t, "actual-final-failure", o.AcquisitionID)
	require.Equal(t, "read-failed", o.QualityReason)
	require.Empty(t, observations)
}
