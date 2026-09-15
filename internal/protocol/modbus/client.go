package modbus

import (
	"encoding/binary"
	"fmt"
	"sync"
)

// ModbusClient Modbus 客戶端
type ModbusClient struct { //nolint:revive // Preserve the exported Go name and its existing callers during lint maintenance.
	transport Transport
	unitID    byte
	mu        sync.Mutex
}

type transportUnwrapper interface {
	GetOriginalTransport() interface{}
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

// WithUnitID returns a client that reuses the same transport with a different unit ID.
func (c *ModbusClient) WithUnitID(unitID byte) *ModbusClient {
	return NewClient(c.transport, unitID)
}

// Connect 建立連線
func (c *ModbusClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.transport.Connect()
}

// Close 關閉連線
func (c *ModbusClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.transport.Close()
}

// IsConnected 檢查底層傳輸是否已連線
func (c *ModbusClient) IsConnected() bool {
	if c == nil || c.transport == nil {
		return false
	}
	if checker, ok := c.transport.(interface{ IsConnected() bool }); ok {
		return checker.IsConnected()
	}
	return true
}

// sendRequest 發送請求並接收回應 (TCP/UDP)
func (c *ModbusClient) sendTCPRequest(functionCode byte, data []byte) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var transactionID uint16
	// 嘗試通過接口調用（支持 WrappedTransport）
	if transportWithID, ok := c.transport.(interface{ GetNextTransactionID() uint16 }); ok {
		transactionID = transportWithID.GetNextTransactionID()
	} else {
		// 向後兼容：嘗試類型斷言
		switch transport := c.transport.(type) {
		case *TCPTransport:
			transactionID = transport.GetNextTransactionID()
		case *UDPTransport:
			transactionID = transport.GetNextTransactionID()
		default:
			return nil, fmt.Errorf("不支援的傳輸類型")
		}
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
	if transportWithOriginal, ok := c.transport.(transportUnwrapper); ok {
		switch transportWithOriginal.GetOriginalTransport().(type) {
		case *RTUTransport:
			return c.sendRTURequest(functionCode, data)
		case *TCPTransport, *UDPTransport:
			return c.sendTCPRequest(functionCode, data)
		}
	}

	// 檢查是否支持 GetNextTransactionID（TCP/UDP）或直接是 RTU
	_, hasTransactionID := c.transport.(interface{ GetNextTransactionID() uint16 })
	_, isRTU := c.transport.(*RTUTransport)

	switch {
	case hasTransactionID:
		return c.sendTCPRequest(functionCode, data)
	case isRTU:
		return c.sendRTURequest(functionCode, data)
	default:
		// 嘗試通過類型斷言判斷（向後兼容）
		switch c.transport.(type) {
		case *TCPTransport, *UDPTransport:
			return c.sendTCPRequest(functionCode, data)
		case *RTUTransport:
			return c.sendRTURequest(functionCode, data)
		default:
			return nil, fmt.Errorf("不支援的傳輸類型")
		}
	}
}

// ReadCoils 讀取線圈狀態
// address: 起始地址 (0-65535)
// quantity: 讀取數量 (1-2000)
func (c *ModbusClient) ReadCoils(address, quantity uint16) ([]bool, error) {
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
func (c *ModbusClient) ReadDiscreteInputs(address, quantity uint16) ([]bool, error) {
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
func (c *ModbusClient) ReadHoldingRegisters(address, quantity uint16) ([]uint16, error) {
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
func (c *ModbusClient) ReadInputRegisters(address, quantity uint16) ([]uint16, error) {
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

	// The response should contain the same address and value as the request.
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
func (c *ModbusClient) WriteSingleRegister(address, value uint16) error {
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
	quantity := len(values)
	if quantity < 1 || quantity > int(CoilMaxWriteQuantity) {
		return ErrInvalidQuantity
	}
	expectedQuantity := uint16(quantity)

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

	if respQuantity != expectedQuantity {
		return fmt.Errorf("回應數量不匹配: 預期 %d, 實際 %d", expectedQuantity, respQuantity)
	}

	return nil
}

// WriteMultipleRegisters 寫入多個暫存器
func (c *ModbusClient) WriteMultipleRegisters(address uint16, values []uint16) error {
	quantity := len(values)
	if quantity < 1 || quantity > int(HoldingRegisterMaxWriteQuantity) {
		return ErrInvalidQuantity
	}
	expectedQuantity := uint16(quantity)

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

	if respQuantity != expectedQuantity {
		return fmt.Errorf("回應數量不匹配: 預期 %d, 實際 %d", expectedQuantity, respQuantity)
	}

	return nil
}
