package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/schema"
)

// ErrSetupTransactionUnavailable reports repositories that cannot join a
// caller-owned local transaction.
var ErrSetupTransactionUnavailable = errors.New("database setup transaction is unavailable")

type connectorTxBinder interface {
	WithTx(tx *sql.Tx) ConnectorRepository
}

type preparedConnectorMode int

const (
	preparedConnectorCreate preparedConnectorMode = iota + 1
	preparedConnectorUpdate
	preparedConnectorExisting
)

// PreparedConnector is a validated connector change that is not persisted
// yet. Connection probes run while preparing, before a local transaction
// opens, so the configuration database is never locked during a probe.
type PreparedConnector struct {
	Connector        *schema.DatabaseConnector
	mode             preparedConnectorMode
	expectedRevision string
}

// PrepareExisting keeps a saved connector unchanged; persisting it re-reads
// the connector and requires the identity revision it was resolved with.
func PrepareExisting(connector *schema.DatabaseConnector) *PreparedConnector {
	return &PreparedConnector{
		Connector:        connector,
		mode:             preparedConnectorExisting,
		expectedRevision: strings.TrimSpace(connector.IdentityRevision),
	}
}

// PrepareCreate validates and probes a new connector without writing it.
func (s *ConnectorService) PrepareCreate(ctx context.Context, req CreateConnectorRequest) (*PreparedConnector, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, validationError("資料庫連接器名稱不可為空")
	}
	connectionConfigJSON, err := serializeConnectionConfig(req.ConnectionConfig)
	if err != nil {
		return nil, err
	}
	id, identityRevision, err := newConnectorIDs()
	if err != nil {
		return nil, err
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	defaultIntervalSeconds := defaultWriteIntervalSeconds(req.DefaultWriteIntervalSeconds)
	status, lastCheckAt, lastCheckError := probeConnector(ctx, req.Kind, req.ConnectionConfig)
	now := time.Now()
	connector := &schema.DatabaseConnector{
		ID:                          id,
		Name:                        name,
		Kind:                        req.Kind,
		ConnectionConfig:            connectionConfigJSON,
		IdentityRevision:            identityRevision,
		Status:                      status,
		LastCheckAt:                 lastCheckAt,
		LastCheckError:              lastCheckError,
		Enabled:                     enabled,
		DefaultWriteIntervalSeconds: defaultIntervalSeconds,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}
	return &PreparedConnector{Connector: connector, mode: preparedConnectorCreate}, nil
}

// PrepareUpdate validates and probes a connector update without writing it.
func (s *ConnectorService) PrepareUpdate(ctx context.Context, id string, req UpdateConnectorRequest) (*PreparedConnector, error) {
	connector, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("取得資料庫連接器失敗: %w", err)
	}
	expectedIdentityRevision, previousKind, err := validateConnectorIdentityForUpdate(connector, req.ExpectedIdentityRevision)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, validationError("資料庫連接器名稱不可為空")
		}
		connector.Name = name
	}
	if req.Kind != nil {
		connector.Kind = *req.Kind
	}
	if req.Enabled != nil {
		connector.Enabled = *req.Enabled
	}
	if req.DefaultWriteIntervalSeconds != nil {
		connector.DefaultWriteIntervalSeconds = defaultWriteIntervalSeconds(req.DefaultWriteIntervalSeconds)
	}
	existingConnectionConfig, err := parseConnectionConfig(connector.ConnectionConfig)
	if err != nil {
		return nil, err
	}
	clearPassword := req.ClearPassword != nil && *req.ClearPassword
	var connectionConfig ConnectionConfig
	if req.ConnectionConfig != nil {
		connectionConfig = preserveSensitiveConnectionConfigValues(
			existingConnectionConfig,
			*req.ConnectionConfig,
			clearPassword,
		)
		connector.ConnectionConfig, err = serializeConnectionConfig(connectionConfig)
		if err != nil {
			return nil, err
		}
	} else {
		connectionConfig = cloneConnectionConfig(existingConnectionConfig)
		if clearPassword {
			delete(connectionConfig, "password")
			connector.ConnectionConfig, err = serializeConnectionConfig(connectionConfig)
			if err != nil {
				return nil, err
			}
		}
	}
	if err := finalizeConnectorIdentityUpdate(connector, previousKind, existingConnectionConfig, connectionConfig, req.ConnectionConfig); err != nil {
		return nil, err
	}
	connector.Status, connector.LastCheckAt, connector.LastCheckError = probeConnector(ctx, connector.Kind, connectionConfig)
	connector.UpdatedAt = time.Now()
	return &PreparedConnector{Connector: connector, mode: preparedConnectorUpdate, expectedRevision: expectedIdentityRevision}, nil
}

// PersistInTx writes a prepared connector inside tx. Updates keep the
// identity-revision compare-and-swap; an existing connector is re-read and
// must still carry the revision it was prepared from.
func (s *ConnectorService) PersistInTx(ctx context.Context, tx *sql.Tx, prepared *PreparedConnector) error {
	if prepared == nil || prepared.Connector == nil {
		return validationError("資料庫連接器變更不可為空")
	}
	binder, ok := s.repo.(connectorTxBinder)
	if !ok || tx == nil {
		return ErrSetupTransactionUnavailable
	}
	repo := binder.WithTx(tx)
	switch prepared.mode {
	case preparedConnectorCreate:
		if err := repo.Create(ctx, prepared.Connector); err != nil {
			return fmt.Errorf("建立資料庫連接器失敗: %w", err)
		}
	case preparedConnectorUpdate:
		if err := repo.UpdateWithExpectedIdentityRevision(ctx, prepared.Connector, prepared.expectedRevision); err != nil {
			return fmt.Errorf("更新資料庫連接器失敗: %w", err)
		}
	case preparedConnectorExisting:
		current, err := repo.GetByID(ctx, prepared.Connector.ID)
		if err != nil {
			return fmt.Errorf("取得資料庫連接器失敗: %w", err)
		}
		if strings.TrimSpace(current.IdentityRevision) != prepared.expectedRevision {
			return fmt.Errorf("%w: 資料庫連接器 identity revision 已變更", ErrConnectorRevisionConflict)
		}
	default:
		return validationError("未知的資料庫連接器變更")
	}
	return nil
}
