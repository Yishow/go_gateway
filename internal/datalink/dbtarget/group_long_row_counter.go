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
	// Raw is mixed-only typed evidence; legacy JSON remains unchanged.
	Raw           *measurement.ExactValue `json:"raw,omitempty"`
	RawSampleID   string                  `json:"raw_sample_id,omitempty"`
	RawObservedAt string                  `json:"raw_observed_at,omitempty"`
	Kind          string                  `json:"kind"`
	Basis         string                  `json:"basis,omitempty"`
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

	// The fields below are set only by a wide row, which has no delta, status,
	// unit or series columns of its own: its provenance keeps the window's status,
	// the typed exact difference with the actual readings it spans, and the
	// measurement evidence they rest on. A long record leaves them empty and keeps
	// those facts in its columns, so its JSON is unchanged.
	Status        string `json:"status,omitempty"`
	Reason        string `json:"reason,omitempty"`
	Delta         string `json:"delta,omitempty"`
	DeltaFrom     string `json:"delta_from,omitempty"`
	DeltaTo       string `json:"delta_to,omitempty"`
	Unit          string `json:"unit,omitempty"`
	SeriesEpoch   string `json:"series_epoch,omitempty"`
	MeasurementID string `json:"measurement_id,omitempty"`
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
	if l.spec.ResultShape == MixedResultShapeV1 {
		return
	}
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

func counterProvenance(binding *aggregation.CounterBinding, counter *aggregation.CounterWindowResult) *longCounterProvenance {
	if binding == nil {
		return nil
	}
	entry := &longCounterProvenance{
		Kind: string(binding.Kind), DefinitionRevision: binding.DefinitionRevision,
		SourceBindingRevision: binding.SourceBindingRevision, ProfileRevision: binding.ProfileRevision,
	}
	if counter.Raw != nil && counter.Raw.Value.Type() == measurement.ExactDecimal {
		entry.RawText, _ = counter.Raw.Value.DecimalDigits()
	}
	if counter.Delta != nil {
		entry.Basis = string(counter.Delta.Basis)
		entry.BaselineSampleID, entry.EndSampleID = counter.Delta.BaselineSampleID, counter.Delta.EndSampleID
		entry.SpanTotal = counter.Delta.Span
	}
	if binding.RolloverModulus != nil {
		entry.RolloverModulus, _ = binding.RolloverModulus.CanonicalDigits()
	}
	if binding.MaxRatePerSecond != nil {
		entry.MaxRatePerSecond, _ = binding.MaxRatePerSecond.CanonicalDigits()
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

// wideCounterProvenance is the counter object of one member of a wide row: the
// shared cumulative evidence of counterProvenance plus the status, the typed exact
// difference as digits, its actual bounds and the measurement evidence. A window
// without a trustworthy difference leaves the delta and its bounds out; it never
// writes a zero.
func wideCounterProvenance(binding *aggregation.CounterBinding, counter *aggregation.CounterWindowResult) *longCounterProvenance {
	if counter == nil {
		return nil
	}
	entry := counterProvenance(binding, counter)
	if entry == nil {
		return nil
	}
	entry.Status, entry.Reason = string(counter.Status), counter.Reason
	entry.Unit, entry.SeriesEpoch, entry.MeasurementID = binding.Unit, binding.SeriesEpoch, binding.MeasurementID
	if counter.Delta != nil {
		if digits, err := deltaDecimal(counter.Delta.Value); err == nil {
			entry.Delta, _ = digits.DecimalDigits()
			entry.DeltaFrom, entry.DeltaTo = eventTime(counter.Delta.From), eventTime(counter.Delta.To)
		}
	}
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
