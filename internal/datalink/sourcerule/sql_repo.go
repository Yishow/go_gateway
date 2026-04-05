package sourcerule

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
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
			target_data_type, scale_multiplier, scale_offset, data_format, revision_id,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var targetDataType interface{}
	if rule.TargetDataType != nil {
		targetDataType = *rule.TargetDataType
	}

	var scaleMultiplier interface{}
	if rule.ScaleMultiplier != nil {
		scaleMultiplier = *rule.ScaleMultiplier
	}

	var scaleOffset interface{}
	if rule.ScaleOffset != nil {
		scaleOffset = *rule.ScaleOffset
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
	query := `
		UPDATE source_rules
		SET start_address = ?, count = ?, data_type = ?, naming_prefix = ?,
		    enabled = ?, locked = ?, origin = ?, template_name = ?, skipped_addresses = ?,
		    target_data_type = ?, scale_multiplier = ?, scale_offset = ?, data_format = ?,
		    updated_at = ?, revision_id = ?
		WHERE id = ?
	`

	var targetDataType interface{}
	if rule.TargetDataType != nil {
		targetDataType = *rule.TargetDataType
	}

	var scaleMultiplier interface{}
	if rule.ScaleMultiplier != nil {
		scaleMultiplier = *rule.ScaleMultiplier
	}

	var scaleOffset interface{}
	if rule.ScaleOffset != nil {
		scaleOffset = *rule.ScaleOffset
	}

	result, err := r.db.ExecContext(ctx, query,
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
		rule.UpdatedAt,
		rule.RevisionID,
		rule.ID,
	)
	if err != nil {
		return fmt.Errorf("更新來源規則失敗: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("取得更新影響列數失敗: %w", err)
	}
	if rows == 0 {
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
		       created_at, updated_at, revision_id
		FROM source_rules WHERE id = ?
	`, id)
	return scanRule(row)
}

func (r *SQLRepository) List(ctx context.Context, filter ListFilter) ([]*schema.SourceRule, error) {
	query := `
		SELECT id, device_id, start_address, count, data_type, naming_prefix,
		       enabled, locked, origin, template_name, skipped_addresses,
		       target_data_type, scale_multiplier, scale_offset, data_format,
		       created_at, updated_at, revision_id
		FROM source_rules
		WHERE 1 = 1
	`
	args := make([]interface{}, 0, 2)

	if filter.DeviceID != nil {
		query += ` AND device_id = ?`
		args = append(args, *filter.DeviceID)
	}
	if filter.Enabled != nil {
		query += ` AND enabled = ?`
		args = append(args, *filter.Enabled)
	}
	query += ` ORDER BY created_at ASC, id ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("查詢來源規則失敗: %w", err)
	}
	defer rows.Close()

	items := make([]*schema.SourceRule, 0)
	for rows.Next() {
		rule, scanErr := scanRuleRows(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("讀取來源規則結果失敗: %w", err)
	}
	return items, nil
}

func (r *SQLRepository) CreateLinks(ctx context.Context, links []*schema.SourceRuleLink) error {
	if len(links) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("開啟來源規則連結交易失敗: %w", err)
	}
	defer tx.Rollback()

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

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交來源規則連結交易失敗: %w", err)
	}
	return nil
}

func (r *SQLRepository) ListLinks(ctx context.Context, ruleID string) ([]*schema.SourceRuleLink, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, rule_id, address, point_id, tag_id, mapping_id, created_at, updated_at
		FROM source_rule_links
		WHERE rule_id = ?
		ORDER BY address ASC, id ASC
	`, ruleID)
	if err != nil {
		return nil, fmt.Errorf("查詢來源規則連結失敗: %w", err)
	}
	defer rows.Close()

	items := make([]*schema.SourceRuleLink, 0)
	for rows.Next() {
		link, scanErr := scanLinkRows(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("讀取來源規則連結結果失敗: %w", err)
	}
	return items, nil
}

func (r *SQLRepository) DeleteLinks(ctx context.Context, ruleID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM source_rule_links WHERE rule_id = ?`, ruleID)
	if err != nil {
		return fmt.Errorf("刪除來源規則連結失敗: %w", err)
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanRule(row rowScanner) (*schema.SourceRule, error) {
	var rule schema.SourceRule
	var templateName sql.NullString
	var skipped sql.NullString
	var targetDataType sql.NullString
	var scaleMultiplier sql.NullFloat64
	var scaleOffset sql.NullFloat64
	var dataFormat sql.NullString
	var createdAt string
	var updatedAt string
	err := row.Scan(
		&rule.ID,
		&rule.DeviceID,
		&rule.StartAddress,
		&rule.Count,
		&rule.DataType,
		&rule.NamingPrefix,
		&rule.Enabled,
		&rule.Locked,
		&rule.Origin,
		&templateName,
		&skipped,
		&targetDataType,
		&scaleMultiplier,
		&scaleOffset,
		&dataFormat,
		&createdAt,
		&updatedAt,
		&rule.RevisionID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSourceRuleNotFound
		}
		return nil, fmt.Errorf("掃描來源規則失敗: %w", err)
	}

	rule.TemplateName = strings.TrimSpace(templateName.String)
	rule.SkippedAddresses = skipped.String

	if targetDataType.Valid {
		dt := schema.DataType(targetDataType.String)
		rule.TargetDataType = &dt
	}
	if scaleMultiplier.Valid {
		rule.ScaleMultiplier = &scaleMultiplier.Float64
	}
	if scaleOffset.Valid {
		rule.ScaleOffset = &scaleOffset.Float64
	}
	if dataFormat.Valid {
		rule.DataFormat = strings.TrimSpace(dataFormat.String)
	}

	rule.CreatedAt, err = common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析來源規則建立時間失敗: %w", err)
	}
	rule.UpdatedAt, err = common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析來源規則更新時間失敗: %w", err)
	}
	return &rule, nil
}

func scanRuleRows(rows *sql.Rows) (*schema.SourceRule, error) {
	return scanRule(rows)
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
