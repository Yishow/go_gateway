package api

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
)

// scopeAPIRig stages the existing public cutover matrix without background
// workers. Configuration and delivery share the real single-connection DB.
type scopeAPIRig struct {
	*cutoverRig
	store           *groupdelivery.Store
	reports         chan error
	dispatcher      *groupdelivery.Dispatcher
	connectors      *dbtarget.ConnectorService
	reconcileBudget time.Duration
}

func newScopeAPIRig(t *testing.T, options cutoverOptions) *scopeAPIRig {
	t.Helper()
	r := newCutoverRigWith(t, "sqlite", options)
	require.NoError(t, r.pipe.Stop(time.Second))
	r.anchor = r.wideEffective.Add(time.Minute).Truncate(time.Minute)
	return wireScopeAPIRig(t, r)
}

func wireScopeAPIRig(t *testing.T, r *cutoverRig) *scopeAPIRig {
	t.Helper()
	return wireScopeAPIWithReconcileBudget(t, r, 5*time.Second)
}

func wireScopeAPIWithReconcileBudget(t *testing.T, r *cutoverRig, budget time.Duration) *scopeAPIRig {
	t.Helper()
	now := func() time.Time { return time.Unix(0, r.clock.Load()).UTC() }
	r.f.groups.WithClock(now)
	s := &scopeAPIRig{cutoverRig: r, store: groupdelivery.NewStore(r.f.db).WithClock(now), reports: make(chan error, 32), reconcileBudget: budget}
	s.connectors = dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(r.f.db))
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(r.f.db))
	r.pipe = grouppipeline.New(grouppipeline.Dependencies{
		Groups: r.f.groups, Tags: tag.NewService(tag.NewSQLRepository(r.f.db)), Destinations: s.connectors,
		Inspector: dbtarget.NewReadOnlyTableInspector(s.connectors), Store: s.store, Ledger: ledger,
		Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(r.f.db)),
	}, grouppipeline.Config{Now: now, MaxFutureSkew: time.Minute, OnError: func(err error) {
		select {
		case s.reports <- err:
		default:
			t.Errorf("unexpected pipeline error overflow: %v", err)
		}
	}})
	pipe := r.pipe
	t.Cleanup(func() { require.NoError(t, pipe.Stop(time.Second)) })
	r.f.router = NewRouter(&DatalinkServices{Workspace: r.f.workspace, WriteGroups: r.f.groups, DBTarget: s.connectors, WriteGroupDelivery: pipe, RecordingPlan: ledger})
	recovery := grouppipeline.NewStorageRecovery(s.store, ledger, s.connectors)
	sender := groupdelivery.NewSender(s.store, s, groupdelivery.SenderConfig{Owner: "public-scope", Recovery: recovery, Now: now})
	s.dispatcher = groupdelivery.NewDispatcher(s.store, sender, groupdelivery.DispatcherConfig{Now: now, MaxPartitions: 1})
	s.reconcile()
	return s
}

func (s *scopeAPIRig) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	allowed := false
	for _, group := range append([]map[string]any{s.wide}, s.cutoverRig.targets...) {
		destination := destinationOf(group)
		allowed = allowed || (item.Destination.ConnectorID == destination["connector_id"] &&
			item.Destination.ConnectorRevision == destination["connector_revision"] && item.Destination.TableName == destination["table_name"])
	}
	if !allowed {
		return groupdelivery.Target{}, fmt.Errorf("unexpected frozen destination table")
	}
	opened, err := s.connectors.OpenDestination(ctx, item.Destination.ConnectorID, item.Destination.ConnectorRevision)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	return groupdelivery.Target{DB: opened.DB, Kind: opened.Kind, Close: opened.Close}, nil
}

func (s *scopeAPIRig) advance(at time.Time) {
	s.t.Helper()
	require.False(s.t, at.Before(time.Unix(0, s.clock.Load())), "controlled time must never run backwards")
	s.cutoverRig.advance(at)
}

func (s *scopeAPIRig) offer(id string, member workspace.WriteGroupMember, observed time.Time, value any) error {
	return s.pipe.AcceptSample(s.t.Context(), measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: s.f.record.ID, DeviceID: member.DeviceID, PointID: member.PointID,
		TagID: member.TagID, SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
		ObservedAt: observed, ReceivedAt: time.Unix(0, s.clock.Load()).UTC(), TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
	})
}

func (s *scopeAPIRig) reconcile() {
	s.t.Helper()
	ctx, cancel := context.WithTimeout(s.t.Context(), s.reconcileBudget)
	defer cancel()
	started := time.Now()
	err := s.pipe.Reconcile(ctx)
	elapsed := time.Since(started)
	if s.reconcileBudget > 5*time.Second || err != nil {
		s.t.Logf("explicit reconcile elapsed=%s budget=%s context=%v deadline_chain=%t error=%T:%v", elapsed, s.reconcileBudget, ctx.Err(), errors.Is(err, context.DeadlineExceeded), err, err)
	}
	require.NoError(s.t, err, "explicit reconcile elapsed=%s context=%v deadline_chain=%t", elapsed, ctx.Err(), errors.Is(err, context.DeadlineExceeded))
	require.Empty(s.t, s.reports)
}

func (s *scopeAPIRig) deliver() groupdelivery.CycleReport {
	s.t.Helper()
	ctx, cancel := context.WithTimeout(s.t.Context(), 10*time.Second)
	defer cancel()
	report, err := s.dispatcher.RunOnce(ctx)
	require.NoError(s.t, err)
	require.Zero(s.t, report.Retrying+report.Blocked+report.Quarantined+report.Unknown)
	return report
}

func (s *scopeAPIRig) flush(at time.Time) groupdelivery.CycleReport {
	s.t.Helper()
	s.advance(at)
	s.pipe.TickAll(s.t.Context())
	require.Empty(s.t, s.reports)
	return s.deliver()
}
