package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/point"

	"github.com/gin-gonic/gin"
)

type ModbusShareHandler struct {
	svc        *modbusshare.Service
	pointSvc   *point.Service
	mappingSvc *mapping.Service
}

func NewModbusShareHandler(svc *modbusshare.Service, pointSvc *point.Service, mappingSvc *mapping.Service) *ModbusShareHandler {
	return &ModbusShareHandler{
		svc:        svc,
		pointSvc:   pointSvc,
		mappingSvc: mappingSvc,
	}
}

// Register 使用指標：若為值型別 uint16，validator 的 required 會把合法位址 0（對應顯示 40001）當成零值而拒絕。
type UpsertMirrorMappingRequest struct {
	Register *uint16 `json:"register" binding:"required"`
}

type WriteTagValueRequest struct {
	TagID string      `json:"tag_id" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

type StartModbusShareRequest struct {
	Port int `json:"port"`
}

type SyncSummary struct {
	Updated int      `json:"updated"`
	Skipped int      `json:"skipped"`
	Errors  []string `json:"errors"`
}

func (h *ModbusShareHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": h.svc.Status()})
}

func (h *ModbusShareHandler) ListMappings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": h.svc.ListMappings()})
}

func (h *ModbusShareHandler) Start(c *gin.Context) {
	var req StartModbusShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	port := req.Port
	if port == 0 {
		port = 5020
	}

	if err := h.svc.Start(port); err != nil {
		code := http.StatusUnprocessableEntity
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "already in use") || strings.Contains(msg, "already running") || strings.Contains(msg, "運行中") {
			code = http.StatusConflict
		}
		c.JSON(code, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": h.svc.Status()})
}

func (h *ModbusShareHandler) Stop(c *gin.Context) {
	if err := h.svc.Stop(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": h.svc.Status()})
}

func (h *ModbusShareHandler) UpsertMapping(c *gin.Context) {
	tagID := c.Param("tagId")
	var req UpsertMirrorMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	m, err := h.svc.UpsertMapping(c.Request.Context(), tagID, *req.Register)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": m})
}

func (h *ModbusShareHandler) DeleteMapping(c *gin.Context) {
	tagID := c.Param("tagId")
	h.svc.RemoveMapping(tagID)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"deleted": true}})
}

func (h *ModbusShareHandler) WriteTagValue(c *gin.Context) {
	var req WriteTagValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	if err := h.svc.WriteTagValue(c.Request.Context(), req.TagID, req.Value); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"written": true}})
}

// SyncFromMappings syncs latest point values to modbus mirror by enabled mappings.
func (h *ModbusShareHandler) SyncFromMappings(c *gin.Context) {
	if h.pointSvc == nil || h.mappingSvc == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "error": gin.H{"message": "sync dependencies unavailable"}})
		return
	}

	enabled := true
	mappings, err := h.mappingSvc.List(c.Request.Context(), mapping.ListFilter{Enabled: &enabled, Limit: 10000})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	summary := SyncSummary{Errors: make([]string, 0)}

	for _, m := range mappings {
		if !h.svc.HasMapping(m.TagID) {
			summary.Skipped++
			continue
		}

		pt, err := h.pointSvc.GetByID(c.Request.Context(), m.PointID)
		if err != nil {
			summary.Errors = append(summary.Errors, "point not found: "+m.PointID)
			continue
		}
		if pt.LastValue == nil {
			summary.Skipped++
			continue
		}

		val := parsePointLastValue(*pt.LastValue)
		if err := h.svc.WriteTagValue(c.Request.Context(), m.TagID, val); err != nil {
			summary.Errors = append(summary.Errors, "tag "+m.TagID+": "+err.Error())
			continue
		}
		summary.Updated++
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": summary})
}

func parsePointLastValue(raw string) interface{} {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return raw
	}

	// Prefer scalar parsing first to avoid JSON number -> float64 precision loss.
	if i, err := strconv.ParseInt(trimmed, 10, 64); err == nil {
		return i
	}
	if u, err := strconv.ParseUint(trimmed, 10, 64); err == nil {
		return u
	}
	if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return f
	}
	if strings.EqualFold(trimmed, "true") {
		return true
	}
	if strings.EqualFold(trimmed, "false") {
		return false
	}

	var decoded interface{}
	if err := json.Unmarshal([]byte(trimmed), &decoded); err == nil {
		return decoded
	}
	return raw
}
