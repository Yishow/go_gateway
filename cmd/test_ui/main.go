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
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"go-gateway/internal/api"
	"go-gateway/internal/config"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters" // 導入所有適配器以觸發 init() 註冊協議
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/web"
)

//go:embed static
var staticFiles embed.FS

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
	dbPath := "datalink.db" // Default to local file
	log.Printf("資料庫路徑: %s", dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("無法開啟資料庫: %v", err)
	}
	defer db.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := db.PingContext(pingCtx); err != nil {
		log.Fatalf("無法連接資料庫: %v", err)
	}

	// 此處啟用 WAL journal_mode 等優化可選

	// =========================================================================
	// Migrations
	// =========================================================================
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

	// Device
	devRepo := device.NewSQLRepository(db)
	devSvc := device.NewService(devRepo, connMgr)

	// Point
	pointRepo := point.NewSQLRepository(db)
	pointSvc := point.NewService(pointRepo, nil) // PollingGroupRepository 暫為 nil

	// Tag
	tagRepo := tag.NewSQLRepository(db)
	tagSvc := tag.NewService(tagRepo)

	// Mapping
	mappingRepo := mapping.NewSQLRepository(db)
	mappingSvc := mapping.NewServiceWithTagResolver(mappingRepo, tagSvc.GetByID)

	// Local Modbus Share (Tag -> Virtual Modbus Memory Grid)
	modbusShareSvc := modbusshare.NewService(tagSvc, 65536)
	if err := modbusShareSvc.Start(5020); err != nil {
		log.Printf("本機 Modbus 分享服務啟動失敗 (port 5020): %v", err)
	} else {
		log.Printf("本機 Modbus 分享服務已啟動: %s", modbusShareSvc.Status().Address)
	}
	defer func() {
		if err := modbusShareSvc.Stop(); err != nil {
			log.Printf("關閉本機 Modbus 分享服務失敗: %v", err)
		}
	}()

	// PollingGroup
	pgRepo := pollinggroup.NewSQLRepository(db)
	pgSvc := pollinggroup.NewService(pgRepo)

	// Settings
	settingsRepo := settings.NewSQLRepository(db)
	settingsSvc := settings.NewService(settingsRepo)

	// Container
	datalinkServices := &api.DatalinkServices{
		Device:       devSvc,
		Point:        pointSvc,
		Tag:          tagSvc,
		Mapping:      mappingSvc,
		PollingGroup: pgSvc,
		Settings:     settingsSvc,
		ModbusShare:  modbusShareSvc,
	}

	// 建立 API 路由器
	router := api.NewRouter(datalinkServices)

	// 設定靜態檔案服務（使用 embed）
	web.SetupStaticFiles(router, staticFiles)

	// 設定伺服器
	serverAddr = cfg.GetServerAddr()
	server = &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  0,                 // SSE 連接需要無讀取超時
		WriteTimeout: 0,                 // SSE 連接需要無寫入超時
		IdleTimeout:  120 * time.Second, // 空閒超時設為 120 秒
	}

	// 建立關閉通道
	shutdownCh = make(chan struct{})

	// 啟動伺服器
	go startServer()

	// 處理系統信號（Ctrl+C 等）
	go handleSignals()

	// 阻塞主線程，等待信號
	select {}
}

// startServer 在背景啟動 HTTP 伺服器
func startServer() {
	fullURL := "http://localhost" + serverAddr
	log.Printf("🚀 測試工具伺服器啟動於 %s", fullURL)
	log.Printf("📝 開啟瀏覽器訪問 %s 開始使用", fullURL)

	// 檢查是否自動開啟瀏覽器（預設為 false）
	autoOpenBrowser := os.Getenv("AUTO_OPEN_BROWSER")
	shouldOpen := autoOpenBrowser == "true" || autoOpenBrowser == "1"

	if shouldOpen {
		// 等待一小段時間確保伺服器已啟動，然後自動打開瀏覽器
		go func() {
			time.Sleep(500 * time.Millisecond)
			openBrowser(serverAddr)
		}()
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}

// handleSignals 處理系統信號（Ctrl+C 等）
func handleSignals() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("收到系統信號，正在關閉伺服器...")
	shutdownServer()
	os.Exit(0)
}

// shutdownServer 優雅關閉伺服器
func shutdownServer() {
	if server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("伺服器強制關閉: %v", err)
	} else {
		log.Println("伺服器已優雅關閉")
	}
}

// openBrowser 在預設瀏覽器中打開指定 URL
//
// Args:
//   - serverAddr: 伺服器地址（格式如 ":8080" 或 "localhost:8080"）
func openBrowser(serverAddr string) {
	// 構建完整的 URL
	url := buildURL(serverAddr)
	log.Printf("正在打開瀏覽器: %s", url)

	openCtx, openCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer openCancel()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(openCtx, "cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.CommandContext(openCtx, "open", url)
	default: // linux
		cmd = exec.CommandContext(openCtx, "xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		log.Printf("無法打開瀏覽器: %v", err)
	}
}

// buildURL 構建完整的 HTTP URL
//
// Args:
//   - serverAddr: 伺服器地址（格式如 ":8080" 或 "localhost:8080"）
//
// Returns:
//   - string: 完整的 URL，例如 "http://localhost:8080"
func buildURL(serverAddr string) string {
	// 檢查是否已經包含協議
	if strings.HasPrefix(serverAddr, "http://") || strings.HasPrefix(serverAddr, "https://") {
		return serverAddr
	}

	// 移除前導的冒號（如果有的話）
	addr := strings.TrimPrefix(serverAddr, ":")

	// 如果地址不包含主機名（只有端口號），添加 localhost
	if !strings.Contains(addr, ":") {
		// 只有端口號
		return "http://localhost:" + addr
	}

	// 如果地址包含主機名，添加 http:// 前綴
	return "http://" + addr
}
