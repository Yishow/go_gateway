// Package api 提供 Datalink 模組的 REST API 路由註冊功能。
//
// 本套件整合所有 Datalink 子模組的 HTTP 處理器，並統一管理路由配置。
// 支援標準 RESTful 風格的 CRUD 操作以及特殊端點（連線測試、映射預覽等）。
package api

import (
	"net/http"
	"strings"

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

// registerRoutes 註冊所有路由
func (r *Router) registerRoutes() {
	prefix := r.config.PathPrefix

	// 設備管理 API
	r.mux.HandleFunc(prefix+"/devices", r.handleDevices)
	r.mux.HandleFunc(prefix+"/devices/", r.handleDeviceByID)
	r.mux.HandleFunc(prefix+"/devices/test-batch", r.handleDeviceBatchTest)

	// 點位管理 API
	r.mux.HandleFunc(prefix+"/points", r.handlePoints)
	r.mux.HandleFunc(prefix+"/points/", r.handlePointByID)
	r.mux.HandleFunc(prefix+"/points/poll", r.handlePointsBatchPoll)

	// 輪詢群組 API
	r.mux.HandleFunc(prefix+"/polling-groups", r.handlePollingGroups)
	r.mux.HandleFunc(prefix+"/polling-groups/", r.handlePollingGroupByID)

	// 標籤管理 API
	r.mux.HandleFunc(prefix+"/tags", r.handleTags)
	r.mux.HandleFunc(prefix+"/tags/", r.handleTagByID)
	r.mux.HandleFunc(prefix+"/tags/batch", r.handleTagsBatch)
	r.mux.HandleFunc(prefix+"/tags/validate-key", r.handleTagsValidateKey)

	// 映射管理 API
	r.mux.HandleFunc(prefix+"/mappings", r.handleMappings)
	r.mux.HandleFunc(prefix+"/mappings/", r.handleMappingByID)
	r.mux.HandleFunc(prefix+"/mappings/preview", r.handleMappingPreview)
	r.mux.HandleFunc(prefix+"/mappings/validate-pipeline", r.handleMappingValidatePipeline)

	// 系統設定 API
	r.mux.HandleFunc(prefix+"/settings", r.handleSettings)
	r.mux.HandleFunc(prefix+"/settings/", r.handleSettingByKey)

	// 協議資訊 API
	r.mux.HandleFunc(prefix+"/protocols", r.handleProtocols)

	// 健康檢查
	r.mux.HandleFunc(prefix+"/health", r.handleHealth)
}

// ServeHTTP 實作 http.Handler 介面
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// CORS 處理
	if r.config.EnableCORS {
		r.setCORSHeaders(w)
		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	r.mux.ServeHTTP(w, req)
}

// setCORSHeaders 設定 CORS 標頭
func (r *Router) setCORSHeaders(w http.ResponseWriter) {
	origins := "*"
	if len(r.config.CORSOrigins) > 0 && r.config.CORSOrigins[0] != "*" {
		origins = r.config.CORSOrigins[0]
	}
	w.Header().Set("Access-Control-Allow-Origin", origins)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "86400")
}

// =============================================================================
// 路由處理
// =============================================================================

// handleDevices 處理設備列表/建立
func (r *Router) handleDevices(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.deviceHandler.List(w, req)
	case http.MethodPost:
		r.deviceHandler.Create(w, req)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handleDeviceByID 處理單一設備操作
func (r *Router) handleDeviceByID(w http.ResponseWriter, req *http.Request) {
	// 解析路徑中的 ID 和操作
	path := req.URL.Path
	id, action := parseIDAndAction(path, r.config.PathPrefix+"/devices/")

	switch {
	case action == "test":
		if req.Method == http.MethodPost {
			r.deviceHandler.TestConnection(w, req, id)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	case action == "activate":
		if req.Method == http.MethodPost {
			r.deviceHandler.Activate(w, req, id)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	case action == "disable":
		if req.Method == http.MethodPost {
			r.deviceHandler.Disable(w, req, id)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	default:
		switch req.Method {
		case http.MethodGet:
			r.deviceHandler.Get(w, req, id)
		case http.MethodPut, http.MethodPatch:
			r.deviceHandler.Update(w, req, id)
		case http.MethodDelete:
			r.deviceHandler.Delete(w, req, id)
		default:
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	}
}

// handleDeviceBatchTest 批量測試設備連線
func (r *Router) handleDeviceBatchTest(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.deviceHandler.BatchTestConnections(w, req)
		return
	}

	writeError(w, http.StatusMethodNotAllowed, "方法不允許")
}

// handlePoints 處理點位列表/建立
func (r *Router) handlePoints(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.pointHandler.List(w, req)
	case http.MethodPost:
		r.pointHandler.Create(w, req)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handlePointByID 處理單一點位操作
func (r *Router) handlePointByID(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	id, action := parseIDAndAction(path, r.config.PathPrefix+"/points/")

	switch {
	case action == "poll":
		if req.Method == http.MethodPost {
			r.pointHandler.PollNow(w, req, id)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	default:
		switch req.Method {
		case http.MethodGet:
			r.pointHandler.Get(w, req, id)
		case http.MethodPut, http.MethodPatch:
			r.pointHandler.Update(w, req, id)
		case http.MethodDelete:
			r.pointHandler.Delete(w, req, id)
		default:
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	}
}

// handlePointsBatchPoll 批量輪詢點位
func (r *Router) handlePointsBatchPoll(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.pointHandler.BatchPoll(w, req)
		return
	}

	writeError(w, http.StatusMethodNotAllowed, "方法不允許")
}

// handlePollingGroups 處理輪詢群組列表/建立
func (r *Router) handlePollingGroups(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.pointHandler.ListPollingGroups(w, req)
	case http.MethodPost:
		r.pointHandler.CreatePollingGroup(w, req)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handlePollingGroupByID 處理單一輪詢群組操作
func (r *Router) handlePollingGroupByID(w http.ResponseWriter, req *http.Request) {
	relativePath := strings.TrimPrefix(req.URL.Path, r.config.PathPrefix+"/polling-groups/")
	relativePath = strings.Trim(relativePath, "/")
	parts := strings.Split(relativePath, "/")

	if len(parts) < 1 || parts[0] == "" {
		writeError(w, http.StatusNotFound, "群組不存在")
		return
	}

	groupID := parts[0]
	if len(parts) >= 2 && parts[1] == "points" {
		switch req.Method {
		case http.MethodPost:
			r.pointHandler.AssignToGroup(w, req, groupID)
		case http.MethodDelete:
			if len(parts) < 3 || parts[2] == "" {
				writeError(w, http.StatusBadRequest, "缺少點位 ID")
				return
			}
			r.pointHandler.RemoveFromGroup(w, req, parts[2])
		default:
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
		return
	}

	switch req.Method {
	case http.MethodPut, http.MethodPatch:
		r.pointHandler.UpdatePollingGroup(w, req, groupID)
	case http.MethodDelete:
		r.pointHandler.DeletePollingGroup(w, req, groupID)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handleTags 處理標籤列表/建立
func (r *Router) handleTags(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.tagHandler.List(w, req)
	case http.MethodPost:
		r.tagHandler.Create(w, req)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handleTagByID 處理單一標籤操作
func (r *Router) handleTagByID(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	id, action := parseIDAndAction(path, r.config.PathPrefix+"/tags/")

	switch {
	case action == "activate":
		if req.Method == http.MethodPost {
			r.tagHandler.Activate(w, req, id)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	case action == "retire":
		if req.Method == http.MethodPost {
			r.tagHandler.Retire(w, req, id)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	default:
		switch req.Method {
		case http.MethodGet:
			r.tagHandler.Get(w, req, id)
		case http.MethodPut, http.MethodPatch:
			r.tagHandler.Update(w, req, id)
		case http.MethodDelete:
			r.tagHandler.Delete(w, req, id)
		default:
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	}
}

// handleTagsBatch 批量建立標籤
func (r *Router) handleTagsBatch(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.tagHandler.BatchCreate(w, req)
		return
	}

	writeError(w, http.StatusMethodNotAllowed, "方法不允許")
}

// handleTagsValidateKey 驗證標籤鍵格式與唯一性
func (r *Router) handleTagsValidateKey(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.tagHandler.ValidateKey(w, req)
		return
	}

	writeError(w, http.StatusMethodNotAllowed, "方法不允許")
}

// handleMappings 處理映射列表/建立
func (r *Router) handleMappings(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.mappingHandler.List(w, req)
	case http.MethodPost:
		r.mappingHandler.Create(w, req)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handleMappingByID 處理單一映射操作
func (r *Router) handleMappingByID(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	id, _ := parseIDAndAction(path, r.config.PathPrefix+"/mappings/")

	switch req.Method {
	case http.MethodGet:
		r.mappingHandler.Get(w, req, id)
	case http.MethodPut, http.MethodPatch:
		r.mappingHandler.Update(w, req, id)
	case http.MethodDelete:
		r.mappingHandler.Delete(w, req, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handleMappingPreview 處理映射預覽
func (r *Router) handleMappingPreview(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.mappingHandler.Preview(w, req)
	} else {
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handleMappingValidatePipeline 驗證轉換管線
func (r *Router) handleMappingValidatePipeline(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.mappingHandler.ValidatePipeline(w, req)
		return
	}

	writeError(w, http.StatusMethodNotAllowed, "方法不允許")
}

// handleSettings 處理系統設定
func (r *Router) handleSettings(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.settingsHandler.List(w, req)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handleSettingByKey 處理單一設定操作
func (r *Router) handleSettingByKey(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	key, _ := parseIDAndAction(path, r.config.PathPrefix+"/settings/")

	switch req.Method {
	case http.MethodGet:
		r.settingsHandler.Get(w, req, key)
	case http.MethodPut, http.MethodPatch:
		r.settingsHandler.Update(w, req, key)
	default:
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handleProtocols 處理協議資訊查詢
func (r *Router) handleProtocols(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		r.deviceHandler.ListProtocols(w, req)
	} else {
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
	}
}

// handleHealth 健康檢查
func (r *Router) handleHealth(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"service": "datalink",
	})
}
