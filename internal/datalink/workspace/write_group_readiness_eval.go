package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"go-gateway/internal/datalink/schema"
)

type writeGroupReadinessSnapshot struct {
	record          *Record
	group           *WriteGroup
	validated       *WriteGroup
	tagTypes        map[string]schema.DataType
	connectorKind   string
	connectorConfig string
	validateErr     error
}

func (s *WriteGroupService) readiness(ctx context.Context, id string) (*WriteGroupReadiness, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, writeGroupNotFound("read write-group readiness")
	}
	snapshot, err := s.readinessSnapshot(ctx, id)
	if err != nil {
		return nil, err
	}
	result := newWriteGroupReadiness(snapshot.record, snapshot.group)
	if snapshot.validateErr != nil {
		if issue, ok := writeGroupReadinessValidationIssue(snapshot.validateErr, id); ok {
			result.Issues = append(result.Issues, issue)
		} else {
			return nil, normalizeWriteGroupServiceError("evaluate write-group readiness", snapshot.validateErr)
		}
	} else {
		result.Issues = append(result.Issues, validateBasicWriteGroupConfig(snapshot.validated)...)
	}
	result.ConfigReady = len(result.Issues) == 0
	if s.tableInspector == nil {
		result.SchemaReady = false
		result.Issues = append(result.Issues, writeGroupReadinessIssue(
			"schema-unverified", "destination schema has not been verified", id,
		))
		return finalizeWriteGroupReadiness(result), nil
	}
	if !result.ConfigReady {
		result.SchemaReady = false
		result.Issues = append(result.Issues, writeGroupReadinessIssue(
			"schema-check-blocked", "destination schema check is blocked by configuration issues", id,
		))
		return finalizeWriteGroupReadiness(result), nil
	}

	allowed, issue := allowWriteGroupTableInspection(snapshot)
	if !allowed {
		result.SchemaReady = false
		result.Issues = append(result.Issues, issue)
		return finalizeWriteGroupReadiness(result), nil
	}
	inspection, inspectErr := s.tableInspector.InspectTable(ctx,
		snapshot.validated.Destination.ConnectorID,
		snapshot.validated.Destination.TableSchema,
		snapshot.validated.Destination.TableName,
	)
	postSnapshot, err := s.readinessSnapshot(ctx, id)
	if err != nil {
		return nil, err
	}
	if !sameWriteGroupReadinessSnapshot(snapshot, postSnapshot) {
		result = newWriteGroupReadiness(postSnapshot.record, postSnapshot.group)
		if postSnapshot.validateErr != nil {
			result.ConfigReady = false
			result.SchemaReady = false
			result.Issues = append(result.Issues, readinessValidationIssueOrInternal(postSnapshot.validateErr, id))
			return writeGroupReadinessResult(finalizeWriteGroupReadiness(result))
		}
		result.Issues = append(result.Issues, writeGroupReadinessIssue(
			"readiness-snapshot-changed", "configuration changed during schema verification", id,
		))
		return finalizeWriteGroupReadiness(result), nil
	}
	if postSnapshot.validateErr != nil {
		result.ConfigReady = false
		result.SchemaReady = false
		result.Issues = append(result.Issues, readinessValidationIssueOrInternal(postSnapshot.validateErr, id))
		return writeGroupReadinessResult(finalizeWriteGroupReadiness(result))
	}
	if inspectErr != nil {
		result.Issues = append(result.Issues, writeGroupReadinessIssue(
			"schema-inspection-failed", "destination schema could not be verified", id,
		))
		return writeGroupReadinessResult(finalizeWriteGroupReadiness(result))
	}
	schemaReady, schemaIssues := evaluateWriteGroupInspection(
		postSnapshot.validated, postSnapshot.tagTypes, inspection,
	)
	result.SchemaReady = schemaReady
	result.Issues = append(result.Issues, schemaIssues...)
	if schemaReady {
		result.SchemaDigest, err = digestWriteGroupInspection(postSnapshot.validated, inspection)
		if err != nil {
			return nil, fmt.Errorf("digest write-group schema readiness: %w", err)
		}
	}
	return finalizeWriteGroupReadiness(result), nil
}

func (s *WriteGroupService) readinessSnapshot(ctx context.Context, id string) (result *writeGroupReadinessSnapshot, err error) {
	store, ok := s.workspaceSvc.repo.(DatabaseSetupStore)
	if !ok {
		return nil, ErrWriteGroupReadinessUnavailable
	}
	setupTx, err := store.BeginDatabaseSetup(ctx)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("begin write-group readiness", err)
	}
	defer func() {
		if rollbackErr := setupTx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			result = nil
			err = errors.Join(err, fmt.Errorf("rollback write-group readiness read: %w", rollbackErr))
		}
	}()
	record, err := setupTx.Get(ctx)
	if errors.Is(err, ErrNotFound) {
		return nil, writeGroupNotFound("read write-group readiness")
	}
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read write-group readiness", err)
	}
	group, err := s.repo.GetInTx(ctx, setupTx.SQLTx(), record.ID, id)
	if err != nil {
		return nil, normalizeWriteGroupServiceError("read write-group readiness", err)
	}
	validated := cloneWriteGroup(group)
	validationErr := s.repo.validate(ctx, setupTx.SQLTx(), validated)
	tagTypes := map[string]schema.DataType{}
	var connectorKind, connectorConfig string
	connectorErr := setupTx.SQLTx().QueryRowContext(ctx, s.repo.query(
		"SELECT kind, connection_config FROM database_connectors WHERE id = $1",
	), group.Destination.ConnectorID).Scan(&connectorKind, &connectorConfig)
	if connectorErr != nil && !errors.Is(connectorErr, sql.ErrNoRows) {
		return nil, fmt.Errorf("read write-group connector scope: %w", connectorErr)
	}
	connectorKind = strings.ToLower(strings.TrimSpace(connectorKind))
	if validationErr == nil {
		tagTypes, err = loadWriteGroupTagTypes(ctx, setupTx.SQLTx(), s.repo.query, validated)
		if err != nil {
			validationErr = err
		}
	}
	return &writeGroupReadinessSnapshot{
		record: record, group: group, validated: validated,
		tagTypes: tagTypes, connectorKind: connectorKind, connectorConfig: connectorConfig,
		validateErr: validationErr,
	}, nil
}

func newWriteGroupReadiness(record *Record, group *WriteGroup) *WriteGroupReadiness {
	return &WriteGroupReadiness{
		WorkspaceID: record.ID, WorkspaceRevision: record.DatabaseSetupRevision,
		GroupID: group.ID, GroupRevision: group.Revision, AppliedRevision: group.AppliedRevision,
		Issues: []ReadinessIssue{},
	}
}

func finalizeWriteGroupReadiness(readiness *WriteGroupReadiness) *WriteGroupReadiness {
	readiness.Ready = readiness.ConfigReady && readiness.SchemaReady && len(readiness.Issues) == 0
	return readiness
}

func writeGroupReadinessResult(readiness *WriteGroupReadiness) (*WriteGroupReadiness, error) {
	return readiness, nil
}

func validateBasicWriteGroupConfig(group *WriteGroup) []ReadinessIssue {
	if group == nil {
		return []ReadinessIssue{writeGroupReadinessIssue("group-invalid", "write group is unavailable", "")}
	}
	issues := make([]ReadinessIssue, 0)
	destination := group.Destination
	if strings.TrimSpace(destination.ConnectorID) == "" ||
		strings.TrimSpace(destination.ConnectorRevision) == "" ||
		strings.TrimSpace(destination.Database) == "" ||
		strings.TrimSpace(destination.TableSchema) == "" ||
		strings.TrimSpace(destination.TableName) == "" {
		issues = append(issues, writeGroupReadinessIssue(
			"storage-binding-missing", "write group storage binding is incomplete", group.ID,
		))
	}
	if destination.StorageStrategy != WriteGroupStorageStrategyManaged &&
		destination.StorageStrategy != WriteGroupStorageStrategyCustom {
		issues = append(issues, writeGroupReadinessIssue(
			"storage-strategy-unsupported", "write group storage strategy is unsupported", group.ID,
		))
	}
	if group.RowPolicy.IntervalSeconds <= 0 {
		issues = append(issues, writeGroupReadinessIssue(
			"interval-required", "write group interval must be positive", group.ID,
		))
	}
	if group.RowPolicy.AllowedLatenessSeconds < 0 {
		issues = append(issues, writeGroupReadinessIssue(
			"lateness-invalid", "write group allowed lateness cannot be negative", group.ID,
		))
	}
	incompletePolicy := strings.ToLower(strings.TrimSpace(group.RowPolicy.IncompletePolicy))
	if incompletePolicy != "" && incompletePolicy != writeGroupIncompletePolicySkipRow {
		issues = append(issues, writeGroupReadinessIssue(
			"incomplete-policy-blocked", "selected incomplete-row policy is not supported by basic readiness", group.ID,
		))
	}
	writeMode := strings.ToLower(strings.TrimSpace(group.WritePolicy.Mode))
	if writeMode != "" && writeMode != writeGroupWriteModeAppend {
		issues = append(issues, writeGroupReadinessIssue(
			"write-policy-blocked", "selected write policy is not supported by basic readiness", group.ID,
		))
	}
	for index, member := range group.Members {
		if member.MaxAgeSeconds != nil && *member.MaxAgeSeconds < 0 {
			issues = append(issues, writeGroupReadinessIssue(
				"member-max-age-invalid", fmt.Sprintf("member %d max age cannot be negative", index), group.ID,
			))
		}
	}
	return issues
}

func writeGroupReadinessValidationIssue(err error, scope string) (ReadinessIssue, bool) {
	switch {
	case errors.Is(err, ErrWriteGroupConnectorRevisionConflict):
		return writeGroupReadinessIssue("connector-revision-stale", "saved connector identity revision is stale", scope), true
	case errors.Is(err, ErrWriteGroupSourceRevisionConflict):
		return writeGroupReadinessIssue("source-revision-stale", "saved source or mapping revision is stale", scope), true
	case errors.Is(err, ErrWriteGroupConnectorNotFound):
		return writeGroupReadinessIssue("connector-missing", "saved connector is no longer available", scope), true
	case errors.Is(err, ErrWriteGroupUnsupportedConnectorKind):
		return writeGroupReadinessIssue("connector-kind-unsupported", "saved connector kind is unsupported", scope), true
	case errors.Is(err, ErrWriteGroupValidation):
		return writeGroupReadinessIssue("group-invalid", "write group configuration is invalid", scope), true
	case errors.Is(err, ErrWriteGroupNotFound), errors.Is(err, ErrNotFound):
		return writeGroupReadinessIssue("group-source-missing", "write group source is no longer available", scope), true
	default:
		return ReadinessIssue{}, false
	}
}

func readinessValidationIssueOrInternal(err error, scope string) ReadinessIssue {
	if issue, ok := writeGroupReadinessValidationIssue(err, scope); ok {
		return issue
	}
	return writeGroupReadinessIssue("readiness-validation-failed", "write group readiness could not be verified", scope)
}

func writeGroupReadinessIssue(code, message, scope string) ReadinessIssue {
	return ReadinessIssue{
		Code: code, Severity: ReadinessSeverityBlocking, Step: ReadinessStep4,
		Scope: strings.TrimSpace(scope), Message: message,
	}
}

func sameWriteGroupReadinessSnapshot(left, right *writeGroupReadinessSnapshot) bool {
	if left == nil || right == nil || left.record == nil || right.record == nil || left.group == nil || right.group == nil {
		return false
	}
	if left.record.ID != right.record.ID ||
		left.record.DatabaseSetupRevision != right.record.DatabaseSetupRevision ||
		left.group.ID != right.group.ID ||
		left.group.WorkspaceID != right.group.WorkspaceID ||
		left.group.Revision != right.group.Revision ||
		left.group.AppliedRevision != right.group.AppliedRevision ||
		left.connectorKind != right.connectorKind ||
		left.connectorConfig != right.connectorConfig {
		return false
	}
	return reflect.DeepEqual(left.group, right.group) &&
		reflect.DeepEqual(left.validated, right.validated) &&
		reflect.DeepEqual(left.tagTypes, right.tagTypes)
}

func loadWriteGroupTagTypes(
	ctx context.Context,
	tx *sql.Tx,
	query func(string) string,
	group *WriteGroup,
) (map[string]schema.DataType, error) {
	types := make(map[string]schema.DataType, len(group.Members))
	for _, member := range group.Members {
		var dataType string
		err := tx.QueryRowContext(ctx, query("SELECT data_type FROM tags WHERE id = $1"), member.TagID).Scan(&dataType)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
		}
		if err != nil {
			return nil, fmt.Errorf("read write-group tag data type: %w", err)
		}
		types[member.TagID] = schema.DataType(strings.TrimSpace(dataType))
	}
	return types, nil
}
