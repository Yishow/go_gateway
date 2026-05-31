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

func validateDeviceName(name string) error {
	if strings.TrimSpace(name) == "" {
		return validationError("設備名稱不能為空")
	}

	return nil
}
