package groupdelivery

import (
	"context"
	"database/sql"
	"time"

	"go-gateway/internal/datalink/snapshot"
)

// compactHandoffs runs only inside normal Reclaim. Transferred intermediates do
// not become fake consumers: a consumed leaf terminally completes its ancestry.
func compactHandoffs(ctx context.Context, tx *sql.Tx) error {
	rows, err := tx.QueryContext(ctx, `SELECT DISTINCT source_seq FROM wg_delivery_sample_handoffs WHERE state='accepted' ORDER BY source_seq`)
	if err != nil {
		return ErrOwnershipUnavailable
	}
	var seqs []int64
	for rows.Next() {
		var seq int64
		if rows.Scan(&seq) != nil {
			rows.Close()
			return ErrOwnershipUnavailable
		}
		seqs = append(seqs, seq)
	}
	if rows.Err() != nil {
		rows.Close()
		return ErrOwnershipUnavailable
	}
	if rows.Close() != nil {
		return ErrOwnershipUnavailable
	}
	for _, seq := range seqs {
		if _, err := finalizeHandoffSource(ctx, tx, seq, map[int64]bool{}); err != nil {
			return err
		}
	}
	return nil
}

func finalizeHandoffSource(ctx context.Context, tx *sql.Tx, seq int64, seen map[int64]bool) (bool, error) {
	if seen[seq] {
		return false, ErrOwnershipUnavailable
	}
	seen[seq] = true
	defer delete(seen, seq)
	receipts, err := outgoingHandoffs(ctx, tx, seq)
	if err != nil {
		return false, err
	}
	if len(receipts) == 0 {
		return false, nil
	}
	for _, r := range receipts {
		if err = verifyHandoffReceipt(ctx, tx, r, map[string]bool{}); err != nil {
			return false, err
		}
		switch r.State {
		case "pending", "quarantined":
			return false, nil
		case "finalized":
			continue
		}
		target, err := readHandoffSample(ctx, tx, r.TargetSeq)
		if err != nil {
			return false, err
		}
		if !target.Consumed {
			done, err := finalizeHandoffSource(ctx, tx, r.TargetSeq, seen)
			if err != nil {
				return false, err
			}
			if !done {
				return false, nil
			}
		}
	}
	for _, r := range receipts {
		if r.State == "finalized" {
			continue
		}
		r.State = "finalized"
		if err = writeHandoffReceipt(ctx, tx, r); err != nil {
			return false, err
		}
	}
	return true, nil
}

// A row remains pinned while any receipt needs its raw bytes. Compact receipts
// retain original identity/proof authority after eligible payload reclamation.
const handoffReclaimSamplesSQL = `DELETE FROM wg_delivery_samples
 WHERE (consumed=1 OR EXISTS(SELECT 1 FROM wg_delivery_sample_handoffs h WHERE h.source_seq=wg_delivery_samples.seq AND h.state='finalized'))
 AND NOT EXISTS(SELECT 1 FROM wg_delivery_sample_handoffs h WHERE
 (h.source_seq=wg_delivery_samples.seq OR h.target_seq=wg_delivery_samples.seq OR h.origin_seq=wg_delivery_samples.seq)
 AND h.state<>'finalized')`

func verifyHandoffTerminal(ctx context.Context, q ownershipQuery, r handoffReceipt) error {
	// Existing target bytes must agree even after compaction; state alone is
	// never a substitute for consumption or a complete terminal descendant set.
	var consumed bool
	err := q.QueryRowContext(ctx, "SELECT consumed FROM wg_delivery_samples WHERE seq=?", r.TargetSeq).Scan(&consumed)
	if err == nil {
		target, e := readHandoffSample(ctx, q, r.TargetSeq)
		if e != nil {
			return e
		}
		if target.Key != r.Target || target.Kind != "handoff" || target.HandoffID != canonicalTransferID(r) || target.Origin != r.Origin || target.Bucket != r.TargetBucket || !receiptMatchesRow(r, target) {
			return ErrHandoffProofConflict
		}
		if consumed {
			return nil
		}
	} else if err != sql.ErrNoRows {
		return ErrOwnershipUnavailable
	}
	children, e := outgoingHandoffs(ctx, q, r.TargetSeq)
	if e != nil {
		return e
	}
	if len(children) > 0 {
		for _, child := range children {
			if child.State != "finalized" || !handoffSourceMatches(r, child) {
				return ErrOwnershipUnavailable
			}
			if e = handoffVerifier(q).terminal(ctx, child); e != nil {
				return e
			}
		}
		return nil
	}
	if err == nil {
		return ErrOwnershipUnavailable
	} // still physically open
	// A predecessor may legitimately close an empty truncated window whose
	// aligned checkpoint is later than this bucket. Only an actual owner can
	// use that checkpoint as terminal consumption evidence for this observation.
	observed, e := time.Parse(time.RFC3339Nano, r.ObservedAt)
	if e != nil {
		return ErrOwnershipUnavailable
	}
	ownership, e := resolveCutoverOwnership(ctx, q, sourceOf(r.Target, snapshot.Sample{MemberKey: r.Member}), observed, r.Target)
	if e != nil {
		return e
	}
	if !ownership.Owners[r.Target] {
		return ErrOwnershipUnavailable
	}
	var next string
	if q.QueryRowContext(ctx, "SELECT next_close FROM wg_delivery_checkpoints WHERE group_id=? AND group_revision=?", r.Target.GroupID, r.Target.GroupRevision).Scan(&next) != nil {
		return ErrOwnershipUnavailable
	}
	closed, e := time.Parse(time.RFC3339Nano, next)
	if e != nil {
		return ErrOwnershipUnavailable
	}
	bucket, e := time.Parse(time.RFC3339Nano, r.TargetBucket)
	if e != nil || !bucket.Before(closed) {
		return ErrOwnershipUnavailable
	}
	return nil
}

// Forward-only terminal validation never calls compactLineage. Separate active
// groups distinguish real descendant cycles from legitimate backward references.
func (v *handoffVerification) terminal(ctx context.Context, r handoffReceipt) error {
	if v.terminalActive[r.TargetSeq] {
		return ErrOwnershipUnavailable
	}
	if v.terminals[r.TargetSeq] {
		return nil
	}
	v.terminalActive[r.TargetSeq] = true
	defer delete(v.terminalActive, r.TargetSeq)
	peers, err := queryHandoffs(ctx, v, "target_seq=?", r.TargetSeq)
	if err != nil {
		return err
	}
	if _, err = handoffCanonical(ctx, v, r.TargetSeq, peers); err != nil {
		return err
	}
	for _, peer := range peers {
		if err = verifyHandoffBranch(ctx, v, peer); err != nil {
			return err
		}
		if err = verifyCompactHandoffMeaning(ctx, v, peer); err != nil {
			return err
		}
		if err = verifyHandoffTerminal(ctx, v, peer); err != nil {
			return err
		}
	}
	v.terminals[r.TargetSeq] = true
	return nil
}

// An accepted receipt can already have a consumed target before normal GC.
// Check only this target's converging obligations: an independent pending
// profile is not required to finish merely because another profile consumed.
func (v *handoffVerification) consumed(ctx context.Context, r handoffReceipt) error {
	if v.consumedGroups[r.TargetSeq] {
		return nil
	}
	var consumed bool
	if err := v.QueryRowContext(ctx, "SELECT consumed FROM wg_delivery_samples WHERE seq=?", r.TargetSeq).Scan(&consumed); err != nil {
		return &handoffStorageError{cause: err}
	}
	if !consumed {
		return nil
	}
	observed, err := time.Parse(time.RFC3339Nano, r.ObservedAt)
	if err != nil {
		return ErrOwnershipUnavailable
	}
	ownership, err := resolveCutoverOwnership(ctx, v, sourceOf(r.Origin, snapshot.Sample{MemberKey: r.Member}), observed, GroupKey{})
	if err != nil {
		return err
	}
	if !ownership.Owners[r.Target] {
		return ErrOwnershipUnavailable
	}
	receipts, err := queryHandoffs(ctx, v, "origin_seq=?", r.OriginSeq)
	if err != nil {
		return err
	}
	outgoing := map[int64][]handoffReceipt{}
	for _, receipt := range receipts {
		outgoing[receipt.SourceSeq] = append(outgoing[receipt.SourceSeq], receipt)
	}
	active, done := map[int64]bool{}, map[int64]bool{}
	var walk func(int64, GroupKey) error
	walk = func(seq int64, key GroupKey) error {
		if seq == r.TargetSeq {
			if key != r.Target {
				return ErrHandoffProofConflict
			}
			return nil
		}
		if active[seq] || key == r.Target {
			return ErrOwnershipUnavailable
		}
		if done[seq] {
			return nil
		}
		active[seq] = true
		defer delete(active, seq)
		required := false
		for _, edge := range outgoing[seq] {
			if edge.Source != key {
				return ErrHandoffProofConflict
			}
			if !handoffDescendant(ownership, r.Target, edge.Target, map[GroupKey]bool{}) {
				continue
			}
			required = true
			if edge.State != "accepted" && edge.State != "finalized" {
				return ErrOwnershipUnavailable
			}
			if !sameHandoffOriginal(r, edge) {
				return ErrHandoffProofConflict
			}
			// Validate the edge's own raw/proof/lineage authority without
			// re-entering this same consumed-target obligation walk.
			if err := verifyHandoffReceiptEdge(ctx, v, edge, map[string]bool{}); err != nil {
				return err
			}
			if err := walk(edge.TargetSeq, edge.Target); err != nil {
				return err
			}
		}
		if !required {
			return ErrOwnershipUnavailable
		}
		done[seq] = true
		return nil
	}
	if err = walk(r.OriginSeq, r.Origin); err != nil {
		return err
	}
	v.consumedGroups[r.TargetSeq] = true
	return nil
}
