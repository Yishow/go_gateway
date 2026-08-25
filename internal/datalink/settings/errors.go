package settings

import (
	"errors"
	"fmt"
)

// ErrRevisionConflict identifies an expected-revision mismatch, not a storage failure.
var ErrRevisionConflict = errors.New("settings revision conflict")

// ErrSettingNotFound identifies an absent durable setting.
var ErrSettingNotFound = errors.New("setting not found")

// StorageError identifies a durable settings operation failure.
type StorageError struct {
	Operation string
	Err       error
}

func (e *StorageError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("settings storage %s: %v", e.Operation, e.Err)
}

func (e *StorageError) Unwrap() error { return e.Err }

func storageError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return &StorageError{Operation: operation, Err: err}
}
