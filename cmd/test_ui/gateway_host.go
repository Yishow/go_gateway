package main

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"go-gateway/internal/apphost"
	"go-gateway/internal/config"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/modbusshare"
	datalinkruntime "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/desktop"
	"go-gateway/internal/diagnostics"
)

type gatewayHost struct {
	launch                    apphost.Launch
	cfg                       *config.Config
	data                      apphost.DataSelection
	owner                     *desktop.Owner
	shell                     desktop.Shell
	newShell                  func(desktop.Options) desktop.Shell
	trayReady                 bool
	sink                      *diagnostics.FileSink
	logs                      *diagnostics.Broker
	logMu                     sync.RWMutex
	listener                  net.Listener
	server                    *http.Server
	requests                  *requestGate
	httpCancel                context.CancelFunc
	db                        *sql.DB
	connections               *connector.ConnectionManager
	runtime                   *datalinkruntime.Service
	pipeline                  *grouppipeline.Pipeline
	share                     *modbusshare.Service
	fixtureCleanup            func() error
	restoreLogging            func()
	consoleWriter             io.Writer
	quit                      chan struct{}
	quitOnce                  sync.Once
	stateMu                   sync.Mutex
	state                     desktop.State
	shellErr                  error
	browserOpening            atomic.Bool
	shutdown                  *apphost.Shutdown
	shutdownDone, startupDone chan struct{}
	startupOnce               sync.Once
	startupCtx                context.Context
	startupCancel             context.CancelFunc
}

func (h *gatewayHost) requestQuit() {
	h.quitOnce.Do(func() {
		if h.startupCancel != nil {
			h.startupCancel()
		}
		close(h.quit)
		if h.shutdown != nil {
			h.shutdown.Begin()
		} else {
			h.updateShutdownState("stopping", "startup")
		}
	})
}
func (h *gatewayHost) emit(code string, fields map[string]any) {
	h.logMu.RLock()
	defer h.logMu.RUnlock()
	if h.logs != nil {
		h.logs.Emit(diagnostics.Input{Code: code, Fields: fields})
	}
}
func (h *gatewayHost) updateState(process string) {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	if h.isStopping() {
		return
	}
	h.state.Process = process
	h.state.PendingPhase = ""
	if h.runtime != nil && h.runtime.IsRunning() {
		h.state.Acquisition = "running"
	} else {
		h.state.Acquisition = "unavailable"
	}
	if h.listener != nil {
		base, local, err := apphost.ListenerURL(h.listener.Addr().String())
		if err == nil {
			h.state.SetupURL = base + "/studio/v2"
			if local {
				h.state.LogsURL = base + "/studio/logs"
			} else {
				h.state.LogsUnavailable = "此listener不能使用日誌的loopback Host白名單。請使用127.0.0.1、::1或wildcard bind。"
			}
		}
	}
	if h.shell != nil {
		h.shell.Update(h.state)
	}
}
func (h *gatewayHost) updateShutdownState(process, phase string) {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	h.state.Process = process
	h.state.PendingPhase = phase
	h.state.SetupURL = ""
	h.state.LogsURL = ""
	h.state.LogsUnavailable = "Gateway 正在關閉。"
	if phase == "runtime" {
		h.state.Acquisition = "stopping"
	}
	if phase == "pipeline" || phase == "share" || phase == "connections" || phase == "database" || phase == "diagnostics" {
		h.state.Acquisition = "stopped"
	}
	if h.shell != nil {
		h.shell.Update(h.state)
	}
}
func finishGatewayResult(runErr, closeErr error) error {
	if errors.Is(runErr, errLaunchHandled) {
		runErr = nil
	}
	return errors.Join(runErr, closeErr)
}

// isStopping is safe before the asynchronous coordinator starts its first phase.
func (h *gatewayHost) isStopping() bool {
	select {
	case <-h.quit:
		return true
	default:
	}
	return h.shutdown != nil && h.shutdown.IsStopping()
}

// recordShellFailure accepts only a fixed native fault classification.
func (h *gatewayHost) recordShellFailure() {
	h.stateMu.Lock()
	h.shellErr = apphost.NewFault("startup.tray_failed", nil)
	h.stateMu.Unlock()
	h.emit("startup.tray_failed", nil)
}
func (h *gatewayHost) shellFailure() error {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	return h.shellErr
}
