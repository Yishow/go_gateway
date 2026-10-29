package handlers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"
	"go-gateway/internal/datalink/workspace"
)

type runtimeDeliveryReaderStub struct {
	view                   *grouppipeline.DeliveryView
	err                    error
	group, revision        string
	before, after          []grouppipeline.GroupStatus
	statusCalls, readCalls int
}

func (r *runtimeDeliveryReaderStub) DeliveryForRevision(_ context.Context, group, revision string) (*grouppipeline.DeliveryView, error) {
	r.readCalls++
	r.group, r.revision = group, revision
	return r.view, r.err
}
func (r *runtimeDeliveryReaderStub) Status() []grouppipeline.GroupStatus {
	r.statusCalls++
	if r.statusCalls == 1 {
		return r.before
	}
	return r.after
}

func TestRuntimeRuleDeliveryUsesEffectiveRevisionNotRequestedB(t *testing.T) {
	saved, snapshot, scope, raw, at := runtimeDeliveryProjectionFixture()
	r := &runtimeDeliveryReaderStub{view: raw, before: []grouppipeline.GroupStatus{{GroupID: saved.ID, Revision: "A", State: "active"}, {GroupID: saved.ID, Revision: "B", State: "blocked", Reason: "recording-ownership-conflict"}}}
	r.after = append([]grouppipeline.GroupStatus(nil), r.before...)
	h := (&RuntimeHandler{}).WithRuleDeliveryServices(r)
	got, err := h.readRuntimeRuleDelivery(t.Context(), saved, snapshot, scope, at)
	require.NoError(t, err)
	require.Equal(t, "A", r.revision)
	require.Equal(t, saved.ID, r.group)
	require.Equal(t, &runtimeview.IntakeEvidence{GroupRevision: "A", State: "active"}, got.Current.Intake)
	require.Equal(t, 1, r.readCalls)
	require.Equal(t, 2, r.statusCalls)
	runtimeDeliveryMarshal(t, saved, snapshot, scope, got, at)
}

func TestRuntimeRuleDeliveryExactRevisionIntakeAndEndedSnapshot(t *testing.T) {
	cases := map[string]struct {
		statuses []grouppipeline.GroupStatus
		ended    bool
		want     *runtimeview.IntakeEvidence
	}{
		"missing":             {},
		"healthy-reason-none": {statuses: []grouppipeline.GroupStatus{{GroupID: "group-G", Revision: "A", State: "active", Reason: "none"}}, want: &runtimeview.IntakeEvidence{GroupRevision: "A", State: "active"}},
		"other-revision":      {statuses: []grouppipeline.GroupStatus{{GroupID: "group-G", Revision: "B", State: "active"}}},
		"conflicting":         {statuses: []grouppipeline.GroupStatus{{GroupID: "group-G", Revision: "A", State: "active"}, {GroupID: "group-G", Revision: "A", State: "blocked"}}},
		"ended-active":        {statuses: []grouppipeline.GroupStatus{{GroupID: "group-G", Revision: "A", State: "active"}}, ended: true},
		"unknown-state":       {statuses: []grouppipeline.GroupStatus{{GroupID: "group-G", Revision: "A", State: "private-state"}}},
		"blocked-typed":       {statuses: []grouppipeline.GroupStatus{{GroupID: "group-G", Revision: "A", State: "blocked", Reason: "recording-handoff-proof-conflict"}}, want: &runtimeview.IntakeEvidence{GroupRevision: "A", State: "blocked", Reason: runtimeDeliveryReasonPointer(runtimeview.ReasonHandoffConflict)}},
		"blocked-unknown":     {statuses: []grouppipeline.GroupStatus{{GroupID: "group-G", Revision: "A", State: "blocked", Reason: "password=private-secret"}}, want: &runtimeview.IntakeEvidence{GroupRevision: "A", State: "blocked", Reason: runtimeDeliveryReasonPointer(runtimeview.ReasonUnavailable)}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			saved, snapshot, scope, raw, at := runtimeDeliveryProjectionFixture()
			raw.Intake = grouppipeline.IntakeView{State: "active"}
			if tc.ended {
				snapshot.EndsAt = at
			}
			r := &runtimeDeliveryReaderStub{view: raw, before: tc.statuses, after: tc.statuses}
			got, err := (&RuntimeHandler{}).WithRuleDeliveryServices(r).readRuntimeRuleDelivery(t.Context(), saved, snapshot, scope, at)
			require.NoError(t, err)
			require.Equal(t, tc.want, got.Current.Intake)
			runtimeDeliveryMarshal(t, saved, snapshot, scope, got, at)
		})
	}
}
func runtimeDeliveryReasonPointer(reason runtimeview.Reason) *runtimeview.Reason { return &reason }

func TestRuntimeRuleDeliveryMidreadRelevantIntakeChange(t *testing.T) {
	saved, snapshot, scope, raw, at := runtimeDeliveryProjectionFixture()
	for _, relevant := range []bool{true, false} {
		r := &runtimeDeliveryReaderStub{view: raw, before: []grouppipeline.GroupStatus{{GroupID: saved.ID, Revision: "A", State: "active"}}, after: []grouppipeline.GroupStatus{{GroupID: saved.ID, Revision: "A", State: "active"}}}
		if relevant {
			r.after[0].State = "blocked"
		} else {
			r.after = append(r.after, grouppipeline.GroupStatus{GroupID: "other", Revision: "foreign", State: "blocked"})
		}
		got, err := (&RuntimeHandler{}).WithRuleDeliveryServices(r).readRuntimeRuleDelivery(t.Context(), saved, snapshot, scope, at)
		if relevant {
			require.True(t, runtimeview.IsContractError(err, runtimeview.ReasonStale))
			require.Nil(t, got)
		} else {
			require.NoError(t, err)
			require.Equal(t, "active", got.Current.Intake.State)
		}
	}
}

func TestRuntimeRuleDeliveryProviderErrorNotEmptySuccess(t *testing.T) {
	saved, snapshot, scope, _, at := runtimeDeliveryProjectionFixture()
	for _, r := range []*runtimeDeliveryReaderStub{{err: fmt.Errorf("driver password=private-secret")}, {}} {
		got, err := (&RuntimeHandler{}).WithRuleDeliveryServices(r).readRuntimeRuleDelivery(t.Context(), saved, snapshot, scope, at)
		require.True(t, runtimeview.IsContractError(err, runtimeview.ReasonUnavailable))
		require.Nil(t, got)
		require.NotContains(t, err.Error(), "private-secret")
	}
	got, err := (&RuntimeHandler{}).readRuntimeRuleDelivery(t.Context(), saved, snapshot, scope, at)
	require.NoError(t, err)
	require.Nil(t, got, "optional absent provider preserves configuration-only callers")
}

func TestRuntimeRuleDeliveryMissingProofRequestsEmptyRevision(t *testing.T) {
	saved, _, scope, raw, at := runtimeDeliveryProjectionFixture()
	raw.RevisionStages = nil
	r := &runtimeDeliveryReaderStub{view: raw}
	got, err := (&RuntimeHandler{}).WithRuleDeliveryServices(r).readRuntimeRuleDelivery(t.Context(), saved, nil, scope, at.Add(time.Second))
	require.NoError(t, err)
	require.Empty(t, r.revision)
	require.Nil(t, got.Current)
}

func TestRuntimeRuleDeliveryDisabledSnapshotNeverGrantsActiveIntake(t *testing.T) {
	for _, savedDisabled := range []bool{true, false} {
		saved, snapshot, scope, raw, at := runtimeDeliveryProjectionFixture()
		if savedDisabled {
			saved.Status = workspace.WriteGroupStatusDisabled
		} else {
			snapshot.Group.Status = workspace.WriteGroupStatusDisabled
		}
		statuses := []grouppipeline.GroupStatus{{GroupID: saved.ID, Revision: "A", State: "active"}}
		r := &runtimeDeliveryReaderStub{view: raw, before: statuses, after: statuses}
		got, err := (&RuntimeHandler{}).WithRuleDeliveryServices(r).readRuntimeRuleDelivery(t.Context(), saved, snapshot, scope, at)
		require.NoError(t, err)
		require.Nil(t, got.Current.Intake)
	}
}

func TestRuntimeRuleDeliveryStatusOrderDoesNotMutateOrCauseDrift(t *testing.T) {
	saved, snapshot, scope, raw, at := runtimeDeliveryProjectionFixture()
	statuses := []grouppipeline.GroupStatus{{GroupID: saved.ID, Revision: "A", State: "retiring"}, {GroupID: saved.ID, Revision: "A", State: "blocked"}}
	r := &runtimeDeliveryReaderStub{view: raw, before: statuses, after: []grouppipeline.GroupStatus{statuses[1], statuses[0]}}
	got, err := (&RuntimeHandler{}).WithRuleDeliveryServices(r).readRuntimeRuleDelivery(t.Context(), saved, snapshot, scope, at)
	require.NoError(t, err)
	require.Nil(t, got.Current.Intake)
	require.Equal(t, "retiring", statuses[0].State, "borrowed provider slice order is immutable to reader")
}

type runtimeDeliveryGroupsStub struct {
	list *workspace.WriteGroupListResult
}

func (r runtimeDeliveryGroupsStub) List(context.Context) (*workspace.WriteGroupListResult, error) {
	return r.list, nil
}

func TestRuntimeRuleDeliveryOnlyReadsSelectedDeviceGroups(t *testing.T) {
	h, record := runtimeEvidenceScopeFixture(t)
	saved, _, _, raw, at := runtimeDeliveryProjectionFixture()
	saved.WorkspaceID = record.ID
	saved.Members[0].DeviceID = "device-1"
	unrelated := *saved
	unrelated.ID = "other-group"
	unrelated.Members = []workspace.WriteGroupMember{{DeviceID: "other-device", PointID: "other-point", TagID: "other-tag", TargetColumn: "other"}}
	scope := runtimeview.Scope{WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, DeviceID: "device-1"}
	raw.RevisionStages = nil
	reader := &runtimeDeliveryReaderStub{view: raw}
	h.WithRuleEvidenceServices(runtimeDeliveryGroupsStub{&workspace.WriteGroupListResult{WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision, Groups: []*workspace.WriteGroup{saved, &unrelated}}}, runtimeSourceApplied{}).
		WithRuleDeliveryServices(reader).WithRuleEvidenceClock(func() time.Time { return at })
	evidence, err := h.readRuntimeRuleEvidence(t.Context(), record, scope, nil)
	require.NoError(t, err)
	require.Len(t, evidence.Groups, 1)
	require.Equal(t, saved.ID, reader.group)
	require.Equal(t, 1, reader.readCalls, "unrelated device group must not add a delivery read")
}
