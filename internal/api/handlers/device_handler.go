package handlers

import (
	"net/http"
	"strings"
	"time"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	svc *device.Service
}

type TestDraftConnectionRequest struct {
	Protocol         schema.ProtocolType    `json:"protocol"`
	ConnectionConfig map[string]interface{} `json:"connection_config"`
}

func NewDeviceHandler(svc *device.Service) *DeviceHandler {
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
			Success:     false,
			Error:       err.Error(),
			Timestamp:   time.Now(),
			Connect:     device.TestConnectionStageResult{Status: device.TestConnectionStageFailed, Error: err.Error()},
			Probe:       device.TestConnectionStageResult{Status: device.TestConnectionStageSkipped},
			CanActivate: false,
			CanCollect:  false,
		}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (h *DeviceHandler) TestDraftConnection(c *gin.Context) {
	var req TestDraftConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	result, err := h.svc.TestDraftConnectionWithResult(c.Request.Context(), device.DraftTestConnectionRequest{
		Protocol:         req.Protocol,
		ConnectionConfig: req.ConnectionConfig,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (h *DeviceHandler) Activate(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Activate(c.Request.Context(), id); err != nil {
		// 區分錯誤類型
		if strings.Contains(err.Error(), "不存在") || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		}
		return
	}

	dev, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": "Device not found"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": dev})
}

func (h *DeviceHandler) Disable(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Disable(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": "Device not found"}})
		return
	}

	dev, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": "Device not found"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": dev})
}

// TestConnectionBatchRequest 批量測試連線請求
type TestConnectionBatchRequest struct {
	DeviceIDs *[]string `json:"device_ids"`
}

// TestConnectionBatch 批量測試連線
// POST /datalink/devices/test-batch
func (h *DeviceHandler) TestConnectionBatch(c *gin.Context) {
	var req TestConnectionBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	// 驗證必填欄位（device_ids 欄位必須存在，即使是空陣列）
	if req.DeviceIDs == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": "device_ids field is required"}})
		return
	}

	results := make([]device.TestConnectionResult, 0, len(*req.DeviceIDs))

	for _, id := range *req.DeviceIDs {
		result, err := h.svc.TestConnectionWithResult(c.Request.Context(), id)
		if err != nil {
			results = append(results, device.TestConnectionResult{
				LatencyMs:   0,
				Success:     false,
				Error:       err.Error(),
				Timestamp:   time.Now(),
				Connect:     device.TestConnectionStageResult{Status: device.TestConnectionStageFailed, Error: err.Error()},
				Probe:       device.TestConnectionStageResult{Status: device.TestConnectionStageSkipped},
				CanActivate: false,
				CanCollect:  false,
			})
		} else {
			results = append(results, *result)
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}
