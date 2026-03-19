package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/gin-gonic/gin"
)

type SourceRuleHandler struct {
	svc *sourcerule.Service
}

type sourceRuleResponse struct {
	ID               string   `json:"id"`
	DeviceID         string   `json:"device_id"`
	StartAddress     string   `json:"start_address"`
	Count            int      `json:"count"`
	DataType         string   `json:"data_type"`
	NamingPrefix     string   `json:"naming_prefix"`
	Enabled          bool     `json:"enabled"`
	Locked           bool     `json:"locked"`
	Origin           string   `json:"origin"`
	TemplateName     string   `json:"template_name,omitempty"`
	SkippedAddresses []string `json:"skipped_addresses"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

func NewSourceRuleHandler(svc *sourcerule.Service) *SourceRuleHandler {
	return &SourceRuleHandler{svc: svc}
}

func (h *SourceRuleHandler) List(c *gin.Context) {
	filter := sourcerule.ListFilter{}
	if deviceID := c.Query("device_id"); deviceID != "" {
		filter.DeviceID = &deviceID
	}

	rules, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	payload := make([]sourceRuleResponse, 0, len(rules))
	for _, rule := range rules {
		payload = append(payload, mapSourceRuleResponse(rule))
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": payload})
}

func (h *SourceRuleHandler) Get(c *gin.Context) {
	rule, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, sourcerule.ErrSourceRuleNotFound) {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": mapSourceRuleResponse(rule)})
}

func (h *SourceRuleHandler) Create(c *gin.Context) {
	var req sourcerule.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	rule, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": mapSourceRuleResponse(rule)})
}

func (h *SourceRuleHandler) Update(c *gin.Context) {
	var req sourcerule.UpdateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	rule, err := h.svc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": mapSourceRuleResponse(rule)})
}

func (h *SourceRuleHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *SourceRuleHandler) Enable(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Enable(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	rule, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": mapSourceRuleResponse(rule)})
}

func (h *SourceRuleHandler) Disable(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Disable(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	rule, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": mapSourceRuleResponse(rule)})
}

func mapSourceRuleResponse(rule *schema.SourceRule) sourceRuleResponse {
	var skipped []string
	if rule != nil && rule.SkippedAddresses != "" {
		_ = json.Unmarshal([]byte(rule.SkippedAddresses), &skipped)
	}

	return sourceRuleResponse{
		ID:               rule.ID,
		DeviceID:         rule.DeviceID,
		StartAddress:     rule.StartAddress,
		Count:            rule.Count,
		DataType:         string(rule.DataType),
		NamingPrefix:     rule.NamingPrefix,
		Enabled:          rule.Enabled,
		Locked:           rule.Locked,
		Origin:           rule.Origin,
		TemplateName:     rule.TemplateName,
		SkippedAddresses: skipped,
		CreatedAt:        rule.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:        rule.UpdatedAt.Format(time.RFC3339Nano),
	}
}
