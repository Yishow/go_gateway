package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	svc *tag.Service
}

func NewTagHandler(svc *tag.Service) *TagHandler {
	// Seed some tags if empty (optional logic moved here or kept in main, for now keeping to match logic)
	// But usually seeding belongs to main or a separate seeding func.
	// We will duplicate seeding logic check if we want to keep behavior, or remove it.
	// For Refactoring: keep it simple. Remove seeding from Handler.
	return &TagHandler{svc: svc}
}

func (h *TagHandler) List(c *gin.Context) {
	filter := tag.ListFilter{}
	// TODO: 支援 status 查詢篩選。
	if s := c.Query("key_prefix"); s != "" {
		filter.KeyPrefix = s
	}

	tags, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: tags})
}

func (h *TagHandler) Get(c *gin.Context) {
	id := c.Param("id")
	t, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: "Tag not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: t})
}

func (h *TagHandler) Create(c *gin.Context) {
	var req tag.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	t, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: t})
}

func (h *TagHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req tag.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	t, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: t})
}

func (h *TagHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true})
}

// Activate 啟用標籤
// POST /datalink/tags/:id/activate
func (h *TagHandler) Activate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Activate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	t, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: t})
}

// Retire 退役標籤
// POST /datalink/tags/:id/retire
func (h *TagHandler) Retire(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Retire(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	t, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: t})
}

// BatchCreateRequest 批量建立請求
type BatchCreateRequest struct {
	Tags *[]tag.CreateTagRequest `json:"tags"`
}

// BatchCreateResponse 批量建立回應
type BatchCreateResponse struct {
	Created []string           `json:"created"`
	Errors  []BatchCreateError `json:"errors"`
}

// BatchCreateError 批量建立錯誤
type BatchCreateError struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

// BatchCreate 批量建立標籤
// POST /datalink/tags/batch
func (h *TagHandler) BatchCreate(c *gin.Context) {
	var req BatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	// 驗證必填欄位（tags 欄位必須存在，即使是空陣列）
	if req.Tags == nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: "tags field is required"}})
		return
	}

	created, batchErrors := h.svc.BatchCreate(c.Request.Context(), *req.Tags)

	resp := BatchCreateResponse{
		Created: created,
		Errors:  make([]BatchCreateError, 0, len(batchErrors)),
	}
	for _, be := range batchErrors {
		resp.Errors = append(resp.Errors, BatchCreateError{Key: be.Key, Error: be.Error})
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: resp})
}

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

// ValidateKey 驗證標籤鍵
// POST /datalink/tags/validate-key
func (h *TagHandler) ValidateKey(c *gin.Context) {
	var req ValidateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	// 驗證必填欄位（key 欄位必須存在且非空）
	if req.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: "key field is required"}})
		return
	}

	resp := ValidateKeyResponse{
		Valid:      true,
		Normalized: tag.NormalizeTagKey(req.Key),
	}

	// 驗證格式
	if err := tag.ValidateTagKey(req.Key); err != nil {
		resp.Valid = false
		resp.Error = err.Error()
	} else {
		// 檢查是否已存在
		_, err := h.svc.GetByKey(c.Request.Context(), req.Key)
		resp.Exists = (err == nil)
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: resp})
}
