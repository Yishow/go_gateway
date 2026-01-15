package modbus

import "fmt"

// ModbusError 是 Modbus 協議錯誤的基礎介面
type ModbusError interface {
	error
	Code() byte
}

// ProtocolError 表示 Modbus 協議層級錯誤
type ProtocolError struct {
	code    byte
	message string
}

func (e *ProtocolError) Error() string {
	return e.message
}

func (e *ProtocolError) Code() byte {
	return e.code
}

// NewProtocolError 建立新的協議錯誤
func NewProtocolError(code byte, format string, args ...interface{}) *ProtocolError {
	msg := fmt.Sprintf(format, args...)
	if exceptionMsg, ok := ExceptionMessages[code]; ok {
		msg = fmt.Sprintf("%s: %s", msg, exceptionMsg)
	}
	return &ProtocolError{
		code:    code,
		message: msg,
	}
}

// CommunicationError 表示通訊層級錯誤
type CommunicationError struct {
	message string
}

func (e *CommunicationError) Error() string {
	return e.message
}

// NewCommunicationError 建立新的通訊錯誤
func NewCommunicationError(format string, args ...interface{}) *CommunicationError {
	return &CommunicationError{
		message: fmt.Sprintf(format, args...),
	}
}

// 預定義錯誤
var (
	ErrConnectionClosed    = NewCommunicationError("連線已關閉")
	ErrTimeout             = NewCommunicationError("操作逾時")
	ErrInvalidResponse     = NewCommunicationError("無效的回應")
	ErrCRCError            = NewCommunicationError("CRC 校驗錯誤")
	ErrInvalidFrame        = NewCommunicationError("無效的封包格式")
	ErrInvalidFunctionCode  = NewCommunicationError("無效的功能碼")
	ErrInvalidAddress      = NewCommunicationError("無效的地址")
	ErrInvalidQuantity     = NewCommunicationError("無效的數量")
	ErrResponseTooShort    = NewCommunicationError("回應長度不足")
	ErrResponseTooLong     = NewCommunicationError("回應長度過長")
	ErrTransactionIDMismatch = NewCommunicationError("交易 ID 不匹配")
)
