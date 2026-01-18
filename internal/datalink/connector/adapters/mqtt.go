package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

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

	// 初始化停止通道
	c.stopCh = make(chan struct{})

	// TODO: 實際連線邏輯 (需要 paho.mqtt.golang)
	// 以下為骨架程式碼，展示預期的連線流程
	/*
		opts := mqtt.NewClientOptions()
		opts.AddBroker(c.config.BrokerURL)
		opts.SetClientID(c.config.ClientID)

		if c.config.Username != "" {
			opts.SetUsername(c.config.Username)
			opts.SetPassword(c.config.Password)
		}

		if c.config.UseTLS {
			// 設定 TLS
		}

		opts.SetDefaultPublishHandler(c.messageHandler)
		opts.SetOnConnectHandler(c.onConnect)
		opts.SetConnectionLostHandler(c.onDisconnect)

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			return token.Error()
		}

		c.client = client
	*/

	c.connected = true
	return nil
}

// Close 關閉連線
func (c *MQTTConnector) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopCh != nil {
		close(c.stopCh)
	}

	// TODO: 實際斷線邏輯
	/*
		if c.client != nil {
			c.client.Disconnect(250)
		}
	*/

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
	if !c.IsConnected() {
		return fmt.Errorf("未連線")
	}

	// 地址格式為主題名稱
	topic := req.Address

	// 將值序列化為 JSON
	payload, err := json.Marshal(req.Value)
	if err != nil {
		return fmt.Errorf("序列化寫入值失敗: %w", err)
	}

	// TODO: 實際發布邏輯
	/*
		token := c.client.Publish(topic, c.config.QoS, false, payload)
		token.Wait()
		return token.Error()
	*/

	_ = topic
	_ = payload

	return fmt.Errorf("MQTT 發布功能尚未實作")
}

// =============================================================================
// MQTT 訊息處理 (內部方法)
// =============================================================================

// handleMessage 處理接收到的 MQTT 訊息
func (c *MQTTConnector) handleMessage(topic string, payload []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	// 嘗試解析 JSON
	var value interface{}
	if err := json.Unmarshal(payload, &value); err != nil {
		// 非 JSON 格式，直接使用字串
		value = string(payload)
	}

	// 更新快取
	c.valueCache[topic] = cachedValue{
		Value:     value,
		Timestamp: now,
		Topic:     topic,
		RawBytes:  payload,
	}

	// 通知訂閱者
	msg := MQTTMessage{
		Topic:     topic,
		Payload:   payload,
		Timestamp: now,
	}

	for _, ch := range c.subscribers {
		select {
		case ch <- msg:
		default:
			// 通道已滿，跳過
		}
	}
}

// Subscribe 訂閱訊息 (用於即時處理)
func (c *MQTTConnector) Subscribe(bufferSize int) <-chan MQTTMessage {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan MQTTMessage, bufferSize)
	c.subscribers = append(c.subscribers, ch)
	return ch
}

// Unsubscribe 取消訂閱
func (c *MQTTConnector) Unsubscribe(ch <-chan MQTTMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, sub := range c.subscribers {
		if sub == ch {
			c.subscribers = append(c.subscribers[:i], c.subscribers[i+1:]...)
			close(sub)
			break
		}
	}
}

// GetCachedTopics 取得所有已快取的主題
func (c *MQTTConnector) GetCachedTopics() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	topics := make([]string, 0, len(c.valueCache))
	for topic := range c.valueCache {
		topics = append(topics, topic)
	}
	return topics
}

// GetCachedValue 取得指定主題的快取值
func (c *MQTTConnector) GetCachedValue(topic string) (cachedValue, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.valueCache[topic]
	return val, ok
}

// =============================================================================
// MQTT 主題映射輔助
// =============================================================================

// TopicMapping MQTT 主題到點位的映射配置
type TopicMapping struct {
	// Topic MQTT 主題 (支援萬用字元)
	Topic string `json:"topic"`

	// PayloadPath JSON 路徑 (用於從 payload 擷取值，如 "data.temperature")
	PayloadPath string `json:"payload_path,omitempty"`

	// DataType 資料型別
	DataType schema.DataType `json:"data_type"`
}

// ExtractValueFromPayload 從 MQTT payload 擷取值
func ExtractValueFromPayload(payload []byte, path string, dataType schema.DataType) (interface{}, error) {
	// 如果沒有路徑，直接解析整個 payload
	if path == "" {
		var value interface{}
		if err := json.Unmarshal(payload, &value); err != nil {
			// 非 JSON，嘗試直接解析為字串
			return string(payload), nil
		}
		return convertMQTTValue(value, dataType), nil
	}

	// 解析 JSON 並按路徑擷取
	var obj map[string]interface{}
	if err := json.Unmarshal(payload, &obj); err != nil {
		return nil, fmt.Errorf("payload 不是 JSON 物件: %w", err)
	}

	// 簡單的點號路徑解析
	parts := splitPath(path)
	current := interface{}(obj)

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[part]
			if !ok {
				return nil, fmt.Errorf("路徑 '%s' 中找不到欄位 '%s'", path, part)
			}
		default:
			return nil, fmt.Errorf("路徑 '%s' 無效，'%s' 不是物件", path, part)
		}
	}

	return convertMQTTValue(current, dataType), nil
}

// splitPath 分割 JSON 路徑
func splitPath(path string) []string {
	// 簡單實作：按點號分割
	result := make([]string, 0)
	current := ""
	for _, ch := range path {
		if ch == '.' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// convertMQTTValue 將 MQTT 值轉換為指定型別
func convertMQTTValue(value interface{}, dataType schema.DataType) interface{} {
	switch dataType {
	case schema.DataTypeBool:
		switch v := value.(type) {
		case bool:
			return v
		case float64:
			return v != 0
		case string:
			return v == "true" || v == "1" || v == "on"
		}
	case schema.DataTypeInt16, schema.DataTypeInt32, schema.DataTypeInt64:
		switch v := value.(type) {
		case float64:
			return int64(v)
		case string:
			if n, err := json.Number(v).Int64(); err == nil {
				return n
			}
		}
	case schema.DataTypeUint16, schema.DataTypeUint32, schema.DataTypeUint64:
		switch v := value.(type) {
		case float64:
			return uint64(v)
		case string:
			if n, err := json.Number(v).Int64(); err == nil {
				return uint64(n)
			}
		}
	case schema.DataTypeFloat32, schema.DataTypeFloat64:
		switch v := value.(type) {
		case float64:
			return v
		case string:
			if n, err := json.Number(v).Float64(); err == nil {
				return n
			}
		}
	case schema.DataTypeString:
		switch v := value.(type) {
		case string:
			return v
		default:
			data, _ := json.Marshal(v)
			return string(data)
		}
	}

	return value
}

// =============================================================================
// 註冊連接器
// =============================================================================

func init() {
	connector.Register(schema.ProtocolMQTT, NewMQTTConnector)
}
