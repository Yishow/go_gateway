package api

import (
	"net/http"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// PointHandler 點位管理 API 處理器
// =============================================================================

// PointHandler 點位 API 處理器
type PointHandler struct {
	svc       *point.Service
	scheduler *collector.Scheduler
}

// NewPointHandler 建立新的點位處理器
func NewPointHandler(svc *point.Service, scheduler *collector.Scheduler) *PointHandler {
	return &PointHandler{
		svc:       svc,
		scheduler: scheduler,
	}
}

// =============================================================================
// 點位 CRUD 操作
// =============================================================================

// List 列出點位
// GET /points?device_id=&polling_group_id=&enabled=&limit=&offset=
func (h *PointHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := point.ListFilter{
		Limit:  getQueryParamInt(r, "limit", 100),
		Offset: getQueryParamInt(r, "offset", 0),
	}

	// 篩選設備
	if deviceID := getQueryParam(r, "device_id", ""); deviceID != "" {
		filter.DeviceID = &deviceID
	}

	// 篩選輪詢群組
	if groupID := getQueryParam(r, "polling_group_id", ""); groupID != "" {
		filter.PollingGroupID = &groupID
	}

	// 篩選啟用狀態
	if enabledStr := getQueryParam(r, "enabled", ""); enabledStr != "" {
		enabled := getQueryParamBool(r, "enabled", true)
		filter.Enabled = &enabled
	}

	// 篩選資料型別
	if dataType := getQueryParam(r, "data_type", ""); dataType != "" {
		dt := schema.DataType(dataType)
		filter.DataType = &dt
	}

	points, err := h.svc.List(ctx, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, points)
}

// Get 取得單一點位
// GET /points/{id}
func (h *PointHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	pt, err := h.svc.GetByID(ctx, id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, pt)
}

// Create 建立點位
// POST /points
func (h *PointHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req point.CreatePointRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	pt, err := h.svc.Create(ctx, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// 如果排程器存在且點位已啟用，將其加入排程
	if h.scheduler != nil && pt.Enabled {
		h.scheduler.AddPoint(pt)
	}

	writeJSON(w, http.StatusCreated, pt)
}

// BatchCreateRequest 批量建立點位請求
type BatchCreateRequest struct {
	DeviceID       string          `json:"device_id"`
	PollingGroupID string          `json:"polling_group_id"`
	DataType       schema.DataType `json:"data_type"`
	Enabled        bool            `json:"enabled"`
	DryRun         bool            `json:"dry_run,omitempty"`
	ApplyIfClean   bool            `json:"apply_if_clean,omitempty"`
	Points         []struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Function string `json:"function,omitempty"`
	} `json:"points"`
}

// BatchCreate 批量建立點位
// POST /points/batch
func (h *PointHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req BatchCreateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	points := make([]point.BatchPointItem, 0, len(req.Points))
	for _, p := range req.Points {
		points = append(points, point.BatchPointItem{
			Name:     p.Name,
			Address:  p.Address,
			Function: p.Function,
		})
	}

	result, err := h.svc.BatchCreate(ctx, point.BatchCreatePointsRequest{
		DeviceID:       req.DeviceID,
		PollingGroupID: req.PollingGroupID,
		DataType:       req.DataType,
		Enabled:        req.Enabled,
		DryRun:         req.DryRun,
		ApplyIfClean:   req.ApplyIfClean,
		Points:         points,
	})
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if h.scheduler != nil && result.Applied {
		for _, pt := range result.Points {
			if pt.Enabled {
				h.scheduler.AddPoint(pt)
			}
		}
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// Update 更新點位
// PUT/PATCH /points/{id}
func (h *PointHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req point.UpdatePointRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	pt, err := h.svc.Update(ctx, id, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// 更新排程器中的點位
	if h.scheduler != nil {
		if pt.Enabled {
			h.scheduler.AddPoint(pt)
		} else {
			h.scheduler.RemovePoint(pt.ID)
		}
	}

	writeJSON(w, http.StatusOK, pt)
}

// Delete 刪除點位
// DELETE /points/{id}
func (h *PointHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	// 從排程器移除
	if h.scheduler != nil {
		h.scheduler.RemovePoint(id)
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

// =============================================================================
// 即時輪詢
// =============================================================================

// PollNow 立即輪詢點位
// POST /points/{id}/poll
func (h *PointHandler) PollNow(w http.ResponseWriter, r *http.Request, id string) {
	if h.scheduler == nil {
		writeError(w, http.StatusServiceUnavailable, "排程器未啟用")
		return
	}

	results := h.scheduler.PollNow([]string{id})
	if len(results) == 0 {
		writeError(w, http.StatusNotFound, "點位不存在或未配置")
		return
	}

	writeJSON(w, http.StatusOK, results[0])
}

// BatchPollRequest 批量輪詢請求
type BatchPollRequest struct {
	PointIDs []string `json:"point_ids"`
}

// BatchPoll 批量輪詢點位
// POST /points/poll
func (h *PointHandler) BatchPoll(w http.ResponseWriter, r *http.Request) {
	if h.scheduler == nil {
		writeError(w, http.StatusServiceUnavailable, "排程器未啟用")
		return
	}

	var req BatchPollRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式")
		return
	}

	results := h.scheduler.PollNow(req.PointIDs)
	writeJSON(w, http.StatusOK, results)
}

// =============================================================================
// 輪詢群組 CRUD 操作
// =============================================================================

// ListPollingGroups 列出輪詢群組
// GET /polling-groups
func (h *PointHandler) ListPollingGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	groups, err := h.svc.ListPollingGroups(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, groups)
}

// CreatePollingGroup 建立輪詢群組
// POST /polling-groups
func (h *PointHandler) CreatePollingGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req point.CreatePollingGroupRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	group, err := h.svc.CreatePollingGroup(ctx, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// 將群組加入排程器
	if h.scheduler != nil {
		h.scheduler.AddPollingGroup(group)
	}

	writeJSON(w, http.StatusCreated, group)
}

// UpdatePollingGroup 更新輪詢群組
// PUT/PATCH /polling-groups/{id}
func (h *PointHandler) UpdatePollingGroup(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req point.UpdatePollingGroupRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	group, err := h.svc.UpdatePollingGroup(ctx, id, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// 更新排程器中的群組
	if h.scheduler != nil {
		h.scheduler.RemovePollingGroup(id)
		if group.Enabled {
			h.scheduler.AddPollingGroup(group)
		}
	}

	writeJSON(w, http.StatusOK, group)
}

// DeletePollingGroup 刪除輪詢群組
// DELETE /polling-groups/{id}
func (h *PointHandler) DeletePollingGroup(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	// 從排程器移除
	if h.scheduler != nil {
		h.scheduler.RemovePollingGroup(id)
	}

	if err := h.svc.DeletePollingGroup(ctx, id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

// =============================================================================
// 點位分組操作
// =============================================================================

// AssignToGroupRequest 分配點位到群組請求
type AssignToGroupRequest struct {
	PointID string `json:"point_id"`
	GroupID string `json:"group_id"`
}

// AssignToGroup 將點位分配到輪詢群組
// POST /polling-groups/{id}/points
func (h *PointHandler) AssignToGroup(w http.ResponseWriter, r *http.Request, groupID string) {
	ctx := r.Context()

	var req struct {
		PointID string `json:"point_id"`
	}
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式")
		return
	}

	if err := h.svc.AssignToGroup(ctx, req.PointID, groupID); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"assigned": true})
}

// RemoveFromGroup 將點位從輪詢群組移除
// DELETE /polling-groups/{group_id}/points/{point_id}
func (h *PointHandler) RemoveFromGroup(w http.ResponseWriter, r *http.Request, pointID string) {
	ctx := r.Context()

	if err := h.svc.RemoveFromGroup(ctx, pointID); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"removed": true})
}
