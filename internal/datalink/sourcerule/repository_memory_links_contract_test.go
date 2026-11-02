package sourcerule

import (
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryCreateLinksKeepsFirstGroupOnlyOrdering(t *testing.T) {
	repo := NewMemoryRepository()
	require.NoError(t, repo.CreateLinks(t.Context(), nil))
	require.Empty(t, repo.links)
	one := &schema.SourceRuleLink{RuleID: "one", Address: "b"}
	require.NoError(t, repo.CreateLinks(t.Context(), []*schema.SourceRuleLink{one}))
	require.Equal(t, []*schema.SourceRuleLink{one}, repo.links["one"])
	links := []*schema.SourceRuleLink{{RuleID: "first", Address: "b"}, {RuleID: "second", Address: "z"}, {RuleID: "first", Address: "a"}, {RuleID: "second", Address: "c"}}
	require.NoError(t, repo.CreateLinks(t.Context(), links))
	require.Equal(t, []string{"a", "b"}, []string{repo.links["first"][0].Address, repo.links["first"][1].Address})
	require.Equal(t, []string{"z", "c"}, []string{repo.links["second"][0].Address, repo.links["second"][1].Address})
	require.Equal(t, "b", links[0].Address, "caller order stays unchanged")
}
