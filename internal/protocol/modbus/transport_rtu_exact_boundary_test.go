package modbus

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
	"go.bug.st/serial"
)

type plannedSerialPort struct {
	serial.Port
	chunks         [][]byte
	timeout        time.Duration
	timeouts       []time.Duration
	gapDelay       time.Duration
	eof, earlyZero bool
	closed         chan struct{}
	once           sync.Once
	writes, closes atomic.Int32
	entered        chan struct{}
}

func (p *plannedSerialPort) SetReadTimeout(timeout time.Duration) error {
	p.timeout = timeout
	p.timeouts = append(p.timeouts, timeout)
	return nil
}
func (p *plannedSerialPort) Write(b []byte) (int, error) { p.writes.Add(1); return len(b), nil }
func (p *plannedSerialPort) Close() error {
	p.closes.Add(1)
	p.once.Do(func() { close(p.closed) })
	return nil
}
func (p *plannedSerialPort) Read(b []byte) (int, error) {
	if p.entered != nil {
		close(p.entered)
		p.entered = nil
	}
	if len(p.chunks) > 0 {
		if p.gapDelay > 0 {
			timer := time.NewTimer(p.gapDelay)
			defer timer.Stop()
			select {
			case <-timer.C:
			case <-p.closed:
				return 0, io.EOF
			}
		}
		n := copy(b, p.chunks[0])
		p.chunks[0] = p.chunks[0][n:]
		if len(p.chunks[0]) == 0 {
			p.chunks = p.chunks[1:]
		}
		return n, nil
	}
	if p.eof {
		return 0, io.EOF
	}
	if p.earlyZero {
		return 0, nil
	}
	timer := time.NewTimer(p.timeout)
	defer timer.Stop()
	select {
	case <-timer.C:
		return 0, nil
	case <-p.closed:
		return 0, io.EOF
	}
}
func plannedPhysicalRTU(chunks ...[]byte) (*ModbusClient, *RTUTransport, *plannedSerialPort) {
	port := &plannedSerialPort{chunks: chunks, closed: make(chan struct{})}
	transport := &RTUTransport{port: port, BaudRate: 115200, DataBits: 8, Timeout: time.Second}
	return NewClient(transport, 1), transport, port
}
func TestPlannedActualRTURejectsContinuousTail(t *testing.T) {
	for _, frame := range [][]byte{{1, 3, 2, 0x12, 0x34, 0xb5, 0x33, 0, 0}, {1, 0x83, 1, 0x80, 0xf0, 0, 0}} {
		for _, split := range []bool{false, true} {
			t.Run(string(rune(frame[1]))+map[bool]string{false: "single-read", true: "chunked-tail"}[split], func(t *testing.T) {
				chunks := [][]byte{frame}
				if split {
					chunks = [][]byte{frame[:len(frame)-2], frame[len(frame)-2 : len(frame)-1], frame[len(frame)-1:]}
				}
				client, _, port := plannedPhysicalRTU(chunks...)
				if split {
					port.gapDelay = 100 * time.Microsecond
				}
				ctx := wirebudget.With(t.Context(), 1, 0)
				words, err := client.WithContext(ctx).ReadWordsExact(3, 0, 1)
				require.Error(t, err)
				require.Nil(t, words)
				require.False(t, IsCommandRejected(err))
				require.Equal(t, 1, wirebudget.From(ctx).Used())
				require.EqualValues(t, 1, port.writes.Load())
			})
		}
	}
}
func TestPlannedActualRTUCompleteFrameBoundary(t *testing.T) {
	for _, tc := range []struct {
		name      string
		chunks    [][]byte
		exception bool
	}{
		{"normal", [][]byte{{1, 3, 2, 0x12, 0x34, 0xb5, 0x33}}, false},
		{"partial-reads", [][]byte{{1}, {3, 2}, {0x12}, {0x34, 0xb5}, {0x33}}, false},
		{"exception", [][]byte{{1, 0x83, 1, 0x80, 0xf0}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, _, port := plannedPhysicalRTU(tc.chunks...)
			ctx := wirebudget.With(t.Context(), 1, 0)
			words, err := client.WithContext(ctx).ReadWordsExact(3, 0, 1)
			if tc.exception {
				require.Error(t, err)
				require.True(t, IsCommandRejected(err))
				require.Nil(t, words)
			} else {
				require.NoError(t, err)
				require.Equal(t, []uint16{0x1234}, words)
			}
			require.Equal(t, 1, wirebudget.From(ctx).Used())
			require.GreaterOrEqual(t, len(port.timeouts), 3)
			require.Equal(t, port.timeouts[0], port.timeouts[len(port.timeouts)-1], "restore the captured request timeout")
		})
	}
}
func TestPlannedActualRTUUnknownBoundaryFailsClosed(t *testing.T) {
	for _, fault := range []string{"eof", "early-zero", "overflow", "deadline", "cancel"} {
		t.Run(fault, func(t *testing.T) {
			frame := []byte{1, 3, 2, 0x12, 0x34, 0xb5, 0x33}
			if fault == "overflow" {
				frame = make([]byte, 257)
			}
			client, tr, port := plannedPhysicalRTU(frame)
			port.eof = fault == "eof"
			port.earlyZero = fault == "early-zero"
			ctx := wirebudget.With(t.Context(), 1, 0)
			var cancel context.CancelFunc
			if fault == "deadline" {
				tr.BaudRate = 300
				ctx, cancel = context.WithTimeout(ctx, 20*time.Millisecond)
				defer cancel()
			}
			if fault == "cancel" {
				ctx, cancel = context.WithCancel(ctx)
				defer cancel()
				entered := make(chan struct{})
				port.entered = entered
				go func() { <-entered; cancel() }()
			}
			words, err := client.WithContext(ctx).ReadWordsExact(3, 0, 1)
			require.Error(t, err)
			require.Nil(t, words)
			require.False(t, IsCommandRejected(err))
			require.Equal(t, 1, wirebudget.From(ctx).Used())
			if fault == "deadline" {
				require.ErrorIs(t, err, context.DeadlineExceeded)
			}
			if fault == "cancel" {
				require.ErrorIs(t, err, context.Canceled)
			}
			if fault != "overflow" {
				require.GreaterOrEqual(t, len(port.timeouts), 2)
			}
		})
	}
}
