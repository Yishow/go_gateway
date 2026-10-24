package api

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// The guards of the timescale mode: nothing existing is converted, stale
// evidence is refused before DDL, the scope is not caller-specifiable, and a
// drifted hypertable is never ready.

// TestTimescaleNeverConvertsAnExistingTable: a table that already holds the
// timescale scope's name is never turned into a hypertable. A foreign table, a
// populated or empty owned-shaped plain table are all refused, nothing is
// altered, converted or deleted, and a table that appears after preview is
// refused at confirmation with a stale/conflict answer and zero DDL.
func TestTimescaleNeverConvertsAnExistingTable(t *testing.T) {
	setup := func(t *testing.T) (managedGroupFixture, string, string) {
		f := newManagedDestinationFixture(t, "postgres")
		if !timescaleInstalled(t, f) {
			t.Skip("this destination database has no TimescaleDB extension; the case is NOT RUN here")
		}
		_, created := ensureTimescaleGroup(t, f)
		_, group := groupSaveData(t, created)
		return f, group["id"].(string), destinationOf(group)["table_name"].(string)
	}
	state := func(t *testing.T, f managedGroupFixture, table string) (rows int, hypertable bool, columns int) {
		t.Helper()
		qualified := `"` + f.namespace + `"."` + table + `"`
		require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM `+qualified).Scan(&rows))
		require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) > 0 FROM timescaledb_information.hypertables WHERE hypertable_schema = $1 AND hypertable_name = $2`, f.namespace, table).Scan(&hypertable))
		require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2`, f.namespace, table).Scan(&columns))
		return rows, hypertable, columns
	}
	requireRefused := func(t *testing.T, f managedGroupFixture, groupID, table string) {
		t.Helper()
		rows, _, columns := state(t, f, table)
		response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-preview", scopeBody(t, f, groupID))
		require.Equal(t, http.StatusUnprocessableEntity, response.Code, response.Body.String())
		require.Equal(t, "RECORDING_SCHEMA_INCOMPATIBLE", errorOf(t, decodeJSONBody(t, response))["code"])
		afterRows, hypertable, afterColumns := state(t, f, table)
		require.Equal(t, [3]any{rows, false, columns}, [3]any{afterRows, hypertable, afterColumns}, "the table was neither converted nor altered")
	}
	qualifiedOf := func(f managedGroupFixture, table string) string { return `"` + f.namespace + `"."` + table + `"` }

	t.Run("a populated foreign table", func(t *testing.T) {
		f, groupID, table := setup(t)
		_, err := f.target.ExecContext(t.Context(), `CREATE TABLE `+qualifiedOf(f, table)+` (id integer primary key, note text)`)
		require.NoError(t, err)
		_, err = f.target.ExecContext(t.Context(), `INSERT INTO `+qualifiedOf(f, table)+` VALUES (1, 'foreign')`)
		require.NoError(t, err)
		requireRefused(t, f, groupID, table)
	})

	for name, populated := range map[string]bool{"an empty owned-shaped plain table": false, "a populated owned-shaped plain table": true} {
		t.Run(name, func(t *testing.T) {
			f, groupID, table := setup(t)
			token := previewLong(t, f, groupID)
			require.Len(t, token.Statements, 4)
			// Only the CREATE TABLE of the layout: exactly the owned columns, but never made a hypertable.
			_, err := f.target.ExecContext(t.Context(), `SET search_path TO "`+f.namespace+`"; `+token.Statements[0])
			require.NoError(t, err)
			if populated {
				insertRequiredColumns(t, f, table)
			}
			requireRefused(t, f, groupID, table)
		})
	}

	t.Run("a table that appears after preview", func(t *testing.T) {
		f, groupID, table := setup(t)
		token := previewLong(t, f, groupID)
		require.NotEmpty(t, token.Statements)
		_, err := f.target.ExecContext(t.Context(), `CREATE TABLE `+qualifiedOf(f, table)+` (id integer primary key)`)
		require.NoError(t, err)
		request := scopeBody(t, f, groupID)
		request["token"], request["operation_id"] = token.Token, token.OperationID
		response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-apply", request)
		require.GreaterOrEqual(t, response.Code, 400, response.Body.String())
		_, hypertable, columns := state(t, f, table)
		require.False(t, hypertable, "no if_not_exists success and no conversion of a table that appeared")
		require.Equal(t, 1, columns, "the table that appeared is untouched")
		require.Empty(t, destinationOf(getGroup(t, f, groupID))["schema_revision"], "an unconfirmed result must not be saved as confirmed")
	})
}

// insertRequiredColumns writes one row that fills every NOT NULL column without
// a default, with a value of the column's own type, so a populated table can be
// built from the layout without naming its 29 columns here.
func insertRequiredColumns(t *testing.T, f managedGroupFixture, table string) {
	t.Helper()
	rows, err := f.target.QueryContext(t.Context(), `SELECT column_name, data_type FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2 AND is_nullable = 'NO' AND column_default IS NULL ORDER BY ordinal_position`, f.namespace, table)
	require.NoError(t, err)
	var names, values []string
	for rows.Next() {
		var name, dataType string
		require.NoError(t, rows.Scan(&name, &dataType))
		literal := map[string]string{
			"text": "'x'", "character varying": "'x'", "timestamp with time zone": "now()", "integer": "0", "bigint": "0", "smallint": "0",
			"double precision": "0", "numeric": "0", "boolean": "false", "jsonb": "'{}'::jsonb",
		}[dataType]
		require.NotEmpty(t, literal, "no sample literal for %s %s", name, dataType)
		names, values = append(names, `"`+name+`"`), append(values, literal)
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	_, err = f.target.ExecContext(t.Context(), `INSERT INTO "`+f.namespace+`"."`+table+`" (`+strings.Join(names, ",")+`) VALUES (`+strings.Join(values, ",")+`)`)
	require.NoError(t, err)
}

// TestTimescaleStaleEvidenceRejectedBeforeMutation: changing anything the
// preview froze — the workspace, the source, or the connector revision —
// rejects the confirmation with a conflict before any DDL runs, and the
// hypertable is never created. The storage mode cannot change on a saved group
// (the next test), so it cannot go stale between preview and confirmation.
func TestTimescaleStaleEvidenceRejectedBeforeMutation(t *testing.T) {
	for _, scenario := range []string{"workspace", "source", "connector"} {
		t.Run(scenario, func(t *testing.T) {
			f := newManagedDestinationFixture(t, "postgres")
			if !timescaleInstalled(t, f) {
				t.Skip("this destination database has no TimescaleDB extension; the case is NOT RUN here")
			}
			_, created := ensureTimescaleGroup(t, f)
			_, group := groupSaveData(t, created)
			groupID, table := group["id"].(string), destinationOf(group)["table_name"].(string)
			token := previewLong(t, f, groupID)
			request := scopeBody(t, f, groupID)
			request["token"], request["operation_id"] = token.Token, token.OperationID
			switch scenario {
			case "workspace":
				_, err := f.workspace.UpdateDatabaseSetup(t.Context(), request["expected_workspace_revision"].(string), func(context.Context, *sql.Tx, *workspace.Record) error { return nil })
				require.NoError(t, err)
			case "source":
				_, err := f.db.ExecContext(t.Context(), `UPDATE points SET address='40100' WHERE id='point-A'`)
				require.NoError(t, err)
			case "connector":
				_, err := f.db.ExecContext(t.Context(), `UPDATE database_connectors SET identity_revision='connector-2' WHERE id='connector-A'`)
				require.NoError(t, err)
			}
			response := performJSONRequest(t, f.router, http.MethodPost, writeGroupsPath+"/"+groupID+"/schema-apply", request)
			require.Equal(t, http.StatusConflict, response.Code, response.Body.String())
			require.Empty(t, columnsOf(t, f.target, "postgres", f.namespace, table), "no statement ran")
			var hypertables int
			require.NoError(t, f.target.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM timescaledb_information.hypertables WHERE hypertable_schema = $1`, f.namespace).Scan(&hypertables))
			require.Zero(t, hypertables)
		})
	}
}

// TestTimescaleScopeIsNotCallerSpecifiable: the table of a timescale group is
// derived from its storage scope, which includes the mode. A caller can echo
// the derived name but can neither pick another table nor reuse the plain
// scope's table for a timescale group.
func TestTimescaleScopeIsNotCallerSpecifiable(t *testing.T) {
	f := newManagedDestinationFixture(t, "postgres")
	_, created := ensureTimescaleGroup(t, f)
	_, group := groupSaveData(t, created)
	groupID, table := group["id"].(string), destinationOf(group)["table_name"].(string)
	current, err := f.groups.Get(t.Context(), groupID)
	require.NoError(t, err)

	plain := *current.Group
	plain.StorageMode = workspace.WriteGroupStorageModePlain
	plainTable, plainOwner, ok := workspace.LongStorageIdentity(&plain)
	require.True(t, ok)
	tsTable, tsOwner, _ := workspace.LongStorageIdentity(current.Group)
	require.Equal(t, table, tsTable)
	require.NotEqual(t, plainTable, tsTable, "the mode is part of the table identity")
	require.NotEqual(t, plainOwner, tsOwner, "and of the owner marker, so the two modes can never own each other's table")

	for name, chosen := range map[string]string{"another table": "foreign_table", "the plain scope's table": plainTable} {
		changed := *current.Group
		changed.Destination.TableName = chosen
		_, err := f.groups.Update(t.Context(), groupID, workspace.WriteGroupMutation{WorkspaceID: f.record.ID, ExpectedWorkspaceRevision: current.WorkspaceRevision, ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: "connector-1", Group: &changed})
		require.ErrorIs(t, err, workspace.ErrWriteGroupValidation, name)
	}
	echoed := *current.Group
	_, err = f.groups.Update(t.Context(), groupID, workspace.WriteGroupMutation{WorkspaceID: f.record.ID, ExpectedWorkspaceRevision: current.WorkspaceRevision, ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: "connector-1", Group: &echoed})
	require.NoError(t, err, "echoing the derived name back is accepted")
	require.Equal(t, table, destinationOf(getGroup(t, f, groupID))["table_name"])
}

// TestTimescaleStorageModeCannotChangeAfterSave: the mode is part of the table
// identity, so a saved group can neither be switched to plain nor between modes;
// the refused update leaves the group and its derived table untouched.
func TestTimescaleStorageModeCannotChangeAfterSave(t *testing.T) {
	f := newManagedDestinationFixture(t, "postgres")
	_, created := ensureTimescaleGroup(t, f)
	_, group := groupSaveData(t, created)
	groupID, table := group["id"].(string), destinationOf(group)["table_name"].(string)
	current, err := f.groups.Get(t.Context(), groupID)
	require.NoError(t, err)
	for _, mode := range []workspace.WriteGroupStorageMode{workspace.WriteGroupStorageModePlain, "native"} {
		changed := *current.Group
		changed.StorageMode, changed.Destination.TableName = mode, ""
		_, err := f.groups.Update(t.Context(), groupID, workspace.WriteGroupMutation{WorkspaceID: f.record.ID, ExpectedWorkspaceRevision: current.WorkspaceRevision, ExpectedGroupRevision: current.Group.Revision, ExpectedConnectorRevision: "connector-1", Group: &changed})
		require.ErrorIs(t, err, workspace.ErrWriteGroupValidation, string(mode))
	}
	saved := getGroup(t, f, groupID)
	require.Equal(t, "timescale", saved["storage_mode"])
	require.Equal(t, table, destinationOf(saved)["table_name"])
	require.Equal(t, group["revision"], saved["revision"])
}
