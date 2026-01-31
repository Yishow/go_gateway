package memory

import (
	"encoding/binary"
	"errors"
	"sync"
)

// =============================================================================
// 類型定義
// =============================================================================

// ByteOrder 字節序
type ByteOrder int

const (
	// BigEndian 大端序 (Modbus 標準)
	BigEndian ByteOrder = iota
	// LittleEndian 小端序
	LittleEndian
)

// 錯誤定義
var (
	ErrOutOfBounds = errors.New("記憶體存取越界")
)

// =============================================================================
// MemoryBank 虛擬記憶體庫
// =============================================================================

// MemoryBank 虛擬記憶體庫，模擬 PLC 記憶體區域
type MemoryBank struct {
	mu        sync.RWMutex
	data      []byte
	size      int
	byteOrder ByteOrder
}

// NewMemoryBank 建立新的記憶體庫 (預設 BigEndian)
func NewMemoryBank(size int) *MemoryBank {
	return NewMemoryBankWithOrder(size, BigEndian)
}

// NewMemoryBankWithOrder 建立指定字節序的記憶體庫
func NewMemoryBankWithOrder(size int, order ByteOrder) *MemoryBank {
	return &MemoryBank{
		data:      make([]byte, size),
		size:      size,
		byteOrder: order,
	}
}

// =============================================================================
// 基礎 Byte 操作
// =============================================================================

// ReadByte 讀取單個字節
func (m *MemoryBank) ReadByte(offset int) (byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if offset < 0 || offset >= m.size {
		return 0, ErrOutOfBounds
	}
	return m.data[offset], nil
}

// WriteByte 寫入單個字節
func (m *MemoryBank) WriteByte(offset int, value byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if offset < 0 || offset >= m.size {
		return ErrOutOfBounds
	}
	m.data[offset] = value
	return nil
}

// =============================================================================
// Word (16-bit) 操作
// =============================================================================

// ReadWord 讀取 16 位無符號整數
func (m *MemoryBank) ReadWord(offset int) (uint16, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if offset < 0 || offset+2 > m.size {
		return 0, ErrOutOfBounds
	}

	if m.byteOrder == BigEndian {
		return binary.BigEndian.Uint16(m.data[offset:]), nil
	}
	return binary.LittleEndian.Uint16(m.data[offset:]), nil
}

// WriteWord 寫入 16 位無符號整數
func (m *MemoryBank) WriteWord(offset int, value uint16) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if offset < 0 || offset+2 > m.size {
		return ErrOutOfBounds
	}

	if m.byteOrder == BigEndian {
		binary.BigEndian.PutUint16(m.data[offset:], value)
	} else {
		binary.LittleEndian.PutUint16(m.data[offset:], value)
	}
	return nil
}

// =============================================================================
// DWord (32-bit) 操作
// =============================================================================

// ReadDWord 讀取 32 位無符號整數
func (m *MemoryBank) ReadDWord(offset int) (uint32, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if offset < 0 || offset+4 > m.size {
		return 0, ErrOutOfBounds
	}

	if m.byteOrder == BigEndian {
		return binary.BigEndian.Uint32(m.data[offset:]), nil
	}
	return binary.LittleEndian.Uint32(m.data[offset:]), nil
}

// WriteDWord 寫入 32 位無符號整數
func (m *MemoryBank) WriteDWord(offset int, value uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if offset < 0 || offset+4 > m.size {
		return ErrOutOfBounds
	}

	if m.byteOrder == BigEndian {
		binary.BigEndian.PutUint32(m.data[offset:], value)
	} else {
		binary.LittleEndian.PutUint32(m.data[offset:], value)
	}
	return nil
}

// =============================================================================
// Slice 操作
// =============================================================================

// ReadSlice 讀取連續字節
func (m *MemoryBank) ReadSlice(offset, length int) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if offset < 0 || offset+length > m.size {
		return nil, ErrOutOfBounds
	}

	result := make([]byte, length)
	copy(result, m.data[offset:offset+length])
	return result, nil
}

// WriteSlice 寫入連續字節
func (m *MemoryBank) WriteSlice(offset int, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if offset < 0 || offset+len(data) > m.size {
		return ErrOutOfBounds
	}

	copy(m.data[offset:], data)
	return nil
}

// =============================================================================
// 輔助方法
// =============================================================================

// Dump 返回記憶體的完整快照 (複製)
func (m *MemoryBank) Dump() []byte {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]byte, m.size)
	copy(result, m.data)
	return result
}

// Size 返回記憶體大小
func (m *MemoryBank) Size() int {
	return m.size
}

// Clear 清空記憶體
func (m *MemoryBank) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.data {
		m.data[i] = 0
	}
}

// SetByteOrder 設定字節序
func (m *MemoryBank) SetByteOrder(order ByteOrder) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.byteOrder = order
}
