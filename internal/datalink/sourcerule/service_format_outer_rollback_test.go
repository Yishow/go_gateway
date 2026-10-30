package sourcerule

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
)

type formatRollbackLinkFault struct {
	Repository
	fail func()
}

func (r *formatRollbackLinkFault) ReplaceLinks(ctx context.Context, id string, links []*schema.SourceRuleLink) error {
	if r.fail != nil {
		f := r.fail
		r.fail = nil
		f()
		return errors.New("injected link persistence failure")
	}
	return r.Repository.ReplaceLinks(ctx, id, links)
}

func TestRuleFormatOuterRollbackDoesNotBypassRevisionConflict(t *testing.T) {
	for _, operation := range []string{"update", "disable"} {
		t.Run(operation, func(t *testing.T) { testFormatOuterRollback(t, operation) })
	}
}

func testFormatOuterRollback(t *testing.T, operation string) {
	ctx := t.Context()
	db := setupSQLRepoDB(t)
	t.Cleanup(func() { _ = db.Close() })
	dr := device.NewSQLRepository(db)
	require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "d"))
	ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
	repo := &formatRollbackLinkFault{Repository: NewSQLRepository(db)}
	syncer := &stubRuntimeSync{}
	svc := NewService(repo, device.NewService(dr, nil), ps, syncer)
	fc := &rollbackFailCloser{}
	svc.SetShareRuntimeReconciler(fc)
	rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 1, DataType: schema.DataTypeInt32, NamingPrefix: "BEFORE"})
	require.NoError(t, err)
	links, err := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	var accepted *schema.Point
	repo.fail = func() {
		var e error
		accepted, e = ps.Update(ctx, links[0].PointID, point.UpdatePointRequest{Description: common.Ptr("operator-accepted")})
		require.NoError(t, e)
		syncer.upserted = nil
	}
	if operation == "update" {
		_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("AFTER")})
	} else {
		_, err = svc.DisableWithRuntimeReconcile(ctx, rule.ID)
	}
	var typed *modbusshare.Error
	require.ErrorAs(t, err, &typed)
	require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
	got, err := ps.GetByID(ctx, accepted.ID)
	require.NoError(t, err)
	require.Equal(t, accepted.RevisionID, got.RevisionID)
	require.Equal(t, "operator-accepted", got.Description)
	require.Empty(t, syncer.upserted)
	require.Positive(t, fc.failClosedCalls)
}

type formatPlanWindowRepository struct {
	Repository
	beforeUpdate func()
}

func (r *formatPlanWindowRepository) Update(ctx context.Context, rule *schema.SourceRule) error {
	if r.beforeUpdate != nil {
		f := r.beforeUpdate
		r.beforeUpdate = nil
		f()
	}
	return r.Repository.Update(ctx, rule)
}
func TestRuleFormatPlanRevisionRejectsInterveningEdit(t *testing.T) {
	for _, operation := range []string{"update", "disable"} {
		t.Run(operation, func(t *testing.T) {
			ctx := t.Context()
			db := setupSQLRepoDB(t)
			t.Cleanup(func() { _ = db.Close() })
			dr := device.NewSQLRepository(db)
			require.NoError(t, seedSQLSourceRuleDevice(ctx, dr, "d"))
			ps := point.NewService(point.NewSQLRepository(db), pollinggroup.NewSQLRepository(db))
			repo := &formatPlanWindowRepository{Repository: NewSQLRepository(db)}
			syncer := &stubRuntimeSync{}
			svc := NewService(repo, device.NewService(dr, nil), ps, syncer)
			fc := &rollbackFailCloser{}
			svc.SetShareRuntimeReconciler(fc)
			rule, err := svc.Create(ctx, CreateRuleRequest{DeviceID: "d", StartAddress: "40001", Count: 2, DataType: schema.DataTypeInt32, NamingPrefix: "BEFORE"})
			require.NoError(t, err)
			links, err := svc.ListLinks(ctx, rule.ID)
			require.NoError(t, err)
			first, err := ps.GetByID(ctx, links[0].PointID)
			require.NoError(t, err)
			var accepted *schema.Point
			repo.beforeUpdate = func() {
				var e error
				accepted, e = ps.Update(ctx, links[1].PointID, point.UpdatePointRequest{Description: common.Ptr("operator-accepted")})
				require.NoError(t, e)
				syncer.upserted = nil
			}
			if operation == "disable" {
				hook := repo.beforeUpdate
				window := &formatPointReadWindow{SQLRepository: point.NewSQLRepository(db), target: links[1].PointID}
				svc.pointSvc = point.NewService(window, pollinggroup.NewSQLRepository(db))
				repo.beforeUpdate = func() { window.afterRead = hook }
			}
			if operation == "update" {
				_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: common.Ptr("AFTER")})
			} else {
				_, err = svc.DisableWithRuntimeReconcile(ctx, rule.ID)
			}
			var typed *modbusshare.Error
			require.ErrorAs(t, err, &typed)
			require.Equal(t, modbusshare.ErrCodeDirtyUnknown, typed.Code)
			got, err := ps.GetByID(ctx, accepted.ID)
			require.NoError(t, err)
			require.Equal(t, accepted.RevisionID, got.RevisionID)
			require.Equal(t, "operator-accepted", got.Description)
			restored, err := ps.GetByID(ctx, first.ID)
			require.NoError(t, err)
			require.Equal(t, first.RevisionID, restored.RevisionID)
			require.NotContains(t, syncer.upserted, accepted.ID)
		})
	}
}

// Inject a successful ordinary edit after the materializer reads its plan.
type formatPointReadWindow struct {
	*point.SQLRepository
	target    string
	afterRead func()
}

func (r *formatPointReadWindow) GetByID(ctx context.Context, id string) (*schema.Point, error) {
	p, err := r.SQLRepository.GetByID(ctx, id)
	if err == nil && id == r.target && r.afterRead != nil {
		f := r.afterRead
		r.afterRead = nil
		f()
	}
	return p, err
}
