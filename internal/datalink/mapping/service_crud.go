package mapping

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// CRUD 操作
// =============================================================================

// Create 建立新映射
func (s *Service) Create(ctx context.Context, req CreateMappingRequest) (*schema.Mapping, error) {
	if err := ValidateTransformPipeline(req.TransformPipeline); err != nil {
		return nil, fmt.Errorf("轉換管線無效: %w", err)
	}

	pipelineJSON, err := json.Marshal(req.TransformPipeline)
	if err != nil {
		return nil, fmt.Errorf("序列化轉換管線失敗: %w", err)
	}
	if err := s.validateMappingPreviewGate(ctx, req.TagID, string(pipelineJSON), req.PreviewRawValue); err != nil {
		return nil, fmt.Errorf("映射預覽驗證失敗: %w", err)
	}

	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立映射 ID 失敗: %w", err)
	}

	mapping := &schema.Mapping{
		ID:                id,
		PointID:           req.PointID,
		TagID:             req.TagID,
		TransformPipeline: string(pipelineJSON),
		Enabled:           true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.repo.Create(ctx, mapping); err != nil {
		return nil, fmt.Errorf("建立映射失敗: %w", err)
	}

	return mapping, nil
}

// CreateMappingRequest 建立映射請求
type CreateMappingRequest struct {
	PointID           string                 `json:"point_id"`
	TagID             string                 `json:"tag_id"`
	TransformPipeline []schema.TransformStep `json:"transform_pipeline"`
	PreviewRawValue   interface{}            `json:"preview_raw_value,omitempty"`
}

// Update 更新映射
func (s *Service) Update(ctx context.Context, id string, req UpdateMappingRequest) (*schema.Mapping, error) {
	mapping, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得映射失敗: %w", err)
	}

	nextPipelineJSON := mapping.TransformPipeline
	if req.TransformPipeline != nil {
		if err := ValidateTransformPipeline(req.TransformPipeline); err != nil {
			return nil, fmt.Errorf("轉換管線無效: %w", err)
		}
		pipelineJSON, err := json.Marshal(req.TransformPipeline)
		if err != nil {
			return nil, fmt.Errorf("序列化轉換管線失敗: %w", err)
		}
		nextPipelineJSON = string(pipelineJSON)
		mapping.TransformPipeline = nextPipelineJSON
	}

	nextEnabled := mapping.Enabled
	if req.Enabled != nil {
		nextEnabled = *req.Enabled
	}

	shouldValidate := nextEnabled && (req.TransformPipeline != nil || (req.Enabled != nil && *req.Enabled))
	if shouldValidate {
		if err := s.validateMappingPreviewGate(ctx, mapping.TagID, nextPipelineJSON, req.PreviewRawValue); err != nil {
			return nil, fmt.Errorf("映射預覽驗證失敗: %w", err)
		}
	}
	mapping.Enabled = nextEnabled

	mapping.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, mapping); err != nil {
		return nil, fmt.Errorf("更新映射失敗: %w", err)
	}

	return mapping, nil
}

// UpdateMappingRequest 更新映射請求
type UpdateMappingRequest struct {
	TransformPipeline []schema.TransformStep `json:"transform_pipeline,omitempty"`
	Enabled           *bool                  `json:"enabled,omitempty"`
	PreviewRawValue   interface{}            `json:"preview_raw_value,omitempty"`
}

// Delete 刪除映射
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除映射失敗: %w", err)
	}
	return nil
}

// GetByID 根據 ID 取得映射
func (s *Service) GetByID(ctx context.Context, id string) (*schema.Mapping, error) {
	mapping, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得映射失敗: %w", err)
	}
	return mapping, nil
}

// List 列出映射
func (s *Service) List(ctx context.Context, filter ListFilter) ([]*schema.Mapping, error) {
	mappings, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出映射失敗: %w", err)
	}
	return mappings, nil
}
