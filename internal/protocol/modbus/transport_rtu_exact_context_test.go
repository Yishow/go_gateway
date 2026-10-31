package modbus

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
)

type plainExactRTU struct{ wrappedRTUTransport }

func (p *plainExactRTU) SendReceiveExact(b []byte) ([]byte, error) { return p.SendReceive(b) }
func TestPlannedExactRTUWrapperMustPreserveContextAndBoundary(t *testing.T) {
	plain := &plainExactRTU{}
	ctx := wirebudget.With(t.Context(), 1, 0)
	words, err := NewClient(plain, 1).WithContext(ctx).ReadWordsExact(3, 0, 1)
	require.ErrorIs(t, err, ErrInvalidFrame)
	require.Nil(t, words)
	require.Empty(t, plain.sent)
	require.Zero(t, wirebudget.From(ctx).Used())
	wrapper := &wrappedRTUTransport{}
	words, err = NewClient(wrapper, 1).ReadWordsExact(3, 0, 1)
	require.ErrorIs(t, err, ErrInvalidFrame)
	require.Nil(t, words)
	require.Empty(t, wrapper.sent)
}
func TestPlannedExactRTUSilenceUsesOfficialBaudThreshold(t *testing.T) {
	for _, tc := range []struct {
		baud int
		gap  time.Duration
	}{{0, 4010417}, {9600, 4010417}, {19200, 2005209}, {19201, 1750000}, {115200, 1750000}} {
		tr := &RTUTransport{BaudRate: tc.baud}
		require.Equal(t, tc.gap, tr.exactFrameSilence())
	}
}
func TestPlannedExactRTUCancellationCannotRetireReplacement(t *testing.T) {
	client, tr, port := plannedPhysicalRTU([]byte{1, 3, 2, 0x12, 0x34, 0xb5, 0x33})
	port.entered = make(chan struct{})
	entered := port.entered
	base, cancel := context.WithCancel(t.Context())
	defer cancel()
	ctx := wirebudget.With(base, 2, 0)
	done := make(chan error, 1)
	go func() { _, err := client.WithContext(ctx).ReadWordsExact(3, 0, 1); done <- err }()
	<-entered
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)
	require.EqualValues(t, 1, port.closes.Load())
	require.Nil(t, tr.port)
	require.Equal(t, port.timeouts[0], port.timeouts[len(port.timeouts)-1])
	oldCalls := len(port.timeouts)
	_, _, replacement := plannedPhysicalRTU([]byte{1, 3, 2, 0x12, 0x34, 0xb5, 0x33})
	replacement.entered = make(chan struct{})
	nextEntered := replacement.entered
	tr.mu.Lock()
	tr.port = replacement
	tr.mu.Unlock()
	go func() { <-nextEntered; cancel() }()
	words, err := client.WithContext(t.Context()).ReadWordsExact(3, 0, 1)
	require.NoError(t, err)
	require.Equal(t, []uint16{0x1234}, words)
	require.Zero(t, replacement.closes.Load())
	require.Equal(t, oldCalls, len(port.timeouts))
	require.Equal(t, 1, wirebudget.From(ctx).Used())
	require.NoError(t, tr.Close())
}
func TestPlannedExactRTUBitFrameAndTotalDeadline(t *testing.T) {
	client, _, _ := plannedPhysicalRTU([]byte{1, 1, 2, 0xa5, 2, 0x43, 0x6d})
	bits, err := client.WithContext(t.Context()).ReadBitsExact(1, 0, 10)
	require.NoError(t, err)
	require.Equal(t, []bool{true, false, true, false, false, true, false, true, false, true}, bits)
	var chunks [][]byte
	for range 100 {
		chunks = append(chunks, []byte{1})
	}
	client, tr, port := plannedPhysicalRTU(chunks...)
	tr.BaudRate = 300
	port.gapDelay = time.Millisecond
	ctx, cancel := context.WithTimeout(wirebudget.With(t.Context(), 1, 0), 10*time.Millisecond)
	defer cancel()
	words, err := client.WithContext(ctx).ReadWordsExact(3, 0, 1)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, words)
	require.Equal(t, 1, wirebudget.From(ctx).Used())
	require.NotEmpty(t, port.chunks, "chunks cannot extend the total deadline")
}

type failedRestoreSerial struct {
	*plannedSerialPort
	calls int
}

func (p *failedRestoreSerial) SetReadTimeout(timeout time.Duration) error {
	p.calls++
	if p.calls == 4 {
		return errors.New("restore failed")
	}
	return p.plannedSerialPort.SetReadTimeout(timeout)
}
func TestPlannedExactRTURestoreFailureIsNotGood(t *testing.T) {
	_, tr, port := plannedPhysicalRTU([]byte{1, 3, 2, 0x12, 0x34, 0xb5, 0x33})
	failing := &failedRestoreSerial{plannedSerialPort: port}
	tr.port = failing
	ctx := wirebudget.With(t.Context(), 1, 0)
	words, err := NewClient(tr, 1).WithContext(ctx).ReadWordsExact(3, 0, 1)
	require.ErrorContains(t, err, "restore")
	require.Nil(t, words)
	require.False(t, IsCommandRejected(err))
	require.Equal(t, 1, wirebudget.From(ctx).Used())
	require.EqualValues(t, 1, port.closes.Load())
	require.Nil(t, tr.port)
}

// The supported Windows serial dependency floors SetReadTimeout to milliseconds
// (go.bug.st/serial v1.6.4 serial_windows.go). This model keeps that real contract.
type millisExactSerial struct {
	*plannedSerialPort
	requested []time.Duration
}

func (p *millisExactSerial) SetReadTimeout(timeout time.Duration) error {
	p.requested = append(p.requested, timeout)
	return p.plannedSerialPort.SetReadTimeout(time.Duration(timeout.Milliseconds()) * time.Millisecond)
}
func (p *millisExactSerial) Read(b []byte) (int, error) {
	if len(p.chunks) > 0 {
		return p.plannedSerialPort.Read(b)
	}
	// A short deterministic driver wait avoids host timer rounding concealing the
	// 1.750ms -> 1ms regression. This is test-only, not a production polling loop.
	until := time.Now().Add(p.timeout)
	for time.Now().Before(until) {
	}
	return 0, nil
}
func TestPlannedExactRTUMillisecondDriverStillProvesSilence(t *testing.T) {
	for _, baud := range []int{9600, 19200, 115200} {
		t.Run(fmt.Sprint(baud), func(t *testing.T) {
			_, tr, port := plannedPhysicalRTU([]byte{1, 3, 2, 0x12, 0x34, 0xb5, 0x33})
			tr.BaudRate = baud
			quantized := &millisExactSerial{plannedSerialPort: port}
			tr.port = quantized
			ctx := wirebudget.With(t.Context(), 1, 0)
			words, err := NewClient(tr, 1).WithContext(ctx).ReadWordsExact(3, 0, 1)
			require.NoError(t, err)
			require.Equal(t, []uint16{0x1234}, words)
			require.Equal(t, 1, wirebudget.From(ctx).Used())
			silence := quantized.requested[2]
			require.Zero(t, silence%time.Millisecond)
			require.GreaterOrEqual(t, silence, tr.exactFrameSilence())
			require.Equal(t, quantized.requested[0], quantized.requested[len(quantized.requested)-1])
		})
	}
}
