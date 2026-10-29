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
	candidate, err := r.prepareNewWriteGroupDraftInTx(ctx, tx, group)
	if err != nil {
		return err
	}
	if err := r.insertGroup(ctx, tx, candidate); err != nil {
		return err
	}
	if candidate.RecordingProfile != nil {
		if err := r.insertProfileKeyInTx(ctx, tx, candidate); err != nil {
			return err
		}
	}
	if err := r.ensureStorageScopeInTx(ctx, tx, candidate); err != nil {
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
	profile := profileColumnValues(group)
	_, err = runner.ExecContext(ctx, r.query(`
		INSERT INTO write_groups (
			id, workspace_id, revision, applied_revision, name, status,
			destination_connector_id, destination_connector_revision,
			destination_database, destination_table_schema, destination_table_name, destination_storage_strategy,
			destination_schema_revision, destination_schema_digest,
			row_policy, write_policy, migration, created_at, updated_at,
			layout_version, profile_kind, profile_interval, profile_unit, storage_mode, sample_reduction
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25)
	`), group.ID, group.WorkspaceID, group.Revision, group.AppliedRevision, group.Name, group.Status,
		group.Destination.ConnectorID, group.Destination.ConnectorRevision,
		group.Destination.Database, group.Destination.TableSchema, group.Destination.TableName, group.Destination.StorageStrategy,
		group.Destination.SchemaRevision, group.Destination.SchemaDigest,
		string(rowPolicy), string(writePolicy), string(migration), group.CreatedAt, group.UpdatedAt,
		profile.layout, profile.kind, profile.interval, profile.unit, profile.mode, profile.reduction)
	if err != nil {
		return fmt.Errorf("insert write group: %w", err)
	}
	for memberIndex, member := range group.Members {
		var measurementID any
		if member.MeasurementID != nil {
			measurementID = *member.MeasurementID
		}
		var exactRecordingSource any
		if member.ExactRecordingSource != nil {
			payload, err := json.Marshal(member.ExactRecordingSource)
			if err != nil {
				return fmt.Errorf("encode exact recording source member %d: %w", memberIndex, err)
			}
			exactRecordingSource = string(payload)
		}
		_, err = runner.ExecContext(ctx, r.query(`
			INSERT INTO write_group_members (
				group_id, member_index, device_id, point_id, tag_id,
				entity_key, source_revision, mapping_revision, measurement_id, target_column,
				required, max_age_seconds, circuit_id, circuit_name,
				value_role, source_value_type, result_value_type, exact_recording_source
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		`), group.ID, memberIndex, member.DeviceID, member.PointID, member.TagID,
			member.EntityKey, member.SourceRevision, member.MappingRevision, measurementID, member.TargetColumn,
			member.Required, member.MaxAgeSeconds, nullableText(member.CircuitID), nullableText(member.CircuitName),
			nullableText(string(member.ValueRole)), nullableText(string(member.SourceValueType)), nullableText(string(member.ResultValueType)),
			exactRecordingSource)
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
			row_policy, write_policy, migration, created_at, updated_at,
			layout_version, profile_kind, profile_interval, profile_unit, storage_mode, sample_reduction
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
	if err := r.loadBasicManagedDeviceID(ctx, runner, group); err != nil {
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
			row_policy, write_policy, migration, created_at, updated_at,
			layout_version, profile_kind, profile_interval, profile_unit, storage_mode, sample_reduction
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
		if err := r.loadBasicManagedDeviceID(ctx, runner, group); err != nil {
			return nil, err
		}
	}
	cloned := make([]*WriteGroup, len(groups))
	for i, group := range groups {
		cloned[i] = cloneWriteGroup(group)
	}
	return cloned, nil
}

func scanWriteGroup(row interface{ Scan(...any) error }) (*WriteGroup, error) {
	var group WriteGroup
	var status string
	var rowPolicy, writePolicy, migration string
	var layoutVersion, profileKind, profileUnit, storageMode, sampleReduction sql.NullString
	var profileInterval sql.NullInt64
	if err := row.Scan(&group.ID, &group.WorkspaceID, &group.Revision, &group.AppliedRevision,
		&group.Name, &status, &group.Destination.ConnectorID, &group.Destination.ConnectorRevision,
		&group.Destination.Database, &group.Destination.TableSchema, &group.Destination.TableName, &group.Destination.StorageStrategy,
		&group.Destination.SchemaRevision, &group.Destination.SchemaDigest,
		&rowPolicy, &writePolicy, &migration, &group.CreatedAt, &group.UpdatedAt,
		&layoutVersion, &profileKind, &profileInterval, &profileUnit, &storageMode, &sampleReduction); err != nil {
		return nil, err
	}
	group.Status = WriteGroupStatus(status)
	if err := applyProfileColumns(&group, layoutVersion, profileKind, profileUnit, storageMode, profileInterval, sampleReduction); err != nil {
		return nil, err
	}
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
