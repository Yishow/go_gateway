package main

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductionEntityRowsSQLAndFrozenRestart(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		for _, scenario := range []string{"shared", "different", "case-alias", "quoted-distinct", "missing", "partial", "silent"} {
			t.Run(kind+"/"+scenario, func(t *testing.T) {
				if kind == "sqlite" && scenario == "quoted-distinct" {
					t.Skip("SQLite identifiers cannot differ only by case")
				}
				target := newEntityTarget(t, kind)
				if scenario == "quoted-distinct" {
					_, err := target.db.ExecContext(t.Context(), `ALTER TABLE `+target.name+` ADD COLUMN "Temperature" DOUBLE PRECISION`)
					require.NoError(t, err)
				}
				env, clock := target.env, &testClock{}
				// Keep injected bucket time behind the real store clock so retry due
				// checks remain eligible after concurrent SQLite writer contention.
				base := time.Now().UTC().Truncate(10 * time.Second).Add(-time.Minute)
				clock.set(base.Add(-time.Second))
				env.services.writeGroups.WithClock(clock.now)
				group := target.create(t, func(g *workspace.WriteGroup) {
					switch scenario {
					case "different":
						g.Members[2].TargetColumn = "other_temperature"
						g.Members[3].TargetColumn = "other_pressure"
					case "case-alias":
						g.RowPolicy.EntityKeyColumn = "ENTITY"
						for i := range g.Members {
							g.Members[i].TargetColumn = strings.ToUpper(g.Members[i].TargetColumn)
						}
					case "quoted-distinct":
						g.Members[2].TargetColumn = "Temperature"
					case "partial":
						g.RowPolicy.IncompletePolicy = "partial"
						g.Members[1].Required = false
					}
				})
				applied, err := env.services.writeGroups.Apply(t.Context(), group.ID, lifecycleMutation(t, env, group.ID))
				require.NoError(t, err)
				group = applied.Group
				clock.set(base.Add(time.Second))
				first := env.pipeline(clock, "node-1/entity-before-restart")
				require.NoError(t, first.Reconcile(t.Context()))
				if scenario != "silent" {
					values := []any{21.5, int64(101), 42.5, int64(202)}
					at := base.Add(2 * time.Second)
					clock.set(at)
					for i, value := range values {
						if i == 1 && (scenario == "missing" || scenario == "partial") {
							continue
						}
						require.NoError(t, first.AcceptSample(t.Context(), env.envelope(group, i, fmt.Sprintf("entity-%d", i), at, value)))
					}
				}
				var frozenBefore string
				require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_runtime_versions WHERE group_id = ?`, group.ID).Scan(&frozenBefore))
				// Reopen the real configuration store before any bucket has closed.
				reopenLifecycleConfiguration(t, env, clock)
				recovered := env.pipeline(clock, "node-1/entity-after-restart")
				require.NoError(t, recovered.Reconcile(t.Context()))
				var frozenAfter string
				require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT payload FROM wg_runtime_versions WHERE group_id = ?`, group.ID).Scan(&frozenAfter))
				require.Equal(t, frozenBefore, frozenAfter, "restart preserves the applied member partition and inspected SQL identifiers")
				clock.set(base.Add(11 * time.Second))
				recovered.TickAll(t.Context())
				want := 2
				switch scenario {
				case "missing":
					want = 1
				case "silent":
					want = 0
				}
				require.Equal(t, want, outboxCount(t, env.db, "connector-A", "1=1"))
				require.Zero(t, lifecycleJournalCount(t, env, group.AppliedRevision, true))
				var buckets, checkpoints int
				require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_buckets`).Scan(&buckets))
				require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_checkpoints`).Scan(&checkpoints))
				require.Equal(t, 2, buckets, "one outcome per entity in the same group and bucket")
				require.Equal(t, 1, checkpoints)
				assertEntityPayloads(t, env, group, scenario, want, base)
				require.NoError(t, recovered.Start(t.Context()))
				t.Cleanup(func() { require.NoError(t, recovered.Stop(2*time.Second)) })
				require.EventuallyWithT(t, func(c *assert.CollectT) {
					var states string
					if err := env.db.QueryRowContext(t.Context(), `SELECT COALESCE(group_concat(entity_key || ':' || state || ':' || last_error_code), '') FROM wg_delivery_outbox`).Scan(&states); err != nil {
						assert.NoError(c, err)
					}
					assert.Equal(c, want, outboxCount(t, env.db, "connector-A", "state = 'sql_committed'"), states)
				}, 10*time.Second, 20*time.Millisecond)
				assertEntitySQLRows(t, target, scenario, want)
				var receipts int
				require.NoError(t, env.db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM wg_delivery_receipts r JOIN wg_delivery_outbox o ON r.effect_key = o.effect_key AND r.payload_digest = o.payload_digest`).Scan(&receipts))
				require.Equal(t, want, receipts)
				require.NoError(t, recovered.Stop(2*time.Second))
				reopenLifecycleConfiguration(t, env, clock)
				again := env.pipeline(clock, "node-1/entity-settled-restart")
				require.NoError(t, again.Start(t.Context()))
				again.TickAll(t.Context())
				require.NoError(t, again.Stop(2*time.Second))
				assertEntitySQLRows(t, target, scenario, want)
				require.Equal(t, want, outboxCount(t, env.db, "connector-A", "1=1"), "closed rows retain their identities across another restart")
			})
		}
	}
}

func assertEntityPayloads(t *testing.T, env *outageEnv, group *workspace.WriteGroup, scenario string, want int, base time.Time) {
	t.Helper()
	rows, err := env.db.QueryContext(t.Context(), `SELECT record_id, effect_key, payload, payload_digest FROM wg_delivery_outbox ORDER BY entity_key`)
	require.NoError(t, err)
	defer rows.Close()
	identities := make(map[string]bool)
	for rows.Next() {
		var record, effect, payload, digest string
		require.NoError(t, rows.Scan(&record, &effect, &payload, &digest))
		row, err := groupdelivery.DecodeRowPayload([]byte(payload))
		require.NoError(t, err)
		require.Equal(t, record, row.RecordID)
		require.Equal(t, effect, row.EffectKey)
		require.False(t, identities[record])
		identities[record] = true
		require.Equal(t, base, row.BucketStart)
		actualDigest, err := groupdelivery.RowPayloadDigest(row)
		require.NoError(t, err)
		require.Equal(t, digest, actualDigest)
		cells := make(map[string]any)
		for _, cell := range row.Cells {
			cells[cell.Column] = cell.Value
		}
		require.Equal(t, row.EntityKey, cells["entity"])
		start, temp, pressure := 0, 21.5, int64(101)
		if row.EntityKey == "B" {
			start, temp, pressure = 2, 42.5, 202
		}
		tempColumn, pressureColumn := "temperature", "pressure"
		if scenario == "quoted-distinct" && row.EntityKey == "B" {
			tempColumn = "Temperature"
			require.NotContains(t, cells, "temperature")
		}
		if scenario == "different" && row.EntityKey == "B" {
			tempColumn, pressureColumn = "other_temperature", "other_pressure"
			require.NotContains(t, cells, "temperature")
			require.NotContains(t, cells, "pressure")
		} else {
			require.NotContains(t, cells, "other_temperature")
			require.NotContains(t, cells, "other_pressure")
		}
		require.Equal(t, temp, cells[tempColumn])
		if scenario == "partial" && row.EntityKey == "A" {
			require.Nil(t, cells[pressureColumn])
			require.True(t, row.Partial)
			provenance, ok := cells["prov"].(string)
			require.True(t, ok)
			require.Contains(t, provenance, group.Members[start].PointID)
			require.NotContains(t, provenance, group.Members[2].PointID)
		} else {
			require.Equal(t, pressure, cells[pressureColumn])
			require.False(t, row.Partial)
		}
	}
	require.NoError(t, rows.Err())
	require.Len(t, identities, want)
}

func assertEntitySQLRows(t *testing.T, target *entityTarget, scenario string, want int) {
	t.Helper()
	otherColumn := "other_temperature"
	if scenario == "quoted-distinct" {
		otherColumn = `"Temperature"`
	}
	rows, err := target.db.QueryContext(t.Context(), `SELECT entity, temperature, pressure, `+otherColumn+`, other_pressure, prov FROM `+target.name+` ORDER BY entity`)
	require.NoError(t, err)
	defer rows.Close()
	count := 0
	for rows.Next() {
		var entity string
		var temp, otherTemp sql.NullFloat64
		var pressure, otherPressure sql.NullInt64
		var provenance sql.NullString
		require.NoError(t, rows.Scan(&entity, &temp, &pressure, &otherTemp, &otherPressure, &provenance))
		count++
		if entity == "A" {
			require.NotEqual(t, "missing", scenario)
			require.Equal(t, sql.NullFloat64{Float64: 21.5, Valid: true}, temp)
			require.False(t, otherTemp.Valid)
			require.False(t, otherPressure.Valid)
			if scenario == "partial" {
				require.False(t, pressure.Valid)
				require.True(t, provenance.Valid)
			} else {
				require.Equal(t, sql.NullInt64{Int64: 101, Valid: true}, pressure)
			}
		} else {
			require.Equal(t, "B", entity)
			switch scenario {
			case "quoted-distinct":
				require.False(t, temp.Valid)
				require.Equal(t, sql.NullFloat64{Float64: 42.5, Valid: true}, otherTemp)
				require.Equal(t, sql.NullInt64{Int64: 202, Valid: true}, pressure)
				require.False(t, otherPressure.Valid)
			case "different":
				require.False(t, temp.Valid)
				require.False(t, pressure.Valid)
				require.Equal(t, sql.NullFloat64{Float64: 42.5, Valid: true}, otherTemp)
				require.Equal(t, sql.NullInt64{Int64: 202, Valid: true}, otherPressure)
			default:
				require.Equal(t, sql.NullFloat64{Float64: 42.5, Valid: true}, temp)
				require.Equal(t, sql.NullInt64{Int64: 202, Valid: true}, pressure)
				require.False(t, otherTemp.Valid)
				require.False(t, otherPressure.Valid)
			}
		}
	}
	require.NoError(t, rows.Err())
	require.Equal(t, want, count)
}
