package main

import (
	"context"
	"errors"
	"time"
)

// Native ShellExecute is not cancellable. At most one host opener owns its own
// OS thread; waiting for it never delays service shutdown or admits another one.
func (h *gatewayHost) openBrowser(target string, opener func(string) error) error {
	if !h.browserOpening.CompareAndSwap(false, true) {
		return errors.New("browser open already pending")
	}
	ctx := h.startupCtx
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if ctx.Err() != nil {
		h.browserOpening.Store(false)
		return ctx.Err()
	}
	result := make(chan error, 1)
	go func() { defer h.browserOpening.Store(false); result <- opener(target) }()
	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
