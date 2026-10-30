package datalink

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func recordProfileDB(t *testing.T) *sql.DB {
	t.Helper()
	db := openLineageMigrationDB(t, ":memory:")
	createPre047LineageSchema(t, db)
	seedLegacyLineageRows(t, db)
	require.NoError(t, ensureSQLiteWriteGroupLineages(db))
	_, err := db.Exec("PRAGMA foreign_keys=ON")
	require.NoError(t, err)
	return db
}

func TestRecordProfileKeysMigrationPreservesLegacyAndAcceptsOnlyRecordTuple(t *testing.T) {
	db := recordProfileDB(t)
	columns := lineageSnapshotColumns(t, db, "write_group_profile_keys", "write_group_lineages", "write_group_generations")
	before := lineageMigrationSnapshot(t, db, columns)
	require.NoError(t, runSQLiteConfigurationExpansion(db, migrateMixedRecordProfileKeysTx))
	require.Equal(t, before, lineageMigrationSnapshot(t, db, columns))
	_, err := db.Exec(`INSERT INTO write_group_profile_keys(workspace_id,selection_scope,profile_kind,canonical_role,group_id) VALUES ('workspace','12345678-1234-1234-1234-123456789abc','mixed','record','unkeyed')`)
	require.NoError(t, err)
	require.NoError(t, validateSQLiteLineageKeyScopes(context.Background(), db))
	for _, tuple := range []struct{ kind, role, circuit, scope string }{
		{"mixed", "circuit", "circuit", ""}, {"mixed", "unassigned", "", "device"},
		{"instantaneous", "record", "", "12345678-1234-1234-1234-123456789abc"},
		{"mixed", "record", "circuit", "12345678-1234-1234-1234-123456789abc"},
		{"mixed", "record", "", "device"}, {"mixed", "record", "", "00000000-0000-0000-0000-000000000000"}, {"mixed", "record", "", "12345678-1234-1234-1234-123456789ABC"},
	} {
		_, err = db.Exec(`UPDATE write_group_profile_keys SET profile_kind=?,canonical_role=?,circuit_id=?,selection_scope=? WHERE group_id='unkeyed'`, tuple.kind, tuple.role, tuple.circuit, tuple.scope)
		require.Error(t, err, "%+v", tuple)
	}
	require.NoError(t, runSQLiteConfigurationExpansion(db, migrateMixedRecordProfileKeysTx))
}
