package workspace

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func atomicActivationFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorActivationConfirmation) {
	t.Helper()
	h, id, req := activationReadyFixture(t)
	review, err := h.service.SuccessorActivationReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, review.Blockers)
	confirmed := WriteGroupSuccessorActivationConfirmation{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}
	ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, confirmed)
	require.NoError(t, err)
	require.True(t, ready.Readiness.Ready, "%+v", ready.Readiness.Issues)
	confirmed.ReadinessDigest = ready.ReadinessDigest
	return h, id, confirmed
}

func TestSuccessorActivationAtomicCommitAndExactReplay(t *testing.T) {
	h, id, req := atomicActivationFixture(t)
	pred, err := h.service.repo.Get(t.Context(), req.Envelope.Request.Stage.Envelope.Request.WorkspaceID, id)
	require.NoError(t, err)
	event, err := h.service.ConfirmSuccessorActivation(t.Context(), id, req)
	require.NoError(t, err)
	require.Equal(t, WriteGroupCutoverForward, event.Kind)
	require.Equal(t, req.Envelope.Request.Boundary, event.Boundary)
	require.Equal(t, []CutoverEnd{{GroupID: id, Revision: pred.AppliedRevision}}, event.Ended)
	require.Equal(t, []CutoverStart{{GroupID: req.Envelope.Request.Stage.Current.GroupID, Revision: req.Envelope.Request.Stage.Current.GroupRevision}}, event.Started)
	after, err := h.service.repo.Get(t.Context(), pred.WorkspaceID, id)
	require.NoError(t, err)
	require.Equal(t, pred, after, "schedule never stops predecessor early")
	target, err := h.service.repo.Get(t.Context(), pred.WorkspaceID, event.Started[0].GroupID)
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusReady, target.Status)
	require.Equal(t, target.Revision, target.AppliedRevision)
	end, err := h.service.repo.revisionEndInTx(t.Context(), h.db, id, pred.AppliedRevision)
	require.NoError(t, err)
	require.True(t, end.Equal(event.Boundary))
	version, err := h.service.repo.versionByRevisionInTx(t.Context(), h.db, pred.WorkspaceID, target.ID, target.AppliedRevision)
	require.NoError(t, err)
	require.True(t, version.EffectiveAt.Equal(event.Boundary))
	require.Equal(t, len(event.Moved)*2, h.countRows("SELECT COUNT(*) FROM write_group_ownership_fences"))
	var head string
	var pending sql.NullString
	require.NoError(t, h.db.QueryRow(`SELECT current_group_id,reserved_group_id FROM write_group_lineages WHERE lineage_id=?`, req.Envelope.Request.Stage.Envelope.Request.LineageID).Scan(&head, &pending))
	require.Equal(t, target.ID, head)
	require.False(t, pending.Valid)
	before := recreateRepositorySnapshot(t, h)
	h.service.now = func() time.Time { return req.Envelope.Request.Boundary.Add(time.Hour) }
	replay, err := h.service.ConfirmSuccessorActivation(t.Context(), id, req)
	require.NoError(t, err)
	require.Equal(t, event, replay)
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
	changed := req
	changed.ReadinessDigest = "another-readiness"
	_, err = h.service.ConfirmSuccessorActivation(t.Context(), id, changed)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	require.Equal(t, before, recreateRepositorySnapshot(t, h))
}

func TestSuccessorActivationRefusesThirdPartyWideOwner(t *testing.T) {
	for _, futureEnd := range []bool{false, true} {
		t.Run(fmt.Sprint(futureEnd), func(t *testing.T) {
			h, id, req := atomicActivationFixture(t)
			pred, err := h.service.repo.Get(t.Context(), req.Envelope.Request.Stage.Envelope.Request.WorkspaceID, id)
			require.NoError(t, err)
			third := cloneWriteGroup(pred)
			third.ID = "third-party-owner"
			third.Revision = "third-party-revision"
			third.AppliedRevision = third.Revision
			third.BasicManagedDeviceID = ""
			third.RecordingProfile = nil
			third.LayoutVersion = WriteGroupLayoutManagedWideV1
			third.Destination.ConnectorID = "another-destination"
			third.Destination.TableName = "third_party_values"
			tx, err := h.db.BeginTx(context.Background(), nil)
			require.NoError(t, err)
			require.NoError(t, h.service.repo.insertGroup(t.Context(), tx, third))
			require.NoError(t, h.service.repo.insertVersionInTx(t.Context(), tx, third, time.Date(2026, 10, 28, 10, 0, 0, 0, time.UTC)))
			if futureEnd {
				require.NoError(t, h.service.repo.insertRevisionEndInTx(t.Context(), tx, third.WorkspaceID, third.ID, third.Revision, "third-future-end", req.Envelope.Request.Boundary.Add(time.Hour)))
			}
			require.NoError(t, tx.Commit())
			before := recreateRepositorySnapshot(t, h)
			_, err = h.service.ConfirmSuccessorActivation(t.Context(), id, req)
			require.Error(t, err, "only the verified predecessor may overlap; another active wide owner must still block")
			require.Equal(t, before, recreateRepositorySnapshot(t, h))
		})
	}
}

func activationTestService(source *WriteGroupService) *WriteGroupService {
	return &WriteGroupService{workspaceSvc: source.workspaceSvc, repo: source.repo, tableInspector: source.tableInspector, managedTableInspector: source.managedTableInspector, backlogGuard: source.backlogGuard, writerOwnershipActivation: source.writerOwnershipActivation, now: source.now}
}
