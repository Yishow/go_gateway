package groupdelivery

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOperatorResolutionRealFaultsKeepFrozenIdentityAndDedupe(t *testing.T) {
	for _, fault := range []string{"missing-table", "permission", "poison"} {
		for _, resolution := range []QuarantineResolution{ResolutionRetry, ResolutionSkip} {
			t.Run(fault+"/"+string(resolution), func(t *testing.T) {
				f := newDeliveryFixture(t, plainTargetTable, true)
				first := f.enqueue(t, DedupeReceipt, "", at(0), 13)
				second := f.enqueue(t, DedupeReceipt, "", at(10), 2)
				before := f.state(t, first)
				var statement, repair string
				switch fault {
				case "missing-table":
					statement = `DROP TABLE readings`
					repair = plainTargetTable
				case "permission":
					f.target.SetMaxOpenConns(1)
					statement = `PRAGMA query_only=ON`
					repair = `PRAGMA query_only=OFF`
				case "poison":
					statement = `CREATE TRIGGER poison BEFORE INSERT ON readings WHEN NEW.counter=13 BEGIN SELECT RAISE(ABORT,'rejected row'); END`
					repair = `DROP TRIGGER poison`
				}
				_, err := f.target.ExecContext(t.Context(), statement)
				require.NoError(t, err)
				result, err := f.sender.Deliver(t.Context(), first)
				require.NoError(t, err)
				expected := StateBlocked
				if fault == "poison" {
					expected = StateQuarantined
				}
				require.Equal(t, expected, result.State)
				heads, err := f.store.ReadyHeads(t.Context(), time.Now().Add(time.Minute))
				require.NoError(t, err)
				require.Empty(t, heads, "head prevents later row")
				_, err = f.target.ExecContext(t.Context(), repair)
				require.NoError(t, err)
				stateRevision := f.state(t, first).ClaimEpoch
				decision := OperatorDecision{ExpectedStateRevision: &stateRevision, DecisionID: fault + "-" + string(resolution), EffectKey: first, ExpectedState: expected, PayloadDigest: before.PayloadDigest, Resolution: resolution, Reason: "Owned test destination repaired", ConfirmSkip: resolution == ResolutionSkip}
				failOn(t, f.local, "workspace_audit_history")
				_, err = f.store.ResolveAttention(t.Context(), testKey.WorkspaceID, testKey.GroupID, decision)
				require.Error(t, err)
				require.Equal(t, expected, f.state(t, first).State, "failed audit rolls back state")
				_, err = f.local.ExecContext(t.Context(), `DROP TRIGGER fail_workspace_audit_history`)
				require.NoError(t, err)
				_, err = f.store.ResolveAttention(t.Context(), testKey.WorkspaceID, testKey.GroupID, decision)
				require.NoError(t, err)
				after := f.state(t, first)
				require.Equal(t, before.Destination, after.Destination)
				require.Equal(t, before.Payload, after.Payload)
				require.Equal(t, before.PayloadDigest, after.PayloadDigest)
				require.Equal(t, before.RecordID, after.RecordID)
				if resolution == ResolutionRetry {
					result, err = f.sender.Deliver(t.Context(), first)
					require.NoError(t, err)
					require.Equal(t, StateCommitted, result.State)
				}
				heads, err = f.store.ReadyHeads(t.Context(), time.Now().Add(time.Minute))
				require.NoError(t, err)
				require.Len(t, heads, 1)
				require.Equal(t, second, heads[0].EffectKey)
				result, err = f.sender.Deliver(t.Context(), second)
				require.NoError(t, err)
				require.Equal(t, StateCommitted, result.State)
				rows := 1
				if resolution == ResolutionRetry {
					rows = 2
				}
				require.Equal(t, rows, f.targetRows(t))
				require.Equal(t, rows, countRows(t, f.target, "gw_effect_receipts"))
				_, err = f.store.ResolveAttention(t.Context(), testKey.WorkspaceID, testKey.GroupID, decision)
				require.NoError(t, err)
				require.Equal(t, rows, f.targetRows(t), "duplicate operator decision never resends")
			})
		}
	}
}
