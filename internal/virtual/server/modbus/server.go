package modbus

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"go-gateway/internal/virtual/memory"
)

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
	port     int
	running  bool
	wg       sync.WaitGroup
	done     chan struct{}
}

// NewServer 建立新的 Modbus TCP 伺服器
func NewServer(bank *memory.MemoryBank) *Server {
	return &Server{
		bank: bank,
		done: make(chan struct{}),
	}
}

// Start 啟動伺服器
func (s *Server) Start(port int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return errors.New("伺服器已在運行中")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("監聽端口失敗: %w", err)
	}

	s.listener = listener
	s.port = listener.Addr().(*net.TCPAddr).Port
	s.running = true
	s.done = make(chan struct{})

	s.wg.Add(1)
	go s.acceptLoop()

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
	s.listener.Close()
	s.port = 0
	s.mu.Unlock()

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
	return fmt.Sprintf("127.0.0.1:%d", s.Port())
}
