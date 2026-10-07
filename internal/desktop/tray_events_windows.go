//go:build windows

package desktop

import "unsafe"

func trayWindowProcedure(window uintptr, message uint32, wparam, lparam uintptr) uintptr {
	value, exists := trayWindows.Load(window)
	if !exists {
		return winCall(defaultWindowProc, window, uintptr(message), wparam, lparam)
	}
	shell, valid := value.(*nativeShell)
	if !valid {
		return winCall(defaultWindowProc, window, uintptr(message), wparam, lparam)
	}
	if message == shell.taskbarMessage {
		if err := shell.addIcon(); err != nil {
			shell.fail("TRAY_RESTORE_FAILED", "無法恢復系統列圖示。")
		}
		return 0
	}
	switch message {
	case wmTray:
		switch uint32(lparam) & 0xffff {
		case wmContextMenu, wmRightUp:
			shell.popup()
		case wmLeftDoubleClick, 0x401:
			shell.openSetup()
		}
		return 0
	case wmUpdate:
		data := shell.iconData()
		winCall(notifyIcon, 1, uintptr(unsafe.Pointer(&data)))
		return 0
	case wmBrowserError:
		shell.mu.Lock()
		active := !shell.closing && !shell.uiClosing
		shell.mu.Unlock()
		if active {
			showBox(window, "Go Gateway", "BROWSER_OPEN_FAILED\n無法開啟瀏覽器。\n請檢查預設瀏覽器設定後再試一次。", 0x10)
		}
		return 0
	case wmTimeout:
		shell.confirmForce()
		return 0
	case wmQueryEndSession:
		// Other applications can still veto; a query is not committed shutdown.
		return 1
	case wmEndSession:
		if wparam != 0 {
			shell.endSession()
		}
		return 0
	case wmClose:
		shell.mu.Lock()
		shell.uiClosing = true
		shell.mu.Unlock()
		popup := winCall(getLastActivePopup, window)
		if popup != 0 && popup != window {
			winCall(postMessage, popup, wmClose, 0, 0)
		}
		winCall(destroyWindow, window)
		return 0
	case wmDestroy:
		shell.mu.Lock()
		shell.uiClosing = true
		shell.mu.Unlock()
		winCall(postQuit, 0)
		return 0
	}
	return winCall(defaultWindowProc, window, uintptr(message), wparam, lparam)
}
