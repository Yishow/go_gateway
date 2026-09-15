package api

import (
	"encoding/json"
	"net/http"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// MappingHandler 映射管理 API 處理器
// =============================================================================

// MappingHandler 映射 API 處理器
type MappingHandler struct {
	svc       *mapping.Service
	scheduler *collector.Scheduler
}

// NewMappingHandler 建立新的映射處理器
func NewMappingHandler(svc *mapping.Service, scheduler *collector.Scheduler) *MappingHandler {
	return &MappingHandler{
		svc:       svc,
		scheduler: scheduler,
	}
}

// =============================================================================
// CRUD 操作
// =============================================================================

// List 列出映射
// GET /mappings?point_id=&tag_id=&enabled=&limit=&offset=
func (h *MappingHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := mapping.ListFilter{
		Limit:  getQueryParamInt(r, "limit", 100),
		Offset: getQueryParamInt(r, "offset", 0),
	}

	// 篩選點位
	if pointID := getQueryParam(r, "point_id"); pointID != "" {
		filter.PointID = &pointID
	}

	// 篩選標籤
	if tagID := getQueryParam(r, "tag_id"); tagID != "" {
		filter.TagID = &tagID
	}

	// 篩選啟用狀態
	if enabledStr := getQueryParam(r, "enabled"); enabledStr != "" {
		enabled := getQueryParamBool(r, "enabled", true)
		filter.Enabled = &enabled
	}

	mappings, err := h.svc.List(ctx, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, mappings)
}

// Get 取得單一映射
// GET /mappings/{id}
func (h *MappingHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	m, err := h.svc.GetByID(ctx, id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, m)
}

// Create 建立映射
// POST /mappings
func (h *MappingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req mapping.CreateMappingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	m, err := h.svc.Create(ctx, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, m)
}

// Update 更新映射
// PUT/PATCH /mappings/{id}
func (h *MappingHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req mapping.UpdateMappingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	m, err := h.svc.Update(ctx, id, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, m)
}

// Delete 刪除映射
// DELETE /mappings/{id}
func (h *MappingHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := h.svc.Delete(ctx, id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{responseDeletedKey: true})
}

// =============================================================================
// 映射預覽
// =============================================================================

// PreviewRequest 預覽請求
type PreviewRequest struct {
	// RawValue 原始值 (用於測試轉換)
	RawValue interface{} `json:"raw_value"`

	// PointID 點位 ID (可選，用於取得實際值)
	PointID string `json:"point_id,omitempty"`

	// TransformPipeline 轉換管線 (可選，直接測試管線)
	TransformPipeline []schema.TransformStep `json:"transform_pipeline,omitempty"`

	// MappingID 映射 ID (可選，使用現有映射的管線)
	MappingID string `json:"mapping_id,omitempty"`
}

// PreviewResponse 預覽回應
type PreviewResponse struct {
	RawValue    interface{}            `json:"raw_value"`
	FinalValue  interface{}            `json:"final_value"`
	StepResults []mapping.StepResult   `json:"step_results"`
	Pipeline    []schema.TransformStep `json:"pipeline,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// Preview 預覽映射轉換結果
// POST /mappings/preview
func (h *MappingHandler) Preview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req PreviewRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	resp := PreviewResponse{}

	// 取得原始值
	rawValue := req.RawValue
	if rawValue == nil && req.PointID != "" && h.scheduler != nil {
		// 從排程器取得實際值
		results := h.scheduler.PollNowContext(r.Context(), []string{req.PointID})
		if len(results) > 0 && results[0].Error == "" {
			rawValue = results[0].Value
		}
	}
	resp.RawValue = rawValue

	// 取得轉換管線
	var pipelineJSON string
	if len(req.TransformPipeline) > 0 {
		// 使用請求中的管線
		if err := mapping.ValidateTransformPipeline(req.TransformPipeline); err != nil {
			resp.Error = err.Error()
			writeJSON(w, http.StatusOK, resp)
			return
		}
		serializedPipeline, marshalErr := mustMarshalJSON(req.TransformPipeline)
		if marshalErr != nil {
			resp.Error = "序列化轉換管線失敗: " + marshalErr.Error()
			writeJSON(w, http.StatusOK, resp)
			return
		}
		pipelineJSON = serializedPipeline
		resp.Pipeline = req.TransformPipeline
	} else if req.MappingID != "" {
		// 使用現有映射的管線
		m, err := h.svc.GetByID(ctx, req.MappingID)
		if err != nil {
			resp.Error = "映射不存在: " + err.Error()
			writeJSON(w, http.StatusOK, resp)
			return
		}
		pipelineJSON = m.TransformPipeline
	}

	if rawValue == nil {
		resp.Error = "未提供原始值"
		writeJSON(w, http.StatusOK, resp)
		return
	}

	if pipelineJSON == "" {
		resp.FinalValue = rawValue
		writeJSON(w, http.StatusOK, resp)
		return
	}

	// 執行轉換
	result, err := mapping.ExecutePipeline(rawValue, pipelineJSON)
	if err != nil {
		resp.Error = err.Error()
		resp.StepResults = result.StepResults
		writeJSON(w, http.StatusOK, resp)
		return
	}

	resp.FinalValue = result.CurrentValue
	resp.StepResults = result.StepResults

	writeJSON(w, http.StatusOK, resp)
}

// =============================================================================
// 管線驗證
// =============================================================================

// ValidatePipelineRequest 驗證管線請求
type ValidatePipelineRequest struct {
	Pipeline []schema.TransformStep `json:"pipeline"`
}

// ValidatePipelineResponse 驗證管線回應
type ValidatePipelineResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// ValidatePipeline 驗證轉換管線
// POST /mappings/validate-pipeline
func (h *MappingHandler) ValidatePipeline(w http.ResponseWriter, r *http.Request) {
	var req ValidatePipelineRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式")
		return
	}

	resp := ValidatePipelineResponse{}
	if err := mapping.ValidateTransformPipeline(req.Pipeline); err != nil {
		resp.Valid = false
		resp.Error = err.Error()
	} else {
		resp.Valid = true
	}

	writeJSON(w, http.StatusOK, resp)
}

// =============================================================================
// 輔助函數
// =============================================================================

func mustMarshalJSON(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
