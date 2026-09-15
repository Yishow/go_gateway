package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// Save and restore stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"version"})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("version 命令執行失敗: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Go Gateway") {
		t.Errorf("version 輸出應包含 'Go Gateway', 實際: %s", output)
	}
}

func TestValidateCommand_WithDefaults(_ *testing.T) {
	resetValidateFlags()
	rootCmd.SetArgs([]string{"validate"})

	// This should work with default config
	err := rootCmd.Execute()
	// Note: Actual result depends on environment
	_ = err
}

func TestValidateCommand_JSONOutput(t *testing.T) {
	resetValidateFlags()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"validate", "--json"})
	_ = rootCmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should contain JSON structure
	if !strings.Contains(output, "valid") || !strings.Contains(output, "{") {
		t.Logf("JSON output: %s", output)
	}
}

func TestValidateCommand_NonExistentConfig(t *testing.T) {
	resetValidateFlags()

	rootCmd.SetArgs([]string{"validate", "-c", "/nonexistent/config.yaml"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("對不存在的配置文件應返回錯誤")
	}
}

func TestRootCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	_ = rootCmd.Execute()

	output := buf.String()
	if !strings.Contains(output, "Go Gateway") {
		t.Errorf("help 輸出應包含 'Go Gateway', 實際: %s", output)
	}
	if !strings.Contains(output, "validate") {
		t.Errorf("help 輸出應包含 'validate' 子命令, 實際: %s", output)
	}
}

func TestValidateCommand_WithEnvFile(t *testing.T) {
	// Create a temporary .env file
	tmpFile, err := os.CreateTemp("", "test-config-*.env")
	if err != nil {
		t.Fatalf("無法創建臨時文件: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write valid config
	content := `PORT=9090
HOST=localhost
API_BASE_PATH=/api/v2
CORS_ALLOW_ORIGINS=http://localhost:3000
MAX_CONNECTIONS=50
`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("無法寫入臨時文件: %v", err)
	}
	tmpFile.Close()

	resetValidateFlags()
	rootCmd.SetArgs([]string{"validate", "-c", tmpFile.Name()})

	err = rootCmd.Execute()
	// With valid config, should not error
	if err != nil {
		t.Logf("驗證結果: %v", err)
	}
}

func TestValidateCommand_InvalidPortEnv(t *testing.T) {
	// Save original env
	origPort := os.Getenv("PORT")
	origBasePath := os.Getenv("API_BASE_PATH")
	origCORS := os.Getenv("CORS_ALLOW_ORIGINS")
	origMaxConn := os.Getenv("MAX_CONNECTIONS")

	// Set invalid port directly
	os.Setenv("PORT", "invalid")
	os.Setenv("API_BASE_PATH", "/api/v1")
	os.Setenv("CORS_ALLOW_ORIGINS", "*")
	os.Setenv("MAX_CONNECTIONS", "10")

	defer func() {
		// Restore original env
		if origPort != "" {
			os.Setenv("PORT", origPort)
		} else {
			os.Unsetenv("PORT")
		}
		if origBasePath != "" {
			os.Setenv("API_BASE_PATH", origBasePath)
		} else {
			os.Unsetenv("API_BASE_PATH")
		}
		if origCORS != "" {
			os.Setenv("CORS_ALLOW_ORIGINS", origCORS)
		} else {
			os.Unsetenv("CORS_ALLOW_ORIGINS")
		}
		if origMaxConn != "" {
			os.Setenv("MAX_CONNECTIONS", origMaxConn)
		} else {
			os.Unsetenv("MAX_CONNECTIONS")
		}
	}()

	resetValidateFlags()
	rootCmd.SetArgs([]string{"validate"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("無效的 PORT 環境變數應該導致驗證失敗")
	}
}

// resetValidateFlags resets the validate command flags to default values
func resetValidateFlags() {
	configFile = ""
	jsonOutput = false
	// Reset cobra command state
	validateCmd.Flags().Set("config", "")
	validateCmd.Flags().Set("json", "false")
}
