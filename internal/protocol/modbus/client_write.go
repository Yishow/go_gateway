package modbus

import (
	"encoding/binary"
	"fmt"
)

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
