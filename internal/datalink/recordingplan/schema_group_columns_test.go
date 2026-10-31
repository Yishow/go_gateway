package recordingplan

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFrozenManagedGroupColumnsDoesNotInventIndexEvidence(t *testing.T) {
	layout, err := ManagedLongGroupLayout("sqlite", "long_rows", "_gw_owner_0123456789abcdef01234567")
	require.NoError(t, err)
	state := TargetTableInspection{Status: "exists", ColumnTypes: map[string]string{}, ColumnNullable: map[string]bool{}, ColumnPrimaryKey: map[string]bool{}}
	for _, c := range layout.Columns {
		state.Columns = append(state.Columns, c.Name)
		state.ColumnTypes[c.Name] = c.SQLType
		state.ColumnNullable[c.Name] = c.Nullable
		state.ColumnPrimaryKey[c.Name] = c.PrimaryKey
	}
	require.NoError(t, ValidateManagedGroupColumns(layout, state))
	require.Error(t, ValidateManagedGroupTableInspection(layout, state), "columns-only must not manufacture missing index inspection")
	state.ColumnNullable[layout.OwnerColumn] = true
	require.Error(t, ValidateManagedGroupColumns(layout, state))
	state.ColumnNullable[layout.OwnerColumn] = false
	state.Columns = state.Columns[1:]
	require.Error(t, ValidateManagedGroupColumns(layout, state))
}
