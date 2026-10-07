package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-gateway/internal/diagnostics"

	"github.com/gin-gonic/gin"
)

const (
	runtimeLogWriteTimeout = 5 * time.Second
	runtimeLogReset        = "reset"
)

func (h *runtimeLogHandler) stream(c *gin.Context) {
	query, after, err := parseRuntimeLogQuery(c.Request, true)
	if err != nil {
		runtimeLogError(c, err)
		return
	}
	if h.broker == nil {
		runtimeLogError(c, diagnostics.ErrClosed)
		return
	}
	sub, err := h.broker.Subscribe(query, after)
	if err != nil {
		runtimeLogError(c, err)
		return
	}
	defer sub.Close()
	// Gin's Flush cannot return transport errors. Unwrap its writer so the
	// controller reaches net/http's FlushError while writes retain Gin bookkeeping.
	var raw http.ResponseWriter = c.Writer
	if wrapped, ok := raw.(interface{ Unwrap() http.ResponseWriter }); ok {
		raw = wrapped.Unwrap()
	}
	controller := http.NewResponseController(raw)
	if err := controller.SetWriteDeadline(time.Time{}); err != nil {
		runtimeLogError(c, diagnostics.ErrClosed)
		return
	}
	c.Header("Content-Type", "text/event-stream")
	c.Header("X-Accel-Buffering", "no")
	writer := runtimeLogStreamWriter{writer: c.Writer, controller: controller, gaps: make(map[string]logGapState, 6)}
	if err := writer.metadata("handshake", sub.Metadata()); err != nil {
		return
	}
	nextHeartbeat := time.Now().Add(h.heartbeat)
	for {
		ctx, cancel := context.WithDeadline(c.Request.Context(), nextHeartbeat)
		event, err := sub.Next(ctx)
		cancel()
		if c.Request.Context().Err() != nil {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			if err := writer.metadata("heartbeat", sub.Metadata()); err != nil {
				return
			}
			nextHeartbeat = time.Now().Add(h.heartbeat)
			continue
		}
		if err != nil {
			return
		} // Overflow/cancellation remains visible on reconnect.
		if err := writer.losses(sub.Metadata()); err != nil {
			return
		}
		if err := writeRuntimeLogFrame(c.Writer, controller, "log", event.Cursor(), event); err != nil {
			return
		}
		sub.Ack(event.Cursor())
		if !time.Now().Before(nextHeartbeat) {
			if err := writer.metadata("heartbeat", sub.Metadata()); err != nil {
				return
			}
			nextHeartbeat = time.Now().Add(h.heartbeat)
		}
	}
}

type logGapState struct {
	gap                diagnostics.Gap
	counter, recovered uint64
}
type runtimeLogStreamWriter struct {
	writer     http.ResponseWriter
	controller *http.ResponseController
	gaps       map[string]logGapState
}

func (w *runtimeLogStreamWriter) losses(meta diagnostics.Metadata) error {
	for _, gap := range meta.Gaps {
		state := logGapState{gap: gap}
		switch gap.Reason {
		case "admission":
			state.counter = meta.CaptureDropped
		case "subscriber_overflow":
			state.counter = meta.SubscriberOverflow
		case "file_loss":
			state.counter = meta.FileDropped
			state.recovered = meta.FileRecoveredGaps
		case "retention", "replay_limit", runtimeLogReset:
		default:
			continue
		}
		if previous, ok := w.gaps[gap.Reason]; ok && previous == state {
			continue
		}
		name := "gap"
		if gap.Reason == runtimeLogReset {
			name = runtimeLogReset
		}
		if err := writeRuntimeLogFrame(w.writer, w.controller, name, "", gap); err != nil {
			return err
		}
		w.gaps[gap.Reason] = state
	}
	return nil
}
func (w *runtimeLogStreamWriter) metadata(name string, meta diagnostics.Metadata) error {
	if err := w.losses(meta); err != nil {
		return err
	}
	return writeRuntimeLogFrame(w.writer, w.controller, name, meta.Progress, meta)
}

func writeRuntimeLogFrame(w http.ResponseWriter, controller *http.ResponseController, name, id string, value any) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := controller.SetWriteDeadline(time.Now().Add(runtimeLogWriteTimeout)); err != nil {
		return err
	}
	// Idle heartbeats are 15 seconds apart: never leave a five-second write
	// deadline armed between successful frames.
	defer func() { err = errors.Join(err, controller.SetWriteDeadline(time.Time{})) }()
	frame := fmt.Appendf(nil, "event: %s\n", name)
	if id != "" {
		frame = fmt.Appendf(frame, "id: %s\n", id)
	}
	frame = append(frame, "data: "...)
	frame = append(frame, data...)
	frame = append(frame, '\n', '\n')
	if n, err := w.Write(frame); err != nil {
		return err
	} else if n != len(frame) {
		return io.ErrShortWrite
	}
	return controller.Flush()
}
