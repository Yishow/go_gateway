package handlers

import (
	"context"
	"net/http"
	"strings"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

type PointHandler struct {
	svc             *point.Service
	runtimeSync     pointRuntimeSyncer
	manualPoller    pointManualPoller
	directReader    pointDirectReader
	mappingLister   pointMappingLister
	pollingGroupSvc pointPollingGroupGetter
}

type pointDirectReader interface {
	PollDirect(ctx context.Context, pointIDs []string) []collector.CollectedValue
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
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: points})
}

func (h *PointHandler) Get(c *gin.Context) {
	id := c.Param("id")
	p, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: pointNotFoundMessage}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: p})
}

func (h *PointHandler) Create(c *gin.Context) {
	var req point.CreatePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	p, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	if h.runtimeSync != nil {
		h.runtimeSync.UpsertPoint(p)
	}
	c.JSON(http.StatusCreated, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: p})
}

// BatchCreate 批次建立點位
// POST /datalink/points/batch
func (h *PointHandler) BatchCreate(c *gin.Context) {
	var req point.BatchCreatePointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	result, err := h.svc.BatchCreate(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

func (h *PointHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req point.UpdatePointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	p, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	if h.runtimeSync != nil {
		h.runtimeSync.UpsertPoint(p)
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: p})
}

func (h *PointHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}
	if h.runtimeSync != nil {
		h.runtimeSync.RemovePoint(id)
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true})
}

// PollResult 輪詢結果
type PollResult struct {
	PointID          string      `json:"point_id"`
	Value            interface{} `json:"value"`
	TransformedValue interface{} `json:"transformed_value"`
	Timestamp        string      `json:"timestamp"`
	Quality          int         `json:"quality"`
	Stale            bool        `json:"stale"`
	Error            string      `json:"error"`
}

// Poll 單點輪詢
// POST /datalink/points/:id/poll
func (h *PointHandler) Poll(c *gin.Context) {
	result, err := h.pollPoint(c.Request.Context(), c.Param("id"))
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := err.Error()
		if isPointPollNotFound(err) {
			statusCode = http.StatusNotFound
			message = pointNotFoundMessage
		}
		c.JSON(statusCode, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: message}})
		return
	}
	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: result})
}

// PollBatchRequest 批量輪詢請求
type PollBatchRequest struct {
	PointIDs *[]string `json:"point_ids"`
}

func isPointPollNotFound(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "點位不存在") || strings.Contains(message, "point not found")
}

// PollBatch 批量輪詢
// POST /datalink/points/poll
func (h *PointHandler) PollBatch(c *gin.Context) {
	var req PollBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	// 驗證必填欄位（point_ids 欄位必須存在，即使是空陣列）
	if req.PointIDs == nil {
		c.JSON(http.StatusBadRequest, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: "point_ids field is required"}})
		return
	}

	results, err := h.pollBatch(c.Request.Context(), *req.PointIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{apiResponseSuccessKey: false, apiResponseErrorKey: gin.H{apiResponseMessageKey: err.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: results})
}
