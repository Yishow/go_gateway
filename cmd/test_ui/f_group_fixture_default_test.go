//go:build !f_write_group_fixture

package main

import "testing"

func TestNormalBinaryFixtureHooksAreNoOp(t *testing.T) {
	hooks := newGroupFixture()
	if _, ok := hooks.(defaultGroupFixture); !ok {
		t.Fatalf("normal build returned %T, want default no-op fixture", hooks)
	}
	config := hooks.pipelineConfig("node", "owner")
	if config.Now != nil {
		t.Fatal("normal pipeline must use the production wall clock")
	}
	if err := preflightGroupFixtureDatabase(); err != nil {
		t.Fatalf("normal preflight returned error: %v", err)
	}
	cleanup, err := hooks.configure(t.Context(), nil, nil)
	if err != nil {
		t.Fatalf("normal no-op configure returned error: %v", err)
	}
	if err := cleanup(); err != nil {
		t.Fatalf("normal no-op cleanup returned error: %v", err)
	}
}
