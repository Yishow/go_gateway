package pollinggroup

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestMemoryRepositoryRejectsDuplicateNames(t *testing.T) {
	repo := NewMemoryRepository()
	first := &schema.PollingGroup{ID: "first", Name: "same-name", IntervalMs: 1000}
	require.NoError(t, repo.Create(t.Context(), first))
	require.Error(t, repo.Create(t.Context(), &schema.PollingGroup{ID: "duplicate", Name: first.Name, IntervalMs: 2000}))
	require.Equal(t, 1, repo.Count())
	second := &schema.PollingGroup{ID: "second", Name: "other-name", IntervalMs: 1000}
	require.NoError(t, repo.Create(t.Context(), second))
	second.Name = first.Name
	second.IntervalMs = 2000
	require.Error(t, repo.Update(t.Context(), second))
	stored, err := repo.GetByID(t.Context(), second.ID)
	require.NoError(t, err)
	require.Equal(t, "other-name", stored.Name)
	require.Equal(t, 1000, stored.IntervalMs)
	first.IntervalMs = 3000
	require.NoError(t, repo.Update(t.Context(), first), "unchanged own name remains valid")
	stored, err = repo.GetByID(t.Context(), first.ID)
	require.NoError(t, err)
	require.Equal(t, 3000, stored.IntervalMs)
}
