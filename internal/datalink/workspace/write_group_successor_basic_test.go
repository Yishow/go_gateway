package workspace

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func recreateBasicFixture(t *testing.T, engine string) (*circuitProfileHarness, writeGroupRecreateReservationRequest) {
	t.Helper()
	return recreateBasicFixtureAt(t, engine, ":memory:")
}

func recreateBasicFixtureAt(t *testing.T, engine, path string) (*circuitProfileHarness, writeGroupRecreateReservationRequest) {
	t.Helper()
	h := newCircuitProfileHarness(t, path)
	h.db.SetMaxOpenConns(1)
	if engine == "postgres" {
		_, err := h.db.Exec(`UPDATE database_connectors SET kind='postgres',connection_config='{"database":"persisted-db","schema":"persisted-schema"}'`)
		require.NoError(t, err)
	}
	group := lineageTestDraft(h, "basic")
	second := h.seedSource(h.fixture.deviceID, "second")
	member := group.Members[0]
	member.PointID, member.TagID, member.TargetColumn = second.pointID, second.tagID, "pressure"
	group.Members = append(group.Members, member)
	saved, err := h.service.EnsureBasicManaged(t.Context(), h.fixture.deviceID, h.mutation(t.Context(), group))
	require.NoError(t, err)
	lineage := requireLineageGenesis(t, h.db, saved.Group)
	_, err = h.db.Exec(`UPDATE write_groups SET status='deleted' WHERE id=?`, saved.Group.ID)
	require.NoError(t, err)
	return h, writeGroupRecreateReservationRequest{WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision,
		LineageID: lineage, ExpectedLineageRevision: 1, PredecessorGroupID: saved.Group.ID, ExpectedPredecessorRevision: saved.Group.Revision,
		ExpectedPredecessorConnectorRevision: saved.Group.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: saved.Group.Destination.ConnectorRevision,
		OperationID: "O2", Group: cloneWriteGroup(saved.Group)}
}

func TestRecreateBasicSameCompleteShape(t *testing.T) {
	for _, engine := range []string{"sqlite", "postgres"} {
		t.Run(engine, func(t *testing.T) {
			h, req := recreateBasicFixture(t, engine)
			before, err := h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
			require.NoError(t, err)
			preserved := recreateRepositorySnapshot(t, h)
			req.Group.Destination.SchemaRevision, req.Group.Destination.SchemaDigest = "old-proof", "old-digest"
			first, err := h.service.reserveRecreatedDraft(t.Context(), req)
			if err != nil {
				var rejected *WriteGroupTransitionError
				if errors.As(err, &rejected) {
					t.Logf("reason: %s source=%#v %#v req=%#v %#v sourceMembers=%#v reqMembers=%#v", rejected.ReasonCode, before.RowPolicy, before.WritePolicy, req.Group.RowPolicy, req.Group.WritePolicy, before.Members, req.Group.Members)
				}
				for cause := err; cause != nil; cause = errors.Unwrap(cause) {
					t.Logf("cause: %v", cause)
				}
			}
			require.NoError(t, err)
			require.EqualValues(t, 2, first.Generation)
			require.Equal(t, req.LineageID, first.LineageID)
			require.NotEqual(t, before.ID, first.Group.ID)
			require.NotEqual(t, before.Destination.TableName, first.Group.Destination.TableName)
			require.NotEmpty(t, first.Group.Destination.TableName)
			require.Equal(t, WriteGroupStatusDraft, first.Group.Status)
			require.Nil(t, first.Group.RecordingProfile)
			require.Empty(t, first.Group.AppliedRevision)
			require.Empty(t, first.Group.Destination.SchemaRevision)
			require.Empty(t, first.Group.Destination.SchemaDigest)
			require.ElementsMatch(t, before.Members, first.Group.Members)
			require.Equal(t, before.RowPolicy, first.Group.RowPolicy)
			require.Equal(t, before.WritePolicy, first.Group.WritePolicy)
			after, err := h.service.repo.Get(t.Context(), req.WorkspaceID, req.PredecessorGroupID)
			require.NoError(t, err)
			require.Equal(t, before, after)
			current := recreateRepositorySnapshot(t, h)
			for _, table := range []string{"write_group_basic_keys", "write_group_profile_keys", "write_group_versions", "wg_runtime_versions", "write_group_ownership_fences"} {
				require.Equal(t, preserved[table], current[table], table)
			}
			var head string
			require.NoError(t, h.db.QueryRow(`SELECT current_group_id FROM write_group_lineages`).Scan(&head))
			require.Equal(t, before.ID, head)
			replay, err := h.service.reserveRecreatedDraft(t.Context(), req)
			require.NoError(t, err)
			require.Equal(t, first, replay)
			req.Group.Destination.TableName = ""
			replay, err = h.service.reserveRecreatedDraft(t.Context(), req)
			require.NoError(t, err, "old-table echo and empty table are one intent")
			require.Equal(t, first, replay)
			mutation := lifecycleReadMutation(t, h.service, first.Group)
			_, err = h.service.Apply(t.Context(), first.Group.ID, mutation)
			requireSuccessorReserved(t, err)
		})
	}
}

func TestRecreateBasicRepositoryAlwaysAssignsNewTable(t *testing.T) {
	h, req := recreateBasicFixture(t, "sqlite")
	result, err := reserveRepositoryTest(t, h, req)
	require.NoError(t, err)
	require.NotEqual(t, req.Group.Destination.TableName, result.Group.Destination.TableName)
	require.Equal(t, managedGroupTableName(result.Group.ID), result.Group.Destination.TableName)
}
