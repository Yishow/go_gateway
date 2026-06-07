package runtime

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
)

// Config runtime 設定。
type Config struct {
	// UpdatePointState 控制是否更新 points.last_* 欄位。
	UpdatePointState bool

	// Writer 舊有相容欄位（snapshot 模式）。
	Writer storage.Writer

	// Snapshot 舊有相容欄位（snapshot 模式）。
	Snapshot Snapshot
}

// DefaultConfig 預設設定。
func DefaultConfig() Config {
	return Config{UpdatePointState: true}
}

// Dependencies runtime 依賴（service 模式）。
type Dependencies struct {
	Scheduler           *collector.Scheduler
	Writer              storage.Writer
	TargetWriter        TargetWriter
	WorkspaceProjection WorkspaceProjectionReader
	DeviceService       *device.Service
	PointService        *point.Service
	MappingService      *mapping.Service
	TagService          *tag.Service
	PollingGroupService *pollinggroup.Service
}

type TargetWriter interface {
	WriteTagValue(ctx context.Context, tagID string, value any, observedAt time.Time) error
}

// WorkspaceProjectionReader rebuilds persisted workspace scope for runtime bootstrap.
type WorkspaceProjectionReader interface {
	RuntimeProjection(ctx context.Context) (*workspace.RuntimeProjection, error)
}

// Stats runtime 指標。
type Stats struct {
	CollectedTotal  uint64 `json:"collected_total"`
	WriteSuccess    uint64 `json:"write_success_total"`
	WriteError      uint64 `json:"write_error_total"`
	MappingError    uint64 `json:"mapping_error_total"`
	PointStateError uint64 `json:"point_state_error_total"`
}

type mappingBinding struct {
	TagID             string
	TagDataType       schema.DataType
	TransformPipeline string
}

// Service 負責串接 scheduler -> mapping -> storage。
type Service struct {
	config Config

	scheduler *collector.Scheduler
	writer    storage.Writer
	target    TargetWriter

	deviceSvc  *device.Service
	pointSvc   *point.Service
	mappingSvc *mapping.Service
	tagSvc     *tag.Service
	groupSvc   *pollinggroup.Service
	workspace  WorkspaceProjectionReader

	// snapshot mode 使用
	snapshot Snapshot

	mappingMu      sync.RWMutex
	mappingIndex   map[string][]mappingBinding // pointID -> mappings
	pointMetaMu    sync.RWMutex
	pointMetaIndex map[string]pointMeta
	projectionMu   sync.RWMutex
	projections    map[string]runtimeProjectionDeviceState

	subscriberMu       sync.RWMutex
	subscribers        map[int64]valueSubscriber
	statusSubscriberMu sync.RWMutex
	statusSubscribers  map[int64]statusSubscriber
	nextSubscriberID   atomic.Int64
	lastStatusMu       sync.Mutex
	lastStatuses       map[string]DeviceRuntimeStatus

	stopCh    chan struct{}
	wg        sync.WaitGroup
	running   atomic.Bool
	startedAt atomic.Int64

	collectedTotal  atomic.Uint64
	writeSuccess    atomic.Uint64
	writeError      atomic.Uint64
	mappingError    atomic.Uint64
	pointStateError atomic.Uint64
}

// NewService 建立 runtime service。
// 支援兩種模式：
// 1) service mode: NewService(cfg, deps)
// 2) snapshot mode: NewService(Config{Writer:..., Snapshot:...})
func NewService(config Config, depsOpt ...Dependencies) (*Service, error) {
	s := &Service{
		config:            config,
		mappingIndex:      make(map[string][]mappingBinding),
		pointMetaIndex:    make(map[string]pointMeta),
		projections:       make(map[string]runtimeProjectionDeviceState),
		subscribers:       make(map[int64]valueSubscriber),
		statusSubscribers: make(map[int64]statusSubscriber),
		lastStatuses:      make(map[string]DeviceRuntimeStatus),
		stopCh:            make(chan struct{}),
	}

	if len(depsOpt) > 0 {
		deps := depsOpt[0]
		if deps.Scheduler == nil {
			return nil, fmt.Errorf("scheduler 不可為 nil")
		}
		if deps.Writer == nil {
			return nil, fmt.Errorf("writer 不可為 nil")
		}
		if deps.DeviceService == nil || deps.PointService == nil || deps.MappingService == nil || deps.TagService == nil || deps.PollingGroupService == nil {
			return nil, fmt.Errorf("service 依賴不可為 nil")
		}
		s.scheduler = deps.Scheduler
		s.writer = deps.Writer
		s.target = deps.TargetWriter
		s.workspace = deps.WorkspaceProjection
		s.deviceSvc = deps.DeviceService
		s.pointSvc = deps.PointService
		s.mappingSvc = deps.MappingService
		s.tagSvc = deps.TagService
		s.groupSvc = deps.PollingGroupService
		s.pointSvc.SetPollingGroupSyncer(s.scheduler)
		return s, nil
	}

	// snapshot mode
	if config.Writer == nil {
		return nil, fmt.Errorf("writer 不可為 nil")
	}
	s.writer = config.Writer
	s.snapshot = config.Snapshot

	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	s.scheduler = collector.NewScheduler(collector.DefaultSchedulerConfig(), connMgr)

	return s, nil
}

// Start 啟動 runtime。
func (s *Service) Start(ctx context.Context) error {
	if !s.running.CompareAndSwap(false, true) {
		return fmt.Errorf("runtime 已在運行")
	}
	s.startedAt.Store(time.Now().UnixNano())

	if err := s.bootstrap(ctx); err != nil {
		s.running.Store(false)
		s.startedAt.Store(0)
		return err
	}

	s.wg.Add(1)
	go s.consumeLoop()
	s.wg.Add(1)
	go s.statusLoop()
	return nil
}

// Stop 停止 runtime。
func (s *Service) Stop(ctx context.Context) error {
	if !s.running.CompareAndSwap(true, false) {
		return nil
	}
	s.startedAt.Store(0)

	close(s.stopCh)
	s.scheduler.Stop()
	s.wg.Wait()

	if err := s.writer.Flush(ctx); err != nil {
		return err
	}
	if closer, ok := s.target.(interface{ Close(context.Context) error }); ok && closer != nil {
		if err := closer.Close(ctx); err != nil {
			return err
		}
	}
	return s.writer.Close()
}

// RefreshMappings 重新載入 mapping 快取。
func (s *Service) RefreshMappings(ctx context.Context) error {
	return s.refreshMappings(ctx)
}

// Snapshot 取得統計快照。
func (s *Service) Snapshot() Stats {
	return Stats{
		CollectedTotal:  s.collectedTotal.Load(),
		WriteSuccess:    s.writeSuccess.Load(),
		WriteError:      s.writeError.Load(),
		MappingError:    s.mappingError.Load(),
		PointStateError: s.pointStateError.Load(),
	}
}

// Metrics 舊有相容指標。
func (s *Service) Metrics() Metrics {
	st := s.Snapshot()
	errTotal := int64(st.WriteError + st.MappingError + st.PointStateError)
	return Metrics{
		TotalReads:        int64(st.CollectedTotal),
		TotalWrites:       int64(st.WriteSuccess),
		ErrorCount:        errTotal,
		MappingErrorCount: int64(st.MappingError),
		WriteErrorCount:   int64(st.WriteError),
	}
}

// IsRunning reports whether runtime service is actively running.
func (s *Service) IsRunning() bool {
	return s.running.Load()
}

// SetWorkspaceProjectionReader sets the persisted workspace projection reader used by runtime restart.
func (s *Service) SetWorkspaceProjectionReader(reader WorkspaceProjectionReader) {
	if s == nil {
		return
	}
	s.workspace = reader
}

// UptimeSeconds returns elapsed runtime seconds since the latest successful start.
func (s *Service) UptimeSeconds() int64 {
	if !s.IsRunning() {
		return 0
	}

	startedAt := s.startedAt.Load()
	if startedAt <= 0 {
		return 0
	}

	return int64(time.Since(time.Unix(0, startedAt)).Seconds())
}

// UpsertPoint updates runtime metadata and scheduler state for a point.
func (s *Service) UpsertPoint(point *schema.Point) {
	if point == nil {
		return
	}

	s.registerPointMeta(point.ID, pointMeta{
		DeviceID: point.DeviceID,
		Address:  point.Address,
	})

	if s.scheduler == nil {
		return
	}

	if point.Enabled {
		s.scheduler.AddPoint(point)
		return
	}
	s.scheduler.RemovePoint(point.ID)
}

// RemovePoint removes runtime metadata and scheduler state for a point.
func (s *Service) RemovePoint(pointID string) {
	s.pointMetaMu.Lock()
	delete(s.pointMetaIndex, pointID)
	s.pointMetaMu.Unlock()

	if s.scheduler != nil {
		s.scheduler.RemovePoint(pointID)
	}
}

func (s *Service) bootstrap(ctx context.Context) error {
	if s.deviceSvc != nil {
		return s.bootstrapFromServices(ctx)
	}
	return s.bootstrapFromSnapshot(ctx)
}

func (s *Service) bootstrapFromServices(ctx context.Context) error {
	if s.workspace != nil {
		projection, err := s.workspace.RuntimeProjection(ctx)
		if err != nil {
			return fmt.Errorf("載入 workspace runtime projection 失敗: %w", err)
		}
		return s.bootstrapFromWorkspaceProjection(ctx, projection)
	}

	active := schema.DeviceStatusActive
	devices, err := s.deviceSvc.List(ctx, device.ListFilter{Status: &active, Limit: 100000})
	if err != nil {
		return fmt.Errorf("載入 devices 失敗: %w", err)
	}
	for _, d := range devices {
		s.scheduler.AddDevice(d)
	}

	points, err := s.pointSvc.List(ctx, point.ListFilter{Limit: 100000})
	if err != nil {
		return fmt.Errorf("載入 points 失敗: %w", err)
	}
	for _, p := range points {
		s.registerPointMeta(p.ID, pointMeta{
			DeviceID: p.DeviceID,
			Address:  p.Address,
		})
		if !p.Enabled {
			continue
		}
		s.scheduler.AddPoint(p)
	}

	if err := s.refreshMappings(ctx); err != nil {
		return err
	}

	groups, err := s.groupSvc.List(ctx)
	if err != nil {
		return fmt.Errorf("載入 polling groups 失敗: %w", err)
	}
	if err := s.scheduler.Start(groups); err != nil {
		return fmt.Errorf("啟動 scheduler 失敗: %w", err)
	}
	return nil
}

func (s *Service) bootstrapFromSnapshot(_ context.Context) error {
	for _, d := range s.snapshot.Devices {
		if d.Status != schema.DeviceStatusActive {
			continue
		}
		s.scheduler.AddDevice(d)
	}
	for _, p := range s.snapshot.Points {
		s.registerPointMeta(p.ID, pointMeta{
			DeviceID: p.DeviceID,
			Address:  p.Address,
		})
		if !p.Enabled {
			continue
		}
		s.scheduler.AddPoint(p)
	}

	if err := s.refreshMappings(context.Background()); err != nil {
		return err
	}

	if err := s.scheduler.Start(s.snapshot.PollingGroups); err != nil {
		return fmt.Errorf("啟動 scheduler 失敗: %w", err)
	}
	return nil
}

func (s *Service) refreshMappings(ctx context.Context) error {
	if s.mappingSvc != nil {
		enabled := true
		mappings, err := s.mappingSvc.List(ctx, mapping.ListFilter{Enabled: &enabled, Limit: 100000})
		if err != nil {
			return fmt.Errorf("載入 mappings 失敗: %w", err)
		}

		next := make(map[string][]mappingBinding)
		for _, m := range mappings {
			tg, err := s.tagSvc.GetByID(ctx, m.TagID)
			if err != nil {
				return fmt.Errorf("載入 tag(%s) 失敗: %w", m.TagID, err)
			}
			next[m.PointID] = append(next[m.PointID], mappingBinding{
				TagID:             m.TagID,
				TagDataType:       tg.DataType,
				TransformPipeline: m.TransformPipeline,
			})
		}

		s.mappingMu.Lock()
		s.mappingIndex = next
		s.mappingMu.Unlock()
		return nil
	}

	// snapshot mode
	tagType := make(map[string]schema.DataType, len(s.snapshot.Tags))
	for _, t := range s.snapshot.Tags {
		tagType[t.ID] = t.DataType
	}

	next := make(map[string][]mappingBinding)
	for _, m := range s.snapshot.Mappings {
		if !m.Enabled {
			continue
		}
		next[m.PointID] = append(next[m.PointID], mappingBinding{
			TagID:             m.TagID,
			TagDataType:       tagType[m.TagID],
			TransformPipeline: m.TransformPipeline,
		})
	}

	s.mappingMu.Lock()
	s.mappingIndex = next
	s.mappingMu.Unlock()
	return nil
}

func (s *Service) consumeLoop() {
	defer s.wg.Done()

	for {
		select {
		case <-s.stopCh:
			return
		case cv := <-s.scheduler.ValueChannel():
			s.collectedTotal.Add(1)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			s.handleCollectedValue(ctx, cv)
			cancel()
		}
	}
}
