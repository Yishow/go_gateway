package decoding

import (
	"errors"
	"fmt"
	"strings"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

// FormatContext supplies defaults and revisions without becoming point intent.
type FormatContext struct {
	Protocol           schema.ProtocolType `json:"protocol"`
	RuleOwned          bool                `json:"rule_owned"`
	RuleFormat         string              `json:"rule_format"`
	ConnectionFormat   string              `json:"connection_format"`
	WorkspaceRevision  string              `json:"workspace_revision"`
	RuleRevision       string              `json:"rule_revision"`
	ConnectionRevision string              `json:"connection_revision"`
}

// PointResolution is a read-only projection. An empty effective format with
// NativeCompatibility means component assembly, not an assumed ABCD override.
type PointResolution struct {
	Resolved            bool                   `json:"resolved"`
	Mode                schema.PointFormatMode `json:"format_mode"`
	ConfiguredFormat    string                 `json:"configured_format"`
	EffectiveFormat     string                 `json:"effective_format"`
	Source              Source                 `json:"effective_source"`
	Policy              Policy                 `json:"effective_policy"`
	NativeCompatibility bool                   `json:"native_compatibility"`
	Applicable          bool                   `json:"applicable"`
	FormatVersion       int                    `json:"format_version"`
	PointRevision       string                 `json:"point_revision"`
	WorkspaceRevision   string                 `json:"workspace_revision"`
	RuleRevision        string                 `json:"rule_revision"`
	ConnectionRevision  string                 `json:"connection_revision"`
}

// IntentEdit is an internal candidate edit, not an authorized public mutation.
// Nil fields omit intent changes. Selecting inherit explicitly resets override.
type IntentEdit struct {
	Mode             *schema.PointFormatMode
	ConfiguredFormat *string
}

// ErrInvalidIntent rejects ambiguous or malformed new intent without raw input.
var ErrInvalidIntent = errors.New("invalid point format intent")

// ResolvePoint describes intent and effective policy without changing a point.
func ResolvePoint(point schema.Point, defaults FormatContext) (PointResolution, error) {
	if point.FormatVersion < 0 || point.FormatVersion > schema.PointFormatVersion {
		return PointResolution{}, ErrInvalidIntent
	}
	mode := point.FormatMode
	if mode == "" {
		mode = schema.PointFormatLegacy
	}
	if mode != schema.PointFormatLegacy && mode != schema.PointFormatInherit && mode != schema.PointFormatOverride {
		return PointResolution{}, ErrInvalidIntent
	}
	if mode != schema.PointFormatOverride && point.FormatOverride != "" {
		return PointResolution{}, ErrInvalidIntent
	}
	_, capabilityErr := pointFormatCapability(defaults.Protocol, point)
	applicable := capabilityErr == nil
	multiword := multiwordNumeric(point.DataType)
	legacyInput := LegacyInput{Protocol: defaults.Protocol, PointFormat: point.DataFormat, ConnectionFormat: defaults.ConnectionFormat}
	if mode != schema.PointFormatLegacy {
		legacyInput.PointFormat = ""
	}
	policy, err := ResolveLegacy(legacyInput)
	if err != nil {
		return PointResolution{}, err
	}
	if mode == schema.PointFormatOverride {
		if !canonicalFormat(point.FormatOverride) {
			return PointResolution{}, ErrInvalidIntent
		}
		if capabilityErr != nil {
			return PointResolution{}, capabilityErr
		}
		policy.EffectiveFormat, policy.Source, policy.Policy = point.FormatOverride, SourceOverride, PolicyOverride
		policy.NativeCompatibility = false
	} else if mode == schema.PointFormatInherit {
		if multiword && capabilityErr != nil {
			return PointResolution{}, capabilityErr
		}
		// Non-multiword MC/FATEK reads keep their native driver policy. Modbus
		// still carries its default metadata, without changing register16 semantics.
		useRule := multiword || policy.Policy == LegacyModbus
		if useRule && defaults.RuleOwned && strings.TrimSpace(defaults.RuleFormat) != "" {
			format := strings.ToUpper(strings.TrimSpace(defaults.RuleFormat))
			if !canonicalFormat(format) {
				return PointResolution{}, ErrInvalidIntent
			}
			policy.EffectiveFormat, policy.Source = format, SourceRule
			policy.NativeCompatibility = false
		} else if multiword && defaults.Protocol == schema.ProtocolFatekFBs && strings.TrimSpace(defaults.ConnectionFormat) != "" {
			format := strings.ToUpper(strings.TrimSpace(defaults.ConnectionFormat))
			if !canonicalFormat(format) {
				return PointResolution{}, ErrInvalidIntent
			}
			policy.EffectiveFormat, policy.Source = format, SourceConnection
			policy.NativeCompatibility = false
		}
		if !policy.NativeCompatibility {
			policy.Policy = PolicyInherit
		}
	}
	result := PointResolution{
		Resolved: true,
		Mode:     mode, ConfiguredFormat: point.FormatOverride,
		EffectiveFormat: policy.EffectiveFormat, Source: policy.Source, Policy: policy.Policy,
		NativeCompatibility: policy.NativeCompatibility, Applicable: applicable,
		FormatVersion: point.FormatVersion, PointRevision: point.RevisionID,
		WorkspaceRevision: defaults.WorkspaceRevision, ConnectionRevision: defaults.ConnectionRevision,
	}
	if defaults.RuleOwned {
		result.RuleRevision = defaults.RuleRevision
	}
	return result, nil
}

// WithIntent returns an edited candidate. Authorization and atomic persistence
// belong to the caller; it must not be exposed as an unguarded mutation endpoint.
func WithIntent(point schema.Point, edit IntentEdit) (schema.Point, error) {
	if edit.Mode == nil && edit.ConfiguredFormat == nil {
		return point, nil
	}
	if point.FormatVersion < 0 || point.FormatVersion > schema.PointFormatVersion {
		return point, ErrInvalidIntent
	}
	candidate := point
	if edit.Mode != nil {
		candidate.FormatMode = *edit.Mode
	}
	switch candidate.FormatMode {
	case schema.PointFormatInherit:
		if edit.ConfiguredFormat != nil && *edit.ConfiguredFormat != "" {
			return point, ErrInvalidIntent
		}
		candidate.FormatOverride = ""
	case schema.PointFormatOverride:
		if edit.ConfiguredFormat != nil {
			candidate.FormatOverride = *edit.ConfiguredFormat
		}
		if !multiwordNumeric(candidate.DataType) {
			return point, ErrInapplicableFormat
		}
		if !canonicalFormat(candidate.FormatOverride) {
			return point, ErrInvalidIntent
		}
	default:
		return point, ErrInvalidIntent
	}
	candidate.FormatVersion = schema.PointFormatVersion
	if candidate.FormatMode == point.FormatMode && candidate.FormatOverride == point.FormatOverride && candidate.FormatVersion == point.FormatVersion {
		return point, nil
	}
	revision, err := common.NewUUID()
	if err != nil {
		return point, fmt.Errorf("create point format revision: %w", err)
	}
	candidate.RevisionID = revision
	return candidate, nil
}

// SourceRule and SourceOverride distinguish inherited defaults from operator intent.
const (
	SourceRule     Source = "rule"
	SourceOverride Source = "override"
	// PolicyInherit resolves defaults on each read, rather than freezing a format.
	PolicyInherit Policy = "inherit_scalar"
	// PolicyOverride selects the point's explicit canonical scalar format.
	PolicyOverride Policy = "override_scalar"
)

func canonicalFormat(format string) bool {
	switch format {
	case "ABCD", "BADC", "CDAB", "DCBA":
		return true
	}
	return false
}

// multiwordNumeric is the type-only gate for context-free candidate edits.
// ResolvePoint also validates the protocol/address before mutation can be applied.
func multiwordNumeric(dataType schema.DataType) bool {
	switch dataType {
	case schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32,
		schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64:
		return true
	}
	return false
}
