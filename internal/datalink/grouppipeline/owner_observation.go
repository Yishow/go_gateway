package grouppipeline

import (
	"bytes"
	"context"
	"errors"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"math"
	"strings"
)

// OwnerObservationKind separates accepted input quarantine from an ownership refusal.
type OwnerObservationKind string

const (
	OwnerCounterQuarantine OwnerObservationKind = "counter-quarantine"
	OwnerRecordingRefusal  OwnerObservationKind = "recording-refusal"
)

// OwnerObservation contains only verified diagnostic identity and input counts.
// ObservationID is an ephemeral observation reference, never a durable operation ID.
type OwnerObservation struct {
	Available                      bool
	Key                            groupdelivery.GroupKey
	WorkspaceRevision              string
	ConnectorID, ConnectorRevision string
	ObservationID                  string
	Kind                           OwnerObservationKind
	Reason                         string
	PendingSamples                 uint64
}

// OwnerObservationConfig is immutable after construction. Callbacks must be
// nonblocking; panics and proof failures cannot alter the native outcome.
type OwnerObservationConfig struct {
	// RecordLimit is injected from the existing diagnostics admission cap.
	RecordLimit   int
	OnDrop        func(uint64)
	ScopeRevision func(context.Context, groupdelivery.GroupKey) (string, bool)
	Observe       func(OwnerObservation)
	// ObserveCaptured is preferred to Observe (exactly one is called). A true
	// result acknowledges local coalescer acceptance, not broker/sink admission,
	// persistence, delivery, or native success.
	ObserveCaptured func(OwnerObservation) bool
}

// ownerDescriptor is an operation-local fact, not a replay/ownership decision.
type ownerDescriptor struct {
	intake         *intakeDescriptor
	key            groupdelivery.GroupKey
	kind           OwnerObservationKind
	reason         string
	owner          groupdelivery.OwnershipView
	ownerKnown     bool
	pendingGeneric bool
}

func ownerRefusalReason(err error) string {
	for _, reason := range []*groupdelivery.OwnershipError{groupdelivery.ErrOwnershipUnavailable, groupdelivery.ErrOwnershipConflict, groupdelivery.ErrHandoffProofConflict} {
		if errors.Is(err, reason) {
			return reason.Error()
		}
	}
	return ""
}

func safeOwnerID(v string) bool {
	if v == "" || len(v) > 128 {
		return false
	}
	for i, c := range []byte(v) {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '-' || i > 0 && c == '.' {
			continue
		}
		return false
	}
	return true
}

// verifyOwnerObservation reads only after the native operation released its locks.
// WorkspaceRevision brackets this observation, not the historical configuration.
func (p *Pipeline) verifyOwnerObservation(ctx context.Context, d ownerDescriptor) (out OwnerObservation) {
	defer func() {
		if recover() != nil {
			out = OwnerObservation{}
		}
	}()
	cfg := p.deps.OwnerObservation
	if !cfg.enabled() || cfg.ScopeRevision == nil || p.deps.Store == nil || !d.key.Valid() {
		return out
	}
	if d.kind == OwnerRecordingRefusal && !d.ownerKnown {
		return out
	}
	if d.kind != OwnerCounterQuarantine && d.kind != OwnerRecordingRefusal {
		return out
	}
	if d.kind == OwnerCounterQuarantine && !isCounterProofReason(d.reason) {
		return out
	}
	if d.kind == OwnerRecordingRefusal && ownerRefusalReason(&groupdelivery.OwnershipError{Reason: d.reason}) == "" {
		return out
	}
	for _, id := range []string{d.key.WorkspaceID, d.key.GroupID, d.key.GroupRevision} {
		if !safeOwnerID(id) {
			return out
		}
	}
	revision, ok := cfg.ScopeRevision(ctx, d.key)
	if !ok || !safeOwnerID(revision) || !p.ownerQuarantineMatches(d) {
		return out
	}
	raw, err := p.deps.Store.RuntimeVersion(ctx, d.key)
	if err != nil || len(raw) == 0 {
		return out
	}
	// Storage-only verification preserves original connector identity even when
	// counter replay proof is damaged. It never authorizes rebuilding or replay.
	storage, err := counterproof.SnapshotStorageProjection(raw)
	if err != nil {
		return out
	}
	frozen, err := decodeFrozenStorageSnapshot(storage, d.key)
	if err != nil || frozen == nil || !safeOwnerID(frozen.Group.Destination.ConnectorID) || !safeOwnerID(frozen.Group.Destination.ConnectorRevision) {
		return out
	}
	beforeOwner, err := p.deps.Store.OwnershipView(ctx, d.key)
	if err != nil || d.kind == OwnerRecordingRefusal && beforeOwner != d.owner {
		return out
	}
	count, err := p.deps.Store.PendingSampleCount(ctx, d.key)
	if err != nil || count < 0 {
		return out
	}
	afterRevision, ok := cfg.ScopeRevision(ctx, d.key)
	if !ok || revision != afterRevision || !p.ownerQuarantineMatches(d) {
		return out
	}
	afterRaw, err := p.deps.Store.RuntimeVersion(ctx, d.key)
	if err != nil || !bytes.Equal(raw, afterRaw) {
		return out
	}
	afterOwner, err := p.deps.Store.OwnershipView(ctx, d.key)
	if err != nil || afterOwner != beforeOwner {
		return out
	}
	if !p.ownerQuarantineMatches(d) {
		return out
	}
	id, err := common.NewUUID()
	if err != nil {
		return out
	}
	return OwnerObservation{Available: true, Key: d.key, WorkspaceRevision: revision, ConnectorID: frozen.Group.Destination.ConnectorID, ConnectorRevision: frozen.Group.Destination.ConnectorRevision, ObservationID: id, Kind: d.kind, Reason: d.reason, PendingSamples: uint64(count)}
}

func (p *Pipeline) ownerQuarantineMatches(d ownerDescriptor) bool {
	if d.kind != OwnerCounterQuarantine {
		return true
	}
	p.mu.RLock()
	q, ok := p.quarantines[d.key]
	p.mu.RUnlock()
	return ok && q.Reason == d.reason
}

type ownerBatchContextKey struct{}
type ownerObservationBatch struct {
	descriptors []ownerDescriptor
	dropped     uint64
	limit       int
	loss        func(uint64)
}

// ownerOperation defers provider reads and callbacks until the outermost native
// operation has released its locks and intake lifetime. Nested calls reuse it.
func (p *Pipeline) ownerOperation(ctx context.Context) (operationContext context.Context, flush func()) {
	limit, loss, enabled := p.observationBatchPolicy()
	if !enabled {
		return ctx, func() {}
	}
	if _, ok := ctx.Value(ownerBatchContextKey{}).(*ownerObservationBatch); ok {
		return ctx, func() {}
	}
	b := &ownerObservationBatch{limit: limit, loss: loss}
	return context.WithValue(ctx, ownerBatchContextKey{}, b), func() { p.flushOwnerObservations(ctx, b) }
}

func (p *Pipeline) recordOwnerObservation(ctx context.Context, d ownerDescriptor) {
	if !p.deps.OwnerObservation.enabled() {
		return
	}
	b, ok := ctx.Value(ownerBatchContextKey{}).(*ownerObservationBatch)
	if !ok { // Only standalone, already-unlocked wrappers may use this path.
		ctx, finish := p.ownerOperation(ctx)
		defer finish()
		p.recordOwnerObservation(ctx, d)
		return
	}
	if b.limit <= 0 || len(b.descriptors) >= b.limit {
		b.drop()
		return
	}
	// A descriptor cannot retain arbitrary-sized/hostile storage identity bytes.
	if !safeOwnerID(d.key.WorkspaceID) || !safeOwnerID(d.key.GroupID) || !safeOwnerID(d.key.GroupRevision) ||
		d.kind == OwnerCounterQuarantine && !isCounterProofReason(d.reason) ||
		d.kind == OwnerRecordingRefusal && ownerRefusalReason(&groupdelivery.OwnershipError{Reason: d.reason}) == "" {
		d = ownerDescriptor{}
	} else {
		d.key.WorkspaceID = strings.Clone(d.key.WorkspaceID)
		d.key.GroupID = strings.Clone(d.key.GroupID)
		d.key.GroupRevision = strings.Clone(d.key.GroupRevision)
		d.reason = strings.Clone(d.reason)
	}
	b.descriptors = append(b.descriptors, d)
}

func (b *ownerObservationBatch) drop() {
	if b.dropped < math.MaxUint64 {
		b.dropped++
	}
}
func (p *Pipeline) flushOwnerObservations(ctx context.Context, b *ownerObservationBatch) {
	for _, d := range b.descriptors {
		if d.intake != nil {
			p.flushIntakeObservation(ctx, b, *d.intake)
			continue
		}
		accepted := p.flushOwnerObservation(ctx, b, d)
		if d.pendingGeneric && !accepted {
			reportGenericDiagnostic(p.config)
		}
	}
	if b.dropped != 0 && b.loss != nil {
		func() {
			defer func() { recover() }() //nolint:errcheck // Loss reporting cannot affect the native operation.
			b.loss(b.dropped)
		}()
	}
}

func (p *Pipeline) observedCounterQuarantine(ctx context.Context, key groupdelivery.GroupKey, reason string) bool {
	if !p.setCounterQuarantineNative(key, reason) {
		return false
	}
	p.recordOwnerObservation(ctx, ownerDescriptor{key: key, kind: OwnerCounterQuarantine, reason: reason})
	return true
}

func (p *Pipeline) observedProofGate(ctx context.Context, m *managed, operation func() error) error {
	var original groupdelivery.OwnershipView
	invoked := false
	err := p.withCounterProofGate(m, func() error {
		// This is the immutable view held by the boundary actually executing the
		// operation, not a newly read or substituted routing hint.
		if p.deps.OwnerObservation.enabled() {
			original = m.boundary.ReplayState().OwnershipView
		}
		invoked = true
		return operation()
	}) // Native deferred unlock has finished.
	if invoked {
		p.recordQuotaRefusal(ctx, m, err)
	}
	key := groupdelivery.GroupKey{WorkspaceID: m.workspaceID, GroupID: m.groupID, GroupRevision: m.revision}
	var proof *counterproof.CounterBindingProofError
	if errors.As(err, &proof) && proof.Reason != nil && isCounterProofReason(proof.Reason.Error()) {
		p.recordOwnerObservation(ctx, ownerDescriptor{key: key, kind: OwnerCounterQuarantine, reason: proof.Reason.Error()})
	} else {
		p.recordOwnerRefusalsWithView(ctx, key, err, original, invoked)
	}
	return err
}

// Joined native refusals can contain distinct original classified reasons.
func (p *Pipeline) recordOwnerRefusals(ctx context.Context, key groupdelivery.GroupKey, err error) {
	p.recordOwnerRefusalsWithView(ctx, key, err, groupdelivery.OwnershipView{}, false)
}

func (p *Pipeline) recordOwnerRefusalsWithView(ctx context.Context, key groupdelivery.GroupKey, err error, original groupdelivery.OwnershipView, known bool) {
	for _, reason := range []*groupdelivery.OwnershipError{groupdelivery.ErrOwnershipUnavailable, groupdelivery.ErrOwnershipConflict, groupdelivery.ErrHandoffProofConflict} {
		if errors.Is(err, reason) {
			p.recordOwnerObservation(ctx, ownerDescriptor{key: key, kind: OwnerRecordingRefusal, reason: reason.Error(), owner: original, ownerKnown: known})
		}
	}
}
