package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/modbusshare"

	"github.com/gin-gonic/gin"
)

type ModbusShareHandler struct {
	svc *modbusshare.Service
}

func NewModbusShareHandler(svc *modbusshare.Service) *ModbusShareHandler {
	return &ModbusShareHandler{svc: svc}
}

type UpsertMirrorMappingRequest struct {
	Register uint16 `json:"register" binding:"required"`
}

type WriteTagValueRequest struct {
	TagID string      `json:"tag_id" binding:"required"`
	Value interface{} `json:"value" binding:"required"`
}

func (h *ModbusShareHandler) Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": h.svc.Status()})
}

func (h *ModbusShareHandler) ListMappings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": h.svc.ListMappings()})
}

func (h *ModbusShareHandler) UpsertMapping(c *gin.Context) {
	tagID := c.Param("tagId")
	var req UpsertMirrorMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	m, err := h.svc.UpsertMapping(c.Request.Context(), tagID, req.Register)
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
