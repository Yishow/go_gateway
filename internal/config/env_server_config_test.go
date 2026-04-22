package config

import "testing"

func TestNewConfigFromEnv_UsesServerHostAndPort(t *testing.T) {
	t.Setenv("HOST", "0.0.0.0")
	t.Setenv("PORT", "3333")

	cfg := newConfigFromEnv()

	if cfg.Server.Host != "0.0.0.0" {
		t.Fatalf("expected host from env, got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != "3333" {
		t.Fatalf("expected port from env, got %q", cfg.Server.Port)
	}
}
