package grouppipeline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

// sourceEvidenceJSON is the scoped delivery response's per source evidence of the applied revision exactly as the
// API serializes it, or nil when the group has none. The assertions below read the wire form, so they pin what an
// operator's client can see rather than a Go shape.
func (f *longPipeline) sourceEvidenceJSON() (document map[string]any, raw string) {
	f.t.Helper()
	view, err := f.pipe.DeliveryForRevision(f.t.Context(), f.group.ID, f.group.AppliedRevision)
	require.NoError(f.t, err)
	encoded, err := json.Marshal(view.SourceEvidence)
	require.NoError(f.t, err)
	require.NoError(f.t, json.Unmarshal(encoded, &document))
	return document, string(encoded)
}

func recordOf(t *testing.T, source any, which string) map[string]any {
	t.Helper()
	record, _ := source.(map[string]any)[which].(map[string]any)
	return record
}

func evidenceSources(t *testing.T, document map[string]any) []any {
	t.Helper()
	require.NotNil(t, document, "the group's applied revision has per source evidence")
	require.Equal(t, "available", document["state"])
	sources, _ := document["sources"].([]any)
	return sources
}

// closedAndDelivered closes window 0 and delivers it, so every record is SQL committed with its receipt.
func closedAndDelivered(t *testing.T, shape averageShape) *longPipeline {
	t.Helper()
	f := newAveragePipeline(t, shape)
	f.feed(shape, 0)
	f.closeUntil(shape.windowEnd(0))
	f.deliver()
	return f
}

// tamperPayloads rewrites the frozen payload of every outbox row with the digests and receipts kept consistent, so
// only the metadata inside it is wrong.
func (f *longPipeline) tamperPayloads(old, replacement string) {
	f.t.Helper()
	local := f.plainLocal()
	rows, err := local.QueryContext(f.t.Context(), `SELECT effect_key, payload FROM wg_delivery_outbox`)
	require.NoError(f.t, err)
	type row struct{ key, payload string }
	var found []row
	for rows.Next() {
		var r row
		require.NoError(f.t, rows.Scan(&r.key, &r.payload))
		found = append(found, r)
	}
	require.NoError(f.t, rows.Err())
	require.NoError(f.t, rows.Close())
	require.NotEmpty(f.t, found)
	changed := 0
	for _, r := range found {
		next := strings.ReplaceAll(r.payload, old, replacement)
		if next == r.payload {
			continue
		}
		changed++
		sum := sha256.Sum256([]byte(next))
		digest := hex.EncodeToString(sum[:])
		_, err := local.ExecContext(f.t.Context(), `UPDATE wg_delivery_outbox SET payload = ?, payload_digest = ? WHERE effect_key = ?`, next, digest, r.key)
		require.NoError(f.t, err)
		_, err = local.ExecContext(f.t.Context(), `UPDATE wg_delivery_receipts SET payload_digest = ? WHERE effect_key = ?`, digest, r.key)
		require.NoError(f.t, err)
	}
	require.Positive(f.t, changed, "the corruption must really reach a stored payload: %q", old)
}

// TestAverageDeliveryEvidence pins the value-free, read-only projection of the average and counter results through the
// real durable store: the pipeline closes real windows, the destination commits them, and the scoped delivery view is
// read back. It must serve both new managed layouts from the frozen applied layout, payload and matching receipt, keep
// source and result types apart, keep unknown coverage unknown, and never reveal a value, a sum or a difference.
func TestAverageDeliveryEvidence(t *testing.T) {
	for _, shape := range []averageShape{shapeMeanLong, shapeMeanWide} {
		t.Run(string(shape), func(t *testing.T) {
			t.Run("each source reports method, counts, separate types and unknown coverage without any value", func(t *testing.T) {
				f := closedAndDelivered(t, shape)
				document, raw := f.sourceEvidenceJSON()
				sources := evidenceSources(t, document)
				require.Len(t, sources, 2)

				for _, source := range sources {
					for _, which := range []string{"latest", "last_sql_committed"} {
						record := recordOf(t, source, which)
						require.NotNil(t, record, which)
						require.Equal(t, true, record["available"], which)
						aggregation, _ := record["aggregation"].(map[string]any)
						require.NotNil(t, aggregation, "%s carries its verified aggregation evidence", which)
						require.Equal(t, "arithmetic-mean", aggregation["method"])
						require.Equal(t, "arithmetic-mean-v1", aggregation["version"])
						require.Equal(t, "exact", aggregation["basis"])
						require.Equal(t, "uint64", aggregation["source_type"])
						require.Equal(t, "decimal", aggregation["result_type"], "the result type is not the Tag's type")
						require.EqualValues(t, 2, aggregation["valid_count"])
						require.EqualValues(t, 0, aggregation["rejected_count"])
						require.Equal(t, "unknown", aggregation["coverage"], "counts do not prove completeness")
					}
				}
				for _, leaked := range []string{"9007199254740992", "18014398509481985", "15.5", `"sum"`, "scale", "rounding"} {
					require.NotContains(t, raw, leaked, "evidence never carries a value, a sum or its encoding")
				}
			})

			t.Run("a receipt that does not match the frozen payload confirms nothing", func(t *testing.T) {
				f := closedAndDelivered(t, shape)
				_, err := f.plainLocal().ExecContext(t.Context(), `UPDATE wg_delivery_receipts SET payload_digest = ?`, strings.Repeat("f", 64))
				require.NoError(t, err)

				document, _ := f.sourceEvidenceJSON()
				for _, source := range evidenceSources(t, document) {
					require.Nil(t, recordOf(t, source, "last_sql_committed"), "no destination evidence without a matching receipt")
					latest := recordOf(t, source, "latest")
					require.Equal(t, "unknown", latest["delivery"])
					require.NotContains(t, latest, "committed_at")
				}
			})

			t.Run("an applied proof that no longer holds leaves the evidence unavailable", func(t *testing.T) {
				f := closedAndDelivered(t, shape)
				// Corrupt the real frozen runtime version Apply wrote; a fake in-memory
				// snapshot would not be the authority the production pipeline reads.
				var payload string
				require.NoError(t, f.plainLocal().QueryRowContext(t.Context(),
					`SELECT payload FROM wg_runtime_versions WHERE workspace_id=? AND group_id=? AND group_revision=?`,
					f.group.WorkspaceID, f.group.ID, f.group.AppliedRevision).Scan(&payload))
				var frozen workspace.WriteGroupAppliedSnapshot
				require.NoError(t, json.Unmarshal([]byte(payload), &frozen))
				frozen.RuntimeLayout = nil
				corrupted, err := json.Marshal(frozen)
				require.NoError(t, err)
				_, err = f.plainLocal().ExecContext(t.Context(),
					`UPDATE wg_runtime_versions SET payload=? WHERE workspace_id=? AND group_id=? AND group_revision=?`,
					string(corrupted), f.group.WorkspaceID, f.group.ID, f.group.AppliedRevision)
				require.NoError(t, err)

				document, _ := f.sourceEvidenceJSON()
				require.NotNil(t, document)
				require.Equal(t, "unavailable", document["state"], "no frozen layout, no per source claim")
				require.Equal(t, "applied_snapshot_unavailable", document["unavailable_reason"])
				require.Empty(t, document["sources"])
			})

			t.Run("a destination other than the applied one proves nothing", func(t *testing.T) {
				f := closedAndDelivered(t, shape)
				// Keep the frozen proof valid and move the real delivered record to
				// another destination revision; the current draft is not the applied one.
				_, err := f.plainLocal().ExecContext(t.Context(),
					`UPDATE wg_delivery_outbox SET connector_revision='crev-other'`)
				require.NoError(t, err)

				document, raw := f.sourceEvidenceJSON()
				for _, source := range evidenceSources(t, document) {
					latest := recordOf(t, source, "latest")
					require.Equal(t, false, latest["available"])
					require.Equal(t, "identity_mismatch", latest["cause"])
					require.NotContains(t, latest, "aggregation")
				}
				require.NotContains(t, raw, "arithmetic-mean")
			})

			for _, corruption := range []struct{ name, old, replacement string }{
				{"a coverage that is not one of the four", `\"coverage\":\"unknown\"`, `\"coverage\":\"bogus\"`},
				{"a result type that contradicts an exact basis", `\"result_type\":\"decimal\"`, `\"result_type\":\"float64\"`},
				{"a negative valid count", `\"valid_count\":2`, `\"valid_count\":-2`},
			} {
				t.Run("unavailable rather than guessed for "+corruption.name, func(t *testing.T) {
					f := closedAndDelivered(t, shape)
					f.tamperPayloads(corruption.old, corruption.replacement)

					document, raw := f.sourceEvidenceJSON()
					for _, source := range evidenceSources(t, document) {
						latest := recordOf(t, source, "latest")
						require.Equal(t, false, latest["available"], "invalid metadata is never shown as evidence")
						require.Equal(t, "payload_invalid", latest["cause"])
						require.NotContains(t, latest, "aggregation")
					}
					require.NotContains(t, raw, "complete")
				})
			}

			t.Run("a source type that differs from the frozen applied layout is not evidence, even when it is self consistent", func(t *testing.T) {
				f := closedAndDelivered(t, shape)
				// uint64 -> decimal is what the applied layout froze. int64 -> decimal is a legal mean of another
				// source, so it only fails against the layout, never against the record's own fields.
				f.tamperPayloads(`\"source_type\":\"uint64\"`, `\"source_type\":\"int64\"`)

				document, raw := f.sourceEvidenceJSON()
				for _, source := range evidenceSources(t, document) {
					latest := recordOf(t, source, "latest")
					require.Equal(t, false, latest["available"], "metadata the applied layout does not carry is not shown")
					require.Contains(t, []any{"identity_mismatch", "payload_invalid"}, latest["cause"])
					require.NotContains(t, latest, "aggregation")
				}
				require.NotContains(t, raw, "int64")
			})

			t.Run("complete coverage without any whole-window proof is not evidence", func(t *testing.T) {
				f := closedAndDelivered(t, shape)
				// All four coverage words are legal; none of the stored evidence proves a whole window, so a record
				// that says complete is a claim nothing can back and must not be shown as one.
				f.tamperPayloads(`\"coverage\":\"unknown\"`, `\"coverage\":\"complete\"`)

				document, raw := f.sourceEvidenceJSON()
				for _, source := range evidenceSources(t, document) {
					latest := recordOf(t, source, "latest")
					require.Equal(t, false, latest["available"])
					require.Equal(t, "payload_invalid", latest["cause"])
					require.NotContains(t, latest, "aggregation")
				}
				require.NotContains(t, raw, "complete")
			})
		})
	}

	t.Run("a cumulative wide group reports a baseline as a first reading with an unknown difference", func(t *testing.T) {
		f := closedAndDelivered(t, shapeCumulativeWide)
		document, raw := f.sourceEvidenceJSON()
		for _, source := range evidenceSources(t, document) {
			latest := recordOf(t, source, "latest")
			require.Equal(t, true, latest["available"])
			require.Equal(t, "baseline", latest["status"], "a baseline is a status, not a failure")
			require.Equal(t, false, latest["delta_known"], "an unknown difference is false, never zero")
			require.Equal(t, "sql_committed", latest["delivery"])
			require.NotContains(t, latest, "aggregation", "a counter has no mean")
		}
		for _, leaked := range []string{"9007199254740992", `"delta"`} {
			require.NotContains(t, raw, leaked)
		}
	})
}

// TestAverageDeliveryEvidenceCounters pins what a counted wide row may report about a source: the types the applied
// layout froze for it, and a difference that is a number with both of its bounds. The counted wide group is closed for
// two windows so the newest row carries a real difference for each source.
func TestAverageDeliveryEvidenceCounters(t *testing.T) {
	closedTwice := func(t *testing.T) *longPipeline {
		t.Helper()
		f := newAveragePipeline(t, shapeCumulativeWide)
		for window := range 2 {
			f.feed(shapeCumulativeWide, window)
			f.closeUntil(shapeCumulativeWide.windowEnd(window))
		}
		f.deliver()
		return f
	}

	t.Run("a whole counted row reports its differences as known", func(t *testing.T) {
		f := closedTwice(t)
		document, raw := f.sourceEvidenceJSON()
		for _, source := range evidenceSources(t, document) {
			latest := recordOf(t, source, "latest")
			require.Equal(t, true, latest["available"])
			require.Equal(t, "ok", latest["status"])
			require.Equal(t, true, latest["delta_known"])
		}
		require.NotContains(t, raw, `"delta"`, "the difference itself is never returned")
	})

	for _, tc := range []struct{ name, old, replacement string }{
		{"a source type that is legal but not the one the applied layout froze", `\"source_type\":\"uint64\"`, `\"source_type\":\"int64\"`},
		{"a result type that is legal but not the one the applied layout froze", `\"result_type\":\"uint64\"`, `\"result_type\":\"int64\"`},
	} {
		t.Run("unavailable for "+tc.name, func(t *testing.T) {
			f := closedTwice(t)
			f.tamperPayloads(tc.old, tc.replacement)

			document, _ := f.sourceEvidenceJSON()
			for _, source := range evidenceSources(t, document) {
				latest := recordOf(t, source, "latest")
				require.Equal(t, false, latest["available"])
				require.Contains(t, []any{"identity_mismatch", "payload_invalid"}, latest["cause"])
			}
		})
	}

	t.Run("unavailable for a difference that is not a number, only for the source that carries it", func(t *testing.T) {
		f := closedTwice(t)
		f.tamperPayloads(`\"delta\":\"1\"`, `\"delta\":\"one\"`)

		sources := evidenceSources(t, firstOf(f.sourceEvidenceJSON()))
		require.Len(t, sources, 2)
		first, second := recordOf(t, sources[0], "latest"), recordOf(t, sources[1], "latest")
		require.Equal(t, false, first["available"])
		require.Equal(t, "payload_invalid", first["cause"])
		require.Equal(t, true, second["available"], "another source of the same row is judged by its own evidence")
		require.Equal(t, true, second["delta_known"])
	})

	t.Run("unavailable for a difference whose bounds are missing", func(t *testing.T) {
		f := closedTwice(t)
		f.tamperPayloads(`\"delta_from\":`, `\"delta_from_lost\":`)

		document, _ := f.sourceEvidenceJSON()
		for _, source := range evidenceSources(t, document) {
			latest := recordOf(t, source, "latest")
			require.Equal(t, false, latest["available"], "a difference without both bounds claims nothing")
			require.Equal(t, "payload_invalid", latest["cause"])
		}
	})
}

func firstOf(document map[string]any, _ string) map[string]any { return document }
