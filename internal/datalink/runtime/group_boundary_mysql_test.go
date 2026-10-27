package runtime

import (
	"database/sql"
	"math"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// mysqlIntakeFixture exercises the real journal/selector with declared MySQL
// columns. It does not claim catalog readiness or destination SQL delivery.
func mysqlIntakeFixture(t *testing.T, decimalType string) (*longFixture, *GroupBoundary, *sql.DB) {
	t.Helper()
	db := openDurableDB(t, filepath.Join(t.TempDir(), "journal.db"))
	f := newLongFixture(t)
	f.config.Ledger = longLedger(db)
	b := f.build(t)
	layout, err := recordingplan.ManagedMySQLLongGroupLayout(longFixtureTable, longFixtureOwner, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: "mysql", ServerVersion: "8.4.13", CapabilityRevision: "unit-declaration",
		Engine: "InnoDB", Collation: "utf8mb4_0900_bin",
	})
	require.NoError(t, err)
	columns := make([]dbtarget.ColumnInfo, 0, len(layout.Columns))
	for _, column := range layout.Columns {
		declared := column.SQLType
		if column.Name == "value_decimal" {
			declared = decimalType
		}
		columns = append(columns, dbtarget.ColumnInfo{
			Name: column.Name, DataType: declared, Nullable: column.Nullable, PrimaryKey: column.PrimaryKey,
		})
	}
	members := make([]dbtarget.LongRowMember, 0, len(f.members))
	for _, member := range f.members {
		kind, ok := measurement.ExactTypeForTag(f.config.TagTypes[member.TagID])
		require.True(t, ok)
		members = append(members, dbtarget.LongRowMember{
			MemberKey: memberKey(member.DeviceID, member.PointID, member.TagID), DeviceID: member.DeviceID,
			PointID: member.PointID, TagID: member.TagID, CircuitID: member.CircuitID, Type: kind,
		})
	}
	encoder, issues := dbtarget.NewLongRowLayout(dbtarget.LongRowSpec{
		Dialect: dbtarget.SQLDialectMySQL, Columns: columns, WorkspaceID: "workspace-A", GroupID: "group-L",
		GroupRevision: "rev-1", ProfileKind: "instantaneous", ProfileInterval: 15, ProfileUnit: "second", Members: members,
	})
	require.Empty(t, issues)
	b.layout = encoder
	return f, b, db
}

func TestMySQLLongIntakeRejectsBeforeJournal(t *testing.T) {
	// A narrower declaration deliberately exercises the refusal using a real
	// supported uint64 source; it is not an accepted managed_long_v1 catalog.
	f, boundary, db := mysqlIntakeFixture(t, "DECIMAL(19,0)")
	envelope := f.envelope(4, "retry-same-id", at(2*time.Second), uint64(math.MaxUint64))
	err := boundary.AcceptSample(t.Context(), envelope)
	require.ErrorIs(t, err, dbtarget.ErrExactValueUnrepresentable)
	var refused *GroupSampleError
	require.ErrorAs(t, err, &refused)
	require.Equal(t, "exact_value_unrepresentable", refused.Reason)
	require.Zero(t, countJournal(t, db))
	require.Zero(t, baselineCount(t, db))

	// Retrying the same identity with a representable value proves refusal did
	// not put the first value into either the journal or the selector's seen set.
	envelope.Value = uint64(1)
	require.NoError(t, boundary.AcceptSample(t.Context(), envelope))
	require.Equal(t, 1, countJournal(t, db))
}

func TestMySQLLongIntakePreservesSharedRouting(t *testing.T) {
	t.Run("full managed decimal accepts maximum uint64", func(t *testing.T) {
		f, boundary, db := mysqlIntakeFixture(t, "DECIMAL(65,30)")
		require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(4, "maximum", at(2*time.Second), uint64(math.MaxUint64))))
		require.Equal(t, 1, countJournal(t, db))
	})
	t.Run("bad reading is journaled as unknown rather than rejected by value range", func(t *testing.T) {
		f, boundary, db := mysqlIntakeFixture(t, "DECIMAL(19,0)")
		envelope := f.envelope(4, "bad", at(2*time.Second), uint64(math.MaxUint64))
		envelope.Quality = schema.QualityBad
		require.NoError(t, boundary.AcceptSample(t.Context(), envelope))
		require.Equal(t, 1, countJournal(t, db))
	})
	t.Run("workspace refusal precedes destination value validation", func(t *testing.T) {
		f, boundary, db := mysqlIntakeFixture(t, "DECIMAL(19,0)")
		envelope := f.envelope(4, "foreign", at(2*time.Second), uint64(math.MaxUint64))
		envelope.WorkspaceID = "another-workspace"
		err := boundary.AcceptSample(t.Context(), envelope)
		var refused *GroupSampleError
		require.ErrorAs(t, err, &refused)
		require.Equal(t, reasonWorkspaceMismatch, refused.Reason)
		require.NotErrorIs(t, err, dbtarget.ErrExactValueUnrepresentable)
		require.Zero(t, countJournal(t, db))
	})
	t.Run("SQLite intake is unchanged", func(t *testing.T) {
		db := openDurableDB(t, filepath.Join(t.TempDir(), "journal.db"))
		f := newLongFixture(t)
		f.config.Ledger = longLedger(db)
		boundary := f.build(t)
		require.NoError(t, boundary.AcceptSample(t.Context(), f.envelope(4, "sqlite-maximum", at(2*time.Second), uint64(math.MaxUint64))))
		require.Equal(t, 1, countJournal(t, db))
	})
}
