package api

import "net/http"

// registerRoutes 註冊所有路由
func (r *Router) registerRoutes() {
	prefix := r.config.PathPrefix

	// 設備管理 API
	r.mux.HandleFunc(prefix+"/devices", r.handleDevices)
	r.mux.HandleFunc(prefix+"/devices/", r.handleDeviceByID)
	r.mux.HandleFunc(prefix+"/devices/test-batch", r.handleDeviceBatchTest)

	// 點位管理 API
	r.mux.HandleFunc(prefix+"/points/poll", r.handlePointsBatchPoll)
	r.mux.HandleFunc(prefix+"/points/batch", r.handlePointsBatch)
	r.mux.HandleFunc(prefix+"/points", r.handlePoints)
	r.mux.HandleFunc(prefix+"/points/", r.handlePointByID)

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
