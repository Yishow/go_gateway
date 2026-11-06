package api

import (
	"go-gateway/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// registerRecordingHomeSourceRoutes always mounts the readonly observation route.
// Missing dependencies are reported as unavailable rather than an absent route.
func registerRecordingHomeSourceRoutes(group *gin.RouterGroup, services *DatalinkServices) {
	if services == nil {
		services = &DatalinkServices{}
	}
	handler := handlers.NewStudioV2RecordingHomeSourcesHandler(
		services.Workspace, services.Device, services.SourceRule,
		services.Point, services.Mapping, services.Tag,
	)
	group.GET("/studio-v2/recording-home/sources", handler.Get)
}
