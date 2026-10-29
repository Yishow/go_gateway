package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/recordingplan"
)

func TestSuccessorSQLServerUnitProofRejectsCorruptLedger(t *testing.T) {
	queries := map[string]string{
		"missing-operation": `DELETE FROM managed_schema_operations`,
		"unknown-operation": `UPDATE managed_schema_operations SET status='unknown'`,
		"operation-digest":  `UPDATE managed_schema_operations SET verified_digest='foreign'`,
		"operation-scope":   `UPDATE managed_schema_operations SET scope_key='foreign'`,
		"foreign-group":     `UPDATE managed_schema_preview_tokens SET plan_id='foreign'`,
		"foreign-workspace": `UPDATE managed_schema_preview_tokens SET workspace_id='foreign'`,
		"stale-connector":   `UPDATE managed_schema_preview_tokens SET connector_revision='stale'`,
		"source":            `UPDATE managed_schema_preview_tokens SET tables=json_set(tables,'$.source_digest','foreign')`,
		"statements":        `UPDATE managed_schema_preview_tokens SET statements='["DROP TABLE wrong"]'`,
		"owner":             `UPDATE managed_schema_preview_tokens SET tables=replace(tables,'_gw_owner_','_notowner_')`,
	}
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		for name, query := range queries {
			t.Run(string(mode)+"/"+name, func(t *testing.T) {
				h, id, req, p := successorSQLServerUnitConfirmed(t, mode)
				_, err := h.db.Exec(query)
				require.NoError(t, err)
				calls, catalogs := p.calls, p.catalogCalls
				before := lineageReadSnapshot(t, h)
				view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
				require.NoError(t, err)
				require.Contains(t, view.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "schema-confirmation-unverified", NextAction: "review_lifecycle"})
				require.Equal(t, before, lineageReadSnapshot(t, h))
				require.Equal(t, calls, p.calls)
				require.Equal(t, catalogs, p.catalogCalls)
			})
		}
	}
}

func TestSuccessorSQLServerUnitValidTokenStillBindsCurrentScope(t *testing.T) {
	for _, field := range []string{"workspace", "group", "group-revision", "connector", "connector-revision", "database", "schema", "reservation", "source"} {
		t.Run(field, func(t *testing.T) {
			h, id, req, _ := successorSQLServerUnitConfirmed(t, WriteGroupStorageModePlain)
			view, snap, err := h.service.readSuccessorActivation(t.Context(), id, req)
			require.NoError(t, err)
			require.Empty(t, view.Blockers)
			tx, err := h.db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			defer tx.Rollback()
			g := snap.reserved.Group
			_, err = recordingplan.SucceededSchemaTokenInTx(t.Context(), tx, g.WorkspaceID, g.Destination.SchemaRevision, g.Destination.SchemaDigest)
			require.NoError(t, err, "token and operation themselves remain valid")
			switch field {
			case "workspace":
				g.WorkspaceID = "other-workspace"
			case "group":
				g.ID = "other-group"
			case "group-revision":
				g.Revision = "other-revision"
			case "connector":
				g.Destination.ConnectorID = snap.source.Group.Destination.ConnectorID
			case "connector-revision":
				g.Destination.ConnectorRevision = "other-revision"
			case "database":
				g.Destination.Database = "other-database"
			case "schema":
				g.Destination.TableSchema = "other_schema"
			case "reservation":
				snap.reserved.IntentDigest = "other-reservation"
			case "source":
				g.Members[0].SourceRevision = "other-source"
			}
			require.Error(t, h.service.activationSchemaProofInTx(t.Context(), tx, snap))
		})
	}
}

func TestSuccessorSQLServerUnitRejectsBothModeProofExchanges(t *testing.T) {
	for _, mode := range []WriteGroupStorageMode{WriteGroupStorageModePlain, WriteGroupStorageModeNativeMonthly} {
		t.Run(string(mode), func(t *testing.T) {
			h, id, stage, p := successorSQLServerUnitReserved(t, mode)
			scope, err := h.service.SuccessorSchemaScope(t.Context(), id, stage)
			require.NoError(t, err)
			options := *scope.GroupLayout.SQLServer
			if mode == WriteGroupStorageModeNativeMonthly {
				options.NativeMonthly = nil
			} else {
				function, scheme, err := recordingplan.SQLServerNativeMonthlyObjectNames(scope.GroupLayout.TableName)
				require.NoError(t, err)
				options.NativeMonthly, err = recordingplan.NewSQLServerNativeMonthlyPartition(h.service.clockNow(), function, scheme)
				require.NoError(t, err)
			}
			wrong, err := recordingplan.ManagedSQLServerLongGroupLayout(scope.GroupLayout.TableName, scope.GroupLayout.OwnerColumn, options)
			require.NoError(t, err)
			scope.GroupLayout = &wrong
			token, digest := successorSQLServerUnitSucceeded(t, h, scope, p)
			before := lineageReadSnapshot(t, h)
			_, err = h.service.ConfirmSuccessorSchema(t.Context(), id, stage, token.OperationID, digest, &StorageEvidence{Layout: *token.GroupLayout, CreatedTable: true})
			require.Error(t, err)
			require.Equal(t, before, lineageReadSnapshot(t, h))
			// Even a locally injected schema reference cannot promote the valid
			// opposite-mode token into activation authority.
			_, err = h.db.Exec(`UPDATE write_groups SET destination_schema_revision=?,destination_schema_digest=? WHERE id=?`, token.OperationID, digest, stage.Current.GroupID)
			require.NoError(t, err)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, WriteGroupSuccessorActivationRequest{Stage: stage, Boundary: stage.Envelope.Request.Boundary})
			require.NoError(t, err)
			require.Contains(t, view.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "schema-confirmation-unverified", NextAction: "review_lifecycle"})
		})
	}
}

func TestSuccessorSQLServerUnitNativeProofBindsStorageConfirmation(t *testing.T) {
	for _, mutation := range []string{"missing", "corrupt-digest", "different-partitions"} {
		t.Run(mutation, func(t *testing.T) {
			h, id, req, _ := successorSQLServerUnitConfirmed(t, WriteGroupStorageModeNativeMonthly)
			group, err := h.service.Get(t.Context(), req.Stage.Current.GroupID)
			require.NoError(t, err)
			connector := group.Group.Destination.ConnectorID
			switch mutation {
			case "missing":
				_, err = h.db.Exec(`UPDATE write_group_storage_scopes SET confirmed_layout='' WHERE connector_id=?`, connector)
			case "corrupt-digest":
				_, err = h.db.Exec(`UPDATE write_group_storage_scopes SET layout_digest='corrupt' WHERE connector_id=?`, connector)
			case "different-partitions":
				lineage, e := recordingplan.ReadStorageLineage(t.Context(), h.db, StorageScopeOf(group.Group))
				require.NoError(t, e)
				lineage.Layout.SQLServer.NativeMonthly, e = recordingplan.NewSQLServerNativeMonthlyPartition(h.service.clockNow().AddDate(0, 1, 0), lineage.Layout.SQLServer.NativeMonthly.FunctionName, lineage.Layout.SQLServer.NativeMonthly.SchemeName)
				require.NoError(t, e)
				raw, e := json.Marshal(lineage.Layout)
				require.NoError(t, e)
				sum := sha256.Sum256(raw)
				_, err = h.db.Exec(`UPDATE write_group_storage_scopes SET confirmed_layout=?,layout_digest=? WHERE connector_id=?`, string(raw), hex.EncodeToString(sum[:]), connector)
			}
			require.NoError(t, err)
			view, err := h.service.SuccessorActivationReview(t.Context(), id, req)
			require.NoError(t, err)
			require.Contains(t, view.Blockers, WriteGroupSuccessorBlocker{ReasonCode: "schema-confirmation-unverified", NextAction: "review_lifecycle"})
		})
	}
}
