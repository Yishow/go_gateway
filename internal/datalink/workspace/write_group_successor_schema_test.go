package workspace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

type successorSchemaAPI interface {
	SuccessorSchemaScope(context.Context, string, WriteGroupSuccessorStageRequest) (recordingplan.SchemaPreviewScope, error)
	ConfirmSuccessorSchema(context.Context, string, WriteGroupSuccessorStageRequest, string, string, *StorageEvidence) (*WriteGroupSaveResult, error)
}

type successorPreparationAPI interface {
	SuccessorReview(context.Context, string, WriteGroupSuccessorReviewRequest) (*WriteGroupSuccessorReview, error)
	ReserveSuccessor(context.Context, string, WriteGroupSuccessorReservationRequest) (*WriteGroupSuccessorPreparation, error)
}

func successorSchemaFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorStageRequest, successorSchemaAPI) {
	t.Helper()
	h, old := recreateBasicFixture(t, "sqlite")
	_, fixtureErr := h.db.Exec(`UPDATE write_groups SET destination_table_schema='main' WHERE id=?`, old.PredecessorGroupID)
	require.NoError(t, fixtureErr)
	old.Group.Destination.TableSchema = "main"
	_, fixtureErr = h.db.Exec(`UPDATE database_connectors SET connection_config=json_set(connection_config,'$.schema','main') WHERE id=?`, old.Group.Destination.ConnectorID)
	require.NoError(t, fixtureErr)
	api, ok := any(h.service).(successorSchemaAPI)
	require.True(t, ok, "reserved successor needs an operation-bound schema preparation service")
	prep, ok := any(h.service).(successorPreparationAPI)
	require.True(t, ok)
	review, err := prep.SuccessorReview(t.Context(), old.PredecessorGroupID, WriteGroupSuccessorReviewRequest{
		WorkspaceID: old.WorkspaceID, ExpectedWorkspaceRevision: old.ExpectedWorkspaceRevision, LineageID: old.LineageID, ExpectedLineageRevision: old.ExpectedLineageRevision,
		ExpectedPredecessorRevision: old.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: old.ExpectedPredecessorConnectorRevision,
		ExpectedTargetConnectorRevision: old.ExpectedTargetConnectorRevision, OperationID: old.OperationID, Boundary: time.Date(2026, 10, 28, 4, 0, 0, 0, time.UTC), Group: old.Group})
	require.NoError(t, err)
	require.Empty(t, review.Blockers)
	saved, err := prep.ReserveSuccessor(t.Context(), old.PredecessorGroupID, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
	require.NoError(t, err)
	return h, old.PredecessorGroupID, WriteGroupSuccessorStageRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{
		WorkspaceRevision: saved.WorkspaceRevision, LineageRevision: saved.LineageRevision, PredecessorRevision: old.ExpectedPredecessorRevision,
		PredecessorConnectorRevision: old.ExpectedPredecessorConnectorRevision, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision, TargetConnectorRevision: saved.Group.Destination.ConnectorRevision}}, api
}

func TestSuccessorSchemaUsesReservedAuthorityWithoutEnablingGroup(t *testing.T) {
	h, id, req, api := successorSchemaFixture(t)
	before := recreateRepositorySnapshot(t, h)
	scope, err := api.SuccessorSchemaScope(t.Context(), id, req)
	require.NoError(t, err)
	require.Equal(t, req.Current.GroupID, scope.PlanID)
	require.Equal(t, "sqlite", scope.Dialect)
	require.NotNil(t, scope.GroupLayout)
	require.NotEmpty(t, scope.SourceDigest)
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
	_, err = h.service.ManagedSchemaScope(t.Context(), req.Current.GroupID, WriteGroupMutation{WorkspaceID: req.Envelope.Request.WorkspaceID, ExpectedWorkspaceRevision: req.Current.WorkspaceRevision, ExpectedGroupRevision: req.Current.GroupRevision, ExpectedConnectorRevision: req.Current.TargetConnectorRevision})
	requireSuccessorReserved(t, err)
	for _, mutate := range []func(*WriteGroupSuccessorStageRequest){
		func(r *WriteGroupSuccessorStageRequest) { r.Current.WorkspaceRevision = "stale" },
		func(r *WriteGroupSuccessorStageRequest) { r.Current.GroupRevision = "stale" },
		func(r *WriteGroupSuccessorStageRequest) { r.Envelope.Request.OperationID = "foreign" },
		func(r *WriteGroupSuccessorStageRequest) { r.ReviewDigest = "foreign" },
	} {
		bad := req
		mutate(&bad)
		_, err = api.SuccessorSchemaScope(t.Context(), id, bad)
		require.Error(t, err)
	}
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
}

func TestSuccessorSchemaRefusesUnverifiedOrForeignProof(t *testing.T) {
	h, id, req, api := successorSchemaFixture(t)
	before := recreateRepositorySnapshot(t, h)
	_, err := api.ConfirmSuccessorSchema(t.Context(), id, req, "unverified-schema-operation", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", &StorageEvidence{})
	require.Error(t, err)
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
	_, err = api.ConfirmSuccessorSchema(t.Context(), id, req, req.Envelope.Request.OperationID, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", &StorageEvidence{})
	require.Error(t, err, "transition ID cannot be a schema operation")
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
}
