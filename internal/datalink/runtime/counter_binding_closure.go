package runtime

import (
	"maps"

	"go-gateway/internal/datalink/counterproof"
)

// counterClosureProofs carries the activation-time references into the same
// transaction that will persist the closure and consume accepted inputs.
func (b *GroupBoundary) counterClosureProofs() []counterproof.CounterBindingReference {
	var refs []counterproof.CounterBindingReference
	for member := range maps.Values(b.members) {
		if !member.counterProof.IsZero() {
			refs = append(refs, member.counterProof)
		}
	}
	return refs
}
