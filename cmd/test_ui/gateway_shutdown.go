package main

import (
	"context"
	"errors"
	"time"

	"go-gateway/internal/apphost"
)

func (h *gatewayHost) prepareShutdown() {
	phases := []apphost.Phase{
		{Name: "startup", Stop: func(ctx context.Context) error {
			select {
			case <-h.startupDone:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}, Wait: func(context.Context) error { <-h.startupDone; return nil }},
		{Name: "http", Stop: h.stopHTTP, Wait: func(context.Context) error {
			if h.requests != nil {
				h.requests.wait()
			}
			return nil
		}},
		{Name: "runtime", Stop: func(ctx context.Context) error {
			if h.runtime != nil {
				return h.runtime.Stop(ctx)
			}
			return nil
		}, Wait: func(ctx context.Context) error {
			if h.runtime != nil {
				return h.runtime.WaitStopped(ctx)
			}
			return nil
		}},
		{Name: "pipeline", Stop: func(ctx context.Context) error {
			if h.pipeline != nil {
				return h.pipeline.StopContext(ctx)
			}
			return nil
		}, Wait: func(ctx context.Context) error {
			if h.pipeline != nil {
				return h.pipeline.WaitStopped(ctx)
			}
			return nil
		}},
		{Name: "share", Stop: func(context.Context) error {
			if h.share != nil {
				return h.share.CloseRuntime()
			}
			return nil
		}},
		{Name: "connections", Stop: func(context.Context) error {
			var errs []error
			if h.fixtureCleanup != nil {
				errs = append(errs, h.fixtureCleanup())
			}
			if h.connections != nil {
				errs = append(errs, h.connections.CloseAll())
			}
			return errors.Join(errs...)
		}},
		{Name: "database", Stop: func(context.Context) error {
			if h.db != nil {
				return h.db.Close()
			}
			return nil
		}},
		{Name: "diagnostics", Stop: func(ctx context.Context) error {
			h.emit("shutdown.finalizing", nil)
			if h.restoreLogging != nil {
				h.restoreLogging()
			}
			return h.closeDiagnostics(ctx)
		}, Wait: func(ctx context.Context) error {
			if !h.mustObserveDiagnosticCompletion() {
				return nil
			}
			return h.closeDiagnostics(ctx)
		}},
		{Name: "tray", Stop: func(context.Context) error {
			close(h.shutdownDone)
			if h.shell != nil {
				return h.shell.Close()
			}
			return nil
		}},
		{Name: "owner", Stop: func(context.Context) error {
			if h.owner != nil {
				return h.owner.Close()
			}
			return nil
		}},
	}
	h.shutdown = apphost.NewShutdown(15*time.Second, phases, func(state apphost.ShutdownState) {
		h.updateShutdownState(state.Process, state.Phase)
		if state.Process == "stop-timeout" {
			h.emit("shutdown.timeout", map[string]any{"phase": state.Phase})
		}
		if state.Failed {
			h.emit("shutdown.failed", map[string]any{"phase": state.Phase})
		}
	})
}
func (h *gatewayHost) close() error {
	h.startupOnce.Do(func() { close(h.startupDone) })
	h.emit("shutdown.begin", nil)
	h.shutdown.Begin()
	return h.shutdown.Wait(context.Background())
}
func (h *gatewayHost) closeDiagnostics(ctx context.Context) error {
	if h.logs != nil {
		return h.logs.Close(ctx)
	}
	if h.sink != nil {
		return h.sink.Close(ctx)
	}
	return nil
}

// Only a pre-service fatal without a usable tray may leave a timed-out log writer
// to process termination. Its handles remain worker-owned; no live DB is closed.
func (h *gatewayHost) mustObserveDiagnosticCompletion() bool {
	return h.trayReady || h.db != nil || h.runtime != nil || h.pipeline != nil || h.share != nil || h.connections != nil || h.fixtureCleanup != nil
}
