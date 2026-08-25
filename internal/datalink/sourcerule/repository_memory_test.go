package sourcerule

import (
	"context"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestMemoryRepository_CreateLinks_EmptySlice(t *testing.T) {
	t.Parallel()

	repo := NewMemoryRepository()

	require.NoError(t, repo.CreateLinks(context.Background(), nil))
	require.NoError(t, repo.CreateLinks(context.Background(), []*schema.SourceRuleLink{}))
}

func TestMemoryRepository_CandidateSnapshotCASRejectsStaleRevision(t *testing.T) {
	repo := NewMemoryRepository()
	require.NoError(t, repo.Create(context.Background(), &schema.SourceRule{ID: "rule-1", RevisionID: "r2"}))
	err := repo.ReplaceCandidateSnapshotsAtRevision(context.Background(), []*schema.SourceRuleCandidateSnapshot{{
		SourceRuleID: "rule-1", RevisionID: "r1", CandidateType: schema.SourceRuleCandidateTypeTags,
	}}, "r1")
	require.Error(t, err)
	require.Empty(t, repo.snapshots)
}
