//go:build windows

package desktop

import (
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

//go:embed gateway.ico
var embeddedIcon []byte

func (s *nativeShell) initialize() error {
	if len(embeddedIcon) < 22 {
		return errors.New("embedded tray icon missing")
	}
	offset := int(binary.LittleEndian.Uint32(embeddedIcon[18:22]))
	length := int(binary.LittleEndian.Uint32(embeddedIcon[14:18]))
	if offset < 22 || length < 40 || offset > len(embeddedIcon)-length {
		return errors.New("embedded tray icon invalid")
	}
	icon, _, err := createIcon.Call(uintptr(unsafe.Pointer(&embeddedIcon[offset])), uintptr(length), 1, 0x00030000, 32, 32, 0)
	if icon == 0 {
		return winError("create tray icon", err)
	}
	s.icon = icon
	var module windows.Handle
	if err := windows.GetModuleHandleEx(2, nil, &module); err != nil {
		return err
	}
	s.module = uintptr(module)
	className := utf16(fmt.Sprintf("GoGatewayTray%d", trayNumber.Add(1)))
	class := winClass{Procedure: trayProcedure, Instance: s.module, Icon: icon, ClassName: className}
	class.Size = uint32(unsafe.Sizeof(class))
	atom, _, err := registerClass.Call(uintptr(unsafe.Pointer(&class)))
	if atom == 0 {
		return winError("register tray window", err)
	}
	s.classAtom = atom
	// A hidden top-level window receives shell broadcasts; HWND_MESSAGE does not.
	window, _, err := createWindow.Call(0x80, uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(utf16("Go Gateway"))), 0, 0, 0, 0, 0, 0, 0, s.module, 0)
	if window == 0 {
		return winError("create tray window", err)
	}
	s.mu.Lock()
	s.window = window
	s.mu.Unlock()
	trayWindows.Store(window, s)
	message := winCall(registerMessage, uintptr(unsafe.Pointer(utf16("TaskbarCreated"))))
	s.taskbarMessage = uint32(message)
	if message == 0 {
		return errors.New("register taskbar recovery message failed")
	}
	// This single fixed recovery message can cross UIPI from normal Explorer to an elevated gateway.
	if winCall(changeMessageFilter, window, message, 1, 0) == 0 {
		return errors.New("allow taskbar recovery message failed")
	}
	return s.addIcon()
}
func (s *nativeShell) iconData() notifyData {
	s.mu.Lock()
	state, window := s.state, s.window
	s.mu.Unlock()
	data := notifyData{Window: window, ID: 1, Flags: 1 | 2 | 4 | 0x80, Callback: wmTray, Icon: s.icon}
	data.Size = uint32(unsafe.Sizeof(data))
	copy(data.Tip[:len(data.Tip)-1], windows.StringToUTF16("Go Gateway | "+state.Process+" | "+state.Acquisition))
	return data
}
func (s *nativeShell) addIcon() error {
	data := s.iconData()
	result, _, err := notifyIcon.Call(0, uintptr(unsafe.Pointer(&data)))
	if result == 0 {
		return winError("register notification icon", err)
	}
	data.Version = 4
	result, _, err = notifyIcon.Call(4, uintptr(unsafe.Pointer(&data)))
	if result == 0 {
		winCall(notifyIcon, 2, uintptr(unsafe.Pointer(&data)))
		return winError("set notification icon version", err)
	}
	return nil
}
func (s *nativeShell) dispose() {
	s.mu.Lock()
	window := s.window
	s.uiClosing = true
	s.window = 0
	s.mu.Unlock()
	if window != 0 {
		data := notifyData{Window: window, ID: 1}
		data.Size = uint32(unsafe.Sizeof(data))
		winCall(notifyIcon, 2, uintptr(unsafe.Pointer(&data)))
		trayWindows.Delete(window)
		winCall(destroyWindow, window)
	}
	if s.classAtom != 0 {
		winCall(unregisterClass, s.classAtom, s.module)
		s.classAtom = 0
	}
	if s.icon != 0 {
		winCall(destroyIcon, s.icon)
		s.icon = 0
	}
}
