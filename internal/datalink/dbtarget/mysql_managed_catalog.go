package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	mysqldriver "github.com/go-sql-driver/mysql"
)

type sqlMySQLManagedSession struct {
	manager   *datalinkbase.DBManager
	connector *schema.DatabaseConnector
	family    string
}

type mysqlManagedQuerier interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func (s *sqlMySQLManagedSession) ServerFacts(ctx context.Context, database string) (mysqlManagedServerFacts, error) {
	if s == nil || s.manager == nil || s.manager.DB() == nil {
		return mysqlManagedServerFacts{}, fmt.Errorf("mysql managed session is unavailable")
	}
	metadata := inspectMySQLRecordingMetadata(ctx, s.manager.DB(), database)
	capability := projectMySQLRecordingCapability(s.connector, metadata, true)
	facts := mysqlManagedServerFacts{
		family:              capability.ServerFamily,
		version:             capability.ServerVersion,
		capabilityRevision:  capability.CapabilityRevision,
		database:            defaultString(metadata.actualDatabase, database),
		metadataStatus:      metadata.metadataStatus,
		transactionalEngine: metadata.transactionalEngine,
	}
	s.family = facts.family
	return facts, nil
}

func (s *sqlMySQLManagedSession) TableFacts(ctx context.Context, database, table string) (mysqlManagedTableFacts, error) {
	return readMySQLManagedTableFacts(ctx, s.manager.DB(), s.family, database, table)
}

func mergeMySQLManagedPartitions(facts mysqlManagedTableFacts, partitions *MySQLPartitionEvidence, err error) (mysqlManagedTableFacts, error) {
	if err != nil {
		if isPermissionDeniedError(err) {
			facts.status, facts.reason = TableInspectionForbidden, inspectionReasonPermissionDenied
			return facts, nil
		}
		return facts, err
	}
	facts.partitions = partitions
	return facts, nil
}

func readMySQLManagedTableFacts(ctx context.Context, q mysqlManagedQuerier, family, database, table string) (mysqlManagedTableFacts, error) {
	row, err := readTableRow(ctx, q, database, table)
	if err != nil {
		if isPermissionDeniedError(err) {
			return mysqlManagedTableFacts{status: TableInspectionForbidden, reason: inspectionReasonPermissionDenied}, nil
		}
		return mysqlManagedTableFacts{}, err
	}
	if row == nil {
		return absentTableFacts(ctx, q, database, table), nil
	}
	facts := *row
	facts.status = TableInspectionExists
	facts.charset = charsetOfCollation(facts.collation)

	columns, err := readColumns(ctx, q, family, database, table)
	if err != nil {
		if isPermissionDeniedError(err) {
			facts.status, facts.reason = TableInspectionForbidden, inspectionReasonPermissionDenied
			return facts, nil
		}
		return facts, err
	}
	facts.columns = columns
	if !strings.EqualFold(facts.tableType, "BASE TABLE") {
		return facts, nil
	}

	partitions, partitionErr := readMySQLManagedPartitions(ctx, q, database, table, facts.engine)
	facts, err = mergeMySQLManagedPartitions(facts, partitions, partitionErr)
	if err != nil {
		return facts, err
	}

	primaryKey, err := readPrimaryKey(ctx, q, database, table)
	if err != nil {
		if isPermissionDeniedError(err) {
			facts.status, facts.reason = TableInspectionForbidden, inspectionReasonPermissionDenied
			return facts, nil
		}
		return facts, err
	}
	facts.primaryKey = primaryKey

	indexes, err := readIndexes(ctx, q, family, database, table)
	if err != nil {
		if isPermissionDeniedError(err) {
			facts.status, facts.reason = TableInspectionForbidden, inspectionReasonPermissionDenied
			return facts, nil
		}
		return facts, err
	}
	facts.indexes = indexes

	if err := markJSONColumns(ctx, q, family, database, table, &facts); err != nil {
		if isPermissionDeniedError(err) {
			facts.status, facts.reason = TableInspectionForbidden, inspectionReasonPermissionDenied
			return facts, nil
		}
		return facts, err
	}
	return facts, nil
}

func (s *sqlMySQLManagedSession) Exec(ctx context.Context, statement string) error {
	if s == nil || s.manager == nil || s.manager.DB() == nil {
		return fmt.Errorf("mysql managed session is unavailable")
	}
	_, err := s.manager.DB().ExecContext(ctx, statement)
	return err
}

func (s *sqlMySQLManagedSession) Close() error {
	if s == nil || s.manager == nil {
		return nil
	}
	return s.manager.Close()
}

func readTableRow(ctx context.Context, q mysqlManagedQuerier, database, table string) (*mysqlManagedTableFacts, error) {
	var tableType, engine, collation sql.NullString
	err := q.QueryRowContext(ctx, `
		SELECT TABLE_TYPE, ENGINE, TABLE_COLLATION
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`, database, table,
	).Scan(&tableType, &engine, &collation)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &mysqlManagedTableFacts{
		tableType: tableType.String, engine: engine.String, collation: collation.String,
	}, nil
}

func absentTableFacts(ctx context.Context, q mysqlManagedQuerier, database, table string) mysqlManagedTableFacts {
	rows, err := q.QueryContext(ctx, "SHOW CREATE TABLE "+qualifiedTableName(schema.DatabaseConnectorKindMySQL, database, table))
	if err == nil {
		_ = rows.Close()
		// A concurrent creator appeared between the catalog read and SHOW CREATE.
		// Presence is not proven by this probe, so require a fresh inspection.
		return mysqlManagedTableFacts{status: TableInspectionFailed, reason: inspectionReasonExistenceUnconfirmed}
	}
	return mysqlManagedAbsenceFacts(err)
}

func mysqlManagedAbsenceFacts(err error) mysqlManagedTableFacts {
	if number, ok := mysqlErrorNumber(err); ok {
		switch number {
		case 1146:
			return mysqlManagedTableFacts{status: TableInspectionMissing}
		case 1044, 1045, 1142, 1143, 1145, 1227:
			return mysqlManagedTableFacts{status: TableInspectionForbidden, reason: inspectionReasonPermissionDenied}
		}
	}
	if isPermissionDeniedError(err) {
		return mysqlManagedTableFacts{status: TableInspectionForbidden, reason: inspectionReasonPermissionDenied}
	}
	return mysqlManagedTableFacts{status: TableInspectionFailed, reason: inspectionReasonExistenceUnconfirmed}
}

func readColumns(ctx context.Context, q mysqlManagedQuerier, family, database, table string) ([]mysqlManagedColumnFacts, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT, EXTRA, COLLATION_NAME
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION`, database, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []mysqlManagedColumnFacts
	for rows.Next() {
		var name, dataType, columnType, nullable, extra, collation sql.NullString
		var rawDefault sql.NullString
		if err := rows.Scan(&name, &dataType, &columnType, &nullable, &rawDefault, &extra, &collation); err != nil {
			return nil, err
		}
		var defaultInput *string
		if rawDefault.Valid {
			value := rawDefault.String
			defaultInput = &value
		}
		decodedDefault, err := decodeMySQLManagedDefault(family, defaultInput)
		if err != nil {
			return nil, err
		}
		column := mysqlManagedColumnFacts{
			name: name.String, dataType: dataType.String, columnType: columnType.String,
			collation: collation.String, nullable: strings.EqualFold(nullable.String, "YES"),
			hasDefault:   decodedDefault != nil || hasGeneratedDefault(extra.String),
			defaultValue: decodedDefault,
		}
		columns = append(columns, column)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return columns, nil
}

func hasGeneratedDefault(extra string) bool {
	extra = strings.ToLower(extra)
	return strings.Contains(extra, "auto_increment") ||
		strings.Contains(extra, "default_generated") ||
		strings.Contains(extra, "generated")
}

func readPrimaryKey(ctx context.Context, q mysqlManagedQuerier, database, table string) ([]string, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT COLUMN_NAME, SUB_PART, INDEX_TYPE, NON_UNIQUE
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = 'PRIMARY'
		ORDER BY SEQ_IN_INDEX`, database, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var primaryKey []mysqlManagedPrimaryKeyRow
	for rows.Next() {
		var row mysqlManagedPrimaryKeyRow
		if err := rows.Scan(&row.column, &row.subPart, &row.indexType, &row.nonUnique); err != nil {
			return nil, err
		}
		primaryKey = append(primaryKey, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return mysqlManagedPrimaryKeyColumns(primaryKey)
}

func readIndexes(ctx context.Context, q mysqlManagedQuerier, family, database, table string) ([]mysqlManagedIndexFacts, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT INDEX_NAME, SEQ_IN_INDEX, COLUMN_NAME, NON_UNIQUE, INDEX_TYPE, SUB_PART
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME <> 'PRIMARY'
		ORDER BY INDEX_NAME, SEQ_IN_INDEX`, database, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byName := map[string]*mysqlManagedIndexFacts{}
	var order []string
	for rows.Next() {
		var name, indexType sql.NullString
		var column sql.NullString
		var seq, nonUnique int
		var subPart sql.NullInt64
		if err := rows.Scan(&name, &seq, &column, &nonUnique, &indexType, &subPart); err != nil {
			return nil, err
		}
		entry, ok := byName[name.String]
		if !ok {
			entry = &mysqlManagedIndexFacts{name: name.String, unique: nonUnique == 0, indexType: indexType.String, visible: true}
			byName[name.String] = entry
			order = append(order, name.String)
		}
		entry.functional = entry.functional || !column.Valid
		entry.prefix = entry.prefix || subPart.Valid
		if column.Valid {
			entry.columns = append(entry.columns, column.String)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(byName) == 0 {
		return []mysqlManagedIndexFacts{}, nil
	}
	if err := applyIndexVisibility(ctx, q, family, database, table, byName); err != nil {
		return nil, err
	}
	indexes := make([]mysqlManagedIndexFacts, 0, len(order))
	for _, name := range order {
		entry := byName[name]
		entry.columns = slices.Clone(entry.columns)
		indexes = append(indexes, *entry)
	}
	return indexes, nil
}

func applyIndexVisibility(ctx context.Context, q mysqlManagedQuerier, family, database, table string, indexes map[string]*mysqlManagedIndexFacts) error {
	rows, err := q.QueryContext(ctx, "SHOW INDEX FROM "+qualifiedTableName(schema.DatabaseConnectorKindMySQL, database, table))
	if err != nil {
		return err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	keyIndex, visibleIndex, ignoredIndex := -1, -1, -1
	for i, name := range columns {
		switch strings.ToLower(name) {
		case "key_name":
			keyIndex = i
		case "visible":
			visibleIndex = i
		case mysqlManagedIndexIgnored:
			ignoredIndex = i
		}
	}
	if keyIndex < 0 {
		return applyMySQLManagedIndexVisibility(family, indexes, columns, nil)
	}
	values := make([]sql.NullString, len(columns))
	destinations := make([]any, len(columns))
	for i := range values {
		destinations[i] = &values[i]
	}
	var visibility []mysqlManagedVisibilityRow
	for rows.Next() {
		if err := rows.Scan(destinations...); err != nil {
			return err
		}
		row := mysqlManagedVisibilityRow{name: values[keyIndex].String}
		if visibleIndex >= 0 {
			row.visible = values[visibleIndex]
		}
		if ignoredIndex >= 0 {
			row.ignored = values[ignoredIndex]
		}
		visibility = append(visibility, row)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return applyMySQLManagedIndexVisibility(family, indexes, columns, visibility)
}

func markJSONColumns(ctx context.Context, q mysqlManagedQuerier, family, database, table string, facts *mysqlManagedTableFacts) error {
	if family != mysqlManagedFamilyMariaDB {
		facts.jsonColumns = map[string]bool{}
		for i := range facts.columns {
			if strings.EqualFold(facts.columns[i].dataType, "json") {
				facts.columns[i].jsonValid = true
				facts.jsonColumns[facts.columns[i].name] = true
			}
		}
		return nil
	}
	facts.jsonColumns = map[string]bool{}
	rows, err := q.QueryContext(ctx, `
		SELECT CHECK_CLAUSE
		FROM information_schema.CHECK_CONSTRAINTS cc
		JOIN information_schema.TABLE_CONSTRAINTS tc
		  ON tc.CONSTRAINT_SCHEMA = cc.CONSTRAINT_SCHEMA
		 AND tc.CONSTRAINT_NAME = cc.CONSTRAINT_NAME
		 AND tc.TABLE_NAME = cc.TABLE_NAME
		WHERE tc.TABLE_SCHEMA = ? AND tc.TABLE_NAME = ? AND tc.CONSTRAINT_TYPE = 'CHECK'`, database, table)
	if err != nil {
		return err
	}
	defer rows.Close()
	var clauses []string
	for rows.Next() {
		var clause string
		if err := rows.Scan(&clause); err != nil {
			return err
		}
		clauses = append(clauses, strings.ToLower(clause))
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for i := range facts.columns {
		column := &facts.columns[i]
		if !strings.EqualFold(column.dataType, "longtext") {
			continue
		}
		for _, clause := range clauses {
			if mysqlManagedCheckIsJSONValid(clause, column.name) {
				column.jsonValid = true
				facts.jsonColumns[column.name] = true
				break
			}
		}
	}
	return nil
}

func charsetOfCollation(collation string) string {
	charset, _, found := strings.Cut(collation, "_")
	if !found {
		return ""
	}
	return charset
}

func mysqlErrorNumber(err error) (uint16, bool) {
	var mysqlErr *mysqldriver.MySQLError
	if !errors.As(err, &mysqlErr) {
		return 0, false
	}
	return mysqlErr.Number, true
}
