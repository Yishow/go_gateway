package groupdelivery

// OwnershipError is a typed not-accepted result; its reason determines recovery.
// Its public text is a safe scoped code, never a database error, destination,
// or observed sample value.
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
	// ErrIntakeStopped means canonical stop/remove forbids new acceptance.
	ErrIntakeStopped = &OwnershipError{Reason: "recording-intake-stopped"}
)
