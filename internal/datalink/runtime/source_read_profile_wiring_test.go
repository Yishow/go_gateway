package runtime

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
	"testing"
)

type profileBoundaryRepository struct {
	*sourcerule.MemoryRepository
	fault   string
	inspect func()
}

func (r *profileBoundaryRepository) ReplaceLinks(ctx context.Context, id string, links []*schema.SourceRuleLink) error {
	if r.fault == "links" {
		r.fault = ""
		r.inspect()
		return errors.New("injected links failure")
	}
	return r.MemoryRepository.ReplaceLinks(ctx, id, links)
}
func (r *profileBoundaryRepository) ReplaceCandidateSnapshotsAtRevision(ctx context.Context, snapshots []*schema.SourceRuleCandidateSnapshot, rev string) error {
	if r.fault == "candidate" {
		r.fault = ""
		r.inspect()
		return errors.New("injected candidate failure")
	}
	return r.MemoryRepository.ReplaceCandidateSnapshotsAtRevision(ctx, snapshots, rev)
}
func TestReadProfileSourceCreateAndUpdateInstallOnlyAfterLinksAndCandidate(t *testing.T) {
	for _, fault := range []string{"", "links", "candidate"} {
		t.Run("boundary-"+fault, func(t *testing.T) {
			ctx := t.Context()
			dr := device.NewMemoryRepository()
			ready := true
			d := &schema.Device{ID: "D", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusActive, ConnectionConfig: `{"host":"127.0.0.1","data_format":"ABCD"}`, LastTestSuccess: &ready}
			require.NoError(t, dr.Create(ctx, d))
			ps := point.NewService(point.NewMemoryRepository(), nil)
			scheduler := collector.NewScheduler(collector.DefaultSchedulerConfig(), nil)
			live := &Service{scheduler: scheduler, pointSvc: ps, mappingIndex: map[string][]mappingBinding{}, pointMetaIndex: map[string]pointMeta{}}
			repo := &profileBoundaryRepository{MemoryRepository: sourcerule.NewMemoryRepository()}
			rules := sourcerule.NewService(repo, device.NewService(dr, nil), ps, live)
			live.workspace = compensationWorkspaceFunc(func(ctx context.Context) (*workspace.RuntimeProjection, error) {
				return persistedRuleProjection(ctx, rules, ps, d)
			})
			// Real device entry first; rule entry then creates persisted points and links.
			require.NoError(t, live.UpsertDevice(ctx, d))
			live.running.Store(true)
			rule, outcome, err := rules.CreateWithRuntimeReconcile(ctx, sourcerule.CreateRuleRequest{ID: "R", DeviceID: d.ID, StartAddress: "40001", Count: 1, DataType: schema.DataTypeUint32, DataFormat: "CDAB", NamingPrefix: "P", Enabled: true})
			require.NoError(t, err)
			require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
			links, err := rules.ListLinks(ctx, rule.ID)
			require.NoError(t, err)
			require.Len(t, links, 1)
			originalID := links[0].PointID
			before, err := ps.GetByID(ctx, originalID)
			require.NoError(t, err)
			lease := live.CapturePointLease(before.ID, before.RevisionID)
			require.NotZero(t, lease)
			originalProfile, ok := live.CapturePointReadProfile(before.ID, lease)
			require.True(t, ok)
			require.Equal(t, "CDAB", originalProfile.Resolution().EffectiveFormat)
			repo.inspect = func() {
				points, err := ps.ListByDevice(ctx, d.ID)
				require.NoError(t, err)
				for _, p := range points {
					if p.ID != originalID {
						require.Zero(t, live.CapturePointLease(p.ID, p.RevisionID), "new point must not install before links AND candidates commit")
					}
				}
			}
			repo.fault = fault
			count := 2
			_, outcome, err = rules.UpdateWithRuntimeReconcile(ctx, rule.ID, sourcerule.UpdateRuleRequest{Count: &count})
			if fault != "" {
				require.Error(t, err)
				points, err := ps.ListByDevice(ctx, d.ID)
				require.NoError(t, err)
				require.Len(t, points, 1)
				restored := live.CapturePointLease(before.ID, before.RevisionID)
				require.NotZero(t, restored)
				profile, ok := live.CapturePointReadProfile(before.ID, restored)
				require.True(t, ok)
				require.Equal(t, originalProfile, profile)
			} else {
				require.NoError(t, err)
				require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
				points, err := ps.ListByDevice(ctx, d.ID)
				require.NoError(t, err)
				require.Len(t, points, 2)
				for _, p := range points {
					lease := live.CapturePointLease(p.ID, p.RevisionID)
					require.NotZero(t, lease)
					profile, ok := live.CapturePointReadProfile(p.ID, lease)
					require.True(t, ok)
					require.Equal(t, "CDAB", profile.Resolution().EffectiveFormat)
				}
			}
		})
	}
}
