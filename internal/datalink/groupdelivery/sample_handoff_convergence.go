package groupdelivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

func canonicalTransferID(r handoffReceipt) string {
	if r.CanonicalTransferID != "" {
		return r.CanonicalTransferID
	}
	return r.ID
}

// Shared meaning excludes legitimate branch-specific fence, source and fanout
// facts. Each edge must independently prove those facts before this is authority.
func sameHandoffTarget(a, b handoffReceipt) bool {
	return sameHandoffOriginal(a, b) && a.Target == b.Target && a.TargetBucket == b.TargetBucket &&
		a.TargetProof == b.TargetProof && a.TargetRuntimeDigest == b.TargetRuntimeDigest
}

// Verification caches live only within one validation walk and its query
// snapshot. Group recursion checks remain separate from completed-edge memoizing
// so converged ancestors are linear work, never an authority cache across reads.
type handoffVerification struct {
	ownershipQuery
	groups         map[int64]*handoffReceipt
	active         map[int64]bool
	edges          map[string]string
	lineage        map[string]bool
	lineageActive  map[string]bool
	terminals      map[int64]bool
	terminalActive map[int64]bool
	consumedGroups map[int64]bool
}

func handoffVerifier(q ownershipQuery) *handoffVerification {
	if v, ok := q.(*handoffVerification); ok {
		return v
	}
	return &handoffVerification{
		ownershipQuery: q,
		groups:         map[int64]*handoffReceipt{}, active: map[int64]bool{},
		edges: map[string]string{}, lineage: map[string]bool{}, lineageActive: map[string]bool{},
		terminals: map[int64]bool{}, terminalActive: map[int64]bool{},
		consumedGroups: map[int64]bool{},
	}
}

func verifyHandoffReceipt(ctx context.Context, q ownershipQuery, r handoffReceipt, seen map[string]bool) error {
	v := handoffVerifier(q)
	if r.TargetSeq > 0 {
		if _, err := v.incoming(ctx, r.TargetSeq, seen); err != nil {
			return err
		}
		raw, _ := json.Marshal(r)
		if v.edges[r.ID] != string(raw) {
			return ErrOwnershipUnavailable
		}
		return nil
	}
	return v.edge(ctx, r, seen)
}

func (v *handoffVerification) edge(ctx context.Context, r handoffReceipt, seen map[string]bool) error {
	raw, err := json.Marshal(r)
	if err != nil {
		return ErrOwnershipUnavailable
	}
	if previous, ok := v.edges[r.ID]; ok {
		if previous != string(raw) || seen[r.ID] {
			return ErrOwnershipUnavailable
		}
		return nil
	}
	if err = verifyHandoffReceiptEdge(ctx, v, r, seen); err != nil {
		return err
	}
	if r.State == "accepted" {
		if err = v.consumed(ctx, r); err != nil {
			return err
		}
	}
	if r.State == "finalized" {
		if err = v.compactLineage(ctx, r); err != nil {
			return err
		}
	}
	v.edges[r.ID] = string(raw)
	return nil
}

func incomingHandoff(ctx context.Context, q ownershipQuery, seq int64) (*handoffReceipt, error) {
	return handoffVerifier(q).incoming(ctx, seq, map[string]bool{})
}

func (v *handoffVerification) incoming(ctx context.Context, seq int64, seen map[string]bool) (*handoffReceipt, error) {
	if v.active[seq] {
		return nil, ErrOwnershipUnavailable
	}
	if canonical, ok := v.groups[seq]; ok {
		return canonical, nil
	}
	v.active[seq] = true
	defer delete(v.active, seq)
	peers, err := queryHandoffs(ctx, v, "target_seq=?", seq)
	if err != nil || len(peers) == 0 {
		return nil, err
	}
	canonical, err := handoffCanonical(ctx, v, seq, peers)
	if err != nil {
		return nil, err
	}
	for _, r := range peers {
		if err = v.edge(ctx, r, seen); err != nil {
			return nil, err
		}
	}
	v.groups[seq] = canonical
	return canonical, nil
}

// Select authority without traversing either lineage direction. The caller
// must validate every peer, not just the returned primary.
func handoffCanonical(ctx context.Context, q ownershipQuery, seq int64, peers []handoffReceipt) (*handoffReceipt, error) {
	// Live authority is the journal's handoff_id, never SQL result order.
	row, rowErr := readHandoffSample(ctx, q, seq)
	canonicalID := row.HandoffID
	if rowErr != nil {
		if !errors.Is(rowErr, sql.ErrNoRows) {
			return nil, rowErr
		}
		// After GC the unique primary audit must still exist. References must
		// resolve to it; an arbitrary surviving alias is never a replacement.
		for _, peer := range peers {
			if peer.CanonicalTransferID == "" {
				if canonicalID != "" {
					return nil, ErrOwnershipUnavailable
				}
				canonicalID = peer.ID
			}
		}
	}
	var canonical *handoffReceipt
	for i := range peers {
		r := &peers[i]
		if canonicalTransferID(*r) != canonicalID || r.TargetSeq != seq ||
			(r.State != "accepted" && r.State != "finalized") ||
			(rowErr != nil && r.State != "finalized") {
			return nil, ErrOwnershipUnavailable
		}
		if r.ID == canonicalID {
			if r.CanonicalTransferID != "" || canonical != nil {
				return nil, ErrOwnershipUnavailable
			}
			canonical = r
		} else if r.CanonicalTransferID == "" {
			return nil, ErrOwnershipUnavailable
		}
	}
	if canonical == nil {
		return nil, ErrOwnershipUnavailable
	}
	for _, peer := range peers {
		if !sameHandoffTarget(*canonical, peer) {
			return nil, ErrHandoffProofConflict
		}
	}
	return canonical, nil
}

// Compact ancestry is an independent backward-only walk. Combining it with the
// terminal forward walk would turn an ordinary parent/child link into a false
// cycle. Both walks retain their own cycle checks and snapshot-local memoizing.
func (v *handoffVerification) compactLineage(ctx context.Context, r handoffReceipt) error {
	if v.lineageActive[r.ID] {
		return ErrOwnershipUnavailable
	}
	if v.lineage[r.ID] {
		return nil
	}
	v.lineageActive[r.ID] = true
	defer delete(v.lineageActive, r.ID)
	if r.State == "finalized" {
		// A finalized source disposes its entire immutable fanout. Validate
		// siblings' forward obligations, not only the path back to this edge.
		siblings, err := outgoingHandoffs(ctx, v, r.SourceSeq)
		if err != nil {
			return err
		}
		for _, sibling := range siblings {
			if sibling.State != "finalized" {
				return ErrOwnershipUnavailable
			}
			if err = v.terminal(ctx, sibling); err != nil {
				return err
			}
		}
	}
	if r.SourceSeq == r.OriginSeq {
		if r.Source != r.Origin {
			return ErrHandoffProofConflict
		}
		v.lineage[r.ID] = true
		return nil
	}
	peers, err := queryHandoffs(ctx, v, "target_seq=?", r.SourceSeq)
	if err != nil {
		return err
	}
	if _, err = handoffCanonical(ctx, v, r.SourceSeq, peers); err != nil {
		return err
	}
	for _, prior := range peers {
		if !handoffSourceMatches(prior, r) {
			return ErrHandoffProofConflict
		}
		if err = verifyHandoffBranch(ctx, v, prior); err != nil {
			return err
		}
		if err = verifyCompactHandoffMeaning(ctx, v, prior); err != nil {
			return err
		}
		if err = v.compactLineage(ctx, prior); err != nil {
			return err
		}
	}
	v.lineage[r.ID] = true
	return nil
}

// A forward continuation and backward predecessor use the same complete
// original identity, while each edge's target proof remains its own authority.
func handoffSourceMatches(prior, next handoffReceipt) bool {
	return prior.Target == next.Source && prior.TargetSeq == next.SourceSeq &&
		prior.TargetBucket == next.SourceBucket && sameHandoffOriginal(prior, next)
}

func sameHandoffOriginal(a, b handoffReceipt) bool {
	return a.Origin == b.Origin && a.OriginSeq == b.OriginSeq && a.SampleID == b.SampleID &&
		a.Member == b.Member && a.ObservedAt == b.ObservedAt && a.Digest == b.Digest &&
		a.RawDigest == b.RawDigest && a.AcceptedAt == b.AcceptedAt && a.SampleProof == b.SampleProof &&
		a.OriginalProof == b.OriginalProof && a.OriginRuntimeDigest == b.OriginRuntimeDigest
}

// The target acceptance and independently verified source edge are committed in
// the caller's transaction. Coalescing never changes a canonical row or proof.
func acceptHandoffTarget(ctx context.Context, tx *sql.Tx, row, origin handoffSample, r *handoffReceipt) error {
	var seq int64
	err := tx.QueryRowContext(ctx, `SELECT seq FROM wg_delivery_samples WHERE workspace_id=? AND group_id=? AND group_revision=? AND acceptance_kind='handoff' AND origin_workspace_id=? AND origin_group_id=? AND origin_revision=? AND origin_sample_id=?`, r.Target.WorkspaceID, r.Target.GroupID, r.Target.GroupRevision, origin.Key.WorkspaceID, origin.Key.GroupID, origin.Key.GroupRevision, row.SampleID).Scan(&seq)
	if err == nil {
		canonical, e := incomingHandoff(ctx, tx, seq)
		if e != nil {
			return e
		}
		if canonical == nil || !sameHandoffTarget(*canonical, *r) {
			return ErrHandoffProofConflict
		}
		target, e := readHandoffSample(ctx, tx, seq)
		if e != nil {
			return e
		}
		if target.Payload != row.Payload || target.Consumed {
			return ErrHandoffProofConflict
		}
		r.TargetSeq, r.CanonicalTransferID = seq, canonical.ID
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return &handoffStorageError{cause: err}
	}
	address, err := allocateJournalAddress(ctx, tx, r.Target, r.ID)
	if err != nil {
		return err
	}
	insert, err := tx.ExecContext(ctx, `INSERT INTO wg_delivery_samples(workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,consumed,created_at,acceptance_kind,origin_workspace_id,origin_group_id,origin_revision,origin_sample_id,handoff_id) VALUES(?,?,?,?,?,?,?,?,?,0,?,'handoff',?,?,?,?,?)`, r.Target.WorkspaceID, r.Target.GroupID, r.Target.GroupRevision, address, row.Member, row.ObservedAt, r.TargetBucket, row.Payload, row.Digest, row.AcceptedAt, origin.Key.WorkspaceID, origin.Key.GroupID, origin.Key.GroupRevision, row.SampleID, r.ID)
	if err != nil {
		return err
	}
	r.TargetSeq, err = insert.LastInsertId()
	return err
}
