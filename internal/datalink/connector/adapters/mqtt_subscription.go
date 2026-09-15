package adapters

import (
	"encoding/json"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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

func (c *MQTTConnector) onConnect(client mqtt.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.stopCh == nil {
		return
	}

	for _, topic := range c.config.Topics {
		token := client.Subscribe(topic, c.config.QoS, func(_ mqtt.Client, msg mqtt.Message) {
			c.handleMessage(msg.Topic(), msg.Payload())
		})
		token.Wait()
	}
}

func (c *MQTTConnector) onDisconnect(_ error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected = false
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
func (c *MQTTConnector) GetCachedValue(topic string) (cachedValue, bool) { //nolint:revive // Preserve the existing public cache accessor signature.
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.valueCache[topic]
	return val, ok
}
