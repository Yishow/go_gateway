package api

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func mysqlSuccessorRuntimeCounterRows(t *testing.T, r *mysqlSuccessorRig, g *workspace.WriteGroup) []counterWindow {
	t.Helper()
	target := r.targetDB
	if g.ID == r.groupID {
		target = r.f.target
	}
	from := qualified("mysql", g.Destination.TableSchema, g.Destination.TableName)
	rows, err := target.QueryContext(t.Context(), `SELECT status, CAST(value_decimal AS CHAR), CAST(value_int AS CHAR), CAST(value_float AS CHAR), CAST(delta_value AS CHAR), CAST(delta_from AS CHAR), CAST(delta_to AS CHAR), CAST(provenance AS CHAR) FROM `+from+` ORDER BY bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	var windows []counterWindow
	for rows.Next() {
		var status string
		var decimal, integer, float, delta, deltaFrom, deltaTo sql.NullString
		var provenance string
		require.NoError(t, rows.Scan(&status, &decimal, &integer, &float, &delta, &deltaFrom, &deltaTo, &provenance))
		var document struct {
			Counter *counterProvenanceView `json:"counter"`
		}
		require.NoError(t, json.Unmarshal([]byte(provenance), &document), provenance)
		require.NotNil(t, document.Counter)
		require.Equal(t, delta.Valid, deltaFrom.Valid)
		require.Equal(t, delta.Valid, deltaTo.Valid)
		window := counterWindowOf(t, document.Counter)
		window.status, window.provenance = status, provenance
		for _, raw := range []sql.NullString{decimal, integer, float} {
			if raw.Valid {
				window.raw = plainDecimal(raw.String)
			}
		}
		if delta.Valid {
			window.delta, window.hasDelta = plainDecimal(delta.String), true
		}
		if deltaFrom.Valid && deltaTo.Valid {
			window.deltaFrom, window.deltaTo = mysqlSuccessorStoredInstant(t, deltaFrom.String), mysqlSuccessorStoredInstant(t, deltaTo.String)
		}
		windows = append(windows, window)
	}
	require.NoError(t, rows.Err())
	return windows
}

// MySQL DATETIME has no timezone suffix; production stores UTC explicitly.
func mysqlSuccessorStoredInstant(t *testing.T, raw string) time.Time {
	t.Helper()
	value, err := time.Parse("2006-01-02 15:04:05.999999", raw)
	require.NoError(t, err)
	return value.UTC()
}
