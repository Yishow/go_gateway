//go:build windows

package desktop

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

func TestOwnerCaseAliasesShareOneActualFileGuard(t *testing.T) {
	path := filepath.Join(t.TempDir(), "gateway.db")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	owner, err := AcquireOwner(path, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close()
	for _, candidate := range []string{path, strings.ToUpper(path)} {
		identity, err := CanonicalDatabase(candidate)
		if err != nil {
			t.Fatal(err)
		}
		if identity.ID != owner.Identity().ID {
			t.Fatalf("alias identity changed: %s", candidate)
		}
		second, err := AcquireOwner(candidate, false, nil)
		if second != nil {
			_ = second.Close()
		}
		if !errors.Is(err, ErrAlreadyRunning) {
			t.Fatalf("duplicate owner accepted: %v", err)
		}
	}
	if err = owner.Close(); err != nil {
		t.Fatal(err)
	}
	next, err := AcquireOwner(path, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := next.Close(); err != nil {
		t.Fatal(err)
	}
}
func TestOwnerReservationDoesNotExtendDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new.db")
	owner, err := AcquireOwner(path, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close()
	info, err := os.Stat(path)
	if err != nil || info.Size() != 0 {
		t.Fatalf("reservation modified file data: %v", err)
	}
	actual, err := CanonicalDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	if actual != owner.Identity() {
		t.Fatal("missing to created identity split")
	}
	if RequestOpenSetup(path) == nil {
		t.Fatal("headless owner exposed setup handoff")
	}
}
func TestOwnerDeathReleasesOSLockDespiteStalePID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crash.db")
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	child := exec.CommandContext(t.Context(), executable, "-test.run=^TestOwnerSubprocessHelper$")
	child.Env = append(os.Environ(), "GATEWAY_OWNER_TEST_PATH="+path)
	output, err := child.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err = child.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if child.ProcessState == nil {
			_ = child.Process.Kill()
			_ = child.Wait()
		}
	})
	lines := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(output)
		if scanner.Scan() {
			lines <- scanner.Text()
		} else {
			lines <- ""
		}
	}()
	select {
	case line := <-lines:
		if line != "locked" {
			t.Fatalf("helper not ready: %q", line)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("helper startup timeout")
	}
	second, err := AcquireOwner(path, false, nil)
	if second != nil {
		_ = second.Close()
	}
	if !errors.Is(err, ErrAlreadyRunning) {
		t.Fatalf("live owner unprotected: %v", err)
	}
	if err = child.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = child.Wait()
	if err = os.WriteFile(path+".pid", []byte("1234"), 0o600); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for {
		owner, err := AcquireOwner(path, false, nil)
		if err == nil {
			_ = owner.Close()
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("OS did not release owner: %v", err)
		}
		time.Sleep(20 * time.Millisecond)
	}
}
func TestOwnerSubprocessHelper(t *testing.T) {
	path := os.Getenv("GATEWAY_OWNER_TEST_PATH")
	if path == "" {
		return
	}
	owner, err := AcquireOwner(path, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close()
	fmt.Println("locked")
	<-time.After(time.Hour)
}

func TestHardlinksAreIdentifiedButFreshWritableOwnershipIsRejected(t *testing.T) {
	path := filepath.Join(t.TempDir(), "original.db")
	alias := filepath.Join(filepath.Dir(path), "alias.db")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(path, alias); err != nil {
		t.Fatal(err)
	}
	original, err := CanonicalDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	linked, err := CanonicalDatabase(alias)
	if err != nil {
		t.Fatal(err)
	}
	if original.ID != linked.ID {
		t.Fatal("hardlink identity was not recognized")
	}
	// Model an existing OS owner: the lock must win before the new-link policy.
	handle, err := openIdentityHandle(path, windows.OPEN_EXISTING, windows.GENERIC_READ|windows.GENERIC_WRITE)
	if err != nil {
		t.Fatal(err)
	}
	offset := ownerLockOffset()
	if err = windows.LockFileEx(handle, windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &offset); err != nil {
		_ = windows.CloseHandle(handle)
		t.Fatal(err)
	}
	owner, ownerErr := AcquireOwner(alias, false, nil)
	if owner != nil {
		_ = owner.Close()
	}
	if err = windows.CloseHandle(handle); err != nil {
		t.Fatal(err)
	}
	if !errors.Is(ownerErr, ErrAlreadyRunning) {
		t.Fatalf("active owner was not authoritative: %v", ownerErr)
	}
	for _, candidate := range []string{path, alias} {
		owner, err := AcquireOwner(candidate, false, nil)
		if owner != nil {
			_ = owner.Close()
		}
		if !errors.Is(err, ErrHardlinkedDatabase) {
			t.Fatalf("ambiguous WAL namespace accepted: %v", err)
		}
	}
	// The rejection must have closed its handles and released its lock.
	if err = os.Remove(alias); err != nil {
		t.Fatal(err)
	}
	owner, err = AcquireOwner(path, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = owner.Close(); err != nil {
		t.Fatal(err)
	}
}
