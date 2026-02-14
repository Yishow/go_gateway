package tag

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

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

// Update 更新標籤
func (s *Service) Update(ctx context.Context, id string, req UpdateTagRequest) (*schema.Tag, error) {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}

	// 標籤鍵不可修改 (若提供則驗證一致性)
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
