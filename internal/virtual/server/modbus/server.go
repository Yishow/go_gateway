package modbus

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"

	"go-gateway/internal/virtual/memory"
)

// Config controls the network and holding-register boundary of a Modbus TCP server.
type Config struct {
	BindAddress       string
	Port              int
	SlaveID           uint8
	CapacityRegisters int
}

const defaultBindAddress = "127.0.0.1"

// =============================================================================
// Modbus TCP Server
// =============================================================================

const (
	// Modbus 功能碼
	FuncReadCoils              = 0x01
	FuncReadDiscreteInputs     = 0x02
	FuncReadHoldingRegisters   = 0x03
	FuncReadInputRegisters     = 0x04
	FuncWriteSingleCoil        = 0x05
	FuncWriteSingleRegister    = 0x06
	FuncWriteMultipleCoils     = 0x0F
	FuncWriteMultipleRegisters = 0x10

	// Modbus 異常碼
	ExceptionIllegalFunction    = 0x01
	ExceptionIllegalDataAddress = 0x02
	ExceptionIllegalDataValue   = 0x03
	ExceptionSlaveDeviceFailure = 0x04

	// MBAP Header 長度
	MBAPHeaderLength = 7
)

// Server Modbus TCP 伺服器
type Server struct {
	mu       sync.RWMutex
	bank     *memory.MemoryBank
	listener net.Listener
	conns    map[net.Conn]struct{}
	bindAddr string
	port     int
	slaveID  uint8
	capacity int
	running  bool
	wg       sync.WaitGroup
	done     chan struct{}
}

// NewServer 建立新的 Modbus TCP 伺服器
func NewServer(bank *memory.MemoryBank) *Server {
	return &Server{
		bank:     bank,
		done:     make(chan struct{}),
		conns:    make(map[net.Conn]struct{}),
		bindAddr: defaultBindAddress,
		slaveID:  1,
		capacity: bank.Size() / 2,
	}
}

// Start 啟動伺服器
func (s *Server) Start(port int) error {
	return s.StartWithConfig(context.Background(), Config{BindAddress: defaultBindAddress, Port: port, SlaveID: 1, CapacityRegisters: s.bank.Size() / 2})
}

// StartWithConfig starts the server on the configured address and port.
func (s *Server) StartWithConfig(ctx context.Context, config Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return errors.New("伺服器已在運行中")
	}

	if config.BindAddress == "" {
		config.BindAddress = defaultBindAddress
	}
	if config.Port < 0 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d", config.Port)
	}
	if config.SlaveID == 0 {
		config.SlaveID = 1
	}
	if config.CapacityRegisters <= 0 {
		config.CapacityRegisters = s.bank.Size() / 2
	}
	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(ctx, "tcp", net.JoinHostPort(config.BindAddress, strconv.Itoa(config.Port)))
	if err != nil {
		return fmt.Errorf("監聽端口失敗: %w", err)
	}

	s.listener = listener
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		_ = listener.Close()
		return fmt.Errorf("監聽地址類型錯誤: %T", listener.Addr())
	}
	s.port = tcpAddr.Port
	s.bindAddr = config.BindAddress
	s.slaveID = config.SlaveID
	s.capacity = config.CapacityRegisters
	s.running = true
	s.done = make(chan struct{})
	s.conns = make(map[net.Conn]struct{})

	s.wg.Add(1)
	go s.acceptLoop(listener)

	return nil
}

// Stop 停止伺服器
func (s *Server) Stop() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}

	s.running = false
	close(s.done)
	listener := s.listener
	s.listener = nil
	conns := make([]net.Conn, 0, len(s.conns))
	for conn := range s.conns {
		conns = append(conns, conn)
	}
	s.port = 0
	s.mu.Unlock()

	if listener != nil {
		_ = listener.Close()
	}
	for _, conn := range conns {
		_ = conn.Close()
	}

	s.wg.Wait()
	return nil
}

// Port 返回伺服器端口
func (s *Server) Port() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.running {
		return 0
	}
	return s.port
}

// Address 返回伺服器地址
func (s *Server) Address() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return net.JoinHostPort(s.bindAddr, strconv.Itoa(s.port))
}

// BindAddress returns the address used by the active or last configured listener.
func (s *Server) BindAddress() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.bindAddr
}

// SlaveID returns the configured Modbus unit identifier.
func (s *Server) SlaveID() uint8 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.slaveID
}

// CapacityRegisters returns the configured holding-register capacity.
func (s *Server) CapacityRegisters() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.capacity
}

func (s *Server) registerRangeValid(start, quantity uint16) bool {
	s.mu.RLock()
	capacity := s.capacity
	s.mu.RUnlock()
	if capacity <= 0 {
		capacity = s.bank.Size() / 2
	}
	return capacity > 0 && int(start)+int(quantity) <= capacity
}
