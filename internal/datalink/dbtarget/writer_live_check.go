package dbtarget

import (
	"context"
	"errors"
	"fmt"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

// A worker may wake while refresh is still reading. Resolve just the pending
// binding from the repositories so a stale index cannot send it to an old target.
func (w *Writer) mappingIsLive(ctx context.Context, mapping *schema.DatabaseTargetMapping, connector *schema.DatabaseConnector) (bool, error) {
	var current *schema.DatabaseTargetMapping
	err := retryWriterConfigBusy(ctx, func() error {
		var err error
		current, err = w.mappingRepo.GetByID(ctx, mapping.ID)
		return err
	})
	if err != nil {
		return false, fmt.Errorf("確認目前資料庫目標映射失敗: %w", err)
	}
	if current == nil || !current.Enabled {
		return false, nil
	}
	var currentConnector *schema.DatabaseConnector
	err = retryWriterConfigBusy(ctx, func() error {
		var err error
		currentConnector, err = w.connectorRepo.GetByID(ctx, current.ConnectorID)
		return err
	})
	if errors.Is(err, ErrConnectorNotFound) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("確認目前資料庫連接器失敗: %w", err)
	}
	if currentConnector == nil || !currentConnector.Enabled || !sameSnapshotBinding(mapping, connector, current, currentConnector) {
		return false, nil
	}
	if w.tagReader != nil {
		err = retryWriterConfigBusy(ctx, func() error {
			_, err := w.tagReader.GetByID(ctx, current.TagID)
			return err
		})
		if errors.Is(err, tag.ErrTagNotFound) {
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf("確認目前標籤失敗: %w", err)
		}
	}
	return true, nil
}
