package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
)

// =============================================================================
// FATEK FBs 連接器
// =============================================================================

// FatekConnector FATEK FBs 協議連接器適配器
type FatekConnector struct {
	client         *fatek.FatekClient
	transport      fatek.Transport
	config         schema.ConnectionConfigFatekFBs
	connected      bool
	persistentMode bool // 長連接模式標誌
}

// NewFatekConnector 建立新的 FATEK 連接器
func NewFatekConnector() connector.Protocol {
	return &FatekConnector{}
}

// Connect 建立連線
func (c *FatekConnector) Connect(ctx context.Context, configJSON string) error {
	// 解析配置
	if err := json.Unmarshal([]byte(configJSON), &c.config); err != nil {
		return fmt.Errorf("解析 FATEK 配置失敗: %w", err)
	}

	// 設定預設值
	if c.config.Timeout == 0 {
		c.config.Timeout = 5
	}

	// 根據模式建立傳輸層
	switch strings.ToLower(c.config.Mode) {
	case "tcp":
		if c.config.Port == 0 {
			c.config.Port = 500
		}
		transport := fatek.NewTCPTransport(c.config.Host, c.config.Port)
		transport.Timeout = time.Duration(c.config.Timeout) * time.Second
		c.transport = transport
	case "serial":
		if c.config.BaudRate == 0 {
			c.config.BaudRate = 9600
		}
		c.transport = fatek.NewSerialTransport(
			c.config.SerialPort,
			c.config.BaudRate,
			7,      // DataBits (FATEK 預設 7)
			1,      // StopBits
			"even", // Parity (FATEK 預設 Even)
			time.Duration(c.config.Timeout)*time.Second,
		)
	default:
		return fmt.Errorf("不支援的 FATEK 連線模式: %s", c.config.Mode)
	}

	// 建立客戶端
	c.client = fatek.NewClient(c.transport, int(c.config.StationNo))

	// 預設使用長連接模式
	if !c.persistentMode {
		c.persistentMode = true
	}

	// 連線（僅在長連接模式下立即連線）
	if c.persistentMode {
		if err := c.client.Connect(); err != nil {
			return fmt.Errorf("FATEK 連線失敗: %w", err)
		}
		c.connected = true
	}

	return nil
}

// Close 關閉連線
func (c *FatekConnector) Close() error {
	if c.client != nil {
		c.connected = false
		return c.client.Close()
	}
	return nil
}

// IsConnected 檢查連線狀態
func (c *FatekConnector) IsConnected() bool {
	return c.connected
}

// ProtocolType 取得協議類型
func (c *FatekConnector) ProtocolType() schema.ProtocolType {
	return schema.ProtocolFatekFBs
}
