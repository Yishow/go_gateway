package adapters

import (
	"context"
	"testing"
	"time"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
)

func TestMQTTConnector_SubscribeAndCache(t *testing.T) {
	conn := NewMQTTConnector().(*MQTTConnector)
	conn.connected = true

	ch := conn.Subscribe(1)
	payload := []byte(`{"value":42}`)
	conn.handleMessage("sensors/temp", payload)

	select {
	case msg := <-ch:
		if msg.Topic != "sensors/temp" {
			t.Fatalf("expected topic sensors/temp, got %s", msg.Topic)
		}
		if string(msg.Payload) != string(payload) {
			t.Fatalf("payload mismatch")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("expected subscriber message")
	}

	if _, ok := conn.GetCachedValue("sensors/temp"); !ok {
		t.Fatal("expected cached value")
	}

	conn.Unsubscribe(ch)
}

func TestMQTTConnector_ReadUnknownTopic(t *testing.T) {
	conn := NewMQTTConnector().(*MQTTConnector)
	conn.connected = true

	res, err := conn.Read(context.Background(), connector.ReadRequest{Address: "unknown"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Error == "" {
		t.Fatal("expected error message for unknown topic")
	}
}

func TestExtractValueFromPayload(t *testing.T) {
	payload := []byte(`{"data":{"temp":12.5},"flag":true}`)

	value, err := ExtractValueFromPayload(payload, "data.temp", schema.DataTypeFloat64)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value.(float64) != 12.5 {
		t.Fatalf("expected 12.5, got %v", value)
	}

	value, err = ExtractValueFromPayload(payload, "flag", schema.DataTypeBool)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != true {
		t.Fatalf("expected true, got %v", value)
	}
}
