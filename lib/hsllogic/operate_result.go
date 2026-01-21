package hsllogic

// =============================================================================
// OperateResult 統一結果封裝 (仿 HslCommunication OperateResult<T>)
// =============================================================================

// OperateResult 統一的操作結果封裝
// 用於標準化所有協議操作的返回值
type OperateResult[T any] struct {
	// IsSuccess 操作是否成功
	IsSuccess bool `json:"is_success"`

	// ErrorCode 錯誤代碼 (0 = 成功)
	ErrorCode int `json:"error_code"`

	// Message 錯誤或狀態訊息
	Message string `json:"message,omitempty"`

	// Content 操作結果內容
	Content T `json:"content,omitempty"`
}

// =============================================================================
// 標準錯誤代碼定義
// =============================================================================

const (
	// ErrCodeSuccess 成功
	ErrCodeSuccess = 0
	// ErrCodeConnectionFailed 連線失敗
	ErrCodeConnectionFailed = 1001
	// ErrCodeTimeout 超時
	ErrCodeTimeout = 1002
	// ErrCodeDisconnected 連線中斷
	ErrCodeDisconnected = 1003
	// ErrCodeInvalidAddress 無效地址
	ErrCodeInvalidAddress = 2001
	// ErrCodeInvalidDataType 無效數據類型
	ErrCodeInvalidDataType = 2002
	// ErrCodeInvalidValue 無效數值
	ErrCodeInvalidValue = 2003
	// ErrCodeReadFailed 讀取失敗
	ErrCodeReadFailed = 3001
	// ErrCodeWriteFailed 寫入失敗
	ErrCodeWriteFailed = 3002
	// ErrCodeDeviceError PLC/設備錯誤
	ErrCodeDeviceError = 4001
	// ErrCodeUnknown 未知錯誤
	ErrCodeUnknown = 9999
)

// =============================================================================
// 建構函數
// =============================================================================

// Success 建立成功結果
func Success[T any](content T) OperateResult[T] {
	return OperateResult[T]{
		IsSuccess: true,
		ErrorCode: ErrCodeSuccess,
		Message:   "操作成功",
		Content:   content,
	}
}

// SuccessWithMessage 建立成功結果 (自訂訊息)
func SuccessWithMessage[T any](content T, message string) OperateResult[T] {
	return OperateResult[T]{
		IsSuccess: true,
		ErrorCode: ErrCodeSuccess,
		Message:   message,
		Content:   content,
	}
}

// Fail 建立失敗結果
func Fail[T any](code int, message string) OperateResult[T] {
	var zero T
	return OperateResult[T]{
		IsSuccess: false,
		ErrorCode: code,
		Message:   message,
		Content:   zero,
	}
}

// FailFromError 從 error 建立失敗結果
func FailFromError[T any](err error) OperateResult[T] {
	var zero T
	if err == nil {
		return Success(zero)
	}
	return OperateResult[T]{
		IsSuccess: false,
		ErrorCode: ErrCodeUnknown,
		Message:   err.Error(),
		Content:   zero,
	}
}

// =============================================================================
// 輔助方法
// =============================================================================

// Then 鏈式操作：如果成功則執行 fn，否則返回失敗
func (r OperateResult[T]) Then(fn func(T) error) OperateResult[T] {
	if !r.IsSuccess {
		return r
	}
	if err := fn(r.Content); err != nil {
		return Fail[T](ErrCodeUnknown, err.Error())
	}
	return r
}

// Map 映射結果內容到新類型
func Map[T any, U any](r OperateResult[T], fn func(T) U) OperateResult[U] {
	if !r.IsSuccess {
		return Fail[U](r.ErrorCode, r.Message)
	}
	return Success(fn(r.Content))
}

// GetContentOrDefault 取得內容，如果失敗則返回預設值
func (r OperateResult[T]) GetContentOrDefault(defaultValue T) T {
	if !r.IsSuccess {
		return defaultValue
	}
	return r.Content
}

// =============================================================================
// 非泛型版本 (用於不需要返回值的操作)
// =============================================================================

// SimpleResult 簡單操作結果 (不包含內容)
type SimpleResult struct {
	IsSuccess bool   `json:"is_success"`
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message,omitempty"`
}

// SimpleSuccess 建立簡單成功結果
func SimpleSuccess() SimpleResult {
	return SimpleResult{
		IsSuccess: true,
		ErrorCode: ErrCodeSuccess,
		Message:   "操作成功",
	}
}

// SimpleFail 建立簡單失敗結果
func SimpleFail(code int, message string) SimpleResult {
	return SimpleResult{
		IsSuccess: false,
		ErrorCode: code,
		Message:   message,
	}
}

// SimpleFailFromError 從 error 建立簡單失敗結果
func SimpleFailFromError(err error) SimpleResult {
	if err == nil {
		return SimpleSuccess()
	}
	return SimpleResult{
		IsSuccess: false,
		ErrorCode: ErrCodeUnknown,
		Message:   err.Error(),
	}
}
