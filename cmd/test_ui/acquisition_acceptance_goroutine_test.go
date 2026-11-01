package main

import (
	"bytes"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

// Diagnostic only: no extra workers, wait/drain policy, or leak assertion.
// Register before every fixture cleanup so the last count follows their joins.
func startAcceptanceGoroutineDiagnostic(t *testing.T) *acceptanceGoroutineDiagnostic {
	t.Helper()
	d := &acceptanceGoroutineDiagnostic{t: t}
	d.count("before-fixture", "no fixture owners created")
	t.Cleanup(func() {
		d.count("after-owned-cleanups", "scheduler/runtime stopped; sample drain and ValueEvent reader joined; manager/DB closed; peer listener/connections joined; global/test workers may remain")
	})
	return d
}

type acceptanceGoroutineDiagnostic struct {
	t     *testing.T
	taken bool
}

func (d *acceptanceGoroutineDiagnostic) count(stage, owners string) {
	d.t.Helper()
	d.t.Logf("goroutine diagnostic stage=%s count=%d owners=%s; count alone is not leak evidence", stage, runtime.NumGoroutine(), owners)
}

func (d *acceptanceGoroutineDiagnostic) observe(count int) {
	if d.taken || count < 256 {
		return
	}
	d.taken = true
	started := time.Now()
	// Aggregate identical stacks; cap retained text at 1MiB. Write still reports
	// the full input length so the profile can finish without changing runtime.
	output := &acceptanceStackText{}
	err := pprof.Lookup("goroutine").WriteTo(output, 1)
	d.t.Logf("goroutine diagnostic first-threshold=256 sampled_count=%d capture_elapsed=%s retained_bytes=%d truncated=%t error=%v; one profile included in this case CPU/elapsed, not a CPU-causality percentage\n%s", count, time.Since(started), output.Len(), output.truncated, err, output.String())
}

type acceptanceStackText struct {
	bytes.Buffer
	truncated bool
}

func (b *acceptanceStackText) Write(p []byte) (int, error) {
	n := len(p)
	remaining := (1 << 20) - b.Len()
	if len(p) > remaining {
		b.truncated = true
		p = p[:remaining]
	}
	_, _ = b.Buffer.Write(p)
	return n, nil
}
