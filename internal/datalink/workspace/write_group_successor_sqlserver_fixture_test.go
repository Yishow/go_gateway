package workspace

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// Unit-only fake SQL Server: SQLite stores configuration/ledger state; no SQL
// Server socket, DDL, write, readback, or real-database acceptance occurs here.
type successorSQLServerUnitInspector struct {
	fakeSQLServerReadinessInspector
	calls, catalogCalls int
	err                 error
	during              func()
}

func (p *successorSQLServerUnitInspector) ProbeSQLServerRecordingAtRevision(context.Context, string, string) (*dbtarget.SQLServerRecordingCapability, error) {
	p.calls++
	return p.capability, p.err
}

func (p *successorSQLServerUnitInspector) InspectSQLServerManagedTableAtRevision(ctx context.Context, connector, revision, database, schema, capability, table string) (*dbtarget.SQLServerManagedInspection, error) {
	p.catalogCalls++
	if p.during != nil {
		p.during()
	}
	return p.fakeSQLServerReadinessInspector.InspectSQLServerManagedTableAtRevision(ctx, connector, revision, database, schema, capability, table)
}

func successorSQLServerUnitFixture(t *testing.T, mode WriteGroupStorageMode) (*circuitProfileHarness, string, WriteGroupSuccessorReviewRequest, *successorSQLServerUnitInspector) {
	t.Helper()
	f := newSQLServerApplyCompatibilityFixture(t, mode)
	h := f.h
	applied, err := h.service.Apply(t.Context(), f.group.ID, f.applyMutation(t.Context(), t))
	require.NoError(t, err)
	h.db.SetMaxOpenConns(1)
	saved := applied.Group
	lineage := requireLineageGenesis(t, h.db, saved)
	b := seedSQLServerConnector(t, h.db, `{"host":"unit-only.invalid","database":"target_recordings","schema":"dbo","username":"unit","password":"unit"}`)
	candidate := cloneWriteGroup(saved)
	candidate.Destination.ConnectorID = b
	candidate.Destination.Database = "target_recordings"
	candidate.Destination.TableName = ""
	candidate.Destination.SchemaRevision, candidate.Destination.SchemaDigest = "", ""
	capability := sqlServerTestCapability(b, saved.Destination.ConnectorRevision)
	capability.Database, capability.CapabilityRevision = "target_recordings", "target-sqlserver-capability"
	p := &successorSQLServerUnitInspector{fakeSQLServerReadinessInspector: fakeSQLServerReadinessInspector{capability: capability}}
	h.service.WithManagedTableInspector(p)
	h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	record, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	return h, saved.ID, WriteGroupSuccessorReviewRequest{WorkspaceID: saved.WorkspaceID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, LineageID: lineage, ExpectedLineageRevision: 1, ExpectedPredecessorRevision: saved.Revision, ExpectedPredecessorConnectorRevision: saved.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: candidate.Destination.ConnectorRevision, OperationID: "sqlserver-unit-O2", Boundary: time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC), Group: candidate}, p
}

func successorSQLServerUnitReserved(t *testing.T, mode WriteGroupStorageMode) (*circuitProfileHarness, string, WriteGroupSuccessorStageRequest, *successorSQLServerUnitInspector) {
	t.Helper()
	h, id, req, p := successorSQLServerUnitFixture(t, mode)
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, view.Blockers)
	saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
	require.NoError(t, err)
	stage := WriteGroupSuccessorStageRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{WorkspaceRevision: saved.WorkspaceRevision, LineageRevision: saved.LineageRevision, PredecessorRevision: req.ExpectedPredecessorRevision, PredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision, TargetConnectorRevision: saved.Group.Destination.ConnectorRevision}}
	return h, id, stage, p
}

func successorSQLServerUnitSucceeded(t *testing.T, h *circuitProfileHarness, scope recordingplan.SchemaPreviewScope, p *successorSQLServerUnitInspector) (*recordingplan.SchemaPreviewToken, string) {
	t.Helper()
	repo := recordingplan.NewSQLRepository(h.db)
	svc := recordingplan.NewService(repo)
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
		options := scope.GroupLayout.SQLServer
		return recordingplan.TargetTableInspection{Status: "missing", SQLServer: &recordingplan.TargetSQLServerTableEvidence{Database: options.Database, Schema: options.Schema, ServerVersion: options.ServerVersion, ServerEdition: options.ServerEdition, EngineEdition: options.EngineEdition, CompatibilityLevel: options.CompatibilityLevel, CapabilityRevision: options.CapabilityRevision, MetadataVerified: true, PartitionInspected: options.NativeMonthly != nil}}, nil
	})
	require.NoError(t, err)
	op, _, err := svc.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	layout := *scope.GroupLayout
	p.group, _ = sqlServerTestRuntimeCatalog(layout.Columns, layout.Indexes, layout.TableName, layout.OwnerColumn, 401, *layout.SQLServer)
	receiptColumns := []recordingplan.GroupSchemaColumn{{Name: "effect_key", SQLType: "CHAR(64)", PrimaryKey: true}, {Name: "payload_digest", SQLType: "CHAR(64)"}, {Name: "committed_at", SQLType: "DATETIME2(6)"}, {Name: "_gw_effect_receipts_v1", SQLType: "NVARCHAR(16)"}}
	p.receipt, _ = sqlServerTestRuntimeCatalog(receiptColumns, nil, dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1", 402, *layout.SQLServer)
	if layout.SQLServer.NativeMonthly != nil {
		p.group.Evidence.NativeMonthly = sqlServerManagedPartitionEvidenceForTest(layout.SQLServer.NativeMonthly)
		p.group.Evidence.PartitionInspected, p.group.Evidence.Plain = true, false
	}
	digest, err := recordingplan.VerifyManagedSQLServerCatalog(layout, scope.Database, recordingPlanSQLServerInspection(p.group), recordingPlanSQLServerInspection(p.receipt))
	require.NoError(t, err)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	return token, digest
}

func successorSQLServerUnitConfirmed(t *testing.T, mode WriteGroupStorageMode) (*circuitProfileHarness, string, WriteGroupSuccessorActivationRequest, *successorSQLServerUnitInspector) {
	t.Helper()
	h, id, stage, p := successorSQLServerUnitReserved(t, mode)
	scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
	require.NoError(t, err)
	token, digest := successorSQLServerUnitSucceeded(t, h, scope, p)
	done, err := h.service.ConfirmSuccessorSchema(t.Context(), id, stage, token.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	require.NoError(t, err)
	stage.Current.WorkspaceRevision = done.WorkspaceRevision
	h.service.now = func() time.Time { return time.Date(2026, 10, 28, 11, 0, 0, 0, time.UTC) }
	return h, id, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: time.Date(2026, 10, 28, 12, 0, 12, 500000000, time.UTC)}, p
}
