package hsllogic

import (
	"fmt"
	"io"
	"os"
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
	// HexString 十六進位字串表示（延遲生成）
	HexString string `json:"hex_string"`
	// Description 描述說明
	Description string `json:"description,omitempty"`
	// Duration 通訊耗時 (僅用於接收報文)
	Duration time.Duration `json:"duration,omitempty"`
	// hexCached 標記 HexString 是否已快取
	hexCached bool
}

// GetHexString 取得十六進位字串（延遲生成）
func (p *PacketLog) GetHexString() string {
	if !p.hexCached {
		p.HexString = formatHexBytes(p.RawData)
		p.hexCached = true
	}
	return p.HexString
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
	sb.WriteString(p.GetHexString())
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
// PacketLog 物件池
// =============================================================================

// packetLogPool PacketLog 物件池，用於減少 GC 壓力
var packetLogPool = sync.Pool{
	New: func() interface{} {
		return &PacketLog{}
	},
}

// acquirePacketLog 從物件池取得 PacketLog
func acquirePacketLog() *PacketLog {
	log, ok := packetLogPool.Get().(*PacketLog)
	if !ok {
		return &PacketLog{}
	}
	return log
}

// releasePacketLog 歸還 PacketLog 至物件池
func releasePacketLog(log *PacketLog) {
	// 重置敏感欄位，避免資料洩漏
	log.Timestamp = time.Time{}
	log.Direction = ""
	log.Protocol = ""
	log.DeviceID = ""
	log.RawData = nil
	log.HexString = ""
	log.Description = ""
	log.Duration = 0
	log.hexCached = false
	packetLogPool.Put(log)
}
