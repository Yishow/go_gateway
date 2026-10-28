package workspace

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

// This is a local configuration/ledger test, not PostgreSQL DDL acceptance.
// The applied source fixture is synthetic; real SQL coverage lives in API/runtime.
func postgresActivationLocalFixture(t *testing.T, confirmed bool) (*circuitProfileHarness, string, WriteGroupSuccessorActivationRequest) {
	t.Helper()
	h, id, req := activationFileSourceFixture(t)
	_, err := h.db.Exec(`UPDATE database_connectors SET kind='postgres' WHERE id=?`, req.Group.Destination.ConnectorID)
	require.NoError(t, err)
	var payload string
	require.NoError(t, h.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=?`, id).Scan(&payload))
	var source WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal([]byte(payload), &source))
	source.RuntimeLayout.Dialect = "postgres"
	raw, err := json.Marshal(source)
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=?`, string(raw), id)
	require.NoError(t, err)
	old := writeGroupRecreateReservationRequest{Group: cloneWriteGroup(req.Group)}
	recreateRepositoryNewTarget(t, h, &old)
	_, err = h.db.Exec(`UPDATE database_connectors SET kind='postgres' WHERE id=?`, old.Group.Destination.ConnectorID)
	require.NoError(t, err)
	req.Group = old.Group
	req.Group.Destination.Database = "gateway_test"
	req.Group.Destination.TableSchema = "pg_successor_unit"
	req.Group.Destination.TableName = ""
	req.ExpectedTargetConnectorRevision = old.ExpectedTargetConnectorRevision
	view, err := h.service.SuccessorReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Empty(t, view.Blockers)
	saved, err := h.service.ReserveSuccessor(t.Context(), id, WriteGroupSuccessorReservationRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest})
	require.NoError(t, err)
	stage := WriteGroupSuccessorStageRequest{Envelope: view.Envelope, ReviewDigest: view.ReviewDigest, Current: WriteGroupSuccessorCurrentCAS{WorkspaceRevision: saved.WorkspaceRevision, LineageRevision: saved.LineageRevision, PredecessorRevision: req.ExpectedPredecessorRevision, PredecessorConnectorRevision: req.ExpectedPredecessorConnectorRevision, GroupID: saved.Group.ID, GroupRevision: saved.Group.Revision, TargetConnectorRevision: saved.Group.Destination.ConnectorRevision}}
	if confirmed {
		scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
		require.NoError(t, err)
		repo := recordingplan.NewSQLRepository(h.db)
		ledger := recordingplan.NewService(repo)
		token, err := ledger.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
			return recordingplan.TargetTableInspection{Status: "missing"}, nil
		})
		require.NoError(t, err)
		op, _, err := ledger.ClaimSchemaApply(t.Context(), token)
		require.NoError(t, err)
		digest := strings.Repeat("b", 64)
		_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
		require.NoError(t, err)
		done, err := h.service.ConfirmSuccessorSchema(t.Context(), id, stage, op.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
		require.NoError(t, err)
		stage.Current.WorkspaceRevision = done.WorkspaceRevision
	}
	h.service.now = func() time.Time { return time.Date(2026, 10, 28, 11, 0, 0, 0, time.UTC) }
	return h, id, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: time.Date(2026, 10, 28, 12, 0, 12, 500000000, time.UTC)}
}

func TestSuccessorActivationPostgresLocalReviewRequiresSchemaProof(t *testing.T) {
	h, id, req := postgresActivationLocalFixture(t, false)
	before := lineageReadSnapshot(t, h)
	view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
	require.NoError(t, err)
	require.Equal(t, []WriteGroupSuccessorBlocker{{ReasonCode: "schema-confirmation-required", NextAction: "review_lifecycle"}}, view.Blockers)
	require.Equal(t, before, lineageReadSnapshot(t, h))
}

func TestSuccessorActivationPostgresLocalReviewVerifiesLedger(t *testing.T) {
	for mutation, query := range map[string]string{
		"none":            "",
		"digest":          `UPDATE managed_schema_operations SET verified_digest='tampered'`,
		"missing":         `DELETE FROM managed_schema_operations`,
		"scope":           `UPDATE managed_schema_operations SET scope_key='foreign'`,
		"token-connector": `UPDATE managed_schema_preview_tokens SET connector_id='foreign'`,
		"token-schema":    `UPDATE managed_schema_preview_tokens SET schema_name='foreign'`,
		"token-dialect":   `UPDATE managed_schema_preview_tokens SET dialect='sqlite'`,
		"token-group":     `UPDATE managed_schema_preview_tokens SET plan_id='foreign'`,
		"token-revision":  `UPDATE managed_schema_preview_tokens SET connector_revision='foreign'`,
	} {
		t.Run(mutation, func(t *testing.T) {
			h, id, req := postgresActivationLocalFixture(t, true)
			if query != "" {
				_, err := h.db.Exec(query)
				require.NoError(t, err)
			}
			before := lineageReadSnapshot(t, h)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
			require.NoError(t, err)
			if mutation == "none" {
				require.Empty(t, view.Blockers)
			} else {
				require.Contains(t, view.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "schema-confirmation-unverified", NextAction: "review_lifecycle"})
			}
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}
