package modbus

import "encoding/binary"

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

		val, err := s.bank.ReadByteAt(byteOffset)
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

	currentVal, err := s.bank.ReadByteAt(byteOffset)
	if err != nil {
		return s.exceptionResponse(FuncWriteSingleCoil, ExceptionIllegalDataAddress)
	}

	if value == 0xFF00 {
		currentVal |= 1 << bitOffset
	} else {
		currentVal &^= 1 << bitOffset
	}

	err = s.bank.WriteByteAt(byteOffset, currentVal)
	if err != nil {
		return s.exceptionResponse(FuncWriteSingleCoil, ExceptionIllegalDataAddress)
	}

	// 回應與請求相同
	return pdu[:5]
}

// exceptionResponse 構建異常回應
func (s *Server) exceptionResponse(funcCode, exceptionCode byte) []byte {
	return []byte{funcCode | 0x80, exceptionCode}
}
