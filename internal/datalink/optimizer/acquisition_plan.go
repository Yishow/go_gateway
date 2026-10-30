package optimizer

import (
	"go-gateway/internal/datalink/schema"
	"time"
)

// AcquisitionIdentity is the compatibility boundary; simultaneous due times do
// not erase distinct polling groups or cadences.
type AcquisitionIdentity struct {
	DeviceID, DeviceRevision, Station, ConnectionRevision string
	Protocol, Frame, ProtocolRevision, PollingGroup       string
	Cadence                                               time.Duration
}

// ApplicationFrameBudget counts the complete application frame, excluding
// network transport headers. Modbus RTU includes unit and CRC; FATEK includes
// ASCII framing/checksum. Zero means no declared budget in that direction.
type ApplicationFrameBudget struct{ RequestBytes, ResponseBytes int }

// NativeAddressBound is one optional authority declaration for one area.
// An absent declaration adds no bound; it is not proof of device readability.
type NativeAddressBound struct {
	Area                         string
	StartInclusive, EndExclusive int
}

// AcquisitionCapability is a revision-bound declaration, never a probing request.
// Zero limits are unspecified; at least one positive native-unit limit is required.
// Limits are additionally clamped to the selected command's actual client cap.
type AcquisitionCapability struct {
	Revision                                                        string
	RandomCrossArea                                                 bool   // Explicit support for mixed-area random commands.
	RandomSupport                                                   string // "supported", "unsupported", or empty/"unknown"
	ProtocolLimit, ClientLimit, DeviceLimit                         int
	ProtocolFrameBudget, ClientFrameBudget, DeviceFrameBudget       ApplicationFrameBudget
	ProtocolAddressBounds, ClientAddressBounds, DeviceAddressBounds []NativeAddressBound
}

// NativeRange is half-open and expressed in the area's native units.
type NativeRange struct {
	Area       string
	Start, End int
}

// ReadableCoverage is copied into the plan's fingerprint, not inferred from gaps.
// Excluded marks disabled/skipped gaps; active point spans can cover them.
// Forbidden is absolute, including active scalar spans.
type ReadableCoverage struct {
	Revision                      string
	Readable, Forbidden, Excluded []NativeRange
}

// EffectivePoint is an acquisition snapshot, not recording membership.
// NativeWidth must describe the complete value in protocol-native components.
type EffectivePoint struct {
	ID, SourceIdentity, SourceRevision, Address, Function string
	DataType                                              schema.DataType
	NativeWidth                                           int
	Format, FormatSource, FormatRevision                  string
}

// AcquisitionSource binds points to their transport, capability and coverage.
type AcquisitionSource struct {
	Identity   AcquisitionIdentity
	Capability AcquisitionCapability
	Coverage   ReadableCoverage
	Points     []EffectivePoint
}

// AcquisitionInput intentionally has no recording-selection field.
type AcquisitionInput struct {
	ConfigurationRevision string
	Sources               []AcquisitionSource
}

// AcquisitionGroup is the full request compatibility key.
type AcquisitionGroup struct {
	Identity                                             AcquisitionIdentity
	CapabilityRevision, CoverageRevision, Function, Area string
}

// PointSlice preserves a point's decoder policy independently of shared wire data.
type PointSlice struct {
	Point                                EffectivePoint
	Offset, Length, NativeBits, BitIndex int
}

// NativeUnit identifies one random-read component in wire order.
type NativeUnit struct {
	Area          string
	Address, Bits int
}

// AcquisitionRequest contains complete scalar slices in a safe contiguous span
// or an ordered random Units list. Random PointSlice.Offset indexes that list.
// Limit is the effective native-unit cap, including known frame budgets.
// FrameCostKnown=false means no encoded-byte compliance claim can be made.
type AcquisitionRequest struct {
	Group                AcquisitionGroup
	Kind                 string
	Start, Length, Limit int
	FallbackReason       string
	FrameBytes           ApplicationFrameBudget
	FrameCostKnown       bool
	Points               []PointSlice
	Units                []NativeUnit
}

// AcquisitionPointError is explicit, stable and attributable to a source point.
type AcquisitionPointError struct {
	Identity      AcquisitionIdentity
	PointID, Code string
}

// AcquisitionPlan owns its storage. Accessors return defensive copies.
type AcquisitionPlan struct {
	fingerprint string
	requests    []AcquisitionRequest
	errors      []AcquisitionPointError
}

// Fingerprint identifies all revision-bound acquisition inputs.
func (p AcquisitionPlan) Fingerprint() string { return p.fingerprint }

// Requests returns detached request/decoder snapshots.
func (p AcquisitionPlan) Requests() []AcquisitionRequest {
	result := append([]AcquisitionRequest(nil), p.requests...)
	for i := range result {
		result[i].Points = append([]PointSlice(nil), result[i].Points...)
		result[i].Units = append([]NativeUnit(nil), result[i].Units...)
	}
	return result
}

// Errors returns detached point errors, including rejected points.
func (p AcquisitionPlan) Errors() []AcquisitionPointError {
	return append([]AcquisitionPointError(nil), p.errors...)
}
