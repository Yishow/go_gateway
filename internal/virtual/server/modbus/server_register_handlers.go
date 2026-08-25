package modbus

import "encoding/binary"

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
	if !s.registerRangeValid(startAddr, quantity) {
		return s.exceptionResponse(FuncReadHoldingRegisters, ExceptionIllegalDataAddress)
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
	if !s.registerRangeValid(startAddr, quantity) {
		return s.exceptionResponse(FuncReadInputRegisters, ExceptionIllegalDataAddress)
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
	if !s.registerRangeValid(regAddr, 1) {
		return s.exceptionResponse(FuncWriteSingleRegister, ExceptionIllegalDataAddress)
	}

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
	if !s.registerRangeValid(startAddr, quantity) {
		return s.exceptionResponse(FuncWriteMultipleRegisters, ExceptionIllegalDataAddress)
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
