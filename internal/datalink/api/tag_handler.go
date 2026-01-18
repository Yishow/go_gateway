package api

import (
	"net/http"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

// =============================================================================
// TagHandler 標籤管理 API 處理器
// =============================================================================

// TagHandler 標籤 API 處理器
type TagHandler struct {
	svc *tag.Service
}

// NewTagHandler 建立新的標籤處理器
func NewTagHandler(svc *tag.Service) *TagHandler {
	return &TagHandler{svc: svc}
}

// =============================================================================
// CRUD 操作
// =============================================================================

// List 列出標籤
// GET /tags?status=&data_type=&key_prefix=&limit=&offset=
func (h *TagHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter := tag.ListFilter{
		Limit:  getQueryParamInt(r, "limit", 100),
		Offset: getQueryParamInt(r, "offset", 0),
	}

	// 篩選狀態
	if status := getQueryParam(r, "status", ""); status != "" {
		s := schema.TagStatus(status)
		filter.Status = &s
	}

	// 篩選資料型別
	if dataType := getQueryParam(r, "data_type", ""); dataType != "" {
		dt := schema.DataType(dataType)
		filter.DataType = &dt
	}

	// 鍵前綴篩選
	filter.KeyPrefix = getQueryParam(r, "key_prefix", "")

	tags, err := h.svc.List(ctx, filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tags)
}

// Get 取得單一標籤
// GET /tags/{id}
func (h *TagHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	// 嘗試以 ID 取得
	t, err := h.svc.GetByID(ctx, id)
	if err != nil {
		// 嘗試以 Key 取得
		t, err = h.svc.GetByKey(ctx, id)
		if err != nil {
			writeError(w, http.StatusNotFound, "標籤不存在")
			return
		}
	}

	writeJSON(w, http.StatusOK, t)
}

// Create 建立標籤
// POST /tags
func (h *TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req tag.CreateTagRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	t, err := h.svc.Create(ctx, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

// Update 更新標籤
// PUT/PATCH /tags/{id}
func (h *TagHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	var req tag.UpdateTagRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式: "+err.Error())
		return
	}

	t, err := h.svc.Update(ctx, id, req)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, t)
}

// Delete 刪除標籤
// DELETE /tags/{id}
func (h *TagHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := h.svc.Delete(ctx, id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"deleted": true})
}

// =============================================================================
// 生命週期管理
// =============================================================================

// Activate 啟用標籤
// POST /tags/{id}/activate
func (h *TagHandler) Activate(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := h.svc.Activate(ctx, id); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	t, _ := h.svc.GetByID(ctx, id)
	writeJSON(w, http.StatusOK, t)
}

// Retire 退役標籤
// POST /tags/{id}/retire
func (h *TagHandler) Retire(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	if err := h.svc.Retire(ctx, id); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	t, _ := h.svc.GetByID(ctx, id)
	writeJSON(w, http.StatusOK, t)
}

// =============================================================================
// 批量操作
// =============================================================================

// BatchCreateRequest 批量建立標籤請求
type BatchCreateTagRequest struct {
	Tags []tag.CreateTagRequest `json:"tags"`
}

// BatchCreateResponse 批量建立結果
type BatchCreateTagResponse struct {
	Created []string       `json:"created"`
	Errors  []BatchTagError `json:"errors,omitempty"`
}

// BatchTagError 批量操作錯誤
type BatchTagError struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

// BatchCreate 批量建立標籤
// POST /tags/batch
func (h *TagHandler) BatchCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req BatchCreateTagRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式")
		return
	}

	result := BatchCreateTagResponse{
		Created: make([]string, 0),
		Errors:  make([]BatchTagError, 0),
	}

	for _, tagReq := range req.Tags {
		t, err := h.svc.Create(ctx, tagReq)
		if err != nil {
			result.Errors = append(result.Errors, BatchTagError{
				Key:   tagReq.Key,
				Error: err.Error(),
			})
		} else {
			result.Created = append(result.Created, t.ID)
		}
	}

	status := http.StatusCreated
	if len(result.Errors) > 0 {
		if len(result.Created) > 0 {
			status = http.StatusMultiStatus
		} else {
			status = http.StatusUnprocessableEntity
		}
	}

	writeJSON(w, status, result)
}

// =============================================================================
// 搜尋和驗證
// =============================================================================

// ValidateKeyRequest 驗證標籤鍵請求
type ValidateKeyRequest struct {
	Key string `json:"key"`
}

// ValidateKeyResponse 驗證標籤鍵回應
type ValidateKeyResponse struct {
	Valid      bool   `json:"valid"`
	Normalized string `json:"normalized"`
	Exists     bool   `json:"exists"`
	Error      string `json:"error,omitempty"`
}

// ValidateKey 驗證標籤鍵格式和唯一性
// POST /tags/validate-key
func (h *TagHandler) ValidateKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req ValidateKeyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "無效的請求格式")
		return
	}

	resp := ValidateKeyResponse{
		Normalized: tag.NormalizeTagKey(req.Key),
	}

	// 驗證格式
	if err := tag.ValidateTagKey(req.Key); err != nil {
		resp.Valid = false
		resp.Error = err.Error()
	} else {
		resp.Valid = true

		// 檢查是否已存在
		_, err := h.svc.GetByKey(ctx, req.Key)
		resp.Exists = err == nil
	}

	writeJSON(w, http.StatusOK, resp)
}
