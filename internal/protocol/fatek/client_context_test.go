package fatek

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"testing"
)

type contextClientTransport struct{ sends int }

func (*contextClientTransport) Connect() error                       { return nil }
func (*contextClientTransport) Close() error                         { return nil }
func (t *contextClientTransport) SendReceive([]byte) ([]byte, error) { t.sends++; return nil, io.EOF }
func TestFatekContextViewReplacesOldContext(t *testing.T) {
	transport := &contextClientTransport{}
	client := NewClient(transport, 3)
	old, cancel := context.WithCancel(t.Context())
	cancel()
	fresh := wirebudget.With(t.Context(), 2, 0)
	view := client.WithContext(old).WithContext(fresh)
	_, err := view.ReadRegisters("D", 0, 1)
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, 1, wirebudget.From(fresh).Used())
	require.Equal(t, 1, transport.sends)
	require.Same(t, transport, client.transport)
	require.Equal(t, 3, view.station)
}
