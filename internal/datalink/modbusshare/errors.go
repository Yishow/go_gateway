package modbusshare

import (
	"fmt"
)

// Typed error codes for Modbus Share lifecycle
const (
	ErrCodeHydrationRequired    = "modbus_share_hydration_required"
	ErrCodeDisabled             = "modbus_share_disabled"
	ErrCodeWorkspaceScope       = "modbus_share_workspace_scope"
	ErrCodeRangeCollision       = "modbus_share_range_collision"
	ErrCodeCapacityExceeded     = "modbus_share_capacity_exceeded"
	ErrCodeReconcileFailed      = "modbus_share_reconcile_failed"
	ErrCodeDirtyUnknown         = "modbus_share_dirty_unknown"
	ErrCodeRevisionConflict     = "modbus_share_revision_conflict"
	ErrCodeSaveIncomplete       = "modbus_share_save_incomplete"
	ErrCodeProjectionRequired   = "modbus_share_projection_required"
	ErrCodeListenerBindFailed   = "modbus_share_listener_bind_failed"
	ErrCodeInvalidGeometry      = "modbus_share_invalid_geometry"
	ErrCodeDeliveryFailed       = "modbus_share_delivery_failed"
	ErrCodeSettingsUpdateFailed = "settings_update_failed"
	ErrCodeStorageFailure       = "modbus_share_storage_failed"
	ErrCodeRevisionRequired     = "modbus_share_revision_required"
)

// Error represents an operator-safe typed error envelope for Modbus Share.
type Error struct {
	Code              string `json:"code"`
	Message           string `json:"message"`
	Retryable         bool   `json:"retryable"`
	WorkspaceRevision string `json:"workspace_revision,omitempty"`
	SettingsRevision  string `json:"settings_revision,omitempty"`
	DirtyState        string `json:"dirty_state,omitempty"`
	Action            string `json:"action,omitempty"`
}

// StorageError preserves the durable operation that failed for diagnostics.
type StorageError struct {
	Operation string
	Err       error
}

func (e *StorageError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("modbus share storage %s: %v", e.Operation, e.Err)
}

func (e *StorageError) Unwrap() error { return e.Err }

func newStorageError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return &StorageError{Operation: operation, Err: err}
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewError creates a new Modbus Share Error.
func NewError(code, message string, retryable bool) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Retryable: retryable,
	}
}
