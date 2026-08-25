package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TemplateHandler 處理測試模板相關的 API 請求
type TemplateHandler struct {
	templates map[string]TestTemplate
	mu        sync.RWMutex
}

// TestTemplate 測試模板
type TestTemplate struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Protocol    string                 `json:"protocol"`
	Config      map[string]interface{} `json:"config"`
	Operations  []Operation            `json:"operations"`
	Description string                 `json:"description,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// Operation 操作定義
type Operation struct {
	Type    string                 `json:"type"` // "read", "write", "wait"
	Params  map[string]interface{} `json:"params"`
	Comment string                 `json:"comment,omitempty"`
}

// NewTemplateHandler 建立新的模板處理器
func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{
		templates: make(map[string]TestTemplate),
	}
}

// List 取得測試模板列表
func (h *TemplateHandler) List(c *gin.Context) {
	h.mu.RLock()
	templates := make([]TestTemplate, 0, len(h.templates))
	for _, template := range h.templates {
		templates = append(templates, template)
	}
	h.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

// Get 取得單一測試模板
func (h *TemplateHandler) Get(c *gin.Context) {
	id := c.Param("id")

	h.mu.RLock()
	template, exists := h.templates[id]
	h.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

// SaveTemplateRequest 保存測試模板請求
type SaveTemplateRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Protocol    string                 `json:"protocol" binding:"required"`
	Config      map[string]interface{} `json:"config" binding:"required"`
	Operations  []Operation            `json:"operations" binding:"required"`
	Description string                 `json:"description,omitempty"`
}

// Save 保存測試模板
func (h *TemplateHandler) Save(c *gin.Context) {
	var req SaveTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template := TestTemplate{
		ID:          generateTemplateID(),
		Name:        req.Name,
		Protocol:    req.Protocol,
		Config:      req.Config,
		Operations:  req.Operations,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	h.mu.Lock()
	h.templates[template.ID] = template
	h.mu.Unlock()

	c.JSON(http.StatusOK, template)
}

// Delete 刪除測試模板
func (h *TemplateHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	h.mu.Lock()
	_, exists := h.templates[id]
	if exists {
		delete(h.templates, id)
	}
	h.mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
