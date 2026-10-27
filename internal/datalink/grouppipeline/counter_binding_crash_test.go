package grouppipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/groupdelivery"
	"go-gateway/internal/datalink/workspace"
	"go-gateway/internal/testutil/lostcommit"
)

func crashCounterJournal(t *testing.T, f *longPipeline) [][3]string {
	t.Helper()
	rows, err := f.local.Query(`SELECT sample_id,payload,payload_digest FROM wg_delivery_samples ORDER BY seq`)
	require.NoError(t, err)
	defer rows.Close()
	var result [][3]string
	for rows.Next() {
		var r [3]string
		require.NoError(t, rows.Scan(&r[0], &r[1], &r[2]))
		result = append(result, r)
	}
	require.NoError(t, rows.Err())
	return result
}

// This matrix models lost commit replies and a discarded process's runtime
// memory, not power-loss hardware. Existing lostcommit executes real SQLite
// commits; reopening the local file proves recovery is not an in-memory retry.
func TestCounterBindingCrashResendAcrossJournalFormats(t *testing.T) {
	for _, format := range []string{"versioned-native", "versioned-exact", "legacy-exact"} {
		for _, stage := range []string{"sample-before", "sample-after", "sample-acked", "closure-before", "closure-after"} {
			t.Run(format+"/"+stage, func(t *testing.T) {
				scaled, legacy := format != "versioned-native", format == "legacy-exact"
				f := newLongPipelineOn(t, nil, scaled)
				key := groupdelivery.GroupKey{WorkspaceID: "workspace-A", GroupID: "group-C", GroupRevision: "rev-1"}
				descriptor, err := f.store.RuntimeVersion(t.Context(), key)
				require.NoError(t, err)
				if legacy {
					// Original legacy descriptor construction, before any accepted input.
					// The complete exact-source evidence remains; no current definition helps.
					var old workspace.WriteGroupAppliedSnapshot
					require.NoError(t, json.Unmarshal(descriptor, &old))
					old.CounterBindings = nil
					descriptor, err = json.Marshal(old)
					require.NoError(t, err)
					_, err = f.local.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_revision='rev-1'`, string(descriptor))
					require.NoError(t, err)
				}
				reader := &unavailableExactCounterReader{}
				restart := func() {
					require.NoError(t, f.local.Close())
					f.local = lostcommit.OpenFaults(f.localPath+"?_pragma=busy_timeout(5000)", f.localFault)
					f.store = groupdelivery.NewStore(f.local)
					f.pipe = New(Dependencies{Groups: f.groups, Destinations: f.dest, Store: f.store, Counters: reader}, Config{Now: func() time.Time { return f.clock }, Delivery: groupdelivery.DispatcherConfig{MaxPartitions: 1}, OnError: func(err error) { f.errs = append(f.errs, err) }})
					require.NoError(t, f.pipe.Reconcile(t.Context()))
				}
				accept := func(id string, observed time.Time, value uint64) error {
					f.clock = observed.Add(time.Second)
					if !legacy {
						return f.reading(0, id, observed, value)
					}
					return appendPreUpgradeCounterFixture(t, f, id, observed)
				}
				first, second := uint64(100), uint64(105)
				if scaled {
					first, second = 1000, 1050
				}
				if stage == "sample-before" {
					hit := false
					reject := func(query string) error {
						if strings.Contains(query, "INSERT INTO wg_delivery_samples") {
							hit = true
							return errors.New("injected precommit sample failure")
						}
						return nil
					}
					f.localFault.Reject.Store(&reject)
					require.Error(t, accept("c-1", at(1), first))
					require.True(t, hit)
					f.localFault.Reject.Store(nil)
					require.Empty(t, crashCounterJournal(t, f))
					require.Zero(t, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_checkpoints`))
					require.Zero(t, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_counter_baselines`))
				}
				if stage == "sample-after" {
					f.localFault.Lose.Store(true)
				}
				err = accept("c-1", at(1), first)
				if stage == "sample-after" {
					require.ErrorIs(t, err, lostcommit.ErrResponseLost)
					f.localFault.Lose.Store(false)
				} else {
					require.NoError(t, err)
				}
				originalFirst := crashCounterJournal(t, f)
				require.Len(t, originalFirst, 1, "the real commit happened, whether its ACK arrived or not")
				if legacy {
					const payload = `{"sample_id":"c-1","member_key":"[\"device-1\",\"point-1\",\"tag-kwh-1\"]","observed_at":"2026-01-01T00:00:01Z","quality":"good","value":{"type":"decimal","encoding":"decimal-string","value":"100.0"},"source_revision":"src-1","mapping_revision":"map-1"}`
					require.Equal(t, payload, originalFirst[0][1])
					sum := sha256.Sum256([]byte(payload))
					require.Equal(t, hex.EncodeToString(sum[:]), originalFirst[0][2])
				}
				if legacy {
					require.NoError(t, accept("c-2", at(601), second))
				}
				beforeUpgrade := crashCounterJournal(t, f)
				restart()
				// New reader accepts the identical raw resend only after original evidence
				// resolution (and durable association for legacy), never rewrites d-1.
				require.NoError(t, f.reading(0, "c-1", at(1), first))
				require.Equal(t, beforeUpgrade, crashCounterJournal(t, f))
				if legacy {
					require.Equal(t, 2, f.count(f.local, `SELECT COUNT(*) FROM wg_counter_proof_associations`))
				}
				if !legacy {
					require.NoError(t, accept("c-2", at(601), second))
				}
				originals := crashCounterJournal(t, f)
				require.Len(t, originals, 2)
				restart()
				require.NoError(t, f.reading(0, "c-2", at(601), second))
				require.Equal(t, originals, crashCounterJournal(t, f))
				before, err := f.store.Restore(t.Context(), key)
				require.NoError(t, err)
				require.Len(t, before.Samples, 2)
				if stage == "closure-before" {
					f.failClosureBaselineWrite()
				}
				if stage == "closure-after" {
					f.localFault.Lose.Store(true)
				}
				closeAt := at(1200)
				if stage == "closure-after" {
					closeAt = at(602)
				}
				f.closeUntil(closeAt)
				f.localFault.Reject.Store(nil)
				f.localFault.Lose.Store(false)
				if stage == "closure-before" {
					require.NotEmpty(t, f.errs)
					after, err := f.store.Restore(t.Context(), key)
					require.NoError(t, err)
					require.Equal(t, before, after)
					require.Empty(t, f.outboxDigests())
				}
				if stage == "closure-after" {
					require.NotEmpty(t, f.errs)
					require.Len(t, f.outboxDigests(), 2, "first bucket committed before its reply was lost")
					after, err := f.store.Restore(t.Context(), key)
					require.NoError(t, err)
					require.Len(t, after.Samples, 1)
					require.Len(t, after.Baselines, 1)
				}
				committed := f.outboxDigests()
				restart()
				f.closeUntil(at(1200))
				final := f.outboxDigests()
				require.Len(t, final, 4)
				for id, digest := range committed {
					require.Equal(t, digest, final[id], "finalized output is never recomputed")
				}
				require.Equal(t, originals, crashCounterJournal(t, f))
				require.Equal(t, 2, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_samples WHERE consumed=1`))
				if legacy {
					require.Equal(t, 2, f.count(f.local, `SELECT COUNT(*) FROM wg_counter_proof_associations`))
				}
				restored, err := f.store.Restore(t.Context(), key)
				require.NoError(t, err)
				require.Empty(t, restored.Samples)
				require.Len(t, restored.Baselines, 1)
				f.deliver()
				require.Equal(t, 4, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
				require.Equal(t, 4, f.count(f.local, `SELECT COUNT(*) FROM wg_delivery_receipts`))
				row := f.targetRows()[windowKey(longKwhFirst, 1)]
				wantRaw, wantDelta := "105", "5"
				if scaled {
					wantRaw, wantDelta = "105.0", "5.0"
				}
				require.Equal(t, wantRaw, row.raw.String)
				require.Equal(t, wantDelta, row.delta.String)
				require.Contains(t, row.provenance.String, `"definition_revision":"def-1"`)
				var unit, epoch string
				require.NoError(t, f.plainTarget().QueryRow(`SELECT unit,series_epoch FROM `+f.table+` WHERE tag_id=? AND bucket_start=?`, longKwhFirst, at(600).Format(longTime)).Scan(&unit, &epoch))
				require.Equal(t, "kWh", unit)
				require.Equal(t, "epoch-1", epoch)
				restart()
				f.closeUntil(at(1200))
				f.deliver()
				require.Equal(t, final, f.outboxDigests())
				require.Equal(t, 4, f.count(f.plainTarget(), `SELECT COUNT(*) FROM `+f.table))
				retained, err := f.store.RuntimeVersion(t.Context(), key)
				require.NoError(t, err)
				require.Equal(t, descriptor, retained)
				require.Zero(t, reader.calls)
			})
		}
	}
}

// appendPreUpgradeCounterFixture reproduces the old SQL transaction with a
// literal old wire format. The current AppendSample intentionally refuses NEW
// counter inputs without a reference; this helper is never used after upgrade.
func appendPreUpgradeCounterFixture(t *testing.T, f *longPipeline, id string, observed time.Time) error {
	t.Helper()
	digits := "100.0"
	if id == "c-2" {
		digits = "105.0"
	}
	payload := fmt.Sprintf(`{"sample_id":"%s","member_key":"[\"device-1\",\"point-1\",\"tag-kwh-1\"]","observed_at":"%s","quality":"good","value":{"type":"decimal","encoding":"decimal-string","value":"%s"},"source_revision":"src-1","mapping_revision":"map-1"}`, id, observed.Format(time.RFC3339Nano), digits)
	sum := sha256.Sum256([]byte(payload))
	tx, err := f.local.BeginTx(t.Context(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	const stamp = "2006-01-02T15:04:05.000000000Z"
	_, err = tx.ExecContext(t.Context(), `INSERT INTO wg_delivery_samples (workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,created_at) VALUES (?,?,?,?,?,?,?,?,?,?)`, "workspace-A", "group-C", "rev-1", id, `["device-1","point-1","tag-kwh-1"]`, observed.Format(stamp), observed.Truncate(longWindow).Format(stamp), payload, hex.EncodeToString(sum[:]), observed.Add(time.Second).Format(stamp))
	if err != nil {
		return err
	}
	return tx.Commit()
}
