package modbus

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go-gateway/internal/protocol/wirebudget"
	"io"
	"time"

	"go.bug.st/serial"
)

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
	mu       wirebudget.Mutex
	closeErr error
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

func (r *RTUTransport) Connect() error { return r.ConnectContext(context.Background()) }
func (r *RTUTransport) ConnectContext(ctx context.Context) error {
	ctx, cancel := transportContext(ctx, r.Timeout)
	defer cancel()
	if err := r.mu.LockContext(ctx); err != nil {
		return err
	}
	defer r.mu.Unlock()

	if err := r.closePort(); err != nil {
		return err
	}
	// Closing the retired handle may outlive the caller. Never start a new
	// non-context-aware OS open after cancellation during that cleanup.
	if err := ctx.Err(); err != nil {
		return err
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
	if err := ctx.Err(); err != nil {
		return errors.Join(err, r.closePort())
	}
	return nil
}

func (r *RTUTransport) Close() error { return r.CloseContext(context.Background()) }
func (r *RTUTransport) CloseContext(ctx context.Context) error {
	if err := r.mu.LockContext(ctx); err != nil {
		return err
	}
	defer r.mu.Unlock()
	return r.closePort()
}
func (r *RTUTransport) closePort() error {
	if r.port == nil {
		return nil
	}
	if err := r.port.Close(); err != nil {
		r.closeErr = err
		return err
	}
	r.port = nil
	r.reader = nil
	r.closeErr = nil
	return nil
}

// IsConnected 檢查 RTU 序列埠連線狀態
func (r *RTUTransport) IsConnected() bool {
	live, _ := r.IsConnectedContext(context.Background())
	return live
}
func (r *RTUTransport) IsConnectedContext(ctx context.Context) (bool, error) {
	if err := r.mu.LockContext(ctx); err != nil {
		return false, err
	}
	defer r.mu.Unlock()
	return r.port != nil && r.closeErr == nil, nil
}

func (r *RTUTransport) SendReceive(data []byte) ([]byte, error) {
	return r.SendReceiveContext(context.Background(), data)
}
func (r *RTUTransport) SendReceiveContext(ctx context.Context, data []byte) (response []byte, resultErr error) {
	return r.sendReceiveContext(ctx, data, false)
}

func (r *RTUTransport) sendReceiveContext(ctx context.Context, data []byte, exact bool) (response []byte, resultErr error) {
	ctx, cancel := transportContext(ctx, r.Timeout)
	defer cancel()
	if err := r.mu.LockContext(ctx); err != nil {
		return nil, err
	}
	defer r.mu.Unlock()

	if r.closeErr != nil {
		return nil, r.closeErr
	}
	if r.port == nil {
		return nil, ErrConnectionClosed
	}

	if b := wirebudget.From(ctx); b != nil && b.Remaining() == 0 {
		return nil, wirebudget.ErrExhausted
	}
	port := r.port
	finish := serialContext(ctx, port)
	defer func() {
		closeErr := finish(resultErr != nil)
		if closeErr != nil {
			r.closeErr = closeErr
			resultErr = errors.Join(resultErr, closeErr)
		}
		if closeErr == nil && (resultErr != nil || ctx.Err() != nil) {
			r.port = nil
			r.reader = nil
		}
		if ctx.Err() != nil {
			resultErr = errors.Join(ctx.Err(), closeErr)
		}
	}()
	reader := contextSerialReader{ctx: ctx, reader: port}
	deadline, _ := ctx.Deadline()
	readTimeout := time.Until(deadline)
	if err := port.SetReadTimeout(readTimeout); err != nil {
		return nil, fmt.Errorf("RTU 設定讀取逾時失敗: %w", err)
	}

	// 寫入請求
	if err := wirebudget.Take(ctx); err != nil {
		return nil, err
	}
	if n, err := port.Write(data); err != nil || n != len(data) {
		if err == nil {
			err = io.ErrShortWrite
		}
		return nil, fmt.Errorf("RTU 寫入失敗: %w", err)
	}

	if exact {
		return r.readExactFrame(ctx, port, readTimeout)
	}

	// RTU 需要等待一段時間 (根據波特率計算)
	// 通常為 3.5 個字符時間
	baudRate := r.BaudRate
	if baudRate <= 0 {
		baudRate = DefaultBaudRate
	}
	charTime := time.Duration(10000000/baudRate) * time.Microsecond
	timer := time.NewTimer(charTime * 35 / 10)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
	}

	// Read the response.
	// 先讀取地址和功能碼
	header := make([]byte, 2)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, fmt.Errorf("RTU 讀取標頭失敗: %w", err)
	}

	functionCode := header[1]

	// 檢查是否為異常回應
	if functionCode&FuncExceptionOffset != 0 {
		// 異常回應: Address(1) + Function(1) + Exception(1) + CRC(2)
		exceptionAndCRC := make([]byte, 3)
		if _, err := io.ReadFull(reader, exceptionAndCRC); err != nil {
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
		if _, err := io.ReadFull(reader, byteCount); err != nil {
			return nil, fmt.Errorf("RTU 讀取 ByteCount 失敗: %w", err)
		}
		// 讀取數據和 CRC
		dataAndCRC := make([]byte, int(byteCount[0])+2)
		if _, err := io.ReadFull(reader, dataAndCRC); err != nil {
			return nil, fmt.Errorf("RTU 讀取數據失敗: %w", err)
		}
		response := make([]byte, 0, len(header)+len(byteCount)+len(dataAndCRC))
		response = append(response, header...)
		response = append(response, byteCount...)
		response = append(response, dataAndCRC...)
		return response, nil

	case FuncWriteSingleCoil, FuncWriteSingleRegister:
		// 寫入單個回應: Address(1) + Function(1) + Address(2) + Value(2) + CRC(2)
		expectedLength = 2 + 4 + 2 // Header + Address(2) + Value(2) + CRC

	case FuncWriteMultipleCoils, FuncWriteMultipleRegisters:
		// 寫入多個回應: Address(1) + Function(1) + Address(2) + Quantity(2) + CRC(2)
		expectedLength = 2 + 4 + 2 // Header + Address(2) + Quantity(2) + CRC

	default:
		return nil, ErrInvalidFunctionCode
	}

	// Read the remaining address, value/quantity, and CRC bytes.
	remaining := expectedLength - 2
	if remaining > 0 {
		rest := make([]byte, remaining)
		if _, err := io.ReadFull(reader, rest); err != nil {
			return nil, fmt.Errorf("RTU 讀取數據失敗: %w", err)
		}
		return append(header, rest...), nil
	}

	return header, nil
}

func (r *RTUTransport) internalClose() { _ = r.closePort() }
