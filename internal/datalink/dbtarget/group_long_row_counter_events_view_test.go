package dbtarget

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// provenanceView is the cumulative part of a stored provenance document.
type provenanceView struct {
	Reason  string `json:"reason"`
	Counter struct {
		SpanTotal        bool   `json:"span_total"`
		RolloverModulus  string `json:"rollover_modulus"`
		MaxRatePerSecond string `json:"max_rate_per_second"`
		EventsOmitted    int    `json:"events_omitted"`
		Events           []struct {
			Kind      string `json:"kind"`
			Reason    string `json:"reason"`
			FromValue string `json:"from_value"`
			ToValue   string `json:"to_value"`
		} `json:"events"`
	} `json:"counter"`
}

func decodeProvenance(t *testing.T, cell any) provenanceView {
	t.Helper()
	var view provenanceView
	require.NoError(t, json.Unmarshal([]byte(fmt.Sprint(cell)), &view))
	return view
}
