package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/point"

	"github.com/gin-gonic/gin"
)

type PointHandler struct {
	svc *point.Service
}

func NewPointHandler() *PointHandler {
	repo := point.NewMemoryRepository()
    // TODO: Pass actual group repo
	svc := point.NewService(repo, nil) 
	return &PointHandler{svc: svc}
}

func (h *PointHandler) List(c *gin.Context) {
	deviceID := c.Query("device_id")
    filter := point.ListFilter{}
    if deviceID != "" {
        filter.DeviceID = &deviceID
    }
    
	points, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": points})
}

func (h *PointHandler) Get(c *gin.Context) {
	id := c.Param("id")
	p, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": "Point not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": p})
}

func (h *PointHandler) Create(c *gin.Context) {
	var req point.CreatePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	p, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": p})
}

func (h *PointHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req point.UpdatePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	p, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": p})
}

func (h *PointHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
