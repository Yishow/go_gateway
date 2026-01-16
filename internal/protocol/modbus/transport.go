package modbus

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"go.bug.st/serial"
)

// Transport 定義 Modbus 傳輸介面
type Transport interface {
	Connect() error
	Close() error
	SendReceive(data []byte) ([]byte, error)
}

// TCPTransport Modbus TCP 傳輸實作
type TCPTransport struct {
	Host          string
	Port          int
	Timeout       time.Duration
	conn          net.Conn
	mu            sync.Mutex
	transactionID uint16
}

// NewTCPTransport 建立新的 TCP 傳輸實例
func NewTCPTransport(host string, port int) *TCPTransport {
	if port == 0 {
		port = TCPDefaultPort
	}
	return &TCPTransport{
		Host:          host,
		Port:          port,
		Timeout:       2 * time.Second,
		transactionID: 0,
	}
}

func (t *TCPTransport) Connect() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn != nil {
		t.conn.Close()
	}

	// 使用 net.JoinHostPort 支援 IPv6
	addr := net.JoinHostPort(t.Host, strconv.Itoa(t.Port))
	conn, err := net.DialTimeout("tcp", addr, t.Timeout)
	if err != nil {
		return fmt.Errorf("TCP 連線失敗: %w", err)
	}

	t.conn = conn
	return nil
}

func (t *TCPTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn != nil {
		err := t.conn.Close()
		t.conn = nil
		return err
	}
	return nil
}

func (t *TCPTransport) SendReceive(data []byte) ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.conn == nil {
		return nil, ErrConnectionClosed
	}

	t.conn.SetDeadline(time.Now().Add(t.Timeout))

	// 寫入請求
	if _, err := t.conn.Write(data); err != nil {
		t.internalClose()
		return nil, fmt.Errorf("TCP 寫入失敗: %w", err)
	}

	// 讀取 MBAP 標頭 (7 bytes)
	header := make([]byte, MBAPHeaderLength)
	if _, err := io.ReadFull(t.conn, header); err != nil {
		t.internalClose()
		return nil, fmt.Errorf("TCP 讀取標頭失敗: %w", err)
	}

	mbap, err := ParseMBAPHeader(header)
	if err != nil {
		t.internalClose()
		return nil, err
	}

	// 讀取 PDU (Length - 1, 因為 Length 包含 Unit ID)
	pduLength := int(mbap.Length) - 1
	if pduLength < 0 || pduLength > 256 {
		t.internalClose()
		return nil, ErrInvalidFrame
	}

	pdu := make([]byte, pduLength)
	if _, err := io.ReadFull(t.conn, pdu); err != nil {
		t.internalClose()
		return nil, fmt.Errorf("TCP 讀取 PDU 失敗: %w", err)
	}

	// 返回完整回應 (MBAP + PDU)
	return append(header, pdu...), nil
}

func (t *TCPTransport) internalClose() {
	if t.conn != nil {
		t.conn.Close()
		t.conn = nil
	}
}

// GetNextTransactionID 獲取下一個交易 ID
func (t *TCPTransport) GetNextTransactionID() uint16 {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.transactionID++
	if t.transactionID == 0 {
		t.transactionID = 1
	}
	return t.transactionID
}

// UDPTransport Modbus UDP 傳輸實作
type UDPTransport struct {
	Host          string
	Port          int
	Timeout       time.Duration
	conn          *net.UDPConn
	mu            sync.Mutex
	transactionID uint16
}

// NewUDPTransport 建立新的 UDP 傳輸實例
func NewUDPTransport(host string, port int) *UDPTransport {
	if port == 0 {
		port = UDPDefaultPort
	}
	return &UDPTransport{
		Host:          host,
		Port:          port,
		Timeout:       2 * time.Second,
		transactionID: 0,
	}
}

func (u *UDPTransport) Connect() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.conn != nil {
		u.conn.Close()
	}

	// 使用 net.JoinHostPort 支援 IPv6
	addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(u.Host, strconv.Itoa(u.Port)))
	if err != nil {
		return fmt.Errorf("UDP 地址解析失敗: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("UDP 連線失敗: %w", err)
	}

	u.conn = conn
	return nil
}

func (u *UDPTransport) Close() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.conn != nil {
		err := u.conn.Close()
		u.conn = nil
		return err
	}
	return nil
}

func (u *UDPTransport) SendReceive(data []byte) ([]byte, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if u.conn == nil {
		return nil, ErrConnectionClosed
	}

	u.conn.SetDeadline(time.Now().Add(u.Timeout))

	// 寫入請求
	if _, err := u.conn.Write(data); err != nil {
		return nil, fmt.Errorf("UDP 寫入失敗: %w", err)
	}

	// 讀取回應
	buffer := make([]byte, 256)
	n, err := u.conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("UDP 讀取失敗: %w", err)
	}

	return buffer[:n], nil
}

func (u *UDPTransport) GetNextTransactionID() uint16 {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.transactionID++
	if u.transactionID == 0 {
		u.transactionID = 1
	}
	return u.transactionID
}

// RTUTransport Modbus RTU 傳輸實作 (串列埠)
type RTUTransport struct {
	Port     string
	BaudRate int
	DataBits int
	Parity   serial.Parity
	StopBits serial.StopBits
	Timeout  time.Duration
	port     serial.Port
	reader   *bufio.Reader
	mu       sync.Mutex
}

// NewRTUTransport 建立新的 RTU 傳輸實例
func NewRTUTransport(port string, baudRate, dataBits, stopBits int, parity string, timeout time.Duration) *RTUTransport {
	if baudRate == 0 {
		baudRate = DefaultBaudRate
	}
	if dataBits == 0 {
		dataBits = DefaultDataBits
	}
	if timeout == 0 {
		timeout = DefaultTimeout * time.Second
	}

	var p serial.Parity
	switch parity {
	case "N", "n":
		p = serial.NoParity
	case "E", "e":
		p = serial.EvenParity
	case "O", "o":
		p = serial.OddParity
	default:
		p = serial.NoParity
	}

	var sb serial.StopBits
	switch stopBits {
	case 1:
		sb = serial.OneStopBit
	case 2:
		sb = serial.TwoStopBits
	default:
		sb = serial.OneStopBit
	}

	return &RTUTransport{
		Port:     port,
		BaudRate: baudRate,
		DataBits: dataBits,
		Parity:   p,
		StopBits: sb,
		Timeout:  timeout,
	}
}

func (r *RTUTransport) Connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.port != nil {
		r.port.Close()
	}

	mode := &serial.Mode{
		BaudRate: r.BaudRate,
		DataBits: r.DataBits,
		Parity:   r.Parity,
		StopBits: r.StopBits,
	}

	port, err := serial.Open(r.Port, mode)
	if err != nil {
		return fmt.Errorf("串列埠開啟失敗: %w", err)
	}

	r.port = port
	r.reader = bufio.NewReader(port)
	return nil
}

func (r *RTUTransport) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.port != nil {
		err := r.port.Close()
		r.port = nil
		r.reader = nil
		return err
	}
	return nil
}

func (r *RTUTransport) SendReceive(data []byte) ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.port == nil {
		return nil, ErrConnectionClosed
	}

	r.port.SetReadTimeout(r.Timeout)

	// 寫入請求
	if _, err := r.port.Write(data); err != nil {
		r.internalClose()
		return nil, fmt.Errorf("RTU 寫入失敗: %w", err)
	}

	// RTU 需要等待一段時間 (根據波特率計算)
	// 通常為 3.5 個字符時間
	charTime := time.Duration(10000000/r.BaudRate) * time.Microsecond
	time.Sleep(charTime * 35 / 10) // 3.5 字符時間

	// 讀取回應
	// 先讀取地址和功能碼
	header := make([]byte, 2)
	if _, err := io.ReadFull(r.reader, header); err != nil {
		r.internalClose()
		return nil, fmt.Errorf("RTU 讀取標頭失敗: %w", err)
	}

	functionCode := header[1]

	// 檢查是否為異常回應
	if functionCode&FuncExceptionOffset != 0 {
		// 異常回應: Address(1) + Function(1) + Exception(1) + CRC(2)
		exceptionAndCRC := make([]byte, 3)
		if _, err := io.ReadFull(r.reader, exceptionAndCRC); err != nil {
			r.internalClose()
			return nil, fmt.Errorf("RTU 讀取異常回應失敗: %w", err)
		}
		return append(header, exceptionAndCRC...), nil
	}

	// 根據功能碼判斷後續數據長度
	var expectedLength int
	switch functionCode {
	case FuncReadCoils, FuncReadDiscreteInputs, FuncReadHoldingRegisters, FuncReadInputRegisters:
		// 讀取回應: Address(1) + Function(1) + ByteCount(1) + Data(N) + CRC(2)
		// 需要先讀取 ByteCount
		byteCount := make([]byte, 1)
		if _, err := io.ReadFull(r.reader, byteCount); err != nil {
			r.internalClose()
			return nil, fmt.Errorf("RTU 讀取 ByteCount 失敗: %w", err)
		}
		expectedLength = 2 + 1 + int(byteCount[0]) + 2 // Header + ByteCount + Data + CRC

		// 讀取數據和 CRC
		dataAndCRC := make([]byte, int(byteCount[0])+2)
		if _, err := io.ReadFull(r.reader, dataAndCRC); err != nil {
			r.internalClose()
			return nil, fmt.Errorf("RTU 讀取數據失敗: %w", err)
		}
		return append(header, append(byteCount, dataAndCRC...)...), nil

	case FuncWriteSingleCoil, FuncWriteSingleRegister:
		// 寫入單個回應: Address(1) + Function(1) + Address(2) + Value(2) + CRC(2)
		expectedLength = 2 + 4 + 2 // Header + Address(2) + Value(2) + CRC

	case FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		// 寫入多個回應: Address(1) + Function(1) + Address(2) + Quantity(2) + CRC(2)
		expectedLength = 2 + 4 + 2 // Header + Address(2) + Quantity(2) + CRC

	default:
		r.internalClose()
		return nil, ErrInvalidFunctionCode
	}

	// 讀取剩餘數據 (Address + Value/Quantity + CRC)
	remaining := expectedLength - 2
	if remaining > 0 {
		rest := make([]byte, remaining)
		if _, err := io.ReadFull(r.reader, rest); err != nil {
			r.internalClose()
			return nil, fmt.Errorf("RTU 讀取數據失敗: %w", err)
		}
		return append(header, rest...), nil
	}

	return header, nil
}

func (r *RTUTransport) internalClose() {
	if r.port != nil {
		r.port.Close()
		r.port = nil
		r.reader = nil
	}
}
