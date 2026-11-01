package runtime

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type independentInstallReader struct {
	projection *workspace.RuntimeProjection
	contexts   []context.Context
}

func (r *independentInstallReader) RuntimeProjection(ctx context.Context) (*workspace.RuntimeProjection, error) {
	r.contexts = append(r.contexts, ctx)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return r.projection, nil
}

type independentInstallFixture struct {
	service *Service
	reader  *independentInstallReader
	record  *schema.Point
	device  *schema.Device
}

func newIndependentInstallFixture(t *testing.T) independentInstallFixture {
	t.Helper()
	groupID := "group"
	record := &schema.Point{ID: "P", DeviceID: "D", Address: "40001", Function: "03", DataType: schema.DataTypeUint32, Enabled: true, RevisionID: "persisted-point", FormatMode: schema.PointFormatOverride, FormatOverride: "CDAB", FormatVersion: 1, PollingGroupID: &groupID}
	device := &schema.Device{ID: "D", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"data_format":"ABCD"}`}
	projection := &workspace.RuntimeProjection{WorkspaceID: "W", DeviceIDs: []string{"D"}, Devices: []*schema.Device{device}, Points: []*schema.Point{record}}
	reader := &independentInstallReader{projection: projection, contexts: make([]context.Context, 0, 2)}
	scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
	t.Cleanup(func() { require.NoError(t, scheduler.StopContext(context.WithoutCancel(t.Context()))) })
	service := &Service{scheduler: scheduler, workspace: reader, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{}, projections: map[string]runtimeProjectionDeviceState{}}
	return independentInstallFixture{service: service, reader: reader, record: record, device: device}
}

func assertIndependentProfileInstalled(t *testing.T, f independentInstallFixture) {
	t.Helper()
	require.NotEmpty(t, f.reader.contexts)
	require.NoError(t, f.reader.contexts[0].Err())
	require.Nil(t, f.reader.contexts[0].Value(readProfilesKey{}), "installation must not borrow a prepared/foreign cached profile")
	lease := f.service.CapturePointLease(f.record.ID, f.record.RevisionID)
	require.NotZero(t, lease)
	profile, ok := f.service.CapturePointReadProfile(f.record.ID, lease)
	require.True(t, ok)
	require.Equal(t, "CDAB", profile.Resolution().EffectiveFormat)
}

func TestSnapshotInstallKeepsIndependentFreshProfileBoundary(t *testing.T) {
	for _, mode := range []string{"canceled", "foreign-cache", "both"} {
		t.Run(mode, func(t *testing.T) {
			f := newIndependentInstallFixture(t)
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if mode != "canceled" {
				ctx = context.WithValue(ctx, readProfilesKey{}, map[string]connector.ReadProfile{"P": {}})
			}
			if mode != "foreign-cache" {
				cancel()
			}
			// The original later mapping stage rejects the canceled/foreign request,
			// so this unit test never starts scheduler tickers or network reads.
			f.service.snapshot = Snapshot{WorkspaceID: "W", Devices: []*schema.Device{f.device}, Points: []*schema.Point{f.record}, Mappings: []*schema.Mapping{{ID: "M", PointID: "P", TagID: "T", Enabled: true}}}
			require.ErrorIs(t, f.service.bootstrapFromSnapshot(ctx), connector.ErrReadProfile)
			assertIndependentProfileInstalled(t, f)
			if mode == "canceled" {
				require.Len(t, f.reader.contexts, 2)
				require.ErrorIs(t, f.reader.contexts[1].Err(), context.Canceled, "later mapping retains caller cancellation")
			} else {
				require.Len(t, f.reader.contexts, 1, "later mapping rejects the foreign cached profile without another projection read")
			}
		})
	}
}

type cancelAfterPointListRepository struct {
	*point.MemoryRepository
	cancel context.CancelFunc
}

func (r *cancelAfterPointListRepository) List(ctx context.Context, filter point.ListFilter) ([]*schema.Point, error) {
	records, err := r.MemoryRepository.List(ctx, filter)
	r.cancel()
	return records, err
}

func TestDeviceUpsertInstallKeepsIndependentFreshProfileBoundary(t *testing.T) {
	f := newIndependentInstallFixture(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctx = context.WithValue(ctx, readProfilesKey{}, map[string]connector.ReadProfile{"P": {}})
	repo := &cancelAfterPointListRepository{MemoryRepository: point.NewMemoryRepository(), cancel: cancel}
	require.NoError(t, repo.Create(t.Context(), f.record))
	f.service.pointSvc = point.NewService(repo, nil)
	require.NoError(t, f.service.UpsertDevice(ctx, f.device), "preserve the original return; later alignment read may fail")
	require.ErrorIs(t, ctx.Err(), context.Canceled)
	assertIndependentProfileInstalled(t, f)
	require.Len(t, f.reader.contexts, 2)
	require.ErrorIs(t, f.reader.contexts[1].Err(), context.Canceled, "later alignment still uses the original canceled caller")
}
