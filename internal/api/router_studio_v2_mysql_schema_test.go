package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestMySQLSchemaAPIRetainsOperationOnPartialConfirmation(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			f := newMySQLSchemaRouterFixture(t, family)
			group := ensureLong(t, f, "instantaneous", "second", 15, "")
			id := group["id"].(string)
			table := destinationOf(group)["table_name"].(string)
			var existing int
			require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME IN (?,?)", table, "gw_effect_receipts").Scan(&existing))
			require.Zero(t, existing, "fixture objects must be absent before this test takes ownership")
			var created []string
			cleanupMySQLSchemaTables(t, f.target, &created)
			token := previewLong(t, f, id)
			require.Equal(t, family, token.GroupLayout.MySQL.ServerFamily)
			require.NotEmpty(t, token.GroupLayout.MySQL.CapabilityRevision)
			require.Equal(t, f.namespace, token.Database)
			require.Equal(t, token.Database, token.Schema)
			require.Len(t, token.Statements, 3)

			// Lose the local progress acknowledgement immediately before the index.
			// The real table is already committed at this point; no SQL rollback can
			// remove it. The test then persists a partial outcome through the ledger
			// owner to isolate HTTP replay routing from the service fault tests.
			_, err := f.db.ExecContext(t.Context(), `CREATE TRIGGER fail_mysql_index_progress BEFORE UPDATE OF detail ON managed_schema_operations
				WHEN json_extract(NEW.detail, '$.steps[1].status') = 'started'
				BEGIN SELECT RAISE(ABORT, 'injected local progress unavailable'); END`)
			require.NoError(t, err)
			request := scopeBody(t, f, id)
			request["token"], request["operation_id"] = token.Token, token.OperationID
			first := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/schema-apply", request)
			var tableExists int
			require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME=?", table).Scan(&tableExists))
			if tableExists == 1 {
				created = append(created, table)
			}
			require.Equal(t, http.StatusServiceUnavailable, first.Code, first.Body.String())
			require.Equal(t, 1, tableExists)
			_, err = f.db.ExecContext(t.Context(), `DROP TRIGGER fail_mysql_index_progress`)
			require.NoError(t, err)
			repo := recordingplan.NewSQLRepository(f.db)
			stored, err := repo.GetSchemaOperation(t.Context(), f.record.ID, token.OperationID)
			require.NoError(t, err)
			require.Equal(t, recordingplan.SchemaOperationRunning, stored.Status)
			require.Contains(t, stored.Detail, `"verified"`)
			_, err = repo.FinishSchemaOperation(t.Context(), stored.OperationID, stored.Owner, recordingplan.SchemaOperationResult{
				Status: recordingplan.SchemaOperationPartial, ExecutedStatements: 1, Detail: stored.Detail,
				Reason: recordingplan.SchemaReasonPartiallyApplied,
			})
			require.NoError(t, err)
			result := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+id+"/schema-apply", request)
			var receiptExists int
			require.NoError(t, f.target.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='gw_effect_receipts'").Scan(&receiptExists))
			if receiptExists == 1 {
				created = append(created, "gw_effect_receipts")
			}
			require.Equal(t, http.StatusOK, result.Code, result.Body.String())
			var body struct {
				Data recordingplan.SchemaOperation `json:"data"`
			}
			require.NoError(t, json.Unmarshal(result.Body.Bytes(), &body))
			require.Equal(t, recordingplan.SchemaOperationSucceeded, body.Data.Status)
			require.Equal(t, token.OperationID, body.Data.OperationID)
			require.Equal(t, 1, receiptExists)
			saved, err := f.groups.Get(t.Context(), id)
			require.NoError(t, err)
			lineage, err := recordingplan.ReadStorageLineage(t.Context(), f.db, workspace.StorageScopeOf(saved.Group))
			require.NoError(t, err)
			require.True(t, lineage.Confirmed())
			encoded, err := json.Marshal(token.GroupLayout)
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf("%x", sha256.Sum256(encoded)), lineage.LayoutDigest)
			require.Equal(t, 1, lineage.Generation)
			require.False(t, lineage.HistoryMayBeLost)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, id, token).Status)
			require.Empty(t, previewLong(t, f, id).Statements)
			seedProfileSource(t, f.writeGroupRouterFixture, "energy")
			energy := ensureLong(t, f, "cumulative", "minute", 10, "energy")
			require.Equal(t, table, destinationOf(energy)["table_name"], "both profiles retain the same verified MySQL storage scope")
			require.Empty(t, destinationOf(energy)["schema_revision"], "another group never inherits the first group's confirmation")
			energyToken := previewLong(t, f, energy["id"].(string))
			require.Empty(t, energyToken.Statements)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, energy["id"].(string), energyToken).Status)
		})
	}
}
