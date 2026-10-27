package runtime

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"

	"github.com/stretchr/testify/require"
)

// Only the journal is a real local database. SQL Server columns below are
// declarations for testing the intake/codec, not catalog or delivery evidence.
func sqlServerScaledIntakeFixture(t *testing.T, cumulative bool, scale string) (*longFixture, *GroupBoundary, *sql.DB) {
	t.Helper()
	f := scaledCounterFixture(t, cumulative)
	if scale != "" {
		f.members[0].ExactRecordingSource.Scale = scale
	}
	db := openDurableDB(t, filepath.Join(t.TempDir(), "sqlserver-journal.db"))
	f.config.Ledger = longLedger(db)
	layout, err := recordingplan.ManagedSQLServerLongGroupLayout(longFixtureTable, longFixtureOwner, recordingplan.SQLServerGroupSchemaOptions{Schema: "dbo"})
	require.NoError(t, err)
	f.config.Columns = make([]dbtarget.ColumnInfo, 0, len(layout.Columns))
	for _, column := range layout.Columns {
		f.config.Columns = append(f.config.Columns, dbtarget.ColumnInfo{
			Name: column.Name, DataType: column.SQLType, Nullable: column.Nullable, PrimaryKey: column.PrimaryKey,
		})
	}
	f.config.Dialect = dbtarget.SQLDialectSQLServer
	f.config.Group.Destination.TableSchema = "dbo"
	boundary := f.build(t)
	return f, boundary, db
}

func TestSQLServerDecimalRejectedBeforeJournal(t *testing.T) {
	f, boundary, db := sqlServerScaledIntakeFixture(t, false, "0.0000000000000000001")
	envelope := f.envelope(0, "same-sample", at(time.Second), int64(1))
	err := boundary.AcceptSample(t.Context(), envelope)
	require.ErrorIs(t, err, dbtarget.ErrExactValueUnrepresentable)
	var refused *GroupSampleError
	require.ErrorAs(t, err, &refused)
	require.Equal(t, "exact_value_unrepresentable", refused.Reason)
	require.Zero(t, countJournal(t, db))
	require.Zero(t, baselineCount(t, db))
	envelope.Value = int64(0)
	require.NoError(t, boundary.AcceptSample(t.Context(), envelope))
	require.Equal(t, 1, countJournal(t, db), "refusal must not consume the sample identity")
}

func TestSQLServerProfilesUseFrozenExactSource(t *testing.T) {
	for _, cumulative := range []bool{false, true} {
		name, interval := "instantaneous", 15*time.Second
		if cumulative {
			name, interval = "cumulative", counterWindow
		}
		t.Run(name, func(t *testing.T) {
			f, boundary, db := sqlServerScaledIntakeFixture(t, cumulative, "")
			f.accept(t, boundary, f.envelope(0, "first", at(interval-time.Second), int64(120000)))
			var payload string
			require.NoError(t, db.QueryRowContext(t.Context(), "SELECT payload FROM wg_delivery_samples").Scan(&payload))
			var sample struct {
				Value measurement.ExactValue `json:"value"`
			}
			require.NoError(t, json.Unmarshal([]byte(payload), &sample))
			require.Equal(t, measurement.ExactDecimal, sample.Value.Type())
			require.Equal(t, measurement.Decimal("12000.0"), sample.Value.Value())

			restored := f.build(t)
			require.NoError(t, restored.Tick(t.Context(), at(interval)))
			first := outboxRows(t, db)["tag-int"]
			require.Equal(t, "12000.0", first["value_decimal"])
			require.Nil(t, first["delta_value"])
			f.clock = at(2 * interval)
			f.accept(t, restored, f.envelope(0, "second", at(2*interval-time.Second), int64(120004)))
			require.NoError(t, restored.Tick(t.Context(), at(2*interval)))
			require.NoError(t, db.QueryRowContext(t.Context(), "SELECT payload FROM wg_delivery_outbox ORDER BY bucket_start DESC LIMIT 1").Scan(&payload))
			row, err := groupdelivery.DecodeRowPayload([]byte(payload))
			require.NoError(t, err)
			cells := make(map[string]any)
			for _, cell := range row.Cells {
				cells[cell.Column] = cell.Value
			}
			require.Equal(t, "12000.4", cells["value_decimal"])
			if cumulative {
				require.Equal(t, "0.4", cells["delta_value"])
			} else {
				require.Nil(t, cells["delta_value"])
			}
		})
	}
}
