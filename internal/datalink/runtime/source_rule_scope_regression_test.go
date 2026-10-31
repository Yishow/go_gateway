package runtime

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"testing"
)

func TestSourceRuleMutationPreservesOtherRuleSelectedScope(t *testing.T) {
	for _, operation := range []sourcerule.RuntimeReconcileOperation{sourcerule.RuntimeReconcileOperationDisable, sourcerule.RuntimeReconcileOperationUpdate, sourcerule.RuntimeReconcileOperationDelete} {
		t.Run(string(operation), func(t *testing.T) {
			f := newSourceRuleReconcileRuntime(t, false, false)
			defer func() { require.NoError(t, f.service.Stop(t.Context())) }()
			f.service.RemovePoint("point-B")
			p := cloneRuntimeProjection(f.projection)
			p.Points = []*schema.Point{f.projection.Points[0]}
			p.RuleLinks = []*schema.SourceRuleLink{{RuleID: "rule-A", PointID: "point-A"}}
			for n := 1; n < 16; n++ {
				point := *p.Points[0]
				point.ID = fmt.Sprintf("extra-%d", n)
				point.Address = fmt.Sprintf("%d", 40001+n)
				p.Points = append(p.Points, &point)
				rule := "rule-A"
				if n >= 8 {
					rule = "rule-B"
				}
				p.RuleLinks = append(p.RuleLinks, &schema.SourceRuleLink{RuleID: rule, PointID: point.ID})
			}
			p.Mappings = p.Mappings[:1]
			p.Tags = p.Tags[:1]
			f.reader.projection.Store(p)
			req := sourcerule.RuntimeReconcileRequest{Operation: sourcerule.RuntimeReconcileOperationCreate, Scope: sourcerule.RuntimeReconcileScope{RuleID: "rule-B", DeviceID: f.device.ID}}
			require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, f.service.ReconcileSourceRule(t.Context(), req).Status)
			require.Len(t, f.service.scheduler.CollectionStatus(), 9)
			// Persisted removal scope must survive disappearance of the rule links.
			req.AffectedPointIDs = sourceRuleProjectionPointIDs(p, "rule-B")
			req.Operation = operation
			expected := 1
			switch operation {
			case sourcerule.RuntimeReconcileOperationDisable:
				for n := 8; n < 16; n++ {
					p.Points[n].Enabled = false
				}
			case sourcerule.RuntimeReconcileOperationUpdate:
				p.Points = p.Points[:10]
				p.RuleLinks = p.RuleLinks[:10]
				expected = 3
			case sourcerule.RuntimeReconcileOperationDelete:
				p.Points = p.Points[:8]
				p.RuleLinks = p.RuleLinks[:8]
			}
			require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, f.service.ReconcileSourceRule(t.Context(), req).Status)
			require.Len(t, f.service.scheduler.CollectionStatus(), expected)
			require.Empty(t, f.service.scheduler.PollNow([]string{"extra-1"}), "A2 must not become readable")
			if operation == sourcerule.RuntimeReconcileOperationDisable {
				for n := 8; n < 16; n++ {
					p.Points[n].Enabled = true
				}
				req.Operation = sourcerule.RuntimeReconcileOperationEnable
				require.Equal(t, sourcerule.RuntimeReconcileStatusAligned, f.service.ReconcileSourceRule(t.Context(), req).Status)
				require.Len(t, f.service.scheduler.CollectionStatus(), 9)
			}
		})
	}
}
