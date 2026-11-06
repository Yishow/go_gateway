package dbtarget

import (
	"context"
	"errors"
	"net"
	"syscall"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestWriter_FailedOutcomeProvidesSafeActionableReason(t *testing.T) {
	cases := []struct {
		name, want, secret, column string
		failure                    error
	}{
		{"permission", "權限", "fake-token-permission", "d0", &pq.Error{Code: "42501", Message: "permission denied fake-token-permission"}},
		{"connection", "連線", "192.0.2.42", "d0", &net.OpError{Op: "dial", Net: "tcp", Addr: &net.TCPAddr{IP: net.ParseIP("192.0.2.42"), Port: 5432}, Err: syscall.ECONNREFUSED}},
		{"schema", "結構", "secret_column_6789", "secret_column_6789", nil},
		{"unknown", "資料庫寫入失敗", "fake-token-unknown", "d0", errors.New("unclassified fake-token-unknown")},
	}
	for _, tc := range cases {
		for _, grouped := range []bool{false, true} {
			name := tc.name
			if grouped {
				name += "-grouped"
			}
			t.Run(name, func(t *testing.T) {
				connectors, mappings, connector, _ := backoffTestTarget(t)
				mapping := &schema.DatabaseTargetMapping{
					ID: "reason-mapping", TagID: "reason-tag", ConnectorID: connector.ID,
					TableName: "snapshot_values", ColumnName: tc.column, WriteMode: schema.DatabaseWriteModeInsert,
					Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
				}
				if grouped {
					group := "meter"
					mapping.GroupKey = &group
				}
				require.NoError(t, mappings.Create(t.Context(), mapping))
				sample := time.Date(2026, 6, 7, 10, 0, 1, 0, time.UTC)
				clock := &backoffTestClock{now: sample}
				tick := make(chan time.Time, 1)
				writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
					FlushTick: tick, Now: clock.Now,
					OpenDestination: func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
						if tc.failure != nil {
							return nil, tc.failure
						}
						return openExternalDBManagerFunc(kind, config)
					},
				})
				t.Cleanup(func() { _ = writer.Close(context.Background()) })
				require.ErrorIs(t, writer.WriteTagValue(t.Context(), "reason-tag", 12, sample), ErrSnapshotAccepted)
				clock.Set(sample.Add(time.Minute))
				tick <- clock.Now()
				service := NewConnectorService(connectors, mappings)
				var outcome *schema.DatabaseConnector
				var history []WriteHistoryRecord
				require.Eventually(t, func() bool {
					var err error
					outcome, err = connectors.GetByID(t.Context(), connector.ID)
					if err != nil || outcome.LastWriteStatus != "failed" {
						return false
					}
					history, err = service.ListWriteHistory(t.Context(), connector.ID, 200)
					if err != nil {
						return false
					}
					for _, record := range history {
						if record.Status == "failed" {
							return true
						}
					}
					return false
				}, 2*time.Second, 10*time.Millisecond)
				require.Contains(t, outcome.LastWriteError, tc.want)
				require.Contains(t, outcome.LastWriteError, "1 秒後重試")
				require.NotContains(t, outcome.LastWriteError, tc.secret)
				if grouped {
					require.NotContains(t, outcome.LastFlushError, tc.secret)
				}
				for _, record := range history {
					if record.Status == "failed" {
						require.Contains(t, record.ErrorSummary, tc.secret, "history retains the diagnostic detail")
						return
					}
				}
				t.Fatal("failure history was not recorded")
			})
		}
	}
}

func TestWriter_CloseFailureKeepsOperatorOutcomeSafe(t *testing.T) {
	for _, grouped := range []bool{false, true} {
		name := "ungrouped"
		if grouped {
			name = "grouped"
		}
		t.Run(name, func(t *testing.T) {
			connectors, mappings, connector, _ := backoffTestTarget(t)
			mapping := &schema.DatabaseTargetMapping{
				ID: "close-reason-mapping", TagID: "close-reason-tag", ConnectorID: connector.ID,
				TableName: "snapshot_values", ColumnName: "d0", WriteMode: schema.DatabaseWriteModeInsert,
				Enabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}
			if grouped {
				group := "meter"
				mapping.GroupKey = &group
			}
			require.NoError(t, mappings.Create(t.Context(), mapping))
			const secret = "fake-password-close"
			writer := NewWriterWithConfig(connectors, mappings, WriterConfig{
				FlushTick: make(chan time.Time),
				OpenDestination: func(schema.DatabaseConnectorKind, ConnectionConfig) (*datalinkbase.DBManager, error) {
					return nil, errors.New("unclassified " + secret)
				},
			})
			t.Cleanup(func() { _ = writer.Close(context.Background()) })
			require.ErrorIs(t, writer.WriteTagValue(t.Context(), "close-reason-tag", 12, time.Now()), ErrSnapshotAccepted)
			require.Error(t, writer.Close(t.Context()))
			outcome, err := connectors.GetByID(t.Context(), connector.ID)
			require.NoError(t, err)
			require.Equal(t, "failed", outcome.LastWriteStatus)
			require.Contains(t, outcome.LastWriteError, "資料庫寫入失敗")
			require.NotContains(t, outcome.LastWriteError, secret)
			require.NotContains(t, outcome.LastFlushError, secret)
			history, err := NewConnectorService(connectors, mappings).ListWriteHistory(t.Context(), connector.ID, 200)
			require.NoError(t, err)
			require.NotEmpty(t, history)
			require.Contains(t, history[0].ErrorSummary, secret, "history retains the diagnostic detail")
		})
	}
}
