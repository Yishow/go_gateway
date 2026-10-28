package grouppipeline

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"
	"go-gateway/internal/datalink/workspace"
)

// Every descriptor below is an original, legal applied fixture established
// before acceptance. No replay edits, current measurement reads or new layout
// transition are used to make handoff succeed.
func reductionHandoffFixture(t *testing.T, shape averageShape, reduction string) (*ownershipPipelineFixture, *changedCounterReader) {
	t.Helper()
	f := newHandoffPipelineFixture(t, shape)
	// The existing forward matrix starts from managed wide, never long-to-long.
	group := averageGroup(shapeMeanWide)
	group.ID, group.Revision, group.AppliedRevision = f.source.GroupID, f.source.AppliedRevision, f.source.AppliedRevision
	group.Destination.TableName = "source_reduction_wide"
	statements, err := recordingplan.GenerateGroupSchemaDDL("sqlite", wideAverageLayout(group))
	require.NoError(t, err)
	_, err = f.plainTarget().Exec(statements[0]) // Receipt table already exists.
	require.NoError(t, err)
	f.source = snapshotOf(group, at(0))
	digest := sha256.Sum256([]byte(group.Destination.TableName))
	f.source.RuntimeLayout = &workspace.WriteGroupRuntimeLayout{
		Dialect: dbtarget.SQLDialectSQLite, Columns: sqliteColumnInfo(t, f.plainTarget(), group.Destination.TableName),
		TagTypes:          map[string]schema.DataType{longKwhFirst: schema.DataTypeUint64, longKwhSecond: schema.DataTypeUint64},
		ReceiptTableReady: true, SchemaDigest: hex.EncodeToString(digest[:]),
	}
	f.groups.groups = []*workspace.WriteGroup{group}
	f.groups.snapshots = map[string][]*workspace.WriteGroupAppliedSnapshot{group.ID: {f.source}}
	f.boundary, f.clock, f.target.EffectiveAt = at(10), at(7), at(10)
	for _, s := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
		s.GroupRevision = s.AppliedRevision
		g := s.Group
		g.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionLatest
		sources := map[string]runtime.CounterSource{}
		for i := range g.Members {
			m := &g.Members[i]
			m.ResultValueType = measurement.ExactUint64
			switch reduction {
			case "mean":
				g.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionMean
				m.ResultValueType = measurement.ExactDecimal
			case "state":
				g.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionMean
				m.ValueRole, m.SourceValueType, m.ResultValueType = workspace.WriteGroupValueRoleState, schema.DataTypeString, measurement.ExactText
				s.RuntimeLayout.TagTypes[m.TagID] = schema.DataTypeString
			case "counter_raw", "cumulative":
				id := fmt.Sprintf("measurement-%d", i)
				m.MeasurementID = &id
				m.ValueRole = workspace.WriteGroupValueRoleCounterRaw
				g.RecordingProfile.SampleReduction = workspace.WriteGroupSampleReductionMean
				if reduction == "cumulative" {
					g.RecordingProfile = &workspace.WriteGroupRecordingProfile{Kind: workspace.WriteGroupProfileCumulative, Interval: 1, Unit: workspace.WriteGroupIntervalUnitMinute}
					g.RowPolicy.IntervalSeconds = 60 // Existing legal cumulative unit.
					m.ValueRole = ""
					if g.LayoutVersion == workspace.WriteGroupLayoutManagedLongV1 {
						m.SourceValueType, m.ResultValueType = "", ""
					}
				}
				sources[m.TagID] = runtime.CounterSource{MeasurementID: id, DefinitionRevision: "def-original", SourceBindingRevision: "binding-original", SeriesEpoch: "epoch-original", Unit: "kWh", Kind: aggregation.CounterKindMonotonic}
			}
		}
		if reduction == "cumulative" && s == f.source {
			s.EffectiveAt = at(-60)
		}
		if reduction == "mean" && s == f.source {
			age := 2
			g.Members[0].MaxAgeSeconds = &age
		}
		if len(sources) > 0 {
			freezeFixtureCounters(t, s, sources)
		}
		requireLegalProfile(t, g)
		raw, err := json.Marshal(s)
		require.NoError(t, err)
		_, err = f.local.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=?`, string(raw), s.GroupID)
		require.NoError(t, err)
	}
	reader := &changedCounterReader{err: errors.New("current definitions must not participate in replay")}
	deps := f.pipe.deps
	deps.Counters = reader
	f.pipe = New(deps, f.pipe.config)
	require.NoError(t, f.pipe.Reconcile(t.Context()))
	require.NotNil(t, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)], "%+v", f.pipe.Status())
	return f, reader
}

type reductionRowEvidence struct {
	SampleID   string                                 `json:"sample_id"`
	Role       string                                 `json:"role"`
	ObservedAt string                                 `json:"observed_at"`
	Window     struct{ Start, End, Truncated string } `json:"window"`
	Mean       *struct {
		Sum              string
		ValidCount       int64 `json:"valid_count"`
		Coverage, Reason string
	} `json:"mean"`
	Counter *struct {
		Status           string
		Delta            *string `json:"delta"`
		DeltaFrom        *string `json:"delta_from"`
		DeltaTo          *string `json:"delta_to"`
		BaselineSampleID string  `json:"baseline_sample_id"`
		EndSampleID      string  `json:"end_sample_id"`
	} `json:"counter"`
	status                    string
	delta, deltaFrom, deltaTo sql.NullString
}

func readReductionRow(t *testing.T, f *ownershipPipelineFixture, owner *workspace.WriteGroupAppliedSnapshot, start time.Time, reduction string) (string, reductionRowEvidence) {
	t.Helper()
	column := "value_decimal"
	if reduction == "state" {
		column = "value_text"
	}
	wide := owner.Group.LayoutVersion == workspace.WriteGroupLayoutManagedWideProfileV1
	stamp := start.Format("2006-01-02T15:04:05.000000Z")
	query := `SELECT CAST(` + column + ` AS TEXT),provenance,observed_at,status,delta_value,delta_from,delta_to FROM ` + owner.Group.Destination.TableName + ` WHERE group_id=? AND tag_id='tag-kwh-1' AND bucket_start=?`
	if wide {
		query = `SELECT kwh_1,provenance FROM ` + owner.Group.Destination.TableName + ` WHERE group_id=? AND bucket_start=?`
		stamp = start.Format(time.RFC3339Nano)
	}
	var value sql.NullString
	var provenance string
	var columns reductionRowEvidence
	args := []any{&value, &provenance}
	if !wide {
		args = append(args, &columns.ObservedAt, &columns.status, &columns.delta, &columns.deltaFrom, &columns.deltaTo)
	}
	require.NoError(t, f.plainTarget().QueryRow(query, owner.GroupID, stamp).Scan(args...))
	require.True(t, value.Valid, "eligible input must remain usable at the actual owner boundary")
	var entries []reductionRowEvidence
	if !wide {
		provenance = "[" + provenance + "]"
	}
	require.NoError(t, json.Unmarshal([]byte(provenance), &entries))
	require.NotEmpty(t, entries)
	evidence := entries[0]
	if !wide {
		evidence.ObservedAt, evidence.status = columns.ObservedAt, columns.status
		evidence.delta, evidence.deltaFrom, evidence.deltaTo = columns.delta, columns.deltaFrom, columns.deltaTo
	} else if evidence.Counter != nil {
		evidence.status = evidence.Counter.Status
		for column, value := range map[*sql.NullString]*string{&evidence.delta: evidence.Counter.Delta, &evidence.deltaFrom: evidence.Counter.DeltaFrom, &evidence.deltaTo: evidence.Counter.DeltaTo} {
			if value != nil {
				*column = sql.NullString{String: *value, Valid: true}
			}
		}
	}
	return value.String, evidence
}

func TestCutoverReconcileReductionPartitionAndRecoveryMatrix(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		for _, reduction := range []string{"latest", "mean", "state", "counter_raw", "cumulative"} {
			t.Run(string(shape)+"/"+reduction, func(t *testing.T) {
				f, reader := reductionHandoffFixture(t, shape, reduction)
				if reduction == "cumulative" {
					f.clock = at(-2)
					require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, "previous-baseline", at(-1), uint64(5))))
					require.NoError(t, f.pipe.boundaries[boundaryKey(f.source.GroupID, f.source.AppliedRevision)].boundary.Tick(t.Context(), at(0)))
					f.clock = at(7)
				}
				// Exact spec inputs, in reverse arrival order: 10@8,20@9,100@11.
				for _, input := range []struct {
					id     string
					second int
					value  uint64
					state  string
				}{
					{"after", 11, 100, "future"}, {"later", 9, 20, "running"}, {"earlier", 8, 10, "idle"},
				} {
					var value any = input.value
					if reduction == "state" {
						value = input.state
					}
					require.NoError(t, f.pipe.AcceptSample(t.Context(), envelope(f.source.Group, input.id, at(input.second), value)))
				}
				before, err := f.store.Restore(t.Context(), groupKey(f.source))
				require.NoError(t, err)
				require.Len(t, before.Samples, 3)
				originalRows := handoffOriginalRows(t, f)
				f.commit()
				f.source.EndsAt = f.boundary
				f.groups.groups = append(f.groups.groups, f.target.Group)
				f.groups.snapshots[f.target.GroupID] = []*workspace.WriteGroupAppliedSnapshot{f.target}
				f.clock = at(70)
				require.NoError(t, f.pipe.Reconcile(t.Context()))
				require.Equal(t, originalRows, handoffOriginalRows(t, f))
				old, err := f.store.Restore(t.Context(), groupKey(f.source))
				require.NoError(t, err)
				next, err := f.store.Restore(t.Context(), groupKey(f.target))
				require.NoError(t, err)
				require.Len(t, old.Samples, 2)
				require.Len(t, next.Samples, 1)
				original := before.Samples[0]
				require.Equal(t, "after", original.SampleID)
				original.AcceptanceOrigin = snapshot.AcceptanceOrigin{WorkspaceID: f.source.WorkspaceID, GroupID: f.source.GroupID, GroupRevision: f.source.AppliedRevision}
				require.Equal(t, original, next.Samples[0], "owner changes, original proof/value/timestamps do not")
				for _, owner := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
					require.NoError(t, f.pipe.boundaries[boundaryKey(owner.GroupID, owner.AppliedRevision)].boundary.Tick(t.Context(), f.clock))
				}
				f.deliver()
				oldValue, oldEvidence := readReductionRow(t, f, f.source, at(0), reduction)
				newValue, newEvidence := readReductionRow(t, f, f.target, at(10), reduction)
				wantOld, wantNew := "20", "100"
				if reduction == "mean" {
					wantOld = "15"
				}
				if reduction == "state" {
					wantOld, wantNew = "running", "future"
				}
				require.Equal(t, wantOld, oldValue)
				require.Equal(t, wantNew, newValue)
				require.Equal(t, "later", oldEvidence.SampleID)
				require.Equal(t, "after", newEvidence.SampleID)
				for _, selected := range []struct {
					actual string
					want   time.Time
				}{{oldEvidence.ObservedAt, at(9)}, {newEvidence.ObservedAt, at(11)}} {
					observed, err := time.Parse(time.RFC3339Nano, selected.actual)
					require.NoError(t, err)
					require.Equal(t, selected.want, observed, "SQLite readback preserves selected observation time")
				}
				require.Equal(t, at(10).Format("2006-01-02T15:04:05.000000Z"), oldEvidence.Window.End)
				require.Equal(t, oldEvidence.Window.End, newEvidence.Window.Start)
				require.Equal(t, "end", oldEvidence.Window.Truncated)
				require.Equal(t, "start", newEvidence.Window.Truncated)
				if reduction == "mean" {
					require.NotNil(t, oldEvidence.Mean)
					require.Equal(t, "30", oldEvidence.Mean.Sum)
					require.Equal(t, int64(2), oldEvidence.Mean.ValidCount)
					require.Equal(t, "partial", oldEvidence.Mean.Coverage)
					require.Equal(t, "boundary-truncated", oldEvidence.Mean.Reason)
				} else {
					require.Nil(t, oldEvidence.Mean)
				}
				if reduction == "cumulative" {
					require.NotNil(t, newEvidence.Counter)
					require.Equal(t, "baseline", newEvidence.status)
					require.False(t, newEvidence.delta.Valid || newEvidence.deltaFrom.Valid || newEvidence.deltaTo.Valid, "new owner baseline has no invented delta or delta bounds")
					require.Equal(t, "ok", oldEvidence.status)
					require.Equal(t, sql.NullString{String: "15", Valid: true}, oldEvidence.delta, "old owner compares 20 to its original baseline 5, never successor 100")
					require.Equal(t, "previous-baseline", oldEvidence.Counter.BaselineSampleID)
					require.Equal(t, "later", oldEvidence.Counter.EndSampleID)
					for _, owner := range []*workspace.WriteGroupAppliedSnapshot{f.source, f.target} {
						restored, err := f.store.Restore(t.Context(), groupKey(owner))
						require.NoError(t, err)
						require.Len(t, restored.Baselines, 1)
						baseline, err := aggregation.DecodeCounterCheckpoint(restored.Baselines[0].Payload)
						require.NoError(t, err)
						wantID := "later"
						if owner == f.target {
							wantID = "after"
						}
						require.Equal(t, wantID, baseline.Reading.SampleID)
					}
				} else {
					require.Nil(t, newEvidence.Counter)
				}
				var consumedByOld, consumedByNew int
				require.NoError(t, f.local.QueryRow(`SELECT consumed FROM wg_delivery_samples WHERE group_id=? AND json_extract(payload,'$.sample_id')='after'`, f.source.GroupID).Scan(&consumedByOld))
				require.NoError(t, f.local.QueryRow(`SELECT consumed FROM wg_delivery_samples WHERE group_id=?`, f.target.GroupID).Scan(&consumedByNew))
				require.Zero(t, consumedByOld)
				require.Equal(t, 1, consumedByNew)
				count := f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_outbox`)
				digests := f.outboxDigests()
				f.pipe = New(f.pipe.deps, f.pipe.config)
				require.NoError(t, f.pipe.Reconcile(t.Context()))
				f.pipe.TickAll(t.Context())
				require.Equal(t, count, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_outbox`), "recovery cannot recalculate or duplicate finalized outputs")
				require.Equal(t, digests, f.outboxDigests(), "finalized effect identities and payload digests remain frozen")
				require.Zero(t, reader.calls, "original frozen descriptors are the sole measurement authority")
			})
		}
	}
}
