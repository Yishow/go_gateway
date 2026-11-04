package dbtarget

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

// This test is opt-in because it connects to an operator-provided disposable
// PostgreSQL instance. The ordinary package suite never opens a real database.
func TestReadCatalogPostgresIntegration(t *testing.T) {
	dsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("POSTGRES_DSN is not set")
	}

	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, db.PingContext(ctx))

	nameSuffix := time.Now().UTC().UnixNano()
	userSchema := fmt.Sprintf("gw_catalog_test_%d", nameSuffix)
	emptySchema := fmt.Sprintf("gw_catalog_empty_%d", nameSuffix)
	quotedUserSchema := quoteIdentifier(schema.DatabaseConnectorKindPostgres, userSchema)
	quotedEmptySchema := quoteIdentifier(schema.DatabaseConnectorKindPostgres, emptySchema)
	_, err = db.ExecContext(ctx, "CREATE SCHEMA "+quotedUserSchema)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, dropErr := db.ExecContext(context.Background(), "DROP SCHEMA "+quotedUserSchema+" CASCADE")
		require.NoError(t, dropErr)
	})
	_, err = db.ExecContext(ctx, "CREATE SCHEMA "+quotedEmptySchema)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, dropErr := db.ExecContext(context.Background(), "DROP SCHEMA "+quotedEmptySchema+" CASCADE")
		require.NoError(t, dropErr)
	})
	for _, table := range []string{"gw_catalog_table", "user_catalog_table"} {
		_, err = db.ExecContext(ctx, "CREATE TABLE "+quotedUserSchema+"."+quoteIdentifier(schema.DatabaseConnectorKindPostgres, table)+" (id INTEGER PRIMARY KEY)")
		require.NoError(t, err)
	}

	catalog, err := readCatalogByKind(ctx, db, schema.DatabaseConnectorKindPostgres, ConnectionConfig{"schema": userSchema})
	require.NoError(t, err)
	require.Equal(t, userSchema, catalog.DefaultSchema)
	var userCatalog, emptyCatalog *CatalogSchema
	for index := range catalog.Schemas {
		entry := &catalog.Schemas[index]
		if entry.Name == userSchema {
			userCatalog = entry
		}
		if entry.Name == emptySchema {
			emptyCatalog = entry
		}
		require.NotEqual(t, "pg_catalog", entry.Name)
		require.NotEqual(t, "information_schema", entry.Name)
		require.False(t, strings.HasPrefix(entry.Name, "pg_"), entry.Name)
	}
	require.NotNil(t, userCatalog)
	require.ElementsMatch(t, []string{"gw_catalog_table", "user_catalog_table"}, userCatalog.Tables)
	require.NotNil(t, emptyCatalog)
	require.Empty(t, emptyCatalog.Tables)
}
