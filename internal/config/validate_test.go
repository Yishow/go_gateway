package config

import (
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

func TestValidationResult_ErrorString(t *testing.T) {
	result := &ValidationResult{Valid: true}
	result.AddError("field1", "error1")
	result.AddError("field2", "error2")

	errStr := result.ErrorString()
	if errStr == "" {
		t.Error("應該返回錯誤字串")
	}
}
