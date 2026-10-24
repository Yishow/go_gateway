package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"modernc.org/sqlite"

	"go-gateway/internal/datalink/schema"
)

// IsSchemaShapeError reports whether a failed insert says the destination's
// table is not the shape the insert expects: the relation is gone, a column is
// missing or a column's type is not the expected one. Only such an error can
// start a look at the destination catalog for managed storage. A refused
// credential, a missing privilege, a violated constraint, a lock or a lost
// connection never can. PostgreSQL is told apart by SQLSTATE; SQLite reports
// all of these as its generic error code, so only there the statement error
// text is read, and only to decide whether to look, never to decide a repair.
func IsSchemaShapeError(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgSchemaShapeState(pgErr.Code)
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pgSchemaShapeState(string(pqErr.Code))
	}
	var liteErr *sqlite.Error
	if errors.As(err, &liteErr) && liteErr.Code()&0xff == sqliteError {
		message := strings.ToLower(liteErr.Error())
		return strings.Contains(message, "no such table") || strings.Contains(message, "no such column") ||
			strings.Contains(message, "has no column named")
	}
	return false
}

// IsMissingRelationError reports whether a failed statement says the relation
// it names does not exist. It is the narrowest of the schema-shape errors and
// says nothing about which relation, or why; the catalog decides that.
func IsMissingRelationError(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgUndefinedTable
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == pgUndefinedTable
	}
	var liteErr *sqlite.Error
	return errors.As(err, &liteErr) && liteErr.Code()&0xff == sqliteError && strings.Contains(strings.ToLower(liteErr.Error()), "no such table")
}

// pgUndefinedTable is the SQLSTATE of a relation that does not exist.
const pgUndefinedTable = "42P01"

func pgSchemaShapeState(code string) bool {
	switch code {
	case pgUndefinedTable, "42703", "42804":
		// undefined table, undefined column, datatype mismatch
		return true
	}
	return false
}

// InspectOpenManagedTable inspects one table of a destination that delivery
// already opened. It never opens a connection of its own, so it cannot create a
// missing SQLite file, and it only reads. A name held by something that is not
// a table (a view, an index, a sequence) is reported as an existing table
// without the managed columns, which no layout accepts.
func InspectOpenManagedTable(ctx context.Context, db *sql.DB, kind schema.DatabaseConnectorKind, schemaName, tableName string, withIndexes bool) (*TableInspection, error) {
	result := &TableInspection{Schema: strings.TrimSpace(schemaName), Table: strings.TrimSpace(tableName), Columns: []ColumnInfo{}}
	switch kind {
	case schema.DatabaseConnectorKindSQLite:
		result.Schema = managedSQLiteSchemaName
		var objectType string
		err := db.QueryRowContext(ctx, `SELECT type FROM sqlite_master WHERE lower(name) = lower(?) LIMIT 1`, result.Table).Scan(&objectType)
		if errors.Is(err, sql.ErrNoRows) {
			return result.withStatus(TableInspectionMissing, ""), nil
		}
		if err != nil {
			return result.classified(err), nil
		}
		if objectType != "table" {
			result.Status, result.Reason = TableInspectionExists, ""
			return result, nil
		}
		tables, err := inspectSQLiteTables(ctx, db)
		if err != nil {
			return result.classified(err), nil
		}
		result = result.fromListing(tables, TableInspectionFailed, inspectionReasonExistenceUnconfirmed)
		if withIndexes && result.Status == TableInspectionExists {
			indexes, err := inspectSQLiteIndexes(ctx, db, result.Table)
			if err != nil {
				return nil, err
			}
			result.Indexes, result.IndexesInspected = indexes, true
		}
		return result, nil
	case schema.DatabaseConnectorKindPostgres:
		result.Schema = defaultString(result.Schema, "public")
		var relkind string
		err := db.QueryRowContext(ctx, `
			SELECT c.relkind::text FROM pg_catalog.pg_class c JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
			WHERE n.nspname = $1 AND c.relname = $2`, result.Schema, result.Table).Scan(&relkind)
		if errors.Is(err, sql.ErrNoRows) {
			return result.withStatus(TableInspectionMissing, ""), nil
		}
		if err != nil {
			return result.classified(err), nil
		}
		if relkind != "r" && relkind != "p" {
			result.Status, result.Reason = TableInspectionExists, ""
			return result, nil
		}
		tables, err := inspectPostgresTables(ctx, db)
		if err != nil {
			return result.classified(err), nil
		}
		result = result.fromListing(tables, TableInspectionForbidden, inspectionReasonPermissionDenied)
		if withIndexes && result.Status == TableInspectionExists {
			indexes, err := inspectPostgresIndexes(ctx, db, result.Schema, result.Table)
			if err != nil {
				return nil, err
			}
			result.Indexes, result.IndexesInspected = indexes, true
		}
		return result, nil
	}
	return result.withStatus(TableInspectionFailed, inspectionReasonUnsupportedKind), nil
}

// CanCreateManagedTable reports whether the connected role may create a table
// in the schema: for PostgreSQL both USAGE and CREATE on it. SQLite has no
// per-role privileges, so a refusal there shows up when the statement runs.
func CanCreateManagedTable(ctx context.Context, db *sql.DB, kind schema.DatabaseConnectorKind, schemaName string) (bool, error) {
	if kind != schema.DatabaseConnectorKindPostgres {
		return true, nil
	}
	var allowed bool
	err := db.QueryRowContext(ctx, `
		SELECT to_regnamespace($1::text) IS NOT NULL
		   AND has_schema_privilege(current_user, $1::text, 'USAGE') AND has_schema_privilege(current_user, $1::text, 'CREATE')`,
		defaultString(strings.TrimSpace(schemaName), "public")).Scan(&allowed)
	if err != nil {
		return false, fmt.Errorf("read schema privileges: %w", err)
	}
	return allowed, nil
}

// RunManagedSchemaBatch runs confirmed managed statements on an already opened
// destination as one transaction: all of them commit, or none stays applied. An
// error from the commit itself leaves the outcome unknown (RolledBack false).
func RunManagedSchemaBatch(ctx context.Context, db *sql.DB, kind schema.DatabaseConnectorKind, schemaName string, statements []string) (*SchemaStatementExecution, error) {
	if !schemaExecutionVerifiedKind(kind) {
		return &SchemaStatementExecution{RolledBack: true}, fmt.Errorf("%w: %s", ErrSchemaExecutionUnsupported, kind)
	}
	if len(statements) == 0 {
		return &SchemaStatementExecution{}, nil
	}
	return executeSchemaBatch(ctx, db, kind, schemaName, statements)
}
