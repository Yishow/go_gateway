package handlers

import (
	"fmt"
	"sync"
)

// PacketRecorder 數據包記錄器接口
type PacketRecorder interface {
	RecordPacket(connectionID, protocol, direction string, data []byte)
}

// WrappedTransport 包裝的 Transport，用於記錄數據包
type WrappedTransport struct {
	transport     interface{} // 原始 Transport
	connectionID  string
	protocol      string
	packetRecorder PacketRecorder
	mu            sync.Mutex
}

// NewWrappedTransport 創建包裝的 Transport
func NewWrappedTransport(transport interface{}, connectionID, protocol string, recorder PacketRecorder) *WrappedTransport {
	return &WrappedTransport{
		transport:     transport,
		connectionID:  connectionID,
		protocol:      protocol,
		packetRecorder: recorder,
	}
}

// Connect 轉發到原始 Transport
func (w *WrappedTransport) Connect() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 使用類型斷言調用原始 Transport 的 Connect 方法
	switch t := w.transport.(type) {
	case interface{ Connect() error }:
		return t.Connect()
	default:
		return fmt.Errorf("transport does not support Connect")
	}
}

// Close 轉發到原始 Transport
func (w *WrappedTransport) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	switch t := w.transport.(type) {
	case interface{ Close() error }:
		return t.Close()
	default:
		return fmt.Errorf("transport does not support Close")
	}
}

// SendReceive 攔截並記錄數據包
func (w *WrappedTransport) SendReceive(data []byte) ([]byte, error) {
	w.mu.Lock()
	transport := w.transport
	connectionID := w.connectionID
	protocol := w.protocol
	recorder := w.packetRecorder
	w.mu.Unlock()

	// 記錄發送的數據包
	if recorder != nil {
		recorder.RecordPacket(connectionID, protocol, "request", data)
	}

	// 調用原始 Transport 的 SendReceive
	var response []byte
	var err error
	switch t := transport.(type) {
	case interface{ SendReceive([]byte) ([]byte, error) }:
		response, err = t.SendReceive(data)
	default:
		return nil, fmt.Errorf("transport does not support SendReceive")
	}

	// 記錄接收的數據包
	if recorder != nil && err == nil && len(response) > 0 {
		recorder.RecordPacket(connectionID, protocol, "response", response)
	}

	return response, err
}

// GetOriginalTransport 獲取原始 Transport（用於需要訪問原始 Transport 的情況）
func (w *WrappedTransport) GetOriginalTransport() interface{} {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.transport
}

// GetNextTransactionID 轉發給原始 Transport（用於 Modbus TCP/UDP）
func (w *WrappedTransport) GetNextTransactionID() uint16 {
	w.mu.Lock()
	transport := w.transport
	w.mu.Unlock()

	// 嘗試類型斷言並調用原始 Transport 的方法
	switch t := transport.(type) {
	case interface{ GetNextTransactionID() uint16 }:
		return t.GetNextTransactionID()
	default:
		// 如果原始 Transport 不支持，返回 0（這不應該發生）
		return 0
	}
}
