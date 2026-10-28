package workspace

import (
	"context"
	"strings"
)

const successorActivationIntentVersion = "go-gateway/successor-activation-intent/v1"

// isSuccessorActivationIntent is a denial classifier, not proof of authority.
// Unknown or malformed successor-origin events must not enter legacy rollback.
func isSuccessorActivationIntent(intent string) bool {
	return strings.Contains(intent, "go-gateway/successor-activation-intent/")
}

// requireCurrentCanonicalOwnerInTx prevents a historical canonical generation
// from taking back ownership through ordinary mutation paths. Unkeyed legacy
// groups keep their existing contract; historical reads/delivery do not call it.
func (r *SQLWriteGroupRepository) requireCurrentCanonicalOwnerInTx(ctx context.Context, runner writeGroupSQLRunner, workspaceID, groupID string) error {
	meta, found, err := r.loadWriteGroupLineage(ctx, runner, workspaceID, groupID)
	if err != nil {
		return err
	}
	if found && meta.CurrentGroupID != groupID {
		return &WriteGroupLifecycleActionError{Action: "apply", ReasonCode: "historical_generation", NextAction: "review_lifecycle", cause: ErrWriteGroupLifecycleBlocked}
	}
	return nil
}
