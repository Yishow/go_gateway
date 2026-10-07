//go:build !linux && !darwin && !freebsd && !openbsd && !netbsd && !dragonfly && !windows

package diagnostics

import "os"

func lockDiagnosticFile(*os.File) error { return ErrFileUnsafe }
func singleLink(*os.File) bool          { return false }
func reparse(info os.FileInfo) bool     { return info.Mode()&os.ModeSymlink != 0 }
