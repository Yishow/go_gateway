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
		wireRecordingPlanMembership(planHandler, datalinkServices)
		datalinkGroup.GET("/studio-v2/workspace/recording-plans", planHandler.List)
		datalinkGroup.POST("/studio-v2/workspace/recording-plans", planHandler.Create)
		datalinkGroup.GET("/studio-v2/workspace/recording-plans/:id", planHandler.Get)
		datalinkGroup.PUT("/studio-v2/workspace/recording-plans/:id", planHandler.Update)
		datalinkGroup.DELETE("/studio-v2/workspace/recording-plans/:id", planHandler.Delete)
		datalinkGroup.GET("/studio-v2/workspace/recording-plans/capabilities", planHandler.Capabilities)
		datalinkGroup.POST("/studio-v2/workspace/recording-plans/schema-preview", planHandler.SchemaPreview)
		datalinkGroup.POST("/studio-v2/workspace/recording-plans/schema-apply", planHandler.SchemaApply)
		datalinkGroup.POST("/studio-v2/workspace/recording-plans/test-write", planHandler.TestWrite)

		operationsHandler := handlers.NewStudioV2WorkspaceDatabaseOperationsHandler(datalinkServices.Workspace, datalinkServices.RecordingPlan)
		datalinkGroup.GET("/studio-v2/workspace/database-operations/:operation_id", operationsHandler.Get)
	}

	// History Reports & Export
	registerStudioV2HistoryRoutes(datalinkGroup, datalinkServices)
}

// wireRecordingPlanMembership passes only configured readers. An absent
// service stays a nil interface, so plan creation fails closed instead of
// calling a nil service pointer.
func wireRecordingPlanMembership(handler *handlers.StudioV2WorkspaceRecordingPlansHandler, services *DatalinkServices) {
	var measurementReader handlers.RecordingPlanMeasurementService
	if services.Measurement != nil {
		measurementReader = services.Measurement
	}
	var pointReader handlers.RecordingPlanPointService
	if services.Point != nil {
		pointReader = services.Point
	}
	handler.SetRecordingMembershipServices(measurementReader, pointReader)
}

func registerStudioV2HistoryRoutes(datalinkGroup *gin.RouterGroup, datalinkServices *DatalinkServices) {
	if datalinkServices.History != nil {
		historyHandler := handlers.NewStudioV2WorkspaceHistoryHandler(datalinkServices.Workspace, datalinkServices.History)
		datalinkGroup.POST("/studio-v2/workspace/history/query", historyHandler.Query)
		datalinkGroup.POST("/studio-v2/workspace/history/export", historyHandler.Export)
	}
}
