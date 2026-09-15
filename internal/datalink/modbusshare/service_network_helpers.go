package modbusshare

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

// ReadHoldingWords is a helper for tests and diagnostics.
func (s *Service) ReadHoldingWords(startRegister, quantity uint16) ([]uint16, error) {
	if quantity == 0 {
		return []uint16{}, nil
	}
	buf := make([]byte, quantity*2)
	for i := uint16(0); i < quantity; i++ {
		offset := int(startRegister+i) * 2
		w, err := s.bank.ReadWord(offset)
		if err != nil {
			return nil, err
		}
		binary.BigEndian.PutUint16(buf[i*2:], w)
	}
	out := make([]uint16, quantity)
	for i := range out {
		out[i] = binary.BigEndian.Uint16(buf[i*2:])
	}
	return out, nil
}

func preflightPortAvailable(ctx context.Context, bindAddress string, port int) error {
	listenCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var lc net.ListenConfig
	ln, err := lc.Listen(listenCtx, "tcp", net.JoinHostPort(bindAddress, fmt.Sprintf("%d", port)))
	if err != nil {
		if isAddressInUseError(err) {
			return fmt.Errorf("port %d is already in use, please stop the conflicting process or choose another port", port)
		}
		return fmt.Errorf("cannot bind port %d: %w", port, err)
	}
	if err := ln.Close(); err != nil {
		return err
	}
	// Keep the legacy conflict check for wildcard listeners too. A process
	// already bound on all interfaces must not be reported as a successful
	// Share bind merely because the configured loopback probe is narrower.
	wildcard, err := lc.Listen(listenCtx, "tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		if isAddressInUseError(err) {
			return fmt.Errorf("port %d is already in use, please stop the conflicting process or choose another port", port)
		}
		return fmt.Errorf("cannot bind port %d: %w", port, err)
	}
	return wildcard.Close()
}

func isAddressInUseError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "address already in use") ||
		strings.Contains(msg, "only one usage of each socket address")
}
