//go:build windows

package desktop

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/sys/windows"
)

type namespacePins struct{ handles []windows.Handle }

// pinOwnerNamespace freezes each ordinary canonical directory object before the
// later SQLite path open. Denying write sharing also prevents an in-place reparse
// mutation; denying delete sharing prevents ancestor rename/replacement.
func pinOwnerNamespace(identity Identity) (*namespacePins, error) {
	paths, err := windowsAncestorPaths(identity.Path)
	if err != nil {
		return nil, err
	}
	pins := &namespacePins{handles: make([]windows.Handle, 0, len(paths))}
	for _, path := range paths {
		handle, err := openPinnedDirectory(path)
		if err != nil {
			return nil, errors.Join(fmt.Errorf("pin database namespace: %w", err), pins.close())
		}
		pins.handles = append(pins.handles, handle)
	}
	// A parent could have moved before its pin was established. Reopen only after
	// the whole chain is pinned and prove this namespace still denotes our file.
	current, err := CanonicalDatabase(identity.Path)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("recheck pinned database: %w", err), pins.close())
	}
	if current.ID != identity.ID {
		return nil, errors.Join(errors.New("database namespace changed while acquiring ownership"), pins.close())
	}
	return pins, nil
}
func openPinnedDirectory(path string) (windows.Handle, error) {
	pointer, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	handle, err := windows.CreateFile(pointer, windows.FILE_READ_ATTRIBUTES, windows.FILE_SHARE_READ, nil, windows.OPEN_EXISTING,
		windows.FILE_FLAG_BACKUP_SEMANTICS|windows.FILE_FLAG_OPEN_REPARSE_POINT, 0)
	if err != nil {
		return 0, err
	}
	if err = validatePinnedDirectory(handle, path); err != nil {
		return 0, errors.Join(err, windows.CloseHandle(handle))
	}
	return handle, nil
}
func validatePinnedDirectory(handle windows.Handle, path string) error {
	var attributes windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &attributes); err != nil {
		return err
	}
	if attributes.FileAttributes&windows.FILE_ATTRIBUTE_DIRECTORY == 0 || attributes.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
		return errors.New("canonical database ancestor is not an ordinary directory")
	}
	buffer := make([]uint16, 32768)
	n, err := windows.GetFinalPathNameByHandle(handle, &buffer[0], 32768, 0)
	if err != nil {
		return err
	}
	if n == 0 || n >= 32768 {
		return errors.New("invalid directory identity path")
	}
	canonical := windows.UTF16ToString(buffer[:n])
	if suffix, ok := strings.CutPrefix(canonical, `\\?\UNC\`); ok {
		canonical = `\\` + suffix
	} else {
		canonical = strings.TrimPrefix(canonical, `\\?\`)
	}
	if !strings.EqualFold(strings.TrimRight(canonical, `\`), strings.TrimRight(path, `\`)) {
		return errors.New("database ancestor namespace changed during pinning")
	}
	return nil
}
func (p *namespacePins) close() error {
	var result error
	for i := len(p.handles) - 1; i >= 0; i-- {
		result = errors.Join(result, windows.CloseHandle(p.handles[i]))
	}
	p.handles = nil
	return result
}
