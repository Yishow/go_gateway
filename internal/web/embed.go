package web

import (
	"embed"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

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
		// 如果找不到 dist 目錄，可能是開發模式，使用空處理
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "前端檔案未找到，請先執行 npm run build",
			})
		})
		return
	}

	// 取得 assets 子目錄的檔案系統
	assetsFS, err := fs.Sub(fsys, "assets")
	if err != nil {
		// 如果找不到 assets 目錄，可能是構建配置問題
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "前端資源檔案未找到",
			})
		})
		return
	}

	// 靜態檔案服務處理函數
	serveAsset := func(c *gin.Context) {
		// 取得請求的檔案路徑
		// Gin 的 *filepath 參數會包含前導斜線，例如 "/index-CFH3XUSf.css"
		filePath := strings.TrimPrefix(c.Param("filepath"), "/")
		if filePath == "" {
			c.Status(http.StatusNotFound)
			return
		}

		// 開啟檔案
		file, err := assetsFS.Open(filePath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer file.Close()

		// 取得檔案資訊
		stat, err := file.Stat()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		// 如果是目錄，回傳 404
		if stat.IsDir() {
			c.Status(http.StatusNotFound)
			return
		}

		// 根據副檔名設定 MIME type
		ext := filepath.Ext(filePath)
		contentType := mime.TypeByExtension(ext)
		
		// 確保 CSS 和 JS 檔案有正確的 MIME type 和 charset
		if strings.HasSuffix(filePath, ".css") {
			contentType = "text/css; charset=utf-8"
		} else if strings.HasSuffix(filePath, ".js") || strings.HasSuffix(filePath, ".mjs") {
			contentType = "application/javascript; charset=utf-8"
		} else if contentType == "" {
			// 如果無法從副檔名判斷，使用預設值
			contentType = "application/octet-stream"
		} else if !strings.Contains(contentType, "charset") && (strings.HasPrefix(contentType, "text/") || strings.HasPrefix(contentType, "application/javascript")) {
			// 為文字類型的檔案添加 charset
			contentType += "; charset=utf-8"
		}

		// 設定 Cache-Control 標頭（可選）
		c.Header("Cache-Control", "public, max-age=31536000")

		// 回傳檔案內容
		c.DataFromReader(http.StatusOK, stat.Size(), contentType, file, nil)
	}

	// 註冊 GET 和 HEAD 方法（瀏覽器可能會先發送 HEAD 請求）
	router.GET("/assets/*filepath", serveAsset)
	router.HEAD("/assets/*filepath", serveAsset)

	// 處理根目錄的靜態資源（如 vite.svg）
	router.GET("/vite.svg", func(c *gin.Context) {
		file, err := fsys.Open("vite.svg")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.DataFromReader(http.StatusOK, stat.Size(), "image/svg+xml", file, nil)
	})

	// SPA 路由處理：所有非 API 路徑都回傳 index.html
	router.NoRoute(func(c *gin.Context) {
		// 如果是 API 路徑，回傳 404
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "API endpoint not found",
			})
			return
		}

		// 如果是靜態資源路徑，回傳 404
		if strings.HasPrefix(c.Request.URL.Path, "/assets") {
			c.Status(http.StatusNotFound)
			return
		}

		// 讀取 index.html
		indexFile, err := fsys.Open("index.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "無法載入前端頁面",
			})
			return
		}
		defer indexFile.Close()

		// 讀取檔案內容
		stat, err := indexFile.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "無法讀取前端檔案",
			})
			return
		}

		c.DataFromReader(http.StatusOK, stat.Size(), "text/html", indexFile, nil)
	})
}
