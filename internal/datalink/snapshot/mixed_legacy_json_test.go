package snapshot

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMixedFieldsDoNotExpandLegacyJSON(t *testing.T) {
	for _, v := range []any{Member{Key: "legacy"}, MemberResult{MemberKey: "legacy", Status: MemberMissing}, Outcome{Kind: OutcomeNoData}} {
		data, err := json.Marshal(v)
		require.NoError(t, err)
		for _, field := range []string{`"Calculation":`, `"Interval":`, `"Window":`, `"DueAt":`} {
			require.NotContains(t, string(data), field)
		}
	}
	data, err := json.Marshal(Outcome{DueAt: mixedTime(15), Members: []MemberResult{{MemberKey: "mixed", Window: &MemberWindow{Calculation: CalculationLatest}}}})
	require.NoError(t, err)
	require.Contains(t, string(data), `"DueAt":`)
	require.Contains(t, string(data), `"Window":`)
}
