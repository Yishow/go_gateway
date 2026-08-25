package modbusshare

import (
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/assert"
)

func TestCanonicalSharePlanSignatureCoversGeometryAndRuleRevision(t *testing.T) {
	first := DesiredMapping{
		WorkspaceID: "workspace-1", SourceRuleID: "rule-a", SourceRuleRevision: "rev-a", TagID: "tag-a",
		DataType: schema.DataTypeInt16, ShareStartRegister: 40001, ZeroBasedRegister: 0, SpanRegisters: 1, StrideRegisters: 2,
	}
	second := DesiredMapping{
		WorkspaceID: "workspace-1", SourceRuleID: "rule-b", SourceRuleRevision: "rev-b", TagID: "tag-b",
		DataType: schema.DataTypeInt32, ShareStartRegister: 40012, ZeroBasedRegister: 11, SpanRegisters: 2, StrideRegisters: 2,
	}

	plan := NewCanonicalSharePlan("workspace-1", "workspace-rev", "settings-rev", []DesiredMapping{second, first})
	reordered := NewCanonicalSharePlan("workspace-1", "workspace-rev", "settings-rev", []DesiredMapping{first, second})
	assert.Equal(t, plan.Signature, reordered.Signature)
	assert.True(t, plan.Matches(plan.Signature))

	changedGeometry := second
	changedGeometry.ZeroBasedRegister++
	changed := NewCanonicalSharePlan("workspace-1", "workspace-rev", "settings-rev", []DesiredMapping{first, changedGeometry})
	assert.NotEqual(t, plan.Signature, changed.Signature)
	assert.False(t, plan.Matches(changed.Signature))
	changedRevision := first
	changedRevision.SourceRuleRevision = "rev-a-next"
	revised := NewCanonicalSharePlan("workspace-1", "workspace-rev", "settings-rev", []DesiredMapping{changedRevision, second})
	assert.NotEqual(t, plan.Signature, revised.Signature)
}
