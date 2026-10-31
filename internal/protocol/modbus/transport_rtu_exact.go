package modbus

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"go.bug.st/serial"
)

// SendReceiveExact uses the existing RTU transaction gate and actual-send budget,
// but establishes the complete frame boundary instead of trusting a byte count.
func (r *RTUTransport) SendReceiveExact(data []byte) ([]byte, error) {
	return r.SendReceiveExactContext(context.Background(), data)
}

// SendReceiveExactContext preserves one total deadline, including final silence.
// Legacy scalar/write calls continue using SendReceiveContext unchanged.
func (r *RTUTransport) SendReceiveExactContext(ctx context.Context, data []byte) ([]byte, error) {
	return r.sendReceiveContext(ctx, data, true)
}

// The serial-line specification uses 11-bit RTU characters through 19200 baud,
// and recommends fixed 1.750 ms t3.5 above it. Waiting the full protocol interval
// is conservative for existing nonstandard 8N1 configurations; none are changed.
func (r *RTUTransport) exactFrameSilence() time.Duration {
	baud := r.BaudRate
	if baud <= 0 {
		baud = DefaultBaudRate
	}
	if baud > 19200 {
		return 1750 * time.Microsecond
	}
	return time.Duration((385*int64(time.Second) + 10*int64(baud) - 1) / (10 * int64(baud)))
}

func (r *RTUTransport) readExactFrame(ctx context.Context, port serial.Port, restoreTimeout time.Duration) (response []byte, resultErr error) {
	// Only this captured port is touched. The existing operation gate and joined
	// serialContext callback prevent a replacement from inheriting this cleanup.
	defer func() {
		if err := port.SetReadTimeout(restoreTimeout); err != nil {
			response = nil
			resultErr = errors.Join(resultErr, fmt.Errorf("RTU restore read timeout: %w", err))
		}
	}()
	deadline, _ := ctx.Deadline()
	gap := r.exactFrameSilence()
	// The supported Windows serial driver floors durations to milliseconds;
	// round up so its timeout still establishes at least the protocol interval.
	// The parent deadline remains the upper bound on every physical Read.
	silenceTimeout := (gap + time.Millisecond - 1) / time.Millisecond * time.Millisecond
	buffer := make([]byte, 257) // 256-byte ADU plus a sentinel, never a valid prefix.
	used := 0
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, context.DeadlineExceeded
		}
		timeout := remaining
		if used > 0 && silenceTimeout < timeout {
			timeout = silenceTimeout
		}
		if err := port.SetReadTimeout(timeout); err != nil {
			return nil, fmt.Errorf("RTU set frame timeout: %w", err)
		}
		started := time.Now()
		n, err := port.Read(buffer[used:])
		if contextErr := transportContextError(ctx, nil); contextErr != nil {
			return nil, contextErr
		}
		if n < 0 || n > len(buffer)-used {
			return nil, ErrInvalidFrame
		}
		used += n
		if used > 256 {
			return nil, ErrInvalidFrame
		}
		if err != nil {
			// EOF is not a frame boundary. Once any bytes arrived, an unknown boundary
			// is malformed, not an unsupported exception or a retryable partial value.
			if used > 0 {
				return nil, fmt.Errorf("%w: RTU frame boundary: %v", ErrInvalidFrame, err)
			}
			return nil, err
		}
		if n > 0 {
			continue
		}
		if used == 0 {
			return nil, io.ErrNoProgress
		}
		// A driver returning early zero bytes did not prove t3.5. Do not spin or guess
		// that EOF/zero progress means a complete response.
		if timeout < gap || time.Since(started) < gap {
			return nil, ErrInvalidFrame
		}
		return buffer[:used], nil
	}
}
