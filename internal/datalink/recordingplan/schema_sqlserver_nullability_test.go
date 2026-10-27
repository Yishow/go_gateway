package recordingplan

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLServerDDLDeclaresNullableSlots(t *testing.T) {
	layout, err := ManagedSQLServerLongGroupLayout("sqlserver_long", "_gw_owner_0123456789abcdef01234567", SQLServerGroupSchemaOptions{Schema: "dbo"})
	require.NoError(t, err)
	ddl, err := GenerateGroupSchemaDDL("sqlserver", layout)
	require.NoError(t, err)
	// Omitted nullability depends on SQL Server session/database settings.
	// Baseline, bad and missing readings must work under either default.
	require.Contains(t, ddl[0], "[value_decimal] DECIMAL(38,18) NULL")
	require.Contains(t, ddl[0], "[delta_value] DECIMAL(38,18) NULL")
	require.Contains(t, ddl[0], "[provenance] NVARCHAR(MAX) COLLATE Latin1_General_100_BIN2 NOT NULL CHECK (ISJSON([provenance]) = 1)")
}
