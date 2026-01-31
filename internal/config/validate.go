package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// =============================================================================
// 配置驗證
// =============================================================================

// ValidationError 驗證錯誤
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult 驗證結果
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// AddError 新增錯誤
func (r *ValidationResult) AddError(field, message string) {
	r.Valid = false
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

// ErrorString 返回所有錯誤的字串表示
func (r *ValidationResult) ErrorString() string {
	if len(r.Errors) == 0 {
		return ""
	}

	var msgs []string
	for _, e := range r.Errors {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "; ")
}

// =============================================================================
// 配置驗證器
// =============================================================================

// Validator 配置驗證器
type Validator struct {
	cfg *Config
}

// NewValidator 建立新的驗證器
func NewValidator(cfg *Config) *Validator {
	return &Validator{cfg: cfg}
}

// Validate 驗證配置
func (v *Validator) Validate() *ValidationResult {
	result := &ValidationResult{Valid: true}

	// 驗證伺服器配置
	v.validateServer(result)

	// 驗證 API 配置
	v.validateAPI(result)

	// 驗證 CORS 配置
	v.validateCORS(result)

	// 驗證連接池配置
	v.validateConnectionPool(result)

	return result
}

// validateServer 驗證伺服器配置
func (v *Validator) validateServer(result *ValidationResult) {
	port, err := strconv.Atoi(v.cfg.Server.Port)
	if err != nil {
		result.AddError("server.port", "埠號必須是數字")
		return
	}

	if port < 1 || port > 65535 {
		result.AddError("server.port", "埠號必須在 1-65535 之間")
	}
}

// validateAPI 驗證 API 配置
func (v *Validator) validateAPI(result *ValidationResult) {
	if v.cfg.API.BasePath == "" {
		result.AddError("api.base_path", "API 基礎路徑不能為空")
	}

	if !strings.HasPrefix(v.cfg.API.BasePath, "/") {
		result.AddError("api.base_path", "API 基礎路徑必須以 / 開頭")
	}
}

// validateCORS 驗證 CORS 配置
func (v *Validator) validateCORS(result *ValidationResult) {
	if len(v.cfg.CORS.AllowOrigins) == 0 {
		result.AddError("cors.allow_origins", "CORS 允許來源不能為空")
	}
}

// validateConnectionPool 驗證連接池配置
func (v *Validator) validateConnectionPool(result *ValidationResult) {
	if v.cfg.ConnectionPool.MaxConnections < 1 {
		result.AddError("connection_pool.max_connections", "最大連接數必須至少為 1")
	}

	if v.cfg.ConnectionPool.ConnectionTimeout < 0 {
		result.AddError("connection_pool.connection_timeout", "連接超時不能為負數")
	}
}

// =============================================================================
// 快捷驗證函數
// =============================================================================

// Validate 驗證當前配置
func Validate() (*ValidationResult, error) {
	cfg := Get()
	if cfg == nil {
		return nil, errors.New("配置尚未載入")
	}

	validator := NewValidator(cfg)
	return validator.Validate(), nil
}

// ValidateConfig 驗證指定配置
func ValidateConfig(cfg *Config) *ValidationResult {
	validator := NewValidator(cfg)
	return validator.Validate()
}
