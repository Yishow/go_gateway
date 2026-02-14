// Package api 提供 Datalink 模組的 REST API 路由註冊功能。
//
// 本套件整合所有 Datalink 子模組的 HTTP 處理器，並統一管理路由配置。
// 支援標準 RESTful 風格的 CRUD 操作以及特殊端點（連線測試、映射預覽等）。
package api

import (
	"net/http"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
)

// =============================================================================
// 路由器配置
// =============================================================================

// Config API 路由器配置
type Config struct {
	// PathPrefix API 路徑前綴 (預設 "/api/v1/datalink")
	PathPrefix string

	// EnableCORS 是否啟用 CORS
	EnableCORS bool

	// CORSOrigins 允許的 CORS 來源
	CORSOrigins []string
}

// DefaultConfig 預設配置
func DefaultConfig() Config {
	return Config{
		PathPrefix:  "/api/v1/datalink",
		EnableCORS:  true,
		CORSOrigins: []string{"*"},
	}
}

// =============================================================================
// 路由器
// =============================================================================

// Router Datalink API 路由器
type Router struct {
	config Config
	mux    *http.ServeMux

	// Handlers
	deviceHandler   *DeviceHandler
	pointHandler    *PointHandler
	tagHandler      *TagHandler
	mappingHandler  *MappingHandler
	settingsHandler *SettingsHandler
}

// Services 服務依賴
type Services struct {
	Device    *device.Service
	Point     *point.Service
	Tag       *tag.Service
	Mapping   *mapping.Service
	Scheduler *collector.Scheduler
	ConnMgr   *connector.ConnectionManager
	Storage   *storage.BatchWriter
}

// NewRouter 建立新的 API 路由器
func NewRouter(config Config, services Services) *Router {
	if config.PathPrefix == "" {
		config.PathPrefix = "/api/v1/datalink"
	}

	r := &Router{
		config: config,
		mux:    http.NewServeMux(),
	}

	// 初始化 Handlers
	r.deviceHandler = NewDeviceHandler(services.Device, services.ConnMgr)
	r.pointHandler = NewPointHandler(services.Point, services.Scheduler)
	r.tagHandler = NewTagHandler(services.Tag)
	r.mappingHandler = NewMappingHandler(services.Mapping, services.Scheduler)
	r.settingsHandler = NewSettingsHandler()

	// 註冊路由
	r.registerRoutes()

	return r
}
