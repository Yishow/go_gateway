package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/api/handlers"
	"go-gateway/internal/config"
)

// NewRouter 建立並配置 Gin 路由器
//
// Returns:
//   - 配置好的 Gin Engine 實例
func NewRouter() *gin.Engine {
	cfg := config.Get()

	// 設定 Gin 模式：根據配置決定是否使用 debug 模式
	// 如果 DEBUG=true 或 GIN_MODE=debug，則使用 debug 模式
	// 否則使用 release 模式（不顯示路由註冊資訊）
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 中間件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg))

	// API 路由群組（使用配置中的路徑）
	apiV1 := router.Group(cfg.API.BasePath)
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
func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 設定允許的來源
		origin := c.Request.Header.Get("Origin")
		allowOrigin := "*"
		
		// 如果配置了特定的來源列表，檢查是否匹配
		if len(cfg.CORS.AllowOrigins) > 0 && cfg.CORS.AllowOrigins[0] != "*" {
			allowOrigin = ""
			for _, allowed := range cfg.CORS.AllowOrigins {
				if allowed == "*" || allowed == origin {
					allowOrigin = origin
					break
				}
			}
		}
		
		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		}
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
