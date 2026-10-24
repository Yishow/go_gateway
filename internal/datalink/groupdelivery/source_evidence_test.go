package groupdelivery

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"maps"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

var (
	sourceA = SourceRef{CircuitID: "cir-a", DeviceID: "dev-1", PointID: "pt-a", TagID: "tag-a"}
	sourceB = SourceRef{CircuitID: "cir-a", DeviceID: "dev-1", PointID: "pt-b", TagID: "tag-b"}
	scope   = SourceScope{
		WorkspaceID: testKey.WorkspaceID, GroupID: testKey.GroupID, GroupRevision: testKey.GroupRevision, ProfileKind: "instantaneous",
		ConnectorID: testDestination.ConnectorID, ConnectorRevision: testDestination.ConnectorRevision,
	}
)

// longColumns is the fixed column order of a managed long row, as dbtarget encodes it.
var longColumns = []string{
	"bucket_start", "bucket_end", "record_id", "workspace_id", "group_id", "group_revision", "circuit_id", "device_id", "point_id",
	"tag_id", "measurement_id", "series_epoch", "profile_kind", "value_type", "unit", "status", "quality", "observed_at", "received_at",
	"value_int", "value_decimal", "value_float", "value_bool", "value_text", "delta_value", "delta_from", "delta_to", "provenance",
}

// longSpec describes the stored record of one source and window; unset fields stay NULL.
type longSpec struct {
	source             SourceRef
	start              int
	status, reason     string
	profile            string
	slots              map[string]any
	delta              bool
	workspace, group   string
	revision           string
	provenance, window string
}

func timeText(at time.Time) string { return at.UTC().Format("2006-01-02T15:04:05.000000Z") }

func (l longSpec) row(t *testing.T) (dbtarget.EncodedRow, snapshot.Outcome) {
	t.Helper()
	start, end := at(l.start), at(l.start+15) // one 15 second window
	profile := l.profile
	if profile == "" {
		profile = "instantaneous"
	}
	workspace, group, revision := l.workspace, l.group, l.revision
	if workspace == "" {
		workspace = testKey.WorkspaceID
	}
	if group == "" {
		group = testKey.GroupID
	}
	if revision == "" {
		revision = testKey.GroupRevision
	}
	entity := SourceEntityKey(l.source.DeviceID, l.source.PointID, l.source.TagID)
	recordID := fmt.Sprintf("record-%s-%d", l.source.TagID, l.start)
	valueType := "float64"
	for slot, typ := range map[string]string{"value_int": "int64", "value_bool": "bool", "value_decimal": "decimal", "value_text": "text"} {
		if _, ok := l.slots[slot]; ok {
			valueType = typ
		}
	}
	provenance := l.provenance
	if provenance == "" {
		provenance = fmt.Sprintf(`{"type":%q,"member":%q,"reason":%q,"profile":{"kind":%q,"interval":15,"unit":"second"}}`, valueType, entity, l.reason, profile)
	}
	values := map[string]any{
		"bucket_start": timeText(start), "bucket_end": timeText(end), "record_id": recordID, "workspace_id": workspace, "group_id": group,
		"group_revision": revision, "circuit_id": l.source.CircuitID, "device_id": l.source.DeviceID, "point_id": l.source.PointID,
		"tag_id": l.source.TagID, "profile_kind": profile, "value_type": valueType, "status": l.status, "provenance": provenance,
	}
	maps.Copy(values, l.slots)
	if l.delta {
		values["delta_value"] = "1"
		values["delta_from"], values["delta_to"] = timeText(start), timeText(end)
	}
	if l.window != "" {
		values["bucket_end"] = l.window
	}
	row := dbtarget.EncodedRow{RecordID: recordID, EffectKey: "effect-" + recordID, EntityKey: entity, BucketStart: start}
	for _, column := range longColumns {
		row.Cells = append(row.Cells, dbtarget.EncodedCell{Column: column, Value: values[column]})
	}
	return row, snapshot.Outcome{
		Kind: snapshot.OutcomeRow, EntityKey: entity, RecordID: recordID, EffectKey: row.EffectKey, BucketStart: start, BucketEnd: end,
	}
}

// enqueueLong closes one record into the outbox, as the pipeline's closure does.
func enqueueLong(t *testing.T, store *Store, spec longSpec) string {
	t.Helper()
	row, outcome := spec.row(t)
	key := testKey
	if spec.revision != "" {
		key.GroupRevision = spec.revision
	}
	require.NoError(t, store.CommitClosure(t.Context(), key, Closure{
		Destination: testDestination, NextClose: outcome.BucketEnd, Buckets: []ClosedBucket{{Outcome: outcome, Row: &row}},
	}))
	return row.EffectKey
}

// confirm completes a delivery with the receipt the destination would have written.
func confirm(t *testing.T, store *Store, effect string) {
	t.Helper()
	item, err := store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	claim, err := store.BeginDelivery(t.Context(), effect, "writer", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.CompleteDelivery(t.Context(), claim, item.PayloadDigest))
}

func rewritePayload(t *testing.T, db *sql.DB, effect, payload string) {
	t.Helper()
	sum := sha256.Sum256([]byte(payload))
	_, err := db.ExecContext(t.Context(), `UPDATE wg_delivery_outbox SET payload = ?, payload_digest = ? WHERE effect_key = ?`, payload, hex.EncodeToString(sum[:]), effect)
	require.NoError(t, err)
}

func readSource(t *testing.T, store *Store, ref SourceRef) SourceRecords {
	t.Helper()
	got, err := store.SourceRecords(t.Context(), scope, ref)
	require.NoError(t, err)
	return got
}

func TestSourceEntityKeyIsTheMemberKeyTheBoundaryUses(t *testing.T) {
	require.Equal(t, `["dev-1","pt-a","tag-a"]`, SourceEntityKey("dev-1", "pt-a", "tag-a"))
}

func TestSourceRecordsKeepRawPresentAndDeltaUnknownApart(t *testing.T) {
	store, _, _ := newTestStore(t)
	cumulative := scope
	cumulative.ProfileKind = "cumulative"
	enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "baseline", reason: "first-reading", profile: "cumulative", slots: map[string]any{"value_decimal": "100"}})
	got, err := store.SourceRecords(t.Context(), cumulative, sourceA)
	require.NoError(t, err)
	require.NotNil(t, got.Latest)
	require.True(t, got.Latest.Available)
	require.Equal(t, SourceStatusBaseline, got.Latest.Status)
	require.Equal(t, "first-reading", got.Latest.Reason)
	require.True(t, got.Latest.RawPresent, "the raw reading was recorded")
	require.NotNil(t, got.Latest.DeltaKnown)
	require.False(t, *got.Latest.DeltaKnown, "no comparable baseline: the difference is unknown, not zero")
	require.Equal(t, SourceDeliveryQueued, got.Latest.Delivery)
	require.Nil(t, got.Latest.CommittedAt)
	require.Nil(t, got.LastSQLCommitted, "a queued record is not SQL evidence")
}

func TestSourceRecordsCountZeroAndFalseAsPresent(t *testing.T) {
	store, _, _ := newTestStore(t)
	enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "ok", slots: map[string]any{"value_int": int64(0)}})
	enqueueLong(t, store, longSpec{source: sourceB, start: 0, status: "ok", slots: map[string]any{"value_bool": false}})
	require.True(t, readSource(t, store, sourceA).Latest.RawPresent)
	require.True(t, readSource(t, store, sourceB).Latest.RawPresent)
	require.Nil(t, readSource(t, store, sourceA).Latest.DeltaKnown, "an instantaneous record has no difference to report")
}

func TestSourceRecordsKeepSourcesOfOneGroupIndependent(t *testing.T) {
	store, _, _ := newTestStore(t)
	enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "missing", reason: "no-samples"})
	enqueueLong(t, store, longSpec{source: sourceB, start: 0, status: "ok", slots: map[string]any{"value_float": 21.5}})
	a, b := readSource(t, store, sourceA).Latest, readSource(t, store, sourceB).Latest
	require.Equal(t, SourceStatusMissing, a.Status)
	require.Equal(t, "no-samples", a.Reason)
	require.False(t, a.RawPresent, "a missing source is stored as NULL")
	require.Equal(t, SourceStatusOK, b.Status)
	require.True(t, b.RawPresent)
	require.Empty(t, b.Reason)
}

func TestSourceRecordsSeparateTheNewestQueuedFromTheEarlierConfirmedRecord(t *testing.T) {
	store, _, _ := newTestStore(t)
	older := enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "ok", slots: map[string]any{"value_float": 1.0}})
	confirm(t, store, older)
	enqueueLong(t, store, longSpec{source: sourceA, start: 15, status: "ok", slots: map[string]any{"value_float": 2.0}})
	got := readSource(t, store, sourceA)
	require.Equal(t, SourceDeliveryQueued, got.Latest.Delivery)
	require.Equal(t, at(15), got.Latest.WindowStart)
	require.Nil(t, got.Latest.CommittedAt)
	require.NotNil(t, got.LastSQLCommitted)
	require.Equal(t, SourceDeliveryCommitted, got.LastSQLCommitted.Delivery)
	require.Equal(t, at(0), got.LastSQLCommitted.WindowStart)
	require.Equal(t, older, got.LastSQLCommitted.EffectKey)
	require.NotNil(t, got.LastSQLCommitted.CommittedAt)
}

func TestSourceRecordsPickTheLatestWindowRegardlessOfInsertionOrder(t *testing.T) {
	store, _, _ := newTestStore(t)
	enqueueLong(t, store, longSpec{source: sourceA, start: 30, status: "ok", slots: map[string]any{"value_float": 3.0}})
	enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "ok", slots: map[string]any{"value_float": 1.0}})
	enqueueLong(t, store, longSpec{source: sourceA, start: 15, status: "ok", slots: map[string]any{"value_float": 2.0}})
	require.Equal(t, at(30), readSource(t, store, sourceA).Latest.WindowStart)
}

func TestSourceRecordsNeverPromoteWithoutAMatchingReceipt(t *testing.T) {
	store, _, _ := newTestStore(t)
	effect := enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "ok", slots: map[string]any{"value_float": 1.0}})
	claim, err := store.BeginDelivery(t.Context(), effect, "writer", time.Minute)
	require.NoError(t, err)
	require.NoError(t, store.CompleteDelivery(t.Context(), claim, "a-different-digest"))
	got := readSource(t, store, sourceA)
	require.NotNil(t, got.Latest)
	require.NotEqual(t, SourceDeliveryCommitted, got.Latest.Delivery, "a mismatched receipt is not evidence for this payload")
	require.Nil(t, got.Latest.CommittedAt)
	require.Nil(t, got.LastSQLCommitted)
}

func TestSourceRecordsExposeOnlyStableReasonCodes(t *testing.T) {
	store, _, _ := newTestStore(t)
	enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "bad", reason: "driver said password=hunter2"})
	enqueueLong(t, store, longSpec{source: sourceB, start: 0, status: "bad", reason: "quality-bad"})
	got := readSource(t, store, sourceA).Latest
	require.Equal(t, SourceReasonUnavailable, got.Reason)
	require.NotContains(t, fmt.Sprintf("%+v", got.VerifiedSourceRecord), "hunter2")
	require.Equal(t, "quality-bad", readSource(t, store, sourceB).Latest.Reason)
}

func TestSourceRecordsStayWithinTheAppliedRevisionAndSource(t *testing.T) {
	store, _, _ := newTestStore(t)
	enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "ok", slots: map[string]any{"value_float": 1.0}})
	enqueueLong(t, store, longSpec{source: sourceA, start: 60, status: "missing", reason: "no-samples", revision: "rev-old"})
	enqueueLong(t, store, longSpec{source: sourceB, start: 120, status: "ok", slots: map[string]any{"value_float": 9.0}})
	got := readSource(t, store, sourceA)
	require.Equal(t, at(0), got.Latest.WindowStart, "a drained older revision and another source never replace the applied record")
	require.Equal(t, SourceStatusOK, got.Latest.Status)
}

func TestSourceRecordsReportNothingWhenTheLocalRecordsWereReclaimed(t *testing.T) {
	store, db, _ := newTestStore(t)
	confirm(t, store, enqueueLong(t, store, longSpec{source: sourceA, start: 0, status: "ok", slots: map[string]any{"value_float": 1.0}}))
	_, err := db.ExecContext(t.Context(), `DELETE FROM wg_delivery_outbox`)
	require.NoError(t, err)
	got := readSource(t, store, sourceA)
	require.Nil(t, got.Latest, "reclaimed evidence is unavailable, not a baseline, zero or failure")
	require.Nil(t, got.LastSQLCommitted)
}
