package modbus

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
)

func plannedUDPServer(t *testing.T, reply func([]byte) ([]byte, error)) (*UDPTransport, <-chan error) {
	t.Helper()
	peer, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	require.NoError(t, err)
	t.Cleanup(func() { _ = peer.Close() })
	done := make(chan error, 1)
	go func() {
		_ = peer.SetReadDeadline(time.Now().Add(2 * time.Second))
		data := make([]byte, 300)
		n, addr, err := peer.ReadFromUDP(data)
		if err != nil {
			done <- err
			return
		}
		response, err := reply(data[:n])
		if err != nil {
			done <- err
			return
		}
		if response != nil {
			_, err = peer.WriteToUDP(response, addr)
		}
		done <- err
	}()
	tr := NewUDPTransport("127.0.0.1", peer.LocalAddr().(*net.UDPAddr).Port)
	tr.Timeout = time.Second
	require.NoError(t, tr.ConnectContext(t.Context()))
	t.Cleanup(func() { _ = tr.Close() })
	return tr, done
}
func TestPlannedExactUDP125WordsAreNotTruncated(t *testing.T) {
	peer := plannedBasicPeer(false)
	for i := uint16(0); i < 125; i++ {
		peer.words[plannedCell{3, i}] = i + 1
	}
	tr, done := plannedUDPServer(t, peer.SendReceive)
	ctx := wirebudget.With(t.Context(), 1, 0)
	values, err := NewClient(tr, 1).WithContext(ctx).ReadWordsExact(3, 0, 125)
	require.NoError(t, err)
	require.Len(t, values, 125)
	for i, v := range values {
		require.Equal(t, uint16(i+1), v)
	}
	require.NoError(t, <-done)
	require.Equal(t, 1, wirebudget.From(ctx).Used())
	require.Equal(t, []byte{0, 1, 0, 0, 0, 6, 1, 3, 0, 0, 0, 125}, peer.frames[0])
}
func TestPlannedUDPDatagramFrameBoundary(t *testing.T) {
	for _, size := range []int{260, 261, 512} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			tr, done := plannedUDPServer(t, func(req []byte) ([]byte, error) {
				if len(req) != 12 || req[7] != 3 {
					return nil, fmt.Errorf("unexpected UDP request")
				}
				response := make([]byte, size)
				copy(response, req[:8])
				binary.BigEndian.PutUint16(response[4:6], uint16(size-6))
				return response, nil
			})
			ctx := wirebudget.With(t.Context(), 1, 0)
			data, err := tr.SendReceiveContext(ctx, []byte{0, 1, 0, 0, 0, 6, 1, 3, 0, 0, 0, 1})
			if size == 260 {
				require.NoError(t, err)
				require.Len(t, data, 260)
			} else {
				require.ErrorIs(t, err, ErrInvalidFrame)
				require.Nil(t, data)
			}
			require.NoError(t, <-done)
			require.Equal(t, 1, wirebudget.From(ctx).Used())
		})
	}
}
func TestPlannedExactUDPCancelAndDeadlinePreserved(t *testing.T) {
	for _, cancelActive := range []bool{false, true} {
		t.Run(fmt.Sprint(cancelActive), func(t *testing.T) {
			received := make(chan struct{})
			tr, done := plannedUDPServer(t, func(req []byte) ([]byte, error) { close(received); return nil, nil })
			budget := wirebudget.With(t.Context(), 1, 0)
			ctx, cancel := context.WithTimeout(budget, 25*time.Millisecond)
			defer cancel()
			if cancelActive {
				go func() { <-received; cancel() }()
			}
			_, err := NewClient(tr, 1).WithContext(ctx).ReadWordsExact(3, 0, 125)
			if cancelActive {
				require.ErrorIs(t, err, context.Canceled)
			} else {
				require.ErrorIs(t, err, context.DeadlineExceeded)
			}
			require.NoError(t, <-done)
			require.Equal(t, 1, wirebudget.From(ctx).Used())
			_, err = NewClient(tr, 1).WithContext(ctx).ReadWordsExact(3, 0, 125)
			require.Error(t, err)
			require.Equal(t, 1, wirebudget.From(ctx).Used())
		})
	}
}
