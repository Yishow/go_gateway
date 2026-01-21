package hsllogic

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// =============================================================================
// PacketDirection 報文方向
// =============================================================================

// PacketDirection 報文方向常數
type PacketDirection string

const (
	// PacketSend 發送報文
	PacketSend PacketDirection = "SEND"
	// PacketRecv 接收報文
	PacketRecv PacketDirection = "RECV"
)

// =============================================================================
// PacketLog 報文日誌條目
// =============================================================================

// PacketLog 報文日誌條目
type PacketLog struct {
	// Timestamp 時間戳記
	Timestamp time.Time `json:"timestamp"`
	// Direction 方向 (SEND/RECV)
	Direction PacketDirection `json:"direction"`
	// Protocol 協議類型
	Protocol string `json:"protocol"`
	// DeviceID 設備識別碼
	DeviceID string `json:"device_id"`
	// RawData 原始位元組資料
	RawData []byte `json:"raw_data"`
	// HexString 十六進位字串表示
	HexString string `json:"hex_string"`
	// Description 描述說明
	Description string `json:"description,omitempty"`
	// Duration 通訊耗時 (僅用於接收報文)
	Duration time.Duration `json:"duration,omitempty"`
}

// String 取得報文日誌的字串表示
func (p *PacketLog) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] ", p.Timestamp.Format("2006-01-02 15:04:05.000")))
	sb.WriteString(fmt.Sprintf("[%s] ", p.Direction))
	sb.WriteString(fmt.Sprintf("[%s] ", p.Protocol))
	if p.DeviceID != "" {
		sb.WriteString(fmt.Sprintf("[%s] ", p.DeviceID))
	}
	sb.WriteString(p.HexString)
	if p.Description != "" {
		sb.WriteString(fmt.Sprintf(" ; %s", p.Description))
	}
	if p.Duration > 0 {
		sb.WriteString(fmt.Sprintf(" (%v)", p.Duration))
	}
	return sb.String()
}

// =============================================================================
// PacketLoggerOptions 報文日誌器選項
// =============================================================================

// PacketLoggerOptions 報文日誌器選項
type PacketLoggerOptions struct {
	// Enabled 是否啟用
	Enabled bool
	// LogToFile 是否寫入檔案
	LogToFile bool
	// LogDir 日誌目錄
	LogDir string
	// MaxFileSize 最大檔案大小 (位元組)
	MaxFileSize int64
	// MaxFiles 最大檔案數量
	MaxFiles int
	// LogToConsole 是否輸出到控制台
	LogToConsole bool
	// ConsoleWriter 控制台輸出器
	ConsoleWriter io.Writer
	// BufferSize 緩衝區大小
	BufferSize int
	// FlushInterval 刷新間隔
	FlushInterval time.Duration
}

// DefaultPacketLoggerOptions 預設報文日誌器選項
var DefaultPacketLoggerOptions = PacketLoggerOptions{
	Enabled:       true,
	LogToFile:     true,
	LogDir:        "logs/packets",
	MaxFileSize:   10 * 1024 * 1024, // 10MB
	MaxFiles:      10,
	LogToConsole:  false,
	ConsoleWriter: os.Stdout,
	BufferSize:    1000,
	FlushInterval: 5 * time.Second,
}

// =============================================================================
// PacketLogger 報文日誌器
// =============================================================================

// PacketLogger 報文日誌器，用於記錄通訊報文軌跡
type PacketLogger struct {
	options       PacketLoggerOptions
	mu            sync.RWMutex
	buffer        chan *PacketLog
	currentFile   *os.File
	currentSize   int64
	stopChan      chan struct{}
	wg            sync.WaitGroup
	fileIndex     int
	lastFlushTime time.Time
}

// NewPacketLogger 建立新的報文日誌器
func NewPacketLogger(options PacketLoggerOptions) *PacketLogger {
	if options.BufferSize <= 0 {
		options.BufferSize = DefaultPacketLoggerOptions.BufferSize
	}
	if options.FlushInterval <= 0 {
		options.FlushInterval = DefaultPacketLoggerOptions.FlushInterval
	}
	if options.MaxFileSize <= 0 {
		options.MaxFileSize = DefaultPacketLoggerOptions.MaxFileSize
	}
	if options.MaxFiles <= 0 {
		options.MaxFiles = DefaultPacketLoggerOptions.MaxFiles
	}

	pl := &PacketLogger{
		options:       options,
		buffer:        make(chan *PacketLog, options.BufferSize),
		stopChan:      make(chan struct{}),
		lastFlushTime: time.Now(),
	}

	if options.Enabled && options.LogToFile {
		pl.wg.Add(1)
		go pl.writeLoop()
	}

	return pl
}

// =============================================================================
// 公開方法
// =============================================================================

// LogSend 記錄發送報文
func (pl *PacketLogger) LogSend(protocol, deviceID string, data []byte, desc string) {
	pl.log(PacketSend, protocol, deviceID, data, desc, 0)
}

// LogRecv 記錄接收報文
func (pl *PacketLogger) LogRecv(protocol, deviceID string, data []byte, desc string, duration time.Duration) {
	pl.log(PacketRecv, protocol, deviceID, data, desc, duration)
}

// LogTransaction 記錄完整的通訊交易 (發送 + 接收)
func (pl *PacketLogger) LogTransaction(
	protocol, deviceID string,
	sendData, recvData []byte,
	sendDesc, recvDesc string,
	duration time.Duration,
) {
	pl.LogSend(protocol, deviceID, sendData, sendDesc)
	pl.LogRecv(protocol, deviceID, recvData, recvDesc, duration)
}

// Close 關閉日誌器
func (pl *PacketLogger) Close() error {
	close(pl.stopChan)
	pl.wg.Wait()

	pl.mu.Lock()
	defer pl.mu.Unlock()

	if pl.currentFile != nil {
		return pl.currentFile.Close()
	}
	return nil
}

// Flush 手動刷新緩衝區
func (pl *PacketLogger) Flush() {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	if pl.currentFile != nil {
		pl.currentFile.Sync()
	}
}

// =============================================================================
// 內部方法
// =============================================================================

// log 記錄報文
func (pl *PacketLogger) log(
	direction PacketDirection,
	protocol, deviceID string,
	data []byte,
	desc string,
	duration time.Duration,
) {
	if !pl.options.Enabled {
		return
	}

	entry := &PacketLog{
		Timestamp:   time.Now(),
		Direction:   direction,
		Protocol:    protocol,
		DeviceID:    deviceID,
		RawData:     data,
		HexString:   formatHexBytes(data),
		Description: desc,
		Duration:    duration,
	}

	// 控制台輸出
	if pl.options.LogToConsole && pl.options.ConsoleWriter != nil {
		fmt.Fprintln(pl.options.ConsoleWriter, entry.String())
	}

	// 寫入緩衝區
	if pl.options.LogToFile {
		select {
		case pl.buffer <- entry:
		default:
			// 緩衝區已滿，丟棄最舊的條目
			select {
			case <-pl.buffer:
				pl.buffer <- entry
			default:
			}
		}
	}
}

// writeLoop 寫入迴圈
func (pl *PacketLogger) writeLoop() {
	defer pl.wg.Done()

	ticker := time.NewTicker(pl.options.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-pl.stopChan:
			// 處理剩餘的緩衝區
			for len(pl.buffer) > 0 {
				entry := <-pl.buffer
				pl.writeEntry(entry)
			}
			return
		case entry := <-pl.buffer:
			pl.writeEntry(entry)
		case <-ticker.C:
			pl.Flush()
		}
	}
}

// writeEntry 寫入單筆條目
func (pl *PacketLogger) writeEntry(entry *PacketLog) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	// 確保檔案已開啟
	if err := pl.ensureFile(); err != nil {
		return
	}

	// 寫入日誌
	line := entry.String() + "\n"
	n, err := pl.currentFile.WriteString(line)
	if err != nil {
		return
	}

	pl.currentSize += int64(n)

	// 檢查是否需要輪替
	if pl.currentSize >= pl.options.MaxFileSize {
		pl.rotateFile()
	}
}

// ensureFile 確保日誌檔案已開啟
func (pl *PacketLogger) ensureFile() error {
	if pl.currentFile != nil {
		return nil
	}

	// 建立目錄
	if err := os.MkdirAll(pl.options.LogDir, 0755); err != nil {
		return err
	}

	// 開啟新檔案
	filename := filepath.Join(
		pl.options.LogDir,
		fmt.Sprintf("packet_%s.log", time.Now().Format("20060102_150405")),
	)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	pl.currentFile = f
	pl.currentSize = 0
	return nil
}

// rotateFile 輪替檔案
func (pl *PacketLogger) rotateFile() {
	if pl.currentFile != nil {
		pl.currentFile.Close()
		pl.currentFile = nil
	}

	pl.fileIndex++
	pl.cleanOldFiles()
}

// cleanOldFiles 清理舊檔案
func (pl *PacketLogger) cleanOldFiles() {
	files, err := filepath.Glob(filepath.Join(pl.options.LogDir, "packet_*.log"))
	if err != nil {
		return
	}

	if len(files) <= pl.options.MaxFiles {
		return
	}

	// 刪除最舊的檔案
	for i := 0; i < len(files)-pl.options.MaxFiles; i++ {
		os.Remove(files[i])
	}
}

// =============================================================================
// 輔助函數
// =============================================================================

// formatHexBytes 格式化位元組為十六進位字串
func formatHexBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	return strings.ToUpper(hex.EncodeToString(data))
}

// FormatHexWithSpaces 格式化位元組為帶空格的十六進位字串
func FormatHexWithSpaces(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	parts := make([]string, len(data))
	for i, b := range data {
		parts[i] = fmt.Sprintf("%02X", b)
	}
	return strings.Join(parts, " ")
}

// ParseHexString 解析十六進位字串為位元組
func ParseHexString(hexStr string) ([]byte, error) {
	// 移除空格和前綴
	hexStr = strings.ReplaceAll(hexStr, " ", "")
	hexStr = strings.TrimPrefix(hexStr, "0x")
	hexStr = strings.TrimPrefix(hexStr, "0X")

	return hex.DecodeString(hexStr)
}

// =============================================================================
// 全域日誌器
// =============================================================================

var (
	globalLogger     *PacketLogger
	globalLoggerOnce sync.Once
)

// GetGlobalPacketLogger 取得全域報文日誌器
func GetGlobalPacketLogger() *PacketLogger {
	globalLoggerOnce.Do(func() {
		globalLogger = NewPacketLogger(DefaultPacketLoggerOptions)
	})
	return globalLogger
}

// SetGlobalPacketLogger 設定全域報文日誌器
func SetGlobalPacketLogger(logger *PacketLogger) {
	globalLogger = logger
}

// LogPacketSend 使用全域日誌器記錄發送報文
func LogPacketSend(protocol, deviceID string, data []byte, desc string) {
	GetGlobalPacketLogger().LogSend(protocol, deviceID, data, desc)
}

// LogPacketRecv 使用全域日誌器記錄接收報文
func LogPacketRecv(protocol, deviceID string, data []byte, desc string, duration time.Duration) {
	GetGlobalPacketLogger().LogRecv(protocol, deviceID, data, desc, duration)
}
