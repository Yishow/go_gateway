package main

// @title Go Gateway API
// @version 1.0
// @description 工業數據採集閘道系統 API 文檔
// @host localhost:8080
// @BasePath /api

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"go-gateway/internal/api"
	"go-gateway/internal/apphost"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters" // 導入所有適配器以觸發 init() 註冊協議
	"go-gateway/internal/datalink/history"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/desktop"
	"go-gateway/internal/web"
)

//go:embed static
var staticFiles embed.FS

const (
	shareOutcomeAligned = "aligned"
	shareOutcomeApplied = "applied"
)

func main() { os.Exit(runCommand(os.Args[1:])) }

func runGateway(launch apphost.Launch) (runErr error) {
	startupCtx, startupCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer startupCancel()
	h := newGatewayHost(startupCtx, launch, startupCancel)
	go func() { <-startupCtx.Done(); h.requestQuit() }() //nolint:contextcheck // Shutdown has an independent shared notification deadline.
	defer func() {
		if recovered := recover(); recovered != nil {
			runErr = apphost.NewFault("startup.failed", fmt.Errorf("gateway panic: %v", recovered))
		}
		if errors.Is(runErr, context.Canceled) && startupCtx.Err() != nil {
			runErr = nil
		}
		runErr = finishGatewayResult(runErr, nil)
		reportErr := errors.Join(runErr, h.shellFailure())
		dialog := h.reportAsync(reportErr) //nolint:contextcheck // Fatal diagnostics outlive startup cancellation.
		closeErr := h.close()              //nolint:contextcheck // Cleanup observes its own single deadline.
		<-dialog
		runErr = errors.Join(runErr, closeErr, h.shellFailure())
		if reportErr == nil && runErr != nil && launch.Mode == apphost.Desktop {
			<-h.reportAsync(apphost.NewFault("shutdown.failed", runErr)) //nolint:contextcheck // Report actual cleanup failure after completion.
		}
	}()
	if err := prepareHost(h); err != nil {
		return err
	}
	if startupCtx.Err() != nil {
		return nil
	}
	return initializeGateway(h)
}

func newGatewayHost(ctx context.Context, launch apphost.Launch, cancel context.CancelFunc) *gatewayHost {
	h := &gatewayHost{launch: launch, startupCtx: ctx, startupCancel: cancel, quit: make(chan struct{}), startupDone: make(chan struct{}), shutdownDone: make(chan struct{}), newShell: desktop.New}
	h.prepareShutdown()
	return h
}

func initializeGateway(h *gatewayHost) (runErr error) {
	// =========================================================================
	// Database Setup (SQLite)
	// =========================================================================
	if err := preflightGroupFixtureDatabase(); err != nil {
		return fmt.Errorf("驗收 fixture 資料庫前置檢查失敗: %w", err)
	}
	sqliteDSN := embeddedSQLiteDSNForPath(h.owner.Identity().Path)
	log.Printf("SQLite DSN: %s", sqliteDSN)

	db, err := sql.Open("sqlite", sqliteDSN)
	if err != nil {
		return apphost.NewFault("startup.db_open_failed", err)
	}
	h.db = db
	datalink.ApplySQLitePoolDefaults(db)

	pingCtx, pingCancel := context.WithTimeout(h.startupCtx, 5*time.Second)
	defer pingCancel()
	if err := db.PingContext(pingCtx); err != nil {
		return apphost.NewFault("startup.db_ping_failed", err)
	}

	if err := verifyOpenedDatabase(h.startupCtx, db, h.owner.Identity()); err != nil {
		return apphost.NewFault("startup.path_unusable", err)
	}

	// =========================================================================
	// Migrations
	// =========================================================================
	migrator := datalink.NewMigrator()           // Remove db arg
	if err := migrator.Migrate(db); err != nil { // Add db arg
		return apphost.NewFault("startup.migration_failed", err)
	}

	// =========================================================================
	// Connector Manager
	// =========================================================================
	connMgr := connector.GetConnectionManager()
	h.connections = connMgr
	log.Println("ConnectionManager 已初始化")

	// 輸出已註冊的協議列表
	registeredProtocols := connector.ListProtocols()
	log.Printf("已註冊的協議: %v", registeredProtocols)
	if len(registeredProtocols) == 0 {
		log.Println("警告: 沒有協議被註冊！請檢查適配器包的導入。")
	}

	// =========================================================================
	// Repository & Service Wiring
	// =========================================================================

	services, err := wireGatewayServicesContext(h.startupCtx, db, connMgr, h.logs)
	if err != nil {
		return apphost.NewFault("startup.service_failed", err)
	}
	h.runtime, h.pipeline, h.share = services.runtime, services.groupPipe, services.modbusShare
	if h.startupCtx.Err() != nil {
		return nil
	}
	fixtureCleanup, err := services.fixture.configure(h.startupCtx, db, &services)
	if err != nil {
		return fmt.Errorf("設定驗收 fixture 失敗: %w", err)
	}
	h.fixtureCleanup = fixtureCleanup
	devSvc, pgSvc, pointSvc, tagSvc, mappingSvc := services.device, services.pollingGroup, services.point, services.tag, services.mapping
	settingsSvc, modbusShareSvc := services.settings, services.modbusShare
	workspaceSvc, auditSvc := services.workspace, services.audit
	dbTargetConnectorSvc, dbTargetMappingSvc := services.dbTarget, services.dbMapping
	scheduler, runtimeSvc, sourceRuleSvc := services.scheduler, services.runtime, services.sourceRule

	revisionStore := modbusshare.NewSQLWorkspaceRevisionStore(db)
	reconciler := modbusshare.NewReconciler(modbusShareSvc, revisionStore)
	configureShareRuntimeReconciler(sourceRuleSvc, workspaceSvc, modbusShareSvc, revisionStore, reconciler)
	hydrateStartupShare(h.startupCtx, services, revisionStore, reconciler)
	if h.startupCtx.Err() != nil {
		return nil
	}
	startConfiguredShareListener(h.startupCtx, modbusShareSvc)
	if hydration, err := modbusShareSvc.CheckHydration(h.startupCtx); err != nil || hydration.State != modbusshare.HydrationStateReady {
		h.emit("runtime.share_degraded", nil)
	}
	// Applied groups must be hydrated and their delivery recovered before any
	// sample can be accepted; running without the pipeline would silently drop
	// the data of groups that own their outputs, so a failed start is fatal.
	if err := services.groupPipe.Start(h.startupCtx); err != nil {
		return apphost.NewFault("startup.pipeline_failed", err)
	}

	if h.startupCtx.Err() != nil {
		return nil
	}
	h.startAcquisition(h.startupCtx, runtimeSvc)

	measurementSvc := measurement.NewService(measurement.NewSQLRepository(db))
	recordingPlanSvc := services.recordingPlan
	historySvc := history.NewService(history.NewMemoryHistoryRepository())

	datalinkServices := &api.DatalinkServices{
		Diagnostics:                h.logs,
		Device:                     devSvc,
		Point:                      pointSvc,
		Tag:                        tagSvc,
		Mapping:                    mappingSvc,
		PollingGroup:               pgSvc,
		Settings:                   settingsSvc,
		ModbusShare:                modbusShareSvc,
		ModbusShareReconciler:      reconciler,
		Scheduler:                  scheduler,
		Runtime:                    runtimeSvc,
		DBTarget:                   dbTargetConnectorSvc,
		DBMapping:                  dbTargetMappingSvc,
		SourceRule:                 sourceRuleSvc,
		Measurement:                measurementSvc,
		RecordingPlan:              recordingPlanSvc,
		History:                    historySvc,
		Workspace:                  workspaceSvc,
		WriteGroups:                services.writeGroups,
		WriteGroupDelivery:         services.groupPipe,
		WriteGroupTestWrite:        services.groupTestWrite,
		Audit:                      auditSvc,
		ShareRestore:               newWorkspaceShareRestoreBarrier(workspaceSvc, sourceRuleSvc, modbusShareSvc, reconciler),
		PointFormatRecovery:        newPointFormatRecovery(services, reconciler),
		RecordingStartShareBarrier: newRecordingStartShareBarrier(sourceRuleSvc, modbusShareSvc),
	}

	if h.startupCtx.Err() != nil {
		return nil
	}
	router := api.NewRouter(datalinkServices)

	web.SetupStaticFiles(router, staticFiles)

	return h.serve(router)
}
