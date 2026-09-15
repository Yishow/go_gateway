package measurement

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
)

// Service 提供量測定義的業務邏輯、版本修訂與範本套用。
type Service struct {
	repo Repository
}

// NewService 建立新的 Service。
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateMeasurement 建立量測定義。
func (s *Service) CreateMeasurement(ctx context.Context, def *MeasurementDefinition) error {
	if strings.TrimSpace(def.ID) == "" {
		uuid, err := common.NewUUID()
		if err != nil {
			return fmt.Errorf("create measurement id: %w", err)
		}
		def.ID = "meas-" + uuid
	}
	if strings.TrimSpace(def.DefinitionRevision) == "" {
		def.DefinitionRevision = initialDefinitionRevision
	}
	if strings.TrimSpace(def.SourceBindingRevision) == "" {
		def.SourceBindingRevision = initialDefinitionRevision
	}
	if strings.TrimSpace(def.SeriesEpoch) == "" {
		def.SeriesEpoch = initialSeriesEpoch
	}
	if err := def.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, def)
}

// UpdateMeasurement 更新量測定義，自動檢查可比性並維護 series_epoch 與 revision。
func (s *Service) UpdateMeasurement(ctx context.Context, def *MeasurementDefinition) (bool, error) {
	if err := def.Validate(); err != nil {
		return false, err
	}

	existing, err := s.repo.GetByID(ctx, def.ID)
	if err != nil {
		return false, err
	}

	needsNewEpoch, nextEpoch, _ := CheckEpochTransition(*existing, *def)
	if needsNewEpoch {
		def.SeriesEpoch = nextEpoch
	} else {
		def.SeriesEpoch = existing.SeriesEpoch
	}

	def.DefinitionRevision = incrementRevision(existing.DefinitionRevision)
	def.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, def); err != nil {
		return false, err
	}
	return needsNewEpoch, nil
}

// GetMeasurement 取得單一量測定義。
func (s *Service) GetMeasurement(ctx context.Context, id string) (*MeasurementDefinition, error) {
	return s.repo.GetByID(ctx, id)
}

// ListByWorkspace 取得 Workspace 底下所有量測定義。
func (s *Service) ListByWorkspace(ctx context.Context, workspaceID string) ([]MeasurementDefinition, error) {
	return s.repo.ListByWorkspace(ctx, workspaceID)
}

// ListByDevice 取得指定設備的所有量測定義。
func (s *Service) ListByDevice(ctx context.Context, workspaceID, deviceID string) ([]MeasurementDefinition, error) {
	return s.repo.ListByDevice(ctx, workspaceID, deviceID)
}

// DeleteMeasurement 刪除指定量測定義。
func (s *Service) DeleteMeasurement(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// ApplyTemplate 接受經使用者確認的範本預覽並建立量測定義。
func (s *Service) ApplyTemplate(ctx context.Context, workspaceID string, preview *TemplateApplyPreview) error {
	if !preview.Confirmed {
		return fmt.Errorf("template preview must be confirmed before applying")
	}

	for _, def := range preview.Definitions {
		defCopy := def
		defCopy.WorkspaceID = workspaceID
		if err := s.CreateMeasurement(ctx, &defCopy); err != nil {
			return fmt.Errorf("failed to create measurement for item %s: %w", def.Name, err)
		}
	}
	return nil
}

func incrementRevision(rev string) string {
	rev = strings.TrimSpace(rev)
	if strings.HasPrefix(rev, "rev-") {
		var n int
		if _, err := fmt.Sscanf(rev, "rev-%d", &n); err == nil {
			return fmt.Sprintf("rev-%d", n+1)
		}
	}
	return "rev-2"
}
