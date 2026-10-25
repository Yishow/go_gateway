package dbtarget

import (
	"testing"

	"go-gateway/internal/datalink/snapshot"

	"github.com/stretchr/testify/require"
)

// TestLongRowLegacyLatestProvenanceKeepsItsBytes pins that a long member saved
// before reductions and roles existed still writes the original provenance
// document: no mean block and no source/result/role keys are added, and nothing
// is invented for it.
func TestLongRowLegacyLatestProvenanceKeepsItsBytes(t *testing.T) {
	cases := longCases(t)
	layout, issues := NewLongRowLayout(longSpec(SQLDialectSQLite, longColumnsFor(t, "sqlite"), cases))
	require.Empty(t, issues)
	legacy := cases[0]
	require.Empty(t, legacy.member.Reduction+legacy.member.ValueRole+string(legacy.member.SourceValueType), "the fixture member is a legacy shape")

	cells := cellMap(mustEncode(t, layout, legacy.member, snapshot.MemberOK, "", longSample("s-legacy", legacy.value)))
	// The stored document is compared byte for byte: the original key order, no
	// formatting whitespace and no new key.
	require.Equal(t,
		`{"type":"int64","member":"m-int","sample_id":"s-legacy","source_revision":"src-1","mapping_revision":"map-1",`+
			`"profile":{"kind":"instantaneous","interval":15,"unit":"second"}}`,
		cells["provenance"].(string))
}
