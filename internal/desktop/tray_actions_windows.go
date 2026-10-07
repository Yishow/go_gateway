//go:build windows

package desktop

import (
	"cmp"
	"time"
	"unsafe"
)

func (s *nativeShell) popup() {
	menu := winCall(createPopup)
	if menu == 0 {
		return
	}
	defer winCall(destroyMenu, menu)
	s.mu.Lock()
	state, window := s.state, s.window
	s.mu.Unlock()
	for _, item := range menuItems(state) {
		flags := uintptr(0)
		if !item.Enabled {
			flags = 1
		}
		if winCall(appendMenu, menu, flags, item.ID, uintptr(unsafe.Pointer(utf16(item.Label)))) == 0 {
			showBox(window, "Go Gateway", "TRAY_MENU_FAILED\n無法建立完整系統列選單，請重試。", 0x10)
			return
		}
	}
	var point winPoint
	winCall(getCursorPosition, uintptr(unsafe.Pointer(&point)))
	winCall(foregroundWindow, window)
	command := winCall(trackPopup, menu, 0x100|0x2, uintptr(point.X), uintptr(point.Y), 0, window, 0)
	winCall(postMessage, window, 0, 0, 0)
	switch command {
	case 1:
		s.openSetup()
	case 2:
		if state.LogsURL != "" {
			s.open(state.LogsURL)
		} else {
			showBox(window, "Go Gateway", cmp.Or(state.LogsUnavailable, "本機執行日誌尚未就緒。"), 0x40)
		}
	case 3:
		showBox(window, "版本資訊", versionText(s.options), 0x40)
	case 4:
		s.confirmQuit()
	}
	data := s.iconData()
	winCall(notifyIcon, 3, uintptr(unsafe.Pointer(&data)))
}
func (s *nativeShell) openSetup() {
	s.mu.Lock()
	url := s.state.SetupURL
	s.mu.Unlock()
	if url != "" {
		s.open(url)
	}
}
func (s *nativeShell) open(url string) {
	s.mu.Lock()
	active := !s.closing && !s.uiClosing && s.window != 0
	s.mu.Unlock()
	if !active {
		return
	}
	s.browser.start(OpenURL, url, func() {
		// Post only while the original native window is still alive. WM_CLOSE marks
		// uiClosing before destruction; no worker ever calls into a disposed window.
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.closing && !s.uiClosing && s.window != 0 {
			winCall(postMessage, s.window, wmBrowserError, 0, 0)
		}
	})
}
func (s *nativeShell) confirmQuit() {
	s.mu.Lock()
	state, window := s.state, s.window
	s.mu.Unlock()
	if state.Process == processStopTimeout {
		s.confirmForce()
		return
	}
	if state.Process == processStopping {
		return
	}
	if showBox(window, "結束 Go Gateway", "確定停止採集及網頁服務並結束閘道？\n\n結束不代表所有待送資料已送達。", 0x21|0x100) != 1 {
		return
	}
	s.mu.Lock()
	if s.state.Process == processStopping || s.state.Process == processStopTimeout || s.sessionEnding || s.closing {
		s.mu.Unlock()
		return
	}
	s.state.Process = processStopping
	s.mu.Unlock()
	if s.options.OnQuit != nil {
		go s.options.OnQuit()
	}
}
func (s *nativeShell) confirmForce() {
	s.mu.Lock()
	state, window, sessionEnding := s.state, s.window, s.sessionEnding
	s.mu.Unlock()
	if state.Process != processStopTimeout || sessionEnding {
		return
	}
	text := "仍在關閉，尚未確認完成。\n等待階段：" + cmp.Or(state.PendingPhase, "unknown") + "\n\n按取消繼續等待原關閉程序。\n只有按確定才會強制結束；可能存在未完成或結果未知的資料，尚未持久保存的資料可能遺失。"
	if showBox(window, "強制結束 Go Gateway？", text, 0x31|0x100) != 1 {
		return
	}
	s.mu.Lock()
	active := s.state.Process == processStopTimeout && !s.sessionEnding && !s.closing
	s.mu.Unlock()
	if active && s.options.OnForceQuit != nil {
		go s.options.OnForceQuit()
	}
}
func (s *nativeShell) endSession() {
	s.mu.Lock()
	if s.sessionEnding {
		s.mu.Unlock()
		return
	}
	s.sessionEnding = true
	s.mu.Unlock()
	if s.options.OnSessionEnd == nil {
		return
	}
	go s.options.OnSessionEnd()
	// Windows can terminate after handlers return. Give the existing coordinator
	// bounded time, without claiming the OS grants a full shutdown deadline.
	deadline := time.NewTimer(4 * time.Second)
	defer deadline.Stop()
	select {
	case <-s.options.ShutdownDone:
	case <-deadline.C:
	}
}
