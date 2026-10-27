package dbtarget

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerManagedSchemaExecutionKeepsUnverifiedGatesClosed(t *testing.T) {
	for _, native := range []bool{false, true} {
		t.Run(map[bool]string{false: "plain", true: "native"}[native], func(t *testing.T) {
			script := sqlServerManagedCatalogFixture()
			if native {
				script.metadata.alterAnyDataspace = false
			} else {
				script.metadata.createTable = false
			}
			inspector, connector := sqlServerManagedInspectionTestSetup(t, script, "sql-managed-schema-gate")
			result, err := inspector.ExecuteSQLServerManagedStatementsAtRevision(t.Context(), connector.ID,
				connector.IdentityRevision, "fixture_db", "dbo", sqlServerManagedExpectedCapabilityRevision(connector),
				native, []string{"CREATE TABLE [dbo].[fixture_only] ([id] BIGINT)"})
			require.Error(t, err, "unverified permissions cannot release managed schema operations")
			require.NotNil(t, result)
			require.True(t, result.RolledBack, "the capability refusal occurs before any transaction starts")
			require.Zero(t, result.Committed)
			require.Zero(t, script.beginCalls)
			require.Zero(t, script.execCalls)
			require.NotEmpty(t, script.queries, "the decision uses actual read-only metadata")
			require.True(t, script.closed)
		})
	}
}

func TestSQLServerManagedSchemaExecutionRejectsStaleScopeBeforeMetadata(t *testing.T) {
	for _, scope := range []struct {
		name, revision, database, schema string
	}{
		{name: "stale revision", revision: "stale", database: "fixture_db", schema: "dbo"},
		{name: "changed database", database: "other_database", schema: "dbo"},
		{name: "changed schema", database: "fixture_db", schema: "foreign_schema"},
	} {
		t.Run(scope.name, func(t *testing.T) {
			script := sqlServerManagedCatalogFixture()
			inspector, connector := sqlServerManagedInspectionTestSetup(t, script, "sql-managed-schema-scope")
			revision := connector.IdentityRevision
			if scope.revision != "" {
				revision = scope.revision
			}
			_, err := inspector.ExecuteSQLServerManagedStatementsAtRevision(t.Context(), connector.ID,
				revision, scope.database, scope.schema, sqlServerManagedExpectedCapabilityRevision(connector),
				false, []string{"CREATE TABLE [dbo].[fixture_only] ([id] BIGINT)"})
			require.Error(t, err)
			require.Empty(t, script.queries)
			require.Zero(t, script.beginCalls)
			require.Zero(t, script.execCalls)
		})
	}
}
