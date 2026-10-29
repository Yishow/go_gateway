package handlers

import "go-gateway/internal/datalink/workspace"

// Runtime readiness is a detached read model. Shared setup readiness may retain
// authoring diagnostics; none of its free-form messages enter the monitor.
type runtimeReadinessSummary struct {
	Ready         bool                    `json:"ready"`
	BlockingCount int                     `json:"blocking_count"`
	WarningCount  int                     `json:"warning_count"`
	Issues        []runtimeReadinessIssue `json:"issues"`
}
type runtimeReadinessIssue struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Step     string `json:"step"`
	Scope    string `json:"scope"`
	Message  string `json:"message"`
}

func projectRuntimeReadiness(summary *workspace.ReadinessSummary) *runtimeReadinessSummary {
	if summary == nil {
		return nil
	}
	result := &runtimeReadinessSummary{Ready: summary.Ready, BlockingCount: summary.BlockingCount, WarningCount: summary.WarningCount, Issues: make([]runtimeReadinessIssue, 0, len(summary.Issues))}
	for _, issue := range summary.Issues {
		code := runtimeSafeID(issue.Code)
		if code == "" {
			code = "unavailable"
			result.Ready = false
		}
		severity := runtimeSafeEnum(string(issue.Severity), "blocking", "warning")
		if severity == "" {
			severity = "blocking"
			result.Ready = false
		}
		result.Issues = append(result.Issues, runtimeReadinessIssue{
			Code: code, Severity: severity, Step: runtimeSafeEnum(string(issue.Step), "Step 1", "Step 2", "Step 3", "Step 4"), Scope: runtimeSafeID(issue.Scope),
			Message: "Review saved configuration before activation",
		})
	}
	return result
}
