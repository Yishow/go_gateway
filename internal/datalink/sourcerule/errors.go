package sourcerule

import (
	"errors"
	"fmt"
)

var ErrSourceRuleNotFound = errors.New("source rule not found")
var ErrValidation = errors.New("source rule validation failed")

// FieldError reports a request field that failed validation. ItemIDs names
// the mixed layout items involved, such as both sides of an overlap.
type FieldError struct {
	Field   string
	ItemIDs []string
	Message string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s: %s", ErrValidation, e.Field, e.Message)
}

// Unwrap lets callers classify a FieldError as ErrValidation.
func (e *FieldError) Unwrap() error { return ErrValidation }

func fieldError(field, message string, itemIDs ...string) error {
	return &FieldError{Field: field, ItemIDs: itemIDs, Message: message}
}
