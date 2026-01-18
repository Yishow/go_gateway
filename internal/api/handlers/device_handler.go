package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	svc *device.Service
}

func NewDeviceHandler() *DeviceHandler {
	// In a real app, you'd inject the service. For now we use the memory repo default.
	// We need to ensure the service is initialized properly in main or here.
	// For simplicity in this prototype, we'll create a new service instance with memory repo.

	repo := device.NewMemoryRepository()
	// Pass nil for connection manager, service will use default singleton
	svc := device.NewService(repo, nil)

	// Pre-seed some data for demo if empty
	list, _ := svc.List(context.Background(), device.ListFilter{})
	if len(list) == 0 {
		_, err := svc.Create(context.Background(), device.CreateDeviceRequest{
			Name:             "Demo Modbus Device",
			Description:      "A simulated Modbus TCP device",
			Protocol:         "modbus_tcp",
			ConnectionConfig: map[string]interface{}{"host": "localhost", "port": 502, "slave_id": 1},
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to seed device: %v", err))
		}
	}

	return &DeviceHandler{svc: svc}
}

func (h *DeviceHandler) List(c *gin.Context) {
	devices, err := h.svc.List(c.Request.Context(), device.ListFilter{
		Protocol: func() *schema.ProtocolType {
			s := c.Query("protocol")
			if s == "" {
				return nil
			}
			p := schema.ProtocolType(s)
			return &p
		}(),
		Status: func() *schema.DeviceStatus {
			s := c.Query("status")
			if s == "" {
				return nil
			}
			st := schema.DeviceStatus(s)
			return &st
		}(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": devices})
}

func (h *DeviceHandler) Get(c *gin.Context) {
	id := c.Param("id")
	dev, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": "Device not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": dev})
}

func (h *DeviceHandler) Create(c *gin.Context) {
	var req device.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	dev, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": dev})
}

func (h *DeviceHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req device.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	dev, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": dev})
}

func (h *DeviceHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *DeviceHandler) TestConnection(c *gin.Context) {
	id := c.Param("id")
	result, err := h.svc.TestConnectionWithResult(c.Request.Context(), id)
	if err != nil {
		// Return structured failure
		c.JSON(http.StatusOK, gin.H{"success": true, "data": device.TestConnectionResult{
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (h *DeviceHandler) Activate(c *gin.Context) {
	id := c.Param("id")
	dev, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	// Quick hack: update status using Update method
	// Real implementation would have State Machine logic
	dev, _ = h.svc.Update(c.Request.Context(), id, device.UpdateDeviceRequest{
		// We need a way to specific status updates in the service, assuming service handles it or we add a method
	})

	c.JSON(http.StatusOK, gin.H{"success": true, "data": dev})
}

func (h *DeviceHandler) Disable(c *gin.Context) {
	id := c.Param("id")
	dev, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": "Device not found"}})
		return
	}

	// TODO: 實現真正的狀態更新邏輯
	c.JSON(http.StatusOK, gin.H{"success": true, "data": dev})
}

// TestConnectionBatchRequest 批量測試連線請求
type TestConnectionBatchRequest struct {
	DeviceIDs []string `json:"device_ids"`
}

// TestConnectionBatch 批量測試連線
// POST /datalink/devices/test-batch
func (h *DeviceHandler) TestConnectionBatch(c *gin.Context) {
	var req TestConnectionBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	results := make([]device.TestConnectionResult, 0, len(req.DeviceIDs))

	for _, id := range req.DeviceIDs {
		result, err := h.svc.TestConnectionWithResult(c.Request.Context(), id)
		if err != nil {
			results = append(results, device.TestConnectionResult{
				Success:   false,
				Error:     err.Error(),
				Timestamp: time.Now(),
			})
		} else {
			results = append(results, *result)
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}

