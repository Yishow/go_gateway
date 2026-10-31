package measurement

import (
	"encoding/json"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// Preserve old identities for actual legacy installations. Modern identities
// require the very profile consumed by the read, never a new default lookup.
func sourceFormatIdentity(device schema.Device, point schema.Point, profiles []connector.ReadProfile) (string, error) {
	var profile connector.ReadProfile
	if len(profiles) > 0 {
		profile = profiles[0]
	}
	if err := connector.ValidatePointReadProfile(&point, profile); err != nil {
		return "", err
	}
	if !profile.Present() || profile.Resolution().Mode == schema.PointFormatLegacy {
		return "", nil
	}
	if !profile.MatchesConnection(device.Protocol, device.ConnectionConfig) || profile.DeviceID() != device.ID {
		return "", connector.ErrReadProfile
	}
	payload, err := json.Marshal(profile.IdentityView())
	if err != nil {
		return "", err
	}
	return string(payload), nil
}
