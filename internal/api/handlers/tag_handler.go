package handlers

import (
	"context"
	"net/http"

	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	svc *tag.Service
}

func NewTagHandler() *TagHandler {
	repo := tag.NewMemoryRepository()
	svc := tag.NewService(repo)

	// Seed some tags
	list, _ := svc.List(context.Background(), tag.ListFilter{})
	if len(list) == 0 {
		svc.Create(context.Background(), tag.CreateTagRequest{
			Key:         "temp_c",
			DisplayName: "Temperature (Celsius)",
			DataType:    "float64",
			Unit:        "°C",
		})
		svc.Create(context.Background(), tag.CreateTagRequest{
			Key:         "pressure_bar",
			DisplayName: "Pressure (Bar)",
			DataType:    "float64",
			Unit:        "bar",
		})
	}

	return &TagHandler{svc: svc}
}

func (h *TagHandler) List(c *gin.Context) {
	filter := tag.ListFilter{}
	if s := c.Query("status"); s != "" {
		// filter.Status =
	}
	if s := c.Query("key_prefix"); s != "" {
		filter.KeyPrefix = s
	}

	tags, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tags})
}

func (h *TagHandler) Get(c *gin.Context) {
	id := c.Param("id")
	t, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": gin.H{"message": "Tag not found"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": t})
}

func (h *TagHandler) Create(c *gin.Context) {
	var req tag.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	t, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": t})
}

func (h *TagHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req tag.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}

	t, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": t})
}

func (h *TagHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
