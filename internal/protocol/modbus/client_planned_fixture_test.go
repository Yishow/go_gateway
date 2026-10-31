package modbus

import (
	"context"
	"encoding/binary"
	"fmt"
	"go-gateway/internal/protocol/wirebudget"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type exactReader interface {
	ReadWordsExact(byte, uint16, uint16) ([]uint16, error)
	ReadBitsExact(byte, uint16, uint16) ([]bool, error)
}

func plannedExactReader(t *testing.T, c *ModbusClient) exactReader {
	t.Helper()
	r, ok := any(c).(exactReader)
	require.True(t, ok, "Modbus client must expose exact raw planned reads")
	return r
}

type plannedCell struct {
	function byte
	address  uint16
}
type plannedPeer struct {
	mu     sync.Mutex
	rtu    bool
	ids    uint16
	frames [][]byte
	words  map[plannedCell]uint16
	bits   map[plannedCell]bool
	mutate func([]byte) []byte
	fault  error
}

func (*plannedPeer) Connect() error { return nil }
func (*plannedPeer) Close() error   { return nil }
func (p *plannedPeer) GetOriginalTransport() any {
	if p.rtu {
		return &RTUTransport{}
	}
	return &TCPTransport{}
}
func (p *plannedPeer) GetNextTransactionID() uint16 {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ids++
	return p.ids
}
func (p *plannedPeer) SendReceive(req []byte) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.frames = append(p.frames, append([]byte(nil), req...))
	if p.fault != nil {
		return nil, p.fault
	}
	off := 7
	if p.rtu {
		off = 1
		if len(req) != 8 || CalculateCRC16(req[:6]) != binary.LittleEndian.Uint16(req[6:]) {
			return nil, fmt.Errorf("invalid RTU request")
		}
	} else if len(req) != 12 || binary.BigEndian.Uint16(req[2:4]) != 0 || binary.BigEndian.Uint16(req[4:6]) != 6 {
		return nil, fmt.Errorf("invalid MBAP request")
	}
	fn := req[off]
	start, count := binary.BigEndian.Uint16(req[off+1:]), binary.BigEndian.Uint16(req[off+3:])
	var payload []byte
	switch fn {
	case 1, 2:
		payload = make([]byte, 1+(int(count)+7)/8)
		payload[0] = byte(len(payload) - 1)
		for i := uint16(0); i < count; i++ {
			v, ok := p.bits[plannedCell{fn, start + i}]
			if !ok {
				return nil, fmt.Errorf("unmapped bit %d:%d", fn, start+i)
			}
			if v {
				payload[1+i/8] |= 1 << (i % 8)
			}
		}
	case 3, 4:
		payload = make([]byte, 1+int(count)*2)
		payload[0] = byte(len(payload) - 1)
		for i := uint16(0); i < count; i++ {
			v, ok := p.words[plannedCell{fn, start + i}]
			if !ok {
				return nil, fmt.Errorf("unmapped word %d:%d", fn, start+i)
			}
			binary.BigEndian.PutUint16(payload[1+int(i)*2:], v)
		}
	default:
		return nil, fmt.Errorf("unexpected function %d", fn)
	}
	response := append(append([]byte(nil), req[:off+1]...), payload...)
	if p.rtu {
		response = append(response, 0, 0)
		binary.LittleEndian.PutUint16(response[len(response)-2:], CalculateCRC16(response[:len(response)-2]))
	} else {
		binary.BigEndian.PutUint16(response[4:6], uint16(len(response)-6))
	}
	if p.mutate != nil {
		response = p.mutate(response)
	}
	return response, nil
}
func plannedResponsePayload(r []byte, rtu bool, pdu []byte) []byte {
	off := 7
	if rtu {
		off = 1
	}
	out := append(append([]byte(nil), r[:off]...), pdu...)
	if rtu {
		out = append(out, 0, 0)
		binary.LittleEndian.PutUint16(out[len(out)-2:], CalculateCRC16(out[:len(out)-2]))
	} else {
		binary.BigEndian.PutUint16(out[4:6], uint16(len(out)-6))
	}
	return out
}
func plannedBasicPeer(rtu bool) *plannedPeer {
	p := &plannedPeer{rtu: rtu, words: map[plannedCell]uint16{}, bits: map[plannedCell]bool{}}
	for _, fn := range []byte{3, 4} {
		p.words[plannedCell{fn, 0}] = 0x0102
		p.words[plannedCell{fn, 1}] = 0x0304
	}
	for _, fn := range []byte{1, 2} {
		for i := uint16(0); i < 9; i++ {
			p.bits[plannedCell{fn, i}] = i == 0 || i == 3 || i == 8
		}
	}
	return p
}

// This fake peer returns one complete parsed response frame, never a serial prefix.
func (p *plannedPeer) SendReceiveExact(request []byte) ([]byte, error) { return p.SendReceive(request) }

func (p *plannedPeer) SendReceiveExactContext(ctx context.Context, request []byte) ([]byte, error) {
	if err := wirebudget.Take(ctx); err != nil {
		return nil, err
	}
	return p.SendReceive(request)
}
