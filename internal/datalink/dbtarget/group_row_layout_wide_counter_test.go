package dbtarget

import (
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func wideCounterLayout(t *testing.T, member LongRowMember, source schema.DataType) *GroupRowLayout {
	t.Helper()
	spec := wideProfileSpec(func(s *GroupRowSpec) {
		s.Members = []GroupRowMember{{
			MemberKey: member.MemberKey, Column: "v_mean", Type: member.Type, SourceValueType: source, Counter: member.Counter,
		}}
	})
	layout, issues := NewGroupRowLayout(spec)
	require.Empty(t, issues)
	return layout
}

type wideCounterEntry struct {
	Status  string `json:"status"`
	Counter *struct {
		Kind                  string `json:"kind"`
		Basis                 string `json:"basis"`
		RawText               string `json:"raw_text"`
		Status                string `json:"status"`
		Reason                string `json:"reason"`
		Delta                 string `json:"delta"`
		DeltaFrom             string `json:"delta_from"`
		DeltaTo               string `json:"delta_to"`
		Unit                  string `json:"unit"`
		SeriesEpoch           string `json:"series_epoch"`
		MeasurementID         string `json:"measurement_id"`
		DefinitionRevision    string `json:"definition_revision"`
		SourceBindingRevision string `json:"source_binding_revision"`
		ProfileRevision       string `json:"profile_revision"`
		BaselineSampleID      string `json:"baseline_sample_id"`
		EndSampleID           string `json:"end_sample_id"`
	} `json:"counter"`
}

func entryOf(t *testing.T, row EncodedRow) (entry wideCounterEntry, provenance string) {
	t.Helper()
	provenance = cellMap(row)["provenance"].(string)
	var entries []wideCounterEntry
	require.NoError(t, json.Unmarshal([]byte(provenance), &entries), provenance)
	require.Len(t, entries, 1)
	return entries[0], provenance
}

// TestWideCounterRowKeepsRawInTheColumnAndTheDeltaInProvenance pins the typed wide
// counter encoding on the real exact engine: the value column holds the source's
// own raw digits, and provenance.counter holds the status, the exact delta as
// digits, the actual readings it spans and the policy evidence, all as strings.
func TestWideCounterRowKeepsRawInTheColumnAndTheDeltaInProvenance(t *testing.T) {
	member := counterMemberFor(measurement.ExactDecimal)
	layout := wideCounterLayout(t, member, schema.DataTypeFloat64)
	outcomes := counterOutcomes(t, member, exactDecimal(t, "12000.0"), exactDecimal(t, "12000.4"), exactDecimal(t, "12000.7"))
	rows := make([]EncodedRow, len(outcomes))
	for i, outcome := range outcomes {
		var err error
		rows[i], err = layout.EncodeRow(outcome)
		require.NoError(t, err)
	}

	first, firstJSON := entryOf(t, rows[0])
	require.Equal(t, "12000.0", cellMap(rows[0])["v_mean"], "the source's own digits are the stored value")
	require.Equal(t, [3]string{"baseline", "", ""}, [3]string{first.Counter.Status, first.Counter.Delta, first.Counter.DeltaFrom},
		"a baseline has no difference: absent, never zero")
	require.Contains(t, firstJSON, `"reason":"first-reading"`)

	second, secondJSON := entryOf(t, rows[1])
	baseline := counterWindowStart.Add(9*time.Minute + 123456*time.Microsecond)
	require.Equal(t, "12000.4", cellMap(rows[1])["v_mean"])
	require.Equal(t, [3]string{"ok", "0.4", "12000.4"}, [3]string{second.Counter.Status, second.Counter.Delta, second.Counter.RawText})
	require.Equal(t, [2]string{baseline.Format(longTimeText), baseline.Add(counterWindowLength).Format(longTimeText)},
		[2]string{second.Counter.DeltaFrom, second.Counter.DeltaTo}, "the delta spans the real readings, not the nominal window")
	require.Equal(t, [4]string{"counter", "exact", "kWh", "epoch-1"}, [4]string{second.Counter.Kind, second.Counter.Basis, second.Counter.Unit, second.Counter.SeriesEpoch})
	require.Equal(t, [4]string{"meas-kwh", "def-1", "bind-1", "r1"},
		[4]string{second.Counter.MeasurementID, second.Counter.DefinitionRevision, second.Counter.SourceBindingRevision, second.Counter.ProfileRevision})
	require.Equal(t, [2]string{"s-0", "s-1"}, [2]string{second.Counter.BaselineSampleID, second.Counter.EndSampleID})
	require.Contains(t, secondJSON, `"delta":"0.4"`)
	require.NotRegexp(t, `"delta":[-0-9]`, secondJSON, "no exact value passes through a JSON number")
}

func TestWideCounterRowKeepsExactIntegersAboveTwoToThe53(t *testing.T) {
	member := counterMemberFor(measurement.ExactUint64)
	layout := wideCounterLayout(t, member, schema.DataTypeUint64)
	outcomes := counterOutcomes(t, member, exactUint(9007199254740992), exactUint(9007199254740993))
	rows := make([]EncodedRow, 0, len(outcomes))
	for _, outcome := range outcomes {
		row, err := layout.EncodeRow(outcome)
		require.NoError(t, err)
		rows = append(rows, row)
	}
	require.Equal(t, "9007199254740993", cellMap(rows[1])["v_mean"], "the raw reading is exact digits, never a float")
	entry, provenance := entryOf(t, rows[1])
	require.Equal(t, "1", entry.Counter.Delta)
	require.Contains(t, provenance, `"delta":"1"`)
}

// TestWideCounterRowNeverInventsAZeroOrLosesADelta pins the NULL and refusal
// rules: a decrease the engine cannot prove has no difference (and is not
// absolute-valued), a silent window has no counter object and a NULL cell, and a
// counter result for a member without a binding is refused rather than dropped.
func TestWideCounterRowNeverInventsAZeroOrLosesADelta(t *testing.T) {
	member := counterMemberFor(measurement.ExactUint64)
	layout := wideCounterLayout(t, member, schema.DataTypeUint64)

	reset := counterOutcomes(t, member, exactUint(100), exactUint(40))[1]
	row, err := layout.EncodeRow(reset)
	require.NoError(t, err)
	entry, provenance := entryOf(t, row)
	require.Equal(t, "40", cellMap(row)["v_mean"], "the raw reading is always kept")
	require.Equal(t, [2]string{"reset_or_uncertain", ""}, [2]string{entry.Counter.Status, entry.Counter.Delta})
	require.NotContains(t, provenance, `"delta"`)

	silent := counterOutcomes(t, member, exactUint(100), nil)[1]
	row, err = layout.EncodeRow(silent)
	require.NoError(t, err)
	entry, _ = entryOf(t, row)
	require.Nil(t, cellMap(row)["v_mean"], "a silent window is NULL")
	require.Nil(t, entry.Counter, "and has no counter object")
	require.Equal(t, "missing", entry.Status)

	unbound := wideCounterLayout(t, LongRowMember{MemberKey: member.MemberKey, Type: member.Type}, schema.DataTypeUint64)
	_, err = unbound.EncodeRow(counterOutcomes(t, member, exactUint(100))[0])
	require.Equal(t, rowErrCounterUnbound, rowCode(t, err), "a cumulative result is never silently written without its delta")
}
