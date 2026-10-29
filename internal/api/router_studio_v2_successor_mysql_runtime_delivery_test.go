package api

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

// Unlike the SQLite/PG fixture, the MySQL writer needs the production-frozen
// managed row proof. Read it from the real pipeline's immutable runtime version;
// never manufacture catalog evidence or relabel the connector driver as a family.
func (r *mysqlSuccessorRig) Resolve(ctx context.Context, item groupdelivery.OutboxItem) (groupdelivery.Target, error) {
	raw, err := r.store.RuntimeVersion(ctx, item.Key)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	var frozen workspace.WriteGroupAppliedSnapshot
	if err := json.Unmarshal(raw, &frozen); err != nil {
		return groupdelivery.Target{}, err
	}
	if err := workspace.ValidateMySQLRuntimeLayout(frozen.Group, frozen.RuntimeLayout); err != nil {
		return groupdelivery.Target{}, err
	}
	opened, err := r.targets.OpenDestination(ctx, item.Destination.ConnectorID, item.Destination.ConnectorRevision)
	if err != nil {
		return groupdelivery.Target{}, err
	}
	return groupdelivery.Target{DB: opened.DB, Kind: opened.Kind, Close: opened.Close,
		MySQL: frozen.RuntimeLayout.MySQL, MySQLCapability: opened.MySQL}, nil
}

func (r *mysqlSuccessorRig) deliver(t *testing.T) {
	t.Helper()
	now := func() time.Time { return time.Unix(0, r.clock.Load()).UTC() }
	recovery := grouppipeline.NewStorageRecovery(r.store, recordingplan.NewService(recordingplan.NewSQLRepository(r.f.db)), r.targets)
	sender := groupdelivery.NewSender(r.store, r, groupdelivery.SenderConfig{Owner: "mysql-activation-runtime", Recovery: recovery, Now: now})
	report, err := groupdelivery.NewDispatcher(r.store, sender, groupdelivery.DispatcherConfig{Now: now, MaxPartitions: 1}).RunOnce(t.Context())
	require.NoError(t, err)
	require.Zero(t, report.Blocked)
	require.Zero(t, report.Quarantined)
	require.Zero(t, report.Unknown)
	require.Zero(t, report.Retrying)
}
