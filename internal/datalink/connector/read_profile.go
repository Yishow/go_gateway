package connector

import (
	"errors"

	"go-gateway/internal/datalink/decoding"
	"go-gateway/internal/datalink/schema"
)

// ReadProfile is a detached, immutable decoding decision for one logical point.
// Raw batch requests do not carry it: each scatter slice retains its own profile.
// Installation leases and acquisition proof belong to the collector envelope.
type ReadProfile struct {
	version    uint8
	protocol   schema.ProtocolType
	deviceID   string
	pointID    string
	address    string
	function   string
	dataType   schema.DataType
	resolution decoding.PointResolution
}

var ErrReadProfile = errors.New("invalid or mismatched point read profile")

// ResolveReadProfile never returns a usable profile when resolution fails.
// Callers must fail closed on error, not install the zero result as legacy.
func ResolveReadProfile(point schema.Point, defaults decoding.FormatContext) (ReadProfile, error) {
	resolution, err := decoding.ResolvePoint(point, defaults)
	if err != nil {
		return ReadProfile{}, err
	}
	if point.ID == "" || !resolution.Resolved {
		return ReadProfile{}, ErrReadProfile
	}
	return ReadProfile{version: 1, protocol: defaults.Protocol, pointID: point.ID, deviceID: point.DeviceID, address: point.Address, function: point.Function, dataType: point.DataType, resolution: resolution}, nil
}

// Present distinguishes an explicit profile from the historical request API.
func (p ReadProfile) Present() bool                        { return p != (ReadProfile{}) }
func (p ReadProfile) Protocol() schema.ProtocolType        { return p.protocol }
func (p ReadProfile) DeviceID() string                     { return p.deviceID }
func (p ReadProfile) PointID() string                      { return p.pointID }
func (p ReadProfile) Resolution() decoding.PointResolution { return p.resolution }

// ValidateReadProfile compares logical scalar identity, never batch wire layout.
// Empty profiles are only the compatibility path for existing legacy callers.
func ValidateReadProfile(req ReadRequest) error {
	p := req.Profile
	if !p.Present() {
		return nil
	}
	if p.version != 1 || !p.resolution.Resolved || p.pointID == "" || p.address != req.Address || p.function != req.Function || p.dataType != req.DataType {
		return ErrReadProfile
	}
	return nil
}

// ValidatePointReadProfile prevents a detached policy from being attached to a
// different configuration revision at installation. Legacy callers need no profile.
func ValidatePointReadProfile(point *schema.Point, profile ReadProfile) error {
	if point == nil {
		return ErrReadProfile
	}
	mode := point.FormatMode
	if mode == "" {
		mode = schema.PointFormatLegacy
	}
	if !profile.Present() {
		if mode == schema.PointFormatLegacy {
			return nil
		}
		return ErrReadProfile
	}
	if err := ValidateReadProfile(ReadRequest{Address: point.Address, Function: point.Function, DataType: point.DataType, Profile: profile}); err != nil {
		return err
	}
	r := profile.Resolution()
	if profile.PointID() != point.ID || profile.DeviceID() != point.DeviceID || r.PointRevision != point.RevisionID || r.Mode != mode || r.ConfiguredFormat != point.FormatOverride || r.FormatVersion != point.FormatVersion {
		return ErrReadProfile
	}
	return nil
}
