package workspace

import (
	"context"
	"errors"
	"go-gateway/internal/datalink/dbtarget"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

func successorMySQLReservedFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorStageRequest, *successorMySQLProbe) {
	t.Helper()
	return successorMySQLReservedShapeFixture(t, "mysql", WriteGroupStorageModePlain)
}

func successorMySQLReservedShapeFixture(t *testing.T, family string, mode WriteGroupStorageMode) (*circuitProfileHarness, string, WriteGroupSuccessorStageRequest, *successorMySQLProbe) {
	t.Helper()
	h, id, req, p := successorMySQLShapeFixture(t, family, mode)
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
	require.NoError(t, err)
	return h, id, WriteGroupSuccessorStageRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{WorkspaceRevision: saved.WorkspaceRevision, LineageRevision: saved.LineageRevision, PredecessorRevision: req.ExpectedPredecessorRevision, PredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision, TargetConnectorRevision: saved.Group.Destination.ConnectorRevision}}, p
}

func TestSuccessorMySQLSchemaFreshFamilyAndCapabilities(t *testing.T) {
	h, id, stage, p := successorMySQLReservedFixture(t)
	before := lineageReadSnapshot(t, h)
	scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
	require.NoError(t, err)
	require.Equal(t, "mysql", scope.GroupLayout.MySQL.ServerFamily)
	require.Equal(t, 2, p.calls, "schema must freshly probe after reserve")
	for _, change := range []func(){
		func() { p.capability.ServerFamily = "mariadb"; p.capability.ServerVersion = "11.4.13-MariaDB-ubu2404" },
		func() {
			p.capability.ServerFamily = "mysql"
			p.capability.ServerVersion = "8.4.13"
			p.capability.SupportsManagedLong = false
		},
		func() { p.capability.SupportsManagedLong = true; p.capability.Database = "foreign" },
	} {
		change()
		_, err = h.service.SuccessorSchemaScope(t.Context(), id, stage)
		require.Error(t, err)
		require.Equal(t, before, lineageReadSnapshot(t, h))
	}
}

func confirmMySQLSuccessorFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorActivationRequest, *successorMySQLProbe) {
	t.Helper()
	return confirmMySQLSuccessorShapeFixture(t, "mysql", WriteGroupStorageModePlain)
}

func confirmMySQLSuccessorShapeFixture(t *testing.T, family string, mode WriteGroupStorageMode) (*circuitProfileHarness, string, WriteGroupSuccessorActivationRequest, *successorMySQLProbe) {
	t.Helper()
	h, id, stage, p := successorMySQLReservedShapeFixture(t, family, mode)
	h.service.now = func() time.Time { return time.Date(2026, 10, 25, 0, 0, 0, 0, time.UTC) }
	scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
	require.NoError(t, err)
	repo := recordingplan.NewSQLRepository(h.db)
	svc := recordingplan.NewService(repo)
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
		return recordingplan.TargetTableInspection{Status: "missing", MySQL: &recordingplan.TargetMySQLTableEvidence{ServerFamily: p.capability.ServerFamily, ServerVersion: p.capability.ServerVersion, CapabilityRevision: p.capability.CapabilityRevision, Database: p.capability.Database}}, nil
	})
	require.NoError(t, err)
	op, _, err := svc.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	physical, receipt := successorMySQLCatalog(scope)
	digest, err := recordingplan.VerifyManagedMySQLCatalog(*scope.GroupLayout, scope.Database, recordingPlanMySQLInspection(physical), recordingPlanMySQLInspection(receipt))
	require.NoError(t, err)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	done, err := h.service.ConfirmSuccessorSchema(t.Context(), id, stage, op.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	require.NoError(t, err)
	stage.Current.WorkspaceRevision = done.WorkspaceRevision
	h.service.WithManagedTableInspector(&successorMySQLRuntimeInspector{successorMySQLProbe: p, group: physical, receipt: receipt})
	h.service.now = func() time.Time { return time.Date(2026, 10, 28, 11, 0, 0, 0, time.UTC) }
	return h, id, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: time.Date(2026, 10, 28, 12, 0, 12, 500000000, time.UTC)}, p
}

func TestSuccessorMySQLActivationReviewUsesLocalSucceededToken(t *testing.T) {
	for _, mutation := range []string{"none", "digest", "source", "family", "status", "owner"} {
		t.Run(mutation, func(t *testing.T) {
			h, id, req, p := confirmMySQLSuccessorFixture(t)
			switch mutation {
			case "digest":
				_, err := h.db.Exec(`UPDATE managed_schema_operations SET verified_digest='foreign'`)
				require.NoError(t, err)
			case "source":
				_, err := h.db.Exec(`UPDATE managed_schema_preview_tokens SET tables=json_set(tables,'$.source_digest','foreign')`)
				require.NoError(t, err)
			case "family":
				_, err := h.db.Exec(`UPDATE managed_schema_preview_tokens SET tables=replace(tables,'mysql','mariadb')`)
				require.NoError(t, err)
			case "status":
				_, err := h.db.Exec(`UPDATE managed_schema_operations SET status='running'`)
				require.NoError(t, err)
			case "owner":
				_, err := h.db.Exec(`UPDATE managed_schema_preview_tokens SET plan_id='foreign'`)
				require.NoError(t, err)
			}
			calls := p.calls
			p.err = errors.New("offline local review must not probe")
			before := lineageReadSnapshot(t, h)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
			require.NoError(t, err)
			if mutation == "none" {
				require.Empty(t, view.Blockers)
			} else {
				require.Contains(t, view.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "schema-confirmation-unverified", NextAction: "review_lifecycle"})
			}
			require.Equal(t, calls, p.calls)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

type successorMySQLRuntimeInspector struct {
	*successorMySQLProbe
	group, receipt *dbtarget.MySQLManagedInspection
	during         func()
}

func (p *successorMySQLRuntimeInspector) InspectMySQLManagedTableAtRevision(_ context.Context, _, _, _, _, table string) (*dbtarget.MySQLManagedInspection, error) {
	if p.during != nil {
		p.during()
	}
	if table == dbtarget.EffectReceiptTable {
		return p.receipt, nil
	}
	return p.group, nil
}
func successorMySQLCatalog(scope recordingplan.SchemaPreviewScope) (*dbtarget.MySQLManagedInspection, *dbtarget.MySQLManagedInspection) {
	layout := scope.GroupLayout
	physical := makeMySQLPhysical(layout.TableName, layout.OwnerColumn, *layout, layout.MySQL.NativeMonthly)
	receipt := makeMySQLReceiptPhysical()
	for _, p := range []*dbtarget.MySQLManagedInspection{physical, receipt} {
		p.Inspection.Schema = scope.Database
		p.Evidence.Database = scope.Database
		p.Evidence.CapabilityRevision = layout.MySQL.CapabilityRevision
		p.Evidence.ServerFamily = layout.MySQL.ServerFamily
		p.Evidence.ServerVersion = layout.MySQL.ServerVersion
		p.Evidence.Collation = layout.MySQL.Collation
		for key := range p.Evidence.ColumnCollations {
			p.Evidence.ColumnCollations[key] = layout.MySQL.Collation
		}
		if layout.MySQL.ServerFamily == "mariadb" {
			for i := range p.Inspection.Columns {
				if p.Inspection.Columns[i].DataType == "JSON" {
					p.Inspection.Columns[i].DataType = "LONGTEXT"
					p.Evidence.ColumnCollations[p.Inspection.Columns[i].Name] = layout.MySQL.Collation
				}
			}
		}
	}
	return physical, receipt
}

func TestSuccessorMySQLActivationCommitsOnlyConfirmedFamily(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
			t.Run(family+"/"+string(mode), func(t *testing.T) {
				h, id, req, _ := confirmMySQLSuccessorShapeFixture(t, family, mode)
				view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
				require.NoError(t, err)
				require.Empty(t, view.Blockers)
				confirmation := WriteGroupSuccessorActivationConfirmation{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest}
				ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, confirmation)
				require.NoError(t, err)
				require.True(t, ready.Readiness.Ready, "%+v", ready.Readiness.Issues)
				confirmation.ReadinessDigest = ready.ReadinessDigest
				event, err := h.service.ConfirmSuccessorActivation(t.Context(), id, confirmation)
				require.NoError(t, err)
				replay, err := h.service.ConfirmSuccessorActivation(t.Context(), id, confirmation)
				require.NoError(t, err)
				require.Equal(t, event.ID, replay.ID)
			})
		}
	}
}

func TestSuccessorMySQLSchemaConfirmReprobesBeforeSavingProof(t *testing.T) {
	h, id, stage, p := successorMySQLReservedFixture(t)
	scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
	require.NoError(t, err)
	repo := recordingplan.NewSQLRepository(h.db)
	svc := recordingplan.NewService(repo)
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
		return recordingplan.TargetTableInspection{Status: "missing", MySQL: &recordingplan.TargetMySQLTableEvidence{ServerFamily: p.capability.ServerFamily, ServerVersion: p.capability.ServerVersion, CapabilityRevision: p.capability.CapabilityRevision, Database: p.capability.Database}}, nil
	})
	require.NoError(t, err)
	op, _, err := svc.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	physical, receipt := successorMySQLCatalog(scope)
	digest, err := recordingplan.VerifyManagedMySQLCatalog(*scope.GroupLayout, scope.Database, recordingPlanMySQLInspection(physical), recordingPlanMySQLInspection(receipt))
	require.NoError(t, err)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	p.capability.ServerFamily = "mariadb"
	p.capability.ServerVersion = "11.4.13-MariaDB-ubu2404"
	before := lineageReadSnapshot(t, h)
	_, err = h.service.ConfirmSuccessorSchema(t.Context(), id, stage, op.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	require.Error(t, err, "durable succeeded token does not replace fresh target family verification")
	require.Equal(t, before, lineageReadSnapshot(t, h))
}
