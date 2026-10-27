package grouppipeline

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type changedCounterReader struct {
	calls   int
	sources map[string]runtime.CounterSource
	err     error
}

func (r *changedCounterReader) CounterSources(context.Context, *workspace.WriteGroup) (map[string]runtime.CounterSource, error) {
	r.calls++
	return r.sources, r.err
}

// freezeFixtureCounters represents facts captured before apply/acceptance, not
// a recovery fallback. All expected output below is independently literal.
func freezeFixtureCounters(t *testing.T, snap *workspace.WriteGroupAppliedSnapshot, sources map[string]runtime.CounterSource) {
	t.Helper()
	snap.CounterBindings = workspace.CounterBindingProofs{}
	for _, member := range snap.Group.Members {
		if member.ExactRecordingSource != nil {
			proof, err := workspace.CounterBindingProofFromExactSource(workspace.CounterBindingMemberIdentity{WorkspaceID: snap.WorkspaceID, GroupID: snap.GroupID, AppliedRevision: snap.AppliedRevision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID}, member, snap.RuntimeLayout.TagTypes[member.TagID])
			require.NoError(t, err)
			snap.CounterBindings[proof.Member.Key()] = proof
			continue
		}
		source := sources[member.TagID]
		kind, ok := measurement.ExactTypeForTag(snap.RuntimeLayout.TagTypes[member.TagID])
		require.True(t, ok)
		proof, err := workspace.NewCounterBindingProof(workspace.CounterBindingFacts{
			Member:        workspace.CounterBindingMemberIdentity{WorkspaceID: snap.WorkspaceID, GroupID: snap.GroupID, AppliedRevision: snap.AppliedRevision, DeviceID: member.DeviceID, PointID: member.PointID, TagID: member.TagID},
			MeasurementID: source.MeasurementID, DefinitionRevision: source.DefinitionRevision, SourceBindingRevision: source.SourceBindingRevision, SeriesEpoch: source.SeriesEpoch, Unit: source.Unit, SemanticKind: measurement.SemanticKind(source.Kind),
			SourceRevision: member.SourceRevision, MappingRevision: member.MappingRevision, SourceValueType: snap.RuntimeLayout.TagTypes[member.TagID], ResultValueType: kind, ScaleMode: workspace.CounterScaleUnscaled, ExactScale: "1",
		})
		require.NoError(t, err)
		snap.CounterBindings[proof.Member.Key()] = proof
	}
}

func TestCounterBindingReplayIgnoresDefinitionEditDeleteAndOutage(t *testing.T) {
	for _, scaled := range []bool{false, true} {
		for _, orphan := range []bool{false, true} {
			for _, change := range []string{"edit", "delete", "outage"} {
				t.Run(fmt.Sprintf("%s/scaled=%t/open-journal=%t", change, scaled, orphan), func(t *testing.T) {
					f := newLongPipelineOn(t, nil, scaled)
					key := groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-C", GroupRevision: "rev-1"}
					original, err := f.store.RuntimeVersion(t.Context(), key)
					require.NoError(t, err)
					multiplier := uint64(1)
					if scaled {
						multiplier = 10
					}
					require.NoError(t, f.reading(0, "c-1", at(1), 100*multiplier))
					require.NoError(t, f.reading(0, "c-2", at(601), 105*multiplier))
					before, err := f.store.Restore(t.Context(), key)
					require.NoError(t, err)
					require.Len(t, before.Samples, 2)
					reader := &changedCounterReader{}
					switch change {
					case "edit":
						reader.sources = map[string]runtime.CounterSource{}
						for tag, source := range f.counters.sources {
							source.Unit = "MWh"
							source.SeriesEpoch = "epoch-2"
							source.DefinitionRevision = "def-2"
							reader.sources[tag] = source
						}
					case "delete":
						reader.sources = map[string]runtime.CounterSource{}
					case "outage":
						reader.err = errors.New("private definition store unavailable")
					}
					// A fresh pipeline has no boundary memory, and must load accepted facts
					// from the actual durable descriptor/journal, not this current provider.
					if orphan {
						f.groups.groups = nil
					}
					f.pipe = New(Dependencies{Groups: f.groups, Destinations: f.dest, Store: f.store, Counters: reader}, Config{Now: func() time.Time { return f.clock }})
					require.NoError(t, f.pipe.Reconcile(t.Context()))
					if !orphan {
						require.Equal(t, StateActive, f.statusOf().State, f.statusOf().Reason)
					}
					replayed, err := f.store.Restore(t.Context(), key)
					require.NoError(t, err)
					require.Equal(t, before, replayed)
					require.Zero(t, reader.calls)
					f.closeUntil(at(1200))
					rows, err := f.store.DB().QueryContext(t.Context(), `SELECT payload, table_name FROM wg_delivery_outbox ORDER BY bucket_start`)
					require.NoError(t, err)
					defer rows.Close()
					var got []map[string]any
					for rows.Next() {
						var payload, table string
						require.NoError(t, rows.Scan(&payload, &table))
						require.Equal(t, f.table, table)
						row, err := groupdelivery.DecodeRowPayload([]byte(payload))
						require.NoError(t, err)
						cells := map[string]any{}
						for _, cell := range row.Cells {
							cells[cell.Column] = cell.Value
						}
						if cells["tag_id"] == longKwhFirst {
							got = append(got, cells)
						}
					}
					require.NoError(t, rows.Err())
					require.Len(t, got, 2)
					require.Equal(t, "baseline", got[0]["status"])
					require.Nil(t, got[0]["delta_value"])
					wantDelta := "5"
					if scaled {
						wantDelta = "5.0"
					}
					require.Equal(t, wantDelta, got[1]["delta_value"])
					require.Equal(t, "kWh", got[1]["unit"])
					require.Equal(t, "epoch-1", got[1]["series_epoch"])
					require.Contains(t, got[1]["provenance"], `"definition_revision":"def-1"`)
					after, err := f.store.RuntimeVersion(t.Context(), key)
					require.NoError(t, err)
					require.Equal(t, original, after)
					require.Zero(t, reader.calls)
				})
			}
		}
	}
}

func scaleFrozenFixture(snap *workspace.WriteGroupAppliedSnapshot) {
	for i := range snap.Group.Members {
		m := &snap.Group.Members[i]
		m.SourceValueType = snap.RuntimeLayout.TagTypes[m.TagID]
		m.ResultValueType = measurement.ExactDecimal
		m.ExactRecordingSource = &workspace.ExactRecordingSource{MeasurementID: *m.MeasurementID, DefinitionRevision: "def-1", SourceBindingRevision: "bind-1", SeriesEpoch: "epoch-1", OriginalDataType: m.SourceValueType, Scale: "0.1", Unit: "kWh", SemanticKind: measurement.SemanticKindCounter}
	}
}
