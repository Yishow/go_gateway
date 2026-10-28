package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"testing"
	"time"
)

type successorMySQLProbe struct {
	fakeMySQLProbeInspector
	calls  int
	during func()
}

func (p *successorMySQLProbe) ProbeMySQLRecordingAtRevision(ctx context.Context, id, rev string) (*dbtarget.MySQLRecordingCapability, error) {
	p.calls++
	if p.during != nil {
		p.during()
	}
	return p.fakeMySQLProbeInspector.ProbeMySQLRecordingAtRevision(ctx, id, rev)
}

func successorMySQLFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorReviewRequest, *successorMySQLProbe) {
	return successorMySQLShapeFixture(t, "mysql", WriteGroupStorageModePlain)
}

func successorMySQLShapeFixture(t *testing.T, family string, mode WriteGroupStorageMode) (*circuitProfileHarness, string, WriteGroupSuccessorReviewRequest, *successorMySQLProbe) {
	t.Helper()
	version, collation := "8.4.13", "utf8mb4_0900_bin"
	if family == "mariadb" {
		version, collation = "11.4.13-MariaDB-ubu2404", "utf8mb4_nopad_bin"
	}
	var partitions *recordingplan.MySQLPartitionEvidence
	if mode == WriteGroupStorageModeNativeMonthly {
		partitions = initialNativePartitions(t)
	}
	h := newCircuitProfileHarness(t, ":memory:")
	h.db.SetMaxOpenConns(1)
	a := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/partition_test?tls=true"}`)
	b := seedMySQLConnector(t, h.db, `{"dsn":"fixture:fixture-only@tcp(127.0.0.1:3306)/target_test?tls=true"}`)
	g := h.longDraft(WriteGroupProfileInstantaneous, "mysql-successor", h.primary())
	g.StorageMode = mode
	g.Destination.ConnectorID = a
	g.Destination.ConnectorRevision = "connector-revision-mysql"
	saved := mustEnsureMySQLGroup(t, h, g)
	lineage := requireLineageGenesis(t, h.db, saved)
	table, owner, ok := LongStorageIdentity(saved)
	require.True(t, ok)
	layout, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, recordingplan.MySQLGroupSchemaOptions{ServerFamily: family, ServerVersion: version, CapabilityRevision: "partition-cap-1", Engine: "InnoDB", Charset: "utf8mb4", Collation: collation, NativeMonthly: partitions})
	require.NoError(t, err)
	physical := makeMySQLPhysical(table, owner, layout, partitions)
	receipt := makeMySQLReceiptPhysical()
	for _, physical := range []*dbtarget.MySQLManagedInspection{physical, receipt} {
		physical.Evidence.ServerFamily = family
		physical.Evidence.ServerVersion = version
		physical.Evidence.Collation = collation
		if family == "mariadb" {
			for i := range physical.Inspection.Columns {
				if physical.Inspection.Columns[i].DataType == "JSON" {
					physical.Inspection.Columns[i].DataType = "LONGTEXT"
					physical.Evidence.ColumnCollations[physical.Inspection.Columns[i].Name] = "utf8mb4_bin"
				}
			}
		}
		for key := range physical.Evidence.ColumnCollations {
			physical.Evidence.ColumnCollations[key] = collation
		}
	}
	digest, err := recordingplan.VerifyManagedMySQLCatalog(layout, "partition_test", recordingPlanMySQLInspection(physical), recordingPlanMySQLInspection(receipt))
	require.NoError(t, err)
	saved.Destination.SchemaRevision = "schema-A"
	saved.Destination.SchemaDigest = digest
	saved.WritePolicy.DedupeCapability = "receipt"
	saved.AppliedRevision = saved.Revision
	saved.Status = WriteGroupStatusReady
	tx, err := h.db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	tags, err := loadWriteGroupTagTypes(t.Context(), tx, h.service.repo.query, saved)
	require.NoError(t, err)
	frozen := &WriteGroupRuntimeLayout{Dialect: dbtarget.SQLDialectMySQL, ReceiptTableReady: true, SchemaDigest: digest, Columns: physical.Inspection.Columns, TagTypes: tags, MySQL: &dbtarget.MySQLManagedRowProof{ConnectorID: a, ConnectorRevision: saved.Destination.ConnectorRevision, Group: physical, Receipt: receipt}}
	require.NoError(t, ValidateMySQLRuntimeLayout(saved, frozen))
	require.NoError(t, h.service.repo.insertVersionInTx(t.Context(), tx, saved, time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)))
	require.NoError(t, tx.Commit())
	payload, err := json.Marshal(&WriteGroupAppliedSnapshot{WorkspaceID: saved.WorkspaceID, GroupID: saved.ID, GroupRevision: saved.Revision, AppliedRevision: saved.Revision, Group: saved, RuntimeLayout: frozen})
	require.NoError(t, err)
	_, err = h.db.Exec(`INSERT INTO wg_runtime_versions(workspace_id,group_id,group_revision,payload) VALUES(?,?,?,?)`, saved.WorkspaceID, saved.ID, saved.Revision, string(payload))
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE write_groups SET applied_revision=revision,status='ready',destination_schema_revision=?,destination_schema_digest=? WHERE id=?`, saved.Destination.SchemaRevision, digest, saved.ID)
	require.NoError(t, err)
	h.service.WithBacklogOwnershipGuard(writeGroupBacklogGuardFunc(func(context.Context, *sql.Tx, *WriteGroup) error { return nil }))
	record, err := h.workspace.GetOrCreate(t.Context())
	require.NoError(t, err)
	candidate := cloneWriteGroup(saved)
	candidate.Destination.ConnectorID = b
	candidate.Destination.Database = "target_test"
	candidate.Destination.TableSchema = "target_test"
	candidate.Destination.SchemaRevision = ""
	candidate.Destination.SchemaDigest = ""
	candidate.Destination.TableName = ""
	p := &successorMySQLProbe{fakeMySQLProbeInspector: fakeMySQLProbeInspector{capability: &dbtarget.MySQLRecordingCapability{ConnectorID: b, ConnectorRevision: "connector-revision-mysql", Database: "target_test", ServerFamily: family, ServerVersion: version, CapabilityRevision: "target-cap", Connected: true, MetadataStatus: "verified", TransactionalEngine: true, SupportsManagedLong: true, SupportsNativePartition: true}}}
	h.service.WithManagedTableInspector(p)
	return h, saved.ID, WriteGroupSuccessorReviewRequest{WorkspaceID: saved.WorkspaceID, ExpectedWorkspaceRevision: record.DatabaseSetupRevision, LineageID: lineage, ExpectedLineageRevision: 1, ExpectedPredecessorRevision: saved.Revision, ExpectedPredecessorConnectorRevision: saved.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: candidate.Destination.ConnectorRevision, OperationID: "mysql-O2", Boundary: time.Date(2026, 10, 28, 12, 0, 0, 0, time.UTC), Group: candidate}, p
}

func TestSuccessorMySQLReviewRequiresFamilyWithoutProbe(t *testing.T) {
	h, id, req, p := successorMySQLFixture(t)
	before := lineageReadSnapshot(t, h)
	review, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, review.Blockers)
	raw, err := json.Marshal(review)
	require.NoError(t, err)
	require.Contains(t, string(raw), `"required_checks":[{"reason_code":"target-family-check-required"`)
	require.Zero(t, p.calls)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}
func TestSuccessorMySQLReserveWrongFamilyHasNoEffects(t *testing.T) {
	h, id, req, p := successorMySQLFixture(t)
	review, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	p.capability.ServerFamily = "mariadb"
	p.capability.ServerVersion = "11.4.13-MariaDB-ubu2404"
	before := lineageReadSnapshot(t, h)
	_, err = h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
	require.Error(t, err)
	require.Equal(t, 1, p.calls)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}
func TestSuccessorMySQLExactReplayOffline(t *testing.T) {
	h, id, req, p := successorMySQLFixture(t)
	review, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	request := WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}
	saved, err := h.service.ReserveSuccessor(t.Context(), id, request)
	require.NoError(t, err)
	require.Equal(t, 1, p.calls)
	p.err = errors.New("offline")
	replay, err := h.service.ReserveSuccessor(t.Context(), id, request)
	require.NoError(t, err)
	require.Equal(t, saved.Group.ID, replay.Group.ID)
	require.Equal(t, 1, p.calls)
}
