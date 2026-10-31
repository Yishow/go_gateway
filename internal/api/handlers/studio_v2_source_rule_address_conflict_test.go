package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceSourceRuleAddressConflictIsReadOnlyValidation(t *testing.T) {
	cases := []struct {
		name, method    string
		explicitEnabled bool
	}{{http.MethodPost, http.MethodPost, false}, {http.MethodPut, http.MethodPut, false}, {"PUT-enabled", http.MethodPut, true}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			method := tc.method
			f := newWorkspaceConcurrentFixture(t, 1, 2, nil)
			_, err := f.rules.Create(f.ctx, sourcerule.CreateRuleRequest{
				ID: "other-rule", DeviceID: "diag-dev", StartAddress: "40009", Count: 2,
				DataType: schema.DataTypeInt16, NamingPrefix: "OTHER", Enabled: true,
			})
			require.NoError(t, err)
			h := NewStudioV2WorkspaceSourceRulesHandler(f.h.workspaceSvc, f.h.deviceSvc, f.rules)
			router := gin.New()
			router.POST("/workspace/source-rules", h.Create)
			router.PUT("/workspace/source-rules/:id", h.Update)
			path := "/workspace/source-rules"
			body := `{"id":"conflicting-create","device_id":"diag-dev","start_address":"40001","count":2,"data_type":"int16","naming_prefix":"OTHER","enabled":true}`
			if method == http.MethodPut {
				path += "/other-rule"
				body = `{"start_address":"40001"}`
				if tc.explicitEnabled {
					body = `{"device_id":"diag-dev","start_address":"40001","count":2,"data_type":"int16","naming_prefix":"OTHER","enabled":true,"skipped_addresses":[],"scale_multiplier":1,"scale_offset":0,"data_format":"","share_enabled":false,"share_start_register":null,"share_stride":null}`
				}
			}
			before := addressConflictDatabaseSnapshot(t, f.db)
			var changesBefore int64
			require.NoError(t, f.db.QueryRow("SELECT total_changes()").Scan(&changesBefore))
			response := httptest.NewRecorder()
			request := httptest.NewRequestWithContext(f.ctx, method, path, strings.NewReader(body))
			request.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(response, request)
			// Validation is not an uncertain write: neither compensation nor a
			// timestamp trigger may run for a rejected address plan.
			after := addressConflictDatabaseSnapshot(t, f.db)
			var changesAfter int64
			require.NoError(t, f.db.QueryRow("SELECT total_changes()").Scan(&changesAfter))
			require.Equal(t, before, after, "rejected address plan changed configuration")
			require.Equal(t, changesBefore, changesAfter, "preflight ran a write or compensation")
			require.Equal(t, http.StatusBadRequest, response.Code, response.Body.String())
			var envelope map[string]any
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &envelope))
			problem := envelope["error"].(map[string]any)
			require.Equal(t, apiValidationErrorCode, problem["code"])
			require.Contains(t, problem["message"], "40001")
		})
	}
}

// The fixture has one connection and no background runtime. Preserve every
// durable row, including timestamps, source links, candidates and receipts.
func addressConflictDatabaseSnapshot(t *testing.T, db *sql.DB) map[string][]string {
	t.Helper()
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`)
	require.NoError(t, err)
	var tables []string
	for rows.Next() {
		var name string
		require.NoError(t, rows.Scan(&name))
		tables = append(tables, name)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	snapshot := make(map[string][]string, len(tables))
	for _, table := range tables {
		rows, err := db.Query(fmt.Sprintf(`SELECT * FROM "%s"`, strings.ReplaceAll(table, `"`, `""`)))
		require.NoError(t, err)
		columns, err := rows.Columns()
		require.NoError(t, err)
		var encoded []string
		for rows.Next() {
			values := make([]any, len(columns))
			pointers := make([]any, len(columns))
			for i := range values {
				pointers[i] = &values[i]
			}
			require.NoError(t, rows.Scan(pointers...))
			data, err := json.Marshal(values)
			require.NoError(t, err)
			encoded = append(encoded, string(data))
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		sort.Strings(encoded)
		snapshot[table] = encoded
	}
	return snapshot
}
