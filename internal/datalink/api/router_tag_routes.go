package api

import "net/http"

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
