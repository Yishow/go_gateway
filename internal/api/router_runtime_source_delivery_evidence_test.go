package api

import (
	"net/http"
	"testing"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/runtimeview"

	"github.com/stretchr/testify/require"
)

func TestRuntimeSourceDeliveryEvidenceMeanPartialAndUnknownCoverageNeverComplete(t *testing.T) {
	for _, wide := range []bool{false, true} {
		t.Run(map[bool]string{false: "long", true: "wide"}[wide], func(t *testing.T) {
			r := newRuntimeSourceDeliveryRig(t, wide)
			r.offer(t, "mean", 0, false)
			r.offer(t, "mean", 0, true)
			r.offer(t, "other", 0, false)
			r.closeWindow(t, 0)
			raw := r.raw(t)
			require.NotNil(t, raw.SourceEvidence)
			require.Equal(t, "available", raw.SourceEvidence.State)
			require.Len(t, raw.SourceEvidence.Sources, 3)
			var partial bool
			for _, s := range raw.SourceEvidence.Sources {
				if s.TagID == "tag-mean" {
					require.True(t, s.Latest.Available)
					require.Equal(t, "partial", s.Latest.Aggregation.Coverage)
					partial = true
				}
			}
			require.True(t, partial, "production authority must really contain mean partial evidence")
			g := r.evidence(t)
			require.Equal(t, "available", g.Delivery.Current.SourceState, "scoped GET must preserve verified source evidence, not hardcode unavailable")
			require.Len(t, g.Delivery.Current.Sources, 3)
			for _, tc := range []struct{ tag, state, coverage string }{{"tag-mean", "partial", "partial"}, {"tag-other", "unknown", "unknown"}, {"tag-quiet", "missing", ""}} {
				s := runtimeDeliveredSource(t, g, tc.tag)
				require.NotNil(t, s.Latest)
				require.Equal(t, tc.state, s.Latest.State)
				require.Equal(t, r.snapshot.EffectiveAt.Format(time.RFC3339Nano), *s.Latest.WindowStart)
				require.Equal(t, r.snapshot.EffectiveAt.Add(r.window).Format(time.RFC3339Nano), *s.Latest.WindowEnd)
				if tc.coverage == "" {
					require.Nil(t, s.Latest.Coverage)
				} else {
					require.Equal(t, tc.coverage, *s.Latest.Coverage)
				}
				require.Nil(t, s.Latest.Receipt)
				require.Nil(t, s.LastCommitted)
			}
		})
	}
}

func TestRuntimeSourceDeliveryEvidenceDeltaUnknownNeverZero(t *testing.T) {
	for _, wide := range []bool{false, true} {
		t.Run(map[bool]string{false: "long", true: "wide"}[wide], func(t *testing.T) {
			r := newRuntimeSourceDeliveryOptions(t, runtimeSourceDeliveryOptions{wide: wide, cumulative: true})
			r.offer(t, "mean", 0, false)
			r.closeWindow(t, 0)
			r.deliver(t)
			raw := r.raw(t)
			require.Equal(t, "available", raw.SourceEvidence.State)
			s := runtimeDeliveredSource(t, r.evidence(t), "tag-mean")
			require.NotNil(t, s.Latest)
			require.Equal(t, "unknown", s.Latest.State)
			require.NotNil(t, s.Latest.DeltaKnown)
			require.False(t, *s.Latest.DeltaKnown)
			require.NotNil(t, s.LastCommitted)
			require.Equal(t, "unknown", s.LastCommitted.State)
			require.False(t, *s.LastCommitted.DeltaKnown)
			require.NotNil(t, s.LastCommitted.Receipt)
			require.Equal(t, "runtime", s.LastCommitted.Receipt.Provenance)
			require.Nil(t, s.Latest.Coverage)
		})
	}
}

func TestRuntimeSourceDeliveryEvidenceLatestQueuedSeparateFromEarlierCommittedWindow(t *testing.T) {
	r := newRuntimeSourceDeliveryRig(t, false)
	r.offer(t, "mean", 0, false)
	r.offer(t, "other", 0, false)
	r.closeWindow(t, 0)
	r.deliver(t)
	first := r.evidence(t)
	source := runtimeDeliveredSource(t, first, "tag-mean")
	require.NotNil(t, source.LastCommitted)
	require.NotNil(t, source.LastCommitted.Receipt)
	earlier := source.LastCommitted
	// Different source rows legitimately have distinct receipts; none is forced
	// to equal the latest group-wide receipt.
	ids := map[string]bool{}
	for _, s := range first.Delivery.Current.Sources {
		require.NotNil(t, s.LastCommitted.Receipt)
		ids[s.LastCommitted.Receipt.RecordID] = true
	}
	require.Len(t, ids, 3)
	r.offer(t, "mean", 1, false)
	r.closeWindow(t, 1)
	g := r.evidence(t)
	s := runtimeDeliveredSource(t, g, "tag-mean")
	require.EqualValues(t, 3, g.Delivery.Current.Stages.Queued)
	require.EqualValues(t, 3, g.Delivery.Current.Stages.SQLCommitted)
	require.Equal(t, earlier, s.LastCommitted)
	require.Equal(t, *earlier.WindowEnd, *s.Latest.WindowStart)
	require.Nil(t, s.Latest.Receipt)
	require.NotEqual(t, *s.Latest.WindowStart, *s.LastCommitted.WindowStart)
	require.Nil(t, g.Delivery.TestWrite, "runtime receipt is not test-write/readback/cleanup proof")
}

func TestRuntimeSourceDeliveryEvidenceSourceReceiptMismatchNeverBorrowsGroupCommit(t *testing.T) {
	for _, damage := range []string{"missing", "digest", "time"} {
		t.Run(damage, func(t *testing.T) {
			r := newRuntimeSourceDeliveryRig(t, false)
			r.offer(t, "mean", 0, false)
			r.offer(t, "other", 0, false)
			r.closeWindow(t, 0)
			r.deliver(t)
			prior := runtimeDeliveredSource(t, r.evidence(t), "tag-mean")
			effect := prior.LastCommitted.Receipt.EffectKey
			var err error
			switch damage {
			case "missing":
				_, err = r.f.db.Exec(`DELETE FROM wg_delivery_receipts WHERE effect_key=?`, effect)
			case "digest":
				_, err = r.f.db.Exec(`UPDATE wg_delivery_receipts SET payload_digest='mismatch' WHERE effect_key=?`, effect)
			case "time":
				_, err = r.f.db.Exec(`UPDATE wg_delivery_receipts SET committed_at='2026-10-01T00:00:00Z' WHERE effect_key=?`, effect)
			}
			require.NoError(t, err)
			raw := r.raw(t)
			require.NotNil(t, raw.LastSQLCommittedEffect, "other actual source receipts remain group-level proof")
			g := r.evidence(t)
			s := runtimeDeliveredSource(t, g, "tag-mean")
			require.NotNil(t, g.Delivery.Current.LastCommit)
			require.Nil(t, s.LastCommitted)
			require.Nil(t, s.Latest.Receipt)
			require.NotNil(t, runtimeDeliveredSource(t, g, "tag-other").LastCommitted.Receipt)
		})
	}
}

func TestRuntimeSourceDeliveryEvidenceUnverifiedRecordAndMissingRecordHaveNoInventedProof(t *testing.T) {
	t.Run("no record", func(t *testing.T) {
		r := newRuntimeSourceDeliveryRig(t, false)
		g := r.evidence(t)
		for _, s := range g.Delivery.Current.Sources {
			require.Nil(t, s.Latest)
			require.Nil(t, s.LastCommitted)
		}
		require.Equal(t, "available", g.Delivery.Current.SourceState)
		require.Len(t, g.Delivery.Current.Sources, 3)
	})
	t.Run("invalid durable payload", func(t *testing.T) {
		r := newRuntimeSourceDeliveryRig(t, false)
		r.offer(t, "mean", 0, false)
		r.closeWindow(t, 0)
		_, err := r.f.db.Exec(`UPDATE wg_delivery_outbox SET payload='private-invalid-payload' WHERE entity_key=?`, `["device-A","point-mean","tag-mean"]`)
		require.NoError(t, err)
		g := r.evidence(t)
		s := runtimeDeliveredSource(t, g, "tag-mean")
		require.Equal(t, &runtimeview.SourceResult{State: "unavailable"}, s.Latest)
		require.Nil(t, s.LastCommitted)
	})
}

func TestRuntimeSourceDeliveryEvidenceSameGroupOtherDeviceNeverGetsSelectedSourceSuccess(t *testing.T) {
	r := newRuntimeSourceDeliveryOptions(t, runtimeSourceDeliveryOptions{multi: true})
	// Only B has an actual reading; A's missing rows may be committed, but that
	// proves recording of missing data, not a successful sample or mean for A.
	r.offer(t, "other", 0, false)
	r.closeWindow(t, 0)
	r.deliver(t)
	g := r.evidence(t)
	require.True(t, g.Delivery.Current.Shared)
	a := runtimeDeliveredSource(t, g, "tag-mean")
	b := runtimeDeliveredSource(t, g, "tag-other")
	require.Equal(t, "device-A", a.DeviceID)
	require.Equal(t, "missing", a.Latest.State)
	require.Nil(t, a.Latest.Coverage)
	require.Equal(t, "device-B", b.DeviceID)
	require.Equal(t, "unknown", b.Latest.State)
	require.Equal(t, "unknown", *b.Latest.Coverage)
	require.NotEqual(t, a.LastCommitted.Receipt.RecordID, b.LastCommitted.Receipt.RecordID)
	q := r.query(t)
	q.Set("device_id", "device-B")
	response := performJSONRequest(t, r.router, http.MethodGet, runtimeRuleContextPath+"?"+q.Encode(), nil)
	require.Equal(t, 200, response.Code, response.Body.String())
	data := decodeJSONBody(t, response)["data"].(map[string]any)["rule_evidence"].(map[string]any)
	require.Equal(t, "device-B", data["device_id"])
}

func TestRuntimeSourceDeliveryEvidenceUsesHandlerEffectiveAWhenPipelineClockHasB(t *testing.T) {
	r := newRuntimeSourceDeliveryRig(t, false)
	r.offer(t, "mean", 0, false)
	a := r.snapshot.AppliedRevision
	r.clock = r.snapshot.EffectiveAt.Add(2 * time.Second)
	r.update(t, r.snapshot.Group.Destination.TableName, 15)
	b := r.apply(t).Group.AppliedRevision
	asOf := r.snapshot.EffectiveAt.Add(14 * time.Second)
	r.clock = r.snapshot.EffectiveAt.Add(15 * time.Second)
	raw, err := r.pipe.DeliveryForRevision(t.Context(), r.snapshot.GroupID, a)
	require.NoError(t, err)
	require.Equal(t, grouppipeline.SourceEvidenceUnavailable, raw.SourceEvidence.State)
	require.Equal(t, grouppipeline.SourceEvidenceRevisionMismatch, raw.SourceEvidence.UnavailableReason)
	r.useAsOf(asOf)
	g := r.evidence(t)
	require.Equal(t, a, *g.EffectiveRevision)
	require.Equal(t, b, *g.AppliedRevision)
	require.Equal(t, a, g.Delivery.Current.Binding.EntityRevision)
	require.EqualValues(t, 1, g.Delivery.Current.Stages.Collecting)
	require.Equal(t, "unavailable", g.Delivery.Current.SourceState)
	require.Empty(t, g.Delivery.Current.Sources)
	r.useAsOf(r.clock)
	g = r.evidence(t)
	require.Equal(t, b, *g.EffectiveRevision)
	require.Zero(t, g.Delivery.Current.Stages.Collecting)
	require.Equal(t, "available", g.Delivery.Current.SourceState)
	for _, s := range g.Delivery.Current.Sources {
		require.Nil(t, s.Latest)
		require.Nil(t, s.LastCommitted)
	}
}

func TestRuntimeSourceDeliveryEvidenceLegacyCustomMissingFrozenProofUnavailable(t *testing.T) {
	t.Run("legacy wide", func(t *testing.T) {
		r := newRuntimeDeliveryRig(t)
		r.accept(t, "legacy")
		r.tick(t)
		raw, err := r.pipe.DeliveryForRevision(t.Context(), r.f.group["id"].(string), r.f.group["applied_revision"].(string))
		require.NoError(t, err)
		require.Nil(t, raw.SourceEvidence)
		g := r.evidence(t)
		require.Equal(t, "unavailable", g.Delivery.Current.SourceState)
		require.Empty(t, g.Delivery.Current.Sources)
	})
	t.Run("custom", func(t *testing.T) {
		r := newRuntimeSourceCustomRig(t)
		r.accept(t, "custom")
		r.tick(t)
		raw, err := r.pipe.DeliveryForRevision(t.Context(), r.f.group["id"].(string), r.f.group["applied_revision"].(string))
		require.NoError(t, err)
		require.Nil(t, raw.SourceEvidence)
		g := r.evidence(t)
		require.Nil(t, g.Effective)
		require.Nil(t, g.Delivery.Current)
	})
	t.Run("missing frozen layout", func(t *testing.T) {
		r := newRuntimeSourceDeliveryRig(t, true)
		r.offer(t, "mean", 0, false)
		r.closeWindow(t, 0)
		_, err := r.f.db.Exec(`DELETE FROM wg_runtime_versions`)
		require.NoError(t, err)
		g := r.evidence(t)
		require.Nil(t, g.Effective)
		require.Nil(t, g.Delivery.Current)
	})
}

func TestRuntimeSourceDeliveryEvidenceForeignDuplicateTupleAndUnsafeProofFailClosed(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		change func(*grouppipeline.SourceEvidenceView)
	}{
		{"foreign tuple", 403, func(v *grouppipeline.SourceEvidenceView) { v.Sources[0].DeviceID = "foreign-device" }},
		{"duplicate tuple", 503, func(v *grouppipeline.SourceEvidenceView) { v.Sources = append(v.Sources, v.Sources[0]) }},
		{"mixed revision", 409, func(v *grouppipeline.SourceEvidenceView) { v.GroupRevision = "different-revision" }},
		{"unsafe proof", 503, func(v *grouppipeline.SourceEvidenceView) {
			v.Sources[0].Latest.RecordID = "password=private-proof-secret"
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newRuntimeSourceDeliveryRig(t, false)
			r.offer(t, "mean", 0, false)
			r.closeWindow(t, 0)
			r.useReader(&runtimeDeliveryFault{Pipeline: r.pipe, change: func(v *grouppipeline.DeliveryView) { tc.change(v.SourceEvidence) }})
			before := r.durableState(t)
			response := r.get(t)
			require.Equal(t, tc.status, response.Code, response.Body.String())
			require.NotContains(t, response.Body.String(), "private-proof-secret")
			require.NotContains(t, response.Body.String(), `"rule_evidence"`)
			require.Equal(t, before, r.durableState(t))
		})
	}
}

func TestRuntimeSourceDeliveryEvidenceGETPreservesSQLJournalDescriptorsAndRedaction(t *testing.T) {
	r := newRuntimeSourceDeliveryRig(t, true)
	r.offer(t, "mean", 0, false)
	r.offer(t, "mean", 0, true)
	r.closeWindow(t, 0)
	r.deliver(t)
	before := r.durableState(t)
	var changes int
	require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&changes))
	for range 2 {
		response := r.get(t)
		require.Equal(t, 200, response.Code, response.Body.String())
		for _, secret := range []string{"private-driver-secret", "/private/raw-value", r.f.file, `"raw_present"`, `"value"`, `"sum"`, `"provenance_json"`, `"connection_config"`, `"payload"`} {
			require.NotContains(t, response.Body.String(), secret)
		}
	}
	var after int
	require.NoError(t, r.f.db.QueryRow(`SELECT total_changes()`).Scan(&after))
	require.Equal(t, changes, after)
	require.Equal(t, before, r.durableState(t))
	g := r.evidence(t)
	require.Nil(t, g.Delivery.TestWrite)
	require.NotNil(t, g.Delivery.Current.LastCommit)
	require.Equal(t, "runtime", g.Delivery.Current.LastCommit.Provenance)
	require.Nil(t, g.Delivery.Current.NextDueAt)
	require.Nil(t, g.Delivery.Current.NextRetryAt)
}

func TestRuntimeSourceDeliveryEvidenceDamagedLatestReceiptRetainsOnlyOwnEarlierReceipt(t *testing.T) {
	r := newRuntimeSourceDeliveryRig(t, false)
	r.offer(t, "mean", 0, false)
	r.offer(t, "other", 0, false)
	r.closeWindow(t, 0)
	r.deliver(t)
	earlier := runtimeDeliveredSource(t, r.evidence(t), "tag-mean").LastCommitted
	require.NotNil(t, earlier)
	require.NotNil(t, earlier.Receipt)
	r.offer(t, "mean", 1, false)
	r.offer(t, "other", 1, false)
	r.closeWindow(t, 1)
	r.deliver(t)
	newest := runtimeDeliveredSource(t, r.evidence(t), "tag-mean").LastCommitted
	require.NotNil(t, newest)
	require.NotEqual(t, earlier.Receipt.EffectKey, newest.Receipt.EffectKey)
	_, err := r.f.db.Exec(`UPDATE wg_delivery_receipts SET payload_digest='mismatch' WHERE effect_key=?`, newest.Receipt.EffectKey)
	require.NoError(t, err)
	before := r.durableState(t)
	g := r.evidence(t)
	s := runtimeDeliveredSource(t, g, "tag-mean")
	require.Nil(t, s.Latest.Receipt)
	require.Equal(t, *newest.WindowStart, *s.Latest.WindowStart)
	require.Equal(t, earlier, s.LastCommitted, "only this source's own prior matching receipt survives")
	require.NotEqual(t, g.Delivery.Current.LastCommit.RecordID, s.LastCommitted.Receipt.RecordID, "another source's newer group receipt must not be allocated to this source")
	require.Equal(t, *newest.WindowStart, *runtimeDeliveredSource(t, g, "tag-other").LastCommitted.WindowStart)
	require.Equal(t, before, r.durableState(t))
}
