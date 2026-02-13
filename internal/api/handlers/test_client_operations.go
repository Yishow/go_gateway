package handlers

import (
"fmt"

"go-gateway/internal/protocol/fatek"
"go-gateway/internal/protocol/mcprotocol"
"go-gateway/internal/protocol/modbus"
)

func (h *TestHandler) connectClient(client interface{}, protocol string) error {
	switch c := client.(type) {
	case *modbus.ModbusClient:
		return c.Connect()
	case *fatek.FatekClient:
		return c.Connect()
	case *mcprotocol.MCClient:
		return c.Connect()
	default:
		return fmt.Errorf("unknown client type")
	}
}

func (h *TestHandler) closeClient(client interface{}, protocol string) error {
	switch c := client.(type) {
	case *modbus.ModbusClient:
		return c.Close()
	case *fatek.FatekClient:
		return c.Close()
	case *mcprotocol.MCClient:
		return c.Close()
	default:
		return fmt.Errorf("unknown client type")
	}
}

func (h *TestHandler) executeRead(client interface{}, protocol string, req ReadRequest) (interface{}, error) {
	switch c := client.(type) {
	case *modbus.ModbusClient:
		switch req.Operation {
		case "read_coils":
			return c.ReadCoils(req.Address, req.Count)
		case "read_discrete_inputs":
			return c.ReadDiscreteInputs(req.Address, req.Count)
		case "read_holding_registers":
			return c.ReadHoldingRegisters(req.Address, req.Count)
		case "read_input_registers":
			return c.ReadInputRegisters(req.Address, req.Count)
		default:
			return nil, fmt.Errorf("unsupported operation for modbus: %s", req.Operation)
		}

	case *fatek.FatekClient:
		switch req.Operation {
		case "read_status":
			return c.ReadStatus(req.Symbol, int(req.Address), int(req.Count))
		case "read_registers":
			return c.ReadRegisters(req.Symbol, int(req.Address), int(req.Count))
		default:
			return nil, fmt.Errorf("unsupported operation for fatek: %s", req.Operation)
		}

	case *mcprotocol.MCClient:
		switch req.Operation {
		case "batch_read_word":
			return c.BatchReadWord(req.Device, int(req.Address), int(req.Count))
		case "batch_read_bit":
			return c.BatchReadBit(req.Device, int(req.Address), int(req.Count))
		default:
			return nil, fmt.Errorf("unsupported operation for mcprotocol: %s", req.Operation)
		}

	default:
		return nil, fmt.Errorf("unknown client type")
	}
}

func (h *TestHandler) executeWrite(client interface{}, protocol string, req WriteRequest) error {
	// Helper to convert interface{} to []uint16 or []bool
	// This is tricky because JSON unmarshaling might give []interface{}
	
	toUint16Slice := func(v interface{}) ([]uint16, error) {
		arr, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("value is not an array")
		}
		res := make([]uint16, len(arr))
		for i, val := range arr {
			switch v := val.(type) {
			case float64:
				res[i] = uint16(v)
			case int:
				res[i] = uint16(v)
			case int64:
				res[i] = uint16(v)
			default:
				return nil, fmt.Errorf("invalid value type at index %d: expected number, got %T", i, val)
			}
		}
		return res, nil
	}
	
	toIntSlice := func(v interface{}) ([]int, error) {
		arr, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("value is not an array")
		}
		res := make([]int, len(arr))
		for i, val := range arr {
			switch v := val.(type) {
			case float64:
				res[i] = int(v)
			case int:
				res[i] = v
			case int64:
				res[i] = int(v)
			default:
				return nil, fmt.Errorf("invalid value type at index %d: expected number, got %T", i, val)
			}
		}
		return res, nil
	}

	toBoolSlice := func(v interface{}) ([]bool, error) {
		arr, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("value is not an array")
		}
		res := make([]bool, len(arr))
		for i, val := range arr {
			if b, ok := val.(bool); ok {
				res[i] = b
			} else {
				return nil, fmt.Errorf("invalid value type at index %d", i)
			}
		}
		return res, nil
	}

	switch c := client.(type) {
	case *modbus.ModbusClient:
		switch req.Operation {
		case "write_single_coil":
			// 支援數組或單個值：如果是數組，取第一個元素
			var val bool
			if arr, ok := req.Values.([]interface{}); ok && len(arr) > 0 {
				if b, ok := arr[0].(bool); ok {
					val = b
				} else {
					return fmt.Errorf("value must be bool, got %T", arr[0])
				}
			} else if b, ok := req.Values.(bool); ok {
				val = b
			} else {
				return fmt.Errorf("value must be bool or array of bool")
			}
			return c.WriteSingleCoil(req.Address, val)
		case "write_single_register":
			// 支援數組或單個值：如果是數組，取第一個元素
			var val uint16
			if arr, ok := req.Values.([]interface{}); ok && len(arr) > 0 {
				if f, ok := arr[0].(float64); ok {
					val = uint16(f)
				} else if i, ok := arr[0].(int); ok {
					val = uint16(i)
				} else {
					return fmt.Errorf("value must be number, got %T", arr[0])
				}
			} else if f, ok := req.Values.(float64); ok {
				val = uint16(f)
			} else if i, ok := req.Values.(int); ok {
				val = uint16(i)
			} else {
				return fmt.Errorf("value must be number or array of number")
			}
			return c.WriteSingleRegister(req.Address, val)
		case "write_multiple_coils":
			vals, err := toBoolSlice(req.Values)
			if err != nil { return err }
			return c.WriteMultipleCoils(req.Address, vals)
		case "write_multiple_registers":
			vals, err := toUint16Slice(req.Values)
			if err != nil { return err }
			return c.WriteMultipleRegisters(req.Address, vals)
		default:
			return fmt.Errorf("unsupported operation for modbus: %s", req.Operation)
		}

	case *fatek.FatekClient:
		switch req.Operation {
		case "write_status":
			vals, err := toBoolSlice(req.Values)
			if err != nil { return err }
			return c.WriteStatus(req.Symbol, int(req.Address), vals)
		case "write_registers":
			vals, err := toIntSlice(req.Values)
			if err != nil { return err }
			return c.WriteRegisters(req.Symbol, int(req.Address), vals)
		default:
			return fmt.Errorf("unsupported operation for fatek: %s", req.Operation)
		}

	case *mcprotocol.MCClient:
		switch req.Operation {
		case "batch_write_word":
			vals, err := toIntSlice(req.Values)
			if err != nil { return err }
			return c.BatchWriteWord(req.Device, int(req.Address), vals)
		case "batch_write_bit":
			vals, err := toBoolSlice(req.Values)
			if err != nil { return err }
			return c.BatchWriteBit(req.Device, int(req.Address), vals)
		default:
			return fmt.Errorf("unsupported operation for mcprotocol: %s", req.Operation)
		}

	default:
		return fmt.Errorf("unknown client type")
	}
}
