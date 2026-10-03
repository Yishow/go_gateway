package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

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
	Group             *workspace.WriteGroup
	TagTypes          map[string]schema.DataType
	Dialect           dbtarget.SQLDialect
	Columns           []dbtarget.ColumnInfo
	RecordKeyColumn   string
	BucketStartColumn string
	FirstBucket       time.Time
	// Until, when set, ends this boundary's responsibility: samples observed at
	// or after it belong to a newer applied revision and are ignored here.
	Until time.Time
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
	reasonTypeMismatch      = "type-mismatch"
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
	key  string
	kind measurement.ExactType
}

// GroupBoundary is the opt-in typed-sample consumer for one applied group. It
// keeps open buckets in memory only: they are lost on restart until the
// delivery change journals them, so it must not be described as durable.
type GroupBoundary struct {
	mu          sync.Mutex
	assembler   *snapshot.Assembler
	interval    time.Duration
	notBefore   time.Time
	until       time.Time
	ledger      GroupLedger
	destination groupdelivery.FrozenDestination
	layout      *dbtarget.GroupRowLayout
	members     map[string]boundaryMember
	workspaceID string
	sink        GroupRowSink
	clock       func() time.Time
	onTickError func(error)
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
	layout, issues := dbtarget.NewGroupRowLayout(dbtarget.GroupRowSpec{
		Dialect: cfg.Dialect, Columns: cfg.Columns, Members: layoutMembers,
		Partial:           policy == snapshot.IncompletePartial,
		EntityKeyed:       entityKeyed,
		EntityKeyColumn:   group.RowPolicy.EntityKeyColumn,
		RecordKeyColumn:   cfg.RecordKeyColumn,
		BucketStartColumn: cfg.BucketStartColumn,
		ProvenanceColumn:  group.RowPolicy.ProvenanceColumn,
	})
	if len(issues) > 0 {
		return nil, blocked(boundaryLayoutBlocked, issues...)
	}
	firstBucket := cfg.FirstBucket
	var restored groupdelivery.Restored
	if cfg.Ledger != nil {
		var err error
		if restored, err = cfg.Ledger.Restore(ctx); err != nil {
			return nil, &GroupBoundaryError{Code: boundaryRestoreFailed, Cause: err}
		}
		if restored.NextClose.After(firstBucket) {
			firstBucket = restored.NextClose
		}
	}
	scope, err := snapshot.DestinationScope(
		group.Destination.ConnectorID, group.Destination.ConnectorRevision,
		group.Destination.Database, group.Destination.TableSchema, group.Destination.TableName,
	)
	if err != nil {
		return nil, blocked(boundaryScopeInvalid)
	}
	assembler, err := snapshot.NewAssembler(snapshot.Config{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision, DestinationScope: scope,
		Interval:         time.Duration(group.RowPolicy.IntervalSeconds) * time.Second,
		AllowedLateness:  time.Duration(group.RowPolicy.AllowedLatenessSeconds) * time.Second,
		MaxFutureSkew:    cfg.MaxFutureSkew,
		FirstBucket:      firstBucket,
		Until:            cfg.Until,
		IncompletePolicy: policy,
		Members:          snapshotMembers,
		Storage:          layout.Capabilities(),
	})
	if err != nil {
		return nil, blocked(boundarySnapshotInvalid)
	}
	clock := cfg.Clock
	if clock == nil {
		clock = time.Now
	}
	for _, sample := range restored.Samples {
		// Replay rebuilds memory only; samples of already closed buckets stay out.
		assembler.Replay(sample)
	}
	dedupe := strings.ToLower(strings.TrimSpace(group.WritePolicy.DedupeCapability))
	if dedupe == "" {
		dedupe = groupdelivery.DedupeNone
	}
	if !dedupeSupported(dedupe, cfg) {
		return nil, blocked(boundaryDedupeUnsupport)
	}
	return &GroupBoundary{
		assembler: assembler, interval: time.Duration(group.RowPolicy.IntervalSeconds) * time.Second,
		notBefore: cfg.FirstBucket, until: cfg.Until,
		destination: groupdelivery.FrozenDestination{
			Scope: scope, ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
			Database: group.Destination.Database, TableSchema: group.Destination.TableSchema,
			TableName: group.Destination.TableName, DedupeCapability: dedupe,
		},
		ledger: cfg.Ledger, layout: layout, members: members, workspaceID: group.WorkspaceID,
		sink: cfg.Sink, clock: clock, onTickError: cfg.OnTickError,
	}, nil
}

// AcceptSample routes a typed sample to its bucket. Samples that do not
// belong to this group are ignored, not errors. Refusals (foreign workspace,
// stale revision, late, identity conflict) return a GroupSampleError.
func (b *GroupBoundary) AcceptSample(ctx context.Context, envelope measurement.SampleEnvelope) error {
	member, ok := b.members[memberKey(envelope.DeviceID, envelope.PointID, envelope.TagID)]
	if !ok {
		return nil
	}
	if !b.notBefore.IsZero() && envelope.ObservedAt.Before(b.notBefore) {
		return nil // an earlier applied revision owns this sample
	}
	sample := snapshot.Sample{
		SampleID: envelope.SampleID, MemberKey: member.key, ObservedAt: envelope.ObservedAt,
		Quality: envelope.Quality, QualityReason: envelope.QualityReason,
		SourceRevision: envelope.SourceRevision, MappingRevision: envelope.MappingRevision,
	}
	if envelope.Quality == schema.QualityGood && envelope.Value != nil {
		value, err := measurement.ExactFromGo(member.kind, envelope.Value)
		if err != nil {
			// A value that cannot be carried exactly is bad, never coerced or zeroed.
			sample.Quality, sample.QualityReason = schema.QualityBad, reasonTypeMismatch
		} else {
			sample.Value = value
		}
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.until.IsZero() && !envelope.ObservedAt.Before(b.until) {
		return nil // a newer applied revision owns this sample
	}
	if envelope.WorkspaceID != b.workspaceID {
		return &GroupSampleError{Outcome: snapshot.OfferRejected, Reason: reasonWorkspaceMismatch}
	}
	now := b.clock()
	result := b.assembler.Check(sample, now)
	switch result.Outcome {
	case snapshot.OfferDuplicate:
		return nil
	case snapshot.OfferSelected, snapshot.OfferNotSelected:
	default:
		if result.Outcome == snapshot.OfferLate {
			b.assembler.Offer(sample, now) // counts the late sample; no state changes
		}
		return &GroupSampleError{Outcome: result.Outcome, Reason: result.Reason}
	}
	if b.ledger != nil {
		// The ACK is the commit: refuse, and leave memory untouched, if it fails.
		start := snapshot.BucketStart(sample.ObservedAt, b.interval)
		if err := b.ledger.AppendSample(ctx, start, sample); err != nil {
			reason := reasonJournalFailed
			outcome := snapshot.OfferRejected
			var quotaErr *groupdelivery.QuotaError
			switch {
			case errors.Is(err, groupdelivery.ErrSampleConflict):
				reason, outcome = snapshot.ReasonIdentityConflict, snapshot.OfferConflict
			case errors.As(err, &quotaErr):
				// Capacity refusals keep their exact reason and scope in Cause.
				reason = quotaErr.Reason
			}
			return &GroupSampleError{Outcome: outcome, Reason: reason, Cause: err}
		}
	}
	b.assembler.Offer(sample, now)
	return nil
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
	outcomes := b.assembler.Tick(now)
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
	plan := b.assembler.Plan(now)
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
	closure := groupdelivery.Closure{Destination: b.destination, NextClose: plan.NextClose(), Buckets: buckets}
	if err := b.ledger.CommitClosure(ctx, closure); err != nil {
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

// SetUntil ends the boundary's responsibility at an interval-aligned UTC
// boundary: later samples go to the superseding revision (or nowhere, after a
// disable), and buckets before it still close and deliver. It never extends an
// earlier end.
func (b *GroupBoundary) SetUntil(until time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	until = until.UTC()
	if b.until.IsZero() || until.Before(b.until) {
		b.until = until
	}
	b.assembler.SetUntil(until)
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
