// Package protocol 提供工業協議通訊的共用定義
package protocol

import "fmt"

// ProtocolException 是所有協議錯誤的基礎介面
// 所有協議實作（Fatek、Modbus、MC Protocol）都應實作此介面
type ProtocolException interface {
	error
	// IsProtocolError 標識此錯誤為協議層級錯誤
	IsProtocolError() bool
}

// CommunicationException 表示通訊層級錯誤
// 包括：連線逾時、連線中斷、校驗碼錯誤等
type CommunicationException interface {
	error
	// IsCommunicationError 標識此錯誤為通訊層級錯誤
	IsCommunicationError() bool
}

// BaseError 提供基礎錯誤實作
type BaseError struct {
	Message string
	Cause   error
}

func (e *BaseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *BaseError) Unwrap() error {
	return e.Cause
}

// CommunicationError 表示通訊層級錯誤的具體實作
type CommunicationError struct {
	BaseError
}

func (e *CommunicationError) IsCommunicationError() bool {
	return true
}

// NewCommunicationError 建立新的通訊錯誤
func NewCommunicationError(message string, cause error) *CommunicationError {
	return &CommunicationError{
		BaseError: BaseError{Message: message, Cause: cause},
	}
}

// NewCommunicationErrorf 使用格式化字串建立通訊錯誤
func NewCommunicationErrorf(format string, args ...interface{}) *CommunicationError {
	return &CommunicationError{
		BaseError: BaseError{Message: fmt.Sprintf(format, args...)},
	}
}

// DeviceError 表示設備層級錯誤（PLC 回傳的錯誤碼）
type DeviceError struct {
	BaseError
	Code     int    // 設備錯誤碼
	Protocol string // 協議名稱（fatek, modbus, mc）
}

func (e *DeviceError) IsProtocolError() bool {
	return true
}

// NewDeviceError 建立新的設備錯誤
func NewDeviceError(protocol string, code int, message string) *DeviceError {
	return &DeviceError{
		BaseError: BaseError{Message: message},
		Code:      code,
		Protocol:  protocol,
	}
}

// 預定義錯誤類型
var (
	// ErrConnectionClosed 連線已關閉
	ErrConnectionClosed = NewCommunicationErrorf("connection closed")

	// ErrTimeout 操作逾時
	ErrTimeout = NewCommunicationErrorf("operation timeout")

	// ErrInvalidResponse 無效的回應
	ErrInvalidResponse = NewCommunicationErrorf("invalid response")

	// ErrChecksumMismatch 校驗碼不匹配
	ErrChecksumMismatch = NewCommunicationErrorf("checksum mismatch")
)

// IsConnectionClosed 檢查錯誤是否為連線關閉錯誤
func IsConnectionClosed(err error) bool {
	if err == nil {
		return false
	}
	// 檢查是否為預定義錯誤或包含相關訊息
	if err == ErrConnectionClosed {
		return true
	}
	// 檢查錯誤訊息（相容各協議的錯誤定義）
	errMsg := err.Error()
	return errMsg == "connection closed" ||
		errMsg == "連線已關閉" ||
		errMsg == ErrConnectionClosed.Error()
}

// IsTimeout 檢查錯誤是否為逾時錯誤
func IsTimeout(err error) bool {
	if err == nil {
		return false
	}
	if err == ErrTimeout {
		return true
	}
	errMsg := err.Error()
	return errMsg == "operation timeout" ||
		errMsg == "操作逾時" ||
		errMsg == ErrTimeout.Error()
}
