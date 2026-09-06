package measurement

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Repository 定義量測定義的持久化介面。
type Repository interface {
	Create(ctx context.Context, def *MeasurementDefinition) error
	Update(ctx context.Context, def *MeasurementDefinition) error
	GetByID(ctx context.Context, id string) (*MeasurementDefinition, error)
	ListByWorkspace(ctx context.Context, workspaceID string) ([]MeasurementDefinition, error)
	ListByDevice(ctx context.Context, workspaceID, deviceID string) ([]MeasurementDefinition, error)
	Delete(ctx context.Context, id string) error
}

// SQLRepository 實作基於 SQL 的 MeasurementDefinition 儲存。
type SQLRepository struct {
	db *sql.DB
}

// NewSQLRepository 建立新的 SQLRepository。
func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// Create 新增量測定義。
func (r *SQLRepository) Create(ctx context.Context, def *MeasurementDefinition) error {
	if def.CreatedAt.IsZero() {
		def.CreatedAt = time.Now().UTC()
	}
	if def.UpdatedAt.IsZero() {
		def.UpdatedAt = time.Now().UTC()
	}

	counterPolicyJSON, _ := json.Marshal(def.CounterPolicy)
	stateMapJSON, _ := json.Marshal(def.StateMap)
	bitmaskLabelsJSON, _ := json.Marshal(def.BitmaskLabels)

	query := `
		INSERT INTO measurement_definitions (
			id, workspace_id, device_id, point_id, tag_id, equipment_id,
			definition_revision, source_binding_revision, series_epoch,
			name, quantity, unit, semantic_kind, numeric_encoding,
			counter_policy, state_map, bitmask_labels, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`
	query = adaptPlaceholders(query)

	_, err := r.db.ExecContext(ctx, query,
		def.ID, def.WorkspaceID, def.DeviceID, def.PointID, def.TagID, def.EquipmentID,
		def.DefinitionRevision, def.SourceBindingRevision, def.SeriesEpoch,
		def.Name, def.Quantity, def.Unit, string(def.SemanticKind), def.NumericEncoding,
		string(counterPolicyJSON), string(stateMapJSON), string(bitmaskLabelsJSON),
		def.CreatedAt, def.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert measurement_definition: %w", err)
	}
	return nil
}

// Update 更新量測定義。
func (r *SQLRepository) Update(ctx context.Context, def *MeasurementDefinition) error {
	def.UpdatedAt = time.Now().UTC()

	counterPolicyJSON, _ := json.Marshal(def.CounterPolicy)
	stateMapJSON, _ := json.Marshal(def.StateMap)
	bitmaskLabelsJSON, _ := json.Marshal(def.BitmaskLabels)

	query := `
		UPDATE measurement_definitions SET
			workspace_id = $1, device_id = $2, point_id = $3, tag_id = $4, equipment_id = $5,
			definition_revision = $6, source_binding_revision = $7, series_epoch = $8,
			name = $9, quantity = $10, unit = $11, semantic_kind = $12, numeric_encoding = $13,
			counter_policy = $14, state_map = $15, bitmask_labels = $16, updated_at = $17
		WHERE id = $18
	`
	query = adaptPlaceholders(query)

	res, err := r.db.ExecContext(ctx, query,
		def.WorkspaceID, def.DeviceID, def.PointID, def.TagID, def.EquipmentID,
		def.DefinitionRevision, def.SourceBindingRevision, def.SeriesEpoch,
		def.Name, def.Quantity, def.Unit, string(def.SemanticKind), def.NumericEncoding,
		string(counterPolicyJSON), string(stateMapJSON), string(bitmaskLabelsJSON), def.UpdatedAt,
		def.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update measurement_definition: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("measurement_definition not found: %s", def.ID)
	}
	return nil
}

// GetByID 依 ID 取得量測定義。
func (r *SQLRepository) GetByID(ctx context.Context, id string) (*MeasurementDefinition, error) {
	query := `
		SELECT id, workspace_id, device_id, point_id, tag_id, equipment_id,
			definition_revision, source_binding_revision, series_epoch,
			name, quantity, unit, semantic_kind, numeric_encoding,
			counter_policy, state_map, bitmask_labels, created_at, updated_at
		FROM measurement_definitions
		WHERE id = $1
	`
	query = adaptPlaceholders(query)

	row := r.db.QueryRowContext(ctx, query, id)
	return scanMeasurementDefinition(row)
}

// ListByWorkspace 依 Workspace 取得量測定義列表。
func (r *SQLRepository) ListByWorkspace(ctx context.Context, workspaceID string) ([]MeasurementDefinition, error) {
	query := `
		SELECT id, workspace_id, device_id, point_id, tag_id, equipment_id,
			definition_revision, source_binding_revision, series_epoch,
			name, quantity, unit, semantic_kind, numeric_encoding,
			counter_policy, state_map, bitmask_labels, created_at, updated_at
		FROM measurement_definitions
		WHERE workspace_id = $1
		ORDER BY created_at ASC
	`
	query = adaptPlaceholders(query)

	rows, err := r.db.QueryContext(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []MeasurementDefinition
	for rows.Next() {
		def, err := scanMeasurementDefinitionRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *def)
	}
	return list, rows.Err()
}

// ListByDevice 依 Device 取得量測定義列表。
func (r *SQLRepository) ListByDevice(ctx context.Context, workspaceID, deviceID string) ([]MeasurementDefinition, error) {
	query := `
		SELECT id, workspace_id, device_id, point_id, tag_id, equipment_id,
			definition_revision, source_binding_revision, series_epoch,
			name, quantity, unit, semantic_kind, numeric_encoding,
			counter_policy, state_map, bitmask_labels, created_at, updated_at
		FROM measurement_definitions
		WHERE workspace_id = $1 AND device_id = $2
		ORDER BY created_at ASC
	`
	query = adaptPlaceholders(query)

	rows, err := r.db.QueryContext(ctx, query, workspaceID, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []MeasurementDefinition
	for rows.Next() {
		def, err := scanMeasurementDefinitionRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *def)
	}
	return list, rows.Err()
}

// Delete 刪除量測定義。
func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM measurement_definitions WHERE id = $1`
	query = adaptPlaceholders(query)
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanMeasurementDefinition(s scannable) (*MeasurementDefinition, error) {
	var def MeasurementDefinition
	var tagID sql.NullString
	var unit sql.NullString
	var numericEncoding sql.NullString
	var counterPolicyStr sql.NullString
	var stateMapStr sql.NullString
	var bitmaskLabelsStr sql.NullString
	var semanticKindStr string

	err := s.Scan(
		&def.ID, &def.WorkspaceID, &def.DeviceID, &def.PointID, &tagID, &def.EquipmentID,
		&def.DefinitionRevision, &def.SourceBindingRevision, &def.SeriesEpoch,
		&def.Name, &def.Quantity, &unit, &semanticKindStr, &numericEncoding,
		&counterPolicyStr, &stateMapStr, &bitmaskLabelsStr, &def.CreatedAt, &def.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if tagID.Valid {
		def.TagID = &tagID.String
	}
	if unit.Valid {
		def.Unit = unit.String
	}
	if numericEncoding.Valid {
		def.NumericEncoding = numericEncoding.String
	}
	def.SemanticKind = SemanticKind(semanticKindStr)

	if counterPolicyStr.Valid && counterPolicyStr.String != "" {
		_ = json.Unmarshal([]byte(counterPolicyStr.String), &def.CounterPolicy)
	}
	if stateMapStr.Valid && stateMapStr.String != "" {
		_ = json.Unmarshal([]byte(stateMapStr.String), &def.StateMap)
	}
	if bitmaskLabelsStr.Valid && bitmaskLabelsStr.String != "" {
		_ = json.Unmarshal([]byte(bitmaskLabelsStr.String), &def.BitmaskLabels)
	}

	return &def, nil
}

func scanMeasurementDefinitionRow(rows *sql.Rows) (*MeasurementDefinition, error) {
	return scanMeasurementDefinition(rows)
}

func adaptPlaceholders(query string) string {
	// Simple SQLite vs PostgreSQL adaptation: replace $1, $2 with ? if SQLite
	for i := 25; i >= 1; i-- {
		query = strings.ReplaceAll(query, fmt.Sprintf("$%d", i), "?")
	}
	return query
}
