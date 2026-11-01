//go:build !linux

package main

import "time"

func acceptanceProcessCPU() (time.Duration, bool) { return 0, false }
