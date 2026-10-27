package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerReceiptRepairMarkerIsIncompatible(t *testing.T) {
	scope := sqlServerSchemaScopeForTest(t)
	group := sqlServerCompleteGroupForTest(scope, 11)
	receipt := sqlServerReceiptStateForTest()
	receipt.SQLServer.RepairOperationID = "repair-foreign-receipt"

	_, err := VerifyManagedSQLServerCatalog(*scope.GroupLayout, scope.Database, group, receipt)

	require.ErrorIs(t, err, ErrIncompatibleExistingTable, "a receipt proof refused by delivery must also be refused before preparation/intake")
}
