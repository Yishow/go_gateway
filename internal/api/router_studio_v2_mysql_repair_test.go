package api

import (
	"sync"
	"testing"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func TestMySQLLongMissingRepair(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			r := newMySQLRepairRig(t, family, true)
			sender := r.sender()
			original := r.enqueue(t, 0)
			result, err := sender.Deliver(t.Context(), original)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateCommitted, result.State)
			before := r.receipts(t)
			require.Len(t, before, 1)
			r.drop(t)
			effects := []string{r.enqueue(t, 1), r.enqueue(t, 2)}
			results := make([]groupdelivery.DeliveryResult, len(effects))
			errs := make([]error, len(effects))
			var wg sync.WaitGroup
			for i := range effects {
				wg.Go(func() { results[i], errs[i] = sender.Deliver(t.Context(), effects[i]) })
			}
			wg.Wait()
			for i := range effects {
				require.NoError(t, errs[i])
				require.Equal(t, groupdelivery.StateCommitted, results[i].State, "%+v", results[i])
			}
			require.True(t, r.exists(t))
			require.Equal(t, 2, r.count(t, r.group.Destination.TableName), "only the pending rows return, not lost historical rows")
			after := r.receipts(t)
			require.Len(t, after, 3)
			for key, receipt := range before {
				require.Equal(t, receipt, after[key], "original receipt digest, commit time and owner remain unchanged")
			}
			lineage := r.lineage(t)
			require.Equal(t, 2, lineage.Generation)
			require.True(t, lineage.HistoryMayBeLost)
			require.Equal(t, 1, lineage.Reconstructions)
			require.Empty(t, lineage.RepairState)
			var executed int
			require.NoError(t, r.f.db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM managed_schema_operations WHERE action=? AND executed_statements>0", recordingplan.StorageRepairAction).Scan(&executed))
			require.Equal(t, 1, executed, "concurrent deliveries use one effective repair operation")
			var committed int
			require.NoError(t, r.f.db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM wg_delivery_outbox WHERE state=?", groupdelivery.StateCommitted).Scan(&committed))
			require.Equal(t, 3, committed)
			result, err = sender.Deliver(t.Context(), original)
			require.NoError(t, err)
			require.Equal(t, groupdelivery.StateCommitted, result.State)
			require.Equal(t, 2, r.count(t, r.group.Destination.TableName), "an old committed effect is never rewritten into the rebuilt table")
			require.Equal(t, after, r.receipts(t))
			rows, err := r.f.target.QueryContext(t.Context(), "SELECT value_float FROM `"+r.group.Destination.TableName+"` ORDER BY value_float")
			require.NoError(t, err)
			defer rows.Close()
			var values []float64
			for rows.Next() {
				var value float64
				require.NoError(t, rows.Scan(&value))
				values = append(values, value)
			}
			require.NoError(t, rows.Err())
			require.Equal(t, []float64{11.5, 12.5}, values)
		})
	}
}

func TestMySQLLongRepairDenied(t *testing.T) {
	for _, family := range []string{"mysql", "mariadb"} {
		t.Run(family, func(t *testing.T) {
			for _, reason := range []string{"disabled", "foreign", "connector revision"} {
				t.Run(reason, func(t *testing.T) {
					r := newMySQLRepairRig(t, family, reason != "disabled")
					sender := r.sender()
					result, err := sender.Deliver(t.Context(), r.enqueue(t, 0))
					require.NoError(t, err)
					require.Equal(t, groupdelivery.StateCommitted, result.State)
					receipts := r.receipts(t)
					r.drop(t)
					if reason == "foreign" {
						_, err = r.f.target.ExecContext(t.Context(), "CREATE TABLE `"+r.group.Destination.TableName+"` (note VARCHAR(32)) ENGINE=InnoDB")
						require.NoError(t, err)
						_, err = r.f.target.ExecContext(t.Context(), "INSERT INTO `"+r.group.Destination.TableName+"` (note) VALUES (?)", "foreign fixture row")
						require.NoError(t, err)
					}
					if reason == "connector revision" {
						_, err = r.f.db.ExecContext(t.Context(), "UPDATE database_connectors SET identity_revision='changed' WHERE id=?", r.group.Destination.ConnectorID)
						require.NoError(t, err)
					}
					result, err = sender.Deliver(t.Context(), r.enqueue(t, 1))
					require.NoError(t, err)
					require.Equal(t, groupdelivery.StateBlocked, result.State)
					code := recordingplan.StorageCodeRepairBlocked
					if reason == "foreign" {
						code = recordingplan.StorageCodeOwnerMismatch
						var note string
						require.NoError(t, r.f.target.QueryRowContext(t.Context(), "SELECT note FROM `"+r.group.Destination.TableName+"`").Scan(&note))
						require.Equal(t, "foreign fixture row", note)
					} else {
						require.False(t, r.exists(t), "a denied repair cannot create a table")
					}
					require.Equal(t, code, result.Code)
					require.Equal(t, receipts, r.receipts(t))
					lineage := r.lineage(t)
					require.Equal(t, 1, lineage.Generation)
					require.False(t, lineage.HistoryMayBeLost)
					require.Zero(t, lineage.Reconstructions)
				})
			}
		})
	}
}
