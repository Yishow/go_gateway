package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

type managedProvenance struct {
	Member     string `json:"member"`
	Status     string `json:"status"`
	Reason     string `json:"reason"`
	SampleID   string `json:"sample_id"`
	ObservedAt string `json:"observed_at"`
	Quality    string `json:"quality"`
}

func managedSQLTable(f *managedPipelineTarget, table string) string {
	return `"` + f.namespace + `"."` + table + `"`
}

func assertManagedSQL(t *testing.T, f *managedPipelineTarget, group *workspace.WriteGroup, base time.Time) {
	t.Helper()
	columns := make([]string, 0, 5+len(group.Members))
	columns = append(columns, "record_id", "group_id", "device_id", "bucket_start", "provenance")
	for _, member := range group.Members {
		columns = append(columns, `CAST("`+member.TargetColumn+`" AS TEXT)`)
	}
	rows, err := f.db.QueryContext(t.Context(), `SELECT `+strings.Join(columns, ",")+` FROM `+managedSQLTable(f, group.Destination.TableName)+` ORDER BY bucket_start`)
	require.NoError(t, err)
	defer rows.Close()
	wantValues := []string{"1", "-32768", "-2147483648", "-9223372036854775808", "65535", "4294967295", "", "1.25", "2.5", "同名保留字 select"}
	if f.kind == schema.DatabaseConnectorKindPostgres {
		wantValues[0] = "true"
	}
	bigValues := []string{"9007199254740993", "9223372036854775808", "18446744073709551615"}
	seen := make(map[string]bool)
	count := 0
	for rows.Next() {
		var record, groupID, device, bucket, rawProvenance string
		values := make([]string, len(group.Members))
		args := make([]any, 0, 5+len(values))
		args = append(args, &record, &groupID, &device, &bucket, &rawProvenance)
		for i := range values {
			args = append(args, &values[i])
		}
		require.NoError(t, rows.Scan(args...))
		require.Less(t, count, len(bigValues))
		require.NotEmpty(t, record)
		require.False(t, seen[record], "each closed bucket keeps its own record identity")
		seen[record] = true
		require.Equal(t, group.ID, groupID)
		require.Equal(t, "device-1", device)
		start := base.Add(time.Duration(count) * 10 * time.Second)
		require.Equal(t, start.Format(time.RFC3339Nano), bucket, "SQL stores acquisition time rather than flush time")
		wantValues[6] = bigValues[count]
		require.Equal(t, wantValues, values, "SQL text readback retains every integer digit")
		var provenance []managedProvenance
		require.NoError(t, json.Unmarshal([]byte(rawProvenance), &provenance))
		require.Len(t, provenance, len(group.Members))
		for i, entry := range provenance {
			member := group.Members[i]
			identity, err := json.Marshal([]string{member.DeviceID, member.PointID, member.TagID})
			require.NoError(t, err)
			require.Equal(t, string(identity), entry.Member)
			require.Equal(t, "ok", entry.Status)
			require.Empty(t, entry.Reason)
			require.Equal(t, "good", entry.Quality)
			require.Equal(t, "managed-"+strconv.Itoa(count)+"-"+strconv.Itoa(i), entry.SampleID)
			require.Equal(t, start.Add(2*time.Second+123*time.Nanosecond).Format(time.RFC3339Nano), entry.ObservedAt)
		}
		var localRecord, localBucket string
		require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT record_id,bucket_start FROM wg_delivery_outbox WHERE record_id=?`, record).Scan(&localRecord, &localBucket))
		require.Equal(t, localRecord, record)
		parsedBucket, err := time.Parse(time.RFC3339Nano, localBucket)
		require.NoError(t, err)
		require.Equal(t, start, parsedBucket)
		count++
	}
	require.NoError(t, rows.Err())
	require.Equal(t, 3, count)
	require.NoError(t, rows.Close())
	assertManagedReceipts(t, f)
}

func assertManagedReceipts(t *testing.T, f *managedPipelineTarget) {
	t.Helper()
	rows, err := f.db.QueryContext(t.Context(), `SELECT effect_key,payload_digest,committed_at FROM `+managedSQLTable(f, "gw_effect_receipts"))
	require.NoError(t, err)
	defer rows.Close()
	count := 0
	for rows.Next() {
		var effect, digest, committed string
		require.NoError(t, rows.Scan(&effect, &digest, &committed))
		require.Len(t, digest, 64)
		_, err := time.Parse(time.RFC3339Nano, committed)
		require.NoError(t, err)
		var localDigest, state string
		require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT r.payload_digest,o.state FROM wg_delivery_receipts r JOIN wg_delivery_outbox o USING(effect_key) WHERE effect_key=?`, effect).Scan(&localDigest, &state))
		require.Equal(t, digest, localDigest)
		require.Equal(t, groupdelivery.StateCommitted, state)
		count++
	}
	require.NoError(t, rows.Err())
	require.Equal(t, 3, count)
	var localCount int
	require.NoError(t, f.env.db.QueryRowContext(t.Context(), `SELECT count(*) FROM wg_delivery_receipts`).Scan(&localCount))
	require.Equal(t, count, localCount)
}
