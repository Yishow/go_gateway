package diagnostics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestRotationIsThreeBoundedOwnedFiles(t *testing.T) {
	o := fileOptions(t)
	files, err := openOwnedFiles(o)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(Event{Code: "startup.ready", Message: strings.Repeat("safe ", 1400)})
	data = append(data, '\n')
	for range 2400 {
		if err := files.append(data); err != nil {
			files.close()
			t.Fatal(err)
		}
	}
	var total int64
	for _, f := range files.files {
		info, err := f.Stat()
		if err != nil {
			t.Fatal(err)
		}
		if info.Size() > MaxFileBytes {
			t.Fatal("file exceeded5MiB")
		}
		total += info.Size()
	}
	if total > 3*MaxFileBytes {
		t.Fatal("set exceeded15MiB")
	}
	files.close()
	again, err := openOwnedFiles(o)
	if err != nil {
		t.Fatal(err)
	}
	again.close()
}
func TestDiskFailureRecoveryAndSafeMemoryContinue(t *testing.T) {
	o := fileOptions(t)
	var calls atomic.Int32
	var sink *FileSink
	o.write = func(data []byte) error {
		if calls.Add(1) == 1 {
			return errors.New("disk full with secret-private-path")
		}
		return sink.owned.append(data)
	}
	var err error
	sink, err = OpenFileSink(o)
	if err != nil {
		t.Fatal(err)
	}
	b, err := New(Options{File: sink})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = b.Close(t.Context()) }()
	b.Emit(Input{Code: "startup.begin"})
	if err := b.Flush(t.Context()); err == nil {
		t.Fatal("failure hidden")
	}
	snap, _ := b.Snapshot(Query{})
	if len(snap.Records) != 1 || snap.SinkHealth != "degraded" || snap.DiskErrors != 1 {
		t.Fatal("memory affected or disk loss hidden")
	}
	b.Emit(Input{Code: "startup.ready"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	status := sink.Status()
	if !status.Saved || status.Degraded || status.RecoveredGaps != 1 || status.Dropped != 1 {
		t.Fatalf("recovery not declared %#v", status)
	}
	snap, _ = b.Snapshot(Query{})
	found := false
	for _, g := range snap.Gaps {
		found = found || g.Reason == "file_loss" && g.From == b.cursor(1) && g.To == b.cursor(1)
	}
	if !found {
		t.Fatal("lost interval missing")
	}
	raw, _ := os.ReadFile(o.Path)
	if bytes.Contains(raw, []byte("secret-private-path")) {
		t.Fatal("raw disk error entered file")
	}
}
func TestBlockedWriterCannotBlockCaptureAndCloseDeadline(t *testing.T) {
	o := fileOptions(t)
	entered := make(chan struct{})
	release := make(chan struct{})
	var first atomic.Bool
	var sink *FileSink
	o.write = func(data []byte) error {
		if first.CompareAndSwap(false, true) {
			close(entered)
			<-release
		}
		return sink.owned.append(data)
	}
	var err error
	sink, err = OpenFileSink(o)
	if err != nil {
		t.Fatal(err)
	}
	b, err := New(Options{File: sink})
	if err != nil {
		t.Fatal(err)
	}
	b.Emit(Input{Code: "startup.begin"})
	<-entered
	for range 18 {
		for range 100 {
			b.Emit(Input{Code: "startup.ready"})
		}
		for {
			b.mu.Lock()
			done := b.queue.count == 0
			b.mu.Unlock()
			if done {
				break
			}
			time.Sleep(time.Millisecond)
		}
	}
	snap, err := b.Snapshot(Query{})
	if err != nil || len(snap.Records) == 0 || sink.health().Dropped == 0 {
		t.Fatal("file stall blocked capture or queue drop hidden")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cancel()
	if err := b.Close(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("close claimed completion during disk stall: %v", err)
	}
	close(release)
	if err := b.Close(t.Context()); err != nil {
		t.Fatal(err)
	}
}
func TestPathSwapCannotRedirectHeldWritesOrLieAboutSavedPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires platform-specific Windows pathname replacement fixture")
	}
	o := fileOptions(t)
	sink, err := OpenFileSink(o)
	if err != nil {
		t.Fatal(err)
	}
	b, err := New(Options{File: sink})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = b.Close(t.Context()) }()
	b.Emit(Input{Code: "startup.begin"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	if !sink.Status().Saved {
		t.Fatal("initial event not synced")
	}
	if err := os.Rename(o.Path, o.Path+".moved"); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(o.DatabasePath, o.Path); err != nil {
		t.Fatal(err)
	}
	b.Emit(Input{Code: "startup.ready"})
	if err := b.Flush(t.Context()); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(o.DatabasePath)
	if string(got) != "protected database" {
		t.Fatal("swapped pathname redirected write")
	}
	if s := sink.Status(); s.Saved || s.Path != "" {
		t.Fatal("native path falsely claimed saved")
	}
}
func TestUnusableDirectoryAndMissingProtectedBackupRejected(t *testing.T) {
	o := fileOptions(t)
	o.Path = filepath.Join(o.DatabasePath, "runtime.jsonl")
	if sink, err := OpenFileSink(o); err == nil {
		_ = sink.Close(t.Context())
		t.Fatal("non-directory accepted")
	}
	o = fileOptions(t)
	o.ProtectedPaths = []string{o.Path + ".1"}
	if sink, err := OpenFileSink(o); err == nil {
		_ = sink.Close(t.Context())
		t.Fatal("protected nonexistent backup accepted")
	}
}
func TestReadOnlyDirectoryFailsWithoutFallback(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires Windows ACL fixture")
	}
	o := fileOptions(t)
	dir := filepath.Dir(o.Path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
	sink, err := OpenFileSink(o)
	if err == nil {
		_ = sink.Close(t.Context())
		if os.Geteuid() == 0 {
			t.Skip("root bypasses permission fixture")
		}
		t.Fatal("read-only location accepted")
	}
	if _, err := os.Stat(o.Path); !os.IsNotExist(err) {
		t.Fatal("file created despite failed startup")
	}
}

func TestFileLossIntervalOrdersAdmissionDropsAndDelayedDiskFailure(t *testing.T) {
	f := &FileSink{}
	instance := strings.Repeat("a", 32)
	f.lossLocked(record{event: Event{InstanceID: instance, Sequence: "20"}, sequence: 20})
	f.lossLocked(record{event: Event{InstanceID: instance, Sequence: "2"}, sequence: 2})
	if f.status.LostFrom != instance+":2" || f.status.LostTo != instance+":20" {
		t.Fatal("delayed older disk error inverted loss interval")
	}
}
