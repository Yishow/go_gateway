package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
)

type recordingPlanMigrationConnector struct {
	writeGroupMigrationConnector
	status string
}

func (s *WriteGroupService) previewRecordingPlanMigration(
	ctx context.Context,
	workspaceID string,
	sourceIDs []string,
) (result *WriteGroupMigrationPreview, err error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	workspaceID = strings.TrimSpace(workspaceID)
	if workspaceID == "" {
		return nil, fmt.Errorf("recording plan migration workspace id: %w", ErrWriteGroupValidation)
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
		return nil, normalizeWriteGroupServiceError("begin recording plan migration preview", err)
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback recording plan migration preview: %w", rollbackErr))
		}
	}()
	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("preview recording plan migration workspace")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read recording plan migration workspace", err)
	}
	if record.ID != workspaceID || strings.TrimSpace(record.DatabaseConnectorID) == "" {
		return nil, writeGroupNotFound("preview recording plan migration workspace")
	}
	connector, err := readRecordingPlanMigrationConnector(ctx, setupTx.SQLTx(), s.repo, record.DatabaseConnectorID)
	if err != nil {
		return nil, err
	}
	return buildRecordingPlanMigrationPreview(ctx, setupTx.SQLTx(), s.repo, record, connector, ids)
}

func buildRecordingPlanMigrationPreview(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	record *Record,
	connector recordingPlanMigrationConnector,
	ids []string,
) (*WriteGroupMigrationPreview, error) {
	planRepo := recordingplan.NewSQLRepository(repo.db)
	measurementRepo := measurement.NewSQLRepository(repo.db)
	preview := &WriteGroupMigrationPreview{
		WorkspaceID:       record.ID,
		WorkspaceRevision: record.DatabaseSetupRevision,
		ConnectorRevision: connector.identityRevision,
		AdapterVersion:    writeGroupRecordingPlanMigrationAdapterVersion,
		Items:             make([]WriteGroupMigrationPreviewItem, 0, len(ids)),
	}
	for _, id := range ids {
		item, err := buildRecordingPlanMigrationItem(ctx, tx, repo, planRepo, measurementRepo, record, connector, id)
		if err != nil {
			return nil, err
		}
		preview.Items = append(preview.Items, item)
	}
	digest, err := migrationPreviewDigest(preview)
	if err != nil {
		return nil, fmt.Errorf("derive recording plan migration preview digest: %w", err)
	}
	preview.ReviewDigest = digest
	return preview, nil
}

func buildRecordingPlanMigrationItem(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	planRepo *recordingplan.SQLRepository,
	measurementRepo *measurement.SQLRepository,
	record *Record,
	connector recordingPlanMigrationConnector,
	planID string,
) (WriteGroupMigrationPreviewItem, error) {
	plan, err := planRepo.GetPlanByWorkspaceInTx(ctx, tx, planID, record.ID)
	if errors.Is(err, recordingplan.ErrPlanNotFound) {
		return WriteGroupMigrationPreviewItem{}, writeGroupNotFound("preview recording plan migration source")
	}
	if err != nil {
		return WriteGroupMigrationPreviewItem{}, err
	}
	item := WriteGroupMigrationPreviewItem{
		SourceID:     plan.ID,
		Differences:  make([]WriteGroupMigrationFinding, 0, 2),
		Issues:       make([]WriteGroupMigrationFinding, 0, 8),
		RepairAction: "open_write_groups",
	}
	addIssue := func(code, message string) {
		item.Issues = append(item.Issues, WriteGroupMigrationFinding{Code: code, Message: message})
	}
	item.Differences = append(item.Differences, WriteGroupMigrationFinding{
		Code:    "recording-plan-no-equivalent-snapshot",
		Message: "the persisted recording plan does not prove an equivalent canonical snapshot column and row identity",
	})
	item.Issues = append(item.Issues, recordingPlanLayoutIssues(plan, connector)...)
	destinationConnectors, destinationIssues, err := readRecordingPlanDestinationConnectors(ctx, tx, repo, plan)
	if err != nil {
		return WriteGroupMigrationPreviewItem{}, err
	}
	item.Issues = append(item.Issues, destinationIssues...)
	intent := &WriteGroupRecordingPlanMigrationIntent{
		SourceID: plan.ID,
		Plan:     *plan,
		Sources:  make([]WriteGroupRecordingPlanMigrationSource, 0),
	}
	for _, measurementID := range recordingPlanSourceIDs(plan) {
		read, readErr := readRecordingPlanMigrationSource(ctx, tx, repo, measurementRepo, record, measurementID)
		if readErr != nil {
			return WriteGroupMigrationPreviewItem{}, readErr
		}
		intent.Sources = append(intent.Sources, recordingPlanSourceSnapshot(read))
		item.Issues = append(item.Issues, read.issues...)
	}
	if len(intent.Sources) == 0 {
		addIssue("recording-plan-source-missing", "the recording plan does not reference a measurement source")
	}
	item.BeforeRecordingPlanIntent = intent
	item.SourceRevision, err = recordingPlanMigrationSourceRevision(plan, connector, destinationConnectors, intent, item.Issues)
	if err != nil {
		return WriteGroupMigrationPreviewItem{}, err
	}
	item.Status = WriteGroupMigrationPreviewStatusBlocked
	if len(item.Issues) == 0 {
		addIssue("recording-plan-migration-blocked", "recording plan migration remains blocked until an equivalent canonical write group is explicitly designed")
	}
	return item, nil
}

func recordingPlanMigrationSourceRevision(
	plan *recordingplan.RecordingPlan,
	connector recordingPlanMigrationConnector,
	destinationConnectors []recordingPlanMigrationConnector,
	intent *WriteGroupRecordingPlanMigrationIntent,
	issues []WriteGroupMigrationFinding,
) (string, error) {
	planPayload, err := json.Marshal(plan)
	if err != nil {
		return "", fmt.Errorf("encode recording plan migration source revision: %w", err)
	}
	intentPayload, err := json.Marshal(intent)
	if err != nil {
		return "", fmt.Errorf("encode recording plan migration intent revision: %w", err)
	}
	issuePayload, err := json.Marshal(issues)
	if err != nil {
		return "", fmt.Errorf("encode recording plan migration issues: %w", err)
	}
	connectorConfigRevision, err := hashWriteGroupRevision(connector.connectionConfig)
	if err != nil {
		return "", fmt.Errorf("derive recording plan connector revision: %w", err)
	}
	connectorParts := make([]string, 0, len(destinationConnectors)*7)
	slices.SortFunc(destinationConnectors, func(left, right recordingPlanMigrationConnector) int {
		return strings.Compare(left.id, right.id)
	})
	for _, destinationConnector := range destinationConnectors {
		connectorParts = append(connectorParts,
			destinationConnector.id,
			destinationConnector.kind,
			destinationConnector.connectionConfig,
			destinationConnector.identityRevision,
			fmt.Sprint(destinationConnector.defaultWriteIntervalSecs),
			fmt.Sprint(destinationConnector.enabled),
			destinationConnector.status,
		)
	}
	return hashWriteGroupRevision(
		"recording-plan-preview",
		string(planPayload),
		string(intentPayload),
		string(issuePayload),
		connector.id,
		connector.kind,
		connector.identityRevision,
		connectorConfigRevision,
		fmt.Sprint(connector.defaultWriteIntervalSecs),
		fmt.Sprint(connector.enabled),
		connector.status,
		strings.Join(connectorParts, "\x00"),
	)
}

func readRecordingPlanDestinationConnectors(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	plan *recordingplan.RecordingPlan,
) ([]recordingPlanMigrationConnector, []WriteGroupMigrationFinding, error) {
	ids := make([]string, 0, len(plan.Destinations))
	seen := make(map[string]struct{}, len(plan.Destinations))
	for _, destination := range plan.Destinations {
		id := strings.TrimSpace(destination.ConnectorID)
		if id == "" {
			continue
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	slices.Sort(ids)
	connectors := make([]recordingPlanMigrationConnector, 0, len(ids))
	issues := make([]WriteGroupMigrationFinding, 0)
	for _, id := range ids {
		connector, err := readRecordingPlanMigrationConnector(ctx, tx, repo, id)
		if errors.Is(err, ErrWriteGroupNotFound) || errors.Is(err, ErrNotFound) {
			issues = append(issues, WriteGroupMigrationFinding{
				Code:    "recording-plan-destination-unresolved",
				Message: "the recording plan references a saved connector that no longer exists",
			})
			continue
		}
		if err != nil {
			return nil, nil, err
		}
		connectors = append(connectors, connector)
	}
	return connectors, issues, nil
}

func readRecordingPlanMigrationConnector(
	ctx context.Context,
	tx *sql.Tx,
	repo *SQLWriteGroupRepository,
	connectorID string,
) (recordingPlanMigrationConnector, error) {
	var connector recordingPlanMigrationConnector
	err := tx.QueryRowContext(ctx, repo.query(`
		SELECT id, kind, connection_config, identity_revision,
			COALESCE(default_write_interval_seconds, 0), enabled, status
		FROM database_connectors
		WHERE id = $1
	`), connectorID).Scan(
		&connector.id,
		&connector.kind,
		&connector.connectionConfig,
		&connector.identityRevision,
		&connector.defaultWriteIntervalSecs,
		&connector.enabled,
		&connector.status,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return recordingPlanMigrationConnector{}, writeGroupNotFound("preview recording plan migration connector")
	}
	if err != nil {
		return recordingPlanMigrationConnector{}, fmt.Errorf("read recording plan migration connector: %w", err)
	}
	connector.kind = strings.ToLower(strings.TrimSpace(connector.kind))
	connector.identityRevision = strings.TrimSpace(connector.identityRevision)
	connector.status = strings.TrimSpace(connector.status)
	return connector, nil
}
