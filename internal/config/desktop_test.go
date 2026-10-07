package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDesktopConfigUsesOnlySelectedDirectory(t *testing.T) {
	previous := globalConfig
	t.Cleanup(func() { globalConfig = previous })
	root, cwd := t.TempDir(), t.TempDir()
	t.Chdir(cwd)
	for _, key := range []string{"PORT", "HOST", "GIN_MODE"} {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	for name, data := range map[string]string{filepath.Join(root, ".env"): "PORT=8101\n", filepath.Join(cwd, ".env"): "HOST=192.0.2.1\nPORT=9999\n"} {
		if err := os.WriteFile(name, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	cfg, err := LoadFromDirectory(root)
	if err != nil || cfg.Server.Port != "8101" || cfg.Server.Host != "" {
		t.Fatalf("%+v %v", cfg, err)
	}
	t.Setenv("PORT", "8202")
	cfg, err = LoadFromDirectory(root)
	if err != nil || cfg.Server.Port != "8202" {
		t.Fatalf("%+v %v", cfg, err)
	}
}
func TestServerAddressSupportsIPv6(t *testing.T) {
	for _, tc := range []struct{ host, port, want string }{{"", ":8080", ":8080"}, {"::1", "8080", "[::1]:8080"}, {"[::1]", "8080", "[::1]:8080"}, {"127.0.0.1", "8080", "127.0.0.1:8080"}} {
		cfg := Config{Server: ServerConfig{Host: tc.host, Port: tc.port}}
		if got := cfg.GetServerAddr(); got != tc.want {
			t.Fatalf("%s != %s", got, tc.want)
		}
	}
}
