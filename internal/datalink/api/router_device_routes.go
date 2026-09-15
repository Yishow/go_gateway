package api

import "net/http"

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

	switch action {
	case "test":
		if req.Method == http.MethodPost {
			r.deviceHandler.TestConnection(w, req, id)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	case "activate":
		if req.Method == http.MethodPost {
			r.deviceHandler.Activate(w, req, id)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	case "disable":
		if req.Method == http.MethodPost {
			r.deviceHandler.Disable(w, req, id)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "方法不允許")
		}
	case "readiness":
		if req.Method == http.MethodPost {
			r.deviceHandler.CheckReadiness(w, req, id)
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
