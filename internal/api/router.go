package api

import (
	"fmt"
	"net/http"

	"go-gateway/internal/api/handlers"
	"go-gateway/internal/config"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/settings"
	"go-gateway/internal/datalink/tag"

	_ "go-gateway/docs/swagger" // Swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// DatalinkServices 包含所有 Datalink 相關服務
type DatalinkServices struct {
	Device       *device.Service
	Point        *point.Service
	Tag          *tag.Service
	Mapping      *mapping.Service
	PollingGroup *pollinggroup.Service
	Settings     *settings.Service
}

// NewRouter 建立並配置 Gin 路由器
//
// Args:
//   - datalinkServices: Datalink 服務容器
//
// Returns:
//   - 配置好的 Gin Engine 實例
func NewRouter(datalinkServices *DatalinkServices) *gin.Engine {
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
	router.Use(customLoggerMiddleware()) // 使用自定義日誌中間件，過濾頻繁的 debug API 請求
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg))

	// API 路由群組（使用配置中的路徑）
	apiV1 := router.Group(cfg.API.BasePath)
	{
		// WebSocket 端點（保留用於其他用途）
		wsHandler := handlers.NewWebSocketHandler()
		apiV1.GET("/ws", wsHandler.HandleWebSocket)

		// SSE 端點（用於監控數據流）
		sseHandler := handlers.NewSSEHandler()

		// Debug 處理器
		debugHandler := handlers.NewDebugHandler()

		apiV1.GET("/test/monitor/stream", sseHandler.HandleMonitorStream)

		// 共用的 TestHandler 實例
		testHandler := handlers.NewTestHandler(wsHandler, sseHandler, debugHandler)

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

		// Debug 相關 API（使用共用的 debugHandler）
		debugGroup := apiV1.Group("/debug")
		{
			debugGroup.GET("/packets", debugHandler.GetPackets)
			debugGroup.GET("/logs", debugHandler.GetLogs)
			debugGroup.DELETE("/clear", debugHandler.ClearData)
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

		// Datalink API Group
		datalinkGroup := apiV1.Group("/datalink")
		{
			// Health
			healthHandler := handlers.NewDatalinkHealthHandler()
			datalinkGroup.GET("/health", healthHandler.Check)

			// Dashboard
			dashboardHandler := handlers.NewDashboardHandler(
				datalinkServices.Device,
				datalinkServices.Point,
				datalinkServices.Tag,
			)
			datalinkGroup.GET("/dashboard/stats", dashboardHandler.GetStats)
			datalinkGroup.GET("/dashboard/device-statuses", dashboardHandler.GetDeviceStatuses)

			// Protocols
			protocolHandler := handlers.NewProtocolHandler()
			datalinkGroup.GET("/protocols", protocolHandler.List)

			// Polling Groups
			pollingGroupHandler := handlers.NewPollingGroupHandler(datalinkServices.PollingGroup)
			datalinkGroup.GET("/polling-groups", pollingGroupHandler.List)
			datalinkGroup.POST("/polling-groups", pollingGroupHandler.Create)
			datalinkGroup.GET("/polling-groups/:id", pollingGroupHandler.Get)
			datalinkGroup.PUT("/polling-groups/:id", pollingGroupHandler.Update)
			datalinkGroup.DELETE("/polling-groups/:id", pollingGroupHandler.Delete)

			// Devices
			deviceHandler := handlers.NewDeviceHandler(datalinkServices.Device)
			datalinkGroup.GET("/devices", deviceHandler.List)
			datalinkGroup.POST("/devices", deviceHandler.Create)
			datalinkGroup.GET("/devices/:id", deviceHandler.Get)
			datalinkGroup.PUT("/devices/:id", deviceHandler.Update)
			datalinkGroup.DELETE("/devices/:id", deviceHandler.Delete)
			datalinkGroup.POST("/devices/:id/test", deviceHandler.TestConnection)
			datalinkGroup.POST("/devices/:id/activate", deviceHandler.Activate)
			datalinkGroup.POST("/devices/:id/disable", deviceHandler.Disable)
			datalinkGroup.POST("/devices/test-batch", deviceHandler.TestConnectionBatch)

			// Points
			pointHandler := handlers.NewPointHandler(datalinkServices.Point)
			datalinkGroup.GET("/points", pointHandler.List)
			datalinkGroup.POST("/points", pointHandler.Create)
			datalinkGroup.POST("/points/batch", pointHandler.BatchCreate)
			datalinkGroup.GET("/points/:id", pointHandler.Get)
			datalinkGroup.PUT("/points/:id", pointHandler.Update)
			datalinkGroup.DELETE("/points/:id", pointHandler.Delete)
			datalinkGroup.POST("/points/:id/poll", pointHandler.Poll)
			datalinkGroup.POST("/points/poll", pointHandler.PollBatch)

			// Tags
			tagHandler := handlers.NewTagHandler(datalinkServices.Tag)
			datalinkGroup.GET("/tags", tagHandler.List)
			datalinkGroup.POST("/tags", tagHandler.Create)
			datalinkGroup.GET("/tags/:id", tagHandler.Get)
			datalinkGroup.PUT("/tags/:id", tagHandler.Update)
			datalinkGroup.DELETE("/tags/:id", tagHandler.Delete)
			datalinkGroup.POST("/tags/:id/activate", tagHandler.Activate)
			datalinkGroup.POST("/tags/:id/retire", tagHandler.Retire)
			datalinkGroup.POST("/tags/batch", tagHandler.BatchCreate)
			datalinkGroup.POST("/tags/validate-key", tagHandler.ValidateKey)

			// Mappings
			mappingHandler := handlers.NewMappingHandler(datalinkServices.Mapping)
			datalinkGroup.GET("/mappings", mappingHandler.List)
			datalinkGroup.POST("/mappings", mappingHandler.Create)
			datalinkGroup.GET("/mappings/:id", mappingHandler.Get)
			datalinkGroup.PUT("/mappings/:id", mappingHandler.Update)
			datalinkGroup.DELETE("/mappings/:id", mappingHandler.Delete)
			datalinkGroup.POST("/mappings/preview", mappingHandler.Preview)
			datalinkGroup.POST("/mappings/validate-pipeline", mappingHandler.ValidatePipeline)

			// Settings
			settingsHandler := handlers.NewSettingsHandler(datalinkServices.Settings)
			datalinkGroup.GET("/settings", settingsHandler.List)
			datalinkGroup.PUT("/settings/:key", settingsHandler.Update)

			// SSE Preview Stream
			ssePreviewHandler := handlers.NewDatalinkSSEHandler()
			datalinkGroup.GET("/preview/stream", ssePreviewHandler.PreviewStream)
		}
	}

	// Swagger API 文檔
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// customLoggerMiddleware 自定義日誌中間件，過濾頻繁的 debug API 輪詢請求
func customLoggerMiddleware() gin.HandlerFunc {
	// 配置 Logger，跳過頻繁輪詢的 API 路徑
	// 這些是前端定期輪詢的 API，不需要每次都記錄日誌
	skipPaths := []string{
		"/api/v1/debug/packets",
		"/api/v1/debug/logs",
		"/api/v1/test/status", // 狀態檢查也可能頻繁輪詢
	}

	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: skipPaths,
		// 確保所有請求都被記錄（包括 /test/connect）
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[%s] %s %s %s %d %s \"%s\" %s\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.ClientIP,
				param.Method,
				param.Path,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
	})
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
