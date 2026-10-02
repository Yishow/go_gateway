package grouppipeline

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
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
	Groups       GroupSource
	Tags         TagReader
	Destinations Destinations
	Inspector    TableInspector
	Store        *groupdelivery.Store
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
	// ReclaimEvery and ReclaimRetention bound how long finished data (committed
	// rows, consumed samples, bucket status) stays in the local database.
	ReclaimEvery     time.Duration
	ReclaimRetention time.Duration
	Now              func() time.Time
	OnError          func(error)
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
	groupID   string
	revision  string
	boundary  *runtime.GroupBoundary
	members   []member
	connector string
	interval  time.Duration
}

// Pipeline owns the applied write groups of one gateway.
type Pipeline struct {
	deps   Dependencies
	config Config

	mu         sync.RWMutex
	boundaries map[string]*managed
	blocked    map[string]GroupStatus

	lifecycle sync.Mutex
	started   bool
	cancel    context.CancelFunc
	loopDone  chan struct{}
	worker    *groupdelivery.Worker
}

// New builds a pipeline. Start must be called before samples are accepted.
func New(deps Dependencies, config Config) *Pipeline {
	if config.TickInterval <= 0 {
		config.TickInterval = defaultTickInterval
	}
	if config.ReconcileInterval <= 0 {
		config.ReconcileInterval = defaultReconcileInterval
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
	return &Pipeline{
		deps: deps, config: config,
		boundaries: make(map[string]*managed),
		blocked:    make(map[string]GroupStatus),
	}
}

func (p *Pipeline) report(err error) {
	if err != nil && p.config.OnError != nil {
		p.config.OnError(err)
	}
}

// Start reconciles the applied groups, recovers delivery claims left by an
// earlier incarnation and begins ticking and delivering. A group that cannot
// be built yet (destination offline, table missing) is reported as blocked and
// retried; it never prevents the others from starting.
func (p *Pipeline) Start(ctx context.Context) error {
	p.lifecycle.Lock()
	defer p.lifecycle.Unlock()
	if p.started {
		return errors.New("group pipeline already started")
	}
	if err := p.Reconcile(ctx); err != nil {
		p.report(err)
	}
	sender := groupdelivery.NewSender(p.deps.Store, destinationResolver{destinations: p.deps.Destinations}, p.senderConfig())
	dispatcher := groupdelivery.NewDispatcher(p.deps.Store, sender, p.config.Delivery)
	p.worker = groupdelivery.NewWorker(p.deps.Store, dispatcher, groupdelivery.WorkerConfig{
		Interval: p.config.WorkerInterval, NodeID: p.config.NodeID, Owner: p.config.Owner, OnError: p.report,
		ReclaimEvery: p.config.ReclaimEvery, ReclaimRetention: p.config.ReclaimRetention,
	})
	if err := p.worker.Start(ctx); err != nil {
		return err
	}
	loopCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	p.cancel, p.loopDone, p.started = cancel, make(chan struct{}), true
	go p.loop(loopCtx, p.loopDone)
	return nil
}

func (p *Pipeline) senderConfig() groupdelivery.SenderConfig {
	config := p.config.Sender
	config.Owner = p.config.Owner
	return config
}

func (p *Pipeline) loop(ctx context.Context, done chan struct{}) {
	defer close(done)
	tick := time.NewTicker(p.config.TickInterval)
	defer tick.Stop()
	reconcile := time.NewTicker(p.config.ReconcileInterval)
	defer reconcile.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-reconcile.C:
			p.report(p.Reconcile(ctx))
		case <-tick.C:
			p.TickAll(ctx)
		}
	}
}

// TickAll closes every due bucket of every boundary and drops boundaries that
// finished their interval.
func (p *Pipeline) TickAll(ctx context.Context) {
	p.mu.RLock()
	current := make([]*managed, 0, len(p.boundaries))
	for _, m := range p.boundaries {
		current = append(current, m)
	}
	p.mu.RUnlock()
	now := p.config.Now()
	for _, m := range current {
		p.report(m.boundary.Tick(ctx, now))
	}
	p.mu.Lock()
	for key, m := range p.boundaries {
		if m.boundary.Retired() {
			delete(p.boundaries, key)
		}
	}
	p.mu.Unlock()
}

// Stop ends intake ticking and shuts the delivery worker down within deadline.
// Accepted data stays durable either way.
func (p *Pipeline) Stop(deadline time.Duration) error {
	p.lifecycle.Lock()
	defer p.lifecycle.Unlock()
	if !p.started {
		return nil
	}
	p.started = false
	p.cancel()
	<-p.loopDone
	return p.worker.Stop(deadline)
}

// AcceptSample implements the runtime's typed sample sink. Every boundary sees
// the sample and ignores what is not its own. Refusals from several groups are
// joined so a real fault is never hidden behind a benign one.
func (p *Pipeline) AcceptSample(ctx context.Context, envelope measurement.SampleEnvelope) error {
	p.mu.RLock()
	current := make([]*managed, 0, len(p.boundaries))
	for _, m := range p.boundaries {
		current = append(current, m)
	}
	p.mu.RUnlock()
	var errs []error
	for _, m := range current {
		if err := m.boundary.AcceptSample(ctx, envelope); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// WantsSample reports whether any running group consumes this source, so the
// runtime builds typed envelopes only for tags a group actually owns.
func (p *Pipeline) WantsSample(deviceID, pointID, tagID string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, m := range p.boundaries {
		if m.boundary.Wants(deviceID, pointID, tagID) {
			return true
		}
	}
	return false
}

// Owns reports whether an active or still-closing group writes this tag to
// this connector, in which case the legacy writer must not.
func (p *Pipeline) Owns(connectorID, tagID string) bool {
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
		if m.boundary.Retiring() {
			state = StateRetiring
		}
		statuses = append(statuses, GroupStatus{GroupID: m.groupID, Revision: m.revision, State: state})
	}
	for _, status := range p.blocked {
		statuses = append(statuses, status)
	}
	return statuses
}

// destinationResolver opens the frozen destination of an outbox row and maps a
// blocked destination to the delivery layer's blocked error.
type destinationResolver struct{ destinations Destinations }

func (r destinationResolver) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	opened, err := r.destinations.OpenDestination(ctx, item.Destination.ConnectorID, item.Destination.ConnectorRevision)
	if err != nil {
		if errors.Is(err, dbtarget.ErrDestinationBlocked) {
			return groupdelivery.Target{}, fmt.Errorf("%w: %w", groupdelivery.ErrTargetBlocked, err)
		}
		return groupdelivery.Target{}, err
	}
	return groupdelivery.Target{DB: opened.DB, Kind: opened.Kind, Close: opened.Close}, nil
}

var (
	_ runtime.SampleSink     = (*Pipeline)(nil)
	_ runtime.SampleInterest = (*Pipeline)(nil)
)
