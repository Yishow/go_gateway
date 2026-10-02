package workspace

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrWriteGroupNotFound deliberately does not reveal whether the group is
	// absent or outside the requested workspace.
	ErrWriteGroupNotFound                  = errors.New("write group not found")
	ErrWriteGroupValidation                = errors.New("write-group validation failed")
	ErrWriteGroupConnectorNotFound         = errors.New("write-group connector not found")
	ErrWriteGroupConnectorRevisionConflict = errors.New("write-group connector identity revision conflict")
	ErrWriteGroupSourceRevisionConflict    = errors.New("write-group source revision conflict")
	ErrWriteGroupUnsupportedConnectorKind  = errors.New("write-group connector kind is unsupported")
)

type writeGroupSQLRunner interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// SQLWriteGroupRepository persists canonical WriteGroup records.
type SQLWriteGroupRepository struct {
	db       *sql.DB
	postgres bool
	now      func() time.Time
	newID    func() string
}

// NewSQLWriteGroupRepository creates a repository backed by db.
func NewSQLWriteGroupRepository(db *sql.DB) *SQLWriteGroupRepository {
	driverName := ""
	if db != nil {
		driverType := reflect.TypeOf(db.Driver())
		if driverType.Kind() == reflect.Pointer {
			driverType = driverType.Elem()
		}
		driverName = strings.ToLower(driverType.PkgPath() + "." + driverType.Name())
	}
	return &SQLWriteGroupRepository{
		db:       db,
		postgres: strings.Contains(driverName, "pq") || strings.Contains(driverName, "pgx"),
		now:      func() time.Time { return time.Now().UTC() },
		newID:    uuid.NewString,
	}
}

// Create validates and persists a new draft group in one local transaction.
func (r *SQLWriteGroupRepository) Create(ctx context.Context, group *WriteGroup) (err error) {
	if r == nil || r.db == nil {
		return fmt.Errorf("create write group: %w", ErrWriteGroupValidation)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin write-group transaction: %w", err)
	}
	committed := false
	defer func() {
		if committed {
			return
		}
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, fmt.Errorf("rollback write-group transaction: %w", rollbackErr))
		}
	}()
	if err := r.CreateInTx(ctx, tx, group); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit write-group transaction: %w", err)
	}
	committed = true
	return nil
}

// Get returns one group scoped to workspaceID.
func (r *SQLWriteGroupRepository) Get(ctx context.Context, workspaceID, groupID string) (*WriteGroup, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("get write group: %w", ErrWriteGroupNotFound)
	}
	return r.get(ctx, r.db, workspaceID, groupID)
}

// List returns all groups in a workspace in stable creation order.
func (r *SQLWriteGroupRepository) List(ctx context.Context, workspaceID string) ([]*WriteGroup, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("list write groups: %w", ErrWriteGroupNotFound)
	}
	return r.list(ctx, r.db, workspaceID)
}

// CreateInTx persists a draft using the caller's transaction. It never begins,
// commits, or rolls back tx, allowing workspace CAS saves to share one local
// transaction and connection.
func (r *SQLWriteGroupRepository) CreateInTx(ctx context.Context, tx *sql.Tx, group *WriteGroup) error {
	if tx == nil {
		return fmt.Errorf("create write group in transaction: %w", ErrWriteGroupValidation)
	}
	if group == nil {
		return fmt.Errorf("create write group: %w: group is nil", ErrWriteGroupValidation)
	}
	candidate := cloneWriteGroup(group)
	if err := r.validate(ctx, tx, candidate); err != nil {
		return err
	}
	now := r.now()
	candidate.ID = r.newID()
	candidate.Revision = r.newID()
	candidate.AppliedRevision = ""
	candidate.Status = WriteGroupStatusDraft
	candidate.CreatedAt = now
	candidate.UpdatedAt = now
	if err := r.insertGroup(ctx, tx, candidate); err != nil {
		return err
	}
	*group = *cloneWriteGroup(candidate)
	return nil
}

// GetInTx reads one group from the caller's transaction.
func (r *SQLWriteGroupRepository) GetInTx(ctx context.Context, tx *sql.Tx, workspaceID, groupID string) (*WriteGroup, error) {
	if tx == nil {
		return nil, fmt.Errorf("get write group in transaction: %w", ErrWriteGroupNotFound)
	}
	return r.get(ctx, tx, workspaceID, groupID)
}

// ListInTx reads all groups from the caller's transaction.
func (r *SQLWriteGroupRepository) ListInTx(ctx context.Context, tx *sql.Tx, workspaceID string) ([]*WriteGroup, error) {
	if tx == nil {
		return nil, fmt.Errorf("list write groups in transaction: %w", ErrWriteGroupNotFound)
	}
	return r.list(ctx, tx, workspaceID)
}

func (r *SQLWriteGroupRepository) insertGroup(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup) error {
	rowPolicy, err := json.Marshal(group.RowPolicy)
	if err != nil {
		return fmt.Errorf("encode write-group row policy: %w", err)
	}
	writePolicy, err := json.Marshal(group.WritePolicy)
	if err != nil {
		return fmt.Errorf("encode write-group write policy: %w", err)
	}
	migration, err := json.Marshal(group.Migration)
	if err != nil {
		return fmt.Errorf("encode write-group migration: %w", err)
	}
	_, err = runner.ExecContext(ctx, r.query(`
		INSERT INTO write_groups (
			id, workspace_id, revision, applied_revision, name, status,
			destination_connector_id, destination_connector_revision,
			destination_database, destination_table_schema, destination_table_name, destination_storage_strategy,
			destination_schema_revision, destination_schema_digest,
			row_policy, write_policy, migration, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`), group.ID, group.WorkspaceID, group.Revision, group.AppliedRevision, group.Name, group.Status,
		group.Destination.ConnectorID, group.Destination.ConnectorRevision,
		group.Destination.Database, group.Destination.TableSchema, group.Destination.TableName, group.Destination.StorageStrategy,
		group.Destination.SchemaRevision, group.Destination.SchemaDigest,
		string(rowPolicy), string(writePolicy), string(migration), group.CreatedAt, group.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert write group: %w", err)
	}
	for memberIndex, member := range group.Members {
		var measurementID any
		if member.MeasurementID != nil {
			measurementID = *member.MeasurementID
		}
		_, err = runner.ExecContext(ctx, r.query(`
			INSERT INTO write_group_members (
				group_id, member_index, device_id, point_id, tag_id,
				entity_key, source_revision, mapping_revision, measurement_id, target_column,
				required, max_age_seconds
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`), group.ID, memberIndex, member.DeviceID, member.PointID, member.TagID,
			member.EntityKey, member.SourceRevision, member.MappingRevision, measurementID, member.TargetColumn,
			member.Required, member.MaxAgeSeconds)
		if err != nil {
			return fmt.Errorf("insert write-group member %d: %w", memberIndex, err)
		}
	}
	return nil
}

func (r *SQLWriteGroupRepository) get(ctx context.Context, runner writeGroupSQLRunner, workspaceID, groupID string) (*WriteGroup, error) {
	row := runner.QueryRowContext(ctx, r.query(`
		SELECT id, workspace_id, revision, applied_revision, name, status,
			destination_connector_id, destination_connector_revision,
			destination_database, destination_table_schema, destination_table_name, destination_storage_strategy,
			destination_schema_revision, destination_schema_digest,
			row_policy, write_policy, migration, created_at, updated_at
		FROM write_groups
		WHERE workspace_id = $1 AND id = $2
	`), workspaceID, groupID)
	group, err := scanWriteGroup(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", ErrWriteGroupNotFound, ErrNotFound)
		}
		return nil, fmt.Errorf("read write group: %w", err)
	}
	if err := r.loadMembers(ctx, runner, group); err != nil {
		return nil, err
	}
	return cloneWriteGroup(group), nil
}

func (r *SQLWriteGroupRepository) list(ctx context.Context, runner writeGroupSQLRunner, workspaceID string) ([]*WriteGroup, error) {
	if _, err := readWorkspaceScope(ctx, runner, r.query, workspaceID); err != nil {
		return nil, err
	}
	rows, err := runner.QueryContext(ctx, r.query(`
		SELECT id, workspace_id, revision, applied_revision, name, status,
			destination_connector_id, destination_connector_revision,
			destination_database, destination_table_schema, destination_table_name, destination_storage_strategy,
			destination_schema_revision, destination_schema_digest,
			row_policy, write_policy, migration, created_at, updated_at
		FROM write_groups
		WHERE workspace_id = $1
		ORDER BY created_at ASC, id ASC
	`), workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list write groups: %w", err)
	}
	groups := make([]*WriteGroup, 0)
	for rows.Next() {
		group, scanErr := scanWriteGroup(rows)
		if scanErr != nil {
			_ = rows.Close()
			return nil, fmt.Errorf("scan write group: %w", scanErr)
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return nil, fmt.Errorf("iterate write groups: %w", err)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close write groups: %w", err)
	}
	for _, group := range groups {
		if err := r.loadMembers(ctx, runner, group); err != nil {
			return nil, err
		}
	}
	cloned := make([]*WriteGroup, len(groups))
	for i, group := range groups {
		cloned[i] = cloneWriteGroup(group)
	}
	return cloned, nil
}

func (r *SQLWriteGroupRepository) loadMembers(ctx context.Context, runner writeGroupSQLRunner, group *WriteGroup) error {
	rows, err := runner.QueryContext(ctx, r.query(`
		SELECT device_id, point_id, tag_id, entity_key, source_revision, mapping_revision,
			measurement_id, target_column, required, max_age_seconds
		FROM write_group_members
		WHERE group_id = $1
		ORDER BY member_index ASC
	`), group.ID)
	if err != nil {
		return fmt.Errorf("list write-group members: %w", err)
	}
	group.Members = make([]WriteGroupMember, 0)
	for rows.Next() {
		var member WriteGroupMember
		var measurementID sql.NullString
		var maxAgeSeconds sql.NullInt64
		var entityKey sql.NullString
		if err := rows.Scan(&member.DeviceID, &member.PointID, &member.TagID, &entityKey, &member.SourceRevision,
			&member.MappingRevision, &measurementID, &member.TargetColumn, &member.Required, &maxAgeSeconds); err != nil {
			_ = rows.Close()
			return fmt.Errorf("scan write-group member: %w", err)
		}
		if entityKey.Valid {
			member.EntityKey = entityKey.String
		}
		if measurementID.Valid && strings.TrimSpace(measurementID.String) != "" {
			value := measurementID.String
			member.MeasurementID = &value
		}
		if maxAgeSeconds.Valid {
			value := int(maxAgeSeconds.Int64)
			member.MaxAgeSeconds = &value
		}
		group.Members = append(group.Members, member)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("iterate write-group members: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close write-group members: %w", err)
	}
	return nil
}

func scanWriteGroup(row interface{ Scan(...any) error }) (*WriteGroup, error) {
	var group WriteGroup
	var status string
	var rowPolicy, writePolicy, migration string
	if err := row.Scan(&group.ID, &group.WorkspaceID, &group.Revision, &group.AppliedRevision,
		&group.Name, &status, &group.Destination.ConnectorID, &group.Destination.ConnectorRevision,
		&group.Destination.Database, &group.Destination.TableSchema, &group.Destination.TableName, &group.Destination.StorageStrategy,
		&group.Destination.SchemaRevision, &group.Destination.SchemaDigest,
		&rowPolicy, &writePolicy, &migration, &group.CreatedAt, &group.UpdatedAt); err != nil {
		return nil, err
	}
	group.Status = WriteGroupStatus(status)
	for _, field := range []struct {
		name   string
		raw    string
		target any
	}{
		{name: "row policy", raw: rowPolicy, target: &group.RowPolicy},
		{name: "write policy", raw: writePolicy, target: &group.WritePolicy},
		{name: "migration", raw: migration, target: &group.Migration},
	} {
		if strings.TrimSpace(field.raw) == "" {
			continue
		}
		if err := json.Unmarshal([]byte(field.raw), field.target); err != nil {
			return nil, fmt.Errorf("decode write-group %s: %w", field.name, err)
		}
	}
	return &group, nil
}

func (r *SQLWriteGroupRepository) query(query string) string {
	if r.postgres {
		return query
	}
	for i := 32; i >= 1; i-- {
		query = strings.ReplaceAll(query, fmt.Sprintf("$%d", i), "?")
	}
	return query
}

func (r *SQLWriteGroupRepository) trueLiteral() string {
	if r.postgres {
		return "TRUE"
	}
	return "1"
}
