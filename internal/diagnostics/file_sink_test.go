package diagnostics

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fileOptions(t *testing.T) FileOptions {
	t.Helper()
	// Temporary directories may include system symlinks (for example /var on macOS).
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	db := filepath.Join(root, "datalink.db")
	if err := os.WriteFile(db, []byte("protected database"), 0o600); err != nil {
		t.Fatal(err)
	}
	return FileOptions{Path: filepath.Join(root, "logs", "runtime.jsonl"), DatabasePath: db, DatabaseID: strings.Repeat("a", 64)}
}
func TestRejectProtectedAndUnownedFileSet(t *testing.T) {
	for _, suffix := range []string{"", ".1", ".2"} {
		t.Run("unowned"+suffix, func(t *testing.T) {
			o := fileOptions(t)
			if err := os.MkdirAll(filepath.Dir(o.Path), 0o700); err != nil {
				t.Fatal(err)
			}
			path := o.Path + suffix
			original := []byte("unrelated logger content")
			if err := os.WriteFile(path, original, 0o600); err != nil {
				t.Fatal(err)
			}
			sink, err := OpenFileSink(o)
			if err == nil {
				_ = sink.Close(t.Context())
				t.Fatal("unowned file accepted")
			}
			got, _ := os.ReadFile(path)
			if !bytes.Equal(got, original) {
				t.Fatal("unowned file changed")
			}
		})
	}
	o := fileOptions(t)
	o.Path = o.DatabasePath
	if sink, err := OpenFileSink(o); err == nil {
		_ = sink.Close(t.Context())
		t.Fatal("database accepted as diagnostics")
	}
}

func TestDefaultNamespacesOverlapLocksAndOwnedReopen(t *testing.T) {
	o := fileOptions(t)
	if DefaultPath(filepath.Dir(o.DatabasePath), o.DatabaseID) == DefaultPath(filepath.Dir(o.DatabasePath), strings.Repeat("b", 64)) {
		t.Fatal("shared default namespace")
	}
	sink, err := OpenFileSink(o)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sink.Close(t.Context()) }()
	for _, p := range []string{o.Path, o.Path + ".1", o.Path + ".2"} {
		other := o
		other.Path = p
		if duplicate, err := OpenFileSink(other); err == nil {
			_ = duplicate.Close(t.Context())
			t.Fatal("overlapping set admitted")
		}
	}
	b, err := New(Options{File: sink})
	if err != nil {
		t.Fatal(err)
	}
	b.Emit(Input{Code: "startup.begin"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	if s := sink.Status(); !s.Saved || s.Path != o.Path {
		t.Fatalf("saved evidence absent: %#v", s)
	}
	if err := b.Close(t.Context()); err != nil {
		t.Fatal(err)
	}
	again, err := OpenFileSink(o)
	if err != nil {
		t.Fatal(err)
	}
	_ = again.Close(t.Context())
}
func TestAllProtectedAliasesRejectedBeforeModification(t *testing.T) {
	for _, kind := range []string{"hardlink", "symlink"} {
		for _, suffix := range []string{"", ".1", ".2"} {
			for _, dbSuffix := range []string{"", "-wal", "-shm"} {
				t.Run(kind+suffix+dbSuffix, func(t *testing.T) {
					o := fileOptions(t)
					protected := o.DatabasePath + dbSuffix
					if err := os.WriteFile(protected, []byte("protected bytes"), 0o600); err != nil {
						t.Fatal(err)
					}
					if err := os.MkdirAll(filepath.Dir(o.Path), 0o700); err != nil {
						t.Fatal(err)
					}
					link := os.Link
					if kind == "symlink" {
						link = os.Symlink
					}
					if err := link(protected, o.Path+suffix); err != nil {
						t.Skipf("link fixture unavailable: %v", err)
					}
					if sink, err := OpenFileSink(o); err == nil {
						_ = sink.Close(t.Context())
						t.Fatal("protected alias accepted")
					}
					got, _ := os.ReadFile(protected)
					if string(got) != "protected bytes" {
						t.Fatal("protected bytes changed")
					}
				})
			}
		}
	}
}
func TestDirectorySymlinkAndWindowsLeafSyntaxRejected(t *testing.T) {
	o := fileOptions(t)
	actual := filepath.Join(filepath.Dir(o.Path), "actual")
	if err := os.MkdirAll(actual, 0o700); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(filepath.Dir(o.Path), "alias")
	if err := os.Symlink(actual, alias); err != nil {
		t.Skipf("symlink fixture unavailable: %v", err)
	}
	o.Path = filepath.Join(alias, "runtime.jsonl")
	if sink, err := OpenFileSink(o); err == nil {
		_ = sink.Close(t.Context())
		t.Fatal("unverified directory link accepted")
	}
	for _, name := range []string{"db:stream", "NUL", "con.jsonl", "logs.", "logs ", "../unsafe"} {
		if safeLeaf(name) {
			t.Fatalf("unsafe leaf accepted %s", name)
		}
	}
}
func TestMissingProtectedLeafThroughDirectoryAliasNotCreated(t *testing.T) {
	o := fileOptions(t)
	parent := filepath.Dir(o.DatabasePath)
	alias := filepath.Join(parent, "parent-alias")
	if err := os.Symlink(parent, alias); err != nil {
		t.Skipf("symlink fixture unavailable: %v", err)
	}
	o.Path = filepath.Join(parent, "fresh-wal")
	o.ProtectedPaths = []string{filepath.Join(alias, "fresh-wal")}
	if sink, err := OpenFileSink(o); err == nil {
		_ = sink.Close(t.Context())
		t.Fatal("missing protected leaf accepted")
	}
	if _, err := os.Stat(o.Path); !os.IsNotExist(err) {
		t.Fatal("protected leaf created before alias detection")
	}
}
