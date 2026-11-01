package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// Continue the ordinary reserved successor API through actual disposable DDL
// and the production ledger; no seeded successful schema result is accepted.
func mixedSuccessorConfirmAndActivate(ctx context.Context, t *testing.T, h *circuitProfileHarness, source *WriteGroup, review *WriteGroupSuccessorReview, reserved *WriteGroupSuccessorPreparation, boundary time.Time) {
	t.Helper()
	stage := WriteGroupSuccessorStageRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{
		WorkspaceRevision: reserved.WorkspaceRevision, LineageRevision: reserved.LineageRevision,
		PredecessorRevision: source.Revision, PredecessorConnectorRevision: source.Destination.ConnectorRevision,
		GroupID: reserved.Group.ID, GroupRevision: reserved.Group.Revision, TargetConnectorRevision: reserved.Group.Destination.ConnectorRevision,
	}}
	snapshot, err := h.service.readSuccessorSchemaSnapshot(ctx, source.ID, stage)
	require.NoError(t, err)
	if source.Members[0].Calculation == WriteGroupMemberCalculationDelta {
		require.Len(t, snapshot.counterBindings, 1)
		for _, proof := range snapshot.counterBindings {
			require.Equal(t, reserved.Group.ID, proof.Member.GroupID)
			require.Equal(t, reserved.Group.Revision, proof.Member.AppliedRevision)
			require.NotEqual(t, source.ID, proof.Member.GroupID)
		}
	}
	scope, err := h.service.SuccessorSchemaScope(ctx, source.ID, stage)
	require.NoError(t, err)
	inspector := dbtarget.NewManagedTableInspector(dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(h.db)), h.db)
	inspect := func(ctx context.Context, table string) (recordingplan.TargetTableInspection, error) {
		found, err := inspector.InspectTableAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, table)
		if err != nil {
			return recordingplan.TargetTableInspection{}, err
		}

		if table == scope.GroupLayout.TableName && reserved.Group.isLongLayout() && found.Status == dbtarget.TableInspectionExists {
			indexes, err := inspector.InspectTableIndexesAtRevision(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, table)
			if err != nil {
				return recordingplan.TargetTableInspection{}, err
			}
			found.Indexes, found.IndexesInspected = indexes, true
		}
		return recordingPlanTableInspection(found), nil
	}
	ledger := recordingplan.NewService(recordingplan.NewSQLRepository(h.db))
	token, err := ledger.PrepareSchemaPreview(ctx, scope, inspect)
	require.NoError(t, err)
	verified, err := ledger.ValidatePreviewTokenForApply(ctx, token.Token, scope)
	require.NoError(t, err)
	operation, _, err := ledger.ApplySchemaPreview(ctx, verified, recordingplan.SchemaApplyTarget{Inspect: inspect, Execute: func(ctx context.Context, statements []string) (recordingplan.SchemaExecution, error) {
		done, err := inspector.ExecuteSchemaStatements(ctx, scope.ConnectorID, scope.ConnectorRevision, scope.Schema, statements)
		if done == nil {
			return recordingplan.SchemaExecution{}, err
		}
		return recordingplan.SchemaExecution{Committed: done.Committed, RolledBack: done.RolledBack}, err
	}})
	require.NoError(t, err)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, operation.Status, "%+v", operation)
	storage := &StorageEvidence{Layout: *scope.GroupLayout, CreatedTable: len(token.Statements) > 0}
	confirmed, err := h.service.ConfirmSuccessorSchema(ctx, source.ID, stage, operation.OperationID, operation.VerifiedDigest, storage)
	require.NoError(t, err)
	replay, err := h.service.ConfirmSuccessorSchema(ctx, source.ID, stage, operation.OperationID, operation.VerifiedDigest, storage)
	require.NoError(t, err)
	require.Equal(t, confirmed, replay)
	stage.Current.WorkspaceRevision = confirmed.WorkspaceRevision
	activation, err := h.service.SuccessorActivationReview(ctx, source.ID, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: boundary})
	require.NoError(t, err)
	require.Empty(t, activation.Blockers)
	confirmation := WriteGroupSuccessorActivationConfirmation{Envelope: activation.Envelope, ReviewDigest: activation.ReviewDigest}
	ready, err := h.service.SuccessorActivationReadiness(ctx, source.ID, confirmation)
	require.NoError(t, err)
	require.True(t, ready.Readiness.Ready, "%+v", ready.Readiness.Issues)
	confirmation.ReadinessDigest = ready.ReadinessDigest
	event, err := h.service.ConfirmSuccessorActivation(ctx, source.ID, confirmation)
	require.NoError(t, err)
	require.Equal(t, boundary, event.Boundary)
	unchanged, err := h.service.repo.Get(ctx, source.WorkspaceID, source.ID)
	require.NoError(t, err)
	require.Equal(t, source, unchanged, "future cutover does not rewrite predecessor")
	again, err := h.service.ConfirmSuccessorActivation(ctx, source.ID, confirmation)
	require.NoError(t, err)
	require.Equal(t, event, again)
	current, err := h.service.repo.Get(ctx, source.WorkspaceID, reserved.Group.ID)
	require.NoError(t, err)
	require.Equal(t, source.RecordID, current.RecordID)
	version, err := h.service.repo.versionByRevisionInTx(ctx, h.db, source.WorkspaceID, current.ID, current.AppliedRevision)
	require.NoError(t, err)
	require.Equal(t, boundary, version.EffectiveAt)
	ordinary, err := h.service.Readiness(ctx, current.ID)
	require.NoError(t, err)
	require.True(t, ordinary.Ready, "activated successor retains verifiable schema: %+v", ordinary.Issues)
	record := &Record{ID: current.WorkspaceID, DatabaseSetupRevision: h.mutation(ctx, nil).ExpectedWorkspaceRevision}
	corruptions := map[string]func(*sql.Tx, *WriteGroup){
		"wrong root": func(tx *sql.Tx, _ *WriteGroup) {
			_, err := tx.ExecContext(ctx, `DELETE FROM managed_schema_operations WHERE action='successor_reserve'`)
			require.NoError(t, err)
		},
		"wrong lineage": func(tx *sql.Tx, _ *WriteGroup) {
			_, err := tx.ExecContext(ctx, `UPDATE managed_schema_operations SET detail=json_set(detail,'$.scope.lineage_id','foreign') WHERE action='successor_reserve'`)
			require.NoError(t, err)
		},
		"stale calculation proof": func(_ *sql.Tx, g *WriteGroup) { g.Members[0].Cadence.Interval *= 2 },
	}
	if current.Members[0].Calculation == WriteGroupMemberCalculationDelta {
		corruptions["stale counter proof"] = func(tx *sql.Tx, _ *WriteGroup) {
			_, err := tx.ExecContext(ctx, `UPDATE measurement_definitions SET definition_revision='changed-proof' WHERE id=?`, *current.Members[0].MeasurementID)
			require.NoError(t, err)
		}
	}
	for name, corrupt := range corruptions {
		t.Run(name, func(t *testing.T) {
			tx, err := h.db.BeginTx(ctx, nil)
			require.NoError(t, err)
			defer tx.Rollback()
			candidate := cloneWriteGroup(current)
			corrupt(tx, candidate)
			require.Error(t, h.service.verifyMixedSchemaProofInTx(ctx, tx, record, candidate, operation.OperationID, operation.VerifiedDigest, nil))
		})
	}

	if review.Envelope.Request.OperationID == "next-mixed-successor" {
		return
	}
	t.Run("next generation review", func(t *testing.T) {
		at := boundary.Add(time.Second)
		h.service.WithClock(func() time.Time { return at })
		h.service.repo.now = func() time.Time { return at }
		for name, statement := range map[string]string{
			"unknown historical end":     `DELETE FROM write_group_revision_ends WHERE group_id=?`,
			"future historical end":      `UPDATE write_group_revision_ends SET ends_at='2099-01-01T00:00:00Z' WHERE group_id=?`,
			"foreign historical cutover": `UPDATE write_group_revision_ends SET cutover_id='foreign' WHERE group_id=?`,
		} {
			t.Run(name, func(t *testing.T) {
				tx, err := h.db.BeginTx(ctx, nil)
				require.NoError(t, err)
				defer tx.Rollback()
				_, err = tx.ExecContext(ctx, statement, source.ID)
				require.NoError(t, err)
				require.Error(t, h.service.repo.checkSuccessorProfileOwnersInTx(ctx, tx, current, lifecycleSuccessorCandidate(current, "")))
			})
		}
		lifecycle, err := h.service.Lifecycle(ctx, current.ID)
		require.NoError(t, err)
		require.NotNil(t, lifecycle.SuccessorSeed, "activated successor must support its next reviewed change: %+v", lifecycle.Actions)
		next := writeGroupRecreateReservationRequest{Group: cloneWriteGroup(lifecycle.SuccessorSeed.Group)}
		recreateRepositoryNewTarget(t, h, &next)
		targetPath := filepath.Join(t.TempDir(), "next-target.db")
		db, err := sql.Open("sqlite", targetPath)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, `PRAGMA user_version=0`)
		require.NoError(t, err)
		require.NoError(t, db.Close())
		config, err := json.Marshal(map[string]string{"dsn": targetPath, "database": targetPath, "schema": "main"})
		require.NoError(t, err)
		_, err = h.db.ExecContext(ctx, `UPDATE database_connectors SET connection_config = ? WHERE id = ?`, string(config), next.Group.Destination.ConnectorID)
		require.NoError(t, err)
		nextReview, err := h.service.SuccessorReview(ctx, current.ID, WriteGroupSuccessorReviewRequest{
			WorkspaceID: current.WorkspaceID, ExpectedWorkspaceRevision: h.mutation(ctx, nil).ExpectedWorkspaceRevision,
			LineageID: reserved.LineageID, ExpectedLineageRevision: reserved.LineageRevision + 1,
			ExpectedPredecessorRevision: current.Revision, ExpectedPredecessorConnectorRevision: current.Destination.ConnectorRevision,
			ExpectedTargetConnectorRevision: next.ExpectedTargetConnectorRevision, OperationID: "next-mixed-successor", Boundary: at.Add(time.Minute), Group: next.Group,
		})
		require.NoError(t, err)
		require.Empty(t, nextReview.Blockers)
		nextReserved, err := h.service.ReserveSuccessor(ctx, current.ID, WriteGroupSuccessorReservationRequest{Envelope: nextReview.Envelope, ReviewDigest: nextReview.ReviewDigest})
		require.NoError(t, err)
		require.Equal(t, current.RecordID, nextReserved.Group.RecordID)
		require.Equal(t, reserved.LineageID, nextReserved.LineageID)
		mixedSuccessorConfirmAndActivate(ctx, t, h, current, nextReview, nextReserved, at.Add(time.Minute))
	})

}
