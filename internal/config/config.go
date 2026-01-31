package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config 應用程式配置
type Config struct {
	// 服務器配置
	Server ServerConfig

	// API 配置
	API APIConfig

	// CORS 配置
	CORS CORSConfig

	// 日誌配置
	Log LogConfig

	// WebSocket 配置
	WebSocket WebSocketConfig

	// 連接池配置
	ConnectionPool ConnectionPoolConfig

	// 開發環境配置
	Debug bool
}

// ServerConfig 服務器配置
type ServerConfig struct {
	Port string
	Host string
}

// APIConfig API 配置
type APIConfig struct {
	BasePath string
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins []string
}

// LogConfig 日誌配置
type LogConfig struct {
	Level  string
	Output string
	File   string
}

// WebSocketConfig WebSocket 配置
type WebSocketConfig struct {
	ReadTimeout  int
	WriteTimeout int
	PingInterval int
}

// ConnectionPoolConfig 連接池配置
type ConnectionPoolConfig struct {
	MaxConnections  int
	ConnectionTimeout int
	ReadTimeout     int
	WriteTimeout    int
}

var globalConfig *Config

// Load 載入配置
// 優先順序：環境變數 > .env 文件 > 預設值
func Load() (*Config, error) {
	// 嘗試載入 .env 文件（如果存在）
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", ""),
		},
		API: APIConfig{
			BasePath: getEnv("API_BASE_PATH", "/api/v1"),
		},
		CORS: CORSConfig{
			AllowOrigins: parseStringSlice(getEnv("CORS_ALLOW_ORIGINS", "*")),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Output: getEnv("LOG_OUTPUT", "console"),
			File:   getEnv("LOG_FILE", "./logs/gateway.log"),
		},
		WebSocket: WebSocketConfig{
			ReadTimeout:  getEnvAsInt("WS_READ_TIMEOUT", 60),
			WriteTimeout: getEnvAsInt("WS_WRITE_TIMEOUT", 10),
			PingInterval: getEnvAsInt("WS_PING_INTERVAL", 30),
		},
		ConnectionPool: ConnectionPoolConfig{
			MaxConnections:    getEnvAsInt("MAX_CONNECTIONS", 100),
			ConnectionTimeout: getEnvAsInt("CONNECTION_TIMEOUT", 5),
			ReadTimeout:       getEnvAsInt("READ_TIMEOUT", 2),
			WriteTimeout:      getEnvAsInt("WRITE_TIMEOUT", 2),
		},
		Debug: getEnvAsBool("DEBUG", false),
	}

	// 設定 Gin 模式
	ginMode := getEnv("GIN_MODE", "release")
	os.Setenv("GIN_MODE", ginMode)

	globalConfig = cfg
	return cfg, nil
}

// Get 取得全域配置
func Get() *Config {
	if globalConfig == nil {
		// 如果未載入，使用預設配置
		cfg, _ := Load()
		return cfg
	}
	return globalConfig
}

// getEnv 取得環境變數，如果不存在則返回預設值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 取得環境變數並轉換為整數
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsBool 取得環境變數並轉換為布林值
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// parseStringSlice 解析逗號分隔的字符串為字符串切片
func parseStringSlice(value string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

// GetServerAddr 取得服務器監聽地址
func (c *Config) GetServerAddr() string {
	port := c.Server.Port
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	if c.Server.Host != "" {
		return fmt.Sprintf("%s%s", c.Server.Host, port)
	}
	return port
}

// LoadFromFile 從指定文件載入配置
// 支援 .env 格式的配置文件
func LoadFromFile(path string) (*Config, error) {
	// 載入指定的 .env 文件
	if err := godotenv.Load(path); err != nil {
		return nil, fmt.Errorf("無法載入配置文件 %s: %w", path, err)
	}

	// 使用 Load 函數載入配置（會讀取剛載入的環境變數）
	return Load()
}
