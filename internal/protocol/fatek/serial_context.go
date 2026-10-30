package fatek

import (
	"bufio"
	"context"
	"errors"
	"go-gateway/internal/protocol/wirebudget"
	"go.bug.st/serial"
	"io"
	"time"
)

// ConnectContext checks cancellation around the OS serial-open operation. The OS
// open itself has no context API; no late handle is published after cancellation.
func (s *SerialTransport) ConnectContext(ctx context.Context) error {
	if err := s.mu.LockContext(ctx); err != nil {
		return err
	}
	defer s.mu.Unlock()
	if err := s.closeLocked(); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	port, err := serial.Open(s.Port, &serial.Mode{BaudRate: s.BaudRate, DataBits: s.DataBits, Parity: s.Parity, StopBits: s.StopBits})
	if err != nil {
		return err
	}
	s.port = port
	s.reader = bufio.NewReader(port)
	if err = ctx.Err(); err != nil {
		return errors.Join(err, s.closeLocked())
	}
	return nil
}
func (s *SerialTransport) closeLocked() error {
	if s.port == nil {
		return nil
	}
	if err := s.port.Close(); err != nil {
		s.closeErr = err
		return err
	}
	s.port = nil
	s.reader = nil
	s.closeErr = nil
	return nil
}
func (s *SerialTransport) CloseContext(ctx context.Context) error {
	if err := s.mu.LockContext(ctx); err != nil {
		return err
	}
	defer s.mu.Unlock()
	return s.closeLocked()
}
func (s *SerialTransport) IsConnectedContext(ctx context.Context) (bool, error) {
	if err := s.mu.LockContext(ctx); err != nil {
		return false, err
	}
	defer s.mu.Unlock()
	return s.port != nil && s.closeErr == nil, nil
}

// serialContextReader treats a timeout's zero,nil as termination, not an endless
// ReadFull retry. A cancellation callback closes this exact port to unblock I/O.
type serialContextReader struct {
	ctx  context.Context
	port serial.Port
}

func (r serialContextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	n, err := r.port.Read(p)
	if err == nil && n == 0 {
		return 0, io.ErrNoProgress
	}
	return n, err
}
func (s *SerialTransport) SendReceiveContext(ctx context.Context, data []byte) (response []byte, err error) {
	if err = s.mu.LockContext(ctx); err != nil {
		return nil, err
	}
	defer s.mu.Unlock()
	ctx, cancel := boundedContext(ctx, s.Timeout)
	defer cancel()
	if s.port == nil {
		return nil, ErrConnectionClosed
	}
	if s.closeErr != nil {
		return nil, s.closeErr
	}
	port := s.port
	deadline, _ := ctx.Deadline()
	if err = port.SetReadTimeout(time.Until(deadline)); err != nil {
		return nil, errors.Join(err, s.closeLocked())
	}
	done := make(chan struct{})
	var closeErr error
	stop := context.AfterFunc(ctx, func() { defer close(done); closeErr = port.Close() })
	defer func() {
		if !stop() {
			<-done
			s.closeErr = closeErr
			if closeErr == nil {
				s.port = nil
				s.reader = nil
			}
			err = errors.Join(ctx.Err(), closeErr)
			return
		}
		if err != nil {
			err = errors.Join(err, s.closeLocked())
		}
	}()
	if err = wirebudget.Take(ctx); err != nil {
		return nil, err
	}
	n, err := port.Write(data)
	if err != nil {
		return nil, err
	}
	if n != len(data) {
		return nil, io.ErrShortWrite
	}
	// Each request owns its reader; no stale bytes cross a canceled generation.
	s.reader = bufio.NewReader(serialContextReader{ctx, port})
	return s.reader.ReadBytes(ETX)
}
