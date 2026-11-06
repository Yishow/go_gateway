package device

import (
	"encoding/json"
	"errors"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// The declaration remains part of the existing device configuration authority.
// Optional properties need explicit validation; the historical schema validator
// only visits required fields and must not silently accept acquisition policy.
func validateAcquisitionDeclaration(protocol schema.ProtocolType, config map[string]any) error {
	if _, declared := config["acquisition"]; !declared {
		return nil
	}
	raw, err := json.Marshal(config)
	if err != nil {
		return connector.ErrAcquisitionConfig
	}
	_, err = connector.ResolveAcquisitionConfig(protocol, string(raw))
	if errors.Is(err, connector.ErrAcquisitionRandomPoints) {
		return fieldValidationError("connection_config.acquisition.limits.random_points", "random_points 超過所選設定檔的上限")
	}
	return err
}
