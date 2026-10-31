package groupdelivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"go-gateway/internal/datalink/snapshot"
)

// HandoffBatch reports durable progress, not destination delivery.
type HandoffBatch struct {
	// Participants is advisory rebuild coordination, never ownership authority.
	// It includes only verified scopes whose per-original transaction committed.
	Participants []GroupKey
	// ReplayParticipants excludes pending/quarantined targets so independent
	// proven fanout can rebuild without claiming blocked input was accepted.
	ReplayParticipants []GroupKey
	NextSeq            int64
	More               bool
	Pending            int
	Quarantined        int
	Accepted           int
}

// Discovery is only an interest hint. Keep damaged or orphaned evidence for
// decisive transaction validation, and never require a currently applied group:
// pending receipts and accepted intermediate copies can name historical scopes.
const handoffMemberCandidateSQL = `json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.device_id')=json_extract(CASE WHEN json_valid(s.member_key) THEN s.member_key ELSE '[]' END,'$[0]')
 AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.point_id')=json_extract(CASE WHEN json_valid(s.member_key) THEN s.member_key ELSE '[]' END,'$[1]')
 AND json_extract(CASE WHEN json_valid(m.value) THEN m.value ELSE '{}' END,'$.tag_id')=json_extract(CASE WHEN json_valid(s.member_key) THEN s.member_key ELSE '[]' END,'$[2]')`

const handoffCandidateSQL = `(s.acceptance_kind='handoff'
 OR EXISTS(SELECT 1 FROM wg_delivery_sample_handoffs h WHERE h.source_seq=s.seq OR h.target_seq=s.seq)
 OR EXISTS(SELECT 1 FROM wg_delivery_samples origin_copy WHERE origin_copy.acceptance_kind='handoff'
  AND origin_copy.origin_workspace_id=s.workspace_id AND origin_copy.origin_group_id=s.group_id AND origin_copy.origin_revision=s.group_revision)
 OR EXISTS(SELECT 1 FROM write_group_revision_ends e WHERE e.group_id=s.group_id AND e.group_revision=s.group_revision)
 OR EXISTS(SELECT 1 FROM write_group_revision_ends e JOIN write_group_versions v ON v.group_id=e.group_id AND v.group_revision=e.group_revision,
  json_each(CASE WHEN json_valid(v.payload) THEN v.payload ELSE '{}' END,'$.members') m WHERE e.workspace_id=s.workspace_id AND ` + handoffMemberCandidateSQL + `)
 OR EXISTS(SELECT 1 FROM write_group_ownership_fences f WHERE f.group_id=s.group_id AND f.group_revision=s.group_revision)
 OR EXISTS(SELECT 1 FROM write_group_ownership_fences f WHERE f.workspace_id=s.workspace_id
  AND f.device_id=json_extract(CASE WHEN json_valid(s.member_key) THEN s.member_key ELSE '[]' END,'$[0]')
  AND f.point_id=json_extract(CASE WHEN json_valid(s.member_key) THEN s.member_key ELSE '[]' END,'$[1]')
  AND f.tag_id=json_extract(CASE WHEN json_valid(s.member_key) THEN s.member_key ELSE '[]' END,'$[2]'))
 OR EXISTS(SELECT 1 FROM write_group_cutovers c WHERE c.workspace_id=s.workspace_id AND (
  c.source_group_id=s.group_id
	 OR EXISTS(SELECT 1 FROM json_each(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.moved') m WHERE ` + handoffMemberCandidateSQL + `)
	 OR EXISTS(SELECT 1 FROM write_group_versions v,json_each(CASE WHEN json_valid(v.payload) THEN v.payload ELSE '{}' END,'$.members') m
	  WHERE v.group_id=c.source_group_id AND ` + handoffMemberCandidateSQL + `)
  OR json_type(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.ended') IS NOT 'array'
  OR json_type(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.started') IS NOT 'array'
  OR EXISTS(SELECT 1 FROM json_each(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.ended') r
   WHERE json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.group_id')=s.group_id)
  OR EXISTS(SELECT 1 FROM json_each(CASE WHEN json_valid(c.detail) THEN c.detail ELSE '{}' END,'$.started') r
   WHERE json_extract(CASE WHEN json_valid(r.value) THEN r.value ELSE '{}' END,'$.group_id')=s.group_id))))`

// HandoffSources pages original and pending journals in one explicit workspace.
func (s *Store) HandoffSources(ctx context.Context, workspaceID string, after GroupKey, limit int) ([]GroupKey, error) {
	if workspaceID == "" || limit < 1 || limit > 128 || (after != (GroupKey{}) && (!after.Valid() || after.WorkspaceID != workspaceID)) {
		return nil, ErrOwnershipUnavailable
	}
	rows, err := s.db.QueryContext(ctx, `SELECT DISTINCT s.group_id,s.group_revision FROM wg_delivery_samples s WHERE s.workspace_id=? AND s.consumed=0 AND `+handoffCandidateSQL+` AND (s.group_id>? OR (s.group_id=? AND s.group_revision>?)) ORDER BY s.group_id,s.group_revision LIMIT ?`, workspaceID, after.GroupID, after.GroupID, after.GroupRevision, limit)
	if err != nil {
		return nil, ErrOwnershipUnavailable
	}
	defer rows.Close()
	var keys []GroupKey
	for rows.Next() {
		key := GroupKey{WorkspaceID: workspaceID}
		if rows.Scan(&key.GroupID, &key.GroupRevision) != nil {
			return nil, ErrOwnershipUnavailable
		}
		keys = append(keys, key)
	}
	if rows.Err() != nil {
		return nil, ErrOwnershipUnavailable
	}
	return keys, nil
}

// PrepareCutoverHandoffs atomically partitions each original accepted input.
// No caller-supplied sample, target proof or owner is adopted as authority.
func (s *Store) PrepareCutoverHandoffs(ctx context.Context, source GroupKey, afterSeq int64, limit int) (HandoffBatch, error) {
	var batch HandoffBatch
	if !source.Valid() || afterSeq < 0 || limit < 1 || limit > 128 {
		return batch, ErrOwnershipUnavailable
	}
	rows, err := s.db.QueryContext(ctx, `SELECT s.seq FROM wg_delivery_samples s WHERE s.workspace_id=? AND s.group_id=? AND s.group_revision=? AND s.consumed=0 AND `+handoffCandidateSQL+` AND s.seq>? ORDER BY s.seq LIMIT ?`, source.WorkspaceID, source.GroupID, source.GroupRevision, afterSeq, limit+1)
	if err != nil {
		return batch, ErrOwnershipUnavailable
	}
	var seqs []int64
	for rows.Next() {
		var seq int64
		if rows.Scan(&seq) != nil {
			rows.Close()
			return batch, ErrOwnershipUnavailable
		}
		seqs = append(seqs, seq)
	}
	err = errors.Join(rows.Err(), rows.Close())
	if err != nil {
		return batch, ErrOwnershipUnavailable
	}
	batch.More = len(seqs) > limit
	if batch.More {
		seqs = seqs[:limit]
	}
	var blocked error
	for _, seq := range seqs {
		var states []string
		var semanticErr error
		var participants, replayParticipants []GroupKey
		err = s.inTx(ctx, "handoff", func(ctx context.Context, tx *sql.Tx) error {
			var e error
			states, semanticErr, e = s.prepareHandoffInTx(ctx, tx, source, seq)
			if e != nil {
				return e
			}
			receipts, e := outgoingHandoffs(ctx, tx, seq)
			if e != nil {
				return e
			}
			for _, receipt := range receipts {
				if e = verifyHandoffReceipt(ctx, tx, receipt, map[string]bool{}); e != nil {
					return e
				}
				participants = append(participants, receipt.Source, receipt.Target)
				if receipt.State == "accepted" || receipt.State == "finalized" {
					replayParticipants = append(replayParticipants, receipt.Source, receipt.Target)
				}
			}
			return nil
		})
		s.usage.invalidate()
		if err != nil {
			converted := s.capacityError(ctx, err)
			var quota *QuotaError
			if errors.As(converted, &quota) || errors.Is(converted, ErrOwnershipConflict) || errors.Is(converted, ErrOwnershipUnavailable) || errors.Is(converted, ErrHandoffProofConflict) {
				return batch, converted
			}
			return batch, &handoffStorageError{cause: converted}
		}
		batch.Participants = mergeHandoffParticipants(batch.Participants, participants)
		batch.ReplayParticipants = mergeHandoffParticipants(batch.ReplayParticipants, replayParticipants)
		batch.NextSeq = seq
		for _, state := range states {
			switch state {
			case "pending":
				batch.Pending++
			case "quarantined":
				batch.Quarantined++
			case "accepted", "finalized":
				batch.Accepted++
			}
		}
		blocked = errors.Join(blocked, semanticErr)
	}
	return batch, blocked
}

// Publish only after commit acknowledgement; an uncertain transaction is
// rediscovered on retry, while earlier acknowledged batch items remain useful.
func mergeHandoffParticipants(previous, current []GroupKey) []GroupKey {
	seen := make(map[GroupKey]bool, len(previous)+len(current))
	for _, key := range previous {
		seen[key] = true
	}
	for _, key := range current {
		if !seen[key] {
			previous = append(previous, key)
			seen[key] = true
		}
	}
	sort.Slice(previous, func(i, j int) bool {
		if previous[i].WorkspaceID != previous[j].WorkspaceID {
			return previous[i].WorkspaceID < previous[j].WorkspaceID
		}
		if previous[i].GroupID != previous[j].GroupID {
			return previous[i].GroupID < previous[j].GroupID
		}
		return previous[i].GroupRevision < previous[j].GroupRevision
	})
	return previous
}

func sourceOf(key GroupKey, sample snapshot.Sample) OwnershipSource {
	source := OwnershipSource{WorkspaceID: key.WorkspaceID}
	var ids []string
	if json.Unmarshal([]byte(sample.MemberKey), &ids) == nil && len(ids) == 3 {
		source.DeviceID, source.PointID, source.TagID = ids[0], ids[1], ids[2]
	}
	return source
}

func (s *Store) prepareHandoffInTx(ctx context.Context, tx *sql.Tx, key GroupKey, seq int64) ([]string, error, error) {
	row, err := readHandoffSample(ctx, tx, seq)
	if err != nil {
		return nil, nil, err
	}
	if row.Key != key || row.Consumed {
		return nil, nil, ErrOwnershipConflict
	}
	existingReceipts, e := outgoingHandoffs(ctx, tx, seq)
	if e != nil {
		return nil, nil, e
	}
	allAccepted := len(existingReceipts) > 0
	var existingStates []string
	for _, r := range existingReceipts {
		if e = verifyHandoffReceipt(ctx, tx, r, map[string]bool{}); e != nil {
			return nil, nil, e
		}
		if r.State != "accepted" && r.State != "finalized" {
			allAccepted = false
		}
		existingStates = append(existingStates, r.State)
	}
	if allAccepted {
		return existingStates, nil, nil
	}
	sample, err := decodeVerifiedSample(row.Payload, row.Digest)
	if err != nil {
		return nil, nil, ErrHandoffProofConflict
	}
	own, err := resolveCutoverOwnership(ctx, tx, sourceOf(key, sample), sample.ObservedAt, GroupKey{})
	if err != nil {
		return nil, nil, err
	}
	if !own.Affected[key] || own.Owners[key] {
		return nil, nil, nil
	}
	origin := row
	incoming, err := incomingHandoff(ctx, tx, seq)
	if err != nil {
		return nil, nil, err
	}
	if incoming == nil && row.Kind == "handoff" {
		return nil, nil, ErrOwnershipUnavailable
	}
	if incoming != nil {
		if err = verifyHandoffReceipt(ctx, tx, *incoming, map[string]bool{}); err != nil {
			return nil, nil, err
		}
		origin, err = readHandoffSample(ctx, tx, incoming.OriginSeq)
		if err != nil {
			return nil, nil, err
		}
	}
	var targets []GroupKey
	if len(existingReceipts) > 0 {
		// A persisted fanout is an immutable transfer intent, including while
		// pending. Complete its original target first; a later cutover forwards
		// that accepted copy through the normal verified handoff chain.
		targets = append(targets, existingReceipts[0].Targets...)
	} else {
		for target := range own.Owners {
			if handoffDescendant(own, target, key, map[GroupKey]bool{}) {
				targets = append(targets, target)
			}
		}
	}
	sort.Slice(targets, func(i, j int) bool {
		if targets[i].GroupID != targets[j].GroupID {
			return targets[i].GroupID < targets[j].GroupID
		}
		return targets[i].GroupRevision < targets[j].GroupRevision
	})
	if len(targets) == 0 {
		return nil, nil, ErrOwnershipUnavailable
	}
	beforeQuota, err := s.handoffQuotaBefore(ctx, tx, targets)
	if err != nil {
		return nil, nil, err
	}
	var states []string
	var blocked error
	for _, target := range targets {
		fence, err := handoffFence(ctx, tx, key, target, sourceOf(key, sample), sample.ObservedAt)
		if err != nil {
			return nil, nil, err
		}
		r := newHandoffReceipt(row, origin, target, fence)
		r.Targets = targets
		r.ID = handoffID(r)
		old, err := readHandoffID(ctx, tx, r.ID)
		if err != nil {
			return nil, nil, err
		}
		if old != nil {
			if err = verifyHandoffReceipt(ctx, tx, *old, map[string]bool{}); err != nil {
				return nil, nil, err
			}
			if old.State == "accepted" || old.State == "finalized" {
				states = append(states, old.State)
				continue
			}
		}
		proof, proofErr := validateHandoffProof(ctx, tx, handoffProofInput{Origin: origin.Key, Sequence: origin.Seq, Payload: origin.Payload, Digest: origin.Digest, AcceptedAt: origin.AcceptedAt}, target)
		if proofErr != nil {
			r.State = "pending"
			r.Reason = ErrOwnershipUnavailable.Error()
			if errors.Is(proofErr, ErrHandoffProofConflict) {
				r.State = "quarantined"
				r.Reason = ErrHandoffProofConflict.Error()
			}
			blocked = errors.Join(blocked, proofErr)
		} else {
			r.OriginalProof, r.TargetProof = proof.OriginalReference, proof.TargetReference
			r.OriginRuntimeDigest, r.TargetRuntimeDigest = proof.OriginRuntimeDigest, proof.TargetRuntimeDigest
			if proof.TargetMixed {
				start, _, e := snapshot.CheckedWindowBounds(sample.ObservedAt, proof.TargetInterval)
				if e != nil {
					return nil, nil, ErrHandoffProofConflict
				}
				r.TargetBucket = timestamp(start)
				var progress Restored
				if e := restoreMixedProgress(ctx, tx, target, &progress); e != nil {
					return nil, nil, e
				}
				if progress.Mixed != nil {
					member, ok := progress.Mixed.Members[sample.MemberKey]
					if !ok || sample.ObservedAt.Before(member.Start) || sample.ObservedAt.Before(progress.Mixed.StartsAt) || (!progress.Mixed.EndsAt.IsZero() && !sample.ObservedAt.Before(progress.Mixed.EndsAt)) {
						return nil, nil, ErrOwnershipConflict
					}
				}
			} else {
				r.TargetBucket = timestamp(snapshot.BucketStart(sample.ObservedAt, proof.TargetInterval))
			}
			var checkpoint string
			e := tx.QueryRowContext(ctx, `SELECT next_close FROM wg_delivery_checkpoints WHERE group_id=? AND group_revision=?`, target.GroupID, target.GroupRevision).Scan(&checkpoint)
			if e != nil && !errors.Is(e, sql.ErrNoRows) {
				return nil, nil, ErrOwnershipUnavailable
			}
			if e == nil {
				closed, e := time.Parse(time.RFC3339Nano, checkpoint)
				if e != nil {
					return nil, nil, ErrOwnershipUnavailable
				}
				if sample.ObservedAt.Before(closed) {
					return nil, nil, ErrOwnershipConflict
				}
			}

			if err = acceptHandoffTarget(ctx, tx, row, origin, &r); err != nil {
				return nil, nil, err
			}
			r.State = "accepted"
		}
		if err = writeHandoffReceipt(ctx, tx, r); err != nil {
			return nil, nil, err
		}
		states = append(states, r.State)
	}
	if err = s.checkHandoffQuota(ctx, tx, beforeQuota); err != nil {
		return nil, nil, err
	}
	return states, blocked, nil
}

func handoffDescendant(ownership CutoverOwnership, target, source GroupKey, seen map[GroupKey]bool) bool {
	if target == source {
		return true
	}
	if seen[target] {
		return false
	}
	seen[target] = true
	for parent := range ownership.predecessors[target] {
		if handoffDescendant(ownership, parent, source, seen) {
			return true
		}
	}
	return false
}

// handoffStorageError retains diagnostic identity for errors.Is without exposing
// SQL text, original sample values, filesystem paths or connector credentials.
type handoffStorageError struct{ cause error }

func (e *handoffStorageError) Error() string { return ErrOwnershipUnavailable.Error() }
func (e *handoffStorageError) Unwrap() error { return e.cause }
func (e *handoffStorageError) Is(target error) bool {
	return target == ErrOwnershipUnavailable || errors.Is(e.cause, target)
}
