package modbus

import (
	"encoding/binary"
	"io"
	"net"
)

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
