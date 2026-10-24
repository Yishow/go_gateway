package groupdelivery

import (
	"encoding/json"
	"math"
	"regexp"
	"time"

	"go-gateway/internal/datalink/dbtarget"
)

// sourceProvenance is the part of a record's provenance document this projection reads. Nothing else in the
// document, which may carry raw diagnostics, is ever returned.
type sourceProvenance struct {
	Type   string `json:"type"`
	Member string `json:"member"`
	Reason string `json:"reason"`
	// Window is present only for a window a forward cutover cut short.
	Window *struct {
		Truncated    string `json:"truncated"`
		Start        string `json:"start"`
		End          string `json:"end"`
		NominalStart string `json:"nominal_start"`
		NominalEnd   string `json:"nominal_end"`
	} `json:"window"`
}

func parseSourceProvenance(document string) (sourceProvenance, bool) {
	var provenance sourceProvenance
	if document == "" || json.Unmarshal([]byte(document), &provenance) != nil {
		return sourceProvenance{}, false
	}
	return provenance, true
}

// parseLongTime reads a time in either form a long table time is stored in.
func parseLongTime(text string) (time.Time, bool) {
	for _, layout := range []string{sourceTimeText, time.RFC3339Nano} {
		if parsed, err := time.Parse(layout, text); err == nil {
			return parsed.UTC(), true
		}
	}
	return time.Time{}, false
}

// sourceWindow returns the actual window the record covers. The outbox keeps the nominal bucket start while the
// payload and the table keep the real bounds, so the two differ only when a forward cutover cut the window short,
// and then only when the provenance proves exactly which side was cut and where the nominal bucket was.
func sourceWindow(cells map[string]any, encoded dbtarget.EncodedRow, nominalStart time.Time, provenance sourceProvenance) (start, end time.Time, ok bool) {
	start, startOK := timeCell(cells, "bucket_start")
	end, endOK := timeCell(cells, "bucket_end")
	nominalStart = nominalStart.Truncate(time.Microsecond)
	if !startOK || !endOK || !end.After(start) || !encoded.BucketStart.UTC().Truncate(time.Microsecond).Equal(start) {
		return start, end, false
	}
	cut := provenance.Window
	if cut == nil {
		return start, end, start.Equal(nominalStart)
	}
	actualStart, ok1 := parseLongTime(cut.Start)
	actualEnd, ok2 := parseLongTime(cut.End)
	wantStart, ok3 := parseLongTime(cut.NominalStart)
	wantEnd, ok4 := parseLongTime(cut.NominalEnd)
	if !ok1 || !ok2 || !ok3 || !ok4 || !actualStart.Equal(start) || !actualEnd.Equal(end) || !wantStart.Equal(nominalStart) || !wantEnd.After(wantStart) {
		return start, end, false
	}
	switch cut.Truncated {
	case "start":
		ok = start.After(wantStart) && end.Equal(wantEnd)
	case "end":
		ok = start.Equal(wantStart) && end.Before(wantEnd)
	case "both":
		ok = start.After(wantStart) && end.Before(wantEnd)
	}
	return start, end, ok
}

// sourceSlotByType is the one value column each exact type is stored in.
var sourceSlotByType = map[string]string{
	"int64": "value_int", "uint64": sourceSlotDecimal, "decimal": sourceSlotDecimal, "float64": "value_float", "bool": "value_bool", "text": "value_text",
}

var sourceNumericTypes = map[string]bool{"int64": true, "uint64": true, "decimal": true, "float64": true}

var sourceDecimalText = regexp.MustCompile(`^-?\d+(\.\d+)?$`)

// sourceNumeric reports whether a stored difference is a number: digits, an integer or a finite float.
func sourceNumeric(value any) bool {
	switch v := value.(type) {
	case string:
		return sourceDecimalText.MatchString(v)
	case int64:
		return true
	case float64:
		return !math.IsNaN(v) && !math.IsInf(v, 0)
	}
	return false
}

// sourceValues checks that the stored result says one consistent thing and reads what it says: which value slot
// matches the value type, whether a raw reading is present, and whether a difference was stored. It recomputes
// nothing and consults no destination. A contradiction makes the whole record unavailable.
func sourceValues(cells map[string]any, status SourceStatus, profileKind string, provenance sourceProvenance) (rawPresent, deltaKnown bool, cause SourceUnavailableCause, ok bool) {
	valueType := cellText(cells, "value_type")
	slot, known := sourceSlotByType[valueType]
	if !known || (provenance.Type != "" && provenance.Type != valueType) || (profileKind == sourceProfileCumulative && !sourceNumericTypes[valueType]) {
		return false, false, SourceCausePayloadInvalid, false
	}
	for _, name := range sourceValueSlots {
		if cells[name] != nil && name != slot {
			return false, false, SourceCausePayloadInvalid, false
		}
	}
	rawPresent = cells[slot] != nil
	deltaKnown = cells["delta_value"] != nil
	if (cells["delta_from"] != nil) != deltaKnown || (cells["delta_to"] != nil) != deltaKnown || (deltaKnown && !sourceNumeric(cells["delta_value"])) {
		return false, false, SourceCausePayloadInvalid, false
	}
	storesNothing := status == SourceStatusMissing || status == SourceStatusBad || status == SourceStatusStale || status == SourceStatusInvalid
	storesReading := status == SourceStatusOK || status == SourceStatusPartial || status == SourceStatusBaseline
	noDifference := storesNothing || status == SourceStatusBaseline || status == SourceStatusUncertain
	if (storesNothing && rawPresent) || (storesReading && !rawPresent) || (deltaKnown && (profileKind != sourceProfileCumulative || noDifference)) {
		return false, false, SourceCauseStatusUnknown, false
	}
	return rawPresent, deltaKnown, "", true
}
