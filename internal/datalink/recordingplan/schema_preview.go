package recordingplan

import (
	"cmp"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"go-gateway/internal/datalink/common"
)

var (
	ErrPreviewTokenNotFound        = errors.New("schema preview token not found")
	ErrPreviewTokenExpired         = errors.New("schema preview token expired")
	ErrPreviewTokenStale           = errors.New("schema preview token is stale")
	ErrPreviewTokenLegacy          = errors.New("schema preview token lacks required protection")
	ErrInvalidTablePrefix          = errors.New("managed table prefix is invalid")
	ErrIncompatibleExistingTable   = errors.New("existing managed table is incompatible")
	ErrTargetInspectionUnconfirmed = errors.New("target schema could not be confirmed")
	ErrTargetPermissionDenied      = errors.New("target schema permission denied")
)

const (
	// SchemaApplyAction marks a token that may only confirm schema creation.
	SchemaApplyAction = "schema_apply"
	// NoChangeSchemaCompatible explains an empty preview for compatible tables.
	NoChangeSchemaCompatible = "schema_already_compatible"

	// defaultManagedTablePrefix is the value the frontend sends as table_prefix;
	// kept in sync both ways with MANAGED_TABLE_PREFIX
	// (frontend/src/features/datalink/workbench-v2/steps/step4/useRecordingPlanActions.ts).
	defaultManagedTablePrefix = "gw_record_"
	schemaPreviewTTL          = 10 * time.Minute
	// maxManagedTablePrefixLength keeps the longest generated identifier
	// (idx_<prefix>samples_meas_epoch) within PostgreSQL's 63-byte limit.
	maxManagedTablePrefixLength = 40

	inspectionExists    = "exists"
	inspectionMissing   = "missing"
	inspectionForbidden = "forbidden"
)

var managedTablePrefixPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// SchemaPreviewScope is the server-resolved scope a preview is bound to.
type SchemaPreviewScope struct {
	WorkspaceID       string
	WorkspaceRevision string
	PlanID            string
	PlanRevision      string
	ConnectorID       string
	ConnectorRevision string
	Dialect           string
	Database          string
	Schema            string
	TablePrefix       string
	// GroupLayout selects canonical WriteGroup managed preparation. A nil
	// layout retains the legacy RecordingPlan schema contract.
	GroupLayout    *GroupSchemaLayout
	SourceDigest   string
	SchemaRevision string
	SchemaDigest   string
}

// TargetTableInspection is the actual state of one managed table.
type TargetTableInspection struct {
	Status           string
	Columns          []string
	ColumnTypes      map[string]string
	ColumnNullable   map[string]bool
	ColumnPrimaryKey map[string]bool
}

// TargetInspector reads one managed table from the saved target.
type TargetInspector func(ctx context.Context, table string) (TargetTableInspection, error)

// PrepareSchemaPreview inspects the target and persists a revision-bound
// preview. It never writes to the target database.
func (s *Service) PrepareSchemaPreview(ctx context.Context, scope SchemaPreviewScope, inspect TargetInspector) (*SchemaPreviewToken, error) {
	scope, err := normalizePreviewScope(scope)
	if err != nil {
		return nil, err
	}
	if inspect == nil {
		return nil, fmt.Errorf("%w: no target inspector", ErrTargetInspectionUnconfirmed)
	}
	var statements []string
	var tables []SchemaPreviewTable
	if scope.GroupLayout != nil {
		statements, tables, err = prepareGroupSchemaPreview(ctx, scope, inspect)
	} else {
		statements, err = GenerateManagedSchemaDDL(scope.Dialect, scope.TablePrefix)
		if err == nil {
			pending, legacyTables, inspectErr := inspectManagedTables(ctx, scope.TablePrefix, inspect)
			if inspectErr != nil {
				return nil, inspectErr
			}
			statements = statementsForTables(statements, scope.TablePrefix, pending)
			tables = legacyTables
		}
	}
	if err != nil {
		return nil, err
	}

	tokenID, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("create preview token: %w", err)
	}
	operationID, err := common.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("create preview operation: %w", err)
	}
	now := time.Now().UTC()
	token := &SchemaPreviewToken{
		Token: "tok-" + tokenID, OperationID: "op-" + operationID, Action: SchemaApplyAction,
		WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision,
		PlanID: scope.PlanID, PlanRevision: scope.PlanRevision,
		ConnectorID: scope.ConnectorID, ConnectorRevision: scope.ConnectorRevision,
		Dialect: scope.Dialect, Database: scope.Database, Schema: scope.Schema, TablePrefix: scope.TablePrefix,
		GroupLayout: scope.GroupLayout, SourceDigest: scope.SourceDigest, SchemaRevision: scope.SchemaRevision, SchemaDigest: scope.SchemaDigest,
		Statements: statements, Tables: tables, ExpiresAt: now.Add(schemaPreviewTTL), CreatedAt: now,
	}
	if len(statements) == 0 {
		token.NoChangeReason = NoChangeSchemaCompatible
	}
	if token.Digest, err = previewDigest(token); err != nil {
		return nil, err
	}
	if err := s.repo.SavePreviewToken(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to save preview token: %w", err)
	}
	return token, nil
}

// ValidatePreviewTokenForApply re-checks a stored preview against the current
// server-resolved scope. Tokens from another workspace look unknown.
func (s *Service) ValidatePreviewTokenForApply(ctx context.Context, token string, current SchemaPreviewScope) (*SchemaPreviewToken, error) {
	current, err := normalizePreviewScope(current)
	if err != nil {
		return nil, err
	}
	stored, err := s.repo.GetPreviewToken(ctx, strings.TrimSpace(token))
	if err != nil {
		return nil, err
	}
	if stored.WorkspaceID != current.WorkspaceID {
		return nil, fmt.Errorf("%w: %s", ErrPreviewTokenNotFound, token)
	}
	if isLegacyPreviewToken(stored) {
		return nil, ErrPreviewTokenLegacy
	}
	digest, err := previewDigest(stored)
	if err != nil {
		return nil, err
	}
	if subtle.ConstantTimeCompare([]byte(digest), []byte(stored.Digest)) != 1 {
		return nil, fmt.Errorf("%w: preview content changed", ErrPreviewTokenStale)
	}
	if field := previewScopeMismatch(stored, current); field != "" {
		return nil, fmt.Errorf("%w: %s changed", ErrPreviewTokenStale, field)
	}
	if stored.IsExpired() {
		// An expired token cannot authorize a new DDL claim. It remains
		// readable only when the same durable operation already exists, so a
		// restart can reconcile an active/terminal operation without issuing a
		// second operation or executing a fresh batch.
		op, opErr := s.repo.GetSchemaOperation(ctx, current.WorkspaceID, stored.OperationID)
		if opErr != nil || op.Token != stored.Token || op.Action != stored.Action {
			return nil, ErrPreviewTokenExpired
		}
	}
	return stored, nil
}

func normalizePreviewScope(scope SchemaPreviewScope) (SchemaPreviewScope, error) {
	scope.TablePrefix = strings.TrimSpace(scope.TablePrefix)
	scope.Dialect = normalizePreviewDialect(scope.Dialect)
	if scope.GroupLayout != nil {
		layout, err := ValidateGroupSchemaLayout(*scope.GroupLayout)
		if err != nil {
			return scope, err
		}
		scope.GroupLayout = &layout
		// A group target owns its exact table name; the legacy prefix must not
		// accidentally influence its statements or operation scope.
		scope.TablePrefix = ""
		return scope, nil
	}
	scope.TablePrefix = cmp.Or(scope.TablePrefix, defaultManagedTablePrefix)
	if len(scope.TablePrefix) > maxManagedTablePrefixLength || !managedTablePrefixPattern.MatchString(scope.TablePrefix) {
		return scope, ErrInvalidTablePrefix
	}
	return scope, nil
}

func normalizePreviewDialect(dialect string) string {
	switch d := strings.ToLower(strings.TrimSpace(dialect)); d {
	case dialectSQLite3:
		return dialectSQLite
	case dialectPostgreSQL, dialectPgx:
		return dialectPostgres
	default:
		return d
	}
}

// inspectManagedTables reports which managed tables must be created and
// summarizes every table. Any state that cannot be confirmed as missing or
// compatible blocks the preview.
func inspectManagedTables(ctx context.Context, prefix string, inspect TargetInspector) (map[string]bool, []SchemaPreviewTable, error) {
	pending := make(map[string]bool, len(managedTableSpecs))
	tables := make([]SchemaPreviewTable, 0, len(managedTableSpecs))
	for _, spec := range managedTableSpecs {
		table := prefix + spec.name
		state, err := inspect(ctx, table)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %s: %w", ErrTargetInspectionUnconfirmed, table, err)
		}
		switch state.Status {
		case inspectionMissing:
			pending[spec.name] = true
			tables = append(tables, SchemaPreviewTable{Name: table, Action: SchemaTableActionCreate, Columns: slices.Clone(spec.columns)})
		case inspectionExists:
			if missing := missingColumns(spec.columns, state.Columns); len(missing) > 0 {
				return nil, nil, fmt.Errorf("%w: %s lacks %s", ErrIncompatibleExistingTable, table, strings.Join(missing, ", "))
			}
			if badCol, badType := incompatibleColumnType(spec.columnCategories, state.ColumnTypes); badCol != "" {
				return nil, nil, fmt.Errorf("%w: %s column %s has incompatible type %s", ErrIncompatibleExistingTable, table, badCol, badType)
			}
			tables = append(tables, SchemaPreviewTable{Name: table, Action: SchemaTableActionUnchanged})
		case inspectionForbidden:
			return nil, nil, fmt.Errorf("%w: %w: %s", ErrTargetInspectionUnconfirmed, ErrTargetPermissionDenied, table)
		default:
			return nil, nil, fmt.Errorf("%w: %s is %s", ErrTargetInspectionUnconfirmed, table, state.Status)
		}
	}
	return pending, tables, nil
}

func missingColumns(required, actual []string) []string {
	present := make(map[string]bool, len(actual))
	for _, column := range actual {
		present[strings.ToLower(strings.TrimSpace(column))] = true
	}
	var missing []string
	for _, column := range required {
		if !present[column] {
			missing = append(missing, column)
		}
	}
	return missing
}

// isLegacyPreviewToken reports tokens stored before previews carried their
// operation, revisions and digest; those must be previewed again.
func isLegacyPreviewToken(token *SchemaPreviewToken) bool {
	for _, value := range []string{token.OperationID, token.WorkspaceRevision, token.PlanRevision, token.ConnectorID,
		token.ConnectorRevision, token.Dialect, token.Digest} {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	if token.GroupLayout == nil && strings.TrimSpace(token.TablePrefix) == "" {
		return true
	}
	if token.GroupLayout != nil {
		if _, err := ValidateGroupSchemaLayout(*token.GroupLayout); err != nil {
			return true
		}
	}
	return token.Action != SchemaApplyAction
}

func previewScopeMismatch(token *SchemaPreviewToken, current SchemaPreviewScope) string {
	checks := []struct{ field, stored, current string }{
		{"workspace revision", token.WorkspaceRevision, current.WorkspaceRevision},
		{"plan", token.PlanID, current.PlanID},
		{"plan revision", token.PlanRevision, current.PlanRevision},
		{"connector", token.ConnectorID, current.ConnectorID},
		{"connector revision", token.ConnectorRevision, current.ConnectorRevision},
		{"dialect", token.Dialect, current.Dialect},
		{"database", token.Database, current.Database},
		{"schema", token.Schema, current.Schema},
		{"table prefix", token.TablePrefix, current.TablePrefix},
	}
	for _, check := range checks {
		if check.stored != check.current {
			return check.field
		}
	}
	if (token.GroupLayout == nil) != (current.GroupLayout == nil) {
		return "group layout"
	}
	if token.GroupLayout != nil && !groupSchemaLayoutsEqual(*token.GroupLayout, *current.GroupLayout) {
		return "group layout"
	}
	groupChecks := []struct{ field, stored, current string }{
		{"source digest", token.SourceDigest, current.SourceDigest},
		{"schema revision", token.SchemaRevision, current.SchemaRevision},
		{"schema digest", token.SchemaDigest, current.SchemaDigest},
	}
	for _, check := range groupChecks {
		if check.stored != check.current {
			return check.field
		}
	}
	return ""
}

func groupSchemaLayoutsEqual(left, right GroupSchemaLayout) bool {
	left, leftErr := ValidateGroupSchemaLayout(left)
	right, rightErr := ValidateGroupSchemaLayout(right)
	if leftErr != nil || rightErr != nil {
		return false
	}
	if left.TableName != right.TableName || left.OwnerColumn != right.OwnerColumn || len(left.Columns) != len(right.Columns) {
		return false
	}
	for i := range left.Columns {
		if left.Columns[i] != right.Columns[i] {
			return false
		}
	}
	return true
}

// previewDigest binds the confirmed scope, operation, statements and table
// summary together so any stored change is detected before apply.
func previewDigest(token *SchemaPreviewToken) (string, error) {
	parts := []any{
		token.Action, token.OperationID, token.WorkspaceID, token.WorkspaceRevision, token.PlanID, token.PlanRevision,
		token.ConnectorID, token.ConnectorRevision, token.Dialect, token.Database, token.Schema, token.TablePrefix,
		token.Statements, token.Tables, token.NoChangeReason,
	}
	// Preserve the exact legacy digest bytes when no group fields exist, so
	// stored pre-group tokens remain valid under the old contract. Group fields
	// are appended as one extension only when they were actually issued.
	if token.GroupLayout != nil || token.SourceDigest != "" || token.SchemaRevision != "" || token.SchemaDigest != "" {
		parts = append(parts, token.GroupLayout, token.SourceDigest, token.SchemaRevision, token.SchemaDigest)
	}
	payload, err := json.Marshal(parts)
	if err != nil {
		return "", fmt.Errorf("encode preview digest: %w", err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
