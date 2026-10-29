package runtime

import (
	"context"
	"errors"
	"strings"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/snapshot"
)

// rowEncoder turns a closed outcome into destination column values. The wide
// row layout and the long record layout both satisfy it.
type rowEncoder interface {
	EncodeRow(snapshot.Outcome) (dbtarget.EncodedRow, error)
}

// boundaryParts is what a layout-specific constructor decides; the rest of
// activation (restore, identity, assembler, dedupe check) is shared.
type boundaryParts struct {
	members         map[string]boundaryMember
	snapshotMembers []snapshot.Member
	encoder         rowEncoder
	policy          snapshot.IncompletePolicy
	mode            snapshot.Mode
	storage         snapshot.StorageCapabilities
	counters        map[string]aggregation.CounterBinding
}

// assembleBoundary restores durable state, freezes the destination identity and
// builds the assembler for an applied group.
func assembleBoundary(ctx context.Context, cfg GroupBoundaryConfig, parts boundaryParts) (*GroupBoundary, error) {
	blocked := func(code string) error { return &GroupBoundaryError{Code: code} }
	group := cfg.Group
	if err := installCounterProofReferences(cfg, parts.members); err != nil {
		return nil, &GroupBoundaryError{Code: err.Error(), Cause: err}
	}
	firstBucket := cfg.FirstBucket
	var restored groupdelivery.Restored
	if cfg.Ledger != nil {
		var err error
		if restored, err = cfg.Ledger.Restore(ctx); err != nil {
			if errors.Is(err, groupdelivery.ErrBaselineCorrupt) {
				return nil, &GroupBoundaryError{Code: boundaryCounterBaselineInvalid, Cause: err}
			}
			return nil, counterProofBoundaryError(err)
		}
		// A fresh Restore must not authorize bounds read before a concurrent
		// cutover. Validate the original plan config before applying its cursor.
		key := groupdelivery.GroupKey{WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.AppliedRevision}
		if err := restored.OwnershipView.ValidateBounds(key, cfg.FirstBucket, cfg.StartsAt, cfg.Until, cfg.EndsAt, time.Duration(group.RowPolicy.IntervalSeconds)*time.Second); err != nil {
			return nil, &GroupBoundaryError{Code: groupdelivery.ErrOwnershipConflict.Error(), Cause: err}
		}
		if restored.NextClose.After(firstBucket) {
			firstBucket = restored.NextClose
		}
	}
	legacyAssociated, err := prepareCounterAssociations(ctx, cfg, restored)
	if err != nil {
		return nil, counterProofBoundaryError(err)
	}
	scope, err := snapshot.DestinationScope(
		group.Destination.ConnectorID, group.Destination.ConnectorRevision,
		group.Destination.Database, group.Destination.TableSchema, group.Destination.TableName,
	)
	if err != nil {
		return nil, blocked(boundaryScopeInvalid)
	}
	startsAt := cfg.StartsAt
	interval := time.Duration(group.RowPolicy.IntervalSeconds) * time.Second
	if !startsAt.IsZero() && !snapshot.BucketStart(startsAt, interval).Equal(firstBucket) {
		startsAt = time.Time{} // that window already closed; ownership routing still uses the exact instant
	}
	assembler, err := snapshot.NewAssembler(snapshot.Config{
		WorkspaceID: group.WorkspaceID, GroupID: group.ID, GroupRevision: group.Revision, DestinationScope: scope,
		Interval:         time.Duration(group.RowPolicy.IntervalSeconds) * time.Second,
		AllowedLateness:  time.Duration(group.RowPolicy.AllowedLatenessSeconds) * time.Second,
		MaxFutureSkew:    cfg.MaxFutureSkew,
		FirstBucket:      firstBucket,
		Until:            cfg.Until,
		StartsAt:         startsAt,
		EndsAt:           cfg.EndsAt,
		IncompletePolicy: parts.policy,
		Mode:             parts.mode,
		Counters:         parts.counters,
		Members:          parts.snapshotMembers,
		Storage:          parts.storage,
	})
	if err != nil {
		return nil, blocked(boundarySnapshotInvalid)
	}
	clock := cfg.Clock
	if clock == nil {
		clock = time.Now
	}
	if parts.mode.Cumulative() {
		if code := restoreCounterBaselines(assembler, restored); code != "" {
			return nil, blocked(code)
		}
	}
	if err := validateCounterReplay(parts.members, restored, legacyAssociated); err != nil {
		return nil, counterProofBoundaryError(err)
	}
	for _, sample := range restored.Samples {
		// All original proofs were checked above, even for excluded inputs.
		// Filtering memory never deletes or relabels the durable accepted set.
		if !assembler.WithinActualBounds(sample.ObservedAt) {
			// Preserve the structural/revision checks that Replay performed
			// before bounds filtering existed; no memory is mutated here.
			checked := assembler.Check(sample, sample.ObservedAt)
			if checked.Outcome != snapshot.OfferSelected && checked.Outcome != snapshot.OfferNotSelected && checked.Outcome != snapshot.OfferDuplicate && checked.Outcome != snapshot.OfferLate {
				return nil, blocked(boundaryRestoreFailed)
			}
			continue
		}
		// Replay rebuilds memory only; samples of already closed buckets stay out.
		result := assembler.Replay(sample)
		if result.Outcome != snapshot.OfferSelected && result.Outcome != snapshot.OfferNotSelected && result.Outcome != snapshot.OfferDuplicate && result.Outcome != snapshot.OfferLate {
			return nil, blocked(boundaryRestoreFailed)
		}
	}
	notBefore, until := cfg.FirstBucket, cfg.Until
	if !cfg.StartsAt.IsZero() {
		notBefore = cfg.StartsAt.UTC()
	}
	if !cfg.EndsAt.IsZero() {
		until = cfg.EndsAt.UTC()
	}
	dedupe := strings.ToLower(strings.TrimSpace(group.WritePolicy.DedupeCapability))
	if dedupe == "" {
		dedupe = groupdelivery.DedupeNone
	}
	if !dedupeSupported(dedupe, cfg) {
		return nil, blocked(boundaryDedupeUnsupport)
	}
	var recordKeyColumn string
	if cfg.Dialect == dbtarget.SQLDialectMySQL || cfg.Dialect == dbtarget.SQLDialectSQLServer {
		recordKeyColumn = group.RowPolicy.RecordKeyColumn
	}
	return &GroupBoundary{
		assembler: assembler, interval: time.Duration(group.RowPolicy.IntervalSeconds) * time.Second,
		notBefore: notBefore, until: until,
		destination: groupdelivery.FrozenDestination{
			Scope: scope, ConnectorID: group.Destination.ConnectorID, ConnectorRevision: group.Destination.ConnectorRevision,
			Database: group.Destination.Database, TableSchema: group.Destination.TableSchema,
			TableName: group.Destination.TableName, DedupeCapability: dedupe, RecordKeyColumn: recordKeyColumn,
		},
		ledger: cfg.Ledger, ownershipView: restored.OwnershipView, handoffView: restored.HandoffView, layout: parts.encoder, members: parts.members, workspaceID: group.WorkspaceID,
		sink: cfg.Sink, clock: clock, onTickError: cfg.OnTickError,
	}, nil
}
