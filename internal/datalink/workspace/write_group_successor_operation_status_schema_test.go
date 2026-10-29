package workspace

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

func TestSuccessorOperationStatusRejectsCorruptSchemaBindings(t *testing.T) {
	for _, mutation := range []string{"digest", "statements", "connector", "layout", "operation-action", "operation-token", "operation-scope", "operation-workspace", "operation-proof"} {
		t.Run(mutation, func(t *testing.T) {
			h, id, stage, _ := successorSchemaFixture(t)
			token, _, ledger := statusPreview(t, h, id, stage)
			_, _, err := ledger.ClaimSchemaApply(t.Context(), token)
			require.NoError(t, err)
			switch mutation {
			case "digest":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET digest='corrupt' WHERE token=?`, token.Token)
			case "statements":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET statements='[]' WHERE token=?`, token.Token)
			case "connector":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET connector_revision='foreign' WHERE token=?`, token.Token)
			case "layout":
				_, err = h.db.Exec(`UPDATE managed_schema_preview_tokens SET tables=json_set(tables,'$.group_layout.table_name','foreign') WHERE token=?`, token.Token)
			case "operation-action":
				_, err = h.db.Exec(`UPDATE managed_schema_operations SET action='successor_reserve' WHERE operation_id=?`, token.OperationID)
			case "operation-token":
				_, err = h.db.Exec(`UPDATE managed_schema_operations SET token='foreign' WHERE operation_id=?`, token.OperationID)
			case "operation-scope":
				_, err = h.db.Exec(`UPDATE managed_schema_operations SET scope_key='foreign' WHERE operation_id=?`, token.OperationID)
			case "operation-workspace":
				_, err = h.db.Exec(`UPDATE managed_schema_operations SET workspace_id='foreign' WHERE operation_id=?`, token.OperationID)
			case "operation-proof":
				_, err = h.db.Exec(`UPDATE managed_schema_operations SET status='succeeded',verified_digest=? WHERE operation_id=?`, strings.Repeat("z", 64), token.OperationID)
			}
			require.NoError(t, err)
			before := lineageReadSnapshot(t, h)
			status, err := h.service.SuccessorOperationStatus(t.Context(), id, statusRequest(stage))
			require.Error(t, err)
			require.Nil(t, status)
			require.Equal(t, before, lineageReadSnapshot(t, h))
		})
	}
}

func TestSuccessorOperationStatusValidatesPriorSchemaProofWithoutLatestSubstitution(t *testing.T) {
	for _, foreignPrior := range []bool{false, true} {
		t.Run(map[bool]string{false: "exact-prior", true: "foreign-prior"}[foreignPrior], func(t *testing.T) {
			h, id, stage, _ := successorSchemaFixture(t)
			first, repo, ledger := statusPreview(t, h, id, stage)
			op, _, err := ledger.ClaimSchemaApply(t.Context(), first)
			require.NoError(t, err)
			digest := strings.Repeat("a", 64)
			_, err = repo.FinishSchemaOperation(t.Context(), op.OperationID, op.Owner, recordingplan.SchemaOperationResult{Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: digest})
			require.NoError(t, err)
			confirmed, err := h.service.ConfirmSuccessorSchema(t.Context(), id, stage, first.OperationID, digest, &StorageEvidence{Layout: *first.GroupLayout, CreatedTable: true})
			require.NoError(t, err)
			stage.Current.WorkspaceRevision = confirmed.WorkspaceRevision
			scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
			require.NoError(t, err)
			if foreignPrior {
				scope.SchemaRevision = "foreign-schema-operation"
			}
			_, err = ledger.PrepareSchemaPreview(t.Context(), scope, func(context.Context, string) (recordingplan.TargetTableInspection, error) {
				return recordingplan.TargetTableInspection{Status: "missing"}, nil
			})
			require.NoError(t, err)
			status, err := h.service.SuccessorOperationStatus(t.Context(), id, statusRequest(stage))
			if foreignPrior {
				require.ErrorIs(t, err, ErrWriteGroupLineageCorrupt)
				require.Nil(t, status)
			} else {
				require.NoError(t, err)
				require.Len(t, status.SchemaOperations, 2)
				require.Equal(t, "succeeded", status.SchemaOperations[0].Status)
				require.Equal(t, "preview_not_submitted", status.SchemaOperations[1].Status)
			}
		})
	}
}

func TestSuccessorOperationStatusTimescalePreservesAuthenticatedExtensionSchema(t *testing.T) {
	group := &WriteGroup{ID: "G2", WorkspaceID: "W", Revision: "R", LayoutVersion: WriteGroupLayoutManagedLongV1, StorageMode: WriteGroupStorageModeTimescale,
		Members: []WriteGroupMember{{TagID: "T"}}, Destination: WriteGroupDestination{ConnectorID: "C", ConnectorRevision: "CR", Database: "db", TableSchema: "public", StorageStrategy: WriteGroupStorageStrategyManaged}}
	table, _, ok := LongStorageIdentity(group)
	require.True(t, ok)
	group.Destination.TableName = table
	binding := &successorReserveSchemaBinding{WorkspaceID: "W", Dialect: "postgres", SourceDigest: "source", TagTypes: map[string]schema.DataType{"T": schema.DataTypeFloat32}}
	layout, err := managedSchemaLayout(group, binding.TagTypes, binding.Dialect)
	require.NoError(t, err)
	layout.Partition.ExtensionSchema = "timescale_ext"
	token := &recordingplan.SchemaPreviewToken{WorkspaceID: "W", PlanID: "G2", PlanRevision: "R", ConnectorID: "C", ConnectorRevision: "CR", Dialect: "postgres", Database: "db", Schema: "public", Action: recordingplan.SchemaApplyAction, SourceDigest: "source", GroupLayout: layout}
	reserved := &WriteGroupSuccessorPreparation{Group: group}
	require.NoError(t, successorStatusTokenScope(token, reserved, binding), "extension installation namespace is frozen token evidence, not a different storage layout")
	for _, bad := range []string{"", "private; DROP TABLE x"} {
		token.GroupLayout.Partition.ExtensionSchema = bad
		require.Error(t, successorStatusTokenScope(token, reserved, binding))
	}
}

func TestSuccessorOperationStatusFrozenTypesSurviveLiveDrift(t *testing.T) {
	h, id, stage, _ := successorSchemaFixture(t)
	_, _, _ = statusPreview(t, h, id, stage)
	before, err := h.service.SuccessorOperationStatus(t.Context(), id, statusRequest(stage))
	require.NoError(t, err)
	_, err = h.db.Exec(`UPDATE tags SET data_type='int16' WHERE id=?`, stage.Envelope.Request.Group.Members[0].TagID)
	require.NoError(t, err)
	_, err = h.service.SuccessorSchemaScope(t.Context(), id, stage)
	require.Error(t, err, "changed live source types must not issue a differently-bound preview for the same root")
	after, err := h.service.SuccessorOperationStatus(t.Context(), id, statusRequest(stage))
	require.NoError(t, err)
	require.Equal(t, before, after, "historical schema readback uses frozen type evidence")
}
