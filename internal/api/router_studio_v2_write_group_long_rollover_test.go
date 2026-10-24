package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
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
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// rolloverStack is a cumulative long group applied through the real route over a
// counter measurement saved through the measurement service, with the
// production pipeline (including the production counter-source wiring) running
// over it. Nothing about the counter semantics is filled in by the test.
type rolloverStack struct {
	t            *testing.T
	kind         string
	f            managedGroupFixture
	definitionID string
	groupID      string
	table        string
	effective    time.Time
	member       workspace.WriteGroupMember
	clock        atomic.Int64
	definitions  *measurement.Service
	target       *sql.DB
	pipe         *grouppipeline.Pipeline
}

func newRolloverStack(t *testing.T, kind string, policy *measurement.CounterPolicy) *rolloverStack {
	t.Helper()
	s := &rolloverStack{t: t, kind: kind, f: newManagedDestinationFixture(t, kind, "uint64")}
	s.definitions = measurement.NewService(measurement.NewSQLRepository(s.f.db))
	definition := &measurement.MeasurementDefinition{
		ID: "meas-energy", WorkspaceID: s.f.record.ID, DeviceID: "device-A", PointID: "point-A", Name: "Energy", Quantity: "energy",
		Unit: "kWh", SemanticKind: measurement.SemanticKindCounter, NumericEncoding: "uint64", CounterPolicy: policy,
	}
	tagID := "tag-A"
	definition.TagID = &tagID
	require.NoError(t, s.definitions.CreateMeasurement(t.Context(), definition), "the exact policy is saved through the service")
	s.definitionID = definition.ID

	request := profileRequest(t, s.f.writeGroupRouterFixture, "cumulative", 10, "minute", "main", "")
	saved := request["group"].(map[string]any)
	saved["layout_version"] = "managed_long_v1"
	saved["destination"].(map[string]any)["table_schema"] = s.f.namespace
	saved["members"].([]map[string]any)[0]["measurement_id"] = definition.ID
	created := performJSONRequest(t, s.f.router, http.MethodPost, ensureCircuitProfilePath, request)
	require.Equal(t, http.StatusOK, created.Code, created.Body.String())
	_, group := groupSaveData(t, decodeJSONBody(t, created))
	s.groupID = group["id"].(string)
	s.table = destinationOf(group)["table_name"].(string)
	require.Equal(t, recordingplan.SchemaOperationSucceeded, confirmLong(t, s.f, s.groupID, previewLong(t, s.f, s.groupID)).Status)
	applied := performJSONRequest(t, s.f.router, http.MethodPost, writeGroupsPath+"/"+s.groupID+"/apply", scopeBody(t, s.f, s.groupID))
	require.Equal(t, http.StatusOK, applied.Code, applied.Body.String())
	snapshot, err := s.f.groups.ResolveAppliedAt(t.Context(), s.groupID, time.Now().Add(time.Hour))
	require.NoError(t, err)
	s.effective, s.member = snapshot.EffectiveAt, snapshot.Group.Members[0]
	s.clock.Store(s.effective.Add(time.Second).UnixNano())

	s.target = s.f.target
	if kind == "sqlite" {
		s.target, err = sql.Open("sqlite", s.f.file)
		require.NoError(t, err)
		t.Cleanup(func() { _ = s.target.Close() })
	}
	t.Cleanup(func() { s.stop() })
	return s
}

// start runs a new production pipeline over the same local database, the way a
// restarted gateway would.
func (s *rolloverStack) start() {
	s.t.Helper()
	targetSvc := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(s.f.db))
	s.pipe = grouppipeline.New(grouppipeline.Dependencies{
		Groups: s.f.groups, Tags: tag.NewService(tag.NewSQLRepository(s.f.db)), Destinations: targetSvc,
		Inspector: dbtarget.NewReadOnlyTableInspector(targetSvc), Store: groupdelivery.NewStore(s.f.db),
		Counters: grouppipeline.NewMeasurementCounterSources(measurement.NewSQLRepository(s.f.db)),
	}, grouppipeline.Config{
		NodeID: "node", Owner: "node/test", Now: func() time.Time { return time.Unix(0, s.clock.Load()).UTC() },
		TickInterval: 20 * time.Millisecond, WorkerInterval: 20 * time.Millisecond, ReconcileInterval: time.Hour, MaxFutureSkew: time.Minute,
	})
	require.NoError(s.t, s.pipe.Start(s.t.Context()))
}

func (s *rolloverStack) stop() {
	if s.pipe != nil {
		_ = s.pipe.Stop(time.Second)
		s.pipe = nil
	}
}

func (s *rolloverStack) status() grouppipeline.GroupStatus { return pipeStatus(s.pipe, s.groupID) }

// accept delivers one reading through the production intake, nine minutes into the window.
func (s *rolloverStack) accept(id string, window int, value uint64) {
	s.t.Helper()
	observed := s.effective.Add(time.Duration(window)*10*time.Minute + 9*time.Minute)
	s.clock.Store(observed.Add(time.Second).UnixNano())
	require.NoError(s.t, s.pipe.AcceptSample(s.t.Context(), measurement.SampleEnvelope{
		SampleID: id, WorkspaceID: s.f.record.ID, DeviceID: "device-A", PointID: "point-A", TagID: "tag-A",
		SourceRevision: s.member.SourceRevision, MappingRevision: s.member.MappingRevision,
		ObservedAt: observed, ReceivedAt: observed.Add(30 * time.Millisecond), TimeOrigin: "gateway", Value: value, Quality: schema.QualityGood,
	}))
}

// closeWindow moves the clock past the end of a window so the pipeline closes it.
func (s *rolloverStack) closeWindow(window int) {
	s.clock.Store(s.effective.Add(time.Duration(window+1)*10*time.Minute + time.Second).UnixNano())
}

type rolloverRow struct {
	status, provenance string
	raw, delta         sql.NullString
}

// rows waits until the shared long table holds n records and returns them in
// window order, read back from the real destination.
func (s *rolloverStack) rows(n int) []rolloverRow {
	s.t.Helper()
	from := qualified(s.kind, s.f.namespace, s.table)
	require.Eventually(s.t, func() bool {
		var count int
		return s.target.QueryRowContext(s.t.Context(), `SELECT COUNT(*) FROM `+from).Scan(&count) == nil && count == n
	}, 15*time.Second, 50*time.Millisecond, "the closed windows reach the shared long table")
	result, err := s.target.QueryContext(s.t.Context(), `SELECT status, CAST(value_decimal AS TEXT), CAST(delta_value AS TEXT), CAST(provenance AS TEXT) FROM `+from+` ORDER BY bucket_start`)
	require.NoError(s.t, err)
	defer result.Close()
	var rows []rolloverRow
	for result.Next() {
		var row rolloverRow
		require.NoError(s.t, result.Scan(&row.status, &row.raw, &row.delta, &row.provenance))
		// PostgreSQL renders JSONB with spaces; compare the document, not its layout.
		var document any
		require.NoError(s.t, json.Unmarshal([]byte(row.provenance), &document))
		compact, err := json.Marshal(document)
		require.NoError(s.t, err)
		row.provenance = string(compact)
		rows = append(rows, row)
	}
	require.NoError(s.t, result.Err())
	return rows
}

// TestCumulativeLongGroupRolloverFromSavedPolicy proves the production path of
// the rollover evidence: saved through the measurement service, confirmed at
// pipeline start, applied by the 2.3 engine and read back from a real
// destination, with a changed policy starting a new baseline.
func TestCumulativeLongGroupRolloverFromSavedPolicy(t *testing.T) {
	for _, kind := range []string{"sqlite", "postgres"} {
		t.Run(kind, func(t *testing.T) {
			cases := map[string]struct {
				policy        *measurement.CounterPolicy
				status        string
				delta         string
				contains      string
				first, second uint64 // zero values mean 65530 then 4
			}{
				"a proven single wrap is 10":       {&measurement.CounterPolicy{RolloverModulus: "65536", MaxRatePerSecond: "100"}, "ok", "10", `"kind":"rollover"`, 0, 0},
				"a modulus above 2^53 stays exact": {&measurement.CounterPolicy{RolloverModulus: "18446744073709551616", MaxRatePerSecond: "1000000000000000.5"}, "ok", "10", `"rollover_modulus":"18446744073709551616"`, 18446744073709551610, 4},
				"a modulus alone is not evidence":  {&measurement.CounterPolicy{RolloverModulus: "65536"}, "reset_or_uncertain", "", `"reason":"rollover-evidence-insufficient"`, 0, 0},
				"a rate that allows several wraps": {&measurement.CounterPolicy{RolloverModulus: "65536", MaxRatePerSecond: "200"}, "reset_or_uncertain", "", `"reason":"rollover-multiple-wraps-possible"`, 0, 0},
				"no policy at all":                 {&measurement.CounterPolicy{}, "reset_or_uncertain", "", `"reason":"decrease-unproven"`, 0, 0},
				"the legacy float is not evidence": {&measurement.CounterPolicy{MaxRollOver: func() *float64 { v := 65536.0; return &v }()}, "reset_or_uncertain", "", `"reason":"decrease-unproven"`, 0, 0},
			}
			for name, c := range cases {
				t.Run(name, func(t *testing.T) {
					s := newRolloverStack(t, kind, c.policy)
					s.start()
					require.Equal(t, grouppipeline.StateActive, s.status().State, "%+v", s.pipe.Status())
					first, second := c.first, c.second
					if first == 0 {
						first, second = 65530, 4
					}
					s.accept("energy-1", 0, first)
					s.closeWindow(0)
					s.accept("energy-2", 1, second)
					s.closeWindow(1)
					rows := s.rows(2)
					require.Equal(t, "baseline", rows[0].status)
					require.Equal(t, strconv.FormatUint(first, 10), rows[0].raw.String)
					require.Equal(t, c.status, rows[1].status)
					require.Equal(t, strconv.FormatUint(second, 10), rows[1].raw.String, "the raw reading is always kept")
					if c.delta == "" {
						require.False(t, rows[1].delta.Valid, "no abs() and no zero for a drop that is not proven to be a wrap")
					} else {
						require.Equal(t, c.delta, rows[1].delta.String)
					}
					require.Contains(t, rows[1].provenance, c.contains)
				})
			}

			t.Run("a changed policy starts a new baseline after the restart", func(t *testing.T) {
				s := newRolloverStack(t, kind, &measurement.CounterPolicy{RolloverModulus: "65536", MaxRatePerSecond: "100"})
				s.start()
				s.accept("energy-1", 0, 65530)
				s.closeWindow(0)
				s.accept("energy-2", 1, 4)
				s.closeWindow(1)
				require.Equal(t, "10", s.rows(2)[1].delta.String)
				s.stop()

				definition, err := s.definitions.GetMeasurement(t.Context(), s.definitionID)
				require.NoError(t, err)
				definition.CounterPolicy = &measurement.CounterPolicy{RolloverModulus: "65536", MaxRatePerSecond: "50"}
				moved, err := s.definitions.UpdateMeasurement(t.Context(), definition)
				require.NoError(t, err)
				require.True(t, moved, "a changed rate bound is a new series epoch")

				s.start()
				require.Equal(t, grouppipeline.StateActive, s.status().State, "%+v", s.pipe.Status())
				s.accept("energy-3", 2, 9)
				s.closeWindow(2)
				rows := s.rows(3)
				require.Equal(t, "baseline", rows[2].status, "the first window of the revised evidence is only a baseline")
				require.False(t, rows[2].delta.Valid, "9 is never subtracted from a reading taken under other evidence")
				require.Contains(t, rows[2].provenance, `"reason":"binding-changed"`)
				require.Equal(t, "10", rows[1].delta.String, "the closed record keeps its frozen result")
			})
		})
	}
}
