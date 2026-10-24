package measurement

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Subtract returns to - from without a float64 detour. int64, uint64 and
// decimal operands may be mixed; the result is a decimal that keeps the finer
// of the two scales, so the digits shown never claim more precision than a
// reading carried. Two float64 operands give a float64 and a float never mixes
// with an exact operand: the answer would claim a precision one side lacks.
// Anything else, or a float result that overflows, is ErrExactValueInvalid.
func Subtract(to, from ExactValue) (ExactValue, error) { return combine(to, from, -1) }

// Add returns a + b under the same rules as Subtract.
func Add(a, b ExactValue) (ExactValue, error) { return combine(a, b, 1) }

func combine(a, b ExactValue, sign int) (ExactValue, error) {
	invalid := func(reason string) (ExactValue, error) {
		return ExactValue{}, fmt.Errorf("%w: %s", ErrExactValueInvalid, reason)
	}
	aFloat, aIsFloat := a.val.(float64)
	bFloat, bIsFloat := b.val.(float64)
	if aIsFloat || bIsFloat {
		if !aIsFloat || !bIsFloat {
			return invalid("a float cannot be combined with an exact value")
		}
		result := aFloat + float64(sign)*bFloat
		if math.IsInf(result, 0) || math.IsNaN(result) {
			return invalid("float result is not finite")
		}
		return NewFloat64(result)
	}
	left, ok := scaledOf(a)
	if !ok {
		return invalid("first operand is not numeric")
	}
	right, ok := scaledOf(b)
	if !ok {
		return invalid("second operand is not numeric")
	}
	scale := max(left.scale, right.scale)
	result := new(big.Int)
	if sign < 0 {
		result.Sub(left.at(scale), right.at(scale))
	} else {
		result.Add(left.at(scale), right.at(scale))
	}
	return NewDecimal(formatScaled(result, scale))
}

// Rat returns the exact rational value of an integer or decimal. A float64 is
// not exact and has none, so a caller needing exactness refuses it.
func (v ExactValue) Rat() (*big.Rat, bool) {
	scaled, ok := scaledOf(v)
	if !ok {
		return nil, false
	}
	denominator := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scaled.scale)), nil)
	return new(big.Rat).SetFrac(new(big.Int).Set(scaled.unscaled), denominator), true
}

// CanonicalDigits is the normalized decimal text of an integer or decimal, equal
// for equal numbers (no redundant zeros), for use in identities. A float64 has
// none.
func (v ExactValue) CanonicalDigits() (string, bool) {
	digits, ok := v.DecimalDigits()
	if !ok {
		return "", false
	}
	return canonicalDecimal(digits)
}

// Sign reports -1, 0 or 1 for a numeric value. Non-numeric and unset values
// have no sign.
func (v ExactValue) Sign() (int, bool) {
	if f, isFloat := v.val.(float64); isFloat {
		switch {
		case f > 0:
			return 1, true
		case f < 0:
			return -1, true
		}
		return 0, true
	}
	scaled, ok := scaledOf(v)
	if !ok {
		return 0, false
	}
	return scaled.unscaled.Sign(), true
}

// scaledNumber is an exact base-10 number: unscaled / 10^scale.
type scaledNumber struct {
	unscaled *big.Int
	scale    int
}

func scaledOf(v ExactValue) (scaledNumber, bool) {
	switch x := v.val.(type) {
	case int64:
		return scaledNumber{unscaled: big.NewInt(x)}, true
	case uint64:
		return scaledNumber{unscaled: new(big.Int).SetUint64(x)}, true
	case Decimal:
		negative, intPart, fracPart, err := splitDecimal(string(x))
		if err != nil {
			return scaledNumber{}, false
		}
		// splitDecimal guarantees at least one digit and digits only.
		unscaled, ok := new(big.Int).SetString(intPart+fracPart, 10)
		if !ok {
			return scaledNumber{}, false
		}
		if negative {
			unscaled.Neg(unscaled)
		}
		return scaledNumber{unscaled: unscaled, scale: len(fracPart)}, true
	}
	return scaledNumber{}, false
}

// at returns the unscaled integer expressed at a finer or equal scale.
func (n scaledNumber) at(scale int) *big.Int {
	if scale == n.scale {
		return n.unscaled
	}
	factor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(scale-n.scale)), nil)
	return new(big.Int).Mul(n.unscaled, factor)
}

// formatScaled renders unscaled / 10^scale with exactly scale fraction digits
// and no sign on zero.
func formatScaled(unscaled *big.Int, scale int) string {
	digits := new(big.Int).Abs(unscaled).String()
	if scale > 0 {
		if len(digits) <= scale {
			digits = strings.Repeat("0", scale-len(digits)+1) + digits
		}
		digits = digits[:len(digits)-scale] + "." + digits[len(digits)-scale:]
	}
	if unscaled.Sign() < 0 {
		return "-" + digits
	}
	return digits
}

// FloatDecimal renders a finite float64 as the shortest plain decimal that
// reads back to the same float64, never in exponent form. It carries the
// float's own precision and no more; callers keep the float basis alongside.
func FloatDecimal(v float64) (ExactValue, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return ExactValue{}, fmt.Errorf("%w: float64 must be finite", ErrExactValueInvalid)
	}
	return NewDecimal(strconv.FormatFloat(v, 'f', -1, 64))
}
