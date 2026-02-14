package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_IgnoresMalformedDotEnv(t *testing.T) {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("取得工作目錄失敗: %v", err)
	}

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("INVALID LINE WITHOUT EQUALS\n"), 0o600); err != nil {
		t.Fatalf("建立測試 .env 失敗: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("切換工作目錄失敗: %v", err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(originalWD); chdirErr != nil {
			t.Fatalf("還原工作目錄失敗: %v", chdirErr)
		}
	})

	originalGlobal := globalConfig
	globalConfig = nil
	t.Cleanup(func() {
		globalConfig = originalGlobal
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load 應忽略格式錯誤的 .env，實際錯誤: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load 不應回傳 nil config")
	}
}
