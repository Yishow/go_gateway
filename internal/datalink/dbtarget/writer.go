package dbtarget

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
)

// Writer delivers legacy database target mappings. Ingest (WriteTagValue) only
// updates in-memory state: a latest-value table keyed by mapping for ungrouped
// mappings and the grouped buckets for grouped mappings. One background worker
// per enabled connector writes the due snapshots, so a slow or unreachable
// connector never delays acquisition or any other connector.
type Writer struct {
	connectorRepo ConnectorRepository
	mappingRepo   TargetMappingRepository
	tagReader     ConnectorTagReader
	suppress      func(connectorID, tagID string) bool
	now           func() time.Time
	open          func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error)
	flushTick     <-chan time.Time
	stopTick      func()
	stopCh        chan struct{}
	closeOnce     sync.Once
	closeDone     chan struct{}
	closeErr      error
	wg            sync.WaitGroup
	baseCtx       context.Context
	cancelBase    context.CancelFunc
	workerCtx     context.Context
	cancelWorkers context.CancelFunc
	recoveryLocks sync.Map

	// mu guards the in-memory delivery state filled by ingest.
	mu               sync.Mutex
	closed           bool
	grouped          map[groupedWriteKey]*groupedWriteBucket
	groupedLocations map[string]groupedLocation
	pending          map[string]map[string]*pendingWrite // connector ID -> mapping ID
	lastFlush        map[string]time.Time                // mapping ID -> last attempted interval boundary
	observed         map[string]*observedIntervalSet     // mapping ID -> occupied intervals

	indexMu        sync.RWMutex
	index          map[string][]indexedMapping // tag ID -> live mappings
	connectorIDs   map[string]struct{}
	indexErr       error
	refreshMu      sync.Mutex
	refreshing     atomic.Bool
	refreshActive  atomic.Int64
	refreshVersion atomic.Uint64
	refreshReq     chan struct{}
	refreshDone    chan struct{}

	workersMu sync.Mutex
	workers   map[string]*connectorWorker
	openMu    sync.Mutex
	opening   map[string]struct{}
}

func NewWriter(connectorRepo ConnectorRepository, mappingRepo TargetMappingRepository) *Writer {
	return NewWriterWithConfig(connectorRepo, mappingRepo, WriterConfig{})
}

type WriterConfig struct {
	// FlushTick drives mapping refresh and connector workers; nil uses a
	// ticker with FlushInterval (default one second).
	FlushTick     <-chan time.Time
	FlushInterval time.Duration
	Now           func() time.Time
	TagReader     ConnectorTagReader
	// SuppressMapping, when it returns true for a connector and tag, makes the
	// legacy writer skip that mapping because a canonical write group owns the
	// output. It keeps exactly one writer per output.
	SuppressMapping func(connectorID, tagID string) bool
	// OpenDestination opens a connector's target database. Nil uses the
	// production opener; tests inject unreachable connectors through it.
	OpenDestination func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error)
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
	Mappings  map[string]*schema.DatabaseTargetMapping
	DueAt     time.Time
	Retry     bool
}

type groupedLocation struct {
	key    groupedWriteKey
	column string
}

func NewWriterWithConfig(
	connectorRepo ConnectorRepository,
	mappingRepo TargetMappingRepository,
	config WriterConfig,
) *Writer {
	baseCtx, cancelBase := context.WithCancel(context.Background())
	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	writer := &Writer{
		connectorRepo:    connectorRepo,
		mappingRepo:      mappingRepo,
		tagReader:        config.TagReader,
		suppress:         config.SuppressMapping,
		now:              config.Now,
		open:             config.OpenDestination,
		stopCh:           make(chan struct{}),
		closeDone:        make(chan struct{}),
		baseCtx:          baseCtx,
		cancelBase:       cancelBase,
		workerCtx:        workerCtx,
		cancelWorkers:    cancelWorkers,
		grouped:          make(map[groupedWriteKey]*groupedWriteBucket),
		groupedLocations: make(map[string]groupedLocation),
		pending:          make(map[string]map[string]*pendingWrite),
		lastFlush:        make(map[string]time.Time),
		observed:         make(map[string]*observedIntervalSet),
		index:            make(map[string][]indexedMapping),
		connectorIDs:     make(map[string]struct{}),
		refreshReq:       make(chan struct{}, 1),
		refreshDone:      make(chan struct{}, 1),
		workers:          make(map[string]*connectorWorker),
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
	if connectorRepo != nil && mappingRepo != nil {
		if err := writer.RefreshMappings(baseCtx); err != nil {
			slog.Error("load database target mappings", "error", err)
		}
	}
	writer.wg.Go(writer.refreshLoop)
	writer.wg.Go(writer.superviseLoop)
	for id := range writer.indexedConnectorIDs() {
		writer.workerFor(id)
	}
	return writer
}

// ErrOutputOwnedByWriteGroup remains available to legacy callers. The snapshot
// writer checks ownership during the worker's flush, after buffer acceptance.
var ErrOutputOwnedByWriteGroup = errors.New("output is owned by a write group")

// ErrSnapshotAccepted indicates that the value is buffered, not committed.
var ErrSnapshotAccepted = errors.New("snapshot accepted")

// ErrOutputNotApplicable indicates that no enabled mapping accepted the value.
var ErrOutputNotApplicable = errors.New("no database output accepted the value")

// WriteTagValue returns ErrSnapshotAccepted when a value is buffered. Callers
// must not treat acceptance as a successful SQL write.
func (w *Writer) WriteTagValue(ctx context.Context, tagID string, value any, observedAt time.Time) error {
	accepted, err := w.AcceptTagValue(ctx, tagID, value, observedAt)
	if err != nil {
		return err
	}
	if accepted {
		return ErrSnapshotAccepted
	}
	return ErrOutputNotApplicable
}

// AcceptTagValue buffers the latest value without database I/O. The returned
// boolean indicates acceptance, not whether a destination committed it.
func (w *Writer) AcceptTagValue(_ context.Context, tagID string, value any, observedAt time.Time) (bool, error) {
	if strings.TrimSpace(tagID) == "" {
		return false, fmt.Errorf("tagID 不可為空")
	}

	entries, err := w.lookupMappings(tagID)
	if err != nil {
		return false, err
	}
	if len(entries) == 0 {
		return false, nil
	}

	var failures []error
	accepted := false
	for _, entry := range entries {
		mapping := entry.Mapping
		if entry.Err != nil {
			failures = append(failures, fmt.Errorf("mapping %s 取得連接器失敗: %w", mapping.ID, entry.Err))
			continue
		}
		if hasGroupedWriteKey(mapping) {
			if err := w.bufferGroupedWrite(entry.Connector, mapping, value, observedAt); err != nil {
				failures = append(failures, fmt.Errorf("mapping %s 分組緩衝失敗: %w", mapping.ID, err))
			} else {
				accepted = true
			}
			continue
		}
		if err := w.storeLatest(entry, value, observedAt); err != nil {
			failures = append(failures, fmt.Errorf("mapping %s 緩衝失敗: %w", mapping.ID, err))
		} else {
			accepted = true
		}
	}
	return accepted, errors.Join(failures...)
}

func (w *Writer) writeMapping(
	ctx context.Context,
	connector *schema.DatabaseConnector,
	mapping *schema.DatabaseTargetMapping,
	value any,
	observedAt time.Time,
) error {
	connectionConfig, err := writerConnectionConfig(connector)
	if err != nil {
		return err
	}

	// 執行期寫入不得自動建庫：目標資料庫消失屬於必須浮現的錯誤，
	// 靜默重建會讓資料寫進一個空資料庫。
	manager, err := w.openDestination(ctx, connector, connectionConfig)
	if err != nil {
		return err
	}
	defer manager.Close()

	query, args, err := buildWriteStatement(connector.Kind, mapping, value, observedAt)
	if err != nil {
		return err
	}

	if err := w.executeWithTableRecovery(ctx, manager.DB(), connector, mapping.TableSchema, mapping.TableName, query, args); err != nil {
		return fmt.Errorf("執行資料庫寫入失敗: %w", err)
	}

	return nil
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
	if w.closed {
		return errors.New("database target writer is closed")
	}
	interval := time.Duration(key.EffectiveIntervalSeconds) * time.Second
	if previous, ok := w.groupedLocations[mapping.ID]; ok {
		if old := w.grouped[previous.key]; old != nil {
			oldMapping := old.Mappings[previous.column]
			if sameSnapshotBinding(oldMapping, old.Connector, mapping, connector) {
				if old.Retry {
					key = w.movePendingGroup(previous.key, key, connector)
				} else if previous.key != key {
					if !w.retainGroupedInterval(previous, key, mapping, connector) {
						return nil
					}
				}
			} else {
				w.discardChangedGrouped(mapping, connector)
			}
		}
	}

	bucket, ok := w.grouped[key]
	if ok && !sameConnectorIdentity(bucket.Connector, connector) {
		for column, oldMapping := range bucket.Mappings {
			w.recordMappingSkip(bucket.Connector, oldMapping, bucket.Key.BucketStart, skipMappingNotCurrent)
			if latest, exists := w.groupedLocations[oldMapping.ID]; exists && latest == (groupedLocation{key: key, column: column}) {
				delete(w.groupedLocations, oldMapping.ID)
			}
		}
		delete(w.grouped, key)
		bucket, ok = nil, false
	}
	if !ok {
		bucket = &groupedWriteBucket{
			Key:       key,
			Connector: connector,
			Values:    make(map[string]any),
			Mappings:  make(map[string]*schema.DatabaseTargetMapping),
			DueAt:     key.BucketStart.Add(interval),
		}
		w.grouped[key] = bucket
	}
	if previous := bucket.Mappings[mapping.ColumnName]; previous != nil && previous.ID != mapping.ID {
		w.recordMappingSkip(bucket.Connector, previous, bucket.Key.BucketStart, skipColumnReplaced)
		if latest, exists := w.groupedLocations[previous.ID]; exists && latest == (groupedLocation{key: key, column: mapping.ColumnName}) {
			delete(w.groupedLocations, previous.ID)
		}
	}
	bucket.Values[mapping.ColumnName] = normalizeDBValue(value)
	bucket.Mappings[mapping.ColumnName] = mapping
	if previous, ok := w.groupedLocations[mapping.ID]; !ok || !key.BucketStart.Before(previous.key.BucketStart) {
		w.groupedLocations[mapping.ID] = groupedLocation{key: key, column: mapping.ColumnName}
	}
	w.markObserved(mapping.ID, key.EffectiveIntervalSeconds, key.BucketStart)
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

func (w *Writer) flushGroupedBucket(ctx context.Context, bucket *groupedWriteBucket) error {
	connectionConfig, err := writerConnectionConfig(bucket.Connector)
	if err != nil {
		return err
	}

	// 同上：批次寫入路徑同樣不自動建庫。
	manager, err := w.openDestination(ctx, bucket.Connector, connectionConfig)
	if err != nil {
		return err
	}
	defer manager.Close()

	query, args, err := buildGroupedWriteStatement(bucket.Connector.Kind, bucket.Key, bucket.Values)
	if err != nil {
		return err
	}
	if err := w.executeWithTableRecovery(ctx, manager.DB(), bucket.Connector, bucket.Key.SchemaName, bucket.Key.TableName, query, args); err != nil {
		return fmt.Errorf("執行資料庫分組寫入失敗: %w", err)
	}
	return nil
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

var _ interface {
	WriteTagValue(context.Context, string, any, time.Time) error
} = (*Writer)(nil)
