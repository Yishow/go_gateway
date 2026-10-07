//go:build !windows

package desktop

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCanonicalCLIPathDoesNotCreateOrChangeNamespace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new.db")
	before, err := CanonicalDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("canonical lookup created storage")
	}
	if err = os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	after, err := CanonicalDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	if before != after {
		t.Fatal("creation changed CLI namespace")
	}
	owner, err := AcquireOwner(path, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = owner.Close(); err != nil {
		t.Fatal(err)
	}
}
