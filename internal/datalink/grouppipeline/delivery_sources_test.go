package grouppipeline

import (
	"encoding/json"
	"go-gateway/internal/datalink/groupdelivery"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func longApplied(revision string) *workspace.WriteGroup {
	group := groupAt(revision)
	group.LayoutVersion = workspace.WriteGroupLayoutManagedLongV1
	group.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileInstantaneous, Interval: 15, Unit: workspace.WriteGroupIntervalUnitSecond}
	group.Members[0].CircuitID, group.Members[1].CircuitID = "cir-a", "cir-a"
	return group
}

func deliveryPipeline(t *testing.T, f *fixture) *Pipeline {
	t.Helper()
	return New(Dependencies{Groups: f.groups, Tags: f.tags, Destinations: f.dest, Inspector: f.inspect, Store: f.store}, Config{NodeID: "node", Owner: "node/test"})
}

func TestSourceEvidenceIsNotReportedForWideGroups(t *testing.T) {
	f := newFixture(t)
	view, err := deliveryPipeline(t, f).DeliveryForRevision(t.Context(), "group-G", "rev-1")
	require.NoError(t, err)
	require.Nil(t, view.SourceEvidence, "wide, custom and share only groups keep their existing delivery view")
}

func TestSourceEvidenceListsTheAppliedSourcesNotTheDraft(t *testing.T) {
	f := newFixture(t)
	applied := longApplied("rev-1")
	f.groups.snapshots["group-G"] = []*workspace.WriteGroupAppliedSnapshot{snapshotOf(applied, at(0))}
	draft := longApplied("rev-1")
	draft.Members = append(draft.Members[:1], workspace.WriteGroupMember{DeviceID: "device-9", PointID: "point-9", TagID: "tag-draft"})
	draft.Destination.ConnectorID, draft.Destination.ConnectorRevision = "draft-connector", "draft-crev"
	f.groups.groups = []*workspace.WriteGroup{liveGroup(draft, workspace.WriteGroupStatusReady)}

	view, err := deliveryPipeline(t, f).DeliveryForRevision(t.Context(), "group-G", "rev-1")
	require.NoError(t, err)
	require.NotNil(t, view.SourceEvidence)
	require.Equal(t, SourceEvidenceAvailable, view.SourceEvidence.State)
	require.Equal(t, "connector-1", view.SourceEvidence.ConnectorID, "the frozen applied destination, not the draft one")
	require.Equal(t, "crev-1", view.SourceEvidence.ConnectorRevision)
	tags := make([]string, 0, len(view.SourceEvidence.Sources))
	for _, source := range view.SourceEvidence.Sources {
		tags = append(tags, source.TagID)
		require.Nil(t, source.Latest, "no record is reported for a source nothing was closed for, and none is invented")
		require.Nil(t, source.LastSQLCommitted)
		require.Equal(t, "cir-a", source.CircuitID)
	}
	require.Equal(t, []string{"tag-t", "tag-p"}, tags, "a draft source never stands in for an applied one")
}

func TestSourceEvidenceSaysUnavailableWhenTheAppliedRevisionCannotBeProven(t *testing.T) {
	t.Run("the applied snapshot cannot be resolved", func(t *testing.T) {
		f := newFixture(t)
		f.groups.snapshots = map[string][]*workspace.WriteGroupAppliedSnapshot{}
		view, err := deliveryPipeline(t, f).DeliveryForRevision(t.Context(), "group-G", "rev-1")
		require.NoError(t, err)
		require.Equal(t, &SourceEvidenceView{State: SourceEvidenceUnavailable, UnavailableReason: SourceEvidenceNoSnapshot, GroupRevision: "rev-1", Sources: []SourceEvidenceItem{}}, view.SourceEvidence)
	})
	t.Run("the snapshot belongs to another applied revision", func(t *testing.T) {
		f := newFixture(t)
		other := longApplied("rev-2")
		f.groups.snapshots["group-G"] = []*workspace.WriteGroupAppliedSnapshot{snapshotOf(other, at(0))}
		view, err := deliveryPipeline(t, f).DeliveryForRevision(t.Context(), "group-G", "rev-1")
		require.NoError(t, err)
		require.Equal(t, SourceEvidenceRevisionMismatch, view.SourceEvidence.UnavailableReason)
		require.Empty(t, view.SourceEvidence.Sources)
	})
	t.Run("the applied snapshot froze no destination", func(t *testing.T) {
		f := newFixture(t)
		applied := longApplied("rev-1")
		applied.Destination.ConnectorRevision = ""
		f.groups.snapshots["group-G"] = []*workspace.WriteGroupAppliedSnapshot{snapshotOf(applied, at(0))}
		view, err := deliveryPipeline(t, f).DeliveryForRevision(t.Context(), "group-G", "rev-1")
		require.NoError(t, err)
		require.Equal(t, SourceEvidenceUnavailable, view.SourceEvidence.State)
		require.Empty(t, view.SourceEvidence.Sources)
	})
	t.Run("a group that was never applied has none", func(t *testing.T) {
		f := newFixture(t)
		view, err := deliveryPipeline(t, f).DeliveryForRevision(t.Context(), "group-G", "")
		require.NoError(t, err)
		require.Nil(t, view.SourceEvidence)
	})
}

// TestSourceEvidenceUsesTheFrozenRuntimeVersion pins that the applied evidence
// is proven by the runtime version Apply froze, not only by the saved revision
// ResolveAppliedAt returns (which never carries the layout). A missing, corrupt
// or mismatched frozen version must stay unavailable and never be completed from
// the current draft.
func TestSourceEvidenceUsesTheFrozenRuntimeVersion(t *testing.T) {
	prepared := func(t *testing.T) (*fixture, *workspace.WriteGroup, *workspace.WriteGroupRuntimeLayout) {
		t.Helper()
		f := newFixture(t)
		applied, layout := longProof(t, workspace.WriteGroupStorageModePlain, meanLongShape(nil))
		f.groups.groups = []*workspace.WriteGroup{liveGroup(applied, workspace.WriteGroupStatusReady)}
		f.groups.snapshots = map[string][]*workspace.WriteGroupAppliedSnapshot{applied.ID: {snapshotOf(applied, at(0))}}
		return f, applied, layout
	}
	save := func(t *testing.T, f *fixture, snap *workspace.WriteGroupAppliedSnapshot) {
		t.Helper()
		payload, err := json.Marshal(snap)
		require.NoError(t, err)
		require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), groupdelivery.GroupKey{
			WorkspaceID: snap.WorkspaceID, GroupID: snap.GroupID, GroupRevision: snap.AppliedRevision,
		}, payload))
	}
	view := func(t *testing.T, f *fixture, group *workspace.WriteGroup) *SourceEvidenceView {
		t.Helper()
		got, err := deliveryPipeline(t, f).DeliveryForRevision(t.Context(), group.ID, group.AppliedRevision)
		require.NoError(t, err)
		return got.SourceEvidence
	}

	t.Run("a verified frozen runtime version proves the applied snapshot", func(t *testing.T) {
		f, applied, layout := prepared(t)
		frozen := snapshotOf(applied, at(0))
		frozen.RuntimeLayout = layout
		save(t, f, frozen)

		evidence := view(t, f, applied)
		require.NotNil(t, evidence)
		require.Equal(t, SourceEvidenceAvailable, evidence.State,
			"the layout Apply froze, not the draft or a probe, is the proof")
	})

	t.Run("a missing frozen version stays unavailable", func(t *testing.T) {
		f, applied, _ := prepared(t)
		evidence := view(t, f, applied)
		require.Equal(t, SourceEvidenceUnavailable, evidence.State)
		require.Equal(t, SourceEvidenceNoSnapshot, evidence.UnavailableReason)
		require.Empty(t, evidence.Sources)
	})

	t.Run("a corrupt frozen version is not repaired from the draft", func(t *testing.T) {
		f, applied, _ := prepared(t)
		require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), groupdelivery.GroupKey{
			WorkspaceID: applied.WorkspaceID, GroupID: applied.ID, GroupRevision: applied.AppliedRevision,
		}, []byte("{")))
		evidence := view(t, f, applied)
		require.Equal(t, SourceEvidenceUnavailable, evidence.State)
		require.Equal(t, SourceEvidenceNoSnapshot, evidence.UnavailableReason)
	})

	t.Run("a frozen version of another revision is not borrowed", func(t *testing.T) {
		f, applied, layout := prepared(t)
		frozen := snapshotOf(applied, at(0))
		frozen.RuntimeLayout = layout
		frozen.GroupRevision, frozen.AppliedRevision = "rev-other", "rev-other"
		save(t, f, frozen)
		evidence := view(t, f, applied)
		require.Equal(t, SourceEvidenceUnavailable, evidence.State)
		require.Equal(t, SourceEvidenceNoSnapshot, evidence.UnavailableReason)
	})
}
