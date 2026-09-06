package api

import (
	"go-gateway/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// registerStudioV2RecordingRoutes 註冊 Studio V2 遙測語意、記錄方案與歷史報表路由。
func registerStudioV2RecordingRoutes(datalinkGroup *gin.RouterGroup, datalinkServices *DatalinkServices) {
	if datalinkServices == nil || datalinkServices.Workspace == nil {
		return
	}

	// Measurements
	if datalinkServices.Measurement != nil {
		measurementHandler := handlers.NewStudioV2WorkspaceMeasurementsHandler(datalinkServices.Workspace, datalinkServices.Device, datalinkServices.Measurement)
		datalinkGroup.GET("/studio-v2/workspace/measurements", measurementHandler.List)
		datalinkGroup.GET("/studio-v2/workspace/measurements/templates", measurementHandler.ListTemplates)
		datalinkGroup.POST("/studio-v2/workspace/measurements/templates/preview", measurementHandler.PreviewTemplate)
		datalinkGroup.POST("/studio-v2/workspace/measurements/templates/apply", measurementHandler.ApplyTemplate)
		datalinkGroup.POST("/studio-v2/workspace/measurements", measurementHandler.Create)
		datalinkGroup.PUT("/studio-v2/workspace/measurements/:id", measurementHandler.Update)
		datalinkGroup.DELETE("/studio-v2/workspace/measurements/:id", measurementHandler.Delete)
	}

	// Recording Plans
	if datalinkServices.RecordingPlan != nil {
		planHandler := handlers.NewStudioV2WorkspaceRecordingPlansHandler(datalinkServices.Workspace, datalinkServices.RecordingPlan, datalinkServices.DBTarget)
		datalinkGroup.GET("/studio-v2/workspace/recording-plans", planHandler.List)
		datalinkGroup.POST("/studio-v2/workspace/recording-plans", planHandler.Create)
		datalinkGroup.GET("/studio-v2/workspace/recording-plans/:id", planHandler.Get)
		datalinkGroup.PUT("/studio-v2/workspace/recording-plans/:id", planHandler.Update)
		datalinkGroup.DELETE("/studio-v2/workspace/recording-plans/:id", planHandler.Delete)
		datalinkGroup.GET("/studio-v2/workspace/recording-plans/capabilities", planHandler.Capabilities)
		datalinkGroup.POST("/studio-v2/workspace/recording-plans/schema-preview", planHandler.SchemaPreview)
		datalinkGroup.POST("/studio-v2/workspace/recording-plans/schema-apply", planHandler.SchemaApply)
		datalinkGroup.POST("/studio-v2/workspace/recording-plans/test-write", planHandler.TestWrite)
	}

	// History Reports & Export
	if datalinkServices.History != nil {
		historyHandler := handlers.NewStudioV2WorkspaceHistoryHandler(datalinkServices.Workspace, datalinkServices.History)
		datalinkGroup.POST("/studio-v2/workspace/history/query", historyHandler.Query)
		datalinkGroup.POST("/studio-v2/workspace/history/export", historyHandler.Export)
	}
}
