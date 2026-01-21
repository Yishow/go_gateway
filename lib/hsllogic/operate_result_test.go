package hsllogic

import (
	"errors"
	"testing"
)

// =============================================================================
// OperateResult 測試 (目標覆蓋率 90%+)
// =============================================================================

func TestSuccess(t *testing.T) {
	result := Success(42)
	if !result.IsSuccess {
		t.Error("Success result should have IsSuccess = true")
	}
	if result.ErrorCode != ErrCodeSuccess {
		t.Errorf("Success result should have ErrorCode = 0, got %d", result.ErrorCode)
	}
	if result.Content != 42 {
		t.Errorf("Content should be 42, got %d", result.Content)
	}
}

func TestSuccessWithMessage(t *testing.T) {
	result := SuccessWithMessage("hello", "自訂訊息")
	if !result.IsSuccess {
		t.Error("SuccessWithMessage should have IsSuccess = true")
	}
	if result.Message != "自訂訊息" {
		t.Errorf("Message should be '自訂訊息', got '%s'", result.Message)
	}
	if result.Content != "hello" {
		t.Errorf("Content should be 'hello', got '%s'", result.Content)
	}
}

func TestFail(t *testing.T) {
	result := Fail[int](ErrCodeConnectionFailed, "連線失敗")
	if result.IsSuccess {
		t.Error("Fail result should have IsSuccess = false")
	}
	if result.ErrorCode != ErrCodeConnectionFailed {
		t.Errorf("ErrorCode should be %d, got %d", ErrCodeConnectionFailed, result.ErrorCode)
	}
	if result.Message != "連線失敗" {
		t.Errorf("Message should be '連線失敗', got '%s'", result.Message)
	}
	if result.Content != 0 {
		t.Errorf("Content should be zero value, got %d", result.Content)
	}
}

func TestFailFromError(t *testing.T) {
	err := errors.New("測試錯誤")
	result := FailFromError[string](err)
	if result.IsSuccess {
		t.Error("FailFromError should have IsSuccess = false")
	}
	if result.Message != "測試錯誤" {
		t.Errorf("Message should be '測試錯誤', got '%s'", result.Message)
	}
}

func TestFailFromError_NilError(t *testing.T) {
	result := FailFromError[int](nil)
	if !result.IsSuccess {
		t.Error("FailFromError with nil should return success")
	}
}

func TestThen_Success(t *testing.T) {
	result := Success(10)
	called := false
	newResult := result.Then(func(v int) error {
		called = true
		if v != 10 {
			t.Errorf("Then should receive 10, got %d", v)
		}
		return nil
	})
	if !called {
		t.Error("Then function should be called")
	}
	if !newResult.IsSuccess {
		t.Error("Then should return success when fn returns nil")
	}
}

func TestThen_Failure(t *testing.T) {
	result := Fail[int](ErrCodeTimeout, "超時")
	called := false
	newResult := result.Then(func(v int) error {
		called = true
		return nil
	})
	if called {
		t.Error("Then function should not be called on failure")
	}
	if newResult.IsSuccess {
		t.Error("Then should propagate failure")
	}
}

func TestThen_FnReturnsError(t *testing.T) {
	result := Success(10)
	newResult := result.Then(func(v int) error {
		return errors.New("fn error")
	})
	if newResult.IsSuccess {
		t.Error("Then should return failure when fn returns error")
	}
	if newResult.Message != "fn error" {
		t.Errorf("Message should be 'fn error', got '%s'", newResult.Message)
	}
}

func TestMap_Success(t *testing.T) {
	result := Success(10)
	mapped := Map(result, func(v int) string {
		return "value is 10"
	})
	if !mapped.IsSuccess {
		t.Error("Map should return success")
	}
	if mapped.Content != "value is 10" {
		t.Errorf("Mapped content should be 'value is 10', got '%s'", mapped.Content)
	}
}

func TestMap_Failure(t *testing.T) {
	result := Fail[int](ErrCodeReadFailed, "讀取失敗")
	mapped := Map(result, func(v int) string {
		return "should not be called"
	})
	if mapped.IsSuccess {
		t.Error("Map should propagate failure")
	}
	if mapped.ErrorCode != ErrCodeReadFailed {
		t.Errorf("Map should preserve ErrorCode")
	}
}

func TestGetContentOrDefault_Success(t *testing.T) {
	result := Success(42)
	value := result.GetContentOrDefault(0)
	if value != 42 {
		t.Errorf("Should return content, got %d", value)
	}
}

func TestGetContentOrDefault_Failure(t *testing.T) {
	result := Fail[int](ErrCodeUnknown, "error")
	value := result.GetContentOrDefault(99)
	if value != 99 {
		t.Errorf("Should return default value 99, got %d", value)
	}
}

// =============================================================================
// SimpleResult 測試
// =============================================================================

func TestSimpleSuccess(t *testing.T) {
	result := SimpleSuccess()
	if !result.IsSuccess {
		t.Error("SimpleSuccess should have IsSuccess = true")
	}
	if result.ErrorCode != ErrCodeSuccess {
		t.Errorf("ErrorCode should be 0, got %d", result.ErrorCode)
	}
}

func TestSimpleFail(t *testing.T) {
	result := SimpleFail(ErrCodeWriteFailed, "寫入失敗")
	if result.IsSuccess {
		t.Error("SimpleFail should have IsSuccess = false")
	}
	if result.ErrorCode != ErrCodeWriteFailed {
		t.Errorf("ErrorCode should be %d, got %d", ErrCodeWriteFailed, result.ErrorCode)
	}
}

func TestSimpleFailFromError(t *testing.T) {
	err := errors.New("測試")
	result := SimpleFailFromError(err)
	if result.IsSuccess {
		t.Error("SimpleFailFromError should have IsSuccess = false")
	}
	if result.Message != "測試" {
		t.Errorf("Message should be '測試', got '%s'", result.Message)
	}
}

func TestSimpleFailFromError_NilError(t *testing.T) {
	result := SimpleFailFromError(nil)
	if !result.IsSuccess {
		t.Error("SimpleFailFromError with nil should return success")
	}
}

// =============================================================================
// 錯誤代碼常數測試
// =============================================================================

func TestErrorCodeConstants(t *testing.T) {
	codes := map[int]string{
		ErrCodeSuccess:          "Success",
		ErrCodeConnectionFailed: "ConnectionFailed",
		ErrCodeTimeout:          "Timeout",
		ErrCodeDisconnected:     "Disconnected",
		ErrCodeInvalidAddress:   "InvalidAddress",
		ErrCodeInvalidDataType:  "InvalidDataType",
		ErrCodeInvalidValue:     "InvalidValue",
		ErrCodeReadFailed:       "ReadFailed",
		ErrCodeWriteFailed:      "WriteFailed",
		ErrCodeDeviceError:      "DeviceError",
		ErrCodeUnknown:          "Unknown",
	}

	for code, name := range codes {
		if code == 0 && name != "Success" {
			t.Errorf("ErrCodeSuccess should be 0")
		}
		// 確保每個代碼都是唯一的
		for code2, name2 := range codes {
			if code == code2 && name != name2 {
				t.Errorf("Duplicate error code: %d", code)
			}
		}
	}
}
