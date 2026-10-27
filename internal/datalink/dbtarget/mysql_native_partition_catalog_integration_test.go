package dbtarget

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// TestMySQLNativePartitionCatalogIntegration exercises the production managed
// catalog reader against an owned native RANGE COLUMNS table. It is skipped
// unless the private fixture env is present; root runs it serially.
func TestMySQLNativePartitionCatalogIntegration(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			cfg := loadMySQLLongFixtureConfig(t, family)
			db := openMySQLLongFixtureDB(t, cfg)
			requireMySQLFixtureFamily(t, db, family)
			inspector, capability := mysqlManagedLiveInspector(t, cfg)
			table := uniqueLongTestTable(t)
			owner := uniqueLongTestOwnerColumn(t)
			ddl, expected := mysqlNativePartitionFixtureDDL(table, owner)
			_, err := db.ExecContext(t.Context(), ddl)
			require.NoError(t, err)
			t.Cleanup(func() { dropFixtureTablesOwned(t, db, []string{table}) })

			result, err := inspector.InspectMySQLManagedTableAtRevision(t.Context(),
				capability.ConnectorID, capability.ConnectorRevision, cfg.Database, capability.CapabilityRevision, table)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotNil(t, result.Evidence)
			require.NotNil(t, result.Evidence.Partitions, "production catalog reader must expose actual partitions")
			require.Equal(t, "RANGE COLUMNS", result.Evidence.Partitions.Method)
			require.Equal(t, "bucket_start", result.Evidence.Partitions.Expression)
			require.Equal(t, expected, result.Evidence.Partitions.Partitions)
		})
	}
}

func mysqlNativePartitionFixtureDDL(table, owner string) (string, []MySQLPartitionBoundary) {
	now := time.Now().UTC()
	base := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	boundaries := []time.Time{
		base, base.AddDate(0, 1, 0), base.AddDate(0, 2, 0), base.AddDate(0, 3, 0), base.AddDate(0, 4, 0),
	}
	names := []string{"p_past", "p_current", "p_future_1", "p_future_2", "p_future_3"}
	partitions := make([]string, 0, len(boundaries)+1)
	expected := make([]MySQLPartitionBoundary, 0, len(boundaries)+1)
	for i, boundary := range boundaries {
		partitions = append(partitions, "PARTITION "+quoteIdentifier(schema.DatabaseConnectorKindMySQL, names[i])+
			" VALUES LESS THAN ('"+boundary.Format("2006-01-02 15:04:05")+"')")
		expected = append(expected, MySQLPartitionBoundary{
			Name: names[i], LessThan: boundary.Format(time.RFC3339Nano), Ordinal: i + 1,
		})
	}
	partitions = append(partitions, "PARTITION `p_max` VALUES LESS THAN (MAXVALUE)")
	expected = append(expected, MySQLPartitionBoundary{Name: "p_max", LessThan: "MAXVALUE", Ordinal: len(boundaries) + 1})
	return fmt.Sprintf(`CREATE TABLE %s (
		%s DATETIME(6) NOT NULL,
		record_id CHAR(64) NOT NULL,
		%s VARCHAR(64) NOT NULL DEFAULT 'gateway',
		PRIMARY KEY (%s, record_id)
	) ENGINE=InnoDB PARTITION BY RANGE COLUMNS(%s) (%s)`,
		"`"+table+"`", "`bucket_start`", "`"+owner+"`", "`bucket_start`", "`bucket_start`", strings.Join(partitions, ", ")), expected
}
