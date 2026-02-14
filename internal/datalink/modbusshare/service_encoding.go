package modbusshare

import (
	"fmt"
	"math"

	"go-gateway/internal/datalink/schema"
)

func encodeToWords(dataType schema.DataType, value interface{}) ([]uint16, error) {
	switch dataType {
	case schema.DataTypeBool:
		b, err := toBool(value)
		if err != nil {
			return nil, err
		}
		if b {
			return []uint16{1}, nil
		}
		return []uint16{0}, nil

	case schema.DataTypeInt16:
		v, err := toInt64(value)
		if err != nil {
			return nil, err
		}
		return []uint16{uint16(int16(v))}, nil

	case schema.DataTypeUint16:
		v, err := toUint64(value)
		if err != nil {
			return nil, err
		}
		return []uint16{uint16(v)}, nil

	case schema.DataTypeInt32:
		v, err := toInt64(value)
		if err != nil {
			return nil, err
		}
		raw := uint32(int32(v))
		return []uint16{uint16(raw >> 16), uint16(raw)}, nil

	case schema.DataTypeUint32:
		v, err := toUint64(value)
		if err != nil {
			return nil, err
		}
		raw := uint32(v)
		return []uint16{uint16(raw >> 16), uint16(raw)}, nil

	case schema.DataTypeFloat32:
		f, err := toFloat64(value)
		if err != nil {
			return nil, err
		}
		raw := math.Float32bits(float32(f))
		return []uint16{uint16(raw >> 16), uint16(raw)}, nil

	case schema.DataTypeInt64:
		v, err := toInt64(value)
		if err != nil {
			return nil, err
		}
		raw := uint64(v)
		return []uint16{
			uint16(raw >> 48),
			uint16(raw >> 32),
			uint16(raw >> 16),
			uint16(raw),
		}, nil

	case schema.DataTypeUint64:
		raw, err := toUint64(value)
		if err != nil {
			return nil, err
		}
		return []uint16{
			uint16(raw >> 48),
			uint16(raw >> 32),
			uint16(raw >> 16),
			uint16(raw),
		}, nil

	case schema.DataTypeFloat64:
		f, err := toFloat64(value)
		if err != nil {
			return nil, err
		}
		raw := math.Float64bits(f)
		return []uint16{
			uint16(raw >> 48),
			uint16(raw >> 32),
			uint16(raw >> 16),
			uint16(raw),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported data_type for Modbus mirror: %s", dataType)
	}
}

func toInt64(v interface{}) (int64, error) {
	switch n := v.(type) {
	case int:
		return int64(n), nil
	case int8:
		return int64(n), nil
	case int16:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return n, nil
	case uint:
		return int64(n), nil
	case uint8:
		return int64(n), nil
	case uint16:
		return int64(n), nil
	case uint32:
		return int64(n), nil
	case uint64:
		return int64(n), nil
	case float32:
		return int64(n), nil
	case float64:
		return int64(n), nil
	default:
		return 0, fmt.Errorf("value is not numeric")
	}
}

func toUint64(v interface{}) (uint64, error) {
	switch n := v.(type) {
	case int:
		return uint64(n), nil
	case int8:
		return uint64(n), nil
	case int16:
		return uint64(n), nil
	case int32:
		return uint64(n), nil
	case int64:
		return uint64(n), nil
	case uint:
		return uint64(n), nil
	case uint8:
		return uint64(n), nil
	case uint16:
		return uint64(n), nil
	case uint32:
		return uint64(n), nil
	case uint64:
		return n, nil
	case float32:
		return uint64(n), nil
	case float64:
		return uint64(n), nil
	default:
		return 0, fmt.Errorf("value is not numeric")
	}
}

func toFloat64(v interface{}) (float64, error) {
	switch n := v.(type) {
	case int:
		return float64(n), nil
	case int8:
		return float64(n), nil
	case int16:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case uint:
		return float64(n), nil
	case uint8:
		return float64(n), nil
	case uint16:
		return float64(n), nil
	case uint32:
		return float64(n), nil
	case uint64:
		return float64(n), nil
	case float32:
		return float64(n), nil
	case float64:
		return n, nil
	default:
		return 0, fmt.Errorf("value is not numeric")
	}
}

func toBool(v interface{}) (bool, error) {
	switch b := v.(type) {
	case bool:
		return b, nil
	case int:
		return b != 0, nil
	case int64:
		return b != 0, nil
	case float64:
		return b != 0, nil
	default:
		return false, fmt.Errorf("value is not boolean")
	}
}
