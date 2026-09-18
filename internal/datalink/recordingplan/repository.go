package recordingplan

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrPlanNotFound identifies an absent plan or a plan outside the requested
// workspace. Callers should keep both cases indistinguishable at the API.
var ErrPlanNotFound = errors.New("recording plan not found")

// Repository 定義 RecordingPlan 與 SchemaPreviewToken 的持久化介面。
type Repository interface {
	CreatePlan(ctx context.Context, plan *RecordingPlan) error
	UpdatePlan(ctx context.Context, plan *RecordingPlan) error
	GetPlanByID(ctx context.Context, id string) (*RecordingPlan, error)
	GetPlanByWorkspace(ctx context.Context, id, workspaceID string) (*RecordingPlan, error)
	ListPlansByWorkspace(ctx context.Context, workspaceID string) ([]RecordingPlan, error)
	DeletePlan(ctx context.Context, id string) error
	UpdatePlanByWorkspace(ctx context.Context, plan *RecordingPlan, workspaceID string) error
	DeletePlanByWorkspace(ctx context.Context, id, workspaceID string) error

	SavePreviewToken(ctx context.Context, token *SchemaPreviewToken) error
	GetPreviewToken(ctx context.Context, token string) (*SchemaPreviewToken, error)
	DeletePreviewToken(ctx context.Context, token string) error

	// ClaimSchemaOperation atomically records op as the owner of its scope, or
	// reports the existing operation for the same identity or scope.
	ClaimSchemaOperation(ctx context.Context, op *SchemaOperation) (*SchemaOperation, ClaimOutcome, error)
	GetSchemaOperation(ctx context.Context, workspaceID, operationID string) (*SchemaOperation, error)
	FindActiveSchemaOperation(ctx context.Context, workspaceID, scopeKey string) (*SchemaOperation, error)
	FinishSchemaOperation(ctx context.Context, operationID, owner string, result SchemaOperationResult) (*SchemaOperation, error)
	// FinishStaleSchemaOperation records a terminal result for an active
	// operation whose claim was last touched before staleBefore, regardless of
	// its owner. It returns nil when the operation is gone or no longer
	// stale-active, which means the owning execution finished it first.
	FinishStaleSchemaOperation(ctx context.Context, operationID string, result SchemaOperationResult, staleBefore time.Time) (*SchemaOperation, error)
}

// SQLRepository 實作基於 SQL 資料庫的 Repository。
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository 建立新的 SQLRepository。
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// CreatePlan 新增記錄方案。
func (r *SQLRepository) CreatePlan(ctx context.Context, plan *RecordingPlan) error {
	if plan.CreatedAt.IsZero() {
		plan.CreatedAt = time.Now().UTC()
	}
	if plan.UpdatedAt.IsZero() {
		plan.UpdatedAt = time.Now().UTC()
	}

	encoded, err := encodeRecordingPlan(plan)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO recording_plans (
			id, workspace_id, revision, applied_revision, name, status, timezone,
			members, streams, destinations, retention, limits, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	query = adaptPlaceholders(query)

	_, err = r.db.ExecContext(ctx, query,
		plan.ID, plan.WorkspaceID, plan.Revision, plan.AppliedRevision,
		plan.Name, string(plan.Status), plan.Timezone,
		string(encoded.members), string(encoded.streams), string(encoded.destinations),
		string(encoded.retention), string(encoded.limits), plan.CreatedAt, plan.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert recording_plan: %w", err)
	}
	return nil
}

// UpdatePlan 更新記錄方案。
func (r *SQLRepository) UpdatePlan(ctx context.Context, plan *RecordingPlan) error {
	plan.UpdatedAt = time.Now().UTC()

	encoded, err := encodeRecordingPlan(plan)
	if err != nil {
		return err
	}

	query := `
		UPDATE recording_plans SET
			workspace_id = $1, revision = $2, applied_revision = $3, name = $4,
			status = $5, timezone = $6, members = $7, streams = $8,
			destinations = $9, retention = $10, limits = $11, updated_at = $12
		WHERE id = $13
	`
	query = adaptPlaceholders(query)

	res, err := r.db.ExecContext(ctx, query,
		plan.WorkspaceID, plan.Revision, plan.AppliedRevision, plan.Name,
		string(plan.Status), plan.Timezone, string(encoded.members), string(encoded.streams),
		string(encoded.destinations), string(encoded.retention), string(encoded.limits), plan.UpdatedAt,
		plan.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update recording_plan: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: %s", ErrPlanNotFound, plan.ID)
	}
	return nil
}

// GetPlanByID 依 ID 取得記錄方案。
func (r *SQLRepository) GetPlanByID(ctx context.Context, id string) (*RecordingPlan, error) {
	query := `
		SELECT id, workspace_id, revision, applied_revision, name, status, timezone,
			members, streams, destinations, retention, limits, created_at, updated_at
		FROM recording_plans
		WHERE id = $1
	`
	query = adaptPlaceholders(query)

	row := r.db.QueryRowContext(ctx, query, id)
	return scanRecordingPlan(row)
}

// GetPlanByWorkspace returns a plan only when both its ID and workspace match.
func (r *SQLRepository) GetPlanByWorkspace(ctx context.Context, id, workspaceID string) (*RecordingPlan, error) {
	query := `
		SELECT id, workspace_id, revision, applied_revision, name, status, timezone,
			members, streams, destinations, retention, limits, created_at, updated_at
		FROM recording_plans
		WHERE id = $1 AND workspace_id = $2
	`
	query = adaptPlaceholders(query)

	row := r.db.QueryRowContext(ctx, query, id, workspaceID)
	plan, err := scanRecordingPlan(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", ErrPlanNotFound, id)
	}
	return plan, err
}

// ListPlansByWorkspace 依 Workspace 取得記錄方案列表。
func (r *SQLRepository) ListPlansByWorkspace(ctx context.Context, workspaceID string) ([]RecordingPlan, error) {
	query := `
		SELECT id, workspace_id, revision, applied_revision, name, status, timezone,
			members, streams, destinations, retention, limits, created_at, updated_at
		FROM recording_plans
		WHERE workspace_id = $1
		ORDER BY created_at ASC
	`
	query = adaptPlaceholders(query)

	rows, err := r.db.QueryContext(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []RecordingPlan
	for rows.Next() {
		p, err := scanRecordingPlanRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *p)
	}
	return list, rows.Err()
}

// DeletePlan 刪除記錄方案。
func (r *SQLRepository) DeletePlan(ctx context.Context, id string) error {
	query := `DELETE FROM recording_plans WHERE id = $1`
	query = adaptPlaceholders(query)
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: %s", ErrPlanNotFound, id)
	}
	return nil
}

// UpdatePlanByWorkspace updates a plan only within its current workspace.
func (r *SQLRepository) UpdatePlanByWorkspace(ctx context.Context, plan *RecordingPlan, workspaceID string) error {
	if strings.TrimSpace(plan.WorkspaceID) != "" && plan.WorkspaceID != workspaceID {
		return fmt.Errorf("%w: %s", ErrPlanNotFound, plan.ID)
	}
	plan.WorkspaceID = workspaceID
	plan.UpdatedAt = time.Now().UTC()

	encoded, err := encodeRecordingPlan(plan)
	if err != nil {
		return err
	}

	query := `
		UPDATE recording_plans SET
			revision = $1, applied_revision = $2, name = $3, status = $4,
			timezone = $5, members = $6, streams = $7, destinations = $8,
			retention = $9, limits = $10, updated_at = $11
		WHERE id = $12 AND workspace_id = $13
	`
	query = adaptPlaceholders(query)

	res, err := r.db.ExecContext(ctx, query,
		plan.Revision, plan.AppliedRevision, plan.Name, string(plan.Status), plan.Timezone,
		string(encoded.members), string(encoded.streams), string(encoded.destinations),
		string(encoded.retention), string(encoded.limits), plan.UpdatedAt,
		plan.ID, workspaceID,
	)
	if err != nil {
		return fmt.Errorf("failed to update recording_plan: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: %s", ErrPlanNotFound, plan.ID)
	}
	return nil
}

// DeletePlanByWorkspace deletes a plan only within the requested workspace.
func (r *SQLRepository) DeletePlanByWorkspace(ctx context.Context, id, workspaceID string) error {
	query := `DELETE FROM recording_plans WHERE id = $1 AND workspace_id = $2`
	query = adaptPlaceholders(query)
	res, err := r.db.ExecContext(ctx, query, id, workspaceID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("%w: %s", ErrPlanNotFound, id)
	}
	return nil
}

// SavePreviewToken 儲存預覽 Token。
func (r *SQLRepository) SavePreviewToken(ctx context.Context, token *SchemaPreviewToken) error {
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now().UTC()
	}
	stmtsJSON, err := json.Marshal(token.Statements)
	if err != nil {
		return fmt.Errorf("encode preview statements: %w", err)
	}
	tablesJSON, err := json.Marshal(token.Tables)
	if err != nil {
		return fmt.Errorf("encode preview tables: %w", err)
	}

	query := `
		INSERT INTO managed_schema_preview_tokens (
			token, workspace_id, plan_id, plan_revision, connector_id, table_prefix,
			statements, expires_at, created_at, operation_id, action, workspace_revision,
			connector_revision, dialect, database_name, schema_name, no_change_reason, digest, tables
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`
	query = adaptPlaceholders(query)

	_, err = r.db.ExecContext(ctx, query,
		token.Token, token.WorkspaceID, token.PlanID, token.PlanRevision,
		token.ConnectorID, token.TablePrefix, string(stmtsJSON), token.ExpiresAt, token.CreatedAt,
		token.OperationID, token.Action, token.WorkspaceRevision, token.ConnectorRevision, token.Dialect,
		token.Database, token.Schema, token.NoChangeReason, token.Digest, string(tablesJSON),
	)
	return err
}

// GetPreviewToken 取得預覽 Token。
func (r *SQLRepository) GetPreviewToken(ctx context.Context, token string) (*SchemaPreviewToken, error) {
	query := `
		SELECT token, workspace_id, plan_id, plan_revision, connector_id, table_prefix,
			statements, expires_at, created_at, operation_id, action, workspace_revision,
			connector_revision, dialect, database_name, schema_name, no_change_reason, digest, tables
		FROM managed_schema_preview_tokens
		WHERE token = $1
	`
	query = adaptPlaceholders(query)

	var t SchemaPreviewToken
	var stmtsStr, tablesStr string
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&t.Token, &t.WorkspaceID, &t.PlanID, &t.PlanRevision,
		&t.ConnectorID, &t.TablePrefix, &stmtsStr, &t.ExpiresAt, &t.CreatedAt,
		&t.OperationID, &t.Action, &t.WorkspaceRevision, &t.ConnectorRevision, &t.Dialect,
		&t.Database, &t.Schema, &t.NoChangeReason, &t.Digest, &tablesStr,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", ErrPreviewTokenNotFound, token)
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(stmtsStr), &t.Statements); err != nil {
		return nil, fmt.Errorf("decode preview statements: %w", err)
	}
	if err := json.Unmarshal([]byte(tablesStr), &t.Tables); err != nil {
		return nil, fmt.Errorf("decode preview tables: %w", err)
	}
	return &t, nil
}

// DeletePreviewToken 刪除預覽 Token。
func (r *SQLRepository) DeletePreviewToken(ctx context.Context, token string) error {
	query := `DELETE FROM managed_schema_preview_tokens WHERE token = $1`
	query = adaptPlaceholders(query)
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanRecordingPlan(s scannable) (*RecordingPlan, error) {
	var p RecordingPlan
	var statusStr string
	var membersStr, streamsStr, destsStr, retStr, limStr string

	err := s.Scan(
		&p.ID, &p.WorkspaceID, &p.Revision, &p.AppliedRevision, &p.Name, &statusStr, &p.Timezone,
		&membersStr, &streamsStr, &destsStr, &retStr, &limStr, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.Status = Status(statusStr)
	for _, field := range []struct {
		name, raw string
		target    any
	}{
		{"members", membersStr, &p.Members},
		{"streams", streamsStr, &p.Streams},
		{"destinations", destsStr, &p.Destinations},
		{"retention", retStr, &p.Retention},
		{"limits", limStr, &p.Limits},
	} {
		if err := json.Unmarshal([]byte(field.raw), field.target); err != nil {
			return nil, fmt.Errorf("decode recording plan %s: %w", field.name, err)
		}
	}

	return &p, nil
}

func scanRecordingPlanRow(rows *sql.Rows) (*RecordingPlan, error) {
	return scanRecordingPlan(rows)
}

func adaptPlaceholders(query string) string {
	for i := 25; i >= 1; i-- {
		query = strings.ReplaceAll(query, fmt.Sprintf("$%d", i), "?")
	}
	return query
}
