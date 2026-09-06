package dbtarget

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MySQL 不接受索引建立語句的 IF NOT EXISTS；SQLite 與 PostgreSQL 沿用既有守衛式語法。
func TestBuildEnsureUniqueIndexStatement_MySQLSchemaUsesDialectValidSyntax(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		kind          schema.DatabaseConnectorKind
		schemaName    string
		wantStatement string
		wantIndexKey  string
	}{
		{
			name:          "mysql omits if not exists",
			kind:          schema.DatabaseConnectorKindMySQL,
			schemaName:    "gateway_metrics",
			wantStatement: "CREATE UNIQUE INDEX `gateway_metrics_sensor_readings_ts_uniq` ON `gateway_metrics`.`sensor_readings` (`ts`)",
			wantIndexKey:  "gateway_metrics_sensor_readings_ts_uniq",
		},
		{
			name:          "postgres keeps if not exists",
			kind:          schema.DatabaseConnectorKindPostgres,
			schemaName:    "public",
			wantStatement: `CREATE UNIQUE INDEX IF NOT EXISTS "public_sensor_readings_ts_uniq" ON "public"."sensor_readings" ("ts")`,
			wantIndexKey:  "public_sensor_readings_ts_uniq",
		},
		{
			name:          "sqlite keeps if not exists",
			kind:          schema.DatabaseConnectorKindSQLite,
			schemaName:    "main",
			wantStatement: `CREATE UNIQUE INDEX IF NOT EXISTS "main_sensor_readings_ts_uniq" ON "sensor_readings" ("ts")`,
			wantIndexKey:  "main_sensor_readings_ts_uniq",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			statement, indexKey := buildEnsureUniqueIndexStatement(
				testCase.kind,
				testCase.schemaName,
				"sensor_readings",
				"ts",
			)

			assert.Equal(t, testCase.wantStatement, statement)
			assert.Equal(t, testCase.wantIndexKey, indexKey)
		})
	}
}

func TestBuildEnsureUniqueIndexStatement_MySQLSchemaNeverEmitsIfNotExists(t *testing.T) {
	t.Parallel()

	statement, _ := buildEnsureUniqueIndexStatement(
		schema.DatabaseConnectorKindMySQL,
		"gateway_metrics",
		"sensor_readings",
		"ts",
	)

	assert.NotContains(t, statement, "IF NOT EXISTS")
}

func mysqlSchemaConnector() *schema.DatabaseConnector {
	return &schema.DatabaseConnector{
		ID:               "conn-mysql",
		Kind:             schema.DatabaseConnectorKindMySQL,
		ConnectionConfig: `{"host":"127.0.0.1","port":"3306","user":"gateway","password":"secret","database":"gateway_metrics"}`,
	}
}

// 空白 table_schema 必須以連接器連線設定的資料庫名稱解析，不得退回 kind 層預設的 main。
func TestBuildSchemaGenerateStatements_MySQLSchemaResolvesDefaultFromConnector(t *testing.T) {
	t.Parallel()

	timestampColumn := "ts"
	mappings := []*schema.DatabaseTargetMapping{{
		ID:              "mapping-1",
		TagID:           "tag-1",
		ConnectorID:     "conn-mysql",
		TableSchema:     "",
		TableName:       "sensor_readings",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: &timestampColumn,
		Enabled:         true,
	}}

	statements, err := buildSchemaGenerateStatements(
		context.Background(),
		mysqlSchemaConnector(),
		mappings,
		nil,
		nil,
	)

	require.NoError(t, err)
	require.NotEmpty(t, statements)
	for _, statement := range statements {
		assert.NotContains(t, statement, "`main`")
	}
	assert.Contains(t, statements, "CREATE DATABASE IF NOT EXISTS `gateway_metrics`")
	createTableFound := false
	for _, statement := range statements {
		if strings.Contains(statement, "CREATE TABLE") {
			createTableFound = true
			assert.Contains(t, statement, "`gateway_metrics`.`sensor_readings`")
		}
	}
	assert.True(t, createTableFound, "缺表時必須產生建表語句")
}

// 空白 table_schema 的映射必須能對上以真實資料庫名稱檢查出來的既有資料表。
func TestBuildSchemaGenerateStatements_MySQLSchemaMatchesInspectedTables(t *testing.T) {
	t.Parallel()

	timestampColumn := "ts"
	mappings := []*schema.DatabaseTargetMapping{{
		ID:              "mapping-1",
		TagID:           "tag-1",
		ConnectorID:     "conn-mysql",
		TableSchema:     "",
		TableName:       "sensor_readings",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: &timestampColumn,
		Enabled:         true,
	}}
	tables := []TableInfo{{
		Schema: "gateway_metrics",
		Name:   "sensor_readings",
		Columns: []ColumnInfo{
			{Name: "ts", DataType: "datetime", PrimaryKey: true, Unique: true},
			{Name: "value", DataType: "text"},
		},
	}}

	statements, err := buildSchemaGenerateStatements(
		context.Background(),
		mysqlSchemaConnector(),
		mappings,
		tables,
		nil,
	)

	require.NoError(t, err)
	assert.Empty(t, statements)
}

// MySQL upsert 依賴時間戳欄位上的唯一鍵；schema 產生必須保證該鍵存在。
func TestBuildSchemaGenerateStatements_MySQLSchemaGuaranteesUpsertUniqueKey(t *testing.T) {
	t.Parallel()

	timestampColumn := "ts"
	mappings := []*schema.DatabaseTargetMapping{{
		ID:              "mapping-1",
		TagID:           "tag-1",
		ConnectorID:     "conn-mysql",
		TableSchema:     "gateway_metrics",
		TableName:       "sensor_readings",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: &timestampColumn,
		Enabled:         true,
	}}
	tables := []TableInfo{{
		Schema: "gateway_metrics",
		Name:   "sensor_readings",
		Columns: []ColumnInfo{
			{Name: "ts", DataType: "datetime"},
			{Name: "value", DataType: "text"},
		},
	}}

	statements, err := buildSchemaGenerateStatements(
		context.Background(),
		mysqlSchemaConnector(),
		mappings,
		tables,
		nil,
	)

	require.NoError(t, err)
	uniqueIndexStatement := ""
	for _, statement := range statements {
		if strings.Contains(statement, "CREATE UNIQUE INDEX") {
			uniqueIndexStatement = statement
		}
	}
	require.NotEmpty(t, uniqueIndexStatement, "MySQL upsert 映射必須產生時間戳欄位的唯一索引")
	assert.Contains(t, uniqueIndexStatement, "`gateway_metrics`.`sensor_readings`")
	assert.Contains(t, uniqueIndexStatement, "(`ts`)")
	assert.NotContains(t, uniqueIndexStatement, "IF NOT EXISTS")
}

// 時間戳欄位已是主鍵或唯一鍵時不得重複產生索引語句。
func TestBuildSchemaGenerateStatements_MySQLSchemaSkipsExistingUniqueKey(t *testing.T) {
	t.Parallel()

	timestampColumn := "ts"
	mappings := []*schema.DatabaseTargetMapping{{
		ID:              "mapping-1",
		TagID:           "tag-1",
		ConnectorID:     "conn-mysql",
		TableSchema:     "gateway_metrics",
		TableName:       "sensor_readings",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: &timestampColumn,
		Enabled:         true,
	}}
	tables := []TableInfo{{
		Schema: "gateway_metrics",
		Name:   "sensor_readings",
		Columns: []ColumnInfo{
			{Name: "ts", DataType: "datetime", Unique: true},
			{Name: "value", DataType: "text"},
		},
	}}

	statements, err := buildSchemaGenerateStatements(
		context.Background(),
		mysqlSchemaConnector(),
		mappings,
		tables,
		nil,
	)

	require.NoError(t, err)
	for _, statement := range statements {
		assert.NotContains(t, statement, "CREATE UNIQUE INDEX")
	}
}

// 唯一索引建不起來時，schema 產生必須記為失敗並保留原始資料庫錯誤訊息。
func TestConnectorService_GenerateSchema_UniqueIndexFailureIsRecordedWithDBError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mainDB := openMigratedTestDB(t)
	targetDSN := filepath.Join(t.TempDir(), "duplicate-timestamps.db")

	targetDB, err := sql.Open("sqlite", targetDSN)
	require.NoError(t, err)
	defer targetDB.Close()
	_, err = targetDB.ExecContext(ctx, `CREATE TABLE measurements (ts TEXT, line_a REAL)`)
	require.NoError(t, err)
	_, err = targetDB.ExecContext(ctx, `INSERT INTO measurements (ts, line_a) VALUES ('2026-09-06T00:00:00Z', 1), ('2026-09-06T00:00:00Z', 2)`)
	require.NoError(t, err)

	tagSvc := tag.NewService(tag.NewSQLRepository(mainDB))
	connectorRepo := NewSQLConnectorRepository(mainDB)
	mappingRepo := NewSQLTargetMappingRepository(mainDB)
	connectorSvc := NewConnectorService(connectorRepo, mappingRepo)
	connectorSvc.SetTagReader(tagSvc)

	tagEntity, err := tagSvc.Create(ctx, tag.CreateTagRequest{
		Key:         "db.schema.unique.index.failure",
		DisplayName: "DB Schema Unique Index Failure",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)
	connector, err := connectorSvc.Create(ctx, CreateConnectorRequest{
		Name:             "schema-unique-index-failure",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{"dsn": targetDSN},
	})
	require.NoError(t, err)

	timestampColumn := "ts"
	now := time.Now()
	require.NoError(t, mappingRepo.Create(ctx, &schema.DatabaseTargetMapping{
		ID:              "map-schema-unique-index-failure",
		TagID:           tagEntity.ID,
		ConnectorID:     connector.ID,
		TableSchema:     "main",
		TableName:       "measurements",
		ColumnName:      "line_a",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: &timestampColumn,
		Enabled:         true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}))

	_, err = connectorSvc.GenerateSchema(ctx, connector.ID, SchemaGenerateRequest{DryRun: false})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "執行 schema statement 失敗")
	assert.Contains(t, err.Error(), "UNIQUE")

	updated, getErr := connectorRepo.GetByID(ctx, connector.ID)
	require.NoError(t, getErr)
	assert.Equal(t, "failed", updated.LastSchemaEnsureStatus)
	assert.Contains(t, updated.LastSchemaEnsureError, "執行 schema statement 失敗")
}

// 自動建庫只能發生在明確的佈建進入點，唯讀檢查不得改變伺服器狀態。
func TestInspectTables_MySQLSchemaDoesNotAutoCreateDatabase(t *testing.T) {
	originalOpen := openExternalDBManagerFunc
	originalEnsure := ensureMySQLDatabaseIfMissingFunc
	t.Cleanup(func() {
		openExternalDBManagerFunc = originalOpen
		ensureMySQLDatabaseIfMissingFunc = originalEnsure
	})

	ensureCalls := 0
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		return nil, &mysqldriver.MySQLError{Number: 1049, Message: "Unknown database 'gateway_metrics'"}
	}
	ensureMySQLDatabaseIfMissingFunc = func(ctx context.Context, config ConnectionConfig, connectErr error) error {
		ensureCalls++
		return nil
	}

	_, err := inspectTables(context.Background(), mysqlSchemaConnector())

	require.Error(t, err)
	assert.Equal(t, 0, ensureCalls, "唯讀的資料表檢查不得自動建立資料庫")
}

func TestEnsureExternalDatabaseExists_MySQLSchemaOnlyProvisionsMySQL(t *testing.T) {
	originalOpen := openExternalDBManagerFunc
	originalEnsure := ensureMySQLDatabaseIfMissingFunc
	t.Cleanup(func() {
		openExternalDBManagerFunc = originalOpen
		ensureMySQLDatabaseIfMissingFunc = originalEnsure
	})

	openCalls := 0
	ensureCalls := 0
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		openCalls++
		if openCalls == 1 {
			return nil, &mysqldriver.MySQLError{Number: 1049, Message: "Unknown database 'gateway_metrics'"}
		}
		manager := datalinkbase.NewDBManager(datalinkbase.DBConfig{Type: datalinkbase.DBTypeSQLite, DSN: ":memory:"})
		require.NoError(t, manager.Connect())
		return manager, nil
	}
	ensureMySQLDatabaseIfMissingFunc = func(ctx context.Context, config ConnectionConfig, connectErr error) error {
		ensureCalls++
		return nil
	}

	config := ConnectionConfig{"host": "127.0.0.1", "port": "3306", "user": "gateway", "database": "gateway_metrics"}
	require.NoError(t, ensureExternalDatabaseExists(context.Background(), schema.DatabaseConnectorKindMySQL, config))
	assert.Equal(t, 1, ensureCalls)

	ensureCalls = 0
	require.NoError(t, ensureExternalDatabaseExists(context.Background(), schema.DatabaseConnectorKindSQLite, config))
	assert.Equal(t, 0, ensureCalls, "非 MySQL 類型不得走自動建庫路徑")
}

// 執行期寫入路徑不得因目標資料庫消失而靜默重建。
func TestWriter_MySQLSchemaDoesNotAutoCreateDatabaseOnWrite(t *testing.T) {
	originalOpen := openExternalDBManagerFunc
	originalEnsure := ensureMySQLDatabaseIfMissingFunc
	t.Cleanup(func() {
		openExternalDBManagerFunc = originalOpen
		ensureMySQLDatabaseIfMissingFunc = originalEnsure
	})

	ensureCalls := 0
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		return nil, &mysqldriver.MySQLError{Number: 1049, Message: "Unknown database 'gateway_metrics'"}
	}
	ensureMySQLDatabaseIfMissingFunc = func(ctx context.Context, config ConnectionConfig, connectErr error) error {
		ensureCalls++
		return nil
	}

	timestampColumn := "ts"
	writer := NewWriter(nil, nil)
	err := writer.writeMapping(
		context.Background(),
		mysqlSchemaConnector(),
		&schema.DatabaseTargetMapping{
			ID:              "mapping-1",
			TagID:           "tag-1",
			ConnectorID:     "conn-mysql",
			TableSchema:     "gateway_metrics",
			TableName:       "sensor_readings",
			ColumnName:      "value",
			WriteMode:       schema.DatabaseWriteModeUpsert,
			TimestampColumn: &timestampColumn,
			Enabled:         true,
		},
		1.5,
		time.Now(),
	)

	require.Error(t, err)
	assert.Equal(t, 0, ensureCalls, "執行期寫入不得自動建立資料庫")
}

// 這些判定函式是可替換接縫的入口，nil 錯誤必須回傳 false 而不是 panic。
func TestMissingDatabaseErrorClassifiers_MySQLSchemaHandleNilError(t *testing.T) {
	t.Parallel()

	assert.False(t, isMySQLMissingDatabaseError(nil))
	assert.False(t, isPostgresMissingDatabaseError(nil))
	assert.False(t, isPostgresDuplicateDatabaseError(nil))
	assert.NoError(t, ensureMySQLDatabaseIfMissing(context.Background(), ConnectionConfig{}, nil))
}

// MySQL 的 upsert 沒有衝突目標：值欄位帶唯一鍵時會覆寫無關資料列，必須擋下。
func TestValidateMappingAgainstTables_MySQLSchemaFlagsUniqueValueColumn(t *testing.T) {
	t.Parallel()

	timestampColumn := "ts"
	mapping := schema.DatabaseTargetMapping{
		ID: "mapping-1", TagID: "tag-1", ConnectorID: "conn-mysql",
		TableSchema: "gateway_metrics", TableName: "sensor_readings", ColumnName: "serial_no",
		WriteMode: schema.DatabaseWriteModeUpsert, TimestampColumn: &timestampColumn, Enabled: true,
	}
	tagEntity := schema.Tag{ID: "tag-1", DataType: schema.DataTypeString}
	tables := []TableInfo{{
		Schema: "gateway_metrics",
		Name:   "sensor_readings",
		Columns: []ColumnInfo{
			{Name: "ts", DataType: "datetime", PrimaryKey: true, Unique: true},
			{Name: "serial_no", DataType: "varchar", Unique: true},
		},
	}}

	mysqlIssues := validateMappingAgainstTables(schema.DatabaseConnectorKindMySQL, mapping, tagEntity, tables)
	postgresIssues := validateMappingAgainstTables(schema.DatabaseConnectorKindPostgres, mapping, tagEntity, tables)

	assert.True(t, hasIssueCode(mysqlIssues, "upsert_value_column_unique"))
	assert.False(t, hasIssueCode(postgresIssues, "upsert_value_column_unique"),
		"PostgreSQL 的 ON CONFLICT 會直接報錯，不需要這條驗證")
}

func TestValidateMappingAgainstTables_MySQLSchemaAllowsUniqueTimestampOnly(t *testing.T) {
	t.Parallel()

	timestampColumn := "ts"
	mapping := schema.DatabaseTargetMapping{
		ID: "mapping-1", TagID: "tag-1", ConnectorID: "conn-mysql",
		TableSchema: "gateway_metrics", TableName: "sensor_readings", ColumnName: "value",
		WriteMode: schema.DatabaseWriteModeUpsert, TimestampColumn: &timestampColumn, Enabled: true,
	}
	tagEntity := schema.Tag{ID: "tag-1", DataType: schema.DataTypeFloat64}
	tables := []TableInfo{{
		Schema: "gateway_metrics",
		Name:   "sensor_readings",
		Columns: []ColumnInfo{
			{Name: "ts", DataType: "datetime", PrimaryKey: true, Unique: true},
			{Name: "value", DataType: "double"},
		},
	}}

	issues := validateMappingAgainstTables(schema.DatabaseConnectorKindMySQL, mapping, tagEntity, tables)

	assert.False(t, hasIssueCode(issues, "upsert_value_column_unique"))
}

func hasIssueCode(issues []ValidationIssue, code string) bool {
	for _, issue := range issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

// 以 dsn 設定的連接器沒有離散的 user/database 鍵，自動建庫無法進行時
// 必須保留原始連線錯誤，而不是回報「設定缺漏」蓋掉真正原因。
func TestEnsureMySQLDatabaseIfMissing_MySQLSchemaKeepsOriginalErrorForDSNConfig(t *testing.T) {
	t.Parallel()

	connectErr := &mysqldriver.MySQLError{Number: 1049, Message: "Unknown database 'gateway_metrics'"}

	err := ensureMySQLDatabaseIfMissing(
		context.Background(),
		ConnectionConfig{"dsn": "gateway:secret@tcp(127.0.0.1:3306)/gateway_metrics"},
		connectErr,
	)

	require.Error(t, err)
	assert.ErrorIs(t, err, connectErr)
	assert.NotContains(t, err.Error(), "缺少 user/database")
}
