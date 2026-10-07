//go:build windows

package desktop

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows"
	_ "modernc.org/sqlite"
)

func TestOwnerPinsAncestorsUntilCloseAndAllowsSQLite(t *testing.T) {
	root := t.TempDir()
	ancestor := filepath.Join(root, "parent")
	directory := filepath.Join(ancestor, "nested")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	owner, err := AcquireOwner(filepath.Join(directory, "gateway.db"), false, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close()
	if err = os.Rename(ancestor, ancestor+"-moved"); err == nil {
		t.Fatal("held ancestor could be renamed/replaced")
	}
	writer, err := openDirectoryWriter(directory)
	if err == nil {
		_ = windows.CloseHandle(writer)
		t.Fatal("held ancestor permitted a reparse-capable write handle")
	}
	// Namespace protection must not break SQLite's normal child-file WAL lifecycle.
	db, err := sql.Open("sqlite", owner.Identity().Path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.ExecContext(t.Context(), "PRAGMA journal_mode=WAL"); err != nil {
		_ = db.Close()
		t.Fatal(err)
	}
	if _, err = db.ExecContext(t.Context(), "CREATE TABLE pin_test(value INTEGER); INSERT INTO pin_test VALUES(1)"); err != nil {
		_ = db.Close()
		t.Fatal(err)
	}
	if err = db.Close(); err != nil {
		t.Fatal(err)
	}
	if err = owner.Close(); err != nil {
		t.Fatal(err)
	}
	if err = os.Rename(ancestor, ancestor+"-moved"); err != nil {
		t.Fatalf("namespace pins survived Owner.Close: %v", err)
	}
}
func TestFailedNamespacePinReleasesEarlierDirectories(t *testing.T) {
	root := t.TempDir()
	ancestor := filepath.Join(root, "parent")
	directory := filepath.Join(ancestor, "nested")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, "gateway.db")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	identity, err := CanonicalDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	writer, err := openDirectoryWriter(directory)
	if err != nil {
		t.Fatal(err)
	}
	pins, pinErr := pinOwnerNamespace(identity)
	if pins != nil {
		_ = pins.close()
	}
	if err = windows.CloseHandle(writer); err != nil {
		t.Fatal(err)
	}
	if pinErr == nil {
		t.Fatal("namespace pin ignored preexisting mutation handle")
	}
	if err = os.Rename(ancestor, ancestor+"-moved"); err != nil {
		t.Fatalf("failed pin leaked an ancestor handle: %v", err)
	}
}
func openDirectoryWriter(path string) (windows.Handle, error) {
	return windows.CreateFile(windows.StringToUTF16Ptr(path), windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE, nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
}
