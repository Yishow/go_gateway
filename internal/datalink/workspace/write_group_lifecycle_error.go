package workspace

// WriteGroupLifecycleActionError exposes a safe, operation-specific repair
// route while retaining the original error chain for existing callers.
type WriteGroupLifecycleActionError struct {
	Action     string
	ReasonCode string
	NextAction string
	cause      error
}

// Error deliberately excludes storage errors, connector details and payloads.
func (e *WriteGroupLifecycleActionError) Error() string {
	return "write group lifecycle action is unavailable"
}

// Unwrap preserves existing sentinel and concrete cause matching.
func (e *WriteGroupLifecycleActionError) Unwrap() error { return e.cause }

func writeGroupRemoveActionError(group *WriteGroup, reason string, cause error) error {
	next := writeGroupNextReviewConnector
	if reason == writeGroupReasonOwnershipUnavailable {
		next = writeGroupNextInspectOwnership
	}
	if group.Status != WriteGroupStatusDisabled {
		next = writeGroupActionStop
	}
	return &WriteGroupLifecycleActionError{Action: "remove", ReasonCode: reason, NextAction: next, cause: cause}
}
