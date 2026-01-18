package handlers

import (
	"encoding/json"
	"net/http"

	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
)

type MappingHandler struct {
	svc *mapping.Service
}

func NewMappingHandler() *MappingHandler {
	repo := mapping.NewMemoryRepository()
	svc := mapping.NewService(repo)
	return &MappingHandler{svc: svc}
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": m})
}

func (h *MappingHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

type MappingPreviewRequest struct {
    RawValue interface{} `json:"raw_value"`
    TransformPipeline []schema.TransformStep `json:"transform_pipeline"`
}

type MappingPreviewResponse struct {
    RawValue interface{} `json:"raw_value"`
    FinalValue interface{} `json:"final_value"`
    StepResults []mapping.StepResult `json:"step_results"`
    Error string `json:"error,omitempty"`
}

func (h *MappingHandler) Preview(c *gin.Context) {
	var req MappingPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

    // Serialize pipeline to JSON string for ExecutePipeline
    pipelineJSON, err := json.Marshal(req.TransformPipeline)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": "Invalid pipeline format"}})
        return
    }

	// Execute pipeline
    ctx, err := mapping.ExecutePipeline(req.RawValue, string(pipelineJSON))
    
	res := MappingPreviewResponse{
		RawValue: req.RawValue,
        FinalValue: ctx.CurrentValue,
		StepResults: ctx.StepResults,
	}
    if err != nil {
        res.Error = err.Error()
    }
	
	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}
