package tag

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

func cloneTag(record *schema.Tag) *schema.Tag {
	if record == nil {
		return nil
	}
	copyRecord := *record
	return &copyRecord
}

// Create 建立新標籤
func (s *Service) Create(ctx context.Context, req CreateTagRequest) (*schema.Tag, error) {
	// 驗證標籤鍵
	if err := ValidateTagKey(req.Key); err != nil {
		return nil, err
	}

	// 檢查鍵是否已存在
	exists, err := s.repo.ExistsByKey(ctx, req.Key)
	if err != nil {
		return nil, fmt.Errorf("檢查標籤鍵失敗: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("標籤鍵已存在: %s", req.Key)
	}

	// 驗證資料型別
	if !isValidDataType(req.DataType) {
		return nil, fmt.Errorf("不支援的資料型別: %s", req.DataType)
	}

	// 序列化 Labels
	labelsJSON := "{}"
	if req.Labels != nil {
		if data, err := json.Marshal(req.Labels); err == nil {
			labelsJSON = string(data)
		}
	}

	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立標籤 ID 失敗: %w", err)
	}

	tag := &schema.Tag{
		ID:          id,
		Key:         req.Key,
		KeyLower:    NormalizeTagKey(req.Key),
		DisplayName: req.DisplayName,
		Description: req.Description,
		Unit:        req.Unit,
		DataType:    req.DataType,
		Status:      schema.TagStatusDraft,
		Labels:      labelsJSON,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 設定預設顯示名稱
	if tag.DisplayName == "" {
		tag.DisplayName = req.Key
	}

	if err := s.repo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("建立標籤失敗: %w", err)
	}

	return tag, nil
}

// CreateTagRequest 建立標籤請求
type CreateTagRequest struct {
	Key         string            `json:"key"`
	DisplayName string            `json:"display_name,omitempty"`
	Description string            `json:"description,omitempty"`
	Unit        string            `json:"unit,omitempty"`
	DataType    schema.DataType   `json:"data_type"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// BatchCreateError 批量建立時的單筆錯誤
type BatchCreateError struct {
	Key   string
	Error string
}

// BatchCreate 批量建立標籤，使用交易確保原子性
func (s *Service) BatchCreate(ctx context.Context, reqs []CreateTagRequest) ([]string, []BatchCreateError) {
	s.mu.Lock()
	defer s.mu.Unlock()

	created := make([]string, 0, len(reqs))
	var errs []BatchCreateError

	// 預先驗證所有請求
	tags := make([]*schema.Tag, 0, len(reqs))
	for _, req := range reqs {
		if err := ValidateTagKey(req.Key); err != nil {
			errs = append(errs, BatchCreateError{Key: req.Key, Error: err.Error()})
			continue
		}

		exists, err := s.repo.ExistsByKey(ctx, req.Key)
		if err != nil {
			errs = append(errs, BatchCreateError{Key: req.Key, Error: fmt.Sprintf("檢查標籤鍵失敗: %v", err)})
			continue
		}
		if exists {
			errs = append(errs, BatchCreateError{Key: req.Key, Error: fmt.Sprintf("標籤鍵已存在: %s", req.Key)})
			continue
		}

		if !isValidDataType(req.DataType) {
			errs = append(errs, BatchCreateError{Key: req.Key, Error: fmt.Sprintf("不支援的資料型別: %s", req.DataType)})
			continue
		}

		labelsJSON := "{}"
		if req.Labels != nil {
			if data, jsonErr := json.Marshal(req.Labels); jsonErr == nil {
				labelsJSON = string(data)
			}
		}

		id, err := common.NewUUID()
		if err != nil {
			errs = append(errs, BatchCreateError{Key: req.Key, Error: fmt.Sprintf("建立 ID 失敗: %v", err)})
			continue
		}

		tag := &schema.Tag{
			ID:          id,
			Key:         req.Key,
			KeyLower:    NormalizeTagKey(req.Key),
			DisplayName: req.DisplayName,
			Description: req.Description,
			Unit:        req.Unit,
			DataType:    req.DataType,
			Status:      schema.TagStatusDraft,
			Labels:      labelsJSON,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if tag.DisplayName == "" {
			tag.DisplayName = req.Key
		}

		tags = append(tags, tag)
	}

	// 使用交易批量寫入
	if len(tags) > 0 {
		if err := s.repo.BatchCreate(ctx, tags); err != nil {
			for _, tag := range tags {
				errs = append(errs, BatchCreateError{Key: tag.Key, Error: err.Error()})
			}
		} else {
			for _, tag := range tags {
				created = append(created, tag.ID)
			}
		}
	}

	return created, errs
}

// Update 更新標籤
func (s *Service) Update(ctx context.Context, id string, req UpdateTagRequest) (*schema.Tag, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}

	// 標籤鍵不可修改；請求提供時須驗證一致性。
	if req.Key != nil && NormalizeTagKey(*req.Key) != tag.KeyLower {
		return nil, fmt.Errorf("標籤鍵不可修改")
	}

	if req.DisplayName != nil {
		tag.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		tag.Description = *req.Description
	}
	if req.Unit != nil {
		tag.Unit = *req.Unit
	}
	if req.DataType != nil {
		if !isValidDataType(*req.DataType) {
			return nil, fmt.Errorf("不支援的資料型別: %s", *req.DataType)
		}
		tag.DataType = *req.DataType
	}
	if req.Labels != nil {
		if data, err := json.Marshal(req.Labels); err == nil {
			tag.Labels = string(data)
		}
	}

	tag.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, tag); err != nil {
		return nil, fmt.Errorf("更新標籤失敗: %w", err)
	}

	return tag, nil
}

// UpdateTagRequest 更新標籤請求
type UpdateTagRequest struct {
	Key         *string           `json:"key,omitempty"`
	DisplayName *string           `json:"display_name,omitempty"`
	Description *string           `json:"description,omitempty"`
	Unit        *string           `json:"unit,omitempty"`
	DataType    *schema.DataType  `json:"data_type,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// Delete 刪除標籤
func (s *Service) Delete(ctx context.Context, id string) error {
	// TODO: 檢查是否有映射使用此標籤
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除標籤失敗: %w", err)
	}
	return nil
}

// GetByID 根據 ID 取得標籤
func (s *Service) GetByID(ctx context.Context, id string) (*schema.Tag, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}
	return tag, nil
}

// Restore persists an exact tag snapshot during a higher-level mutation compensation.
func (s *Service) Restore(ctx context.Context, record *schema.Tag) error {
	if record == nil {
		return fmt.Errorf("tag snapshot is empty")
	}
	if _, err := s.GetByID(ctx, record.ID); errors.Is(err, ErrTagNotFound) {
		if createErr := s.repo.Create(ctx, cloneTag(record)); createErr != nil {
			return fmt.Errorf("restore tag: %w", createErr)
		}
		return nil
	} else if err != nil {
		return err
	}
	if err := s.repo.Update(ctx, cloneTag(record)); err != nil {
		return fmt.Errorf("restore tag: %w", err)
	}
	return nil
}

// GetByKey 根據鍵取得標籤
func (s *Service) GetByKey(ctx context.Context, key string) (*schema.Tag, error) {
	tag, err := s.repo.GetByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}
	return tag, nil
}

// List 列出標籤
func (s *Service) List(ctx context.Context, filter ListFilter) ([]*schema.Tag, error) {
	tags, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出標籤失敗: %w", err)
	}
	return tags, nil
}
