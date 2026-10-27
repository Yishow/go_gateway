package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type nativeSQLServerCatalogScript struct {
	table        string
	base         *sqlServerManagedCatalogScript
	objects      []driver.Value
	bounds       [][]driver.Value
	destinations [][]driver.Value
	indexes      [][]driver.Value
	nativeReads  int
}

type nativeSQLServerCatalogConnector struct{ script *nativeSQLServerCatalogScript }

func (c nativeSQLServerCatalogConnector) Connect(context.Context) (driver.Conn, error) {
	return &nativeSQLServerCatalogConn{sqlServerManagedCatalogConn: sqlServerManagedCatalogConn{script: c.script.base}, native: c.script}, nil
}

func (nativeSQLServerCatalogConnector) Driver() driver.Driver { return sqlServerManagedCatalogDriver{} }

type nativeSQLServerCatalogConn struct {
	sqlServerManagedCatalogConn
	native *nativeSQLServerCatalogScript
}

func (c *nativeSQLServerCatalogConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	rows := func(count int, values ...[]driver.Value) driver.Rows {
		columns := make([]string, count)
		for i := range columns {
			columns[i] = "catalog_value"
		}
		return sqlServerRecordingRows(columns, values...)
	}
	switch {
	case strings.Contains(query, "gw_native_partition_objects"):
		c.native.nativeReads++
		if len(args) != 3 || args[0].Value != "pf_"+c.native.table || args[1].Value != "ps_"+c.native.table || args[2].Value != "gateway.recording.partition_owner" {
			return nil, errors.New("native object identities must be parameter bound")
		}
		return rows(12, c.native.objects), nil
	case strings.Contains(query, "gw_native_partition_bounds"):
		return rows(3, c.native.bounds...), nil
	case strings.Contains(query, "gw_native_partition_destinations"):
		return rows(3, c.native.destinations...), nil
	case strings.Contains(query, "gw_native_partition_indexes"):
		return rows(6, c.native.indexes...), nil
	default:
		return c.sqlServerManagedCatalogConn.QueryContext(ctx, query, args)
	}
}

func nativeSQLServerCatalogFixture() *nativeSQLServerCatalogScript {
	base := sqlServerManagedCatalogFixture()
	for _, row := range base.indexRows {
		row[9] = "PS"
	}
	script := &nativeSQLServerCatalogScript{
		table: "circuit_long",
		base:  base,
		objects: []driver.Value{
			int64(101), true, "R", "datetime2", int64(6), int64(1),
			int64(201), int64(101), "dbo.circuit_long", "dbo.circuit_long", int64(1), true,
		},
		indexes: [][]driver.Value{
			{int64(91), int64(1), int64(201), int64(1), int64(1), true},
			{int64(91), int64(2), int64(201), int64(1), int64(0), false},
		},
	}
	for i, value := range []string{"2026-10-01", "2026-11-01", "2026-12-01", "2027-01-01", "2027-02-01"} {
		at, _ := time.Parse("2006-01-02", value)
		script.bounds = append(script.bounds, []driver.Value{int64(i + 1), int64(1), at})
	}
	for i := range 6 {
		script.destinations = append(script.destinations, []driver.Value{int64(i + 1), "PRIMARY", "FG"})
	}
	return script
}

func inspectNativeSQLServerCatalog(t *testing.T, script *nativeSQLServerCatalogScript) (*SQLServerManagedInspection, error) {
	t.Helper()
	inspector, connector := sqlServerManagedInspectionTestSetup(t, script.base, "native-catalog")
	db := sql.OpenDB(nativeSQLServerCatalogConnector{script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	openSQLServerManagedSessionFunc = func(schema.DatabaseConnectorKind, ConnectionConfig) (sqlServerManagedSession, error) {
		return &sqlServerManagedCatalogTestSession{db: db, script: script.base}, nil
	}
	return inspector.InspectSQLServerManagedTableAtRevision(t.Context(), connector.ID, connector.IdentityRevision,
		"fixture_db", "dbo", sqlServerManagedExpectedCapabilityRevision(connector), "circuit_long")
}

func nativeSQLServerEvidenceJSON(t *testing.T, physical *SQLServerManagedInspection) map[string]any {
	t.Helper()
	require.NotNil(t, physical)
	raw, err := json.Marshal(physical.Evidence)
	require.NoError(t, err)
	var value map[string]any
	require.NoError(t, json.Unmarshal(raw, &value))
	native, ok := value["NativeMonthly"].(map[string]any)
	require.True(t, ok, "live partition catalog must survive the frozen proof boundary")
	return native
}

func TestSQLServerNativeCatalogReadsOwnedAlignedObjects(t *testing.T) {
	script := nativeSQLServerCatalogFixture()
	physical, err := inspectNativeSQLServerCatalog(t, script)
	require.NoError(t, err)
	native := nativeSQLServerEvidenceJSON(t, physical)
	require.Equal(t, true, native["Exclusive"])
	require.Equal(t, true, native["IndexesAligned"])
	require.Equal(t, "RANGE RIGHT", native["Method"])
	require.Equal(t, "pf_circuit_long", native["FunctionName"])
	require.Len(t, native["Partitions"], 6)
	require.Equal(t, 1, script.nativeReads)
	require.False(t, physical.Evidence.Plain)
	require.Zero(t, script.base.execCalls)
}

func TestSQLServerNativeCatalogRejectsForeignOrUnalignedObjects(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*nativeSQLServerCatalogScript)
		field  string
	}{
		{"foreign marker", func(s *nativeSQLServerCatalogScript) { s.objects[8] = "dbo.foreign" }, "Exclusive"},
		{"missing scheme marker", func(s *nativeSQLServerCatalogScript) { s.objects[9] = nil }, "Exclusive"},
		{"shared function", func(s *nativeSQLServerCatalogScript) { s.objects[10] = int64(2) }, "Exclusive"},
		{"scheme references other function", func(s *nativeSQLServerCatalogScript) { s.objects[7] = int64(999) }, "Exclusive"},
		{"external table", func(s *nativeSQLServerCatalogScript) {
			s.indexes = append(s.indexes, []driver.Value{int64(999), int64(1), int64(201), int64(1), int64(1), true})
		}, "Exclusive"},
		{"unaligned source index", func(s *nativeSQLServerCatalogScript) { s.indexes[1][2] = int64(301) }, "IndexesAligned"},
		{"unique lacks time", func(s *nativeSQLServerCatalogScript) { s.indexes[0][4] = int64(0) }, "IndexesAligned"},
	} {
		t.Run(test.name, func(t *testing.T) {
			script := nativeSQLServerCatalogFixture()
			test.mutate(script)
			physical, err := inspectNativeSQLServerCatalog(t, script)
			require.NoError(t, err)
			native := nativeSQLServerEvidenceJSON(t, physical)
			require.Equal(t, false, native[test.field])
		})
	}
}

func TestSQLServerNativeCatalogConfirmsExclusiveFunctionBeforeSchemeExists(t *testing.T) {
	script := nativeSQLServerCatalogFixture()
	script.objects[6], script.objects[7], script.objects[9], script.objects[10] = nil, nil, nil, int64(0)
	script.destinations = nil
	script.indexes = nil
	physical, err := inspectNativeSQLServerCatalog(t, script)
	require.NoError(t, err)
	native := physical.Evidence.NativeMonthly
	require.NotNil(t, native)
	require.True(t, native.FunctionExists)
	require.False(t, native.SchemeExists)
	require.True(t, native.Exclusive)
	require.Equal(t, "RANGE RIGHT", native.Method)
	require.Empty(t, native.Filegroup, "the missing scheme has no live destination mapping")
	require.Len(t, native.Partitions, 6)
}

func TestSQLServerNativeCatalogDoesNotInferAbsenceFromInvisibleObjects(t *testing.T) {
	script := nativeSQLServerCatalogFixture()
	script.objects[11] = false
	_, err := inspectNativeSQLServerCatalog(t, script)
	require.ErrorIs(t, err, ErrSQLServerManagedTableUnconfirmed)
}

func TestSQLServerNativeCatalogAcceptsFixedCharRangeCode(t *testing.T) {
	script := nativeSQLServerCatalogFixture()
	// sys.partition_functions.type is char(2); the driver preserves padding.
	script.objects[2] = "R "
	physical, err := inspectNativeSQLServerCatalog(t, script)
	require.NoError(t, err)
	native := nativeSQLServerEvidenceJSON(t, physical)
	require.Equal(t, "RANGE RIGHT", native["Method"])
	require.Equal(t, true, native["Exclusive"])
}

func TestSQLServerNativeCatalogRequiresClusteredPrimaryKey(t *testing.T) {
	script := nativeSQLServerCatalogFixture()
	// SQL Server reserves index_id 1 for the clustered index. A nonunique
	// clustered source index plus a nonclustered unique PK is not our layout.
	script.indexes = [][]driver.Value{
		{int64(91), int64(1), int64(201), int64(1), int64(0), false},
		{int64(91), int64(2), int64(201), int64(1), int64(1), true},
	}
	db := sql.OpenDB(nativeSQLServerCatalogConnector{script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	exclusive, aligned, err := readSQLServerNativeIndexAlignment(t.Context(), db, 91, 201)
	require.NoError(t, err)
	require.True(t, exclusive)
	require.False(t, aligned, "native managed storage requires its clustered primary key")
}

func TestSQLServerNativeCatalogValidatesBoundariesAndPrimaryMapping(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*nativeSQLServerCatalogScript)
	}{
		{"duplicate boundary", func(s *nativeSQLServerCatalogScript) { s.bounds[1][2] = s.bounds[0][2] }},
		{"nonmonth boundary", func(s *nativeSQLServerCatalogScript) { s.bounds[0][2] = time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC) }},
		{"missing ordinal", func(s *nativeSQLServerCatalogScript) { s.bounds[1][0] = int64(3) }},
		{"wrong parameter", func(s *nativeSQLServerCatalogScript) { s.bounds[0][1] = int64(2) }},
		{"incomplete horizon", func(s *nativeSQLServerCatalogScript) { s.bounds = s.bounds[:4] }},
	} {
		t.Run(test.name, func(t *testing.T) {
			script := nativeSQLServerCatalogFixture()
			test.mutate(script)
			_, err := inspectNativeSQLServerCatalog(t, script)
			require.ErrorIs(t, err, ErrSQLServerManagedTableUnconfirmed)
		})
	}
	t.Run("existing filegroup must be PRIMARY", func(t *testing.T) {
		script := nativeSQLServerCatalogFixture()
		script.destinations[2][1] = "DBA_FILEGROUP"
		physical, err := inspectNativeSQLServerCatalog(t, script)
		require.NoError(t, err)
		require.Equal(t, "", nativeSQLServerEvidenceJSON(t, physical)["Filegroup"])
	})
	t.Run("optional NEXT USED still maps to PRIMARY", func(t *testing.T) {
		script := nativeSQLServerCatalogFixture()
		script.destinations = append(script.destinations, []driver.Value{int64(7), "PRIMARY", "FG"})
		physical, err := inspectNativeSQLServerCatalog(t, script)
		require.NoError(t, err)
		require.Equal(t, "PRIMARY", nativeSQLServerEvidenceJSON(t, physical)["Filegroup"])
	})
}
