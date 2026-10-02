package runtime

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type sampleSinkCapture struct {
	samples []measurement.SampleEnvelope
	err     error
}

func (s *sampleSinkCapture) AcceptSample(_ context.Context, sample measurement.SampleEnvelope) error {
	if s.err != nil {
		return s.err
	}
	s.samples = append(s.samples, sample)
	return nil
}

func typedRuntimeFixture(t *testing.T, sink SampleSink) (*Service, mappingBinding, *mockWriter) {
	t.Helper()
	deviceRecord := &schema.Device{ID: "device-A", Name: "Device A", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"device-a"}`}
	pointRecord := &schema.Point{ID: "point-A", DeviceID: deviceRecord.ID, Name: "40001", Address: "40001", Function: "03", DataType: schema.DataTypeUint64, Mode: schema.PointModeReadOnly, Enabled: true}
	tagRecord := &schema.Tag{ID: "tag-A", Key: "tag.a", DisplayName: "Tag A", DataType: schema.DataTypeUint64, Status: schema.TagStatusActive}
	mappingRecord := &schema.Mapping{ID: "mapping-A", PointID: pointRecord.ID, TagID: tagRecord.ID, TransformPipeline: "[]", Status: schema.MappingStatusActive, Enabled: true}
	binding, err := mappingBindingForRecords("workspace-A", mappingRecord, deviceRecord, pointRecord, tagRecord)
	require.NoError(t, err)
	writer := &mockWriter{}
	return &Service{
		config:       Config{UpdatePointState: false},
		writer:       writer,
		sampleSink:   sink,
		mappingIndex: map[string][]mappingBinding{pointRecord.ID: {binding}},
		pointMetaIndex: map[string]pointMeta{
			pointRecord.ID: {DeviceID: deviceRecord.ID, Address: pointRecord.Address},
		},
	}, binding, writer
}

func typedRuntimeValue(binding mappingBinding) collector.CollectedValue {
	acquiredAt := time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC)
	return collector.CollectedValue{
		DeviceID:          binding.DeviceID,
		PointID:           binding.PointID,
		Value:             uint64(9007199254740993),
		RawBytes:          []byte{0x01, 0x02},
		Timestamp:         time.Date(2026, 1, 1, 0, 0, 3, 0, time.UTC),
		AcquisitionID:     "acquisition-A",
		ObservedAt:        acquiredAt,
		ReceivedAt:        acquiredAt,
		TimeOrigin:        string(measurement.GatewayTimeOriginGateway),
		ConfigFingerprint: binding.ConfigFingerprint,
		Quality:           schema.QualityGood,
	}
}

func TestHandleCollectedValueAcceptsTypedSampleWithoutChangingLegacyTimestamp(t *testing.T) {
	sink := &sampleSinkCapture{}
	svc, binding, writer := typedRuntimeFixture(t, sink)
	legacyTimestamp := time.Date(2026, 1, 1, 0, 0, 3, 0, time.UTC)
	value := typedRuntimeValue(binding)
	value.Timestamp = legacyTimestamp

	svc.handleCollectedValue(t.Context(), value)
	require.Len(t, sink.samples, 1)
	sample := sink.samples[0]
	require.Equal(t, binding.WorkspaceID, sample.WorkspaceID)
	require.Equal(t, binding.SourceRevision, sample.SourceRevision)
	require.Equal(t, binding.MappingRevision, sample.MappingRevision)
	require.Equal(t, binding.ConfigFingerprint, sample.ConfigFingerprint)
	require.Equal(t, uint64(9007199254740993), sample.Value)
	require.Equal(t, value.ObservedAt, sample.ObservedAt)
	require.Equal(t, value.ReceivedAt, sample.ReceivedAt)
	require.Equal(t, string(measurement.GatewayTimeOriginGateway), sample.TimeOrigin)
	require.Equal(t, legacyTimestamp, writer.records[0].Timestamp)

	svc.handleCollectedValue(t.Context(), value)
	require.Len(t, sink.samples, 2)
	require.Equal(t, sample.SampleID, sink.samples[1].SampleID)
}

func TestHandleCollectedValueRejectsConfigFingerprintMismatchFromTypedSink(t *testing.T) {
	sink := &sampleSinkCapture{}
	svc, binding, writer := typedRuntimeFixture(t, sink)
	value := typedRuntimeValue(binding)
	value.ConfigFingerprint = "different-config"

	svc.handleCollectedValue(t.Context(), value)
	require.Empty(t, sink.samples)
	require.Len(t, writer.records, 1)
}

func TestHandleCollectedValuePreservesUnknownQualityAsBadReason(t *testing.T) {
	sink := &sampleSinkCapture{}
	svc, binding, _ := typedRuntimeFixture(t, sink)
	value := typedRuntimeValue(binding)
	value.Quality = ""

	svc.handleCollectedValue(t.Context(), value)
	require.Len(t, sink.samples, 1)
	require.Equal(t, schema.QualityBad, sink.samples[0].Quality)
	require.Equal(t, "quality-unknown", sink.samples[0].QualityReason)
}

func TestHandleCollectedValueRejectsMissingTypedIdentityAndTiming(t *testing.T) {
	sink := &sampleSinkCapture{}
	svc, binding, _ := typedRuntimeFixture(t, sink)
	value := typedRuntimeValue(binding)
	value.AcquisitionID = ""
	value.ObservedAt = time.Time{}
	value.ReceivedAt = time.Time{}

	svc.handleCollectedValue(t.Context(), value)
	require.Empty(t, sink.samples)
}

func TestHandleCollectedValueSanitizesRawReadErrorInTypedQuality(t *testing.T) {
	sink := &sampleSinkCapture{}
	svc, binding, _ := typedRuntimeFixture(t, sink)
	value := typedRuntimeValue(binding)
	value.Error = `dial tcp user:secret@10.0.0.8:5432: password=secret`
	value.Quality = schema.QualityGood

	svc.handleCollectedValue(t.Context(), value)
	require.Len(t, sink.samples, 1)
	require.Equal(t, schema.QualityBad, sink.samples[0].Quality)
	require.Equal(t, "read-failed", sink.samples[0].QualityReason)
	require.NotContains(t, sink.samples[0].QualityReason, "secret")
}

func TestHandleCollectedValuePreservesSafeConnectionFailureReason(t *testing.T) {
	sink := &sampleSinkCapture{}
	svc, binding, _ := typedRuntimeFixture(t, sink)
	value := typedRuntimeValue(binding)
	value.Quality = schema.QualityBad
	value.Error = "connection failed: password=secret"
	value.QualityReason = "connection-failed"

	svc.handleCollectedValue(t.Context(), value)
	require.Len(t, sink.samples, 1)
	require.Equal(t, schema.QualityBad, sink.samples[0].Quality)
	require.Equal(t, "connection-failed", sink.samples[0].QualityReason)
}

func TestHandleCollectedValueNormalizesUnknownQualityInTypedEnvelope(t *testing.T) {
	sink := &sampleSinkCapture{}
	svc, binding, _ := typedRuntimeFixture(t, sink)
	value := typedRuntimeValue(binding)
	value.Quality = schema.QualityFlag("adapter-private-quality")
	value.QualityReason = "adapter-private-reason"
	value.Error = `dsn=postgres://user:secret@db/password`

	svc.handleCollectedValue(t.Context(), value)
	require.Len(t, sink.samples, 1)
	require.Equal(t, schema.QualityBad, sink.samples[0].Quality)
	require.Equal(t, "quality-unknown", sink.samples[0].QualityReason)
	require.NotContains(t, sink.samples[0].QualityReason, "secret")
}

func TestHandleCollectedValueKeepsSameAddressDevicesDistinctAndFabricatesNothing(t *testing.T) {
	sink := &sampleSinkCapture{}
	svc, bindingA, _ := typedRuntimeFixture(t, sink)

	// Device B exposes the same address 40001 and its tag shares the display
	// name; its connection config holds a credential that must never travel.
	deviceB := &schema.Device{ID: "device-B", Name: "Device B", Protocol: schema.ProtocolModbusTCP, ConnectionConfig: `{"host":"device-b","password":"hunter2"}`}
	pointB := &schema.Point{ID: "point-B", DeviceID: deviceB.ID, Name: "40001", Address: "40001", Function: "03", DataType: schema.DataTypeUint64, Mode: schema.PointModeReadOnly, Enabled: true}
	tagB := &schema.Tag{ID: "tag-B", Key: "tag.b", DisplayName: "Tag A", DataType: schema.DataTypeUint64, Status: schema.TagStatusActive}
	mappingB := &schema.Mapping{ID: "mapping-B", PointID: pointB.ID, TagID: tagB.ID, TransformPipeline: "[]", Status: schema.MappingStatusActive, Enabled: true}
	bindingB, err := mappingBindingForRecords("workspace-A", mappingB, deviceB, pointB, tagB)
	require.NoError(t, err)
	svc.mappingIndex[pointB.ID] = []mappingBinding{bindingB}
	svc.pointMetaIndex[pointB.ID] = pointMeta{DeviceID: deviceB.ID, Address: pointB.Address}

	valueA := typedRuntimeValue(bindingA)
	valueB := typedRuntimeValue(bindingB)
	valueB.ConfigFingerprint = bindingB.ConfigFingerprint
	svc.handleCollectedValue(t.Context(), valueA)
	svc.handleCollectedValue(t.Context(), valueB)
	svc.handleCollectedValue(t.Context(), valueA) // a repeated acquisition

	require.Len(t, sink.samples, 3)
	require.NotEqual(t, sink.samples[0].SampleID, sink.samples[1].SampleID, "persisted IDs keep the two 40001 samples apart")
	require.Equal(t, sink.samples[0].SampleID, sink.samples[2].SampleID, "repeating an acquisition yields the same sample ID")
	for _, sample := range sink.samples {
		require.Equal(t, time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC), sample.ObservedAt, "never the later processing time")
		require.Equal(t, time.Date(2026, 1, 1, 0, 0, 8, 0, time.UTC), sample.ReceivedAt)
		require.Equal(t, string(measurement.GatewayTimeOriginGateway), sample.TimeOrigin)
		require.Empty(t, sample.MeasurementID, "no measurement ID is fabricated for a raw basic sample")
		require.Equal(t, uint64(9007199254740993), sample.Value)
		encoded, err := sample.MarshalJSON()
		require.NoError(t, err)
		require.NotContains(t, string(encoded), "hunter2", "credentials never reach the typed envelope")
		require.NotContains(t, string(encoded), "device-b\"", "connection settings stay out of the envelope")
	}
	require.Equal(t, bindingA.DeviceID, sink.samples[0].DeviceID)
	require.Equal(t, bindingB.DeviceID, sink.samples[1].DeviceID)
	require.NotEqual(t, sink.samples[0].SourceRevision, sink.samples[1].SourceRevision)
}

func scaledFixture(t *testing.T, sink SampleSink) (*Service, mappingBinding) {
	t.Helper()
	svc, binding, _ := typedRuntimeFixture(t, sink)
	binding.TransformPipeline = `[{"type":"scale","params":{"multiplier":10,"offset":0}}]`
	svc.mappingIndex[binding.PointID] = []mappingBinding{binding}
	return svc, binding
}

func TestHandleCollectedValueFailedReadStillReachesTheGroupWhenAPipelineIsConfigured(t *testing.T) {
	sink := &sampleSinkCapture{}
	svc, binding := scaledFixture(t, sink)
	failed := typedRuntimeValue(binding)
	failed.Value = nil
	failed.Quality = schema.QualityBad
	failed.Error = "read failed: password=secret"
	failed.QualityReason = "read-failed"

	svc.handleCollectedValue(t.Context(), failed)
	require.Len(t, sink.samples, 1, "a failed read must not vanish just because the transform cannot run on it")
	require.Equal(t, schema.QualityBad, sink.samples[0].Quality)
	require.Equal(t, "read-failed", sink.samples[0].QualityReason)
	require.Nil(t, sink.samples[0].Value, "no zero or stale value is invented")
}

func TestHandleCollectedValuePipelineFailureOnAGoodReadBecomesABadSample(t *testing.T) {
	sink := &sampleSinkCapture{}
	svc, binding := scaledFixture(t, sink)
	good := typedRuntimeValue(binding)
	good.Value = "not a number" // the scale step cannot convert it

	svc.handleCollectedValue(t.Context(), good)
	require.Len(t, sink.samples, 1)
	require.Equal(t, schema.QualityBad, sink.samples[0].Quality)
	require.Equal(t, "mapping-failed", sink.samples[0].QualityReason)
	require.Nil(t, sink.samples[0].Value)
	require.EqualValues(t, 1, svc.mappingError.Load(), "the mapping failure is still counted")
}

type refusingSink struct{ err error }

func (r refusingSink) AcceptSample(context.Context, measurement.SampleEnvelope) error { return r.err }

func TestHandleCollectedValueTypedRefusalsAreNotCountedAsWriteFailures(t *testing.T) {
	late := &GroupSampleError{Reason: "late-after-close"}
	svc, binding, _ := typedRuntimeFixture(t, refusingSink{err: late})
	svc.handleCollectedValue(t.Context(), typedRuntimeValue(binding))
	require.Zero(t, svc.writeError.Load(), "a late or conflicting sample is normal data, not a storage fault")

	faulty, binding, _ := typedRuntimeFixture(t, refusingSink{err: &GroupSampleError{Reason: "journal-failed"}})
	faulty.handleCollectedValue(t.Context(), typedRuntimeValue(binding))
	require.EqualValues(t, 1, faulty.writeError.Load(), "a journal failure is a real fault")
}

type interestedSink struct {
	wanted   bool
	accepted int
}

func (s *interestedSink) AcceptSample(context.Context, measurement.SampleEnvelope) error {
	s.accepted++
	return nil
}

func (s *interestedSink) WantsSample(string, string, string) bool { return s.wanted }

func TestHandleCollectedValueSkipsTypedPathForTagsNoGroupOwns(t *testing.T) {
	sink := &interestedSink{wanted: false}
	svc, binding, writer := typedRuntimeFixture(t, sink)
	value := typedRuntimeValue(binding)
	value.ConfigFingerprint = "" // would be rejected if the envelope were built

	svc.handleCollectedValue(t.Context(), value)
	require.Zero(t, sink.accepted, "a sample no group wants never reaches the sink")
	require.Zero(t, svc.writeError.Load(), "and building it is not attempted, so it cannot count as a write failure")
	require.Len(t, writer.records, 1, "the regular storage path is unaffected")

	sink.wanted = true
	value = typedRuntimeValue(binding)
	svc.handleCollectedValue(t.Context(), value)
	require.Equal(t, 1, sink.accepted)
}
