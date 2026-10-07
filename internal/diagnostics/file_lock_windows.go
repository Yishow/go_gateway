//go:build windows

package diagnostics

import (
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func lockDiagnosticFile(f *os.File) error {
	return windows.LockFileEx(windows.Handle(f.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &windows.Overlapped{Offset: 0xfffffff0, OffsetHigh: 0x7fffffff})
}
func singleLink(f *os.File) bool {
	var info windows.ByHandleFileInformation
	return windows.GetFileInformationByHandle(windows.Handle(f.Fd()), &info) == nil && info.NumberOfLinks == 1 && info.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT == 0
}
func reparse(info os.FileInfo) bool {
	if info.Mode()&os.ModeSymlink != 0 {
		return true
	}
	data, ok := info.Sys().(*syscall.Win32FileAttributeData)
	return ok && data.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0
}
