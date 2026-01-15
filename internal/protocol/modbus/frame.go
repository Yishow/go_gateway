package modbus

import (
	"encoding/binary"
	"fmt"
)

// MBAPHeader Modbus TCP/UDP MBAP 標頭
type MBAPHeader struct {
	TransactionID uint16 // 交易 ID
	ProtocolID   uint16 // 協議 ID (0x0000 for Modbus)
	Length       uint16 // 後續數據長度 (Unit ID + PDU)
	UnitID       byte   // 單元 ID (站號)
}

// BuildMBAPHeader 構建 MBAP 標頭
func BuildMBAPHeader(transactionID uint16, unitID byte, pduLength int) []byte {
	header := make([]byte, MBAPHeaderLength)
	binary.BigEndian.PutUint16(header[0:], transactionID)
	binary.BigEndian.PutUint16(header[2:], 0x0000) // Protocol ID
	binary.BigEndian.PutUint16(header[4:], uint16(pduLength+1)) // Length = PDU + Unit ID
	header[6] = unitID
	return header
}

// ParseMBAPHeader 解析 MBAP 標頭
func ParseMBAPHeader(data []byte) (*MBAPHeader, error) {
	if len(data) < MBAPHeaderLength {
		return nil, ErrResponseTooShort
	}

	header := &MBAPHeader{
		TransactionID: binary.BigEndian.Uint16(data[0:]),
		ProtocolID:    binary.BigEndian.Uint16(data[2:]),
		Length:        binary.BigEndian.Uint16(data[4:]),
		UnitID:        data[6],
	}

	if header.ProtocolID != 0x0000 {
		return nil, fmt.Errorf("invalid protocol ID: 0x%04X", header.ProtocolID)
	}

	return header, nil
}

// BuildPDU 構建協議數據單元 (PDU)
// PDU = Function Code + Data
func BuildPDU(functionCode byte, data []byte) []byte {
	pdu := make([]byte, 1+len(data))
	pdu[0] = functionCode
	copy(pdu[1:], data)
	return pdu
}

// ParsePDU 解析協議數據單元
func ParsePDU(data []byte) (functionCode byte, pduData []byte, err error) {
	if len(data) < 1 {
		return 0, nil, ErrResponseTooShort
	}

	functionCode = data[0]
	
	// 檢查是否為異常回應
	if functionCode&FuncExceptionOffset != 0 {
		if len(data) < 2 {
			return 0, nil, ErrResponseTooShort
		}
		exceptionCode := data[1]
		return 0, nil, NewProtocolError(exceptionCode, "Modbus 異常回應")
	}

	pduData = data[1:]
	return functionCode, pduData, nil
}

// BuildTCPFrame 構建 Modbus TCP 封包
// Frame = MBAP Header + PDU
func BuildTCPFrame(transactionID uint16, unitID byte, functionCode byte, data []byte) []byte {
	pdu := BuildPDU(functionCode, data)
	mbap := BuildMBAPHeader(transactionID, unitID, len(pdu))
	return append(mbap, pdu...)
}

// ParseTCPFrame 解析 Modbus TCP 封包
func ParseTCPFrame(data []byte) (transactionID uint16, unitID byte, functionCode byte, pduData []byte, err error) {
	header, err := ParseMBAPHeader(data)
	if err != nil {
		return 0, 0, 0, nil, err
	}

	if len(data) < MBAPHeaderLength+1 {
		return 0, 0, 0, nil, ErrResponseTooShort
	}

	pdu := data[MBAPHeaderLength:]
	functionCode, pduData, err = ParsePDU(pdu)
	if err != nil {
		return 0, 0, 0, nil, err
	}

	return header.TransactionID, header.UnitID, functionCode, pduData, nil
}

// CalculateCRC16 計算 Modbus RTU CRC-16 校驗碼
// 使用 CRC-16 (Modbus) 多項式: 0x8005 (反向)
func CalculateCRC16(data []byte) uint16 {
	crc := uint16(0xFFFF)
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if (crc & 0x0001) != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

// BuildRTUFrame 構建 Modbus RTU 封包
// Frame = Address(1) + PDU + CRC(2)
func BuildRTUFrame(address byte, functionCode byte, data []byte) []byte {
	pdu := BuildPDU(functionCode, data)
	frame := make([]byte, 1+len(pdu)+2)
	frame[0] = address
	copy(frame[1:], pdu)
	
	crc := CalculateCRC16(frame[:1+len(pdu)])
	binary.LittleEndian.PutUint16(frame[1+len(pdu):], crc)
	
	return frame
}

// ParseRTUFrame 解析 Modbus RTU 封包
func ParseRTUFrame(data []byte) (address byte, functionCode byte, pduData []byte, err error) {
	if len(data) < RTUFrameMinLength {
		return 0, 0, nil, ErrResponseTooShort
	}

	address = data[0]
	pdu := data[1 : len(data)-2]
	
	// 驗證 CRC
	receivedCRC := binary.LittleEndian.Uint16(data[len(data)-2:])
	calculatedCRC := CalculateCRC16(data[:len(data)-2])
	if receivedCRC != calculatedCRC {
		return 0, 0, nil, ErrCRCError
	}

	functionCode, pduData, err = ParsePDU(pdu)
	if err != nil {
		return 0, 0, nil, err
	}

	return address, functionCode, pduData, nil
}

// BuildReadRequest 構建讀取請求數據
// Data = Starting Address(2) + Quantity(2)
func BuildReadRequest(startAddress uint16, quantity uint16) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint16(data[0:], startAddress)
	binary.BigEndian.PutUint16(data[2:], quantity)
	return data
}

// ParseReadResponse 解析讀取回應數據
// Response = Byte Count(1) + Data(N)
func ParseReadResponse(data []byte) ([]byte, error) {
	if len(data) < 1 {
		return nil, ErrResponseTooShort
	}

	byteCount := int(data[0])
	if len(data) < 1+byteCount {
		return nil, ErrResponseTooShort
	}

	return data[1 : 1+byteCount], nil
}

// BuildWriteSingleCoilRequest 構建寫入單個線圈請求
// Data = Output Address(2) + Output Value(2)
func BuildWriteSingleCoilRequest(address uint16, value bool) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint16(data[0:], address)
	if value {
		binary.BigEndian.PutUint16(data[2:], 0xFF00) // ON
	} else {
		binary.BigEndian.PutUint16(data[2:], 0x0000) // OFF
	}
	return data
}

// BuildWriteSingleRegisterRequest 構建寫入單個暫存器請求
// Data = Register Address(2) + Register Value(2)
func BuildWriteSingleRegisterRequest(address uint16, value uint16) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint16(data[0:], address)
	binary.BigEndian.PutUint16(data[2:], value)
	return data
}

// BuildWriteMultipleCoilsRequest 構建寫入多個線圈請求
// Data = Starting Address(2) + Quantity(2) + Byte Count(1) + Coil Values(N)
func BuildWriteMultipleCoilsRequest(startAddress uint16, values []bool) []byte {
	quantity := uint16(len(values))
	byteCount := (len(values) + 7) / 8
	
	data := make([]byte, 5+byteCount)
	binary.BigEndian.PutUint16(data[0:], startAddress)
	binary.BigEndian.PutUint16(data[2:], quantity)
	data[4] = byte(byteCount)
	
	// 打包位元值
	for i, val := range values {
		if val {
			data[5+i/8] |= 1 << (i % 8)
		}
	}
	
	return data
}

// BuildWriteMultipleRegistersRequest 構建寫入多個暫存器請求
// Data = Starting Address(2) + Quantity(2) + Byte Count(1) + Register Values(2*N)
func BuildWriteMultipleRegistersRequest(startAddress uint16, values []uint16) []byte {
	quantity := uint16(len(values))
	byteCount := len(values) * 2
	
	data := make([]byte, 5+byteCount)
	binary.BigEndian.PutUint16(data[0:], startAddress)
	binary.BigEndian.PutUint16(data[2:], quantity)
	data[4] = byte(byteCount)
	
	for i, val := range values {
		binary.BigEndian.PutUint16(data[5+i*2:], val)
	}
	
	return data
}

// ParseWriteResponse 解析寫入回應
// Single Write Response = Address(2) + Value(2) or Quantity(2)
func ParseWriteResponse(data []byte, isMultiple bool) (address uint16, valueOrQuantity uint16, err error) {
	if len(data) < 4 {
		return 0, 0, ErrResponseTooShort
	}
	
	address = binary.BigEndian.Uint16(data[0:])
	valueOrQuantity = binary.BigEndian.Uint16(data[2:])
	return address, valueOrQuantity, nil
}

// UnpackBits 將位元組解包為布林值陣列
func UnpackBits(data []byte, count int) []bool {
	result := make([]bool, count)
	for i := 0; i < count; i++ {
		byteIdx := i / 8
		bitIdx := i % 8
		if byteIdx < len(data) {
			result[i] = (data[byteIdx]>>bitIdx)&1 != 0
		}
	}
	return result
}

// PackBits 將布林值陣列打包為位元組
func PackBits(values []bool) []byte {
	byteCount := (len(values) + 7) / 8
	result := make([]byte, byteCount)
	for i, val := range values {
		if val {
			result[i/8] |= 1 << (i % 8)
		}
	}
	return result
}
