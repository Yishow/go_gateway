package modbus

import (
	"context"
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
	"testing"
)

type contextClientTransport struct {
	Transport
	ids   uint16
	sends int
}

func (t *contextClientTransport) GetNextTransactionID() uint16 { t.ids++; return t.ids }
func (t *contextClientTransport) SendReceive(req []byte) ([]byte, error) {
	t.sends++
	return BuildTCPFrame(binary.BigEndian.Uint16(req), req[6], req[7], []byte{2, 0, 9}), nil
}
func TestClientContextPreservesTransactionAuthority(t *testing.T) {
	tr := &contextClientTransport{}
	c := NewClient(tr, 1)
	for i := 0; i < 2; i++ {
		values, err := c.WithContext(wirebudget.With(t.Context(), 1, 0)).ReadHoldingRegisters(0, 1)
		require.NoError(t, err)
		require.Equal(t, []uint16{9}, values)
	}
	require.Equal(t, uint16(2), tr.ids)
	require.Equal(t, 2, tr.sends)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := c.WithContext(ctx).ReadHoldingRegisters(0, 1)
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, 2, tr.sends)
}

func TestClientContextPreservesWrappedRTUFraming(t *testing.T) {
	tr := &wrappedRTUTransport{}
	c := NewClient(tr, 1)
	values, err := c.WithContext(wirebudget.With(t.Context(), 1, 0)).ReadHoldingRegisters(0, 1)
	require.NoError(t, err)
	require.Equal(t, []uint16{12345}, values)
	require.Equal(t, BuildRTUFrame(1, FuncReadHoldingRegisters, BuildReadRequest(0, 1)), tr.sent[0])
}

func TestCommandRejectionIsNotTransportFailure(t *testing.T) {
	require.True(t, IsCommandRejected(NewProtocolError(ExceptionIllegalFunction, "rejected")))
	for _, err := range []error{ErrCRCError, ErrTimeout, ErrResponseTooShort, ErrInvalidFrame, NewProtocolError(ExceptionIllegalDataAddress, "bad address")} {
		require.False(t, IsCommandRejected(err))
	}
}

func TestClientContextRebindingDoesNotDoubleConsume(t *testing.T) {
	tr := &contextClientTransport{}
	c := NewClient(tr, 1)
	old, cancel := context.WithCancel(t.Context())
	cancel()
	ctx := wirebudget.With(t.Context(), 1, 0)
	values, err := c.WithContext(old).WithContext(ctx).ReadHoldingRegisters(0, 1)
	require.NoError(t, err)
	require.Equal(t, []uint16{9}, values)
	require.Equal(t, 1, wirebudget.From(ctx).Used())
	require.Equal(t, 1, tr.sends)
	live, err := c.WithContext(old).IsConnectedContext(ctx)
	require.NoError(t, err)
	require.True(t, live)
}
