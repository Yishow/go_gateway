package measurement

import (
	"fmt"
	"math/big"
)

// MeanBasis names the numeric basis a window mean was computed on.
type MeanBasis string

const (
	// MeanBasisExact means an exact sum and an integer count. The mean is a
	// decimal and never passes through float64.
	MeanBasisExact MeanBasis = "exact"
	// MeanBasisFloat64 means a float64 sum over a deterministic observation
	// order. The result claims no digits beyond that basis.
	MeanBasisFloat64 MeanBasis = "float64"
)

// MeanAggregationVersion identifies the reduction algorithm that produced a mean.
const MeanAggregationVersion = "arithmetic-mean-v1"

// MeanAggregationMethod names the reduction method.
const MeanAggregationMethod = "arithmetic-mean"

// MeanRoundingHalfEven is the rule used to project an exact mean to at most
// MeanScaleLimit fractional digits.
const MeanRoundingHalfEven = "half-even"

// MeanScaleLimit is the greatest number of fractional digits an exact mean carries.
const MeanScaleLimit = 18

// MeanResult is a typed arithmetic mean with the evidence behind it. Value is a
// decimal for an exact basis and a float64 for a float basis; it is never one
// source reading passed off as the window's value.
type MeanResult struct {
	// Version and Method identify the aggregation that produced this result.
	Version string
	Method  string
	// Value is the projected mean.
	Value ExactValue
	// Sum is the exact sum of the contributing observations.
	Sum ExactValue
	// Count is the number of contributing observations.
	Count int64
	// Basis says whether Sum and Value are exact or float64.
	Basis MeanBasis
	// Scale is the fractional digit count of Value (exact basis only).
	Scale int
	// Rounding is the projection rule (exact basis only).
	Rounding string
	// Rounded is true when Value was projected and is not the exact mean.
	Rounded bool
}

// Mean returns the arithmetic mean of values, which must be non-empty and share
// one numeric basis: exact integers/decimals or float64. Exact inputs use an
// exact sum and a decimal result projected half-even; float inputs are summed in
// the given order, which the caller supplies deterministically.
func Mean(values []ExactValue) (MeanResult, error) {
	invalid := func(reason string) (MeanResult, error) {
		return MeanResult{}, fmt.Errorf("%w: %s", ErrExactValueInvalid, reason)
	}
	if len(values) == 0 {
		return invalid("no observations")
	}
	floats, exact := 0, 0
	for _, value := range values {
		switch value.Type() {
		case ExactFloat64:
			floats++
		case ExactInt64, ExactUint64, ExactDecimal:
			exact++
		}
	}
	switch {
	case floats > 0 && exact > 0:
		return invalid("a float cannot be averaged with an exact value")
	case floats > 0:
		return meanFloatBasis(values)
	default:
		return meanExactBasis(values)
	}
}

// meanFloatBasis sums in the given order and keeps the float basis.
func meanFloatBasis(values []ExactValue) (MeanResult, error) {
	sum := values[0]
	if sum.Type() != ExactFloat64 {
		return MeanResult{}, fmt.Errorf("%w: a float observation is required", ErrExactValueInvalid)
	}
	for _, value := range values[1:] {
		next, err := Add(sum, value)
		if err != nil {
			return MeanResult{}, err
		}
		sum = next
	}
	total, ok := sum.Value().(float64)
	if !ok {
		return MeanResult{}, fmt.Errorf("%w: the float sum is not a float64", ErrExactValueInvalid)
	}
	mean, err := NewFloat64(total / float64(len(values)))
	if err != nil {
		return MeanResult{}, err
	}
	return MeanResult{
		Version: MeanAggregationVersion, Method: MeanAggregationMethod,
		Value: mean, Sum: sum, Count: int64(len(values)), Basis: MeanBasisFloat64,
	}, nil
}

// meanExactBasis sums exactly and projects the quotient to MeanScaleLimit
// fractional digits, rounding half-even.
func meanExactBasis(values []ExactValue) (MeanResult, error) {
	sum := values[0]
	switch sum.Type() {
	case ExactInt64, ExactUint64, ExactDecimal:
	default:
		return MeanResult{}, fmt.Errorf("%w: %s is not numeric", ErrExactValueInvalid, sum.Type())
	}
	for _, value := range values[1:] {
		next, err := Add(sum, value)
		if err != nil {
			return MeanResult{}, err
		}
		sum = next
	}
	total, ok := sum.Rat()
	if !ok {
		return MeanResult{}, fmt.Errorf("%w: the exact sum has no rational value", ErrExactValueInvalid)
	}
	mean := new(big.Rat).Quo(total, new(big.Rat).SetInt64(int64(len(values))))
	value, scale, rounded, err := projectHalfEven(mean, MeanScaleLimit)
	if err != nil {
		return MeanResult{}, err
	}
	return MeanResult{
		Version: MeanAggregationVersion, Method: MeanAggregationMethod,
		Value: value, Sum: sum, Count: int64(len(values)), Basis: MeanBasisExact,
		Scale: scale, Rounding: MeanRoundingHalfEven, Rounded: rounded,
	}, nil
}

// projectHalfEven renders r as a decimal with at most maxScale fractional digits.
// The quotient is rounded to the nearest digit; an exact tie goes to the even
// digit. A projected value keeps all maxScale digits so a rounded result is
// never presented as exact.
func projectHalfEven(r *big.Rat, maxScale int) (projected ExactValue, outputScale int, didRound bool, projectionErr error) {
	negative := r.Sign() < 0
	abs := new(big.Rat).Abs(r)
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(maxScale)), nil)
	scaled := new(big.Rat).Mul(abs, new(big.Rat).SetInt(factor))
	unscaled, remainder := new(big.Int).QuoRem(scaled.Num(), scaled.Denom(), new(big.Int))
	rounded := remainder.Sign() != 0
	switch new(big.Int).Lsh(remainder, 1).Cmp(scaled.Denom()) {
	case 1:
		unscaled.Add(unscaled, big.NewInt(1))
	case 0:
		if unscaled.Bit(0) == 1 {
			unscaled.Add(unscaled, big.NewInt(1))
		}
	}
	scale := maxScale
	if unscaled.Sign() == 0 {
		scale = 0
	} else if !rounded {
		// A rounded value keeps every limit digit so it is never read as exact;
		// an exact value keeps only its significant fraction digits.
		for scale > 0 {
			quotient, remainder := new(big.Int).QuoRem(unscaled, big.NewInt(10), new(big.Int))
			if remainder.Sign() != 0 {
				break
			}
			unscaled, scale = quotient, scale-1
		}
		if unscaled.Sign() == 0 {
			scale = 0
		}
	}
	digits := formatScaled(unscaled, scale)
	if negative && unscaled.Sign() != 0 {
		digits = "-" + digits
	}
	value, err := NewDecimal(digits)
	if err != nil {
		return ExactValue{}, 0, false, err
	}
	return value, scale, rounded, nil
}
