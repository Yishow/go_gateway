package groupdelivery

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMixedSourceEvidenceReadsRequiredSkipAndNotDue(t *testing.T) {
	store, scope, refs, _ := mixedSourceFixture(t, true)
	for _, tag := range []string{"fast", "slow"} {
		records, err := store.SourceRecords(t.Context(), scope, refs[tag])
		require.NoError(t, err)
		require.NotNil(t, records.Latest, "local skip is evidence even without an outbox row")
		require.True(t, records.Latest.Available)
		require.NotNil(t, records.Latest.Calculation)
		require.Equal(t, SourceDeliveryWindowSkipped, records.Latest.Delivery)
		require.False(t, records.Latest.Calculation.RowEmitted)
		require.Nil(t, records.LastSQLCommitted)
		require.Empty(t, records.Latest.EffectKey)
		require.Empty(t, records.Latest.PayloadDigest)
		if tag == "slow" {
			require.Equal(t, "not_due", records.Latest.Calculation.Status)
			require.False(t, records.Latest.Calculation.ResultPresent)
			require.Nil(t, records.Latest.Calculation.ActualStart)
		} else {
			require.True(t, records.Latest.Calculation.ResultPresent)
		}
	}
}

func TestMixedSourceEvidenceRowReceiptAndRetention(t *testing.T) {
	store, scope, refs, c := mixedSourceFixture(t, false)
	ref := refs["fast"]
	got, err := store.SourceRecords(t.Context(), scope, ref)
	require.NoError(t, err)
	require.NotNil(t, got.Latest)
	require.True(t, got.Latest.Available)
	require.Equal(t, SourceDeliveryQueued, got.Latest.Delivery)
	require.True(t, got.Latest.Calculation.RowEmitted)
	require.True(t, got.Latest.Calculation.ResultPresent)
	require.True(t, got.Latest.Calculation.RawPresent)
	require.Nil(t, got.LastSQLCommitted)
	confirm(t, store, c.Buckets[0].Outcome.EffectKey)
	got, err = store.SourceRecords(t.Context(), scope, ref)
	require.NoError(t, err)
	require.NotNil(t, got.LastSQLCommitted)
	require.Equal(t, SourceDeliveryCommitted, got.LastSQLCommitted.Delivery)
	slow, err := store.SourceRecords(t.Context(), scope, refs["slow"])
	require.NoError(t, err)
	require.NotNil(t, slow.LastSQLCommitted)
	require.Equal(t, "not_due", slow.LastSQLCommitted.Calculation.Status)
	require.False(t, slow.LastSQLCommitted.Calculation.ResultPresent)
	require.False(t, slow.LastSQLCommitted.Calculation.RawPresent)
	require.True(t, slow.LastSQLCommitted.Calculation.RowEmitted)
	require.Nil(t, slow.LastSQLCommitted.Calculation.ActualEnd)
	_, err = store.db.Exec(`DELETE FROM wg_delivery_mixed_outcomes`)
	require.NoError(t, err)
	got, err = store.SourceRecords(t.Context(), scope, ref)
	require.NoError(t, err)
	require.Nil(t, got.Latest)
	require.Nil(t, got.LastSQLCommitted, "outbox alone cannot reconstruct calculation evidence")
}
func TestMixedSourceEvidenceRejectsChangedFrozenIntent(t *testing.T) {
	store, scope, refs, _ := mixedSourceFixture(t, true)
	for _, change := range []func(*SourceScope, *SourceRef){
		func(s *SourceScope, r *SourceRef) { s.RecordIdentity = "foreign" },
		func(s *SourceScope, r *SourceRef) { s.ResultShape = "future" },
		func(s *SourceScope, r *SourceRef) { s.ConnectorRevision = "foreign" },
		func(s *SourceScope, r *SourceRef) { r.Calculation = "mean" },
		func(s *SourceScope, r *SourceRef) { r.CadenceSeconds = 60 },
		func(s *SourceScope, r *SourceRef) { r.MeasurementID = "foreign" },
		func(s *SourceScope, r *SourceRef) { r.Required = true },
		func(s *SourceScope, r *SourceRef) { r.ResultType = "decimal" },
	} {
		s, r := scope, refs["fast"]
		change(&s, &r)
		_, err := store.SourceRecords(t.Context(), s, r)
		require.Error(t, err)
	}
}
func TestMixedSourceEvidenceRejectsCorruptTypedOutcome(t *testing.T) {
	store, scope, refs, _ := mixedSourceFixture(t, true)
	_, err := store.db.Exec(`UPDATE wg_delivery_mixed_outcomes SET outcome_digest='wrong'`)
	require.NoError(t, err)
	got, err := store.SourceRecords(t.Context(), scope, refs["fast"])
	require.NoError(t, err)
	require.NotNil(t, got.Latest)
	require.False(t, got.Latest.Available)
	require.Nil(t, got.Latest.VerifiedSourceRecord)
}
func TestMixedSourceEvidenceLongAndMeanUseActualEncoder(t *testing.T) {
	for _, long := range []bool{false, true} {
		t.Run(map[bool]string{false: "wide", true: "long"}[long], func(t *testing.T) {
			store, scope, refs, _ := mixedSourceFixtureMode(t, false, long, true)
			got, err := store.SourceRecords(t.Context(), scope, refs["fast"])
			require.NoError(t, err)
			require.NotNil(t, got.Latest)
			require.True(t, got.Latest.Available)
			require.Equal(t, "mean", got.Latest.Calculation.Method)
			require.Equal(t, "decimal", got.Latest.Calculation.ResultType)
			require.NotNil(t, got.Latest.Aggregation)
			require.Equal(t, "unknown", got.Latest.Aggregation.Coverage)
			require.EqualValues(t, 3, got.Latest.Aggregation.ValidCount)
			require.NotNil(t, got.Latest.Calculation.Rounded)
			require.True(t, *got.Latest.Calculation.Rounded)
			slow, err := store.SourceRecords(t.Context(), scope, refs["slow"])
			require.NoError(t, err)
			if long {
				require.Nil(t, slow.Latest)
			} else {
				require.Equal(t, "not_due", slow.Latest.Calculation.Status)
			}
		})
	}
}
func TestMixedSourceEvidenceRejectsRehashedWrongSelectedCell(t *testing.T) {
	store, scope, refs, c := mixedSourceFixture(t, false)
	effect := c.Buckets[0].Outcome.EffectKey
	item, err := store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	row, err := DecodeRowPayload(item.Payload)
	require.NoError(t, err)
	for i := range row.Cells {
		if row.Cells[i].Column == "fast" {
			row.Cells[i].Value = int64(7)
		}
	}
	raw, digest, err := encodeRowPayload(row)
	require.NoError(t, err)
	_, err = store.db.Exec(`UPDATE wg_delivery_outbox SET payload=?,payload_digest=? WHERE effect_key=?`, string(raw), digest, effect)
	require.NoError(t, err)
	got, err := store.SourceRecords(t.Context(), scope, refs["fast"])
	require.NoError(t, err)
	require.False(t, got.Latest.Available)
}
func TestMixedSourceEvidenceRejectsRehashedForeignObservation(t *testing.T) {
	store, scope, refs, _ := mixedSourceFixture(t, true)
	var raw, digest, id string
	require.NoError(t, store.db.QueryRow(`SELECT record_id,outcome,outcome_digest FROM wg_delivery_mixed_outcomes`).Scan(&id, &raw, &digest))
	out, err := DecodeMixedOutcome([]byte(raw), digest)
	require.NoError(t, err)
	for i := range out.Members {
		if out.Members[i].MemberKey == SourceEntityKey("d", "p", "fast") {
			out.Members[i].Sample.SourceRevision = "foreign"
		}
	}
	encoded, err := encodeMixedOutcome(out)
	require.NoError(t, err)
	_, err = store.db.Exec(`UPDATE wg_delivery_mixed_outcomes SET outcome=?,outcome_digest=? WHERE record_id=?`, string(encoded), handoffHash(string(encoded)), id)
	require.NoError(t, err)
	got, err := store.SourceRecords(t.Context(), scope, refs["fast"])
	require.NoError(t, err)
	require.False(t, got.Latest.Available)
	require.Equal(t, SourceCauseIdentityMismatch, got.Latest.Cause)
}
