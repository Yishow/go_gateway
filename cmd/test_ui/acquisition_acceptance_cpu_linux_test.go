//go:build linux

package main

import (
	"syscall"
	"time"
)

// Process CPU is not elapsed RTT, per-device CPU, or planner-only CPU. This
// fixture snapshots the test process after setup, outside the compiler process.
func acceptanceProcessCPU() (time.Duration, bool) {
	var usage syscall.Rusage
	if syscall.Getrusage(syscall.RUSAGE_SELF, &usage) != nil {
		return 0, false
	}
	seconds := usage.Utime.Sec + usage.Stime.Sec
	micros := usage.Utime.Usec + usage.Stime.Usec
	return time.Duration(seconds)*time.Second + time.Duration(micros)*time.Microsecond, true
}
