package handlers

import (
	"context"
	"net/http"

	"go-gateway/internal/datalink/pollinggroup"

	"github.com/gin-gonic/gin"
)

// PollingGroupHandler 輪詢群組 API Handler
type PollingGroupHandler struct {
	svc *pollinggroup.Service
}

// NewPollingGroupHandler 建立新的輪詢群組 Handler
func NewPollingGroupHandler(svc *pollinggroup.Service) *PollingGroupHandler {
	// Seed 預設資料 logic removed or moved to main/service seed method if needed.
	// For now we just inject.
	return &PollingGroupHandler{svc: svc}
}

// List 列出所有輪詢群組
// GET /datalink/polling-groups
func (h *PollingGroupHandler) List(c *gin.Context) {
	groups, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: groups})
}

// Get 取得單一輪詢群組
// GET /datalink/polling-groups/:id
func (h *PollingGroupHandler) Get(c *gin.Context) {
	id := c.Param("id")
	group, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: "Polling group not found"},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: group})
}

// Create 建立新輪詢群組
// POST /datalink/polling-groups
func (h *PollingGroupHandler) Create(c *gin.Context) {
	var req pollinggroup.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: err.Error()},
		})
		return
	}

	// 驗證必填欄位
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: "name field is required"},
		})
		return
	}

	group, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: err.Error()},
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: group})
}

// Update 更新輪詢群組
// PUT /datalink/polling-groups/:id
func (h *PollingGroupHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req pollinggroup.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: err.Error()},
		})
		return
	}

	group, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: group})
}

// Delete 刪除輪詢群組
// DELETE /datalink/polling-groups/:id
func (h *PollingGroupHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			apiResponseSuccessKey: false,
			apiResponseErrorKey:   gin.H{apiResponseMessageKey: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true})
}
