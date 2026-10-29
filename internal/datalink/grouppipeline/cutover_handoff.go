package grouppipeline

import (
	"context"
	"errors"

	"go-gateway/internal/datalink/groupdelivery"
)

const handoffPageSize = 128

// PrepareCutoverHandoffs is an explicit, workspace-scoped migration adapter.
// The explicit adapter only prepares durable state. Reconcile additionally
// rebuilds participating workers; successful preparation alone never means
// external delivery or worker readiness.
func (p *Pipeline) PrepareCutoverHandoffs(ctx context.Context, workspaceID string) error {
	_, err := p.prepareWorkspaceHandoffs(ctx, workspaceID)
	return err
}

func (p *Pipeline) prepareWorkspaceHandoffs(ctx context.Context, workspaceID string) ([][]groupdelivery.GroupKey, error) {
	if p.deps.Store == nil || workspaceID == "" {
		return nil, groupdelivery.ErrOwnershipUnavailable
	}
	var after groupdelivery.GroupKey
	var failures []error
	var groups [][]groupdelivery.GroupKey
	for {
		if err := ctx.Err(); err != nil {
			return groups, errors.Join(append(failures, err)...)
		}
		keys, err := p.deps.Store.HandoffSources(ctx, workspaceID, after, handoffPageSize)
		if err != nil {
			return groups, errors.Join(append(failures, err)...)
		}
		for _, key := range keys {
			participants, err := p.prepareHandoffSource(ctx, key)
			if len(participants) > 0 {
				groups = append(groups, participants)
			}
			if err != nil {
				failures = append(failures, scopedHandoffRefusal(key, err))
			}
			after = key
		}
		if len(keys) < handoffPageSize {
			return groups, errors.Join(failures...)
		}
	}
}

func (p *Pipeline) prepareHandoffSource(ctx context.Context, key groupdelivery.GroupKey) ([]groupdelivery.GroupKey, error) {
	var participants []groupdelivery.GroupKey
	var after int64
	var firstFailure error
	for {
		if err := ctx.Err(); err != nil {
			return participants, errors.Join(firstFailure, err)
		}
		batch, err := p.deps.Store.PrepareCutoverHandoffs(ctx, key, after, handoffPageSize)
		participants = append(participants, batch.ReplayParticipants...)
		if firstFailure == nil {
			firstFailure = err
		}
		if !batch.More {
			return participants, firstFailure
		}
		if batch.NextSeq <= after {
			// A failed transaction must not produce an endless retry loop. The
			// next explicit invocation retries durable identity, never a new ID.
			return participants, errors.Join(firstFailure, groupdelivery.ErrOwnershipUnavailable)
		}
		after = batch.NextSeq
	}
}
