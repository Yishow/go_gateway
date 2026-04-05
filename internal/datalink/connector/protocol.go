// Package connector 提供統一的協議連接器介面和註冊機制。
//
// 本套件定義了所有協議連接器必須實作的共用介面，並提供：
//   - Protocol: 協議連接器的核心介面
//   - Registry: 協議連接器註冊表
//   - ConnectionManager: 連線池和連線生命週期管理
//   - ReadRequest/ReadResult: 標準化的讀取請求和結果
//
// # 使用範例
//
//	// 註冊協議連接器
//	connector.Register("modbus_tcp", modbus.NewTCPConnector)
//
//	// 從註冊表取得連接器
//	conn, err := connector.Get("modbus_tcp", config)
//	if err != nil {
//	    return err
//	}
//	defer conn.Close()
//
//	// 讀取資料
//	result, err := conn.Read(ctx, request)
package connector

import (
	"context"
	"time"

	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// 核心介面定義
// =============================================================================

// Protocol 定義所有協議連接器必須實作的介面
type Protocol interface {
	// Connect 建立與設備的連線
	// config 為 JSON 格式的連線配置 (協議特定)
	Connect(ctx context.Context, config string) error

	// Close 關閉連線並釋放資源
	Close() error

	// IsConnected 檢查連線狀態
	IsConnected() bool

	// Read 讀取指定位址的資料
	Read(ctx context.Context, req ReadRequest) (ReadResult, error)

	// Write 寫入資料到指定位址
	Write(ctx context.Context, req WriteRequest) error

	// TestConnection 測試連線是否正常
	TestConnection(ctx context.Context) error

	// ProtocolType 取得協議類型
	ProtocolType() schema.ProtocolType
}

// =============================================================================
// 請求與結果結構
// =============================================================================

// ReadRequest 讀取請求
type ReadRequest struct {
	// Address 協議特定的位址字串
	Address string

	// Function 功能碼或存取方式 (協議特定)
	Function string

	// DataType 預期的資料型別
	DataType schema.DataType

	// DataFormat Modbus 多暫存器數值的字節序（ABCD/BADC/CDAB/DCBA）；空字串表示歷史預設（等同 ABCD）。
	DataFormat string

	// Count 讀取數量 (用於批次讀取)
	Count int
}

// ReadResult 讀取結果
type ReadResult struct {
	// Value 讀取到的值 (已轉換為 Go 原生型別)
	Value interface{}

	// RawBytes 原始位元組 (用於追溯)
	RawBytes []byte

	// Timestamp 讀取時間戳記
	Timestamp time.Time

	// Quality 資料品質標誌
	Quality schema.QualityFlag

	// Error 錯誤訊息 (若有)
	Error string
}

// WriteRequest 寫入請求
type WriteRequest struct {
	// Address 協議特定的位址字串
	Address string

	// Function 功能碼或存取方式 (協議特定)
	Function string

	// DataType 資料型別
	DataType schema.DataType

	// Value 要寫入的值
	Value interface{}
}

// =============================================================================
// 批次讀取支援
// =============================================================================

// BatchReader 支援批次讀取的協議連接器可實作此介面
type BatchReader interface {
	Protocol

	// BatchRead 批次讀取多個位址
	BatchRead(ctx context.Context, reqs []ReadRequest) ([]ReadResult, error)
}

// =============================================================================
// 長連接支援
// =============================================================================

/**
 * PersistentConnection 長連接支援介面
 *
 * 參考 HslCommunication 的 NetworkDoubleBase 雙模式設計，
 * 協議連接器可實作此介面以支援長連接模式。
 *
 * 長連接模式優勢：
 *   - 減少 TCP 握手開銷（高頻輪詢時顯著）
 *   - 降低網路延遲
 *   - 減少設備負擔
 *
 * 短連接模式優勢：
 *   - 不佔用設備連線資源
 *   - 適合低頻率存取
 *   - 更簡單的錯誤恢復
 *
 * 使用範例：
 *   if pc, ok := conn.(connector.PersistentConnection); ok {
 *       pc.SetPersistentConnection(true)  // 啟用長連接模式
 *       defer pc.Disconnect()              // 顯式斷線
 *   }
 */
type PersistentConnection interface {
	Protocol

	// SetPersistentConnection 設定是否使用長連接模式
	// enabled: true 啟用長連接，false 使用短連接（每次操作後斷線）
	SetPersistentConnection(enabled bool)

	// IsPersistentMode 檢查當前是否為長連接模式
	IsPersistentMode() bool

	// Disconnect 顯式斷線（僅在長連接模式下有意義）
	// 短連接模式下此方法等同於 Close()
	Disconnect() error

	// Reconnect 重新連線（用於斷線恢復）
	Reconnect(ctx context.Context) error
}

// =============================================================================
// 連線配置解析器
// =============================================================================

// ConfigParser 連線配置解析器介面
type ConfigParser interface {
	// ParseConfig 解析 JSON 配置為協議特定結構
	ParseConfig(jsonConfig string) (interface{}, error)

	// ValidateConfig 驗證配置是否有效
	ValidateConfig(jsonConfig string) error
}
