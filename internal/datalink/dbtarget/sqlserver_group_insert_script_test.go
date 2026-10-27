package dbtarget

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/require"
)

type sqlServerGroupInsertScript struct {
	mu                  sync.Mutex
	txGate              sync.Mutex
	groupTable          string
	groupColumns        []string
	effectKey           string
	bucketStart         time.Time
	recordID            string
	groupCatalog        *sqlServerManagedCatalogScript
	receiptCatalog      *sqlServerManagedCatalogScript
	rowExists           bool
	rowValues           []driver.Value
	receiptExists       bool
	receiptKey          string
	receiptDigest       string
	receiptOwner        string
	failReceiptInsert   bool
	failGroupInsertOnce bool
	missingSchema       bool
	commitAckLost       bool
	ackLost             bool
	beginCalls          int
	commitCalls         int
	rollbackCalls       int
	groupInsertArgs     [][]driver.NamedValue
	receiptInsertArgs   [][]driver.NamedValue
	queries             []string
	receiptLookupArgs   [][]driver.NamedValue
	groupLookupArgs     [][]driver.NamedValue
	catalogLookups      []string
	execStatements      []string
	forbiddenStatements []string
}

type sqlServerGroupInsertConnector struct{ script *sqlServerGroupInsertScript }
type sqlServerGroupInsertDriver struct{}

type sqlServerGroupInsertConn struct {
	script  *sqlServerGroupInsertScript
	catalog *sqlServerManagedCatalogScript
	tx      *sqlServerGroupInsertTx
}

type sqlServerGroupInsertTx struct {
	script              *sqlServerGroupInsertScript
	rowExists           bool
	rowValues           []driver.Value
	receiptExists       bool
	receiptKey          string
	receiptDigest       string
	receiptOwner        string
	pendingRow          []driver.Value
	pendingReceiptKey   string
	pendingReceiptValue string
	closed              bool
}

var (
	_ driver.Connector      = sqlServerGroupInsertConnector{}
	_ driver.Conn           = (*sqlServerGroupInsertConn)(nil)
	_ driver.ConnBeginTx    = (*sqlServerGroupInsertConn)(nil)
	_ driver.ExecerContext  = (*sqlServerGroupInsertConn)(nil)
	_ driver.QueryerContext = (*sqlServerGroupInsertConn)(nil)
)

func newSQLServerGroupInsertScript(req GroupInsertRequest) *sqlServerGroupInsertScript {
	columns := make([]string, len(req.Row.Cells))
	for index, cell := range req.Row.Cells {
		columns[index] = cell.Column
	}
	return &sqlServerGroupInsertScript{
		groupTable:     req.TableName,
		groupColumns:   columns,
		effectKey:      req.Row.EffectKey,
		bucketStart:    req.Row.BucketStart,
		recordID:       req.Row.RecordID,
		groupCatalog:   sqlServerInsertCatalogScript(req.SQLServer.Group),
		receiptCatalog: sqlServerInsertCatalogScript(req.SQLServer.Receipt),
		receiptOwner:   "gateway",
	}
}

func openSQLServerGroupInsertDB(t *testing.T, script *sqlServerGroupInsertScript) *sql.DB {
	t.Helper()
	db := sql.OpenDB(sqlServerGroupInsertConnector{script: script})
	db.SetMaxOpenConns(4)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return db
}

func (connector sqlServerGroupInsertConnector) Connect(context.Context) (driver.Conn, error) {
	return &sqlServerGroupInsertConn{script: connector.script}, nil
}

func (sqlServerGroupInsertConnector) Driver() driver.Driver { return sqlServerGroupInsertDriver{} }

func (sqlServerGroupInsertDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("scripted SQL Server adapter connector must be used")
}

func (*sqlServerGroupInsertConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepared SQL Server adapter statements are not scripted")
}

func (*sqlServerGroupInsertConn) Close() error { return nil }

func (conn *sqlServerGroupInsertConn) Begin() (driver.Tx, error) {
	return conn.BeginTx(context.Background(), driver.TxOptions{})
}

func (conn *sqlServerGroupInsertConn) BeginTx(ctx context.Context, _ driver.TxOptions) (driver.Tx, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	conn.script.txGate.Lock()
	conn.script.mu.Lock()
	conn.script.beginCalls++
	tx := &sqlServerGroupInsertTx{
		script: conn.script, rowExists: conn.script.rowExists,
		rowValues:     cloneSQLServerDriverValues(conn.script.rowValues),
		receiptExists: conn.script.receiptExists, receiptKey: conn.script.receiptKey,
		receiptDigest: conn.script.receiptDigest, receiptOwner: conn.script.receiptOwner,
	}
	conn.script.mu.Unlock()
	conn.tx = tx
	return tx, nil
}

func (conn *sqlServerGroupInsertConn) QueryContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Rows, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if conn.tx == nil {
		return nil, errors.New("SQL Server adapter query outside transaction")
	}
	normalized := strings.Join(strings.Fields(query), " ")
	lower := strings.ToLower(normalized)
	conn.script.mu.Lock()
	conn.script.queries = append(conn.script.queries, normalized)
	conn.script.mu.Unlock()
	switch {
	case strings.Contains(lower, "with (updlock, holdlock)") && strings.Contains(lower, "[effect_key] = @p1"):
		conn.script.mu.Lock()
		conn.script.receiptLookupArgs = append(conn.script.receiptLookupArgs, cloneSQLServerNamedValues(args))
		conn.script.mu.Unlock()
		if len(args) != 1 || args[0].Value != conn.tx.receiptKey && conn.tx.receiptExists {
			return nil, errors.New("receipt key lookup was not parameterized")
		}
		if !conn.tx.receiptExists || args[0].Value != conn.tx.receiptKey {
			return sqlServerRecordingRows([]string{"payload_digest"}), nil
		}
		return sqlServerRecordingRows(
			[]string{"payload_digest", sqlServerManagedReceiptOwnerColumn},
			[]driver.Value{conn.tx.receiptDigest, conn.tx.receiptOwner},
		), nil
	case strings.Contains(lower, "with (updlock, holdlock)") && strings.Contains(lower, "[bucket_start] = @p1") && strings.Contains(lower, "[record_id] = @p2"):
		conn.script.mu.Lock()
		conn.script.groupLookupArgs = append(conn.script.groupLookupArgs, cloneSQLServerNamedValues(args))
		conn.script.mu.Unlock()
		if len(args) != 2 {
			return nil, errors.New("group key lookup was not parameterized")
		}
		if !conn.tx.rowExists || !sqlServerScriptedGroupIdentityMatches(conn.script, conn.tx.rowValues, args) {
			return sqlServerRecordingRows(conn.script.groupColumns), nil
		}
		return sqlServerRecordingRows(conn.script.groupColumns, conn.tx.rowValues), nil
	case normalized == sqlServerRecordingMetadataQuery:
		return sqlServerRecordingRows(
			[]string{"version", "edition", "engine_edition", "database", "compatibility"},
			[]driver.Value{"16.0.1000.6", "Developer Edition (64-bit)", int64(3), "fixture_db", int64(160)},
		), nil
	case normalized == "SELECT SCHEMA_ID(@p1)":
		if conn.script.missingSchema {
			return sqlServerRecordingRows([]string{"schema_id"}, []driver.Value{nil}), nil
		}
		return sqlServerRecordingRows([]string{"schema_id"}, []driver.Value{int64(1)}), nil
	case strings.Contains(lower, "sys.extended_properties"):
		if conn.catalog == nil {
			return nil, fmt.Errorf("catalog marker lookup preceded object lookup")
		}
		return (&sqlServerManagedCatalogConn{script: conn.catalog}).QueryContext(ctx, query, args)
	case strings.Contains(lower, "sys.objects"):
		if len(args) != 2 {
			return nil, errors.New("catalog object lookup was not parameterized")
		}
		tableName, ok := args[1].Value.(string)
		if !ok {
			return nil, errors.New("catalog table lookup was not text")
		}
		switch tableName {
		case conn.script.groupTable:
			conn.catalog = conn.script.groupCatalog
		case EffectReceiptTable:
			conn.catalog = conn.script.receiptCatalog
		default:
			return nil, errors.New("unexpected catalog table")
		}
		conn.script.mu.Lock()
		conn.script.catalogLookups = append(conn.script.catalogLookups, tableName)
		conn.script.mu.Unlock()
		return (&sqlServerManagedCatalogConn{script: conn.catalog}).QueryContext(ctx, query, args)
	default:
		if conn.catalog == nil {
			return nil, fmt.Errorf("catalog detail query preceded object lookup")
		}
		return (&sqlServerManagedCatalogConn{script: conn.catalog}).QueryContext(ctx, query, args)
	}
}

func (conn *sqlServerGroupInsertConn) ExecContext(
	ctx context.Context,
	query string,
	args []driver.NamedValue,
) (driver.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if conn.tx == nil {
		return nil, errors.New("SQL Server adapter exec outside transaction")
	}
	normalized := strings.Join(strings.Fields(query), " ")
	lower := strings.ToLower(normalized)
	conn.script.mu.Lock()
	conn.script.execStatements = append(conn.script.execStatements, normalized)
	conn.script.mu.Unlock()
	for _, token := range []string{"create table", "alter table", "drop table", "truncate table", "merge ", "on conflict", "on duplicate key", "insert ignore"} {
		if strings.Contains(lower, token) {
			conn.script.mu.Lock()
			conn.script.forbiddenStatements = append(conn.script.forbiddenStatements, normalized)
			conn.script.mu.Unlock()
			return nil, errors.New("forbidden SQL statement")
		}
	}
	if strings.Contains(lower, "insert into [dbo].["+strings.ToLower(conn.script.groupTable)+"]") {
		if conn.script.failGroupInsertOnce {
			conn.script.mu.Lock()
			conn.script.failGroupInsertOnce = false
			conn.script.mu.Unlock()
			return nil, mssql.Error{Number: 1205, Message: "scripted SQL Server deadlock"}
		}
		conn.script.mu.Lock()
		conn.script.groupInsertArgs = append(conn.script.groupInsertArgs, cloneSQLServerNamedValues(args))
		conn.script.mu.Unlock()
		conn.tx.pendingRow = sqlServerNamedValues(args)
		return driver.RowsAffected(1), nil
	}
	if strings.Contains(lower, "insert into [dbo].["+strings.ToLower(EffectReceiptTable)+"]") {
		if conn.script.failReceiptInsert {
			return nil, errors.New("scripted receipt insert failure")
		}
		if len(args) != 4 {
			return nil, errors.New("receipt insert had an unexpected argument count")
		}
		conn.script.mu.Lock()
		conn.script.receiptInsertArgs = append(conn.script.receiptInsertArgs, cloneSQLServerNamedValues(args))
		conn.script.mu.Unlock()
		conn.tx.pendingReceiptKey, _ = args[0].Value.(string)
		conn.tx.pendingReceiptValue, _ = args[1].Value.(string)
		return driver.RowsAffected(1), nil
	}
	return nil, errors.New("unexpected SQL Server adapter statement")
}

func (tx *sqlServerGroupInsertTx) Commit() error {
	if tx.closed {
		return sql.ErrTxDone
	}
	tx.closed = true
	tx.script.mu.Lock()
	if tx.pendingRow != nil {
		tx.script.rowExists = true
		tx.script.rowValues = cloneSQLServerDriverValues(tx.pendingRow)
	}
	if tx.pendingReceiptKey != "" {
		tx.script.receiptExists = true
		tx.script.receiptKey = tx.pendingReceiptKey
		tx.script.receiptDigest = tx.pendingReceiptValue
		tx.script.receiptOwner = "gateway"
	}
	tx.script.commitCalls++
	lostAck := tx.script.commitAckLost && !tx.script.ackLost
	if lostAck {
		tx.script.ackLost = true
	}
	tx.script.mu.Unlock()
	tx.script.txGate.Unlock()
	if lostAck {
		return errors.New("scripted lost SQL Server commit acknowledgment")
	}
	return nil
}

func (tx *sqlServerGroupInsertTx) Rollback() error {
	if tx.closed {
		return sql.ErrTxDone
	}
	tx.closed = true
	tx.script.mu.Lock()
	tx.script.rollbackCalls++
	tx.script.mu.Unlock()
	tx.script.txGate.Unlock()
	return nil
}

func sqlServerScriptedGroupIdentityMatches(
	script *sqlServerGroupInsertScript,
	values []driver.Value,
	args []driver.NamedValue,
) bool {
	bucketIndex, recordIndex := -1, -1
	for index, name := range script.groupColumns {
		switch name {
		case longColBucketStart:
			bucketIndex = index
		case longColRecordID:
			recordIndex = index
		}
	}
	if bucketIndex < 0 || recordIndex < 0 || len(values) != len(script.groupColumns) {
		return false
	}
	return groupCellEqual(values[bucketIndex], args[0].Value) && values[recordIndex] == args[1].Value
}

func sqlServerNamedValues(args []driver.NamedValue) []driver.Value {
	values := make([]driver.Value, len(args))
	for index, arg := range args {
		values[index] = cloneSQLServerDriverValue(arg.Value)
	}
	return values
}

func cloneSQLServerNamedValues(args []driver.NamedValue) []driver.NamedValue {
	cloned := append([]driver.NamedValue(nil), args...)
	for index := range cloned {
		cloned[index].Value = cloneSQLServerDriverValue(cloned[index].Value)
	}
	return cloned
}

func cloneSQLServerDriverValues(values []driver.Value) []driver.Value {
	cloned := make([]driver.Value, len(values))
	for index := range values {
		cloned[index] = cloneSQLServerDriverValue(values[index])
	}
	return cloned
}

func cloneSQLServerDriverValue(value driver.Value) driver.Value {
	if bytes, ok := value.([]byte); ok {
		return append([]byte(nil), bytes...)
	}
	return value
}

func sqlServerInsertCatalogScript(inspection *SQLServerManagedInspection) *sqlServerManagedCatalogScript {
	metadata := sqlServerRecordingMetadataFixture()
	evidence := inspection.Evidence
	script := &sqlServerManagedCatalogScript{
		metadata: metadata, repairOperationID: evidence.RepairOperationID,
		objectRows: [][]driver.Value{{
			evidence.ObjectID, evidence.Schema, inspection.Inspection.Table, evidence.TableType,
			true, false, int64(0), false, false, false, false,
		}},
		columnRows: make([][]driver.Value, 0, len(inspection.Inspection.Columns)),
	}
	columnIDs := make(map[string]int64, len(inspection.Inspection.Columns))
	for index, column := range inspection.Inspection.Columns {
		columnID := int64(index + 1)
		columnIDs[column.Name] = columnID
		typeName, maxLength, precision, scale := sqlServerInsertCatalogType(column.DataType)
		var collation, defaultDefinition driver.Value
		if value, ok := evidence.ColumnCollations[column.Name]; ok {
			collation = value
		}
		if value := evidence.ColumnDefaults[column.Name]; value != nil {
			defaultDefinition = "((N'" + strings.ReplaceAll(*value, "'", "''") + "'))"
		}
		script.columnRows = append(script.columnRows, []driver.Value{
			columnID, column.Name, typeName, "sys", false, maxLength, precision, scale,
			column.Nullable, false, false, int64(0), false, false, false, collation, defaultDefinition,
		})
	}
	appendIndex := func(id int64, name string, indexType int64, unique, primary bool, columns []string) {
		for ordinal, column := range columns {
			columnID, exists := columnIDs[column]
			if !exists {
				continue
			}
			script.indexRows = append(script.indexRows, []driver.Value{
				id, name, indexType, unique, primary, false, false, false, false, "FG",
				int64(ordinal + 1), columnID, false, false, column, false,
			})
		}
	}
	if len(evidence.PrimaryKey) > 0 {
		appendIndex(1, "PK_"+inspection.Inspection.Table, 1, true, true, evidence.PrimaryKey)
	}
	for index, info := range inspection.Indexes {
		appendIndex(int64(index+2), info.Name, 2, info.Unique, false, info.Columns)
	}
	if evidence.ProvenanceJSONChecked {
		script.checkRows = [][]driver.Value{{
			"CK_" + inspection.Inspection.Table + "_provenance", "([ISJSON]([provenance])=(1))", false, false, int64(0),
		}}
	}
	return script
}

func sqlServerInsertCatalogType(declared string) (typeName string, maxLengthValue, precisionValue, scaleValue driver.Value) {
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
		if arguments == "MAX" {
			return base, int64(-1), int64(0), int64(0)
		}
		width, _ := strconv.ParseInt(arguments, 10, 64)
		return base, width * 2, int64(0), int64(0)
	case "decimal":
		precisionText, scaleText, _ := strings.Cut(arguments, ",")
		precision, _ := strconv.ParseInt(precisionText, 10, 64)
		scale, _ := strconv.ParseInt(scaleText, 10, 64)
		length := int64(17)
		switch {
		case precision <= 9:
			length = 5
		case precision <= 19:
			length = 9
		case precision <= 28:
			length = 13
		}
		return base, length, precision, scale
	case "datetime2":
		scale, _ := strconv.ParseInt(arguments, 10, 64)
		return base, int64(8), int64(27), scale
	case "float":
		precision, _ := strconv.ParseInt(arguments, 10, 64)
		return base, int64(8), precision, int64(0)
	default:
		return base, int64(0), int64(0), int64(0)
	}
}

func sqlServerGroupInsertSnapshot(req GroupInsertRequest) []driver.Value {
	values := make([]driver.Value, len(req.Row.Cells))
	for index, cell := range req.Row.Cells {
		values[index] = normalizeDBValue(cell.Value)
	}
	return values
}
