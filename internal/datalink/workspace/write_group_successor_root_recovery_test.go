package workspace

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSuccessorReserveResultSurvivesClearedReservationAndHeadAdvance(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	confirmation := WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest}
	saved, err := h.service.ReserveSuccessor(t.Context(), id, confirmation)
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_group_lineages SET current_group_id=?,reserved_group_id=NULL,reservation_operation_id=NULL,reservation_intent_digest=NULL,revision=revision+1 WHERE lineage_id=?`, saved.Group.ID, saved.LineageID)
	require.NoError(t, err)
	probe := &reservedProbeCounter{}
	h.service.WithManagedTableInspector(probe)
	replay, err := h.service.ReserveSuccessor(t.Context(), id, confirmation)
	require.NoError(t, err, "immutable root must survive consumed reservation independently of current head")
	require.Equal(t, saved, replay)
	require.Zero(t, probe.calls)
}

func TestSuccessorReserveResultIsImmutableAfterSchemaProgress(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	confirmation := WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest}
	saved, err := h.service.ReserveSuccessor(t.Context(), id, confirmation)
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_groups SET revision='later-schema',destination_schema_revision='later-proof' WHERE id=?`, saved.Group.ID)
	require.NoError(t, err)
	replay, err := h.service.ReserveSuccessor(t.Context(), id, confirmation)
	require.NoError(t, err)
	require.Equal(t, saved, replay, "reserve readback must not substitute mutable latest draft")
}

func TestSuccessorReserveCheckpointRollbackAndLostAcknowledgement(t *testing.T) {
	for _, failure := range []string{"ledger", "workspace_save", "before_commit", "lost_ack"} {
		t.Run(failure, func(t *testing.T) {
			h, id, req := successorReviewFixture(t)
			view, err := h.service.SuccessorReview(t.Context(), id, req)
			require.NoError(t, err)
			confirmation := WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest}
			service := activationTestService(h.service)
			if failure == "ledger" {
				_, err = h.db.Exec(`CREATE TRIGGER reject_root BEFORE INSERT ON managed_schema_operations WHEN NEW.action='successor_reserve' BEGIN SELECT RAISE(ABORT,'root checkpoint failed'); END`)
				require.NoError(t, err)
			} else {
				service.workspaceSvc = NewService(&recreateCommitRepository{SQLRepository: NewSQLRepository(h.db), stage: failure})
			}
			before := recreateRepositorySnapshot(t, h)
			saved, err := service.ReserveSuccessor(t.Context(), id, confirmation)
			if failure == "lost_ack" {
				require.NoError(t, err)
				replay, e := h.service.ReserveSuccessor(t.Context(), id, confirmation)
				require.NoError(t, e)
				require.Equal(t, saved, replay)
				require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM managed_schema_operations WHERE action='successor_reserve'`))
			} else {
				require.Error(t, err)
				require.Nil(t, saved)
				require.Equal(t, before, recreateRepositorySnapshot(t, h))
				require.Zero(t, h.countRows(`SELECT COUNT(*) FROM managed_schema_operations WHERE action='successor_reserve'`))
			}
		})
	}
}

func TestSuccessorReserveSameOperationRejectsDifferentValidFrozenIntent(t *testing.T) {
	h, id, req := successorReviewFixture(t)
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
	require.NoError(t, err)
	changed := view.Envelope
	changed.Request.Group = cloneWriteGroup(changed.Request.Group)
	changed.Request.Group.Name = "different intent"
	changed, review, _, err := normalizeSuccessorEnvelope(changed)
	require.NoError(t, err)
	probe := &reservedProbeCounter{}
	h.service.WithManagedTableInspector(probe)
	_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: changed, ReviewDigest: review})
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	require.Zero(t, probe.calls)
	original, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
	require.NoError(t, err)
	require.Equal(t, saved, original)
	require.Equal(t, 2, h.countRows(`SELECT COUNT(*) FROM write_groups`))
}
