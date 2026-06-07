package dbtarget

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
)

type tableKey struct {
	schema string
	name   string
}

type generatedTableState struct {
	schemaName string
	tableName  string
	columns    map[string]ColumnInfo
}

func (s *ConnectorService) GenerateSchema(
	ctx context.Context,
	connectorID string,
	req SchemaGenerateRequest,
) (*SchemaGenerateResult, error) {
	connector, err := s.repo.GetByID(ctx, connectorID)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}
	recordFailure := func(cause error) (*SchemaGenerateResult, error) {
		if !req.DryRun {
			if outcomeErr := s.recordSchemaEnsureOutcome(ctx, connector, time.Now(), deliveryOutcomeFailed, cause.Error()); outcomeErr != nil {
				return nil, fmt.Errorf("%w; %v", cause, outcomeErr)
			}
		}
		return nil, cause
	}
	recordSuccess := func(result *SchemaGenerateResult) (*SchemaGenerateResult, error) {
		if !req.DryRun {
			if outcomeErr := s.recordSchemaEnsureOutcome(ctx, connector, time.Now(), deliveryOutcomeSuccess, ""); outcomeErr != nil {
				return nil, outcomeErr
			}
		}
		return result, nil
	}
	if s.mappingRepo == nil {
		return recordFailure(fmt.Errorf("資料庫目標映射儲存庫未配置"))
	}

	projection, err := listLiveTargetProjection(
		ctx,
		s.mappingRepo,
		TargetMappingListFilter{ConnectorID: &connectorID},
		s.connectorTagReader(),
	)
	if err != nil {
		return recordFailure(err)
	}

	result := &SchemaGenerateResult{
		ConnectorID: connector.ID,
		DryRun:      req.DryRun,
		Statements:  []string{},
		Executed:    0,
	}
	if len(projection.Mappings) == 0 {
		return recordSuccess(result)
	}

	tables, err := inspectTables(ctx, connector)
	if err != nil {
		return recordFailure(fmt.Errorf("檢查資料庫目標表結構失敗: %w", err))
	}

	statements, err := buildSchemaGenerateStatements(
		ctx,
		connector.Kind,
		projection.Mappings,
		tables,
		s.connectorTagReader(),
	)
	if err != nil {
		return recordFailure(err)
	}

	result.Statements = statements
	if req.DryRun || len(statements) == 0 {
		return recordSuccess(result)
	}

	connectionConfig, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return recordFailure(err)
	}
	manager, err := openExternalDBManager(connector.Kind, connectionConfig)
	if err != nil {
		return recordFailure(err)
	}
	defer manager.Close()

	for _, statement := range statements {
		if _, execErr := manager.DB().ExecContext(ctx, statement); execErr != nil {
			return recordFailure(fmt.Errorf("執行 schema statement 失敗: %w", execErr))
		}
	}

	result.Executed = len(statements)
	return recordSuccess(result)
}

func (s *ConnectorService) ListWriteHistory(
	ctx context.Context,
	connectorID string,
	limit int,
) ([]WriteHistoryRecord, error) {
	if _, err := s.repo.GetByID(ctx, connectorID); err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}
	return listWriteHistory(connectorID, limit), nil
}

func (s *MappingService) DryRun(
	ctx context.Context,
	connectorID string,
	req MappingDryRunRequest,
) (*MappingDryRunResult, error) {
	connector, err := s.connectorRepo.GetByID(ctx, connectorID)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}

	projection, err := listLiveTargetProjection(ctx, s.repo, TargetMappingListFilter{ConnectorID: &connectorID}, s.tagService)
	if err != nil {
		return nil, err
	}

	normalizedCandidateIDs, unknownCandidateIDs := normalizeDryRunCandidateIDs(req.CandidateIDs, projection.Mappings)
	selectedMappings := selectMappingsForDryRun(projection.Mappings, normalizedCandidateIDs)
	results := make([]MappingDryRunCandidateResult, 0, len(selectedMappings)+len(unknownCandidateIDs))

	tables, inspectErr := inspectTables(ctx, connector)
	if inspectErr != nil {
		for _, mappingRecord := range selectedMappings {
			results = append(results, MappingDryRunCandidateResult{
				CandidateID: mappingRecord.ID,
				MappingID:   mappingRecord.ID,
				TagID:       mappingRecord.TagID,
				Status:      "blocked",
				Code:        "connector_unavailable",
				Reason:      inspectErr.Error(),
			})
		}
		for _, candidateID := range unknownCandidateIDs {
			results = append(results, MappingDryRunCandidateResult{
				CandidateID: candidateID,
				Status:      "failed",
				Code:        "validation",
				Reason:      "candidate is not found",
			})
		}
		return &MappingDryRunResult{
			ConnectorID: connector.ID,
			Results:     results,
		}, nil
	}

	for _, mappingRecord := range selectedMappings {
		tagEntity := projection.TagsByMappingID[mappingRecord.ID]
		if tagEntity == nil {
			var getTagErr error
			tagEntity, getTagErr = s.tagService.GetByID(ctx, mappingRecord.TagID)
			if getTagErr != nil {
				results = append(results, MappingDryRunCandidateResult{
					CandidateID: mappingRecord.ID,
					MappingID:   mappingRecord.ID,
					TagID:       mappingRecord.TagID,
					Status:      "blocked",
					Code:        "schema_missing",
					Reason:      getTagErr.Error(),
				})
				continue
			}
		}

		issues := validateMappingAgainstTables(*mappingRecord, *tagEntity, tables)
		blockingIssue, hasBlockingIssue := firstBlockingIssue(issues)
		if hasBlockingIssue {
			results = append(results, MappingDryRunCandidateResult{
				CandidateID: mappingRecord.ID,
				MappingID:   mappingRecord.ID,
				TagID:       mappingRecord.TagID,
				Status:      "blocked",
				Code:        mapValidationIssueCodeForDryRun(blockingIssue.Code),
				Reason:      blockingIssue.Message,
			})
			continue
		}

		results = append(results, MappingDryRunCandidateResult{
			CandidateID: mappingRecord.ID,
			MappingID:   mappingRecord.ID,
			TagID:       mappingRecord.TagID,
			Status:      "ready",
		})
	}

	for _, candidateID := range unknownCandidateIDs {
		results = append(results, MappingDryRunCandidateResult{
			CandidateID: candidateID,
			Status:      "failed",
			Code:        "validation",
			Reason:      "candidate is not found",
		})
	}

	return &MappingDryRunResult{
		ConnectorID: connector.ID,
		Results:     results,
	}, nil
}

func buildSchemaGenerateStatements(
	ctx context.Context,
	kind schema.DatabaseConnectorKind,
	mappings []*schema.DatabaseTargetMapping,
	tables []TableInfo,
	tagReader ConnectorTagReader,
) ([]string, error) {
	statements := make([]string, 0, 16)
	tableStates := make(map[tableKey]*generatedTableState, len(tables))
	createdSchemas := make(map[string]struct{})
	createdUniqueIndexes := make(map[string]struct{})

	for _, table := range tables {
		state := &generatedTableState{
			schemaName: table.Schema,
			tableName:  table.Name,
			columns:    make(map[string]ColumnInfo, len(table.Columns)),
		}
		for _, column := range table.Columns {
			state.columns[strings.ToLower(strings.TrimSpace(column.Name))] = column
		}
		tableStates[normalizeTableKey(table.Schema, table.Name)] = state
	}

	for _, mappingRecord := range mappings {
		if mappingRecord == nil {
			continue
		}
		if !mappingRecord.Enabled {
			continue
		}
		tableName := strings.TrimSpace(mappingRecord.TableName)
		columnName := strings.TrimSpace(mappingRecord.ColumnName)
		if tableName == "" || columnName == "" {
			continue
		}
		schemaName := strings.TrimSpace(mappingRecord.TableSchema)
		if schemaName == "" {
			schemaName = defaultSchemaForKind(kind)
		}
		key := normalizeTableKey(schemaName, tableName)

		valueColumnType, resolveErr := resolveMappingColumnType(ctx, tagReader, mappingRecord.TagID, kind)
		if resolveErr != nil {
			if errors.Is(resolveErr, tag.ErrTagNotFound) {
				continue
			}
			return nil, resolveErr
		}

		state, exists := tableStates[key]
		if !exists {
			if kind == schema.DatabaseConnectorKindPostgres {
				schemaDDL := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", quoteIdentifier(kind, schemaName))
				if _, seen := createdSchemas[schemaDDL]; !seen {
					statements = append(statements, schemaDDL)
					createdSchemas[schemaDDL] = struct{}{}
				}
			}
			createStatement := buildCreateTableStatement(kind, schemaName, tableName, columnName, valueColumnType, mappingRecord)
			statements = append(statements, createStatement)
			state = &generatedTableState{
				schemaName: schemaName,
				tableName:  tableName,
				columns:    map[string]ColumnInfo{},
			}
			state.columns[strings.ToLower(columnName)] = ColumnInfo{Name: columnName, DataType: valueColumnType}
			if mappingRecord.WriteMode == schema.DatabaseWriteModeUpsert && mappingRecord.TimestampColumn != nil {
				ts := strings.TrimSpace(*mappingRecord.TimestampColumn)
				if ts != "" {
					state.columns[strings.ToLower(ts)] = ColumnInfo{
						Name:       ts,
						DataType:   timestampColumnType(kind),
						PrimaryKey: true,
						Unique:     true,
					}
				}
			}
			tableStates[key] = state
			continue
		}

		valueColumnKey := strings.ToLower(columnName)
		if _, hasColumn := state.columns[valueColumnKey]; !hasColumn {
			statements = append(statements, buildAddColumnStatement(kind, schemaName, tableName, columnName, valueColumnType))
			state.columns[valueColumnKey] = ColumnInfo{Name: columnName, DataType: valueColumnType}
		}

		if mappingRecord.WriteMode != schema.DatabaseWriteModeUpsert || mappingRecord.TimestampColumn == nil {
			continue
		}

		timestampColumn := strings.TrimSpace(*mappingRecord.TimestampColumn)
		if timestampColumn == "" {
			continue
		}
		timestampColumnKey := strings.ToLower(timestampColumn)
		tsColumn, hasTimestampColumn := state.columns[timestampColumnKey]
		if !hasTimestampColumn {
			statements = append(statements, buildAddColumnStatement(kind, schemaName, tableName, timestampColumn, timestampColumnType(kind)))
			tsColumn = ColumnInfo{Name: timestampColumn, DataType: timestampColumnType(kind)}
			state.columns[timestampColumnKey] = tsColumn
		}

		if tsColumn.PrimaryKey || tsColumn.Unique {
			continue
		}
		indexStatement, indexKey := buildEnsureUniqueIndexStatement(kind, schemaName, tableName, timestampColumn)
		if _, exists := createdUniqueIndexes[indexKey]; exists {
			continue
		}
		statements = append(statements, indexStatement)
		createdUniqueIndexes[indexKey] = struct{}{}
		tsColumn.Unique = true
		state.columns[timestampColumnKey] = tsColumn
	}

	return dedupeStatements(statements), nil
}
