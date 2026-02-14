package modbusshare

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// ReadHoldingWords is a helper for tests and diagnostics.
func (s *Service) ReadHoldingWords(startRegister uint16, quantity uint16) ([]uint16, error) {
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

func preflightPortAvailable(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		if isAddressInUseError(err) {
			return fmt.Errorf("port %d is already in use, please stop the conflicting process or choose another port", port)
		}
		return fmt.Errorf("cannot bind port %d: %w", port, err)
	}
	return ln.Close()
}

func isAddressInUseError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "address already in use") ||
		strings.Contains(msg, "only one usage of each socket address")
}
