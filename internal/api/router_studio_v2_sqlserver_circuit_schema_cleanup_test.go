package api

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func sqlServerCircuitPreviewCreatesGroup(token recordingplan.SchemaPreviewToken) bool {
	if token.GroupLayout == nil {
		return false
	}
	for _, table := range token.Tables {
		if table.Name == token.GroupLayout.TableName && table.Action == recordingplan.SchemaTableActionCreate {
			return true
		}
	}
	return false
}

// Only a table created by this exact confirmed preview can be removed. Receipt
// storage is shared and is retained; effect cleanup separately proves its rows.
func registerSQLServerCircuitCreatedTableCleanup(t *testing.T, f managedGroupFixture, targets *dbtarget.ConnectorService, token recordingplan.SchemaPreviewToken) {
	t.Helper()
	if !sqlServerCircuitPreviewCreatesGroup(token) || token.GroupLayout.SQLServer.NativeMonthly != nil {
		return // Native cleanup is owned by the partition integration fixture.
	}
	capability, err := targets.ProbeSQLServerRecordingAtRevision(t.Context(), token.ConnectorID, token.ConnectorRevision)
	require.NoError(t, err)
	group, err := dbtarget.InspectOpenSQLServerManagedTable(t.Context(), f.target, capability, token.GroupLayout.TableName)
	require.NoError(t, err)
	receipt, err := dbtarget.InspectOpenSQLServerManagedTable(t.Context(), f.target, capability, dbtarget.EffectReceiptTable)
	require.NoError(t, err)
	_, err = recordingplan.VerifyManagedSQLServerCatalog(*token.GroupLayout, token.Database,
		sqlServerCircuitTargetInspection(group), sqlServerCircuitTargetInspection(receipt))
	require.NoError(t, err, "cleanup ownership must be verified from the actual post-confirmation catalog")
	objectID := group.Evidence.ObjectID
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		current, err := dbtarget.InspectOpenSQLServerManagedTable(ctx, f.target, capability, token.GroupLayout.TableName)
		require.NoError(t, err)
		if current.Inspection.Status == dbtarget.TableInspectionMissing {
			return // A deliberate, separately owned missing-table case removed it.
		}
		require.Equal(t, objectID, current.Evidence.ObjectID, "a replacement object is never test-owned")
		receipts, err := dbtarget.InspectOpenSQLServerManagedTable(ctx, f.target, capability, dbtarget.EffectReceiptTable)
		require.NoError(t, err)
		_, err = recordingplan.VerifyManagedSQLServerCatalog(*token.GroupLayout, token.Database,
			sqlServerCircuitTargetInspection(current), sqlServerCircuitTargetInspection(receipts))
		require.NoError(t, err)
		var remaining int
		require.NoError(t, f.target.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+sqlServerQualified(token.Schema, token.GroupLayout.TableName)).Scan(&remaining))
		require.Zero(t, remaining, "leave nonempty storage intact; effect cleanup must finish first")
		_, err = f.target.ExecContext(ctx, "IF OBJECT_ID(@p1,N'U')=@p2 AND NOT EXISTS (SELECT 1 FROM "+
			sqlServerQualified(token.Schema, token.GroupLayout.TableName)+") DROP TABLE "+sqlServerQualified(token.Schema, token.GroupLayout.TableName)+
			" ELSE THROW 50000, 'fixture table identity or contents changed', 1", token.Schema+"."+token.GroupLayout.TableName, objectID)
		require.NoError(t, err)
	})
}
