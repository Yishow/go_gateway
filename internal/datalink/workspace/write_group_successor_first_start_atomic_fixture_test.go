package workspace

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// This fixture uses real file-backed configuration SQLite and a fake target
// catalog. G1 is tombstoned without ever applying it or erasing any history.
func firstStartAtomicFixture(t *testing.T) (*circuitProfileHarness, string, WriteGroupSuccessorActivationConfirmation) {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "first-start-config.db")
	h := newCircuitProfileHarness(t, configPath)
	h.db.SetMaxOpenConns(1)
	setSQLiteConnectorConfig(t, h.db, h.fixture.connectorID, configPath)
	draft := lineageTestDraft(h, "basic")
	draft.RowPolicy.IntervalSeconds = 15
	draft.Destination.Database, draft.Destination.TableSchema = configPath, "main"
	saved, err := h.service.EnsureBasicManaged(t.Context(), h.fixture.deviceID, h.mutation(t.Context(), draft))
	require.NoError(t, err)
	lineage := requireLineageGenesis(t, h.db, saved.Group)
	_, err = h.db.Exec(`UPDATE write_groups SET status='deleted' WHERE id=?`, saved.Group.ID)
	require.NoError(t, err)
	old := writeGroupRecreateReservationRequest{WorkspaceID: saved.Group.WorkspaceID, ExpectedWorkspaceRevision: saved.WorkspaceRevision,
		LineageID: lineage, ExpectedLineageRevision: 1, PredecessorGroupID: saved.Group.ID, ExpectedPredecessorRevision: saved.Group.Revision,
		ExpectedPredecessorConnectorRevision: saved.Group.Destination.ConnectorRevision, ExpectedTargetConnectorRevision: saved.Group.Destination.ConnectorRevision,
		OperationID: "first-start-atomic", Group: cloneWriteGroup(saved.Group)}
	boundary := time.Date(2026, 10, 28, 12, 0, 12, 500000000, time.UTC)
	h.service.now = func() time.Time { return boundary.Add(-time.Hour) }
	review, err := h.service.SuccessorReview(t.Context(), old.PredecessorGroupID, WriteGroupSuccessorReviewRequest{
		WorkspaceID: old.WorkspaceID, ExpectedWorkspaceRevision: old.ExpectedWorkspaceRevision,
		LineageID: old.LineageID, ExpectedLineageRevision: old.ExpectedLineageRevision,
		ExpectedPredecessorRevision: old.ExpectedPredecessorRevision, ExpectedPredecessorConnectorRevision: old.ExpectedPredecessorConnectorRevision,
		ExpectedTargetConnectorRevision: old.ExpectedTargetConnectorRevision, OperationID: old.OperationID, Boundary: boundary, Group: old.Group,
	})
	require.NoError(t, err)
	require.Empty(t, review.Blockers)
	reserved, err := h.service.ReserveSuccessor(t.Context(), old.PredecessorGroupID, WriteGroupSuccessorReservationRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest})
	require.NoError(t, err)
	stage := WriteGroupSuccessorStageRequest{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{
		WorkspaceRevision: reserved.WorkspaceRevision, LineageRevision: reserved.LineageRevision,
		PredecessorRevision: old.ExpectedPredecessorRevision, PredecessorConnectorRevision: old.ExpectedPredecessorConnectorRevision,
		GroupID: reserved.Group.ID, GroupRevision: reserved.Group.Revision, TargetConnectorRevision: reserved.Group.Destination.ConnectorRevision,
	}}
	scope, err := h.service.SuccessorSchemaScope(t.Context(), old.PredecessorGroupID, stage)
	require.NoError(t, err)
	repo := recordingplan.NewSQLRepository(h.db)
	ledger := recordingplan.NewService(repo)
	token, err := ledger.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
		return recordingplan.TargetTableInspection{Status: "missing"}, nil
	})
	require.NoError(t, err)
	op, _, err := ledger.ClaimSchemaApply(t.Context(), token)
	require.NoError(t, err)
	digest := strings.Repeat("a", 64)
	_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
	require.NoError(t, err)
	prepared, err := h.service.ConfirmSuccessorSchema(t.Context(), old.PredecessorGroupID, stage, op.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
	require.NoError(t, err)
	stage.Current.WorkspaceRevision = prepared.WorkspaceRevision
	h.service.WithManagedTableInspector(firstStartAtomicCatalog{layout: token.GroupLayout})
	require.Zero(t, h.countRows("SELECT COUNT(*) FROM write_group_versions"))
	require.Zero(t, h.countRows("SELECT COUNT(*) FROM wg_runtime_versions"))
	return h, old.PredecessorGroupID, firstStartAtomicConfirmation(t, h, old.PredecessorGroupID, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: boundary})
}

type firstStartAtomicCatalog struct {
	layout *recordingplan.GroupSchemaLayout
}

func (i firstStartAtomicCatalog) InspectTableAtRevision(_ context.Context, _, _, schema, table string) (*dbtarget.TableInspection, error) {
	if table == dbtarget.EffectReceiptTable {
		return &dbtarget.TableInspection{Status: dbtarget.TableInspectionExists, Schema: schema, Table: table, Columns: []dbtarget.ColumnInfo{{Name: "effect_key", DataType: "TEXT", Nullable: true, PrimaryKey: true}, {Name: "payload_digest", DataType: "TEXT"}, {Name: "committed_at", DataType: "TEXT"}, {Name: "_gw_effect_receipts_v1", DataType: "TEXT"}}}, nil
	}
	return managedTableInspection(i.layout, ""), nil
}

func firstStartAtomicConfirmation(t *testing.T, h *circuitProfileHarness, id string, req WriteGroupSuccessorActivationRequest) WriteGroupSuccessorActivationConfirmation {
	t.Helper()
	review, err := h.service.SuccessorActivationReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, review.Blockers)
	require.Equal(t, "first_start", review.Envelope.ActivationKind)
	require.Empty(t, review.Envelope.SourceAppliedRevision)
	confirmed := WriteGroupSuccessorActivationConfirmation{Envelope: review.Envelope, ReviewDigest: review.ReviewDigest}
	ready, err := h.service.SuccessorActivationReadiness(t.Context(), id, confirmed)
	require.NoError(t, err)
	require.True(t, ready.Readiness.Ready, "%+v", ready.Readiness.Issues)
	require.NotNil(t, ready.Readiness.runtimeLayout)
	confirmed.ReadinessDigest = ready.ReadinessDigest
	return confirmed
}

func firstStartAtomicSnapshot(t *testing.T, h *circuitProfileHarness) map[string]string {
	t.Helper()
	state := recreateRepositorySnapshot(t, h)
	for _, table := range []string{"managed_schema_operations", "managed_schema_preview_tokens", "write_group_cutovers", "write_group_revision_ends", "wg_delivery_samples", "wg_delivery_outbox", "wg_delivery_receipts", "wg_delivery_checkpoints", "wg_delivery_buckets", "wg_delivery_counter_baselines", "wg_counter_proof_associations", "wg_delivery_sample_handoffs"} {
		rows, err := h.db.Query("SELECT * FROM " + table + " ORDER BY rowid")
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		all := [][]any{}
		for rows.Next() {
			values, pointers := make([]any, len(columns)), make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			require.NoError(t, rows.Scan(pointers...))
			all = append(all, values)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		raw, err := json.Marshal(all)
		require.NoError(t, err)
		state[table] = string(raw)
	}
	return state
}

func requireFirstStartNoPredecessorHistory(t *testing.T, h *circuitProfileHarness, id string) {
	t.Helper()
	for _, table := range []string{"write_group_versions", "wg_runtime_versions"} {
		require.Zero(t, h.countRows("SELECT COUNT(*) FROM "+table+" WHERE group_id=?", id), table)
	}
	for _, table := range []string{"write_group_revision_ends", "write_group_cutovers", "write_group_ownership_fences", "wg_delivery_sample_handoffs"} {
		require.Zero(t, h.countRows("SELECT COUNT(*) FROM "+table), table)
	}
}
