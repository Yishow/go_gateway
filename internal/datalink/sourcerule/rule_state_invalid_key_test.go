package sourcerule

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"math"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestInvalidRuleEncodingCannotAuthorizeConditionalMutation(t *testing.T) {
	for _, kind := range []string{"nan", "infinity", "invalid-time"} {
		t.Run(kind, func(t *testing.T) {
			before := &schema.SourceRule{ID: "rule", DeviceID: "device", RevisionID: "revision", NamingPrefix: "before"}
			switch kind {
			case "nan":
				v := math.NaN()
				before.ScaleMultiplier = &v
			case "infinity":
				v := math.Inf(1)
				before.ScaleOffset = &v
			default:
				before.CreatedAt = time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC)
			}
			_, marshalErr := json.Marshal(before)
			require.Error(t, marshalErr, "fixture must actually fail canonical JSON encoding")
			for _, operation := range []string{"update", "delete"} {
				t.Run(operation, func(t *testing.T) {
					repo := NewMemoryRepository()
					foreign := cloneRule(before)
					foreign.NamingPrefix = "foreign"
					require.NoError(t, repo.Create(t.Context(), foreign))
					var err error
					if operation == "delete" {
						err = repo.DeleteRevision(t.Context(), before)
					} else {
						next := cloneRule(before)
						next.NamingPrefix = "stale-write"
						err = repo.UpdateRevision(t.Context(), next, before)
					}
					require.ErrorIs(t, err, ErrRuleRevisionConflict)
					current, getErr := repo.GetByID(t.Context(), before.ID)
					require.NoError(t, getErr)
					require.Equal(t, "foreign", current.NamingPrefix)
				})
			}
			for _, operation := range []string{"invalid-next", "invalid-delete", "invalid-restore", "invalid-absent-snapshot"} {
				t.Run(operation, func(t *testing.T) {
					repo := &invalidRuleWriteSpy{MemoryRepository: NewMemoryRepository()}
					svc := NewService(repo, nil, nil, nil)
					ctx, release := svc.AcquireRuleMutation(t.Context(), before.ID)
					defer release()
					valid := &schema.SourceRule{ID: before.ID, DeviceID: before.DeviceID, RevisionID: "valid"}
					svc.noteRuleBaseline(ctx, valid)
					var err error
					switch operation {
					case "invalid-next":
						err = svc.updateOwnedRule(ctx, valid, before)
					case "invalid-delete":
						err = svc.deleteOwnedRule(ctx, before)
					case "invalid-restore":
						err = svc.restoreOwnedRule(ctx, valid, before)
					default:
						journal := svc.pointJournal(ctx)
						journal.ruleDeleted = true
						err = svc.restoreRuleSnapshot(ctx, &sourceRuleSnapshot{rule: before, pointWrites: journal})
					}
					require.ErrorIs(t, err, ErrRuleRevisionConflict)
					require.Zero(t, repo.creates)
					require.Zero(t, repo.updates)
					require.Zero(t, repo.deletes)
					require.Zero(t, repo.links)
					require.True(t, svc.ownsRuleState(ctx, valid))
				})
			}
			t.Run("preserves-valid-baseline", func(t *testing.T) {
				svc := NewService(NewMemoryRepository(), nil, nil, nil)
				ctx, release := svc.AcquireRuleMutation(t.Context(), before.ID)
				defer release()
				valid := &schema.SourceRule{ID: before.ID, DeviceID: before.DeviceID, RevisionID: "valid"}
				svc.noteRuleBaseline(ctx, valid)
				require.True(t, svc.ownsRuleState(ctx, valid))
				svc.noteRuleBaseline(ctx, before)
				require.True(t, svc.ownsRuleState(ctx, valid))
				require.False(t, svc.ownsRuleState(ctx, before))
				require.Len(t, svc.pointJournal(ctx).ruleTargets, 1)
			})
			t.Run("journal", func(t *testing.T) {
				svc := NewService(NewMemoryRepository(), nil, nil, nil)
				ctx, release := svc.AcquireRuleMutation(t.Context(), before.ID)
				defer release()
				svc.noteRuleBaseline(ctx, before)
				require.False(t, svc.ownsRuleState(ctx, before))
				foreign := cloneRule(before)
				foreign.NamingPrefix = "foreign"
				require.False(t, svc.ownsRuleState(ctx, foreign))
			})
		})
	}
}

func TestValidRuleOwnershipKeepsCanonicalJSONKey(t *testing.T) {
	scale := 1.25
	original := &schema.SourceRule{ID: "rule", DeviceID: "device", RevisionID: "revision", NamingPrefix: "漢字\"\\<>&", ScaleMultiplier: &scale, CreatedAt: time.Date(2026, time.January, 2, 3, 4, 5, 0, time.FixedZone("offset", 3600)), UpdatedAt: time.Now()}
	normalized := cloneRule(original)
	normalized.CreatedAt = normalized.CreatedAt.UTC()
	normalized.UpdatedAt = time.Time{}
	encoded, err := json.Marshal(normalized)
	require.NoError(t, err)
	expected := sha256.Sum256(encoded)
	svc := NewService(NewMemoryRepository(), nil, nil, nil)
	ctx, release := svc.AcquireRuleMutation(t.Context(), original.ID)
	defer release()
	svc.noteRuleBaseline(ctx, original)
	require.True(t, svc.pointJournal(ctx).ruleTargets[expected])
	require.True(t, svc.ownsRuleState(ctx, original))
	timestampOnly := cloneRule(original)
	timestampOnly.UpdatedAt = original.UpdatedAt.Add(time.Hour)
	require.True(t, svc.ownsRuleState(ctx, timestampOnly))
	foreign := cloneRule(original)
	foreign.NamingPrefix = "changed"
	require.False(t, svc.ownsRuleState(ctx, foreign))
}

// All writes fail loudly so an unexpected call cannot proceed into unrelated services.
type invalidRuleWriteSpy struct {
	*MemoryRepository
	creates, updates, deletes, links int
}

func (r *invalidRuleWriteSpy) Create(context.Context, *schema.SourceRule) error {
	r.creates++
	return errors.New("unexpected rule create")
}
func (r *invalidRuleWriteSpy) UpdateRevision(context.Context, *schema.SourceRule, *schema.SourceRule) error {
	r.updates++
	return errors.New("unexpected rule update")
}
func (r *invalidRuleWriteSpy) DeleteRevision(context.Context, *schema.SourceRule) error {
	r.deletes++
	return errors.New("unexpected rule delete")
}
func (r *invalidRuleWriteSpy) ReplaceLinks(context.Context, string, []*schema.SourceRuleLink) error {
	r.links++
	return errors.New("unexpected rule links")
}
