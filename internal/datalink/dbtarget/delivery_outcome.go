package dbtarget

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

const (
	deliveryOutcomeSuccess = "success"
	deliveryOutcomeFailed  = "failed"
)

func (s *ConnectorService) recordSchemaEnsureOutcome(
	ctx context.Context,
	connector *schema.DatabaseConnector,
	at time.Time,
	status string,
	errorSummary string,
) error {
	if connector == nil {
		return nil
	}
	outcomeAt := at.UTC()
	connector.LastSchemaEnsureAt = &outcomeAt
	connector.LastSchemaEnsureStatus = strings.TrimSpace(status)
	connector.LastSchemaEnsureError = strings.TrimSpace(errorSummary)
	connector.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, connector); err != nil {
		return fmt.Errorf("更新資料庫 schema ensure outcome 失敗: %w", err)
	}
	return nil
}

func (w *Writer) recordWriteOutcome(
	ctx context.Context,
	connector *schema.DatabaseConnector,
	at time.Time,
	status string,
	errorSummary string,
) error {
	if connector == nil {
		return nil
	}
	outcomeAt := at.UTC()
	connector.LastWriteAt = &outcomeAt
	connector.LastWriteStatus = strings.TrimSpace(status)
	connector.LastWriteError = strings.TrimSpace(errorSummary)
	connector.UpdatedAt = time.Now()
	if err := w.connectorRepo.Update(ctx, connector); err != nil {
		return fmt.Errorf("更新資料庫 write outcome 失敗: %w", err)
	}
	return nil
}

func (w *Writer) recordFlushOutcome(
	ctx context.Context,
	connector *schema.DatabaseConnector,
	at time.Time,
	status string,
	errorSummary string,
) error {
	if connector == nil {
		return nil
	}
	outcomeAt := at.UTC()
	connector.LastFlushAt = &outcomeAt
	connector.LastFlushStatus = strings.TrimSpace(status)
	connector.LastFlushError = strings.TrimSpace(errorSummary)
	connector.UpdatedAt = time.Now()
	if err := w.connectorRepo.Update(ctx, connector); err != nil {
		return fmt.Errorf("更新資料庫 flush outcome 失敗: %w", err)
	}
	return nil
}
