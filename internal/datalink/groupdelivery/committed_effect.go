package groupdelivery

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// CommittedEffect identifies a delivery confirmed for one applied revision.
// It does not imply a separate readback verification.
type CommittedEffect struct {
	GroupRevision     string    `json:"group_revision"`
	ConnectorRevision string    `json:"connector_revision"`
	RecordID          string    `json:"record_id"`
	EffectKey         string    `json:"effect_key"`
	PayloadDigest     string    `json:"payload_digest"`
	CommittedAt       time.Time `json:"committed_at"`
}

// RevisionStageCounts excludes accepted samples and rows from older revisions.
func (s *Store) RevisionStageCounts(ctx context.Context, groupID, revision string) (StageCounts, error) {
	if revision == "" {
		return StageCounts{}, nil
	}
	status := GroupStatus{GroupID: groupID}
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM wg_delivery_samples
		WHERE group_id = ? AND group_revision = ? AND consumed = 0`, groupID, revision).Scan(&status.Stages.Collecting); err != nil {
		return StageCounts{}, fmt.Errorf("count revision collecting samples: %w", err)
	}
	if err := s.readRevisionStageCounts(ctx, groupID, revision, &status); err != nil {
		return StageCounts{}, err
	}
	return status.Stages, nil
}

// LastCommittedEffect reads actual receipt evidence for the requested revision.
// An empty revision never selects historical data.
func (s *Store) LastCommittedEffect(ctx context.Context, groupID, revision string) (*CommittedEffect, error) {
	if revision == "" {
		return nil, nil
	}
	var effect CommittedEffect
	var committed string
	err := s.db.QueryRowContext(ctx, `
		SELECT o.group_revision, o.connector_revision, o.record_id, o.effect_key, o.payload_digest, o.committed_at
		FROM wg_delivery_outbox o
		JOIN wg_delivery_receipts r ON r.effect_key = o.effect_key AND r.payload_digest = o.payload_digest
		WHERE o.group_id = ? AND o.group_revision = ? AND o.state = 'sql_committed'
		  AND o.committed_at != '' AND r.committed_at = o.committed_at
		ORDER BY o.committed_at DESC, o.effect_key DESC LIMIT 1`, groupID, revision).Scan(
		&effect.GroupRevision, &effect.ConnectorRevision, &effect.RecordID, &effect.EffectKey, &effect.PayloadDigest, &committed)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read committed effect: %w", err)
	}
	effect.CommittedAt, err = time.Parse(time.RFC3339Nano, committed)
	if err != nil {
		return nil, fmt.Errorf("parse committed effect timestamp: %w", err)
	}
	effect.CommittedAt = effect.CommittedAt.UTC()
	return &effect, nil
}
