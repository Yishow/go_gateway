package device

import (
	"errors"
	"fmt"
	"strings"
)

var ErrValidation = errors.New("device validation failed")

func validationError(message string) error {
	return fmt.Errorf("%w: %s", ErrValidation, message)
}

// FieldError 標示驗證失敗的請求欄位；errors.Is(err, ErrValidation) 仍成立。
type FieldError struct {
	Field   string
	Message string
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", ErrValidation, e.Message)
}

// Unwrap 讓 FieldError 可被識別為 ErrValidation。
func (e *FieldError) Unwrap() error { return ErrValidation }

func fieldValidationError(field, message string) error {
	return &FieldError{Field: field, Message: message}
}

func connectionConfigValidationError(err error) error {
	var fieldErr *FieldError
	if errors.As(err, &fieldErr) {
		return fieldValidationError(fieldErr.Field, "連線配置無效: "+fieldErr.Message)
	}
	return validationError(fmt.Sprintf("連線配置無效: %v", err))
}

func validateDeviceName(name string) error {
	if strings.TrimSpace(name) == "" {
		return validationError("設備名稱不能為空")
	}

	return nil
}
