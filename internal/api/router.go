package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/api/handlers"
)

// NewRouter 建立並配置 Gin 路由器
//
// Returns:
//   - 配置好的 Gin Engine 實例
func NewRouter() *gin.Engine {
	// 設定 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// 中間件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// API 路由群組
	apiV1 := router.Group("/api/v1")
	{
		// WebSocket 端點
		wsHandler := handlers.NewWebSocketHandler()
		apiV1.GET("/ws", wsHandler.HandleWebSocket)
		apiV1.GET("/test/monitor/stream", wsHandler.HandleMonitorStream)

		// 共用的 TestHandler 實例
		testHandler := handlers.NewTestHandler(wsHandler)

		// 測試相關 API
		testGroup := apiV1.Group("/test")
		{
			testGroup.POST("/connect", testHandler.Connect)
			testGroup.POST("/disconnect", testHandler.Disconnect)
			testGroup.GET("/status", testHandler.GetStatus)
			testGroup.POST("/read", testHandler.Read)
			testGroup.POST("/write", testHandler.Write)
			testGroup.POST("/batch", testHandler.Batch)
			testGroup.POST("/script", testHandler.ExecuteScript)
			testGroup.GET("/scripts", testHandler.ListScripts)
			testGroup.POST("/scripts", testHandler.SaveScript)
			testGroup.DELETE("/scripts/:id", testHandler.DeleteScript)
			testGroup.POST("/monitor/start", testHandler.StartMonitor)
			testGroup.POST("/monitor/stop", testHandler.StopMonitor)
		}

		// Debug 相關 API
		debugGroup := apiV1.Group("/debug")
		{
			debugHandler := handlers.NewDebugHandler()
			debugGroup.GET("/packets", debugHandler.GetPackets)
			debugGroup.GET("/logs", debugHandler.GetLogs)
			debugGroup.POST("/send-raw", debugHandler.SendRaw)
			debugGroup.GET("/analyze/:packetId", debugHandler.AnalyzePacket)
		}

		// 配置管理 API
		configGroup := apiV1.Group("/config")
		{
			configHandler := handlers.NewConfigHandler()
			configGroup.GET("/presets", configHandler.GetPresets)
			configGroup.POST("/presets", configHandler.SavePreset)
			configGroup.DELETE("/presets/:id", configHandler.DeletePreset)
		}

		// 測試模板 API
		templateGroup := apiV1.Group("/templates")
		{
			templateHandler := handlers.NewTemplateHandler()
			templateGroup.GET("", templateHandler.List)
			templateGroup.POST("", templateHandler.Save)
			templateGroup.GET("/:id", templateHandler.Get)
			templateGroup.DELETE("/:id", templateHandler.Delete)
		}

		// 連線池 API
		connectionGroup := apiV1.Group("/connections")
		{
			// 使用共用的 testHandler
			connectionHandler := handlers.NewConnectionHandler(testHandler)
			connectionGroup.GET("", connectionHandler.List)
			connectionGroup.GET("/:id", connectionHandler.Get)
		}
	}

	return router
}

// corsMiddleware 處理 CORS 跨域請求
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
