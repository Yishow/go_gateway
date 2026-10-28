package workspace

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

func TestSuccessorFirstStartAtomicCommitHasNoFictitiousPredecessorHistory(t *testing.T) {
	h, id, req := firstStartAtomicFixture(t)
	scope := req.Envelope.Request.Stage.Envelope.Request
	predecessor, err := h.service.repo.Get(t.Context(), scope.WorkspaceID, id)
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusDeleted, predecessor.Status)
	require.Empty(t, predecessor.AppliedRevision)
	result, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, req.Envelope.Request.Boundary, result.Boundary)
	require.Equal(t, recordingplan.SuccessorFirstStartOperationID(scope.WorkspaceID, recordingplan.SuccessorReserveOperationID(scope.WorkspaceID, scope.OperationID)), result.OperationID)
	require.Equal(t, "succeeded", result.Status)
	require.EqualValues(t, 2, result.Generation)
	after, err := h.service.repo.Get(t.Context(), scope.WorkspaceID, id)
	require.NoError(t, err)
	require.Equal(t, predecessor, after)
	successor, err := h.service.repo.Get(t.Context(), scope.WorkspaceID, result.GroupID)
	require.NoError(t, err)
	require.Equal(t, WriteGroupStatusReady, successor.Status)
	require.Equal(t, successor.Revision, successor.AppliedRevision)
	version, err := h.service.repo.versionByRevisionInTx(t.Context(), h.db, scope.WorkspaceID, result.GroupID, result.AppliedRevision)
	require.NoError(t, err)
	require.Equal(t, result.Boundary, version.EffectiveAt)
	var head string
	var reserved, operation, intent sql.NullString
	require.NoError(t, h.db.QueryRow(`SELECT current_group_id,reserved_group_id,reservation_operation_id,reservation_intent_digest FROM write_group_lineages WHERE lineage_id=?`, result.LineageID).Scan(&head, &reserved, &operation, &intent))
	require.Equal(t, result.GroupID, head)
	require.False(t, reserved.Valid)
	require.False(t, operation.Valid)
	require.False(t, intent.Valid)
	requireFirstStartNoPredecessorHistory(t, h, id)
	require.Equal(t, 1, h.countRows("SELECT COUNT(*) FROM write_group_versions"))
	require.Equal(t, 1, h.countRows("SELECT COUNT(*) FROM wg_runtime_versions"))
	require.Equal(t, 1, h.countRows("SELECT COUNT(*) FROM managed_schema_operations WHERE action='successor_first_start'"))
	before := firstStartAtomicSnapshot(t, h)
	h.service.now = func() time.Time { return result.Boundary.Add(time.Hour) }
	h.service.managedTableInspector = nil
	replay, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, req)
	require.NoError(t, err)
	require.Equal(t, result, replay)
	require.Equal(t, before, firstStartAtomicSnapshot(t, h))
	changed := req
	changed.ReadinessDigest = strings.Repeat("c", 64)
	_, err = h.service.ConfirmSuccessorFirstStart(t.Context(), id, changed)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	changed = req
	changed.Envelope.Request.Boundary = result.Boundary.Add(time.Nanosecond)
	changed.ReviewDigest, err = successorActivationDigest(changed.Envelope)
	require.NoError(t, err)
	_, _, err = successorFirstStartIntent(id, changed)
	require.NoError(t, err, "changed boundary must be a valid normalized full intent, not just bad digest input")
	_, err = h.service.ConfirmSuccessorFirstStart(t.Context(), id, changed)
	require.ErrorIs(t, err, ErrWriteGroupRevisionConflict)
	require.Equal(t, before, firstStartAtomicSnapshot(t, h))
}
