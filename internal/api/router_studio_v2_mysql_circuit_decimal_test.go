package api

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

const (
	mysqlCircuitDecimalBaselineID = "mysql-decimal-baseline"
	mysqlCircuitDecimalEndID      = "mysql-decimal-end"
	mysqlCircuitDecimalBaseline   = "12000.0"
	mysqlCircuitDecimalEnd        = "12000.4"
)

var mysqlCircuitDecimalInputs = []mysqlCircuitInput{
	{sampleID: mysqlCircuitDecimalBaselineID, offset: 9 * time.Minute, value: 120000},
	{sampleID: mysqlCircuitDecimalEndID, offset: 19 * time.Minute, value: 120004},
}

type mysqlCircuitDecimalCounterProvenance struct {
	RawText               string `json:"raw_text"`
	BaselineSampleID      string `json:"baseline_sample_id"`
	EndSampleID           string `json:"end_sample_id"`
	DefinitionRevision    string `json:"definition_revision"`
	SourceBindingRevision string `json:"source_binding_revision"`
	ProfileRevision       string `json:"profile_revision"`
}

type mysqlCircuitDecimalProfileProvenance struct {
	Kind     string `json:"kind"`
	Interval int    `json:"interval"`
	Unit     string `json:"unit"`
}

type mysqlCircuitDecimalProvenance struct {
	Type            string                                `json:"type"`
	SampleID        string                                `json:"sample_id"`
	SourceRevision  string                                `json:"source_revision"`
	MappingRevision string                                `json:"mapping_revision"`
	SourceType      string                                `json:"source_type"`
	ResultType      string                                `json:"result_type"`
	Role            string                                `json:"role"`
	Profile         mysqlCircuitDecimalProfileProvenance  `json:"profile"`
	Counter         *mysqlCircuitDecimalCounterProvenance `json:"counter"`
}

type mysqlCircuitDecimalSQLRow struct {
	recordID, groupID, circuitID, deviceID, pointID, tagID string
	profile, valueType, status                             string
	measurementID, seriesEpoch, unit                       sql.NullString
	raw, delta                                             sql.NullString
	bucketStart                                            time.Time
	observedAt, receivedAt, deltaFrom, deltaTo             sql.NullTime
	provenance                                             string
}

// TestMySQLCircuitPlainDecimalIntegration proves the saved exact scale is
// applied once through both API-created profiles and the durable MySQL path.
func TestMySQLCircuitPlainDecimalIntegration(t *testing.T) {
	runMySQLCircuitPlainDecimalIntegration(t, mysqlCircuitFamilyMySQL)
}

func runMySQLCircuitPlainDecimalIntegration(t *testing.T, family string) {
	f := newMySQLCircuitPlainFixture(t, mysqlCircuitPlainFixtureOptions{
		Family:          family,
		SeedMeasurement: measurement.MeasurementDefinition{ExactRecordingScale: "0.1"},
		Inputs:          mysqlCircuitDecimalInputs,
	})
	require.Equal(t, "0.1", f.definition.ExactRecordingScale)
	require.Equal(t, uint64(120000), f.inputs[0].value)
	require.Equal(t, uint64(120004), f.inputs[1].value)
	instantMember, cumulativeMember := f.instant.group.Members[0], f.cumulative.group.Members[0]
	require.NotEmpty(t, instantMember.CircuitID)
	require.Equal(t, instantMember.CircuitID, cumulativeMember.CircuitID)
	require.Equal(t, instantMember.TagID, cumulativeMember.TagID)
	require.Equal(t, *instantMember.MeasurementID, *cumulativeMember.MeasurementID)

	f.accept(0)
	f.closeAfter(10 * time.Minute)
	baselineBeforeRetry := mysqlCircuitBaseline(t, f)
	requireEventuallyCumulativeRetry(t, f)
	waitForMySQLCircuitOutboxes(t, f, 41)
	require.Equal(t, baselineBeforeRetry, mysqlCircuitBaseline(t, f), "delivery retry cannot advance the counter baseline")

	f.accept(1)
	f.closeAfter(20 * time.Minute)
	waitForMySQLCircuitOutboxes(t, f, 82)

	assertMySQLCircuitDecimalJournal(t, f)
	assertMySQLCircuitColumnTypes(t, f)
	outboxRows := assertMySQLCircuitDecimalOutboxAndReceipts(t, f)
	assertMySQLCircuitDecimalSQLRows(t, f, outboxRows)
}

func assertMySQLCircuitDecimalJournal(t *testing.T, f *mysqlCircuitPlainFixture) {
	t.Helper()
	counts := map[string]int{f.instant.id: 0, f.cumulative.id: 0}
	wantDigits := map[string]string{mysqlCircuitDecimalBaselineID: mysqlCircuitDecimalBaseline, mysqlCircuitDecimalEndID: mysqlCircuitDecimalEnd}
	for _, entry := range f.acceptedJournal {
		groupID, sampleID, payload, digest := entry.groupID, entry.sampleID, entry.payload, entry.digest
		require.Contains(t, counts, groupID)
		require.Zero(t, entry.consumed)
		require.Len(t, digest, 64)
		var decoded mysqlCircuitPlainJournalPayload
		require.NoError(t, json.Unmarshal([]byte(payload), &decoded))
		require.Equal(t, sampleID, decoded.SampleID)
		require.Equal(t, measurement.ExactDecimal, decoded.Value.Type())
		digits, ok := decoded.Value.DecimalDigits()
		require.True(t, ok)
		require.Equal(t, wantDigits[sampleID], digits)
		input := mysqlCircuitDecimalInputByID(t, f, sampleID)
		observed := f.base.Add(input.offset)
		require.True(t, observed.Equal(parseMySQLCircuitTime(t, decoded.ObservedAt)))
		require.True(t, observed.Add(30*time.Millisecond).Equal(parseMySQLCircuitTime(t, decoded.ReceivedAt)))
		counts[groupID]++
	}
	require.Equal(t, map[string]int{f.instant.id: 2, f.cumulative.id: 2}, counts)
}

func assertMySQLCircuitDecimalOutboxAndReceipts(t *testing.T, f *mysqlCircuitPlainFixture) map[string]dbtarget.EncodedRow {
	t.Helper()
	rows, err := f.f.db.QueryContext(t.Context(), `
		SELECT effect_key, group_id, payload, payload_digest, state
		FROM wg_delivery_outbox WHERE group_id IN (?,?) ORDER BY group_id, bucket_start`, f.instant.id, f.cumulative.id)
	require.NoError(t, err)
	defer rows.Close()
	digests := make(map[string]string, 82)
	byRecordID := make(map[string]dbtarget.EncodedRow, 82)
	for rows.Next() {
		var effectKey, groupID, payload, digest, state string
		require.NoError(t, rows.Scan(&effectKey, &groupID, &payload, &digest, &state))
		require.Equal(t, "sql_committed", state)
		row, err := groupdelivery.DecodeRowPayload([]byte(payload))
		require.NoError(t, err)
		require.Equal(t, effectKey, row.EffectKey)
		actualDigest, err := groupdelivery.RowPayloadDigest(row)
		require.NoError(t, err)
		require.Equal(t, digest, actualDigest, "the frozen outbox digest covers every encoded decimal cell")
		require.NotContains(t, byRecordID, row.RecordID, "record identities are unique")
		byRecordID[row.RecordID] = row
		digests[effectKey] = digest
		assertMySQLCircuitDecimalFrozenRow(t, f, groupID, row)
	}
	require.NoError(t, rows.Err())
	require.Len(t, digests, 82)
	require.Len(t, byRecordID, 82)

	localRows, err := f.f.db.QueryContext(t.Context(), `SELECT effect_key, payload_digest FROM wg_delivery_receipts`)
	require.NoError(t, err)
	localReceipts := make(map[string]string, len(digests))
	for localRows.Next() {
		var effectKey, digest string
		require.NoError(t, localRows.Scan(&effectKey, &digest))
		localReceipts[effectKey] = digest
	}
	require.NoError(t, localRows.Err())
	require.NoError(t, localRows.Close())
	require.Equal(t, digests, localReceipts, "local acknowledgements match every frozen row digest")

	remoteRows, err := f.f.target.QueryContext(t.Context(), `
		SELECT effect_key, payload_digest, committed_at, _gw_effect_receipts_v1
		FROM `+"`"+dbtarget.EffectReceiptTable+"`"+` ORDER BY effect_key`)
	require.NoError(t, err)
	remoteReceipts := make(map[string]string, len(digests))
	for remoteRows.Next() {
		var effectKey, digest, owner string
		var committedAt time.Time
		require.NoError(t, remoteRows.Scan(&effectKey, &digest, &committedAt, &owner))
		require.NotEmpty(t, owner)
		require.False(t, committedAt.IsZero())
		remoteReceipts[effectKey] = digest
	}
	require.NoError(t, remoteRows.Err())
	require.NoError(t, remoteRows.Close())
	require.Equal(t, digests, remoteReceipts, "destination receipts match every frozen row digest")
	return byRecordID
}

func assertMySQLCircuitDecimalFrozenRow(t *testing.T, f *mysqlCircuitPlainFixture, groupID string, row dbtarget.EncodedRow) {
	t.Helper()
	instant := groupID == f.instant.id
	require.True(t, instant || groupID == f.cumulative.id)
	profile := f.cumulative
	member := f.cumulative.group.Members[0]
	if instant {
		profile, member = f.instant, f.instant.group.Members[0]
	}
	assertMySQLCircuitCell(t, row, "record_id", row.RecordID)
	assertMySQLCircuitCell(t, row, "group_id", groupID)
	assertMySQLCircuitCell(t, row, "circuit_id", member.CircuitID)
	assertMySQLCircuitCell(t, row, "device_id", member.DeviceID)
	assertMySQLCircuitCell(t, row, "point_id", member.PointID)
	assertMySQLCircuitCell(t, row, "tag_id", member.TagID)
	assertMySQLCircuitCell(t, row, "measurement_id", f.definition.ID)
	assertMySQLCircuitCell(t, row, "profile_kind", string(profile.group.RecordingProfile.Kind))
	assertMySQLCircuitCell(t, row, "value_type", string(measurement.ExactDecimal))
	assertMySQLCircuitCell(t, row, "bucket_start", row.BucketStart)

	var provenance mysqlCircuitDecimalProvenance
	require.NoError(t, json.Unmarshal([]byte(mysqlCircuitStringCell(t, row, "provenance")), &provenance))
	require.Equal(t, string(measurement.ExactDecimal), provenance.Type)
	require.Equal(t, profile.group.RecordingProfile.Interval, provenance.Profile.Interval)
	require.Equal(t, string(profile.group.RecordingProfile.Unit), provenance.Profile.Unit)
	require.Equal(t, string(profile.group.RecordingProfile.Kind), provenance.Profile.Kind)
	if instant {
		require.Equal(t, "instantaneous", provenance.Profile.Kind)
		require.Equal(t, 15, provenance.Profile.Interval)
		require.Equal(t, "second", provenance.Profile.Unit)
		require.Equal(t, "uint64", provenance.SourceType)
		require.Equal(t, string(measurement.ExactDecimal), provenance.ResultType)
		require.Equal(t, "counter_raw", provenance.Role)
		input, hasInput := mysqlCircuitDecimalInputAt(f, row.BucketStart)
		if !hasInput {
			require.Equal(t, "missing", mysqlCircuitStringCell(t, row, "status"))
			require.Nil(t, mysqlCircuitCell(t, row, "value_decimal"))
			require.Empty(t, provenance.SampleID)
			require.Nil(t, mysqlCircuitCell(t, row, "observed_at"))
			require.Nil(t, mysqlCircuitCell(t, row, "received_at"))
			return
		}
		expected := mysqlCircuitDecimalDigits(input.sampleID)
		require.Equal(t, "ok", mysqlCircuitStringCell(t, row, "status"))
		assertMySQLCircuitCell(t, row, "value_decimal", expected)
		require.Equal(t, input.sampleID, provenance.SampleID)
		require.Equal(t, member.SourceRevision, provenance.SourceRevision)
		require.Equal(t, member.MappingRevision, provenance.MappingRevision)
		require.Nil(t, mysqlCircuitCell(t, row, "delta_value"))
		assertMySQLCircuitTimeCell(t, row, "observed_at", f.base.Add(input.offset))
		assertMySQLCircuitTimeCell(t, row, "received_at", f.base.Add(input.offset+30*time.Millisecond))
		return
	}

	require.Equal(t, "cumulative", provenance.Profile.Kind)
	require.Equal(t, 10, provenance.Profile.Interval)
	require.Equal(t, "minute", provenance.Profile.Unit)
	require.NotNil(t, provenance.Counter)
	require.Equal(t, f.definition.DefinitionRevision, provenance.Counter.DefinitionRevision)
	require.Equal(t, f.definition.SourceBindingRevision, provenance.Counter.SourceBindingRevision)
	require.Equal(t, profile.group.Revision, provenance.Counter.ProfileRevision)
	require.Equal(t, f.definition.SeriesEpoch, mysqlCircuitStringCell(t, row, "series_epoch"))
	require.Equal(t, "kWh", mysqlCircuitStringCell(t, row, "unit"))
	if row.BucketStart.Equal(f.base) {
		require.Equal(t, "baseline", mysqlCircuitStringCell(t, row, "status"))
		assertMySQLCircuitCell(t, row, "value_decimal", mysqlCircuitDecimalBaseline)
		require.Equal(t, mysqlCircuitDecimalBaseline, provenance.Counter.RawText)
		require.Nil(t, mysqlCircuitCell(t, row, "delta_value"))
		require.Equal(t, mysqlCircuitDecimalBaselineID, provenance.SampleID)
		input := mysqlCircuitDecimalInputByID(t, f, mysqlCircuitDecimalBaselineID)
		assertMySQLCircuitTimeCell(t, row, "observed_at", f.base.Add(input.offset))
		assertMySQLCircuitTimeCell(t, row, "received_at", f.base.Add(input.offset+30*time.Millisecond))
		return
	}
	require.True(t, row.BucketStart.Equal(f.base.Add(mysqlCircuitCumulativeWindow)))
	require.Equal(t, "ok", mysqlCircuitStringCell(t, row, "status"))
	assertMySQLCircuitCell(t, row, "value_decimal", mysqlCircuitDecimalEnd)
	assertMySQLCircuitCell(t, row, "delta_value", "0.4")
	require.Equal(t, mysqlCircuitDecimalEnd, provenance.Counter.RawText)
	require.Equal(t, mysqlCircuitDecimalBaselineID, provenance.Counter.BaselineSampleID)
	require.Equal(t, mysqlCircuitDecimalEndID, provenance.Counter.EndSampleID)
	require.Equal(t, mysqlCircuitDecimalEndID, provenance.SampleID)
	input := mysqlCircuitDecimalInputByID(t, f, mysqlCircuitDecimalEndID)
	assertMySQLCircuitTimeCell(t, row, "observed_at", f.base.Add(input.offset))
	assertMySQLCircuitTimeCell(t, row, "received_at", f.base.Add(input.offset+30*time.Millisecond))
	assertMySQLCircuitTimeCell(t, row, "delta_from", f.base.Add(mysqlCircuitDecimalInputs[0].offset))
	assertMySQLCircuitTimeCell(t, row, "delta_to", f.base.Add(mysqlCircuitDecimalInputs[1].offset))
}

func assertMySQLCircuitDecimalSQLRows(t *testing.T, f *mysqlCircuitPlainFixture, outboxRows map[string]dbtarget.EncodedRow) {
	t.Helper()
	rows, err := f.f.target.QueryContext(t.Context(), `
		SELECT record_id, group_id, circuit_id, device_id, point_id, tag_id, measurement_id, series_epoch, unit,
			profile_kind, value_type, status, value_decimal, delta_value, bucket_start, observed_at, received_at,
			delta_from, delta_to, provenance
		FROM `+"`"+f.table+"`"+` ORDER BY group_id, bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	instantRows, cumulativeRows := 0, 0
	instantValues, missingInstantRows := 0, 0
	cumulativeWindow := 0
	seenRecords := make(map[string]struct{}, 82)
	for rows.Next() {
		var row mysqlCircuitDecimalSQLRow
		require.NoError(t, rows.Scan(&row.recordID, &row.groupID, &row.circuitID, &row.deviceID, &row.pointID, &row.tagID,
			&row.measurementID, &row.seriesEpoch, &row.unit, &row.profile, &row.valueType, &row.status, &row.raw,
			&row.delta, &row.bucketStart, &row.observedAt, &row.receivedAt, &row.deltaFrom, &row.deltaTo, &row.provenance))
		require.Equal(t, time.UTC, row.bucketStart.Location())
		require.Equal(t, "decimal", row.valueType)
		require.Equal(t, f.definition.ID, row.measurementID.String)
		require.True(t, row.measurementID.Valid)
		require.Equal(t, f.instant.group.Members[0].CircuitID, row.circuitID)
		require.Equal(t, f.instant.group.Members[0].DeviceID, row.deviceID)
		require.Equal(t, f.instant.group.Members[0].PointID, row.pointID)
		require.Equal(t, f.instant.group.Members[0].TagID, row.tagID)
		require.Equal(t, f.instant.group.Members[0].CircuitID, f.cumulative.group.Members[0].CircuitID)
		require.Equal(t, f.instant.group.Members[0].TagID, f.cumulative.group.Members[0].TagID)
		frozen, ok := outboxRows[row.recordID]
		require.True(t, ok, "each SQL row is the frozen record in its durable outbox")
		require.NotContains(t, seenRecords, row.recordID)
		seenRecords[row.recordID] = struct{}{}
		require.Equal(t, row.recordID, mysqlCircuitStringCell(t, frozen, "record_id"))
		require.True(t, row.bucketStart.Equal(frozen.BucketStart))
		var frozenProvenance mysqlCircuitDecimalProvenance
		require.NoError(t, json.Unmarshal([]byte(mysqlCircuitStringCell(t, frozen, "provenance")), &frozenProvenance))
		var sqlProvenance mysqlCircuitDecimalProvenance
		require.NoError(t, json.Unmarshal([]byte(row.provenance), &sqlProvenance))
		require.Equal(t, frozenProvenance, sqlProvenance, "SQL provenance matches the frozen outbox payload")

		if row.groupID == f.instant.id {
			require.Equal(t, "instantaneous", row.profile)
			require.False(t, row.seriesEpoch.Valid, "raw instantaneous records do not persist a counter window epoch")
			require.Equal(t, "second", frozenProvenance.Profile.Unit)
			require.Equal(t, 15, frozenProvenance.Profile.Interval)
			require.Equal(t, f.base.Add(time.Duration(instantRows)*15*time.Second), row.bucketStart)
			instantRows++
			if row.raw.Valid {
				instantValues++
				input := mysqlCircuitDecimalInputByID(t, f, frozenProvenance.SampleID)
				require.Equal(t, mysqlCircuitDecimalDigits(input.sampleID), mysqlCircuitStringCell(t, frozen, "value_decimal"))
				assertMySQLCircuitSQLDecimal(t, row.raw, mysqlCircuitDecimalDigits(input.sampleID))
				assertMySQLCircuitSQLReadingTimes(t, f, row.observedAt, row.receivedAt, input)
			} else {
				missingInstantRows++
				require.Equal(t, "missing", row.status)
				require.False(t, row.observedAt.Valid)
				require.False(t, row.receivedAt.Valid)
			}
			continue
		}
		require.Equal(t, f.cumulative.id, row.groupID)
		require.Equal(t, "cumulative", row.profile)
		require.True(t, row.seriesEpoch.Valid)
		require.Equal(t, f.definition.SeriesEpoch, row.seriesEpoch.String)
		require.True(t, row.unit.Valid)
		require.Equal(t, "kWh", row.unit.String)
		require.Equal(t, "minute", frozenProvenance.Profile.Unit)
		require.Equal(t, 10, frozenProvenance.Profile.Interval)
		require.Equal(t, f.base.Add(time.Duration(cumulativeWindow)*mysqlCircuitCumulativeWindow), row.bucketStart)
		if cumulativeWindow == 0 {
			require.Equal(t, "baseline", row.status)
			require.Equal(t, mysqlCircuitDecimalBaselineID, frozenProvenance.SampleID)
			assertMySQLCircuitSQLDecimal(t, row.raw, mysqlCircuitDecimalBaseline)
			require.False(t, row.delta.Valid)
			require.False(t, row.deltaFrom.Valid)
			require.False(t, row.deltaTo.Valid)
			input := mysqlCircuitDecimalInputByID(t, f, mysqlCircuitDecimalBaselineID)
			assertMySQLCircuitSQLReadingTimes(t, f, row.observedAt, row.receivedAt, input)
		} else {
			require.Equal(t, "ok", row.status)
			require.Equal(t, mysqlCircuitDecimalEndID, frozenProvenance.SampleID)
			assertMySQLCircuitSQLDecimal(t, row.raw, mysqlCircuitDecimalEnd)
			assertMySQLCircuitSQLDecimal(t, row.delta, "0.4")
			input := mysqlCircuitDecimalInputByID(t, f, mysqlCircuitDecimalEndID)
			assertMySQLCircuitSQLReadingTimes(t, f, row.observedAt, row.receivedAt, input)
			require.True(t, row.deltaFrom.Valid)
			require.True(t, row.deltaTo.Valid)
			require.Equal(t, time.UTC, row.deltaFrom.Time.Location())
			require.Equal(t, time.UTC, row.deltaTo.Time.Location())
			require.True(t, row.deltaFrom.Time.Equal(f.base.Add(mysqlCircuitDecimalInputs[0].offset)))
			require.True(t, row.deltaTo.Time.Equal(f.base.Add(mysqlCircuitDecimalInputs[1].offset)))
		}
		cumulativeRows++
		cumulativeWindow++
	}
	require.NoError(t, rows.Err())
	require.Equal(t, 80, instantRows)
	require.Equal(t, 2, instantValues)
	require.Equal(t, 78, missingInstantRows)
	require.Equal(t, 2, cumulativeRows)
	require.Len(t, seenRecords, 82)
}

func assertMySQLCircuitSQLDecimal(t *testing.T, got sql.NullString, want string) {
	t.Helper()
	require.True(t, got.Valid)
	actual, err := measurement.NewDecimal(got.String)
	require.NoError(t, err)
	expected, err := measurement.NewDecimal(want)
	require.NoError(t, err)
	actualRat, ok := actual.Rat()
	require.True(t, ok)
	expectedRat, ok := expected.Rat()
	require.True(t, ok)
	require.Zero(t, actualRat.Cmp(expectedRat), "DECIMAL readback is exactly %s", want)
	actualCanonical, ok := actual.CanonicalDigits()
	require.True(t, ok)
	expectedCanonical, ok := expected.CanonicalDigits()
	require.True(t, ok)
	require.Equal(t, expectedCanonical, actualCanonical)
}

func assertMySQLCircuitSQLReadingTimes(t *testing.T, f *mysqlCircuitPlainFixture, observed, received sql.NullTime, input mysqlCircuitInput) {
	t.Helper()
	expectedObserved := f.base.Add(input.offset)
	expectedReceived := expectedObserved.Add(30 * time.Millisecond)
	require.True(t, observed.Valid)
	require.True(t, received.Valid)
	require.Equal(t, time.UTC, observed.Time.Location())
	require.Equal(t, time.UTC, received.Time.Location())
	require.True(t, expectedObserved.Equal(observed.Time))
	require.True(t, expectedReceived.Equal(received.Time))
}

func assertMySQLCircuitTimeCell(t *testing.T, row dbtarget.EncodedRow, column string, expected time.Time) {
	t.Helper()
	value := mysqlCircuitCell(t, row, column)
	actual, ok := value.(time.Time)
	require.True(t, ok, "%s must be a frozen UTC time cell", column)
	require.Equal(t, time.UTC, actual.Location())
	require.True(t, expected.Equal(actual))
}

func assertMySQLCircuitCell(t *testing.T, row dbtarget.EncodedRow, column string, expected any) {
	t.Helper()
	require.Equal(t, expected, mysqlCircuitCell(t, row, column), column)
}

func mysqlCircuitCell(t *testing.T, row dbtarget.EncodedRow, column string) any {
	t.Helper()
	for _, cell := range row.Cells {
		if cell.Column == column {
			return cell.Value
		}
	}
	require.FailNow(t, "frozen outbox row lacks a required cell", column)
	return nil
}

func mysqlCircuitStringCell(t *testing.T, row dbtarget.EncodedRow, column string) string {
	t.Helper()
	value := mysqlCircuitCell(t, row, column)
	text, ok := value.(string)
	require.True(t, ok, "%s must be a string cell", column)
	return text
}

func mysqlCircuitDecimalInputByID(t *testing.T, f *mysqlCircuitPlainFixture, sampleID string) mysqlCircuitInput {
	t.Helper()
	for _, input := range f.inputs {
		if input.sampleID == sampleID {
			return input
		}
	}
	require.FailNow(t, "unexpected decimal fixture sample", sampleID)
	return mysqlCircuitInput{}
}

func mysqlCircuitDecimalInputAt(f *mysqlCircuitPlainFixture, bucket time.Time) (mysqlCircuitInput, bool) {
	for _, input := range f.inputs {
		if f.base.Add(input.offset).Equal(bucket) {
			return input, true
		}
	}
	return mysqlCircuitInput{}, false
}

func mysqlCircuitDecimalDigits(sampleID string) string {
	if sampleID == mysqlCircuitDecimalBaselineID {
		return mysqlCircuitDecimalBaseline
	}
	return mysqlCircuitDecimalEnd
}
