package diagnostics

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// ErrFileUnsafe is a path-free startup error suitable for safe classification.
var ErrFileUnsafe = errors.New("diagnostic file ownership or protected-path validation failed")

// DefaultPath provides one namespace for each canonical database identity.
func DefaultPath(root, id string) string {
	return filepath.Join(root, "logs", id, "gateway-runtime.jsonl")
}

type ownedFiles struct {
	root           *os.Root
	files          [3]*os.File
	names          [3]string
	path           string
	header         []byte
	size, rollback int64
	needsRollback  bool
}

func openOwnedFiles(o FileOptions) (_ *ownedFiles, err error) {
	if len(o.DatabaseID) != 64 || o.Path == "" || o.DatabasePath == "" {
		return nil, ErrFileUnsafe
	}
	if _, err := hex.DecodeString(o.DatabaseID); err != nil {
		return nil, ErrFileUnsafe
	}
	path, err := filepath.Abs(o.Path)
	if err != nil {
		return nil, ErrFileUnsafe
	}
	protected := append([]string{o.DatabasePath, o.DatabasePath + "-wal", o.DatabasePath + "-shm"}, o.ProtectedPaths...)
	for i, p := range protected {
		a, err := filepath.Abs(p)
		if err != nil {
			return nil, ErrFileUnsafe
		}
		protected[i] = a
	}
	names := [3]string{filepath.Base(path), filepath.Base(path) + ".1", filepath.Base(path) + ".2"}
	for _, name := range names {
		if !safeLeaf(name) {
			return nil, ErrFileUnsafe
		}
		for _, p := range protected {
			if strings.EqualFold(filepath.Join(filepath.Dir(path), name), p) {
				return nil, ErrFileUnsafe
			}
		}
	}
	root, err := openSafeDirectory(filepath.Dir(path), true)
	if err != nil {
		return nil, ErrFileUnsafe
	}
	owned := &ownedFiles{root: root, names: names, path: path}
	defer func() {
		if err != nil {
			owned.close()
		}
	}()
	rootInfo, err := root.Stat(".")
	if err != nil {
		return nil, ErrFileUnsafe
	}
	// Parent identities protect nonexistent WAL/SHM names through directory aliases.
	for _, p := range protected {
		parent, err := os.Stat(filepath.Dir(p))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, ErrFileUnsafe
			}
			continue
		}
		if os.SameFile(rootInfo, parent) {
			for _, name := range names {
				if strings.EqualFold(name, filepath.Base(p)) {
					return nil, ErrFileUnsafe
				}
			}
		}
	}
	key := path
	if runtime.GOOS == "windows" {
		key = strings.ToLower(key)
	}
	sum := sha256.Sum256([]byte(key))
	header, err := json.Marshal(struct {
		Owner    string `json:"diagnostics_owner"`
		Database string `json:"database_id"`
		Set      string `json:"set_id"`
	}{"gateway-runtime-v1", o.DatabaseID, hex.EncodeToString(sum[:])})
	if err != nil {
		return nil, ErrFileUnsafe
	}
	owned.header = append(header, '\n')
	for i, name := range names {
		before, statErr := root.Lstat(name)
		if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
			return nil, ErrFileUnsafe
		}
		if statErr == nil && (!before.Mode().IsRegular() || reparse(before)) {
			return nil, ErrFileUnsafe
		}
		fresh := errors.Is(statErr, os.ErrNotExist)
		flags := os.O_RDWR
		if fresh {
			flags |= os.O_CREATE | os.O_EXCL
		}
		file, err := root.OpenFile(name, flags, 0o600)
		if err != nil {
			return nil, ErrFileUnsafe
		}
		owned.files[i] = file
		if err := lockDiagnosticFile(file); err != nil {
			return nil, ErrFileUnsafe
		}
		info, err := file.Stat()
		if err != nil || !info.Mode().IsRegular() || !singleLink(file) {
			return nil, ErrFileUnsafe
		}
		actual, err := root.Lstat(name)
		if err != nil || reparse(actual) || !actual.Mode().IsRegular() || !os.SameFile(info, actual) || !fresh && !os.SameFile(before, info) {
			return nil, ErrFileUnsafe
		}
		for _, p := range protected {
			other, err := os.Stat(p)
			if err == nil {
				if os.SameFile(info, other) {
					return nil, ErrFileUnsafe
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				return nil, ErrFileUnsafe
			}
		}
		for j := range i {
			other, err := owned.files[j].Stat()
			if err != nil || os.SameFile(info, other) {
				return nil, ErrFileUnsafe
			}
		}
		if fresh {
			if n, err := file.WriteAt(owned.header, 0); err != nil || n != len(owned.header) {
				return nil, ErrFileUnsafe
			}
			if err := file.Sync(); err != nil {
				return nil, ErrFileUnsafe
			}
			info, err = file.Stat()
			if err != nil {
				return nil, ErrFileUnsafe
			}
		} else {
			got := make([]byte, len(owned.header))
			if _, err := file.ReadAt(got, 0); err != nil || !bytes.Equal(got, owned.header) || info.Size() > MaxFileBytes {
				return nil, ErrFileUnsafe
			}
		}
		if i == 0 {
			owned.size = info.Size()
		}
	}
	return owned, nil
}

// Validate every component against its opened directory handle. Later operations
// use Root-relative handles, so swapping a pathname cannot redirect truncation.
func openSafeDirectory(path string, create bool) (*os.Root, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	volume := filepath.VolumeName(absolute)
	base := volume + string(os.PathSeparator)
	root, err := os.OpenRoot(base)
	if err != nil {
		return nil, err
	}
	for part := range strings.SplitSeq(strings.TrimPrefix(absolute, base), string(os.PathSeparator)) {
		if part == "" {
			continue
		}
		if !safeLeaf(part) {
			_ = root.Close()
			return nil, ErrFileUnsafe
		}
		before, err := root.Lstat(part)
		if errors.Is(err, os.ErrNotExist) && create {
			if err = root.Mkdir(part, 0o700); err == nil || errors.Is(err, os.ErrExist) {
				before, err = root.Lstat(part)
			}
		}
		if err != nil || !before.IsDir() || reparse(before) {
			_ = root.Close()
			return nil, ErrFileUnsafe
		}
		next, err := root.OpenRoot(part)
		if err != nil {
			_ = root.Close()
			return nil, ErrFileUnsafe
		}
		actual, err := next.Stat(".")
		after, afterErr := root.Lstat(part)
		if err != nil || afterErr != nil || reparse(after) || !os.SameFile(before, actual) || !os.SameFile(after, actual) {
			_ = next.Close()
			_ = root.Close()
			return nil, ErrFileUnsafe
		}
		_ = root.Close()
		root = next
	}
	return root, nil
}
func safeLeaf(name string) bool {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, "\x00:/\\") || strings.HasSuffix(name, ".") || strings.HasSuffix(name, " ") {
		return false
	}
	stem, _, _ := strings.Cut(strings.ToLower(name), ".")
	switch stem {
	case "con", "prn", "aux", "nul":
		return false
	}
	for i := 1; i <= 9; i++ {
		if stem == "com"+strconv.Itoa(i) || stem == "lpt"+strconv.Itoa(i) {
			return false
		}
	}
	return true
}
func (o *ownedFiles) close() {
	for _, f := range o.files {
		if f != nil {
			_ = f.Close()
		}
	}
	if o.root != nil {
		_ = o.root.Close()
	}
}
