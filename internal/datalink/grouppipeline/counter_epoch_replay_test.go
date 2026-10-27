package grouppipeline

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// These are two explicitly persisted applied descriptors. The service's draft /
// review / Apply transaction is covered separately in workspace; no descriptor
// is changed after its first accepted input.
func epochReplaySnapshot(t *testing.T, f *longPipeline, original []byte, revision string, raw bool) *workspace.WriteGroupAppliedSnapshot {
	t.Helper()
	var snap workspace.WriteGroupAppliedSnapshot
	require.NoError(t, json.Unmarshal(original, &snap))
	snap.GroupRevision, snap.AppliedRevision = revision, revision
	snap.Group.Revision, snap.Group.AppliedRevision = revision, revision
	if revision == "rev-2" {
		snap.EffectiveAt = at(1200)
	}
	snap.Group.RecordingProfile.Kind = workspace.WriteGroupProfileCumulative
	for i := range snap.Group.Members {
		m := &snap.Group.Members[i]
		m.ValueRole = ""
		if raw {
			snap.Group.RecordingProfile.Kind = workspace.WriteGroupProfileInstantaneous
			m.ValueRole = workspace.WriteGroupValueRoleCounterRaw
			m.SourceValueType = snap.RuntimeLayout.TagTypes[m.TagID]
			m.ResultValueType = measurement.ExactUint64
			if m.ExactRecordingSource != nil {
				m.ResultValueType = measurement.ExactDecimal
			}
		}
		if revision == "rev-2" && m.ExactRecordingSource != nil {
			m.ExactRecordingSource.DefinitionRevision = "def-2"
			m.ExactRecordingSource.SourceBindingRevision = "bind-2"
			m.ExactRecordingSource.SeriesEpoch = "epoch-2"
			m.ExactRecordingSource.Unit = "MWh"
		}
	}
	sources := map[string]runtime.CounterSource{}
	for tag, source := range f.counters.sources {
		if revision == "rev-2" {
			source.DefinitionRevision, source.SourceBindingRevision, source.SeriesEpoch, source.Unit = "def-2", "bind-2", "epoch-2", "MWh"
		}
		sources[tag] = source
	}
	freezeFixtureCounters(t, &snap, sources)
	return &snap
}

func TestCounterEpochReplayKeepsOldBacklogAndNewOwnerIndependent(t *testing.T) {
	for _, scaled := range []bool{false, true} {
		for _, oldRaw := range []bool{false, true} {
			for _, newRaw := range []bool{false, true} {
				t.Run(fmt.Sprintf("scaled=%t/old-raw=%t/new-raw=%t", scaled, oldRaw, newRaw), func(t *testing.T) {
					f := newLongPipelineOn(t, nil, scaled)
					oldKey := groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-C", GroupRevision: "rev-1"}
					nextKey := oldKey
					nextKey.GroupRevision = "rev-2"
					initial, err := f.store.RuntimeVersion(t.Context(), oldKey)
					require.NoError(t, err)
					old := epochReplaySnapshot(t, f, initial, "rev-1", oldRaw)
					original, err := json.Marshal(old)
					require.NoError(t, err)
					// Fixture construction only: no samples, checkpoints or effects exist yet.
					_, err = f.local.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_revision='rev-1'`, string(original))
					require.NoError(t, err)
					f.group, f.snapshot = old.Group, old
					f.groups.groups, f.groups.snapshots = []*workspace.WriteGroup{old.Group}, map[string][]*workspace.WriteGroupAppliedSnapshot{old.GroupID: {old}}
					f.startPipeline()
					multiplier := uint64(1)
					if scaled {
						multiplier = 10
					}
					require.NoError(t, f.reading(0, "old-baseline", at(1), 100*multiplier))
					f.closeUntil(at(600))
					require.NoError(t, f.reading(0, "same-id", at(601), 105*multiplier))
					oldState, err := f.store.Restore(t.Context(), oldKey)
					require.NoError(t, err)
					require.Len(t, oldState.Samples, 1)
					next := epochReplaySnapshot(t, f, initial, "rev-2", newRaw)
					nextBytes, err := json.Marshal(next)
					require.NoError(t, err)
					require.NoError(t, f.store.SaveRuntimeVersion(t.Context(), nextKey, nextBytes))
					f.group, f.snapshot = next.Group, next
					f.groups.groups, f.groups.snapshots = []*workspace.WriteGroup{next.Group}, map[string][]*workspace.WriteGroupAppliedSnapshot{next.GroupID: {old, next}}
					// This epoch/isolation fixture serializes destination partitions. It does
					// not test SQLite's concurrent-writer retry timing (separate delivery tests).
					reader := &unavailableExactCounterReader{}
					restart := func() {
						f.pipe = New(Dependencies{Groups: f.groups, Destinations: f.dest, Store: f.store, Counters: reader}, Config{Now: func() time.Time { return f.clock }, Delivery: groupdelivery.DispatcherConfig{MaxPartitions: 1}})
						require.NoError(t, f.pipe.Reconcile(t.Context()))
						found := false
						for _, status := range f.pipe.Status() {
							if status.Revision == "rev-2" {
								found = true
								require.Equal(t, StateActive, status.State, status.Reason)
							}
						}
						require.True(t, found)
					}
					f.clock = at(1202)
					restart()
					// The same sample ID is independent across applied owners. Exact digits
					// exceed float's integer precision and are scaled once, never on replay.
					require.NoError(t, f.reading(0, "same-id", at(1201), uint64(9007199254740993)))
					unchanged, err := f.store.Restore(t.Context(), oldKey)
					require.NoError(t, err)
					require.Equal(t, oldState, unchanged, "new acceptance cannot consume old inputs or baseline")
					restart()
					require.NoError(t, f.reading(0, "same-id", at(1201), uint64(9007199254740993)))
					nextState, err := f.store.Restore(t.Context(), nextKey)
					require.NoError(t, err)
					require.Len(t, nextState.Samples, 1, "restart resend is a duplicate within its original owner")
					require.Empty(t, nextState.Baselines, "new owner never inherits old baseline")
					f.closeUntil(at(1800))
					require.Zero(t, reader.calls)
					assertEpochReplayRows(t, f, oldRaw, newRaw, scaled)
					digests := f.outboxDigests()
					require.Len(t, digests, 6, "three windows, two members, no replayed effect")
					f.deliver()
					require.Equal(t, 6, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
					require.Equal(t, 6, f.count(f.plainLocal(), `SELECT COUNT(*) FROM wg_delivery_receipts`))
					// Readback verifies the same immutable row identities/meaning reached SQL.
					var unit *string
					var raw string
					var delta *string
					require.NoError(t, f.plainTarget().QueryRow(`SELECT unit, CAST(value_decimal AS TEXT), CAST(delta_value AS TEXT) FROM `+f.table+` WHERE group_revision='rev-2' AND tag_id=?`, longKwhFirst).Scan(&unit, &raw, &delta))
					if newRaw {
						require.Nil(t, unit)
					} else {
						require.NotNil(t, unit)
						require.Equal(t, "MWh", *unit)
					}
					require.Nil(t, delta)
					want := "9007199254740993"
					if scaled {
						want = "900719925474099.3"
					}
					require.Equal(t, want, raw)
					var oldUnit, oldEpoch, oldDelta *string
					require.NoError(t, f.plainTarget().QueryRow(`SELECT unit,series_epoch,CAST(delta_value AS TEXT) FROM `+f.table+` WHERE group_revision='rev-1' AND tag_id=? AND bucket_start=?`, longKwhFirst, at(600).Format(longTime)).Scan(&oldUnit, &oldEpoch, &oldDelta))
					if oldRaw {
						require.Nil(t, oldUnit)
						require.Nil(t, oldEpoch)
						require.Nil(t, oldDelta)
					} else {
						require.Equal(t, "kWh", *oldUnit)
						require.Equal(t, "epoch-1", *oldEpoch)
						delta := "5"
						if scaled {
							delta = "5.0"
						}
						require.Equal(t, delta, *oldDelta)
					}
					restart()
					f.closeUntil(at(1800))
					f.deliver()
					require.Equal(t, digests, f.outboxDigests())
					require.Equal(t, 6, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
					for _, snap := range []*workspace.WriteGroupAppliedSnapshot{old, next} {
						wantUnit, wantEpoch := "kWh", "epoch-1"
						if snap.AppliedRevision == "rev-2" {
							wantUnit, wantEpoch = "MWh", "epoch-2"
						}
						for _, proof := range snap.CounterBindings {
							require.Equal(t, wantUnit, proof.Unit)
							require.Equal(t, wantEpoch, proof.SeriesEpoch)
						}
					}
					for key, expected := range map[groupdelivery.GroupKey][]byte{oldKey: original, nextKey: nextBytes} {
						got, err := f.store.RuntimeVersion(t.Context(), key)
						require.NoError(t, err)
						require.Equal(t, expected, got)
					}
					require.Zero(t, reader.calls)
				})
			}
		}
	}
}

func assertEpochReplayRows(t *testing.T, f *longPipeline, oldRaw, newRaw, scaled bool) {
	t.Helper()
	rows, err := f.local.Query(`SELECT group_revision,payload FROM wg_delivery_outbox ORDER BY group_revision,bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	got := map[string][]map[string]any{}
	for rows.Next() {
		var revision, payload string
		require.NoError(t, rows.Scan(&revision, &payload))
		row, err := groupdelivery.DecodeRowPayload([]byte(payload))
		require.NoError(t, err)
		cells := map[string]any{}
		for _, cell := range row.Cells {
			cells[cell.Column] = cell.Value
		}
		if cells["tag_id"] == longKwhFirst {
			got[revision] = append(got[revision], cells)
		}
	}
	require.NoError(t, rows.Err())
	require.Len(t, got["rev-1"], 2)
	require.Len(t, got["rev-2"], 1)
	old, next := got["rev-1"][1], got["rev-2"][0]
	if oldRaw {
		require.Nil(t, old["unit"])
	} else {
		require.Equal(t, "kWh", old["unit"])
	}
	if newRaw {
		require.Nil(t, next["unit"])
	} else {
		require.Equal(t, "MWh", next["unit"])
	}
	wantOld, wantDelta, wantNew := "105", "5", "9007199254740993"
	if scaled {
		wantOld, wantDelta, wantNew = "105.0", "5.0", "900719925474099.3"
	}
	require.Equal(t, wantOld, old["value_decimal"])
	require.Equal(t, wantNew, next["value_decimal"])
	if oldRaw {
		require.Nil(t, old["delta_value"])
		require.Nil(t, old["series_epoch"])
	} else {
		require.Equal(t, wantDelta, old["delta_value"])
		require.Equal(t, "epoch-1", old["series_epoch"])
		require.Contains(t, old["provenance"], `"definition_revision":"def-1"`)
	}
	require.Nil(t, next["delta_value"], "new epoch starts without a derived delta")
	if newRaw {
		require.Nil(t, next["series_epoch"])
	} else {
		require.Equal(t, "baseline", next["status"])
		require.Equal(t, "epoch-2", next["series_epoch"])
		require.Contains(t, next["provenance"], `"definition_revision":"def-2"`)
	}
	for revision, raw := range map[string]bool{"rev-1": oldRaw, "rev-2": newRaw} {
		var count int
		require.NoError(t, f.local.QueryRow(`SELECT COUNT(*) FROM wg_delivery_counter_baselines WHERE group_revision=?`, revision).Scan(&count))
		want := 1
		if raw {
			want = 0
		}
		require.Equal(t, want, count, "only cumulative creates a baseline in its own revision")
	}
}
