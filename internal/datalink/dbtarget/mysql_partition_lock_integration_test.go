package dbtarget

import (
	"context"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

// TestMySQLPartitionReorganizeLockTimeoutIntegration proves the database's
// metadata-lock timeout bounds maintenance while retaining the original table.
func TestMySQLPartitionReorganizeLockTimeoutIntegration(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			cfg := loadMySQLLongFixtureConfig(t, family)
			db := openMySQLLongFixtureDB(t, cfg)
			db.SetMaxOpenConns(2)
			inspector, capability := mysqlManagedLiveInspector(t, cfg)
			table := "gw_long_" + randomHex(t, 8)
			require.False(t, fixtureTableExists(t, db, table))
			ddl, _ := mysqlNativePartitionFixtureDDL(table, uniqueLongTestOwnerColumn(t))
			ddl = strings.ReplaceAll(ddl, "`p_max`", "`p_future`")
			_, err := db.ExecContext(t.Context(), ddl)
			require.NoError(t, err)
			t.Cleanup(func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_, err := db.ExecContext(ctx, "DROP TABLE "+quoteIdentifier(schema.DatabaseConnectorKindMySQL, table))
				require.NoError(t, err)
			})
			proof, err := inspector.InspectMySQLManagedTableAtRevision(t.Context(), capability.ConnectorID,
				capability.ConnectorRevision, cfg.Database, capability.CapabilityRevision, table)
			require.NoError(t, err)
			before := proof.Evidence.Partitions
			after := cloneMySQLPartitionEvidence(before)
			last := len(after.Partitions) - 1
			start, err := time.Parse(time.RFC3339, after.Partitions[last-1].LessThan)
			require.NoError(t, err)
			after.Partitions[last] = MySQLPartitionBoundary{
				Name: "p_" + start.Format("200601"), LessThan: start.AddDate(0, 1, 0).Format(time.RFC3339), Ordinal: last + 1,
			}
			after.Partitions = append(after.Partitions, MySQLPartitionBoundary{Name: "p_future", LessThan: "MAXVALUE", Ordinal: last + 2})
			early, err := time.Parse(time.RFC3339, before.Partitions[0].LessThan)
			require.NoError(t, err)
			_, err = db.ExecContext(t.Context(), "INSERT INTO "+quoteIdentifier(schema.DatabaseConnectorKindMySQL, table)+
				" (bucket_start, record_id) VALUES (?, ?), (?, ?)", early.Add(-24*time.Hour), strings.Repeat("a", 64),
				start.Add(24*time.Hour), strings.Repeat("b", 64))
			require.NoError(t, err)

			var original int
			require.NoError(t, db.QueryRowContext(t.Context(), "SELECT @@SESSION.lock_wait_timeout").Scan(&original))
			blocker, err := db.BeginTx(t.Context(), nil)
			require.NoError(t, err)
			defer blocker.Rollback()
			var count int
			require.NoError(t, blocker.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+
				quoteIdentifier(schema.DatabaseConnectorKindMySQL, table)).Scan(&count))
			started := time.Now()
			err = ReorganizeOpenMySQLFuturePartition(t.Context(), db, capability, cfg.Database, proof, after)
			elapsed := time.Since(started)
			var sqlErr *mysqldriver.MySQLError
			require.ErrorAs(t, err, &sqlErr)
			require.Equal(t, uint16(1205), sqlErr.Number, "the session metadata-lock timeout must refuse the DDL")
			require.Less(t, elapsed, 6*time.Second, "one blocked table cannot stall maintenance indefinitely")
			require.NoError(t, blocker.Rollback())
			var restored int
			require.NoError(t, db.QueryRowContext(t.Context(), "SELECT @@SESSION.lock_wait_timeout").Scan(&restored))
			require.Equal(t, original, restored, "maintenance must restore the pooled connection setting")
			unchanged, err := InspectOpenMySQLManagedTable(t.Context(), db, capability, cfg.Database, table)
			require.NoError(t, err)
			require.Equal(t, before, unchanged.Evidence.Partitions)
			t.Logf("%s metadata-lock refusal: %s; original catalog retained", family, elapsed)
			// Releasing the competing transaction permits the same bounded plan.
			// Existing past and catch-all data must survive the actual reorganization.
			require.NoError(t, ReorganizeOpenMySQLFuturePartition(t.Context(), db, capability, cfg.Database, proof, after))
			extended, err := InspectOpenMySQLManagedTable(t.Context(), db, capability, cfg.Database, table)
			require.NoError(t, err)
			require.Equal(t, after, extended.Evidence.Partitions)
			require.NoError(t, db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+
				quoteIdentifier(schema.DatabaseConnectorKindMySQL, table)).Scan(&count))
			require.Equal(t, 2, count)
			require.NoError(t, db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM "+
				quoteIdentifier(schema.DatabaseConnectorKindMySQL, table)+" PARTITION ("+
				quoteIdentifier(schema.DatabaseConnectorKindMySQL, after.Partitions[last].Name)+") WHERE record_id = ?",
				strings.Repeat("b", 64)).Scan(&count))
			require.Equal(t, 1, count, "the existing catch-all row moves into the new month without deletion")
		})
	}
}
