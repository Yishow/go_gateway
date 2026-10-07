//go:build windows

package desktop

import (
	"os"

	"golang.org/x/sys/windows"
)

type peerIdentity struct {
	sid     string
	session uint32
}

func currentPeer() (peerIdentity, error) { return processPeer(windows.GetCurrentProcessId()) }
func processPeer(pid uint32) (peerIdentity, error) {
	var peer peerIdentity
	if e := windows.ProcessIdToSessionId(pid, &peer.session); e != nil {
		return peer, e
	}
	process, e := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if e != nil {
		return peer, e
	}
	defer closeTemporaryHandle(process)
	var token windows.Token
	if e := windows.OpenProcessToken(process, windows.TOKEN_QUERY, &token); e != nil {
		return peer, e
	}
	defer token.Close()
	user, e := token.GetTokenUser()
	if e != nil {
		return peer, e
	}
	peer.sid = user.User.Sid.String()
	return peer, nil
}
func samePeer(pid uint32, want peerIdentity) bool {
	got, e := processPeer(pid)
	return e == nil && got == want
}
func sameExecutable(pid uint32) bool {
	process, e := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if e != nil {
		return false
	}
	defer closeTemporaryHandle(process)
	buffer := make([]uint16, 32768)
	size := uint32(32768)
	if windows.QueryFullProcessImageName(process, 0, &buffer[0], &size) != nil {
		return false
	}
	remote, e := CanonicalDatabase(windows.UTF16ToString(buffer[:size]))
	if e != nil {
		return false
	}
	executable, e := os.Executable()
	if e != nil {
		return false
	}
	local, e := CanonicalDatabase(executable)
	return e == nil && local.ID == remote.ID
}
