package api

import (
	"encoding/json"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/testutil/lostcommit"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

// The server and normal metadata are real fixtures; only the catalog permission
// response is injected. This does not claim acceptance with a restricted role.
func TestMySQLLongRepairDeniedPermission(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			r := newMySQLRepairRig(t, family, true)
			result, err := r.sender().Deliver(t.Context(), r.enqueue(t, 0))
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateCommitted, result.State)
			before := r.receipts(t)
			r.drop(t)
			effect := r.enqueue(t, 1)
			item, err := r.store.GetOutbox(t.Context(), effect)
			require.NoError(t, err)
			target, err := r.Resolve(t.Context(), item)
			require.NoError(t, err)
			t.Cleanup(func() { _ = target.Close() })
			row, err := groupdelivery.DecodeRowPayload(item.Payload)
			require.NoError(t, err)
			_, missing := dbtarget.InsertGroupRow(t.Context(), target.DB, dbtarget.GroupInsertRequest{
				Kind: target.Kind, SchemaName: item.Destination.TableSchema, TableName: item.Destination.TableName,
				Strategy: item.Destination.DedupeCapability, Row: row, PayloadDigest: item.PayloadDigest, MySQL: target.MySQL,
				CommittedAt: time.Now().UTC().Format("2006-01-02T15:04:05.000000Z"),
			})
			require.Error(t, missing)
			require.True(t, dbtarget.IsMissingRelationError(missing), "the trigger is an actual server missing-table response")
			connector, err := r.targets.GetByID(t.Context(), item.Destination.ConnectorID)
			require.NoError(t, err)
			var config struct {
				DSN string `json:"dsn"`
			}
			require.NoError(t, json.Unmarshal([]byte(connector.ConnectionConfig), &config))
			require.True(t, config.DSN != "")
			var ddl atomic.Int32
			reject := func(query string) error {
				upper := strings.ToUpper(strings.TrimSpace(query))
				if strings.HasPrefix(upper, "CREATE ") {
					ddl.Add(1)
				}
				if strings.HasPrefix(upper, "SHOW CREATE TABLE") {
					return &mysql.MySQLError{Number: 1142, Message: "injected catalog permission denial"}
				}
				return nil
			}
			faults := &lostcommit.Faults{}
			faults.Reject.Store(&reject)
			denied := lostcommit.WrapFaults(&mysql.MySQLDriver{}, config.DSN, faults)
			t.Cleanup(func() { _ = denied.Close() })
			decision, err := r.recovery.Restore(t.Context(), item, groupdelivery.Target{
				DB: denied, Kind: target.Kind, MySQL: target.MySQL, MySQLCapability: target.MySQLCapability,
			}, missing)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.RestoreBlock, decision.Action)
			require.Equal(t, recordingplan.StorageCodeRepairBlocked, decision.Code)
			require.Zero(t, ddl.Load(), "forbidden metadata must not be treated as permission to create a missing table")
			require.False(t, r.exists(t))
			require.Equal(t, before, r.receipts(t))
			lineage := r.lineage(t)
			require.Equal(t, 1, lineage.Generation)
			require.False(t, lineage.HistoryMayBeLost)
			events, err := recordingplan.ReadStorageEvents(t.Context(), r.f.db, lineage.StorageScope, r.group.ID)
			require.NoError(t, err)
			var visible bool
			for _, event := range events {
				visible = visible || event.Reason == recordingplan.StorageReasonPermissionDenied
			}
			require.True(t, visible, "permission refusal is a durable safe diagnostic")
		})
	}
}
