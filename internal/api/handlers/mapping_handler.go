package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

type MappingHandler struct {
	svc            *mapping.Service
	runtimeRefresh mappingRuntimeRefresher
}

type mappingRuntimeRefresher interface {
	RefreshMappings(ctx context.Context) error
}

func NewMappingHandler(svc *mapping.Service, runtimeRefresh ...mappingRuntimeRefresher) *MappingHandler {
	var refresher mappingRuntimeRefresher
	if len(runtimeRefresh) > 0 {
		refresher = runtimeRefresh[0]
	}
	return &MappingHandler{svc: svc, runtimeRefresh: refresher}
}

// ... existing List/Get/Create/Update/Delete methods ... (Assuming they will be kept but need struct update)
// Since I am replacing the struct definition, I should verify the rest of the file remains valid or if I need to update it too.
// The tool replaces a block. I will replace the top block.

func (h *MappingHandler) List(c *gin.Context) {
	filter := mapping.ListFilter{}
	if s := c.Query("point_id"); s != "" {
		filter.PointID = &s
	}
	if s := c.Query("tag_id"); s != "" {
		filter.TagID = &s
	}

	mappings, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": mappings})
}

func (h *MappingHandler) Get(c *gin.Context) {
	id := c.Param("id")
	m, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": "Mapping not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": m})
}

func (h *MappingHandler) Create(c *gin.Context) {
	var req mapping.CreateMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	m, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, mapping.ErrValidation) {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	if h.runtimeRefresh != nil {
		if err := h.runtimeRefresh.RefreshMappings(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": m})
}

func (h *MappingHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req mapping.UpdateMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	m, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, mapping.ErrValidation) {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	if h.runtimeRefresh != nil {
		if err := h.runtimeRefresh.RefreshMappings(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": m})
}

func (h *MappingHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	if h.runtimeRefresh != nil {
		if err := h.runtimeRefresh.RefreshMappings(c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

type MappingPreviewRequest struct {
	RawValue          interface{}             `json:"raw_value"`
	TransformPipeline *[]schema.TransformStep `json:"transform_pipeline"`
}

type MappingPreviewResponse struct {
	RawValue    interface{}          `json:"raw_value"`
	FinalValue  interface{}          `json:"final_value"`
	StepResults []mapping.StepResult `json:"step_results"`
	Error       string               `json:"error,omitempty"`
}

func (h *MappingHandler) Preview(c *gin.Context) {
	var req MappingPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	// 驗證必填欄位（transform_pipeline 欄位必須存在）
	if req.TransformPipeline == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": "transform_pipeline field is required"}})
		return
	}

	// Serialize pipeline to JSON string for ExecutePipeline
	pipelineJSON, err := json.Marshal(*req.TransformPipeline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": "Invalid pipeline format"}})
		return
	}

	// Execute pipeline
	ctx, err := mapping.ExecutePipeline(req.RawValue, string(pipelineJSON))

	res := MappingPreviewResponse{
		RawValue:    req.RawValue,
		FinalValue:  ctx.CurrentValue,
		StepResults: ctx.StepResults,
	}
	if err != nil {
		res.Error = err.Error()
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

// ValidatePipelineRequest 驗證管線請求
type ValidatePipelineRequest struct {
	Pipeline *[]schema.TransformStep `json:"pipeline"`
}

// ValidatePipelineResponse 驗證管線回應
type ValidatePipelineResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// ValidatePipeline 驗證轉換管線
// POST /datalink/mappings/validate-pipeline
func (h *MappingHandler) ValidatePipeline(c *gin.Context) {
	var req ValidatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	// 驗證必填欄位（pipeline 欄位必須存在，即使是空陣列）
	if req.Pipeline == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": "pipeline field is required"}})
		return
	}

	resp := ValidatePipelineResponse{Valid: true, Error: ""}

	// 驗證管線格式
	if len(*req.Pipeline) == 0 {
		// 空管線是有效的（pass-through）
		c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
		return
	}

	// 檢查每個步驟的類型是否有效
	validTypes := map[schema.TransformType]bool{
		schema.TransformDecode:      true,
		schema.TransformCast:        true,
		schema.TransformScale:       true,
		schema.TransformLookup:      true,
		schema.TransformConditional: true,
		schema.TransformFormula:     true,
	}

	for i, step := range *req.Pipeline {
		if !validTypes[step.Type] {
			resp.Valid = false
			resp.Error = fmt.Sprintf("步驟 %d: 不支援的轉換類型 '%s'", i+1, step.Type)
			break
		}

		// 驗證 Scale 類型必須有 multiplier 參數
		if step.Type == schema.TransformScale {
			if step.Params == nil {
				resp.Valid = false
				resp.Error = fmt.Sprintf("步驟 %d: scale 類型缺少 params", i+1)
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}
