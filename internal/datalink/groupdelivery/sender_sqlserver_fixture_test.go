package groupdelivery

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strconv"
	"strings"
	"sync/atomic"

	"go-gateway/internal/datalink/dbtarget"
)

const (
	sqlServerTestOwnerColumn = "_gw_owner_0123456789abcdef01234567"
	sqlServerTestCollation   = "Latin1_General_100_BIN2"
)

type sqlServerNumberedError struct {
	number int32
	msg    string
}

func (e sqlServerNumberedError) Error() string         { return e.msg }
func (e sqlServerNumberedError) SQLErrorNumber() int32 { return e.number }

type sqlServerMockRows struct {
	columns []string
	rows    [][]driver.Value
	index   int
}

func (r *sqlServerMockRows) Columns() []string { return r.columns }
func (r *sqlServerMockRows) Close() error      { return nil }
func (r *sqlServerMockRows) Next(dest []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.index])
	r.index++
	return nil
}

type sqlServerTestDriver struct{}

func (sqlServerTestDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("must use sql.OpenDB")
}

type sqlServerTestConnector struct {
	queries  *atomic.Int32
	catalog  *sqlServerCatalogMock
	onQuery  func(query string, args []driver.NamedValue) (driver.Rows, error)
	onExec   func(query string, args []driver.NamedValue) (driver.Result, error)
	onCommit func() error
}

func (c sqlServerTestConnector) Connect(context.Context) (driver.Conn, error) {
	return &sqlServerTestConn{connector: c}, nil
}
func (c sqlServerTestConnector) Driver() driver.Driver { return sqlServerTestDriver{} }

type sqlServerTestConn struct{ connector sqlServerTestConnector }

func (c *sqlServerTestConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepared statement not used")
}
func (c *sqlServerTestConn) Close() error { return nil }
func (c *sqlServerTestConn) Begin() (driver.Tx, error) {
	return &sqlServerTestTx{connector: c.connector}, nil
}
func (c *sqlServerTestConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return &sqlServerTestTx{connector: c.connector}, nil
}

func (c *sqlServerTestConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return handleSQLServerTestQuery(c.connector, query, args)
}

func (c *sqlServerTestConn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if c.connector.onExec != nil {
		return c.connector.onExec(query, args)
	}
	return driver.RowsAffected(1), nil
}

type sqlServerTestTx struct{ connector sqlServerTestConnector }

func (t *sqlServerTestTx) Commit() error {
	if t.connector.onCommit != nil {
		return t.connector.onCommit()
	}
	return nil
}
func (t *sqlServerTestTx) Rollback() error { return nil }

func (t *sqlServerTestTx) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return handleSQLServerTestQuery(t.connector, query, args)
}

func (t *sqlServerTestTx) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if t.connector.onExec != nil {
		return t.connector.onExec(query, args)
	}
	return driver.RowsAffected(1), nil
}

func handleSQLServerTestQuery(c sqlServerTestConnector, query string, args []driver.NamedValue) (driver.Rows, error) {
	if c.queries != nil {
		c.queries.Add(1)
	}
	if c.onQuery != nil {
		rows, err := c.onQuery(query, args)
		if rows != nil || err != nil {
			return rows, err
		}
	}
	if c.catalog != nil {
		if rows, ok := c.catalog.Query(query, args); ok {
			return rows, nil
		}
	}
	return nil, sql.ErrNoRows
}

type sqlServerCatalogMock struct {
	group   *dbtarget.SQLServerManagedInspection
	receipt *dbtarget.SQLServerManagedInspection
	current *dbtarget.SQLServerManagedInspection
}

func newSQLServerCatalogMock(group, receipt *dbtarget.SQLServerManagedInspection) *sqlServerCatalogMock {
	return &sqlServerCatalogMock{group: group, receipt: receipt}
}

func (m *sqlServerCatalogMock) Query(query string, args []driver.NamedValue) (driver.Rows, bool) {
	lower := strings.ToLower(strings.Join(strings.Fields(query), " "))
	switch {
	case strings.Contains(lower, "serverproperty('productversion')"):
		return &sqlServerMockRows{
			columns: []string{"version", "edition", "engine_edition", "database", "compatibility"},
			rows:    [][]driver.Value{{"16.0.4180.1", "Developer Edition (64-bit)", int64(3), "gwtest", int64(160)}},
		}, true
	case strings.Contains(lower, "select schema_id("):
		return &sqlServerMockRows{
			columns: []string{"schema_id"},
			rows:    [][]driver.Value{{int64(1)}},
		}, true
	case strings.Contains(lower, "sys.extended_properties"):
		rows := &sqlServerMockRows{columns: []string{"repair_operation_id"}}
		if m.current != nil && m.current.Evidence.RepairOperationID != "" {
			rows.rows = [][]driver.Value{{m.current.Evidence.RepairOperationID}}
		}
		return rows, true
	case strings.Contains(lower, "sys.objects"):
		tableName := ""
		if len(args) >= 2 {
			if s, ok := args[1].Value.(string); ok {
				tableName = s
			}
		}
		if tableName == m.receipt.Inspection.Table {
			m.current = m.receipt
		} else {
			m.current = m.group
		}
		evidence := m.current.Evidence
		return &sqlServerMockRows{
			columns: []string{"object_id", "schema_name", "object_name", "type_desc", "is_user_table", "is_memory_optimized", "temporal_type", "is_filetable", "is_external", "is_node", "is_edge"},
			rows: [][]driver.Value{{
				evidence.ObjectID, evidence.Schema, m.current.Inspection.Table, evidence.TableType,
				true, false, int64(0), false, false, false, false,
			}},
		}, true
	case strings.Contains(lower, "sys.indexes"):
		if m.current == nil {
			return nil, false
		}
		return m.indexesRows(m.current), true
	case strings.Contains(lower, "sys.columns"):
		if m.current == nil {
			return nil, false
		}
		return m.columnsRows(m.current), true
	case strings.Contains(lower, "sys.check_constraints"):
		if m.current == nil {
			return nil, false
		}
		if m.current.Evidence.ProvenanceJSONChecked {
			return &sqlServerMockRows{
				columns: []string{"name", "definition", "is_disabled", "is_not_trusted", "parent_column_id"},
				rows: [][]driver.Value{{
					"CK_" + m.current.Inspection.Table + "_provenance", "([ISJSON]([provenance])=(1))", false, false, int64(0),
				}},
			}, true
		}
		return &sqlServerMockRows{
			columns: []string{"name", "definition", "is_disabled", "is_not_trusted", "parent_column_id"},
			rows:    [][]driver.Value{},
		}, true
	}
	return nil, false
}

func (m *sqlServerCatalogMock) columnsRows(inspection *dbtarget.SQLServerManagedInspection) driver.Rows {
	evidence := inspection.Evidence
	rows := make([][]driver.Value, 0, len(inspection.Inspection.Columns))
	for i, col := range inspection.Inspection.Columns {
		colID := int64(i + 1)
		typeName, maxLen, prec, scale := sqlServerTypeDetails(col.DataType)
		var collation driver.Value
		if c, ok := evidence.ColumnCollations[col.Name]; ok {
			collation = c
		}
		var defVal driver.Value
		if d, ok := evidence.ColumnDefaults[col.Name]; ok && d != nil {
			defVal = "((N'" + strings.ReplaceAll(*d, "'", "''") + "'))"
		}
		rows = append(rows, []driver.Value{
			colID, col.Name, typeName, "sys", false, maxLen, prec, scale,
			col.Nullable, false, false, int64(0), false, false, false, collation, defVal,
		})
	}
	return &sqlServerMockRows{
		columns: []string{"column_id", "name", "type_name", "type_schema", "is_user_defined", "max_length", "precision", "scale", "is_nullable", "is_computed", "is_identity", "generated_always_type", "is_hidden", "is_sparse", "is_rowguidcol", "collation_name", "default_definition"},
		rows:    rows,
	}
}

func (m *sqlServerCatalogMock) indexesRows(inspection *dbtarget.SQLServerManagedInspection) driver.Rows {
	evidence := inspection.Evidence
	var rows [][]driver.Value
	colIDs := make(map[string]int64, len(inspection.Inspection.Columns))
	for i, col := range inspection.Inspection.Columns {
		colIDs[col.Name] = int64(i + 1)
	}
	appendIdx := func(id int64, name string, idxType int64, unique, primary bool, cols []string) {
		for ord, cName := range cols {
			cID, exists := colIDs[cName]
			if !exists {
				continue
			}
			rows = append(rows, []driver.Value{
				id, name, idxType, unique, primary, false, false, false, false, "FG",
				int64(ord + 1), cID, false, false, cName, false,
			})
		}
	}
	if len(evidence.PrimaryKey) > 0 {
		appendIdx(1, "PK_"+inspection.Inspection.Table, 1, true, true, evidence.PrimaryKey)
	}
	for i, info := range inspection.Indexes {
		appendIdx(int64(i+2), info.Name, 2, info.Unique, false, info.Columns)
	}
	cols := []string{"index_id", "name", "type", "is_unique", "is_primary_key", "has_filter", "is_disabled", "ignore_dup_key", "is_hypothetical", "data_space_type", "key_ordinal", "index_column_id", "is_included_column", "is_descending_key", "column_name", "is_computed"}
	return &sqlServerMockRows{
		columns: cols,
		rows:    rows,
	}
}

func sqlServerTypeDetails(declared string) (typeName string, maxLength, precision, scale driver.Value) {
	base, arguments, hasArguments := strings.Cut(declared, "(")
	base = strings.ToLower(base)
	if !hasArguments {
		switch base {
		case "bigint":
			return base, int64(8), int64(19), int64(0)
		case "bit":
			return base, int64(1), int64(1), int64(0)
		default:
			return base, int64(8), int64(53), int64(0)
		}
	}
	arguments = strings.TrimSuffix(arguments, ")")
	switch base {
	case "char":
		width, _ := strconv.ParseInt(arguments, 10, 64)
		return base, width, int64(0), int64(0)
	case "nvarchar":
		if strings.EqualFold(arguments, "max") {
			return base, int64(-1), int64(0), int64(0)
		}
		width, _ := strconv.ParseInt(arguments, 10, 64)
		return base, width * 2, int64(0), int64(0)
	case "datetime2":
		scale, _ := strconv.ParseInt(arguments, 10, 64)
		return base, int64(8), int64(27), scale
	case "decimal":
		precision, scaleText, _ := strings.Cut(arguments, ",")
		precisionVal, _ := strconv.ParseInt(strings.TrimSpace(precision), 10, 64)
		scaleVal, _ := strconv.ParseInt(strings.TrimSpace(scaleText), 10, 64)
		return base, int64(17), precisionVal, scaleVal
	default:
		return base, int64(8), int64(53), int64(0)
	}
}
