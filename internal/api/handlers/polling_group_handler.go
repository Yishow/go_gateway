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
func NewPollingGroupHandler() *PollingGroupHandler {
	repo := pollinggroup.NewMemoryRepository()
	// 預設資料
	_ = repo.Seed()
	svc := pollinggroup.NewService(repo)

	return &PollingGroupHandler{svc: svc}
}

// List 列出所有輪詢群組
// GET /datalink/polling-groups
func (h *PollingGroupHandler) List(c *gin.Context) {
	groups, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": groups})
}

// Get 取得單一輪詢群組
// GET /datalink/polling-groups/:id
func (h *PollingGroupHandler) Get(c *gin.Context) {
	id := c.Param("id")
	group, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   gin.H{"message": "Polling group not found"},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": group})
}

// Create 建立新輪詢群組
// POST /datalink/polling-groups
func (h *PollingGroupHandler) Create(c *gin.Context) {
	var req pollinggroup.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"message": err.Error()},
		})
		return
	}

	group, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": group})
}

// Update 更新輪詢群組
// PUT /datalink/polling-groups/:id
func (h *PollingGroupHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req pollinggroup.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"message": err.Error()},
		})
		return
	}

	group, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": group})
}

// Delete 刪除輪詢群組
// DELETE /datalink/polling-groups/:id
func (h *PollingGroupHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
