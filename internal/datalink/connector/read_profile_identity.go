package connector

import (
	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
)

// ReadProfileIdentity is a detached, explicitly versioned identity view. It is
// not a mutable decoding policy: only ResolveReadProfile creates usable profiles.
// ConnectionRevision is an opaque digest; raw connection settings are absent.
type ReadProfileIdentity struct {
	Version    uint8                    `json:"version"`
	Protocol   schema.ProtocolType      `json:"protocol"`
	DeviceID   string                   `json:"device_id"`
	PointID    string                   `json:"point_id"`
	Address    string                   `json:"address"`
	Function   string                   `json:"function"`
	DataType   schema.DataType          `json:"data_type"`
	RuleID     string                   `json:"rule_id"`
	Resolution decoding.PointResolution `json:"resolution"`
}

func (p ReadProfile) IdentityView() ReadProfileIdentity {
	// Workspace revisions authorize edits but are not decoding identity.
	// In particular DatabaseSetupRevision and Share revisions must not be mixed.
	resolution := p.resolution
	resolution.WorkspaceRevision = ""
	return ReadProfileIdentity{Version: p.version, Protocol: p.protocol, DeviceID: p.deviceID, PointID: p.pointID, Address: p.address, Function: p.function, DataType: p.dataType, RuleID: p.ruleID, Resolution: resolution}
}
