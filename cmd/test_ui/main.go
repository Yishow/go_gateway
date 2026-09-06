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
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"

	"github.com/google/uuid"

	"go-gateway/internal/api"
	"go-gateway/internal/api/handlers"
	"go-gateway/internal/config"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters" // 導入所有適配器以觸發 init() 註冊協議
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/web"
)

//go:embed static
var staticFiles embed.FS

const (
	shareOutcomeAligned = "aligned"
	shareOutcomeApplied = "applied"
)

// 全域變數用於控制伺服器
var (
	server     *http.Server
	serverAddr string
	shutdownCh chan struct{}
)

func main() {
	// 載入配置（從 .env 文件或環境變數）
	cfg, err := config.Load()
	if err != nil {
		log.Printf("警告: 載入配置失敗，使用預設值: %v", err)
		cfg = config.Get()
	}

	// 設定日誌
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// =========================================================================
	// Database Setup (SQLite)
	// =========================================================================
	sqliteDSN := embeddedSQLiteDSN()
	log.Printf("SQLite DSN: %s", sqliteDSN)

	db, err := sql.Open("sqlite", sqliteDSN)
	if err != nil {
		log.Fatalf("無法開啟資料庫: %v", err)
	}
	defer db.Close()
	datalink.ApplySQLitePoolDefaults(db)

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("無法連接資料庫: %v", err)
	}

	// =========================================================================
	// Migrations
	// =========================================================================
	migrator := datalink.NewMigrator()           // Remove db arg
	if err := migrator.Migrate(db); err != nil { // Add db arg
		log.Fatalf("資料庫遷移失敗: %v", err)
	}

	// =========================================================================
	// Connector Manager
	// =========================================================================
	connMgr := connector.GetConnectionManager()
	defer func() {
		if err := connMgr.CloseAll(); err != nil {
			log.Printf("關閉 ConnectionManager 失敗: %v", err)
		}
	}()
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

	services := wireGatewayServices(db, connMgr)
	devSvc, pgSvc, pointSvc, tagSvc, mappingSvc := services.device, services.pollingGroup, services.point, services.tag, services.mapping
	settingsSvc, modbusShareSvc := services.settings, services.modbusShare
	shareSettingsLoaded, workspaceSvc, auditSvc := services.shareLoaded, services.workspace, services.audit
	dbTargetConnectorSvc, dbTargetMappingSvc := services.dbTarget, services.dbMapping
	scheduler, runtimeSvc, sourceRuleSvc := services.scheduler, services.runtime, services.sourceRule
	defer func() {
		if err := modbusShareSvc.CloseRuntime(); err != nil {
			log.Printf("關閉本機 Modbus 分享服務失敗: %v", err)
		}
	}()

	revisionStore := modbusshare.NewSQLWorkspaceRevisionStore(db)
	reconciler := modbusshare.NewReconciler(modbusShareSvc, revisionStore)
	configureShareRuntimeReconciler(sourceRuleSvc, workspaceSvc, modbusShareSvc, revisionStore, reconciler)
	workspaceRecord, workspaceErr := workspaceSvc.GetOrCreate(context.Background())
	//nolint:gocritic // Startup branches preserve the required hydration failure precedence.
	if !shareSettingsLoaded {
		modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateFailed, Readiness: false})
	} else if workspaceErr != nil {
		log.Printf("讀取 Studio V2 workspace 失敗，Share hydration 維持 failed: %v", workspaceErr)
		modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateFailed, Readiness: false})
	} else {
		reconciler.WithOwnershipValidator(sourcerule.NewShareDesiredMappingOwnershipChecker(workspaceSvc, sourceRuleSvc, tagSvc, modbusShareSvc.Settings, mappingSvc))

		settingsSnapshot := modbusShareSvc.Settings()
		workspaceRevision, _, revisionErr := revisionStore.GetRevision(context.Background(), workspaceRecord.ID)
		//nolint:gocritic // Revision, empty-state, and enabled-state checks intentionally remain ordered.
		if revisionErr != nil {
			log.Printf("讀取 Share workspace revision 失敗，hydration 維持 failed: %v", revisionErr)
			modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateFailed, WorkspaceID: workspaceRecord.ID, SettingsRevision: settingsSnapshot.SettingsRevision, Readiness: false})
		} else if workspaceRevision == "" {
			log.Printf("Share workspace revision 為空，hydration 維持 failed")
			modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateFailed, WorkspaceID: workspaceRecord.ID, SettingsRevision: settingsSnapshot.SettingsRevision, Readiness: false})
		} else if settingsSnapshot.Enabled {
			// Source-rule candidates are the sole restart authority. The durable
			// mapping snapshot is evidence for reconcile rollback/consistency,
			// never a competing desired-state source.
			if out, err := sourceRuleSvc.RestoreLocalModbusProjectionForDevices(context.Background(), workspaceRecord.ID, workspaceRevision, settingsSnapshot, workspaceRecord.OrderedDeviceIDs, reconciler); err != nil || (out.Outcome != shareOutcomeAligned && out.Outcome != shareOutcomeApplied) {
				if err != nil {
					log.Printf("本機 Modbus 分享服務 source-rule projection 還原失敗，hydration 維持 failed: %v", err)
				}
				modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateFailed, WorkspaceID: workspaceRecord.ID, WorkspaceRevision: workspaceRevision, SettingsRevision: settingsSnapshot.SettingsRevision, Readiness: false})
			} else {
				modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, WorkspaceID: workspaceRecord.ID, WorkspaceRevision: workspaceRevision, SettingsRevision: settingsSnapshot.SettingsRevision, Readiness: true, ReadinessToken: uuid.NewString()})
			}
		} else if shareSettingsLoaded {
			modbusShareSvc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, WorkspaceID: workspaceRecord.ID, WorkspaceRevision: workspaceRevision, SettingsRevision: settingsSnapshot.SettingsRevision, Readiness: true, ReadinessToken: uuid.NewString()})
		}
	}
	startConfiguredShareListener(context.Background(), modbusShareSvc)
	if err := runtimeSvc.Start(context.Background()); err != nil {
		log.Printf("datalink runtime 啟動失敗，runtime 功能將不可用: %v", err)
	}
	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := runtimeSvc.Stop(stopCtx); err != nil {
			log.Printf("關閉 datalink runtime 失敗: %v", err)
		}
	}()

	measurementSvc := measurement.NewService(measurement.NewSQLRepository(db))
	recordingPlanSvc := recordingplan.NewService(recordingplan.NewSQLRepository(db))

	datalinkServices := &api.DatalinkServices{
		Device:                devSvc,
		Point:                 pointSvc,
		Tag:                   tagSvc,
		Mapping:               mappingSvc,
		PollingGroup:          pgSvc,
		Settings:              settingsSvc,
		ModbusShare:           modbusShareSvc,
		ModbusShareReconciler: reconciler,
		Scheduler:             scheduler,
		Runtime:               runtimeSvc,
		DBTarget:              dbTargetConnectorSvc,
		DBMapping:             dbTargetMappingSvc,
		SourceRule:            sourceRuleSvc,
		Measurement:           measurementSvc,
		RecordingPlan:         recordingPlanSvc,
		Workspace:             workspaceSvc,
		Audit:                 auditSvc,
		ShareRestore: handlers.ShareRestoreBarrier(func(ctx context.Context, req handlers.ActivateWorkspaceRequest) error {
			hydration, err := modbusShareSvc.CheckHydration(ctx)
			if err != nil {
				return modbusshare.NewError(modbusshare.ErrCodeHydrationRequired, "workspace hydration is required before Share restore", true)
			}
			settingsSnapshot := modbusShareSvc.Settings()
			if !settingsSnapshot.Enabled {
				return modbusshare.NewError(modbusshare.ErrCodeDisabled, "Modbus Share is disabled in global settings", false)
			}
			if req.WorkspaceRevision == "" || req.WorkspaceRevision != hydration.WorkspaceRevision {
				return modbusshare.NewError(modbusshare.ErrCodeRevisionConflict, "workspace revision conflict", true)
			}
			// Restore always rebuilds desired state from persisted source-rule
			// candidates. Durable mapping rows are only consistency/rollback
			// evidence and must not become a competing authority.
			workspaceRecord, workspaceErr := workspaceSvc.GetOrCreate(ctx)
			if workspaceErr != nil {
				return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "workspace membership is unavailable", true)
			}
			var out modbusshare.ReconcileOutcome
			out, err = sourceRuleSvc.RestoreLocalModbusProjectionForDevices(ctx, hydration.WorkspaceID, req.WorkspaceRevision, settingsSnapshot, workspaceRecord.OrderedDeviceIDs, reconciler)
			if err == nil && out.Outcome != "aligned" && out.Outcome != "applied" {
				err = modbusshare.NewError(modbusshare.ErrCodeHydrationRequired, "Share restore did not reach an aligned projection", true)
			}
			return err
		}),
	}

	router := api.NewRouter(datalinkServices)

	web.SetupStaticFiles(router, staticFiles)

	serverAddr = cfg.GetServerAddr()
	server = &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  0,                 // SSE 連接需要無讀取超時
		WriteTimeout: 0,                 // SSE 連接需要無寫入超時
		IdleTimeout:  120 * time.Second, // 空閒超時設為 120 秒
	}

	shutdownCh = make(chan struct{})

	go startServer()

	go handleSignals()

	<-shutdownCh
}
