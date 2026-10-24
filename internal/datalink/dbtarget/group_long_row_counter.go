package dbtarget

import (
	"database/sql/driver"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
)

// rowErrCounterUnbound refuses a cumulative result whose member has no confirmed
// counter binding, so a delta is never silently dropped.
const rowErrCounterUnbound = "counter-unbound"

// longCounterProvenance is the cumulative part of a record's provenance: what
// the delta rests on and which readings and revisions it compares.
type longCounterProvenance struct {
	Kind  string `json:"kind"`
	Basis string `json:"basis,omitempty"`
	// RawText is the decimal text the source gave, kept so a numeric column
	// that normalizes digits never erases the original precision evidence.
	RawText               string `json:"raw_text,omitempty"`
	BaselineSampleID      string `json:"baseline_sample_id,omitempty"`
	EndSampleID           string `json:"end_sample_id,omitempty"`
	DefinitionRevision    string `json:"definition_revision"`
	SourceBindingRevision string `json:"source_binding_revision"`
	ProfileRevision       string `json:"profile_revision"`
	// SpanTotal marks a delta that crosses missing windows: a total over its real
	// interval, not the usage of one window.
	SpanTotal bool `json:"span_total,omitempty"`
	// RolloverModulus and MaxRatePerSecond are the confirmed evidence a wrap
	// compensation rests on.
	RolloverModulus  string `json:"rollover_modulus,omitempty"`
	MaxRatePerSecond string `json:"max_rate_per_second,omitempty"`
	// Events are the discontinuities and compensations of the window, bounded;
	// EventsOmitted counts the ones the bound dropped.
	Events        []longCounterEvent `json:"events,omitempty"`
	EventsOmitted int                `json:"events_omitted,omitempty"`
}

// longCounterEvent is one aggregation.CounterEvent with its exact readings as
// digits, so a value above 2^53 never passes through a JSON number.
type longCounterEvent struct {
	Kind         string `json:"kind"`
	Reason       string `json:"reason,omitempty"`
	FromSampleID string `json:"from_sample_id"`
	ToSampleID   string `json:"to_sample_id"`
	From         string `json:"from"`
	To           string `json:"to"`
	FromValue    string `json:"from_value"`
	ToValue      string `json:"to_value"`
}

// counterCells writes the cells only a cumulative record has: the series epoch
// and unit of its confirmed binding, and the exact difference with the actual
// observation times it spans. A record without a trustworthy difference leaves
// delta_value, delta_from and delta_to NULL; it is never written as zero.
func (l *LongRowLayout) counterCells(values map[string]driver.Value, member LongRowMember, result snapshot.MemberResult) {
	if member.Counter == nil {
		return
	}
	values[longColSeriesEpoch] = l.text(longColSeriesEpoch, member.Counter.SeriesEpoch)
	values[longColUnit] = l.text(longColUnit, member.Counter.Unit)
	if result.Counter == nil || result.Counter.Delta == nil {
		return
	}
	delta := result.Counter.Delta
	digits, err := deltaDecimal(delta.Value)
	if err != nil {
		values[longColDeltaValue] = errCell{column: longColDeltaValue}
		return
	}
	bound, err := EncodeExactValue(l.spec.Dialect, l.columns[longColDeltaValue].DataType, digits)
	if err != nil {
		values[longColDeltaValue] = errCell{column: longColDeltaValue}
		return
	}
	values[longColDeltaValue] = bound
	values[longColDeltaFrom] = l.instant(longColDeltaFrom, delta.From)
	values[longColDeltaTo] = l.instant(longColDeltaTo, delta.To)
}

// deltaDecimal is the decimal form a delta is stored in. An exact delta already
// is one; a float delta is rendered with its own shortest digits and its basis
// stays recorded in provenance, so it is not mistaken for an exact measurement.
func deltaDecimal(value measurement.ExactValue) (measurement.ExactValue, error) {
	if value.Type() == measurement.ExactFloat64 {
		if f, ok := value.Value().(float64); ok {
			return measurement.FloatDecimal(f)
		}
	}
	return value, nil
}

func counterProvenance(member LongRowMember, counter *aggregation.CounterWindowResult) *longCounterProvenance {
	if member.Counter == nil {
		return nil
	}
	entry := &longCounterProvenance{
		Kind: string(member.Counter.Kind), DefinitionRevision: member.Counter.DefinitionRevision,
		SourceBindingRevision: member.Counter.SourceBindingRevision, ProfileRevision: member.Counter.ProfileRevision,
	}
	if counter.Raw != nil && counter.Raw.Value.Type() == measurement.ExactDecimal {
		entry.RawText, _ = counter.Raw.Value.DecimalDigits()
	}
	if counter.Delta != nil {
		entry.Basis = string(counter.Delta.Basis)
		entry.BaselineSampleID, entry.EndSampleID = counter.Delta.BaselineSampleID, counter.Delta.EndSampleID
		entry.SpanTotal = counter.Delta.Span
	}
	if member.Counter.RolloverModulus != nil {
		entry.RolloverModulus, _ = member.Counter.RolloverModulus.CanonicalDigits()
	}
	if member.Counter.MaxRatePerSecond != nil {
		entry.MaxRatePerSecond, _ = member.Counter.MaxRatePerSecond.CanonicalDigits()
	}
	for _, event := range counter.Events {
		entry.Events = append(entry.Events, longCounterEvent{
			Kind: string(event.Kind), Reason: event.Reason, FromSampleID: event.FromSampleID, ToSampleID: event.ToSampleID,
			From: eventTime(event.From), To: eventTime(event.To),
			FromValue: eventValue(event.FromValue), ToValue: eventValue(event.ToValue),
		})
	}
	entry.EventsOmitted = counter.EventsOmitted
	return entry
}

// eventTime is an observation time in the table's UTC microsecond form.
func eventTime(at time.Time) string { return at.UTC().Truncate(time.Microsecond).Format(longTimeText) }

// eventValue is the exact text of a reading: its own digits for integers and
// decimals, the float's shortest digits for a float source.
func eventValue(value measurement.ExactValue) string {
	if digits, ok := value.DecimalDigits(); ok {
		return digits
	}
	if f, ok := value.Value().(float64); ok {
		if decimal, err := measurement.FloatDecimal(f); err == nil {
			digits, _ := decimal.DecimalDigits()
			return digits
		}
	}
	return ""
}
