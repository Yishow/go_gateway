package sourcerule

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type rollbackReadFailureRepository struct {
	Repository
	getByIDCalls     int
	listLinksCalls   int
	getByIDErr       error
	listLinksErr     error
	updateErr        error
	createCalls      int
	deleteCalls      int
	createLinksCalls int
	deleteLinksCalls int
}

func (r *rollbackReadFailureRepository) GetByID(ctx context.Context, id string) (*schema.SourceRule, error) {
	r.getByIDCalls++
	if r.getByIDCalls == 3 && r.getByIDErr != nil {
		return nil, r.getByIDErr
	}
	return r.Repository.GetByID(ctx, id)
}

func (r *rollbackReadFailureRepository) ListLinks(ctx context.Context, ruleID string) ([]*schema.SourceRuleLink, error) {
	r.listLinksCalls++
	if r.listLinksCalls == 3 && r.listLinksErr != nil {
		return nil, r.listLinksErr
	}
	return r.Repository.ListLinks(ctx, ruleID)
}

func (r *rollbackReadFailureRepository) Update(context.Context, *schema.SourceRule) error {
	return r.updateErr
}

func (r *rollbackReadFailureRepository) Create(ctx context.Context, rule *schema.SourceRule) error {
	r.createCalls++
	return r.Repository.Create(ctx, rule)
}

func (r *rollbackReadFailureRepository) Delete(ctx context.Context, id string) error {
	r.deleteCalls++
	return r.Repository.Delete(ctx, id)
}

func (r *rollbackReadFailureRepository) CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) error {
	r.createLinksCalls++
	return r.Repository.CreateLinks(ctx, links)
}

func (r *rollbackReadFailureRepository) DeleteLinks(ctx context.Context, ruleID string) error {
	r.deleteLinksCalls++
	return r.Repository.DeleteLinks(ctx, ruleID)
}

func TestService_UpdateWithRuntimeReconcile_GetByIDTransientReadFailureDoesNotMutateDuringRollback(t *testing.T) {
	testUpdateWithRuntimeReconcileReadFailure(t, func(repo *rollbackReadFailureRepository) {
		repo.getByIDErr = errors.New("database is temporarily unavailable")
	})
}

func TestService_UpdateWithRuntimeReconcile_ListLinksTransientReadFailureDoesNotMutateDuringRollback(t *testing.T) {
	testUpdateWithRuntimeReconcileReadFailure(t, func(repo *rollbackReadFailureRepository) {
		repo.listLinksErr = errors.New("source rule links read timed out")
	})
}

func TestService_RestoreRuleSnapshot_SourceRuleNotFoundAllowsCreate(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo, nil, nil, nil)
	snapshot := &sourceRuleSnapshot{rule: &schema.SourceRule{ID: "missing-rule", RevisionID: "revision-1"}}

	require.NoError(t, svc.restoreRuleSnapshot(context.Background(), snapshot))
	stored, err := repo.GetByID(context.Background(), snapshot.rule.ID)
	require.NoError(t, err)
	assert.Equal(t, snapshot.rule.ID, stored.ID)
}

func testUpdateWithRuntimeReconcileReadFailure(t *testing.T, configure func(*rollbackReadFailureRepository)) {
	t.Helper()
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	pointRepo := point.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	baseRepo := NewMemoryRepository()
	repo := &rollbackReadFailureRepository{Repository: baseRepo, updateErr: errors.New("save failed after validation")}
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	dev, err := seedActiveDevice(ctx, deviceRepo, "rollback-read-device")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID: "rollback-read-rule", DeviceID: dev.ID, StartAddress: "40001", Count: 1,
		DataType: schema.DataTypeInt16, NamingPrefix: "SRC", Enabled: true,
	})
	require.NoError(t, err)
	links, err := baseRepo.ListLinks(ctx, rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)

	repo.getByIDCalls = 0
	repo.listLinksCalls = 0
	repo.createCalls = 0
	repo.deleteCalls = 0
	repo.createLinksCalls = 0
	repo.deleteLinksCalls = 0
	configure(repo)
	_, _, err = svc.UpdateWithRuntimeReconcile(ctx, rule.ID, UpdateRuleRequest{Count: intPtr(2)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "dirty_unknown")
	assert.Zero(t, repo.createCalls)
	assert.Zero(t, repo.deleteCalls)
	assert.Zero(t, repo.createLinksCalls)
	assert.Zero(t, repo.deleteLinksCalls)
	storedPoint, pointErr := pointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, pointErr)
	assert.Equal(t, links[0].PointID, storedPoint.ID)
}
