package device

import (
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

func TestValidateConnectionConfigRejectsBrokenRegisteredSchema(t *testing.T) {
	info, ok := connector.GetProtocolInfo(schema.ProtocolModbusTCP)
	if !ok || len(info.ConfigSchema) == 0 {
		t.Fatal("expected built-in Modbus schema")
	}
	// ProtocolInfo exposes the registered schema bytes. This serial test corrupts
	// and restores only those bytes to exercise a broken adapter definition.
	original := append([]byte(nil), info.ConfigSchema...)
	t.Cleanup(func() { copy(info.ConfigSchema, original) })
	info.ConfigSchema[0] = '!'
	if err := validateConnectionConfig(schema.ProtocolModbusTCP, map[string]interface{}{}); err == nil {
		t.Fatal("invalid registered schema must not bypass configuration validation")
	}
}
