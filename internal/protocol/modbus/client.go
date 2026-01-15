package modbus

import (
	"encoding/binary"
	"fmt"
	"sync"
)

// ModbusClient Modbus 客戶端
type ModbusClient struct {
	transport Transport
	unitID    byte
	mu        sync.Mutex
}

// NewClient 建立新的 Modbus 客戶端
func NewClient(transport Transport, unitID byte) *ModbusClient {
	if unitID == 0 {
		unitID = DefaultUnitID
	}
	return &ModbusClient{
		transport: transport,
		unitID:    unitID,
	}
}

// Connect 建立連線
func (c *ModbusClient) Connect() error {
	return c.transport.Connect()
}

// Close 關閉連線
func (c *ModbusClient) Close() error {
	return c.transport.Close()
}

// sendRequest 發送請求並接收回應 (TCP/UDP)
func (c *ModbusClient) sendTCPRequest(functionCode byte, data []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var transactionID uint16
	if tcpTransport, ok := c.transport.(*TCPTransport); ok {
		transactionID = tcpTransport.GetNextTransactionID()
	} else if udpTransport, ok := c.transport.(*UDPTransport); ok {
		transactionID = udpTransport.GetNextTransactionID()
	} else {
		return nil, fmt.Errorf("不支援的傳輸類型")
	}

	// 構建 TCP/UDP 封包
	frame := BuildTCPFrame(transactionID, c.unitID, functionCode, data)

	// 發送並接收
	response, err := c.transport.SendReceive(frame)
	if err != nil {
		return nil, err
	}

	// 解析回應
	respTransactionID, respUnitID, respFunctionCode, respData, err := ParseTCPFrame(response)
	if err != nil {
		return nil, err
	}

	// 驗證交易 ID (僅 TCP/UDP)
	if respTransactionID != transactionID {
		return nil, ErrTransactionIDMismatch
	}

	if respUnitID != c.unitID {
		return nil, fmt.Errorf("單元 ID 不匹配: 預期 %d, 實際 %d", c.unitID, respUnitID)
	}

	if respFunctionCode != functionCode {
		return nil, fmt.Errorf("功能碼不匹配: 預期 0x%02X, 實際 0x%02X", functionCode, respFunctionCode)
	}

	return respData, nil
}

// sendRTURequest 發送 RTU 請求並接收回應
func (c *ModbusClient) sendRTURequest(functionCode byte, data []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 構建 RTU 封包
	frame := BuildRTUFrame(c.unitID, functionCode, data)

	// 發送並接收
	response, err := c.transport.SendReceive(frame)
	if err != nil {
		return nil, err
	}

	// 解析回應
	respAddress, respFunctionCode, respData, err := ParseRTUFrame(response)
	if err != nil {
		return nil, err
	}

	if respAddress != c.unitID {
		return nil, fmt.Errorf("地址不匹配: 預期 %d, 實際 %d", c.unitID, respAddress)
	}

	if respFunctionCode != functionCode {
		return nil, fmt.Errorf("功能碼不匹配: 預期 0x%02X, 實際 0x%02X", functionCode, respFunctionCode)
	}

	return respData, nil
}

// sendRequest 根據傳輸類型選擇適當的發送方法
func (c *ModbusClient) sendRequest(functionCode byte, data []byte) ([]byte, error) {
	switch c.transport.(type) {
	case *TCPTransport, *UDPTransport:
		return c.sendTCPRequest(functionCode, data)
	case *RTUTransport:
		return c.sendRTURequest(functionCode, data)
	default:
		return nil, fmt.Errorf("不支援的傳輸類型")
	}
}

// ReadCoils 讀取線圈狀態
// address: 起始地址 (0-65535)
// quantity: 讀取數量 (1-2000)
func (c *ModbusClient) ReadCoils(address uint16, quantity uint16) ([]bool, error) {
	if quantity < 1 || quantity > CoilMaxQuantity {
		return nil, ErrInvalidQuantity
	}

	requestData := BuildReadRequest(address, quantity)
	responseData, err := c.sendRequest(FuncReadCoils, requestData)
	if err != nil {
		return nil, err
	}

	// 解析回應: Byte Count + Coil Values
	coilData, err := ParseReadResponse(responseData)
	if err != nil {
		return nil, err
	}

	return UnpackBits(coilData, int(quantity)), nil
}

// ReadDiscreteInputs 讀取離散輸入狀態
func (c *ModbusClient) ReadDiscreteInputs(address uint16, quantity uint16) ([]bool, error) {
	if quantity < 1 || quantity > DiscreteInputMaxQuantity {
		return nil, ErrInvalidQuantity
	}

	requestData := BuildReadRequest(address, quantity)
	responseData, err := c.sendRequest(FuncReadDiscreteInputs, requestData)
	if err != nil {
		return nil, err
	}

	coilData, err := ParseReadResponse(responseData)
	if err != nil {
		return nil, err
	}

	return UnpackBits(coilData, int(quantity)), nil
}

// ReadHoldingRegisters 讀取保持暫存器
func (c *ModbusClient) ReadHoldingRegisters(address uint16, quantity uint16) ([]uint16, error) {
	if quantity < 1 || quantity > HoldingRegisterMaxQuantity {
		return nil, ErrInvalidQuantity
	}

	requestData := BuildReadRequest(address, quantity)
	responseData, err := c.sendRequest(FuncReadHoldingRegisters, requestData)
	if err != nil {
		return nil, err
	}

	registerData, err := ParseReadResponse(responseData)
	if err != nil {
		return nil, err
	}

	if len(registerData) < int(quantity)*2 {
		return nil, ErrResponseTooShort
	}

	registers := make([]uint16, quantity)
	for i := uint16(0); i < quantity; i++ {
		registers[i] = binary.BigEndian.Uint16(registerData[i*2:])
	}

	return registers, nil
}

// ReadInputRegisters 讀取輸入暫存器
func (c *ModbusClient) ReadInputRegisters(address uint16, quantity uint16) ([]uint16, error) {
	if quantity < 1 || quantity > InputRegisterMaxQuantity {
		return nil, ErrInvalidQuantity
	}

	requestData := BuildReadRequest(address, quantity)
	responseData, err := c.sendRequest(FuncReadInputRegisters, requestData)
	if err != nil {
		return nil, err
	}

	registerData, err := ParseReadResponse(responseData)
	if err != nil {
		return nil, err
	}

	if len(registerData) < int(quantity)*2 {
		return nil, ErrResponseTooShort
	}

	registers := make([]uint16, quantity)
	for i := uint16(0); i < quantity; i++ {
		registers[i] = binary.BigEndian.Uint16(registerData[i*2:])
	}

	return registers, nil
}

// WriteSingleCoil 寫入單個線圈
func (c *ModbusClient) WriteSingleCoil(address uint16, value bool) error {
	requestData := BuildWriteSingleCoilRequest(address, value)
	responseData, err := c.sendRequest(FuncWriteSingleCoil, requestData)
	if err != nil {
		return err
	}

	// 驗證回應 (回應應與請求相同)
	if len(responseData) < 4 {
		return ErrResponseTooShort
	}

	respAddress := binary.BigEndian.Uint16(responseData[0:])
	respValue := binary.BigEndian.Uint16(responseData[2:])

	if respAddress != address {
		return fmt.Errorf("回應地址不匹配: 預期 %d, 實際 %d", address, respAddress)
	}

	expectedValue := uint16(0x0000)
	if value {
		expectedValue = 0xFF00
	}
	if respValue != expectedValue {
		return fmt.Errorf("回應值不匹配: 預期 0x%04X, 實際 0x%04X", expectedValue, respValue)
	}

	return nil
}

// WriteSingleRegister 寫入單個暫存器
func (c *ModbusClient) WriteSingleRegister(address uint16, value uint16) error {
	requestData := BuildWriteSingleRegisterRequest(address, value)
	responseData, err := c.sendRequest(FuncWriteSingleRegister, requestData)
	if err != nil {
		return err
	}

	// 驗證回應
	if len(responseData) < 4 {
		return ErrResponseTooShort
	}

	respAddress := binary.BigEndian.Uint16(responseData[0:])
	respValue := binary.BigEndian.Uint16(responseData[2:])

	if respAddress != address {
		return fmt.Errorf("回應地址不匹配: 預期 %d, 實際 %d", address, respAddress)
	}

	if respValue != value {
		return fmt.Errorf("回應值不匹配: 預期 %d, 實際 %d", value, respValue)
	}

	return nil
}

// WriteMultipleCoils 寫入多個線圈
func (c *ModbusClient) WriteMultipleCoils(address uint16, values []bool) error {
	quantity := uint16(len(values))
	if quantity < 1 || quantity > CoilMaxQuantity {
		return ErrInvalidQuantity
	}

	requestData := BuildWriteMultipleCoilsRequest(address, values)
	responseData, err := c.sendRequest(FuncWriteMultipleCoils, requestData)
	if err != nil {
		return err
	}

	// 驗證回應: Address(2) + Quantity(2)
	if len(responseData) < 4 {
		return ErrResponseTooShort
	}

	respAddress := binary.BigEndian.Uint16(responseData[0:])
	respQuantity := binary.BigEndian.Uint16(responseData[2:])

	if respAddress != address {
		return fmt.Errorf("回應地址不匹配: 預期 %d, 實際 %d", address, respAddress)
	}

	if respQuantity != quantity {
		return fmt.Errorf("回應數量不匹配: 預期 %d, 實際 %d", quantity, respQuantity)
	}

	return nil
}

// WriteMultipleRegisters 寫入多個暫存器
func (c *ModbusClient) WriteMultipleRegisters(address uint16, values []uint16) error {
	quantity := uint16(len(values))
	if quantity < 1 || quantity > HoldingRegisterMaxQuantity {
		return ErrInvalidQuantity
	}

	requestData := BuildWriteMultipleRegistersRequest(address, values)
	responseData, err := c.sendRequest(FuncWriteMultipleRegisters, requestData)
	if err != nil {
		return err
	}

	// 驗證回應: Address(2) + Quantity(2)
	if len(responseData) < 4 {
		return ErrResponseTooShort
	}

	respAddress := binary.BigEndian.Uint16(responseData[0:])
	respQuantity := binary.BigEndian.Uint16(responseData[2:])

	if respAddress != address {
		return fmt.Errorf("回應地址不匹配: 預期 %d, 實際 %d", address, respAddress)
	}

	if respQuantity != quantity {
		return fmt.Errorf("回應數量不匹配: 預期 %d, 實際 %d", quantity, respQuantity)
	}

	return nil
}
