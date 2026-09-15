package sourcerule

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

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
	var shareStartRegister sql.NullInt64
	var shareStride sql.NullInt64
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
		&rule.ShareEnabled,
		&shareStartRegister,
		&shareStride,
		&createdAt,
		&updatedAt,
		&rule.RevisionID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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
	if shareStartRegister.Valid {
		value := int(shareStartRegister.Int64)
		rule.ShareStartRegister = &value
	}
	if shareStride.Valid {
		value := int(shareStride.Int64)
		rule.ShareStride = &value
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
