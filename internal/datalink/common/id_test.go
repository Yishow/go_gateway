package common

import "testing"

func TestNewUUID(t *testing.T) {
	id, err := NewUUID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(id) != 36 {
		t.Fatalf("expected UUID length 36, got %d", len(id))
	}
	if id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		t.Fatalf("invalid UUID format: %s", id)
	}
}
