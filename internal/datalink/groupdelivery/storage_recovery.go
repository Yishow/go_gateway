package groupdelivery

import (
	"context"
)

// RestoreAction is what delivery does after asking storage recovery about an
// item whose destination table is not the shape the insert expects.
type RestoreAction int

const (
	// RestoreNotApplicable: storage recovery has nothing to say; the failure is
	// classified as it always was.
	RestoreNotApplicable RestoreAction = iota
	// RestoreRetryNow: the table was found missing and has been restored from
	// its confirmed layout; the insert is tried once more.
	RestoreRetryNow
	// RestoreWait: another execution is restoring the table; the item retries
	// later.
	RestoreWait
	// RestoreBlock: restoring is refused; the item is blocked with Code until a
	// person repairs the destination or confirms the schema again.
	RestoreBlock
)

// RestoreDecision is the answer of storage recovery.
type RestoreDecision struct {
	Action RestoreAction
	// Code is the stable safe code recorded on a blocked item.
	Code string
}

// StorageRecovery lets delivery restore a managed table that was confirmed
// earlier and has since vanished, and nothing else. It is optional: without
// it a destination whose table is gone blocks delivery as it always did.
type StorageRecovery interface {
	// Admit is asked before every attempt. A non-empty code blocks the item
	// without any access to the destination.
	Admit(ctx context.Context, item OutboxItem) (code string, err error)
	// Restore is asked after the destination refused an insert with a
	// schema-shaped error. It decides from the destination catalog and the
	// local evidence frozen at confirmation, never from the error itself.
	Restore(ctx context.Context, item OutboxItem, target Target, cause error) (RestoreDecision, error)
	// BlockedCode names the stable code of an item whose destination could not
	// be opened because its confirmed identity changed; empty keeps the
	// delivery layer's own code.
	BlockedCode(ctx context.Context, item OutboxItem) string
}

// settlement records, in the local database and within its own budget, what an
// attempt came to.
type settlement func(context.Context) (DeliveryResult, error)

// recoverStorage asks storage recovery about an insert the destination refused
// with a schema-shaped error. tryAgain is true when the table was restored and
// the insert should be repeated; otherwise record is what to settle with, or nil
// to settle the refusal like any other. The restore runs inside the attempt
// context; settlement is the caller's, after all destination work.
func (s *Sender) recoverStorage(ctx context.Context, item OutboxItem, target Target, claim Claim, cause error) (tryAgain bool, record settlement) {
	decision, err := s.config.Recovery.Restore(ctx, item, target, cause)
	switch {
	case err != nil:
		return false, func(c context.Context) (DeliveryResult, error) {
			return s.retry(c, claim, codeStorageRecoveryUnavailable)
		}
	case decision.Action == RestoreRetryNow:
		return true, nil
	case decision.Action == RestoreWait:
		return false, func(c context.Context) (DeliveryResult, error) { return s.retry(c, claim, codeStorageRestoring) }
	case decision.Action == RestoreBlock:
		return false, func(c context.Context) (DeliveryResult, error) { return s.blockClaimed(c, claim, decision.Code) }
	}
	return false, nil
}
