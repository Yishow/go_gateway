package adapters

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"testing"
)

type countedMCTransport struct{ sends, closes int }

func (t *countedMCTransport) Connect() error                     { return nil }
func (t *countedMCTransport) Close() error                       { t.closes++; return nil }
func (t *countedMCTransport) SendReceive([]byte) ([]byte, error) { t.sends++; return nil, io.EOF }

func TestMCReadCanceledBeforeEntryDoesNotSend(t *testing.T) {
	c := newRetryTestConnector(0)
	transport := &countedMCTransport{}
	c.client = mcprotocol.NewClientWithTransport(transport)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	_, err := c.Read(ctx, connector.ReadRequest{Address: "D0", Count: 1})
	require.ErrorIs(t, err, context.Canceled)
	require.Zero(t, transport.sends)
	require.ErrorIs(t, c.Reconnect(ctx), context.Canceled)
	require.Zero(t, transport.closes)
}
func TestMCAdapterUsesCallersAttemptBudget(t *testing.T) {
	for _, limit := range []int{1, 2, 4} {
		c := newRetryTestConnector(0)
		transport := &countedMCTransport{}
		c.client = mcprotocol.NewClientWithTransport(transport)
		ctx := wirebudget.With(t.Context(), limit, 0)
		_, err := c.Read(ctx, connector.ReadRequest{Address: "D0", Count: 1})
		require.ErrorIs(t, err, io.EOF)
		require.Equal(t, limit, transport.sends)
		require.Equal(t, limit, wirebudget.From(ctx).Used())
	}
}

func TestMCRetryDoesNotGuessConnectionErrorsFromText(t *testing.T) {
	require.False(t, mcConnectionFailure(errors.New("invalid address EOF")))
	require.False(t, mcConnectionFailure(&mcprotocol.MCError{Code: 0xC059, Msg: "EOF"}))
	require.True(t, mcConnectionFailure(io.EOF))
}
