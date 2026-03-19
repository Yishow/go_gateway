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
