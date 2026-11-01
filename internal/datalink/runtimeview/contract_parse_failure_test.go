package runtimeview

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTimeComparisonsRejectInvalidInputs(t *testing.T) {
	valid := "2026-10-27T10:00:00Z"
	for _, pair := range [][2]string{{"invalid", "invalid"}, {valid, "invalid"}, {"invalid", valid}} {
		if afterTime(pair[0], pair[1]) || sameTime(pair[0], pair[1]) {
			t.Fatalf("invalid timestamp compared as a zero time: %q", pair)
		}
	}
	if !sameTime(valid, "2026-10-27T10:00:00.000Z") || !afterTime("2026-10-27T10:00:00.000000001Z", valid) || afterTime(valid, valid) {
		t.Fatal("valid timestamp ordering changed")
	}
}

func TestNumericSchemaBoundsRejectInvalidTags(t *testing.T) {
	for _, tag := range []reflect.StructTag{`min:"invalid"`, `max:"invalid"`} {
		err := checkJSON([]byte("0"), reflect.TypeFor[float64](), tag)
		require.Same(t, errUnavailable, err, "invalid fixed bound must fail closed: %q", tag)
	}
	tag := reflect.StructTag(`signed:"true" min:"-1.5" max:"2.5"`)
	for _, raw := range []string{"-1.5", "0", "2.5"} {
		if err := checkJSON([]byte(raw), reflect.TypeFor[float64](), tag); err != nil {
			t.Fatalf("valid inclusive bound rejected: %s: %v", raw, err)
		}
	}
	for _, raw := range []string{"-1.5001", "2.5001"} {
		err := checkJSON([]byte(raw), reflect.TypeFor[float64](), tag)
		require.Same(t, errUnavailable, err, "out-of-bound value accepted: %s", raw)
	}
}
