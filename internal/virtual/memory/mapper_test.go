package memory

import (
	"testing"
)

// =============================================================================
// 任務 1.3: AddressMapper 測試
// =============================================================================

func TestAddressMapper_ModbusHoldingRegister(t *testing.T) {
	mapper := NewAddressMapper()

	// Modbus Holding Register: 40001 -> Offset 0
	offset, err := mapper.Map("40001")
	if err != nil {
		t.Fatalf("Map 失敗: %v", err)
	}
	if offset != 0 {
		t.Errorf("40001 預期偏移 0，實際 %d", offset)
	}

	// 40002 -> Offset 2 (每暫存器 2 bytes)
	offset, err = mapper.Map("40002")
	if err != nil {
		t.Fatalf("Map 失敗: %v", err)
	}
	if offset != 2 {
		t.Errorf("40002 預期偏移 2，實際 %d", offset)
	}
}

func TestAddressMapper_ModbusInputRegister(t *testing.T) {
	mapper := NewAddressMapper()

	// Input Register: 30001 -> Offset 0 (在 Input 區)
	offset, err := mapper.Map("30001")
	if err != nil {
		t.Fatalf("Map 失敗: %v", err)
	}
	if offset != 0 {
		t.Errorf("30001 預期偏移 0，實際 %d", offset)
	}
}

func TestAddressMapper_ModbusCoil(t *testing.T) {
	mapper := NewAddressMapper()

	// Coil: 00001 -> Offset 0 (bit 地址)
	offset, err := mapper.Map("00001")
	if err != nil {
		t.Fatalf("Map 失敗: %v", err)
	}
	if offset != 0 {
		t.Errorf("00001 預期偏移 0，實際 %d", offset)
	}
}

func TestAddressMapper_PLCStyle_D(t *testing.T) {
	mapper := NewAddressMapper()

	// D100 -> Offset 200 (D 區基址 0 + 100 * 2 bytes)
	offset, err := mapper.Map("D100")
	if err != nil {
		t.Fatalf("Map 失敗: %v", err)
	}
	if offset != 200 {
		t.Errorf("D100 預期偏移 200，實際 %d", offset)
	}
}

func TestAddressMapper_PLCStyle_M(t *testing.T) {
	mapper := NewAddressMapper()

	// M0 -> Offset 0 (bit 地址，每 8 bits = 1 byte)
	offset, err := mapper.Map("M0")
	if err != nil {
		t.Fatalf("Map 失敗: %v", err)
	}
	if offset != 0 {
		t.Errorf("M0 預期偏移 0，實際 %d", offset)
	}

	// M8 -> Offset 1
	offset, err = mapper.Map("M8")
	if err != nil {
		t.Fatalf("Map 失敗: %v", err)
	}
	if offset != 1 {
		t.Errorf("M8 預期偏移 1，實際 %d", offset)
	}
}

func TestAddressMapper_PLCStyle_R(t *testing.T) {
	mapper := NewAddressMapper()

	// R0 -> Offset 0 (16-bit 暫存器)
	offset, err := mapper.Map("R0")
	if err != nil {
		t.Fatalf("Map 失敗: %v", err)
	}
	if offset != 0 {
		t.Errorf("R0 預期偏移 0，實際 %d", offset)
	}
}

func TestAddressMapper_InvalidAddress(t *testing.T) {
	mapper := NewAddressMapper()

	// 無效地址
	_, err := mapper.Map("INVALID")
	if err == nil {
		t.Error("預期無效地址錯誤，但沒有返回錯誤")
	}
}

func TestAddressMapper_MemoryArea(t *testing.T) {
	mapper := NewAddressMapper()

	// 獲取記憶體區域資訊
	area, offset, err := mapper.MapWithArea("40001")
	if err != nil {
		t.Fatalf("MapWithArea 失敗: %v", err)
	}
	if area != AreaHoldingRegister {
		t.Errorf("預期區域 HoldingRegister，實際 %v", area)
	}
	if offset != 0 {
		t.Errorf("預期偏移 0，實際 %d", offset)
	}

	// D 區
	area, offset, err = mapper.MapWithArea("D100")
	if err != nil {
		t.Fatalf("MapWithArea 失敗: %v", err)
	}
	if area != AreaDataRegister {
		t.Errorf("預期區域 DataRegister，實際 %v", area)
	}
}
