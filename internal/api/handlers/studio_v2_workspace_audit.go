package handlers

import (
	"context"
	"log"
	"strings"

	"go-gateway/internal/datalink/audit"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

type workspaceAuditRecorder interface {
	Record(ctx context.Context, event audit.RecordEvent) error
}

// WithAudit attaches workspace audit recording to activation responses.
func (h *StudioV2WorkspaceActivationHandler) WithAudit(auditSvc *audit.Service) *StudioV2WorkspaceActivationHandler {
	if h == nil {
		return nil
	}
	h.auditSvc = auditSvc
	return h
}

// WithAudit attaches workspace audit recording to persisted database changes.
func (h *StudioV2WorkspaceDatabaseHandler) WithAudit(auditSvc *audit.Service) *StudioV2WorkspaceDatabaseHandler {
	if h == nil {
		return nil
	}
	h.auditSvc = auditSvc
	return h
}

func (h *StudioV2WorkspaceActivationHandler) recordActivationAudit(ctx context.Context, response *workspace.ActivationResponse) {
	if h == nil || h.auditSvc == nil || response == nil || strings.TrimSpace(response.WorkspaceID) == "" {
		return
	}

	results := append([]workspace.ActivationResult(nil), response.Results...)
	if err := h.auditSvc.Record(ctx, audit.RecordEvent{
		WorkspaceID: response.WorkspaceID,
		EventType:   audit.EventTypeWorkspaceActivation,
		Result:      activationAuditResult(results),
		Scope:       activationAuditScope(results),
		Details: map[string]any{
			"results":             results,
			apiResponseMessageKey: response.Message,
		},
	}); err != nil {
		log.Printf("record workspace activation audit failed: %v", err)
	}
}

func activationAuditResult(results []workspace.ActivationResult) audit.Result {
	successCount := 0
	failureCount := 0
	for _, result := range results {
		switch result.Status {
		case workspace.ActivationResultStatusSuccess:
			successCount++
		case workspace.ActivationResultStatusFailed:
			failureCount++
		}
	}
	if failureCount > 0 && successCount > 0 {
		return audit.ResultPartialSuccess
	}
	if failureCount > 0 {
		return audit.ResultFailure
	}
	return audit.ResultSuccess
}

func activationAuditScope(results []workspace.ActivationResult) string {
	deviceIDs := make([]string, 0, len(results))
	for _, result := range results {
		if deviceID := strings.TrimSpace(result.DeviceID); deviceID != "" {
			deviceIDs = append(deviceIDs, deviceID)
		}
	}
	if len(deviceIDs) == 0 {
		return "workspace"
	}
	return "devices:" + strings.Join(deviceIDs, ",")
}

func (h *StudioV2WorkspaceDatabaseHandler) recordDatabaseConfigAudit(
	ctx context.Context,
	workspaceID string,
	connector *schema.DatabaseConnector,
	runtimeApplyStatus string,
) {
	if h == nil || h.auditSvc == nil || connector == nil || strings.TrimSpace(workspaceID) == "" {
		return
	}

	if err := h.auditSvc.Record(ctx, audit.RecordEvent{
		WorkspaceID: workspaceID,
		EventType:   audit.EventTypeDatabaseConfigSaved,
		Result:      audit.ResultSuccess,
		Scope:       "database_connector:" + connector.ID,
		ReferenceID: connector.ID,
		Details: map[string]any{
			"kind":                 connector.Kind,
			"table":                connectorConfigString(connector, "table"),
			"write_mode":           connectorConfigWriteMode(connector),
			"runtime_apply_status": runtimeApplyStatus,
		},
	}); err != nil {
		log.Printf("record workspace database config audit failed: %v", err)
	}
}

func (h *StudioV2WorkspaceDatabaseHandler) recordDatabaseTargetAudit(
	ctx context.Context,
	workspaceID string,
	mapping *schema.DatabaseTargetMapping,
	pointID string,
	runtimeApplyStatus string,
) {
	if h == nil || h.auditSvc == nil || mapping == nil || strings.TrimSpace(workspaceID) == "" {
		return
	}

	if err := h.auditSvc.Record(ctx, audit.RecordEvent{
		WorkspaceID: workspaceID,
		EventType:   audit.EventTypeDatabaseTargetSaved,
		Result:      audit.ResultSuccess,
		Scope:       "database_target:" + strings.TrimSpace(pointID),
		ReferenceID: mapping.ID,
		Details: map[string]any{
			"mapping_id":            mapping.ID,
			"point_id":              strings.TrimSpace(pointID),
			"tag_id":                mapping.TagID,
			"column_name":           mapping.ColumnName,
			"enabled":               mapping.Enabled,
			"runtime_apply_status":  runtimeApplyStatus,
			"database_connector_id": mapping.ConnectorID,
		},
	}); err != nil {
		log.Printf("record workspace database target audit failed: %v", err)
	}
}
