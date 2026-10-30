package grouppipeline

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
	"time"
)

// IntakeQuotaObservationKind identifies a native intake fact or capacity decision.
type IntakeQuotaObservationKind string

const (
	IntakeGroupApplied   IntakeQuotaObservationKind = "group-applied"
	IntakeGroupPaused    IntakeQuotaObservationKind = "group-paused"
	IntakeGroupRetiring  IntakeQuotaObservationKind = "group-retiring"
	IntakeQuotaRefused   IntakeQuotaObservationKind = "quota-refused"
	IntakeQuotaRecovered IntakeQuotaObservationKind = "quota-recovered"
)

// IntakeQuotaObservation contains no data values or destination configuration.
// ObservationID identifies an observation, not a durable operation or receipt.
type IntakeQuotaObservation struct {
	Available                      bool
	Key                            groupdelivery.GroupKey
	WorkspaceRevision              string
	ConnectorID, ConnectorRevision string
	ObservationID                  string
	Kind                           IntakeQuotaObservationKind
	IntakeState                    string
	Quota                          groupdelivery.QuotaStatus
}

// QuotaRecoveryHint is a bounded diagnostic hint, never intake authorization.
// Observe closes over the original opaque diagnostic episode token.
type QuotaRecoveryHint struct {
	Key                                                            groupdelivery.GroupKey
	WorkspaceRevision, ConnectorID, ConnectorRevision, ReferenceID string
	Scope                                                          groupdelivery.QuotaScope
	Observe                                                        func(IntakeQuotaObservation)
}

// IntakeQuotaObservationConfig is immutable after construction. Callbacks must
// be nonblocking; diagnostic failures cannot change native operation outcomes.
type IntakeQuotaObservationConfig struct {
	RecordLimit                 int
	OnDrop                      func(uint64)
	IntakeEnabled, QuotaEnabled bool
	ScopeRevision               func(context.Context, groupdelivery.GroupKey) (string, bool)
	Observe                     func(IntakeQuotaObservation)
	TakeQuotaRecoveryHint       func(context.Context) (QuotaRecoveryHint, bool)
}

// intakeDescriptor retains bounded operation input facts, never a current-draft
// substitute. The managed pointer is only a cache-identity check, not authority.
type intakeDescriptor struct {
	key           groupdelivery.GroupKey
	kind          IntakeQuotaObservationKind
	managed       *managed
	until         time.Time
	appliedDigest [32]byte
	quota         groupdelivery.QuotaStatus
	hint          *QuotaRecoveryHint
}

func intakeGroupDigest(group *workspace.WriteGroup) [32]byte {
	raw, err := json.Marshal(group)
	if err != nil {
		return [32]byte{}
	}
	return sha256.Sum256(raw)
}

func (p *Pipeline) intakeEnabled() bool {
	c := p.deps.IntakeQuotaObservation
	return c.Observe != nil && c.IntakeEnabled
}
func (p *Pipeline) quotaEnabled() bool {
	c := p.deps.IntakeQuotaObservation
	return c.Observe != nil && c.QuotaEnabled
}

func (p *Pipeline) recordIntakeObservation(ctx context.Context, d intakeDescriptor) {
	if d.kind == IntakeQuotaRefused || d.kind == IntakeQuotaRecovered {
		if !p.quotaEnabled() {
			return
		}
	} else if !p.intakeEnabled() {
		return
	}
	b, ok := ctx.Value(ownerBatchContextKey{}).(*ownerObservationBatch)
	if !ok {
		return
	}
	if b.limit <= 0 || len(b.descriptors) >= b.limit {
		b.drop()
		return
	}
	if !safeOwnerID(d.key.WorkspaceID) || !safeOwnerID(d.key.GroupID) || !safeOwnerID(d.key.GroupRevision) {
		d = intakeDescriptor{}
	}
	b.descriptors = append(b.descriptors, ownerDescriptor{intake: &d})
}

// recordAppliedIntake uses the ORIGINAL canonical input already resolved by
// Reconcile. No extra List/Resolve (which could bootstrap) is performed for logs.
func (p *Pipeline) recordAppliedIntake(ctx context.Context, snap *workspace.WriteGroupAppliedSnapshot, now time.Time) {
	if !p.intakeEnabled() || snap == nil || snap.Group == nil || snap.EffectiveAt.After(now) || !snap.EndsAt.IsZero() && !now.Before(snap.EndsAt) {
		return
	}
	p.mu.RLock()
	m := p.boundaries[boundaryKey(snap.GroupID, snap.AppliedRevision)]
	p.mu.RUnlock()
	if m == nil || m.drainOnly || m.boundary.Retiring() {
		return
	}
	p.recordIntakeObservation(ctx, intakeDescriptor{key: groupKey(snap), kind: IntakeGroupApplied, managed: m, appliedDigest: intakeGroupDigest(snap.Group)})
}

// observedIntakeCutoff preserves the native transition exactly. A missing group
// in a list is not a disabled/deleted eligibility decision and cannot be paused.
func (p *Pipeline) observedIntakeCutoff(ctx context.Context, m *managed, operation func(), paused bool) {
	if !p.intakeEnabled() {
		operation()
		return
	}
	before := m.boundary.ReplayState().Until
	operation()
	after := m.boundary.ReplayState().Until
	if after.IsZero() || after.Equal(before) {
		return
	}
	d := intakeDescriptor{key: groupdelivery.GroupKey{WorkspaceID: m.workspaceID, GroupID: m.groupID, GroupRevision: m.revision}, kind: IntakeGroupRetiring, managed: m, until: after}
	p.recordIntakeObservation(ctx, d)
	if paused {
		d.kind = IntakeGroupPaused
		p.recordIntakeObservation(ctx, d)
	}
}

func (p *Pipeline) intakeDescriptorMatches(d intakeDescriptor) bool {
	if d.managed == nil {
		return d.kind == IntakeQuotaRecovered
	}
	p.mu.RLock()
	current := p.boundaries[boundaryKey(d.key.GroupID, d.key.GroupRevision)]
	p.mu.RUnlock()
	if current != d.managed {
		return false
	}
	m := d.managed
	m.proofGate.Lock()
	defer m.proofGate.Unlock()
	if m.superseded || m.proofQuarantined != nil {
		return false
	}
	state := m.boundary.ReplayState()
	if !state.Until.Equal(d.until) {
		return false
	}
	if d.kind == IntakeGroupApplied && (m.drainOnly || !state.Until.IsZero()) {
		return false
	}
	if (d.kind == IntakeGroupPaused || d.kind == IntakeGroupRetiring) && state.Until.IsZero() {
		return false
	}
	return m.workspaceID == d.key.WorkspaceID && m.groupID == d.key.GroupID && m.revision == d.key.GroupRevision
}

// verifyIntakeObservation brackets post-unlock evidence with pure observed
// workspace metadata. It does NOT claim the original input is today's complete
// configuration. Stable WRevision alone never substitutes for frozen/native proof.
func (p *Pipeline) verifyIntakeObservation(ctx context.Context, d intakeDescriptor) (out IntakeQuotaObservation, emit bool) {
	emit = true
	defer func() {
		if recover() != nil {
			out = IntakeQuotaObservation{}
			emit = true
		}
	}()
	cfg := p.deps.IntakeQuotaObservation
	if ctx.Err() != nil || cfg.ScopeRevision == nil || p.deps.Store == nil || !d.key.Valid() || !p.intakeDescriptorMatches(d) {
		return
	}
	revision, ok := cfg.ScopeRevision(ctx, d.key)
	if !ok || !safeOwnerID(revision) {
		return
	}
	if d.hint != nil && (revision != d.hint.WorkspaceRevision || !safeOwnerID(d.hint.ReferenceID)) {
		return
	}
	raw, err := p.deps.Store.RuntimeVersion(ctx, d.key)
	if err != nil || len(raw) == 0 {
		return
	}
	if _, err = counterproof.ResolveSnapshot(raw, counterproof.SnapshotScope{WorkspaceID: d.key.WorkspaceID, GroupID: d.key.GroupID, AppliedRevision: d.key.GroupRevision}); err != nil {
		return
	}
	frozen, err := decodeFrozenStorageSnapshot(raw, d.key)
	if err != nil {
		return
	}
	destination := frozen.Group.Destination
	if !safeOwnerID(destination.ConnectorID) || !safeOwnerID(destination.ConnectorRevision) {
		return
	}
	if d.managed != nil && (destination.ConnectorID != d.managed.connector || !frozen.EffectiveAt.Equal(d.managed.effective)) {
		return
	}
	if d.kind == IntakeGroupApplied && d.appliedDigest != intakeGroupDigest(frozen.Group) {
		return
	}
	if d.hint != nil && (destination.ConnectorID != d.hint.ConnectorID || destination.ConnectorRevision != d.hint.ConnectorRevision) {
		return
	}
	quota := d.quota
	if d.kind == IntakeQuotaRecovered {
		quota, err = p.deps.Store.QuotaStatusForScope(ctx, d.key, d.hint.Scope)
		if err != nil || !validObservedQuota(quota, d.key) {
			return
		}
	} else if d.kind == IntakeQuotaRefused && !validObservedQuota(quota, d.key) {
		return
	}
	after, ok := cfg.ScopeRevision(ctx, d.key)
	if !ok || revision != after {
		return
	}
	afterRaw, err := p.deps.Store.RuntimeVersion(ctx, d.key)
	if err != nil || !bytes.Equal(raw, afterRaw) || !p.intakeDescriptorMatches(d) {
		return
	}
	if d.kind == IntakeQuotaRecovered && quota.IntakeRefused {
		return IntakeQuotaObservation{}, false
	}
	var id string
	if d.hint != nil {
		id = d.hint.ReferenceID
	} else {
		id, err = common.NewUUID()
		if err != nil {
			return
		}
	}
	quota.LossRiskNotice = ""
	state := map[IntakeQuotaObservationKind]string{IntakeGroupApplied: StateActive, IntakeGroupPaused: "not_running", IntakeGroupRetiring: StateRetiring}[d.kind]
	return IntakeQuotaObservation{Available: true, Key: d.key, WorkspaceRevision: revision, ConnectorID: destination.ConnectorID, ConnectorRevision: destination.ConnectorRevision, ObservationID: id, Kind: d.kind, IntakeState: state, Quota: quota}, true
}

func (p *Pipeline) retireObservedGroup(ctx context.Context, group *workspace.WriteGroup, now time.Time) {
	if !p.intakeEnabled() {
		p.retireGroup(group.ID, now)
		return
	}
	p.clearBlocked(group.ID)
	p.mu.RLock()
	var current []*managed
	for _, m := range p.boundaries {
		if m.groupID == group.ID {
			current = append(current, m)
		}
	}
	p.mu.RUnlock()
	for _, m := range current {
		paused := m.workspaceID == group.WorkspaceID && (group.Status == workspace.WriteGroupStatusDisabled || group.Status == workspace.WriteGroupStatusDeleted)
		p.observedIntakeCutoff(ctx, m, func() { m.boundary.SetUntil(nextBoundary(now, m.interval)) }, paused)
	}
}
