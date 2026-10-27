package grouptestwrite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type sqlServerExecutionScript struct {
	columns  []string
	rows     [][]driver.Value
	receipts map[string]string
}

type sqlServerExecutionConnector struct{ script *sqlServerExecutionScript }

func (c sqlServerExecutionConnector) Connect(context.Context) (driver.Conn, error) {
	return &sqlServerExecutionConn{script: c.script}, nil
}
func (sqlServerExecutionConnector) Driver() driver.Driver { return sqlServerExecutionDriver{} }

type sqlServerExecutionDriver struct{}

func (sqlServerExecutionDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("open is not used by Connector")
}

type sqlServerExecutionConn struct{ script *sqlServerExecutionScript }

func (*sqlServerExecutionConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not used")
}
func (*sqlServerExecutionConn) Close() error { return nil }
func (*sqlServerExecutionConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transaction is not used")
}
func (c *sqlServerExecutionConn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if strings.Contains(strings.ToLower(query), "effect_key") {
		effectKey, _ := args[0].Value.(string)
		digest, ok := c.script.receipts[effectKey]
		if !ok {
			return &sqlServerExecutionRows{columns: []string{"payload_digest"}}, nil
		}
		return &sqlServerExecutionRows{columns: []string{"payload_digest"}, rows: [][]driver.Value{{digest}}}, nil
	}
	selected := splitSQLServerExecutionSelect(query)
	if len(selected) == 0 {
		selected = append([]string(nil), c.script.columns...)
	}
	index := make(map[string]int, len(c.script.columns))
	for i, column := range c.script.columns {
		index[strings.ToLower(column)] = i
	}
	rows := make([][]driver.Value, 0, len(c.script.rows))
	for _, source := range c.script.rows {
		row := make([]driver.Value, len(selected))
		for i, column := range selected {
			at, ok := index[strings.ToLower(column)]
			if !ok {
				return nil, fmt.Errorf("test driver missing selected column %q", column)
			}
			row[i] = source[at]
		}
		rows = append(rows, row)
	}
	return &sqlServerExecutionRows{columns: selected, rows: rows}, nil
}

func splitSQLServerExecutionSelect(query string) []string {
	start := strings.Index(query, ") ")
	if start < 0 {
		return nil
	}
	body := query[start+2:]
	end := strings.Index(strings.ToUpper(body), " FROM ")
	if end < 0 {
		return nil
	}
	parts := strings.Split(body[:end], ",")
	columns := make([]string, 0, len(parts))
	for _, part := range parts {
		columns = append(columns, strings.Trim(strings.TrimSpace(part), "[]"))
	}
	return columns
}

type sqlServerExecutionRows struct {
	columns []string
	rows    [][]driver.Value
	index   int
}

func (r *sqlServerExecutionRows) Columns() []string { return r.columns }
func (*sqlServerExecutionRows) Close() error        { return nil }
func (r *sqlServerExecutionRows) Next(dest []driver.Value) error {
	if r.index >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.index])
	r.index++
	return nil
}

func sqlServerExecutionService(t *testing.T) (*Service, *plan) {
	t.Helper()
	group, layout, tagTypes := sqlServerTestWriteFixture(t)
	groups := &sqlServerTestWriteGroups{fakeGroups: &fakeGroups{group: group}, layout: layout, testOnlySupportsGroupTestWrite: true}
	svc := New(Dependencies{
		Groups: groups, Tags: fakeTags(tagTypes), Ledger: recordingplan.NewService(recordingplan.NewMemoryRepository()),
	}, Config{StepTimeout: time.Second, Now: func() time.Time { return time.Date(2026, 10, 26, 3, 4, 5, 123456789, time.UTC) }})
	p, err := svc.buildSQLServerTestPlan(t.Context(), group, &schema.DatabaseConnector{
		ID: "conn-1", Kind: schema.DatabaseConnectorKindSQLServer, IdentityRevision: "crev-1", Enabled: true,
	})
	require.NoError(t, err)
	return svc, p
}

func sqlServerExecutionToken(t *testing.T, svc *Service, p *plan) (*recordingplan.SchemaOperation, time.Time) {
	t.Helper()
	scope := recordingplan.TestWriteScope{
		WorkspaceID: p.group.WorkspaceID, WorkspaceRevision: "setup-1", GroupID: p.group.ID, GroupRevision: p.group.Revision,
		ConnectorID: "conn-1", ConnectorRevision: "crev-1", Dialect: "sqlserver", Database: "db",
		Schema: p.group.Destination.TableSchema, Table: p.group.Destination.TableName,
	}
	ledger := svc.deps.Ledger
	token, err := ledger.PrepareTestWritePreview(t.Context(), scope, p.digest)
	require.NoError(t, err)
	op, _, err := ledger.ClaimTestWrite(t.Context(), token)
	require.NoError(t, err)
	return op, token.CreatedAt
}

func sqlServerExecutionColumns(p *plan) []string {
	columns := make([]string, 0, len(p.columns)+1)
	seen := make(map[string]struct{}, len(p.columns)+1)
	add := func(column string) {
		key := strings.ToLower(column)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		columns = append(columns, column)
	}
	add(mysqlTestRecordColumn)
	for _, column := range p.columns {
		add(column.Name)
	}
	return columns
}

func sqlServerExecutionRowValues(p *plan, state progress, row dbtarget.EncodedRow, columns []string) []driver.Value {
	cells := make(map[string]any, len(row.Cells))
	for _, cell := range row.Cells {
		cells[strings.ToLower(cell.Column)] = cell.Value
	}
	values := make([]driver.Value, len(columns))
	for i, column := range columns {
		key := strings.ToLower(column)
		switch key {
		case "record_id":
			values[i] = row.RecordID
		case "workspace_id":
			values[i] = state.WorkspaceID
		case "group_id":
			values[i] = state.GroupID
		case "group_revision":
			values[i] = state.GroupRevision
		case strings.ToLower(p.ownerColumn):
			values[i] = "gateway"
		default:
			values[i] = cells[key]
		}
	}
	return values
}

func TestSQLServerExactReadbackTypes(t *testing.T) {
	decimal, err := measurement.NewDecimal("12000.4")
	require.NoError(t, err)
	require.True(t, sqlServerTestValueMatches("DECIMAL(38,18)", measurement.ExactDecimal, decimal, []byte("12000.400000000000000000")))
	require.True(t, sqlServerTestValueMatches("BIT", measurement.ExactBool, measurement.NewBool(true), int64(1)))
	require.True(t, sqlServerTestValueMatches("NVARCHAR(MAX)", measurement.ExactText, measurement.NewText("運轉中🟢"), []byte("運轉中🟢")))
	require.True(t, sqlServerTestValueMatches("DECIMAL(38,18)", measurement.ExactUint64, measurement.NewUint64(9007199254740993), []byte("9007199254740993")))
}

func sqlServerExecutionProgress(p *plan, op *recordingplan.SchemaOperation, rows []dbtarget.EncodedRow) progress {
	state := progress{
		Phase: phaseClaimed, OwnerValue: ownerValue(op.OperationID), Schema: p.group.Destination.TableSchema,
		Table: p.group.Destination.TableName, OwnerColumn: p.ownerColumn, TableOwnerValue: "gateway",
		ProvenanceColumn: p.group.RowPolicy.ProvenanceColumn, WorkspaceID: p.group.WorkspaceID,
		GroupID: p.group.ID, GroupRevision: p.group.Revision, RecordIDs: make([]string, len(rows)),
		SampleIDs: make([]string, len(rows)), EffectKeys: make([]string, len(rows)),
		EffectDigests: make(map[string]string, len(rows)),
	}
	for i, row := range rows {
		state.RecordIDs[i], state.SampleIDs[i], state.EffectKeys[i] = row.RecordID, fmt.Sprintf("%s-%d", state.OwnerValue, i), row.EffectKey
		digest, err := groupdelivery.RowPayloadDigest(row)
		if err == nil {
			state.EffectDigests[row.EffectKey] = digest
		}
	}
	if len(rows) > 0 {
		state.EffectKey = rows[0].EffectKey
	}
	return state
}

func sqlServerExecutionScriptForRows(t *testing.T, p *plan, state progress, rows []dbtarget.EncodedRow) *sqlServerExecutionScript {
	t.Helper()
	columns := sqlServerExecutionColumns(p)
	script := &sqlServerExecutionScript{columns: columns, receipts: map[string]string{}}
	for _, row := range rows {
		script.rows = append(script.rows, sqlServerExecutionRowValues(p, state, row, columns))
		digest, err := groupdelivery.RowPayloadDigest(row)
		require.NoError(t, err)
		script.receipts[row.EffectKey] = digest
	}
	return script
}

func TestSQLServerReadbackRejectsWrongBucketTimestamp(t *testing.T) {
	svc, p := sqlServerExecutionService(t)
	op, at := sqlServerExecutionToken(t, svc, p)
	rows, err := p.sqlserverRows(ownerValue(op.OperationID), at)
	require.NoError(t, err)
	state := sqlServerExecutionProgress(p, op, rows)
	script := sqlServerExecutionScriptForRows(t, p, state, rows)
	bucketIndex := -1
	for i, column := range script.columns {
		if strings.EqualFold(column, "bucket_start") {
			bucketIndex = i
			break
		}
	}
	require.NotEqual(t, -1, bucketIndex)
	script.rows[0][bucketIndex] = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	db := sql.OpenDB(sqlServerExecutionConnector{script: script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	result := svc.verifySQLServer(t.Context(), &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindSQLServer}, p, state, rows)
	require.Equal(t, WriteUnverified, result.outcome, "wrong bucket_start must not be silently ignored")
}

func TestSQLServerReadbackRejectsNonNullUnusedSlot(t *testing.T) {
	svc, p := sqlServerExecutionService(t)
	op, at := sqlServerExecutionToken(t, svc, p)
	rows, err := p.sqlserverRows(ownerValue(op.OperationID), at)
	require.NoError(t, err)
	state := sqlServerExecutionProgress(p, op, rows)
	script := sqlServerExecutionScriptForRows(t, p, state, rows)
	unused := ""
	for _, column := range p.columns {
		if strings.HasPrefix(column.Name, "value_") && column.Name != sqlServerValueSlot(p.sqlserver.members[0].kind) {
			unused = column.Name
			break
		}
	}
	require.NotEmpty(t, unused)
	unusedIndex := -1
	for i, column := range script.columns {
		if strings.EqualFold(column, unused) {
			unusedIndex = i
			break
		}
	}
	require.NotEqual(t, -1, unusedIndex)
	script.rows[0][unusedIndex] = "unexpected"
	db := sql.OpenDB(sqlServerExecutionConnector{script: script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	result := svc.verifySQLServer(t.Context(), &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindSQLServer}, p, state, rows)
	require.Equal(t, WriteUnverified, result.outcome, "non-NULL unused slot must not be silently ignored")
}

func TestSQLServerReadbackRejectsMissingExpectedCell(t *testing.T) {
	svc, p := sqlServerExecutionService(t)
	op, at := sqlServerExecutionToken(t, svc, p)
	rows, err := p.sqlserverRows(ownerValue(op.OperationID), at)
	require.NoError(t, err)
	state := sqlServerExecutionProgress(p, op, rows)
	missing := *p
	missing.columns = make([]dbtarget.ColumnInfo, 0, len(p.columns)-1)
	for _, column := range p.columns {
		if !strings.EqualFold(column.Name, "bucket_start") {
			missing.columns = append(missing.columns, column)
		}
	}
	script := sqlServerExecutionScriptForRows(t, &missing, state, rows)
	db := sql.OpenDB(sqlServerExecutionConnector{script: script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	result := svc.verifySQLServer(t.Context(), &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindSQLServer}, &missing, state, rows)
	require.Equal(t, WriteUnverified, result.outcome, "an expected cell absent from readback must fail")
}

func TestSQLServerWritePhasePersistsExactRowsAndReadsBack(t *testing.T) {
	svc, p := sqlServerExecutionService(t)
	op, at := sqlServerExecutionToken(t, svc, p)
	rows, err := p.sqlserverRows(ownerValue(op.OperationID), at)
	require.NoError(t, err)
	require.Len(t, rows, len(p.sqlserver.members))
	state := progress{Phase: phaseClaimed, OwnerValue: ownerValue(op.OperationID)}
	columns := sqlServerExecutionColumns(p)
	script := &sqlServerExecutionScript{columns: columns, receipts: map[string]string{}}
	for _, row := range rows {
		script.rows = append(script.rows, sqlServerExecutionRowValues(p, progress{
			WorkspaceID: p.group.WorkspaceID, GroupID: p.group.ID, GroupRevision: p.group.Revision,
		}, row, columns))
	}
	for _, row := range rows {
		digest, err := groupdelivery.RowPayloadDigest(row)
		require.NoError(t, err)
		script.receipts[row.EffectKey] = digest
	}
	db := sql.OpenDB(sqlServerExecutionConnector{script: script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	inserts := 0
	svc.insertRow = func(_ context.Context, _ *sql.DB, req dbtarget.GroupInsertRequest) (dbtarget.GroupInsertResult, error) {
		inserts++
		require.Equal(t, dbtarget.GroupEffectReceipt, req.Strategy)
		require.Same(t, p.sqlserver.proof, req.SQLServer)
		require.Equal(t, "2026-10-26T03:04:05.123456Z", req.CommittedAt)
		return dbtarget.GroupInsertResult{}, nil
	}
	written, finished, err := svc.writeSQLServerPhase(t.Context(), op, p, at, &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindSQLServer}, &state)
	require.NoError(t, err)
	require.Nil(t, finished)
	require.Equal(t, WriteVerified, written.outcome, written.reason)
	require.Equal(t, len(rows), inserts)
	require.Equal(t, rows[0].RecordID, state.RecordIDs[0])
	require.Equal(t, ownerValue(op.OperationID)+"-0", state.SampleIDs[0])
	require.Equal(t, rows[0].EffectKey, state.EffectKey)
}

func TestSQLServerWritePhaseRestartInWritingReadsWithoutInsert(t *testing.T) {
	svc, p := sqlServerExecutionService(t)
	op, at := sqlServerExecutionToken(t, svc, p)
	rows, err := p.sqlserverRows(ownerValue(op.OperationID), at)
	require.NoError(t, err)
	state := progress{
		Phase: phaseWriting, OwnerValue: ownerValue(op.OperationID), Schema: p.group.Destination.TableSchema,
		Table: p.group.Destination.TableName, OwnerColumn: p.ownerColumn, TableOwnerValue: "gateway",
		ProvenanceColumn: p.group.RowPolicy.ProvenanceColumn, WorkspaceID: p.group.WorkspaceID, GroupID: p.group.ID,
		GroupRevision: p.group.Revision, RecordIDs: make([]string, len(rows)), SampleIDs: make([]string, len(rows)),
		EffectKeys: make([]string, len(rows)), EffectDigests: make(map[string]string, len(rows)),
	}
	columns := sqlServerExecutionColumns(p)
	script := &sqlServerExecutionScript{columns: columns, receipts: map[string]string{}}
	for i, row := range rows {
		state.RecordIDs[i], state.SampleIDs[i], state.EffectKeys[i] = row.RecordID, fmt.Sprintf("%s-%d", state.OwnerValue, i), row.EffectKey
		digest, err := groupdelivery.RowPayloadDigest(row)
		require.NoError(t, err)
		state.EffectDigests[row.EffectKey] = digest
		script.receipts[row.EffectKey] = digest
		script.rows = append(script.rows, sqlServerExecutionRowValues(p, state, row, columns))
	}
	db := sql.OpenDB(sqlServerExecutionConnector{script: script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	svc.insertRow = func(context.Context, *sql.DB, dbtarget.GroupInsertRequest) (dbtarget.GroupInsertResult, error) {
		t.Fatal("restart in writing phase must not insert again")
		return dbtarget.GroupInsertResult{}, nil
	}
	written, finished, err := svc.writeSQLServerPhase(t.Context(), op, p, at, &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindSQLServer}, &state)
	require.NoError(t, err)
	require.Nil(t, finished)
	require.Equal(t, WriteVerified, written.outcome, written.reason)
}

func TestSQLServerFirstInsertPreCommitFailureIsTerminalFailed(t *testing.T) {
	svc, p := sqlServerExecutionService(t)
	op, at := sqlServerExecutionToken(t, svc, p)
	state := progress{Phase: phaseClaimed, OwnerValue: ownerValue(op.OperationID)}
	columns := sqlServerExecutionColumns(p)
	db := sql.OpenDB(sqlServerExecutionConnector{script: &sqlServerExecutionScript{columns: columns, receipts: map[string]string{}}})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	inserts := 0
	svc.insertRow = func(context.Context, *sql.DB, dbtarget.GroupInsertRequest) (dbtarget.GroupInsertResult, error) {
		inserts++
		return dbtarget.GroupInsertResult{}, &dbtarget.GroupInsertError{Phase: dbtarget.InsertPhasePreCommit, Cause: dbtarget.ErrExistingRowDiffers}
	}
	written, finished, err := svc.writeSQLServerPhase(t.Context(), op, p, at, &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindSQLServer}, &state)
	require.NoError(t, err)
	require.Equal(t, 1, inserts)
	require.Zero(t, written.outcome)
	require.NotNil(t, finished)
	require.Equal(t, recordingplan.SchemaOperationFailed, finished.Status)
	require.Equal(t, WriteFailed, finished.WriteOutcome)
	require.Equal(t, reasonRowRejected, finished.Reason)
	require.Equal(t, CleanupNotAttempted, finished.CleanupStatus)
}

func TestSQLServerSecondInsertFailureReportsPartial(t *testing.T) {
	svc, p := sqlServerExecutionService(t)
	op, at := sqlServerExecutionToken(t, svc, p)
	state := progress{Phase: phaseClaimed, OwnerValue: ownerValue(op.OperationID)}
	columns := sqlServerExecutionColumns(p)
	script := &sqlServerExecutionScript{columns: columns, receipts: map[string]string{}}
	db := sql.OpenDB(sqlServerExecutionConnector{script: script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	inserts := 0
	svc.insertRow = func(_ context.Context, _ *sql.DB, req dbtarget.GroupInsertRequest) (dbtarget.GroupInsertResult, error) {
		inserts++
		if inserts == 2 {
			return dbtarget.GroupInsertResult{}, &dbtarget.GroupInsertError{Phase: dbtarget.InsertPhasePreCommit, Cause: errors.New("second row failed")}
		}
		script.rows = append(script.rows, sqlServerExecutionRowValues(p, progress{
			WorkspaceID: p.group.WorkspaceID, GroupID: p.group.ID, GroupRevision: p.group.Revision,
		}, req.Row, columns))
		digest, err := groupdelivery.RowPayloadDigest(req.Row)
		require.NoError(t, err)
		script.receipts[req.Row.EffectKey] = digest
		return dbtarget.GroupInsertResult{}, nil
	}
	written, finished, err := svc.writeSQLServerPhase(t.Context(), op, p, at, &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindSQLServer}, &state)
	require.NoError(t, err)
	require.Nil(t, finished)
	require.Equal(t, 2, inserts)
	require.Equal(t, WriteUnverified, written.outcome)
}

// TestSQLServerSharedWritePhaseDispatch is the behavioral gap gate: the shared
// writePhase must dispatch SQL Server plans to writeSQLServerPhase rather than
// dereferencing the generic single-row plan layout.
func TestSQLServerSharedWritePhaseDispatch(t *testing.T) {
	svc, p := sqlServerExecutionService(t)
	op, at := sqlServerExecutionToken(t, svc, p)
	rows, err := p.sqlserverRows(ownerValue(op.OperationID), at)
	require.NoError(t, err)
	state := progress{Phase: phaseClaimed, OwnerValue: ownerValue(op.OperationID)}
	columns := sqlServerExecutionColumns(p)
	script := &sqlServerExecutionScript{columns: columns, receipts: map[string]string{}}
	for _, row := range rows {
		script.rows = append(script.rows, sqlServerExecutionRowValues(p, progress{
			WorkspaceID: p.group.WorkspaceID, GroupID: p.group.ID, GroupRevision: p.group.Revision,
		}, row, columns))
		digest, err := groupdelivery.RowPayloadDigest(row)
		require.NoError(t, err)
		script.receipts[row.EffectKey] = digest
	}
	db := sql.OpenDB(sqlServerExecutionConnector{script: script})
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	svc.insertRow = func(context.Context, *sql.DB, dbtarget.GroupInsertRequest) (dbtarget.GroupInsertResult, error) {
		return dbtarget.GroupInsertResult{}, nil
	}
	require.NotPanics(t, func() {
		written, finished, err := svc.writePhase(t.Context(), op, p, at, &dbtarget.OpenedDestination{DB: db, Kind: schema.DatabaseConnectorKindSQLServer}, &state)
		require.NoError(t, err)
		require.Nil(t, finished)
		require.Equal(t, WriteVerified, written.outcome, written.reason)
	})
}
