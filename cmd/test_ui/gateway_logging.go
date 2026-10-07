package main

import (
	"context"
	"io"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"

	"go-gateway/internal/diagnostics"
)

type managedWriter struct {
	broker       *diagnostics.Broker
	code, source string
	mirror       io.Writer
}

func (w managedWriter) Write(raw []byte) (int, error) {
	w.broker.Emit(diagnostics.Input{Code: w.code, Fields: map[string]any{"source": w.source}})
	if w.mirror != nil {
		return w.mirror.Write(raw)
	}
	return len(raw), nil
}

type managedSlog struct {
	broker *diagnostics.Broker
	mirror slog.Handler
}

func (h managedSlog) Enabled(ctx context.Context, level slog.Level) bool {
	return h.mirror == nil || h.mirror.Enabled(ctx, level)
}
func (h managedSlog) Handle(ctx context.Context, record slog.Record) error {
	h.broker.Emit(diagnostics.Input{Code: "raw.suppressed", Fields: map[string]any{"source": "slog"}})
	if h.mirror != nil {
		return h.mirror.Handle(ctx, record)
	}
	return nil
}
func (h managedSlog) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h.mirror != nil {
		h.mirror = h.mirror.WithAttrs(attrs)
	}
	return h
}
func (h managedSlog) WithGroup(name string) slog.Handler {
	if h.mirror != nil {
		h.mirror = h.mirror.WithGroup(name)
	}
	return h
}

func (h *gatewayHost) installLogging(mirror bool) {
	oldWriter, oldSlog, oldFlags := log.Writer(), slog.Default(), log.Flags()
	oldGin, oldRecovery := gin.DefaultWriter, gin.DefaultErrorWriter
	var slogMirror slog.Handler
	if mirror {
		h.consoleWriter = oldWriter
		slogMirror = oldSlog.Handler()
	} else {
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
	}
	slog.SetDefault(slog.New(managedSlog{broker: h.logs, mirror: slogMirror}))
	// SetDefault rewires standard log. Install its separate safe adapter afterward.
	log.SetOutput(managedWriter{broker: h.logs, code: "raw.suppressed", source: "standard", mirror: h.consoleWriter})
	log.SetFlags(oldFlags)
	h.restoreLogging = func() {
		slog.SetDefault(oldSlog)
		log.SetOutput(oldWriter)
		log.SetFlags(oldFlags)
		gin.DefaultWriter = oldGin
		gin.DefaultErrorWriter = oldRecovery
	}
}

// mirrorConsole retains the selected output policy after global logger restore.
func (h *gatewayHost) mirrorConsole(format string, args ...any) {
	if h.consoleWriter != nil {
		log.New(h.consoleWriter, "", log.LstdFlags).Printf(format, args...)
	}
}
