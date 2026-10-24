package grouppipeline

import (
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
