package groupdelivery

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"go-gateway/internal/datalink/counterproof"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

func replaceLegacySample(t *testing.T, store *Store, change func(*snapshot.Sample)) {
	t.Helper()
	var payload, digest string
	require.NoError(t, store.DB().QueryRow(`SELECT payload,payload_digest FROM wg_delivery_samples WHERE sample_id='c-2'`).Scan(&payload, &digest))
	sample, err := decodeVerifiedSample(payload, digest)
	require.NoError(t, err)
	change(&sample)
	wire, nextDigest, err := encodeSample(sample)
	require.NoError(t, err)
	_, err = store.DB().Exec(`UPDATE wg_delivery_samples SET payload=?,payload_digest=?,member_key=? WHERE sample_id='c-2'`, string(wire), nextDigest, sample.MemberKey)
	require.NoError(t, err)
}

func TestLegacyCounterProofAssociationChecksEveryOriginalSample(t *testing.T) {
	for name, mutate := range map[string]func(*testing.T, *Store){
		"source gap": func(t *testing.T, s *Store) {
			replaceLegacySample(t, s, func(v *snapshot.Sample) { v.SourceRevision = "" })
		},
		"mapping mismatch": func(t *testing.T, s *Store) {
			replaceLegacySample(t, s, func(v *snapshot.Sample) { v.MappingRevision = "map-2" })
		},
		"result type mismatch": func(t *testing.T, s *Store) {
			replaceLegacySample(t, s, func(v *snapshot.Sample) { v.Value = measurement.NewInt64(120004) })
		},
		"unknown member": func(t *testing.T, s *Store) {
			replaceLegacySample(t, s, func(v *snapshot.Sample) { v.MemberKey = `["device-1","point-2","tag-2"]` })
		},
		"corrupt digest": func(t *testing.T, s *Store) {
			_, err := s.DB().Exec(`UPDATE wg_delivery_samples SET payload_digest='wrong' WHERE sample_id='c-2'`)
			require.NoError(t, err)
		},
		"foreign workspace": func(t *testing.T, s *Store) {
			_, err := s.DB().Exec(`UPDATE wg_delivery_samples SET workspace_id='other' WHERE sample_id='c-2'`)
			require.NoError(t, err)
		},
		"row member disagrees": func(t *testing.T, s *Store) {
			_, err := s.DB().Exec(`UPDATE wg_delivery_samples SET member_key='other' WHERE sample_id='c-2'`)
			require.NoError(t, err)
		},
		"row sample identity disagrees": func(t *testing.T, s *Store) {
			_, err := s.DB().Exec(`UPDATE wg_delivery_samples SET sample_id='other' WHERE sample_id='c-2'`)
			require.NoError(t, err)
		},
		"accepted timestamp corrupt": func(t *testing.T, s *Store) {
			_, err := s.DB().Exec(`UPDATE wg_delivery_samples SET created_at='not-a-time' WHERE sample_id='c-2'`)
			require.NoError(t, err)
		},
	} {
		t.Run(name, func(t *testing.T) {
			store, _, descriptor, refs := legacyAssociationFixture(t)
			mutate(t, store)
			before := associationJournal(t, store)
			err := store.AssociateLegacyCounterProofs(t.Context(), testKey, refs)
			require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
			require.Zero(t, countRows(t, store.DB(), "wg_counter_proof_associations"))
			require.Equal(t, before, associationJournal(t, store))
			raw, err := store.RuntimeVersion(t.Context(), testKey)
			require.NoError(t, err)
			require.Equal(t, descriptor, raw)
		})
	}
}

func TestLegacyCounterProofAssociationRequiresActualImmutableEvidence(t *testing.T) {
	for _, change := range []string{"missing", "incomplete exact", "changed actual evidence", "caller reference", "unsupported proof"} {
		t.Run(change, func(t *testing.T) {
			store, _, descriptor, refs := legacyAssociationFixture(t)
			raw := string(descriptor)
			switch change {
			case "missing":
				_, err := store.DB().Exec(`DELETE FROM wg_runtime_versions`)
				require.NoError(t, err)
			case "incomplete exact":
				raw = strings.Replace(raw, `"source_binding_revision":"binding-1"`, `"source_binding_revision":""`, 1)
			case "changed actual evidence":
				raw = strings.Replace(raw, "kWh", "MWh", 1)
			case "caller reference":
				refs[0].Digest = strings.Repeat("a", 64)
			case "unsupported proof":
				refs[0].Version = 77
			}
			if change != "missing" {
				_, err := store.DB().Exec(`UPDATE wg_runtime_versions SET payload=?`, raw)
				require.NoError(t, err)
			}
			before := associationJournal(t, store)
			err := store.AssociateLegacyCounterProofs(t.Context(), testKey, refs)
			var typed *counterproof.CounterBindingProofError
			require.ErrorAs(t, err, &typed)
			require.Zero(t, countRows(t, store.DB(), "wg_counter_proof_associations"))
			require.Equal(t, before, associationJournal(t, store))
		})
	}
}

func TestLegacyCounterProofAssociationPublicationIsAtomic(t *testing.T) {
	for _, failure := range []string{"second audit insert", "cancelled", "conflicting prior association", "unknown association version"} {
		t.Run(failure, func(t *testing.T) {
			store, _, _, refs := legacyAssociationFixture(t)
			before := associationJournal(t, store)
			ctx := t.Context()
			if failure == "conflicting prior association" || failure == "unknown association version" {
				require.NoError(t, store.AssociateLegacyCounterProofs(ctx, testKey, refs))
				_, err := store.DB().Exec(`DELETE FROM wg_counter_proof_associations WHERE sample_id='c-1'`)
				require.NoError(t, err)
				column := "proof_digest='conflicting'"
				if failure == "unknown association version" {
					column = "association_version=77"
				}
				_, err = store.DB().Exec(`UPDATE wg_counter_proof_associations SET ` + column + ` WHERE sample_id='c-2'`)
				require.NoError(t, err)
			}
			if failure == "second audit insert" {
				_, err := store.DB().Exec(`CREATE TRIGGER refuse_second_association BEFORE INSERT ON wg_counter_proof_associations WHEN NEW.sample_id='c-2' BEGIN SELECT RAISE(ABORT,'injected association failure'); END`)
				require.NoError(t, err)
			}
			if failure == "cancelled" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}
			audit := associationAudit(t, store)
			err := store.AssociateLegacyCounterProofs(ctx, testKey, refs)
			require.Error(t, err)
			if failure == "conflicting prior association" {
				require.ErrorIs(t, err, counterproof.ErrCounterBindingProofConflict)
			}
			if failure == "unknown association version" {
				require.ErrorIs(t, err, counterproof.ErrCounterBindingProofUnsupported)
			}
			require.Equal(t, audit, associationAudit(t, store), "a first valid insert cannot survive a later failure")
			require.Equal(t, before, associationJournal(t, store))
		})
	}
}

func TestLegacyCounterProofAssociationConcurrentRetryIsOneAudit(t *testing.T) {
	store, _, _, refs := legacyAssociationFixture(t)
	var wg sync.WaitGroup
	errors := make(chan error, 6)
	for range 6 {
		wg.Go(func() { errors <- store.AssociateLegacyCounterProofs(t.Context(), testKey, refs) })
	}
	wg.Wait()
	close(errors)
	for err := range errors {
		require.NoError(t, err)
	}
	require.Equal(t, 2, countRows(t, store.DB(), "wg_counter_proof_associations"))
	require.Equal(t, 2, countRows(t, store.DB(), "wg_delivery_samples"))
}

func TestLegacyCounterProofAssociationSkipsNewAndFinalizedInputs(t *testing.T) {
	t.Run("new versioned input", func(t *testing.T) {
		store, db, _ := newTestStore(t)
		proof, sample, payload := journalProofFixture(t)
		require.NoError(t, store.SaveRuntimeVersion(t.Context(), testKey, payload))
		require.NoError(t, store.AppendSample(t.Context(), testKey, t0, sample))
		ref, err := proof.Reference()
		require.NoError(t, err)
		require.NoError(t, store.AssociateLegacyCounterProofs(t.Context(), testKey, []counterproof.CounterBindingReference{ref}))
		require.Zero(t, countRows(t, db, "wg_counter_proof_associations"))
	})
	t.Run("finalized input", func(t *testing.T) {
		store, _, _, refs := legacyAssociationFixture(t)
		_, err := store.DB().Exec(`UPDATE wg_delivery_samples SET consumed=1`)
		require.NoError(t, err)
		before := associationJournal(t, store)
		require.NoError(t, store.AssociateLegacyCounterProofs(t.Context(), testKey, refs))
		require.Zero(t, countRows(t, store.DB(), "wg_counter_proof_associations"))
		require.Equal(t, before, associationJournal(t, store))
	})
	t.Run("ordinary noncounter is not assigned counter meaning", func(t *testing.T) {
		store, _, _, refs := legacyAssociationFixture(t)
		var raw string
		require.NoError(t, store.DB().QueryRow(`SELECT payload FROM wg_runtime_versions`).Scan(&raw))
		var descriptor map[string]any
		require.NoError(t, json.Unmarshal([]byte(raw), &descriptor))
		group := descriptor["group"].(map[string]any)
		group["recording_profile"].(map[string]any)["kind"] = "instantaneous"
		members := group["members"].([]any)
		members[0].(map[string]any)["value_role"] = "counter_raw"
		group["members"] = append(members, map[string]any{"device_id": "device-2", "point_id": "point-2", "tag_id": "tag-2", "value_role": "gauge"})
		payload, err := json.Marshal(descriptor)
		require.NoError(t, err)
		_, err = store.DB().Exec(`UPDATE wg_runtime_versions SET payload=?`, string(payload))
		require.NoError(t, err)
		ordinary := sampleAt("gauge-1", `["device-2","point-2","tag-2"]`, 5, 30)
		require.NoError(t, store.AppendSample(t.Context(), testKey, t0, ordinary))
		require.NoError(t, store.AssociateLegacyCounterProofs(t.Context(), testKey, refs))
		require.Equal(t, 2, countRows(t, store.DB(), "wg_counter_proof_associations"))
		require.Equal(t, 3, countRows(t, store.DB(), "wg_delivery_samples"))
	})
}
