package api

import (
	"context"
	"net/http"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// DeviceHandler 設備管理 API 處理器
// =============================================================================

// DeviceHandler 設備 API 處理器
type DeviceHandler struct {
	svc     *device.Service
	connMgr *connector.ConnectionManager
}

// NewDeviceHandler 建立新的設備處理器
func NewDeviceHandler(svc *device.Service, connMgr *connector.ConnectionManager) *DeviceHandler {
	return &DeviceHandler{
		svc:     svc,
		connMgr: connMgr,
	}
}

// =============================================================================
// CRUD 操作
// =============================================================================

// List 列出設備
// GET /devices?protocol=&status=&limit=&offset=
func (h *DeviceHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := device.ListFilter{
		Limit:  getQueryParamInt(r, "limit", 100),
		Offset: getQueryParamInt(r, "offset", 0),
	}

	// 篩選協議
	if protocol := getQueryParam(r, "protocol", ""); protocol != "" {
		p := schema.ProtocolType(protocol)
		filter.Protocol = &p
	}

	// 篩選狀態
	if status := getQueryParam(r, "status", ""); status != "" {
		s := schema.DeviceStatus(status)
		filter.Status = &s
	}

	devices, err := h.svc.List(ctx, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, devices)
}

// Get 取得單一設備
// GET /devices/{id}
func (h *DeviceHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	dev, err := h.svc.GetByID(ctx, id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, dev)
}

// Create 建立設備
// POST /devices
func (h *DeviceHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req device.CreateDeviceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	dev, err := h.svc.Create(ctx, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, dev)
}

// Update 更新設備
// PUT/PATCH /devices/{id}
func (h *DeviceHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req device.UpdateDeviceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	dev, err := h.svc.Update(ctx, id, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, dev)
}

// Delete 刪除設備
// DELETE /devices/{id}
func (h *DeviceHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := h.svc.Delete(ctx, id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

// =============================================================================
// 狀態管理
// =============================================================================

// TestConnection 測試設備連線
// POST /devices/{id}/test
func (h *DeviceHandler) TestConnection(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	result, err := h.svc.TestConnectionWithResult(ctx, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Activate 啟用設備
// POST /devices/{id}/activate
func (h *DeviceHandler) Activate(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := h.svc.Activate(ctx, id); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	dev, _ := h.svc.GetByID(ctx, id)
	writeJSON(w, http.StatusOK, dev)
}

// Disable 停用設備
// POST /devices/{id}/disable
func (h *DeviceHandler) Disable(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := h.svc.Disable(ctx, id); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	dev, _ := h.svc.GetByID(ctx, id)
	writeJSON(w, http.StatusOK, dev)
}

// =============================================================================
// 協議資訊
// =============================================================================

// ProtocolInfo 協議資訊回應
type ProtocolInfo struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ConfigSchema string `json:"config_schema"`
}

// ListProtocols 列出所有支援的協議
// GET /protocols
func (h *DeviceHandler) ListProtocols(w http.ResponseWriter, _ *http.Request) {
	protocols := connector.ListProtocolInfos()

	result := make([]ProtocolInfo, 0, len(protocols))
	for _, p := range protocols {
		result = append(result, ProtocolInfo{
			Type:         string(p.Type),
			Name:         p.Name,
			Description:  p.Description,
			ConfigSchema: string(p.ConfigSchema),
		})
	}

	writeJSON(w, http.StatusOK, result)
}

// =============================================================================
// 批量連線測試
// =============================================================================

// TestConnectionRequest 批量連線測試請求
type TestConnectionRequest struct {
	DeviceIDs []string `json:"device_ids"`
}

// TestConnectionResponse 連線測試結果
type TestConnectionResponse struct {
	DeviceID string `json:"device_id"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
	Latency  int64  `json:"latency_ms"`
}

// BatchTestConnections 批量測試設備連線
// POST /devices/test-batch
func (h *DeviceHandler) BatchTestConnections(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req TestConnectionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式")
		return
	}

	results := make([]TestConnectionResponse, 0, len(req.DeviceIDs))

	for _, id := range req.DeviceIDs {
		result, err := h.svc.TestConnectionWithResult(context.Background(), id)
		if err != nil {
			results = append(results, TestConnectionResponse{
				DeviceID: id,
				Success:  false,
				Error:    err.Error(),
			})
			continue
		}

		results = append(results, TestConnectionResponse{
			DeviceID: id,
			Success:  result.Success,
			Error:    result.Error,
			Latency:  result.LatencyMs,
		})
	}

	_ = ctx // 避免未使用警告
	writeJSON(w, http.StatusOK, results)
}
