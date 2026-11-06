package optimizer

import (
	"errors"

	"go-gateway/internal/datalink/schema"
)

// ValidateAcquisitionCoverage uses the planner's canonical area and physical
// half-open range authority without building a plan or performing any I/O.
func ValidateAcquisitionCoverage(protocol string, coverage ReadableCoverage) bool {
	return validCoverage(protocol, coverage)
}

// AcquisitionPointShape resolves the same native scalar shape used by planning.
// It does not authorize coverage, apply capabilities or change source membership.
func AcquisitionPointShape(protocol string, point EffectivePoint) (nativeBits, nativeWidth int, err error) {
	return AcquisitionPointShapeForIdentity(AcquisitionIdentity{Protocol: protocol}, point)
}

// AcquisitionPointShapeForIdentity uses the same declared protocol profile as
// the planner, including the Q/iQ-R X/Y address radix.
func AcquisitionPointShapeForIdentity(identity AcquisitionIdentity, point EffectivePoint) (nativeBits, nativeWidth int, err error) {
	resolved := pointAddressForIdentity(identity, point)
	area, start, nativeBits, bitIndex, code := resolved.area, resolved.start, resolved.nativeBits, resolved.bitIndex, resolved.code
	if code != "" {
		return 0, 0, errors.New(code)
	}
	bits := scalarBits(point.DataType)
	if bits == 0 || ((nativeBits == 1 || bitIndex >= 0) && point.DataType != schema.DataTypeBool) {
		return 0, 0, errors.New("unsupported-type")
	}
	nativeWidth = (bits + nativeBits - 1) / nativeBits
	if !validNativeSpan(identity.Protocol, area, start, nativeWidth) {
		return 0, 0, errors.New(acquisitionReasonInvalidAddress)
	}
	return nativeBits, nativeWidth, nil
}

// AcquisitionPointSpan returns the complete physical half-open scalar span using
// the same address/function, native width and address-step authority as planning.
func AcquisitionPointSpan(protocol string, point EffectivePoint) (NativeRange, error) {
	return AcquisitionPointSpanForIdentity(AcquisitionIdentity{Protocol: protocol}, point)
}

// AcquisitionPointSpanForIdentity resolves a physical span using the same
// profile-bound address interpretation as acquisition planning.
func AcquisitionPointSpanForIdentity(identity AcquisitionIdentity, point EffectivePoint) (NativeRange, error) {
	_, width, err := AcquisitionPointShapeForIdentity(identity, point)
	if err != nil {
		return NativeRange{}, err
	}
	resolved := pointAddressForIdentity(identity, point)
	area, start := resolved.area, resolved.start
	step := nativeAddressStep(identity.Protocol, area)
	if step <= 0 || width > (int(^uint(0)>>1)-start)/step {
		return NativeRange{}, errors.New(acquisitionReasonInvalidAddress)
	}
	span := NativeRange{Area: area, Start: start, End: start + width*step}
	if !validNativeRange(identity.Protocol, area, span.Start, span.End) {
		return NativeRange{}, errors.New(acquisitionReasonInvalidAddress)
	}
	return span, nil
}
