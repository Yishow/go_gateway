package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/sqltx"
)

// HandoffView is a Store-verified snapshot of transfer dispositions.
// Only Restore can construct a nonzero view; callers preserve it unchanged.
type HandoffView struct {
	key    GroupKey
	digest string
}

type handoffObservationKey struct {
	sampleID, member string
	observedAt       time.Time
}

type handoffAuthorizationKey struct {
	origin   snapshot.AcceptanceOrigin
	sampleID string
}

func handoffSampleOrigin(key GroupKey) snapshot.AcceptanceOrigin {
	return snapshot.AcceptanceOrigin{WorkspaceID: key.WorkspaceID, GroupID: key.GroupID, GroupRevision: key.GroupRevision}
}

type handoffAuthorization struct {
	member     string
	observedAt time.Time
	digest     string
	target     counterproof.CounterBindingReference
}

// AuthorizesHandoff checks private durable evidence without changing the sample's
// original group-bound proof. It is not an authorization to append new input.
func (r Restored) AuthorizesHandoff(sample snapshot.Sample, target counterproof.CounterBindingReference) bool {
	_, digest, err := encodeSample(sample)
	auth, ok := r.handoffs[handoffAuthorizationKey{sample.AcceptanceOrigin, sample.SampleID}]
	return err == nil && ok && auth.digest == digest && auth.target == target
}

// ValidateHandoffOrigin prevents an erased annotation from borrowing a native
// legacy association. Observation identity is checked independently of payload
// digest so changing a legacy value cannot disguise a known transfer.
func (r Restored) ValidateHandoffOrigin(sample snapshot.Sample, target counterproof.CounterBindingReference) bool {
	if sample.AcceptanceOrigin != (snapshot.AcceptanceOrigin{}) {
		return r.AuthorizesHandoff(sample, target)
	}
	return !r.handoffObservations[handoffObservationKey{sample.SampleID, sample.MemberKey, sample.ObservedAt.UTC()}]
}

// HandoffRecognition distinguishes a durable resend from a new observation.
// A finalized receipt may acknowledge identity after normal payload reclamation;
// Sample is then nil and must never be synthesized or replayed.
type HandoffRecognition struct {
	Found     bool
	Finalized bool
	Sample    *snapshot.Sample
}

func (s *Store) RecognizeTransferredSample(ctx context.Context, key GroupKey, incoming snapshot.Sample) (HandoffRecognition, error) {
	if !key.Valid() {
		return HandoffRecognition{}, ErrInvalidGroupKey
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return HandoffRecognition{}, ErrOwnershipUnavailable
	}
	defer tx.Rollback()
	rs, err := queryHandoffs(ctx, tx, "workspace_id=? AND target_group_id=? AND target_revision=? AND json_extract(record,'$.SampleID')=?", key.WorkspaceID, key.GroupID, key.GroupRevision, incoming.SampleID)
	if err != nil {
		return HandoffRecognition{}, err
	}
	for _, r := range rs {
		at, e := time.Parse(time.RFC3339Nano, r.ObservedAt)
		if e != nil {
			return HandoffRecognition{}, ErrOwnershipUnavailable
		}
		if r.SampleID != incoming.SampleID || r.Member != incoming.MemberKey || !at.Equal(incoming.ObservedAt) {
			continue
		}
		if e = verifyHandoffReceipt(ctx, tx, r, map[string]bool{}); e != nil {
			return HandoffRecognition{}, e
		}
		if r.State != "accepted" && r.State != "finalized" {
			return HandoffRecognition{}, ErrOwnershipUnavailable
		}
		if r.State == "accepted" {
			disposed, e := handoffDisposed(ctx, tx, r.TargetSeq)
			if e != nil {
				return HandoffRecognition{}, e
			}
			if disposed {
				return HandoffRecognition{}, ErrOwnershipConflict
			}
			own, e := resolveCutoverOwnership(ctx, tx, sourceOf(key, incoming), incoming.ObservedAt, key)
			if e != nil {
				return HandoffRecognition{}, e
			}
			if own.Affected[key] && !own.Owners[key] {
				return HandoffRecognition{}, ErrOwnershipConflict
			}
		}
		if incoming.AcceptanceOrigin != (snapshot.AcceptanceOrigin{}) && incoming.AcceptanceOrigin != handoffSampleOrigin(r.Origin) {
			return HandoffRecognition{}, ErrOwnershipConflict
		}
		if incoming.CounterProof != r.SampleProof {
			if e = verifySampleProof(ctx, tx, key, incoming, true); e != nil {
				return HandoffRecognition{}, ErrHandoffProofConflict
			}
		}
		identity := incoming
		identity.CounterProof = r.SampleProof // identity comparison only; never journal/replay this rewritten value
		_, digest, e := encodeSample(identity)
		if e != nil || digest != r.Digest {
			return HandoffRecognition{}, ErrHandoffProofConflict
		}
		result := HandoffRecognition{Found: true, Finalized: r.State == "finalized"}
		if !result.Finalized {
			row, e := readHandoffSample(ctx, tx, r.TargetSeq)
			if e != nil {
				return HandoffRecognition{}, e
			}
			if row.Consumed {
				result.Finalized = true
				if tx.Commit() != nil {
					return HandoffRecognition{}, ErrOwnershipUnavailable
				}
				return result, nil
			}
			original, e := decodeVerifiedSample(row.Payload, row.Digest)
			if e != nil {
				return HandoffRecognition{}, e
			}
			original.AcceptanceOrigin = handoffSampleOrigin(r.Origin)
			result.Sample = &original
		}
		if tx.Commit() != nil {
			return HandoffRecognition{}, ErrOwnershipUnavailable
		}
		return result, nil
	}
	// A verified miss authorizes nothing. The ordinary Append transaction is
	// decisive; avoid introducing an unrelated commit acknowledgement here.
	return HandoffRecognition{}, nil
}

func (s *Store) restoreHandoffView(ctx context.Context, key GroupKey) (result Restored, retErr error) {
	if !key.Valid() {
		return Restored{}, ErrInvalidGroupKey
	}
	owned, err := sqltx.Begin(ctx, s.db, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return Restored{}, ErrOwnershipUnavailable
	}
	tx := owned.Tx
	defer func() { //nolint:contextcheck // sqltx terminal cleanup uses its own bounded context so caller cancellation cannot leak a transaction.
		if e := owned.Rollback(); e != nil && !errors.Is(e, sql.ErrTxDone) {
			retErr = errors.Join(retErr, e)
			result = Restored{}
		}
	}()
	r := Restored{handoffs: map[handoffAuthorizationKey]handoffAuthorization{}, handoffObservations: map[handoffObservationKey]bool{}}
	view, _, err := handoffView(ctx, tx, key)
	if err != nil {
		return r, err
	}
	r.HandoffView = view
	r.OwnershipView, err = captureOwnershipView(ctx, tx, key)
	if err != nil {
		return Restored{}, err
	}
	var next string
	err = tx.QueryRowContext(ctx, `SELECT next_close FROM wg_delivery_checkpoints WHERE group_id=? AND group_revision=?`, key.GroupID, key.GroupRevision).Scan(&next)
	if err == nil {
		r.NextClose, err = time.Parse(time.RFC3339Nano, next)
	} else if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	if err != nil {
		return Restored{}, err
	}
	if err := restoreMixedProgress(ctx, tx, key, &r); err != nil {
		return Restored{}, err
	}
	rows, err := tx.QueryContext(ctx, `SELECT seq FROM wg_delivery_samples WHERE group_id=? AND group_revision=? AND consumed=0 ORDER BY seq`, key.GroupID, key.GroupRevision)
	if err != nil {
		return Restored{}, &handoffStorageError{cause: err}
	}
	var seqs []int64
	for rows.Next() {
		var seq int64
		if rows.Scan(&seq) != nil {
			rows.Close()
			return Restored{}, ErrOwnershipUnavailable
		}
		seqs = append(seqs, seq)
	}
	if rows.Err() != nil {
		rows.Close()
		return Restored{}, ErrOwnershipUnavailable
	}
	if rows.Close() != nil {
		return Restored{}, ErrOwnershipUnavailable
	}
	for _, seq := range seqs {
		disposed, e := handoffDisposed(ctx, tx, seq)
		if e != nil {
			return Restored{}, e
		}
		if disposed {
			continue
		}
		row, e := readHandoffSample(ctx, tx, seq)
		if e != nil {
			return Restored{}, e
		}
		if row.Key != key {
			return Restored{}, sampleProofError(counterproof.ErrCounterBindingProofConflict)
		}
		sample, e := decodeVerifiedSample(row.Payload, row.Digest)
		if e != nil {
			return Restored{}, e
		}
		incoming, e := incomingHandoff(ctx, tx, seq)
		if e != nil {
			return Restored{}, e
		}
		if incoming != nil {
			if e = verifyHandoffReceipt(ctx, tx, *incoming, map[string]bool{}); e != nil {
				return Restored{}, e
			}
			sample.AcceptanceOrigin = handoffSampleOrigin(incoming.Origin)
			r.handoffObservations[handoffObservationKey{sample.SampleID, sample.MemberKey, sample.ObservedAt.UTC()}] = true
			r.handoffs[handoffAuthorizationKey{sample.AcceptanceOrigin, sample.SampleID}] = handoffAuthorization{digest: row.Digest, target: incoming.TargetProof, member: sample.MemberKey, observedAt: sample.ObservedAt}
		} else {
			if row.Kind != "native" {
				return Restored{}, ErrOwnershipUnavailable
			}
			if e = verifySampleProof(ctx, tx, key, sample, false); e != nil {
				return Restored{}, e
			}
			ownership, e := resolveCutoverOwnership(ctx, tx, sourceOf(key, sample), sample.ObservedAt, key)
			if e != nil {
				return Restored{}, e
			}
			if e = rejectAcceptedPredecessor(ctx, tx, ownership, key, sourceOf(key, sample), sample); e != nil {
				return Restored{}, e
			}
		}
		r.Samples = append(r.Samples, sample)
	}
	// Baselines and sample/disposition/checkpoint evidence share this snapshot.
	rows, err = tx.QueryContext(ctx, `SELECT member_key,payload,payload_digest,workspace_id,window_end FROM wg_delivery_counter_baselines WHERE group_id=? AND group_revision=? ORDER BY member_key`, key.GroupID, key.GroupRevision)
	if err != nil {
		return Restored{}, err
	}
	for rows.Next() {
		var member, payload, digest, workspace, end string
		if rows.Scan(&member, &payload, &digest, &workspace, &end) != nil {
			rows.Close()
			return Restored{}, ErrBaselineCorrupt
		}
		at, e := time.Parse(time.RFC3339Nano, end)
		if handoffHash(payload) != digest || workspace != key.WorkspaceID || e != nil || at.After(restoredThrough(r)) {
			rows.Close()
			return Restored{}, ErrBaselineCorrupt
		}
		r.Baselines = append(r.Baselines, CounterBaseline{MemberKey: member, Payload: []byte(payload)})
	}
	if rows.Err() != nil {
		rows.Close()
		return Restored{}, ErrBaselineCorrupt
	}
	if rows.Close() != nil {
		return Restored{}, ErrBaselineCorrupt
	}
	if r.Mixed != nil {
		if err := verifyMixedBaselineProof(ctx, tx, key, r.Baselines); err != nil {
			return Restored{}, err
		}
	}
	if err := owned.Commit(); err != nil { //nolint:contextcheck // sqltx terminal cleanup uses its own bounded context so caller cancellation cannot leak a transaction.
		return Restored{}, errors.Join(ErrOwnershipUnavailable, err)
	}
	return r, nil
}

func refuseOrdinaryHandoffAppend(ctx context.Context, q ownershipQuery, key GroupKey, sample snapshot.Sample) error {
	rs, err := queryHandoffs(ctx, q, "workspace_id=? AND target_group_id=? AND target_revision=? AND json_extract(record,'$.SampleID')=?", key.WorkspaceID, key.GroupID, key.GroupRevision, sample.SampleID)
	if err != nil {
		return err
	}
	for _, r := range rs {
		at, e := time.Parse(time.RFC3339Nano, r.ObservedAt)
		if e != nil {
			return ErrOwnershipUnavailable
		}
		if r.SampleID == sample.SampleID && r.Member == sample.MemberKey && at.Equal(sample.ObservedAt) {
			if e = verifyHandoffReceipt(ctx, q, r, map[string]bool{}); e != nil {
				return e
			}
			return ErrOwnershipConflict
		}
	}
	return nil
}
