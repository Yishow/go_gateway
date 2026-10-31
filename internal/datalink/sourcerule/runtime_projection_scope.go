package sourcerule

import (
	"context"
	"go-gateway/internal/datalink/collector"
)

// RuntimePointProof is an internal mutation receipt, never supplied by HTTP.
type RuntimePointProof struct {
	Revision string
	Lease    uint64
	Deleted  bool
}

// RuntimeProjectionProof binds the rule-owned point receipts to their last contiguous owned epoch.
type RuntimeProjectionProof struct {
	Epoch  uint64
	Points map[string]RuntimePointProof
}

type sourceProjectionBaseline interface {
	CaptureOrdinaryPoint(string) (collector.PointMutationLease, bool)
}

func (s *Service) captureProjectionEpoch(ctx context.Context, id string) {
	j := s.pointJournal(ctx)
	runtime, ok := s.runtimeSync.(sourceProjectionBaseline)
	if j == nil || !ok || j.projectionEpochKnown {
		return
	}
	proof, captured := runtime.CaptureOrdinaryPoint(id)
	if !captured {
		j.runtimeUncertain = true
		return
	}
	j.projectionEpoch, j.projectionEpochKnown = proof.Epoch, true
}

func (s *Service) sourceProjectionProof(ctx context.Context) *RuntimeProjectionProof {
	j := s.pointJournal(ctx)
	if j == nil || !j.projectionEpochKnown {
		return nil
	}
	proof := &RuntimeProjectionProof{Epoch: advanceOwnedEpoch(j, j.projectionEpoch), Points: map[string]RuntimePointProof{}}
	for id, e := range j.points {
		p := e.current
		if p == nil {
			p = e.before
		}
		// Rolled-back creations have neither a before nor current durable row;
		// retain their original owned write identity to clear only their projection keys.
		if p == nil && len(e.attempted) > 0 {
			p = e.attempted[len(e.attempted)-1]
		}
		if p == nil || p.DeviceID != j.deviceID {
			continue
		}
		proof.Points[id] = RuntimePointProof{Revision: p.RevisionID, Lease: e.currentLease, Deleted: e.current == nil}
	}
	return proof
}

func (s *Service) observeProjectionRemoval(ctx context.Context, id string, epoch uint64) {
	j := s.pointJournal(ctx)
	if j == nil {
		return
	}
	e := j.points[id]
	if e == nil || (e.current != nil && e.current.Enabled) {
		return
	}
	if j.ownedEpochs == nil {
		j.ownedEpochs = map[uint64]bool{}
	}
	j.ownedEpochs[epoch] = true
}

// Copy only the retained snapshot identities; no current read can expand this set.
func snapshotPointIDs(snapshot *sourceRuleSnapshot) []string {
	if snapshot == nil {
		return nil
	}
	ids := make([]string, 0, len(snapshot.points))
	for id := range snapshot.points {
		ids = append(ids, id)
	}
	return ids
}

type sourceAffectedPointsKey struct{}
type sourceProjectionSnapshotKey struct{}

func withSourceAffectedPoints(ctx context.Context, snapshot *sourceRuleSnapshot) context.Context {
	if snapshot != nil && snapshot.runtimeProjection != nil {
		ctx = context.WithValue(ctx, sourceProjectionSnapshotKey{}, snapshot.runtimeProjection)
	}
	return context.WithValue(ctx, sourceAffectedPointsKey{}, snapshotPointIDs(snapshot))
}
