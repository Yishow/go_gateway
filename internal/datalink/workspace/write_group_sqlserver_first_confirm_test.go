package workspace

import (
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func TestSQLServerFirstNativeConfirmationDoesNotRequirePriorConfirmation(t *testing.T) {
	h := newCircuitProfileHarness(t, ":memory:")
	initial, err := recordingplan.NewSQLServerNativeMonthlyPartition(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		"pf_gw_long_sqlserver_normalize", "ps_gw_long_sqlserver_normalize")
	require.NoError(t, err)
	layout := sqlServerNativeTestLayout(t, initial)
	scope := recordingplan.SchemaPreviewScope{
		WorkspaceID: "ws-sql", ConnectorID: "conn-sql", Database: "recordings", Schema: "dbo", GroupLayout: &layout,
	}
	actual := sqlServerNativeTestState(layout, initial)
	receipt := sqlServerReceiptTestState(layout)
	_, err = recordingplan.VerifyManagedSQLServerCatalog(layout, scope.Database, actual, receipt)
	require.NoError(t, err, "the just-created physical catalog satisfies the initial confirmed preview")
	result, err := h.service.NormalizeSQLServerPartitionInspection(t.Context(), scope, actual, receipt)
	require.NoError(t, err, "post-DDL verification must work before the first confirmation is recorded")
	require.Equal(t, actual, result, "first confirmation uses the actual initial catalog without projecting history")

	expanded, err := recordingplan.PlanSQLServerNativeMonthlyExtension(initial, time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	_, err = h.service.NormalizeSQLServerPartitionInspection(t.Context(), scope,
		sqlServerNativeTestState(layout, expanded), receipt)
	require.Error(t, err, "an unconfirmed scope cannot normalize an unrecorded extension")
}
