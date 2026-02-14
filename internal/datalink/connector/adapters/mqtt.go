package adapters

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

// =============================================================================
// MQTT 接收連接器
// =============================================================================

// MQTTConnector MQTT 訂閱接收連接器
// 注意：此連接器需要外部 MQTT 客戶端庫 (如 paho.mqtt.golang)
// 目前提供骨架實作，實際使用時需要引入相應依賴
type MQTTConnector struct {
	config    schema.ConnectionConfigMQTT
	connected bool
	mu        sync.RWMutex
	client    mqtt.Client

	// valueCache 快取最新接收到的值
	valueCache map[string]cachedValue

	// subscribers 訂閱者通道
	subscribers []chan MQTTMessage

	// stopCh 停止信號
	stopCh chan struct{}
}

// cachedValue 快取的值結構
type cachedValue struct {
	Value     interface{}
	Timestamp time.Time
	Topic     string
	RawBytes  []byte
}

// MQTTMessage MQTT 訊息結構
type MQTTMessage struct {
	Topic     string
	Payload   []byte
	Timestamp time.Time
}

// NewMQTTConnector 建立新的 MQTT 連接器
func NewMQTTConnector() connector.Protocol {
	return &MQTTConnector{
		valueCache:  make(map[string]cachedValue),
		subscribers: make([]chan MQTTMessage, 0),
	}
}

// Connect 建立連線
// 注意：完整實作需要引入 paho.mqtt.golang
func (c *MQTTConnector) Connect(ctx context.Context, configJSON string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 解析配置
	if err := json.Unmarshal([]byte(configJSON), &c.config); err != nil {
		return fmt.Errorf("解析 MQTT 配置失敗: %w", err)
	}

	// 驗證配置
	if c.config.BrokerURL == "" {
		return fmt.Errorf("MQTT Broker URL 不能為空")
	}
	if c.config.ClientID == "" {
		return fmt.Errorf("MQTT Client ID 不能為空")
	}
	if len(c.config.Topics) == 0 {
		return fmt.Errorf("MQTT 訂閱主題不能為空")
	}

	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
	}

	// 初始化停止通道
	c.stopCh = make(chan struct{})

	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.config.BrokerURL)
	opts.SetClientID(c.config.ClientID)
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetWriteTimeout(10 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetDefaultPublishHandler(func(_ mqtt.Client, msg mqtt.Message) {
		c.handleMessage(msg.Topic(), msg.Payload())
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		c.onConnect(client)
	})
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		c.onDisconnect(err)
	})

	if c.config.Username != "" {
		opts.SetUsername(c.config.Username)
		opts.SetPassword(c.config.Password)
	}
	if c.config.UseTLS {
		opts.SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12})
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("MQTT 連線逾時")
	}
	if err := token.Error(); err != nil {
		return err
	}

	c.client = client
	if c.stopCh == nil {
		c.stopCh = make(chan struct{})
	}

	c.connected = true
	return nil
}

// Close 關閉連線
func (c *MQTTConnector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopCh != nil {
		close(c.stopCh)
		c.stopCh = nil
	}

	if c.client != nil {
		c.client.Disconnect(250)
		c.client = nil
	}

	c.connected = false
	return nil
}

// IsConnected 檢查連線狀態
func (c *MQTTConnector) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// ProtocolType 取得協議類型
func (c *MQTTConnector) ProtocolType() schema.ProtocolType {
	return schema.ProtocolMQTT
}

// TestConnection 測試連線
func (c *MQTTConnector) TestConnection(ctx context.Context) error {
	if !c.IsConnected() {
		return fmt.Errorf("未連線")
	}
	// MQTT 連線狀態由客戶端庫管理
	return nil
}

// Read 讀取資料 (從快取讀取最新接收到的值)
func (c *MQTTConnector) Read(ctx context.Context, req connector.ReadRequest) (connector.ReadResult, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := connector.ReadResult{
		Timestamp: time.Now(),
		Quality:   schema.QualityGood,
	}

	if !c.connected {
		result.Quality = schema.QualityBad
		result.Error = "未連線"
		return result, fmt.Errorf("未連線")
	}

	// 地址格式為主題名稱或主題+路徑
	// 例如: "sensors/temperature" 或 "sensors/+/value"
	topic := req.Address

	cached, ok := c.valueCache[topic]
	if !ok {
		result.Quality = schema.QualityUncertain
		result.Error = "尚未接收到此主題的資料"
		return result, nil
	}

	result.Value = cached.Value
	result.Timestamp = cached.Timestamp
	result.RawBytes = cached.RawBytes

	return result, nil
}

// Write 寫入資料 (發布到 MQTT)
func (c *MQTTConnector) Write(ctx context.Context, req connector.WriteRequest) error {
	c.mu.RLock()
	client := c.client
	connected := c.connected
	qos := c.config.QoS
	c.mu.RUnlock()

	if !connected || client == nil || !client.IsConnected() {
		return fmt.Errorf("未連線")
	}

	// 地址格式為主題名稱
	topic := req.Address

	payload, err := json.Marshal(req.Value)
	if err != nil {
		return fmt.Errorf("序列化寫入值失敗: %w", err)
	}

	token := client.Publish(topic, qos, false, payload)
	if !token.WaitTimeout(10 * time.Second) {
		return fmt.Errorf("MQTT 發布逾時")
	}
	if err := token.Error(); err != nil {
		return err
	}

	return nil
}
