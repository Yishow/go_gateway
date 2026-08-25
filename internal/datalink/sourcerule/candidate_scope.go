package sourcerule

import (
	"context"
	"strings"

	"go-gateway/internal/datalink/modbusshare"
)

// CandidateScopeRequest is the common scope and revision contract for
// candidate review/apply operations.
type CandidateScopeRequest struct {
	WorkspaceID               string `json:"workspace_id"`
	ExpectedWorkspaceRevision string `json:"expected_workspace_revision"`
	RevisionID                string `json:"revision_id"`
}

// SetCandidateScope wires the persisted workspace and Share hydration seams.
// A nil pair keeps legacy unit fixtures usable; production routes always wire
// both readers from DatalinkServices.
func (s *Service) SetCandidateScope(workspace CandidateWorkspaceReader, hydration CandidateHydrationReader) {
	s.candidateScope = candidateScope{workspace: workspace, hydration: hydration}
}

// CandidateScopeConfigured reports whether candidate endpoints are protected
// by the production workspace/hydration barrier.
func (s *Service) CandidateScopeConfigured() bool {
	return s != nil && s.candidateScope.workspace != nil && s.candidateScope.hydration != nil
}

// ValidateCandidateScope validates hydration, workspace/device membership and
// both the expected workspace and persisted candidate revisions before any
// candidate snapshot or tag/link/mapping mutation is attempted.
func (s *Service) ValidateCandidateScope(ctx context.Context, ruleID string, req CandidateScopeRequest) error {
	if !s.CandidateScopeConfigured() {
		return nil
	}
	if strings.TrimSpace(req.WorkspaceID) == "" || strings.TrimSpace(req.ExpectedWorkspaceRevision) == "" || strings.TrimSpace(req.RevisionID) == "" {
		return modbusshare.NewError(modbusshare.ErrCodeHydrationRequired, "candidate workspace hydration scope is required", true)
	}

	hydration, err := s.candidateScope.hydration.CheckHydration(ctx)
	if err != nil || hydration.State != modbusshare.HydrationStateReady || !hydration.Readiness {
		return &modbusshare.Error{
			Code:              modbusshare.ErrCodeHydrationRequired,
			Message:           "workspace hydration is required before candidate review or apply",
			Retryable:         true,
			WorkspaceRevision: hydration.WorkspaceRevision,
			SettingsRevision:  hydration.SettingsRevision,
			Action:            "complete workspace bootstrap and retry",
		}
	}

	workspaceID, deviceIDs, err := s.candidateScope.workspace.ShareOwnershipSnapshot(ctx)
	if err != nil || workspaceID != strings.TrimSpace(req.WorkspaceID) {
		return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "candidate is outside the requested workspace", false)
	}
	if hydration.WorkspaceRevision != strings.TrimSpace(req.ExpectedWorkspaceRevision) {
		return &modbusshare.Error{
			Code:              modbusshare.ErrCodeRevisionConflict,
			Message:           "workspace revision conflict",
			Retryable:         true,
			WorkspaceRevision: hydration.WorkspaceRevision,
			SettingsRevision:  hydration.SettingsRevision,
			Action:            candidateRetryAction,
		}
	}

	rule, err := s.repo.GetByID(ctx, ruleID)
	if err != nil {
		return err
	}
	ownedDevice := false
	for _, deviceID := range deviceIDs {
		if strings.TrimSpace(deviceID) == strings.TrimSpace(rule.DeviceID) {
			ownedDevice = true
			break
		}
	}
	if !ownedDevice {
		return modbusshare.NewError(modbusshare.ErrCodeWorkspaceScope, "source rule device is outside the requested workspace", false)
	}
	if rule.RevisionID != strings.TrimSpace(req.RevisionID) {
		return &modbusshare.Error{
			Code:              modbusshare.ErrCodeRevisionConflict,
			Message:           candidateRevisionConflictMessage,
			Retryable:         true,
			WorkspaceRevision: hydration.WorkspaceRevision,
			SettingsRevision:  hydration.SettingsRevision,
			Action:            candidateRetryAction,
		}
	}
	return nil
}
