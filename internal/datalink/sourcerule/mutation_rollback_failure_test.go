package sourcerule

import (
	"context"
	"errors"
	"testing"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type rollbackSourceRuleRepository struct {
	*MemoryRepository
	createLinksFailures int
}

func (r *rollbackSourceRuleRepository) CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) error {
	if r.createLinksFailures > 0 {
		r.createLinksFailures--
		return errors.New("injected source-rule link write failure")
	}
	return r.MemoryRepository.CreateLinks(ctx, links)
}

type rollbackPointRepository struct {
	point.Repository
	getFailures    map[int]error
	updateFailures map[int]error
	getCalls       int
	updateCalls    int
}

func (r *rollbackPointRepository) GetByID(ctx context.Context, id string) (*schema.Point, error) {
	r.getCalls++
	if err := r.getFailures[r.getCalls]; err != nil {
		return nil, err
	}
	return r.Repository.GetByID(ctx, id)
}

func (r *rollbackPointRepository) Update(ctx context.Context, pointRecord *schema.Point) error {
	r.updateCalls++
	if err := r.updateFailures[r.updateCalls]; err != nil {
		return err
	}
	return r.Repository.Update(ctx, pointRecord)
}

type rollbackTagRepository struct {
	tag.Repository
	deleteFailures int
}

func (r *rollbackTagRepository) Delete(ctx context.Context, id string) error {
	if r.deleteFailures > 0 {
		r.deleteFailures--
		return errors.New("injected tag rollback failure")
	}
	return r.Repository.Delete(ctx, id)
}

type rollbackMappingRepository struct {
	mapping.Repository
	createCalls  int
	failCreateOn int
}

func (r *rollbackMappingRepository) Create(ctx context.Context, record *schema.Mapping) error {
	r.createCalls++
	if r.createCalls == r.failCreateOn {
		return errors.New("injected mapping apply failure")
	}
	return r.Repository.Create(ctx, record)
}

type rollbackPollingGroupRepository struct {
	point.PollingGroupRepository
	listErr error
}

func (r *rollbackPollingGroupRepository) List(context.Context) ([]*schema.PollingGroup, error) {
	return nil, r.listErr
}

type rollbackFailCloser struct {
	failClosedCalls int
}

func (r *rollbackFailCloser) ReconcileSourceRuleShare(context.Context, RuntimeReconcileRequest) (RuntimeReconcileOutcome, error) {
	return RuntimeReconcileOutcome{Status: RuntimeReconcileStatusAligned}, nil
}

func (r *rollbackFailCloser) FailClosed(context.Context) {
	r.failClosedCalls++
}

func newRollbackUpdateService(t *testing.T, pointRepo point.Repository, repo Repository, runtimeSync RuntimeSyncer) (*Service, *schema.SourceRule) {
	t.Helper()
	ctx := context.Background()
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(pointRepo, nil)
	svc := NewService(repo, deviceSvc, pointSvc, runtimeSync)

	dev, err := seedActiveDevice(ctx, deviceRepo, "rollback-device")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "rollback-rule",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		DataFormat:   "ABCD",
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)
	return svc, rule
}

func TestService_Update_PointRollbackFailureReturnsDirtyUnknown(t *testing.T) {
	ctx := context.Background()
	pointRepo := &rollbackPointRepository{
		Repository:     point.NewMemoryRepository(),
		getFailures:    map[int]error{},
		updateFailures: map[int]error{2: errors.New("injected point rollback failure")},
	}
	baseRepo := NewMemoryRepository()
	repo := &rollbackSourceRuleRepository{MemoryRepository: baseRepo}
	runtimeSync := &stubRuntimeSync{}
	svc, rule := newRollbackUpdateService(t, pointRepo, repo, runtimeSync)
	share := &rollbackFailCloser{}
	svc.SetShareRuntimeReconciler(share)
	repo.createLinksFailures = 1
	initialUpserts := len(runtimeSync.upserted)

	prefix := "NEXT"
	_, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dirty_unknown")
	assert.Equal(t, 1, share.failClosedCalls)
	assert.Equal(t, initialUpserts+1, len(runtimeSync.upserted))
}

func TestService_Update_CleanRollbackPreservesOriginalError(t *testing.T) {
	ctx := context.Background()
	pointRepo := &rollbackPointRepository{
		Repository:     point.NewMemoryRepository(),
		getFailures:    map[int]error{},
		updateFailures: map[int]error{},
	}
	baseRepo := NewMemoryRepository()
	repo := &rollbackSourceRuleRepository{MemoryRepository: baseRepo}
	svc, rule := newRollbackUpdateService(t, pointRepo, repo, nil)
	share := &rollbackFailCloser{}
	svc.SetShareRuntimeReconciler(share)
	repo.createLinksFailures = 1

	dataFormat := "CDAB"
	_, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{DataFormat: &dataFormat})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "重建來源規則連結失敗")
	assert.NotContains(t, err.Error(), "dirty_unknown")
	assert.Zero(t, share.failClosedCalls)
	links, linkErr := svc.ListLinks(ctx, rule.ID)
	require.NoError(t, linkErr)
	require.Len(t, links, 1)
	restoredPoint, pointErr := svc.pointSvc.GetByID(ctx, links[0].PointID)
	require.NoError(t, pointErr)
	assert.Equal(t, "ABCD", restoredPoint.DataFormat)
}

func TestService_Update_RepositoryRollbackFailureReturnsDirtyUnknown(t *testing.T) {
	ctx := context.Background()
	pointRepo := &rollbackPointRepository{
		Repository:     point.NewMemoryRepository(),
		getFailures:    map[int]error{},
		updateFailures: map[int]error{},
	}
	baseRepo := NewMemoryRepository()
	repo := &rollbackSourceRuleRepository{MemoryRepository: baseRepo}
	svc, rule := newRollbackUpdateService(t, pointRepo, repo, nil)
	share := &rollbackFailCloser{}
	svc.SetShareRuntimeReconciler(share)
	repo.createLinksFailures = 2

	prefix := "NEXT"
	_, err := svc.Update(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dirty_unknown")
	assert.Equal(t, 1, share.failClosedCalls)
}

func TestService_Update_DisabledNewPointUpdateFailureDoesNotPanic(t *testing.T) {
	ctx := context.Background()
	pointRepo := &rollbackPointRepository{
		Repository:     point.NewMemoryRepository(),
		getFailures:    map[int]error{},
		updateFailures: map[int]error{3: errors.New("injected new point state failure")},
	}
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	svc := NewService(repo, deviceSvc, point.NewService(pointRepo, nil), nil)
	dev, err := seedActiveDevice(ctx, deviceRepo, "disabled-update-device")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "disabled-update-rule",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      false,
	})
	require.NoError(t, err)
	count := 2

	var updateErr error
	assert.NotPanics(t, func() {
		_, updateErr = svc.Update(ctx, rule.ID, UpdateRuleRequest{Count: &count})
	})
	require.Error(t, updateErr)
	assert.Contains(t, updateErr.Error(), "更新新增衍生點位狀態失敗")
}

func TestService_Update_PollingGroupResolutionFailureRestoresRule(t *testing.T) {
	ctx := context.Background()
	pointRepo := &rollbackPointRepository{
		Repository:     point.NewMemoryRepository(),
		getFailures:    map[int]error{},
		updateFailures: map[int]error{},
	}
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, nil)
	repo := NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)
	svc := NewService(repo, deviceSvc, pointSvc, nil)
	dev, err := seedActiveDevice(ctx, deviceRepo, "polling-rollback-device")
	require.NoError(t, err)
	rule, err := svc.Create(ctx, CreateRuleRequest{
		ID:           "polling-rollback-rule",
		DeviceID:     dev.ID,
		StartAddress: "40001",
		Count:        1,
		DataType:     schema.DataTypeInt16,
		NamingPrefix: "SRC",
		Enabled:      true,
	})
	require.NoError(t, err)

	groupRepo := &rollbackPollingGroupRepository{
		PollingGroupRepository: pollinggroup.NewMemoryRepository(),
		listErr:                errors.New("injected polling group read failure"),
	}
	svc.pointSvc = point.NewService(pointRepo, groupRepo)
	share := &rollbackFailCloser{}
	svc.SetShareRuntimeReconciler(share)
	prefix := "NEXT"
	_, err = svc.Update(ctx, rule.ID, UpdateRuleRequest{NamingPrefix: &prefix})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "確保預設輪詢群組失敗")
	assert.Zero(t, share.failClosedCalls)
	stored, getErr := svc.GetByID(ctx, rule.ID)
	require.NoError(t, getErr)
	assert.Equal(t, "SRC", stored.NamingPrefix)
}
