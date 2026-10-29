package workspace

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"testing"
	"time"
)

type successorDualMySQLProbe struct {
	fakeMySQLProbeInspector
	source *dbtarget.MySQLRecordingCapability
	calls  []string
	during func()
}

func (p *successorDualMySQLProbe) ProbeMySQLRecordingAtRevision(_ context.Context, id, rev string) (*dbtarget.MySQLRecordingCapability, error) {
	p.calls = append(p.calls, id+":"+rev)
	if p.during != nil {
		p.during()
	}
	if p.err != nil {
		return nil, p.err
	}
	if id == p.source.ConnectorID {
		v := *p.source
		return &v, nil
	}
	v := *p.capability
	return &v, nil
}

func successorMySQLNeverAppliedFixture(t *testing.T, family string) (*circuitProfileHarness, string, WriteGroupSuccessorReviewRequest, *successorDualMySQLProbe) {
	t.Helper()
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	a := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/partition_test?tls=true"}`)
	b := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/target_test?tls=true"}`)
	g := h.longDraft(WriteGroupProfileInstantaneous, "mysql-first-start", h.primary())
	g.StorageMode = WriteGroupStorageModePlain
	g.Destination.ConnectorID = a
	g.Destination.ConnectorRevision = "connector-revision-mysql"
	saved := mustEnsureMySQLGroup(t, h, g)
	lineage := requireLineageGenesis(t, h.db, saved)
	_, err := h.db.Exec(`UPDATE write_groups SET status='deleted' WHERE id=?`, saved.ID)
	require.NoError(t, err)
	record, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	candidate := cloneWriteGroup(saved)
	candidate.Destination.ConnectorID = b
	candidate.Destination.Database = "target_test"
	candidate.Destination.TableSchema = "target_test"
	candidate.Destination.TableName = ""
	version := "8.4.13"
	if family == "mariadb" {
		version = "11.4.13-MariaDB-ubu2404"
	}
	source := &dbtarget.MySQLRecordingCapability{ConnectorID: a, ConnectorRevision: "connector-revision-mysql", Database: "partition_test", ServerFamily: family, ServerVersion: version, CapabilityRevision: "source-cap", Connected: true, MetadataStatus: "verified", TransactionalEngine: true, SupportsManagedLong: true, SupportsNativePartition: true}
	target := *source
	target.ConnectorID = b
	target.Database = "target_test"
	target.CapabilityRevision = "target-cap"
	probe := &successorDualMySQLProbe{fakeMySQLProbeInspector: fakeMySQLProbeInspector{capability: &target}, source: source}
	h.service.WithManagedTableInspector(probe)
	return h, saved.ID, WriteGroupSuccessorReviewRequest{WorkspaceID: saved.WorkspaceID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, LineageID: lineage, ExpectedLineageRevision: 1, ExpectedPredecessorRevision: saved.Revision, ExpectedPredecessorConnectorRevision: saved.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: candidate.Destination.ConnectorRevision, OperationID: "mysql-first-O", Boundary: time.Now().Add(time.Hour), Group: candidate}, probe
}

func TestSuccessorNeverAppliedMySQLRequiresAndPersistsDualFamilyProof(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			h, id, req, p := successorMySQLNeverAppliedFixture(t, family)
			view, err := h.service.SuccessorReview(t.Context(), id, req)
			require.NoError(t, err)
			require.Empty(t, view.Blockers)
			require.ElementsMatch(t, []WriteGroupSuccessorBlocker{{ReasonCode: "source-family-check-required", NextAction: "confirm_reservation"}, {ReasonCode: "target-family-check-required", NextAction: "confirm_reservation"}}, view.RequiredChecks)
			require.Empty(t, p.calls)
			saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
			require.NoError(t, err)
			require.Len(t, p.calls, 2)
			require.Empty(t, saved.Group.AppliedRevision)
			stage := WriteGroupSuccessorStageRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{WorkspaceRevision: saved.WorkspaceRevision, LineageRevision: saved.LineageRevision, PredecessorRevision: req.ExpectedPredecessorRevision, PredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision, TargetConnectorRevision: saved.Group.Destination.ConnectorRevision}}
			scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
			require.NoError(t, err)
			require.Equal(t, family, scope.GroupLayout.MySQL.ServerFamily)
		})
	}
}

func TestSuccessorNeverAppliedMySQLFirstStartCommitsWithoutOldOwner(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			h, id, req, p := successorMySQLNeverAppliedFixture(t, family)
			view, err := h.service.SuccessorReview(t.Context(), id, req)
			require.NoError(t, err)
			saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
			require.NoError(t, err)
			stage := WriteGroupSuccessorStageRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{WorkspaceRevision: saved.WorkspaceRevision, LineageRevision: saved.LineageRevision, PredecessorRevision: req.ExpectedPredecessorRevision, PredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision, TargetConnectorRevision: saved.Group.Destination.ConnectorRevision}}
			scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
			require.NoError(t, err)
			repo := recordingplan.NewSQLRepository(h.db)
			ledger := recordingplan.NewService(repo)
			token, err := ledger.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
				return recordingplan.TargetTableInspection{Status: "missing", MySQL: &recordingplan.TargetMySQLTableEvidence{ServerFamily: family, ServerVersion: p.capability.ServerVersion, CapabilityRevision: p.capability.CapabilityRevision, Database: p.capability.Database}}, nil
			})
			require.NoError(t, err)
			op, _, err := ledger.ClaimSchemaApply(t.Context(), token)
			require.NoError(t, err)
			physical, receipt := successorMySQLCatalog(scope)
			digest, err := recordingplan.VerifyManagedMySQLCatalog(*scope.GroupLayout, scope.Database, recordingPlanMySQLInspection(physical), recordingPlanMySQLInspection(receipt))
			require.NoError(t, err)
			_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
			require.NoError(t, err)
			done, err := h.service.ConfirmSuccessorSchema(t.Context(), id, stage, op.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
			require.NoError(t, err)
			stage.Current.WorkspaceRevision = done.WorkspaceRevision
			targetProbe := &successorMySQLProbe{fakeMySQLProbeInspector: fakeMySQLProbeInspector{capability: p.capability}}
			h.service.WithManagedTableInspector(&successorMySQLRuntimeInspector{successorMySQLProbe: targetProbe, group: physical, receipt: receipt})
			request := WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: time.Now().Add(time.Hour)}
			review, err := h.service.SuccessorActivationReview(t.Context(), id, request)
			require.NoError(t, err)
			require.Empty(t, review.Blockers)
			confirmation := WriteGroupSuccessorActivationConfirmation{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}
			ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, confirmation)
			require.NoError(t, err)
			require.True(t, ready.Readiness.Ready, "%+v", ready.Readiness.Issues)
			confirmation.ReadinessDigest = ready.ReadinessDigest
			result, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, confirmation)
			require.NoError(t, err)
			require.Equal(t, saved.Group.ID, result.GroupID)
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_versions WHERE group_id=?`, id))
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_cutovers`))
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_revision_ends`))
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_ownership_fences`))
			h.service.WithManagedTableInspector(&reservedProbeCounter{})
			replay, err := h.service.ConfirmSuccessorFirstStart(t.Context(), id, confirmation)
			require.NoError(t, err)
			require.Equal(t, result, replay)
		})
	}
}

func TestSuccessorNeverAppliedMySQLSourceProofFailsClosed(t *testing.T) {
	for _, failure := range []string{"family", "revision", "database", "unverified", "unsupported", "version", "source-revision-during", "history-during"} {
		t.Run(failure, func(t *testing.T) {
			h, id, req, p := successorMySQLNeverAppliedFixture(t, "mysql")
			view, err := h.service.SuccessorReview(t.Context(), id, req)
			require.NoError(t, err)
			require.Empty(t, view.Blockers)
			switch failure {
			case "family":
				p.source.ServerFamily = "mariadb"
				p.source.ServerVersion = "11.4.13-MariaDB-ubu2404"
			case "revision":
				p.source.ConnectorRevision = "foreign"
			case "database":
				p.source.Database = "foreign"
			case "unverified":
				p.source.MetadataStatus = "unknown"
			case "unsupported":
				p.source.SupportsManagedLong = false
			case "version":
				p.source.ServerVersion = "unrecognized"
			case "source-revision-during":
				p.during = func() {
					if len(p.calls) == 1 {
						_, e := h.db.Exec(`UPDATE database_connectors SET identity_revision='changed' WHERE id=?`, p.source.ConnectorID)
						require.NoError(t, e)
					}
				}
			case "history-during":
				p.during = func() {
					if len(p.calls) == 1 {
						_, e := h.db.Exec(`INSERT INTO wg_delivery_checkpoints(group_id,group_revision,next_close,updated_at) VALUES(?,'hidden-R','2026-10-28T12:00:00Z','2026-10-28T12:00:00Z')`, id)
						require.NoError(t, e)
					}
				}
			}
			_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
			require.Error(t, err)
			require.Equal(t, 1, h.countRows(`SELECT COUNT(*) FROM write_groups`))
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM write_group_lineages WHERE reserved_group_id IS NOT NULL`))
			require.Zero(t, h.countRows(`SELECT COUNT(*) FROM managed_schema_operations WHERE action='successor_reserve'`))
		})
	}
}
