package grouppipeline

import (
	"context"
	"go-gateway/internal/datalink/groupdelivery"
)

// ResolveAttention uses the same durable store as intake, workers and delivery reads.
func (p *Pipeline) ResolveAttention(ctx context.Context, workspaceID, groupID string, decision groupdelivery.OperatorDecision) (*groupdelivery.DecisionResult, error) {
	return p.deps.Store.ResolveAttention(ctx, workspaceID, groupID, decision)
}
