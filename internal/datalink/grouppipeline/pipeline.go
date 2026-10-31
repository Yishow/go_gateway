package grouppipeline

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

// GroupSource reads canonical groups and their immutable applied snapshots.
type GroupSource interface {
	List(ctx context.Context) (*workspace.WriteGroupListResult, error)
	ResolveAppliedAt(ctx context.Context, id string, at time.Time) (*workspace.WriteGroupAppliedSnapshot, error)
}

// TagReader resolves persisted tag data types.
type TagReader interface {
	GetByID(ctx context.Context, id string) (*schema.Tag, error)
}

// Destinations is everything the pipeline needs from saved connectors.
type Destinations interface {
	GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error)
	OpenDestination(ctx context.Context, connectorID, expectedRevision string) (*dbtarget.OpenedDestination, error)
}

// TableInspector reads real destination table metadata without changing it.
type TableInspector interface {
	InspectTable(ctx context.Context, connectorID, schemaName, tableName string) (*dbtarget.TableInspection, error)
}

// Dependencies wires the pipeline to the production services.
type Dependencies struct {
	// OwnerObservation is optional diagnostic capture, never ownership authority.
	OwnerObservation       OwnerObservationConfig
	IntakeQuotaObservation IntakeQuotaObservationConfig
	StorageObservation     StorageObservationConfig
	Groups                 GroupSource
	Tags                   TagReader
	Destinations           Destinations
	Inspector              TableInspector
	Store                  *groupdelivery.Store
	// Counters confirms the counter semantics of cumulative groups' members.
	// Without it a cumulative group is not built.
	Counters CounterSources
	// Ledger coordinates confirmed managed-storage recovery across scopes.
	Ledger *recordingplan.Service
	// SQLServerSplit is the typed dbtarget native partition executor. It is
	// wired from the root integration path; nil keeps SQL Server maintenance
	// unavailable without changing other dialects.
	SQLServerSplit SQLServerNativeSplit
}

// Config bounds the pipeline; zero values take the documented defaults.
type Config struct {
	NodeID string
	Owner  string
	// TickInterval drives time-based bucket closure; ReconcileInterval picks up
	// newly applied, superseded and disabled groups.
	TickInterval      time.Duration
	ReconcileInterval time.Duration
	MaxFutureSkew     time.Duration
	Delivery          groupdelivery.DispatcherConfig
	Sender            groupdelivery.SenderConfig
	WorkerInterval    time.Duration
	// PartitionMaintenanceInterval is a background cadence, not a row interval.
	PartitionMaintenanceInterval time.Duration
	// ReclaimEvery and ReclaimRetention bound how long finished data (committed
	// rows, consumed samples, bucket status) stays in the local database.
	ReclaimEvery     time.Duration
	ReclaimRetention time.Duration
	Now              func() time.Time
	OnError          func(error)
	// OnDiagnosticError emits only a fixed generic diagnostic, never a native error.
	OnDiagnosticError func()
}

// Defaults. They are bounds and cadences, not measured performance targets.
const (
	defaultTickInterval      = time.Second
	defaultReconcileInterval = 5 * time.Second
	defaultMaxFutureSkew     = 5 * time.Second
	defaultWorkerInterval    = time.Second
	defaultReclaimEvery      = 10 * time.Minute
	defaultReclaimRetention  = 24 * time.Hour
	receiptTableName         = dbtarget.EffectReceiptTable
)

// Group states reported by Status.
const (
	StateActive   = "active"
	StateRetiring = "retiring"
	StateBlocked  = "blocked"
)

// GroupStatus describes one group's intake side. Delivery state lives in the
// durable outbox and is read from the store.
type GroupStatus struct {
	GroupID  string
	Revision string
	State    string
	// Reason is a safe code, never a value, DSN or driver text.
	Reason string
}

type member struct{ tagID, connectorID string }

type managed struct {
	// proofGate precedes Pipeline.mu; never wait for it while holding Pipeline.mu.
	proofGate        sync.Mutex
	proofQuarantined error
	superseded       bool // protected by proofGate; captured old workers must not run after replacement
	workspaceID      string
	groupID          string
	revision         string
	boundary         *runtime.GroupBoundary
	members          []member
	connector        string
	interval         time.Duration
	mixed            bool
	effective        time.Time
	drainOnly        bool
	localDrain       atomic.Bool // local intake fence, never a persisted revision end
	drainFinished    atomic.Bool
	drain            *mixedDrain // protected by proofGate
}

// Pipeline owns the applied write groups of one gateway.
type Pipeline struct {
	// Serializes cache reconciliation, never sample/closure SQL authority.
	reconcileMu    sync.Mutex
	deps           Dependencies
	config         Config
	recovery       *StorageRecovery
	partitionMu    sync.Mutex
	sqlServerSplit SQLServerNativeSplit

	mu          sync.RWMutex
	boundaries  map[string]*managed
	blocked     map[string]GroupStatus
	quarantines map[groupdelivery.GroupKey]CounterReplayQuarantine

	lifecycle   sync.Mutex
	started     bool
	startDone   chan struct{}
	startCancel context.CancelFunc
	stopping    bool
	stopDone    chan struct{}
	intakeWG    sync.WaitGroup
	cancel      context.CancelFunc
	loopDone    chan struct{}
	worker      *groupdelivery.Worker
}

// New builds a pipeline. Start must be called before samples are accepted.
func New(deps Dependencies, config Config) *Pipeline {
	if config.TickInterval <= 0 {
		config.TickInterval = defaultTickInterval
	}
	if config.ReconcileInterval <= 0 {
		config.ReconcileInterval = defaultReconcileInterval
	}
	if config.PartitionMaintenanceInterval <= 0 {
		config.PartitionMaintenanceInterval = time.Minute
	}
	if config.MaxFutureSkew <= 0 {
		config.MaxFutureSkew = defaultMaxFutureSkew
	}
	if config.WorkerInterval <= 0 {
		config.WorkerInterval = defaultWorkerInterval
	}
	if config.ReclaimEvery <= 0 {
		config.ReclaimEvery = defaultReclaimEvery
	}
	if config.ReclaimRetention <= 0 {
		config.ReclaimRetention = defaultReclaimRetention
	}
	if config.Now == nil {
		config.Now = func() time.Time { return time.Now().UTC() }
	}
	if config.NodeID == "" {
		config.NodeID = "gateway"
	}
	if config.Owner == "" {
		config.Owner = config.NodeID + "/" + fmt.Sprintf("%d", config.Now().UnixNano())
	}
	var recovery *StorageRecovery
	if deps.Store != nil && deps.Store.DB() != nil {
		ledger := deps.Ledger
		if ledger == nil {
			ledger = recordingplan.NewService(recordingplan.NewSQLRepository(deps.Store.DB()))
		}
		recovery = NewStorageRecoveryWithObservation(deps.Store, ledger, deps.Destinations, deps.StorageObservation)
		recovery.report = func(err error) { reportPipelineError(config, err) }
		if config.Sender.Recovery == nil {
			config.Sender.Recovery = recovery
		}
	}
	sqlServerSplit := deps.SQLServerSplit
	if sqlServerSplit == nil {
		sqlServerSplit = dbtarget.SplitOpenSQLServerNativeMonthly
	}
	return &Pipeline{
		deps: deps, config: config, recovery: recovery,
		boundaries:     make(map[string]*managed),
		blocked:        make(map[string]GroupStatus),
		quarantines:    make(map[groupdelivery.GroupKey]CounterReplayQuarantine),
		sqlServerSplit: sqlServerSplit,
	}
}

// Start reconciles the applied groups, recovers delivery claims left by an
// earlier incarnation and begins ticking and delivering. A group that cannot
// be built yet (destination offline, table missing) is reported as blocked and
// retried; it never prevents the others from starting.
func (p *Pipeline) Start(ctx context.Context) error {
	p.lifecycle.Lock()
	if p.started {
		p.lifecycle.Unlock()
		return errors.New("group pipeline already started")
	}
	if p.stopping {
		select {
		case <-p.stopDone:
		default:
			p.lifecycle.Unlock()
			return errors.New("group pipeline is still stopping")
		}
	}
	p.stopDone, p.stopping, p.started = nil, false, true
	p.loopDone, p.cancel, p.worker = nil, nil, nil
	p.startDone = make(chan struct{})
	startCtx, startCancel := context.WithCancel(ctx)
	p.startCancel = startCancel
	p.lifecycle.Unlock()
	defer startCancel()
	if err := p.Reconcile(startCtx); err != nil {
		p.report(err)
	}
	if err := startCtx.Err(); err != nil {
		p.lifecycle.Lock()
		p.started = false
		close(p.startDone)
		p.lifecycle.Unlock()
		return err
	}
	sender := groupdelivery.NewSender(p.deps.Store, destinationResolver{
		destinations: p.deps.Destinations, loadRuntimeSnapshot: p.loadFinalizedSnapshot,
		loadMySQLProof:     p.mysqlDeliveryProof,
		loadSQLServerProof: p.sqlServerDeliveryProof,
	}, p.senderConfig())
	dispatcher := groupdelivery.NewDispatcher(p.deps.Store, sender, p.config.Delivery)
	worker := groupdelivery.NewWorker(p.deps.Store, dispatcher, groupdelivery.WorkerConfig{
		Interval: p.config.WorkerInterval, NodeID: p.config.NodeID, Owner: p.config.Owner, OnError: p.report,
		ReclaimEvery: p.config.ReclaimEvery, ReclaimRetention: p.config.ReclaimRetention,
	})
	p.lifecycle.Lock()
	defer p.lifecycle.Unlock()
	defer close(p.startDone)
	if p.stopping {
		return errors.New("group pipeline stopped during startup")
	}
	// Starting the worker's recovery may perform IO; completion remains represented
	// by startDone while StopContext can still publish the caller's deadline.
	loopCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	p.cancel, p.worker = cancel, worker
	p.lifecycle.Unlock()
	err := worker.StartWithLifetime(startCtx, loopCtx)
	p.lifecycle.Lock()
	if err != nil {
		p.started = false
		cancel()
		return err
	}
	if err := startCtx.Err(); err != nil {
		p.beginStopLocked(startCtx)
		return err
	}
	if p.stopping {
		return errors.New("group pipeline stopped during startup")
	}
	p.loopDone = make(chan struct{})
	go p.loop(loopCtx, p.loopDone)
	return nil
}

func (p *Pipeline) senderConfig() groupdelivery.SenderConfig {
	config := p.config.Sender
	config.Owner = p.config.Owner
	if config.Recovery == nil && p.recovery != nil {
		config.Recovery = p.recovery
	}
	return config
}

func (p *Pipeline) loop(ctx context.Context, done chan struct{}) {
	defer close(done)
	tick := time.NewTicker(p.config.TickInterval)
	defer tick.Stop()
	reconcile := time.NewTicker(p.config.ReconcileInterval)
	defer reconcile.Stop()
	maintenance := time.NewTicker(p.config.PartitionMaintenanceInterval)
	defer maintenance.Stop()
	var flight chan struct{}
	startMaintenance := func() {
		if flight != nil {
			select {
			case <-flight:
			default:
				return
			}
		}
		flight = make(chan struct{})
		go func(done chan struct{}) {
			defer close(done)
			p.report(p.MaintainMySQLPartitions(ctx))
			p.report(p.MaintainSQLServerPartitions(ctx))
		}(flight)
	}
	defer func() {
		if flight != nil {
			<-flight
		}
	}()
	// A restore whose owner never came back is closed from the catalog here, not
	// in Start: it may wait on the local database and must not delay startup.
	p.resolveStaleStorage(ctx)
	startMaintenance()
	for {
		select {
		case <-ctx.Done():
			return
		case <-reconcile.C:
			p.report(p.Reconcile(ctx))
			p.resolveStaleStorage(ctx)
		case <-tick.C:
			p.TickAll(ctx)
		case <-maintenance.C:
			startMaintenance()
		}
	}
}

func (p *Pipeline) resolveStaleStorage(ctx context.Context) {
	if p.recovery != nil {
		p.report(p.recovery.ResolveStale(ctx))
	}
}

// TickAll closes every due bucket of every boundary and drops boundaries that
// finished their interval.
func (p *Pipeline) TickAll(ctx context.Context) {
	ctx, finish := p.ownerOperation(ctx)
	defer finish()
	p.mu.RLock()
	current := make([]*managed, 0, len(p.boundaries))
	for _, m := range p.boundaries {
		current = append(current, m)
	}
	p.mu.RUnlock()
	now := p.config.Now()
	for _, m := range current {
		p.reportTick(ctx, m, func() error { return p.tickManaged(ctx, m, now) })
	}
	p.pruneRetired(current)
}

func (p *Pipeline) pruneRetired(current []*managed) {
	for _, m := range current {
		// Retired can wait for a boundary's SQL-backed Tick; never keep the
		// global map locked while an unrelated worker needs it for intake.
		if !m.drainFinished.Load() && !m.boundary.Retired() {
			continue
		}
		p.mu.Lock()
		key := boundaryKey(m.groupID, m.revision)
		if p.boundaries[key] == m {
			delete(p.boundaries, key)
		}
		p.mu.Unlock()
	}
}

// AcceptSample implements the runtime's typed sample sink. Committed cutover
// authority chooses affected owners independently of reconciliation. Refusals
// are joined while unaffected scopes continue their existing intake contract.
func (p *Pipeline) AcceptSample(ctx context.Context, envelope measurement.SampleEnvelope) error {
	ctx, finish := p.ownerOperation(ctx)
	defer finish()
	p.lifecycle.Lock()
	if p.stopping {
		p.lifecycle.Unlock()
		return errors.New("group pipeline is stopping")
	}
	p.intakeWG.Add(1)
	p.lifecycle.Unlock()
	defer p.intakeWG.Done()
	p.mu.RLock()
	current := make([]*managed, 0, len(p.boundaries))
	for _, m := range p.boundaries {
		current = append(current, m)
	}
	p.mu.RUnlock()
	return p.acceptWithOwnership(ctx, envelope, current)
}

// WantsSample retains typed transport interest for committed cutover sources
// even while their worker is absent. Interest never authorizes acceptance.
func (p *Pipeline) WantsSample(deviceID, pointID, tagID string) bool {
	p.mu.RLock()
	for _, m := range p.boundaries {
		if !m.drainOnly && !m.localDrain.Load() && m.boundary.Wants(deviceID, pointID, tagID) {
			p.mu.RUnlock()
			return true
		}
	}
	p.mu.RUnlock()
	return p.wantsCutoverSample(deviceID, pointID, tagID)
}

// Owns includes durable successful Apply history, independent of active
// intake workers. Ownership read failure suppresses legacy output fail closed.
func (p *Pipeline) Owns(connectorID, tagID string) bool {
	if p.deps.Store != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		owned, err := p.deps.Store.HasAppliedOwner(ctx, connectorID, tagID)
		cancel()
		if err != nil {
			p.report(err)
			return true
		}
		if owned {
			return true
		}
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, m := range p.boundaries {
		for _, mem := range m.members {
			if mem.tagID == tagID && mem.connectorID == connectorID {
				return true
			}
		}
	}
	return false
}

// Status lists the intake state of every group the pipeline knows about.
func (p *Pipeline) Status() []GroupStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()
	statuses := make([]GroupStatus, 0, len(p.boundaries)+len(p.blocked))
	for _, m := range p.boundaries {
		state := StateActive
		if m.localDrain.Load() || m.boundary.Retiring() {
			state = StateRetiring
		}
		statuses = append(statuses, GroupStatus{GroupID: m.groupID, Revision: m.revision, State: state})
	}
	for _, status := range p.blocked {
		statuses = append(statuses, status)
	}
	for _, entry := range p.quarantines {
		statuses = append(statuses, GroupStatus{GroupID: entry.GroupID, Revision: entry.GroupRevision, State: StateBlocked, Reason: entry.Reason})
	}
	return statuses
}

type destinationResolver struct {
	destinations        Destinations
	loadRuntimeSnapshot func(context.Context, groupdelivery.GroupKey) (*workspace.WriteGroupAppliedSnapshot, error)
	loadMySQLProof      func(context.Context, *workspace.WriteGroupAppliedSnapshot) (*dbtarget.MySQLManagedRowProof, error)
	loadSQLServerProof  func(context.Context, *workspace.WriteGroupAppliedSnapshot, *dbtarget.OpenedDestination) (*dbtarget.SQLServerManagedRowProof, error)
}

var (
	_ runtime.SampleSink     = (*Pipeline)(nil)
	_ runtime.SampleInterest = (*Pipeline)(nil)
)
