// Package decoding describes decoder policy independently of transport state.
// Resolving metadata does not install a converter or change persisted intent.
package decoding

import (
	"errors"
	"strings"

	"go-gateway/internal/datalink/schema"
)

// Policy identifies a legacy algorithm, not a frozen resolved format.
type Policy string

const (
	// LegacyModbus preserves point-to-connection-to-ABCD fallback.
	LegacyModbus Policy = "legacy_modbus"
	// LegacyMC3E preserves connection-to-CDAB fallback, ignoring point format.
	LegacyMC3E Policy = "legacy_mc3e"
	// LegacyFatek preserves native address-component-width assembly.
	LegacyFatek Policy = "legacy_fatek_components"
)

// Source identifies the input actually used by the legacy decoder.
type Source string

const (
	// SourcePoint means the stored point format was selected.
	SourcePoint Source = "point"
	// SourceConnection means the connection format was selected.
	SourceConnection Source = "connection"
	// SourceNative means the protocol default or component assembly was selected.
	SourceNative Source = "native"
)

// LegacyInput is a read-only snapshot. Empty Mode means absent old metadata;
// it is never inferred from PointFormat. Nonlegacy modes are not handled here.
type LegacyInput struct {
	Protocol         schema.ProtocolType
	Mode             string
	PointFormat      string
	ConnectionFormat string
}

// Resolution separates policy identity from its current projection. FATEK has
// no scalar EffectiveFormat: NativeCompatibility selects component assembly.
// This is not a converter, applicability decision, or persistence record.
type Resolution struct {
	Mode                string
	Policy              Policy
	EffectiveFormat     string
	Source              Source
	NativeCompatibility bool
}

var (
	// ErrUnsupportedMode prevents silently treating new or corrupt intent as legacy.
	ErrUnsupportedMode = errors.New("unsupported decoding mode")
	// ErrUnsupportedProtocol indicates that no legacy policy exists.
	ErrUnsupportedProtocol = errors.New("unsupported decoding protocol")
)

// ResolveLegacy resolves legacy policy without mutating settings or converters.
func ResolveLegacy(input LegacyInput) (Resolution, error) {
	if input.Mode != "" && input.Mode != "legacy" {
		return Resolution{}, ErrUnsupportedMode
	}
	result := Resolution{Mode: "legacy", Source: SourceNative}
	switch input.Protocol {
	case schema.ProtocolModbusTCP, schema.ProtocolModbusRTU, schema.ProtocolModbusUDP:
		result.Policy = LegacyModbus
		result.EffectiveFormat = "ABCD"
		if strings.TrimSpace(input.PointFormat) != "" {
			result.Source, result.EffectiveFormat = SourcePoint, input.PointFormat
		} else if strings.TrimSpace(input.ConnectionFormat) != "" {
			result.Source, result.EffectiveFormat = SourceConnection, input.ConnectionFormat
		}
		// Modbus historically normalizes its selected input.
		result.EffectiveFormat = strings.ToUpper(strings.TrimSpace(result.EffectiveFormat))
	case schema.ProtocolMC3E:
		result.Policy = LegacyMC3E
		result.EffectiveFormat = "CDAB"
		if input.ConnectionFormat != "" {
			result.Source, result.EffectiveFormat = SourceConnection, input.ConnectionFormat
		}
		// Do not introduce Modbus normalization into MC's exact-string policy.
	case schema.ProtocolFatekFBs:
		result.Policy = LegacyFatek
		result.NativeCompatibility = true
		return result, nil
	default:
		return Resolution{}, ErrUnsupportedProtocol
	}
	switch result.EffectiveFormat {
	case "ABCD", "BADC", "CDAB", "DCBA":
		return result, nil
	default:
		// Legacy adapters accepted unknown formats and decoded them as ABCD.
		// Preserve that fallback and its selected source, not the invalid string.
		// Strict validation belongs at the boundary for new explicit intent.
		result.EffectiveFormat = "ABCD"
		return result, nil
	}
}
