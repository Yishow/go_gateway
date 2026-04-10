package dbtarget

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"go-gateway/internal/datalink/schema"
)

type Writer struct {
	connectorRepo ConnectorRepository
	mappingRepo   TargetMappingRepository
	now           func() time.Time
	flushTick     <-chan time.Time
	stopTick      func()
	stopCh        chan struct{}
	closeOnce     sync.Once
	wg            sync.WaitGroup
	mu            sync.Mutex
	grouped       map[groupedWriteKey]*groupedWriteBucket
}

func NewWriter(connectorRepo ConnectorRepository, mappingRepo TargetMappingRepository) *Writer {
	return NewWriterWithConfig(connectorRepo, mappingRepo, WriterConfig{})
}

type WriterConfig struct {
	FlushTick     <-chan time.Time
	FlushInterval time.Duration
	Now           func() time.Time
}

type groupedWriteKey struct {
	ConnectorID              string
	SchemaName               string
	TableName                string
	GroupKey                 string
	WriteMode                schema.DatabaseWriteMode
	TimestampColumn          string
	EffectiveIntervalSeconds int
	BucketStart              time.Time
}

type groupedWriteBucket struct {
	Key       groupedWriteKey
	Connector *schema.DatabaseConnector
	Values    map[string]any
}

func NewWriterWithConfig(
	connectorRepo ConnectorRepository,
	mappingRepo TargetMappingRepository,
	config WriterConfig,
) *Writer {
	writer := &Writer{
		connectorRepo: connectorRepo,
		mappingRepo:   mappingRepo,
		now:           config.Now,
		stopCh:        make(chan struct{}),
		grouped:       make(map[groupedWriteKey]*groupedWriteBucket),
	}
	if writer.now == nil {
		writer.now = time.Now
	}
	if config.FlushTick != nil {
		writer.flushTick = config.FlushTick
		writer.stopTick = func() {}
	} else {
		interval := config.FlushInterval
		if interval <= 0 {
			interval = time.Second
		}
		ticker := time.NewTicker(interval)
		writer.flushTick = ticker.C
		writer.stopTick = ticker.Stop
	}
	writer.wg.Add(1)
	go writer.flushLoop()
	return writer
}

func (w *Writer) WriteTagValue(ctx context.Context, tagID string, value any, observedAt time.Time) error {
	if strings.TrimSpace(tagID) == "" {
		return fmt.Errorf("tagID 不可為空")
	}

	enabled := true
	mappings, err := w.mappingRepo.List(ctx, TargetMappingListFilter{
		TagID:   &tagID,
		Enabled: &enabled,
	})
	if err != nil {
		return fmt.Errorf("查詢資料庫目標映射失敗: %w", err)
	}
	if len(mappings) == 0 {
		return nil
	}

	var failures []string
	for _, mapping := range mappings {
		connector, err := w.connectorRepo.GetByID(ctx, mapping.ConnectorID)
		if err != nil {
			recordWriteHistory(mapping.ConnectorID, WriteHistoryRecord{
				ObservedAt:   observedAt,
				Status:       "failed",
				RowCount:     0,
				TableName:    mapping.TableName,
				GroupKey:     cloneOptionalString(mapping.GroupKey),
				ErrorSummary: err.Error(),
			})
			failures = append(failures, fmt.Sprintf("mapping %s 取得連接器失敗: %v", mapping.ID, err))
			continue
		}
		if !connector.Enabled {
			continue
		}

		if hasGroupedWriteKey(mapping) {
			if err := w.bufferGroupedWrite(connector, mapping, value, observedAt); err != nil {
				failures = append(failures, fmt.Sprintf("mapping %s 分組緩衝失敗: %v", mapping.ID, err))
			}
			continue
		}

		if err := w.writeMapping(ctx, connector, mapping, value, observedAt); err != nil {
			recordWriteHistory(connector.ID, WriteHistoryRecord{
				ObservedAt:               observedAt,
				Status:                   "failed",
				RowCount:                 0,
				TableName:                mapping.TableName,
				GroupKey:                 cloneOptionalString(mapping.GroupKey),
				EffectiveIntervalSeconds: effectiveWriteIntervalSeconds(mapping, connector),
				ErrorSummary:             err.Error(),
			})
			failures = append(failures, fmt.Sprintf("mapping %s 寫入失敗: %v", mapping.ID, err))
			continue
		}
		recordWriteHistory(connector.ID, WriteHistoryRecord{
			ObservedAt:               observedAt.UTC(),
			Status:                   "success",
			RowCount:                 1,
			TableName:                mapping.TableName,
			GroupKey:                 nil,
			EffectiveIntervalSeconds: effectiveWriteIntervalSeconds(mapping, connector),
		})
	}

	if len(failures) > 0 {
		return fmt.Errorf("%s", strings.Join(failures, "; "))
	}

	return nil
}

func (w *Writer) writeMapping(
	ctx context.Context,
	connector *schema.DatabaseConnector,
	mapping *schema.DatabaseTargetMapping,
	value any,
	observedAt time.Time,
) error {
	connectionConfig, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return err
	}

	manager, err := openExternalDBManager(connector.Kind, connectionConfig)
	if err != nil {
		return err
	}
	defer manager.Close()

	query, args, err := buildWriteStatement(connector.Kind, mapping, value, observedAt)
	if err != nil {
		return err
	}

	if _, err := manager.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("執行資料庫寫入失敗: %w", err)
	}

	return nil
}

func (w *Writer) Close(ctx context.Context) error {
	var closeErr error
	w.closeOnce.Do(func() {
		close(w.stopCh)
		if w.stopTick != nil {
			w.stopTick()
		}
		w.wg.Wait()
		closeErr = w.flushBuckets(ctx, w.now(), true)
	})
	return closeErr
}

func (w *Writer) flushLoop() {
	defer w.wg.Done()
	for {
		select {
		case tick, ok := <-w.flushTick:
			if !ok {
				return
			}
			_ = w.flushBuckets(context.Background(), tick, false)
		case <-w.stopCh:
			return
		}
	}
}

func (w *Writer) bufferGroupedWrite(
	connector *schema.DatabaseConnector,
	mapping *schema.DatabaseTargetMapping,
	value any,
	observedAt time.Time,
) error {
	key, err := buildGroupedWriteKey(connector, mapping, observedAt)
	if err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	bucket, ok := w.grouped[key]
	if !ok {
		bucket = &groupedWriteBucket{
			Key:       key,
			Connector: connector,
			Values:    make(map[string]any),
		}
		w.grouped[key] = bucket
	}
	bucket.Values[mapping.ColumnName] = normalizeDBValue(value)
	return nil
}

func buildGroupedWriteKey(
	connector *schema.DatabaseConnector,
	mapping *schema.DatabaseTargetMapping,
	observedAt time.Time,
) (groupedWriteKey, error) {
	groupKey := strings.TrimSpace(derefOptionalString(mapping.GroupKey))
	if groupKey == "" {
		return groupedWriteKey{}, fmt.Errorf("grouped write requires non-empty group_key")
	}
	intervalSeconds := effectiveWriteIntervalSeconds(mapping, connector)
	if intervalSeconds <= 0 {
		return groupedWriteKey{}, fmt.Errorf("grouped write requires positive interval")
	}
	if mapping.TableName == "" || mapping.ColumnName == "" {
		return groupedWriteKey{}, fmt.Errorf("資料庫目標映射缺少資料表或欄位資訊")
	}
	if mapping.WriteMode == schema.DatabaseWriteModeUpsert && mapping.TimestampColumn == nil {
		return groupedWriteKey{}, fmt.Errorf("upsert 模式缺少 timestamp_column")
	}

	bucketStart := observedAt.UTC().Truncate(time.Duration(intervalSeconds) * time.Second)
	return groupedWriteKey{
		ConnectorID:              mapping.ConnectorID,
		SchemaName:               mapping.TableSchema,
		TableName:                mapping.TableName,
		GroupKey:                 groupKey,
		WriteMode:                mapping.WriteMode,
		TimestampColumn:          derefOptionalString(mapping.TimestampColumn),
		EffectiveIntervalSeconds: intervalSeconds,
		BucketStart:              bucketStart,
	}, nil
}

func (w *Writer) flushBuckets(ctx context.Context, now time.Time, flushAll bool) error {
	buckets := w.takeFlushableBuckets(now.UTC(), flushAll)
	if len(buckets) == 0 {
		return nil
	}

	var failures []string
	for _, bucket := range buckets {
		if err := w.flushGroupedBucket(ctx, bucket); err != nil {
			recordWriteHistory(bucket.Key.ConnectorID, WriteHistoryRecord{
				ObservedAt:               bucket.Key.BucketStart,
				Status:                   "failed",
				RowCount:                 0,
				GroupKey:                 copyStringPointer(bucket.Key.GroupKey),
				TableName:                bucket.Key.TableName,
				EffectiveIntervalSeconds: bucket.Key.EffectiveIntervalSeconds,
				ErrorSummary:             err.Error(),
			})
			failures = append(failures, err.Error())
			continue
		}
		recordWriteHistory(bucket.Key.ConnectorID, WriteHistoryRecord{
			ObservedAt:               bucket.Key.BucketStart,
			Status:                   "success",
			RowCount:                 1,
			GroupKey:                 copyStringPointer(bucket.Key.GroupKey),
			TableName:                bucket.Key.TableName,
			EffectiveIntervalSeconds: bucket.Key.EffectiveIntervalSeconds,
		})
	}
	if len(failures) > 0 {
		return fmt.Errorf("%s", strings.Join(failures, "; "))
	}
	return nil
}

func (w *Writer) takeFlushableBuckets(now time.Time, flushAll bool) []*groupedWriteBucket {
	w.mu.Lock()
	defer w.mu.Unlock()

	buckets := make([]*groupedWriteBucket, 0, len(w.grouped))
	for key, bucket := range w.grouped {
		dueAt := key.BucketStart.Add(time.Duration(key.EffectiveIntervalSeconds) * time.Second)
		if !flushAll && now.Before(dueAt) {
			continue
		}
		buckets = append(buckets, bucket)
		delete(w.grouped, key)
	}
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].Key.BucketStart.Before(buckets[j].Key.BucketStart)
	})
	return buckets
}

func (w *Writer) flushGroupedBucket(ctx context.Context, bucket *groupedWriteBucket) error {
	connectionConfig, err := parseConnectionConfig(bucket.Connector.ConnectionConfig)
	if err != nil {
		return err
	}

	manager, err := openExternalDBManager(bucket.Connector.Kind, connectionConfig)
	if err != nil {
		return err
	}
	defer manager.Close()

	query, args, err := buildGroupedWriteStatement(bucket.Connector.Kind, bucket.Key, bucket.Values)
	if err != nil {
		return err
	}
	if _, err := manager.DB().ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("執行資料庫分組寫入失敗: %w", err)
	}
	return nil
}

func buildGroupedWriteStatement(
	kind schema.DatabaseConnectorKind,
	key groupedWriteKey,
	values map[string]any,
) (string, []any, error) {
	if len(values) == 0 {
		return "", nil, fmt.Errorf("grouped write has no buffered values")
	}

	columnNames := make([]string, 0, len(values))
	for columnName := range values {
		columnNames = append(columnNames, columnName)
	}
	sort.Strings(columnNames)

	columns := make([]string, 0, len(columnNames)+1)
	args := make([]any, 0, len(columnNames)+1)
	if key.WriteMode == schema.DatabaseWriteModeUpsert && strings.TrimSpace(key.TimestampColumn) != "" {
		columns = append(columns, quoteIdentifier(kind, key.TimestampColumn))
		args = append(args, key.BucketStart)
	}
	for _, columnName := range columnNames {
		columns = append(columns, quoteIdentifier(kind, columnName))
		args = append(args, values[columnName])
	}

	placeholders := buildPlaceholders(kind, len(args))
	tableRef := qualifiedTableName(kind, key.SchemaName, key.TableName)
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableRef,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	if key.WriteMode == schema.DatabaseWriteModeUpsert {
		if strings.TrimSpace(key.TimestampColumn) == "" {
			return "", nil, fmt.Errorf("upsert 模式缺少 timestamp_column")
		}
		assignments := make([]string, 0, len(columnNames))
		for _, columnName := range columnNames {
			quoted := quoteIdentifier(kind, columnName)
			assignments = append(assignments, fmt.Sprintf("%s = EXCLUDED.%s", quoted, quoted))
		}
		query += fmt.Sprintf(
			" ON CONFLICT (%s) DO UPDATE SET %s",
			quoteIdentifier(kind, key.TimestampColumn),
			strings.Join(assignments, ", "),
		)
	}
	return query, args, nil
}

func hasGroupedWriteKey(mapping *schema.DatabaseTargetMapping) bool {
	return strings.TrimSpace(derefOptionalString(mapping.GroupKey)) != ""
}

func effectiveWriteIntervalSeconds(mapping *schema.DatabaseTargetMapping, connector *schema.DatabaseConnector) int {
	if mapping != nil && mapping.WriteIntervalSeconds != nil {
		return *mapping.WriteIntervalSeconds
	}
	if connector != nil && connector.DefaultWriteIntervalSeconds > 0 {
		return connector.DefaultWriteIntervalSeconds
	}
	return 15
}

func cloneOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := strings.TrimSpace(*value)
	return &cloned
}

func derefOptionalString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func copyStringPointer(value string) *string {
	trimmed := strings.TrimSpace(value)
	return &trimmed
}

func buildWriteStatement(
	kind schema.DatabaseConnectorKind,
	mapping *schema.DatabaseTargetMapping,
	value any,
	observedAt time.Time,
) (string, []any, error) {
	if mapping.TableName == "" || mapping.ColumnName == "" {
		return "", nil, fmt.Errorf("資料庫目標映射缺少資料表或欄位資訊")
	}

	columns := []string{quoteIdentifier(kind, mapping.ColumnName)}
	args := []any{normalizeDBValue(value)}
	if mapping.WriteMode == schema.DatabaseWriteModeUpsert && mapping.TimestampColumn != nil {
		columns = append(columns, quoteIdentifier(kind, *mapping.TimestampColumn))
		args = append(args, observedAt.UTC())
	}

	placeholders := buildPlaceholders(kind, len(args))
	tableRef := qualifiedTableName(kind, mapping.TableSchema, mapping.TableName)
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableRef,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	if mapping.WriteMode == schema.DatabaseWriteModeUpsert {
		if mapping.TimestampColumn == nil {
			return "", nil, fmt.Errorf("upsert 模式缺少 timestamp_column")
		}

		valueColumn := quoteIdentifier(kind, mapping.ColumnName)
		timestampColumn := quoteIdentifier(kind, *mapping.TimestampColumn)
		query += fmt.Sprintf(
			" ON CONFLICT (%s) DO UPDATE SET %s = EXCLUDED.%s",
			timestampColumn,
			valueColumn,
			valueColumn,
		)
	}

	return query, args, nil
}

func buildPlaceholders(kind schema.DatabaseConnectorKind, count int) []string {
	placeholders := make([]string, 0, count)
	for i := 0; i < count; i++ {
		switch kind {
		case schema.DatabaseConnectorKindPostgres:
			placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		default:
			placeholders = append(placeholders, "?")
		}
	}
	return placeholders
}

func quoteIdentifier(kind schema.DatabaseConnectorKind, name string) string {
	escaped := strings.ReplaceAll(strings.TrimSpace(name), `"`, `""`)
	return `"` + escaped + `"`
}

func qualifiedTableName(kind schema.DatabaseConnectorKind, schemaName string, tableName string) string {
	quotedTable := quoteIdentifier(kind, tableName)
	if strings.TrimSpace(schemaName) == "" || kind == schema.DatabaseConnectorKindSQLite {
		return quotedTable
	}
	return quoteIdentifier(kind, schemaName) + "." + quotedTable
}

func normalizeDBValue(value any) any {
	switch typed := value.(type) {
	case nil:
		return nil
	case map[string]any, []any:
		data, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprintf("%v", typed)
		}
		return string(data)
	case time.Time:
		return typed.UTC()
	default:
		return value
	}
}

var _ interface {
	WriteTagValue(context.Context, string, any, time.Time) error
} = (*Writer)(nil)
