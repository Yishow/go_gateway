package api

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func sqlServerSuccessorRuntimeCounterRows(t *testing.T, r *sqlServerSuccessorRig, g *workspace.WriteGroup) []counterWindow {
	t.Helper()
	target := r.targetDB
	if g.ID == r.groupID {
		target = r.f.target
	}
	from := sqlServerQualified(g.Destination.TableSchema, g.Destination.TableName)
	rows, err := target.QueryContext(t.Context(), `SELECT status, CONVERT(varchar(100),value_decimal), CONVERT(varchar(100),value_int), CONVERT(varchar(100),value_float), CONVERT(varchar(100),delta_value), delta_from, delta_to, provenance FROM `+from+` ORDER BY bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	var windows []counterWindow
	for rows.Next() {
		var status, provenance string
		var decimal, integer, float, delta sql.NullString
		var deltaFrom, deltaTo sql.NullTime
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
			window.deltaFrom, window.deltaTo = deltaFrom.Time.UTC(), deltaTo.Time.UTC()
		}
		windows = append(windows, window)
	}
	require.NoError(t, rows.Err())
	return windows
}
