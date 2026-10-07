//go:build windows

package desktop

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const setupOpcode byte = 1

type setupPipe struct {
	handle    windows.Handle
	stop      chan struct{}
	wg        sync.WaitGroup
	busy      atomic.Bool
	closeOnce func() error
}

func pipeName(identity string, session uint32) *uint16 {
	return windows.StringToUTF16Ptr(fmt.Sprintf(`\\.\pipe\GoGateway.Setup.%s.%d`, identity, session))
}
func newSetupPipe(identity string, openSetup func()) (*setupPipe, error) {
	if openSetup == nil {
		return nil, errors.New("setup callback is required")
	}
	peer, e := currentPeer()
	if e != nil {
		return nil, e
	}
	descriptor, e := windows.SecurityDescriptorFromString("D:P(A;;GA;;;" + peer.sid + ")")
	if e != nil {
		return nil, e
	}
	attributes := windows.SecurityAttributes{Length: uint32(unsafe.Sizeof(windows.SecurityAttributes{})), SecurityDescriptor: descriptor}
	handle, e := windows.CreateNamedPipe(pipeName(identity, peer.session), windows.PIPE_ACCESS_DUPLEX|windows.FILE_FLAG_FIRST_PIPE_INSTANCE,
		windows.PIPE_TYPE_MESSAGE|windows.PIPE_READMODE_MESSAGE|windows.PIPE_NOWAIT|windows.PIPE_REJECT_REMOTE_CLIENTS, 1, 16, 16, 0, &attributes)
	if e != nil {
		return nil, e
	}
	pipe := &setupPipe{handle: handle, stop: make(chan struct{})}
	pipe.closeOnce = sync.OnceValue(func() error { close(pipe.stop); pipe.wg.Wait(); return windows.CloseHandle(handle) })
	pipe.wg.Go(func() { pipe.serve(peer, openSetup) })
	return pipe, nil
}
func (p *setupPipe) close() error { return p.closeOnce() }
func (p *setupPipe) serve(peer peerIdentity, openSetup func()) {
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-p.stop:
			return
		case <-ticker.C:
		}
		err := windows.ConnectNamedPipe(p.handle, nil)
		// In nonblocking mode, success means entering listening state, not connected.
		if err == nil || errors.Is(err, windows.ERROR_PIPE_LISTENING) {
			continue
		}
		if errors.Is(err, windows.ERROR_PIPE_CONNECTED) {
			p.accept(peer, openSetup)
		}
		if !disconnectPipe(p.handle) {
			return
		}
	}
}
func (p *setupPipe) accept(peer peerIdentity, openSetup func()) {
	var pid uint32
	if windows.GetNamedPipeClientProcessId(p.handle, &pid) != nil || !samePeer(pid, peer) {
		return
	}
	command, e := readPipeMessage(p.handle, p.stop, 500*time.Millisecond)
	if e != nil || len(command) != 1 || command[0] != setupOpcode {
		return
	}
	// Keep shell callbacks bounded even if the user's browser opener stalls.
	if p.busy.CompareAndSwap(false, true) {
		go func() { defer p.busy.Store(false); openSetup() }()
	}
	var n uint32
	if windows.WriteFile(p.handle, []byte{setupOpcode}, &n, nil) != nil {
		return
	}
	// DisconnectNamedPipe discards unread bytes. Await a bounded drain acknowledgement.
	if _, e = readPipeMessage(p.handle, p.stop, 500*time.Millisecond); e != nil {
		return
	}
}
func disconnectPipe(handle windows.Handle) bool {
	err := windows.DisconnectNamedPipe(handle)
	return err == nil || errors.Is(err, windows.ERROR_PIPE_NOT_CONNECTED)
}
func readPipeMessage(handle windows.Handle, stop <-chan struct{}, timeout time.Duration) ([]byte, error) {
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		var n uint32
		buffer := make([]byte, 2)
		err := windows.ReadFile(handle, buffer, &n, nil)
		if err == nil && n > 0 {
			return buffer[:n], nil
		}
		if err != nil && !errors.Is(err, windows.ERROR_NO_DATA) {
			return nil, err
		}
		select {
		case <-stop:
			return nil, errors.New("handoff stopped")
		case <-deadline.C:
			return nil, errors.New("handoff timed out")
		case <-ticker.C:
		}
	}
}

// RequestOpenSetup sends no URL or command string, only the fixed setup opcode.
func RequestOpenSetup(path string) error {
	identity, e := CanonicalDatabase(path)
	if e != nil {
		return e
	}
	peer, e := currentPeer()
	if e != nil {
		return e
	}
	handle, e := windows.CreateFile(pipeName(identity.ID, peer.session), windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil, windows.OPEN_EXISTING, windows.SECURITY_SQOS_PRESENT|windows.SECURITY_IDENTIFICATION, 0)
	if e != nil {
		return fmt.Errorf("setup handoff unavailable: %w", e)
	}
	defer closeTemporaryHandle(handle)
	var pid uint32
	if windows.GetNamedPipeServerProcessId(handle, &pid) != nil || !samePeer(pid, peer) || !sameExecutable(pid) {
		return errors.New("setup owner identity not verified")
	}
	mode := uint32(windows.PIPE_READMODE_MESSAGE | windows.PIPE_NOWAIT)
	if err := windows.SetNamedPipeHandleState(handle, &mode, nil, nil); err != nil {
		return err
	}
	var n uint32
	if err := windows.WriteFile(handle, []byte{setupOpcode}, &n, nil); err != nil {
		return err
	}
	ack, e := readPipeMessage(handle, nil, time.Second)
	if e != nil {
		return e
	}
	if len(ack) != 1 || ack[0] != setupOpcode {
		return errors.New("invalid setup acknowledgement")
	}
	return windows.WriteFile(handle, []byte{2}, &n, nil)
}
