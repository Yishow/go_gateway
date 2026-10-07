//go:build windows

package desktop

import (
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sys/windows"
)

func TestRepeatedDesktopHandoff(t *testing.T) {
	path := filepath.Join(t.TempDir(), "desktop.db")
	calls := make(chan struct{}, 3)
	owner, err := AcquireOwner(path, true, func() { calls <- struct{}{} })
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close()
	for range 3 {
		if err = RequestOpenSetup(path); err != nil {
			t.Fatal(err)
		}
		select {
		case <-calls:
		case <-time.After(time.Second):
			t.Fatal("setup callback missing")
		}
		time.Sleep(50 * time.Millisecond)
	}
}
func TestMalformedOrIdleClientCannotControlGateway(t *testing.T) {
	path := filepath.Join(t.TempDir(), "malformed.db")
	var calls atomic.Int32
	owner, err := AcquireOwner(path, true, func() { calls.Add(1) })
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close()
	peer, err := currentPeer()
	if err != nil {
		t.Fatal(err)
	}
	for _, command := range [][]byte{{2}, {1, 1}, []byte("exit"), nil} {
		handle, err := windows.CreateFile(pipeName(owner.Identity().ID, peer.session), windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil, windows.OPEN_EXISTING, windows.SECURITY_SQOS_PRESENT|windows.SECURITY_IDENTIFICATION, 0)
		if err != nil {
			t.Fatal(err)
		}
		if command != nil {
			var n uint32
			if err = windows.WriteFile(handle, command, &n, nil); err != nil {
				t.Fatal(err)
			}
		}
		time.Sleep(700 * time.Millisecond)
		if err = windows.CloseHandle(handle); err != nil {
			t.Fatal(err)
		}
		if calls.Load() != 0 {
			t.Fatal("malformed client invoked setup")
		}
	}
	if err = RequestOpenSetup(path); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for calls.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if calls.Load() != 1 {
		t.Fatal("pipe did not recover after rejected clients")
	}
}
func TestPeerRequiresUserAndSession(t *testing.T) {
	peer, err := currentPeer()
	if err != nil {
		t.Fatal(err)
	}
	pid := windows.GetCurrentProcessId()
	if !samePeer(pid, peer) || !sameExecutable(pid) {
		t.Fatal("self identity not verified")
	}
	different := peer
	different.sid = "S-1-0-0"
	if samePeer(pid, different) {
		t.Fatal("different SID accepted")
	}
	different = peer
	different.session++
	if samePeer(pid, different) {
		t.Fatal("different session accepted")
	}
}
func TestSessionQueryAndCancellationDoNotStopGateway(t *testing.T) {
	calls := make(chan struct{}, 1)
	done := make(chan struct{})
	close(done)
	shell, ok := New(Options{OnSessionEnd: func() { calls <- struct{}{} }, ShutdownDone: done}).(*nativeShell)
	if !ok {
		t.Fatal("native adapter missing")
	}
	const window = uintptr(0x1234)
	trayWindows.Store(window, shell)
	defer trayWindows.Delete(window)
	if trayWindowProcedure(window, wmQueryEndSession, 0, 0) != 1 {
		t.Fatal("session query rejected")
	}
	trayWindowProcedure(window, wmEndSession, 0, 0)
	select {
	case <-calls:
		t.Fatal("canceled session stopped gateway")
	default:
	}
	trayWindowProcedure(window, wmEndSession, 1, 0)
	trayWindowProcedure(window, wmEndSession, 1, 0)
	select {
	case <-calls:
	case <-time.After(time.Second):
		t.Fatal("committed session not forwarded")
	}
	select {
	case <-calls:
		t.Fatal("session coordinator invoked twice")
	default:
	}
}
