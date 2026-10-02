package collector

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector/health"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type acquisitionFactsProtocol struct {
	timeOrigin string
	quality    schema.QualityFlag
	readErr    error
	resultErr  string
}

type retryFactsProtocol struct {
	acquisitionFactsProtocol
	calls           atomic.Int32
	completionNanos *atomic.Int64
	finalCompletion time.Time
}

func (p *retryFactsProtocol) Read(context.Context, connector.ReadRequest) (connector.ReadResult, error) {
	call := p.calls.Add(1)
	if call == 1 {
		return connector.ReadResult{Quality: schema.QualityBad, Error: "transient read failure"}, errors.New("transient read failure")
	}
	if p.completionNanos != nil {
		p.completionNanos.Store(p.finalCompletion.UnixNano())
	}
	return connector.ReadResult{Value: uint64(9007199254740993), Quality: schema.QualityGood}, nil
}

func (acquisitionFactsProtocol) Connect(context.Context, string) error { return nil }
func (acquisitionFactsProtocol) Close() error                          { return nil }
func (acquisitionFactsProtocol) IsConnected() bool                     { return true }
func (acquisitionFactsProtocol) ProtocolType() schema.ProtocolType     { return "acquisition-facts-test" }
func (acquisitionFactsProtocol) TestConnection(context.Context) error  { return nil }
func (acquisitionFactsProtocol) Write(context.Context, connector.WriteRequest) error {
	return nil
}
func (p acquisitionFactsProtocol) Read(context.Context, connector.ReadRequest) (connector.ReadResult, error) {
	timeOrigin := p.timeOrigin
	if timeOrigin == "" {
		timeOrigin = "source"
	}
	quality := p.quality
	if quality == "" {
		quality = schema.QualityGood
	}
	result := connector.ReadResult{
		Value:      uint64(9007199254740993),
		Timestamp:  time.Date(2026, 1, 1, 0, 0, 3, 0, time.UTC),
		TimeOrigin: timeOrigin,
		Quality:    quality,
		Error:      p.resultErr,
	}
	if p.readErr != nil {
		result.Quality = schema.QualityBad
		result.Error = p.readErr.Error()
	}
	return result, p.readErr
}

func TestSchedulerPollNowPreservesTypedAcquisitionFacts(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-test")
	connector.Register(protocolType, func() connector.Protocol { return acquisitionFactsProtocol{} })
	defer connector.Unregister(protocolType)

	scheduler := NewScheduler(DefaultSchedulerConfig(), connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	scheduler.WithClock(func() time.Time { return time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC) })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-A" })
	scheduler.AddDevice(&schema.Device{ID: "device-A", Protocol: protocolType, ConnectionConfig: `{"host":"device-a"}`})
	scheduler.AddPoint(&schema.Point{ID: "point-A", DeviceID: "device-A", Address: "40001", Function: "03", DataType: schema.DataTypeUint64})

	results := scheduler.PollNowContext(t.Context(), []string{"point-A"})
	require.Len(t, results, 1)
	result := results[0]
	require.Equal(t, "acquisition-A", result.AcquisitionID)
	require.Equal(t, uint64(9007199254740993), result.Value)
	require.Equal(t, time.Date(2026, 1, 1, 0, 0, 3, 0, time.UTC), result.ObservedAt)
	require.Equal(t, time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC), result.ReceivedAt)
	require.Equal(t, "source", result.TimeOrigin)
	require.NotEmpty(t, result.ConfigFingerprint)
}

func TestSchedulerUsesManagedConnectionConfigurationForFingerprint(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-managed-config-test")
	connector.Register(protocolType, func() connector.Protocol { return acquisitionFactsProtocol{} })
	defer connector.Unregister(protocolType)

	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
	scheduler.WithClock(func() time.Time { return time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC) })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-managed-config" })
	oldConfig := `{"host":"old-device-a"}`
	newConfig := `{"host":"new-device-a"}`
	scheduler.AddDevice(&schema.Device{ID: "device-managed", Protocol: protocolType, ConnectionConfig: oldConfig})
	point := &schema.Point{ID: "point-managed", DeviceID: "device-managed", Address: "40001", Function: "03", DataType: schema.DataTypeUint64, Mode: schema.PointModeReadOnly}
	scheduler.AddPoint(point)

	first := scheduler.PollNowContext(t.Context(), []string{point.ID})
	require.Len(t, first, 1)

	scheduler.AddDevice(&schema.Device{ID: "device-managed", Protocol: protocolType, ConnectionConfig: newConfig})
	second := scheduler.PollNowContext(t.Context(), []string{point.ID})
	require.Len(t, second, 1)

	oldFingerprint, err := measurement.AcquisitionConfigFingerprint(schema.Device{
		ID:               "device-managed",
		Protocol:         protocolType,
		ConnectionConfig: oldConfig,
	}, *point)
	require.NoError(t, err)
	newFingerprint, err := measurement.AcquisitionConfigFingerprint(schema.Device{
		ID:               "device-managed",
		Protocol:         protocolType,
		ConnectionConfig: newConfig,
	}, *point)
	require.NoError(t, err)
	require.NotEqual(t, oldFingerprint, newFingerprint)
	require.Equal(t, oldFingerprint, second[0].ConfigFingerprint)
}

func TestSchedulerUsesGatewayCompletionForUntrustedSourceTimestamp(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-untrusted-time-test")
	connector.Register(protocolType, func() connector.Protocol {
		return acquisitionFactsProtocol{timeOrigin: "adapter-wall-clock"}
	})
	defer connector.Unregister(protocolType)

	scheduler := NewScheduler(DefaultSchedulerConfig(), connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	completion := time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC)
	scheduler.WithClock(func() time.Time { return completion })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-untrusted-time" })
	scheduler.AddDevice(&schema.Device{ID: "device-time", Protocol: protocolType, ConnectionConfig: `{}`})
	scheduler.AddPoint(&schema.Point{ID: "point-time", DeviceID: "device-time", Address: "40001", DataType: schema.DataTypeUint64})

	results := scheduler.PollNowContext(t.Context(), []string{"point-time"})
	require.Len(t, results, 1)
	require.Equal(t, completion, results[0].ObservedAt)
	require.Equal(t, completion, results[0].ReceivedAt)
	require.Equal(t, "gateway", results[0].TimeOrigin)
}

func TestSchedulerMarksReadFailureBadWithReasonAndGatewayTime(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-error-test")
	readErr := errors.New("read timeout")
	connector.Register(protocolType, func() connector.Protocol {
		return acquisitionFactsProtocol{readErr: readErr}
	})
	defer connector.Unregister(protocolType)

	scheduler := NewScheduler(DefaultSchedulerConfig(), connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	completion := time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC)
	scheduler.WithClock(func() time.Time { return completion })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-error" })
	scheduler.AddDevice(&schema.Device{ID: "device-error", Protocol: protocolType, ConnectionConfig: `{}`})
	scheduler.AddPoint(&schema.Point{ID: "point-error", DeviceID: "device-error", Address: "40001", DataType: schema.DataTypeUint64})

	results := scheduler.PollNowContext(t.Context(), []string{"point-error"})
	require.Len(t, results, 1)
	require.Equal(t, schema.QualityBad, results[0].Quality)
	require.Equal(t, "read timeout", results[0].Error)
	require.Equal(t, "read-failed", results[0].QualityReason)
	require.Equal(t, completion, results[0].ObservedAt)
	require.Equal(t, "gateway", results[0].TimeOrigin)
}

func TestSchedulerConnectionFailureUsesSafeTypedReasonAndCompletion(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-connection-error-test")
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
	completion := time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC)
	scheduler.WithClock(func() time.Time { return completion })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-connection-error" })
	scheduler.AddDevice(&schema.Device{ID: "device-connection-error", Protocol: protocolType, ConnectionConfig: `{"password":"secret"}`})
	scheduler.AddPoint(&schema.Point{ID: "point-connection-error", DeviceID: "device-connection-error", Address: "40001", DataType: schema.DataTypeUint64})

	results := scheduler.PollNowContext(t.Context(), []string{"point-connection-error"})
	require.Len(t, results, 1)
	require.Equal(t, schema.QualityBad, results[0].Quality)
	require.Equal(t, "connection-failed", results[0].QualityReason)
	require.Contains(t, results[0].Error, "建立連線失敗")
	require.Equal(t, completion, results[0].ObservedAt)
	require.Equal(t, completion, results[0].ReceivedAt)
	require.Equal(t, "gateway", results[0].TimeOrigin)
}

func TestSchedulerCircuitBreakerUsesSafeTypedReasonAndCompletion(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-breaker-test")
	connector.Register(protocolType, func() connector.Protocol { return acquisitionFactsProtocol{} })
	defer connector.Unregister(protocolType)
	config := DefaultSchedulerConfig()
	config.BreakerConfig = health.BreakerConfig{ErrorThreshold: 0.2, UnstableThreshold: 0.5, WindowSize: 1, CooldownPeriod: time.Hour}
	scheduler := NewScheduler(config, connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	completion := time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC)
	scheduler.WithClock(func() time.Time { return completion })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-breaker" })
	device := &schema.Device{ID: "device-breaker", Protocol: protocolType, ConnectionConfig: `{}`}
	point := &schema.Point{ID: "point-breaker", DeviceID: "device-breaker", Address: "40001", DataType: schema.DataTypeUint64}
	scheduler.AddDevice(device)
	scheduler.AddPoint(point)
	scheduler.deviceBreakers["device-breaker"].ReportResult(errors.New("first failure"))
	scheduler.deviceBreakers["device-breaker"].ReportResult(errors.New("second failure"))

	results := scheduler.PollNowContext(t.Context(), []string{"point-breaker"})
	require.Len(t, results, 1)
	require.Equal(t, schema.QualityBad, results[0].Quality)
	require.Equal(t, "circuit-open", results[0].QualityReason)
	require.Contains(t, results[0].Error, "設備熔斷中")
	expectedFingerprint, err := measurement.AcquisitionConfigFingerprint(*device, *point)
	require.NoError(t, err)
	require.Equal(t, expectedFingerprint, results[0].ConfigFingerprint)
	require.Equal(t, completion, results[0].ObservedAt)
	require.Equal(t, completion, results[0].ReceivedAt)
	require.Equal(t, "gateway", results[0].TimeOrigin)
}

func TestSchedulerRetryUsesFinalReadCompletionForTypedFacts(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-retry-test")
	beforeCompletion := time.Date(2026, 1, 1, 0, 0, 7, 0, time.UTC)
	finalCompletion := time.Date(2026, 1, 1, 0, 0, 9, 0, time.UTC)
	var completionNanos atomic.Int64
	completionNanos.Store(beforeCompletion.UnixNano())
	protocol := &retryFactsProtocol{completionNanos: &completionNanos, finalCompletion: finalCompletion}
	connector.Register(protocolType, func() connector.Protocol { return protocol })
	defer connector.Unregister(protocolType)
	config := DefaultSchedulerConfig()
	config.DefaultRetryCount = 1
	config.DefaultRetryDelay = 0
	scheduler := NewScheduler(config, connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	scheduler.WithClock(func() time.Time { return time.Unix(0, completionNanos.Load()).UTC() })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-retry" })
	scheduler.AddDevice(&schema.Device{ID: "device-retry", Protocol: protocolType, ConnectionConfig: `{}`})
	scheduler.AddPoint(&schema.Point{ID: "point-retry", DeviceID: "device-retry", Address: "40001", DataType: schema.DataTypeUint64})

	scheduler.pollDevicePoints("device-retry", []pointInfo{scheduler.pointInfos["point-retry"]})
	collected := <-scheduler.ValueChannel()
	require.Equal(t, int32(2), protocol.calls.Load())
	require.Equal(t, "acquisition-retry", collected.AcquisitionID)
	require.Equal(t, schema.QualityGood, collected.Quality)
	require.Equal(t, finalCompletion, collected.ObservedAt)
	require.Equal(t, finalCompletion, collected.ReceivedAt)
	require.Equal(t, "gateway", collected.TimeOrigin)
	require.Empty(t, collected.QualityReason)
}

func TestSchedulerConnectionFailureUsesAttemptedConfigFingerprint(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-attempted-config-test")
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	scheduler := NewScheduler(DefaultSchedulerConfig(), manager)
	completion := time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC)
	scheduler.WithClock(func() time.Time { return completion })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-attempted-config" })
	device := &schema.Device{ID: "device-attempted-config", Protocol: protocolType, ConnectionConfig: `{"password":"secret"}`}
	point := &schema.Point{ID: "point-attempted-config", DeviceID: device.ID, Address: "40001", DataType: schema.DataTypeUint64}
	scheduler.AddDevice(device)
	scheduler.AddPoint(point)

	results := scheduler.PollNowContext(t.Context(), []string{point.ID})
	require.Len(t, results, 1)
	expectedFingerprint, err := measurement.AcquisitionConfigFingerprint(*device, *point)
	require.NoError(t, err)
	require.Equal(t, expectedFingerprint, results[0].ConfigFingerprint)
	require.Equal(t, schema.QualityBad, results[0].Quality)
	require.Equal(t, "connection-failed", results[0].QualityReason)
}

func TestSchedulerGoodQualityWithResultErrorIsBad(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-good-with-error-test")
	connector.Register(protocolType, func() connector.Protocol {
		return acquisitionFactsProtocol{resultErr: "read failed: password=secret"}
	})
	defer connector.Unregister(protocolType)
	scheduler := NewScheduler(DefaultSchedulerConfig(), connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	scheduler.WithClock(func() time.Time { return time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC) })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-good-with-error" })
	scheduler.AddDevice(&schema.Device{ID: "device-good-with-error", Protocol: protocolType, ConnectionConfig: `{}`})
	scheduler.AddPoint(&schema.Point{ID: "point-good-with-error", DeviceID: "device-good-with-error", Address: "40001", DataType: schema.DataTypeUint64})

	results := scheduler.PollNowContext(t.Context(), []string{"point-good-with-error"})
	require.Len(t, results, 1)
	require.Equal(t, schema.QualityBad, results[0].Quality)
	require.Equal(t, "read-failed", results[0].QualityReason)
	require.Contains(t, results[0].Error, "secret")
}

func TestSchedulerUnknownQualityWithRawErrorUsesSafeTypedReason(t *testing.T) {
	protocolType := schema.ProtocolType("acquisition-facts-unknown-quality-test")
	connector.Register(protocolType, func() connector.Protocol {
		return acquisitionFactsProtocol{quality: schema.QualityFlag("adapter-private-quality"), resultErr: "dsn=postgres://user:secret@db/password"}
	})
	defer connector.Unregister(protocolType)
	scheduler := NewScheduler(DefaultSchedulerConfig(), connector.NewConnectionManager(connector.DefaultConnectionManagerConfig()))
	completion := time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC)
	scheduler.WithClock(func() time.Time { return completion })
	scheduler.WithAcquisitionIDFactory(func() string { return "acquisition-unknown-quality" })
	scheduler.AddDevice(&schema.Device{ID: "device-unknown-quality", Protocol: protocolType, ConnectionConfig: `{}`})
	scheduler.AddPoint(&schema.Point{ID: "point-unknown-quality", DeviceID: "device-unknown-quality", Address: "40001", DataType: schema.DataTypeUint64})

	results := scheduler.PollNowContext(t.Context(), []string{"point-unknown-quality"})
	require.Len(t, results, 1)
	require.Equal(t, schema.QualityBad, results[0].Quality)
	require.Equal(t, "quality-unknown", results[0].QualityReason)
	require.NotContains(t, results[0].QualityReason, "secret")
	require.Contains(t, results[0].Error, "secret")
}
