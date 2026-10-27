package dbtarget

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/stretchr/testify/require"
)

func TestSQLServerGroupInsertCommitsRowAndReceiptAtomically(t *testing.T) {
	req := sqlServerGroupInsertTestRequest()
	script := newSQLServerGroupInsertScript(req)
	db := openSQLServerGroupInsertDB(t, script)

	result, err := InsertGroupRow(t.Context(), db, req)
	require.NoError(t, err)
	require.False(t, result.AlreadyCommitted)
	require.Equal(t, 1, script.beginCalls)
	require.Equal(t, 1, script.commitCalls)
	require.Zero(t, script.rollbackCalls)
	require.True(t, script.rowExists)
	require.True(t, script.receiptExists)
	require.Equal(t, req.Row.EffectKey, script.receiptKey)
	require.Equal(t, req.PayloadDigest, script.receiptDigest)
	require.Len(t, script.groupInsertArgs, 1)
	require.Len(t, script.receiptInsertArgs, 1)
	require.Empty(t, script.forbiddenStatements)
	requireSQLServerInsertLockOrder(t, script)
}

func TestSQLServerReceiptConcurrent(t *testing.T) {
	req := sqlServerGroupInsertTestRequest()
	script := newSQLServerGroupInsertScript(req)
	db := openSQLServerGroupInsertDB(t, script)
	start := make(chan struct{})
	type outcome struct {
		result GroupInsertResult
		err    error
	}
	outcomes := make(chan outcome, 2)
	var wait sync.WaitGroup
	for range 2 {
		wait.Go(func() {
			<-start
			result, err := InsertGroupRow(t.Context(), db, req)
			outcomes <- outcome{result: result, err: err}
		})
	}
	close(start)
	wait.Wait()
	close(outcomes)

	alreadyCommitted := 0
	for result := range outcomes {
		require.NoError(t, result.err)
		if result.result.AlreadyCommitted {
			alreadyCommitted++
		}
	}
	require.Equal(t, 1, alreadyCommitted)
	require.Len(t, script.groupInsertArgs, 1)
	require.Len(t, script.receiptInsertArgs, 1)
	require.Equal(t, 1, script.commitCalls)
	require.Equal(t, 1, script.rollbackCalls)
	requireSQLServerInsertLockOrder(t, script)
}

func TestSQLServerCommitAckLost(t *testing.T) {
	req := sqlServerGroupInsertTestRequest()
	script := newSQLServerGroupInsertScript(req)
	script.commitAckLost = true
	db := openSQLServerGroupInsertDB(t, script)

	_, err := InsertGroupRow(t.Context(), db, req)
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.True(t, insertErr.Ambiguous())
	require.True(t, script.rowExists)
	require.True(t, script.receiptExists)

	result, err := InsertGroupRow(t.Context(), db, req)
	require.NoError(t, err)
	require.True(t, result.AlreadyCommitted)
	require.Len(t, script.groupInsertArgs, 1, "retry resolves the ambiguous commit from its receipt")
	require.Len(t, script.receiptInsertArgs, 1)
	require.Equal(t, 1, script.commitCalls)
	require.Equal(t, 1, script.rollbackCalls)
}

func TestSQLServerGroupInsertReceiptAndExistingRowOutcomes(t *testing.T) {
	for _, test := range []struct {
		name            string
		receipt         bool
		digest          string
		row             bool
		receiptOwner    string
		mutateStoredRow bool
		wantAlready     bool
		wantError       error
	}{
		{name: "same digest and matching row", receipt: true, digest: strings.Repeat("c", 64), row: true, wantAlready: true},
		{name: "same digest but foreign owner", receipt: true, digest: strings.Repeat("c", 64), row: true, receiptOwner: "foreign", wantError: ErrGroupInsertUnsupported},
		{name: "same digest but different row", receipt: true, digest: strings.Repeat("c", 64), row: true, mutateStoredRow: true, wantError: ErrExistingRowDiffers},
		{name: "different digest", receipt: true, digest: strings.Repeat("d", 64), row: true, wantError: ErrReceiptDigestMismatch},
		{name: "receipt without row", receipt: true, digest: strings.Repeat("c", 64), wantError: ErrExistingRowDiffers},
		{name: "row without receipt", row: true, wantError: ErrExistingRowDiffers},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := sqlServerGroupInsertTestRequest()
			script := newSQLServerGroupInsertScript(req)
			script.receiptExists, script.receiptKey, script.receiptDigest = test.receipt, req.Row.EffectKey, test.digest
			if test.receiptOwner != "" {
				script.receiptOwner = test.receiptOwner
			}
			script.rowExists = test.row
			if test.row {
				script.rowValues = sqlServerGroupInsertSnapshot(req)
				if test.mutateStoredRow {
					script.rowValues[indexSQLServerTestColumn(script.groupColumns, longColTagID)] = "changed"
				}
			}
			db := openSQLServerGroupInsertDB(t, script)

			result, err := InsertGroupRow(t.Context(), db, req)
			if test.wantError != nil {
				require.ErrorIs(t, err, test.wantError)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.wantAlready, result.AlreadyCommitted)
			}
			require.Empty(t, script.groupInsertArgs)
			require.Empty(t, script.receiptInsertArgs)
			require.Equal(t, 1, script.rollbackCalls)
			require.Zero(t, script.commitCalls)
			requireSQLServerInsertLockOrder(t, script)
		})
	}
}

func TestSQLServerGroupInsertPartialReceiptFailureRollsBackRow(t *testing.T) {
	req := sqlServerGroupInsertTestRequest()
	script := newSQLServerGroupInsertScript(req)
	script.failReceiptInsert = true
	db := openSQLServerGroupInsertDB(t, script)

	_, err := InsertGroupRow(t.Context(), db, req)
	require.Error(t, err)
	require.False(t, script.rowExists)
	require.False(t, script.receiptExists)
	require.Equal(t, 1, script.rollbackCalls)
	require.Zero(t, script.commitCalls)
	require.Len(t, script.groupInsertArgs, 1)
	require.Empty(t, script.receiptInsertArgs)
}

func TestSQLServerAdapterDeadlockRetryRetainsFrozenRequest(t *testing.T) {
	req := sqlServerGroupInsertTestRequest()
	proofBefore, err := json.Marshal(req.SQLServer)
	require.NoError(t, err)
	rowBefore := req.Row
	rowBefore.Cells = slices.Clone(req.Row.Cells)
	script := newSQLServerGroupInsertScript(req)
	script.failGroupInsertOnce = true
	db := openSQLServerGroupInsertDB(t, script)

	_, err = InsertGroupRow(t.Context(), db, req)
	var insertErr *GroupInsertError
	require.ErrorAs(t, err, &insertErr)
	require.Equal(t, InsertPhasePreCommit, insertErr.Phase)
	var serverErr mssql.Error
	require.ErrorAs(t, err, &serverErr)
	require.EqualValues(t, 1205, serverErr.Number)
	require.False(t, errors.Is(err, ErrGroupInsertUnsupported))
	require.False(t, script.rowExists)
	require.False(t, script.receiptExists)

	result, err := InsertGroupRow(t.Context(), db, req)
	require.NoError(t, err)
	require.False(t, result.AlreadyCommitted)
	proofAfter, err := json.Marshal(req.SQLServer)
	require.NoError(t, err)
	require.Equal(t, proofBefore, proofAfter)
	require.Equal(t, rowBefore, req.Row)
	require.True(t, script.rowExists)
	require.True(t, script.receiptExists)
	require.Len(t, script.groupInsertArgs, 1)
}

func TestSQLServerGroupInsertBindsUnicodeAndExactNumericParameters(t *testing.T) {
	for _, test := range []struct {
		name      string
		valueType string
		column    string
		value     string
	}{
		{name: "uint64", valueType: "uint64", column: longColValueDecimal, value: "18446744073709551615"},
		{name: "decimal", valueType: "decimal", column: longColValueDecimal, value: "12345678901234567890.123456789012345678"},
		{name: "unicode text", valueType: "text", column: longColValueText, value: "電力🟢"},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := sqlServerGroupInsertTestRequest()
			sqlServerTestCell(req.Row.Cells, longColValueType).Value = test.valueType
			sqlServerTestCell(req.Row.Cells, longColStatus).Value = "ok"
			sqlServerTestCell(req.Row.Cells, test.column).Value = test.value
			sqlServerTestCell(req.Row.Cells, longColProvenance).Value = `{"type":"` + test.valueType + `","member":"member-1","profile":{"kind":"interval","interval":60,"unit":"minute"}}`
			script := newSQLServerGroupInsertScript(req)
			db := openSQLServerGroupInsertDB(t, script)

			_, err := InsertGroupRow(t.Context(), db, req)
			require.NoError(t, err)
			valueIndex := indexSQLServerTestColumn(script.groupColumns, test.column)
			got, ok := script.groupInsertArgs[0][valueIndex].Value.(string)
			require.True(t, ok, "exact numerics and Unicode text remain bound strings")
			require.Equal(t, test.value, got)
			tagIndex := indexSQLServerTestColumn(script.groupColumns, longColTagID)
			require.Equal(t, "Tag-電力🟢", script.groupInsertArgs[0][tagIndex].Value)
			committedAt, ok := script.receiptInsertArgs[0][2].Value.(time.Time)
			require.True(t, ok)
			require.True(t, sqlServerUTCMicrosecond(committedAt))
		})
	}
}

func TestSQLServerGroupInsertBindsAndReplaysEmptyExactText(t *testing.T) {
	req := sqlServerGroupInsertTestRequest()
	sqlServerTestCell(req.Row.Cells, longColValueType).Value = "text"
	sqlServerTestCell(req.Row.Cells, longColStatus).Value = "ok"
	sqlServerTestCell(req.Row.Cells, longColValueText).Value = ""
	sqlServerTestCell(req.Row.Cells, longColProvenance).Value = `{"type":"text","member":"member-1","profile":{"kind":"interval","interval":60,"unit":"minute"}}`
	script := newSQLServerGroupInsertScript(req)
	db := openSQLServerGroupInsertDB(t, script)

	first, err := InsertGroupRow(t.Context(), db, req)
	require.NoError(t, err)
	require.False(t, first.AlreadyCommitted)
	textIndex := indexSQLServerTestColumn(script.groupColumns, longColValueText)
	boundText, ok := script.groupInsertArgs[0][textIndex].Value.(string)
	require.True(t, ok)
	require.Empty(t, boundText)
	require.Equal(t, "", script.rowValues[textIndex], "scripted readback preserves empty text, not NULL")

	replay, err := InsertGroupRow(t.Context(), db, req)
	require.NoError(t, err)
	require.True(t, replay.AlreadyCommitted)
	require.Len(t, script.groupInsertArgs, 1)
	require.Len(t, script.receiptInsertArgs, 1)
}

func TestSQLServerGroupInsertRejectsRowAndScopeMismatchBeforeBegin(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*GroupInsertRequest)
	}{
		{name: "foreign group owner marker", mutate: func(req *GroupInsertRequest) {
			sqlServerTestCell(req.Row.Cells, sqlServerInsertTestOwner).Value = "foreign"
		}},
		{name: "request schema differs from proof", mutate: func(req *GroupInsertRequest) { req.SchemaName = "other" }},
		{name: "request table differs from proof", mutate: func(req *GroupInsertRequest) { req.TableName = "other_table" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := sqlServerGroupInsertTestRequest()
			test.mutate(&req)
			script := newSQLServerGroupInsertScript(req)
			db := openSQLServerGroupInsertDB(t, script)

			_, err := InsertGroupRow(t.Context(), db, req)
			require.ErrorIs(t, err, ErrGroupInsertUnsupported)
			require.Zero(t, script.beginCalls)
			require.Empty(t, script.execStatements)
		})
	}
}

func TestSQLServerGroupInsertRejectsRepairMarkerDriftInCatalog(t *testing.T) {
	req := sqlServerGroupInsertTestRequest()
	req.SQLServer.Group.Evidence.RepairOperationID = "repair-operation-test"
	script := newSQLServerGroupInsertScript(req)
	script.groupCatalog.repairOperationID = "another-repair-operation"
	db := openSQLServerGroupInsertDB(t, script)

	_, err := InsertGroupRow(t.Context(), db, req)
	require.ErrorIs(t, err, ErrGroupInsertUnsupported)
	require.False(t, script.rowExists)
	require.False(t, script.receiptExists)
	require.Empty(t, script.groupInsertArgs)
	require.Empty(t, script.receiptInsertArgs)
}

func TestSQLServerGroupInsertFailsClosedOnSchemaAndCatalogDrift(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*sqlServerGroupInsertScript)
	}{
		{name: "schema missing", mutate: func(script *sqlServerGroupInsertScript) { script.missingSchema = true }},
		{name: "group object replaced", mutate: func(script *sqlServerGroupInsertScript) { script.groupCatalog.objectRows[0][0] = int64(777) }},
		{name: "receipt object replaced", mutate: func(script *sqlServerGroupInsertScript) { script.receiptCatalog.objectRows[0][0] = int64(778) }},
		{name: "group owner default changed", mutate: func(script *sqlServerGroupInsertScript) {
			for _, row := range script.groupCatalog.columnRows {
				if row[1] == sqlServerInsertTestOwner {
					row[16] = "((N'foreign'))"
					return
				}
			}
			t.Fatal("group catalog fixture is missing the owner column")
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			req := sqlServerGroupInsertTestRequest()
			script := newSQLServerGroupInsertScript(req)
			test.mutate(script)
			db := openSQLServerGroupInsertDB(t, script)

			_, err := InsertGroupRow(t.Context(), db, req)
			require.ErrorIs(t, err, ErrGroupInsertUnsupported)
			require.False(t, script.rowExists)
			require.False(t, script.receiptExists)
			require.Empty(t, script.groupInsertArgs)
			require.Empty(t, script.receiptInsertArgs)
		})
	}
}

func requireSQLServerInsertLockOrder(t *testing.T, script *sqlServerGroupInsertScript) {
	t.Helper()
	require.GreaterOrEqual(t, len(script.queries), 6)
	receiptLookup := strings.ToLower(script.queries[0])
	groupLookup := strings.ToLower(script.queries[1])
	require.Contains(t, receiptLookup, "from [dbo].[gw_effect_receipts] with (updlock, holdlock)")
	require.Contains(t, receiptLookup, "where [effect_key] = @p1")
	require.Contains(t, groupLookup, "from [dbo].[gw_long_fixture] with (updlock, holdlock)")
	require.Contains(t, groupLookup, "where [bucket_start] = @p1 and [record_id] = @p2")
	require.Len(t, script.receiptLookupArgs, script.beginCalls)
	for _, args := range script.receiptLookupArgs {
		require.Len(t, args, 1)
		require.Equal(t, script.effectKey, args[0].Value)
	}
	require.Len(t, script.groupLookupArgs, script.beginCalls)
	for _, args := range script.groupLookupArgs {
		require.Len(t, args, 2)
		bucketStart, ok := args[0].Value.(time.Time)
		require.True(t, ok)
		require.True(t, bucketStart.Equal(script.bucketStart))
		require.Equal(t, script.recordID, args[1].Value)
	}
	require.NotContains(t, receiptLookup, "cast(")
	require.NotContains(t, groupLookup, "convert(")
	firstCatalogLookup := slices.IndexFunc(script.queries, func(query string) bool {
		return strings.Contains(strings.ToLower(query), "sys.objects")
	})
	require.Greater(t, firstCatalogLookup, 1, "both equality locks precede the catalog reads")
	require.Len(t, script.catalogLookups, script.beginCalls*2)
	for index := range len(script.catalogLookups) / 2 {
		require.Equal(t, "gw_long_fixture", script.catalogLookups[index*2])
		require.Equal(t, EffectReceiptTable, script.catalogLookups[index*2+1])
	}
	require.Contains(t, script.groupCatalog.queries[0].query, "sys.objects")
	require.Contains(t, script.receiptCatalog.queries[0].query, "sys.objects")
	require.Empty(t, script.forbiddenStatements)
	for _, query := range script.queries {
		lower := strings.ToLower(query)
		require.NotContains(t, lower, "sp_getapplock")
		require.NotContains(t, lower, "tablock")
		require.NotContains(t, lower, "xlock")
	}
	for _, statement := range script.execStatements {
		lower := strings.ToLower(statement)
		require.Contains(t, lower, "insert into")
		require.NotContains(t, lower, "merge")
		require.NotContains(t, lower, "create table")
		require.NotContains(t, lower, "on conflict")
	}
}

func indexSQLServerTestColumn(columns []string, name string) int {
	for index, column := range columns {
		if column == name {
			return index
		}
	}
	panic("test column not found: " + name)
}
