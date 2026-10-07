//go:build windows

package desktop

import (
	"errors"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows"
)

func TestOwnerRejectsAmbiguousFinalReparseTarget(t *testing.T) {
	for _, suffix := range []string{".", " "} {
		t.Run(suffix, func(t *testing.T) {
			directory := t.TempDir()
			target := `\\?\` + filepath.Join(directory, "target.db") + suffix
			targetPtr := windows.StringToUTF16Ptr(target)
			handle, err := windows.CreateFile(targetPtr, windows.GENERIC_READ|windows.GENERIC_WRITE, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, nil, windows.CREATE_NEW, windows.FILE_ATTRIBUTE_NORMAL, 0)
			if err != nil {
				t.Fatal(err)
			}
			if err = windows.CloseHandle(handle); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := windows.DeleteFile(targetPtr); err != nil {
					t.Errorf("remove test target: %v", err)
				}
			})
			alias := filepath.Join(directory, "alias.db")
			err = windows.CreateSymbolicLink(windows.StringToUTF16Ptr(alias), targetPtr, 2)
			if errors.Is(err, windows.ERROR_PRIVILEGE_NOT_HELD) {
				t.Skip("native symlink test requires permitted symlink creation")
			}
			if err != nil {
				t.Fatal(err)
			}
			owner, err := AcquireOwner(alias, false, nil)
			if owner != nil {
				_ = owner.Close()
			}
			if err == nil {
				t.Fatal("reparse target changed meaning when converted to ordinary path")
			}
		})
	}
}
