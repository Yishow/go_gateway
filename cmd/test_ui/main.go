package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-gateway/internal/api"
	"go-gateway/internal/config"
	"go-gateway/internal/web"

	"github.com/getlantern/systray"
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

	// 建立 API 路由器
	router := api.NewRouter()

	// 設定靜態檔案服務（使用 embed）
	web.SetupStaticFiles(router, staticFiles)

	// 設定伺服器
	serverAddr = cfg.GetServerAddr()
	server = &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 建立關閉通道
	shutdownCh = make(chan struct{})

	// 在背景啟動伺服器
	go startServer()

	// 處理系統信號（Ctrl+C 等）
	go handleSignals()

	// 啟動系統托盤（阻塞主線程）
	// systray.Run 必須在主線程中執行
	systray.Run(onReady, onExit)
}

// startServer 在背景啟動 HTTP 伺服器
func startServer() {
	fullURL := "http://localhost" + serverAddr
	log.Printf("🚀 測試工具伺服器啟動於 %s", fullURL)
	log.Printf("📝 開啟瀏覽器訪問 %s 開始使用", fullURL)
	log.Printf("💡 應用程式已最小化到系統托盤，點擊托盤圖示可打開瀏覽器")

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
	systray.Quit()
}

// onReady 系統托盤準備就緒時的回調函數
func onReady() {
	SetupTray(serverAddr)
}

// onExit 系統托盤退出時的回調函數
func onExit() {
	log.Println("系統托盤退出，正在關閉伺服器...")
	shutdownServer()
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
