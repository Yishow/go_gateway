package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"go-gateway/internal/apphost"
	"go-gateway/internal/desktop"
)

const hostDegraded = "degraded"
const goosWindows = "windows"

func (h *gatewayHost) serve(handler http.Handler) error {
	select {
	case <-h.quit:
		return nil
	default:
	}
	baseCtx, cancel := context.WithCancel(context.Background())
	h.httpCancel = cancel
	h.requests = &requestGate{next: handler, stopping: h.quit}
	h.server = &http.Server{Addr: h.listener.Addr().String(), Handler: h.requests, IdleTimeout: 120 * time.Second, ReadHeaderTimeout: 5 * time.Second,
		BaseContext: func(net.Listener) context.Context { return baseCtx }, ErrorLog: log.New(managedWriter{broker: h.logs, code: "http.server_error", mirror: h.consoleWriter}, "", 0)}
	// All resource fields are now published. Shutdown must not wait on browser UI.
	h.startupOnce.Do(func() { close(h.startupDone) })
	serving := make(chan error, 1)
	go func() { serving <- h.server.Serve(h.listener) }()
	expected, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		return apphost.NewFault("startup.assets_missing", err)
	}
	if err = probeOwnHTTP(h.startupCtx, h.listener.Addr().String(), expected); err != nil {
		return apphost.NewFault("startup.failed", err)
	}
	select {
	case <-h.quit:
		return nil
	default:
	}
	h.emit("startup.ready", nil)
	if origin, _, err := apphost.ListenerURL(h.listener.Addr().String()); err == nil {
		h.mirrorConsole("Gateway ready at %s/studio/v2", origin)
	}
	h.refreshServingState()
	select {
	case <-h.quit:
		return nil
	default:
	}
	if h.launch.AutoOpen(os.Getenv("AUTO_OPEN_BROWSER")) {
		if err := h.openSetup(); err != nil {
			h.emit("http.server_error", nil)
		}
	}
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		select {
		case <-h.quit:
			return nil
		case <-tick.C:
			h.refreshServingState()
		case err := <-serving:
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return apphost.NewFault("startup.bind_failed", err)
		}
	}
}
func (h *gatewayHost) refreshServingState() {
	state := "serving"
	if h.runtime == nil || !h.runtime.IsRunning() || (h.sink != nil && h.sink.Health().Degraded) {
		state = hostDegraded
	}
	h.updateState(state)
}
func (h *gatewayHost) openSetup() error {
	h.stateMu.Lock()
	if h.isStopping() {
		h.stateMu.Unlock()
		return errors.New("gateway is stopping")
	}
	target := h.state.SetupURL
	h.stateMu.Unlock()
	if target == "" {
		return errors.New("setup is not ready")
	}
	return h.openBrowser(target, func(target string) error {
		return openProductURL(target, func(err error) {
			h.emit("http.server_error", nil)
			h.mirrorConsole("browser opener returned an error: %v", err)
		})
	})
}
func openProductURL(target string, onError func(error)) error {
	if runtime.GOOS == goosWindows {
		return desktop.OpenURL(target)
	}
	executable := "xdg-open"
	if runtime.GOOS == "darwin" {
		executable = "open"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cmd := exec.CommandContext(ctx, executable, target)
	if err := cmd.Start(); err != nil {
		cancel()
		return err
	}
	go func() {
		defer cancel()
		if err := cmd.Wait(); err != nil {
			if onError != nil {
				onError(err)
			}
		}
	}()
	return nil
}

type requestGate struct {
	mu       sync.Mutex
	closed   bool
	stopping <-chan struct{}
	workers  sync.WaitGroup
	next     http.Handler
}

func (g *requestGate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	select {
	case <-g.stopping:
		g.closed = true
	default:
	}
	if g.closed {
		g.mu.Unlock()
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, err := w.Write([]byte(`{"code":"gateway_stopping","message":"Gateway is stopping."}`)); err != nil {
			return
		}
		return
	}
	g.workers.Add(1)
	g.mu.Unlock()
	defer g.workers.Done()
	g.next.ServeHTTP(w, r)
}
func (g *requestGate) stop() { g.mu.Lock(); g.closed = true; g.mu.Unlock() }
func (g *requestGate) wait() { g.workers.Wait() }
func (h *gatewayHost) stopHTTP(ctx context.Context) error {
	if h.requests != nil {
		h.requests.stop()
	}
	if h.logs != nil {
		h.logs.CloseSubscriptions()
	}
	if h.httpCancel != nil {
		h.httpCancel()
	}
	if h.server != nil {
		return h.server.Shutdown(ctx)
	}
	if h.listener != nil {
		return h.listener.Close()
	}
	return nil
}

type acquisitionStarter interface {
	Start(context.Context) error
	IsRunning() bool
}

func (h *gatewayHost) startAcquisition(ctx context.Context, starter acquisitionStarter) {
	if err := starter.Start(ctx); err != nil {
		log.Printf("datalink runtime startup failed: %v", err)
		h.emit("runtime.degraded", nil)
		return
	}
	if starter.IsRunning() {
		h.emit("runtime.started", nil)
	} else {
		h.emit("runtime.degraded", nil)
	}
}
