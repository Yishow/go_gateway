package recordingplan

import "fmt"

// longColumn is one column of the fixed managed_long_v1 table. The layout
// contract is the same for every circuit, source and point count: sources add
// rows, never columns.
type longColumn struct {
	name       string
	sqlite     string
	postgres   string
	mysql      string
	nullable   bool
	primaryKey bool
}

// longColumns follows the documented fixed schema. PostgreSQL gets native
// TIMESTAMPTZ, NUMERIC and JSONB; SQLite keeps UTC microsecond text and exact
// decimal text so a uint64 or decimal never passes through a float. MySQL binds
// UTC DATETIME(6), bounded exact DECIMAL(65,30), bounded case-exact IDs and
// JSON; the family collation and engine come from the layout options.
var longColumns = []longColumn{
	{name: "bucket_start", sqlite: managedTextSQLType, postgres: sqlTypeTimestamptz, mysql: mysqlTimeSQLType, primaryKey: true},
	{name: "bucket_end", sqlite: managedTextSQLType, postgres: sqlTypeTimestamptz, mysql: mysqlTimeSQLType},
	{name: "record_id", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlHashSQLType, primaryKey: true},
	{name: "workspace_id", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType},
	{name: "group_id", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType},
	{name: "group_revision", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType},
	{name: "circuit_id", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType, nullable: true},
	{name: "device_id", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType},
	{name: "point_id", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType},
	{name: "tag_id", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType},
	{name: "measurement_id", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType, nullable: true},
	{name: "series_epoch", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType, nullable: true},
	{name: "profile_kind", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType},
	{name: "value_type", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType},
	{name: "unit", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType, nullable: true},
	{name: "status", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType},
	{name: "quality", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlIDSQLType, nullable: true},
	{name: "observed_at", sqlite: managedTextSQLType, postgres: sqlTypeTimestamptz, mysql: mysqlTimeSQLType, nullable: true},
	{name: "received_at", sqlite: managedTextSQLType, postgres: sqlTypeTimestamptz, mysql: mysqlTimeSQLType, nullable: true},
	{name: "value_int", sqlite: sqlTypeInteger, postgres: sqlTypeBigint, mysql: "BIGINT", nullable: true},
	{name: "value_decimal", sqlite: managedTextSQLType, postgres: "NUMERIC", mysql: mysqlDecimalSQLType, nullable: true},
	{name: "value_float", sqlite: "REAL", postgres: sqlTypeDoublePrecision, mysql: mysqlDoubleSQLType, nullable: true},
	{name: "value_bool", sqlite: sqlTypeInteger, postgres: sqlTypeBoolean, mysql: sqlTypeBoolean, nullable: true},
	{name: "value_text", sqlite: managedTextSQLType, postgres: managedTextSQLType, mysql: mysqlLongTextSQLType, nullable: true},
	{name: "delta_value", sqlite: managedTextSQLType, postgres: "NUMERIC", mysql: mysqlDecimalSQLType, nullable: true},
	{name: "delta_from", sqlite: managedTextSQLType, postgres: sqlTypeTimestamptz, mysql: mysqlTimeSQLType, nullable: true},
	{name: "delta_to", sqlite: managedTextSQLType, postgres: sqlTypeTimestamptz, mysql: mysqlTimeSQLType, nullable: true},
	{name: "provenance", sqlite: managedTextSQLType, postgres: "JSONB", mysql: mysqlJSONSQLType},
}

// ManagedLongGroupLayout returns the fixed typed long-table layout for one
// owned table. ownerColumn is the scope-derived ownership marker; it is the
// only column that depends on who owns the table. The composite key
// (bucket_start, record_id) and the source/time index are part of the layout,
// so they are generated, confirmed by digest and checked against the catalog.
func ManagedLongGroupLayout(dialect, tableName, ownerColumn string) (GroupSchemaLayout, error) {
	dialect = normalizePreviewDialect(dialect)
	if dialect != dialectSQLite && dialect != dialectPostgres {
		return GroupSchemaLayout{}, fmt.Errorf("%w: %s", ErrUnsupportedSchemaDialect, dialect)
	}
	layout := GroupSchemaLayout{
		TableName: tableName, OwnerColumn: ownerColumn,
		Columns: make([]GroupSchemaColumn, 0, len(longColumns)+1),
		Indexes: []GroupSchemaIndex{{
			Name: "idx_" + tableName + "_src_time", Columns: []string{"workspace_id", "tag_id", "bucket_start"},
		}},
	}
	for _, column := range longColumns {
		sqlType := column.sqlite
		if dialect == dialectPostgres {
			sqlType = column.postgres
		}
		layout.Columns = append(layout.Columns, GroupSchemaColumn{
			Name: column.name, SQLType: sqlType, Nullable: column.nullable, PrimaryKey: column.primaryKey,
		})
	}
	layout.Columns = append(layout.Columns, GroupSchemaColumn{Name: ownerColumn, SQLType: managedTextSQLType})
	return ValidateGroupSchemaLayout(layout)
}

// ManagedLongTimescaleLayout is the plain long layout on PostgreSQL plus the
// frozen time partitioning that makes the new table an owned hypertable. The
// typed columns, the composite key (bucket_start, record_id), the source/time
// index and the owner marker are exactly the plain contract; only the table and
// owner names differ, because the storage scope includes the mode.
func ManagedLongTimescaleLayout(tableName, ownerColumn string) (GroupSchemaLayout, error) {
	layout, err := ManagedLongGroupLayout(dialectPostgres, tableName, ownerColumn)
	if err != nil {
		return GroupSchemaLayout{}, err
	}
	layout.Partition = &GroupSchemaPartition{Kind: PartitionKindTimescale, Column: TimescalePartitionColumn, ChunkInterval: TimescaleChunkInterval}
	return ValidateGroupSchemaLayout(layout)
}
