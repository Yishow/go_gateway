// Package desktop provides native Windows operations behind an injectable lifecycle boundary.
package desktop

import (
	"cmp"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
)

// Version and Commit are release linker inputs; development builds never invent a release.
var Version = "dev"
var Commit = "unknown"

// ErrAlreadyRunning means an OS-backed owner prevented database acquisition.
var ErrAlreadyRunning = errors.New("database already owned")

// ErrUnsupported means the requested native operation is unavailable.
var ErrUnsupported = errors.New("native desktop unavailable")

// Identity describes the actual database held by the owner guard.
type Identity struct{ Path, ID string }

// Options contains native-only information and callbacks to the single lifecycle coordinator.
type Options struct {
	Version, Commit, BuildMode, DataPath string
	OnQuit, OnSessionEnd, OnForceQuit    func()
	// OnFault records a fixed safe terminal shell failure before shutdown begins.
	OnFault func()
	// ShutdownDone closes after service cleanup, before Shell.Close, to avoid UI deadlock.
	ShutdownDone <-chan struct{}
}

// State projects existing lifecycle/acquisition authorities without controlling them.
type State struct{ Process, Acquisition, SetupURL, LogsURL, LogsUnavailable, PendingPhase string }

// Shell is the injectable native UI boundary. Start must succeed before services start.
type Shell interface {
	Start() error
	Update(State)
	Close() error
}

// ErrorInfo contains safe local text and explicitly verified diagnostic-save status.
type ErrorInfo struct {
	Code, Message, NextStep, DiagnosticPath string
	DiagnosticSaved                         bool
}

const processStopping = "stopping"

func opaqueID(s string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(s))) }
func versionText(o Options) string {
	text := fmt.Sprintf("Go Gateway\nVersion: %s\nCommit: %s\nBuild mode: %s", cmp.Or(o.Version, Version, "dev"), cmp.Or(o.Commit, Commit, "unknown"), cmp.Or(o.BuildMode, "unknown/dev"))
	if o.DataPath != "" {
		text += "\n資料來源：" + o.DataPath
	}
	return text
}
func errorText(info ErrorInfo) string {
	text := info.Code + "\n" + info.Message + "\n\n" + info.NextStep
	if info.DiagnosticSaved && info.DiagnosticPath != "" {
		return text + "\n\n診斷紀錄：" + info.DiagnosticPath
	}
	return text + "\n\n診斷紀錄未保存。"
}
func validateURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "http" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.RawPath != "" {
		return errors.New("invalid product URL")
	}
	if u.Path != "/studio/v2" && u.Path != "/studio/logs" {
		return errors.New("unsupported product route")
	}
	if u.Hostname() != "localhost" {
		address, err := netip.ParseAddr(u.Hostname())
		if err != nil {
			return errors.New("product URL requires a numeric host")
		}
		if strings.ContainsAny(address.Zone(), "%/?#[]@:\\\"<>|\x00\r\n\t") {
			return errors.New("invalid scoped IPv6 interface")
		}
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil || port < 1 || port > 65535 {
		return errors.New("invalid product port")
	}
	if strings.ContainsAny(raw, "\x00\r\n") {
		return errors.New("invalid product URL characters")
	}
	return nil
}

type menuItem struct {
	ID      uintptr
	Label   string
	Enabled bool
}

func menuItems(s State) []menuItem {
	process := cmp.Or(s.Process, "starting")
	if s.PendingPhase != "" {
		process += " (" + s.PendingPhase + ")"
	}
	return []menuItem{{1, "開啟設定頁", s.SetupURL != ""}, {2, "查看執行日誌", true}, {0, "程序：" + process, false}, {0, "採集：" + cmp.Or(s.Acquisition, "unknown"), false}, {3, "版本資訊", true}, {4, "結束…", s.Process != processStopping}}
}

// requestFaultStop records a safe fault synchronously before permitting shutdown.
// It never transmits raw native error data to the managed diagnostics callback.
func requestFaultStop(options Options) {
	if options.OnFault != nil {
		options.OnFault()
	}
	if options.OnQuit != nil {
		go options.OnQuit()
	}
}
