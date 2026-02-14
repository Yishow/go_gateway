package api

import "net/http"

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
