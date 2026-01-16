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
)

//go:embed static
var staticFiles embed.FS

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
	serverAddr := cfg.GetServerAddr()
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 啟動伺服器
	go func() {
		log.Printf("🚀 測試工具伺服器啟動於 http://localhost%s", serverAddr)
		log.Printf("📝 開啟瀏覽器訪問 http://localhost%s 開始使用", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("伺服器啟動失敗: %v", err)
		}
	}()

	// 優雅關閉
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在關閉伺服器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("伺服器強制關閉: %v", err)
	}

	log.Println("伺服器已關閉")
}
