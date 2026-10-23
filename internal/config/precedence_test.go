package config

import (
	"os"
	"path/filepath"
	"testing"
)

// forgetEnvironment removes variables for one test and restores them afterwards,
// so dotenv loading by the code under test cannot leak into later tests.
func forgetEnvironment(t *testing.T, keys ...string) {
	t.Helper()
	for _, key := range keys {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLoadPrefersProcessEnvironmentOverWorkingDirectoryDotenvOverDefaults(t *testing.T) {
	previous := globalConfig
	t.Cleanup(func() { globalConfig = previous })
	t.Chdir(t.TempDir())
	forgetEnvironment(t, "HOST", "LOG_LEVEL", "API_BASE_PATH", "GIN_MODE")
	t.Setenv("PORT", "8202")
	if err := os.WriteFile(".env", []byte("PORT=8101\nHOST=192.0.2.7\nLOG_LEVEL=debug\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != "8202" {
		t.Fatalf("dotenv overrode the process environment: %q", cfg.Server.Port)
	}
	if cfg.Server.Host != "192.0.2.7" || cfg.Log.Level != "debug" {
		t.Fatalf("working-directory dotenv was ignored: host=%q level=%q", cfg.Server.Host, cfg.Log.Level)
	}
	if cfg.API.BasePath != "/api/v1" {
		t.Fatalf("default lost to an unset value: %q", cfg.API.BasePath)
	}
}

func TestLoadFromDirectoryMissingDotenvUsesDefaultsAndMalformedFailsClosed(t *testing.T) {
	previous := globalConfig
	t.Cleanup(func() { globalConfig = previous })
	selected, cwd := t.TempDir(), t.TempDir()
	t.Chdir(cwd)
	forgetEnvironment(t, "PORT", "HOST", "GIN_MODE")
	// A shortcut working directory must never supply desktop configuration.
	if err := os.WriteFile(filepath.Join(cwd, ".env"), []byte("PORT=9999\nHOST=192.0.2.1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadFromDirectory(selected)
	if err != nil || cfg.Server.Port != "8080" || cfg.Server.Host != "" {
		t.Fatalf("missing selected dotenv did not fall back to defaults only: %+v %v", cfg, err)
	}
	if err = os.WriteFile(filepath.Join(selected, ".env"), []byte("BROKEN='unterminated"), 0o600); err != nil {
		t.Fatal(err)
	}
	if cfg, err = LoadFromDirectory(selected); err == nil || cfg != nil {
		t.Fatalf("malformed selected dotenv was tolerated: %+v %v", cfg, err)
	}
}
