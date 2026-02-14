package tag

import (
	"context"
	"fmt"

	"go-gateway/internal/datalink/schema"
)

// Activate 啟用標籤
func (s *Service) Activate(ctx context.Context, id string) error {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得標籤失敗: %w", err)
	}

	if tag.Status == schema.TagStatusActive {
		return nil
	}

	if tag.Status == schema.TagStatusRetired {
		return fmt.Errorf("已退役的標籤不能重新啟用")
	}

	if err := s.repo.UpdateStatus(ctx, id, schema.TagStatusActive); err != nil {
		return fmt.Errorf("更新標籤狀態失敗: %w", err)
	}

	return nil
}

// Retire 退役標籤
func (s *Service) Retire(ctx context.Context, id string) error {
	tag, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("取得標籤失敗: %w", err)
	}

	if tag.Status == schema.TagStatusRetired {
		return nil
	}

	if err := s.repo.UpdateStatus(ctx, id, schema.TagStatusRetired); err != nil {
		return fmt.Errorf("更新標籤狀態失敗: %w", err)
	}

	return nil
}
