package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

type shareReadySourceRuleRepository struct {
	*sourcerule.MemoryRepository
	pointSvc   *point.Service
	tagSvc     *tag.Service
	mappingSvc *mapping.Service
}

func (r *shareReadySourceRuleRepository) CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) error {
	for _, link := range links {
		pointRecord, err := r.pointSvc.GetByID(ctx, link.PointID)
		if err != nil {
			return err
		}
		tagRecord, err := r.tagSvc.Create(ctx, tag.CreateTagRequest{
			Key:         pointRecord.Name,
			DisplayName: pointRecord.Name,
			DataType:    pointRecord.DataType,
		})
		if err != nil {
			return err
		}
		enabled := true
		status := schema.MappingStatusActive
		mappingRecord, err := r.mappingSvc.Create(ctx, mapping.CreateMappingRequest{
			PointID:           pointRecord.ID,
			TagID:             tagRecord.ID,
			Enabled:           &enabled,
			Status:            &status,
			TransformPipeline: nil,
		})
		if err != nil {
			return err
		}
		link.TagID = &tagRecord.ID
		link.MappingID = &mappingRecord.ID
	}
	return r.MemoryRepository.CreateLinks(ctx, links)
}

type freshShareHandlerFixture struct {
	handler     *StudioV2WorkspaceSourceRulesHandler
	share       *modbusshare.Service
	workspaceID string
}

func newFreshShareHandlerFixture(t *testing.T) freshShareHandlerFixture {
	return newFreshShareHandlerFixtureWithRuntime(t, true, nil)
}

func newFreshShareHandlerFixtureWithRuntime(t *testing.T, running bool, runtimeSync sourcerule.RuntimeSyncer) freshShareHandlerFixture {
	t.Helper()
	gin.SetMode(gin.TestMode)
	ctx := t.Context()

	deviceRepo := device.NewMemoryRepository()
	if running {
		seedRunningWorkspaceRuleDevice(t, deviceRepo, "dev-share-fresh")
	} else {
		seedWorkspaceRuleDevice(t, deviceRepo, "dev-share-fresh")
	}
	deviceSvc := device.NewService(deviceRepo, nil)
	pointSvc := point.NewService(point.NewMemoryRepository(), nil)
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	mappingSvc := mapping.NewServiceWithTagResolver(mapping.NewMemoryRepository(), tagSvc.GetByID)
	repo := &shareReadySourceRuleRepository{
		MemoryRepository: sourcerule.NewMemoryRepository(),
		pointSvc:         pointSvc,
		tagSvc:           tagSvc,
		mappingSvc:       mappingSvc,
	}
	ruleSvc := sourcerule.NewService(repo, deviceSvc, pointSvc, runtimeSync)
	ruleSvc.SetTagMappingServices(tagSvc, mappingSvc)
	workspaceSvc := workspace.NewService(workspace.NewMemoryRepository()).WithReadinessServices(deviceSvc, ruleSvc, nil, nil)
	_, err := workspaceSvc.AttachDevice(ctx, "dev-share-fresh")
	require.NoError(t, err)
	record, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	revisionStore := modbusshare.NewSQLWorkspaceRevisionStore(db)
	workspaceRevision, _, err := revisionStore.GetRevision(ctx, record.ID)
	require.NoError(t, err)

	share := modbusshare.NewService(tagSvc, 4096)
	share.SetWorkspaceRevisionStore(revisionStore)
	share.SetHydrationState(modbusshare.HydrationState{
		State:             modbusshare.HydrationStateReady,
		WorkspaceID:       record.ID,
		WorkspaceRevision: workspaceRevision,
		SettingsRevision:  "settings-fresh-share",
		Readiness:         true,
		ReadinessToken:    "hydration-fresh-share",
	})
	require.NoError(t, share.ApplySettings(ctx, modbusshare.Settings{
		Enabled:           true,
		BindAddress:       "127.0.0.1",
		Port:              reserveFreshSharePort(t),
		SlaveID:           1,
		CapacityRegisters: 64,
		SettingsRevision:  "settings-fresh-share",
	}))
	t.Cleanup(func() { require.NoError(t, share.StopListener()) })

	ownership := sourcerule.NewShareDesiredMappingOwnershipChecker(
		workspaceSvc,
		ruleSvc,
		tagSvc,
		share.Settings,
		mappingSvc,
	)
	share.SetDesiredMappingOwnershipChecker(ownership)
	reconciler := modbusshare.NewReconciler(share, revisionStore).WithOwnershipValidator(ownership)
	ruleSvc.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(ctx context.Context, req sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		currentWorkspace, err := workspaceSvc.GetOrCreate(ctx)
		if err != nil {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusFailed, Scope: req.Scope}, err
		}
		currentRevision, _, err := revisionStore.GetRevision(ctx, currentWorkspace.ID)
		if err != nil {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusFailed, Scope: req.Scope}, err
		}
		hydration, err := share.CheckHydration(ctx)
		if err != nil {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusFailed, Scope: req.Scope}, err
		}
		settings := share.Settings()
		desired, err := ruleSvc.BuildDesiredShareMappingsForDevices(ctx, currentWorkspace.ID, settings, currentWorkspace.OrderedDeviceIDs)
		if err != nil {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusFailed, Scope: req.Scope}, err
		}
		outcome, err := reconciler.Reconcile(ctx, modbusshare.ReconcileRequest{
			WorkspaceID:               currentWorkspace.ID,
			ExpectedWorkspaceRevision: currentRevision,
			ExpectedSettingsRevision:  settings.SettingsRevision,
			ReadinessToken:            hydration.ReadinessToken,
			DesiredMappings:           desired,
		})
		if err != nil {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusFailed, Scope: req.Scope}, err
		}
		if outcome.Outcome != "applied" && outcome.Outcome != "aligned" {
			return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusFailed, Scope: req.Scope}, fmt.Errorf("share projection outcome=%s", outcome.Outcome)
		}
		return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusAligned, Scope: req.Scope, Message: "projection aligned"}, nil
	}))

	return freshShareHandlerFixture{
		handler:     NewStudioV2WorkspaceSourceRulesHandler(workspaceSvc, deviceSvc, ruleSvc),
		share:       share,
		workspaceID: record.ID,
	}
}

type freshShareNotRunningRuntime struct{}

func (freshShareNotRunningRuntime) UpsertPoint(*schema.Point) {}

func (freshShareNotRunningRuntime) RemovePoint(string) {}

func (freshShareNotRunningRuntime) ReconcileSourceRule(_ context.Context, req sourcerule.RuntimeReconcileRequest) sourcerule.RuntimeReconcileOutcome {
	return sourcerule.RuntimeReconcileOutcome{
		Status:  sourcerule.RuntimeReconcileStatusNotRunning,
		Scope:   req.Scope,
		Message: "runtime is not running",
	}
}

func reserveFreshSharePort(t *testing.T) int {
	t.Helper()
	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func TestStudioV2WorkspaceSourceRulesHandler_CreateAppliesFreshShareProjection(t *testing.T) {
	fixture := newFreshShareHandlerFixture(t)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/source-rules", bytes.NewBufferString(`{
		"id":"rule-fresh-share",
		"device_id":"dev-share-fresh",
		"start_address":"40001",
		"count":2,
		"data_type":"int16",
		"naming_prefix":"FRESH_",
		"enabled":true,
		"share_enabled":true,
		"share_start_register":40001,
		"share_stride":1
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	fixture.handler.Create(c)

	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())
	mappings, err := fixture.share.ListMappingsForWorkspace(t.Context(), fixture.workspaceID)
	require.NoError(t, err)
	require.Len(t, mappings, 2)
	require.Contains(t, resp.Body.String(), `"runtime_apply_status":"aligned"`)
}

func TestStudioV2WorkspaceSourceRulesHandler_CreateInactiveDeviceStillUpdatesShareProjection(t *testing.T) {
	fixture := newFreshShareHandlerFixtureWithRuntime(t, false, freshShareNotRunningRuntime{})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/api/v1/datalink/studio-v2/workspace/source-rules", bytes.NewBufferString(`{
		"id":"rule-inactive-share",
		"device_id":"dev-share-fresh",
		"start_address":"40001",
		"count":2,
		"data_type":"int16",
		"naming_prefix":"INACTIVE_",
		"enabled":true,
		"share_enabled":true,
		"share_start_register":40001,
		"share_stride":1
	}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	fixture.handler.Create(c)

	require.Equal(t, http.StatusCreated, resp.Code, resp.Body.String())
	mappings, err := fixture.share.ListMappingsForWorkspace(t.Context(), fixture.workspaceID)
	require.NoError(t, err)
	require.Len(t, mappings, 2)
	require.Contains(t, resp.Body.String(), `"runtime_apply_status":"aligned"`)
}
