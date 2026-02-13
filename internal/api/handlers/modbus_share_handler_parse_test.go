package handlers

import "testing"

func TestParsePointLastValue_PreservesLargeIntegerPrecision(t *testing.T) {
	raw := "9007199254740993" // > Number.MAX_SAFE_INTEGER
	got := parsePointLastValue(raw)

	v, ok := got.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T (%v)", got, got)
	}
	if v != 9007199254740993 {
		t.Fatalf("unexpected value: %d", v)
	}
}

func TestParsePointLastValue_JSONFallbackAndBool(t *testing.T) {
	boolVal := parsePointLastValue("true")
	if v, ok := boolVal.(bool); !ok || !v {
		t.Fatalf("expected true bool, got %T (%v)", boolVal, boolVal)
	}

	obj := parsePointLastValue("{\"a\":1}")
	m, ok := obj.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map from JSON object, got %T", obj)
	}
	if _, ok := m["a"]; !ok {
		t.Fatalf("expected key a in parsed object: %#v", m)
	}
}
