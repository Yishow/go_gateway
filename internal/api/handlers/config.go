package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// ConfigHandler 處理配置管理相關的 API 請求
type ConfigHandler struct {
	presets map[string]ConfigPreset
	mu      sync.RWMutex
}

// ConfigPreset 配置預設值
type ConfigPreset struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Protocol    string                 `json:"protocol"`
	Config      map[string]interface{} `json:"config"`
	Description string                 `json:"description,omitempty"`
}

// NewConfigHandler 建立新的配置處理器
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		presets: make(map[string]ConfigPreset),
	}
}

// GetPresets 取得配置預設值列表
func (h *ConfigHandler) GetPresets(c *gin.Context) {
	protocol := c.Query("protocol")

	h.mu.RLock()
	presets := make([]ConfigPreset, 0)
	for _, preset := range h.presets {
		if protocol == "" || preset.Protocol == protocol {
			presets = append(presets, preset)
		}
	}
	h.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{"presets": presets})
}

// SavePresetRequest 保存配置預設值請求
type SavePresetRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Protocol    string                 `json:"protocol" binding:"required"`
	Config      map[string]interface{} `json:"config" binding:"required"`
	Description string                 `json:"description,omitempty"`
}

// SavePreset 保存配置預設值
func (h *ConfigHandler) SavePreset(c *gin.Context) {
	var req SavePresetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseErrorKey: err.Error()})
		return
	}

	preset := ConfigPreset{
		ID:          generatePresetID(),
		Name:        req.Name,
		Protocol:    req.Protocol,
		Config:      req.Config,
		Description: req.Description,
	}

	h.mu.Lock()
	h.presets[preset.ID] = preset
	h.mu.Unlock()

	c.JSON(http.StatusOK, preset)
}

// DeletePreset 刪除配置預設值
func (h *ConfigHandler) DeletePreset(c *gin.Context) {
	id := c.Param("id")

	h.mu.Lock()
	_, exists := h.presets[id]
	if exists {
		delete(h.presets, id)
	}
	h.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{apiResponseErrorKey: "preset not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseStatusKey: "deleted"})
}
