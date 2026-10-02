package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/schema"
)

type writeGroupMigrationTarget struct {
	id                   string
	tagID                string
	connectorID          string
	tableSchema          string
	tableName            string
	columnName           string
	writeMode            string
	timestampColumn      *string
	groupKey             *string
	writeIntervalSeconds *int
	enabled              bool
}

type writeGroupMigrationConnector struct {
	id                       string
	kind                     string
	connectionConfig         string
	identityRevision         string
	defaultWriteIntervalSecs int
	enabled                  bool
}

type writeGroupMigrationSource struct {
	mapping      mappingRevisionSource
	source       sourceRevisionSource
	mappingState string
	deviceStatus string
	tagStatus    string
}

func (s *WriteGroupService) previewSingleMappingMigration(
	ctx context.Context,
	workspaceID string,
	sourceIDs []string,
) (result *WriteGroupMigrationPreview, err error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	workspaceID = strings.TrimSpace(workspaceID)
	if workspaceID == "" {
		return nil, fmt.Errorf("migration preview workspace id: %w", ErrWriteGroupValidation)
	}
	ids, err := normalizeMigrationPreviewIDs(sourceIDs)
	if err != nil {
		return nil, err
	}
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupServiceUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("begin migration preview", err)
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback migration preview read: %w", rollbackErr))
		}
	}()

	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("preview migration workspace")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read migration preview workspace", err)
	}
	if record.ID != workspaceID || strings.TrimSpace(record.DatabaseConnectorID) == "" {
		return nil, writeGroupNotFound("preview migration workspace")
	}
	tx := setupTx.SQLTx()
	connector, err := readMigrationPreviewConnector(ctx, tx, s.repo, record.DatabaseConnectorID)
	if err != nil {
		return nil, err
	}
	return buildSingleMappingMigrationPreview(ctx, tx, s.repo, record, connector, ids)
}

func buildSingleMappingMigrationPreview(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	record *Record,
	connector writeGroupMigrationConnector,
	ids []string,
) (*WriteGroupMigrationPreview, error) {
	preview := &WriteGroupMigrationPreview{
		WorkspaceID:       record.ID,
		WorkspaceRevision: record.DatabaseSetupRevision,
		ConnectorRevision: connector.identityRevision,
		AdapterVersion:    writeGroupMigrationPreviewAdapterVersion,
		Items:             make([]WriteGroupMigrationPreviewItem, 0, len(ids)),
	}
	for _, id := range ids {
		item, err := readMigrationPreviewItem(ctx, tx, repo, record, connector, id)
		if err != nil {
			return nil, err
		}
		preview.Items = append(preview.Items, item)
	}
	digest, err := migrationPreviewDigest(preview)
	if err != nil {
		return nil, fmt.Errorf("derive migration preview digest: %w", err)
	}
	preview.ReviewDigest = digest
	return preview, nil
}

func normalizeMigrationPreviewIDs(sourceIDs []string) ([]string, error) {
	if len(sourceIDs) == 0 {
		return nil, fmt.Errorf("migration preview source ids: %w", ErrWriteGroupValidation)
	}
	ids := make([]string, 0, len(sourceIDs))
	seen := make(map[string]struct{}, len(sourceIDs))
	for _, rawID := range sourceIDs {
		id := strings.TrimSpace(rawID)
		if id == "" {
			return nil, fmt.Errorf("migration preview source id: %w", ErrWriteGroupValidation)
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	slices.Sort(ids)
	if len(ids) == 0 {
		return nil, fmt.Errorf("migration preview source ids: %w", ErrWriteGroupValidation)
	}
	return ids, nil
}

func readMigrationPreviewConnector(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	connectorID string,
) (writeGroupMigrationConnector, error) {
	var connector writeGroupMigrationConnector
	err := tx.QueryRowContext(ctx, repo.query(`
		SELECT id, kind, connection_config, identity_revision,
			COALESCE(default_write_interval_seconds, 0), enabled
		FROM database_connectors
		WHERE id = $1
	`), connectorID).Scan(
		&connector.id,
		&connector.kind,
		&connector.connectionConfig,
		&connector.identityRevision,
		&connector.defaultWriteIntervalSecs,
		&connector.enabled,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return writeGroupMigrationConnector{}, writeGroupNotFound("preview migration connector")
	}
	if err != nil {
		return writeGroupMigrationConnector{}, fmt.Errorf("read migration preview connector: %w", err)
	}
	if !connector.enabled {
		return writeGroupMigrationConnector{}, fmt.Errorf("preview migration connector: %w", ErrWriteGroupValidation)
	}
	connector.kind = strings.ToLower(strings.TrimSpace(connector.kind))
	connector.identityRevision = strings.TrimSpace(connector.identityRevision)
	return connector, nil
}

func readMigrationPreviewItem(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	record *Record,
	connector writeGroupMigrationConnector,
	sourceID string,
) (WriteGroupMigrationPreviewItem, error) {
	target, err := readMigrationPreviewTarget(ctx, tx, repo, sourceID)
	if err != nil {
		return WriteGroupMigrationPreviewItem{}, err
	}
	if target.connectorID != connector.id {
		return WriteGroupMigrationPreviewItem{}, writeGroupNotFound("preview migration source")
	}
	allTargets, err := readMigrationPreviewTargetsForTag(ctx, tx, repo, target.tagID)
	if err != nil {
		return WriteGroupMigrationPreviewItem{}, err
	}
	connectorTargets, err := readMigrationPreviewTargetsForConnector(ctx, tx, repo, connector.id)
	if err != nil {
		return WriteGroupMigrationPreviewItem{}, err
	}
	allSources, err := readMigrationPreviewSources(ctx, tx, repo, target.tagID)
	if err != nil {
		return WriteGroupMigrationPreviewItem{}, err
	}

	item := WriteGroupMigrationPreviewItem{
		SourceID:    sourceID,
		Differences: make([]WriteGroupMigrationFinding, 0),
		Issues:      make([]WriteGroupMigrationFinding, 0),
	}
	addIssue := func(code, message string) {
		item.Issues = append(item.Issues, WriteGroupMigrationFinding{Code: code, Message: message})
	}
	if !target.enabled {
		addIssue("target-mapping-disabled", "the selected legacy target mapping is disabled")
	}
	if target.writeMode != string(schema.DatabaseWriteModeInsert) && target.writeMode != string(schema.DatabaseWriteModeUpsert) {
		addIssue("write-mode-unsupported", "the persisted legacy write mode is unsupported")
	}
	if countEnabledMigrationPreviewTargets(allTargets) > 1 {
		addIssue("multiple-target-mappings", "more than one legacy target mapping exists for this tag")
	}
	item.Issues = append(item.Issues, migrationPreviewTargetIssues(target, connectorTargets)...)
	if connector.kind != writeGroupConnectorKindSQLite && connector.kind != writeGroupConnectorKindPostgres {
		addIssue("connector-kind-out-of-batch", "the saved connector kind is outside this preview adapter")
	}
	if destination, scopeErr := migrationPreviewDestination(target, connector); scopeErr != nil {
		addIssue("connector-scope-invalid", "the saved connector scope cannot be interpreted safely")
	} else if destination.TableSchema != target.tableSchema {
		addIssue("connector-scope-mismatch", "the saved connector scope cannot preserve the legacy destination schema")
	}
	if target.groupKey != nil && strings.TrimSpace(*target.groupKey) != "" {
		addIssue("group-key-out-of-batch", "grouped legacy mappings require a separate migration review")
	}
	if strings.EqualFold(strings.TrimSpace(target.writeMode), string(schema.DatabaseWriteModeUpsert)) {
		addIssue("upsert-semantic-difference", "legacy upsert behavior is not equivalent to a basic snapshot append")
	}

	source, sourceIssues, sourceErr := selectMigrationPreviewSource(record, allSources)
	if sourceErr != nil {
		return WriteGroupMigrationPreviewItem{}, sourceErr
	}
	item.Issues = append(item.Issues, sourceIssues...)
	if source != nil {
		mappingRevision, err := source.mapping.revision()
		if err != nil {
			return WriteGroupMigrationPreviewItem{}, fmt.Errorf("derive migration mapping revision: %w", err)
		}
		sourceRevision, err := source.source.revision()
		if err != nil {
			return WriteGroupMigrationPreviewItem{}, fmt.Errorf("derive migration source revision: %w", err)
		}
		sourceDigest, err := migrationPreviewSourceDigest(source, target, connector)
		if err != nil {
			return WriteGroupMigrationPreviewItem{}, err
		}
		item.SourceRevision = sourceDigest
		item.BeforeIntent = migrationPreviewIntent(source, target, connector)
		if len(item.Issues) == 0 {
			item.Differences = migrationPreviewDifferences(target, connector)
			candidate := migrationPreviewCandidate(record, source, target, connector, sourceDigest, sourceRevision, mappingRevision)
			ownershipIssues, ownershipErr := migrationSourceOwnershipIssues(ctx, tx, repo, record.ID, candidate.Migration.SourceKind, sourceID, candidate.Migration)
			if ownershipErr != nil {
				return WriteGroupMigrationPreviewItem{}, ownershipErr
			}
			item.Issues = append(item.Issues, ownershipIssues...)
			if len(ownershipIssues) == 0 {
				item.CandidateGroup = candidate
			}
		}
	}
	if len(item.Issues) == 0 {
		item.Status = WriteGroupMigrationPreviewStatusNeedsReview
	} else {
		item.Status = WriteGroupMigrationPreviewStatusBlocked
	}
	return item, nil
}

func selectMigrationPreviewSource(record *Record, sources []writeGroupMigrationSource) (*writeGroupMigrationSource, []WriteGroupMigrationFinding, error) {
	enabled := make([]writeGroupMigrationSource, 0, len(sources))
	active := make([]writeGroupMigrationSource, 0, len(sources))
	hasWorkspaceSource := false
	for _, source := range sources {
		if slices.Contains(record.OrderedDeviceIDs, source.source.DeviceID) {
			hasWorkspaceSource = true
		}
		if source.mapping.Enabled {
			enabled = append(enabled, source)
		}
		if source.mapping.Enabled && source.mappingState == string(schema.MappingStatusActive) {
			active = append(active, source)
		}
	}
	if len(sources) > 0 && !hasWorkspaceSource {
		return nil, nil, writeGroupNotFound("preview migration source")
	}
	if len(enabled) > 1 {
		return nil, []WriteGroupMigrationFinding{{Code: "multiple-enabled-source-mappings", Message: "more than one enabled source mapping reaches this tag"}}, nil
	}
	if len(active) == 0 {
		return nil, []WriteGroupMigrationFinding{{Code: "source-mapping-missing-or-disabled", Message: "no enabled active source mapping reaches this tag"}}, nil
	}
	source := active[0]
	if !slices.Contains(record.OrderedDeviceIDs, source.source.DeviceID) {
		return nil, nil, writeGroupNotFound("preview migration source")
	}
	issues := make([]WriteGroupMigrationFinding, 0, 4)
	if source.deviceStatus == string(schema.DeviceStatusDisabled) {
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-device-disabled", Message: "the source device is disabled"})
	}
	if !source.source.PointEnabled {
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-point-disabled", Message: "the source point is disabled"})
	}
	if source.tagStatus != string(schema.TagStatusActive) {
		issues = append(issues, WriteGroupMigrationFinding{Code: "source-tag-missing-or-disabled", Message: "the source tag is not active"})
	}
	if pipeline := strings.TrimSpace(source.mapping.TransformPipeline); pipeline != "" && pipeline != "[]" {
		issues = append(issues, WriteGroupMigrationFinding{Code: "advanced-transform-pipeline", Message: "the mapping has a non-empty transform pipeline"})
	}
	if source.mapping.RuleCandidateID != "" || source.mapping.ProposedSignature != "" || source.mapping.LastAppliedSignature != "" || source.mapping.BlockingReason != "" {
		issues = append(issues, WriteGroupMigrationFinding{Code: "advanced-mapping-rule", Message: "the mapping carries rule or blocking metadata outside this adapter"})
	}
	return &source, issues, nil
}

func countEnabledMigrationPreviewTargets(targets []writeGroupMigrationTarget) int {
	count := 0
	for _, target := range targets {
		if target.enabled {
			count++
		}
	}
	return count
}

func migrationPreviewTargetIssues(target writeGroupMigrationTarget, connectorTargets []writeGroupMigrationTarget) []WriteGroupMigrationFinding {
	issues := make([]WriteGroupMigrationFinding, 0, 3)
	if strings.TrimSpace(target.tableSchema) == "" || strings.TrimSpace(target.tableName) == "" || strings.TrimSpace(target.columnName) == "" {
		issues = append(issues, WriteGroupMigrationFinding{Code: "target-scope-incomplete", Message: "the legacy target mapping has an incomplete table or column scope"})
	}
	if target.writeIntervalSeconds != nil && *target.writeIntervalSeconds < 0 {
		issues = append(issues, WriteGroupMigrationFinding{Code: "negative-write-interval", Message: "the legacy target mapping has a negative write interval"})
	}
	for _, other := range connectorTargets {
		if other.id == target.id || !other.enabled {
			continue
		}
		if other.tableSchema != target.tableSchema || other.tableName != target.tableName || other.tagID == target.tagID {
			continue
		}
		if other.columnName == target.columnName {
			issues = append(issues, WriteGroupMigrationFinding{Code: "shared-target-column", Message: "multiple enabled legacy mappings share one destination column"})
		} else if target.groupKey == nil && other.groupKey == nil {
			issues = append(issues, WriteGroupMigrationFinding{Code: "row-identity-ambiguous", Message: "multiple enabled legacy mappings share a table without a persisted row identity"})
		}
	}
	return issues
}
