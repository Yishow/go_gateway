package grouppipeline

import (
	"errors"

	"go-gateway/internal/datalink/groupdelivery"
)

// HandoffRefusal identifies a failed source handoff without exposing raw inputs
// or storage errors. It is a diagnostic, never owner or retry authorization.
// Operator presentation and pending counts belong to the runtime-view provider.
type HandoffRefusal struct {
	scope  groupdelivery.GroupKey
	reason string
	cause  error
}

// Error returns only a closed, safe reason code. Inspect the original error
// chain internally through errors.Is/As, not by serializing its raw causes.
func (e *HandoffRefusal) Error() string { return e.reason }

// Unwrap preserves the underlying typed refusal and internal failure evidence.
func (e *HandoffRefusal) Unwrap() error { return e.cause }

// HandoffScope identifies the known source revision; it does not guess a target
// when ownership evidence is unreadable or conflicting.
func (e *HandoffRefusal) HandoffScope() groupdelivery.GroupKey { return e.scope }

// RepairAction requests read-only inspection of original immutable evidence.
// Preserve accepted inputs; never relabel them using current definitions, reset
// baselines, or resend finalized output. Unrecoverable evidence needs maintenance.
func (e *HandoffRefusal) RepairAction() string { return "inspect-immutable-evidence" }

func scopedHandoffRefusal(key groupdelivery.GroupKey, err error) error {
	if err == nil || !key.Valid() {
		return err
	}
	// Do not classify cancellation, quota, or arbitrary errors as proof faults.
	for _, reason := range []*groupdelivery.OwnershipError{
		groupdelivery.ErrHandoffProofConflict,
		groupdelivery.ErrOwnershipConflict,
		groupdelivery.ErrOwnershipUnavailable,
	} {
		if errors.Is(err, reason) {
			return &HandoffRefusal{scope: key, reason: reason.Error(), cause: err}
		}
	}
	return err
}
