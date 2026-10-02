package workspace

import (
	"context"
	"errors"
)

// ErrWriteGroupPlanMigrationBlocked identifies a recording plan that cannot
// be represented by the current canonical basic WriteGroup contract.
var ErrWriteGroupPlanMigrationBlocked = errors.New("recording plan migration is blocked")

// PreviewRecordingPlanMigration previews a recording-plan migration without
// creating a canonical group or changing the persisted plan.
func (s *WriteGroupService) PreviewRecordingPlanMigration(
	ctx context.Context,
	workspaceID string,
	sourceIDs []string,
) (*WriteGroupMigrationPreview, error) {
	return s.previewRecordingPlanMigration(ctx, workspaceID, sourceIDs)
}

// ReviewRecordingPlanMigration validates a reviewed recording-plan migration
// envelope and returns the explicit blocked result without persisting changes.
func (s *WriteGroupService) ReviewRecordingPlanMigration(
	ctx context.Context,
	request WriteGroupMigrationReviewRequest,
) (*WriteGroupMigrationReviewResult, error) {
	return s.reviewRecordingPlanMigration(ctx, request)
}
