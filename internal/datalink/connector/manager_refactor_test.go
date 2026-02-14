package connector

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionManagerGetOrCreateReusesConnectedConnection(t *testing.T) {
	protoType := schema.ProtocolType("test_manager_reuse")
	var factoryCalls int32

	Register(protoType, func() Protocol {
		atomic.AddInt32(&factoryCalls, 1)
		return &stubProtocol{connected: true}
	})
	defer Unregister(protoType)

	cm := NewConnectionManager(DefaultConnectionManagerConfig())

	conn1, err := cm.GetOrCreate(context.Background(), "dev-1", protoType, "{}")
	require.NoError(t, err)
	conn2, err := cm.GetOrCreate(context.Background(), "dev-1", protoType, "{}")
	require.NoError(t, err)

	assert.Same(t, conn1, conn2)
	assert.Equal(t, int32(1), atomic.LoadInt32(&factoryCalls))
}

func TestConnectionManagerCleanupIdle(t *testing.T) {
	protoType := schema.ProtocolType("test_manager_cleanup")
	proto := &stubProtocol{connected: true}

	Register(protoType, func() Protocol { return proto })
	defer Unregister(protoType)

	cfg := DefaultConnectionManagerConfig()
	cfg.IdleTimeout = 10 * time.Millisecond
	cm := NewConnectionManager(cfg)

	conn, err := cm.GetOrCreate(context.Background(), "dev-clean", protoType, "{}")
	require.NoError(t, err)

	conn.mu.Lock()
	conn.inUse = false
	conn.LastUsed = time.Now().Add(-time.Minute)
	conn.mu.Unlock()

	cleaned := cm.CleanupIdle()
	assert.Equal(t, 1, cleaned)
	_, ok := cm.GetStatus("dev-clean")
	assert.False(t, ok)
	assert.Equal(t, 1, proto.closeCalls)
}

func TestManagedConnectionReadRecordsLastError(t *testing.T) {
	expectedErr := fmt.Errorf("read failed")
	mc := &ManagedConnection{
		DeviceID: "dev-read",
		Protocol: &stubProtocol{
			connected: true,
			readErr:   expectedErr,
		},
	}

	_, err := mc.Read(context.Background(), ReadRequest{Address: "D100"})
	require.Error(t, err)
	assert.Equal(t, expectedErr, mc.LastError)
}

type stubProtocol struct {
	connected  bool
	readErr    error
	writeErr   error
	testErr    error
	closeCalls int
}

func (s *stubProtocol) Connect(ctx context.Context, config string) error {
	s.connected = true
	return nil
}

func (s *stubProtocol) Close() error {
	s.closeCalls++
	s.connected = false
	return nil
}

func (s *stubProtocol) IsConnected() bool {
	return s.connected
}

func (s *stubProtocol) Read(ctx context.Context, req ReadRequest) (ReadResult, error) {
	if s.readErr != nil {
		return ReadResult{}, s.readErr
	}
	return ReadResult{Timestamp: time.Now(), Quality: schema.QualityGood}, nil
}

func (s *stubProtocol) Write(ctx context.Context, req WriteRequest) error {
	return s.writeErr
}

func (s *stubProtocol) TestConnection(ctx context.Context) error {
	return s.testErr
}

func (s *stubProtocol) ProtocolType() schema.ProtocolType {
	return schema.ProtocolType("test")
}
