package modbus

import "encoding/binary"

// ReadWordsExact reads a complete FC03/04 vector for a planned request. Unlike
// the legacy scalar API, a short, oversized or miscounted payload is never used.
func (c *ModbusClient) ReadWordsExact(function byte, address, quantity uint16) ([]uint16, error) {
	if function != FuncReadHoldingRegisters && function != FuncReadInputRegisters {
		return nil, ErrInvalidFunctionCode
	}
	if err := validateExactReadBounds(address, quantity, HoldingRegisterMaxQuantity); err != nil {
		return nil, err
	}
	data, err := c.readExact(function, address, quantity, int(quantity)*2)
	if err != nil {
		return nil, err
	}
	words := make([]uint16, quantity)
	for i := range words {
		words[i] = binary.BigEndian.Uint16(data[2*i:])
	}
	return words, nil
}

// ReadBitsExact reads every FC01/02 bit in ascending wire address order. Padding
// bits must be zero, as required by Modbus Application Protocol §§6.1 and 6.2.
func (c *ModbusClient) ReadBitsExact(function byte, address, quantity uint16) ([]bool, error) {
	if function != FuncReadCoils && function != FuncReadDiscreteInputs {
		return nil, ErrInvalidFunctionCode
	}
	if err := validateExactReadBounds(address, quantity, CoilMaxQuantity); err != nil {
		return nil, err
	}
	data, err := c.readExact(function, address, quantity, (int(quantity)+7)/8)
	if err != nil {
		return nil, err
	}
	if remainder := quantity % 8; remainder != 0 && data[len(data)-1]>>remainder != 0 {
		return nil, ErrInvalidResponse
	}
	return UnpackBits(data, int(quantity)), nil
}
func validateExactReadBounds(address, quantity uint16, limit int) error {
	if quantity < 1 || int(quantity) > limit {
		return ErrInvalidQuantity
	}
	if int(address)+int(quantity) > 65536 {
		return ErrInvalidAddress
	}
	return nil
}
func (c *ModbusClient) readExact(function byte, address, quantity uint16, bytes int) ([]byte, error) {
	// The decorator retains the original context transport, transaction counter,
	// serializer and parser. It adds validation, not a send or retry owner.
	strict := NewClient(&exactReadTransport{Transport: c.transport}, c.unitID)
	data, err := strict.sendRequest(function, BuildReadRequest(address, quantity))
	if err != nil {
		return nil, err
	}
	if len(data) < 1+bytes {
		return nil, ErrResponseTooShort
	}
	if len(data) > 1+bytes {
		return nil, ErrResponseTooLong
	}
	if int(data[0]) != bytes {
		return nil, ErrInvalidResponse
	}
	return data[1:], nil
}
