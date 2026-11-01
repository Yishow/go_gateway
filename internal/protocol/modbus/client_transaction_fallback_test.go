package modbus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Both concrete legacy transports already match the preceding interface case.
var (
	_ interface{ GetNextTransactionID() uint16 } = (*TCPTransport)(nil)
	_ interface{ GetNextTransactionID() uint16 } = (*UDPTransport)(nil)
)

type unsupportedTransactionPeer struct {
	Transport
	sends int
}

func (p *unsupportedTransactionPeer) SendReceive([]byte) ([]byte, error) {
	p.sends++
	return nil, ErrInvalidFrame
}
func TestTransactionIDUnsupportedTransportFailsBeforeSend(t *testing.T) {
	peer := &unsupportedTransactionPeer{}
	body, err := NewClient(peer, 1).sendTCPRequest(FuncReadHoldingRegisters, BuildReadRequest(0, 1))
	require.Nil(t, body)
	require.EqualError(t, err, "不支援的傳輸類型")
	require.Zero(t, peer.sends)
}
