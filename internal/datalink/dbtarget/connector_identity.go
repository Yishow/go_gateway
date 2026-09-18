package dbtarget

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

var connectorNonIdentityConfigKeys = []string{
	"schema",
	"table",
	"table_schema",
	"table_name",
	"column_name",
	"write_mode",
	"timestamp_column",
	"group_key",
	"write_interval_seconds",
}

func newConnectorIDs() (id, identityRevision string, err error) {
	id, err = common.NewUUID()
	if err != nil {
		return "", "", fmt.Errorf("建立資料庫連接器 ID 失敗: %w", err)
	}
	identityRevision, err = newConnectorIdentityRevision()
	if err != nil {
		return "", "", err
	}
	return id, identityRevision, nil
}

func newConnectorIdentityRevision() (string, error) {
	revision, err := common.NewUUID()
	if err != nil {
		return "", fmt.Errorf("建立資料庫連接器 identity revision 失敗: %w", err)
	}
	return revision, nil
}

func validateConnectorIdentityForUpdate(connector *schema.DatabaseConnector, expectedRevision *string) (string, schema.DatabaseConnectorKind, error) {
	currentRevision, err := validateConnectorIdentityRevision(connector, expectedRevision)
	if err != nil {
		return "", "", err
	}
	return currentRevision, connector.Kind, nil
}

func finalizeConnectorIdentityUpdate(
	connector *schema.DatabaseConnector,
	previousKind schema.DatabaseConnectorKind,
	previousConfig, nextConfig ConnectionConfig,
	requestedConfig *ConnectionConfig,
) error {
	if previousKind != connector.Kind {
		requestedPassword := ""
		if requestedConfig != nil {
			requestedPassword, _ = passwordConfigValue(*requestedConfig)
		}
		if requestedPassword == "" {
			delete(nextConfig, "password")
			serialized, err := serializeConnectionConfig(nextConfig)
			if err != nil {
				return err
			}
			connector.ConnectionConfig = serialized
		}
	}
	if connectorIdentityChanged(previousKind, previousConfig, connector.Kind, nextConfig) {
		revision, err := newConnectorIdentityRevision()
		if err != nil {
			return err
		}
		connector.IdentityRevision = revision
	}
	return nil
}

func connectorIdentityChanged(
	previousKind schema.DatabaseConnectorKind,
	previousConfig ConnectionConfig,
	nextKind schema.DatabaseConnectorKind,
	nextConfig ConnectionConfig,
) bool {
	if previousKind != nextKind {
		return true
	}

	previousJSON, previousErr := json.Marshal(identityConnectionConfig(previousConfig))
	nextJSON, nextErr := json.Marshal(identityConnectionConfig(nextConfig))
	if previousErr != nil || nextErr != nil {
		return true
	}
	return !bytes.Equal(previousJSON, nextJSON)
}

func identityConnectionConfig(config ConnectionConfig) ConnectionConfig {
	identity := cloneConnectionConfig(config)
	for _, key := range connectorNonIdentityConfigKeys {
		delete(identity, key)
	}
	if password, ok := passwordConfigValue(identity); ok && password == "" {
		delete(identity, "password")
	}
	return identity
}

func endpointConnectionConfig(config ConnectionConfig) ConnectionConfig {
	endpoint := identityConnectionConfig(config)
	delete(endpoint, "password")
	return endpoint
}

func connectorEndpointChanged(previousConfig, nextConfig ConnectionConfig) bool {
	previousJSON, previousErr := json.Marshal(endpointConnectionConfig(previousConfig))
	nextJSON, nextErr := json.Marshal(endpointConnectionConfig(nextConfig))
	if previousErr != nil || nextErr != nil {
		return true
	}
	return !bytes.Equal(previousJSON, nextJSON)
}

func mergeConnectorConnectionConfig(
	existing ConnectionConfig,
	next ConnectionConfig,
	clearPassword bool,
) ConnectionConfig {
	merged := cloneConnectionConfig(next)
	identityChanged := connectorEndpointChanged(existing, next)
	if clearPassword {
		delete(merged, "password")
		return merged
	}

	currentPassword, hasCurrentPassword := passwordConfigValue(existing)
	nextPassword, hasNextPassword := passwordConfigValue(next)
	if identityChanged {
		if !hasNextPassword || nextPassword == "" {
			delete(merged, "password")
			return merged
		}
		return merged
	}
	// identityChanged is false here: the branch above already returned.
	if hasCurrentPassword && currentPassword != "" && (!hasNextPassword || nextPassword == "") {
		merged["password"] = currentPassword
	}
	return merged
}

func passwordConfigValue(config ConnectionConfig) (string, bool) {
	value, ok := config["password"]
	if !ok || value == nil {
		return "", false
	}
	if password, ok := value.(string); ok {
		return password, true
	}
	return fmt.Sprintf("%v", value), true
}

func validateConnectorIdentityRevision(connector *schema.DatabaseConnector, expectedRevision *string) (string, error) {
	currentRevision := strings.TrimSpace(connector.IdentityRevision)
	if currentRevision == "" {
		return "", fmt.Errorf("%w: 資料庫連接器 identity revision 不可為空", ErrConnectorRevisionConflict)
	}
	if expectedRevision == nil {
		return currentRevision, nil
	}
	expected := strings.TrimSpace(*expectedRevision)
	if expected == "" || expected != currentRevision {
		return "", fmt.Errorf("%w: 資料庫連接器 identity revision 已變更", ErrConnectorRevisionConflict)
	}
	return currentRevision, nil
}

// ResolveSavedTarget returns a persisted, enabled connector only when the
// caller presents its current opaque identity revision. It never probes the
// external database.
func (s *ConnectorService) ResolveSavedTarget(
	ctx context.Context,
	connectorID string,
	expectedRevision string,
) (*schema.DatabaseConnector, error) {
	if strings.TrimSpace(expectedRevision) == "" {
		return nil, fmt.Errorf("%w: 資料庫連接器 identity revision 不可為空", ErrConnectorRevisionConflict)
	}

	connector, err := s.repo.GetByID(ctx, strings.TrimSpace(connectorID))
	if err != nil {
		if errors.Is(err, ErrConnectorNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("取得已保存資料庫連接器失敗: %w", err)
	}
	if !connector.Enabled {
		return nil, fmt.Errorf("%w: 資料庫連接器已停用", ErrConnectorDisabled)
	}
	if strings.TrimSpace(connector.IdentityRevision) != strings.TrimSpace(expectedRevision) {
		return nil, fmt.Errorf("%w: 資料庫連接器 identity revision 已變更", ErrConnectorRevisionConflict)
	}
	return connector, nil
}
