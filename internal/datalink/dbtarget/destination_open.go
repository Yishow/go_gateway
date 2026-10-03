package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"maps"
	"net/url"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// ErrDestinationBlocked means a saved destination cannot be used until it is
// repaired: its identity revision changed, it is disabled or gone, or its
// credentials are rejected. Delivery keeps the data and waits instead of
// retargeting.
var ErrDestinationBlocked = errors.New("destination blocked")

// ErrDestinationUnreachable means the destination is offline right now. It is
// transient: the same identity is retried later.
var ErrDestinationUnreachable = errors.New("destination unreachable")

const (
	sqliteDeliveryModeReadOnly            = "ro"
	sqliteDeliveryModeReadWrite           = "rw"
	sqliteDeliveryModeReadWriteCreate     = "rwc"
	sqliteDeliveryModeMemory              = "memory"
	sqliteDeliveryBusyTimeoutMilliseconds = 15000
)

// OpenedDestination is an open connection to a saved destination.
type OpenedDestination struct {
	DB    *sql.DB
	Kind  schema.DatabaseConnectorKind
	Close func() error
}

// OpenDestination opens the saved connector for delivery, but only while its
// identity revision still equals expectedRevision, so accepted rows can never
// be sent to an edited endpoint. It never creates a missing database file: a
// vanished SQLite destination is an outage to surface, not an empty database
// to fill.
func (s *ConnectorService) OpenDestination(ctx context.Context, connectorID, expectedRevision string) (*OpenedDestination, error) {
	connector, err := s.repo.GetByID(ctx, connectorID)
	if errors.Is(err, ErrConnectorNotFound) {
		return nil, fmt.Errorf("%w: connector no longer exists", ErrDestinationBlocked)
	}
	if err != nil {
		return nil, fmt.Errorf("read destination connector: %w", err)
	}
	if connector.IdentityRevision != expectedRevision {
		return nil, fmt.Errorf("%w: connector identity changed", ErrDestinationBlocked)
	}
	if !connector.Enabled {
		return nil, fmt.Errorf("%w: connector is disabled", ErrDestinationBlocked)
	}
	if connector.Kind != schema.DatabaseConnectorKindSQLite && connector.Kind != schema.DatabaseConnectorKindPostgres {
		return nil, fmt.Errorf("%w: connector kind is not supported for group delivery", ErrDestinationBlocked)
	}
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: connector configuration is invalid", ErrDestinationBlocked)
	}
	if connector.Kind == schema.DatabaseConnectorKindSQLite {
		config, err = prepareSQLiteDeliveryConfig(config)
		if err != nil {
			return nil, fmt.Errorf("%w: connector configuration is invalid", ErrDestinationBlocked)
		}
	}
	manager, err := openExternalDBManagerFunc(connector.Kind, config)
	if err != nil {
		if classifyConnectorError(err) == schema.DatabaseConnectorStatusAuthFailed {
			return nil, fmt.Errorf("%w: credentials were rejected", ErrDestinationBlocked)
		}
		return nil, fmt.Errorf("%w: %w", ErrDestinationUnreachable, err)
	}
	fixtureDB, fixtureClose, err := openFixtureDestination(connector.ID, connector.Kind, config, manager)
	if err != nil {
		_ = manager.Close()
		return nil, fmt.Errorf("%w: fixture destination unavailable", ErrDestinationUnreachable)
	}
	if fixtureDB != nil {
		return &OpenedDestination{DB: fixtureDB, Kind: connector.Kind, Close: fixtureClose}, nil
	}
	return &OpenedDestination{DB: manager.DB(), Kind: connector.Kind, Close: manager.Close}, nil
}

// prepareSQLiteDeliveryConfig changes only the delivery DSN. The SQLite rw
// mode makes opening an existing file atomic, so a file that disappears after
// connector lookup cannot be replaced by a new empty database.
func prepareSQLiteDeliveryConfig(config ConnectionConfig) (ConnectionConfig, error) {
	dsn := stringConfigValue(config, "dsn")
	if dsn == "" {
		dsn = stringConfigValue(config, "path")
	}
	if dsn == "" {
		return nil, errors.New("sqlite connector configuration lacks dsn/path")
	}
	preparedDSN, err := prepareSQLiteDeliveryDSN(dsn)
	if err != nil {
		return nil, err
	}
	prepared := maps.Clone(config)
	if prepared == nil {
		prepared = make(ConnectionConfig)
	}
	prepared["dsn"] = preparedDSN
	return prepared, nil
}

func prepareSQLiteDeliveryDSN(rawDSN string) (string, error) {
	dsn := strings.TrimSpace(rawDSN)
	if dsn == "" || isSQLiteMemoryDSN(dsn) {
		return dsn, nil
	}
	if strings.HasPrefix(strings.ToLower(dsn), "file:") {
		parsed, err := url.Parse(dsn)
		if err != nil {
			return "", fmt.Errorf("parse sqlite file dsn: %w", err)
		}
		query, err := url.ParseQuery(parsed.RawQuery)
		if err != nil {
			return "", fmt.Errorf("parse sqlite dsn query: %w", err)
		}
		mode, err := normalizeSQLiteDeliveryMode(query)
		if err != nil {
			return "", err
		}
		if mode == sqliteDeliveryModeReadOnly || mode == sqliteDeliveryModeMemory {
			return dsn, nil
		}
		ensureSQLiteDeliveryBusyTimeout(query)
		parsed.RawQuery = query.Encode()
		return parsed.String(), nil
	}

	base, rawQuery, hasQuery := dsn, "", false
	if candidateBase, candidateQuery, found := strings.Cut(dsn, "?"); found && strings.Contains(candidateQuery, "=") {
		base, rawQuery, hasQuery = candidateBase, candidateQuery, true
	}
	query := make(url.Values)
	if hasQuery {
		var err error
		query, err = url.ParseQuery(rawQuery)
		if err != nil {
			return "", fmt.Errorf("parse sqlite dsn query: %w", err)
		}
	}
	mode, err := normalizeSQLiteDeliveryMode(query)
	if err != nil {
		return "", err
	}
	if mode == sqliteDeliveryModeMemory {
		return "", errors.New("sqlite memory mode requires an explicit file URI")
	}
	if mode == sqliteDeliveryModeReadOnly {
		return sqliteDeliveryFileDSN(base, query), nil
	}
	ensureSQLiteDeliveryBusyTimeout(query)
	return sqliteDeliveryFileDSN(base, query), nil
}

func ensureSQLiteDeliveryBusyTimeout(query url.Values) {
	for _, pragma := range query["_pragma"] {
		normalized := strings.ToLower(strings.Join(strings.Fields(pragma), ""))
		if strings.HasPrefix(normalized, "busy_timeout(") || strings.HasPrefix(normalized, "busy_timeout=") {
			return
		}
	}
	query.Add("_pragma", fmt.Sprintf("busy_timeout(%d)", sqliteDeliveryBusyTimeoutMilliseconds))
}

func sqliteDeliveryFileDSN(path string, query url.Values) string {
	if strings.HasPrefix(path, "/") {
		return (&url.URL{Scheme: "file", Path: path, RawQuery: query.Encode()}).String()
	}
	return (&url.URL{Scheme: "file", Opaque: url.PathEscape(path), RawQuery: query.Encode()}).String()
}

func normalizeSQLiteDeliveryMode(query url.Values) (string, error) {
	mode := strings.ToLower(strings.TrimSpace(query.Get("mode")))
	switch mode {
	case "":
		query.Set("mode", sqliteDeliveryModeReadWrite)
		return sqliteDeliveryModeReadWrite, nil
	case sqliteDeliveryModeReadWriteCreate, sqliteDeliveryModeReadWrite:
		query.Set("mode", sqliteDeliveryModeReadWrite)
		return sqliteDeliveryModeReadWrite, nil
	case sqliteDeliveryModeReadOnly, sqliteDeliveryModeMemory:
		return mode, nil
	default:
		return "", fmt.Errorf("unsupported sqlite mode %q", mode)
	}
}

func isSQLiteMemoryDSN(dsn string) bool {
	trimmed := strings.TrimSpace(dsn)
	if strings.EqualFold(trimmed, ":memory:") {
		return true
	}
	if !strings.HasPrefix(strings.ToLower(trimmed), "file:") {
		return false
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
	}
	if strings.EqualFold(parsed.Opaque, ":memory:") || strings.EqualFold(parsed.Path, ":memory:") {
		return true
	}
	query, err := url.ParseQuery(parsed.RawQuery)
	return err == nil && strings.EqualFold(strings.TrimSpace(query.Get("mode")), sqliteDeliveryModeMemory)
}
