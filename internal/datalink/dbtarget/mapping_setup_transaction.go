package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

type mappingTxBinder interface {
	WithTx(tx *sql.Tx) TargetMappingRepository
}

// PreparedTargetMapping is a validated target mapping change that is not
// persisted yet; validation and tag reads happen before the transaction.
type PreparedTargetMapping struct {
	Mapping *schema.DatabaseTargetMapping
	create  bool
}

// PrepareCreate validates a new target mapping without writing it.
func (s *MappingService) PrepareCreate(ctx context.Context, req CreateTargetMappingRequest) (*PreparedTargetMapping, error) {
	connector, err := s.connectorRepo.GetByID(ctx, req.ConnectorID)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}

	tagEntity, err := s.tagService.GetByID(ctx, req.TagID)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}
	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立資料庫目標映射 ID 失敗: %w", err)
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	writeMode := req.WriteMode
	if writeMode == "" {
		writeMode = schema.DatabaseWriteModeInsert
	}
	tableSchema := normalizeOptionalString(req.TableSchema)
	if tableSchema == "" {
		tableSchema = defaultSchemaForConnector(connector)
	}
	mapping := &schema.DatabaseTargetMapping{
		ID:                   id,
		TagID:                strings.TrimSpace(req.TagID),
		ConnectorID:          connector.ID,
		TableSchema:          tableSchema,
		TableName:            strings.TrimSpace(req.TableName),
		ColumnName:           strings.TrimSpace(req.ColumnName),
		WriteMode:            writeMode,
		TimestampColumn:      normalizeTimestampColumnForWriteMode(writeMode, req.TimestampColumn),
		GroupKey:             normalizeOptionalPointer(req.GroupKey),
		WriteIntervalSeconds: normalizeOptionalIntPointer(req.WriteIntervalSeconds),
		Enabled:              enabled,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	if err := validateMappingDefinition(ctx, connector, tagEntity, mapping, req.AllowMissingTable); err != nil {
		return nil, err
	}
	return &PreparedTargetMapping{Mapping: mapping, create: true}, nil
}

// PrepareUpdate validates a target mapping update without writing it.
func (s *MappingService) PrepareUpdate(ctx context.Context, id string, req UpdateTargetMappingRequest) (*PreparedTargetMapping, error) {
	mapping, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫目標映射失敗: %w", err)
	}

	connector, err := s.connectorRepo.GetByID(ctx, mapping.ConnectorID)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}

	tagEntity, err := s.tagService.GetByID(ctx, mapping.TagID)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}
	if req.TableSchema != nil {
		tableSchema := normalizeOptionalString(*req.TableSchema)
		if tableSchema == "" {
			tableSchema = defaultSchemaForConnector(connector)
		}
		mapping.TableSchema = tableSchema
	}
	if req.TableName != nil {
		mapping.TableName = strings.TrimSpace(*req.TableName)
	}
	if req.ColumnName != nil {
		mapping.ColumnName = strings.TrimSpace(*req.ColumnName)
	}
	if req.WriteMode != nil {
		mapping.WriteMode = *req.WriteMode
	}
	if mapping.WriteMode != schema.DatabaseWriteModeUpsert {
		mapping.TimestampColumn = nil
	} else if req.TimestampColumn != nil {
		mapping.TimestampColumn = normalizeOptionalPointer(req.TimestampColumn)
	}
	if req.Enabled != nil {
		mapping.Enabled = *req.Enabled
	}
	if req.GroupKey != nil {
		mapping.GroupKey = normalizeOptionalPointer(req.GroupKey)
	}
	if req.WriteIntervalSeconds != nil {
		mapping.WriteIntervalSeconds = normalizeOptionalIntPointer(req.WriteIntervalSeconds)
	}
	mapping.UpdatedAt = time.Now()
	if err := validateMappingDefinition(ctx, connector, tagEntity, mapping, req.AllowMissingTable); err != nil {
		return nil, err
	}
	return &PreparedTargetMapping{Mapping: mapping}, nil
}

// PersistInTx writes a prepared target mapping inside tx.
func (s *MappingService) PersistInTx(ctx context.Context, tx *sql.Tx, prepared *PreparedTargetMapping) error {
	if prepared == nil || prepared.Mapping == nil {
		return validationError("資料庫目標映射變更不可為空")
	}
	binder, ok := s.repo.(mappingTxBinder)
	if !ok || tx == nil {
		return ErrSetupTransactionUnavailable
	}
	repo := binder.WithTx(tx)
	if prepared.create {
		if err := repo.Create(ctx, prepared.Mapping); err != nil {
			return fmt.Errorf("建立資料庫目標映射失敗: %w", err)
		}
		return nil
	}
	if err := repo.Update(ctx, prepared.Mapping); err != nil {
		return fmt.Errorf("更新資料庫目標映射失敗: %w", err)
	}
	return nil
}
