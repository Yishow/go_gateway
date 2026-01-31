package config

import (
	"strings"
	"testing"
)

func TestValidator_ValidConfig(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "localhost",
		},
		API: APIConfig{
			BasePath: "/api/v1",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections:    10,
			ConnectionTimeout: 30,
		},
	}

	result := ValidateConfig(cfg)
	if !result.Valid {
		t.Errorf("有效配置驗證失敗: %s", result.ErrorString())
	}
}

func TestValidator_ValidConfig_MultipleOrigins(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "3000",
		},
		API: APIConfig{
			BasePath: "/api/v2",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"http://localhost:3000", "https://example.com"},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 50,
		},
	}

	result := ValidateConfig(cfg)
	if !result.Valid {
		t.Errorf("多個 CORS 來源配置驗證應該通過: %s", result.ErrorString())
	}
}

func TestValidator_InvalidPort(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "0",
		},
		API: APIConfig{
			BasePath: "/api/v1",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 10,
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("無效埠號應該驗證失敗")
	}

	hasPortError := false
	for _, e := range result.Errors {
		if e.Field == "server.port" {
			hasPortError = true
		}
	}
	if !hasPortError {
		t.Error("應該包含 server.port 錯誤")
	}
}

func TestValidator_PortTooHigh(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "70000", // > 65535
		},
		API: APIConfig{
			BasePath: "/api/v1",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 10,
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("埠號超過 65535 應該驗證失敗")
	}
}

func TestValidator_NonNumericPort(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "invalid",
		},
		API: APIConfig{
			BasePath: "/api/v1",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 10,
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("非數字埠號應該驗證失敗")
	}
}

func TestValidator_EmptyPort(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "",
		},
		API: APIConfig{
			BasePath: "/api/v1",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 10,
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("空埠號應該驗證失敗")
	}
}

func TestValidator_InvalidBasePath(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		API: APIConfig{
			BasePath: "api/v1", // 缺少 /
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 10,
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("無效基礎路徑應該驗證失敗")
	}
}

func TestValidator_EmptyBasePath(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		API: APIConfig{
			BasePath: "",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 10,
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("空基礎路徑應該驗證失敗")
	}
}

func TestValidator_EmptyCORS(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		API: APIConfig{
			BasePath: "/api/v1",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 10,
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("空 CORS 允許來源應該驗證失敗")
	}
}

func TestValidator_InvalidConnectionPool(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		API: APIConfig{
			BasePath: "/api/v1",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 0, // 無效
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("無效連接池配置應該驗證失敗")
	}
}

func TestValidator_NegativeConnectionPool(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "8080",
		},
		API: APIConfig{
			BasePath: "/api/v1",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{"*"},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: -5,
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("負數連接池應該驗證失敗")
	}
}

func TestValidator_MultipleErrors(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "invalid",
		},
		API: APIConfig{
			BasePath: "no-slash",
		},
		CORS: CORSConfig{
			AllowOrigins: []string{},
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections: 0,
		},
	}

	result := ValidateConfig(cfg)
	if result.Valid {
		t.Error("應該有多個驗證錯誤")
	}
	if len(result.Errors) < 3 {
		t.Errorf("應該至少有 3 個錯誤，實際: %d", len(result.Errors))
	}
}

func TestValidationResult_ErrorString(t *testing.T) {
	result := &ValidationResult{Valid: true}
	result.AddError("field1", "error1")
	result.AddError("field2", "error2")

	errStr := result.ErrorString()
	if errStr == "" {
		t.Error("應該返回錯誤字串")
	}
	if !strings.Contains(errStr, "field1") {
		t.Error("錯誤字串應包含 field1")
	}
	if !strings.Contains(errStr, "field2") {
		t.Error("錯誤字串應包含 field2")
	}
}

func TestValidationResult_AddError_SetsInvalid(t *testing.T) {
	result := &ValidationResult{Valid: true}
	if !result.Valid {
		t.Error("初始狀態應為 valid")
	}

	result.AddError("test", "test error")
	if result.Valid {
		t.Error("添加錯誤後應為 invalid")
	}
}

func TestValidationResult_NoErrors(t *testing.T) {
	result := &ValidationResult{Valid: true}
	errStr := result.ErrorString()
	if errStr != "" {
		t.Errorf("無錯誤時應返回空字串，實際: %s", errStr)
	}
}

func TestValidator_SpecialPorts(t *testing.T) {
	testCases := []struct {
		port  string
		valid bool
		name  string
	}{
		{"1", true, "最小有效埠號"},
		{"80", true, "HTTP 預設埠"},
		{"443", true, "HTTPS 預設埠"},
		{"8080", true, "常用開發埠"},
		{"65535", true, "最大有效埠號"},
		{"65536", false, "超過最大埠號"},
		{"-1", false, "負數埠號"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{Port: tc.port},
				API:    APIConfig{BasePath: "/api"},
				CORS:   CORSConfig{AllowOrigins: []string{"*"}},
				ConnectionPool: ConnectionPoolConfig{MaxConnections: 10},
			}

			result := ValidateConfig(cfg)
			if result.Valid != tc.valid {
				t.Errorf("埠號 %s 預期 valid=%v，實際 valid=%v", tc.port, tc.valid, result.Valid)
			}
		})
	}
}
