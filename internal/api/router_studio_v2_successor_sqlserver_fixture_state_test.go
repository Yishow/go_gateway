package api

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
)

// Hash full local rows before assertions: connector JSON must never appear in
// failure output. No content is omitted from equality checks.
func sqlServerSuccessorLocalState(t *testing.T, db *sql.DB) map[string][]string {
	t.Helper()
	state := successorAPIAllState(t, db)
	for name, rows := range state {
		for i, row := range rows {
			digest := sha256.Sum256([]byte(row))
			rows[i] = hex.EncodeToString(digest[:])
		}
		sort.Strings(rows)
		state[name] = rows
	}
	return state
}

func (r *sqlServerSuccessorRig) targetState(t *testing.T) map[string][]string {
	t.Helper()
	state := map[string][]string{}
	for i, db := range []*sql.DB{r.f.target, r.targetDB} {
		fixture := r.fixtures[i]
		scope := fixture.Manifest.Database + "." + fixture.Manifest.Schema
		rows, err := db.QueryContext(t.Context(), `SELECT name FROM sys.tables WHERE schema_id=SCHEMA_ID(@p1) ORDER BY name`, fixture.Manifest.Schema)
		require.NoError(t, err)
		var names []string
		for rows.Next() {
			var name string
			require.NoError(t, rows.Scan(&name))
			names = append(names, name)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
		for _, name := range names {
			require.Regexp(t, `^(gw_long_[0-9a-f]{16}|gw_effect_receipts)$`, name)
			key := scope + "." + name
			state[key] = sqlServerSuccessorQueryRows(t, db, "SELECT * FROM "+sqlServerQualified(fixture.Manifest.Schema, name))
			// Columns/defaults, keys/indexes and partition catalog are part of the
			// snapshot, so a row-count-only false success cannot hide DDL damage.
			for label, query := range map[string]string{
				"columns":       `SELECT c.*,d.definition FROM sys.columns c LEFT JOIN sys.default_constraints d ON d.object_id=c.default_object_id WHERE c.object_id=OBJECT_ID(@p1)`,
				"indexes":       `SELECT * FROM sys.indexes WHERE object_id=OBJECT_ID(@p1)`,
				"index_columns": `SELECT * FROM sys.index_columns WHERE object_id=OBJECT_ID(@p1)`,
			} {
				state[key+"/"+label] = sqlServerSuccessorQueryRows(t, db, query, fixture.Manifest.Schema+"."+name)
			}
		}
		state[scope+"/partition_functions"] = sqlServerSuccessorQueryRows(t, db, `SELECT * FROM sys.partition_functions`)
		state[scope+"/partition_schemes"] = sqlServerSuccessorQueryRows(t, db, `SELECT * FROM sys.partition_schemes`)
		state[scope+"/partition_ranges"] = sqlServerSuccessorQueryRows(t, db, `SELECT * FROM sys.partition_range_values`)
	}
	return state
}

func sqlServerSuccessorQueryRows(t *testing.T, db *sql.DB, query string, args ...any) []string {
	t.Helper()
	rows, err := db.QueryContext(t.Context(), query, args...)
	require.NoError(t, err)
	defer rows.Close()
	columns, err := rows.Columns()
	require.NoError(t, err)
	bodies := []string{}
	for rows.Next() {
		values, refs := make([]any, len(columns)), make([]any, len(columns))
		for i := range values {
			refs[i] = &values[i]
		}
		require.NoError(t, rows.Scan(refs...))
		raw, err := json.Marshal(values)
		require.NoError(t, err)
		bodies = append(bodies, string(raw))
	}
	require.NoError(t, rows.Err())
	sort.Strings(bodies)
	return bodies
}

// The existing circuit cleanup owns plain tables. Native cleanup additionally
// proves the exact created table and exclusive PF/PS before dropping any of it.
// Receipt storage is shared and never dropped. Runtime tests clean their exact
// effect rows first via cleanupSQLServerCircuitEffects (LIFO cleanup ordering).
func registerSQLServerSuccessorSchemaCleanup(t *testing.T, f managedGroupFixture, token recordingplan.SchemaPreviewToken) {
	t.Helper()
	require.True(t, sqlServerCircuitPreviewCreatesGroup(token), "cleanup only owns newly created schema")
	targets := dbtarget.NewConnectorService(dbtarget.NewSQLConnectorRepository(f.db))
	registerSQLServerSuccessorReceiptCleanup(t, f, targets, token)
	if token.GroupLayout.SQLServer.NativeMonthly == nil {
		registerSQLServerCircuitCreatedTableCleanup(t, f, targets, token)
		return
	}
	capability, err := targets.ProbeSQLServerRecordingAtRevision(t.Context(), token.ConnectorID, token.ConnectorRevision)
	require.NoError(t, err)
	group, err := dbtarget.InspectOpenSQLServerManagedTable(t.Context(), f.target, capability, token.GroupLayout.TableName)
	require.NoError(t, err)
	receipt, err := dbtarget.InspectOpenSQLServerManagedTable(t.Context(), f.target, capability, dbtarget.EffectReceiptTable)
	require.NoError(t, err)
	_, err = recordingplan.VerifyManagedSQLServerCatalog(*token.GroupLayout, token.Database, sqlServerCircuitTargetInspection(group), sqlServerCircuitTargetInspection(receipt))
	require.NoError(t, err)
	objectID := group.Evidence.ObjectID
	native := token.GroupLayout.SQLServer.NativeMonthly
	require.Equal(t, "pf_"+token.GroupLayout.TableName, native.FunctionName)
	require.Equal(t, "ps_"+token.GroupLayout.TableName, native.SchemeName)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		current, err := dbtarget.InspectOpenSQLServerManagedTable(ctx, f.target, capability, token.GroupLayout.TableName)
		require.NoError(t, err)
		require.Equal(t, objectID, current.Evidence.ObjectID, "a replacement object is not owned")
		receipts, err := dbtarget.InspectOpenSQLServerManagedTable(ctx, f.target, capability, dbtarget.EffectReceiptTable)
		require.NoError(t, err)
		_, err = recordingplan.VerifyManagedSQLServerCatalog(*token.GroupLayout, token.Database, sqlServerCircuitTargetInspection(current), sqlServerCircuitTargetInspection(receipts))
		require.NoError(t, err)
		table := sqlServerQualified(token.Schema, token.GroupLayout.TableName)
		var count int
		require.NoError(t, f.target.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count))
		require.Zero(t, count, "leave nonempty storage intact until owned effects are removed")
		_, err = f.target.ExecContext(ctx, "IF OBJECT_ID(@p1,N'U')=@p2 AND NOT EXISTS(SELECT 1 FROM "+table+") BEGIN DROP TABLE "+table+"; DROP PARTITION SCHEME ["+native.SchemeName+"]; DROP PARTITION FUNCTION ["+native.FunctionName+"]; END ELSE THROW 50000,'fixture identity changed',1", token.Schema+"."+token.GroupLayout.TableName, objectID)
		require.NoError(t, err)
	})
}

// Unknown external outcomes are not an ownership grant. Register before DDL:
// later proven cleanup runs first; a partial/unknown leftover is reported and
// preserved for explicit ledger reconciliation, never guessed away.
func registerSQLServerSuccessorPreparationAudit(t *testing.T, f managedGroupFixture, token recordingplan.SchemaPreviewToken) {
	t.Helper()
	require.True(t, sqlServerCircuitPreviewCreatesGroup(token))
	var before int
	require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM sys.objects WHERE schema_id=SCHEMA_ID(@p1) AND name=@p2`, token.Schema, token.GroupLayout.TableName).Scan(&before))
	require.Zero(t, before, "a foreign/preexisting target cannot be acquired by the fixture")
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var remaining int
		require.NoError(t, f.target.QueryRowContext(ctx, `SELECT COUNT(*) FROM sys.objects WHERE schema_id=SCHEMA_ID(@p1) AND name=@p2`, token.Schema, token.GroupLayout.TableName).Scan(&remaining))
		require.Zero(t, remaining, "CLEANUP BLOCKED: unproven/changed schema was retained; reconcile the exact schema ledger before removal")
	})
}

func registerSQLServerSuccessorReceiptCleanup(t *testing.T, f managedGroupFixture, targets *dbtarget.ConnectorService, token recordingplan.SchemaPreviewToken) {
	t.Helper()
	created := false
	for _, table := range token.Tables {
		created = created || table.Name == dbtarget.EffectReceiptTable && table.Action == recordingplan.SchemaTableActionCreate
	}
	if !created {
		return
	} // A shared preexisting receipt table is never dropped.
	capability, err := targets.ProbeSQLServerRecordingAtRevision(t.Context(), token.ConnectorID, token.ConnectorRevision)
	require.NoError(t, err)
	initial, err := dbtarget.InspectOpenSQLServerManagedTable(t.Context(), f.target, capability, dbtarget.EffectReceiptTable)
	require.NoError(t, err)
	require.NotNil(t, initial.Evidence)
	require.Positive(t, initial.Evidence.ObjectID)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		current, err := dbtarget.InspectOpenSQLServerManagedTable(ctx, f.target, capability, dbtarget.EffectReceiptTable)
		require.NoError(t, err)
		require.Equal(t, initial, current, "changed/replaced receipt storage is not safe to remove")
		table := sqlServerQualified(token.Schema, dbtarget.EffectReceiptTable)
		var count int
		require.NoError(t, f.target.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table).Scan(&count))
		require.Zero(t, count, "foreign or undelivered effects retain their receipt storage")
		_, err = f.target.ExecContext(ctx, "IF OBJECT_ID(@p1,N'U')=@p2 AND NOT EXISTS(SELECT 1 FROM "+table+") DROP TABLE "+table+" ELSE THROW 50000,'fixture receipt changed',1", token.Schema+"."+dbtarget.EffectReceiptTable, initial.Evidence.ObjectID)
		require.NoError(t, err)
	})
}
