package dbtarget

import (
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestSQLServerManagedInspectionDoesNotCallMissingSchemaATable(t *testing.T) {
	script := sqlServerManagedCatalogFixture()
	script.metadata.visibleSchema = ""
	script.metadata.schemaViewDefinition = false
	inspector, connector := sqlServerManagedInspectionTestSetup(t, script, "sqlserver-managed-missing-schema")

	result, err := inspector.InspectSQLServerManagedTableAtRevision(
		t.Context(), connector.ID, connector.IdentityRevision, "fixture_db", "dbo",
		sqlServerManagedExpectedCapabilityRevision(connector), "circuit_long",
	)
	require.ErrorIs(t, err, ErrSQLServerManagedTableUnconfirmed)
	require.NotNil(t, result)
	require.Equal(t, TableInspectionFailed, result.Inspection.Status)
	require.Equal(t, "schema_missing", result.Inspection.Reason)
	require.False(t, result.Evidence.MetadataVerified)
	require.False(t, sqlServerManagedHasQuery(script, "sys.objects"), "an absent schema is not evidence of a missing table")
	require.Zero(t, script.execCalls, "inspection must not create the missing schema")
}

func TestSQLServerManagedInspectionRejectsApproximateJSONChecks(t *testing.T) {
	script := sqlServerManagedCatalogFixture()
	script.checkRows = [][]driver.Value{{
		"CK_circuit_long_provenance", "ISJSON([provenance]) = 1 OR 1 = 1", false, false, int64(0),
	}}
	inspector, connector := sqlServerManagedInspectionTestSetup(t, script, "sqlserver-managed-json-check")

	result, err := inspector.InspectSQLServerManagedTableAtRevision(
		t.Context(), connector.ID, connector.IdentityRevision, "fixture_db", "dbo",
		sqlServerManagedExpectedCapabilityRevision(connector), "circuit_long",
	)
	require.NoError(t, err)
	require.False(t, result.Evidence.ProvenanceJSONChecked, "a predicate containing ISJSON is not the canonical constraint")
}

func TestSQLServerManagedInspectionRejectsNonCanonicalPrimaryKeyIndexes(t *testing.T) {
	for _, test := range []struct {
		name   string
		column int
		value  bool
	}{
		{name: "ignore duplicate key", column: 7, value: true},
		{name: "not unique", column: 3, value: false},
		{name: "filtered primary key", column: 5, value: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			script := sqlServerManagedCatalogFixture()
			for _, row := range script.indexRows {
				if row[0] == int64(1) {
					row[test.column] = test.value
				}
			}
			inspector, connector := sqlServerManagedInspectionTestSetup(t, script, "sqlserver-managed-pk-"+test.name)
			result, err := inspector.InspectSQLServerManagedTableAtRevision(
				t.Context(), connector.ID, connector.IdentityRevision, "fixture_db", "dbo",
				sqlServerManagedExpectedCapabilityRevision(connector), "circuit_long",
			)
			require.ErrorIs(t, err, ErrSQLServerManagedIndexUnusable)
			require.NotNil(t, result)
			require.False(t, result.Evidence.Plain)
		})
	}
}

func sqlServerManagedInspectionTestSetup(
	t *testing.T,
	script *sqlServerManagedCatalogScript,
	connectorID string,
) (*ManagedTableInspector, *schema.DatabaseConnector) {
	t.Helper()
	internalDB := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(internalDB)
	service := NewConnectorService(repo)
	connector := createSQLServerRecordingSavedConnector(
		t, repo, connectorID, schema.DatabaseConnectorKindSQLServer, sqlServerConnectorConfigForTest(), true,
	)
	session := openSQLServerManagedCatalogTestSession(t, script)
	originalSessionOpen := openSQLServerManagedSessionFunc
	originalManagerOpen := openExternalDBManagerFunc
	openSQLServerManagedSessionFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (sqlServerManagedSession, error) {
		return session, nil
	}
	openExternalDBManagerFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
		return nil, errors.New("unexpected generic open")
	}
	t.Cleanup(func() {
		openSQLServerManagedSessionFunc = originalSessionOpen
		openExternalDBManagerFunc = originalManagerOpen
	})
	return NewManagedTableInspector(service, internalDB), connector
}

func sqlServerManagedHasQuery(script *sqlServerManagedCatalogScript, fragment string) bool {
	for _, query := range script.queries {
		if strings.Contains(query.query, fragment) {
			return true
		}
	}
	return false
}
