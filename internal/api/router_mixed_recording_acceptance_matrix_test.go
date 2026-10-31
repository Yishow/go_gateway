package api

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func mixedNativeRequest(t *testing.T, wide bool, intervals []int64, required []bool) (fixture managedGroupFixture, payload map[string]any) {
	t.Helper()
	f := newManagedDestinationFixture(t, "sqlite")
	sources := []meanSource{{name: "a", tagType: "float64", role: "gauge"}, {name: "b", tagType: "float64", role: "gauge"}, {name: "c", tagType: "uint64", measurementID: "mixed-counter"}}
	seedMeanSources(t, f, sources)
	for _, s := range sources {
		_, err := f.db.ExecContext(t.Context(), `UPDATE points SET data_type=? WHERE id=?`, s.tagType, "point-"+s.name)
		require.NoError(t, err)
	}
	insertMeasurement(t, f, "mixed-counter", "c", "epoch-1", "counter")
	_, err := f.db.ExecContext(t.Context(), `UPDATE measurement_definitions SET numeric_encoding='uint64' WHERE id='mixed-counter'`)
	require.NoError(t, err)
	req := f.createRequest()
	g := req["group"].(map[string]any)
	g["record_id"] = uuid.NewString()
	g["layout_version"] = "managed_long_v1"
	if wide {
		g["layout_version"] = "managed_wide_profile_v1"
	}
	g["storage_mode"] = "plain"
	if wide {
		delete(g, "storage_mode")
	}
	g["recording_profile"] = map[string]any{"kind": "mixed", "result_shape": "member-calculation-v1"}
	g["row_policy"] = map[string]any{"incomplete_policy": "partial"}
	dest := g["destination"].(map[string]any)
	dest["storage_strategy"] = "managed"
	dest["table_schema"] = "main"
	dest["table_name"] = ""
	if wide {
		dest["table_name"] = "mixed_results"
	}
	members := make([]map[string]any, 0, len(sources))
	for i, s := range sources {
		method := "mean"
		if i == 2 {
			method = "delta"
		}
		m := map[string]any{"device_id": "device-A", "point_id": "point-" + s.name, "tag_id": "tag-" + s.name, "calculation": method, "required": required[i], "cadence": map[string]any{"interval": intervals[i], "unit": "second", "alignment": "utc_epoch"}}
		if s.role != "" {
			m["value_role"] = s.role
		}
		if s.measurementID != "" {
			m["measurement_id"] = s.measurementID
		}
		if wide {
			m["target_column"] = s.name
		}
		members = append(members, m)
	}
	g["members"] = members
	return f, req
}
func TestMixedAcceptanceCommonBoundaryWrites22_12_8(t *testing.T) {
	f, request := mixedNativeRequest(t, true, []int64{60, 60, 60}, []bool{true, true, true})
	base := time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)
	r := applyMixedAcceptance(t, f, request, base.Add(7*time.Second))
	r.offer(t, 2, "baseline", base.Add(10*time.Second), uint64(1000))
	r.tick(t, base.Add(time.Minute))
	require.Empty(t, r.deliver(t))
	for i, v := range []float64{20, 22, 24} {
		r.offer(t, 0, "a"+string(rune('0'+i)), base.Add(time.Duration(70+10*i)*time.Second), v)
	}
	for i, v := range []float64{10, 14} {
		r.offer(t, 1, "b"+string(rune('0'+i)), base.Add(time.Duration(72+10*i)*time.Second), v)
	}
	r.offer(t, 2, "counter", base.Add(100*time.Second), uint64(1008))
	r.tick(t, base.Add(2*time.Minute))
	require.Len(t, r.deliver(t), 1)
	var a, b float64
	var delta string
	require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT "`+r.frozen.Group.Members[0].TargetColumn+`","`+r.frozen.Group.Members[1].TargetColumn+`",CAST("`+r.frozen.Group.Members[2].TargetColumn+`" AS TEXT) FROM mixed_results`).Scan(&a, &b, &delta))
	require.Equal(t, float64(22), a)
	require.Equal(t, float64(12), b)
	require.Equal(t, "8", delta)
	var count int
	require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT count(*) FROM mixed_results`).Scan(&count))
	require.Equal(t, 1, count)
}
func TestMixedAcceptanceIndependentWideAndLongDueRows(t *testing.T) {
	for _, wide := range []bool{true, false} {
		name := "long"
		if wide {
			name = "wide"
		}
		t.Run(name, func(t *testing.T) {
			f, request := mixedNativeRequest(t, wide, []int64{15, 60, 60}, []bool{true, true, false})
			base := time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)
			r := applyMixedAcceptance(t, f, request, base.Add(7*time.Second))
			r.offer(t, 0, "a10", base.Add(10*time.Second), float64(20))
			r.offer(t, 1, "b10", base.Add(10*time.Second), float64(10))
			r.offer(t, 2, "c10", base.Add(10*time.Second), uint64(1000))
			r.tick(t, base.Add(15*time.Second))
			first := r.deliver(t)
			require.Len(t, first, 1, "not-due required slow members must not suppress fast output")
			table := r.frozen.Group.Destination.TableName
			if wide {
				var a float64
				var b, c sql.NullString
				require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT "`+r.frozen.Group.Members[0].TargetColumn+`","`+r.frozen.Group.Members[1].TargetColumn+`","`+r.frozen.Group.Members[2].TargetColumn+`" FROM "`+table+`"`).Scan(&a, &b, &c))
				require.Equal(t, float64(20), a)
				require.False(t, b.Valid)
				require.False(t, c.Valid)
			} else {
				var tag string
				require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT tag_id FROM "`+table+`"`).Scan(&tag))
				require.Equal(t, "tag-a", tag)
			}
			// SQL actual coverage remains partial while immutable row identity is aligned.
			var rawEvidence string
			require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT provenance FROM "`+table+`"`).Scan(&rawEvidence))
			var entry map[string]any
			if wide {
				var entries []map[string]any
				require.NoError(t, json.Unmarshal([]byte(rawEvidence), &entries))
				entry = entries[0]
				require.Equal(t, "not_due", entries[1]["status"])
			} else {
				require.NoError(t, json.Unmarshal([]byte(rawEvidence), &entry))
			}
			window := entry["window"].(map[string]any)
			require.Equal(t, base.Format("2006-01-02T15:04:05.000000Z"), window["aligned_start"])
			require.Equal(t, base.Add(7*time.Second).Format("2006-01-02T15:04:05.000000Z"), window["actual_start"])
			// Exactly T belongs to the next fast window, never to the already accepted row.
			r.offer(t, 0, "a15-boundary", base.Add(15*time.Second), float64(30))
			r.tick(t, base.Add(30*time.Second))
			require.Len(t, r.deliver(t), 1)
			var boundaryValue float64
			if wide {
				require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT "`+r.frozen.Group.Members[0].TargetColumn+`" FROM "`+table+`" ORDER BY due_at DESC LIMIT 1`).Scan(&boundaryValue))
			} else {
				require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT value_float FROM "`+table+`" ORDER BY bucket_start DESC LIMIT 1`).Scan(&boundaryValue))
			}
			require.Equal(t, float64(30), boundaryValue, "sample exactly :15 belongs only to next member window")
			r.offer(t, 0, "a40", base.Add(40*time.Second), float64(40))
			r.tick(t, base.Add(45*time.Second))
			require.Len(t, r.deliver(t), 1)
			r.offer(t, 0, "a50", base.Add(50*time.Second), float64(50))
			r.tick(t, base.Add(60*time.Second))
			last := r.deliver(t)
			if wide {
				require.Len(t, last, 1)
				var count int
				require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT count(*) FROM "`+table+`"`).Scan(&count))
				require.Equal(t, 4, count)
			} else {
				require.Len(t, last, 3, "only the three due members produce outputs, including optional null baseline")
				var count int
				require.NoError(t, r.f.target.QueryRowContext(t.Context(), `SELECT count(*) FROM "`+table+`"`).Scan(&count))
				require.Equal(t, 6, count)
			}
		})
	}
}
