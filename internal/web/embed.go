package web

import (
	"embed"
	"io/fs"
	"net/http"

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

	// 靜態檔案服務
	router.StaticFS("/assets", http.FS(fsys))

	// SPA 路由處理：所有非 API 路徑都回傳 index.html
	router.NoRoute(func(c *gin.Context) {
		// 如果是 API 路徑，回傳 404
		if c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "API endpoint not found",
			})
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
