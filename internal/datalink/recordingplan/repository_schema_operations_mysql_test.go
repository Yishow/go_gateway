package recordingplan

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMySQLSchemaApplyReclaimRequiresMatchingTerminalAndFreeScope(t *testing.T) {
	for name, base := range map[string]Repository{
		"sql":    NewSQLRepository(recordingRepositoryDB(t)),
		"memory": NewMemoryRepository(),
	} {
		t.Run(name, func(t *testing.T) {
			resume, ok := base.(mysqlSchemaOperationResumeRepository)
			require.True(t, ok)
			ctx := t.Context()
			detail := `{"version":1,"token_digest":"digest-1","scope_key":"scope-a","steps":[]}`
			original := ledgerOperation("op-mysql", "tok-mysql", "scope-a", "owner-old")
			original.Detail = detail
			claimed, outcome, err := base.ClaimSchemaOperation(ctx, original)
			require.NoError(t, err)
			require.Equal(t, ClaimAcquired, outcome)
			require.Equal(t, "owner-old", claimed.Owner)
			partial, err := base.FinishSchemaOperation(ctx, original.OperationID, original.Owner, SchemaOperationResult{
				Status: SchemaOperationPartial, Detail: detail,
			})
			require.NoError(t, err)

			other := ledgerOperation("op-other", "tok-other", original.ScopeKey, "owner-other")
			_, outcome, err = base.ClaimSchemaOperation(ctx, other)
			require.NoError(t, err)
			require.Equal(t, ClaimAcquired, outcome)
			busy, outcome, err := resume.ReclaimMySQLSchemaApplyOperation(
				ctx, original.OperationID, original.WorkspaceID, original.Token, original.ScopeKey, detail, "owner-retry",
			)
			require.NoError(t, err)
			require.Equal(t, ClaimScopeBusy, outcome)
			require.Equal(t, other.OperationID, busy.OperationID)
			unchanged, err := base.GetSchemaOperation(ctx, original.WorkspaceID, original.OperationID)
			require.NoError(t, err)
			require.Equal(t, SchemaOperationPartial, unchanged.Status)
			require.Equal(t, partial.Owner, unchanged.Owner)

			_, err = base.FinishSchemaOperation(ctx, other.OperationID, other.Owner, SchemaOperationResult{Status: SchemaOperationSucceeded})
			require.NoError(t, err)
			reclaimed, outcome, err := resume.ReclaimMySQLSchemaApplyOperation(
				ctx, original.OperationID, original.WorkspaceID, original.Token, original.ScopeKey, detail, "owner-retry",
			)
			require.NoError(t, err)
			require.Equal(t, ClaimAcquired, outcome)
			require.Equal(t, SchemaOperationRunning, reclaimed.Status)
			require.Equal(t, "owner-retry", reclaimed.Owner)
			require.Equal(t, detail, reclaimed.Detail)
			if progress, ok := base.(operationExecutionRepository); ok {
				require.ErrorIs(t, progress.SaveSchemaOperationProgress(ctx, original.OperationID, original.Owner, detail), ErrSchemaOperationNotOwned)
			}
			require.ErrorIs(t, func() error {
				_, err := base.FinishSchemaOperation(ctx, original.OperationID, original.Owner, SchemaOperationResult{Status: SchemaOperationFailed})
				return err
			}(), ErrSchemaOperationNotOwned)

			_, err = base.FinishSchemaOperation(ctx, original.OperationID, "owner-retry", SchemaOperationResult{
				Status: SchemaOperationSucceeded, Detail: detail,
			})
			require.NoError(t, err)
			terminal, outcome, err := resume.ReclaimMySQLSchemaApplyOperation(
				ctx, original.OperationID, original.WorkspaceID, original.Token, original.ScopeKey, detail, "owner-again",
			)
			require.NoError(t, err)
			require.Equal(t, ClaimCompleted, outcome)
			require.Equal(t, SchemaOperationSucceeded, terminal.Status)
			require.NotEqual(t, "owner-again", terminal.Owner)
		})
	}
}

func TestMySQLSchemaApplyReclaimRejectsOtherActionsAndScopeDrift(t *testing.T) {
	for name, base := range map[string]Repository{
		"sql":    NewSQLRepository(recordingRepositoryDB(t)),
		"memory": NewMemoryRepository(),
	} {
		t.Run(name, func(t *testing.T) {
			resume := base.(mysqlSchemaOperationResumeRepository)
			ctx := t.Context()
			operation := ledgerOperation("op-other-action", "tok-other-action", "scope-action", "owner-old")
			operation.Action = "test_write"
			operation.Detail = `{"marker":"not-a-schema-apply"}`
			_, outcome, err := base.ClaimSchemaOperation(ctx, operation)
			require.NoError(t, err)
			require.Equal(t, ClaimAcquired, outcome)
			_, err = base.FinishSchemaOperation(ctx, operation.OperationID, operation.Owner, SchemaOperationResult{
				Status: SchemaOperationPartial, Detail: operation.Detail,
			})
			require.NoError(t, err)

			_, _, err = resume.ReclaimMySQLSchemaApplyOperation(
				ctx, operation.OperationID, operation.WorkspaceID, operation.Token, operation.ScopeKey, operation.Detail, "owner-new",
			)
			require.ErrorIs(t, err, ErrSchemaOperationMismatch)
			_, _, err = resume.ReclaimMySQLSchemaApplyOperation(
				ctx, operation.OperationID, operation.WorkspaceID, operation.Token, "other-scope", operation.Detail, "owner-new",
			)
			require.ErrorIs(t, err, ErrSchemaOperationMismatch)
			stored, err := base.GetSchemaOperation(ctx, operation.WorkspaceID, operation.OperationID)
			require.NoError(t, err)
			require.Equal(t, SchemaOperationPartial, stored.Status)
			require.Equal(t, "owner-old", stored.Owner)

		})
	}
}

func TestMySQLSchemaApplyBaselineMatchesOnlyVerifiedFamilyVersions(t *testing.T) {
	cases := []struct {
		name    string
		options MySQLGroupSchemaOptions
		want    bool
	}{
		{name: "mysql 8.4 patch", options: mysqlTestOptions(), want: true},
		{name: "mysql 8.3", options: MySQLGroupSchemaOptions{
			ServerFamily: mysqlFamilyMySQL, ServerVersion: "8.3.12", CapabilityRevision: "cap-1",
			Engine: mysqlEngineInnoDB, Collation: mysqlCollationMySQL,
		}},
		{name: "MariaDB 11.4 patch", options: mariadbTestOptions(), want: true},
		{name: "MariaDB 10.11", options: MySQLGroupSchemaOptions{
			ServerFamily: mysqlFamilyMariaDB, ServerVersion: "10.11.8-MariaDB", CapabilityRevision: "cap-1",
			Engine: mysqlEngineInnoDB, Collation: mysqlCollationMariaDB,
		}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			layout, err := ManagedMySQLLongGroupLayout("gw_baseline", "_gw_owner_0123456789abcdef01234567", testCase.options)
			require.NoError(t, err)
			require.Equal(t, testCase.want, mysqlSchemaApplyBaselineSupported(&layout))
		})
	}
}

func TestMySQLSchemaApplyReclaimNeverIgnoresRepositoryFailure(t *testing.T) {
	var resume mysqlSchemaOperationResumeRepository = NewSQLRepository(recordingRepositoryDB(t))
	_, _, err := resume.ReclaimMySQLSchemaApplyOperation(
		t.Context(), "missing", "ws-1", "token", "scope", `{"version":1}`, "owner-new",
	)
	require.True(t, errors.Is(err, ErrSchemaOperationNotFound) || errors.Is(err, ErrSchemaOperationMismatch))
}
