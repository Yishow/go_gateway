package grouptestwrite

import (
	"context"
	"strings"

	"go-gateway/internal/datalink/recordingplan"
)

// repairCleanup runs cleanup-only retry for an operation whose write was verified
// but cleanup failed or was unconfirmed. It never inserts rows.
func (s *Service) repairCleanup(ctx context.Context, ws Workspace, req Confirmation) (*Outcome, error) {
	if err := s.checkGroup(ctx, ws, req); err != nil {
		return nil, err
	}
	groupID := strings.TrimSpace(req.GroupID)
	if groupID == "" {
		token, err := s.deps.Ledger.PreviewTokenForWorkspace(ctx, ws.ID, req.Token)
		if err != nil {
			return nil, err
		}
		groupID = token.PlanID
	}
	p, err := s.buildPlan(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if p.sqlserver == nil {
		return nil, unsupported("cleanup-repair-unsupported")
	}

	// 1. Resolve existing operation without mutation
	op, claim, err := s.deps.Ledger.ResolveTestWriteReplay(ctx, ws.ID, req.Token, req.OperationID)
	if err != nil {
		return nil, err
	}
	if op == nil || claim == recordingplan.ClaimNone {
		return nil, unsupported("cleanup-repair-operation-unconfirmed")
	}
	if claim == recordingplan.ClaimCompleted && op.CleanupStatus == CleanupCleaned {
		return &Outcome{Operation: op, Claim: claim}, nil
	}
	if claim == recordingplan.ClaimInProgress {
		return &Outcome{Operation: op, Claim: claim}, nil
	}
	if op.WriteOutcome != WriteVerified {
		return nil, unsupported("cleanup-repair-write-unverified")
	}
	if p.digest != op.PayloadDigest {
		return nil, recordingplan.ErrPreviewTokenStale
	}

	// 2. Validate original progress BEFORE reclaim so op is never stuck on malformed proof
	state := decodeProgress(op.Detail)
	if state.Phase != phaseCleaning || state.WriteOutcome != WriteVerified || state.OwnerValue != ownerValue(op.OperationID) ||
		len(state.RecordIDs) == 0 || len(state.EffectKeys) == 0 || len(state.EffectDigests) == 0 ||
		len(state.RecordIDs) != len(state.EffectKeys) ||
		state.OwnerColumn == "" || state.TableOwnerValue == "" ||
		state.Schema != p.group.Destination.TableSchema || state.Table != p.group.Destination.TableName ||
		state.WorkspaceID != ws.ID || state.GroupID != p.group.ID || state.GroupRevision != p.group.Revision {
		return nil, unsupported(ReasonProofMalformed)
	}

	// 3. Atomically claim execution rights
	claimedOp, claim, err := s.deps.Ledger.ClaimTestWriteCleanup(ctx, ws.ID, req.Token, req.OperationID, s.scopeFor(ws, p), p.digest)
	if err != nil {
		return nil, err
	}
	if claim != recordingplan.ClaimAcquired {
		return &Outcome{Operation: claimedOp, Claim: claim}, nil
	}

	// 4. Open destination (NEVER insert!)
	token, err := s.deps.Ledger.PreviewTokenForWorkspace(ctx, ws.ID, req.Token)
	if err != nil {
		return nil, err
	}
	dest, err := s.openDestination(ctx, token)
	if err != nil {
		return s.finishUnreachable(ctx, claimedOp, state, recordingplan.SchemaOperationResult{
			PayloadDigest: claimedOp.PayloadDigest,
			WriteOutcome:  claimedOp.WriteOutcome,
		}, err)
	}
	if dest.Close != nil {
		defer func() { _ = dest.Close() }() //nolint:errcheck // releasing a finished connection; nothing to recover
	}

	// 5. Run cleanupPhase only
	cleanup := s.cleanupPhase(ctx, claimedOp, p, dest, &state)
	result := recordingplan.SchemaOperationResult{
		Status:        recordingplan.SchemaOperationSucceeded,
		PayloadDigest: claimedOp.PayloadDigest,
		WriteOutcome:  claimedOp.WriteOutcome,
		Reason:        claimedOp.Reason,
		CleanupStatus: cleanup.status,
		CleanupReason: cleanup.reason,
		Detail:        state.encode(),
	}
	return s.finish(ctx, claimedOp, result)
}
