package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// LegacyWriteCoordinator guards legacy target mapping persistence against
// canonical write-group ownership without importing the workspace package.
type LegacyWriteCoordinator interface {
	PreflightLegacyTargetWrite(context.Context, string, string, string) error
	RunLegacyTargetWrite(context.Context, string, *schema.DatabaseTargetMapping, func(context.Context, *sql.Tx) error) error
	CheckLegacyTargetWriteInTx(context.Context, *sql.Tx, string, string, string) error
}

// WithLegacyWriteCoordinator wires the workspace ownership guard into a
// target mapping service while keeping dbtarget independent of workspace.
func WithLegacyWriteCoordinator(service *MappingService, coordinator LegacyWriteCoordinator) *MappingService {
	if service != nil {
		service.legacyWriteCoordinator = coordinator
	}
	return service
}

func (s *MappingService) Create(ctx context.Context, req CreateTargetMappingRequest) (*schema.DatabaseTargetMapping, error) {
	prepared, err := s.PrepareCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	if s.legacyWriteCoordinator != nil {
		err := s.legacyWriteCoordinator.RunLegacyTargetWrite(ctx, prepared.Mapping.ID, prepared.Mapping, func(ctx context.Context, tx *sql.Tx) error {
			return s.PersistInTx(ctx, tx, prepared)
		})
		if err != nil {
			return nil, fmt.Errorf("建立資料庫目標映射失敗: %w", err)
		}
		return prepared.Mapping, nil
	}
	if err := s.repo.Create(ctx, prepared.Mapping); err != nil {
		return nil, fmt.Errorf("建立資料庫目標映射失敗: %w", err)
	}
	return prepared.Mapping, nil
}

func (s *MappingService) Update(ctx context.Context, id string, req UpdateTargetMappingRequest) (*schema.DatabaseTargetMapping, error) {
	prepared, err := s.PrepareUpdate(ctx, id, req)
	if err != nil {
		return nil, err
	}
	if s.legacyWriteCoordinator != nil {
		err := s.legacyWriteCoordinator.RunLegacyTargetWrite(ctx, prepared.Mapping.ID, prepared.Mapping, func(ctx context.Context, tx *sql.Tx) error {
			return s.PersistInTx(ctx, tx, prepared)
		})
		if err != nil {
			return nil, fmt.Errorf("更新資料庫目標映射失敗: %w", err)
		}
		return prepared.Mapping, nil
	}
	if err := s.repo.Update(ctx, prepared.Mapping); err != nil {
		return nil, fmt.Errorf("更新資料庫目標映射失敗: %w", err)
	}
	return prepared.Mapping, nil
}

func (s *MappingService) Delete(ctx context.Context, id string) error {
	if s.legacyWriteCoordinator != nil {
		id = strings.TrimSpace(id)
		if err := s.legacyWriteCoordinator.PreflightLegacyTargetWrite(ctx, id, "", ""); err != nil {
			return err
		}
		mapping, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("取得資料庫目標映射失敗: %w", err)
		}
		err = s.legacyWriteCoordinator.RunLegacyTargetWrite(ctx, id, mapping, func(ctx context.Context, tx *sql.Tx) error {
			binder, ok := s.repo.(mappingTxBinder)
			if !ok || tx == nil {
				return ErrSetupTransactionUnavailable
			}
			if err := binder.WithTx(tx).Delete(ctx, id); err != nil {
				return fmt.Errorf("刪除資料庫目標映射失敗: %w", err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("刪除資料庫目標映射失敗: %w", err)
		}
		return nil
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除資料庫目標映射失敗: %w", err)
	}
	return nil
}
