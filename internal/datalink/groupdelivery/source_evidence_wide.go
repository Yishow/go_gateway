package groupdelivery

import (
	"encoding/json"
	"strings"
	"time"

	"go-gateway/internal/datalink/dbtarget"
)

// sourceWideEntry is the member level evidence of one wide profile row, the part of it this projection reads.
type sourceWideEntry struct {
	Member     string      `json:"member"`
	Status     string      `json:"status"`
	Reason     string      `json:"reason"`
	Role       string      `json:"role"`
	SourceType string      `json:"source_type"`
	ResultType string      `json:"result_type"`
	Mean       *sourceMean `json:"mean"`
	Counter    *struct {
		Status    string `json:"status"`
		Reason    string `json:"reason"`
		Delta     string `json:"delta"`
		DeltaFrom string `json:"delta_from"`
		DeltaTo   string `json:"delta_to"`
	} `json:"counter"`
	// Window is present only for a window a forward cutover cut short.
	Window *sourceWindowCut `json:"window"`
}

func parseWideProvenance(document string) ([]sourceWideEntry, bool) {
	var entries []sourceWideEntry
	if document == "" || json.Unmarshal([]byte(document), &entries) != nil {
		return nil, false
	}
	return entries, true
}

// verifyWideRecord keeps what one wide profile row proves about one of its sources. The row is the whole window of
// the group, so the source is found by its member key in the frozen provenance and its reading by its frozen column.
// Everything is read from the frozen payload with the columns and window the applied layout froze; nothing is queried.
func verifyWideRecord(scope SourceScope, ref SourceRef, row sourceRow, encoded dbtarget.EncodedRow, cells map[string]any, delivery SourceDelivery) SourceRecord {
	column := strings.ToLower
	if cellText(cells, column(scope.RecordKeyColumn)) != row.recordID || cellText(cells, column(scope.GroupIDColumn)) != scope.GroupID {
		return unavailable(delivery, SourceCauseIdentityMismatch)
	}
	entries, documented := parseWideProvenance(cellText(cells, column(scope.ProvenanceColumn)))
	if !documented {
		return unavailable(delivery, SourceCausePayloadInvalid)
	}
	member := SourceEntityKey(ref.DeviceID, ref.PointID, ref.TagID)
	var entry *sourceWideEntry
	for i := range entries {
		if entries[i].Member == member {
			entry = &entries[i]
			break
		}
	}
	if entry == nil || (ref.Role != "" && entry.Role != "" && entry.Role != ref.Role) {
		return unavailable(delivery, SourceCauseIdentityMismatch)
	}
	nominalStart, frozenOK := frozenTime(row.bucketStart)
	start, startOK := timeCell(cells, column(scope.BucketStartColumn))
	if !frozenOK || !startOK || scope.IntervalSeconds <= 0 {
		return unavailable(delivery, SourceCauseWindowInvalid)
	}
	// A whole window is the aligned interval. A cut window names its real end, and its aligned bucket has to be the
	// interval the applied layout froze, so the same rule that judges a long window judges this one.
	interval := time.Duration(scope.IntervalSeconds) * time.Second
	end := start.Add(interval)
	if entry.Window != nil {
		cutEnd, cutOK := parseLongTime(entry.Window.End)
		wantEnd, wantOK := parseLongTime(entry.Window.NominalEnd)
		if !cutOK || !wantOK || !wantEnd.Equal(nominalStart.Truncate(time.Microsecond).Add(interval)) {
			return unavailable(delivery, SourceCauseWindowInvalid)
		}
		end = cutEnd
	}
	// A wide row keeps the aligned bucket as its identity and the real window start in its bucket cell, so the row's
	// own start is the aligned one and the cell start is what the shared window rule judges.
	start, end, windowOK := resolveSourceWindow(start, end, start, nominalStart, entry.Window)
	if !windowOK || !encoded.BucketStart.UTC().Truncate(time.Microsecond).Equal(nominalStart.Truncate(time.Microsecond)) {
		return unavailable(delivery, SourceCauseWindowInvalid)
	}

	statusText, reason := entry.Status, entry.Reason
	if entry.Counter != nil {
		statusText = entry.Counter.Status
		if entry.Counter.Reason != "" {
			reason = entry.Counter.Reason
		}
	}
	status := SourceStatus(statusText)
	if !sourceStatuses[status] {
		return unavailable(delivery, SourceCauseStatusUnknown)
	}
	cumulative := scope.ProfileKind == sourceProfileCumulative
	rawPresent := ref.Column != "" && cells[column(ref.Column)] != nil
	deltaKnown := entry.Counter != nil && entry.Counter.Delta != ""
	if deltaKnown && !wideDifferenceWhole(entry.Counter.Delta, entry.Counter.DeltaFrom, entry.Counter.DeltaTo) ||
		entry.Counter != nil && !deltaKnown && (entry.Counter.DeltaFrom != "" || entry.Counter.DeltaTo != "") {
		return unavailable(delivery, SourceCausePayloadInvalid)
	}
	// A member that is not a mean states its frozen types itself; they have to be the ones the applied layout froze.
	if entry.Mean == nil && ((ref.SourceType != "" && entry.SourceType != ref.SourceType) || (ref.ResultType != "" && entry.ResultType != ref.ResultType)) {
		return unavailable(delivery, SourceCauseIdentityMismatch)
	}
	// A counted member carries its counter window unless it had no usable reading; any other member has none.
	if (entry.Counter != nil) != cumulative && (!cumulative || entry.Counter != nil || !sourceStoresNothing(status)) {
		return unavailable(delivery, SourceCausePayloadInvalid)
	}
	if cause, consistent := sourceStatusConsistent(status, scope.ProfileKind, rawPresent, deltaKnown); !consistent {
		return unavailable(delivery, cause)
	}
	evidence, verified := verifySourceAggregation(ref, status, entry.Mean)
	if !verified {
		return unavailable(delivery, SourceCausePayloadInvalid)
	}
	record := &VerifiedSourceRecord{
		GroupRevision: row.groupRevision, ConnectorRevision: row.connectorRevision, RecordID: row.recordID, EffectKey: row.effectKey,
		PayloadDigest: row.digest, ProfileKind: scope.ProfileKind, WindowStart: start, WindowEnd: end, Status: status,
		Reason: sourceReason(reason), RawPresent: rawPresent, Aggregation: evidence,
	}
	if cumulative {
		record.DeltaKnown = &deltaKnown
	}
	if delivery == SourceDeliveryCommitted {
		committed, ok := frozenTime(row.committedAt)
		if !ok {
			return unavailable(SourceDeliveryUnknown, SourceCausePayloadInvalid)
		}
		record.CommittedAt = &committed
	}
	return SourceRecord{Available: true, Delivery: delivery, VerifiedSourceRecord: record}
}

// wideDifferenceWhole reports whether a stored difference is a number with both of its bounds in order. A difference
// without its bounds, or with bounds that do not span time, claims nothing.
func wideDifferenceWhole(delta, from, to string) bool {
	start, startOK := parseLongTime(from)
	end, endOK := parseLongTime(to)
	return sourceNumeric(delta) && startOK && endOK && end.After(start)
}
