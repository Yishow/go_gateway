package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"go-gateway/internal/datalink/schema"
)

var (
	// ErrManagedDestinationInternal means the saved destination aliases an
	// internal SQLite database used by the gateway.
	ErrManagedDestinationInternal = errors.New("managed destination is internal")
	// ErrManagedDestinationInvalid means the saved destination cannot be
	// inspected safely with the managed destination contract.
	ErrManagedDestinationInvalid = errors.New("managed destination is invalid")
)

const managedSQLiteSchemaName = "main"

// ManagedTableInspector validates a saved database connector before managed
// schema preview and inspects its existing tables without provisioning them.
type ManagedTableInspector struct {
	service    *ConnectorService
	internalDB *sql.DB
}

// NewManagedTableInspector creates an inspector bound to the saved connector
// service and the gateway's internal SQLite database.
func NewManagedTableInspector(service *ConnectorService, internalDB *sql.DB) *ManagedTableInspector {
	return &ManagedTableInspector{service: service, internalDB: internalDB}
}

// ValidateDestination resolves the current saved connector and applies the
// destination path guard. It does not connect to or provision the target.
func (i *ManagedTableInspector) ValidateDestination(ctx context.Context, connectorID, expectedRevision string) (*schema.DatabaseConnector, error) {
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	connector, err := i.service.ResolveSavedTarget(ctx, strings.TrimSpace(connectorID), strings.TrimSpace(expectedRevision))
	if err != nil {
		return nil, err
	}
	if connector == nil {
		return nil, managedDestinationInvalid("saved connector is unavailable")
	}

	switch connector.Kind {
	case schema.DatabaseConnectorKindSQLite:
		if _, err := i.inspectSQLiteDestination(ctx, connector); err != nil {
			return nil, err
		}
	case schema.DatabaseConnectorKindPostgres:
		if err := validateManagedPostgresConfig(connector); err != nil {
			return nil, err
		}
	default:
		return nil, managedDestinationInvalid("saved connector kind is unsupported")
	}
	return connector, nil
}

// InspectTable reports the requested table from the saved connector. A
// missing SQLite file is reported as missing only after its parent directory
// passes the managed path guard; no SQLite connection is opened in that case.
func (i *ManagedTableInspector) InspectTable(ctx context.Context, connectorID, schemaName, tableName string) (*TableInspection, error) {
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	connector, err := i.service.GetByID(ctx, strings.TrimSpace(connectorID))
	if err != nil {
		return nil, err
	}
	if connector == nil {
		return nil, managedDestinationInvalid("saved connector is unavailable")
	}
	if !connector.Enabled {
		return nil, fmt.Errorf("%w: saved connector is disabled", ErrConnectorDisabled)
	}
	return i.inspectManagedConnectorTable(ctx, connector, schemaName, tableName)
}

func (i *ManagedTableInspector) inspectManagedConnectorTable(ctx context.Context, connector *schema.DatabaseConnector, schemaName, tableName string) (*TableInspection, error) {

	switch connector.Kind {
	case schema.DatabaseConnectorKindSQLite:
		if strings.TrimSpace(schemaName) != "" && strings.TrimSpace(schemaName) != managedSQLiteSchemaName {
			return nil, managedDestinationInvalid("managed SQLite schema must be main")
		}
		missing, err := i.inspectSQLiteDestination(ctx, connector)
		if err != nil {
			return nil, err
		}
		if missing {
			return managedMissingTableInspection(schemaName, tableName), nil
		}
	case schema.DatabaseConnectorKindPostgres:
		if err := validateManagedPostgresConfig(connector); err != nil {
			return nil, err
		}
	default:
		return nil, managedDestinationInvalid("saved connector kind is unsupported")
	}

	return inspectReadOnlySavedConnectorTable(ctx, connector, schemaName, tableName)
}

// InspectTableAtRevision inspects the confirmed saved connector revision.
func (i *ManagedTableInspector) InspectTableAtRevision(ctx context.Context, connectorID, expectedRevision, schemaName, tableName string) (*TableInspection, error) {
	if i == nil || i.service == nil || i.service.repo == nil {
		return nil, managedDestinationInvalid("saved connector service is unavailable")
	}
	connector, err := i.service.ResolveSavedTarget(ctx, connectorID, expectedRevision)
	if err != nil {
		return nil, err
	}
	return i.inspectManagedConnectorTable(ctx, connector, schemaName, tableName)
}

func validateManagedPostgresConfig(connector *schema.DatabaseConnector) error {
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return managedDestinationInvalid("saved PostgreSQL connector configuration is invalid")
	}
	if _, err := buildExternalDBConfig(connector.Kind, config); err != nil {
		return managedDestinationInvalid("saved PostgreSQL connector configuration is incomplete")
	}
	return nil
}

type managedSQLiteFile struct {
	path     string
	resolved string
	info     os.FileInfo
}

func (i *ManagedTableInspector) inspectSQLiteDestination(ctx context.Context, connector *schema.DatabaseConnector) (bool, error) {
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return false, managedDestinationInvalid("saved SQLite connector configuration is invalid")
	}
	targetPath, err := managedSQLitePath(config)
	if err != nil {
		return false, err
	}
	internalFiles, err := managedInternalSQLiteFiles(ctx, i.internalDB)
	if err != nil {
		return false, err
	}
	if managedSQLitePathAliasesInternal(targetPath, internalFiles) {
		return false, managedDestinationInternal()
	}

	_, err = os.Lstat(targetPath)
	if errors.Is(err, os.ErrNotExist) {
		if err := validateManagedMissingParent(targetPath); err != nil {
			return false, err
		}
		return true, nil
	}
	if err != nil {
		return false, managedDestinationInvalid("saved SQLite destination cannot be inspected")
	}

	targetInfo, err := os.Stat(targetPath)
	if errors.Is(err, os.ErrNotExist) {
		// An existing symlink whose target vanished is not a safe missing
		// destination: it is an unverified alias.
		return false, managedDestinationInvalid("saved SQLite destination cannot be verified")
	}
	if err != nil || !targetInfo.Mode().IsRegular() {
		return false, managedDestinationInvalid("saved SQLite destination is not a regular file")
	}
	targetResolved, err := filepath.EvalSymlinks(targetPath)
	if err != nil {
		return false, managedDestinationInvalid("saved SQLite destination cannot be verified")
	}
	targetResolved, err = filepath.Abs(targetResolved)
	if err != nil {
		return false, managedDestinationInvalid("saved SQLite destination cannot be verified")
	}
	for _, internal := range internalFiles {
		if os.SameFile(targetInfo, internal.info) || filepath.Clean(targetResolved) == filepath.Clean(internal.resolved) || filepath.Clean(targetPath) == filepath.Clean(internal.path) {
			return false, managedDestinationInternal()
		}
	}
	return false, nil
}

func managedSQLitePathAliasesInternal(targetPath string, internalFiles []managedSQLiteFile) bool {
	for _, internal := range internalFiles {
		if filepath.Clean(targetPath) == filepath.Clean(internal.path) {
			return true
		}
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(targetPath))
	if err != nil {
		return false
	}
	resolvedCandidate, err := filepath.Abs(filepath.Join(parent, filepath.Base(targetPath)))
	if err != nil {
		return false
	}
	for _, internal := range internalFiles {
		if filepath.Clean(resolvedCandidate) == filepath.Clean(internal.resolved) {
			return true
		}
	}
	return false
}

func managedInternalSQLiteFiles(ctx context.Context, db *sql.DB) ([]managedSQLiteFile, error) {
	if db == nil {
		return nil, managedDestinationInvalid("internal SQLite database is unavailable")
	}
	rows, err := db.QueryContext(ctx, `PRAGMA database_list`)
	if err != nil {
		return nil, managedDestinationInvalid("internal SQLite database identity cannot be confirmed")
	}
	defer rows.Close()

	files := make([]managedSQLiteFile, 0, 1)
	for rows.Next() {
		var sequence int
		var name, rawPath string
		if err := rows.Scan(&sequence, &name, &rawPath); err != nil {
			return nil, managedDestinationInvalid("internal SQLite database identity cannot be confirmed")
		}
		if strings.TrimSpace(rawPath) == "" {
			return nil, managedDestinationInvalid("internal SQLite database identity cannot be confirmed")
		}
		path, err := normalizeManagedSQLitePath(rawPath)
		if err != nil {
			return nil, managedDestinationInvalid("internal SQLite database identity cannot be confirmed")
		}
		info, err := os.Stat(path)
		if err != nil || !info.Mode().IsRegular() {
			return nil, managedDestinationInvalid("internal SQLite database identity cannot be confirmed")
		}
		resolved, err := filepath.EvalSymlinks(path)
		if err != nil {
			return nil, managedDestinationInvalid("internal SQLite database identity cannot be confirmed")
		}
		resolved, err = filepath.Abs(resolved)
		if err != nil {
			return nil, managedDestinationInvalid("internal SQLite database identity cannot be confirmed")
		}
		files = append(files, managedSQLiteFile{path: path, resolved: resolved, info: info})
	}
	if err := rows.Err(); err != nil {
		return nil, managedDestinationInvalid("internal SQLite database identity cannot be confirmed")
	}
	if err := rows.Close(); err != nil {
		return nil, managedDestinationInvalid("internal SQLite database identity cannot be confirmed")
	}
	if len(files) == 0 {
		return nil, managedDestinationInvalid("internal SQLite database identity cannot be confirmed")
	}
	return files, nil
}

func managedSQLitePath(config ConnectionConfig) (string, error) {
	raw := strings.TrimSpace(stringConfigValue(config, "dsn"))
	if raw == "" {
		raw = strings.TrimSpace(stringConfigValue(config, "path"))
	}
	if raw == "" {
		return "", managedDestinationInvalid("saved SQLite connector lacks a file destination")
	}
	path, err := managedSQLitePathFromRaw(raw)
	if err != nil {
		return "", managedDestinationInvalid("saved SQLite connector must use a local file destination")
	}
	return path, nil
}

func managedSQLitePathFromRaw(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("empty sqlite path")
	}
	var path string
	if strings.HasPrefix(strings.ToLower(raw), "file:") {
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Host != "" || parsed.Fragment != "" {
			return "", errors.New("non-local sqlite uri")
		}
		query, err := url.ParseQuery(parsed.RawQuery)
		if err != nil || !validManagedSQLiteQuery(query) {
			return "", errors.New("unsafe sqlite uri query")
		}
		if parsed.Opaque != "" {
			path, err = url.PathUnescape(parsed.Opaque)
		} else {
			path, err = url.PathUnescape(parsed.Path)
		}
		if err != nil {
			return "", errors.New("invalid sqlite uri path")
		}
	} else {
		base, rawQuery, hasQuery := strings.Cut(raw, "?")
		path = base
		if hasQuery {
			query, err := url.ParseQuery(rawQuery)
			if err != nil || !validManagedSQLiteQuery(query) {
				return "", errors.New("unsafe sqlite path query")
			}
		}
	}
	path = strings.TrimSpace(path)
	if path == "" || strings.EqualFold(path, ":memory:") || strings.EqualFold(path, "memory") || strings.ContainsRune(path, '\x00') {
		return "", errors.New("sqlite memory or empty path")
	}
	return normalizeManagedSQLitePath(path)
}

func normalizeManagedSQLitePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" || strings.ContainsRune(path, '\x00') {
		return "", errors.New("empty sqlite path")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(absPath), nil
}

func validManagedSQLiteQuery(query url.Values) bool {
	for key, values := range query {
		if len(values) != 1 {
			return false
		}
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "mode":
			mode := strings.ToLower(strings.TrimSpace(values[0]))
			if mode != "" && mode != "ro" && mode != "rw" && mode != "rwc" {
				return false
			}
		case "cache":
			cache := strings.ToLower(strings.TrimSpace(values[0]))
			if cache != "" && cache != "shared" && cache != "private" {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func validateManagedMissingParent(path string) error {
	parent := filepath.Dir(path)
	info, err := os.Stat(parent)
	if err != nil || !info.IsDir() {
		return managedDestinationInvalid("saved SQLite destination parent is unavailable")
	}
	permissions := info.Mode().Perm()
	if permissions&0o222 == 0 || permissions&0o111 == 0 {
		return managedDestinationInvalid("saved SQLite destination parent is not writable")
	}
	if _, err := filepath.EvalSymlinks(parent); err != nil {
		return managedDestinationInvalid("saved SQLite destination parent cannot be verified")
	}
	handle, err := os.OpenFile(parent, os.O_WRONLY, 0)
	if err == nil {
		if err := handle.Close(); err != nil {
			return managedDestinationInvalid("saved SQLite destination parent is not writable")
		}
		return nil
	}
	if errors.Is(err, syscall.EISDIR) {
		return nil
	}
	return managedDestinationInvalid("saved SQLite destination parent is not writable")
}

func managedMissingTableInspection(schemaName, tableName string) *TableInspection {
	schemaName = strings.TrimSpace(schemaName)
	if schemaName == "" {
		schemaName = managedSQLiteSchemaName
	}
	return &TableInspection{
		Status: TableInspectionMissing, Schema: schemaName, Table: strings.TrimSpace(tableName), Columns: []ColumnInfo{},
	}
}

func managedDestinationInvalid(reason string) error {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = "saved destination requires repair"
	}
	return fmt.Errorf("%w: %s", ErrManagedDestinationInvalid, reason)
}

func managedDestinationInternal() error {
	return fmt.Errorf("%w: saved destination is reserved for gateway storage", ErrManagedDestinationInternal)
}
