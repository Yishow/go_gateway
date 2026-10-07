//go:build linux || darwin || freebsd || openbsd || netbsd || dragonfly

package diagnostics

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func lockDiagnosticFile(f *os.File) error { return unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB) }
func singleLink(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	return ok && stat.Nlink == 1
}
func reparse(info os.FileInfo) bool { return info.Mode()&os.ModeSymlink != 0 }
