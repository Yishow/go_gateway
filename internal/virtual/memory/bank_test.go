package memory

import (
	"testing"
)

// =============================================================================
// 任務 1.1: MemoryBank 行為測試
// =============================================================================

func TestMemoryBank_ByteReadWrite(t *testing.T) {
	bank := NewMemoryBank(1024)

	// 寫入 0xFF 到偏移量 0
	err := bank.WriteByte(0, 0xFF)
	if err != nil {
		t.Fatalf("WriteByte 失敗: %v", err)
	}

	// 讀取並驗證
	val, err := bank.ReadByte(0)
	if err != nil {
		t.Fatalf("ReadByte 失敗: %v", err)
	}
	if val != 0xFF {
		t.Errorf("預期 0xFF，實際 0x%X", val)
	}
}

func TestMemoryBank_WordReadWrite_BigEndian(t *testing.T) {
	bank := NewMemoryBankWithOrder(1024, BigEndian)

	// 寫入 0x1234
	err := bank.WriteWord(0, 0x1234)
	if err != nil {
		t.Fatalf("WriteWord 失敗: %v", err)
	}

	// BigEndian: 高字節在低地址
	high, _ := bank.ReadByte(0)
	low, _ := bank.ReadByte(1)

	if high != 0x12 {
		t.Errorf("高字節預期 0x12，實際 0x%X", high)
	}
	if low != 0x34 {
		t.Errorf("低字節預期 0x34，實際 0x%X", low)
	}

	// 驗證 ReadWord
	word, _ := bank.ReadWord(0)
	if word != 0x1234 {
		t.Errorf("ReadWord 預期 0x1234，實際 0x%X", word)
	}
}

func TestMemoryBank_WordReadWrite_LittleEndian(t *testing.T) {
	bank := NewMemoryBankWithOrder(1024, LittleEndian)

	// 寫入 0x1234
	err := bank.WriteWord(0, 0x1234)
	if err != nil {
		t.Fatalf("WriteWord 失敗: %v", err)
	}

	// LittleEndian: 低字節在低地址
	low, _ := bank.ReadByte(0)
	high, _ := bank.ReadByte(1)

	if low != 0x34 {
		t.Errorf("低字節預期 0x34，實際 0x%X", low)
	}
	if high != 0x12 {
		t.Errorf("高字節預期 0x12，實際 0x%X", high)
	}
}

func TestMemoryBank_DWord(t *testing.T) {
	bank := NewMemoryBankWithOrder(1024, BigEndian)

	// 寫入 32 位值
	err := bank.WriteDWord(0, 0x12345678)
	if err != nil {
		t.Fatalf("WriteDWord 失敗: %v", err)
	}

	val, _ := bank.ReadDWord(0)
	if val != 0x12345678 {
		t.Errorf("預期 0x12345678，實際 0x%X", val)
	}
}

func TestMemoryBank_OutOfBounds(t *testing.T) {
	bank := NewMemoryBank(100)

	// 超出範圍的寫入
	err := bank.WriteByte(100, 0xFF)
	if err == nil {
		t.Error("預期越界錯誤，但沒有返回錯誤")
	}

	// 超出範圍的讀取
	_, err = bank.ReadByte(100)
	if err == nil {
		t.Error("預期越界錯誤，但沒有返回錯誤")
	}

	// Word 操作越界
	err = bank.WriteWord(99, 0x1234)
	if err == nil {
		t.Error("預期 Word 越界錯誤，但沒有返回錯誤")
	}
}

func TestMemoryBank_ReadSlice(t *testing.T) {
	bank := NewMemoryBank(1024)

	// 寫入連續數據
	for i := 0; i < 10; i++ {
		bank.WriteByte(i, byte(i))
	}

	// 讀取 slice
	data, err := bank.ReadSlice(0, 10)
	if err != nil {
		t.Fatalf("ReadSlice 失敗: %v", err)
	}

	for i := 0; i < 10; i++ {
		if data[i] != byte(i) {
			t.Errorf("索引 %d: 預期 %d，實際 %d", i, i, data[i])
		}
	}
}

func TestMemoryBank_WriteSlice(t *testing.T) {
	bank := NewMemoryBank(1024)

	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	err := bank.WriteSlice(10, data)
	if err != nil {
		t.Fatalf("WriteSlice 失敗: %v", err)
	}

	// 驗證
	for i, expected := range data {
		val, _ := bank.ReadByte(10 + i)
		if val != expected {
			t.Errorf("索引 %d: 預期 0x%X，實際 0x%X", i, expected, val)
		}
	}
}

func TestMemoryBank_Dump(t *testing.T) {
	bank := NewMemoryBank(100)

	bank.WriteByte(0, 0xAB)
	bank.WriteByte(99, 0xCD)

	dump := bank.Dump()

	if len(dump) != 100 {
		t.Errorf("Dump 長度預期 100，實際 %d", len(dump))
	}
	if dump[0] != 0xAB {
		t.Errorf("索引 0 預期 0xAB，實際 0x%X", dump[0])
	}
	if dump[99] != 0xCD {
		t.Errorf("索引 99 預期 0xCD，實際 0x%X", dump[99])
	}
}
