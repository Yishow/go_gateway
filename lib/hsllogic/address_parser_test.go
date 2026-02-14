package hsllogic

import (
	"testing"
)

// =============================================================================
// Modbus 地址解析測試
// =============================================================================

func TestParseModbusAddress_HoldingRegister(t *testing.T) {
	tests := []struct {
		name       string
		address    string
		wantDevice string
		wantOffset int
		wantBit    int
	}{
		{"帶前綴 HR 地址", "HR100", "HR", 100, -1},
		{"五位數 HR 地址", "40100", "HR", 100, -1},
		{"HR 地址帶位元索引", "HR100.5", "HR", 100, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAddress(ProtocolModbus, tt.address)
			if err != nil {
				t.Fatalf("ParseAddress() 錯誤 = %v", err)
			}
			if got.DeviceType != tt.wantDevice {
				t.Errorf("DeviceType = %v, 預期 %v", got.DeviceType, tt.wantDevice)
			}
			if got.Offset != tt.wantOffset {
				t.Errorf("Offset = %v, 預期 %v", got.Offset, tt.wantOffset)
			}
			if got.BitIndex != tt.wantBit {
				t.Errorf("BitIndex = %v, 預期 %v", got.BitIndex, tt.wantBit)
			}
		})
	}
}

func TestParseModbusAddress_OtherRegisters(t *testing.T) {
	tests := []struct {
		name       string
		address    string
		wantDevice string
		wantOffset int
		wantIsBit  bool
	}{
		{"Coil 地址", "COIL100", "COIL", 100, true},
		{"五位數 Coil", "00100", "COIL", 100, true},
		{"Discrete Input", "DI50", "DI", 50, true},
		{"五位數 DI", "10050", "DI", 50, true},
		{"Input Register", "IR200", "IR", 200, false},
		{"五位數 IR", "30200", "IR", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAddress(ProtocolModbus, tt.address)
			if err != nil {
				t.Fatalf("ParseAddress() 錯誤 = %v", err)
			}
			if got.DeviceType != tt.wantDevice {
				t.Errorf("DeviceType = %v, 預期 %v", got.DeviceType, tt.wantDevice)
			}
			if got.Offset != tt.wantOffset {
				t.Errorf("Offset = %v, 預期 %v", got.Offset, tt.wantOffset)
			}
			if got.IsBitDevice != tt.wantIsBit {
				t.Errorf("IsBitDevice = %v, 預期 %v", got.IsBitDevice, tt.wantIsBit)
			}
		})
	}
}

// =============================================================================
// 三菱 MC Protocol 地址解析測試
// =============================================================================

func TestParseMitsubishiAddress(t *testing.T) {
	tests := []struct {
		name       string
		address    string
		wantDevice string
		wantOffset int
		wantIsBit  bool
	}{
		// 十進位設備
		{"D 暫存器", "D100", "D", 100, false},
		{"M 線圈", "M0", "M", 0, true},
		{"R 檔案暫存器", "R500", "R", 500, false},
		{"ZR 擴展暫存器", "ZR1000", "ZR", 1000, false},
		{"SM 特殊繼電器", "SM400", "SM", 400, true},

		// 八進位設備 (X, Y)
		{"X 輸入 (八進位)", "X17", "X", 15, true}, // 八進位 17 = 十進位 15
		{"Y 輸出 (八進位)", "Y10", "Y", 8, true},  // 八進位 10 = 十進位 8

		// 十六進位設備 (W, B, SW, SB)
		{"W 鏈路暫存器 (十六進位)", "W10", "W", 16, false},   // 十六進位 10 = 十進位 16
		{"W 鏈路暫存器 (大數)", "W1000", "W", 4096, false}, // 十六進位 1000 = 十進位 4096
		{"B 鏈路繼電器 (十六進位)", "B20", "B", 32, true},    // 十六進位 20 = 十進位 32
		{"B 鏈路繼電器 (英文)", "BFF", "B", 255, true},     // 十六進位 FF = 十進位 255
		{"SW 特殊鏈路暫存器", "SW100", "SW", 256, false},   // 十六進位 100 = 十進位 256
		{"SB 特殊鏈路繼電器", "SB1A", "SB", 26, true},      // 十六進位 1A = 十進位 26
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAddress(ProtocolMitsubishi, tt.address)
			if err != nil {
				t.Fatalf("ParseAddress() 錯誤 = %v", err)
			}
			if got.DeviceType != tt.wantDevice {
				t.Errorf("DeviceType = %v, 預期 %v", got.DeviceType, tt.wantDevice)
			}
			if got.Offset != tt.wantOffset {
				t.Errorf("Offset = %v, 預期 %v", got.Offset, tt.wantOffset)
			}
			if got.IsBitDevice != tt.wantIsBit {
				t.Errorf("IsBitDevice = %v, 預期 %v", got.IsBitDevice, tt.wantIsBit)
			}
		})
	}
}

// =============================================================================
// FATEK 地址解析測試
// =============================================================================

func TestParseFatekAddress(t *testing.T) {
	tests := []struct {
		name       string
		address    string
		wantDevice string
		wantOffset int
		wantIsBit  bool
	}{
		{"D 暫存器", "D100", "D", 100, false},
		{"R 暫存器", "R0", "R", 0, false},
		{"M 線圈", "M100", "M", 100, true},
		{"X 輸入", "X0", "X", 0, true},
		{"Y 輸出", "Y0", "Y", 0, true},
		{"DD 雙字組", "DD100", "DD", 100, false},
		{"WM 字組", "WM0", "WM", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAddress(ProtocolFatek, tt.address)
			if err != nil {
				t.Fatalf("ParseAddress() 錯誤 = %v", err)
			}
			if got.DeviceType != tt.wantDevice {
				t.Errorf("DeviceType = %v, 預期 %v", got.DeviceType, tt.wantDevice)
			}
			if got.Offset != tt.wantOffset {
				t.Errorf("Offset = %v, 預期 %v", got.Offset, tt.wantOffset)
			}
			if got.IsBitDevice != tt.wantIsBit {
				t.Errorf("IsBitDevice = %v, 預期 %v", got.IsBitDevice, tt.wantIsBit)
			}
		})
	}
}

// =============================================================================
// Siemens S7 地址解析測試
// =============================================================================

func TestParseSiemensAddress(t *testing.T) {
	tests := []struct {
		name       string
		address    string
		wantDevice string
		wantDB     int
		wantOffset int
		wantBit    int
		wantIsBit  bool
	}{
		{"DB 位元", "DB1.DBX0.0", "DBX", 1, 0, 0, true},
		{"DB 字組", "DB1.DBW10", "DBW", 1, 10, -1, false},
		{"DB 雙字組", "DB10.DBD100", "DBD", 10, 100, -1, false},
		{"M 位元", "M0.0", "M", 0, 0, 0, true},
		{"MW 字組", "MW100", "MW", 0, 100, -1, false},
		{"I 輸入", "I0.0", "I", 0, 0, 0, true},
		{"IW 輸入字組", "IW0", "IW", 0, 0, -1, false},
		{"Q 輸出", "Q0.0", "Q", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAddress(ProtocolSiemens, tt.address)
			if err != nil {
				t.Fatalf("ParseAddress() 錯誤 = %v", err)
			}
			if got.DeviceType != tt.wantDevice {
				t.Errorf("DeviceType = %v, 預期 %v", got.DeviceType, tt.wantDevice)
			}
			if got.DBNumber != tt.wantDB {
				t.Errorf("DBNumber = %v, 預期 %v", got.DBNumber, tt.wantDB)
			}
			if got.Offset != tt.wantOffset {
				t.Errorf("Offset = %v, 預期 %v", got.Offset, tt.wantOffset)
			}
			if got.BitIndex != tt.wantBit {
				t.Errorf("BitIndex = %v, 預期 %v", got.BitIndex, tt.wantBit)
			}
			if got.IsBitDevice != tt.wantIsBit {
				t.Errorf("IsBitDevice = %v, 預期 %v", got.IsBitDevice, tt.wantIsBit)
			}
		})
	}
}

// =============================================================================
// 錯誤處理測試
// =============================================================================

func TestParseAddress_Errors(t *testing.T) {
	tests := []struct {
		name     string
		protocol ProtocolType
		address  string
	}{
		{"空地址", ProtocolModbus, ""},
		{"無效 Modbus 地址", ProtocolModbus, "ABC"},
		{"無效三菱地址", ProtocolMitsubishi, "Z"},
		{"無效三菱地址數字", ProtocolMitsubishi, "DXYZ"},
		{"無效 FATEK 地址", ProtocolFatek, "Z100"},
		{"無效 Siemens 地址", ProtocolSiemens, "INVALID"},
		{"不支援的協議", "unknown", "D100"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseAddress(tt.protocol, tt.address)
			if err == nil {
				t.Error("預期應該返回錯誤，但沒有")
			}
		})
	}
}

// =============================================================================
// ParsedAddress.String() 測試
// =============================================================================

func TestParsedAddress_String(t *testing.T) {
	tests := []struct {
		name string
		addr ParsedAddress
		want string
	}{
		{
			name: "簡單地址",
			addr: ParsedAddress{DeviceType: "D", Offset: 100, BitIndex: -1},
			want: "D100",
		},
		{
			name: "帶位元索引",
			addr: ParsedAddress{DeviceType: "M", Offset: 10, BitIndex: 5},
			want: "M10.5",
		},
		{
			name: "帶 DB 編號",
			addr: ParsedAddress{DeviceType: "DBW", DBNumber: 1, Offset: 10, BitIndex: -1},
			want: "DB1.DBW10",
		},
		{
			name: "帶 DB 編號和位元",
			addr: ParsedAddress{DeviceType: "DBX", DBNumber: 1, Offset: 0, BitIndex: 0},
			want: "DB1.DBX0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.addr.String()
			if got != tt.want {
				t.Errorf("String() = %v, 預期 %v", got, tt.want)
			}
		})
	}
}

func TestParsedAddress_ToByteOffset(t *testing.T) {
	siemens := &ParsedAddress{Protocol: ProtocolSiemens, Offset: 10}
	if siemens.ToByteOffset() != 10 {
		t.Fatalf("Siemens ToByteOffset 預期 10，實際 %d", siemens.ToByteOffset())
	}

	modbus := &ParsedAddress{Protocol: ProtocolModbus, Offset: 10}
	if modbus.ToByteOffset() != 20 {
		t.Fatalf("Modbus ToByteOffset 預期 20，實際 %d", modbus.ToByteOffset())
	}
}
