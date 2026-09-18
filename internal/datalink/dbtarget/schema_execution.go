package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// SchemaStatementExecution reports how many confirmed schema statements stayed applied.
type SchemaStatementExecution struct {
	Committed  int
	RolledBack bool
}

var (
	// ErrSchemaExecutionUnsupported marks a database kind without verified managed DDL execution.
	ErrSchemaExecutionUnsupported = errors.New("schema execution is not supported for this database kind")
	// ErrSchemaExecutionPermissionDenied marks a target that refused the statements.
	ErrSchemaExecutionPermissionDenied = errors.New("schema execution was denied by the target")
)

// ExecuteSchemaStatements runs confirmed statements on the saved target in
// schemaName. SQLite and PostgreSQL roll DDL back as one batch, so it either
// commits whole or leaves nothing applied. Other kinds are refused before any
// connection is opened.
func (s *ConnectorService) ExecuteSchemaStatements(ctx context.Context, connectorID, schemaName string, statements []string) (*SchemaStatementExecution, error) {
	connector, err := s.repo.GetByID(ctx, strings.TrimSpace(connectorID))
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}
	if !schemaExecutionVerifiedKind(connector.Kind) {
		return nil, fmt.Errorf("%w: %s", ErrSchemaExecutionUnsupported, connector.Kind)
	}
	if len(statements) == 0 {
		return &SchemaStatementExecution{}, nil
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return &SchemaStatementExecution{RolledBack: true}, err
	}
	manager, err := openExternalDBManagerFunc(connector.Kind, config)
	if err != nil {
		return &SchemaStatementExecution{RolledBack: true}, err
	}
	defer manager.Close()

	return executeSchemaBatch(ctx, manager.DB(), connector.Kind, schemaName, statements)
}

// schemaExecutionVerifiedKind lists adapters whose managed DDL execution and
// post-execution verification have been tested.
func schemaExecutionVerifiedKind(kind schema.DatabaseConnectorKind) bool {
	return kind == schema.DatabaseConnectorKindSQLite || kind == schema.DatabaseConnectorKindPostgres
}

// executeSchemaBatch runs the batch in one transaction. When a rollback or
// commit fails the applied state is unknown, so it is not reported as rolled
// back and the caller must verify the actual schema.
func executeSchemaBatch(ctx context.Context, db *sql.DB, kind schema.DatabaseConnectorKind, schemaName string, statements []string) (*SchemaStatementExecution, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return &SchemaStatementExecution{RolledBack: true}, fmt.Errorf("begin schema transaction: %w", err)
	}
	if err := bindSchemaSearchPath(ctx, tx, kind, schemaName); err != nil {
		return rollbackSchemaBatch(tx, err)
	}
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return rollbackSchemaBatch(tx, fmt.Errorf("execute schema statement: %w", err))
		}
	}
	if err := tx.Commit(); err != nil {
		return &SchemaStatementExecution{}, fmt.Errorf("commit schema transaction: %w", err)
	}
	return &SchemaStatementExecution{Committed: len(statements)}, nil
}

// bindSchemaSearchPath keeps unqualified PostgreSQL DDL in the confirmed
// schema, the same one the preview inspected.
func bindSchemaSearchPath(ctx context.Context, tx *sql.Tx, kind schema.DatabaseConnectorKind, schemaName string) error {
	name := strings.TrimSpace(schemaName)
	if kind != schema.DatabaseConnectorKindPostgres || name == "" {
		return nil
	}
	if _, err := tx.ExecContext(ctx, `SET LOCAL search_path TO `+quoteIdentifier(kind, name)); err != nil {
		return fmt.Errorf("bind schema search path: %w", err)
	}
	return nil
}

func rollbackSchemaBatch(tx *sql.Tx, cause error) (*SchemaStatementExecution, error) {
	cause = classifySchemaExecutionError(cause)
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return &SchemaStatementExecution{}, errors.Join(cause, fmt.Errorf("rollback schema transaction: %w", rollbackErr))
	}
	return &SchemaStatementExecution{RolledBack: true}, cause
}

// classifySchemaExecutionError marks a refusal by the target so the caller can
// report a permission problem instead of an unexplained failure.
func classifySchemaExecutionError(err error) error {
	if isPermissionDeniedError(err) {
		return fmt.Errorf("%w: %w", ErrSchemaExecutionPermissionDenied, err)
	}
	return err
}
