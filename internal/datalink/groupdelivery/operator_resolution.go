package groupdelivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidDecision  = errors.New("invalid delivery decision")
	ErrDecisionConflict = errors.New("delivery decision conflicts with current state")
)

// AttentionItem excludes values and credentials; identity stays bound to the accepted row.
type AttentionItem struct {
	EffectKey         string `json:"effect_key"`
	State             string `json:"state"`
	StateRevision     int64  `json:"state_revision"`
	PayloadDigest     string `json:"payload_digest"`
	GroupRevision     string `json:"group_revision"`
	ConnectorID       string `json:"connector_id"`
	ConnectorRevision string `json:"connector_revision"`
	TableSchema       string `json:"table_schema"`
	TableName         string `json:"table_name"`
	ErrorCode         string `json:"error_code"`
}

// Attention lists a bounded, oldest-first set; unknown is visible but never resolvable here.
func (s *Store) Attention(ctx context.Context, groupID string) ([]AttentionItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT effect_key,state,payload_digest,group_revision,connector_id,connector_revision,table_schema,table_name,last_error_code,claim_epoch FROM wg_delivery_outbox WHERE group_id=? AND state IN ('blocked','quarantined','unknown') ORDER BY bucket_start,effect_key LIMIT 20`, groupID)
	if err != nil {
		return nil, fmt.Errorf("read delivery attention: %w", err)
	}
	defer rows.Close()
	items := []AttentionItem{}
	for rows.Next() {
		var item AttentionItem
		if err := rows.Scan(&item.EffectKey, &item.State, &item.PayloadDigest, &item.GroupRevision, &item.ConnectorID, &item.ConnectorRevision, &item.TableSchema, &item.TableName, &item.ErrorCode, &item.StateRevision); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// OperatorDecision is a CAS decision, scoped by the trusted current-workspace group.
type OperatorDecision struct {
	DecisionID            string               `json:"decision_id"`
	EffectKey             string               `json:"effect_key"`
	ExpectedState         string               `json:"expected_state"`
	ExpectedStateRevision *int64               `json:"expected_state_revision"`
	PayloadDigest         string               `json:"payload_digest"`
	Resolution            QuarantineResolution `json:"resolution"`
	Reason                string               `json:"reason"`
	ConfirmSkip           bool                 `json:"confirm_skip"`
}

type DecisionResult struct {
	DecisionID string `json:"decision_id"`
	EffectKey  string `json:"effect_key"`
	State      string `json:"state"`
	Duplicate  bool   `json:"duplicate"`
}

type decisionAudit struct {
	OperatorDecision
	Actor         string            `json:"actor"`
	GroupID       string            `json:"group_id"`
	GroupRevision string            `json:"group_revision"`
	Destination   FrozenDestination `json:"frozen_destination"`
}

// ResolveAttention commits state and audit in one local transaction. A repeated
// decision ID with identical content is safe even after delivery has continued.
// Unknown effects require destination evidence, outside this repair/skip API.
func (s *Store) ResolveAttention(ctx context.Context, workspaceID, groupID string, d OperatorDecision) (*DecisionResult, error) {
	if strings.TrimSpace(d.DecisionID) != d.DecisionID || d.DecisionID == "" || len(d.DecisionID) > 128 ||
		d.ExpectedStateRevision == nil || (d.ExpectedStateRevision != nil && *d.ExpectedStateRevision < 0) ||
		d.EffectKey == "" || d.PayloadDigest == "" || strings.TrimSpace(d.Reason) == "" || utf8.RuneCountInString(d.Reason) > 200 ||
		(d.Resolution != ResolutionRetry && d.Resolution != ResolutionSkip) || (d.Resolution == ResolutionSkip && !d.ConfirmSkip) {
		return nil, ErrInvalidDecision
	}
	result := &DecisionResult{DecisionID: d.DecisionID, EffectKey: d.EffectKey}
	err := s.inTx(ctx, "operator delivery resolution", func(tx *sql.Tx) error {
		item, err := scanOutbox(tx.QueryRowContext(ctx, `SELECT `+outboxColumns+` FROM wg_delivery_outbox WHERE effect_key=? AND workspace_id=? AND group_id=?`, d.EffectKey, workspaceID, groupID))
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOutboxItemNotFound
		}
		if err != nil {
			return err
		}
		details, err := json.Marshal(decisionAudit{OperatorDecision: d, Actor: "local_api", GroupID: groupID, GroupRevision: item.Key.GroupRevision, Destination: item.Destination})
		if err != nil {
			return err
		}
		var oldDetails, oldWorkspace, eventType string
		err = tx.QueryRowContext(ctx, `SELECT details,workspace_id,event_type FROM workspace_audit_history WHERE id=?`, d.DecisionID).Scan(&oldDetails, &oldWorkspace, &eventType)
		if err == nil {
			if oldDetails != string(details) || oldWorkspace != workspaceID || eventType != "write_group_delivery_resolution" {
				return ErrDecisionConflict
			}
			result.State = item.State
			result.Duplicate = true
			return nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if (item.State != StateBlocked && item.State != StateQuarantined) || d.ExpectedState != item.State || d.PayloadDigest != item.PayloadDigest || *d.ExpectedStateRevision != item.ClaimEpoch {
			return ErrDecisionConflict
		}
		result.State = StatePending
		set := `state='pending',retry_count=0,next_retry_at=?,last_error_code='',claim_owner='',claim_epoch=claim_epoch+1,updated_at=?`
		args := make([]any, 0, 8)
		args = append(args, timestamp(s.now()), timestamp(s.now()))
		if d.Resolution == ResolutionSkip {
			result.State = StateSkipped
			set = `state='operator_skipped',claim_owner='',claim_epoch=claim_epoch+1,updated_at=?`
			args = []any{timestamp(s.now())}
		}
		args = append(args, d.EffectKey, workspaceID, groupID, d.ExpectedState, d.PayloadDigest, *d.ExpectedStateRevision)
		updated, err := tx.ExecContext(ctx, `UPDATE wg_delivery_outbox SET `+set+` WHERE effect_key=? AND workspace_id=? AND group_id=? AND state=? AND payload_digest=? AND claim_epoch=?`, args...)
		if err != nil {
			return err
		}
		affected, err := updated.RowsAffected()
		if err != nil {
			return err
		}
		if affected != 1 {
			return ErrDecisionConflict
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO workspace_audit_history(id,workspace_id,event_type,result,scope,reference_id,details,occurred_at,created_at) VALUES(?,?,'write_group_delivery_resolution','success','write_group',?,?,?,?)`, d.DecisionID, workspaceID, d.EffectKey, string(details), s.now(), s.now())
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
