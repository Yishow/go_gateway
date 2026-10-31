package sourcerule

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sqltx"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, rule *schema.SourceRule) error {
	query := `
		INSERT INTO source_rules (
			id, device_id, start_address, count, data_type, naming_prefix,
			enabled, locked, origin, template_name, skipped_addresses,
			target_data_type, scale_multiplier, scale_offset, data_format,
			share_enabled, share_start_register, share_stride, revision_id,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var targetDataType any
	if rule.TargetDataType != nil {
		targetDataType = *rule.TargetDataType
	}

	var scaleMultiplier any
	if rule.ScaleMultiplier != nil {
		scaleMultiplier = *rule.ScaleMultiplier
	}

	var scaleOffset any
	if rule.ScaleOffset != nil {
		scaleOffset = *rule.ScaleOffset
	}

	var shareStartRegister any
	if rule.ShareStartRegister != nil {
		shareStartRegister = *rule.ShareStartRegister
	}

	var shareStride any
	if rule.ShareStride != nil {
		shareStride = *rule.ShareStride
	}

	_, err := r.db.ExecContext(ctx, query,
		rule.ID,
		rule.DeviceID,
		rule.StartAddress,
		rule.Count,
		rule.DataType,
		rule.NamingPrefix,
		rule.Enabled,
		rule.Locked,
		rule.Origin,
		rule.TemplateName,
		rule.SkippedAddresses,
		targetDataType,
		scaleMultiplier,
		scaleOffset,
		rule.DataFormat,
		rule.ShareEnabled,
		shareStartRegister,
		shareStride,
		rule.RevisionID,
		rule.CreatedAt,
		rule.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("建立來源規則失敗: %w", err)
	}
	return nil
}

func (r *SQLRepository) Update(ctx context.Context, rule *schema.SourceRule) error {
	return r.updateRuleRevision(ctx, rule, nil)
}

func (r *SQLRepository) updateRuleRevision(ctx context.Context, rule *schema.SourceRule, expected *schema.SourceRule) error {
	query := `
		UPDATE source_rules
		SET start_address = ?, count = ?, data_type = ?, naming_prefix = ?,
		    enabled = ?, locked = ?, origin = ?, template_name = ?, skipped_addresses = ?,
		    target_data_type = ?, scale_multiplier = ?, scale_offset = ?, data_format = ?,
		    share_enabled = ?, share_start_register = ?, share_stride = ?,
		    updated_at = ?, revision_id = ?
		WHERE id = ?
	`

	var targetDataType any
	if rule.TargetDataType != nil {
		targetDataType = *rule.TargetDataType
	}

	var scaleMultiplier any
	if rule.ScaleMultiplier != nil {
		scaleMultiplier = *rule.ScaleMultiplier
	}

	var scaleOffset any
	if rule.ScaleOffset != nil {
		scaleOffset = *rule.ScaleOffset
	}

	var shareStartRegister any
	if rule.ShareStartRegister != nil {
		shareStartRegister = *rule.ShareStartRegister
	}

	var shareStride any
	if rule.ShareStride != nil {
		shareStride = *rule.ShareStride
	}

	args := []any{
		rule.StartAddress,
		rule.Count,
		rule.DataType,
		rule.NamingPrefix,
		rule.Enabled,
		rule.Locked,
		rule.Origin,
		rule.TemplateName,
		rule.SkippedAddresses,
		targetDataType,
		scaleMultiplier,
		scaleOffset,
		rule.DataFormat,
		rule.ShareEnabled,
		shareStartRegister,
		shareStride,
		rule.UpdatedAt,
		rule.RevisionID,
		rule.ID,
	}
	if expected != nil {
		var rawTimestamp string
		if err := r.db.QueryRowContext(ctx, `SELECT CAST(updated_at AS TEXT) FROM source_rules WHERE id=? AND device_id=? AND revision_id=?`, expected.ID, expected.DeviceID, expected.RevisionID).Scan(&rawTimestamp); err != nil {
			return ErrRuleRevisionConflict
		}
		parsed, err := common.ParseTimeString(rawTimestamp)
		if err != nil || !parsed.Equal(expected.UpdatedAt) {
			return ErrRuleRevisionConflict
		}
		query += " AND device_id = ? AND revision_id = ? AND updated_at = ?"
		args = append(args, expected.DeviceID, expected.RevisionID, rawTimestamp)
		predicate, stateArgs := ruleStatePredicate(expected)
		query += " AND " + predicate
		args = append(args, stateArgs...)
	}
	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("更新來源規則失敗: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得更新影響列數失敗: %w", err)
	}
	if rows == 0 {
		if expected != nil {
			return ErrRuleRevisionConflict
		}
		return fmt.Errorf("%w: %s", ErrSourceRuleNotFound, rule.ID)
	}
	return nil
}

func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM source_rules WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("刪除來源規則失敗: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得刪除影響列數失敗: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: %s", ErrSourceRuleNotFound, id)
	}
	return nil
}

func (r *SQLRepository) GetByID(ctx context.Context, id string) (*schema.SourceRule, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, device_id, start_address, count, data_type, naming_prefix,
		       enabled, locked, origin, template_name, skipped_addresses,
		       target_data_type, scale_multiplier, scale_offset, data_format,
		       share_enabled, share_start_register, share_stride,
		       created_at, updated_at, revision_id
		FROM source_rules WHERE id = ?
	`, id)
	return scanRule(row)
}

func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]*schema.SourceRule, error) {
	return readShareList(ctx, r.db, filter)
}

func (r *SQLRepository) CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) (err error) {
	if len(links) == 0 {
		return nil
	}

	ownedTx, err := sqltx.Begin(ctx, r.db, nil)
	if err != nil {
		return fmt.Errorf("開啟來源規則連結交易失敗: %w", err)
	}
	tx := ownedTx.Tx
	defer func() {
		if rollbackErr := ownedTx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			err = errors.Join(err, fmt.Errorf("回滾來源規則連結交易失敗: %w", rollbackErr))
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO source_rule_links (
			id, rule_id, address, point_id, tag_id, mapping_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("準備來源規則連結語句失敗: %w", err)
	}
	defer stmt.Close()

	for _, link := range links {
		_, err = stmt.ExecContext(ctx,
			link.ID,
			link.RuleID,
			link.Address,
			link.PointID,
			link.TagID,
			link.MappingID,
			link.CreatedAt,
			link.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("建立來源規則連結失敗: %w", err)
		}
	}

	if err := ownedTx.Commit(); err != nil {
		return fmt.Errorf("提交來源規則連結交易失敗: %w", err)
	}
	return nil
}

func (r *SQLRepository) ListLinks(ctx context.Context, ruleID string) ([]*schema.SourceRuleLink, error) {
	return readShareListLinks(ctx, r.db, ruleID)
}

func (r *SQLRepository) DeleteLinks(ctx context.Context, ruleID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM source_rule_links WHERE rule_id = ?`, ruleID)
	if err != nil {
		return fmt.Errorf("刪除來源規則連結失敗: %w", err)
	}
	return nil
}

func (r *SQLRepository) ReplaceCandidateSnapshots(ctx context.Context, snapshots []*schema.SourceRuleCandidateSnapshot) (err error) {
	if len(snapshots) == 0 {
		return nil
	}

	ruleID := snapshots[0].SourceRuleID
	revisionID := snapshots[0].RevisionID
	ownedTx, err := sqltx.Begin(ctx, r.db, nil)
	if err != nil {
		return fmt.Errorf("開啟候選快照交易失敗: %w", err)
	}
	tx := ownedTx.Tx
	defer func() {
		if rollbackErr := ownedTx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			err = errors.Join(err, fmt.Errorf("回滾候選快照交易失敗: %w", rollbackErr))
		}
	}()

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM source_rule_candidate_snapshots
		WHERE source_rule_id = ? AND revision_id = ?
	`, ruleID, revisionID); err != nil {
		return fmt.Errorf("清除既有候選快照失敗: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO source_rule_candidate_snapshots (
			source_rule_id, revision_id, candidate_type, payload, status, reason, generated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
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
		if _, err := stmt.ExecContext(
			ctx,
			snapshot.SourceRuleID,
			snapshot.RevisionID,
			snapshot.CandidateType,
			snapshot.Payload,
			snapshot.Status,
			snapshot.Reason,
			snapshot.GeneratedAt,
		); err != nil {
			return fmt.Errorf("建立候選快照失敗: %w", err)
		}
	}

	if err := ownedTx.Commit(); err != nil {
		return fmt.Errorf("提交候選快照交易失敗: %w", err)
	}
	return nil
}

func (r *SQLRepository) ListCandidateSnapshots(ctx context.Context, ruleID, revisionID string) ([]*schema.SourceRuleCandidateSnapshot, error) {
	return readShareListCandidateSnapshots(ctx, r.db, ruleID, revisionID)
}

func (r *SQLRepository) DeleteCandidateSnapshots(ctx context.Context, ruleID, revisionID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM source_rule_candidate_snapshots
		WHERE source_rule_id = ? AND revision_id = ?
	`, ruleID, revisionID)
	if err != nil {
		return fmt.Errorf("刪除候選快照失敗: %w", err)
	}
	return nil
}

func scanLinkRows(rows *sql.Rows) (*schema.SourceRuleLink, error) {
	var link schema.SourceRuleLink
	var tagID sql.NullString
	var mappingID sql.NullString
	var createdAt string
	var updatedAt string
	err := rows.Scan(
		&link.ID,
		&link.RuleID,
		&link.Address,
		&link.PointID,
		&tagID,
		&mappingID,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("掃描來源規則連結失敗: %w", err)
	}
	if tagID.Valid {
		value := tagID.String
		link.TagID = &value
	}
	if mappingID.Valid {
		value := mappingID.String
		link.MappingID = &value
	}
	link.CreatedAt, err = common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析來源規則連結建立時間失敗: %w", err)
	}
	link.UpdatedAt, err = common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析來源規則連結更新時間失敗: %w", err)
	}
	return &link, nil
}

func scanCandidateSnapshot(rows *sql.Rows) (*schema.SourceRuleCandidateSnapshot, error) {
	var snapshot schema.SourceRuleCandidateSnapshot
	var generatedAt string
	if err := rows.Scan(
		&snapshot.SourceRuleID,
		&snapshot.RevisionID,
		&snapshot.CandidateType,
		&snapshot.Payload,
		&snapshot.Status,
		&snapshot.Reason,
		&generatedAt,
	); err != nil {
		return nil, fmt.Errorf("掃描候選快照失敗: %w", err)
	}

	parsedTime, err := common.ParseTimeString(generatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析候選快照建立時間失敗: %w", err)
	}
	snapshot.GeneratedAt = parsedTime
	return &snapshot, nil
}
