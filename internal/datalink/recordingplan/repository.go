package recordingplan

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Repository 定義 RecordingPlan 與 SchemaPreviewToken 的持久化介面。
type Repository interface {
	CreatePlan(ctx context.Context, plan *RecordingPlan) error
	UpdatePlan(ctx context.Context, plan *RecordingPlan) error
	GetPlanByID(ctx context.Context, id string) (*RecordingPlan, error)
	ListPlansByWorkspace(ctx context.Context, workspaceID string) ([]RecordingPlan, error)
	DeletePlan(ctx context.Context, id string) error

	SavePreviewToken(ctx context.Context, token *SchemaPreviewToken) error
	GetPreviewToken(ctx context.Context, token string) (*SchemaPreviewToken, error)
	DeletePreviewToken(ctx context.Context, token string) error
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

	membersJSON, _ := json.Marshal(plan.Members)
	streamsJSON, _ := json.Marshal(plan.Streams)
	destinationsJSON, _ := json.Marshal(plan.Destinations)
	retentionJSON, _ := json.Marshal(plan.Retention)
	limitsJSON, _ := json.Marshal(plan.Limits)

	query := `
		INSERT INTO recording_plans (
			id, workspace_id, revision, applied_revision, name, status, timezone,
			members, streams, destinations, retention, limits, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	query = adaptPlaceholders(query)

	_, err := r.db.ExecContext(ctx, query,
		plan.ID, plan.WorkspaceID, plan.Revision, plan.AppliedRevision,
		plan.Name, string(plan.Status), plan.Timezone,
		string(membersJSON), string(streamsJSON), string(destinationsJSON),
		string(retentionJSON), string(limitsJSON), plan.CreatedAt, plan.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert recording_plan: %w", err)
	}
	return nil
}

// UpdatePlan 更新記錄方案。
func (r *SQLRepository) UpdatePlan(ctx context.Context, plan *RecordingPlan) error {
	plan.UpdatedAt = time.Now().UTC()

	membersJSON, _ := json.Marshal(plan.Members)
	streamsJSON, _ := json.Marshal(plan.Streams)
	destinationsJSON, _ := json.Marshal(plan.Destinations)
	retentionJSON, _ := json.Marshal(plan.Retention)
	limitsJSON, _ := json.Marshal(plan.Limits)

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
		string(plan.Status), plan.Timezone, string(membersJSON), string(streamsJSON),
		string(destinationsJSON), string(retentionJSON), string(limitsJSON), plan.UpdatedAt,
		plan.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update recording_plan: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("recording_plan not found: %s", plan.ID)
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
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// SavePreviewToken 儲存預覽 Token。
func (r *SQLRepository) SavePreviewToken(ctx context.Context, token *SchemaPreviewToken) error {
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now().UTC()
	}
	stmtsJSON, _ := json.Marshal(token.Statements)

	query := `
		INSERT INTO managed_schema_preview_tokens (
			token, workspace_id, plan_id, plan_revision, connector_id, table_prefix,
			statements, expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	query = adaptPlaceholders(query)

	_, err := r.db.ExecContext(ctx, query,
		token.Token, token.WorkspaceID, token.PlanID, token.PlanRevision,
		token.ConnectorID, token.TablePrefix, string(stmtsJSON), token.ExpiresAt, token.CreatedAt,
	)
	return err
}

// GetPreviewToken 取得預覽 Token。
func (r *SQLRepository) GetPreviewToken(ctx context.Context, token string) (*SchemaPreviewToken, error) {
	query := `
		SELECT token, workspace_id, plan_id, plan_revision, connector_id, table_prefix,
			statements, expires_at, created_at
		FROM managed_schema_preview_tokens
		WHERE token = $1
	`
	query = adaptPlaceholders(query)

	var t SchemaPreviewToken
	var stmtsStr string
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&t.Token, &t.WorkspaceID, &t.PlanID, &t.PlanRevision,
		&t.ConnectorID, &t.TablePrefix, &stmtsStr, &t.ExpiresAt, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal([]byte(stmtsStr), &t.Statements)
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
	_ = json.Unmarshal([]byte(membersStr), &p.Members)
	_ = json.Unmarshal([]byte(streamsStr), &p.Streams)
	_ = json.Unmarshal([]byte(destsStr), &p.Destinations)
	_ = json.Unmarshal([]byte(retStr), &p.Retention)
	_ = json.Unmarshal([]byte(limStr), &p.Limits)

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
