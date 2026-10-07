//go:build !windows

package desktop

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Owner preserves non-Windows CLI behavior without creating storage or a native lock.
type Owner struct{ identity Identity }

// AcquireOwner resolves a namespace only; native Windows locking is platform isolated.
func AcquireOwner(path string, _ bool, _ func()) (*Owner, error) {
	i, e := CanonicalDatabase(path)
	if e != nil {
		return nil, e
	}
	return &Owner{identity: i}, nil
}

// Identity returns the selected canonical namespace.
func (o *Owner) Identity() Identity { return o.identity }

// Close has no native resource to release on non-Windows.
func (o *Owner) Close() error { return nil }

// CanonicalDatabase resolves the existing parent without creating a missing database.
func CanonicalDatabase(path string) (Identity, error) {
	absolute, e := filepath.Abs(path)
	if e != nil {
		return Identity{}, e
	}
	canonical, e := filepath.EvalSymlinks(absolute)
	if errors.Is(e, os.ErrNotExist) {
		parent, err := filepath.EvalSymlinks(filepath.Dir(absolute))
		if err != nil {
			return Identity{}, err
		}
		canonical = filepath.Join(parent, filepath.Base(absolute))
	} else if e != nil {
		return Identity{}, e
	}
	return Identity{Path: canonical, ID: opaqueID(canonical)}, nil
}

// RequestOpenSetup requires a Windows native owner.
func RequestOpenSetup(string) error { return ErrUnsupported }

type unsupportedShell struct{}

// New returns a native-unavailable boundary for non-Windows builds.
func New(Options) Shell               { return unsupportedShell{} }
func (unsupportedShell) Start() error { return ErrUnsupported }
func (unsupportedShell) Update(State) {}
func (unsupportedShell) Close() error { return nil }

// ConfirmCreateDB never creates native dialogs on a CLI platform.
func ConfirmCreateDB(string) bool { return false }

// ShowInfo never creates native dialogs on a CLI platform.
func ShowInfo(string, string) {}

// ShowError never creates native dialogs on a CLI platform.
func ShowError(ErrorInfo) {}

// OpenURL uses the platform opener directly, without command-shell interpolation.
func OpenURL(raw string) error {
	if e := validateURL(raw); e != nil {
		return e
	}
	command := "xdg-open"
	if runtime.GOOS == "darwin" {
		command = "open"
	}
	child := exec.CommandContext(context.Background(), command, raw)
	if e := child.Start(); e != nil {
		return e
	}
	go func() {
		_ = child.Wait() //nolint:errcheck // Reap the asynchronous opener; its exit does not control the gateway.
	}()
	return nil
}
