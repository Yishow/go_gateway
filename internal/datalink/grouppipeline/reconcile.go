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
	list, err := p.deps.Groups.List(ctx)
	if err != nil {
		return err
	}
	now := p.config.Now()
	known := make(map[string]struct{}, len(list.Groups))
	for _, group := range list.Groups {
		known[group.ID] = struct{}{}
		if group.Status == workspace.WriteGroupStatusDisabled || group.Status == workspace.WriteGroupStatusDeleted || group.AppliedRevision == "" {
			p.retireGroup(group.ID, now)
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
		m.boundary.SetUntil(nextBoundary(now, m.interval))
	}
	return p.restoreOpenJournals(ctx)
}

// reconcileGroup never fails the whole reconcile: a group that cannot be built
// is recorded as blocked with a safe reason and retried on the next pass.
func (p *Pipeline) reconcileGroup(ctx context.Context, group *workspace.WriteGroup, now time.Time) {
	current, err := p.deps.Groups.ResolveAppliedAt(ctx, group.ID, now)
	if err != nil {
		p.setBlocked(group, reasonSnapshotUnavailable)
		return
	}
	snapshots := []*workspace.WriteGroupAppliedSnapshot{current}
	if upcoming, err := p.deps.Groups.ResolveAppliedAt(ctx, group.ID, now.Add(lookaheadFactor*p.config.ReconcileInterval)); err == nil &&
		upcoming.AppliedRevision != current.AppliedRevision {
		snapshots = append(snapshots, upcoming)
	}
	// A reconcile can miss the lookahead window. End every older boundary at
	// the current snapshot's start even when no future snapshot was observed.
	p.mu.RLock()
	for _, previous := range p.boundaries {
		if previous.groupID == group.ID && previous.effective.Before(current.EffectiveAt) {
			previous.boundary.SetUntil(alignUp(current.EffectiveAt, previous.interval))
		}
	}
	p.mu.RUnlock()
	built := false
	for _, snap := range snapshots {
		key := boundaryKey(group.ID, snap.AppliedRevision)
		p.mu.RLock()
		_, exists := p.boundaries[key]
		p.mu.RUnlock()
		if exists {
			built = true
			continue
		}
		m, reason := p.build(ctx, snap)
		if reason != "" {
			p.setBlocked(group, reason)
			continue
		}
		p.mu.Lock()
		p.boundaries[key] = m
		p.mu.Unlock()
		built = true
	}
	if built {
		p.clearBlocked(group.ID)
	}
	if len(snapshots) == 2 {
		// The superseded revision owns samples only until the new one is effective.
		p.mu.RLock()
		old := p.boundaries[boundaryKey(group.ID, current.AppliedRevision)]
		p.mu.RUnlock()
		if old != nil {
			old.boundary.SetUntil(alignUp(snapshots[1].EffectiveAt, old.interval))
		}
	}
}

// build creates the boundary for one applied snapshot from real destination
// metadata, or returns the safe reason it cannot be built yet.
func (p *Pipeline) build(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot) (built *managed, blockedReason string) {
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
	if !ok {
		return nil, reasonDestinationUnsupported
	}
	frozen, err := p.loadRuntimeSnapshot(ctx, groupKey(snap))
	if err != nil {
		return nil, reasonSnapshotUnavailable
	}
	if frozen != nil {
		if !sameAppliedSnapshot(snap, frozen) || frozen.RuntimeLayout.Dialect != dialect {
			return nil, reasonSnapshotUnavailable
		}
		return p.buildFrozen(ctx, frozen, time.Time{}, time.Time{})
	}
	if canonicalManagedGroup(group) {
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
	built, blockedReason = p.buildFrozen(ctx, &verified, time.Time{}, time.Time{})
	if blockedReason != "" {
		return nil, blockedReason
	}
	if err := p.saveRuntimeSnapshot(ctx, &verified); err != nil {
		return nil, reasonSnapshotUnavailable
	}
	return built, ""
}

func (p *Pipeline) buildFrozen(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot, first, until time.Time) (built *managed, blockedReason string) {
	group, layout := snap.Group, snap.RuntimeLayout
	interval := time.Duration(group.RowPolicy.IntervalSeconds) * time.Second
	if interval < time.Second {
		return nil, reasonIntervalInvalid
	}
	if first.IsZero() {
		first = alignUp(snap.EffectiveAt, interval)
	}
	members := make([]member, 0, len(group.Members))
	for _, m := range group.Members {
		members = append(members, member{tagID: m.TagID, connectorID: group.Destination.ConnectorID})
	}
	key := groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: snap.AppliedRevision}
	boundary, err := runtime.NewGroupBoundaryContext(ctx, runtime.GroupBoundaryConfig{
		Group: group, TagTypes: layout.TagTypes, Dialect: layout.Dialect, Columns: layout.Columns,
		FirstBucket: first, Until: until, MaxFutureSkew: p.config.MaxFutureSkew,
		Ledger: p.deps.Store.Ledger(key), Clock: p.config.Now, ReceiptTableReady: layout.ReceiptTableReady,
	})
	if err != nil {
		var blocked *runtime.GroupBoundaryError
		if errors.As(err, &blocked) {
			return nil, blocked.Code
		}
		return nil, reasonDestinationUnavailable
	}
	return &managed{
		groupID: group.ID, revision: snap.AppliedRevision, boundary: boundary, members: members,
		connector: group.Destination.ConnectorID, interval: interval, effective: snap.EffectiveAt,
	}, ""
}

func dialectFor(kind schema.DatabaseConnectorKind) (dbtarget.SQLDialect, bool) {
	switch kind {
	case schema.DatabaseConnectorKindSQLite:
		return dbtarget.SQLDialectSQLite, true
	case schema.DatabaseConnectorKindPostgres:
		return dbtarget.SQLDialectPostgres, true
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
		m.boundary.SetUntil(nextBoundary(now, m.interval))
	}
}

func (p *Pipeline) setBlocked(group *workspace.WriteGroup, reason string) {
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
