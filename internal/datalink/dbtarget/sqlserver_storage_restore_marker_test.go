package dbtarget

import (
	"database/sql/driver"
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerRepairOperationMarkerJSONRemainsBackwardCompatible(t *testing.T) {
	evidence := SQLServerManagedTableEvidence{MetadataVerified: true}
	encoded, err := json.Marshal(evidence)
	require.NoError(t, err)
	require.NotContains(t, string(encoded), "RepairOperationID")

	evidence.RepairOperationID = "repair-operation-test"
	encoded, err = json.Marshal(evidence)
	require.NoError(t, err)
	require.Contains(t, string(encoded), `"RepairOperationID":"repair-operation-test"`)
	require.NotContains(t, string(encoded), `"repair_operation_id"`)
}

func TestSQLServerOpenCatalogProjectsBoundRepairOperationMarker(t *testing.T) {
	script := sqlServerManagedCatalogFixture()
	script.repairOperationID = "repair-operation-test"
	db, capability := sqlServerOpenManagedInspectionTestDB(t, script)

	inspection, err := InspectOpenSQLServerManagedTable(t.Context(), db, capability, "circuit_long")
	require.NoError(t, err)
	require.Equal(t, script.repairOperationID, inspection.Evidence.RepairOperationID)

	markerCallIndex := slices.IndexFunc(script.queries, func(call sqlServerManagedQueryCall) bool {
		return strings.Contains(call.query, "sys.extended_properties")
	})
	require.NotEqual(t, -1, markerCallIndex)
	markerCall := script.queries[markerCallIndex]
	require.Equal(t, []driver.NamedValue{
		{Ordinal: 1, Value: sqlServerManagedRepairOperationPropertyName},
		{Ordinal: 2, Value: "dbo"},
		{Ordinal: 3, Value: "circuit_long"},
	}, markerCall.args)
}

func TestSQLServerAbsentOrDifferentMarkerCannotAttestRepair(t *testing.T) {
	for _, test := range []struct {
		name   string
		marker string
	}{
		{name: "absent"},
		{name: "different operation", marker: "another-repair-operation"},
	} {
		t.Run(test.name, func(t *testing.T) {
			script := sqlServerManagedCatalogFixture()
			script.repairOperationID = test.marker
			db, capability := sqlServerOpenManagedInspectionTestDB(t, script)
			inspection, err := InspectOpenSQLServerManagedTable(t.Context(), db, capability, "circuit_long")
			require.NoError(t, err)
			require.Equal(t, test.marker, inspection.Evidence.RepairOperationID)

			facts, err := inspectSQLServerManagedCatalog(t.Context(), db, "dbo", "circuit_long")
			require.NoError(t, err)
			expected := *inspection
			expectedEvidence := *inspection.Evidence
			expectedEvidence.RepairOperationID = "requested-repair-operation"
			expected.Evidence = &expectedEvidence
			require.False(t, sqlServerManagedCatalogMatches(&expected, facts))
		})
	}
}

func TestSQLServerReceiptShapeRejectsRepairOperationMarker(t *testing.T) {
	receipt := sqlServerGroupInsertTestRequest().SQLServer.Receipt
	require.True(t, sqlServerManagedReceiptShapeMatches(receipt))

	marked := *receipt
	evidence := *receipt.Evidence
	evidence.RepairOperationID = "repair-operation-test"
	marked.Evidence = &evidence
	require.False(t, sqlServerManagedReceiptShapeMatches(&marked))
}
