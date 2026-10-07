//go:build windows

package desktop

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32              = windows.NewLazySystemDLL("user32.dll")
	shell32             = windows.NewLazySystemDLL("shell32.dll")
	ole32               = windows.NewLazySystemDLL("ole32.dll")
	registerClass       = user32.NewProc("RegisterClassExW")
	unregisterClass     = user32.NewProc("UnregisterClassW")
	createWindow        = user32.NewProc("CreateWindowExW")
	destroyWindow       = user32.NewProc("DestroyWindow")
	defaultWindowProc   = user32.NewProc("DefWindowProcW")
	getMessage          = user32.NewProc("GetMessageW")
	translateMessage    = user32.NewProc("TranslateMessage")
	dispatchMessage     = user32.NewProc("DispatchMessageW")
	postMessage         = user32.NewProc("PostMessageW")
	postQuit            = user32.NewProc("PostQuitMessage")
	registerMessage     = user32.NewProc("RegisterWindowMessageW")
	changeMessageFilter = user32.NewProc("ChangeWindowMessageFilterEx")
	createPopup         = user32.NewProc("CreatePopupMenu")
	appendMenu          = user32.NewProc("AppendMenuW")
	destroyMenu         = user32.NewProc("DestroyMenu")
	trackPopup          = user32.NewProc("TrackPopupMenu")
	foregroundWindow    = user32.NewProc("SetForegroundWindow")
	getCursorPosition   = user32.NewProc("GetCursorPos")
	getLastActivePopup  = user32.NewProc("GetLastActivePopup")
	createIcon          = user32.NewProc("CreateIconFromResourceEx")
	destroyIcon         = user32.NewProc("DestroyIcon")
	messageBox          = user32.NewProc("MessageBoxW")
	notifyIcon          = shell32.NewProc("Shell_NotifyIconW")
	shellExecute        = shell32.NewProc("ShellExecuteW")
	coInitialize        = ole32.NewProc("CoInitializeEx")
	coUninitialize      = ole32.NewProc("CoUninitialize")
)

const (
	wmBrowserError     = 0x8004
	wmDestroy          = 0x0002
	wmClose            = 0x0010
	wmQueryEndSession  = 0x0011
	wmEndSession       = 0x0016
	wmContextMenu      = 0x007b
	wmLeftDoubleClick  = 0x0203
	wmRightUp          = 0x0205
	wmTray             = 0x8001
	wmUpdate           = 0x8002
	wmTimeout          = 0x8003
	processStopTimeout = "stop-timeout"
)

type winPoint struct{ X, Y int32 }
type winMessage struct {
	Window         uintptr
	Message        uint32
	WParam, LParam uintptr
	Time           uint32
	Point          winPoint
	Private        uint32
}
type winClass struct {
	Size, Style                        uint32
	Procedure                          uintptr
	ClassExtra, WindowExtra            int32
	Instance, Icon, Cursor, Background uintptr
	MenuName, ClassName                *uint16
	SmallIcon                          uintptr
}
type notifyData struct {
	Size                uint32
	Window              uintptr
	ID, Flags, Callback uint32
	Icon                uintptr
	Tip                 [128]uint16
	State, StateMask    uint32
	Info                [256]uint16
	Version             uint32
	InfoTitle           [64]uint16
	InfoFlags           uint32
	GUID                windows.GUID
	BalloonIcon         uintptr
}

func utf16(text string) *uint16 { return windows.StringToUTF16Ptr(text) }
func winError(operation string, err error) error {
	if err == nil || errors.Is(err, windows.ERROR_SUCCESS) {
		return errors.New(operation + " failed")
	}
	return fmt.Errorf("%s: %w", operation, err)
}

// Win32's native result defines success; GetLastError is not meaningful on
// success or for procedures that do not document it. This directive preserves
// pointer lifetime and stack-movement safety through our uintptr forwarding.
//
//go:uintptrescapes
func winCall(proc *windows.LazyProc, args ...uintptr) uintptr {
	result, _, _ := proc.Call(args...) //nolint:errcheck // Callers inspect the native result rather than potentially stale LastError.
	return result
}
func showBox(owner uintptr, title, text string, flags uintptr) uintptr {
	return winCall(messageBox, owner, uintptr(unsafe.Pointer(utf16(text))), uintptr(unsafe.Pointer(utf16(title))), flags)
}

// ShowInfo presents native read-only information; callers gate this by desktop mode.
func ShowInfo(title, text string) { showBox(0, title, text, 0x40) }

// ShowError displays safe local diagnostics, never arbitrary error strings.
func ShowError(info ErrorInfo) { showBox(0, "Go Gateway", errorText(info), 0x10) }

// ConfirmCreateDB defaults to Cancel and performs no storage mutation itself.
func ConfirmCreateDB(path string) bool {
	return showBox(0, "Go Gateway", "在這個位置建立新資料庫？\n\n"+path+"\n\n這不會匯入或搬移歷史資料。", 0x21|0x100) == 1
}

// OpenURL calls the native shell from a COM STA without a command-shell process.
func OpenURL(raw string) error {
	if err := validateURL(raw); err != nil {
		return err
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	result := winCall(coInitialize, 0, 2|4)
	if int32(result) < 0 {
		return fmt.Errorf("native shell COM initialization failed (%x)", result)
	}
	defer winCall(coUninitialize)
	result = winCall(shellExecute, 0, uintptr(unsafe.Pointer(utf16("open"))), uintptr(unsafe.Pointer(utf16(raw))), 0, 0, 1)
	if result <= 32 {
		return fmt.Errorf("native browser open failed (%d)", result)
	}
	return nil
}
