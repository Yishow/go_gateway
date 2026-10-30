package decoding

import (
	"errors"
	"go-gateway/internal/datalink/schema"
	"math"
)

// ErrShortScalar rejects incomplete samples before a numeric value is produced.
var ErrShortScalar = errors.New("incomplete scalar response; check point type and address width")

// DecodeScalarWords decodes one complete scalar without mutating raw words or a
// shared converter. Nil on error cannot masquerade as a good numeric zero.
// Single-register legacy/inherit reads deliberately retain their driver path.
func DecodeScalarWords(words []uint16, dtype schema.DataType, format string) (any, error) {
	if !canonicalFormat(format) {
		return nil, ErrInvalidIntent
	}
	count := scalarBits(dtype) / 16
	if count == 0 {
		return nil, ErrInapplicableFormat
	}
	if len(words) < count {
		return nil, ErrShortScalar
	}
	var bits uint64
	for i := 0; i < count; i++ {
		index := i
		if format == "CDAB" || format == "DCBA" {
			index = count - 1 - i
		}
		word := words[index]
		if format == "BADC" || format == "DCBA" {
			word = word>>8 | word<<8
		}
		bits = bits<<16 | uint64(word)
	}
	switch dtype {
	case schema.DataTypeInt32:
		return int32(uint32(bits)), nil // #nosec G115 -- Reinterpret fixed-width wire bits.
	case schema.DataTypeUint32:
		return uint32(bits), nil // #nosec G115 -- Exactly two words were assembled.
	case schema.DataTypeFloat32:
		return math.Float32frombits(uint32(bits)), nil // #nosec G115 -- Exactly two words were assembled.
	case schema.DataTypeInt64:
		return int64(bits), nil // #nosec G115 -- Reinterpret fixed-width wire bits.
	case schema.DataTypeUint64:
		return bits, nil
	case schema.DataTypeFloat64:
		return math.Float64frombits(bits), nil
	default:
		return nil, ErrInapplicableFormat
	}
}

// NormalizeFatekComponents extracts high-to-low 16-bit words from the existing
// native components. It neither changes address strides nor doubles read count.
func NormalizeFatekComponents(values []int, nativeBits int, dtype schema.DataType) ([]uint16, error) {
	bits := scalarBits(dtype)
	if bits == 0 || (nativeBits != 16 && nativeBits != 32) {
		return nil, ErrInapplicableFormat
	}
	count := bits / nativeBits
	if len(values) < count {
		return nil, ErrShortScalar
	}
	words := make([]uint16, 0, bits/16)
	for _, value := range values[:count] {
		raw := uint64(value) // #nosec G115 -- Preserve signed native component bit patterns.
		if nativeBits == 32 {
			words = append(words, uint16(raw>>16))
		} // #nosec G115 -- Extract high word.
		words = append(words, uint16(raw)) // #nosec G115 -- Extract low word.
	}
	return words, nil
}
