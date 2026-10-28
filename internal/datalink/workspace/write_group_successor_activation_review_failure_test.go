package workspace

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
)

func TestSuccessorActivationReviewBindsEveryActivationField(t *testing.T) {
	h, id, req := activationReadyFixture(t)
	view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, view.Blockers)
	changes := map[string]func(*WriteGroupSuccessorActivationEnvelope){
		"boundary": func(e *WriteGroupSuccessorActivationEnvelope) {
			e.Request.Boundary = e.Request.Boundary.Add(time.Nanosecond)
		},
		"workspace": func(e *WriteGroupSuccessorActivationEnvelope) { e.Request.Stage.Current.WorkspaceRevision = "foreign" },
		"head": func(e *WriteGroupSuccessorActivationEnvelope) {
			e.Request.Stage.Envelope.PredecessorGroupID = "foreign"
		},
		"lineage":           func(e *WriteGroupSuccessorActivationEnvelope) { e.Request.Stage.Current.LineageRevision++ },
		"source":            func(e *WriteGroupSuccessorActivationEnvelope) { e.SourceAppliedRevision = "foreign" },
		"source-definition": func(e *WriteGroupSuccessorActivationEnvelope) { e.SourceDefinitionDigest = "foreign" },
		"target":            func(e *WriteGroupSuccessorActivationEnvelope) { e.Request.Stage.Current.GroupID = "foreign" },
		"target-revision":   func(e *WriteGroupSuccessorActivationEnvelope) { e.Request.Stage.Current.GroupRevision = "foreign" },
		"target-definition": func(e *WriteGroupSuccessorActivationEnvelope) { e.TargetDefinitionDigest = "foreign" },
		"counter":           func(e *WriteGroupSuccessorActivationEnvelope) { e.CounterProofDigest = "foreign" },
		"schema":            func(e *WriteGroupSuccessorActivationEnvelope) { e.SchemaRevision = "foreign" },
		"schema-digest":     func(e *WriteGroupSuccessorActivationEnvelope) { e.SchemaDigest = "foreign" },
		"connector": func(e *WriteGroupSuccessorActivationEnvelope) {
			e.Request.Stage.Current.TargetConnectorRevision = "foreign"
		},
		"version": func(e *WriteGroupSuccessorActivationEnvelope) { e.Version = "foreign" },
	}
	for name, mutate := range changes {
		t.Run(name, func(t *testing.T) {
			raw, err := json.Marshal(view.Envelope)
			require.NoError(t, err)
			var e WriteGroupSuccessorActivationEnvelope
			require.NoError(t, json.Unmarshal(raw, &e))
			mutate(&e)
			probe := &reservedProbeCounter{}
			h.service.WithManagedTableInspector(probe)
			before := lineageReadSnapshot(t, h)
			_, err = h.service.SuccessorActivationReadiness(t.Context(), id, WriteGroupSuccessorActivationConfirmation{Envelope: e, ReviewDigest: view.ReviewDigest})
			require.Error(t, err)
			require.Zero(t, probe.calls)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorActivationReviewNeverInventsFirstAppliedOwner(t *testing.T) {
	h, id, stage, _ := successorSchemaFixture(t)
	before := lineageReadSnapshot(t, h)
	view, err := h.service.SuccessorActivationReview(t.Context(), id, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: time.Now().Add(time.Hour)})
	require.NoError(t, err)
	require.Equal(t, "first_start", view.Envelope.ActivationKind)
	require.Empty(t, view.Envelope.SourceAppliedRevision)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestSuccessorActivationReviewRejectsTamperedSchemaLedgerAndPastBoundary(t *testing.T) {
	for _, mode := range []string{"schema-proof", "past-boundary", "ended-source"} {
		t.Run(mode, func(t *testing.T) {
			h, id, req := activationReadyFixture(t)
			if mode == "schema-proof" {
				_, err := h.db.Exec(`UPDATE managed_schema_operations SET verified_digest='tampered'`)
				require.NoError(t, err)
			}
			if mode == "past-boundary" {
				req.Boundary = h.service.clockNow().Add(-time.Second)
			}
			if mode == "ended-source" {
				_, err := h.db.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?,?)`, req.Stage.Envelope.Request.WorkspaceID, id, req.Stage.Envelope.SourceAppliedRevision, h.service.clockNow().Format(time.RFC3339Nano), "prior-cut")
				require.NoError(t, err)
			}
			before := lineageReadSnapshot(t, h)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
			require.NoError(t, err)
			require.NotEmpty(t, view.Blockers)
			probe := &reservedProbeCounter{}
			h.service.WithManagedTableInspector(probe)
			ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, WriteGroupSuccessorActivationConfirmation{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
			require.NoError(t, err)
			require.False(t, ready.Readiness.Ready)
			require.Zero(t, probe.calls)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

type activationMutationInspector struct {
	delegate ManagedWriteGroupTableInspector
	mutate   func()
	calls    int
}

func (i *activationMutationInspector) InspectTableAtRevision(ctx context.Context, c, r, s, t string) (*dbtarget.TableInspection, error) {
	result, err := i.delegate.InspectTableAtRevision(ctx, c, r, s, t)
	i.calls++
	if i.calls == 1 {
		i.mutate()
	}
	return result, err
}

func TestSuccessorActivationReadinessRejectsChangesDuringProbe(t *testing.T) {
	h, id, req := activationReadyFixture(t)
	view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, view.Blockers)
	probe := &activationMutationInspector{delegate: h.service.managedTableInspector, mutate: func() {
		_, e := h.db.Exec(`UPDATE database_connectors SET identity_revision='changed' WHERE id=?`, req.Stage.Envelope.Request.Group.Destination.ConnectorID)
		require.NoError(t, e)
	}}
	h.service.WithManagedTableInspector(probe)
	_, err = h.service.SuccessorActivationReadiness(t.Context(), id, WriteGroupSuccessorActivationConfirmation{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
	require.Error(t, err)
	require.Positive(t, probe.calls)
}

func TestSuccessorActivationReviewUsesCallerTransaction(t *testing.T) {
	h, id, req := activationReadyFixture(t)
	h.db.SetMaxOpenConns(1)
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	record, err := h.service.workspaceSvc.repo.Get(ctx)
	require.NoError(t, err)
	tx, err := h.db.BeginTx(ctx, nil)
	require.NoError(t, err)
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `UPDATE managed_schema_operations SET verified_digest='uncommitted-invalid'`)
	require.NoError(t, err)
	review, snap, err := h.service.activationReviewInTx(ctx, tx, record, id, req)
	require.NoError(t, err)
	require.Same(t, tx, snap.tx)
	require.Contains(t, review.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "schema-confirmation-unverified", NextAction: "review_lifecycle"})
}

func TestSuccessorActivationReadinessBlockedDoesNotIssueDigest(t *testing.T) {
	h, id, req := activationReadyFixture(t)
	view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
	require.NoError(t, err)
	h.service.managedTableInspector = nil
	before := lineageReadSnapshot(t, h)
	ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, WriteGroupSuccessorActivationConfirmation{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
	require.NoError(t, err)
	require.False(t, ready.Readiness.Ready)
	require.Empty(t, ready.ReadinessDigest)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestSuccessorActivationReadinessHelperRequiresReloader(t *testing.T) {
	h, id, req := activationReadyFixture(t)
	_, snap, err := h.service.readSuccessorActivation(t.Context(), id, req)
	require.NoError(t, err)
	probe := &reservedProbeCounter{}
	h.service.WithManagedTableInspector(probe)
	_, err = h.service.evaluateReadinessSnapshot(t.Context(), snap.reserved.Group.ID, nil, snap.ready, nil)
	require.ErrorIs(t, err, ErrWriteGroupReadinessUnavailable)
	require.Zero(t, probe.calls)
}
