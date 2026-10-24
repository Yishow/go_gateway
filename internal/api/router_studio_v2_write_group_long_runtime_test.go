package api

import (
	"database/sql"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

// TestLongGroupAppliedByTheServiceRunsInTheProductionPipeline applies a long
// group through the real route, then runs the production pipeline over the proof
// Apply froze: nothing about the applied state, the table identity or the layout
// is built by the test. A sample taken after the cutover reaches the shared
// long table, and the legacy writer is suppressed for the source.
func TestLongGroupAppliedByTheServiceRunsInTheProductionPipeline(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedDestinationFixture(t, kind)
			group := ensureLong(t, f, "instantaneous", "second", 15, "")
			groupID := group["id"].(string)
			table := destinationOf(group)["table_name"].(string)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, groupID, previewLong(t, f, groupID)).Status)

			ready, err := f.groups.Readiness(t.Context(), groupID)
			require.NoError(t, err)
			require.True(t, ready.Ready, "a confirmed plain instantaneous long group is ready: %+v", ready.Issues)

			applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
			require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
			snapshot, err := f.groups.ResolveAppliedAt(t.Context(), groupID, time.Now().Add(time.Minute))
			require.NoError(t, err)
			store := groupdelivery.NewStore(f.db)
			frozen, err := store.RuntimeVersion(t.Context(), groupdelivery.GroupKey{
				WorkspaceID: snapshot.WorkspaceID, GroupID: groupID, GroupRevision: snapshot.AppliedRevision,
			})
			require.NoError(t, err)
			require.NotEmpty(t, frozen, "Apply froze the verified layout for recovery")

			var clock atomic.Int64
			now := func() time.Time { return time.Unix(0, clock.Load()).UTC() }
			clock.Store(snapshot.EffectiveAt.Add(time.Second).UnixNano())
			targetSvc := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
			pipe := grouppipeline.New(grouppipeline.Dependencies{
				Groups: f.groups, Tags: tag.NewService(tag.NewSQLRepository(f.db)), Destinations: targetSvc,
				Inspector: dbtarget.NewReadOnlyTableInspector(targetSvc), Store: store,
			}, grouppipeline.Config{
				NodeID: "node", Owner: "node/test", Now: now, TickInterval: 20 * time.Millisecond, WorkerInterval: 20 * time.Millisecond,
				ReconcileInterval: time.Hour, MaxFutureSkew: time.Minute,
			})
			require.NoError(t, pipe.Start(t.Context()))
			t.Cleanup(func() { _ = pipe.Stop(time.Second) })
			var active bool
			for _, status := range pipe.Status() {
				active = active || status.GroupID == groupID && status.State == grouppipeline.StateActive
			}
			require.True(t, active, "the production pipeline builds the group from Apply's own proof: %+v", pipe.Status())
			owned, err := store.HasAppliedOwner(t.Context(), "connector-A", "tag-A")
			require.NoError(t, err)
			require.True(t, owned, "Apply's durable ownership suppresses the legacy writer for the source")
			require.True(t, pipe.Owns("connector-A", "tag-A"))

			member := snapshot.Group.Members[0] // the source revisions Apply pinned
			observed := snapshot.EffectiveAt.Add(2 * time.Second)
			require.NoError(t, pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{
				SampleID: "sample-1", WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: "point-A", TagID: "tag-A",
				SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, ObservedAt: observed, ReceivedAt: observed.Add(30 * time.Millisecond),
				TimeOrigin: "gateway", Value: 21.5, Quality: schema.QualityGood,
			}))
			clock.Store(snapshot.EffectiveAt.Add(16 * time.Second).UnixNano())

			target := f.target
			query := `SELECT status, value_float, tag_id, profile_kind FROM ` + table
			if kind == "sqlite" {
				var err error
				target, err = sql.Open("sqlite", f.file)
				require.NoError(t, err)
				t.Cleanup(func() { _ = target.Close() })
			} else {
				query = `SELECT status, value_float, tag_id, profile_kind FROM "` + f.namespace + `"."` + table + `"`
			}
			require.Eventually(t, func() bool {
				var n int
				return target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+qualified(kind, f.namespace, table)).Scan(&n) == nil && n == 1
			}, 10*time.Second, 50*time.Millisecond, "the closed window reaches the shared long table")
			var status, tagID, profile string
			var value float64
			require.NoError(t, target.QueryRowContext(t.Context(), query).Scan(&status, &value, &tagID, &profile))
			require.Equal(t, [4]any{"ok", 21.5, "tag-A", "instantaneous"}, [4]any{status, value, tagID, profile})
		})
	}
}

func qualified(kind, namespace, table string) string {
	if kind == "postgres" {
		return `"` + namespace + `"."` + table + `"`
	}
	return table
}

// TestCumulativeLongGroupCountsFromASavedMeasurement applies a cumulative long
// group whose source names a saved counter measurement, then runs the
// production pipeline with the production counter-source wiring: the exact
// difference above 2^53 reaches the shared long table with the measurement's
// confirmed unit and series epoch, and a source whose measurement does not
// confirm counter semantics is not built.
func TestCumulativeLongGroupCountsFromASavedMeasurement(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			f := newManagedDestinationFixture(t, kind, "uint64")
			_, err := f.db.ExecContext(t.Context(), `
				INSERT INTO measurement_definitions (
					id, workspace_id, device_id, point_id, tag_id, equipment_id, definition_revision, source_binding_revision,
					series_epoch, name, quantity, unit, semantic_kind
				) VALUES ('meas-energy', ?, 'device-A', 'point-A', 'tag-A', 'device-A', 'def-5', 'bind-5', 'epoch-3', 'Energy', 'energy', 'kWh', 'counter')`,
				f.record.ID)
			require.NoError(t, err)

			request := profileRequest(t, f.writeGroupRouterFixture, "cumulative", 10, "minute", "main", "")
			saved := request["group"].(map[string]any)
			saved["layout_version"] = "managed_long_v1"
			saved["destination"].(map[string]any)["table_schema"] = f.namespace
			saved["members"].([]map[string]any)[0]["measurement_id"] = "meas-energy"
			created := performJSONRequest(t, f.router, http.MethodPost, ensureCircuitProfilePath, request)
			require.Equal(t, http.StatusOK, created.Code, created.Body.String())
			_, group := groupSaveData(t, decodeJSONBody(t, created))
			groupID := group["id"].(string)
			table := destinationOf(group)["table_name"].(string)
			require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, f, groupID, previewLong(t, f, groupID)).Status)

			ready, err := f.groups.Readiness(t.Context(), groupID)
			require.NoError(t, err)
			require.True(t, ready.Ready, "%+v", ready.Issues)
			applied := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/apply", scopeBody(t, f, groupID))
			require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
			snapshot, err := f.groups.ResolveAppliedAt(t.Context(), groupID, time.Now().Add(time.Hour))
			require.NoError(t, err)

			var clock atomic.Int64
			clock.Store(snapshot.EffectiveAt.Add(time.Second).UnixNano())
			targetSvc := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
			pipe := grouppipeline.New(grouppipeline.Dependencies{
				Groups: f.groups, Tags: tag.NewService(tag.NewSQLRepository(f.db)), Destinations: targetSvc,
				Inspector: dbtarget.NewReadOnlyTableInspector(targetSvc), Store: groupdelivery.NewStore(f.db),
				Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(f.db)),
			}, grouppipeline.Config{
				NodeID: "node", Owner: "node/test", Now: func() time.Time { return time.Unix(0, clock.Load()).UTC() },
				TickInterval: 20 * time.Millisecond, WorkerInterval: 20 * time.Millisecond, ReconcileInterval: time.Hour, MaxFutureSkew: time.Minute,
			})
			require.NoError(t, pipe.Start(t.Context()))
			t.Cleanup(func() { _ = pipe.Stop(time.Second) })
			require.Equal(t, grouppipeline.StateActive, pipeStatus(pipe, groupID).State, "%+v", pipe.Status())

			member := snapshot.Group.Members[0]
			accept := func(id string, observed time.Time, value uint64) {
				clock.Store(observed.Add(time.Second).UnixNano())
				require.NoError(t, pipe.AcceptSample(t.Context(), measurement.SampleEnvelope{
					SampleID: id, WorkspaceID: f.record.ID, DeviceID: "device-A", PointID: "point-A", TagID: "tag-A",
					SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision,
					ObservedAt: observed, ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway",
					Value: value, Quality: schema.QualityGood,
				}))
			}
			effective := snapshot.EffectiveAt
			accept("kwh-1", effective.Add(9*time.Minute), 9007199254740992)
			clock.Store(effective.Add(10*time.Minute + time.Second).UnixNano())
			accept("kwh-2", effective.Add(19*time.Minute), 9007199254740993)
			clock.Store(effective.Add(20*time.Minute + time.Second).UnixNano())

			target := f.target
			if kind == "sqlite" {
				target, err = sql.Open("sqlite", f.file)
				require.NoError(t, err)
				t.Cleanup(func() { _ = target.Close() })
			}
			from := qualified(kind, f.namespace, table)
			require.Eventually(t, func() bool {
				var n int
				return target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+from).Scan(&n) == nil && n == 2
			}, 15*time.Second, 50*time.Millisecond, "both windows reach the shared long table")
			rows, err := target.QueryContext(t.Context(), `SELECT status, CAST(value_decimal AS TEXT), CAST(delta_value AS TEXT), unit, series_epoch, measurement_id, profile_kind FROM `+from+` ORDER BY bucket_start`)
			require.NoError(t, err)
			defer rows.Close()
			type stored struct {
				status, raw, unit, epoch, measurement, profile string
				delta                                          sql.NullString
			}
			var got []stored
			for rows.Next() {
				var row stored
				require.NoError(t, rows.Scan(&row.status, &row.raw, &row.delta, &row.unit, &row.epoch, &row.measurement, &row.profile))
				got = append(got, row)
			}
			require.NoError(t, rows.Err())
			require.Len(t, got, 2)
			require.Equal(t, stored{status: "baseline", raw: "9007199254740992", unit: "kWh", epoch: "epoch-3", measurement: "meas-energy", profile: "cumulative"}, got[0])
			require.Equal(t, "ok", got[1].status)
			require.Equal(t, "9007199254740993", got[1].raw, "exact above 2^53 in the destination")
			require.Equal(t, "1", got[1].delta.String, "the exact difference")
		})
	}
}

func pipeStatus(pipe *grouppipeline.Pipeline, groupID string) grouppipeline.GroupStatus {
	for _, status := range pipe.Status() {
		if status.GroupID == groupID {
			return status
		}
	}
	return grouppipeline.GroupStatus{}
}
