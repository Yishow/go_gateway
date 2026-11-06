package dbtarget

import (
	"context"
	"errors"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestWriter_LastWriteTimeReflectsDelayedAttempt(t *testing.T) {
	for _, grouped := range []bool{false, true} {
		for _, fail := range []bool{false, true} {
			name := "ungrouped"
			if grouped {
				name = "grouped"
			}
			if fail {
				name += "-failed"
			} else {
				name += "-success"
			}
			t.Run(name, func(t *testing.T) {
				connectors, mappings, connector, _ := backoffTestTarget(t)
				mapping := &schema.DatabaseTargetMapping{
					ID: "delayed-mapping", TagID: "delayed-tag", ConnectorID: connector.ID,
					TableName: "snapshot_values", ColumnName: "d0", WriteMode: schema.DatabaseWriteModeInsert,
					Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
				}
				if grouped {
					group := "meter"
					mapping.GroupKey = &group
				}
				require.NoError(t, mappings.Create(t.Context(), mapping))
				sample := time.Date(2026, 6, 7, 10, 0, 1, 0, time.UTC)
				attempt := time.Date(2026, 6, 7, 10, 3, 0, 0, time.UTC)
				clock := &backoffTestClock{now: sample}
				tick := make(chan time.Time, 1)
				writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
					FlushTick: tick, Now: clock.Now,
					OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
						if fail {
							return nil, errors.New("permission denied")
						}
						return openExternalDBManagerFunc(kind, config)
					},
				})
				t.Cleanup(func() { _ = writer.Close(context.Background()) })
				require.ErrorIs(t, writer.WriteTagValue(t.Context(), "delayed-tag", 12, sample), ErrSnapshotAccepted)
				clock.Set(attempt)
				tick <- attempt
				status := "success"
				if fail {
					status = "failed"
				}
				var got *schema.DatabaseConnector
				service := NewConnectorService(connectors, mappings)
				require.Eventually(t, func() bool {
					var err error
					got, err = connectors.GetByID(t.Context(), connector.ID)
					if err != nil || got.LastWriteStatus != status {
						return false
					}
					history, err := service.ListWriteHistory(t.Context(), connector.ID, 200)
					if err != nil {
						return false
					}
					for _, record := range history {
						if record.Status == status {
							return true
						}
					}
					return false
				}, 2*time.Second, 10*time.Millisecond)
				require.NotNil(t, got.LastWriteAt)
				require.Equal(t, attempt, got.LastWriteAt.UTC(), "last write must describe the attempt, not the sample")
				if grouped {
					require.NotNil(t, got.LastFlushAt)
					require.Equal(t, sample.Truncate(15*time.Second), got.LastFlushAt.UTC())
				}
				history, err := service.ListWriteHistory(t.Context(), connector.ID, 200)
				require.NoError(t, err)
				for _, record := range history {
					if record.Status != status {
						continue
					}
					observed := sample
					if grouped {
						observed = sample.Truncate(15 * time.Second)
					}
					require.Equal(t, observed, record.ObservedAt)
					require.NotNil(t, record.AttemptedAt)
					require.Equal(t, attempt, record.AttemptedAt.UTC())
					return
				}
				t.Fatal("write history is missing the attempt")
			})
		}
	}
}
