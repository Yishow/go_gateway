package handlers

import (
	"context"

	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"
)

type workspaceReadinessReader interface {
	Readiness(ctx context.Context) (*workspace.ReadinessSummary, error)
}

const runtimeApplyFailedStatus = "apply_failed"

type studioV2RuntimeApplyOutcome struct {
	Status  string
	Code    string
	Message string
	Issues  []workspace.ReadinessIssue
}

type studioV2RuntimeApplyResponse struct {
	RuntimeApplyStatus  string                     `json:"runtime_apply_status,omitempty"`
	RuntimeApplyCode    string                     `json:"runtime_apply_code,omitempty"`
	RuntimeApplyMessage string                     `json:"runtime_apply_message,omitempty"`
	RuntimeApplyIssues  []workspace.ReadinessIssue `json:"runtime_apply_issues,omitempty"`
}

func mapStudioV2RuntimeApplyResponse(outcome studioV2RuntimeApplyOutcome) studioV2RuntimeApplyResponse {
	return studioV2RuntimeApplyResponse{
		RuntimeApplyStatus:  outcome.Status,
		RuntimeApplyCode:    outcome.Code,
		RuntimeApplyMessage: outcome.Message,
		RuntimeApplyIssues:  outcome.Issues,
	}
}

func mapSourceRuleRuntimeReconcileOutcome(outcome sourcerule.RuntimeReconcileOutcome) studioV2RuntimeApplyOutcome {
	status := string(outcome.Status)
	if outcome.Status == sourcerule.RuntimeReconcileStatusNotRunning {
		status = "not_running"
	}
	if status == "" {
		status = string(sourcerule.RuntimeReconcileStatusStale)
	}
	return studioV2RuntimeApplyOutcome{
		Status:  status,
		Code:    outcome.Code,
		Message: safeRuntimeApplyMessage(outcome.Message, status),
	}
}

func resolveStudioV2RuntimeApplyStatus(ctx context.Context, deviceSvc *device.Service, deviceID string) (string, string) {
	if deviceSvc == nil || deviceID == "" {
		return "not_running", ""
	}

	record, err := deviceSvc.GetByID(ctx, deviceID)
	if err != nil {
		return runtimeApplyFailedStatus, "runtime apply status is unavailable"
	}
	if record.Status != schema.DeviceStatusActive {
		return "not_running", ""
	}
	return "applied", ""
}

func resolveStudioV2WorkspaceRuntimeApplyStatus(ctx context.Context, deviceSvc *device.Service, deviceIDs []string) (string, string) {
	if len(deviceIDs) == 0 {
		return "not_running", ""
	}

	for _, deviceID := range deviceIDs {
		status, message := resolveStudioV2RuntimeApplyStatus(ctx, deviceSvc, deviceID)
		if status == runtimeApplyFailedStatus {
			return status, message
		}
		if status == "applied" {
			return status, message
		}
	}

	return "not_running", ""
}

func resolveStudioV2ScopedRuntimeApplyOutcome(ctx context.Context, workspaceSvc workspaceReadinessReader, deviceSvc *device.Service, deviceIDs []string, relevantScopes []string) studioV2RuntimeApplyOutcome {
	status, message := resolveStudioV2WorkspaceRuntimeApplyStatus(ctx, deviceSvc, deviceIDs)
	if status != "applied" {
		return studioV2RuntimeApplyOutcome{Status: status, Message: message}
	}
	if workspaceSvc == nil {
		return studioV2RuntimeApplyOutcome{Status: status}
	}

	summary, err := workspaceSvc.Readiness(ctx)
	if err != nil {
		return studioV2RuntimeApplyOutcome{Status: runtimeApplyFailedStatus, Code: "workspace_readiness_unavailable", Message: "workspace readiness is unavailable"}
	}
	issues := filterRuntimeApplyIssues(summary, relevantScopes)
	if len(issues) == 0 {
		return studioV2RuntimeApplyOutcome{Status: status}
	}

	blocking := make([]workspace.ReadinessIssue, 0, len(issues))
	warnings := make([]workspace.ReadinessIssue, 0, len(issues))
	for _, issue := range issues {
		if issue.Severity == workspace.ReadinessSeverityBlocking {
			blocking = append(blocking, issue)
			continue
		}
		warnings = append(warnings, issue)
	}
	if len(blocking) > 0 {
		return studioV2RuntimeApplyOutcome{
			Status:  "deferred",
			Message: firstRuntimeApplyIssueMessage(blocking, "workspace readiness blocked live apply"),
			Issues:  blocking,
		}
	}

	return studioV2RuntimeApplyOutcome{
		Status:  status,
		Message: firstRuntimeApplyIssueMessage(warnings, ""),
		Issues:  warnings,
	}
}

func safeRuntimeApplyMessage(message, status string) string {
	switch status {
	case runtimeApplyFailedStatus, "failed":
		return "runtime apply failed"
	case "stale":
		return "runtime projection is stale"
	}
	return message
}

func filterRuntimeApplyIssues(summary *workspace.ReadinessSummary, relevantScopes []string) []workspace.ReadinessIssue {
	if summary == nil || len(summary.Issues) == 0 {
		return nil
	}
	if len(relevantScopes) == 0 {
		issues := make([]workspace.ReadinessIssue, 0, len(summary.Issues))
		issues = append(issues, summary.Issues...)
		return issues
	}

	scopeSet := make(map[string]struct{}, len(relevantScopes))
	for _, scope := range relevantScopes {
		if scope == "" {
			continue
		}
		scopeSet[scope] = struct{}{}
	}
	issues := make([]workspace.ReadinessIssue, 0, len(summary.Issues))
	for _, issue := range summary.Issues {
		if issue.Scope == "" {
			continue
		}
		if _, ok := scopeSet[issue.Scope]; ok {
			issues = append(issues, issue)
		}
	}
	return issues
}

func firstRuntimeApplyIssueMessage(issues []workspace.ReadinessIssue, fallback string) string {
	for _, issue := range issues {
		if issue.Message != "" {
			return issue.Message
		}
	}
	return fallback
}
