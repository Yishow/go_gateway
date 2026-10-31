package groupdelivery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// SourceStatus is what a long record stored for its source.
type SourceStatus string

// Source statuses a verified record can carry; any other stored value makes the record unavailable.
const (
	SourceStatusOK        SourceStatus = "ok"
	SourceStatusPartial   SourceStatus = "partial"
	SourceStatusBaseline  SourceStatus = "baseline"
	SourceStatusMissing   SourceStatus = "missing"
	SourceStatusBad       SourceStatus = "bad"
	SourceStatusStale     SourceStatus = "stale"
	SourceStatusInvalid   SourceStatus = "invalid"
	SourceStatusUncertain SourceStatus = "reset_or_uncertain"
	SourceStatusGap       SourceStatus = "gap"
)

// SourceDelivery is the local stage of a record. Only SourceDeliveryCommitted is destination evidence, and
// only for a record whose receipt matched its payload.
type SourceDelivery string

// Local delivery stages of one record.
const (
	SourceDeliveryQueued      SourceDelivery = "queued"
	SourceDeliveryRetrying    SourceDelivery = "retrying"
	SourceDeliveryBlocked     SourceDelivery = "blocked"
	SourceDeliveryQuarantined SourceDelivery = "quarantined"
	SourceDeliveryUnknown     SourceDelivery = "unknown"
	SourceDeliveryCommitted   SourceDelivery = "sql_committed"
	SourceDeliverySkipped     SourceDelivery = "operator_skipped"
)

// SourceUnavailableCause says why a stored record could not be verified.
type SourceUnavailableCause string

// Causes of an unverifiable record. None of them is a statement about the destination.
const (
	SourceCausePayloadInvalid         SourceUnavailableCause = "payload_invalid"
	SourceCauseDigestMismatch         SourceUnavailableCause = "digest_mismatch"
	SourceCauseIdentityMismatch       SourceUnavailableCause = "identity_mismatch"
	SourceCauseStatusUnknown          SourceUnavailableCause = "status_unknown"
	SourceCauseWindowInvalid          SourceUnavailableCause = "window_invalid"
	SourceCauseCalculationUnsupported SourceUnavailableCause = "calculation_projection_unsupported"
)

// SourceReasonUnavailable replaces a stored reason that is not one of the stable codes.
const SourceReasonUnavailable = "reason_unavailable"

// Names the projection shares with the long row encoder: the cumulative profile and the decimal value column.
const (
	sourceProfileCumulative = "cumulative"
	sourceSlotDecimal       = "value_decimal"
)

// sourceCandidateLimit bounds how many of a source's newest SQL committed rows are inspected for a matching receipt.
const sourceCandidateLimit = 8

// SourceScope is the applied group revision whose records are read.
type SourceScope struct {
	WorkspaceID   string
	GroupID       string
	GroupRevision string
	// ProfileKind is the recording profile the applied group froze.
	ProfileKind string
	// ConnectorID and ConnectorRevision are the destination the applied revision froze. A record is proof only for
	// this destination; an editable draft's destination is never consulted.
	ConnectorID       string
	ConnectorRevision string
	// Wide marks a managed_wide_profile_v1 group, whose one row per window carries every source. The columns and
	// the window length are the ones the applied layout froze, so the payload is read the way it was written.
	Wide                                                                bool
	IntervalSeconds                                                     int
	RecordKeyColumn, GroupIDColumn, BucketStartColumn, ProvenanceColumn string
}

// SourceRef is one applied source of a managed group with what the applied layout froze about it. The recording
// facts are what a record's own evidence is checked against: a record that merely agrees with itself proves nothing.
type SourceRef struct {
	CircuitID string
	DeviceID  string
	PointID   string
	TagID     string
	// Column is the wide column that holds the source; Reduction, Role, SourceType and ResultType are the frozen
	// recording facts (empty for a group saved before they existed).
	Column     string
	Reduction  string
	Role       string
	SourceType string
	ResultType string
}

// VerifiedSourceRecord holds the value-free result of a record whose frozen payload was verified.
type VerifiedSourceRecord struct {
	GroupRevision     string       `json:"group_revision"`
	ConnectorRevision string       `json:"connector_revision"`
	RecordID          string       `json:"record_id"`
	EffectKey         string       `json:"effect_key"`
	PayloadDigest     string       `json:"payload_digest"`
	ProfileKind       string       `json:"profile_kind"`
	WindowStart       time.Time    `json:"window_start"`
	WindowEnd         time.Time    `json:"window_end"`
	Status            SourceStatus `json:"status"`
	Reason            string       `json:"reason,omitempty"`
	RawPresent        bool         `json:"raw_present"`
	// Aggregation is the verified evidence of a mean window: method, counts and coverage, never a value.
	Aggregation *SourceAggregation `json:"aggregation,omitempty"`
	// DeltaKnown is set for cumulative records only; false is an unknown difference, never zero.
	DeltaKnown *bool `json:"delta_known,omitzero"`
	// CommittedAt is set only when the destination receipt matched this payload.
	CommittedAt *time.Time `json:"committed_at,omitzero"`
}

// SourceRecord is one source's record, or the fact that it could not be verified.
type SourceRecord struct {
	Available bool                   `json:"available"`
	Delivery  SourceDelivery         `json:"delivery"`
	Cause     SourceUnavailableCause `json:"cause,omitempty"`
	*VerifiedSourceRecord
}

// SourceRecords is the newest local record of one source and the newest one the destination confirmed.
type SourceRecords struct {
	Latest           *SourceRecord
	LastSQLCommitted *SourceRecord
}

// SourceEntityKey is the outbox entity key of a long source: the JSON array the group boundary builds for it.
func SourceEntityKey(deviceID, pointID, tagID string) string {
	encoded, err := json.Marshal([]string{deviceID, pointID, tagID})
	if err != nil {
		return deviceID + "\x1f" + pointID + "\x1f" + tagID
	}
	return string(encoded)
}

// sourceRow is one outbox row read for verification; its payload is decoded once and never returned.
type sourceRow struct {
	effectKey, recordID, workspaceID, groupID, groupRevision, entityKey              string
	connectorID, connectorRevision, bucketStart, payload, digest, state, committedAt string
}

// SourceRecords reads, without values, the newest record of one applied source and the newest SQL confirmed one.
// At most one candidate is decoded for the former and a bounded few for the latter; history is never read.
func (s *Store) SourceRecords(ctx context.Context, scope SourceScope, ref SourceRef) (SourceRecords, error) {
	var out SourceRecords
	if scope.WorkspaceID == "" || scope.GroupID == "" || scope.GroupRevision == "" || scope.ConnectorID == "" || scope.ConnectorRevision == "" {
		return out, ErrInvalidGroupKey
	}
	// A long source has its own partition; a wide group has one row per window and so one partition.
	entity := SourceEntityKey(ref.DeviceID, ref.PointID, ref.TagID)
	if scope.Wide {
		entity = ""
	}
	partition := partitionKey(GroupKey{GroupID: scope.GroupID}, entity)
	latest, err := s.sourceRows(ctx, scope, partition, "", 1)
	if err != nil {
		return out, err
	}
	for _, row := range latest {
		record, err := s.verifySource(ctx, scope, ref, row)
		if err != nil {
			return out, err
		}
		out.Latest = &record
	}
	committed, err := s.sourceRows(ctx, scope, partition, StateCommitted, sourceCandidateLimit)
	if err != nil {
		return out, err
	}
	for _, row := range committed {
		record, err := s.verifySource(ctx, scope, ref, row)
		if err != nil {
			return out, err
		}
		if record.Available && record.Delivery == SourceDeliveryCommitted {
			out.LastSQLCommitted = &record
			break
		}
	}
	return out, nil
}

// sourceRows returns the newest rows of one partition and applied revision. bucket_start is fixed width UTC, so
// its text order is time order.
func (s *Store) sourceRows(ctx context.Context, scope SourceScope, partition, state string, limit int) ([]sourceRow, error) {
	query := `SELECT effect_key, record_id, workspace_id, group_id, group_revision, entity_key, connector_id, connector_revision,
			bucket_start, payload, payload_digest, state, committed_at
		FROM wg_delivery_outbox WHERE partition_key = ? AND group_revision = ?`
	args := []any{partition, scope.GroupRevision}
	if state != "" {
		query += ` AND state = ?`
		args = append(args, state)
	}
	rows, err := s.db.QueryContext(ctx, query+` ORDER BY COALESCE(due_at,bucket_start) DESC, effect_key DESC LIMIT ?`, append(args, limit)...)
	if err != nil {
		return nil, fmt.Errorf("read source records: %w", err)
	}
	defer rows.Close()
	var found []sourceRow
	for rows.Next() {
		var row sourceRow
		if err := rows.Scan(&row.effectKey, &row.recordID, &row.workspaceID, &row.groupID, &row.groupRevision, &row.entityKey,
			&row.connectorID, &row.connectorRevision, &row.bucketStart, &row.payload, &row.digest, &row.state, &row.committedAt); err != nil {
			return nil, fmt.Errorf("scan source record: %w", err)
		}
		found = append(found, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read source records: %w", err)
	}
	return found, nil
}
