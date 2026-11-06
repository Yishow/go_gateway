package api

import (
	"go-gateway/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// registerAcquisitionAnomalyRoutes exposes the anomaly history and per-device
// summary. The routes are always registered; without a repository or device
// service they answer that the capability is unavailable instead of reporting
// zero anomalies.
func registerAcquisitionAnomalyRoutes(group *gin.RouterGroup, services *DatalinkServices) {
	var devices handlers.AcquisitionAnomalyDeviceReader
	if services.Device != nil {
		devices = services.Device
	}
	var points handlers.AcquisitionAnomalyPointReader
	if services.Point != nil {
		points = services.Point
	}
	h := handlers.NewAcquisitionAnomalyHandler(services.AcquisitionAnomalies, devices, points)
	group.GET("/devices/:id/acquisition-anomalies", h.ListDeviceAnomalies)
	group.GET("/acquisition-anomalies/summary", h.Summary)
}
