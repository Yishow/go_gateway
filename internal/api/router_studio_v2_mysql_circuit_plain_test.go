package api

import (
	"database/sql"
	"encoding/json"
	"slices"
	"strconv"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/measurement"

	"github.com/stretchr/testify/require"
)

type mysqlCircuitPlainSQLRow struct {
	groupID, profile, valueType, status string
	raw, delta                          sql.NullString
	bucketStart                         time.Time
	deltaFrom, deltaTo                  sql.NullTime
}

type mysqlCircuitPlainJournalPayload struct {
	SampleID   string                 `json:"sample_id"`
	ObservedAt string                 `json:"observed_at"`
	ReceivedAt string                 `json:"received_at"`
	Value      measurement.ExactValue `json:"value"`
}

// TestMySQLCircuitPlainIntegration exercises the selected MySQL 8.4 uint64
// case end to end. It seeds only the
// fixture-side snapshots needed to run the production pipeline.
func TestMySQLCircuitPlainIntegration(t *testing.T) {
	runMySQLCircuitPlainIntegration(t, mysqlCircuitFamilyMySQL)
}

func runMySQLCircuitPlainIntegration(t *testing.T, family string) {
	f := newMySQLCircuitPlainFixture(t, mysqlCircuitPlainFixtureOptions{Family: family})

	f.accept(0)
	f.closeAfter(10 * time.Minute)
	baselineBeforeRetry := mysqlCircuitBaseline(t, f)
	requireEventuallyCumulativeRetry(t, f)
	waitForMySQLCircuitOutboxes(t, f, 41)
	require.Equal(t, baselineBeforeRetry, mysqlCircuitBaseline(t, f), "delivery retries cannot advance a committed counter baseline")

	// The empty ten-minute window is durably recorded but must not move the
	// checkpoint; the following sample therefore closes as an exact gap span.
	f.closeAfter(20 * time.Minute)
	require.Equal(t, baselineBeforeRetry, mysqlCircuitBaseline(t, f), "a missing counter window cannot advance its baseline")
	waitForMySQLCircuitOutboxes(t, f, 82)

	for i := range mysqlCircuitInputs[1:] {
		inputIndex := i + 1
		f.accept(inputIndex)
		f.closeAfter(time.Duration((inputIndex+2)*10) * time.Minute)
	}
	waitForMySQLCircuitOutboxes(t, f, 205)

	assertMySQLCircuitJournal(t, f)
	outboxDigests := assertMySQLCircuitOutboxAndReceipts(t, f)
	assertMySQLCircuitColumnTypes(t, f)
	assertMySQLCircuitSQLRows(t, f)
	require.Len(t, outboxDigests, 205)
}

func mysqlCircuitBaseline(t *testing.T, f *mysqlCircuitPlainFixture) []groupdelivery.CounterBaseline {
	t.Helper()
	restored, err := f.store.Restore(t.Context(), f.cumulative.key)
	require.NoError(t, err)
	require.Len(t, restored.Baselines, 1)
	return restored.Baselines
}

func requireEventuallyCumulativeRetry(t *testing.T, f *mysqlCircuitPlainFixture) {
	t.Helper()
	require.Eventually(t, func() bool {
		var retries int
		err := f.f.db.QueryRowContext(t.Context(), `
			SELECT COALESCE(MAX(retry_count),0) FROM wg_delivery_outbox WHERE group_id=?`, f.cumulative.id).Scan(&retries)
		return err == nil && retries > 0
	}, 15*time.Second, 25*time.Millisecond, "the cumulative row must exercise a real destination retry")
}

func waitForMySQLCircuitOutboxes(t *testing.T, f *mysqlCircuitPlainFixture, expected int) {
	t.Helper()
	require.Eventually(t, func() bool {
		var total, unfinished int
		err := f.f.db.QueryRowContext(t.Context(), `
			SELECT COUNT(*), SUM(CASE WHEN state <> 'sql_committed' THEN 1 ELSE 0 END)
			FROM wg_delivery_outbox WHERE group_id IN (?,?)`, f.instant.id, f.cumulative.id).Scan(&total, &unfinished)
		return err == nil && total == expected && unfinished == 0
	}, 20*time.Second, 25*time.Millisecond, "all closed rows must reach a durable SQL receipt")
}

func assertMySQLCircuitJournal(t *testing.T, f *mysqlCircuitPlainFixture) {
	t.Helper()
	counts := map[string]int{f.instant.id: 0, f.cumulative.id: 0}
	for _, entry := range f.acceptedJournal {
		groupID, sampleID, payload, digest := entry.groupID, entry.sampleID, entry.payload, entry.digest
		require.Contains(t, counts, groupID)
		require.Zero(t, entry.consumed, "journal evidence is captured before window closure")
		require.Len(t, digest, 64)
		var decoded mysqlCircuitPlainJournalPayload
		require.NoError(t, json.Unmarshal([]byte(payload), &decoded))
		require.Equal(t, sampleID, decoded.SampleID)
		input, ok := mysqlCircuitInputByID(sampleID)
		require.True(t, ok, "only the four fixture-selected samples are journaled")
		require.Equal(t, measurement.ExactUint64, decoded.Value.Type())
		digits, ok := decoded.Value.DecimalDigits()
		require.True(t, ok)
		require.Equal(t, input.value, parseMySQLCircuitUint64(t, digits))
		observed := f.base.Add(input.offset)
		require.True(t, observed.Equal(parseMySQLCircuitTime(t, decoded.ObservedAt)))
		require.True(t, observed.Add(30*time.Millisecond).Equal(parseMySQLCircuitTime(t, decoded.ReceivedAt)))
		counts[groupID]++
	}
	require.Equal(t, map[string]int{f.instant.id: 4, f.cumulative.id: 4}, counts)
}

func mysqlCircuitInputByID(sampleID string) (mysqlCircuitInput, bool) {
	index := slices.IndexFunc(mysqlCircuitInputs[:], func(input mysqlCircuitInput) bool {
		return input.sampleID == sampleID
	})
	if index < 0 {
		return mysqlCircuitInput{}, false
	}
	return mysqlCircuitInputs[index], true
}

func parseMySQLCircuitUint64(t *testing.T, digits string) uint64 {
	t.Helper()
	value, err := strconv.ParseUint(digits, 10, 64)
	require.NoError(t, err)
	return value
}

func parseMySQLCircuitTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339Nano, value)
	require.NoError(t, err)
	return parsed.UTC()
}

func assertMySQLCircuitOutboxAndReceipts(t *testing.T, f *mysqlCircuitPlainFixture) map[string]string {
	t.Helper()
	rows, err := f.f.db.QueryContext(t.Context(), `
		SELECT effect_key, group_id, payload, payload_digest, state, retry_count
		FROM wg_delivery_outbox WHERE group_id IN (?,?) ORDER BY group_id, bucket_start`, f.instant.id, f.cumulative.id)
	require.NoError(t, err)
	defer rows.Close()
	digests := make(map[string]string, 205)
	cumulativeRetried := false
	for rows.Next() {
		var effectKey, groupID, payload, digest, state string
		var retries int
		require.NoError(t, rows.Scan(&effectKey, &groupID, &payload, &digest, &state, &retries))
		require.Equal(t, "sql_committed", state)
		if groupID == f.cumulative.id && retries > 0 {
			cumulativeRetried = true
		}
		row, err := groupdelivery.DecodeRowPayload([]byte(payload))
		require.NoError(t, err)
		require.Equal(t, effectKey, row.EffectKey)
		actualDigest, err := groupdelivery.RowPayloadDigest(row)
		require.NoError(t, err)
		require.Equal(t, digest, actualDigest, "the durable outbox digest covers its closed payload")
		digests[effectKey] = digest
	}
	require.NoError(t, rows.Err())
	require.NoError(t, rows.Close())
	require.Len(t, digests, 205)
	require.True(t, cumulativeRetried, "the cumulative row was retried without changing its checkpoint")

	localRows, err := f.f.db.QueryContext(t.Context(), `SELECT effect_key, payload_digest FROM wg_delivery_receipts`)
	require.NoError(t, err)
	defer localRows.Close()
	localReceipts := make(map[string]string, len(digests))
	for localRows.Next() {
		var effectKey, digest string
		require.NoError(t, localRows.Scan(&effectKey, &digest))
		localReceipts[effectKey] = digest
	}
	require.NoError(t, localRows.Err())
	require.NoError(t, localRows.Close())
	require.Equal(t, digests, localReceipts, "each local ACK is bound to the frozen outbox digest")

	remoteRows, err := f.f.target.QueryContext(t.Context(), `
		SELECT effect_key, payload_digest, committed_at, _gw_effect_receipts_v1
		FROM `+"`"+dbtarget.EffectReceiptTable+"`"+` ORDER BY effect_key`)
	require.NoError(t, err)
	defer remoteRows.Close()
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
	require.Equal(t, digests, remoteReceipts, "every committed outbox row has a matching destination receipt")
	return digests
}

func assertMySQLCircuitColumnTypes(t *testing.T, f *mysqlCircuitPlainFixture) {
	t.Helper()
	rows, err := f.f.target.QueryContext(t.Context(), `
		SELECT COLUMN_NAME, DATA_TYPE, COLUMN_TYPE FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME=?
		AND COLUMN_NAME IN ('value_decimal','delta_value','value_float','value_type')`, f.table)
	require.NoError(t, err)
	defer rows.Close()
	types := make(map[string][2]string, 4)
	for rows.Next() {
		var name, dataType, columnType string
		require.NoError(t, rows.Scan(&name, &dataType, &columnType))
		types[name] = [2]string{dataType, columnType}
	}
	require.NoError(t, rows.Err())
	require.Equal(t, [2]string{"decimal", "decimal(65,30)"}, types["value_decimal"])
	require.Equal(t, [2]string{"decimal", "decimal(65,30)"}, types["delta_value"])
	require.Equal(t, "double", types["value_float"][0])
	require.Equal(t, "varchar", types["value_type"][0])
}

func assertMySQLCircuitSQLRows(t *testing.T, f *mysqlCircuitPlainFixture) {
	t.Helper()
	require.Regexp(t, longTablePattern, f.table)
	rows, err := f.f.target.QueryContext(t.Context(), `
		SELECT group_id, profile_kind, value_type, status,
			CAST(value_decimal AS DECIMAL(65,0)), CAST(delta_value AS DECIMAL(65,0)),
			bucket_start, delta_from, delta_to
		FROM `+"`"+f.table+"`"+` ORDER BY group_id, bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	var stored []mysqlCircuitPlainSQLRow
	for rows.Next() {
		var row mysqlCircuitPlainSQLRow
		require.NoError(t, rows.Scan(&row.groupID, &row.profile, &row.valueType, &row.status,
			&row.raw, &row.delta, &row.bucketStart, &row.deltaFrom, &row.deltaTo))
		stored = append(stored, row)
	}
	require.NoError(t, rows.Err())
	require.Len(t, stored, 205, "the ten-minute jump closes forty 15-second instantaneous rows per window plus cumulative rows")
	var instant, cumulative []mysqlCircuitPlainSQLRow
	for _, row := range stored {
		switch row.groupID {
		case f.instant.id:
			instant = append(instant, row)
		case f.cumulative.id:
			cumulative = append(cumulative, row)
		default:
			require.FailNow(t, "unexpected SQL group", row.groupID)
		}
	}
	require.Len(t, instant, 200)
	wantRaw := []string{"9007199254740992", "9007199254740993", "18446744073709551615", "7"}
	var instantaneousValues []string
	for _, row := range instant {
		require.Equal(t, "instantaneous", row.profile)
		require.Equal(t, "uint64", row.valueType)
		require.False(t, row.delta.Valid)
		if row.raw.Valid {
			require.Equal(t, "ok", row.status)
			instantaneousValues = append(instantaneousValues, row.raw.String)
		} else {
			require.Equal(t, "missing", row.status)
		}
	}
	require.Equal(t, wantRaw, instantaneousValues)

	require.Len(t, cumulative, 5)
	require.Equal(t, "cumulative", cumulative[0].profile)
	require.Equal(t, "uint64", cumulative[0].valueType)
	require.Equal(t, "baseline", cumulative[0].status)
	require.Equal(t, "9007199254740992", cumulative[0].raw.String)
	require.False(t, cumulative[0].delta.Valid)

	require.Equal(t, "missing", cumulative[1].status)
	require.False(t, cumulative[1].raw.Valid)
	require.False(t, cumulative[1].delta.Valid)

	require.Equal(t, "gap", cumulative[2].status)
	require.Equal(t, "9007199254740993", cumulative[2].raw.String)
	require.Equal(t, "1", cumulative[2].delta.String)
	require.True(t, cumulative[2].deltaFrom.Valid)
	require.True(t, cumulative[2].deltaTo.Valid)
	require.True(t, cumulative[2].deltaFrom.Time.Equal(f.base.Add(mysqlCircuitInputs[0].offset)))
	require.True(t, cumulative[2].deltaTo.Time.Equal(f.base.Add(mysqlCircuitInputs[1].offset)))

	require.Equal(t, "ok", cumulative[3].status)
	require.Equal(t, "18446744073709551615", cumulative[3].raw.String)
	require.Equal(t, "18437736874454810622", cumulative[3].delta.String)

	require.Equal(t, "reset_or_uncertain", cumulative[4].status)
	require.Equal(t, "7", cumulative[4].raw.String)
	require.False(t, cumulative[4].delta.Valid)
}
