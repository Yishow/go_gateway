package grouppipeline

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

// Safe blocked reasons.
const (
	reasonSnapshotUnavailable    = "snapshot-unavailable"
	reasonDestinationUnavailable = "destination-unavailable"
	reasonDestinationMissing     = "destination-missing"
	reasonDestinationRevision    = "destination-revision-changed"
	reasonDestinationUnsupported = "destination-unsupported"
	reasonTableMissing           = "destination-table-missing"
	reasonTableUnavailable       = "destination-table-unavailable"
	reasonTagUnavailable         = "tag-unavailable"
	reasonIntervalInvalid        = "interval-invalid"
)

// lookaheadFactor looks this many reconcile intervals ahead so a revision that
// becomes effective at the next bucket boundary is ready before its first
// sample, and the superseded boundary stops exactly there.
const lookaheadFactor = 2

func boundaryKey(groupID, revision string) string { return groupID + "@" + revision }

// Reconcile aligns the running boundaries with the canonical groups: it builds
// boundaries for applied revisions, ends superseded revisions at the new
// effective boundary, and ends disabled or deleted groups at the next bucket
// boundary. Accepted buckets always keep closing and delivering.
func (p *Pipeline) Reconcile(ctx context.Context) error {
	ctx, finish := p.ownerOperation(ctx)
	defer finish()
	defer p.observeQuotaRecovery(ctx)
	p.reconcileMu.Lock()
	defer p.reconcileMu.Unlock()
	list, err := p.deps.Groups.List(ctx)
	if err != nil {
		return err
	}
	now := p.config.Now()
	known := make(map[string]struct{}, len(list.Groups))
	for _, group := range list.Groups {
		known[group.ID] = struct{}{}
		if group.Status == workspace.WriteGroupStatusDisabled || group.Status == workspace.WriteGroupStatusDeleted || group.AppliedRevision == "" {
			p.retireObservedGroup(ctx, group, now)
			continue
		}
		p.reconcileGroup(ctx, group, now)
	}
	// A group that disappeared entirely is retired the same way.
	p.mu.RLock()
	var gone []*managed
	for _, m := range p.boundaries {
		if _, ok := known[m.groupID]; !ok {
			gone = append(gone, m)
		}
	}
	p.mu.RUnlock()
	for _, m := range gone {
		p.observedIntakeCutoff(ctx, m, func() error { return m.boundary.SetUntil(nextBoundary(now, m.interval)) }, false)
	}
	handoffErr := p.reconcileHandoffs(ctx, list)
	return errors.Join(handoffErr, p.restoreOpenJournals(ctx))
}

// reconcileGroup never fails the whole reconcile: a group that cannot be built
// is recorded as blocked with a safe reason and retried on the next pass.
func (p *Pipeline) reconcileGroup(ctx context.Context, group *workspace.WriteGroup, now time.Time) {
	current, err := p.deps.Groups.ResolveAppliedAt(ctx, group.ID, now)
	var snapshots []*workspace.WriteGroupAppliedSnapshot
	if err == nil {
		snapshots = append(snapshots, current)
	}
	// A revision that becomes effective soon is built before its first sample.
	// That includes the first revision of a group a cutover starts at an exact
	// instant: until then it has no current snapshot to build from.
	if upcoming, upErr := p.deps.Groups.ResolveAppliedAt(ctx, group.ID, now.Add(lookaheadFactor*p.config.ReconcileInterval)); upErr == nil &&
		(current == nil || upcoming.AppliedRevision != current.AppliedRevision) {
		snapshots = append(snapshots, upcoming)
	}
	if len(snapshots) == 0 {
		p.setBlocked(ctx, group, reasonSnapshotUnavailable)
		return
	}
	// A reconcile can miss the lookahead window. End every older boundary at
	// the current snapshot's start even when no future snapshot was observed.
	if current != nil {
		p.mu.RLock()
		var previousRevisions []*managed
		for _, previous := range p.boundaries {
			if previous.groupID == group.ID && previous.effective.Before(current.EffectiveAt) {
				previousRevisions = append(previousRevisions, previous)
			}
		}
		p.mu.RUnlock()
		for _, previous := range previousRevisions {
			p.observedIntakeCutoff(ctx, previous, func() error { return previous.boundary.SetEndsAt(current.EffectiveAt) }, false)
		}
	}
	built := false
	for _, snap := range snapshots {
		key := boundaryKey(group.ID, snap.AppliedRevision)
		p.mu.RLock()
		existing := p.boundaries[key]
		p.mu.RUnlock()
		if existing != nil {
			if !snap.EndsAt.IsZero() {
				// A cutover took the sources away from this revision at an exact instant.
				p.observedIntakeCutoff(ctx, existing, func() error { return existing.boundary.SetEndsAt(snap.EndsAt) }, false)
			}
			built = true
			continue
		}
		m, reason := p.build(ctx, snap)
		if reason != "" {
			p.setBlocked(ctx, snap.Group, reason)
			continue
		}
		p.installManaged(m)
		built = true
	}
	if built {
		p.clearBlocked(group.ID)
	}
	if current != nil && len(snapshots) == 2 {
		// The superseded revision owns samples only until the new one is effective.
		p.mu.RLock()
		old := p.boundaries[boundaryKey(group.ID, current.AppliedRevision)]
		p.mu.RUnlock()
		if old != nil {
			p.observedIntakeCutoff(ctx, old, func() error { return old.boundary.SetEndsAt(snapshots[1].EffectiveAt) }, false)
		}
	}
	p.recordAppliedIntake(ctx, current, now)
}

// build creates the boundary for one applied snapshot from real destination
// metadata, or returns the safe reason it cannot be built yet.
func (p *Pipeline) build(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot) (*managed, string) {
	return p.buildAt(ctx, snap, time.Time{}, time.Time{})
}

func (p *Pipeline) buildAt(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot, first, until time.Time) (built *managed, blockedReason string) {
	group := snap.Group
	connector, err := p.deps.Destinations.GetByID(ctx, group.Destination.ConnectorID)
	if errors.Is(err, dbtarget.ErrConnectorNotFound) {
		return nil, reasonDestinationMissing
	}
	if err != nil {
		return nil, reasonDestinationUnavailable
	}
	if connector.IdentityRevision != group.Destination.ConnectorRevision {
		return nil, reasonDestinationRevision
	}
	dialect, ok := dialectFor(connector.Kind)
	managedLongOnly := connector.Kind == schema.DatabaseConnectorKindMySQL || connector.Kind == schema.DatabaseConnectorKindSQLServer
	if !ok || managedLongOnly && !longGroup(group) {
		return nil, reasonDestinationUnsupported
	}
	frozen, err := p.loadRuntimeSnapshot(ctx, groupKey(snap))
	if err != nil {
		return nil, frozenFailureReason(err)
	}
	if frozen != nil {
		if !sameAppliedSnapshot(snap, frozen) || frozen.RuntimeLayout == nil || frozen.RuntimeLayout.Dialect != dialect ||
			managedLongOnly && !validLongRuntimeProof(group, frozen.RuntimeLayout) {
			return nil, reasonSnapshotUnavailable
		}
		frozen.EndsAt = snap.EndsAt
		return p.buildFrozen(ctx, frozen, first, until)
	}
	if needsCounterSources(group) {
		return nil, "counter-replay-proof-missing"
	}
	if canonicalManagedGroup(group) || longGroup(group) {
		// Only the verified snapshot frozen at apply can start these groups; the
		// inspected table is never used to rebuild one.
		return nil, reasonSnapshotUnavailable
	}
	inspection, err := p.deps.Inspector.InspectTable(ctx, group.Destination.ConnectorID, group.Destination.TableSchema, group.Destination.TableName)
	if err != nil || inspection == nil {
		return nil, reasonTableUnavailable
	}
	if inspection.Status == dbtarget.TableInspectionMissing {
		return nil, reasonTableMissing
	}
	if inspection.Status != dbtarget.TableInspectionExists {
		return nil, reasonTableUnavailable
	}
	receiptReady := false
	if strings.EqualFold(strings.TrimSpace(group.WritePolicy.DedupeCapability), groupdelivery.DedupeReceipt) {
		receipts, err := p.deps.Inspector.InspectTable(ctx, group.Destination.ConnectorID, group.Destination.TableSchema, receiptTableName)
		receiptReady = err == nil && receipts != nil && receipts.Status == dbtarget.TableInspectionExists
	}
	tagTypes := make(map[string]schema.DataType, len(group.Members))
	if source, ok := p.deps.Groups.(interface {
		AppliedRuntimeTagTypes(context.Context, *workspace.WriteGroupAppliedSnapshot) (map[string]schema.DataType, error)
	}); ok {
		tagTypes, err = source.AppliedRuntimeTagTypes(ctx, snap)
		if err != nil {
			return nil, reasonTagUnavailable
		}
	} else {
		for _, m := range group.Members {
			tag, err := p.deps.Tags.GetByID(ctx, m.TagID)
			if err != nil || tag == nil {
				return nil, reasonTagUnavailable
			}
			tagTypes[m.TagID] = tag.DataType
		}
	}
	verified := *snap
	verified.RuntimeLayout = &workspace.WriteGroupRuntimeLayout{
		Dialect: dialect, Columns: inspection.Columns, TagTypes: tagTypes, ReceiptTableReady: receiptReady,
	}
	built, blockedReason = p.buildFrozen(ctx, &verified, first, until)
	if blockedReason != "" {
		return nil, blockedReason
	}
	if err := p.saveRuntimeSnapshot(ctx, &verified); err != nil {
		return nil, reasonSnapshotUnavailable
	}
	return built, ""
}

func (p *Pipeline) buildFrozen(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot, first, until time.Time) (built *managed, reason string) {
	built, reason, err := p.buildFrozenWithCause(ctx, snap, first, until)
	if err != nil {
		return nil, reason
	}
	return built, reason
}

// buildFrozenWithCause preserves activation causes for durable recovery without
// exposing them through the safe blocked reason or changing layout refusals.
func (p *Pipeline) buildFrozenWithCause(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot, first, until time.Time) (built *managed, blockedReason string, cause error) {
	group, layout := snap.Group, snap.RuntimeLayout
	interval := time.Duration(group.RowPolicy.IntervalSeconds) * time.Second
	if interval < time.Second {
		return nil, reasonIntervalInvalid, nil
	}
	effective := snap.EffectiveAt.UTC()
	if first.IsZero() {
		// An ordinary revision takes effect on a window boundary. A cutover revision
		// takes effect at an exact instant: its first window is the one containing it.
		first = snapshot.BucketStart(effective, interval)
	}
	var startsAt, endsAt time.Time
	if !effective.Equal(first) && snapshot.BucketStart(effective, interval).Equal(first) {
		startsAt = effective
	}
	if !snap.EndsAt.IsZero() {
		end := alignUp(snap.EndsAt, interval)
		if until.IsZero() || until.Equal(end) {
			until, endsAt = end, snap.EndsAt.UTC()
		}
	}
	members := make([]member, 0, len(group.Members))
	for _, m := range group.Members {
		members = append(members, member{tagID: m.TagID, connectorID: group.Destination.ConnectorID})
	}
	key := groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: snap.AppliedRevision}
	boundary, err := runtime.NewGroupBoundaryContext(ctx, runtime.GroupBoundaryConfig{
		Group: group, TagTypes: layout.TagTypes, CounterBindings: snap.CounterBindings, Dialect: layout.Dialect, Columns: layout.Columns,
		FirstBucket: first, Until: until, StartsAt: startsAt, EndsAt: endsAt, MaxFutureSkew: p.config.MaxFutureSkew,
		Ledger: p.deps.Store.Ledger(key), Clock: p.config.Now, ReceiptTableReady: layout.ReceiptTableReady,
	})
	if err != nil {
		var blocked *runtime.GroupBoundaryError
		if errors.As(err, &blocked) {
			if blocked.Cause != nil {
				return nil, blocked.Code, err
			}
			return nil, blocked.Code, nil
		}
		return nil, reasonDestinationUnavailable, err
	}
	return &managed{
		workspaceID: group.WorkspaceID, groupID: group.ID, revision: snap.AppliedRevision, boundary: boundary, members: members,
		connector: group.Destination.ConnectorID, interval: interval, effective: snap.EffectiveAt,
	}, "", nil
}

func dialectFor(kind schema.DatabaseConnectorKind) (dbtarget.SQLDialect, bool) {
	switch kind {
	case schema.DatabaseConnectorKindSQLite:
		return dbtarget.SQLDialectSQLite, true
	case schema.DatabaseConnectorKindPostgres:
		return dbtarget.SQLDialectPostgres, true
	case schema.DatabaseConnectorKindMySQL:
		return dbtarget.SQLDialectMySQL, true
	case schema.DatabaseConnectorKindSQLServer:
		return dbtarget.SQLDialectSQLServer, true
	}
	return "", false
}

func (p *Pipeline) retireGroup(groupID string, now time.Time) {
	p.clearBlocked(groupID)
	p.mu.RLock()
	var current []*managed
	for _, m := range p.boundaries {
		if m.groupID == groupID {
			current = append(current, m)
		}
	}
	p.mu.RUnlock()
	for _, m := range current {
		p.report(m.boundary.SetUntil(nextBoundary(now, m.interval)))
	}
}

func (p *Pipeline) setBlocked(ctx context.Context, group *workspace.WriteGroup, reason string) {
	if p.observedCounterQuarantine(ctx, groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.AppliedRevision}, reason) {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.blocked[group.ID] = GroupStatus{GroupID: group.ID, Revision: group.AppliedRevision, State: StateBlocked, Reason: reason}
}

func (p *Pipeline) clearBlocked(groupID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.blocked, groupID)
}

func nextBoundary(now time.Time, interval time.Duration) time.Time {
	return snapshot.BucketStart(now, interval).Add(interval)
}

// alignUp rounds an effective time up to the next interval boundary (it is
// already aligned for applied snapshots).
func alignUp(t time.Time, interval time.Duration) time.Time {
	start := snapshot.BucketStart(t, interval)
	if start.Equal(t.UTC()) {
		return start
	}
	return start.Add(interval)
}
