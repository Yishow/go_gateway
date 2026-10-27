package grouppipeline

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// blockedLong rebuilds a pipeline over the fixture and returns the group's status.
func (f *longPipeline) reconcileStatus() GroupStatus {
	f.t.Helper()
	pipe := New(Dependencies{
		Groups: f.groups, Tags: fakeTags{missing: map[string]bool{}}, Destinations: f.dest, Inspector: &fakeInspector{},
		Store: f.store, Counters: f.counters,
	}, Config{NodeID: "node-1", Owner: "node-1/other", Now: func() time.Time { return f.clock }, ReconcileInterval: 5 * time.Second})
	require.NoError(f.t, pipe.Reconcile(f.t.Context()))
	f.pipe = pipe
	return f.statusOf()
}

// saveSnapshot replaces the frozen runtime version the way a (possibly forged)
// writer of the local database would.
func (f *longPipeline) saveSnapshot(mutate func(*workspace.WriteGroupAppliedSnapshot)) {
	f.t.Helper()
	snap, err := f.pipe.loadRuntimeSnapshot(f.t.Context(), groupdelivery.GroupKey{
		WorkspaceID: f.group.WorkspaceID, GroupID: f.group.ID, GroupRevision: f.group.AppliedRevision,
	})
	require.NoError(f.t, err)
	mutate(snap)
	payload, err := json.Marshal(snap)
	require.NoError(f.t, err)
	_, err = f.plainLocal().ExecContext(f.t.Context(), `UPDATE wg_runtime_versions SET payload = ?`, string(payload))
	require.NoError(f.t, err)
}

func TestLongGroupIsBuiltOnlyFromItsVerifiedFrozenProof(t *testing.T) {
	t.Run("a verified frozen snapshot starts the group", func(t *testing.T) {
		f := newLongPipeline(t)
		require.Equal(t, StateActive, f.statusOf().State)
		require.True(t, f.pipe.Owns("connector-1", longKwhFirst), "the group owns its outputs so the legacy writer is suppressed")
		require.True(t, f.pipe.WantsSample("device-1", "point-1", longKwhFirst))
	})

	t.Run("without the frozen snapshot nothing is rebuilt from the inspected table", func(t *testing.T) {
		f := newLongPipeline(t)
		_, err := f.plainLocal().ExecContext(t.Context(), `DELETE FROM wg_runtime_versions`)
		require.NoError(t, err)
		status := f.reconcileStatus()
		require.Equal(t, StateBlocked, status.State)
		require.Equal(t, "counter-replay-proof-missing", status.Reason)
		require.False(t, f.pipe.Owns("connector-1", longKwhFirst))
	})

	proofFailures := map[string]func(*workspace.WriteGroupAppliedSnapshot){
		"a table that is not the one derived from the storage scope": func(s *workspace.WriteGroupAppliedSnapshot) {
			s.Group.Destination.TableName = "gw_long_0000000000000000"
		},
		"a missing owner marker column": func(s *workspace.WriteGroupAppliedSnapshot) {
			columns := s.RuntimeLayout.Columns[:0:0]
			for _, column := range s.RuntimeLayout.Columns {
				if len(column.Name) < 10 || column.Name[:10] != "_gw_owner_" {
					columns = append(columns, column)
				}
			}
			s.RuntimeLayout.Columns = columns
		},
		"a foreign owner marker": func(s *workspace.WriteGroupAppliedSnapshot) {
			for i, column := range s.RuntimeLayout.Columns {
				if len(column.Name) >= 10 && column.Name[:10] == "_gw_owner_" {
					s.RuntimeLayout.Columns[i].Name = "_gw_owner_ffffffffffffffffffffffff"
				}
			}
		},
		"a column of the fixed layout missing": func(s *workspace.WriteGroupAppliedSnapshot) {
			s.RuntimeLayout.Columns = s.RuntimeLayout.Columns[1:]
		},
		"a storage mode the runtime does not write": func(s *workspace.WriteGroupAppliedSnapshot) {
			s.Group.StorageMode = workspace.WriteGroupStorageMode("timescale")
		},
		"a profile that does not derive the interval": func(s *workspace.WriteGroupAppliedSnapshot) {
			s.Group.RowPolicy.IntervalSeconds = 60
		},
		"no schema digest":                    func(s *workspace.WriteGroupAppliedSnapshot) { s.RuntimeLayout.SchemaDigest = "" },
		"a source type without an exact form": func(s *workspace.WriteGroupAppliedSnapshot) { delete(s.RuntimeLayout.TagTypes, longKwhFirst) },
		"receipt dedupe without a receipt table": func(s *workspace.WriteGroupAppliedSnapshot) {
			s.RuntimeLayout.ReceiptTableReady = false
		},
	}
	for name, mutate := range proofFailures {
		t.Run("a forged proof is refused: "+name, func(t *testing.T) {
			f := newLongPipeline(t)
			f.saveSnapshot(mutate)
			status := f.reconcileStatus()
			require.Equal(t, StateBlocked, status.State)
			want := reasonSnapshotUnavailable
			if name == "a source type without an exact form" {
				want = "counter-replay-proof-conflict"
			}
			require.Equal(t, want, status.Reason)
			require.False(t, f.pipe.Owns("connector-1", longKwhFirst), "a group that cannot run owns nothing")
		})
	}
}

type failingDefinitions struct{}

func (failingDefinitions) GetByID(context.Context, string) (*measurement.MeasurementDefinition, error) {
	return nil, errors.New("measurement store offline")
}

func TestCumulativeLongGroupNeedsOriginalCounterProof(t *testing.T) {
	t.Run("missing original proof refuses rather than fetching current definitions", func(t *testing.T) {
		f := newLongPipeline(t)
		f.saveSnapshot(func(s *workspace.WriteGroupAppliedSnapshot) { s.CounterBindings = nil })
		status := f.reconcileStatus()
		require.Equal(t, StateBlocked, status.State)
		require.Equal(t, "counter-replay-proof-missing", status.Reason)
		require.False(t, f.pipe.Owns("connector-1", longKwhFirst))
	})
	for _, change := range []string{"deleted", "nil", "outage"} {
		t.Run("valid frozen proof survives provider "+change, func(t *testing.T) {
			f := newLongPipeline(t)
			var provider CounterSources
			if change == "deleted" {
				provider = confirmedCounters{}
			} else if change == "outage" {
				provider = NewMeasurementCounterSources(failingDefinitions{})
			}
			f.pipe = New(Dependencies{Groups: f.groups, Destinations: f.dest, Store: f.store, Counters: provider}, Config{Now: func() time.Time { return f.clock }})
			require.NoError(t, f.pipe.Reconcile(t.Context()))
			require.Equal(t, StateActive, f.statusOf().State)
		})
	}
}

type memoryDefinitions map[string]*measurement.MeasurementDefinition

func (m memoryDefinitions) GetByID(_ context.Context, id string) (*measurement.MeasurementDefinition, error) {
	if definition, ok := m[id]; ok {
		return definition, nil
	}
	return nil, sql.ErrNoRows
}

func TestMeasurementCounterSourcesConfirmOnlyTheMembersOwnDefinitions(t *testing.T) {
	f := newLongPipeline(t)
	definition := func(id, device, point string) *measurement.MeasurementDefinition {
		return &measurement.MeasurementDefinition{
			ID: id, WorkspaceID: "workspace-A", DeviceID: device, PointID: point, DefinitionRevision: "def-9", SourceBindingRevision: "bind-9",
			SeriesEpoch: "epoch-7", Unit: "kWh", SemanticKind: measurement.SemanticKindCounter,
		}
	}
	sources, err := NewMeasurementCounterSources(memoryDefinitions{
		"meas-1": definition("meas-1", "device-1", "point-1"),
		"meas-2": definition("meas-2", "device-other", "point-2"), // belongs to another source
	}).CounterSources(t.Context(), f.group)
	require.NoError(t, err)
	require.Equal(t, runtime.CounterSource{
		MeasurementID: "meas-1", DefinitionRevision: "def-9", SourceBindingRevision: "bind-9", SeriesEpoch: "epoch-7",
		Unit: "kWh", Kind: "counter",
	}, sources[longKwhFirst])
	require.NotContains(t, sources, longKwhSecond, "a definition of another device is never borrowed")

	other := *f.group
	other.WorkspaceID = "workspace-B"
	sources, err = NewMeasurementCounterSources(memoryDefinitions{"meas-1": definition("meas-1", "device-1", "point-1")}).CounterSources(t.Context(), &other)
	require.NoError(t, err)
	require.Empty(t, sources, "nor one of another workspace")

	_, err = NewMeasurementCounterSources(failingDefinitions{}).CounterSources(t.Context(), f.group)
	require.Error(t, err)
}

func TestLongGroupBlocksOnDamagedDurableBaselines(t *testing.T) {
	f := newLongPipeline(t)
	f.closeFirstWindow()
	f.deliver()
	rewrite := func(payload string) {
		sum := sha256.Sum256([]byte(payload))
		_, err := f.plainLocal().ExecContext(t.Context(), `UPDATE wg_delivery_counter_baselines SET payload = ?, payload_digest = ? WHERE rowid = (SELECT MIN(rowid) FROM wg_delivery_counter_baselines)`,
			payload, hex.EncodeToString(sum[:]))
		require.NoError(t, err)
	}
	var payload string
	require.NoError(t, f.plainLocal().QueryRowContext(t.Context(), `SELECT payload FROM wg_delivery_counter_baselines ORDER BY rowid LIMIT 1`).Scan(&payload))

	rewrite("{")
	require.Equal(t, "counter-baseline-invalid", f.reconcileStatus().Reason, "damaged state blocks the group with a code an operator can act on")
	require.False(t, f.pipe.WantsSample("device-1", "point-1", longKwhFirst), "no sample is counted against a baseline that cannot be trusted")

	rewrite(`{"version":1}`)
	require.Equal(t, "counter-baseline-version-unsupported", f.reconcileStatus().Reason)

	rewrite(payload)
	require.Equal(t, StateActive, f.reconcileStatus().State, "intact state starts again")
}
