package groupdelivery

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/snapshot"
)

func closureThrough(c Closure) time.Time {
	if c.Mixed != nil {
		return c.Mixed.After.Progress.ClosedThrough
	}
	return c.NextClose
}

func decodeMixedCheckpoint(raw, digest string) (snapshot.MixedCheckpoint, error) {
	var cp snapshot.MixedCheckpoint
	if handoffHash(raw) != digest || !counterproof.UniqueJSONObject([]byte(raw), false) {
		return cp, mixedInvalid("checkpoint digest")
	}
	dec := json.NewDecoder(bytes.NewBufferString(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cp); err != nil {
		return cp, mixedInvalid("checkpoint JSON")
	}
	return cp, nil
}
func restoreMixedProgress(ctx context.Context, q ownershipQuery, key GroupKey, r *Restored) error {
	d, digest, err := readMixedDescriptor(ctx, q, key)
	if err != nil {
		return err
	}
	var raw, hash, storedDescriptor, through string
	var version int
	err = q.QueryRowContext(ctx, `SELECT version,descriptor_digest,closed_through,checkpoint,checkpoint_digest FROM wg_delivery_mixed_progress WHERE workspace_id=? AND group_id=? AND group_revision=?`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&version, &storedDescriptor, &through, &raw, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		if d != nil {
			if !r.NextClose.IsZero() {
				return mixedInvalid("legacy cursor coexistence")
			}
			r.DescriptorDigest = digest
		}
		return nil
	}
	if err != nil {
		return err
	}
	if d == nil || !r.NextClose.IsZero() || storedDescriptor != digest || version != 1 {
		return mixedInvalid("descriptor or cursor conflict")
	}
	cp, err := decodeMixedCheckpoint(raw, hash)
	if err != nil {
		return err
	}
	intervals, err := d.intervals()
	if err != nil {
		return err
	}
	if err := validateMixedCheckpoint(key, cp, intervals); err != nil {
		return err
	}
	if cp.Progress.ClosedThrough.IsZero() || through != timestamp(cp.Progress.ClosedThrough) {
		return mixedInvalid("watermark mismatch")
	}
	r.Mixed = &cp
	r.DescriptorDigest = digest
	return nil
}
func (s *Store) commitMixedClosure(ctx context.Context, key GroupKey, c Closure) error {
	if !key.Valid() {
		return ErrInvalidGroupKey
	}
	if c.Mixed == nil || !c.NextClose.IsZero() || c.Mixed.DescriptorDigest == "" || c.Destination.Scope == "" || c.Destination.ConnectorID == "" || c.Destination.ConnectorRevision == "" || c.Destination.TableName == "" {
		return mixedInvalid("closure identity")
	}
	if !c.Mixed.After.Progress.ClosedThrough.After(c.Mixed.Before.Progress.ClosedThrough) || len(c.Buckets) == 0 {
		return mixedInvalid("no forward closure")
	}
	baselines, err := validateBaselines(key, closureThrough(c), c.Baselines)
	if err != nil {
		return err
	}
	rawClosure, err := json.Marshal(c)
	if err != nil {
		return err
	}
	closureDigest := handoffHash(string(rawClosure))
	err = s.inTx(ctx, "mixed-closure", func(ctx context.Context, tx *sql.Tx) error {
		d, digest, err := readMixedDescriptor(ctx, tx, key)
		if err != nil {
			return err
		}
		if d == nil || digest != c.Mixed.DescriptorDigest {
			return mixedInvalid("frozen descriptor changed")
		}
		intervals, err := d.intervals()
		if err != nil {
			return err
		}
		for _, cp := range []snapshot.MixedCheckpoint{c.Mixed.Before, c.Mixed.After} {
			if err := validateMixedCheckpoint(key, cp, intervals); err != nil {
				return err
			}
		}
		if !c.Mixed.Before.StartsAt.Equal(c.Mixed.After.StartsAt) || !c.Mixed.Before.EndsAt.Equal(c.Mixed.After.EndsAt) {
			return mixedInvalid("coverage changed within plan")
		}
		var legacy int
		if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM wg_delivery_checkpoints WHERE group_id=? AND group_revision=?`, key.GroupID, key.GroupRevision).Scan(&legacy); err != nil {
			return err
		}
		if legacy != 0 {
			return mixedInvalid("legacy cursor coexistence")
		}
		var prior Restored
		if err := restoreMixedProgress(ctx, tx, key, &prior); err != nil {
			return err
		}
		if prior.Mixed != nil {
			var last string
			if err := tx.QueryRowContext(ctx, `SELECT last_closure_digest FROM wg_delivery_mixed_progress WHERE workspace_id=? AND group_id=? AND group_revision=?`, key.WorkspaceID, key.GroupID, key.GroupRevision).Scan(&last); err != nil {
				return err
			}
			if reflect.DeepEqual(*prior.Mixed, c.Mixed.After) && last == closureDigest {
				return nil
			}
			before := c.Mixed.Before
			// A stop may narrow coverage between committed ticks, never extend it.
			if !prior.Mixed.EndsAt.Equal(before.EndsAt) {
				if before.EndsAt.IsZero() || (!prior.Mixed.EndsAt.IsZero() && before.EndsAt.After(prior.Mixed.EndsAt)) {
					return mixedInvalid("cutoff extension")
				}
				before.EndsAt = prior.Mixed.EndsAt
				before.Members = make(map[string]snapshot.MemberProgress, len(before.Members))
				for key, old := range prior.Mixed.Members {
					now, ok := c.Mixed.Before.Members[key]
					if !ok || !now.Start.Equal(old.Start) {
						return mixedInvalid("cutoff changed member phase")
					}
					before.Members[key] = old
				}
			}
			if !reflect.DeepEqual(*prior.Mixed, before) {
				return mixedInvalid("stale progress")
			}
		} else if !c.Mixed.Before.Progress.ClosedThrough.IsZero() {
			return mixedInvalid("missing prior progress")
		}
		if err := validateMixedSchedule(d, c, intervals); err != nil {
			return err
		}
		if err := validateMixedOutcomes(key, c, intervals); err != nil {
			return err
		}
		if err := validateOwnershipClosure(ctx, tx, key, c); err != nil {
			return err
		}
		if err := validateHandoffClosure(ctx, tx, key, c); err != nil {
			return err
		}
		if err := verifyClosureCounterProofs(ctx, tx, key, c.CounterProofs, closureThrough(c)); err != nil {
			return err
		}
		if err := verifyMixedBaselineProof(ctx, tx, key, c.Baselines); err != nil {
			return err
		}
		now := timestamp(s.now())
		if err := storeBaselines(ctx, tx, key, baselines, now); err != nil {
			return err
		}
		for _, b := range c.Buckets {
			raw, err := encodeMixedOutcome(b.Outcome)
			if err != nil {
				return err
			}
			if _, err = tx.ExecContext(ctx, `INSERT INTO wg_delivery_mixed_outcomes(workspace_id,group_id,group_revision,record_id,due_at,kind,outcome,outcome_digest,created_at) VALUES(?,?,?,?,?,?,?,?,?)`, key.WorkspaceID, key.GroupID, key.GroupRevision, b.Outcome.RecordID, timestamp(b.Outcome.DueAt), string(b.Outcome.Kind), string(raw), handoffHash(string(raw)), now); err != nil {
				return fmt.Errorf("mixed outcome: %w", err)
			}
			if b.Row != nil {
				if err := s.enqueueRow(ctx, tx, key, c.Destination, b, now); err != nil {
					return err
				}
			}
		}
		for member, after := range c.Mixed.After.Members {
			before := c.Mixed.Before.Members[member]
			if !after.Start.After(before.Start) {
				continue
			}
			end := after.Start
			if !c.Mixed.After.EndsAt.IsZero() && end.After(c.Mixed.After.EndsAt) {
				end = c.Mixed.After.EndsAt
			}
			if _, err = tx.ExecContext(ctx, `UPDATE wg_delivery_samples SET consumed=1 WHERE workspace_id=? AND group_id=? AND group_revision=? AND member_key=? AND consumed=0 AND observed_at < ? AND NOT EXISTS(SELECT 1 FROM wg_delivery_sample_handoffs h WHERE h.source_seq=wg_delivery_samples.seq)`, key.WorkspaceID, key.GroupID, key.GroupRevision, member, timestamp(end)); err != nil {
				return err
			}
		}
		raw, err := json.Marshal(c.Mixed.After)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO wg_delivery_mixed_progress(workspace_id,group_id,group_revision,version,descriptor_digest,closed_through,checkpoint,checkpoint_digest,last_closure_digest,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?) ON CONFLICT(workspace_id,group_id,group_revision) DO UPDATE SET descriptor_digest=excluded.descriptor_digest,closed_through=excluded.closed_through,checkpoint=excluded.checkpoint,checkpoint_digest=excluded.checkpoint_digest,last_closure_digest=excluded.last_closure_digest,updated_at=excluded.updated_at`, key.WorkspaceID, key.GroupID, key.GroupRevision, 1, digest, timestamp(closureThrough(c)), string(raw), handoffHash(string(raw)), closureDigest, now)
		return err
	})
	s.usage.invalidate()
	return s.capacityError(ctx, err)
}
