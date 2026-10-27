package recordingplan

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMySQLStaleRunningPreflightsWholeScopeBeforeDDL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mysql-stale-foreign-ledger.db")
	ledger := openMySQLApplyLedger(t, path)
	applyRecordingDDL(t, ledger)
	baseRepo := NewSQLRepository(ledger)
	repo := &mysqlApplyFinishOnceRepository{SQLRepository: baseRepo, fail: true}
	fixture := newMySQLDDLApplyFixture(t, repo)
	group := cloneTargetTableInspection(fixture.groupExact)
	group.Indexes = nil
	fixture.states[fixture.layout.TableName] = group
	fixture.states[managedEffectReceiptTable] = cloneTargetTableInspection(fixture.receiptExact)
	token := *fixture.token
	token.Token, token.OperationID = "tok-preflight-scope", "op-preflight-scope"
	token.Statements = []string{fixture.statements[1]}
	var err error
	token.Digest, err = previewDigest(&token)
	require.NoError(t, err)
	require.NoError(t, repo.SavePreviewToken(t.Context(), &token))
	fixture.token = &token
	fixture.failStatement, fixture.failErr = fixture.token.Statements[0], ErrTargetPermissionDenied
	service := NewService(repo)

	_, _, err = service.ApplySchemaPreview(t.Context(), fixture.token, fixture.target())
	require.ErrorIs(t, err, ErrSchemaOperationUnacknowledged)
	require.Len(t, fixture.calls, 1)
	stored, err := baseRepo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
	require.NoError(t, err)
	require.Equal(t, SchemaOperationRunning, stored.Status)
	manifest := mysqlApplyTestDetail(t, stored.Detail)
	require.Equal(t, []string{mysqlSchemaStepFailed}, []string{manifest.Steps[0].Status})

	foreignReceipt := cloneTargetTableInspection(fixture.receiptExact)
	foreignReceipt.MySQL.Engine = "MyISAM"
	fixture.states[managedEffectReceiptTable] = foreignReceipt
	require.NoError(t, ledger.Close())
	reopened := openMySQLApplyLedger(t, path)
	_, err = reopened.ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE operation_id = ?`,
		time.Now().UTC().Add(-2*SchemaOperationLease), fixture.token.OperationID)
	require.NoError(t, err)
	restarted := NewService(NewSQLRepository(reopened))

	resolved, outcome, err := restarted.ApplySchemaPreview(t.Context(), fixture.token, fixture.target())

	require.NoError(t, err)
	require.Equal(t, ClaimAcquired, outcome)
	require.Equal(t, SchemaOperationFailed, resolved.Status)
	require.Equal(t, SchemaReasonVerificationGap, resolved.Reason)
	require.Len(t, fixture.calls, 1, "a foreign receipt must prevent retrying the otherwise-missing index")
}

func TestMySQLApplyForbiddenCatalogNeverMeansMissing(t *testing.T) {
	repo := NewMemoryRepository()
	fixture := newMySQLDDLApplyFixture(t, repo)
	forbidden := fixture.states[fixture.layout.TableName]
	forbidden.Status = inspectionForbidden
	target := fixture.target()
	groupInspections := 0
	target.Inspect = func(ctx context.Context, table string) (TargetTableInspection, error) {
		if table == fixture.layout.TableName {
			groupInspections++
			if groupInspections == 3 {
				return forbidden, nil
			}
		}
		return fixture.inspect(ctx, table)
	}
	service := NewService(repo)

	_, _, err := service.ApplySchemaPreview(t.Context(), fixture.token, target)

	require.ErrorIs(t, err, ErrTargetInspectionUnconfirmed)
	require.ErrorIs(t, err, ErrTargetPermissionDenied)
	require.Empty(t, fixture.calls)
	_, err = repo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
	require.ErrorIs(t, err, ErrSchemaOperationNotFound)
}

func TestMySQLForeignObjectRejected(t *testing.T) {
	repo := NewMemoryRepository()
	fixture := newMySQLDDLApplyFixture(t, repo)
	foreign := cloneTargetTableInspection(fixture.groupExact)
	foreign.MySQL.Engine = "MyISAM"
	target := fixture.target()
	groupInspections := 0
	target.Inspect = func(ctx context.Context, table string) (TargetTableInspection, error) {
		if table == fixture.layout.TableName {
			groupInspections++
			if groupInspections == 3 {
				return foreign, nil
			}
		}
		return fixture.inspect(ctx, table)
	}
	service := NewService(repo)

	_, _, err := service.ApplySchemaPreview(t.Context(), fixture.token, target)

	require.ErrorIs(t, err, ErrIncompatibleExistingTable)
	require.Empty(t, fixture.calls, "a foreign same-name object must never receive DDL")
	_, err = repo.GetSchemaOperation(t.Context(), fixture.token.WorkspaceID, fixture.token.OperationID)
	require.ErrorIs(t, err, ErrSchemaOperationNotFound)
}
