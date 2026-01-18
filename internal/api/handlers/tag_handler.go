package handlers

import (
	"context"
	"net/http"

	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	svc *tag.Service
}

func NewTagHandler() *TagHandler {
	repo := tag.NewMemoryRepository()
	svc := tag.NewService(repo)

	// Seed some tags
	list, _ := svc.List(context.Background(), tag.ListFilter{})
	if len(list) == 0 {
		svc.Create(context.Background(), tag.CreateTagRequest{
			Key:         "temp_c",
			DisplayName: "Temperature (Celsius)",
			DataType:    "float64",
			Unit:        "°C",
		})
		svc.Create(context.Background(), tag.CreateTagRequest{
			Key:         "pressure_bar",
			DisplayName: "Pressure (Bar)",
			DataType:    "float64",
			Unit:        "bar",
		})
	}

	return &TagHandler{svc: svc}
}

func (h *TagHandler) List(c *gin.Context) {
	filter := tag.ListFilter{}
	if s := c.Query("status"); s != "" {
		// filter.Status =
	}
	if s := c.Query("key_prefix"); s != "" {
		filter.KeyPrefix = s
	}

	tags, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tags})
}

func (h *TagHandler) Get(c *gin.Context) {
	id := c.Param("id")
	t, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": "Tag not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": t})
}

func (h *TagHandler) Create(c *gin.Context) {
	var req tag.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	t, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": t})
}

func (h *TagHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req tag.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	t, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": t})
}

func (h *TagHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// Activate 啟用標籤
// POST /datalink/tags/:id/activate
func (h *TagHandler) Activate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Activate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	t, _ := h.svc.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": t})
}

// Retire 退役標籤
// POST /datalink/tags/:id/retire
func (h *TagHandler) Retire(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Retire(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	t, _ := h.svc.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": t})
}

// BatchCreateRequest 批量建立請求
type BatchCreateRequest struct {
	Tags []tag.CreateTagRequest `json:"tags"`
}

// BatchCreateResponse 批量建立回應
type BatchCreateResponse struct {
	Created []string            `json:"created"`
	Errors  []BatchCreateError  `json:"errors"`
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
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	resp := BatchCreateResponse{
		Created: make([]string, 0),
		Errors:  make([]BatchCreateError, 0),
	}

	for _, tagReq := range req.Tags {
		t, err := h.svc.Create(c.Request.Context(), tagReq)
		if err != nil {
			resp.Errors = append(resp.Errors, BatchCreateError{
				Key:   tagReq.Key,
				Error: err.Error(),
			})
		} else {
			resp.Created = append(resp.Created, t.ID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
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
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
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

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

