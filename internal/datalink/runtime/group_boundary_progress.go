package runtime

import "go-gateway/internal/datalink/snapshot"

// MixedProgress is a detached typed cursor projection, never a shared window.
func (b *GroupBoundary) MixedProgress() (snapshot.DispatchWatermark, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.assembler.MixedProgress()
}

// NeedsDurableRestore reports an ambiguous durable acknowledgment. The old
// boundary remains fail-closed; only a newly constructed restored one may run.
func (b *GroupBoundary) NeedsDurableRestore() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.uncertainCommit != nil
}
