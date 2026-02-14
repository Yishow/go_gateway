package api

import (
	"net/http"
	"strings"
)

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

// handlePointsBatch 處理批量點位建立
func (r *Router) handlePointsBatch(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.pointHandler.BatchCreate(w, req)
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
