package point

import (
	"database/sql"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// scanPoint 從單一行掃描點位
func (r *SQLRepository) scanPoint(row *sql.Row) (*schema.Point, error) {
	var point schema.Point
	var description, function, lastValue, lastError sql.NullString
	var pollingGroupID sql.NullString
	var lastReadAt sql.NullString
	var createdAt, updatedAt string

	err := row.Scan(
		&point.ID,
		&point.DeviceID,
		&point.Name,
		&description,
		&point.Address,
		&function,
		&point.DataType,
		&point.Mode,
		&pollingGroupID,
		&lastReadAt,
		&lastValue,
		&lastError,
		&point.Enabled,
		&createdAt,
		&updatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("點位不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("掃描點位失敗: %w", err)
	}

	if description.Valid {
		point.Description = description.String
	}
	if function.Valid {
		point.Function = function.String
	}
	if pollingGroupID.Valid {
		point.PollingGroupID = &pollingGroupID.String
	}
	if lastReadAt.Valid {
		parsed, err := common.ParseTimeString(lastReadAt.String)
		if err != nil {
			return nil, fmt.Errorf("解析最後讀取時間失敗: %w", err)
		}
		point.LastReadAt = &parsed
	}
	if lastValue.Valid {
		point.LastValue = &lastValue.String
	} else {
		point.LastValue = nil
	}
	if lastError.Valid {
		point.LastError = lastError.String
	} else {
		point.LastError = ""
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析更新時間失敗: %w", err)
	}
	point.CreatedAt = parsedCreatedAt
	point.UpdatedAt = parsedUpdatedAt

	return &point, nil
}

// scanPointFromRows 從多行結果掃描點位
func (r *SQLRepository) scanPointFromRows(rows *sql.Rows) (*schema.Point, error) {
	var point schema.Point
	var description, function, lastValue, lastError sql.NullString
	var pollingGroupID sql.NullString
	var lastReadAt sql.NullString
	var createdAt, updatedAt string

	err := rows.Scan(
		&point.ID,
		&point.DeviceID,
		&point.Name,
		&description,
		&point.Address,
		&function,
		&point.DataType,
		&point.Mode,
		&pollingGroupID,
		&lastReadAt,
		&lastValue,
		&lastError,
		&point.Enabled,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("掃描點位失敗: %w", err)
	}

	if description.Valid {
		point.Description = description.String
	}
	if function.Valid {
		point.Function = function.String
	}
	if pollingGroupID.Valid {
		point.PollingGroupID = &pollingGroupID.String
	}
	if lastReadAt.Valid {
		parsed, err := common.ParseTimeString(lastReadAt.String)
		if err != nil {
			return nil, fmt.Errorf("解析最後讀取時間失敗: %w", err)
		}
		point.LastReadAt = &parsed
	}
	if lastValue.Valid {
		point.LastValue = &lastValue.String
	} else {
		point.LastValue = nil
	}
	if lastError.Valid {
		point.LastError = lastError.String
	} else {
		point.LastError = ""
	}

	parsedCreatedAt, err := common.ParseTimeString(createdAt)
	if err != nil {
		return nil, fmt.Errorf("解析建立時間失敗: %w", err)
	}
	parsedUpdatedAt, err := common.ParseTimeString(updatedAt)
	if err != nil {
		return nil, fmt.Errorf("解析更新時間失敗: %w", err)
	}
	point.CreatedAt = parsedCreatedAt
	point.UpdatedAt = parsedUpdatedAt

	return &point, nil
}
