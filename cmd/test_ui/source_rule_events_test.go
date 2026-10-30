package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/diagnostics"
	"testing"
)

func TestProductionSourceEventsAppliedUsesActualSQLiteRuleAndObservedWorkspaceRevision(t *testing.T) {
	f := newSourceEventFixture(t, true)
	outcome, err := f.services.sourceRule.DisableWithRuntimeReconcile(t.Context(), f.rule.ID)
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	records := sourceEventRecords(t, f.broker)
	require.Len(t, records, 1)
	e := records[0]
	require.Equal(t, "domain.source.applied", e.Code)
	require.Equal(t, f.rule.ID, e.Fields["rule_id"])
	require.Equal(t, f.rule.RevisionID, e.Fields["entity_revision"])
	require.Equal(t, f.record.ID, e.Fields["workspace_id"])
	require.Equal(t, f.record.DatabaseSetupRevision, e.Fields["workspace_revision"])
	require.Equal(t, "accepted", e.Fields["stage"])
	require.NotEmpty(t, e.Fields["reference_id"])
}

func TestProductionSourceEventsNotRunningIsUnattributedUnavailable(t *testing.T) {
	f := newSourceEventFixture(t, false)
	outcome, err := f.services.sourceRule.DisableWithRuntimeReconcile(t.Context(), f.rule.ID)
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusNotRunning, outcome.Status)
	records := sourceEventRecords(t, f.broker)
	require.Len(t, records, 1)
	require.Equal(t, "domain.unavailable", records[0].Code)
	require.Empty(t, records[0].Fields)
}

func TestProductionSourceEventsRejectAfterActualShareCompensation(t *testing.T) {
	f := newSourceEventFixture(t, true)
	f.services.sourceRule.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(_ context.Context, req sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		require.Empty(t, sourceEventRecords(t, f.broker), "intermediate generic apply must not escape")
		return sourcerule.RuntimeReconcileOutcome{Scope: req.Scope}, errors.New("password_fixture@private-host")
	}))
	_, err := f.services.sourceRule.DisableWithRuntimeReconcile(t.Context(), f.rule.ID)
	require.Error(t, err)
	restored, getErr := f.services.sourceRule.GetByID(t.Context(), f.rule.ID)
	require.NoError(t, getErr)
	require.Equal(t, f.rule.Enabled, restored.Enabled)
	records := sourceEventRecords(t, f.broker)
	require.Len(t, records, 1)
	require.Equal(t, "domain.source.rejected", records[0].Code)
	require.Equal(t, f.rule.RevisionID, records[0].Fields["entity_revision"])
	raw, e := json.Marshal(records)
	require.NoError(t, e)
	require.NotContains(t, string(raw), "password_fixture")
	require.NotContains(t, string(raw), "private-host")
}

func TestProductionSourceEventsClosedCapturePreservesMutation(t *testing.T) {
	f := newSourceEventFixture(t, true)
	require.NoError(t, f.broker.Close(t.Context()))
	outcome, err := f.services.sourceRule.DisableWithRuntimeReconcile(t.Context(), f.rule.ID)
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	rule, err := f.services.sourceRule.GetByID(t.Context(), f.rule.ID)
	require.NoError(t, err)
	require.False(t, rule.Enabled)
}

func TestProductionSourceEventsFullCapturePreservesMutation(t *testing.T) {
	f := newSourceEventFixture(t, true)
	scope := runtimeview.Scope{WorkspaceID: f.record.ID, WorkspaceRevision: f.record.DatabaseSetupRevision, DeviceID: f.rule.DeviceID}
	for i := 0; i < diagnostics.AdmissionRecords; i++ {
		id := fmt.Sprintf("filler-%d", i)
		ref := diagnostics.DomainEventRef{ContractVersion: 1, Code: "domain.source.rejected", Binding: runtimeview.Binding{WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision, EntityRevision: "R1"}, Locator: runtimeview.EntityLocator{Step: 2, WorkspaceID: scope.WorkspaceID, DeviceID: &scope.DeviceID, RuleID: &id}, Stage: "blocked", Reason: "unavailable", Action: "refresh", Evidence: diagnostics.DomainEvidenceRef{Kind: "source-runtime-transition", ReferenceID: id}}
		raw, err := json.Marshal(ref)
		require.NoError(t, err)
		f.broker.EmitDomain(t.Context(), raw, scope)
	}
	require.NoError(t, f.broker.Flush(t.Context()))
	before, err := f.broker.Snapshot(diagnostics.Query{})
	require.NoError(t, err)
	outcome, err := f.services.sourceRule.DisableWithRuntimeReconcile(t.Context(), f.rule.ID)
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, outcome.Status)
	require.NoError(t, f.broker.Flush(t.Context()))
	after, err := f.broker.Snapshot(diagnostics.Query{})
	require.NoError(t, err)
	require.Greater(t, after.CaptureDropped, before.CaptureDropped)
	rule, err := f.services.sourceRule.GetByID(t.Context(), f.rule.ID)
	require.NoError(t, err)
	require.False(t, rule.Enabled)
}

func TestProductionSourceEventsEmptyShareScopeRemainsUnattributed(t *testing.T) {
	f := newSourceEventFixture(t, true)
	f.services.sourceRule.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(context.Context, sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		return sourcerule.RuntimeReconcileOutcome{}, errors.New("unattributed-share-failure")
	}))
	_, err := f.services.sourceRule.DisableWithRuntimeReconcile(t.Context(), f.rule.ID)
	require.Error(t, err)
	records := sourceEventRecords(t, f.broker)
	require.Len(t, records, 1)
	require.Equal(t, "domain.unavailable", records[0].Code)
	require.Empty(t, records[0].Fields)
}

func TestProductionSourceEventsNotRunningIsNotPromotedByAlignedShare(t *testing.T) {
	f := newSourceEventFixture(t, false)
	f.services.sourceRule.SetShareRuntimeReconciler(sourcerule.ShareRuntimeReconcilerFunc(func(_ context.Context, req sourcerule.RuntimeReconcileRequest) (sourcerule.RuntimeReconcileOutcome, error) {
		return sourcerule.RuntimeReconcileOutcome{Status: sourcerule.RuntimeReconcileStatusAligned, Scope: req.Scope}, nil
	}))
	out, err := f.services.sourceRule.DisableWithRuntimeReconcile(t.Context(), f.rule.ID)
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, out.Status)
	records := sourceEventRecords(t, f.broker)
	require.Len(t, records, 1)
	require.Equal(t, "domain.unavailable", records[0].Code)
}

func TestProductionSourceEventsDeleteRetainsPredeleteRevision(t *testing.T) {
	f := newSourceEventFixture(t, true)
	out, err := f.services.sourceRule.DeleteWithRuntimeReconcile(t.Context(), f.rule.ID)
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, out.Status)
	_, err = f.services.sourceRule.GetByID(t.Context(), f.rule.ID)
	require.Error(t, err)
	records := sourceEventRecords(t, f.broker)
	require.Len(t, records, 1)
	require.Equal(t, "domain.source.applied", records[0].Code)
	require.Equal(t, f.rule.RevisionID, records[0].Fields["entity_revision"])
}

func TestProductionSourceEventsZeroPointCreateIsNotApplied(t *testing.T) {
	f := newSourceEventFixture(t, true)
	_, out, err := f.services.sourceRule.CreateWithRuntimeReconcile(t.Context(), sourcerule.CreateRuleRequest{ID: "no-points", DeviceID: f.rule.DeviceID, StartAddress: "40002", Count: 1, DataType: f.rule.DataType, NamingPrefix: "Empty", SkippedAddresses: []string{"40002"}})
	require.NoError(t, err)
	require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, out.Status)
	records := sourceEventRecords(t, f.broker)
	require.Len(t, records, 1)
	require.Equal(t, "domain.unavailable", records[0].Code)
}
