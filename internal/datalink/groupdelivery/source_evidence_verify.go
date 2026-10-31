package groupdelivery

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-gateway/internal/datalink/aggregation"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/snapshot"
)

// sourceTimeText is the SQLite form of a long table time.
const sourceTimeText = "2006-01-02T15:04:05.000000Z"

var sourceStatuses = map[SourceStatus]bool{
	SourceStatusOK: true, SourceStatusPartial: true, SourceStatusBaseline: true, SourceStatusMissing: true, SourceStatusBad: true,
	SourceStatusStale: true, SourceStatusInvalid: true, SourceStatusUncertain: true, SourceStatusGap: true,
}

// sourceReasons are the typed reasons a long record's provenance can carry. Anything else, such as free text a
// collector attached to a bad sample, is shown as SourceReasonUnavailable.
var sourceReasons = map[string]bool{
	snapshot.ReasonNoSamples: true, snapshot.ReasonTypeMismatch: true, snapshot.ReasonCounterUnavailable: true,
	"missing-value": true, "max-age-exceeded": true, "mapping-failed": true,
	"quality-bad": true, "quality-missing": true, "quality-stale": true, "quality-invalid": true, "quality-uncertain": true,
	snapshot.ReasonCoverageUnverified: true, snapshot.MeanReasonBad: true, snapshot.MeanReasonInvalid: true,
	snapshot.MeanReasonStale: true, snapshot.MeanReasonNonNumeric: true,
	aggregation.ReasonFirstReading: true, aggregation.ReasonBindingChanged: true, aggregation.ReasonNoEligibleReading: true,
	aggregation.ReasonDecreaseUnproven: true, aggregation.ReasonContradictory: true, aggregation.ReasonNonMonotonicTime: true,
	aggregation.ReasonMixedBasis: true, aggregation.ReasonRolloverEvidence: true, aggregation.ReasonRolloverRange: true,
	aggregation.ReasonRolloverMultiWrap: true, aggregation.ReasonRolloverRate: true, aggregation.ReasonNonNumericReading: true,
	aggregation.ReasonDeltaNotRepresented: true,
}

var sourceValueSlots = []string{"value_int", sourceSlotDecimal, "value_float", "value_bool", "value_text"}

func sourceReason(reason string) string {
	switch {
	case reason == "":
		return ""
	case sourceReasons[reason]:
		return reason
	default:
		return SourceReasonUnavailable
	}
}

// sourceDelivery maps an outbox state to a local stage. A SQL committed row without a matching receipt is unknown.
func sourceDelivery(state string, receiptMatched bool) SourceDelivery {
	switch state {
	case StatePending, StateSending:
		return SourceDeliveryQueued
	case StateRetrying:
		return SourceDeliveryRetrying
	case StateBlocked:
		return SourceDeliveryBlocked
	case StateQuarantined:
		return SourceDeliveryQuarantined
	case StateSkipped:
		return SourceDeliverySkipped
	case StateCommitted:
		if receiptMatched {
			return SourceDeliveryCommitted
		}
	}
	return SourceDeliveryUnknown
}

func (s *Store) receiptMatches(ctx context.Context, row sourceRow) (bool, error) {
	if row.state != StateCommitted || row.committedAt == "" {
		return false, nil
	}
	var digest, committed string
	err := s.db.QueryRowContext(ctx, `SELECT payload_digest, committed_at FROM wg_delivery_receipts WHERE effect_key = ?`, row.effectKey).Scan(&digest, &committed)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read source receipt: %w", err)
	}
	return digest == row.digest && committed == row.committedAt, nil
}

func unavailable(delivery SourceDelivery, cause SourceUnavailableCause) SourceRecord {
	return SourceRecord{Available: false, Delivery: delivery, Cause: cause}
}

// cellText reads a text cell; a NULL or non text cell is empty.
func cellText(cells map[string]any, column string) string {
	if text, ok := cells[column].(string); ok {
		return text
	}
	return ""
}

// decodeSourcePayload reports whether a stored payload is a readable row.
func decodeSourcePayload(payload string) (dbtarget.EncodedRow, bool) {
	row, err := DecodeRowPayload([]byte(payload))
	return row, err == nil
}

// frozenTime parses a time the outbox stored in its fixed width form.
func frozenTime(text string) (time.Time, bool) {
	parsed, err := time.Parse(timestampLayout, text)
	return parsed.UTC(), err == nil
}

// cellTime reads a long table time in either dialect form.
func timeCell(cells map[string]any, column string) (time.Time, bool) {
	switch v := cells[column].(type) {
	case time.Time:
		return v.UTC(), true
	case string:
		for _, layout := range []string{sourceTimeText, time.RFC3339Nano} {
			if parsed, err := time.Parse(layout, v); err == nil {
				return parsed.UTC(), true
			}
		}
	}
	return time.Time{}, false
}

// openFrozenRow checks what every managed record must satisfy before anything in it is read: the payload still has
// the digest the outbox froze, it decodes, and its identity, applied revision and destination are the scope's.
func openFrozenRow(scope SourceScope, row sourceRow) (dbtarget.EncodedRow, map[string]any, SourceUnavailableCause, bool) {
	sum := sha256.Sum256([]byte(row.payload))
	if hex.EncodeToString(sum[:]) != row.digest {
		return dbtarget.EncodedRow{}, nil, SourceCauseDigestMismatch, false
	}
	encoded, readable := decodeSourcePayload(row.payload)
	if !readable {
		return dbtarget.EncodedRow{}, nil, SourceCausePayloadInvalid, false
	}
	if encoded.RecordID != row.recordID || encoded.EffectKey != row.effectKey || encoded.EntityKey != row.entityKey ||
		row.workspaceID != scope.WorkspaceID || row.groupID != scope.GroupID || row.groupRevision != scope.GroupRevision ||
		row.connectorID != scope.ConnectorID || row.connectorRevision != scope.ConnectorRevision {
		return dbtarget.EncodedRow{}, nil, SourceCauseIdentityMismatch, false
	}
	cells := make(map[string]any, len(encoded.Cells))
	for _, cell := range encoded.Cells {
		cells[strings.ToLower(cell.Column)] = cell.Value
	}
	return encoded, cells, "", true
}

// verifySource decodes one frozen payload and keeps only what it proves about the source. Nothing about a record
// that fails any check is reported, and a committed row needs a receipt that matches its payload.
func (s *Store) verifySource(ctx context.Context, scope SourceScope, ref SourceRef, row sourceRow) (SourceRecord, error) {
	matched, err := s.receiptMatches(ctx, row)
	if err != nil {
		return SourceRecord{}, err
	}
	delivery := sourceDelivery(row.state, matched)
	encoded, cells, cause, opened := openFrozenRow(scope, row)
	if !opened {
		return unavailable(delivery, cause), nil
	}
	// Mixed stores a selected calculation result, never legacy raw evidence.
	// Its explicit result/provenance projection is a separate DTO capability.
	if !encoded.DueAt.IsZero() {
		return unavailable(delivery, SourceCauseCalculationUnsupported), nil
	}
	if scope.Wide {
		return verifyWideRecord(scope, ref, row, encoded, cells, delivery), nil
	}
	if cellText(cells, "record_id") != row.recordID || cellText(cells, "workspace_id") != scope.WorkspaceID ||
		cellText(cells, "group_id") != scope.GroupID || cellText(cells, "group_revision") != scope.GroupRevision ||
		cellText(cells, "device_id") != ref.DeviceID || cellText(cells, "point_id") != ref.PointID ||
		cellText(cells, "tag_id") != ref.TagID || cellText(cells, "circuit_id") != ref.CircuitID ||
		cellText(cells, "profile_kind") != scope.ProfileKind {
		return unavailable(delivery, SourceCauseIdentityMismatch), nil
	}
	status := SourceStatus(cellText(cells, "status"))
	if !sourceStatuses[status] {
		return unavailable(delivery, SourceCauseStatusUnknown), nil
	}
	provenance, documented := parseSourceProvenance(cellText(cells, "provenance"))
	if !documented {
		return unavailable(delivery, SourceCausePayloadInvalid), nil
	}
	if provenance.Member != SourceEntityKey(ref.DeviceID, ref.PointID, ref.TagID) {
		return unavailable(delivery, SourceCauseIdentityMismatch), nil
	}
	nominalStart, frozenOK := frozenTime(row.bucketStart)
	start, end, windowOK := sourceWindow(cells, encoded, nominalStart, provenance)
	if !frozenOK || !windowOK {
		return unavailable(delivery, SourceCauseWindowInvalid), nil
	}
	rawPresent, deltaKnown, cause, consistent := sourceValues(cells, status, scope.ProfileKind, provenance)
	if !consistent {
		return unavailable(delivery, cause), nil
	}
	evidence, verified := verifySourceAggregation(ref, status, provenance.Mean)
	if !verified {
		return unavailable(delivery, SourceCausePayloadInvalid), nil
	}
	record := &VerifiedSourceRecord{
		GroupRevision: row.groupRevision, ConnectorRevision: row.connectorRevision, RecordID: row.recordID, EffectKey: row.effectKey,
		PayloadDigest: row.digest, ProfileKind: scope.ProfileKind, WindowStart: start, WindowEnd: end, Status: status,
		Reason: sourceReason(provenance.Reason), RawPresent: rawPresent, Aggregation: evidence,
	}
	if scope.ProfileKind == sourceProfileCumulative {
		record.DeltaKnown = &deltaKnown
	}
	if delivery == SourceDeliveryCommitted {
		committed, ok := frozenTime(row.committedAt)
		if !ok {
			return unavailable(SourceDeliveryUnknown, SourceCausePayloadInvalid), nil
		}
		record.CommittedAt = &committed
	}
	return SourceRecord{Available: true, Delivery: delivery, VerifiedSourceRecord: record}, nil
}
