package modbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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
	s.mu.Unlock()

	s.wg.Wait()
	return nil
}

// Port 返回伺服器端口
func (s *Server) Port() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.port
}

// Address 返回伺服器地址
func (s *Server) Address() string {
	return fmt.Sprintf("127.0.0.1:%d", s.Port())
}

// acceptLoop 接受連接循環
func (s *Server) acceptLoop() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.done:
				return
			default:
				continue
			}
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

// handleConnection 處理單一連接
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	for {
		select {
		case <-s.done:
			return
		default:
		}

		// 讀取 MBAP Header
		header := make([]byte, MBAPHeaderLength)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			return
		}

		// 解析 MBAP Header
		transactionID := binary.BigEndian.Uint16(header[0:2])
		protocolID := binary.BigEndian.Uint16(header[2:4])
		length := binary.BigEndian.Uint16(header[4:6])
		unitID := header[6]

		// 驗證協議 ID (必須為 0)
		if protocolID != 0 {
			continue
		}

		// 讀取 PDU
		pdu := make([]byte, length-1) // length 包含 Unit ID
		_, err = io.ReadFull(conn, pdu)
		if err != nil {
			return
		}

		// 處理請求
		response := s.handleRequest(pdu)

		// 構建完整回應 (MBAP Header + PDU)
		respLen := len(response) + 1 // +1 for Unit ID
		fullResponse := make([]byte, MBAPHeaderLength+len(response))
		binary.BigEndian.PutUint16(fullResponse[0:2], transactionID)
		binary.BigEndian.PutUint16(fullResponse[2:4], 0) // Protocol ID
		binary.BigEndian.PutUint16(fullResponse[4:6], uint16(respLen))
		fullResponse[6] = unitID
		copy(fullResponse[7:], response)

		// 發送完整回應
		conn.Write(fullResponse)
	}
}

// handleRequest 處理 Modbus 請求
func (s *Server) handleRequest(pdu []byte) []byte {
	if len(pdu) < 1 {
		return s.exceptionResponse(0, ExceptionIllegalFunction)
	}

	funcCode := pdu[0]

	switch funcCode {
	case FuncReadHoldingRegisters:
		return s.handleReadHoldingRegisters(pdu)
	case FuncReadInputRegisters:
		return s.handleReadInputRegisters(pdu)
	case FuncWriteSingleRegister:
		return s.handleWriteSingleRegister(pdu)
	case FuncWriteMultipleRegisters:
		return s.handleWriteMultipleRegisters(pdu)
	case FuncReadCoils:
		return s.handleReadCoils(pdu)
	case FuncWriteSingleCoil:
		return s.handleWriteSingleCoil(pdu)
	default:
		return s.exceptionResponse(funcCode, ExceptionIllegalFunction)
	}
}

// handleReadHoldingRegisters 處理讀取保持暫存器 (FC 03)
func (s *Server) handleReadHoldingRegisters(pdu []byte) []byte {
	if len(pdu) < 5 {
		return s.exceptionResponse(FuncReadHoldingRegisters, ExceptionIllegalDataValue)
	}

	startAddr := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])

	if quantity == 0 || quantity > 125 {
		return s.exceptionResponse(FuncReadHoldingRegisters, ExceptionIllegalDataValue)
	}

	// 讀取數據
	byteCount := quantity * 2
	data := make([]byte, byteCount)

	for i := uint16(0); i < quantity; i++ {
		offset := int(startAddr+i) * 2
		val, err := s.bank.ReadWord(offset)
		if err != nil {
			return s.exceptionResponse(FuncReadHoldingRegisters, ExceptionIllegalDataAddress)
		}
		binary.BigEndian.PutUint16(data[i*2:], val)
	}

	// 構建回應
	response := make([]byte, 2+len(data))
	response[0] = FuncReadHoldingRegisters
	response[1] = byte(byteCount)
	copy(response[2:], data)

	return response
}

// handleReadInputRegisters 處理讀取輸入暫存器 (FC 04)
func (s *Server) handleReadInputRegisters(pdu []byte) []byte {
	// 與 Holding Registers 相同處理
	if len(pdu) < 5 {
		return s.exceptionResponse(FuncReadInputRegisters, ExceptionIllegalDataValue)
	}

	startAddr := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])

	if quantity == 0 || quantity > 125 {
		return s.exceptionResponse(FuncReadInputRegisters, ExceptionIllegalDataValue)
	}

	byteCount := quantity * 2
	data := make([]byte, byteCount)

	for i := uint16(0); i < quantity; i++ {
		offset := int(startAddr+i) * 2
		val, err := s.bank.ReadWord(offset)
		if err != nil {
			return s.exceptionResponse(FuncReadInputRegisters, ExceptionIllegalDataAddress)
		}
		binary.BigEndian.PutUint16(data[i*2:], val)
	}

	response := make([]byte, 2+len(data))
	response[0] = FuncReadInputRegisters
	response[1] = byte(byteCount)
	copy(response[2:], data)

	return response
}

// handleWriteSingleRegister 處理寫入單一暫存器 (FC 06)
func (s *Server) handleWriteSingleRegister(pdu []byte) []byte {
	if len(pdu) < 5 {
		return s.exceptionResponse(FuncWriteSingleRegister, ExceptionIllegalDataValue)
	}

	regAddr := binary.BigEndian.Uint16(pdu[1:3])
	value := binary.BigEndian.Uint16(pdu[3:5])

	offset := int(regAddr) * 2
	err := s.bank.WriteWord(offset, value)
	if err != nil {
		return s.exceptionResponse(FuncWriteSingleRegister, ExceptionIllegalDataAddress)
	}

	// 回應與請求相同
	return pdu[:5]
}

// handleWriteMultipleRegisters 處理寫入多個暫存器 (FC 16)
func (s *Server) handleWriteMultipleRegisters(pdu []byte) []byte {
	if len(pdu) < 6 {
		return s.exceptionResponse(FuncWriteMultipleRegisters, ExceptionIllegalDataValue)
	}

	startAddr := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])
	byteCount := pdu[5]

	if quantity == 0 || quantity > 123 || int(byteCount) != int(quantity)*2 {
		return s.exceptionResponse(FuncWriteMultipleRegisters, ExceptionIllegalDataValue)
	}

	if len(pdu) < 6+int(byteCount) {
		return s.exceptionResponse(FuncWriteMultipleRegisters, ExceptionIllegalDataValue)
	}

	// 寫入數據
	for i := uint16(0); i < quantity; i++ {
		offset := int(startAddr+i) * 2
		value := binary.BigEndian.Uint16(pdu[6+i*2:])
		err := s.bank.WriteWord(offset, value)
		if err != nil {
			return s.exceptionResponse(FuncWriteMultipleRegisters, ExceptionIllegalDataAddress)
		}
	}

	// 構建回應
	response := make([]byte, 5)
	response[0] = FuncWriteMultipleRegisters
	binary.BigEndian.PutUint16(response[1:3], startAddr)
	binary.BigEndian.PutUint16(response[3:5], quantity)

	return response
}

// handleReadCoils 處理讀取線圈 (FC 01)
func (s *Server) handleReadCoils(pdu []byte) []byte {
	if len(pdu) < 5 {
		return s.exceptionResponse(FuncReadCoils, ExceptionIllegalDataValue)
	}

	startAddr := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])

	if quantity == 0 || quantity > 2000 {
		return s.exceptionResponse(FuncReadCoils, ExceptionIllegalDataValue)
	}

	// 計算需要的字節數
	byteCount := (quantity + 7) / 8
	data := make([]byte, byteCount)

	for i := uint16(0); i < quantity; i++ {
		bitAddr := startAddr + i
		byteOffset := int(bitAddr / 8)
		bitOffset := bitAddr % 8

		val, err := s.bank.ReadByte(byteOffset)
		if err != nil {
			return s.exceptionResponse(FuncReadCoils, ExceptionIllegalDataAddress)
		}

		if (val & (1 << bitOffset)) != 0 {
			data[i/8] |= 1 << (i % 8)
		}
	}

	response := make([]byte, 2+len(data))
	response[0] = FuncReadCoils
	response[1] = byte(byteCount)
	copy(response[2:], data)

	return response
}

// handleWriteSingleCoil 處理寫入單一線圈 (FC 05)
func (s *Server) handleWriteSingleCoil(pdu []byte) []byte {
	if len(pdu) < 5 {
		return s.exceptionResponse(FuncWriteSingleCoil, ExceptionIllegalDataValue)
	}

	coilAddr := binary.BigEndian.Uint16(pdu[1:3])
	value := binary.BigEndian.Uint16(pdu[3:5])

	// Modbus 規範: 0xFF00 = ON, 0x0000 = OFF
	byteOffset := int(coilAddr / 8)
	bitOffset := coilAddr % 8

	currentVal, err := s.bank.ReadByte(byteOffset)
	if err != nil {
		return s.exceptionResponse(FuncWriteSingleCoil, ExceptionIllegalDataAddress)
	}

	if value == 0xFF00 {
		currentVal |= 1 << bitOffset
	} else {
		currentVal &^= 1 << bitOffset
	}

	err = s.bank.WriteByte(byteOffset, currentVal)
	if err != nil {
		return s.exceptionResponse(FuncWriteSingleCoil, ExceptionIllegalDataAddress)
	}

	// 回應與請求相同
	return pdu[:5]
}

// exceptionResponse 構建異常回應
func (s *Server) exceptionResponse(funcCode byte, exceptionCode byte) []byte {
	return []byte{funcCode | 0x80, exceptionCode}
}
