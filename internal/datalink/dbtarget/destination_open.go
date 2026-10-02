package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
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
		if err := requireExistingSQLiteFile(config); err != nil {
			return nil, err
		}
	}
	manager, err := openExternalDBManagerFunc(connector.Kind, config)
	if err != nil {
		if classifyConnectorError(err) == schema.DatabaseConnectorStatusAuthFailed {
			return nil, fmt.Errorf("%w: credentials were rejected", ErrDestinationBlocked)
		}
		return nil, fmt.Errorf("%w: %w", ErrDestinationUnreachable, err)
	}
	return &OpenedDestination{DB: manager.DB(), Kind: connector.Kind, Close: manager.Close}, nil
}

// requireExistingSQLiteFile refuses to open a plain-path SQLite destination
// whose file is gone, because opening it would silently create an empty one.
func requireExistingSQLiteFile(config ConnectionConfig) error {
	dsn := strings.TrimSpace(stringConfigValue(config, "dsn"))
	if dsn == "" {
		dsn = strings.TrimSpace(stringConfigValue(config, "path"))
	}
	if dsn == "" {
		return fmt.Errorf("%w: connector configuration is invalid", ErrDestinationBlocked)
	}
	if strings.HasPrefix(dsn, "file:") || strings.EqualFold(dsn, ":memory:") {
		return nil
	}
	if _, err := os.Stat(dsn); err != nil {
		return fmt.Errorf("%w: database file is not available", ErrDestinationUnreachable)
	}
	return nil
}
