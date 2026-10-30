package mcprotocol

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

type exactWordTransport struct {
	body    []byte
	request []byte
	sends   int
}

func (*exactWordTransport) Connect() error { return nil }
func (*exactWordTransport) Close() error   { return nil }
func (t *exactWordTransport) SendReceive(request []byte) ([]byte, error) {
	t.sends++
	t.request = append([]byte(nil), request...)
	return append([]byte(nil), t.body...), nil
}
func TestBatchReadWordExactKeepsWireAndLegacyPolicy(t *testing.T) {
	legacy := &exactWordTransport{body: []byte{0, 0, 0x34, 0x12}}
	exact := &exactWordTransport{body: append([]byte(nil), legacy.body...)}
	a, err := NewClientWithTransport(legacy).BatchReadWord("D", 100, 1)
	require.NoError(t, err)
	b, err := NewClientWithTransport(exact).BatchReadWordExact("D", 100, 1)
	require.NoError(t, err)
	require.Equal(t, a, b)
	require.Equal(t, legacy.request, exact.request)
	for _, payload := range [][]byte{{}, {1}, {1, 2, 3}, {1, 2, 3, 4}} {
		exact.body = append([]byte{0, 0}, payload...)
		_, err = NewClientWithTransport(exact).BatchReadWordExact("D", 100, 1)
		require.Error(t, err)
	}
	legacy.body = []byte{0, 0, 0x34, 0x12, 0xff}
	a, err = NewClientWithTransport(legacy).BatchReadWord("D", 100, 1)
	require.NoError(t, err)
	require.Equal(t, []int{0x1234}, a)
}
func TestBatchReadWordExactPreservesEndCodeAndCancellation(t *testing.T) {
	tr := &exactWordTransport{body: make([]byte, 3)}
	binary.LittleEndian.PutUint16(tr.body, 0xc051)
	_, err := NewClientWithTransport(tr).BatchReadWordExact("D", 0, 1)
	var mcerr *MCError
	require.True(t, errors.As(err, &mcerr))
	require.Equal(t, 0xc051, mcerr.Code)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err = NewClientWithTransport(tr).WithContext(ctx).BatchReadWordExact("D", 0, 1)
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, 1, tr.sends)
}
