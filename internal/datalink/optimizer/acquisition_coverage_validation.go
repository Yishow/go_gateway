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
	area, start, nativeBits, bitIndex, _, code := pointNativeAddress(protocol, point)
	if code != "" {
		return 0, 0, errors.New(code)
	}
	bits := scalarBits(point.DataType)
	if bits == 0 || ((nativeBits == 1 || bitIndex >= 0) && point.DataType != schema.DataTypeBool) {
		return 0, 0, errors.New("unsupported-type")
	}
	nativeWidth = (bits + nativeBits - 1) / nativeBits
	if !validNativeSpan(protocol, area, start, nativeWidth) {
		return 0, 0, errors.New(acquisitionReasonInvalidAddress)
	}
	return nativeBits, nativeWidth, nil
}

// AcquisitionPointSpan returns the complete physical half-open scalar span using
// the same address/function, native width and address-step authority as planning.
func AcquisitionPointSpan(protocol string, point EffectivePoint) (NativeRange, error) {
	_, width, err := AcquisitionPointShape(protocol, point)
	if err != nil {
		return NativeRange{}, err
	}
	area, start, _, _, _, _ := pointNativeAddress(protocol, point)
	step := nativeAddressStep(protocol, area)
	if step <= 0 || width > (int(^uint(0)>>1)-start)/step {
		return NativeRange{}, errors.New(acquisitionReasonInvalidAddress)
	}
	span := NativeRange{Area: area, Start: start, End: start + width*step}
	if !validNativeRange(protocol, area, span.Start, span.End) {
		return NativeRange{}, errors.New(acquisitionReasonInvalidAddress)
	}
	return span, nil
}
