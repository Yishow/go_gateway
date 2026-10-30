package runtime

import (
	"testing"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type mappingObservationFixture struct {
	service     *Service
	reader      *sourceObservationReader
	projection  *workspace.RuntimeProjection
	seen        *[]MappingObservation
	writer      *mockWriter
	sink        *sampleSinkCapture
	mappingRepo *mapping.MemoryRepository
}

func newMappingObservationFixture(t *testing.T, observe bool) mappingObservationFixture {
	t.Helper()
	d := &schema.Device{ID: "D", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`}
	p := &schema.Point{ID: "P", DeviceID: d.ID, Address: "40001", Function: "03", DataType: schema.DataTypeUint16, Mode: schema.PointModeReadOnly, Enabled: true}
	tg := &schema.Tag{ID: "T", Key: "tag.one", DataType: schema.DataTypeInt16, Status: schema.TagStatusActive}
	m := &schema.Mapping{ID: "M", PointID: p.ID, TagID: tg.ID, TransformPipeline: `[{"type":"cast","params":{"target_type":"int16"}}]`, Enabled: true, Status: schema.MappingStatusActive}
	projection := &workspace.RuntimeProjection{WorkspaceID: "W", Version: "projection-input-not-workspace-revision", DeviceIDs: []string{d.ID}, Devices: []*schema.Device{d}, Points: []*schema.Point{p}, Tags: []*schema.Tag{tg}, Mappings: []*schema.Mapping{m}}
	reader := &sourceObservationReader{record: &workspace.Record{ID: "W", DatabaseSetupRevision: "WR1", OrderedDeviceIDs: []string{d.ID}}, projection: projection}
	dr, pr, gr, tr, mr := device.NewMemoryRepository(), point.NewMemoryRepository(), pollinggroup.NewMemoryRepository(), tag.NewMemoryRepository(), mapping.NewMemoryRepository()
	require.NoError(t, dr.Create(t.Context(), d))
	require.NoError(t, pr.Create(t.Context(), p))
	require.NoError(t, tr.Create(t.Context(), tg))
	require.NoError(t, mr.Create(t.Context(), m))
	seen := []MappingObservation{}
	var observer MappingObserver
	if observe {
		observer = func(v MappingObservation) { seen = append(seen, v) }
	}
	writer, sink := &mockWriter{}, &sampleSinkCapture{}
	s, err := NewService(Config{UpdatePointState: false}, Dependencies{Scheduler: collector.NewScheduler(collector.DefaultSchedulerConfig(), nil), Writer: writer, DeviceService: device.NewService(dr, nil), PointService: point.NewService(pr, gr), MappingService: mapping.NewService(mr), TagService: tag.NewService(tr), PollingGroupService: pollinggroup.NewService(gr), WorkspaceProjection: reader, SampleSink: sink, MappingObserver: observer})
	require.NoError(t, err)
	return mappingObservationFixture{s, reader, projection, &seen, writer, sink, mr}
}

func (f mappingObservationFixture) install(t *testing.T) {
	t.Helper()
	require.NoError(t, f.service.ApplyWorkspaceProjection(t.Context(), f.projection))
}
func (f mappingObservationFixture) value(t *testing.T, value any) collector.CollectedValue {
	t.Helper()
	bindings := f.service.mappingIndex["P"]
	require.NotEmpty(t, bindings)
	cv := typedRuntimeValue(bindings[0])
	cv.Value = value
	return cv
}
func (f mappingObservationFixture) execute(t *testing.T, value any) MappingObservation {
	t.Helper()
	before := len(*f.seen)
	f.service.handleCollectedValue(t.Context(), f.value(t, value))
	require.Len(t, *f.seen, before+1, "one closed actual transform result")
	return (*f.seen)[before]
}
