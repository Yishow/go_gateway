package grouppipeline

import (
	"encoding/json"
	"testing"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

func TestCounterFinalizedResolversKeepStorageProofIndependentOfReplay(t *testing.T) {
	for _, dialect := range []string{"mysql", "sqlserver"} {
		for _, corruptStorage := range []bool{false, true} {
			t.Run(dialect+map[bool]string{false: " valid storage", true: " forged storage"}[corruptStorage], func(t *testing.T) {
				fixture := newFixture(t)
				pipe := fixture.pipeline()
				snap := mysqlRuntimeSnapshot(t)
				if dialect == "sqlserver" {
					snap = sqlServerRuntimeSnapshot(t)
				}
				snap.Group.RecordingProfile.Kind = workspace.WriteGroupProfileCumulative
				snap.Group.Members[0].MeasurementID = ptr("original-counter")
				var resolver destinationResolver
				var item groupdelivery.OutboxItem
				if dialect == "mysql" {
					resolver, item, _ = newMySQLResolverFixture(t, snap, mysqlRuntimeCapability(snap.RuntimeLayout.MySQL))
				} else {
					resolver, item, _ = sqlServerResolverFixture(t, snap)
				}
				if corruptStorage {
					if dialect == "mysql" {
						snap.RuntimeLayout.MySQL.ConnectorRevision = "foreign"
					} else {
						snap.RuntimeLayout.SQLServer.ConnectorRevision = "foreign"
					}
				}
				raw, err := json.Marshal(snap)
				require.NoError(t, err)
				var document map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(raw, &document))
				document["counter_bindings"] = json.RawMessage(`{"future":{"version":77}}`)
				raw, err = json.Marshal(document)
				require.NoError(t, err)
				require.NoError(t, fixture.store.SaveRuntimeVersion(t.Context(), groupKey(snap), raw))
				_, err = pipe.loadRuntimeSnapshot(t.Context(), groupKey(snap))
				require.ErrorIs(t, err, counterproof.ErrCounterBindingProofUnsupported, "replay remains refused")
				resolver.loadRuntimeSnapshot = pipe.loadFinalizedSnapshot
				target, err := resolver.Resolve(t.Context(), item)
				if corruptStorage {
					require.ErrorIs(t, err, groupdelivery.ErrTargetBlocked)
					require.Nil(t, target.DB)
				} else {
					require.NoError(t, err)
					require.NotNil(t, target.DB)
					// Compare all persisted facts, not transient json:"-" inspection flags.
					var expectedProof, actualProof any = snap.RuntimeLayout.MySQL, target.MySQL
					if dialect == "sqlserver" {
						expectedProof, actualProof = snap.RuntimeLayout.SQLServer, target.SQLServer
					}
					expectedBytes, err := json.Marshal(expectedProof)
					require.NoError(t, err)
					actualBytes, err := json.Marshal(actualProof)
					require.NoError(t, err)
					require.Equal(t, expectedBytes, actualBytes)
					require.NoError(t, target.Close())
				}
				stored, err := fixture.store.RuntimeVersion(t.Context(), groupKey(snap))
				require.NoError(t, err)
				require.Equal(t, raw, stored)
			})
		}
	}
}
