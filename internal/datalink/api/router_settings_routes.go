package api

import "net/http"

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
