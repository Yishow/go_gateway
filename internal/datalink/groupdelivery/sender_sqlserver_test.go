package groupdelivery

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"sync/atomic"
	"testing"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/stretchr/testify/require"
)

func TestSQLServerSenderPassesProofAndMicrosecondTimestamp(t *testing.T) {
	fixture := newDeliveryFixture(t, plainTargetTable, true)
	table := "gw_long_fixture"
	effectKey := strings.Repeat("b", 64)
	bucket := enqueueCanonicalSQLServerLong(t, fixture.store, table, effectKey)

	var queries atomic.Int32
	targetDB := sql.OpenDB(sqlServerTestConnector{
		queries: &queries,
		onQuery: func(query string, args []driver.NamedValue) (driver.Rows, error) {
			return nil, errors.New("stop after reaching SQL Server query inside transaction")
		},
	})
	t.Cleanup(func() { _ = targetDB.Close() })

	fixture.resolver.target = Target{
		DB:        targetDB,
		Kind:      schema.DatabaseConnectorKindSQLServer,
		SQLServer: senderTestSQLServerProof(table),
	}

	result, err := fixture.sender.Deliver(t.Context(), bucket.Outcome.EffectKey)
	require.NoError(t, err)
	require.EqualValues(t, 1, queries.Load(), "a frozen proof and a SQL Server DATETIME2(6) receipt timestamp must reach the insert transaction")
	require.NotEqual(t, codeUnsupported, result.Code)
}

func TestSQLServerSenderMissingProofBlocked(t *testing.T) {
	fixture := newDeliveryFixture(t, plainTargetTable, true)
	table := "gw_long_fixture"
	effectKey := strings.Repeat("b", 64)
	bucket := enqueueCanonicalSQLServerLong(t, fixture.store, table, effectKey)

	var queries atomic.Int32
	targetDB := sql.OpenDB(sqlServerTestConnector{
		queries: &queries,
		onQuery: func(query string, args []driver.NamedValue) (driver.Rows, error) {
			queries.Add(1)
			return nil, errors.New("unexpected query")
		},
	})
	t.Cleanup(func() { _ = targetDB.Close() })

	fixture.resolver.target = Target{
		DB:        targetDB,
		Kind:      schema.DatabaseConnectorKindSQLServer,
		SQLServer: nil, // Missing proof
	}

	result, err := fixture.sender.Deliver(t.Context(), bucket.Outcome.EffectKey)
	require.NoError(t, err)
	require.Equal(t, StateBlocked, result.State)
	require.Equal(t, codeUnsupported, result.Code)
	require.EqualValues(t, 0, queries.Load(), "missing proof must block before reaching destination query")
}

func TestSQLServerSenderDeadlockRetriesWithUnchangedPayload(t *testing.T) {
	fixture := newDeliveryFixture(t, plainTargetTable, true)
	table := "gw_long_fixture"
	effectKey := strings.Repeat("b", 64)
	bucket := enqueueCanonicalSQLServerLong(t, fixture.store, table, effectKey)

	beforeItem, err := fixture.store.GetOutbox(t.Context(), effectKey)
	require.NoError(t, err)
	require.Equal(t, StatePending, beforeItem.State)

	var queries atomic.Int32
	targetDB := sql.OpenDB(sqlServerTestConnector{
		queries: &queries,
		onQuery: func(query string, args []driver.NamedValue) (driver.Rows, error) {
			queries.Add(1)
			return nil, sqlServerNumberedError{number: 1205, msg: "Transaction (Process ID 54) was deadlocked"}
		},
	})
	t.Cleanup(func() { _ = targetDB.Close() })

	fixture.resolver.target = Target{
		DB:        targetDB,
		Kind:      schema.DatabaseConnectorKindSQLServer,
		SQLServer: senderTestSQLServerProof(table),
	}

	result, err := fixture.sender.Deliver(t.Context(), bucket.Outcome.EffectKey)
	require.NoError(t, err)
	require.Equal(t, StateRetrying, result.State)
	require.Equal(t, codeInsertFailed, result.Code)

	afterItem, err := fixture.store.GetOutbox(t.Context(), effectKey)
	require.NoError(t, err)
	require.Equal(t, StateRetrying, afterItem.State)
	require.Equal(t, 1, afterItem.RetryCount)
	require.Equal(t, beforeItem.Payload, afterItem.Payload, "frozen payload must remain unchanged on deadlock retry")
	require.Equal(t, beforeItem.PayloadDigest, afterItem.PayloadDigest, "frozen payload digest must remain unchanged on deadlock retry")
	require.Equal(t, beforeItem.Destination, afterItem.Destination)
}

func TestSQLServerSenderConflictIsolation(t *testing.T) {
	t.Run("digest_mismatch_blocks_affected_item", func(t *testing.T) {
		fixture := newDeliveryFixture(t, plainTargetTable, true)
		table := "gw_long_fixture"
		effectKey := strings.Repeat("b", 64)
		bucket := enqueueCanonicalSQLServerLong(t, fixture.store, table, effectKey)

		proof := senderTestSQLServerProof(table)
		targetDB := sql.OpenDB(sqlServerTestConnector{
			catalog: newSQLServerCatalogMock(proof.Group, proof.Receipt),
			onQuery: func(query string, args []driver.NamedValue) (driver.Rows, error) {
				if strings.Contains(query, "WITH (UPDLOCK, HOLDLOCK)") && strings.Contains(query, "gw_effect_receipts") {
					foreignDigest := strings.Repeat("f", 64)
					return &sqlServerMockRows{
						columns: []string{"payload_digest", "_gw_effect_receipts_v1"},
						rows:    [][]driver.Value{{foreignDigest, "gateway"}},
					}, nil
				}
				if strings.Contains(query, "WITH (UPDLOCK, HOLDLOCK)") && strings.Contains(query, table) {
					cells := make([]driver.Value, len(bucket.Row.Cells))
					for i, cell := range bucket.Row.Cells {
						cells[i] = cell.Value
					}
					cols := make([]string, len(bucket.Row.Cells))
					for i, cell := range bucket.Row.Cells {
						cols[i] = cell.Column
					}
					return &sqlServerMockRows{
						columns: cols,
						rows:    [][]driver.Value{cells},
					}, nil
				}
				return nil, nil
			},
		})
		t.Cleanup(func() { _ = targetDB.Close() })

		fixture.resolver.target = Target{
			DB:        targetDB,
			Kind:      schema.DatabaseConnectorKindSQLServer,
			SQLServer: proof,
		}

		result, err := fixture.sender.Deliver(t.Context(), bucket.Outcome.EffectKey)
		require.NoError(t, err)
		require.Equal(t, StateBlocked, result.State)
		require.Equal(t, codeIdentityConflict, result.Code)

		item, err := fixture.store.GetOutbox(t.Context(), effectKey)
		require.NoError(t, err)
		require.Equal(t, StateBlocked, item.State)
		require.Equal(t, codeIdentityConflict, item.LastErrorCode)
	})

	t.Run("duplicate_row_quarantines_affected_item", func(t *testing.T) {
		fixture := newDeliveryFixture(t, plainTargetTable, true)
		table := "gw_long_fixture"
		effectKey := strings.Repeat("b", 64)
		bucket := enqueueCanonicalSQLServerLong(t, fixture.store, table, effectKey)

		targetDB := sql.OpenDB(sqlServerTestConnector{
			onQuery: func(query string, args []driver.NamedValue) (driver.Rows, error) {
				return nil, sqlServerNumberedError{number: 2627, msg: "Violation of PRIMARY KEY constraint"}
			},
		})
		t.Cleanup(func() { _ = targetDB.Close() })

		fixture.resolver.target = Target{
			DB:        targetDB,
			Kind:      schema.DatabaseConnectorKindSQLServer,
			SQLServer: senderTestSQLServerProof(table),
		}

		result, err := fixture.sender.Deliver(t.Context(), bucket.Outcome.EffectKey)
		require.NoError(t, err)
		require.Equal(t, StateQuarantined, result.State)
		require.Equal(t, codeRowRejected, result.Code)

		item, err := fixture.store.GetOutbox(t.Context(), effectKey)
		require.NoError(t, err)
		require.Equal(t, StateQuarantined, item.State)
		require.Equal(t, codeRowRejected, item.LastErrorCode)
	})
}

func TestSQLServerSenderCommitAckLostReceiptRetryOnce(t *testing.T) {
	fixture := newDeliveryFixture(t, plainTargetTable, true)
	table := "gw_long_fixture"
	effectKey := strings.Repeat("b", 64)
	bucket := enqueueCanonicalSQLServerLong(t, fixture.store, table, effectKey)

	beforeItem, err := fixture.store.GetOutbox(t.Context(), effectKey)
	require.NoError(t, err)

	var committed atomic.Bool
	var commits atomic.Int32
	var inserts atomic.Int32

	proof := senderTestSQLServerProof(table)
	targetDB := sql.OpenDB(sqlServerTestConnector{
		catalog: newSQLServerCatalogMock(proof.Group, proof.Receipt),
		onExec: func(_ string, _ []driver.NamedValue) (driver.Result, error) {
			inserts.Add(1)
			return driver.RowsAffected(1), nil
		},
		onCommit: func() error {
			commits.Add(1)
			if commits.Load() == 1 {
				committed.Store(true)
				return lostcommit.ErrResponseLost
			}
			return nil
		},
		onQuery: func(query string, args []driver.NamedValue) (driver.Rows, error) {
			if strings.Contains(query, "WITH (UPDLOCK, HOLDLOCK)") && strings.Contains(query, "gw_effect_receipts") {
				if committed.Load() {
					return &sqlServerMockRows{
						columns: []string{"payload_digest", "_gw_effect_receipts_v1"},
						rows:    [][]driver.Value{{beforeItem.PayloadDigest, "gateway"}},
					}, nil
				}
				return nil, sql.ErrNoRows
			}
			if strings.Contains(query, "WITH (UPDLOCK, HOLDLOCK)") && strings.Contains(query, table) {
				if committed.Load() {
					cells := make([]driver.Value, len(bucket.Row.Cells))
					for i, cell := range bucket.Row.Cells {
						cells[i] = cell.Value
					}
					cols := make([]string, len(bucket.Row.Cells))
					for i, cell := range bucket.Row.Cells {
						cols[i] = cell.Column
					}
					return &sqlServerMockRows{
						columns: cols,
						rows:    [][]driver.Value{cells},
					}, nil
				}
				return nil, sql.ErrNoRows
			}
			return nil, nil
		},
	})
	t.Cleanup(func() { _ = targetDB.Close() })

	fixture.resolver.target = Target{
		DB:        targetDB,
		Kind:      schema.DatabaseConnectorKindSQLServer,
		SQLServer: proof,
	}

	// First attempt: commit succeeded on destination, but client ACK lost
	firstResult, err := fixture.sender.Deliver(t.Context(), bucket.Outcome.EffectKey)
	require.NoError(t, err)
	require.Equal(t, StateRetrying, firstResult.State)
	require.Equal(t, codeCommitAmbiguous, firstResult.Code)

	item, err := fixture.store.GetOutbox(t.Context(), effectKey)
	require.NoError(t, err)
	require.Equal(t, StateRetrying, item.State)
	require.Equal(t, codeCommitAmbiguous, item.LastErrorCode)
	require.Zero(t, countRows(t, fixture.local, "wg_delivery_receipts"))

	// Second attempt: deduped by destination receipt, completes once
	secondResult, err := fixture.sender.Deliver(t.Context(), bucket.Outcome.EffectKey)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, secondResult.State)

	item2, err := fixture.store.GetOutbox(t.Context(), effectKey)
	require.NoError(t, err)
	require.Equal(t, StateCommitted, item2.State)
	require.Equal(t, 1, countRows(t, fixture.local, "wg_delivery_receipts"))
	require.EqualValues(t, 1, commits.Load(), "the receipt replay never commits another effect")
	require.EqualValues(t, 2, inserts.Load(), "exactly one row and one receipt insert, both before the lost ACK")
	require.Equal(t, beforeItem.Payload, item2.Payload)
	require.Equal(t, beforeItem.PayloadDigest, item2.PayloadDigest)
	require.Equal(t, beforeItem.Destination, item2.Destination)
}
