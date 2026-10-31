package groupdelivery

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"
	"reflect"
	"strings"
	"time"
)

func (s *Store) verifyMixedSource(ctx context.Context, scope SourceScope, ref SourceRef, facts mixedSourceFacts, stored mixedSourceOutcome) (SourceRecord, error) {
	reject := func(cause SourceUnavailableCause) (SourceRecord, error) {
		return unavailable(SourceDeliveryUnknown, cause), nil
	}
	o, err := DecodeMixedOutcome([]byte(stored.raw), stored.digest)
	if err != nil {
		return reject(SourceCausePayloadInvalid)
	}
	if o.RecordID != stored.id || timestamp(o.DueAt) != stored.due || string(o.Kind) != stored.kind || !mixedSourceIdentity(scope, o) {
		return reject(SourceCauseIdentityMismatch)
	}
	var member *snapshot.MemberResult
	seen := map[string]bool{}
	for i := range o.Members {
		m := &o.Members[i]
		if seen[m.MemberKey] {
			return reject(SourceCauseIdentityMismatch)
		}
		intent, owned := facts.intents[m.MemberKey]
		if !owned || m.Window == nil || m.Window.Calculation != intent.method || m.Window.Interval != intent.interval || (o.Kind == snapshot.OutcomeRow && intent.required && !m.Usable() && m.Status != snapshot.MemberNotDue) {
			return reject(SourceCauseIdentityMismatch)
		}
		seen[m.MemberKey] = true
		if m.MemberKey == facts.member.MemberKey {
			member = m
		}
	}
	if member == nil || (scope.Wide && len(seen) != len(facts.intents)) {
		return reject(SourceCauseIdentityMismatch)
	}
	if member.Sample != nil && (member.Sample.SourceRevision != facts.sourceRevision || member.Sample.MappingRevision != facts.mappingRevision || member.Sample.CounterProof != facts.counterRef) {
		return reject(SourceCauseIdentityMismatch)
	}
	evidence, err := dbtarget.MixedSourceEvidence(facts.member, o, *member, scope.Wide)
	if err != nil {
		return reject(SourceCausePayloadInvalid)
	}
	if ref.Calculation == sourceReductionMean && member.Status != snapshot.MemberNotDue {
		e := member.Evidence
		if e == nil || e.Method != measurement.MeanAggregationMethod || e.Version != measurement.MeanAggregationVersion || string(e.SourceValueType) != ref.SourceType || string(e.ResultValueType) != ref.ResultType || e.ValidCount < 0 || e.RejectedCount < 0 || !sourceMeanCoverages[string(e.Coverage)] {
			return reject(SourceCausePayloadInvalid)
		}
		if member.Usable() && (member.Reduction == nil || member.Reduction.Count != e.ValidCount || member.Reduction.Basis != e.Basis) {
			return reject(SourceCausePayloadInvalid)
		}
	}
	record := projectMixedSource(scope, ref, o, *member, stored.digest)
	row, err := s.mixedSourceRow(ctx, scope, o.RecordID)
	if o.Kind != snapshot.OutcomeRow {
		if !errors.Is(err, sql.ErrNoRows) {
			if err != nil {
				return SourceRecord{}, err
			}
			return reject(SourceCauseIdentityMismatch)
		}
		if o.EffectKey != "" {
			return reject(SourceCauseIdentityMismatch)
		}
		return SourceRecord{Available: true, Delivery: SourceDeliveryWindowSkipped, VerifiedSourceRecord: record}, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return reject(SourceCausePayloadInvalid)
	}
	if err != nil {
		return SourceRecord{}, err
	}
	encoded, cells, cause, ok := openFrozenRow(scope, row)
	if !ok {
		return reject(cause)
	}
	if row.dueAt != timestamp(o.DueAt) || encoded.EffectKey != o.EffectKey || !encoded.DueAt.Equal(o.DueAt) || !encoded.BucketStart.Equal(o.BucketStart) || len(cells) != len(encoded.Cells) {
		return reject(SourceCauseIdentityMismatch)
	}
	if !mixedSourceCells(scope, ref, facts, o, *member, evidence, cells) {
		return reject(SourceCausePayloadInvalid)
	}
	matched, err := s.receiptMatches(ctx, row)
	if err != nil {
		return SourceRecord{}, err
	}
	delivery := sourceDelivery(row.state, matched)
	record.EffectKey, record.PayloadDigest = row.effectKey, row.digest
	if delivery == SourceDeliveryCommitted {
		at, ok := frozenTime(row.committedAt)
		if !ok {
			return reject(SourceCausePayloadInvalid)
		}
		record.CommittedAt = &at
	}
	return SourceRecord{Available: true, Delivery: delivery, VerifiedSourceRecord: record}, nil
}
func mixedSourceCells(scope SourceScope, ref SourceRef, facts mixedSourceFacts, o snapshot.Outcome, m snapshot.MemberResult, evidence []byte, cells map[string]any) bool {
	provenanceColumn, recordColumn, dueColumn, valueColumn := "provenance", "record_id", "due_at", sourceSlotByType[ref.ResultType]
	if scope.Wide {
		provenanceColumn, recordColumn, dueColumn, valueColumn = scope.ProvenanceColumn, scope.RecordKeyColumn, scope.DueAtColumn, ref.Column
	}
	if cellText(cells, strings.ToLower(recordColumn)) != o.RecordID {
		return false
	}
	if scope.Wide {
		due, ok := timeCell(cells, strings.ToLower(dueColumn))
		if !ok || !due.Equal(o.DueAt) {
			return false
		}
	}
	if !scope.Wide {
		for col, want := range map[string]string{"workspace_id": scope.WorkspaceID, "group_id": scope.GroupID, "group_revision": scope.GroupRevision, "device_id": ref.DeviceID, "point_id": ref.PointID, "tag_id": ref.TagID, "profile_kind": mixedProfileKind} {
			if cellText(cells, col) != want {
				return false
			}
		}
	}
	raw := cellText(cells, strings.ToLower(provenanceColumn))
	if scope.Wide {
		var entries []json.RawMessage
		if json.Unmarshal([]byte(raw), &entries) != nil {
			return false
		}
		raw = ""
		found := false
		for _, entry := range entries {
			var identity struct {
				Member string `json:"member"`
			}
			if json.Unmarshal(entry, &identity) != nil {
				return false
			}
			if identity.Member == m.MemberKey {
				if found {
					return false
				}
				found = true
				raw = string(entry)
			}
		}
	}
	if !equalMixedJSON([]byte(raw), evidence) {
		return false
	}
	actual, exists := cells[strings.ToLower(valueColumn)]
	if !exists {
		return false
	}
	if !m.Usable() {
		return actual == nil
	}
	selected, err := dbtarget.MixedSelectedValue(facts.member.Calculation, m)
	if err != nil || selected.Type() != facts.member.Type {
		return false
	}
	expected, err := dbtarget.EncodeExactValue(facts.dialect, facts.declared, selected)
	return err == nil && reflect.DeepEqual(expected, actual)
}
func equalMixedJSON(a, b []byte) bool {
	if !counterproof.UniqueJSONObject(a, false) || !counterproof.UniqueJSONObject(b, false) {
		return false
	}
	decode := func(raw []byte) (any, error) {
		d := json.NewDecoder(bytes.NewReader(raw))
		d.UseNumber()
		var v any
		err := d.Decode(&v)
		return v, err
	}
	av, ae := decode(a)
	bv, be := decode(b)
	return ae == nil && be == nil && reflect.DeepEqual(av, bv)
}
func timePointer(at time.Time) *time.Time {
	if at.IsZero() {
		return nil
	}
	cloned := at
	return &cloned
}
func projectMixedSource(scope SourceScope, ref SourceRef, o snapshot.Outcome, m snapshot.MemberResult, digest string) *VerifiedSourceRecord {
	c := &SourceCalculation{ResultShape: scope.ResultShape, Method: ref.Calculation, IntervalSeconds: ref.CadenceSeconds, Alignment: ref.Alignment, DueAt: o.DueAt, OutcomeDigest: digest, OutcomeKind: string(o.Kind), Status: string(m.Status), Reason: sourceReason(m.Reason), ResultPresent: m.Usable(), RowEmitted: o.Kind == snapshot.OutcomeRow, RawPresent: m.Sample != nil, SourceType: ref.SourceType, ResultType: ref.ResultType}
	if m.Status == snapshot.MemberNotDue {
		c.Reason = snapshot.ReasonNotDue
	}
	r := &VerifiedSourceRecord{Calculation: c, GroupRevision: scope.GroupRevision, ConnectorRevision: scope.ConnectorRevision, RecordID: o.RecordID, ProfileKind: mixedProfileKind, Status: SourceStatus(m.Status), Reason: c.Reason, RawPresent: c.RawPresent}
	if w := m.Window; w != nil {
		c.NextDue = timePointer(w.NextDue)
		c.AlignedStart = timePointer(w.AlignedStart)
		c.AlignedEnd = timePointer(w.AlignedEnd)
		c.ActualStart = timePointer(w.ActualStart)
		c.ActualEnd = timePointer(w.ActualEnd)
		r.WindowStart = w.ActualStart
		r.WindowEnd = w.ActualEnd
	}
	if m.Counter != nil {
		r.Status = SourceStatus(m.Counter.Status)
		r.Reason = sourceReason(m.Counter.Reason)
	}
	if ref.Calculation == mixedDeltaMethod {
		known := m.Usable()
		r.DeltaKnown = &known
	}
	if m.Evidence != nil {
		e := m.Evidence
		r.Aggregation = &SourceAggregation{Method: e.Method, Version: e.Version, Basis: string(e.Basis), SourceType: ref.SourceType, ResultType: ref.ResultType, ValidCount: e.ValidCount, RejectedCount: e.RejectedCount, Coverage: string(e.Coverage), Reason: sourceReason(e.Reason)}
	}
	if m.Reduction != nil {
		rounded := m.Reduction.Rounded
		c.Rounded = &rounded
	}
	return r
}
