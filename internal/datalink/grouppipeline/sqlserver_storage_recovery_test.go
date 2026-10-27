package grouppipeline

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/require"
)

func TestSQLServerStorageRecoverySelectsSQLDialect(t *testing.T) {
	dialect, supported := dialectName(schema.DatabaseConnectorKindSQLServer)
	require.True(t, supported, "SQL Server schema errors must reach the confirmed repair ledger")
	require.Equal(t, "sqlserver", dialect)
}

func TestSQLServerMissingRepairPreservesFrozenOutboxAndReceipt(t *testing.T) {
	for _, lostACK := range []bool{false, true} {
		t.Run(map[bool]string{false: "committed", true: "ack lost"}[lostACK], func(t *testing.T) {
			f := newSQLServerRepairFixture(t)
			if lostACK {
				f.executeErr = errSQLServerFixtureAckLost
			}
			original, err := f.store.GetOutbox(t.Context(), f.item.EffectKey)
			require.NoError(t, err)
			decision, err := f.recovery.Restore(t.Context(), f.item, f.target, mssql.Error{Number: 208})
			require.NoError(t, err)
			require.Equal(t, groupdelivery.RestoreRetryNow, decision.Action)
			require.Equal(t, 1, f.executions)
			require.Equal(t, int64(1001), f.target.SQLServer.Group.Evidence.ObjectID)
			require.Equal(t, f.requests[0].OperationID, f.target.SQLServer.Group.Evidence.RepairOperationID)
			require.Equal(t, int64(101), f.snapshot.RuntimeLayout.SQLServer.Group.Evidence.ObjectID, "the accepted runtime proof remains original")
			require.Empty(t, f.snapshot.RuntimeLayout.SQLServer.Group.Evidence.RepairOperationID)
			require.Equal(t, f.snapshot.RuntimeLayout.SQLServer.Receipt, f.target.SQLServer.Receipt)
			lineage := f.lineage()
			require.Equal(t, 2, lineage.Generation)
			require.True(t, lineage.HistoryMayBeLost)
			require.Equal(t, 1, lineage.Reconstructions)
			require.Equal(t, f.snapshot.Group.Destination.SchemaDigest, lineage.ConfirmedDigest)
			after, err := f.store.GetOutbox(t.Context(), f.item.EffectKey)
			require.NoError(t, err)
			require.Equal(t, original, after, "schema repair itself never rewrites, settles or recalculates pending data")
		})
	}
}

func TestSQLServerRepairRefreshesOnlyAuthorizedGenerationAfterRestart(t *testing.T) {
	f := newSQLServerRepairFixture(t)
	decision, err := f.recovery.Restore(t.Context(), f.item, f.target, mssql.Error{Number: 208})
	require.NoError(t, err)
	require.Equal(t, groupdelivery.RestoreRetryNow, decision.Action)
	restarted := NewStorageRecovery(f.store, f.recovery.ledger, nil)
	restarted.sqlServerOps = f.recovery.sqlServerOps
	proof, err := restarted.sqlServerCurrentProof(t.Context(), f.snapshot, f.target, true)
	require.NoError(t, err)
	require.Equal(t, int64(1001), proof.Group.Evidence.ObjectID)
	require.Equal(t, 1, f.executions, "reading the authorized generation after restart runs no DDL")
	f.missing = true
	f.target.SQLServer = proof
	decision, err = restarted.Restore(t.Context(), f.item, f.target, mssql.Error{Number: 208})
	require.NoError(t, err)
	require.Equal(t, groupdelivery.RestoreRetryNow, decision.Action)
	require.Equal(t, 2, f.executions)
	require.Equal(t, int64(1002), f.target.SQLServer.Group.Evidence.ObjectID)
	require.NotEqual(t, f.requests[0].OperationID, f.requests[1].OperationID)
	require.Equal(t, 3, f.lineage().Generation)
	require.Equal(t, 2, f.lineage().Reconstructions)

	for _, scenario := range []string{"replaced group", "replaced receipt", "changed marker"} {
		t.Run(scenario, func(t *testing.T) {
			before := dbtarget.CloneSQLServerManagedRowProof(f.current)
			t.Cleanup(func() { f.current = before })
			switch scenario {
			case "replaced group":
				f.current.Group.Evidence.ObjectID++
			case "replaced receipt":
				f.current.Receipt.Evidence.ObjectID++
			case "changed marker":
				f.current.Group.Evidence.RepairOperationID = "foreign-operation"
			}
			_, err := restarted.sqlServerCurrentProof(t.Context(), f.snapshot, f.target, true)
			require.Error(t, err)
			require.Equal(t, 2, f.executions)
		})
	}
}

func TestSQLServerRepairScopeDeniedBeforeAnyDDL(t *testing.T) {
	for _, policy := range []string{"absent", "disabled"} {
		t.Run(policy, func(t *testing.T) {
			f := newSQLServerRepairFixtureWith(t, func(s *workspace.WriteGroupAppliedSnapshot) {
				s.Group.WritePolicy.RecreateMissing = nil
				if policy == "disabled" {
					disabled := false
					s.Group.WritePolicy.RecreateMissing = &disabled
				}
			})
			decision, err := f.recovery.Restore(t.Context(), f.item, f.target, mssql.Error{Number: 208})
			require.NoError(t, err)
			require.Equal(t, groupdelivery.RestoreBlock, decision.Action)
			require.Zero(t, f.executions)
		})
	}
	for _, scenario := range []string{"gate closed", "edited connector", "foreign table", "receipt missing", "metadata forbidden", "frozen alias"} {
		t.Run(scenario, func(t *testing.T) {
			f := newSQLServerRepairFixture(t)
			switch scenario {
			case "gate closed":
				f.target.SQLServerCapability.SupportsMissingRepair = false
			case "edited connector":
				f.item.Destination.ConnectorRevision = "edited"
			case "foreign table":
				f.missing = false
				f.current.Group.Evidence.ObjectID = 999
			case "receipt missing":
				f.receiptMissing = true
			case "metadata forbidden":
				f.inspectErr = dbtarget.ErrSQLServerManagedTableForbidden
			case "frozen alias":
				f.recovery.snapshots.Store(f.item.Key, f.snapshot)
				f.target.SQLServer = f.snapshot.RuntimeLayout.SQLServer
			}
			decision, err := f.recovery.Restore(t.Context(), f.item, f.target, mssql.Error{Number: 208})
			require.NoError(t, err)
			require.Equal(t, groupdelivery.RestoreBlock, decision.Action)
			require.Zero(t, f.executions)
		})
	}
}

func TestSQLServerResolverGenerationFailureClosesOrRetries(t *testing.T) {
	for _, transient := range []bool{false, true} {
		snap := sqlServerRuntimeSnapshot(t)
		resolver, item, opened := sqlServerResolverFixture(t, snap)
		resolver.loadSQLServerProof = func(context.Context, *workspace.WriteGroupAppliedSnapshot, *dbtarget.OpenedDestination) (*dbtarget.SQLServerManagedRowProof, error) {
			if transient {
				return nil, errSQLServerProofReadTransient
			}
			return nil, recordingplan.ErrPreviewTokenStale
		}
		_, err := resolver.Resolve(t.Context(), item)
		require.Error(t, err)
		require.Equal(t, !transient, errors.Is(err, groupdelivery.ErrTargetBlocked))
		require.Equal(t, 1, opened.closeCalls)
	}
}

func TestSQLServerStaleRepairUsesOpenedScopeWithoutRepeatingDDL(t *testing.T) {
	f := newSQLServerRepairFixture(t)
	_, err := f.store.DB().ExecContext(t.Context(), `CREATE TRIGGER reject_sql_repair_finish BEFORE UPDATE OF status
		ON managed_schema_operations WHEN NEW.action = 'storage_repair' AND NEW.status = 'succeeded'
		BEGIN SELECT RAISE(FAIL, 'scripted local ledger acknowledgement failure'); END`)
	require.NoError(t, err)
	_, err = f.recovery.Restore(t.Context(), f.item, f.target, mssql.Error{Number: 208})
	require.ErrorIs(t, err, recordingplan.ErrSchemaOperationUnacknowledged)
	require.Equal(t, 1, f.executions)
	require.Equal(t, 1, f.lineage().Generation, "failed local operation completion also rolls back generation")
	_, err = f.store.DB().ExecContext(t.Context(), `DROP TRIGGER reject_sql_repair_finish`)
	require.NoError(t, err)
	_, err = f.store.DB().ExecContext(t.Context(), `UPDATE managed_schema_operations SET updated_at = ? WHERE action = 'storage_repair'`,
		time.Now().UTC().Add(-2*recordingplan.SchemaOperationLease))
	require.NoError(t, err)
	f.recovery.destinations = &mysqlResolverDestinations{opened: &dbtarget.OpenedDestination{
		DB: f.target.DB, Kind: schema.DatabaseConnectorKindSQLServer, SQLServer: f.target.SQLServerCapability,
	}}
	require.NoError(t, f.recovery.ResolveStale(t.Context()))
	require.Equal(t, 1, f.executions, "restart reconciliation inspects only; it never repeats DDL")
	require.Equal(t, 2, f.lineage().Generation)
	require.Equal(t, 1, f.lineage().Reconstructions)
	proof, err := f.recovery.sqlServerCurrentProof(t.Context(), f.snapshot, f.target, true)
	require.NoError(t, err)
	require.Equal(t, int64(1001), proof.Group.Evidence.ObjectID)
	require.Equal(t, f.requests[0].OperationID, proof.Group.Evidence.RepairOperationID)
}

func TestSQLServerConcurrentRepairUsesSingleOperationAcrossCoordinators(t *testing.T) {
	f := newSQLServerRepairFixture(t)
	f.store.DB().SetMaxOpenConns(1)
	started, finish := make(chan struct{}), make(chan struct{})
	var startOnce, finishOnce sync.Once
	release := func() { finishOnce.Do(func() { close(finish) }) }
	defer release()
	ops := *f.recovery.sqlServerOps
	ops.Execute = func(ctx context.Context, db *sql.DB, capability *dbtarget.SQLServerRecordingCapability, req dbtarget.SQLServerManagedRepairRequest) (*dbtarget.SchemaStatementExecution, error) {
		startOnce.Do(func() { close(started) })
		select {
		case <-finish:
			return f.execute(ctx, db, capability, req)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	type outcome struct {
		decision groupdelivery.RestoreDecision
		err      error
	}
	results := make(chan outcome, 8)
	var workers sync.WaitGroup
	t.Cleanup(workers.Wait)
	for range 8 {
		workers.Go(func() {
			// Separate recovery instances have separate process-local scope locks;
			// only the real local operation ledger coordinates these attempts.
			recovery := NewStorageRecovery(f.store, f.recovery.ledger, nil)
			recovery.sqlServerOps = &ops
			target := f.target
			target.SQLServer = dbtarget.CloneSQLServerManagedRowProof(f.snapshot.RuntimeLayout.SQLServer)
			decision, err := recovery.Restore(t.Context(), f.item, target, mssql.Error{Number: 208})
			results <- outcome{decision, err}
		})
	}
	select {
	case <-started:
	case <-time.After(10 * time.Second):
		t.Fatal("no coordinator reached the claimed DDL")
	}
	for range 7 {
		select {
		case result := <-results:
			require.NoError(t, result.err)
			require.Equal(t, groupdelivery.RestoreWait, result.decision.Action)
		case <-time.After(10 * time.Second):
			t.Fatal("another coordinator did not observe the in-progress operation")
		}
	}
	release()
	workers.Wait()
	result := <-results
	require.NoError(t, result.err)
	require.Equal(t, groupdelivery.RestoreRetryNow, result.decision.Action)
	require.Equal(t, 1, f.executions)
	require.Equal(t, 2, f.lineage().Generation)
	require.Equal(t, 1, f.lineage().Reconstructions)
}

func TestSQLServerCurrentProofAllowsMultipleGroupsSharingConfirmedScope(t *testing.T) {
	f := newSQLServerRepairFixture(t)
	scope := workspace.StorageScopeOf(f.snapshot.Group)
	repo := recordingplan.NewSQLRepository(f.store.DB())

	var scopeKey string
	require.NoError(t, f.store.DB().QueryRowContext(t.Context(),
		"SELECT scope_key FROM managed_schema_operations WHERE operation_id = ?", f.snapshot.Group.Destination.SchemaRevision).Scan(&scopeKey))

	// Simulate a second write group confirming the exact same scope and layout afterwards.
	secondOpID := "op-second-group-confirmation"
	now := time.Now().UTC()
	_, _, err := repo.ClaimSchemaOperation(t.Context(), &recordingplan.SchemaOperation{
		OperationID: secondOpID, Token: "second-token", WorkspaceID: f.snapshot.WorkspaceID,
		ScopeKey: scopeKey, Owner: "claim-second", Action: "schema_apply",
		Status: recordingplan.SchemaOperationRunning, CreatedAt: now, UpdatedAt: now,
	})
	require.NoError(t, err)
	_, err = repo.FinishSchemaOperation(t.Context(), secondOpID, "claim-second", recordingplan.SchemaOperationResult{
		Status: recordingplan.SchemaOperationSucceeded, VerifiedDigest: f.snapshot.Group.Destination.SchemaDigest,
	})
	require.NoError(t, err)

	// Update the storage scope lineage to record the second confirmation's operation ID.
	table, owner, ok := workspace.LongStorageIdentity(f.snapshot.Group)
	require.True(t, ok)
	e := f.snapshot.RuntimeLayout.SQLServer.Group.Evidence
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(table, owner, recordingplan.SQLServerGroupSchemaOptions{
		Database: e.Database, Schema: e.Schema, Collation: "Latin1_General_100_BIN2", ServerVersion: e.ServerVersion,
		ServerEdition: e.ServerEdition, EngineEdition: e.EngineEdition, CompatibilityLevel: e.CompatibilityLevel, CapabilityRevision: e.CapabilityRevision,
	})
	require.NoError(t, err)
	require.NoError(t, recordingplan.RecordStorageConfirmation(t.Context(), f.store.DB(), scope, recordingplan.StorageConfirmation{
		Layout: layout, ConnectorRevision: f.snapshot.Group.Destination.ConnectorRevision, OperationID: secondOpID,
		VerifiedDigest: f.snapshot.Group.Destination.SchemaDigest, GroupID: "second-group-id",
	}, time.Now()))

	// Now check proof for the first group (whose Destination.SchemaRevision is still the first op ID).
	// Under shared physical storage, the first group must remain valid and not be deemed stale.
	restarted := NewStorageRecovery(f.store, f.recovery.ledger, nil)
	restarted.sqlServerOps = f.recovery.sqlServerOps
	proof, err := restarted.sqlServerCurrentProof(t.Context(), f.snapshot, f.target, true)
	require.NoError(t, err, "subsequent confirmation of another group on the same physical table must not invalidate earlier verified group")
	require.NotNil(t, proof)
}
