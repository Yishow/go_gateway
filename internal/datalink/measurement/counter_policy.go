package measurement

import (
	"errors"
	"fmt"
	"strings"
)

// maxExactPolicyText bounds the decimal text of the exact rollover evidence, so
// a pathological value cannot make later exact arithmetic expensive. It is well
// above a 128-bit modulus.
const maxExactPolicyText = 64

// ErrCounterPolicyInvalid marks exact rollover evidence that cannot be used.
var ErrCounterPolicyInvalid = errors.New("counter policy invalid")

// ExactRollover parses the exact rollover evidence. Both results are nil when
// the policy carries none. The legacy float MaxRollOver is deliberately not a
// source: it is not exact and has no maximum rate beside it.
func (p *CounterPolicy) ExactRollover() (modulus, maxRate *ExactValue, err error) {
	if p == nil || (p.RolloverModulus == "" && p.MaxRatePerSecond == "") {
		return nil, nil, nil
	}
	invalid := func(field string) error {
		return fmt.Errorf("%w: %s must be a positive exact decimal", ErrCounterPolicyInvalid, field)
	}
	parse := func(field, text string) (*ExactValue, error) {
		if text == "" {
			return nil, nil
		}
		if len(text) > maxExactPolicyText {
			return nil, invalid(field)
		}
		value, err := NewDecimal(text)
		if err != nil {
			return nil, invalid(field)
		}
		if rat, exact := value.Rat(); !exact || rat.Sign() <= 0 {
			return nil, invalid(field)
		}
		return &value, nil
	}
	if modulus, err = parse("rollover_modulus", p.RolloverModulus); err != nil {
		return nil, nil, err
	}
	if maxRate, err = parse("max_rate_per_second", p.MaxRatePerSecond); err != nil {
		return nil, nil, err
	}
	if modulus == nil {
		return nil, nil, fmt.Errorf("%w: max_rate_per_second needs the rollover_modulus it protects", ErrCounterPolicyInvalid)
	}
	return modulus, maxRate, nil
}

// validateExactRollover checks the exact evidence against the definition it
// belongs to: only a monotonic counter wraps, a counter that allows negative
// differences contradicts a wrap, and a float source cannot prove one.
func (p *CounterPolicy) validateExactRollover(kind SemanticKind, numericEncoding string) error {
	modulus, _, err := p.ExactRollover()
	if err != nil {
		return err
	}
	if modulus == nil {
		return nil
	}
	if kind != SemanticKindCounter {
		return fmt.Errorf("%w: rollover applies to a counter only", ErrCounterPolicyInvalid)
	}
	if p.AllowNegative {
		return fmt.Errorf("%w: a counter that allows negative differences has no rollover", ErrCounterPolicyInvalid)
	}
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(numericEncoding)), "float") {
		return fmt.Errorf("%w: a float source cannot prove a rollover", ErrCounterPolicyInvalid)
	}
	return nil
}

// sameExactText reports whether two exact policy texts are the same number: the
// same value written with another scale is the same evidence.
func sameExactText(a, b string) bool {
	if a == b {
		return true
	}
	if a == "" || b == "" {
		return false
	}
	left, leftErr := NewDecimal(a)
	right, rightErr := NewDecimal(b)
	return leftErr == nil && rightErr == nil && left.Equal(right)
}
