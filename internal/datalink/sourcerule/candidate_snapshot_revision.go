package sourcerule

import (
	"context"
	"database/sql"
	"fmt"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/schema"
)

// ReplaceCandidateSnapshotsAtRevision performs the revision check in the same
// transaction as the snapshot replacement, preventing an old build from
// replacing a newer source-rule revision.
func (r *SQLRepository) ReplaceCandidateSnapshotsAtRevision(ctx context.Context, snapshots []*schema.SourceRuleCandidateSnapshot, expectedRevision string) error {
	if len(snapshots) == 0 {
		return nil
	}
	ruleID := snapshots[0].SourceRuleID
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("開啟候選快照交易失敗: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return
		}
	}()
	var revision string
	if err := tx.QueryRowContext(ctx, `SELECT revision_id FROM source_rules WHERE id = ?`, ruleID).Scan(&revision); err != nil {
		return fmt.Errorf("讀取來源規則 revision 失敗: %w", err)
	}
	if revision != expectedRevision {
		return &modbusshare.Error{Code: modbusshare.ErrCodeRevisionConflict, Message: candidateRevisionConflictMessage, Retryable: true}
	}
	if err := replaceCandidateSnapshotsTx(ctx, tx, snapshots); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交候選快照交易失敗: %w", err)
	}
	return nil
}

func replaceCandidateSnapshotsTx(ctx context.Context, tx *sql.Tx, snapshots []*schema.SourceRuleCandidateSnapshot) error {
	if len(snapshots) == 0 {
		return nil
	}
	ruleID, revisionID := snapshots[0].SourceRuleID, snapshots[0].RevisionID
	if _, err := tx.ExecContext(ctx, `DELETE FROM source_rule_candidate_snapshots WHERE source_rule_id = ? AND revision_id = ?`, ruleID, revisionID); err != nil {
		return fmt.Errorf("清除既有候選快照失敗: %w", err)
	}
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO source_rule_candidate_snapshots (source_rule_id, revision_id, candidate_type, payload, status, reason, generated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("準備候選快照語句失敗: %w", err)
	}
	defer stmt.Close()
	for _, snapshot := range snapshots {
		if snapshot == nil {
			continue
		}
		if snapshot.SourceRuleID != ruleID || snapshot.RevisionID != revisionID {
			return fmt.Errorf("candidate snapshots must share the same source rule revision")
		}
		if _, err := stmt.ExecContext(ctx, snapshot.SourceRuleID, snapshot.RevisionID, snapshot.CandidateType, snapshot.Payload, snapshot.Status, snapshot.Reason, snapshot.GeneratedAt); err != nil {
			return fmt.Errorf("建立候選快照失敗: %w", err)
		}
	}
	return nil
}
