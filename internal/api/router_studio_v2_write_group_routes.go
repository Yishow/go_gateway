package api

import (
	"go-gateway/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func registerStudioV2WriteGroupRoutes(group *gin.RouterGroup, services *DatalinkServices) {
	registerStudioV2RecordingStartRoutes(group, services)
	handler := handlers.NewStudioV2WorkspaceWriteGroupsHandler(services.WriteGroups).WithDelivery(services.WriteGroupDelivery)
	group.GET("/studio-v2/workspace/write-groups", handler.List)
	group.POST("/studio-v2/workspace/write-groups", handler.Create)
	group.POST("/studio-v2/workspace/write-groups/migrations/single-mappings/preview", handler.PreviewSingleMappingMigration)
	group.POST("/studio-v2/workspace/write-groups/migrations/single-mappings/review", handler.ReviewSingleMappingMigration)
	group.POST("/studio-v2/workspace/write-groups/migrations/row-groups/preview", handler.PreviewRowGroupMigration)
	group.POST("/studio-v2/workspace/write-groups/migrations/row-groups/review", handler.ReviewRowGroupMigration)
	group.POST("/studio-v2/workspace/write-groups/migrations/recording-plans/preview", handler.PreviewRecordingPlanMigration)
	group.POST("/studio-v2/workspace/write-groups/migrations/recording-plans/review", handler.ReviewRecordingPlanMigration)
	group.GET("/studio-v2/workspace/write-groups/:id", handler.Get)
	group.GET("/studio-v2/workspace/write-groups/:id/readiness", handler.Readiness)
	group.GET("/studio-v2/workspace/write-groups/:id/delivery", handler.Delivery)
	group.POST("/studio-v2/workspace/write-groups/:id/delivery/resolve", handler.ResolveDelivery)
	group.POST("/studio-v2/workspace/write-groups/:id/disable", handler.Disable)
	group.POST("/studio-v2/workspace/write-groups/:id/apply", handler.Apply)
	group.PUT("/studio-v2/workspace/write-groups/:id", handler.Update)
	group.DELETE("/studio-v2/workspace/write-groups/:id", handler.Delete)
	schemaHandler := handlers.NewStudioV2WorkspaceWriteGroupSchemaHandler(services.Workspace, services.WriteGroups, services.RecordingPlan, services.DBTarget)
	group.POST("/studio-v2/workspace/write-groups/:id/schema-preview", schemaHandler.Preview)
	group.POST("/studio-v2/workspace/write-groups/:id/schema-apply", schemaHandler.Confirm)

	var plans handlers.RecordingPlanGroupResolver
	if services.WriteGroups != nil {
		plans = services.WriteGroups
	}
	testWrite := handlers.NewStudioV2WorkspaceWriteGroupTestWriteHandler(services.Workspace, services.WriteGroupTestWrite, plans)
	group.POST("/studio-v2/workspace/write-groups/:id/test-write-preview", testWrite.Preview)
	group.POST("/studio-v2/workspace/write-groups/:id/test-write", testWrite.Confirm)
	// Legacy plan routes resolve to the same group and service; they keep no
	// token repository or operation ledger of their own.
	group.POST("/studio-v2/workspace/recording-plans/test-write-preview", testWrite.LegacyPreview)
	group.POST("/studio-v2/workspace/recording-plans/test-write", testWrite.LegacyConfirm)
}
