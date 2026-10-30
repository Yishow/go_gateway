package sourcerule

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"testing"
)

type formatShareFailure struct {
	hook   func()
	closed int
}

func (r *formatShareFailure) ReconcileSourceRuleShare(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
	if r.hook != nil {
		f := r.hook
		r.hook = nil
		f()
	}
	return RuntimeReconcileOutcome{Status: RuntimeReconcileStatusFailed}, errors.New("share unavailable")
}
func (r *formatShareFailure) FailClosed(context.Context) { r.closed++ }

func TestRuleFormatCompensationCannotOverwriteLaterPointRevision(t *testing.T) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	t.Cleanup(func() { _ = db.Close() })
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "d"))
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	syncer := &stubRuntimeSync{}
	svc := NewService(NewSQLRepository(db), device.NewService(dr, nil), ps, syncer)
	rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "OLD"})
	require.NoError(t, err)
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	var accepted *schema.Point
	share := &formatShareFailure{hook: func() {
		var e error
		accepted, e = ps.Update(ctx, links[0].PointID, point.UpdatePointRequest{Description: common.Ptr("accepted-newer")})
		require.NoError(t, e)
		syncer.upserted = nil
	}}
	svc.SetShareRuntimeReconciler(share)
	_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("NEW")})
	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
	got, err := ps.GetByID(ctx, accepted.ID)
	require.NoError(t, err)
	require.Equal(t, accepted.RevisionID, got.RevisionID)
	require.Equal(t, accepted.Description, got.Description)
	require.Empty(t, syncer.upserted)
	require.Positive(t, share.closed)
}

func TestRuleFormatCompensationDoesNotReplaceForeignRuleRevision(t *testing.T) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	t.Cleanup(func() { _ = db.Close() })
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "d"))
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	repo := NewSQLRepository(db)
	svc := NewService(repo, device.NewService(dr, nil), ps, nil)
	rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "OLD"})
	require.NoError(t, err)
	share := &formatShareFailure{hook: func() {
		current, e := repo.GetByID(ctx, rule.ID)
		require.NoError(t, e)
		current.RevisionID = "foreign-rule-revision"
		current.NamingPrefix = "EXTERNAL"
		require.NoError(t, repo.Update(ctx, current))
		require.NoError(t, repo.ReplaceCandidateSnapshots(ctx, []*schema.SourceRuleCandidateSnapshot{{SourceRuleID: rule.ID, RevisionID: current.RevisionID, CandidateType: schema.SourceRuleCandidateTypeTags, Status: schema.SourceRuleCandidateStatusReady, Payload: `{"candidates":[]}`}}))
	}}
	svc.SetShareRuntimeReconciler(share)
	_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("NEW")})
	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
	current, err := repo.GetByID(ctx, rule.ID)
	require.NoError(t, err)
	require.Equal(t, "foreign-rule-revision", current.RevisionID)
	candidates, err := repo.ListCandidateSnapshots(ctx, rule.ID, current.RevisionID)
	require.NoError(t, err)
	require.Len(t, candidates, 1)
}
