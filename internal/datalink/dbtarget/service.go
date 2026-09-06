package dbtarget

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

type tagGetter interface {
	GetByID(ctx context.Context, id string) (*schema.Tag, error)
}

type ConnectorService struct {
	repo        ConnectorRepository
	mappingRepo TargetMappingRepository
}

func NewConnectorService(repo ConnectorRepository, mappingRepoOpt ...TargetMappingRepository) *ConnectorService {
	service := &ConnectorService{repo: repo}
	if len(mappingRepoOpt) > 0 {
		service.mappingRepo = mappingRepoOpt[0]
	}
	return service
}

func (s *ConnectorService) Create(ctx context.Context, req CreateConnectorRequest) (*schema.DatabaseConnector, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, validationError("資料庫連接器名稱不可為空")
	}
	connectionConfigJSON, err := serializeConnectionConfig(req.ConnectionConfig)
	if err != nil {
		return nil, err
	}
	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立資料庫連接器 ID 失敗: %w", err)
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	defaultIntervalSeconds := defaultWriteIntervalSeconds(req.DefaultWriteIntervalSeconds)
	status, lastCheckAt, lastCheckError := probeConnector(ctx, req.Kind, req.ConnectionConfig)
	now := time.Now()
	connector := &schema.DatabaseConnector{
		ID:                          id,
		Name:                        name,
		Kind:                        req.Kind,
		ConnectionConfig:            connectionConfigJSON,
		Status:                      status,
		LastCheckAt:                 lastCheckAt,
		LastCheckError:              lastCheckError,
		Enabled:                     enabled,
		DefaultWriteIntervalSeconds: defaultIntervalSeconds,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}
	if err := s.repo.Create(ctx, connector); err != nil {
		return nil, fmt.Errorf("建立資料庫連接器失敗: %w", err)
	}
	return connector, nil
}

func (s *ConnectorService) Update(ctx context.Context, id string, req UpdateConnectorRequest) (*schema.DatabaseConnector, error) {
	connector, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, validationError("資料庫連接器名稱不可為空")
		}
		connector.Name = name
	}
	if req.Kind != nil {
		connector.Kind = *req.Kind
	}
	if req.Enabled != nil {
		connector.Enabled = *req.Enabled
	}
	if req.DefaultWriteIntervalSeconds != nil {
		connector.DefaultWriteIntervalSeconds = defaultWriteIntervalSeconds(req.DefaultWriteIntervalSeconds)
	}
	existingConnectionConfig, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, err
	}
	clearPassword := req.ClearPassword != nil && *req.ClearPassword
	var connectionConfig ConnectionConfig
	if req.ConnectionConfig != nil {
		connectionConfig = preserveSensitiveConnectionConfigValues(
			existingConnectionConfig,
			*req.ConnectionConfig,
			clearPassword,
		)
		connector.ConnectionConfig, err = serializeConnectionConfig(connectionConfig)
		if err != nil {
			return nil, err
		}
	} else {
		connectionConfig = cloneConnectionConfig(existingConnectionConfig)
		if clearPassword {
			delete(connectionConfig, "password")
			connector.ConnectionConfig, err = serializeConnectionConfig(connectionConfig)
			if err != nil {
				return nil, err
			}
		}
	}
	connector.Status, connector.LastCheckAt, connector.LastCheckError = probeConnector(ctx, connector.Kind, connectionConfig)
	connector.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, connector); err != nil {
		return nil, fmt.Errorf("更新資料庫連接器失敗: %w", err)
	}
	return connector, nil
}

func (s *ConnectorService) Delete(ctx context.Context, id string) error {
	if s.mappingRepo != nil {
		if err := s.mappingRepo.DeleteByConnectorID(ctx, id); err != nil {
			return fmt.Errorf("刪除資料庫連接器映射失敗: %w", err)
		}
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除資料庫連接器失敗: %w", err)
	}
	return nil
}

func (s *ConnectorService) GetByID(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	connector, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}
	return connector, nil
}

func (s *ConnectorService) List(ctx context.Context, filter ConnectorListFilter) ([]*schema.DatabaseConnector, error) {
	connectors, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出資料庫連接器失敗: %w", err)
	}
	return connectors, nil
}

func (s *ConnectorService) TestConnection(ctx context.Context, id string) (*schema.DatabaseConnector, error) {
	connector, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}
	connectionConfig, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, err
	}
	connector.Status, connector.LastCheckAt, connector.LastCheckError = probeConnector(ctx, connector.Kind, connectionConfig)
	connector.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, connector); err != nil {
		return nil, fmt.Errorf("更新資料庫連接器測試狀態失敗: %w", err)
	}
	return connector, nil
}

func (s *ConnectorService) ListTables(ctx context.Context, id string) ([]TableInfo, error) {
	connector, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}

	return inspectTables(ctx, connector)
}

type MappingService struct {
	repo          TargetMappingRepository
	connectorRepo ConnectorRepository
	tagService    tagGetter
}

func NewMappingService(repo TargetMappingRepository, connectorRepo ConnectorRepository, tagService tagGetter) *MappingService {
	return &MappingService{
		repo:          repo,
		connectorRepo: connectorRepo,
		tagService:    tagService,
	}
}

func (s *MappingService) Create(ctx context.Context, req CreateTargetMappingRequest) (*schema.DatabaseTargetMapping, error) {
	connector, err := s.connectorRepo.GetByID(ctx, req.ConnectorID)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}

	tagEntity, err := s.tagService.GetByID(ctx, req.TagID)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}
	id, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("建立資料庫目標映射 ID 失敗: %w", err)
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	writeMode := req.WriteMode
	if writeMode == "" {
		writeMode = schema.DatabaseWriteModeInsert
	}
	tableSchema := normalizeOptionalString(req.TableSchema)
	if tableSchema == "" {
		tableSchema = defaultSchemaForConnector(connector)
	}
	mapping := &schema.DatabaseTargetMapping{
		ID:                   id,
		TagID:                strings.TrimSpace(req.TagID),
		ConnectorID:          connector.ID,
		TableSchema:          tableSchema,
		TableName:            strings.TrimSpace(req.TableName),
		ColumnName:           strings.TrimSpace(req.ColumnName),
		WriteMode:            writeMode,
		TimestampColumn:      normalizeTimestampColumnForWriteMode(writeMode, req.TimestampColumn),
		GroupKey:             normalizeOptionalPointer(req.GroupKey),
		WriteIntervalSeconds: normalizeOptionalIntPointer(req.WriteIntervalSeconds),
		Enabled:              enabled,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	if err := validateMappingDefinition(ctx, connector, tagEntity, mapping, req.AllowMissingTable); err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, mapping); err != nil {
		return nil, fmt.Errorf("建立資料庫目標映射失敗: %w", err)
	}
	return mapping, nil
}

func (s *MappingService) Update(ctx context.Context, id string, req UpdateTargetMappingRequest) (*schema.DatabaseTargetMapping, error) {
	mapping, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫目標映射失敗: %w", err)
	}

	connector, err := s.connectorRepo.GetByID(ctx, mapping.ConnectorID)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}

	tagEntity, err := s.tagService.GetByID(ctx, mapping.TagID)
	if err != nil {
		return nil, fmt.Errorf("取得標籤失敗: %w", err)
	}
	if req.TableSchema != nil {
		tableSchema := normalizeOptionalString(*req.TableSchema)
		if tableSchema == "" {
			tableSchema = defaultSchemaForConnector(connector)
		}
		mapping.TableSchema = tableSchema
	}
	if req.TableName != nil {
		mapping.TableName = strings.TrimSpace(*req.TableName)
	}
	if req.ColumnName != nil {
		mapping.ColumnName = strings.TrimSpace(*req.ColumnName)
	}
	if req.WriteMode != nil {
		mapping.WriteMode = *req.WriteMode
	}
	if mapping.WriteMode != schema.DatabaseWriteModeUpsert {
		mapping.TimestampColumn = nil
	} else if req.TimestampColumn != nil {
		mapping.TimestampColumn = normalizeOptionalPointer(req.TimestampColumn)
	}
	if req.Enabled != nil {
		mapping.Enabled = *req.Enabled
	}
	if req.GroupKey != nil {
		mapping.GroupKey = normalizeOptionalPointer(req.GroupKey)
	}
	if req.WriteIntervalSeconds != nil {
		mapping.WriteIntervalSeconds = normalizeOptionalIntPointer(req.WriteIntervalSeconds)
	}
	mapping.UpdatedAt = time.Now()
	if err := validateMappingDefinition(ctx, connector, tagEntity, mapping, req.AllowMissingTable); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, mapping); err != nil {
		return nil, fmt.Errorf("更新資料庫目標映射失敗: %w", err)
	}
	return mapping, nil
}

func (s *MappingService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("刪除資料庫目標映射失敗: %w", err)
	}
	return nil
}

func (s *MappingService) GetByID(ctx context.Context, id string) (*schema.DatabaseTargetMapping, error) {
	mapping, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫目標映射失敗: %w", err)
	}
	return mapping, nil
}

func (s *MappingService) List(ctx context.Context, filter TargetMappingListFilter) ([]*schema.DatabaseTargetMapping, error) {
	mappings, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("列出資料庫目標映射失敗: %w", err)
	}
	return mappings, nil
}

func serializeConnectionConfig(config ConnectionConfig) (string, error) {
	if config == nil {
		config = ConnectionConfig{}
	}
	data, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("序列化資料庫連接設定失敗: %w", err)
	}
	return string(data), nil
}

func parseConnectionConfig(raw string) (ConnectionConfig, error) {
	var config ConnectionConfig
	if err := json.Unmarshal([]byte(raw), &config); err != nil {
		return nil, fmt.Errorf("解析資料庫連接設定失敗: %w", err)
	}
	return config, nil
}

func inspectTables(ctx context.Context, connector *schema.DatabaseConnector) ([]TableInfo, error) {
	return inspectTablesForSchemas(ctx, connector, nil)
}

// inspectTablesForSchemas 除了連接器自身的資料庫，額外納入映射引用到的 schema，
// 讓跨資料庫的映射也能對上既有資料表。
func inspectTablesForSchemas(
	ctx context.Context,
	connector *schema.DatabaseConnector,
	extraSchemas []string,
) ([]TableInfo, error) {
	connectionConfig, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, err
	}

	// 檢查資料表是唯讀動作：不得順手建立目標資料庫。自動建庫僅限探測與
	// 非 dry-run 的 schema 產生這兩個明確的佈建進入點。
	manager, err := openExternalDBManagerFunc(connector.Kind, connectionConfig)
	if err != nil {
		return nil, err
	}
	defer manager.Close()

	return inspectTablesByKind(ctx, manager.DB(), connector.Kind, connectionConfig, extraSchemas)
}

// mappingSchemaNames 取出一組映射實際引用的 schema 名稱。
func mappingSchemaNames(mappings []*schema.DatabaseTargetMapping) []string {
	names := make([]string, 0, len(mappings))
	for _, mappingRecord := range mappings {
		if mappingRecord == nil {
			continue
		}
		if name := strings.TrimSpace(mappingRecord.TableSchema); name != "" {
			names = append(names, name)
		}
	}
	return names
}

func openExternalDBManager(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
	dbConfig, err := buildExternalDBConfig(kind, config)
	if err != nil {
		return nil, err
	}

	manager := datalinkbase.NewDBManager(dbConfig)
	if err := manager.Connect(); err != nil {
		return nil, err
	}
	return manager, nil
}

func buildExternalDBConfig(kind schema.DatabaseConnectorKind, config ConnectionConfig) (datalinkbase.DBConfig, error) {
	switch kind {
	case schema.DatabaseConnectorKindSQLite:
		dsn := stringConfigValue(config, "dsn")
		if dsn == "" {
			dsn = stringConfigValue(config, "path")
		}
		if dsn == "" {
			return datalinkbase.DBConfig{}, fmt.Errorf("sqlite 連接設定缺少 dsn/path")
		}
		return datalinkbase.DBConfig{
			Type:         datalinkbase.DBTypeSQLite,
			DSN:          dsn,
			MaxOpenConns: 1,
			MaxIdleConns: 1,
		}, nil
	case schema.DatabaseConnectorKindPostgres:
		dsn := strings.TrimSpace(stringConfigValue(config, "dsn"))
		if dsn == "" {
			host := defaultString(stringConfigValue(config, "host"), "127.0.0.1")
			port := defaultString(stringConfigValue(config, "port"), "5432")
			user := stringConfigValue(config, "user")
			password := stringConfigValue(config, "password")
			databaseName := defaultString(stringConfigValue(config, "database"), stringConfigValue(config, "dbname"))
			sslMode := defaultString(stringConfigValue(config, "sslmode"), "disable")
			if user == "" || databaseName == "" {
				return datalinkbase.DBConfig{}, fmt.Errorf("postgres 連接設定缺少 user/database")
			}
			dsn = fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				host,
				port,
				user,
				password,
				databaseName,
				sslMode,
			)
		}
		return datalinkbase.DBConfig{
			Type:         datalinkbase.DBTypePostgres,
			DSN:          dsn,
			MaxOpenConns: 3,
			MaxIdleConns: 1,
		}, nil
	case schema.DatabaseConnectorKindMySQL:
		dsn := strings.TrimSpace(stringConfigValue(config, "dsn"))
		if dsn == "" {
			user := defaultString(
				stringConfigValue(config, "user"),
				stringConfigValue(config, "username"),
			)
			databaseName := defaultString(
				stringConfigValue(config, "database"),
				stringConfigValue(config, "dbname"),
			)
			if user == "" || databaseName == "" {
				return datalinkbase.DBConfig{}, fmt.Errorf("mysql 連接設定缺少 user/database")
			}

			mysqlConfig := newMySQLDriverConfig(config, databaseName)

			if timeout := strings.TrimSpace(stringConfigValue(config, "timeout")); timeout != "" {
				duration, err := time.ParseDuration(timeout)
				if err != nil {
					return datalinkbase.DBConfig{}, fmt.Errorf("mysql timeout 格式錯誤: %w", err)
				}
				mysqlConfig.Timeout = duration
				mysqlConfig.ReadTimeout = duration
				mysqlConfig.WriteTimeout = duration
			}

			dsn = mysqlConfig.FormatDSN()
		}
		return datalinkbase.DBConfig{
			Type:         datalinkbase.DBTypeMySQL,
			DSN:          dsn,
			MaxOpenConns: 3,
			MaxIdleConns: 1,
		}, nil
	default:
		return datalinkbase.DBConfig{}, fmt.Errorf("不支援的資料庫類型: %s", kind)
	}
}

func classifyConnectorError(err error) schema.DatabaseConnectorStatus {
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "authentication") ||
		strings.Contains(message, "password") ||
		strings.Contains(message, "access denied"):
		return schema.DatabaseConnectorStatusAuthFailed
	case strings.Contains(message, "不支援的資料庫類型"), strings.Contains(message, "unsupported"):
		return schema.DatabaseConnectorStatusError
	default:
		return schema.DatabaseConnectorStatusUnreachable
	}
}

func connectorErrorMessage(kind schema.DatabaseConnectorKind, err error) string {
	if kind == schema.DatabaseConnectorKindMySQL && errors.Is(err, mysqldriver.ErrUnknownPlugin) {
		return "資料庫連接測試失敗: 伺服器要求未支援的 MySQL/MariaDB 驗證插件；若為 MariaDB PAM / dialog，請在伺服器啟用 pam_use_cleartext_plugin，或改用 mysql_native_password / mysql_clear_password 帳號"
	}
	return err.Error()
}

func inspectSQLiteTables(ctx context.Context, db *sql.DB) ([]TableInfo, error) {
	rows, err := db.QueryContext(
		ctx,
		`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' ORDER BY name`,
	)
	if err != nil {
		return nil, fmt.Errorf("查詢 sqlite 資料表失敗: %w", err)
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("掃描 sqlite 資料表失敗: %w", err)
		}
		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍歷 sqlite 資料表失敗: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("關閉 sqlite 資料表查詢失敗: %w", err)
	}

	tables := make([]TableInfo, 0, len(tableNames))
	for _, tableName := range tableNames {
		uniqueColumns, err := inspectSQLiteUniqueColumns(ctx, db, tableName)
		if err != nil {
			return nil, err
		}

		pragmaQuery := fmt.Sprintf(`PRAGMA table_info('%s')`, escapeSQLiteStringLiteral(tableName))
		columnRows, err := db.QueryContext(ctx, pragmaQuery)
		if err != nil {
			return nil, fmt.Errorf("查詢 sqlite 欄位失敗: %w", err)
		}

		table := TableInfo{
			Schema: "main",
			Name:   tableName,
		}
		for columnRows.Next() {
			var cid int
			var name string
			var dataType string
			var notNull int
			var defaultValue sql.NullString
			var primaryKey int
			if err := columnRows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey); err != nil {
				columnRows.Close()
				return nil, fmt.Errorf("掃描 sqlite 欄位失敗: %w", err)
			}
			table.Columns = append(table.Columns, ColumnInfo{
				Name:       name,
				DataType:   dataType,
				Nullable:   notNull == 0,
				PrimaryKey: primaryKey > 0,
				Unique:     primaryKey > 0 || uniqueColumns[strings.ToLower(name)],
			})
		}
		if err := columnRows.Close(); err != nil {
			return nil, fmt.Errorf("關閉 sqlite 欄位查詢失敗: %w", err)
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func inspectSQLiteUniqueColumns(
	ctx context.Context,
	db *sql.DB,
	tableName string,
) (map[string]bool, error) {
	uniqueColumns := map[string]bool{}
	indexListQuery := fmt.Sprintf(`PRAGMA index_list('%s')`, escapeSQLiteStringLiteral(tableName))
	indexRows, err := db.QueryContext(ctx, indexListQuery)
	if err != nil {
		return nil, fmt.Errorf("查詢 sqlite 索引失敗: %w", err)
	}

	indexNames := make([]string, 0, 4)

	for indexRows.Next() {
		var seq int
		var name string
		var unique int
		var origin string
		var partial int
		if err := indexRows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			indexRows.Close()
			return nil, fmt.Errorf("掃描 sqlite 索引失敗: %w", err)
		}
		if unique == 0 {
			continue
		}
		indexNames = append(indexNames, name)
	}

	if err := indexRows.Err(); err != nil {
		indexRows.Close()
		return nil, fmt.Errorf("遍歷 sqlite 索引失敗: %w", err)
	}
	if err := indexRows.Close(); err != nil {
		return nil, fmt.Errorf("關閉 sqlite 索引查詢失敗: %w", err)
	}

	for _, indexName := range indexNames {
		indexInfoQuery := fmt.Sprintf(`PRAGMA index_info('%s')`, escapeSQLiteStringLiteral(indexName))
		infoRows, err := db.QueryContext(ctx, indexInfoQuery)
		if err != nil {
			return nil, fmt.Errorf("查詢 sqlite 索引欄位失敗: %w", err)
		}

		columnNames := make([]string, 0, 1)
		for infoRows.Next() {
			var seqno int
			var cid int
			var columnName string
			if err := infoRows.Scan(&seqno, &cid, &columnName); err != nil {
				infoRows.Close()
				return nil, fmt.Errorf("掃描 sqlite 索引欄位失敗: %w", err)
			}
			columnNames = append(columnNames, columnName)
		}
		if err := infoRows.Close(); err != nil {
			return nil, fmt.Errorf("關閉 sqlite 索引欄位查詢失敗: %w", err)
		}
		if len(columnNames) == 1 {
			uniqueColumns[strings.ToLower(columnNames[0])] = true
		}
	}

	return uniqueColumns, nil
}

func inspectPostgresTables(ctx context.Context, db *sql.DB) ([]TableInfo, error) {
	rows, err := db.QueryContext(
		ctx,
		`
			SELECT table_schema, table_name
			FROM information_schema.tables
			WHERE table_type = 'BASE TABLE'
			  AND table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY table_schema, table_name
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("查詢 postgres 資料表失敗: %w", err)
	}
	defer rows.Close()

	var tables []TableInfo
	for rows.Next() {
		var schemaName string
		var tableName string
		if err := rows.Scan(&schemaName, &tableName); err != nil {
			return nil, fmt.Errorf("掃描 postgres 資料表失敗: %w", err)
		}

		columnRows, err := db.QueryContext(
			ctx,
			`
				SELECT
					c.column_name,
					c.data_type,
					(c.is_nullable = 'YES') AS is_nullable,
					EXISTS (
						SELECT 1
						FROM information_schema.table_constraints tc
						JOIN information_schema.key_column_usage kcu
						  ON tc.constraint_name = kcu.constraint_name
						 AND tc.table_schema = kcu.table_schema
						 AND tc.table_name = kcu.table_name
						WHERE tc.constraint_type = 'PRIMARY KEY'
						  AND tc.table_schema = c.table_schema
						  AND tc.table_name = c.table_name
						  AND kcu.column_name = c.column_name
					) AS is_primary_key,
					EXISTS (
						SELECT 1
						FROM information_schema.table_constraints tc
						JOIN information_schema.key_column_usage kcu
						  ON tc.constraint_name = kcu.constraint_name
						 AND tc.table_schema = kcu.table_schema
						 AND tc.table_name = kcu.table_name
						WHERE tc.constraint_type IN ('PRIMARY KEY', 'UNIQUE')
						  AND tc.table_schema = c.table_schema
						  AND tc.table_name = c.table_name
						  AND kcu.column_name = c.column_name
						  AND (
							SELECT COUNT(*)
							FROM information_schema.key_column_usage kcu_count
							WHERE kcu_count.constraint_name = tc.constraint_name
							  AND kcu_count.table_schema = tc.table_schema
							  AND kcu_count.table_name = tc.table_name
						  ) = 1
					) AS is_unique
				FROM information_schema.columns c
				WHERE c.table_schema = $1
				  AND c.table_name = $2
				ORDER BY c.ordinal_position
			`,
			schemaName,
			tableName,
		)
		if err != nil {
			return nil, fmt.Errorf("查詢 postgres 欄位失敗: %w", err)
		}

		table := TableInfo{
			Schema: schemaName,
			Name:   tableName,
		}
		for columnRows.Next() {
			var column ColumnInfo
			if err := columnRows.Scan(
				&column.Name,
				&column.DataType,
				&column.Nullable,
				&column.PrimaryKey,
				&column.Unique,
			); err != nil {
				columnRows.Close()
				return nil, fmt.Errorf("掃描 postgres 欄位失敗: %w", err)
			}
			table.Columns = append(table.Columns, column)
		}
		if err := columnRows.Close(); err != nil {
			return nil, fmt.Errorf("關閉 postgres 欄位查詢失敗: %w", err)
		}

		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍歷 postgres 資料表失敗: %w", err)
	}

	return tables, nil
}

func validateMappingDefinition(ctx context.Context, connector *schema.DatabaseConnector, tagEntity *schema.Tag, mapping *schema.DatabaseTargetMapping, allowMissingTable bool) error {
	if mapping.TagID == "" {
		return validationError("標籤 ID 不可為空")
	}
	if mapping.ConnectorID == "" {
		return validationError("資料庫連接器 ID 不可為空")
	}
	if mapping.TableName == "" {
		return validationError("資料表名稱不可為空")
	}
	if mapping.ColumnName == "" {
		return validationError("資料欄位名稱不可為空")
	}
	if mapping.WriteMode == "" {
		mapping.WriteMode = schema.DatabaseWriteModeInsert
	}
	if mapping.WriteMode == schema.DatabaseWriteModeUpsert && normalizeOptionalPointer(mapping.TimestampColumn) == nil {
		return validationError("upsert 模式必須設定 timestamp_column")
	}

	tables, err := inspectTablesForSchemas(ctx, connector, mappingSchemaNames([]*schema.DatabaseTargetMapping{mapping}))
	if err != nil {
		return fmt.Errorf("檢查資料庫目標表結構失敗: %w", err)
	}

	issues := validateMappingAgainstTables(connector.Kind, *mapping, *tagEntity, tables)
	for _, issue := range issues {
		if issue.Severity != "error" {
			continue
		}
		// allowMissingTable 時，表/欄位尚未建立並非阻擋條件，將由後續 schema 生成補建。
		if allowMissingTable && isMissingSchemaObjectIssue(issue.Code) {
			continue
		}
		return validationError(issue.Message)
	}

	return nil
}

// isMissingSchemaObjectIssue 判斷該驗證問題是否為「目標物件尚未建立」類型，
// 這類問題可由 schema 生成補建，不應在 allowMissingTable 模式下阻擋儲存。
func isMissingSchemaObjectIssue(code string) bool {
	switch code {
	case "table_missing", "column_missing", "timestamp_column_missing":
		return true
	default:
		return false
	}
}

func validateMappingAgainstTables(
	kind schema.DatabaseConnectorKind,
	mapping schema.DatabaseTargetMapping,
	tagEntity schema.Tag,
	tables []TableInfo,
) []ValidationIssue {
	var issues []ValidationIssue

	table, ok := findTable(tables, mapping.TableSchema, mapping.TableName)
	if !ok {
		return append(issues, ValidationIssue{
			Severity:  "error",
			MappingID: mapping.ID,
			TagID:     mapping.TagID,
			Code:      "table_missing",
			Message:   fmt.Sprintf("找不到資料表: %s.%s", mapping.TableSchema, mapping.TableName),
		})
	}

	column, ok := findColumn(table.Columns, mapping.ColumnName)
	if !ok {
		issues = append(issues, ValidationIssue{
			Severity:  "error",
			MappingID: mapping.ID,
			TagID:     mapping.TagID,
			Code:      "column_missing",
			Message:   fmt.Sprintf("找不到資料欄位: %s", mapping.ColumnName),
		})
	} else if !isTypeCompatible(tagEntity.DataType, column) {
		issues = append(issues, ValidationIssue{
			Severity:  "error",
			MappingID: mapping.ID,
			TagID:     mapping.TagID,
			Code:      "column_type_mismatch",
			Message:   fmt.Sprintf("欄位 %s 型別 %s 與標籤資料型別 %s 不相容", mapping.ColumnName, column.DataType, tagEntity.DataType),
		})
	}

	if mapping.WriteMode == schema.DatabaseWriteModeUpsert {
		if mapping.TimestampColumn == nil {
			issues = append(issues, ValidationIssue{
				Severity:  "error",
				MappingID: mapping.ID,
				TagID:     mapping.TagID,
				Code:      "timestamp_missing",
				Message:   "upsert 模式必須設定 timestamp_column",
			})
		} else if timestampColumn, ok := findColumn(table.Columns, *mapping.TimestampColumn); !ok {
			issues = append(issues, ValidationIssue{
				Severity:  "error",
				MappingID: mapping.ID,
				TagID:     mapping.TagID,
				Code:      "timestamp_column_missing",
				Message:   fmt.Sprintf("找不到 timestamp 欄位: %s", *mapping.TimestampColumn),
			})
		} else if !timestampColumn.PrimaryKey && !timestampColumn.Unique {
			issues = append(issues, ValidationIssue{
				Severity:  "error",
				MappingID: mapping.ID,
				TagID:     mapping.TagID,
				Code:      "timestamp_column_not_unique",
				Message:   fmt.Sprintf("timestamp 欄位 %s 必須是 primary key 或 single-column unique constraint，upsert 才能正確運作", *mapping.TimestampColumn),
			})
		}
		// MySQL 的 ON DUPLICATE KEY UPDATE 沒有衝突目標，任何唯一鍵撞上都會觸發
		// 覆寫；PostgreSQL 的 ON CONFLICT (ts) 則會直接報錯。值欄位本身帶唯一鍵時，
		// 寫入撞到的會是一列與此時間戳無關的資料，屬於靜默覆寫。
		if kind == schema.DatabaseConnectorKindMySQL {
			if valueColumn, ok := findColumn(table.Columns, mapping.ColumnName); ok &&
				(valueColumn.PrimaryKey || valueColumn.Unique) &&
				!strings.EqualFold(strings.TrimSpace(mapping.ColumnName), strings.TrimSpace(optionalStringValue(mapping.TimestampColumn))) {
				issues = append(issues, ValidationIssue{
					Severity:  "error",
					MappingID: mapping.ID,
					TagID:     mapping.TagID,
					Code:      "upsert_value_column_unique",
					Message:   fmt.Sprintf("值欄位 %s 帶有唯一鍵，MySQL 的 upsert 會在該鍵衝突時覆寫無關資料列", mapping.ColumnName),
				})
			}
		}
	}

	return issues
}

func findTable(tables []TableInfo, schemaName string, tableName string) (TableInfo, bool) {
	targetSchema := strings.TrimSpace(schemaName)
	targetName := strings.TrimSpace(tableName)
	for _, table := range tables {
		if strings.EqualFold(table.Schema, targetSchema) && strings.EqualFold(table.Name, targetName) {
			return table, true
		}
	}
	return TableInfo{}, false
}

func findColumn(columns []ColumnInfo, columnName string) (ColumnInfo, bool) {
	target := strings.TrimSpace(columnName)
	for _, column := range columns {
		if strings.EqualFold(column.Name, target) {
			return column, true
		}
	}
	return ColumnInfo{}, false
}

func isTypeCompatible(dataType schema.DataType, column ColumnInfo) bool {
	columnType := strings.ToLower(strings.TrimSpace(column.DataType))
	if columnType == "" {
		return true
	}

	switch dataType {
	case schema.DataTypeBool:
		return containsAny(columnType, "bool", "bit", "int", "integer", "smallint", "tinyint", "numeric", "decimal")
	case schema.DataTypeInt16, schema.DataTypeUint16, schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeInt64, schema.DataTypeUint64:
		return containsAny(columnType, "int", "integer", "bigint", "smallint", "tinyint", "numeric", "decimal", "real", "double", "float")
	case schema.DataTypeFloat32, schema.DataTypeFloat64:
		return containsAny(columnType, "real", "double", "float", "numeric", "decimal", "int", "integer")
	case schema.DataTypeString:
		return containsAny(columnType, "char", "text", "clob", "varchar", "json", "uuid")
	default:
		return true
	}
}

func containsAny(value string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(value, keyword) {
			return true
		}
	}
	return false
}

func defaultSchemaForKind(kind schema.DatabaseConnectorKind) string {
	if kind == schema.DatabaseConnectorKindPostgres {
		return "public"
	}
	return "main"
}

func stringConfigValue(config ConnectionConfig, key string) string {
	value, ok := config[key]
	if !ok || value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprintf("%v", value))
}

func booleanConfigValue(config ConnectionConfig, defaultValue bool, keys ...string) bool {
	for _, key := range keys {
		value, ok := config[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case bool:
			return typed
		case string:
			switch strings.ToLower(strings.TrimSpace(typed)) {
			case "1", "true", "yes", "on":
				return true
			case "0", "false", "no", "off":
				return false
			}
		}
	}
	return defaultValue
}

func defaultString(value string, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

// newMySQLDriverConfig 組出 MySQL 連線描述，正式連線與探測用管理連線共用同一份
// 認證與傳輸解析；呼叫端只需另外補上各自的逾時處理。databaseName 留空代表管理連線。
func newMySQLDriverConfig(config ConnectionConfig, databaseName string) *mysqldriver.Config {
	mysqlConfig := mysqldriver.NewConfig()
	mysqlConfig.User = defaultString(
		stringConfigValue(config, "user"),
		stringConfigValue(config, "username"),
	)
	mysqlConfig.Passwd = stringConfigValue(config, "password")
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = fmt.Sprintf(
		"%s:%s",
		defaultString(stringConfigValue(config, "host"), "127.0.0.1"),
		defaultString(stringConfigValue(config, "port"), "3306"),
	)
	mysqlConfig.DBName = databaseName
	mysqlConfig.ParseTime = true
	mysqlConfig.AllowNativePasswords = true
	applyMySQLAuthTransportConfig(mysqlConfig, config)
	return mysqlConfig
}

// applyMySQLAuthTransportConfig 決定 cleartext 認證與 TLS 模式。cleartext 外掛會把
// 密碼以明文送給伺服器，因此預設關閉；明確啟用時必須走加密通道，且不得在伺服器
// 未提供 TLS 時退回明文（驅動的 preferred 模式即隱含該退回）。
func applyMySQLAuthTransportConfig(mysqlConfig *mysqldriver.Config, config ConnectionConfig) {
	mysqlConfig.AllowCleartextPasswords = booleanConfigValue(
		config,
		false,
		"allow_cleartext_passwords",
		"allowCleartextPasswords",
	)
	if tlsValue := strings.TrimSpace(stringConfigValue(config, "tls")); tlsValue != "" {
		mysqlConfig.TLSConfig = tlsValue
		return
	}
	if booleanConfigValue(config, false, "use_tls", "useTLS") || mysqlConfig.AllowCleartextPasswords {
		// skip-verify 實際加密但不驗證憑證，適用自簽憑證的工控現場；
		// 明確要求加密（或啟用 cleartext）時不得退回明文。
		mysqlConfig.TLSConfig = "skip-verify"
		mysqlConfig.AllowFallbackToPlaintext = false
		return
	}
	// 未指定時採機會性加密：伺服器支援 TLS 就走 TLS，不支援才退回明文。
	// 密碼此時走 native / caching_sha2 雜湊交握，不會以明文出現在線路上。
	mysqlConfig.TLSConfig = "preferred"
	mysqlConfig.AllowFallbackToPlaintext = true
}

func normalizeOptionalString(value string) string {
	return strings.TrimSpace(value)
}

func normalizeOptionalPointer(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeTimestampColumnForWriteMode(
	writeMode schema.DatabaseWriteMode,
	timestampColumn *string,
) *string {
	if writeMode != schema.DatabaseWriteModeUpsert {
		return nil
	}
	return normalizeOptionalPointer(timestampColumn)
}

func escapeSQLiteStringLiteral(value string) string {
	return strings.ReplaceAll(value, `'`, `''`)
}

func cloneConnectionConfig(config ConnectionConfig) ConnectionConfig {
	cloned := ConnectionConfig{}
	for key, value := range config {
		cloned[key] = value
	}
	return cloned
}

func preserveSensitiveConnectionConfigValues(
	existing ConnectionConfig,
	next ConnectionConfig,
	clearPassword bool,
) ConnectionConfig {
	merged := cloneConnectionConfig(next)
	for key, value := range next {
		merged[key] = value
	}

	if currentPassword := stringConfigValue(existing, "password"); currentPassword != "" {
		if clearPassword {
			delete(merged, "password")
		} else if stringConfigValue(next, "password") == "" {
			merged["password"] = currentPassword
		}
	}

	return merged
}

func optionalStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
