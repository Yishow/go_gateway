package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecreateRepositoryPreservesFullProfileIdentityAndValidation(t *testing.T) {
	for name, mutate := range map[string]func(*testing.T, *circuitProfileHarness, *writeGroupRecreateReservationRequest){
		"profile_kind": func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.RecordingProfile.Kind = WriteGroupProfileCumulative
		},
		"interval": func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.RecordingProfile.Interval++
		},
		"circuit": func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members[0].CircuitID = "different"
		},
		"mapping_revision": func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members[0].MappingRevision = "stale"
		},
		"source_revision": func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members[0].SourceRevision = "stale"
		},
		"missing_source": func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			r.Group.Members[0].PointID = "missing"
		},
		"added_source": func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			source := h.seedSource(h.fixture.deviceID, "added")
			member := r.Group.Members[0]
			member.PointID, member.TagID, member.TargetColumn = source.pointID, source.tagID, "second"
			member.SourceRevision, member.MappingRevision = "", ""
			r.Group.Members = append(r.Group.Members, member)
		},
		"unproven_mysql_family": func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			_, err := h.db.Exec("UPDATE database_connectors SET kind='mysql',connection_config=?", `{"host":"localhost","port":3306,"database":"test"}`)
			require.NoError(t, err)
			r.Group.Destination.Database = "test"
			r.Group.Destination.TableSchema = "test"
			r.Group.Destination.TableName = ""
		},
		"unknown_type_evidence": func(t *testing.T, h *circuitProfileHarness, r *writeGroupRecreateReservationRequest) {
			_, err := h.db.Exec("PRAGMA ignore_check_constraints=ON; UPDATE tags SET data_type='unknown'")
			require.NoError(t, err)
			r.Group.Members[0].SourceRevision = ""
		},
	} {
		t.Run(name, func(t *testing.T) {
			h, req := recreateRepositoryFixture(t, "profile")
			mutate(t, h, &req)
			before := recreateRepositorySnapshot(t, h)
			_, err := reserveRepositoryTest(t, h, req)
			require.Error(t, err)
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func TestRecreateRepositoryRetainsOrdinarySourceOwnership(t *testing.T) {
	h, req := recreateRepositoryFixture(t, "profile")
	other := cloneWriteGroup(req.Group)
	other.Members[0].CircuitID = "other-circuit"
	owner, err := h.service.EnsureCircuitProfile(t.Context(), h.mutation(t.Context(), other))
	require.NoError(t, err)
	require.NotEqual(t, req.PredecessorGroupID, owner.Group.ID)
	before := recreateRepositorySnapshot(t, h)
	_, err = reserveRepositoryTest(t, h, req)
	require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
	// The generic create path retains the same owner refusal after extraction.
	err = h.service.repo.Create(t.Context(), cloneWriteGroup(other))
	require.ErrorIs(t, err, ErrWriteGroupCircuitProfileConflict)
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
}

func TestRecreateRepositoryWideProfileHasFreshPhysicalIdentity(t *testing.T) {
	for _, tableIntent := range []string{"", "echo", "arbitrary"} {
		t.Run(tableIntent, func(t *testing.T) {
			h := newCircuitProfileHarness(t, ":memory:")
			h.db.SetMaxOpenConns(1)
			group := lineageTestDraft(h, "profile")
			group.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
			group.Members[0].ValueRole = WriteGroupValueRoleGauge
			saved, err := h.service.EnsureCircuitProfile(t.Context(), h.mutation(t.Context(), group))
			require.NoError(t, err)
			lineage := requireLineageGenesis(t, h.db, saved.Group)
			_, err = h.db.Exec("UPDATE write_groups SET status='deleted' WHERE id=?", saved.Group.ID)
			require.NoError(t, err)
			req := writeGroupRecreateReservationRequest{WorkspaceID: saved.Group.WorkspaceID, LineageID: lineage, ExpectedLineageRevision: 1,
				PredecessorGroupID: saved.Group.ID, ExpectedPredecessorRevision: saved.Group.Revision,
				ExpectedPredecessorConnectorRevision: saved.Group.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: saved.Group.Destination.ConnectorRevision,
				OperationID: "O2", Group: cloneWriteGroup(saved.Group)}
			if tableIntent != "echo" {
				req.Group.Destination.TableName = tableIntent
			}
			req.Group.AppliedRevision = "untrusted-applied"
			req.Group.Status = WriteGroupStatusRunning
			req.Group.Destination.SchemaRevision, req.Group.Destination.SchemaDigest = "untrusted-schema", "untrusted-proof"
			before := recreateRepositorySnapshot(t, h)
			result, err := reserveRepositoryTest(t, h, req)
			if tableIntent == "arbitrary" {
				require.ErrorIs(t, err, ErrWriteGroupValidation)
				require.Equal(t, before, recreateRepositorySnapshot(t, h))
				return
			}
			require.NoError(t, err)
			require.NotEqual(t, saved.Group.Destination.TableName, result.Group.Destination.TableName)
			require.Equal(t, managedGroupTableName(result.Group.ID), result.Group.Destination.TableName)
			require.Empty(t, result.Group.AppliedRevision)
			require.Empty(t, result.Group.Destination.SchemaRevision)
			require.Empty(t, result.Group.Destination.SchemaDigest)
			require.Equal(t, WriteGroupStatusDraft, result.Group.Status)
		})
	}
}

func TestRecreateRepositorySourceEngineFallbackRequiresSavedConnectorRevision(t *testing.T) {
	h, req := recreateRepositoryFixture(t, "profile")
	recreateRepositoryNewTarget(t, h, &req)
	_, err := h.db.Exec("UPDATE database_connectors SET identity_revision='source-changed' WHERE id=?", h.fixture.connectorID)
	require.NoError(t, err)
	before := recreateRepositorySnapshot(t, h)
	_, err = reserveRepositoryTest(t, h, req)
	require.ErrorIs(t, err, ErrWriteGroupValidation)
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
}
