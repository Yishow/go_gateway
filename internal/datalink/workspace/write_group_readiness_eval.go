package workspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
)

type writeGroupReadinessSnapshot struct {
	record                   *Record
	group                    *WriteGroup
	validated                *WriteGroup
	tagTypes                 map[string]schema.DataType
	connectorKind            string
	connectorConfig          string
	validateErr              error
	schemaProofErr           error
	counterBindings          CounterBindingProofs
	counterIssues            []ReadinessIssue
	confirmedMySQLLayout     *recordingplan.GroupSchemaLayout
	confirmedSQLServerLayout *recordingplan.GroupSchemaLayout
}

func (s *WriteGroupService) readiness(ctx context.Context, id string) (*WriteGroupReadiness, error) {
	return s.readinessOf(ctx, id, nil)
}

// readinessOf evaluates a group's readiness. With content set, that content is
// evaluated instead of the saved row: the group as it would be after a
// rollback, against the same live sources, tags and destination.
func (s *WriteGroupService) readinessOf(ctx context.Context, id string, content *WriteGroup) (*WriteGroupReadiness, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, writeGroupNotFound("read write-group readiness")
	}
	snapshot, err := s.readinessSnapshot(ctx, id, content)
	if err != nil {
		return nil, err
	}
	if isSuccessorReservedError(snapshot.validateErr) {
		result := newWriteGroupReadiness(snapshot.record, snapshot.group)
		result.Issues = append(result.Issues, writeGroupReadinessIssue("successor_reserved", "reserved successor is unavailable until reviewed activation", id))
		return result, nil
	}
	return s.evaluateReadinessSnapshot(ctx, id, content, snapshot, func(ctx context.Context) (*writeGroupReadinessSnapshot, error) {
		return s.readinessSnapshot(ctx, id, content)
	})
}

// evaluateReadinessSnapshot shares schema evaluation, never ownership authority.
// Every caller must supply its own mandatory, guarded coherent reloader.
func (s *WriteGroupService) evaluateReadinessSnapshot(ctx context.Context, id string, content *WriteGroup, snapshot *writeGroupReadinessSnapshot, reload func(context.Context) (*writeGroupReadinessSnapshot, error)) (*WriteGroupReadiness, error) {
	if snapshot == nil || reload == nil {
		return nil, ErrWriteGroupReadinessUnavailable
	}
	result := s.newSnapshotReadiness(ctx, snapshot)
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
	if !result.ConfigReady {
		result.SchemaReady = false
		result.Issues = append(result.Issues, writeGroupReadinessIssue(
			"schema-check-blocked", "destination schema check is blocked by configuration issues", id,
		))
		return finalizeWriteGroupReadiness(result), nil
	}
	if snapshot.schemaProofErr != nil {
		result.SchemaReady = false
		result.Issues = append(result.Issues, writeGroupReadinessIssue("schema-proof-stale", "calculation evidence changed; preview and confirm the schema again", id))
		return finalizeWriteGroupReadiness(result), nil
	}
	managedCanonical := snapshot.validated.Destination.StorageStrategy == WriteGroupStorageStrategyManaged && snapshot.validated.RowPolicy.RecordKeyColumn != ""
	if managedCanonical {
		if s.managedTableInspector == nil {
			result.SchemaReady = false
			result.Issues = append(result.Issues, writeGroupReadinessIssue(
				"schema-unverified", "managed destination schema has not been verified", id,
			))
			return finalizeWriteGroupReadiness(result), nil
		}
	} else if s.tableInspector == nil {
		result.SchemaReady = false
		result.Issues = append(result.Issues, writeGroupReadinessIssue(
			"schema-unverified", "destination schema has not been verified", id,
		))
		return finalizeWriteGroupReadiness(result), nil
	}
	if snapshot.group.Destination.StorageStrategy == WriteGroupStorageStrategyManaged && snapshot.group.RowPolicy.RecordKeyColumn != "" &&
		(snapshot.group.Destination.SchemaRevision == "" || snapshot.group.Destination.SchemaDigest == "") {
		result.Issues = append(result.Issues, writeGroupReadinessIssue("schema-confirmation-required", "preview and explicitly confirm the saved managed table before applying", id))
		return finalizeWriteGroupReadiness(result), nil
	}

	if managedCanonical {
		allowed, issue := allowManagedWriteGroupSchemaScope(snapshot)
		if !allowed {
			result.SchemaReady = false
			result.Issues = append(result.Issues, issue)
			return finalizeWriteGroupReadiness(result), nil
		}
	}
	if !managedCanonical {
		allowed, issue := allowWriteGroupTableInspection(snapshot)
		if !allowed {
			result.SchemaReady = false
			result.Issues = append(result.Issues, issue)
			return finalizeWriteGroupReadiness(result), nil
		}
	}
	var inspection *dbtarget.TableInspection
	if managedCanonical && snapshot.connectorKind == writeGroupConnectorKindMySQL {
		return s.mysqlReadiness(ctx, id, snapshot, result, reload)
	}
	if managedCanonical && snapshot.connectorKind == writeGroupConnectorKindSQLServer {
		return s.sqlServerReadiness(ctx, id, snapshot, result, reload)
	}
	var inspectErr error
	if managedCanonical {
		inspection, inspectErr = s.managedTableInspector.InspectTableAtRevision(ctx,
			snapshot.validated.Destination.ConnectorID,
			snapshot.validated.Destination.ConnectorRevision,
			snapshot.validated.Destination.TableSchema,
			snapshot.validated.Destination.TableName,
		)
	} else {
		inspection, inspectErr = s.tableInspector.InspectTable(ctx,
			snapshot.validated.Destination.ConnectorID,
			snapshot.validated.Destination.TableSchema,
			snapshot.validated.Destination.TableName,
		)
	}
	if managedCanonical && inspectErr == nil && snapshot.validated.isLongLayout() &&
		inspection != nil && inspection.Status == dbtarget.TableInspectionExists {
		inspectErr = s.attachManagedIndexes(ctx, snapshot.validated, inspection)
		if inspectErr == nil && snapshot.validated.StorageMode == WriteGroupStorageModeTimescale {
			inspectErr = s.attachManagedHypertable(ctx, snapshot.validated, inspection)
		}
	}
	if managedCanonical && inspectErr == nil {
		allowed, issue := allowWriteGroupTableInspection(snapshot)
		if !allowed {
			result.SchemaReady = false
			result.Issues = append(result.Issues, issue)
			return finalizeWriteGroupReadiness(result), nil
		}
	}
	receiptReady := false
	if inspectErr == nil && strings.EqualFold(strings.TrimSpace(snapshot.validated.WritePolicy.DedupeCapability), "receipt") {
		var receipts *dbtarget.TableInspection
		var receiptErr error
		if managedCanonical {
			receipts, receiptErr = s.managedTableInspector.InspectTableAtRevision(ctx,
				snapshot.validated.Destination.ConnectorID,
				snapshot.validated.Destination.ConnectorRevision,
				snapshot.validated.Destination.TableSchema, dbtarget.EffectReceiptTable)
		} else {
			receipts, receiptErr = s.tableInspector.InspectTable(ctx, snapshot.validated.Destination.ConnectorID,
				snapshot.validated.Destination.TableSchema, dbtarget.EffectReceiptTable)
		}
		receiptReady = receiptErr == nil && receipts != nil && receipts.Status == dbtarget.TableInspectionExists
		if managedCanonical && receiptReady {
			receiptReady = recordingplan.ValidateManagedReceiptTableInspection(
				snapshot.connectorKind, recordingPlanTableInspection(receipts),
			) == nil
		}
	}
	postSnapshot, err := reload(ctx)
	if err != nil {
		return nil, err
	}
	if postSnapshot != nil && postSnapshot.schemaProofErr != nil {
		result.SchemaReady = false
		result.Issues = append(result.Issues, writeGroupReadinessIssue("schema-proof-stale", "calculation evidence changed; preview and confirm the schema again", id))
		return finalizeWriteGroupReadiness(result), nil
	}
	if !sameWriteGroupReadinessSnapshot(snapshot, postSnapshot) {
		result = s.newSnapshotReadiness(ctx, postSnapshot)
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
		postSnapshot.validated, postSnapshot.tagTypes, inspection, postSnapshot.connectorKind,
	)
	if managedCanonical && strings.EqualFold(strings.TrimSpace(postSnapshot.validated.WritePolicy.DedupeCapability), "receipt") && !receiptReady {
		schemaReady = false
		schemaIssues = append(schemaIssues, writeGroupReadinessIssue(
			"managed-receipt-unverified", "managed receipt ownership and canonical schema could not be verified", id,
		))
	}
	partialIssues := evaluateWriteGroupPartialInspection(postSnapshot.validated, postSnapshot.tagTypes, inspection, postSnapshot.connectorKind)
	schemaIssues = append(schemaIssues, partialIssues...)
	schemaReady = schemaReady && len(partialIssues) == 0
	result.SchemaReady = schemaReady
	result.Issues = append(result.Issues, schemaIssues...)
	if schemaReady {
		result.runtimeLayout = &WriteGroupRuntimeLayout{
			Dialect: dbtarget.SQLDialect(postSnapshot.connectorKind), Columns: inspection.Columns,
			TagTypes: postSnapshot.tagTypes, ReceiptTableReady: receiptReady,
		}
		result.SchemaDigest, err = digestWriteGroupInspection(postSnapshot.validated, inspection, postSnapshot.counterBindings)
		if err != nil {
			return nil, fmt.Errorf("digest write-group schema readiness: %w", err)
		}
		result.runtimeLayout.SchemaDigest = result.SchemaDigest
		if postSnapshot.validated.StorageMode == WriteGroupStorageModeTimescale {
			result.runtimeLayout.Hypertable = freezeHypertableEvidence(inspection.Hypertable)
		}
	}
	return finalizeWriteGroupReadiness(result), nil
}

func (s *WriteGroupService) readinessSnapshot(ctx context.Context, id string, content *WriteGroup) (result *writeGroupReadinessSnapshot, err error) {
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
	if guardErr := s.repo.requireUnreservedWriteGroupInTx(ctx, setupTx.SQLTx(), record.ID, id); guardErr != nil {
		if !isSuccessorReservedError(guardErr) {
			return nil, guardErr
		}
		return &writeGroupReadinessSnapshot{record: record, group: group, validated: cloneWriteGroup(group), validateErr: guardErr}, nil
	}
	if content != nil {
		group = cloneWriteGroup(content)
	}
	validated := cloneWriteGroup(group)
	definitions := make(map[string]*measurement.MeasurementDefinition)
	validationErr := s.repo.validateWithMeasurements(ctx, setupTx.SQLTx(), validated, definitions)
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
	proofs, issues := s.counterBindingEvidenceInTx(ctx, setupTx.SQLTx(), group, tagTypes, definitions)
	var schemaProofErr error
	if isMixedWriteGroup(validated) && validationErr == nil && validated.Destination.SchemaRevision != "" {
		schemaProofErr = s.verifyMixedSchemaProofInTx(ctx, setupTx.SQLTx(), record, validated, validated.Destination.SchemaRevision, validated.Destination.SchemaDigest, nil)
	}
	return &writeGroupReadinessSnapshot{
		record: record, group: group, validated: validated,
		tagTypes: tagTypes, connectorKind: connectorKind, connectorConfig: connectorConfig,
		validateErr: validationErr, schemaProofErr: schemaProofErr,
		counterBindings: proofs, counterIssues: issues,
	}, nil
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
	if !isMixedWriteGroup(group) && group.RowPolicy.IntervalSeconds <= 0 {
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
	if incompletePolicy != "" && incompletePolicy != writeGroupIncompletePolicySkipRow && incompletePolicy != "partial" {
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
	case isSuccessorReservedError(err):
		return writeGroupReadinessIssue("successor_reserved", "reserved successor is unavailable until reviewed activation", scope), true
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
	if left.schemaProofErr != nil || right.schemaProofErr != nil {
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
		reflect.DeepEqual(left.tagTypes, right.tagTypes) &&
		reflect.DeepEqual(left.counterBindings, right.counterBindings) &&
		reflect.DeepEqual(left.counterIssues, right.counterIssues) &&
		reflect.DeepEqual(left.confirmedMySQLLayout, right.confirmedMySQLLayout) &&
		reflect.DeepEqual(left.confirmedSQLServerLayout, right.confirmedSQLServerLayout)
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
