package api

import (
	"database/sql"
	"net/http"
	"testing"

	"go-gateway/internal/datalink/device"

	"github.com/stretchr/testify/require"
)

// The recording home bootstraps the workspace before its compact source read.
// Neither the initial read nor repeated observations may refresh device rows.
func TestNewRouter_StudioV2WorkspaceBootstrapDoesNotWriteDeviceReadiness(t *testing.T) {
	for _, cached := range []bool{false, true} {
		name := "uncached diagnostics"
		if cached {
			name = "cached diagnostics"
		}
		t.Run(name, func(t *testing.T) {
			f := newWriteGroupRouterFixture(t)
			_, err := f.db.ExecContext(t.Context(), `UPDATE devices SET description='' WHERE id='device-A'`)
			require.NoError(t, err)
			devices := device.NewService(device.NewSQLRepository(f.db), nil)
			f.workspace.WithReadinessServices(devices, nil, nil, nil)
			if cached {
				_, err = devices.CheckReadiness(t.Context(), "device-A")
				require.NoError(t, err)
			}
			router := NewRouter(&DatalinkServices{Workspace: f.workspace, Device: devices})
			before := workspaceBootstrapDeviceRows(t, f.db)
			var changesBefore int64
			require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT total_changes()`).Scan(&changesBefore))
			for range 2 {
				response := performJSONRequest(t, router, http.MethodGet, "/api/v1/datalink/studio-v2/workspace", nil)
				require.Equal(t, http.StatusOK, response.Code, response.Body.String())
				require.Contains(t, response.Body.String(), `"code":"device-connect-required"`)
				require.NotContains(t, response.Body.String(), "readiness-unavailable")
				// SELECT * includes readiness_status, updated_at and future columns.
				require.Equal(t, before, workspaceBootstrapDeviceRows(t, f.db), "workspace GET changed device configuration")
				var changesAfter int64
				require.NoError(t, f.db.QueryRowContext(t.Context(), `SELECT total_changes()`).Scan(&changesAfter))
				require.Equal(t, changesBefore, changesAfter, "even an identical same-second UPDATE is forbidden")
			}
		})
	}
}

func workspaceBootstrapDeviceRows(t *testing.T, db *sql.DB) [][]any {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), `SELECT * FROM devices ORDER BY id`)
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	result := make([][]any, 0, 1)
	for rows.Next() {
		values, pointers := make([]any, len(columns)), make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		require.NoError(t, rows.Scan(pointers...))
		result = append(result, values)
	}
	require.NoError(t, rows.Err())
	return result
}
