package groupdelivery

// OwnershipError is a retryable not-accepted result. Its public text is a safe
// scoped code, never a database error, destination, or observed sample value.
type OwnershipError struct{ Reason string }

func (e *OwnershipError) Error() string { return e.Reason }

// Is supports matching equivalent scoped refusals without exposing raw causes.
func (e *OwnershipError) Is(target error) bool {
	other, ok := target.(*OwnershipError)
	return ok && e.Reason == other.Reason
}

var (
	// ErrOwnershipUnavailable means cutover ownership cannot be proven.
	ErrOwnershipUnavailable = &OwnershipError{Reason: "recording-ownership-unavailable"}
	// ErrOwnershipConflict means the durable cutover assigns this observation to another owner.
	ErrOwnershipConflict = &OwnershipError{Reason: "recording-ownership-conflict"}
)
