package handlers

import (
	"net/http"
	"time"

	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

type PointHandler struct {
	svc         *point.Service
	runtimeSync pointRuntimeSyncer
}

type pointRuntimeSyncer interface {
	UpsertPoint(point *schema.Point)
	RemovePoint(pointID string)
}

func NewPointHandler(svc *point.Service, runtimeSync ...pointRuntimeSyncer) *PointHandler {
	var syncer pointRuntimeSyncer
	if len(runtimeSync) > 0 {
		syncer = runtimeSync[0]
	}
	return &PointHandler{svc: svc, runtimeSync: syncer}
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
	if h.runtimeSync != nil {
		h.runtimeSync.UpsertPoint(p)
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": p})
}

// BatchCreate 批次建立點位
// POST /datalink/points/batch
func (h *PointHandler) BatchCreate(c *gin.Context) {
	var req point.BatchCreatePointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	result, err := h.svc.BatchCreate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
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
	if h.runtimeSync != nil {
		h.runtimeSync.UpsertPoint(p)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": p})
}

func (h *PointHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	if h.runtimeSync != nil {
		h.runtimeSync.RemovePoint(id)
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// PollResult 輪詢結果
type PollResult struct {
	PointID          string      `json:"point_id"`
	Value            interface{} `json:"value"`
	TransformedValue interface{} `json:"transformed_value"`
	Timestamp        string      `json:"timestamp"`
	Quality          int         `json:"quality"`
	Stale            bool        `json:"stale"`
	Error            string      `json:"error,omitempty"`
}

// Poll 單點輪詢
// POST /datalink/points/:id/poll
func (h *PointHandler) Poll(c *gin.Context) {
	id := c.Param("id")

	// 取得點位資訊
	pt, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": "Point not found"}})
		return
	}

	// TODO: 實際呼叫協議連接器讀取值
	// 這裡返回模擬結果
	result := PollResult{
		PointID:          pt.ID,
		Value:            nil,
		TransformedValue: nil,
		Timestamp:        time.Now().Format(time.RFC3339),
		Quality:          192, // Good quality
		Stale:            false,
		Error:            "Poll not implemented - requires protocol connector integration",
	}

	if pt.LastValue != nil {
		result.Value = *pt.LastValue
		result.TransformedValue = *pt.LastValue
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// PollBatchRequest 批量輪詢請求
type PollBatchRequest struct {
	PointIDs *[]string `json:"point_ids"`
}

// PollBatch 批量輪詢
// POST /datalink/points/poll
func (h *PointHandler) PollBatch(c *gin.Context) {
	var req PollBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	// 驗證必填欄位（point_ids 欄位必須存在，即使是空陣列）
	if req.PointIDs == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": "point_ids field is required"}})
		return
	}

	results := make([]PollResult, 0, len(*req.PointIDs))

	for _, id := range *req.PointIDs {
		pt, err := h.svc.GetByID(c.Request.Context(), id)
		if err != nil {
			results = append(results, PollResult{
				PointID: id,
				Error:   "Point not found",
			})
			continue
		}

		result := PollResult{
			PointID:          pt.ID,
			TransformedValue: nil,
			Timestamp:        time.Now().Format(time.RFC3339),
			Quality:          192,
			Stale:            false,
			Error:            "Poll not implemented",
		}

		if pt.LastValue != nil {
			result.Value = *pt.LastValue
			result.TransformedValue = *pt.LastValue
		}

		results = append(results, result)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
}
