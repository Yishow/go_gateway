//go:build windows

package desktop

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"golang.org/x/sys/windows"
)

var trayWindows sync.Map
var trayNumber atomic.Uint64
var trayProcedure = windows.NewCallback(trayWindowProcedure)

type nativeShell struct {
	browser                                       browserTask
	uiClosing                                     bool
	options                                       Options
	mu                                            sync.Mutex
	state                                         State
	window, icon, classAtom, module               uintptr
	taskbarMessage                                uint32
	started, closing, timeoutShown, sessionEnding bool
	ready                                         chan error
	done                                          chan struct{}
}

// New constructs the native boundary without starting its message loop.
func New(options Options) Shell {
	return &nativeShell{options: options, state: State{Process: "starting", Acquisition: "unknown"}, ready: make(chan error, 1), done: make(chan struct{})}
}
func (s *nativeShell) Start() error {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return errors.New("tray already started")
	}
	s.started = true
	s.mu.Unlock()
	go s.run()
	return <-s.ready
}
func (s *nativeShell) Update(state State) {
	s.mu.Lock()
	s.state = state
	window := s.window
	timeout := state.Process == processStopTimeout && !s.timeoutShown
	if timeout {
		s.timeoutShown = true
	}
	s.mu.Unlock()
	if window != 0 {
		winCall(postMessage, window, wmUpdate, 0, 0)
		if timeout {
			winCall(postMessage, window, wmTimeout, 0, 0)
		}
	}
}
func (s *nativeShell) Close() error {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return nil
	}
	s.closing = true
	window := s.window
	s.mu.Unlock()
	if window != 0 && winCall(postMessage, window, wmClose, 0, 0) == 0 {
		return errors.New("request native tray closure failed")
	}
	<-s.done
	return nil
}
func (s *nativeShell) run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	announced := false
	defer func() {
		if recover() != nil {
			if !announced {
				s.ready <- errors.New("native tray initialization failed")
			} else {
				s.fail("TRAY_PANIC", "系統列操作發生錯誤。")
			}
		}
		s.dispose()
		close(s.done)
	}()
	if err := s.initialize(); err != nil {
		s.ready <- err
		return
	}
	announced = true
	s.ready <- nil
	s.mu.Lock()
	closing, window := s.closing, s.window
	s.mu.Unlock()
	if closing {
		winCall(postMessage, window, wmClose, 0, 0)
	}
	var message winMessage
	for {
		result := winCall(getMessage, uintptr(unsafe.Pointer(&message)), 0, 0, 0)
		if result == 0 {
			break
		}
		if int32(result) == -1 {
			s.fail("TRAY_MESSAGE_LOOP", "系統列訊息處理失敗。")
			break
		}
		winCall(translateMessage, uintptr(unsafe.Pointer(&message)))
		winCall(dispatchMessage, uintptr(unsafe.Pointer(&message)))
	}
	s.mu.Lock()
	closing = s.closing
	s.mu.Unlock()
	if !closing && s.options.OnQuit != nil {
		go s.options.OnQuit()
	}
}
func (s *nativeShell) fail(code, message string) {
	requestFaultStop(s.options)
	s.mu.Lock()
	window := s.window
	s.mu.Unlock()
	showBox(window, "Go Gateway", errorText(ErrorInfo{Code: code, Message: message, NextStep: "閘道將安全關閉，請重新啟動。"}), 0x10)
}
