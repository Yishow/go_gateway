package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

// GroupRow is a closed row ready for delivery, with its encoded SQL values.
type GroupRow struct {
	Outcome snapshot.Outcome
	Encoded dbtarget.EncodedRow
}

// GroupRowSink receives closed buckets. It is the seam the durable delivery
// change implements; nothing here persists or retries, so rows handed to a
// failing sink are the sink owner's responsibility.
type GroupRowSink interface {
	AcceptRow(context.Context, GroupRow) error
	ReportOutcome(context.Context, snapshot.Outcome) error
}

// GroupBoundaryConfig fixes one applied group definition. Group must be the
// applied immutable version, not a draft, and Columns must come from real
// table inspection.
type GroupBoundaryConfig struct {
	Group    *workspace.WriteGroup
	TagTypes map[string]schema.DataType
	Counters map[string]CounterSource // confirmed cumulative semantics by tag id
	// CounterBindings freezes durable counter meaning; legacy exact descriptors
	// are the only supported fallback, never current measurement definitions.
	CounterBindings   workspace.CounterBindingProofs
	Dialect           dbtarget.SQLDialect
	Columns           []dbtarget.ColumnInfo
	RecordKeyColumn   string
	BucketStartColumn string
	FirstBucket       time.Time
	// Until, when set, ends this boundary's responsibility: samples observed at
	// or after it belong to a newer applied revision and are ignored here.
	Until time.Time
	// StartsAt and EndsAt are the exact, immutable instants of a forward
	// cutover: this revision owns a source's observations from StartsAt on
	// (FirstBucket is the window containing it) and the one it replaces owns
	// them before it. Half-open: observed at StartsAt belongs to the new owner.
	StartsAt, EndsAt time.Time
	// ReceiptTableReady states that the destination holds a verified effect
	// receipt table; the receipt dedupe capability requires it.
	ReceiptTableReady bool
	MaxFutureSkew     time.Duration
	Sink              GroupRowSink
	// Ledger, when set, makes sample acceptance durable: a sample is ACKed
	// only after its journal row commits and open buckets are rebuilt from the
	// journal at startup.
	Ledger GroupLedger
	// Clock is the gateway clock for future-skew checks; nil uses UTC now.
	Clock func() time.Time
	// OnTickError receives sink failures from RunTicks.
	OnTickError func(error)
}

// Activation failure codes. They never carry values or connection details.
const (
	boundaryGroupMissing     = "group-missing"
	boundarySinkMissing      = "sink-missing"
	boundaryGroupNotApplied  = "group-not-applied"
	boundaryMemberIdentity   = "member-identity-incomplete"
	boundaryTagTypeUnsupport = "tag-type-unsupported"
	boundaryLayoutBlocked    = "layout-blocked"
	boundaryScopeInvalid     = "destination-scope-invalid"
	boundarySnapshotInvalid  = "snapshot-config-invalid"
	boundaryRestoreFailed    = "restore-failed"
	boundaryDedupeUnsupport  = "dedupe-unsupported"

	reasonWorkspaceMismatch = "workspace-mismatch"
	reasonTypeMismatch      = snapshot.ReasonTypeMismatch
	reasonEncodeBlocked     = "encode-blocked:"
	reasonJournalFailed     = "journal-failed"
)

// GroupBoundaryError is a safe activation failure: codes only, no values.
type GroupBoundaryError struct {
	Code   string
	Issues []dbtarget.LayoutIssue
	Cause  error
}

// Unwrap exposes the underlying cause, if any.
func (e *GroupBoundaryError) Unwrap() error { return e.Cause }

func (e *GroupBoundaryError) Error() string { return "group boundary blocked: " + e.Code }

// GroupSampleError is a safe per-sample refusal. Cause keeps the underlying
// error for callers; Error never includes it.
type GroupSampleError struct {
	Outcome snapshot.OfferOutcome
	Reason  string
	Cause   error
}

func (e *GroupSampleError) Error() string { return "group sample refused: " + e.Reason }

// Unwrap exposes the underlying cause, if any.
func (e *GroupSampleError) Unwrap() error { return e.Cause }

type boundaryMember struct {
	key                             string
	interval                        time.Duration // mixed member cadence; never a common window
	kind                            measurement.ExactType
	counter                         *aggregation.CounterBinding // set for a cumulative member
	counterProof                    counterproof.CounterBindingReference
	counterClaims                   counterClaimIdentity
	exactScale                      string // set only from the immutable exact recording source proof
	sourceRevision, mappingRevision string // required frozen identities for a scaled source
}

// GroupBoundary is the opt-in typed-sample consumer for one applied group. It
// rebuilds open windows from its ledger when configured. Without a ledger,
// accepted samples and closed outcomes remain explicitly non-durable.
type GroupBoundary struct {
	mu               sync.Mutex
	descriptorDigest string
	uncertainCommit  error
	assembler        *snapshot.Assembler
	interval         time.Duration
	notBefore        time.Time
	until            time.Time
	ledger           GroupLedger
	handoffView      groupdelivery.HandoffView
	ownershipView    groupdelivery.OwnershipView
	destination      groupdelivery.FrozenDestination
	layout           rowEncoder
	members          map[string]boundaryMember
	workspaceID      string
	sink             GroupRowSink
	clock            func() time.Time
	onTickError      func(error)
}

func memberKey(deviceID, pointID, tagID string) string {
	encoded, err := json.Marshal([]string{deviceID, pointID, tagID})
	if err != nil {
		return deviceID + "\x1f" + pointID + "\x1f" + tagID
	}
	return string(encoded)
}

// NewGroupBoundary is NewGroupBoundaryContext with a background context.
func NewGroupBoundary(cfg GroupBoundaryConfig) (*GroupBoundary, error) {
	return NewGroupBoundaryContext(context.Background(), cfg)
}

// NewGroupBoundaryContext validates the applied group against real column
// metadata and returns a blocked error instead of a half-working boundary.
// The context bounds restoring durable state.
func NewGroupBoundaryContext(ctx context.Context, cfg GroupBoundaryConfig) (*GroupBoundary, error) {
	blocked := func(code string, issues ...dbtarget.LayoutIssue) error {
		return &GroupBoundaryError{Code: code, Issues: issues}
	}
	group := cfg.Group
	switch {
	case group == nil:
		return nil, blocked(boundaryGroupMissing)
	case cfg.Sink == nil && cfg.Ledger == nil:
		return nil, blocked(boundarySinkMissing)
	case group.AppliedRevision == "" || group.Revision != group.AppliedRevision:
		return nil, blocked(boundaryGroupNotApplied)
	case group.StorageMode == workspace.WriteGroupStorageModeNativeMonthly &&
		((cfg.Dialect != dbtarget.SQLDialectMySQL && cfg.Dialect != dbtarget.SQLDialectSQLServer) || group.EffectiveLayoutVersion() != workspace.WriteGroupLayoutManagedLongV1):
		return nil, blocked(boundaryProfileUnsupported)
	}
	var proofErr error
	cfg, proofErr = resolveBoundaryCounters(cfg)
	if proofErr != nil {
		return nil, counterProofBoundaryError(proofErr)
	}
	if mixedBoundaryGroup(group) {
		return newMixedGroupBoundary(ctx, cfg)
	}
	for _, member := range group.Members {
		if member.ExactRecordingSource != nil && group.EffectiveLayoutVersion() != workspace.WriteGroupLayoutManagedLongV1 {
			return nil, blocked(boundaryMemberMetadata)
		}
	}
	if group.EffectiveLayoutVersion() == workspace.WriteGroupLayoutManagedLongV1 {
		return newLongGroupBoundary(ctx, cfg)
	}
	if group.EffectiveLayoutVersion() == workspace.WriteGroupLayoutManagedWideProfileV1 {
		return newWideProfileGroupBoundary(ctx, cfg)
	}

	members := make(map[string]boundaryMember, len(group.Members))
	layoutMembers := make([]dbtarget.GroupRowMember, 0, len(group.Members))
	snapshotMembers := make([]snapshot.Member, 0, len(group.Members))
	entityKeyed := false
	for _, member := range group.Members {
		if member.DeviceID == "" || member.PointID == "" || member.TagID == "" || member.TargetColumn == "" {
			return nil, blocked(boundaryMemberIdentity)
		}
		kind, ok := measurement.ExactTypeForTag(cfg.TagTypes[member.TagID])
		if !ok {
			return nil, blocked(boundaryTagTypeUnsupport)
		}
		key := memberKey(member.DeviceID, member.PointID, member.TagID)
		members[key] = boundaryMember{key: key, kind: kind}
		layoutMembers = append(layoutMembers, dbtarget.GroupRowMember{
			MemberKey: key, EntityKey: member.EntityKey, Column: member.TargetColumn, Type: kind, Required: member.Required,
		})
		var maxAge time.Duration
		if member.MaxAgeSeconds != nil {
			maxAge = time.Duration(*member.MaxAgeSeconds) * time.Second
		}
		snapshotMembers = append(snapshotMembers, snapshot.Member{
			Key: key, EntityKey: member.EntityKey, Required: member.Required, MaxAge: maxAge,
			SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		})
		entityKeyed = entityKeyed || member.EntityKey != ""
	}

	policy := snapshot.IncompletePolicy(strings.ToLower(strings.TrimSpace(group.RowPolicy.IncompletePolicy)))
	layout, issues := dbtarget.NewGroupRowLayout(boundaryRowSpec(cfg, layoutMembers, entityKeyed))
	if len(issues) > 0 {
		return nil, blocked(boundaryLayoutBlocked, issues...)
	}
	return assembleBoundary(ctx, cfg, boundaryParts{
		members: members, snapshotMembers: snapshotMembers, encoder: layout,
		policy: policy, storage: layout.Capabilities(),
	})
}

// Tick closes every due bucket (including silent ones). Without a ledger the
// results go to the sink and are not durable. With a ledger the closure is
// committed atomically first and memory advances only after that commit, so a
// failed commit leaves the buckets open and the same closure is retried.
func (b *GroupBoundary) Tick(ctx context.Context, now time.Time) error {
	if b.ledger != nil {
		return b.tickDurable(ctx, now)
	}
	b.mu.Lock()
	plan := b.assembler.Plan(now)
	if plan.Err != nil {
		b.mu.Unlock()
		return plan.Err
	}
	b.assembler.Commit(plan)
	outcomes := plan.Outcomes
	b.mu.Unlock()

	var errs []error
	for _, outcome := range outcomes {
		if outcome.Kind != snapshot.OutcomeRow {
			errs = append(errs, b.sink.ReportOutcome(ctx, outcome))
			continue
		}
		encoded, err := b.layout.EncodeRow(outcome)
		if err != nil {
			errs = append(errs, b.sink.ReportOutcome(ctx, blockedOutcome(outcome, err)))
			continue
		}
		errs = append(errs, b.sink.AcceptRow(ctx, GroupRow{Outcome: outcome, Encoded: encoded}))
	}
	return errors.Join(errs...)
}

func (b *GroupBoundary) tickDurable(ctx context.Context, now time.Time) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.uncertainCommit != nil {
		return b.uncertainCommit
	}
	plan := b.assembler.Plan(now)
	if plan.Err != nil {
		return plan.Err
	}
	if plan.Empty() {
		return nil
	}
	buckets := make([]groupdelivery.ClosedBucket, 0, len(plan.Outcomes))
	reported := make([]snapshot.Outcome, 0, len(plan.Outcomes))
	for _, outcome := range plan.Outcomes {
		if outcome.Kind != snapshot.OutcomeRow {
			buckets = append(buckets, groupdelivery.ClosedBucket{Outcome: outcome})
			reported = append(reported, outcome)
			continue
		}
		encoded, err := b.layout.EncodeRow(outcome)
		if err != nil {
			if structuralEncodingFailure(err) {
				return &GroupBoundaryError{Code: boundaryLayoutBlocked, Cause: err}
			}
			blocked := blockedOutcome(outcome, err)
			buckets = append(buckets, groupdelivery.ClosedBucket{Outcome: blocked})
			reported = append(reported, blocked)
			continue
		}
		buckets = append(buckets, groupdelivery.ClosedBucket{Outcome: outcome, Row: &encoded})
	}
	baselines, err := closureBaselines(plan)
	if err != nil {
		return &GroupBoundaryError{Code: boundaryLayoutBlocked, Cause: err}
	}
	closure := groupdelivery.Closure{OwnershipView: b.ownershipView, HandoffView: b.handoffView, Destination: b.destination, NextClose: plan.NextClose(), Buckets: buckets, Baselines: baselines, CounterProofs: b.counterClosureProofs()}
	if err := b.attachMixedClosure(&closure, plan); err != nil {
		return err
	}
	if err := b.ledger.CommitClosure(ctx, closure); err != nil {
		b.quarantineUncertainCommit(err)
		return err
	}
	b.assembler.Commit(plan)
	if b.sink == nil {
		return nil
	}
	var errs []error
	for _, outcome := range reported {
		errs = append(errs, b.sink.ReportOutcome(ctx, outcome))
	}
	return errors.Join(errs...)
}

// RunTicks drives Tick from a clock channel until ctx ends or the channel
// closes; production supplies a time.Ticker, tests supply their own.
func (b *GroupBoundary) RunTicks(ctx context.Context, ticks <-chan time.Time) {
	for {
		select {
		case <-ctx.Done():
			return
		case now, ok := <-ticks:
			if !ok {
				return
			}
			if err := b.Tick(ctx, now); err != nil && b.onTickError != nil {
				b.onTickError(err)
			}
		}
	}
}

var _ SampleSink = (*GroupBoundary)(nil)

// GroupLedger is the durable side of one applied group revision.
type GroupLedger interface {
	AppendSample(ctx context.Context, bucketStart time.Time, sample snapshot.Sample) error
	Restore(ctx context.Context) (groupdelivery.Restored, error)
	CommitClosure(ctx context.Context, closure groupdelivery.Closure) error
}

// Retired reports whether the boundary was given an end (Until) and has closed
// every bucket before it, so it has nothing left to do.
func (b *GroupBoundary) Retired() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.assembler.Done()
}

// dedupeSupported checks a declared dedupe capability against verified storage
// facts: a receipt needs a verified receipt table, a unique key needs a record
// key column that the inspected table marks unique. Anything else is refused
// before activation rather than discovered at delivery time.
func dedupeSupported(capability string, cfg GroupBoundaryConfig) bool {
	switch capability {
	case groupdelivery.DedupeNone:
		return true
	case groupdelivery.DedupeReceipt:
		return cfg.ReceiptTableReady
	case groupdelivery.DedupeUniqueKey:
		if cfg.RecordKeyColumn == "" {
			return false
		}
		for _, column := range cfg.Columns {
			if strings.EqualFold(column.Name, cfg.RecordKeyColumn) {
				return column.Unique || column.PrimaryKey
			}
		}
	}
	return false
}

// SetUntil ends legacy coverage at an aligned boundary and mixed coverage at
// an exact instant across its independent member windows. Accepted future
// observations are excluded too. It never extends an earlier end; invalid
// mixed cutoffs return an error without changing intake or window coverage.
func (b *GroupBoundary) SetUntil(until time.Time) error {
	return b.setCoverageEnd(until, false)
}

// SetEndsAt hands the source to a newer owner at an exact instant that need not
// be window-aligned. Invalid mixed cutoffs return an error without changing
// either accepted-window coverage or the intake cutoff.
func (b *GroupBoundary) SetEndsAt(ends time.Time) error {
	return b.setCoverageEnd(ends, true)
}

func (b *GroupBoundary) setCoverageEnd(end time.Time, exact bool) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	end = end.UTC()
	if _, mixed := b.assembler.MixedProgress(); mixed {
		if _, err := snapshot.CheckedWindowAdd(end, 0, "cutoff"); err != nil {
			return &GroupBoundaryError{Code: "cutoff-invalid", Cause: err}
		}
		if !b.until.IsZero() && !end.Before(b.until) {
			return nil // replay/later request cannot extend accepted coverage
		}
		// Mixed windows have independent phases: this is an actual coverage
		// end, never a shared aligned Until. Filter accepted future samples too.
		if err := b.assembler.SetMixedEndsAt(end); err != nil {
			return &GroupBoundaryError{Code: "cutoff-invalid", Cause: err}
		}
	} else if exact {
		b.assembler.SetEndsAt(end)
	} else {
		b.assembler.SetUntil(end)
	}
	if b.until.IsZero() || end.Before(b.until) {
		b.until = end
	}
	return nil
}

// NextClose is the first bucket that has not closed yet.
func (b *GroupBoundary) NextClose() time.Time {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.assembler.NextClose()
}

// Retiring reports whether the boundary has been given an end and is only
// closing the buckets that remain before it.
func (b *GroupBoundary) Retiring() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return !b.until.IsZero()
}

// Wants reports whether the group consumes this persisted device/point/tag.
func (b *GroupBoundary) Wants(deviceID, pointID, tagID string) bool {
	_, ok := b.members[memberKey(deviceID, pointID, tagID)]
	return ok
}
