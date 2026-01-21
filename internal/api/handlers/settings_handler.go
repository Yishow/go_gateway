package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/settings"

	"github.com/gin-gonic/gin"
)

// SettingsHandler 設定 API Handler
type SettingsHandler struct {
	svc *settings.Service
}

// NewSettingsHandler 建立新的設定 Handler
func NewSettingsHandler(svc *settings.Service) *SettingsHandler {
	return &SettingsHandler{svc: svc}
}

// List 列出所有設定
// GET /datalink/settings
func (h *SettingsHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// UpdateRequest 更新設定請求
type UpdateSettingRequest struct {
	Value interface{} `json:"value"`
}

// Update 更新設定
// PUT /datalink/settings/:key
func (h *SettingsHandler) Update(c *gin.Context) {
	key := c.Param("key")

	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   gin.H{"message": err.Error()},
		})
		return
	}

	if err := h.svc.Set(c.Request.Context(), key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": err.Error()},
		})
		return
	}

	// 取得更新後的設定
	item, err := h.svc.Get(c.Request.Context(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   gin.H{"message": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": item})
}
