package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// timescaleDB creates a throw-away database from template0 (so it has no
// extension unless the test creates it) on the disposable PostgreSQL and returns
// a handle on it. The extension is installed here, in fixture setup only; the
// production code under test never does.
func timescaleDB(t *testing.T, setup ...string) *sql.DB {
	t.Helper()
	dsn, fields := postgresKeywordDSN(t)
	admin := postgresAdminDB(t, dsn)
	name := fmt.Sprintf("gwts_%d", time.Now().UnixNano())
	_, err := admin.ExecContext(t.Context(), `CREATE DATABASE "`+name+`" TEMPLATE template0`)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE IF EXISTS "`+name+`" WITH (FORCE)`)
	})
	db, err := sql.Open("pgx", replaceDSNField(dsn, "dbname", name))
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	_ = fields
	for _, statement := range setup {
		_, err := db.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}
	return db
}

func replaceDSNField(dsn, key, value string) string {
	parts := strings.Fields(dsn)
	for i, part := range parts {
		if strings.HasPrefix(part, key+"=") {
			parts[i] = key + "=" + value
		}
	}
	return strings.Join(parts, " ")
}

// timescaleAvailable reports whether the server offers the extension to install.
func timescaleAvailable(t *testing.T, db *sql.DB) bool {
	t.Helper()
	var n int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM pg_available_extensions WHERE name = 'timescaledb'`).Scan(&n))
	return n > 0
}

func requireTimescaleServer(t *testing.T, db *sql.DB) {
	t.Helper()
	if !timescaleAvailable(t, db) {
		t.Skip("this PostgreSQL does not offer the timescaledb extension; the capable-destination case is NOT RUN here")
	}
}

func catalogCounts(t *testing.T, db *sql.DB) (objects, extensions int) {
	t.Helper()
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT (SELECT COUNT(*) FROM pg_class), (SELECT COUNT(*) FROM pg_extension)`).Scan(&objects, &extensions))
	return objects, extensions
}

// TestTimescaleCapabilityGate pins the read-only capability probe: only an
// installed, extension-owned generalized by_range / create_hypertable API that
// the role may execute, with the schema rights preparation needs, is capable;
// everything else names its reason, and a probe never changes the database.
func TestTimescaleCapabilityGate(t *testing.T) {
	t.Run("a destination without the extension is unavailable and nothing is installed", func(t *testing.T) {
		db := timescaleDB(t)
		objects, extensions := catalogCounts(t, db)
		capability, err := ProbeTimescale(t.Context(), db, "public")
		require.NoError(t, err)
		require.False(t, capability.Available)
		require.Equal(t, []string{TimescaleReasonNotInstalled}, capability.Reasons)
		afterObjects, afterExtensions := catalogCounts(t, db)
		require.Equal(t, [2]int{objects, extensions}, [2]int{afterObjects, afterExtensions}, "the probe changed nothing")
	})

	t.Run("functions with the same names are not the extension", func(t *testing.T) {
		db := timescaleDB(t,
			`CREATE FUNCTION public.create_hypertable(relation regclass, dimension anyelement, create_default_indexes boolean DEFAULT true, if_not_exists boolean DEFAULT false, migrate_data boolean DEFAULT false) RETURNS int LANGUAGE sql AS 'SELECT 1'`,
			`CREATE FUNCTION public.by_range(column_name name, partition_interval anyelement DEFAULT NULL, partition_func regproc DEFAULT NULL) RETURNS int LANGUAGE sql AS 'SELECT 1'`,
		)
		capability, err := ProbeTimescale(t.Context(), db, "public")
		require.NoError(t, err)
		require.False(t, capability.Available, "a spoofed function is not a capability")
		require.Equal(t, []string{TimescaleReasonNotInstalled}, capability.Reasons)
	})

	t.Run("an installed extension with the generalized API is capable", func(t *testing.T) {
		admin := postgresAdminDB(t, mustDSN(t))
		requireTimescaleServer(t, admin)
		db := timescaleDB(t, `CREATE EXTENSION timescaledb`)
		objects, extensions := catalogCounts(t, db)
		capability, err := ProbeTimescale(t.Context(), db, "public")
		require.NoError(t, err)
		require.True(t, capability.Available, "reasons: %v", capability.Reasons)
		require.Empty(t, capability.Reasons)
		require.Equal(t, "public", capability.ExtensionSchema)
		require.NotEmpty(t, capability.Version)
		afterObjects, afterExtensions := catalogCounts(t, db)
		require.Equal(t, [2]int{objects, extensions}, [2]int{afterObjects, afterExtensions}, "a capable probe is still read-only")
	})

	t.Run("the extension's own schema is used, never a same-named function elsewhere", func(t *testing.T) {
		admin := postgresAdminDB(t, mustDSN(t))
		requireTimescaleServer(t, admin)
		db := timescaleDB(t,
			`CREATE SCHEMA tsext`,
			`CREATE EXTENSION timescaledb SCHEMA tsext`,
			`CREATE FUNCTION public.create_hypertable(relation regclass, dimension anyelement) RETURNS int LANGUAGE sql AS 'SELECT 1'`,
		)
		capability, err := ProbeTimescale(t.Context(), db, "public")
		require.NoError(t, err)
		require.True(t, capability.Available, "reasons: %v", capability.Reasons)
		require.Equal(t, "tsext", capability.ExtensionSchema)
	})

	t.Run("same-named functions the role cannot execute do not count against the extension", func(t *testing.T) {
		dsn := mustDSN(t)
		admin := postgresAdminDB(t, dsn)
		requireTimescaleServer(t, admin)
		role := fmt.Sprintf("tsrole_%d", time.Now().UnixNano())
		_, err := admin.ExecContext(t.Context(), `CREATE ROLE "`+role+`" LOGIN PASSWORD 'tsrole'`)
		require.NoError(t, err)
		t.Cleanup(func() { _, _ = admin.ExecContext(context.Background(), `DROP ROLE IF EXISTS "`+role+`"`) })
		// Lookalikes with the generalized signature and a locally named dimension
		// type, which the role may not execute; the real extension's functions are
		// executable. Only the extension's own functions decide the capability.
		db := timescaleDB(t,
			`CREATE SCHEMA tsext`, `CREATE EXTENSION timescaledb SCHEMA tsext`,
			`CREATE SCHEMA app`,
			`GRANT USAGE ON SCHEMA tsext TO "`+role+`"`, `GRANT USAGE, CREATE ON SCHEMA app TO "`+role+`"`,
			`CREATE TYPE public.dimension_info AS (x integer)`,
			`CREATE FUNCTION public.create_hypertable(relation regclass, dimension public.dimension_info, create_default_indexes boolean, if_not_exists boolean, migrate_data boolean) RETURNS int LANGUAGE sql AS 'SELECT 1'`,
			`REVOKE EXECUTE ON FUNCTION public.create_hypertable(regclass, public.dimension_info, boolean, boolean, boolean) FROM PUBLIC`,
		)
		restricted, err := sql.Open("pgx", replaceDSNField(replaceDSNField(replaceDSNField(dsn, "dbname", currentDatabase(t, db)), "user", role), "password", "tsrole"))
		require.NoError(t, err)
		t.Cleanup(func() { _ = restricted.Close() })
		capability, err := ProbeTimescale(t.Context(), restricted, "app")
		require.NoError(t, err)
		require.True(t, capability.Available, "reasons: %v", capability.Reasons)
		require.Equal(t, "tsext", capability.ExtensionSchema)
	})

	t.Run("a role that can create in the target schema but not use it is not capable", func(t *testing.T) {
		dsn := mustDSN(t)
		admin := postgresAdminDB(t, dsn)
		requireTimescaleServer(t, admin)
		role := fmt.Sprintf("tsrole_%d", time.Now().UnixNano())
		_, err := admin.ExecContext(t.Context(), `CREATE ROLE "`+role+`" LOGIN PASSWORD 'tsrole'`)
		require.NoError(t, err)
		t.Cleanup(func() { _, _ = admin.ExecContext(context.Background(), `DROP ROLE IF EXISTS "`+role+`"`) })
		// Every other fact holds: the extension (in its own schema) is usable and
		// executable, and the role may CREATE in app. Without USAGE it still cannot
		// resolve a table it creates there, so create_hypertable and the writer
		// could not use it.
		db := timescaleDB(t,
			`CREATE SCHEMA tsext`, `CREATE EXTENSION timescaledb SCHEMA tsext`,
			`CREATE SCHEMA app`,
			`REVOKE ALL ON SCHEMA app FROM PUBLIC`,
			`GRANT USAGE ON SCHEMA tsext TO "`+role+`"`,
			`GRANT CREATE ON SCHEMA app TO "`+role+`"`,
		)
		var usage, create bool
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT has_schema_privilege($1, 'app', 'USAGE'), has_schema_privilege($1, 'app', 'CREATE')`, role).Scan(&usage, &create))
		require.False(t, usage, "the fixture really has no USAGE on app, not even through PUBLIC")
		require.True(t, create)
		restricted, err := sql.Open("pgx", replaceDSNField(replaceDSNField(replaceDSNField(dsn, "dbname", currentDatabase(t, db)), "user", role), "password", "tsrole"))
		require.NoError(t, err)
		t.Cleanup(func() { _ = restricted.Close() })

		objects, extensions := catalogCounts(t, db)
		capability, err := ProbeTimescale(t.Context(), restricted, "app")
		require.NoError(t, err)
		require.False(t, capability.Available, "CREATE without USAGE cannot prepare a hypertable")
		require.Equal(t, []string{TimescaleReasonSchemaUsageDenied}, capability.Reasons)
		afterObjects, afterExtensions := catalogCounts(t, db)
		require.Equal(t, [2]int{objects, extensions}, [2]int{afterObjects, afterExtensions}, "the probe changed nothing")

		_, err = db.ExecContext(t.Context(), `GRANT USAGE ON SCHEMA app TO "`+role+`"`)
		require.NoError(t, err)
		capability, err = ProbeTimescale(t.Context(), restricted, "app")
		require.NoError(t, err)
		require.True(t, capability.Available, "with USAGE restored the same role is capable: %v", capability.Reasons)
	})

	t.Run("missing rights name their reason", func(t *testing.T) {
		dsn := mustDSN(t)
		admin := postgresAdminDB(t, dsn)
		requireTimescaleServer(t, admin)
		// The role is cluster-wide: create it first so its cleanup runs after the
		// database that granted it rights has been dropped.
		role := fmt.Sprintf("tsrole_%d", time.Now().UnixNano())
		_, err := admin.ExecContext(t.Context(), `CREATE ROLE "`+role+`" LOGIN PASSWORD 'tsrole'`)
		require.NoError(t, err)
		t.Cleanup(func() { _, _ = admin.ExecContext(context.Background(), `DROP ROLE IF EXISTS "`+role+`"`) })
		db := timescaleDB(t,
			`CREATE SCHEMA tsext`, `CREATE EXTENSION timescaledb SCHEMA tsext`,
			`CREATE SCHEMA app`,
			`GRANT USAGE ON SCHEMA app TO "`+role+`"`, // usage, but not CREATE
			`REVOKE ALL ON SCHEMA tsext FROM PUBLIC`,
			`REVOKE EXECUTE ON FUNCTION tsext.create_hypertable(regclass, _timescaledb_internal.dimension_info, boolean, boolean, boolean) FROM PUBLIC`,
		)
		restricted, err := sql.Open("pgx", replaceDSNField(replaceDSNField(replaceDSNField(dsn, "dbname", currentDatabase(t, db)), "user", role), "password", "tsrole"))
		require.NoError(t, err)
		t.Cleanup(func() { _ = restricted.Close() })

		capability, err := ProbeTimescale(t.Context(), restricted, "app")
		require.NoError(t, err)
		require.False(t, capability.Available)
		require.ElementsMatch(t, []string{TimescaleReasonExtensionSchemaDenied, TimescaleReasonNotExecutable, TimescaleReasonSchemaCreateDenied}, capability.Reasons)

		// The same role with its rights restored is capable: the reasons were the rights.
		for _, grant := range []string{
			`GRANT USAGE ON SCHEMA tsext TO "` + role + `"`,
			`GRANT EXECUTE ON FUNCTION tsext.create_hypertable(regclass, _timescaledb_internal.dimension_info, boolean, boolean, boolean) TO "` + role + `"`,
			`GRANT CREATE ON SCHEMA app TO "` + role + `"`,
		} {
			_, err := db.ExecContext(t.Context(), grant)
			require.NoError(t, err, grant)
		}
		capability, err = ProbeTimescale(t.Context(), restricted, "app")
		require.NoError(t, err)
		require.True(t, capability.Available, "reasons: %v", capability.Reasons)
	})
}

func mustDSN(t *testing.T) string {
	t.Helper()
	dsn, _ := postgresKeywordDSN(t)
	return dsn
}

func currentDatabase(t *testing.T, db *sql.DB) string {
	t.Helper()
	var name string
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT current_database()`).Scan(&name))
	return name
}

func TestTimescaleSignatureMatching(t *testing.T) {
	create := "relation regclass, dimension _timescaledb_internal.dimension_info, create_default_indexes boolean, if_not_exists boolean, migrate_data boolean"
	byRange := "column_name name, partition_interval anyelement, partition_func regproc"
	legacy := "relation regclass, time_column_name name, partitioning_column name, number_partitions integer, associated_schema_name name, associated_table_prefix name, chunk_time_interval anyelement, create_default_indexes boolean, if_not_exists boolean, partitioning_func regproc, migrate_data boolean, time_partitioning_func regproc"
	require.True(t, hasGeneralizedAPI([]string{legacy, create}, []string{byRange}))
	require.True(t, hasGeneralizedAPI([]string{strings.ReplaceAll(create, "_timescaledb_internal.", "other.")}, []string{byRange}), "the internal type's schema qualifier does not matter")
	require.False(t, hasGeneralizedAPI([]string{legacy}, []string{byRange}), "the legacy signature alone is not the generalized API")
	require.False(t, hasGeneralizedAPI([]string{create}, nil), "by_range is required too")
	require.False(t, hasGeneralizedAPI([]string{"relation regclass, dimension anyelement, create_default_indexes boolean, if_not_exists boolean, migrate_data boolean"}, []string{byRange}),
		"a function that merely shares the names and argument names is not the extension's dimension type")
}

// TestTimescaleHypertableCatalog reads the catalog evidence of an owned
// hypertable created by the layout's own DDL shape.
func TestTimescaleHypertableCatalog(t *testing.T) {
	admin := postgresAdminDB(t, mustDSN(t))
	requireTimescaleServer(t, admin)
	db := timescaleDB(t, `CREATE EXTENSION timescaledb`, `CREATE SCHEMA app`)
	exec := func(statement string) {
		t.Helper()
		_, err := db.ExecContext(t.Context(), statement)
		require.NoError(t, err, statement)
	}
	exec(`CREATE TABLE app.plain_table (bucket_start TIMESTAMPTZ NOT NULL, record_id TEXT NOT NULL, PRIMARY KEY (bucket_start, record_id))`)
	exec(`CREATE TABLE app.long_table (bucket_start TIMESTAMPTZ NOT NULL, record_id TEXT NOT NULL, PRIMARY KEY (bucket_start, record_id))`)
	exec(`SELECT public.create_hypertable('app.long_table'::regclass, public.by_range('bucket_start'::name, INTERVAL '7 days'), create_default_indexes => false, if_not_exists => false, migrate_data => false)`)
	exec(`CREATE TABLE app.native (bucket_start TIMESTAMPTZ NOT NULL, record_id TEXT NOT NULL) PARTITION BY RANGE (bucket_start)`)

	info, err := InspectHypertable(t.Context(), db, "app", "long_table")
	require.NoError(t, err)
	require.NotNil(t, info)
	require.Equal(t, TimescalePartitionColumn, info.TimeColumn)
	require.Equal(t, TimescaleChunkInterval, info.ChunkInterval)
	require.Equal(t, 1, info.Dimensions)
	require.Zero(t, info.Chunks)
	require.False(t, info.CompressionEnabled)
	require.Zero(t, info.Jobs, "no retention, compression or other job is bound to it")

	exec(`INSERT INTO app.long_table VALUES ('2026-01-01T00:00:00Z','a'), ('2026-01-12T00:00:00Z','b'), ('2026-01-30T00:00:00Z','c')`)
	info, err = InspectHypertable(t.Context(), db, "app", "long_table")
	require.NoError(t, err)
	require.GreaterOrEqual(t, info.Chunks, 3, "rows seven days apart land in separate chunks")

	for _, table := range []string{"plain_table", "native", "missing"} {
		info, err := InspectHypertable(t.Context(), db, "app", table)
		require.NoError(t, err)
		require.Nil(t, info, "%s is not a hypertable", table)
	}

	bare := timescaleDB(t)
	_, err = bare.ExecContext(t.Context(), `CREATE TABLE t (bucket_start TIMESTAMPTZ NOT NULL)`)
	require.NoError(t, err)
	info, err = InspectHypertable(t.Context(), bare, "public", "t")
	require.NoError(t, err)
	require.Nil(t, info, "without the extension nothing is a hypertable, and the lookup does not fail")
}
