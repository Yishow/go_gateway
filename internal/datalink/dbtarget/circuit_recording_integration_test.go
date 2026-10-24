package dbtarget

import (
	"database/sql"
	"testing"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// TestCircuitRecordingIntegration covers decimal engine/encoder/SQL precision.
// Decimal is not a PLC Tag type; this does not claim collector support for it.
func TestCircuitRecordingIntegration(t *testing.T) {
	type result struct{ raw, delta, status sql.NullString }
	var sqliteRows []result
	for _, mode := range []string{"sqlite", "postgres", "timescale"} {
		t.Run(mode, func(t *testing.T) {
			dialect, kind, namespace := SQLDialectSQLite, schema.DatabaseConnectorKindSQLite, "main"
			var db *sql.DB
			var fixture managedDestinationFixture
			if mode == "sqlite" {
				path := newInspectionSQLiteTarget(t, longDDL(t, "sqlite")...)
				var err error
				db, err = sql.Open("sqlite", path)
				require.NoError(t, err)
				t.Cleanup(func() { _ = db.Close() })
				fixture = newManagedDestinationFixture(t, kind, ConnectionConfig{"dsn": path})
			} else {
				dialect, kind = SQLDialectPostgres, schema.DatabaseConnectorKindPostgres
				dsn, fields := postgresKeywordDSN(t)
				db = postgresAdminDB(t, dsn)
				namespace = postgresScopedSchema(t, db)
				statements := longDDL(t, "postgres")
				if mode == "timescale" {
					capability, err := ProbeTimescale(t.Context(), db, namespace)
					require.NoError(t, err)
					if !capability.Available {
						t.Skipf("NOT_RUN: Timescale capability unavailable: %v", capability.Reasons)
					}
					layout, err := recordingplan.ManagedLongTimescaleLayout(indexTestTable, indexTestOwner)
					require.NoError(t, err)
					layout.Partition.ExtensionSchema = capability.ExtensionSchema
					statements, err = recordingplan.GenerateGroupSchemaDDL("postgres", layout)
					require.NoError(t, err)
				}
				tx, err := db.BeginTx(t.Context(), nil)
				require.NoError(t, err)
				_, err = tx.ExecContext(t.Context(), `SET LOCAL search_path TO `+pgIdentifier(t, namespace))
				require.NoError(t, err)
				for _, statement := range statements {
					_, err = tx.ExecContext(t.Context(), statement)
					require.NoError(t, err)
				}
				require.NoError(t, tx.Commit())
				fixture = newManagedDestinationFixture(t, kind, postgresConnectorConfig(fields, namespace, fields["user"], fields["password"]))
			}
			inspection, err := fixture.inspector.InspectTable(t.Context(), fixture.connector.ID, namespace, indexTestTable)
			require.NoError(t, err)
			require.Equal(t, TableInspectionExists, inspection.Status)
			member := counterMemberFor(measurement.ExactDecimal)
			layout, issues := NewLongRowLayout(counterSpecFor(dialect, inspection.Columns, member))
			require.Empty(t, issues)
			outcomes := counterOutcomes(t, member, exactDecimal(t, "12000.0"), exactDecimal(t, "12000.4"), exactDecimal(t, "12000.7"))
			got := make([]result, 0, len(outcomes))
			table, placeholder := indexTestTable, "?"
			if mode != "sqlite" {
				table, placeholder = pgIdentifier(t, namespace)+"."+pgIdentifier(t, indexTestTable), "$1"
			}
			for _, outcome := range outcomes {
				request := GroupInsertRequest{
					Kind: kind, SchemaName: namespace, TableName: indexTestTable, Row: outcomesEncoded(t, layout, outcome),
					PayloadDigest: "decimal-" + outcome.RecordID, Strategy: GroupEffectReceipt, CommittedAt: "2026-10-24T12:00:00Z",
				}
				first, err := InsertGroupRow(t.Context(), db, request)
				require.NoError(t, err)
				require.False(t, first.AlreadyCommitted)
				again, err := InsertGroupRow(t.Context(), db, request)
				require.NoError(t, err)
				require.True(t, again.AlreadyCommitted)
				var row result
				require.NoError(t, db.QueryRowContext(t.Context(), `SELECT CAST(value_decimal AS TEXT), CAST(delta_value AS TEXT), status FROM `+table+` WHERE record_id=`+placeholder, outcome.RecordID).Scan(&row.raw, &row.delta, &row.status))
				got = append(got, row)
				t.Logf("SQL decimal record=%s raw=%s delta=%q known=%t status=%s", outcome.RecordID, row.raw.String, row.delta.String, row.delta.Valid, row.status.String)
			}
			require.Equal(t, "baseline", got[0].status.String)
			require.False(t, got[0].delta.Valid)
			for i, digits := range []string{"12000.0", "12000.4", "12000.7"} {
				require.Equal(t, digits, got[i].raw.String)
			}
			require.Equal(t, "0.4", got[1].delta.String)
			require.Equal(t, "0.3", got[2].delta.String)
			if mode == "sqlite" {
				sqliteRows = got
			} else {
				require.Equal(t, sqliteRows, got, "identical decimal inputs have equivalent SQL records")
			}
		})
	}
}
