package dbtarget

import (
	"context"
	"strings"
	"testing"

	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

func TestMySQLManagedCatalogRejectsUnsafeDefinitions(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			cfg := loadMySQLLongFixtureConfig(t, family)
			db := openMySQLLongFixtureDB(t, cfg)
			inspector, capability := mysqlManagedLiveInspector(t, cfg)
			var owned []string
			t.Cleanup(func() { dropFixtureTablesOwned(t, db, owned) })
			create := func(t *testing.T, mutation func(string) string) string {
				t.Helper()
				table := uniqueLongTestTable(t)
				layout, err := recordingplan.ManagedMySQLLongGroupLayout(table, uniqueLongTestOwnerColumn(t), recordingplan.MySQLGroupSchemaOptions{
					ServerFamily: family, ServerVersion: capability.ServerVersion, CapabilityRevision: capability.CapabilityRevision,
					Engine: "InnoDB", Collation: mysqlFamilyCollationFor(family),
				})
				require.NoError(t, err)
				statements, err := recordingplan.GenerateGroupSchemaDDL("mysql", layout)
				require.NoError(t, err)
				statement := mutation(statements[0])
				require.NotEqual(t, statements[0], statement, "the unsafe definition must actually differ")
				_, err = db.ExecContext(t.Context(), statement)
				require.NoError(t, err)
				owned = append(owned, table)
				return table
			}
			inspect := func(ctx context.Context, table string) (*MySQLManagedInspection, error) {
				return inspector.InspectMySQLManagedTableAtRevision(ctx, capability.ConnectorID, capability.ConnectorRevision, cfg.Database, capability.CapabilityRevision, table)
			}
			t.Run("prefix_primary_key_is_not_an_exact_identity", func(t *testing.T) {
				table := create(t, func(statement string) string {
					return strings.Replace(statement, "PRIMARY KEY (`bucket_start`, `record_id`)", "PRIMARY KEY (`bucket_start`, `record_id`(10))", 1)
				})
				_, err := inspect(t.Context(), table)
				require.Error(t, err, "a prefix primary key cannot prove exact record identity")
			})
			if family != mariaDBFixtureFamily {
				return
			}
			t.Run("vacuous_json_check_does_not_validate_provenance", func(t *testing.T) {
				table := create(t, func(statement string) string {
					return strings.Replace(statement, "`provenance` JSON NOT NULL", "`provenance` LONGTEXT NOT NULL CHECK (json_valid(`provenance`) OR TRUE)", 1)
				})
				result, err := inspect(t.Context(), table)
				require.NoError(t, err)
				require.False(t, result.Evidence.JSONColumns["provenance"], "OR TRUE accepts invalid JSON and cannot be an alias proof")
			})
			t.Run("another_tables_json_check_is_not_borrowed", func(t *testing.T) {
				other := uniqueLongTestTable(t)
				_, err := db.ExecContext(t.Context(), "CREATE TABLE `"+other+"` (`provenance` JSON) ENGINE=InnoDB")
				require.NoError(t, err)
				owned = append(owned, other)
				table := create(t, func(statement string) string {
					return strings.Replace(statement, "`provenance` JSON NOT NULL", "`provenance` LONGTEXT NOT NULL CHECK (length(`provenance`) > 0)", 1)
				})
				var count int
				require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM information_schema.TABLE_CONSTRAINTS
					WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME IN (?,?) AND CONSTRAINT_NAME='provenance' AND CONSTRAINT_TYPE='CHECK'`, table, other).Scan(&count))
				require.Equal(t, 2, count, "the fixture must reproduce colliding constraint names")
				result, err := inspect(t.Context(), table)
				require.NoError(t, err)
				require.False(t, result.Evidence.JSONColumns["provenance"], "a different table's json_valid constraint is not evidence for this one")
			})
		})
	}
}
