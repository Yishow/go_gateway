package runtime

import (
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"
	"time"
)

// BoundaryReplayState is local cache metadata, never durable ownership authority.
// Until retains an existing stop, including a disabled/draining boundary's stop.
type BoundaryReplayState struct {
	NextClose     time.Time
	Mixed         bool
	Dispatch      snapshot.DispatchWatermark
	NotBefore     time.Time
	Until         time.Time
	HandoffView   groupdelivery.HandoffView
	OwnershipView groupdelivery.OwnershipView
}

// ReplayState captures the cache cursor, stop and verified view together.
func (b *GroupBoundary) ReplayState() BoundaryReplayState {
	b.mu.Lock()
	defer b.mu.Unlock()
	progress, mixed := b.assembler.MixedProgress()
	return BoundaryReplayState{Mixed: mixed, Dispatch: progress, NextClose: b.assembler.NextClose(), NotBefore: b.notBefore, Until: b.until, HandoffView: b.handoffView, OwnershipView: b.ownershipView}
}
