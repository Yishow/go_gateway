//go:build windows

package desktop

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Owner retains the database's actual file handle and optional local setup IPC.
type Owner struct {
	identity  Identity
	handle    windows.Handle
	pipe      *setupPipe
	namespace *namespacePins
	closeOnce func() error
}

// AcquireOwner may reserve an empty missing database only after caller approval.
// The beyond-EOF byte lock neither writes nor extends the SQLite file.
func AcquireOwner(path string, desktopMode bool, openSetup func()) (*Owner, error) {
	h, e := openIdentityHandle(path, windows.OPEN_ALWAYS, windows.GENERIC_READ|windows.GENERIC_WRITE)
	if e != nil {
		return nil, fmt.Errorf("open ownership handle: %w", e)
	}
	offset := ownerLockOffset()
	if e = windows.LockFileEx(h, windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &offset); e != nil {
		closeTemporaryHandle(h)
		if errors.Is(e, windows.ERROR_LOCK_VIOLATION) {
			return nil, ErrAlreadyRunning
		}
		return nil, fmt.Errorf("lock database: %w", e)
	}
	// Active ownership wins before link policy, allowing a duplicate launch to
	// report the actual owner without treating metadata as ownership authority.
	if e := validateWritableDatabase(h); e != nil {
		closeTemporaryHandle(h)
		return nil, e
	}
	identity, e := handleIdentity(h)
	if e != nil {
		closeTemporaryHandle(h)
		return nil, e
	}
	pins, e := pinOwnerNamespace(identity)
	if e != nil {
		closeTemporaryHandle(h)
		return nil, e
	}
	owner := &Owner{identity: identity, handle: h, namespace: pins}
	owner.closeOnce = sync.OnceValue(func() error {
		var pipeErr error
		if owner.pipe != nil {
			pipeErr = owner.pipe.close()
		}
		offset := ownerLockOffset()
		return errors.Join(pipeErr, windows.UnlockFileEx(h, 0, 1, 0, &offset), windows.CloseHandle(h), owner.namespace.close())
	})
	if desktopMode {
		owner.pipe, e = newSetupPipe(identity.ID, openSetup)
		if e != nil {
			return nil, errors.Join(e, owner.Close())
		}
	}
	return owner, nil
}

// Identity returns the locked file identity, never a pre-creation path identity.
func (o *Owner) Identity() Identity { return o.identity }

// Close releases ownership only after actual service cleanup, never merely a timeout.
func (o *Owner) Close() error { return o.closeOnce() }

// CanonicalDatabase obtains trusted existing-file identity without creating storage.
func CanonicalDatabase(path string) (Identity, error) {
	h, e := openIdentityHandle(path, windows.OPEN_EXISTING, windows.FILE_READ_ATTRIBUTES)
	if e != nil {
		return Identity{}, e
	}
	defer closeTemporaryHandle(h)
	return handleIdentity(h)
}
func openIdentityHandle(path string, creation, access uint32) (windows.Handle, error) {
	if err := validateWindowsFilePath(path); err != nil {
		return 0, err
	}
	p, e := windows.UTF16PtrFromString(strings.ReplaceAll(path, "/", `\`))
	if e != nil {
		return 0, e
	}
	// No delete sharing: the held file cannot be replaced underneath its owner.
	return windows.CreateFile(p, access, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, nil, creation, windows.FILE_ATTRIBUTE_NORMAL, 0)
}
func ownerLockOffset() windows.Overlapped {
	return windows.Overlapped{Offset: 0xfffffffe, OffsetHigh: 0x7fffffff}
}
func handleIdentity(h windows.Handle) (Identity, error) {
	var attributes windows.ByHandleFileInformation
	if e := windows.GetFileInformationByHandle(h, &attributes); e != nil {
		return Identity{}, e
	}
	if attributes.FileAttributes&windows.FILE_ATTRIBUTE_DIRECTORY != 0 {
		return Identity{}, errors.New("database must be a regular file")
	}
	// ReFS does not guarantee unique 64-bit legacy file indexes. Require FileIdInfo.
	var fileID struct {
		Volume uint64
		ID     [16]byte
	}
	if e := windows.GetFileInformationByHandleEx(h, windows.FileIdInfo, (*byte)(unsafe.Pointer(&fileID)), uint32(unsafe.Sizeof(fileID))); e != nil {
		return Identity{}, fmt.Errorf("trusted file identity unavailable: %w", e)
	}
	if fileID.ID == [16]byte{} {
		return Identity{}, errors.New("empty file identity")
	}
	path := make([]uint16, 32768)
	n, e := windows.GetFinalPathNameByHandle(h, &path[0], 32768, 0)
	if e != nil {
		return Identity{}, e
	}
	if n == 0 || n >= 32768 {
		return Identity{}, errors.New("invalid canonical path")
	}
	canonical, e := normalizeFinalWindowsPath(windows.UTF16ToString(path[:n]))
	if e != nil {
		return Identity{}, fmt.Errorf("untrusted final database path: %w", e)
	}
	return Identity{Path: canonical, ID: opaqueID(fmt.Sprintf("windows:%016x:%x", fileID.Volume, fileID.ID))}, nil
}

// Temporary query-handle cleanup cannot replace the primary operation's result.
func closeTemporaryHandle(h windows.Handle) {
	_ = windows.CloseHandle(h) //nolint:errcheck // Best-effort cleanup of a temporary handle; owner cleanup propagates errors.
}

func validateWritableDatabase(handle windows.Handle) error {
	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &info); err != nil {
		return err
	}
	return checkWritableLinkCount(info.NumberOfLinks)
}
