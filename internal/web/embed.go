package web

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const responseErrorKey = "error"

// SetupStaticFiles 設定靜態檔案服務
//
// Args:
//   - router: Gin 路由器
//   - staticFiles: 嵌入的靜態檔案系統
func SetupStaticFiles(router *gin.Engine, staticFiles embed.FS) {
	// 取得嵌入的檔案系統
	// embed 路徑是相對於包含 //go:embed 註解的檔案位置
	// 從 cmd/test_ui/main.go 來看，路徑是 static
	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Printf("❌ 無法載入 static 目錄: %v", err)
		return
	}

	// 取得 assets 子目錄的檔案系統
	assetsFS, err := fs.Sub(fsys, "assets")
	if err != nil {
		log.Printf("❌ 無法載入 assets 目錄: %v", err)
		return
	}

	// 手動處理靜態資源請求，確保 MIME type 正確
	serveAsset := func(c *gin.Context) {
		path := c.Param("filepath")
		// 移除前導斜線
		path = strings.TrimPrefix(path, "/")

		if path == "" {
			c.Status(http.StatusNotFound)
			return
		}

		// 嘗試開啟檔案
		file, err := assetsFS.Open(path)
		if err != nil {
			log.Printf("⚠️ 找不到靜態檔案: %s", path)
			c.Status(http.StatusNotFound)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			log.Printf("❌ 無法讀取檔案狀態: %s", path)
			c.Status(http.StatusInternalServerError)
			return
		}

		// 強制設定 Content-Type，不依賴系統偵測
		ext := strings.ToLower(filepath.Ext(path))
		var contentType string

		switch ext {
		case ".css":
			contentType = "text/css; charset=utf-8"
		case ".js", ".mjs":
			contentType = "application/javascript; charset=utf-8"
		case ".svg":
			contentType = "image/svg+xml"
		case ".json":
			contentType = "application/json"
		case ".png":
			contentType = "image/png"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		default:
			// 其他檔案類型回退到自動偵測
			contentType = "application/octet-stream"
		}

		c.Header("Content-Type", contentType)
		// 設定快取 (1年)
		c.Header("Cache-Control", "public, max-age=31536000")

		log.Printf("📦 Serving Asset: %s (%s)", path, contentType)

		readSeeker, ok := file.(io.ReadSeeker)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{responseErrorKey: "檔案不支援讀取"})
			return
		}
		http.ServeContent(c.Writer, c.Request, path, stat.ModTime(), readSeeker)
	}

	// 註冊 /assets 路由
	router.GET("/assets/*filepath", serveAsset)
	router.HEAD("/assets/*filepath", serveAsset)

	// 處理根目錄的 vite.svg
	router.GET("/vite.svg", func(c *gin.Context) {
		file, err := fsys.Open("vite.svg")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer file.Close()
		stat, err := file.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{responseErrorKey: "無法取得 vite.svg 檔案資訊"})
			return
		}
		readSeeker, ok := file.(io.ReadSeeker)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{responseErrorKey: "vite.svg 檔案不支援讀取"})
			return
		}
		c.Header("Content-Type", "image/svg+xml")
		http.ServeContent(c.Writer, c.Request, "vite.svg", stat.ModTime(), readSeeker)
	})

	// SPA 路由處理：所有非 API 且非資源的路徑都回傳 index.html
	router.NoRoute(func(c *gin.Context) {
		// API 路徑不處理
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{responseErrorKey: "API endpoint not found"})
			return
		}

		// Assets 路徑如果不匹配上面的 handler，則 404 (避免回傳 index.html 給 css)
		if strings.HasPrefix(c.Request.URL.Path, "/assets") {
			c.Status(http.StatusNotFound)
			return
		}

		// 讀取 index.html
		indexFile, err := fsys.Open("index.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{responseErrorKey: "無法載入 index.html"})
			return
		}
		defer indexFile.Close()

		stat, err := indexFile.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{responseErrorKey: "無法取得 index.html 檔案資訊"})
			return
		}
		readSeeker, ok := indexFile.(io.ReadSeeker)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{responseErrorKey: "index.html 檔案不支援讀取"})
			return
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), readSeeker)
	})
}
