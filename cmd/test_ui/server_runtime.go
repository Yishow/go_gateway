package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func startServer() {
	fullURL := "http://localhost" + serverAddr
	log.Printf("🚀 測試工具伺服器啟動於 %s", fullURL)
	log.Printf("📝 開啟瀏覽器訪問 %s 開始使用", fullURL)
	shouldOpen := os.Getenv("AUTO_OPEN_BROWSER") == "true" || os.Getenv("AUTO_OPEN_BROWSER") == "1"
	if shouldOpen {
		go func() {
			time.Sleep(500 * time.Millisecond)
			openBrowser(serverAddr)
		}()
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}

func handleSignals() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("收到系統信號，正在關閉伺服器...")
	shutdownServer()
	close(shutdownCh)
}

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

func openBrowser(serverAddr string) {
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
	default:
		cmd = exec.CommandContext(openCtx, "xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		log.Printf("無法打開瀏覽器: %v", err)
	}
}

func buildURL(serverAddr string) string {
	if strings.HasPrefix(serverAddr, "http://") || strings.HasPrefix(serverAddr, "https://") {
		return serverAddr
	}
	addr := strings.TrimPrefix(serverAddr, ":")
	if !strings.Contains(addr, ":") {
		return "http://localhost:" + addr
	}
	return "http://" + addr
}
