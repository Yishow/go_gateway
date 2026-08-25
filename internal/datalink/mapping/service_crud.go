package mapping

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

func cloneMapping(record *schema.Mapping) *schema.Mapping {
	if record == nil {
		return nil
	}
	copyRecord := *record
	return &copyRecord
}

func validationError(message string) error {
	return fmt.Errorf("%w: %s", ErrValidation, message)
}

// =============================================================================
// CRUD 操作
// =============================================================================

// Create 建立新映射
func (s *Service) Create(ctx context.Context, req CreateMappingRequest) (*schema.Mapping, error) {
	if err := s.validateCardinality(ctx, req.PointID, req.TagID, ""); err != nil {
		return nil, err
	}
	if err := ValidateTransformPipeline(req.TransformPipeline); err != nil {
		return nil, errors.Join(ErrValidation, fmt.Errorf("轉換管線無效: %w", err))
	}

	pipelineJSON, err := json.Marshal(req.TransformPipeline)
	if err != nil {
		return nil, fmt.Errorf("序列化轉換管線失敗: %w", err)
	}
	if err := s.validateMappingPreviewGate(ctx, req.TagID, string(pipelineJSON), req.PreviewRawValue); err != nil {
		return nil, errors.Join(ErrValidation, fmt.Errorf("映射預覽驗證失敗: %w", err))
	}

	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立映射 ID 失敗: %w", err)
	}

	mapping := &schema.Mapping{
		ID:                   id,
		PointID:              req.PointID,
		TagID:                req.TagID,
		TransformPipeline:    string(pipelineJSON),
		Status:               schema.MappingStatusActive,
		RuleCandidateID:      req.RuleCandidateID,
		ProposedSignature:    req.ProposedSignature,
		LastAppliedSignature: req.LastAppliedSignature,
		BlockingReason:       req.BlockingReason,
		Enabled:              true,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	if req.Enabled != nil {
		mapping.Enabled = *req.Enabled
	}
	if !mapping.Enabled {
		mapping.Status = schema.MappingStatusDraft
	}
	if req.Status != nil {
		mapping.Status = *req.Status
	}

	if err := s.repo.Create(ctx, mapping); err != nil {
		return nil, fmt.Errorf("建立映射失敗: %w", err)
	}

	return mapping, nil
}

// CreateMappingRequest 建立映射請求
type CreateMappingRequest struct {
	PointID              string                 `json:"point_id"`
	TagID                string                 `json:"tag_id"`
	Enabled              *bool                  `json:"enabled,omitempty"`
	TransformPipeline    []schema.TransformStep `json:"transform_pipeline"`
	Status               *schema.MappingStatus  `json:"status,omitempty"`
	RuleCandidateID      string                 `json:"rule_candidate_id,omitempty"`
	ProposedSignature    string                 `json:"proposed_signature,omitempty"`
	LastAppliedSignature string                 `json:"last_applied_signature,omitempty"`
	BlockingReason       string                 `json:"blocking_reason,omitempty"`
	PreviewRawValue      interface{}            `json:"preview_raw_value,omitempty"`
}

// Update 更新映射
func (s *Service) Update(ctx context.Context, id string, req UpdateMappingRequest) (*schema.Mapping, error) {
	mapping, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得映射失敗: %w", err)
	}

	nextTagID := mapping.TagID
	if req.TagID != nil {
		if err := s.validateCardinality(ctx, mapping.PointID, *req.TagID, mapping.ID); err != nil {
			return nil, err
		}
		nextTagID = *req.TagID
	}

	nextPipelineJSON := mapping.TransformPipeline
	if req.TransformPipeline != nil {
		if err := ValidateTransformPipeline(req.TransformPipeline); err != nil {
			return nil, errors.Join(ErrValidation, fmt.Errorf("轉換管線無效: %w", err))
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

	shouldValidate := nextEnabled && (req.TransformPipeline != nil || req.TagID != nil || (req.Enabled != nil && *req.Enabled))
	if shouldValidate {
		if err := s.validateMappingPreviewGate(ctx, nextTagID, nextPipelineJSON, req.PreviewRawValue); err != nil {
			return nil, errors.Join(ErrValidation, fmt.Errorf("映射預覽驗證失敗: %w", err))
		}
	}
	mapping.TagID = nextTagID
	mapping.Enabled = nextEnabled
	if req.Status != nil {
		mapping.Status = *req.Status
	}
	if req.RuleCandidateID != nil {
		mapping.RuleCandidateID = *req.RuleCandidateID
	}
	if req.ProposedSignature != nil {
		mapping.ProposedSignature = *req.ProposedSignature
	}
	if req.LastAppliedSignature != nil {
		mapping.LastAppliedSignature = *req.LastAppliedSignature
	}
	if req.BlockingReason != nil {
		mapping.BlockingReason = *req.BlockingReason
	}

	mapping.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, mapping); err != nil {
		return nil, fmt.Errorf("更新映射失敗: %w", err)
	}

	return mapping, nil
}

// UpdateMappingRequest 更新映射請求
type UpdateMappingRequest struct {
	TagID                *string                `json:"tag_id,omitempty"`
	TransformPipeline    []schema.TransformStep `json:"transform_pipeline,omitempty"`
	Enabled              *bool                  `json:"enabled,omitempty"`
	Status               *schema.MappingStatus  `json:"status,omitempty"`
	RuleCandidateID      *string                `json:"rule_candidate_id,omitempty"`
	ProposedSignature    *string                `json:"proposed_signature,omitempty"`
	LastAppliedSignature *string                `json:"last_applied_signature,omitempty"`
	BlockingReason       *string                `json:"blocking_reason,omitempty"`
	PreviewRawValue      interface{}            `json:"preview_raw_value,omitempty"`
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

// Restore persists an exact mapping snapshot during a higher-level mutation compensation.
func (s *Service) Restore(ctx context.Context, record *schema.Mapping) error {
	if record == nil {
		return fmt.Errorf("mapping snapshot is empty")
	}
	if _, err := s.GetByID(ctx, record.ID); errors.Is(err, ErrMappingNotFound) {
		if createErr := s.repo.Create(ctx, cloneMapping(record)); createErr != nil {
			return fmt.Errorf("restore mapping: %w", createErr)
		}
		return nil
	} else if err != nil {
		return err
	}
	if err := s.repo.Update(ctx, cloneMapping(record)); err != nil {
		return fmt.Errorf("restore mapping: %w", err)
	}
	return nil
}

// List 列出映射
func (s *Service) List(ctx context.Context, filter ListFilter) ([]*schema.Mapping, error) {
	mappings, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出映射失敗: %w", err)
	}
	return mappings, nil
}

func (s *Service) validateCardinality(ctx context.Context, pointID, tagID, ignoreMappingID string) error {
	pointMappings, err := s.repo.GetByPointID(ctx, pointID)
	if err != nil {
		return fmt.Errorf("檢查 point 既有映射失敗: %w", err)
	}
	for _, existing := range pointMappings {
		if existing.ID == ignoreMappingID {
			continue
		}
		return validationError(fmt.Sprintf("point %s 已綁定其他 tag", pointID))
	}

	tagMappings, err := s.repo.GetByTagID(ctx, tagID)
	if err != nil {
		return fmt.Errorf("檢查 tag 既有映射失敗: %w", err)
	}
	for _, existing := range tagMappings {
		if existing.ID == ignoreMappingID {
			continue
		}
		return validationError(fmt.Sprintf("tag %s 已綁定其他 point", tagID))
	}

	return nil
}
