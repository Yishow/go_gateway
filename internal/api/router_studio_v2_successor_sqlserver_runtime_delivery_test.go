package api

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

// Resolve is a fixture-only adapter for explicit, clock-controlled dispatch.
// It uses the production-frozen descriptor and real SQL Server writer, but is
// NOT grouppipeline.destinationResolver. That resolver's frozen scope/capability
// guards require the separate TestSQLServerDestinationResolver* regressions.
func (r *sqlServerSuccessorRig) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	raw, err := r.store.RuntimeVersion(ctx, item.Key)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	var frozen workspace.WriteGroupAppliedSnapshot
	if err := json.Unmarshal(raw, &frozen); err != nil {
		return groupdelivery.Target{}, err
	}
	if err := workspace.ValidateSQLServerRuntimeLayout(frozen.Group, frozen.RuntimeLayout); err != nil {
		return groupdelivery.Target{}, err
	}
	destination := frozen.Group.Destination
	if frozen.WorkspaceID != item.Key.WorkspaceID || frozen.GroupID != item.Key.GroupID || frozen.GroupRevision != item.Key.GroupRevision ||
		destination.ConnectorID != item.Destination.ConnectorID || destination.ConnectorRevision != item.Destination.ConnectorRevision ||
		destination.Database != item.Destination.Database || destination.TableSchema != item.Destination.TableSchema || destination.TableName != item.Destination.TableName {
		return groupdelivery.Target{}, fmt.Errorf("%w: fixture outbox differs from frozen SQL Server destination", groupdelivery.ErrTargetBlocked)
	}
	opened, err := r.targets.OpenDestination(ctx, destination.ConnectorID, destination.ConnectorRevision)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	return groupdelivery.Target{DB: opened.DB, Kind: opened.Kind, Close: opened.Close,
		SQLServer: dbtarget.CloneSQLServerManagedRowProof(frozen.RuntimeLayout.SQLServer), SQLServerCapability: opened.SQLServer}, nil
}

func (r *sqlServerSuccessorRig) deliver(t *testing.T) {
	t.Helper()
	now := func() time.Time { return time.Unix(0, r.clock.Load()).UTC() }
	recovery := grouppipeline.NewStorageRecovery(r.store, recordingplan.NewService(recordingplan.NewSQLRepository(r.f.db)), r.targets)
	sender := groupdelivery.NewSender(r.store, r, groupdelivery.SenderConfig{Owner: "sqlserver-activation-runtime", Recovery: recovery, Now: now})
	report, err := groupdelivery.NewDispatcher(r.store, sender, groupdelivery.DispatcherConfig{Now: now, MaxPartitions: 1}).RunOnce(t.Context())
	require.NoError(t, err)
	require.Zero(t, report.Blocked)
	require.Zero(t, report.Quarantined)
	require.Zero(t, report.Unknown)
	require.Zero(t, report.Retrying)
}
