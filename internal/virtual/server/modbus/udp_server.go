package modbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net"
	"sync"

	"go-gateway/internal/virtual/memory"
)

// UDPServer provides a local Modbus UDP server backed by virtual memory.
type UDPServer struct {
	mu sync.RWMutex

	bank    *memory.MemoryBank
	conn    *net.UDPConn
	port    int
	done    chan struct{}
	running bool
	wg      sync.WaitGroup

	requestHandler *Server
}

// NewUDPServer creates a Modbus UDP server.
func NewUDPServer(bank *memory.MemoryBank) *UDPServer {
	return &UDPServer{
		bank:           bank,
		done:           make(chan struct{}),
		requestHandler: &Server{bank: bank},
	}
}

// Start launches the UDP server at the target port. port=0 selects a random free port.
func (s *UDPServer) Start(port int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return errors.New("伺服器已在運行中")
	}

	udpAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: port}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("監聽 UDP 端口失敗: %w", err)
	}

	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		_ = conn.Close()
		return fmt.Errorf("監聽地址類型錯誤: %T", conn.LocalAddr())
	}

	s.conn = conn
	s.port = localAddr.Port
	s.done = make(chan struct{})
	s.running = true

	s.wg.Add(1)
	go s.readLoop()

	return nil
}

// Stop stops the UDP server.
func (s *UDPServer) Stop() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}

	s.running = false
	close(s.done)
	if s.conn != nil {
		_ = s.conn.Close()
	}
	s.port = 0
	s.mu.Unlock()

	s.wg.Wait()
	return nil
}

// Port returns current UDP listening port.
func (s *UDPServer) Port() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.running {
		return 0
	}
	return s.port
}

// Address returns endpoint address.
func (s *UDPServer) Address() string {
	return fmt.Sprintf("127.0.0.1:%d", s.Port())
}

func (s *UDPServer) readLoop() {
	defer s.wg.Done()

	buffer := make([]byte, 1024)
	for {
		n, addr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				continue
			}
		}

		req := make([]byte, n)
		copy(req, buffer[:n])

		resp, ok := s.handleFrame(req)
		if !ok {
			continue
		}

		if _, err := s.conn.WriteToUDP(resp, addr); err != nil && !errors.Is(err, net.ErrClosed) {
			slog.Error("write Modbus UDP response", "error", err)
		}
	}
}

func (s *UDPServer) handleFrame(frame []byte) ([]byte, bool) {
	if len(frame) < MBAPHeaderLength+1 {
		return nil, false
	}

	transactionID := binary.BigEndian.Uint16(frame[0:2])
	protocolID := binary.BigEndian.Uint16(frame[2:4])
	length := int(binary.BigEndian.Uint16(frame[4:6]))
	unitID := frame[6]

	if protocolID != 0 || length < 2 {
		return nil, false
	}
	if len(frame) < 6+length {
		return nil, false
	}

	pduEnd := 6 + length
	pdu := frame[MBAPHeaderLength:pduEnd]
	responsePDU := s.requestHandler.handleRequest(pdu)

	pduLength := len(responsePDU)
	if pduLength > math.MaxUint16-1 {
		return nil, false
	}
	response := make([]byte, MBAPHeaderLength+pduLength)
	binary.BigEndian.PutUint16(response[0:2], transactionID)
	binary.BigEndian.PutUint16(response[2:4], 0)
	binary.BigEndian.PutUint16(response[4:6], uint16(pduLength)+1)
	response[6] = unitID
	copy(response[7:], responsePDU)

	return response, true
}
