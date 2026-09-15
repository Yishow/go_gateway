package adapters

import (
	"math"
	"testing"

	"go-gateway/internal/datalink/schema"
)

func TestMQTTUnsignedTextDoesNotWrapNegativeValues(t *testing.T) {
	if value := convertMQTTValue("-1", schema.DataTypeUint64); value != "-1" {
		t.Fatalf("invalid unsigned text must retain conversion-failure fallback, got %v", value)
	}
	if value := convertMQTTValue("18446744073709551615", schema.DataTypeUint64); value != uint64(math.MaxUint64) {
		t.Fatalf("valid unsigned maximum must parse without signed overflow, got %v", value)
	}
}

func TestMQTTRawPayloadRetainsPlainTextFallback(t *testing.T) {
	value, err := ExtractValueFromPayload([]byte("machine ready"), "", schema.DataTypeString)
	if err != nil || value != "machine ready" {
		t.Fatalf("plain text payload must stay available: value=%v err=%v", value, err)
	}
}
