package dbtarget

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriter_RecordsWhyBufferedMappingNoLongerApplies(t *testing.T) {
	for _, grouped := range []bool{false, true} {
		for _, reason := range []string{"mapping_not_current", "write_group_owned"} {
			name := reason
			if grouped {
				name += "_grouped"
			}
			t.Run(name, func(t *testing.T) {
				connectors, mappings, connector, target := backoffTestTarget(t)
				entry := &schema.DatabaseTargetMapping{
					ID: "obsolete-mapping", TagID: "obsolete-tag", ConnectorID: connector.ID,
					TableName: "snapshot_values", ColumnName: "d0", WriteMode: schema.DatabaseWriteModeInsert,
					Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
				}
				if grouped {
					group := "meter"
					entry.GroupKey = &group
				}
				require.NoError(t, mappings.Create(t.Context(), entry))
				var owned atomic.Bool
				writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
					FlushTick: make(chan time.Time), SuppressMapping: func(string, string) bool { return owned.Load() },
				})
				t.Cleanup(func() { _ = writer.Close(context.Background()) })
				at := time.Date(2026, 11, 5, 10, 0, 1, 0, time.UTC)
				require.ErrorIs(t, writer.WriteTagValue(t.Context(), entry.TagID, 12, at), ErrSnapshotAccepted)
				if reason == "mapping_not_current" {
					entry.Enabled = false
					require.NoError(t, mappings.Update(t.Context(), entry))
					require.NoError(t, writer.RefreshMappings(t.Context()))
				} else {
					owned.Store(true)
				}
				require.NoError(t, writer.flushDue(t.Context(), at.Add(time.Minute), false))
				history, err := NewConnectorService(connectors, mappings).ListWriteHistory(t.Context(), connector.ID, 20)
				require.NoError(t, err)
				require.Len(t, history, 1)
				require.Equal(t, "skipped", history[0].Status)
				require.Equal(t, entry.ID, history[0].MappingID)
				require.Equal(t, reason, history[0].SkipReason)
				var rows int
				require.NoError(t, target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM snapshot_values`).Scan(&rows))
				require.Zero(t, rows)
			})
		}
	}
}

func TestWriter_RecordsFailedSnapshotSupersededByNewConnectorRevision(t *testing.T) {
	checkFailedSnapshotAfterConnectorChange(t, false)
}

func TestWriter_RecordsFailedGroupedSnapshotSupersededByNewConnectorRevision(t *testing.T) {
	checkFailedSnapshotAfterConnectorChange(t, true)
}

func checkFailedSnapshotAfterConnectorChange(t *testing.T, grouped bool) {
	t.Helper()
	connectors, mappings, connector, target := backoffTestTarget(t)
	if grouped {
		group := "meter"
		require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
			ID: "mapping-inflight", TagID: "inflight", ConnectorID: connector.ID,
			TableName: "snapshot_values", ColumnName: "d0", GroupKey: &group,
			WriteMode: schema.DatabaseWriteModeInsert, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}))
	} else {
		createSnapshotMapping(t, mappings, connector.ID, "inflight", "d0")
	}
	interval := 1
	require.NoError(t, mappings.Create(t.Context(), &schema.DatabaseTargetMapping{
		ID: "sentinel-idle", TagID: "sentinel-idle", ConnectorID: connector.ID,
		TableName: "snapshot_values", ColumnName: "d1", WriteMode: schema.DatabaseWriteModeInsert,
		WriteIntervalSeconds: &interval, Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}))
	base := time.Date(2026, 11, 5, 10, 0, 0, 0, time.UTC)
	clock := &backoffTestClock{now: base}
	tick := make(chan time.Time, 1)
	entered, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	releaseOpen := func() { releaseOnce.Do(func() { close(release) }) }
	var attempts atomic.Int64
	writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: tick, Now: clock.Now,
		OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
			if attempts.Add(1) == 1 {
				close(entered)
				<-release
				return nil, errors.New("destination unavailable")
			}
			return openExternalDBManagerFunc(kind, config)
		},
	})
	t.Cleanup(func() { releaseOpen(); _ = writer.Close(context.Background()) })
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "inflight", 10, base.Add(time.Second)), ErrSnapshotAccepted)
	clock.Set(base.Add(15 * time.Second))
	tick <- clock.Now()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("first destination attempt did not start")
	}
	previous := connector.IdentityRevision
	connector.IdentityRevision = "new-connector-revision"
	require.NoError(t, connectors.UpdateWithExpectedIdentityRevision(t.Context(), connector, previous))
	require.NoError(t, writer.RefreshMappings(t.Context()))
	require.ErrorIs(t, writer.WriteTagValue(t.Context(), "inflight", 12, base.Add(15200*time.Millisecond)), ErrSnapshotAccepted)
	releaseOpen()
	// An idle-interval entry is recorded after the failed worker iteration has
	// completed, so it is a deterministic barrier for the restoration decision.
	require.Eventually(t, func() bool {
		for _, record := range listWriteHistory(connector.ID, 100) {
			if record.MappingID == "sentinel-idle" && record.SkipReason == "no_new_value" {
				return true
			}
		}
		return false
	}, 3*time.Second, 10*time.Millisecond)
	clock.Set(base.Add(16 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool {
		for _, record := range listWriteHistory(connector.ID, 100) {
			if record.MappingID == "mapping-inflight" && record.SkipReason == "mapping_not_current" {
				return true
			}
		}
		return false
	}, 3*time.Second, 10*time.Millisecond, "the failed old revision must not disappear silently")
	clock.Set(base.Add(30 * time.Second))
	tick <- clock.Now()
	require.Eventually(t, func() bool {
		var value float64
		return target.QueryRowContext(t.Context(), `SELECT d0 FROM snapshot_values`).Scan(&value) == nil && value == 12
	}, 3*time.Second, 10*time.Millisecond)
}

func TestWriter_RecordsSupersededSnapshotAfterTargetChanges(t *testing.T) {
	for _, grouped := range []bool{false, true} {
		for _, change := range []string{"mapping", "connector"} {
			name := change
			if grouped {
				name += "_grouped"
			}
			t.Run(name, func(t *testing.T) {
				connectors, mappings, connector, target := backoffTestTarget(t)
				entry := &schema.DatabaseTargetMapping{
					ID: "changed-mapping", TagID: "changed-tag", ConnectorID: connector.ID,
					TableName: "snapshot_values", ColumnName: "d0", WriteMode: schema.DatabaseWriteModeInsert,
					Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
				}
				if grouped {
					group := "meter"
					entry.GroupKey = &group
				}
				require.NoError(t, mappings.Create(t.Context(), entry))
				writer := NewWriterWithConfig(connectors, mappings, WriterConfig{FlushTick: make(chan time.Time)})
				t.Cleanup(func() { _ = writer.Close(context.Background()) })
				at := time.Date(2026, 11, 5, 10, 0, 1, 0, time.UTC)
				require.ErrorIs(t, writer.WriteTagValue(t.Context(), entry.TagID, 10, at), ErrSnapshotAccepted)
				if change == "mapping" {
					entry.ColumnName = "d1"
					require.NoError(t, mappings.Update(t.Context(), entry))
				} else {
					previous := connector.IdentityRevision
					connector.IdentityRevision = "new-connector-revision"
					require.NoError(t, connectors.UpdateWithExpectedIdentityRevision(t.Context(), connector, previous))
				}
				require.NoError(t, writer.RefreshMappings(t.Context()))
				require.ErrorIs(t, writer.WriteTagValue(t.Context(), entry.TagID, 12, at.Add(time.Second)), ErrSnapshotAccepted)
				require.NoError(t, writer.flushDue(t.Context(), at.Add(time.Minute), false))
				history, err := NewConnectorService(connectors, mappings).ListWriteHistory(t.Context(), connector.ID, 20)
				require.NoError(t, err)
				var skipped, succeeded int
				for _, record := range history {
					if record.Status == "skipped" && record.MappingID == entry.ID && record.SkipReason == "mapping_not_current" {
						skipped++
					}
					if record.Status == "success" {
						succeeded++
					}
				}
				require.Equal(t, 1, skipped, "the old buffered snapshot must not disappear silently")
				require.Equal(t, 1, succeeded, "the new target identity must receive the new value")
				var value float64
				query := `SELECT d0 FROM snapshot_values`
				if change == "mapping" {
					query = `SELECT d1 FROM snapshot_values`
				}
				require.NoError(t, target.QueryRowContext(t.Context(), query).Scan(&value))
				require.Equal(t, 12.0, value)
			})
		}
	}
}
