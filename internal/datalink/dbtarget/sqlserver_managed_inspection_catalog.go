package dbtarget

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"slices"
	"strconv"
	"strings"
)

const (
	sqlServerManagedCharSQLType     = "CHAR"
	sqlServerManagedNvarcharSQLType = "NVARCHAR"
	sqlServerManagedVarcharSQLType  = "VARCHAR"
)

type sqlServerManagedCatalogQuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

const (
	sqlServerManagedObjectQuery = `
		SELECT o.object_id, s.name, o.name, o.type_desc,
		       CONVERT(bit, CASE WHEN t.object_id IS NULL THEN 0 ELSE 1 END),
		       ISNULL(t.is_memory_optimized, 0), ISNULL(t.temporal_type, 0),
		       ISNULL(t.is_filetable, 0),
		       CONVERT(bit, CASE WHEN et.object_id IS NULL THEN 0 ELSE 1 END),
		       ISNULL(t.is_node, 0), ISNULL(t.is_edge, 0)
		FROM sys.objects AS o
		JOIN sys.schemas AS s ON s.schema_id = o.schema_id
		LEFT JOIN sys.tables AS t ON t.object_id = o.object_id
		LEFT JOIN sys.external_tables AS et ON et.object_id = o.object_id
		WHERE s.name = @p1 AND o.name = @p2`
	sqlServerManagedColumnsQuery = `
		SELECT c.column_id, c.name, ty.name, type_schema.name, ty.is_user_defined,
		       c.max_length, c.precision, c.scale, c.is_nullable, c.is_computed,
		       c.is_identity, c.generated_always_type, c.is_hidden, c.is_sparse,
		       c.is_rowguidcol, c.collation_name, dc.definition
		FROM sys.columns AS c
		JOIN sys.types AS ty ON ty.user_type_id = c.user_type_id
		JOIN sys.schemas AS type_schema ON type_schema.schema_id = ty.schema_id
		LEFT JOIN sys.default_constraints AS dc ON dc.object_id = c.default_object_id
		WHERE c.object_id = @p1
		ORDER BY c.column_id`
	sqlServerManagedIndexesQuery = `
		SELECT i.index_id, i.name, i.type, i.is_unique, i.is_primary_key,
		       i.has_filter, i.is_disabled, i.ignore_dup_key, i.is_hypothetical,
		       ds.type, ic.key_ordinal, ic.index_column_id, ic.is_included_column,
		       ic.is_descending_key, c.name, c.is_computed
		FROM sys.indexes AS i
		LEFT JOIN sys.data_spaces AS ds ON ds.data_space_id = i.data_space_id
		LEFT JOIN sys.index_columns AS ic ON ic.object_id = i.object_id AND ic.index_id = i.index_id
		LEFT JOIN sys.columns AS c ON c.object_id = ic.object_id AND c.column_id = ic.column_id
		WHERE i.object_id = @p1 AND i.index_id >= 0
		ORDER BY i.index_id, ic.key_ordinal, ic.index_column_id`
	sqlServerManagedRepairOperationPropertyName  = "gateway.recording.repair_operation_id"
	sqlServerManagedRepairOperationPropertyQuery = `
		SELECT CONVERT(nvarchar(max), ep.value)
		FROM sys.extended_properties AS ep
		JOIN sys.objects AS o ON o.object_id = ep.major_id
		JOIN sys.schemas AS s ON s.schema_id = o.schema_id
		WHERE ep.class = 1 AND ep.minor_id = 0 AND ep.name = @p1
		  AND s.name = @p2 AND o.name = @p3`
)

type sqlServerManagedObjectFacts struct {
	objectID        int64
	schema          string
	name            string
	tableType       string
	userTable       bool
	memoryOptimized bool
	temporalType    int64
	fileTable       bool
	external        bool
	graphNode       bool
	graphEdge       bool
}

type sqlServerManagedCatalogFacts struct {
	object                sqlServerManagedObjectFacts
	columns               []ColumnInfo
	columnCollations      map[string]string
	columnDefaults        map[string]*string
	primaryKey            []string
	indexes               []IndexInfo
	plain                 bool
	indexUnusable         bool
	provenanceJSONChecked bool
	repairOperationID     string
	nativeMonthly         *SQLServerManagedPartitionEvidence
	partitionInspected    bool
}

type sqlServerManagedColumnsResult struct {
	columns            []ColumnInfo
	collations         map[string]string
	defaults           map[string]*string
	specialColumn      bool
	provenanceColumnID int64
}

type sqlServerManagedIndexesResult struct {
	indexes         []IndexInfo
	primaryKey      []string
	plain           bool
	indexUnusable   bool
	hasStorageIndex bool
	partitioned     bool
}

type sqlServerManagedIndexKey struct {
	ordinal  int64
	name     string
	computed bool
	included bool
	desc     bool
}

type sqlServerManagedIndexFacts struct {
	id            int64
	name          string
	indexType     int64
	unique        bool
	primaryKey    bool
	hasFilter     bool
	disabled      bool
	ignoreDupKey  bool
	hypothetical  bool
	dataSpaceType string
	hasDataSpace  bool
	keys          []sqlServerManagedIndexKey
	unusable      bool
}

func inspectSQLServerManagedCatalog(
	ctx context.Context,
	db sqlServerManagedCatalogQuerier,
	schemaName, tableName string,
) (*sqlServerManagedCatalogFacts, error) {
	object, err := readSQLServerManagedObject(ctx, db, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	if object.schema != schemaName || object.name != tableName {
		return nil, ErrSQLServerManagedScopeChanged
	}

	columnResult, err := readSQLServerManagedColumns(ctx, db, object.objectID)
	if err != nil {
		return nil, err
	}
	indexResult, err := readSQLServerManagedIndexes(ctx, db, object.objectID)
	if err != nil {
		return nil, err
	}
	if object.userTable && (len(columnResult.columns) == 0 || !indexResult.hasStorageIndex) {
		return nil, errors.New("user table catalog facts are incomplete")
	}
	columnNames := make(map[string]struct{}, len(columnResult.columns))
	for index := range columnResult.columns {
		columnNames[columnResult.columns[index].Name] = struct{}{}
	}
	for _, name := range indexResult.primaryKey {
		if _, exists := columnNames[name]; !exists {
			return nil, errors.New("primary key references an unlisted column")
		}
		for index := range columnResult.columns {
			if columnResult.columns[index].Name == name {
				columnResult.columns[index].PrimaryKey = true
				columnResult.columns[index].Unique = true
				break
			}
		}
	}

	provenanceJSONChecked, err := readSQLServerManagedProvenanceCheck(ctx, db, object.objectID, columnResult.provenanceColumnID)
	if err != nil {
		return nil, err
	}
	repairOperationID, err := readSQLServerManagedRepairOperationID(ctx, db, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	plain := object.tableType == sqlServerManagedUserTableType && object.userTable && !object.memoryOptimized &&
		object.temporalType == 0 && !object.fileTable && !object.external && !object.graphNode && !object.graphEdge &&
		!columnResult.specialColumn && indexResult.plain
	var native *SQLServerManagedPartitionEvidence
	var partitionInspected bool
	if indexResult.partitioned {
		native, partitionInspected, err = readSQLServerManagedPartitions(ctx, db, schemaName, tableName, object.objectID)
		if err != nil {
			return nil, err
		}
	}
	return &sqlServerManagedCatalogFacts{
		object: object, columns: columnResult.columns, columnCollations: columnResult.collations, columnDefaults: columnResult.defaults,
		primaryKey: indexResult.primaryKey, indexes: indexResult.indexes, plain: plain, indexUnusable: indexResult.indexUnusable,
		provenanceJSONChecked: provenanceJSONChecked, repairOperationID: repairOperationID,
		nativeMonthly: native, partitionInspected: partitionInspected,
	}, nil
}

func readSQLServerManagedRepairOperationID(
	ctx context.Context,
	db sqlServerManagedCatalogQuerier,
	schemaName, tableName string,
) (string, error) {
	var operationID sql.NullString
	err := db.QueryRowContext(ctx, sqlServerManagedRepairOperationPropertyQuery,
		sqlServerManagedRepairOperationPropertyName, schemaName, tableName,
	).Scan(&operationID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	if !operationID.Valid {
		return "", errors.New("SQL Server repair operation marker is null")
	}
	return operationID.String, nil
}

func readSQLServerManagedObject(ctx context.Context, db sqlServerManagedCatalogQuerier, schemaName, tableName string) (sqlServerManagedObjectFacts, error) {
	var object sqlServerManagedObjectFacts
	err := db.QueryRowContext(ctx, sqlServerManagedObjectQuery, schemaName, tableName).Scan(
		&object.objectID, &object.schema, &object.name, &object.tableType,
		&object.userTable, &object.memoryOptimized, &object.temporalType,
		&object.fileTable, &object.external, &object.graphNode, &object.graphEdge,
	)
	return object, err
}

func readSQLServerManagedColumns(ctx context.Context, db sqlServerManagedCatalogQuerier, objectID int64) (*sqlServerManagedColumnsResult, error) {
	rows, err := db.QueryContext(ctx, sqlServerManagedColumnsQuery, objectID)
	if err != nil {
		return nil, err
	}
	columns := make([]ColumnInfo, 0, 16)
	collations := make(map[string]string)
	defaults := make(map[string]*string)
	seenNames := make(map[string]struct{})
	seenIDs := make(map[int64]struct{})
	specialColumn := false
	provenanceColumnID := int64(0)
	for rows.Next() {
		var id, maxLength, precision, scale, generatedAlways int64
		var name, typeName, typeSchema string
		var userDefined, nullable, computed, identity, hidden, sparse, rowGUID bool
		var collation, defaultDefinition sql.NullString
		if err := rows.Scan(
			&id, &name, &typeName, &typeSchema, &userDefined, &maxLength, &precision, &scale,
			&nullable, &computed, &identity, &generatedAlways, &hidden, &sparse, &rowGUID,
			&collation, &defaultDefinition,
		); err != nil {
			_ = rows.Close()
			return nil, err
		}
		if id <= 0 || name == "" {
			_ = rows.Close()
			return nil, errors.New("invalid SQL Server column identity")
		}
		if _, duplicate := seenNames[name]; duplicate {
			_ = rows.Close()
			return nil, errors.New("duplicate SQL Server column name")
		}
		if _, duplicate := seenIDs[id]; duplicate {
			_ = rows.Close()
			return nil, errors.New("duplicate SQL Server column id")
		}
		seenNames[name] = struct{}{}
		seenIDs[id] = struct{}{}
		dataType, valid := sqlServerManagedDeclaredType(typeName, typeSchema, userDefined, maxLength, precision, scale)
		if !valid {
			_ = rows.Close()
			return nil, errors.New("SQL Server column type metadata is invalid")
		}
		column := ColumnInfo{
			Name: name, DataType: dataType, Nullable: nullable,
			HasDefault: defaultDefinition.Valid || identity || computed || generatedAlways != 0,
		}
		columns = append(columns, column)
		if collation.Valid && collation.String != "" {
			collations[name] = collation.String
		}
		if defaultDefinition.Valid {
			canonical := sqlServerCanonicalDefault(defaultDefinition.String)
			defaults[name] = &canonical
		} else {
			defaults[name] = nil
		}
		if name == "provenance" {
			provenanceColumnID = id
		}
		if userDefined || computed || identity || generatedAlways != 0 || hidden || sparse || rowGUID ||
			strings.EqualFold(typeName, "timestamp") || strings.EqualFold(typeName, "rowversion") {
			specialColumn = true
		}
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return &sqlServerManagedColumnsResult{
		columns: columns, collations: collations, defaults: defaults,
		specialColumn: specialColumn, provenanceColumnID: provenanceColumnID,
	}, nil
}

func readSQLServerManagedIndexes(ctx context.Context, db sqlServerManagedCatalogQuerier, objectID int64) (*sqlServerManagedIndexesResult, error) {
	rows, err := db.QueryContext(ctx, sqlServerManagedIndexesQuery, objectID)
	if err != nil {
		return nil, err
	}
	indexesByID := make(map[int64]*sqlServerManagedIndexFacts)
	ordered := make([]*sqlServerManagedIndexFacts, 0, 8)
	for rows.Next() {
		var id, indexType int64
		var name, dataSpaceType, columnName sql.NullString
		var unique, primaryKey, hasFilter, disabled, ignoreDupKey, hypothetical sql.NullBool
		var keyOrdinal, indexColumnID sql.NullInt64
		var included, descending, computed sql.NullBool
		if err := rows.Scan(
			&id, &name, &indexType, &unique, &primaryKey, &hasFilter, &disabled,
			&ignoreDupKey, &hypothetical, &dataSpaceType, &keyOrdinal, &indexColumnID,
			&included, &descending, &columnName, &computed,
		); err != nil {
			_ = rows.Close()
			return nil, err
		}
		if id < 0 || !unique.Valid || !primaryKey.Valid || !hasFilter.Valid || !disabled.Valid ||
			!ignoreDupKey.Valid || !hypothetical.Valid {
			_ = rows.Close()
			return nil, errors.New("SQL Server index metadata is incomplete")
		}
		fact, exists := indexesByID[id]
		if !exists {
			fact = &sqlServerManagedIndexFacts{
				id: id, name: name.String, indexType: indexType, unique: unique.Bool,
				primaryKey: primaryKey.Bool, hasFilter: hasFilter.Bool, disabled: disabled.Bool,
				ignoreDupKey: ignoreDupKey.Bool, hypothetical: hypothetical.Bool,
				dataSpaceType: dataSpaceType.String, hasDataSpace: dataSpaceType.Valid,
			}
			indexesByID[id] = fact
			ordered = append(ordered, fact)
		} else if !sqlServerManagedIndexMetadataMatches(
			fact, name, indexType, unique, primaryKey, hasFilter, disabled, ignoreDupKey, hypothetical, dataSpaceType,
		) {
			_ = rows.Close()
			return nil, errors.New("inconsistent SQL Server index metadata")
		}
		if included.Valid && included.Bool {
			fact.unusable = true
		}
		if keyOrdinal.Valid && keyOrdinal.Int64 > 0 {
			if !columnName.Valid || !computed.Valid {
				fact.unusable = true
			} else {
				fact.keys = append(fact.keys, sqlServerManagedIndexKey{
					ordinal: keyOrdinal.Int64, name: columnName.String, computed: computed.Bool,
					included: included.Valid && included.Bool, desc: descending.Valid && descending.Bool,
				})
			}
			if !descending.Valid || descending.Bool || !computed.Valid || computed.Bool {
				fact.unusable = true
			}
		} else if id > 0 && indexColumnID.Valid && (!included.Valid || !included.Bool) {
			fact.unusable = true
		}
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	indexes := make([]IndexInfo, 0, len(ordered))
	primaryKeyColumns := []string{}
	primaryKeyCount := 0
	plain := len(ordered) > 0
	indexUnusable := false
	hasStorageIndex := false
	partitioned := false
	for _, fact := range ordered {
		partitioned = partitioned || fact.dataSpaceType == "PS"
		slices.SortFunc(fact.keys, func(left, right sqlServerManagedIndexKey) int {
			return cmp.Compare(left.ordinal, right.ordinal)
		})
		if fact.dataSpaceType != "FG" || fact.indexType < 0 || fact.indexType > 2 {
			plain = false
		}
		if fact.id == 0 {
			hasStorageIndex = true
			continue
		}
		keyColumns := make([]string, len(fact.keys))
		for index, key := range fact.keys {
			keyColumns[index] = key.name
		}
		if fact.primaryKey {
			primaryKeyCount++
			if len(keyColumns) == 0 || !fact.unique || fact.hasFilter || fact.disabled || fact.ignoreDupKey ||
				fact.hypothetical || fact.unusable || fact.indexType < 1 || fact.indexType > 2 {
				plain = false
				indexUnusable = true
			}
			primaryKeyColumns = keyColumns
			hasStorageIndex = true
			continue
		}
		if fact.name == "" || len(keyColumns) == 0 || (fact.indexType != 1 && fact.indexType != 2) ||
			fact.unique || fact.hasFilter || fact.disabled || fact.ignoreDupKey || fact.hypothetical || fact.unusable {
			indexUnusable = true
		}
		indexes = append(indexes, IndexInfo{Name: fact.name, Columns: keyColumns, Unique: fact.unique})
		hasStorageIndex = true
	}
	if primaryKeyCount > 1 {
		return nil, errors.New("multiple SQL Server primary key indexes")
	}
	return &sqlServerManagedIndexesResult{
		indexes: indexes, primaryKey: primaryKeyColumns, plain: plain,
		indexUnusable: indexUnusable, hasStorageIndex: hasStorageIndex,
		partitioned: partitioned,
	}, nil
}

func sqlServerManagedIndexMetadataMatches(
	fact *sqlServerManagedIndexFacts,
	name sql.NullString,
	indexType int64,
	unique, primaryKey, hasFilter, disabled, ignoreDupKey, hypothetical sql.NullBool,
	dataSpaceType sql.NullString,
) bool {
	return fact.name == name.String && fact.indexType == indexType && fact.unique == unique.Bool &&
		fact.primaryKey == primaryKey.Bool && fact.hasFilter == hasFilter.Bool && fact.disabled == disabled.Bool &&
		fact.ignoreDupKey == ignoreDupKey.Bool && fact.hypothetical == hypothetical.Bool &&
		fact.dataSpaceType == dataSpaceType.String && fact.hasDataSpace == dataSpaceType.Valid
}

func sqlServerManagedDeclaredType(typeName, typeSchema string, userDefined bool, maxLength, precision, scale int64) (string, bool) {
	name := strings.ToUpper(typeName)
	if userDefined {
		if typeSchema == "" || name == "" {
			return "", false
		}
		return strings.ToUpper(typeSchema + "." + typeName), true
	}
	switch name {
	case "NCHAR", sqlServerManagedNvarcharSQLType:
		if maxLength == -1 && name == sqlServerManagedNvarcharSQLType {
			return sqlServerNvarcharMaxSQLType, true
		}
		if maxLength < 0 || maxLength%2 != 0 {
			return "", false
		}
		return name + "(" + strconv.FormatInt(maxLength/2, 10) + ")", true
	case sqlServerManagedCharSQLType, sqlServerManagedVarcharSQLType, "BINARY", "VARBINARY":
		if maxLength == -1 && (name == sqlServerManagedVarcharSQLType || name == "VARBINARY") {
			return name + "(MAX)", true
		}
		if maxLength < 0 {
			return "", false
		}
		return name + "(" + strconv.FormatInt(maxLength, 10) + ")", true
	case sqlServerDecimalSQLType, affinityNumeric:
		if precision < 1 || scale < 0 || scale > precision {
			return "", false
		}
		return name + "(" + strconv.FormatInt(precision, 10) + "," + strconv.FormatInt(scale, 10) + ")", true
	case "DATETIME2", "DATETIMEOFFSET", "TIME":
		if scale < 0 || scale > 7 {
			return "", false
		}
		return name + "(" + strconv.FormatInt(scale, 10) + ")", true
	case "FLOAT":
		if precision != 24 && precision != 53 {
			return "", false
		}
		return "FLOAT(" + strconv.FormatInt(precision, 10) + ")", true
	default:
		if name == "" {
			return "", false
		}
		return name, true
	}
}
