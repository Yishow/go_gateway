package workspace

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// This fixture executes real disposable SQLite DDL but seeds a succeeded local
// ledger result explicitly. Full ledger execution is verified by API tests.
func activationReadyFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorActivationRequest) {
	t.Helper()
	h, id, req := activationFileSourceFixture(t)
	targetPath := filepath.Join(t.TempDir(), "activation.sqlite")
	target, err := sql.Open("sqlite", targetPath)
	require.NoError(t, err)
	t.Cleanup(func() { target.Close() })
	require.NoError(t, target.Ping())
	old := writeGroupRecreateReservationRequest{Group: cloneWriteGroup(req.Group)}
	recreateRepositoryNewTarget(t, h, &old)
	setSQLiteConnectorConfig(t, h.db, old.Group.Destination.ConnectorID, targetPath)
	old.Group.Destination.Database = targetPath
	old.Group.Destination.TableSchema = "main"
	old.Group.Destination.TableName = ""
	req.Group = old.Group
	req.ExpectedTargetConnectorRevision = old.ExpectedTargetConnectorRevision
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, view.Blockers)
	saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
	require.NoError(t, err)
	stage := WriteGroupSuccessorStageRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{WorkspaceRevision: saved.WorkspaceRevision, LineageRevision: saved.LineageRevision, PredecessorRevision: req.ExpectedPredecessorRevision, PredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision, TargetConnectorRevision: saved.Group.Destination.ConnectorRevision}}
	scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
	require.NoError(t, err)
	repo := recordingplan.NewSQLRepository(h.db)
	ledger := recordingplan.NewService(repo)
	token, err := ledger.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
		return recordingplan.TargetTableInspection{Status: "missing"}, nil
	})
	require.NoError(t, err)
	for _, statement := range token.Statements {
		_, err = target.Exec(statement)
		require.NoError(t, err)
	}
	op, _, err := ledger.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	digest := strings.Repeat("a", 64)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	confirmed, err := h.service.ConfirmSuccessorSchema(t.Context(), id, stage, op.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	require.NoError(t, err)
	stage.Current.WorkspaceRevision = confirmed.WorkspaceRevision
	h.service.WithManagedTableInspector(dbtarget.NewManagedTableInspector(dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(h.db)), h.db))
	h.service.now = func() time.Time { return time.Date(2026, 10, 28, 11, 0, 0, 0, time.UTC) }
	return h, id, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: time.Date(2026, 10, 28, 12, 0, 12, 500000000, time.UTC)}
}

func TestSuccessorActivationReadinessExplicitlyVerifiesPreparedTarget(t *testing.T) {
	h, id, req := activationReadyFixture(t)
	before := lineageReadSnapshot(t, h)
	view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, view.Blockers)
	ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, WriteGroupSuccessorActivationConfirmation{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
	require.NoError(t, err)
	require.True(t, ready.Readiness.Ready, "%+v", ready.Readiness.Issues)
	require.NotNil(t, ready.Readiness.runtimeLayout)
	require.Len(t, ready.ReadinessDigest, 64)
	require.Equal(t, before, lineageReadSnapshot(t, h))
	ordinary, err := h.service.Readiness(t.Context(), req.Stage.Current.GroupID)
	require.NoError(t, err)
	require.False(t, ordinary.Ready)
	require.Contains(t, readinessIssueCodes(ordinary.Issues), "successor_reserved")
}

func activationFileSourceFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorReviewRequest) {
	t.Helper()
	h := newCircuitProfileHarness(t, filepath.Join(t.TempDir(), "activation-config.db"))
	h.db.SetMaxOpenConns(1)
	draft := lineageTestDraft(h, "profile")
	draft.LayoutVersion = WriteGroupLayoutManagedWideProfileV1
	for i := range draft.Members {
		draft.Members[i].ValueRole = WriteGroupValueRoleGauge
	}
	saved, err := ensureLineageKind(h, "profile", h.mutation(t.Context(), draft))
	require.NoError(t, err)
	lineage := requireLineageGenesis(t, h.db, saved.Group)
	old := writeGroupRecreateReservationRequest{WorkspaceID: saved.Group.WorkspaceID, PredecessorGroupID: saved.Group.ID, ExpectedPredecessorRevision: saved.Group.Revision, Group: saved.Group}
	seedRecreateEngineEvidence(t, h, old)
	_, err = h.db.Exec("UPDATE write_groups SET status='ready' WHERE id=?", saved.Group.ID)
	require.NoError(t, err)
	h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	return h, saved.Group.ID, WriteGroupSuccessorReviewRequest{WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision, LineageID: lineage, ExpectedLineageRevision: 1, ExpectedPredecessorRevision: saved.Group.Revision, ExpectedPredecessorConnectorRevision: saved.Group.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: saved.Group.Destination.ConnectorRevision, OperationID: "activation-O", Boundary: time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC), Group: cloneWriteGroup(saved.Group)}
}
