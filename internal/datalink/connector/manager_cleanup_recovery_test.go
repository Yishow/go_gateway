package connector

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

type recoverableCleanupProtocol struct {
	generationProtocol
	failClose atomic.Bool
}

func (p *recoverableCleanupProtocol) Close() error {
	p.closed.Add(1)
	if p.failClose.Load() {
		return errors.New("private-endpoint secret")
	}
	p.live.Store(false)
	return nil
}
func TestGenerationCleanupFailureObservableAndRetryable(t *testing.T) {
	for _, action := range []string{"close", "idle", "failed-dial", "late-dial"} {
		t.Run(action, func(t *testing.T) {
			typ := schema.ProtocolModbusRTU
			p := &recoverableCleanupProtocol{}
			p.failClose.Store(true)
			if action == "failed-dial" {
				p.dial = func(context.Context, string) error { return errors.New("dial private-endpoint") }
			}
			started, unblock := make(chan struct{}), make(chan struct{})
			if action == "late-dial" {
				p.dial = func(context.Context, string) error { close(started); <-unblock; return nil }
			}
			registerGeneration(t, typ, func() Protocol { return p })
			cm := NewConnectionManager(DefaultConnectionManagerConfig())
			if action == "late-dial" {
				ch := acquire(cm, t.Context(), "a", typ, `{"serial_port":"COM9"}`)
				<-started
				require.NoError(t, cm.Close("a"))
				close(unblock)
				require.Error(t, result(t, ch).err)
			} else {
				conn, err := cm.GetOrCreate(t.Context(), "a", typ, `{"serial_port":"COM9"}`)
				if action == "failed-dial" {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					if action == "idle" {
						conn.mu.Lock()
						conn.LastUsed = time.Now().Add(-time.Hour)
						conn.mu.Unlock()
						require.Zero(t, cm.CleanupIdle())
					} else {
						require.Error(t, cm.Close("a"))
					}
				}
			}
			status, ok := cm.GetStatus("a")
			require.True(t, ok, "unresolved cleanup must remain observable")
			require.True(t, status.CleanupPending)
			require.Contains(t, status.LastError, "cleanup")
			require.NotContains(t, status.LastError, "private-endpoint")
			require.NotContains(t, status.LastError, "secret")
			require.Len(t, cm.ListConnections(), 1)
			_, err := cm.GetOrCreate(t.Context(), "b", typ, `{"serial_port":"COM9"}`)
			require.Error(t, err)
			require.Error(t, cm.CloseAll(), "explicit retry must report unresolved cleanup")
			p.failClose.Store(false)
			require.NoError(t, cm.Close("a"))
			_, ok = cm.GetStatus("a")
			require.False(t, ok)
			require.Empty(t, cm.ListConnections())
			p.dial = nil
			_, err = cm.GetOrCreate(t.Context(), "b", typ, `{"serial_port":"COM9"}`)
			require.NoError(t, err)
		})
	}
}
func TestGenerationRelativeAbsoluteSerialAlias(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "port"), nil, 0600))
	require.NoError(t, os.Symlink("port", filepath.Join(dir, "alias")))
	t.Chdir(dir)
	typ := schema.ProtocolModbusRTU
	registerGeneration(t, typ, func() Protocol { return &generationProtocol{} })
	for _, relative := range []string{"port", "alias"} {
		cm := NewConnectionManager(DefaultConnectionManagerConfig())
		_, err := cm.GetOrCreate(t.Context(), "a", typ, `{"serial_port":"`+filepath.Join(dir, "port")+`"}`)
		require.NoError(t, err)
		_, err = cm.GetOrCreate(t.Context(), "b", typ, `{"serial_port":"`+relative+`"}`)
		require.Error(t, err, "same physical port must share reservation identity")
	}
}

func TestGenerationLateFailedCleanupDoesNotCloseNewerConnection(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	started, unblock := make(chan struct{}), make(chan struct{})
	old := &recoverableCleanupProtocol{}
	old.failClose.Store(true)
	old.dial = func(context.Context, string) error { close(started); <-unblock; return nil }
	newer := &generationProtocol{}
	var calls atomic.Int32
	registerGeneration(t, typ, func() Protocol {
		if calls.Add(1) == 1 {
			return old
		}
		return newer
	})
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	ch := acquire(cm, t.Context(), "a", typ, `{"host":"old"}`)
	<-started
	current, err := cm.GetOrCreate(t.Context(), "a", typ, `{"host":"new"}`)
	require.NoError(t, err)
	close(unblock)
	require.Error(t, result(t, ch).err)
	require.Same(t, newer, current.Protocol)
	require.True(t, newer.IsConnected())
	require.Zero(t, newer.closed.Load())
	status, ok := cm.GetStatus("a")
	require.True(t, ok)
	require.True(t, status.Connected)
	require.True(t, status.CleanupPending)
	old.failClose.Store(false)
	require.NoError(t, cm.CloseAll())
	require.EqualValues(t, 1, newer.closed.Load())
	require.EqualValues(t, 2, old.closed.Load())
}
