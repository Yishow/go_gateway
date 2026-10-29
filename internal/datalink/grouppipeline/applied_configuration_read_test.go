package grouppipeline

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/workspace"
)

func TestAppliedConfigurationUsesExplicitAsOfBeforeAndAtBoundary(t *testing.T) {
	f := newAppliedConfigurationFixture(t)
	a := f.apply(t)
	f.now = at(16)
	b := f.update(t, 30, "requested B")
	b = f.apply(t)
	c := f.update(t, 45, "saved C")
	require.NotEqual(t, a.Group.Revision, b.Group.Revision)
	require.NotEqual(t, b.Group.Revision, c.Group.Revision)
	require.Equal(t, b.Group.Revision, c.Group.AppliedRevision)
	// The adapter must not call this clock: the caller fixes one observation.
	f.pipe.config.Now = func() time.Time { t.Fatal("reader consulted another clock"); return at(99) }
	before := appliedConfigurationState(t, f.db)
	for _, tc := range []struct {
		at       time.Time
		revision string
		interval int
	}{
		{at(29), a.Group.Revision, 15}, {at(30), b.Group.Revision, 30},
	} {
		got, err := f.pipe.ResolveAppliedConfigurationAt(t.Context(), a.Group.ID, tc.at)
		require.NoError(t, err)
		require.Equal(t, tc.revision, got.AppliedRevision)
		require.Equal(t, tc.interval, got.Group.RowPolicy.IntervalSeconds)
		require.NotNil(t, got.RuntimeLayout, "effective facts require the verified Apply-frozen layout")
		require.NotEmpty(t, got.RuntimeLayout.SchemaDigest)
	}
	require.Equal(t, before, appliedConfigurationState(t, f.db), "reading must not change any authority or delivery evidence")
}

func TestAppliedConfigurationRejectsMissingOrMismatchedFrozenProof(t *testing.T) {
	for _, mode := range []string{"missing", "layout", "digest", "member", "revision", "future", "malformed"} {
		t.Run(mode, func(t *testing.T) {
			f := newAppliedConfigurationFixture(t)
			saved := f.apply(t)
			key := saved.Group.ID
			var raw string
			require.NoError(t, f.db.QueryRow(`SELECT payload FROM wg_runtime_versions WHERE group_id=?`, key).Scan(&raw))
			var snap workspace.WriteGroupAppliedSnapshot
			require.NoError(t, json.Unmarshal([]byte(raw), &snap))
			expected := ErrAppliedConfigurationUnavailable
			switch mode {
			case "missing":
				_, err := f.db.Exec(`DELETE FROM wg_runtime_versions WHERE group_id=?`, key)
				require.NoError(t, err)
			case "layout":
				snap.RuntimeLayout = nil
			case "digest":
				snap.RuntimeLayout.SchemaDigest = ""
			case "member":
				snap.Group.Name = "different immutable name"
				expected = ErrAppliedConfigurationStale
			case "revision":
				snap.GroupRevision = "different-revision"
			case "future":
				snap.EffectiveAt = at(90)
				expected = ErrAppliedConfigurationStale
			}
			if mode != "missing" {
				payload, err := json.Marshal(snap)
				require.NoError(t, err)
				if mode == "malformed" {
					payload = []byte(`{"private":"credential"}`)
				}
				_, err = f.db.Exec(`UPDATE wg_runtime_versions SET payload=? WHERE group_id=?`, string(payload), key)
				require.NoError(t, err)
			}
			before := appliedConfigurationState(t, f.db)
			got, err := f.pipe.ResolveAppliedConfigurationAt(t.Context(), key, at(16))
			require.Nil(t, got)
			require.ErrorIs(t, err, expected)
			require.NotContains(t, err.Error(), "credential")
			require.Equal(t, before, appliedConfigurationState(t, f.db))
		})
	}
}

func TestAppliedConfigurationPreservesEndedSnapshotWithoutIntake(t *testing.T) {
	f := newAppliedConfigurationFixture(t)
	saved := f.apply(t)
	_, err := f.db.Exec(`INSERT INTO write_group_revision_ends(workspace_id,group_id,group_revision,ends_at,cutover_id) VALUES(?,?,?,?,?)`, saved.Group.WorkspaceID, saved.Group.ID, saved.Group.Revision, at(20).Format(time.RFC3339Nano), "ended-fixture")
	require.NoError(t, err)
	f.now = at(25)
	_, err = f.groups.Disable(t.Context(), saved.Group.ID, f.mutation(t, nil))
	require.NoError(t, err)
	before := appliedConfigurationState(t, f.db)
	got, err := f.pipe.ResolveAppliedConfigurationAt(t.Context(), saved.Group.ID, at(25))
	require.NoError(t, err)
	require.Equal(t, at(20), got.EndsAt)
	require.NotNil(t, got.RuntimeLayout)
	eligibility, err := f.groups.IntakeEligibility(t.Context(), saved.Group.ID)
	require.NoError(t, err)
	require.False(t, eligibility.Accepting)
	require.Empty(t, f.pipe.Status(), "read must never construct a boundary or start intake")
	got.Group.Members[0].TagID = "mutated"
	got.RuntimeLayout.TagTypes["tag-read"] = "string"
	again, err := f.pipe.ResolveAppliedConfigurationAt(t.Context(), saved.Group.ID, at(25))
	require.NoError(t, err)
	require.Equal(t, "tag-read", again.Group.Members[0].TagID)
	require.NotEqual(t, "string", string(again.RuntimeLayout.TagTypes["tag-read"]))
	require.Equal(t, before, appliedConfigurationState(t, f.db))
}

type appliedConfigurationFault struct {
	GroupSource
	mutate func(*workspace.WriteGroupAppliedSnapshot)
}

func (f appliedConfigurationFault) ResolveAppliedAt(ctx context.Context, id string, at time.Time) (*workspace.WriteGroupAppliedSnapshot, error) {
	s, err := f.GroupSource.ResolveAppliedAt(ctx, id, at)
	if err == nil {
		f.mutate(s)
	}
	return s, err
}

func TestAppliedConfigurationRejectsMixedResolvedScope(t *testing.T) {
	f := newAppliedConfigurationFixture(t)
	saved := f.apply(t)
	for name, mutate := range map[string]func(*workspace.WriteGroupAppliedSnapshot){
		"workspace":  func(s *workspace.WriteGroupAppliedSnapshot) { s.WorkspaceID = "foreign" },
		"group":      func(s *workspace.WriteGroupAppliedSnapshot) { s.GroupID = "foreign" },
		"revision":   func(s *workspace.WriteGroupAppliedSnapshot) { s.GroupRevision = "different" },
		"future":     func(s *workspace.WriteGroupAppliedSnapshot) { s.EffectiveAt = at(90) },
		"definition": func(s *workspace.WriteGroupAppliedSnapshot) { s.Group.Name = "raced definition" },
	} {
		t.Run(name, func(t *testing.T) {
			f.pipe.deps.Groups = appliedConfigurationFault{GroupSource: f.groups, mutate: mutate}
			before := appliedConfigurationState(t, f.db)
			got, err := f.pipe.ResolveAppliedConfigurationAt(t.Context(), saved.Group.ID, at(16))
			require.Nil(t, got)
			require.ErrorIs(t, err, ErrAppliedConfigurationStale)
			require.Equal(t, before, appliedConfigurationState(t, f.db))
		})
	}
}

func TestAppliedConfigurationRejectsAbsentScopeAndDependencies(t *testing.T) {
	f := newAppliedConfigurationFixture(t)
	saved := f.apply(t)
	for _, tc := range []struct {
		name string
		pipe *Pipeline
		id   string
		at   time.Time
	}{
		{"nil-pipeline", nil, saved.Group.ID, at(16)},
		{"missing-groups", New(Dependencies{Store: f.pipe.deps.Store}, Config{}), saved.Group.ID, at(16)},
		{"missing-store", New(Dependencies{Groups: f.groups}, Config{}), saved.Group.ID, at(16)},
		{"empty-id", f.pipe, " ", at(16)},
		{"zero-as-of", f.pipe, saved.Group.ID, time.Time{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.pipe.ResolveAppliedConfigurationAt(t.Context(), tc.id, tc.at)
			require.Nil(t, got)
			require.ErrorIs(t, err, ErrAppliedConfigurationUnavailable)
		})
	}
	got, err := f.pipe.ResolveAppliedConfigurationAt(t.Context(), "foreign-group", at(16))
	require.Nil(t, got)
	require.ErrorIs(t, err, workspace.ErrWriteGroupNotFound)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	got, err = f.pipe.ResolveAppliedConfigurationAt(ctx, saved.Group.ID, at(16))
	require.Nil(t, got)
	require.ErrorIs(t, err, context.Canceled)
	_, err = f.groups.Delete(t.Context(), saved.Group.ID, f.mutation(t, nil))
	require.NoError(t, err)
	before := appliedConfigurationState(t, f.db)
	got, err = f.pipe.ResolveAppliedConfigurationAt(t.Context(), saved.Group.ID, at(16))
	require.NoError(t, err)
	require.NotNil(t, got.RuntimeLayout)
	intake, err := f.groups.IntakeEligibility(t.Context(), saved.Group.ID)
	require.NoError(t, err)
	require.False(t, intake.Accepting)
	require.Equal(t, before, appliedConfigurationState(t, f.db))
}

func TestAppliedConfigurationRejectsLegacyCustomProofFallback(t *testing.T) {
	f := newAppliedConfigurationFixture(t)
	saved := f.apply(t)
	for _, statement := range []string{
		`UPDATE write_group_versions SET payload=json_set(payload,'$.destination.storage_strategy','custom') WHERE group_id=?`,
		`UPDATE wg_runtime_versions SET payload=json_set(payload,'$.group.destination.storage_strategy','custom') WHERE group_id=?`,
	} {
		_, err := f.db.Exec(statement, saved.Group.ID)
		require.NoError(t, err)
	}
	got, err := f.pipe.ResolveAppliedConfigurationAt(t.Context(), saved.Group.ID, at(16))
	require.Nil(t, got)
	require.ErrorIs(t, err, ErrAppliedConfigurationUnavailable)
}
