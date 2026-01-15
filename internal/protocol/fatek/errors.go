package fatek

import "fmt"

// FatekException 是所有 FATEK 驅動錯誤的基礎介面
type FatekException interface {
	error
}

// FatekError 是基礎錯誤結構
type FatekError struct {
	message string
}

func (e *FatekError) Error() string {
	return e.message
}

// NewFatekError 建立新的 FATEK 錯誤
func NewFatekError(format string, args ...interface{}) *FatekError {
	return &FatekError{
		message: fmt.Sprintf(format, args...),
	}
}

// FatekCommunicationError 表示底層通訊失敗
//
// 可能的原因包括：
//   - 連線逾時
//   - Checksum (LRC) 不匹配
//   - 格式錯誤的訊框 (缺少 STX/ETX)
//   - 連線中斷
type FatekCommunicationError struct {
	*FatekError
}

// NewFatekCommunicationError 建立新的通訊錯誤
func NewFatekCommunicationError(format string, args ...interface{}) *FatekCommunicationError {
	return &FatekCommunicationError{
		FatekError: NewFatekError(format, args...),
	}
}

// FatekProtocolError 表示 PLC 收到命令但回傳邏輯錯誤碼
//
// Attributes:
//   - ErrorCode: PLC 回傳的錯誤字元 ('2', '4', 'A' 等)
//   - Command: 導致錯誤的命令 ID (例如 '44')
type FatekProtocolError struct {
	*FatekError
	ErrorCode string
	Command   string
}

// NewFatekProtocolError 建立新的協定錯誤
func NewFatekProtocolError(errorCode, command string) *FatekProtocolError {
	var msg string
	switch errorCode {
	case "2":
		msg = fmt.Sprintf("PLC Error %s on Command %s: Illegal Value (Value out of range or non-hex)", errorCode, command)
	case "4":
		msg = fmt.Sprintf("PLC Error %s on Command %s: Illegal Format/Command (Unsupported command or LRC error)", errorCode, command)
	case "A":
		msg = fmt.Sprintf("PLC Error %s on Command %s: Illegal Address (Address boundary exceeded)", errorCode, command)
	default:
		msg = fmt.Sprintf("PLC Error %s on Command %s", errorCode, command)
	}

	return &FatekProtocolError{
		FatekError: &FatekError{message: msg},
		ErrorCode:  errorCode,
		Command:    command,
	}
}

// NewProtocolError 建立協定錯誤（統一介面）
// 此函數為 NewFatekProtocolError 的別名，用於統一 API
func NewProtocolError(code, cmd string) *FatekProtocolError {
	return NewFatekProtocolError(code, cmd)
}
