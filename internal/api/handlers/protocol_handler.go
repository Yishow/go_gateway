package handlers

import (
	"net/http"

	"go-gateway/internal/datalink/connector"

	"github.com/gin-gonic/gin"
)

// ProtocolHandler 協議資訊 API Handler
type ProtocolHandler struct{}

// NewProtocolHandler 建立新的協議 Handler
func NewProtocolHandler() *ProtocolHandler {
	return &ProtocolHandler{}
}

// ProtocolInfo 協議資訊
type ProtocolInfo struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ConfigSchema string `json:"config_schema"`
}

// List 列出所有支援的協議
// GET /datalink/protocols
func (h *ProtocolHandler) List(c *gin.Context) {
	sourceInfos := connector.ListProtocolInfos()
	protocols := make([]ProtocolInfo, 0, len(sourceInfos))
	for _, info := range sourceInfos {
		protocols = append(protocols, ProtocolInfo{
			Type:         string(info.Type),
			Name:         info.Name,
			Description:  info.Description,
			ConfigSchema: string(info.ConfigSchema),
		})
	}

	c.JSON(http.StatusOK, gin.H{apiResponseSuccessKey: true, apiResponseDataKey: protocols})
}
