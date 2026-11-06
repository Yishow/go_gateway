package handlers

import (
	"context"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordingHomeSourcesRuleAndPointIntegrity(t *testing.T) {
	for _, mode := range []string{"nil rule", "duplicate rule", "foreign rule", "malformed skips", "nil point", "duplicate point", "missing point", "foreign point", "wrong address"} {
		t.Run(mode, func(t *testing.T) {
			f := newHomeRealFixture(t)
			before := f.snapshot(t)
			rr, pr := f.handler.rules, f.handler.points
			f.handler.rules = &homeRuleStub{RecordingHomeRuleReader: rr, list: func(c context.Context, ids []string) ([]*schema.SourceRule, error) {
				v, e := rr.ListByDeviceIDs(c, ids)
				switch mode {
				case "nil rule":
					return append(v, nil), e
				case "duplicate rule":
					return append(v, v[0]), e
				case "foreign rule":
					v[0].DeviceID = "foreign"
				case "malformed skips":
					v[0].SkippedAddresses = `[invalid`
				}
				return v, e
			}}
			f.handler.points = homePointFunc(func(c context.Context, filter point.ListFilter) ([]*schema.Point, error) {
				require.Equal(t, point.ListFilter{}, filter)
				v, e := pr.List(c, filter)
				switch mode {
				case "nil point":
					return append(v, nil), e
				case "duplicate point":
					return append(v, v[0]), e
				case "missing point":
					return nil, nil
				case "foreign point":
					v[0].DeviceID = "foreign"
				case "wrong address":
					v[0].Address = "40002"
				}
				return v, e
			})
			requireHomeUnknown(t, f)
			require.Equal(t, before, f.snapshot(t))
		})
	}
}
func TestRecordingHomeSourcesLinkIntegrityNeverRecovers(t *testing.T) {
	for _, mode := range []string{"nil link", "duplicate link", "duplicate point claim", "foreign rule", "missing point", "wrong address", "lost IDs", "lost mapping ID", "lost tag ID", "stale mapping ID", "stale tag ID"} {
		t.Run(mode, func(t *testing.T) {
			f := newHomeRealFixture(t)
			before := f.snapshot(t)
			rr := f.handler.rules
			f.handler.rules = &homeRuleStub{RecordingHomeRuleReader: rr, links: func(c context.Context, id string) ([]*schema.SourceRuleLink, error) {
				v, e := rr.ListLinks(c, id)
				switch mode {
				case "nil link":
					return append(v, nil), e
				case "duplicate link":
					return append(v, v[0]), e
				case "duplicate point claim":
					other := *v[0]
					other.ID = "L2"
					return append(v, &other), e
				case "foreign rule":
					v[0].RuleID = "other"
				case "missing point":
					v[0].PointID = "other"
				case "wrong address":
					v[0].Address = "40002"
				case "lost IDs":
					v[0].MappingID = nil
					v[0].TagID = nil
				case "lost mapping ID":
					v[0].MappingID = nil
				case "lost tag ID":
					v[0].TagID = nil
				case "stale mapping ID":
					*v[0].MappingID = "other"
				case "stale tag ID":
					*v[0].TagID = "other"
				}
				return v, e
			}}
			requireHomeUnknown(t, f)
			require.Equal(t, before, f.snapshot(t))
		})
	}
}
func TestRecordingHomeSourcesMappingAndTagIntegrity(t *testing.T) {
	for _, mode := range []string{"nil listed mapping", "duplicate mapping", "missing mapping", "foreign point", "wrong tag", "nil direct mapping", "wrong direct ID", "changed direct mapping", "malformed pipeline", "nil tag", "wrong returned tag", "foreign tag labels"} {
		t.Run(mode, func(t *testing.T) {
			f := newHomeRealFixture(t)
			before := f.snapshot(t)
			mr, tr := f.handler.mappings, f.handler.tags
			f.handler.mappings = &homeMappingStub{list: func(c context.Context, filter mapping.ListFilter) ([]*schema.Mapping, error) {
				require.NotNil(t, filter.PointID)
				require.Zero(t, filter.Limit)
				require.Zero(t, filter.Offset)
				require.Nil(t, filter.Enabled)
				v, e := mr.List(c, filter)
				switch mode {
				case "nil listed mapping":
					return []*schema.Mapping{nil}, nil
				case "duplicate mapping":
					return append(v, v[0]), e
				case "missing mapping":
					return nil, nil
				case "foreign point":
					v[0].PointID = "other"
				case "wrong tag":
					v[0].TagID = "other"
				case "malformed pipeline":
					v[0].TransformPipeline = "bad"
				}
				return v, e
			}, get: func(c context.Context, id string) (*schema.Mapping, error) {
				v, e := mr.GetByID(c, id)
				switch mode {
				case "nil direct mapping":
					return nil, nil
				case "wrong direct ID":
					v.ID = "other"
				case "changed direct mapping":
					v.Enabled = false
				case "malformed pipeline":
					v.TransformPipeline = "bad"
				}
				return v, e
			}}
			f.handler.tags = homeTagFunc(func(c context.Context, id string) (*schema.Tag, error) {
				v, e := tr.GetByID(c, id)
				switch mode {
				case "nil tag":
					return nil, nil
				case "wrong returned tag":
					v.ID = "other"
				case "foreign tag labels":
					v.Labels = `{"source":"source-rule","source_rule_id":"other","source_rule_address":"40001"}`
				}
				return v, e
			})
			requireHomeUnknown(t, f)
			require.Equal(t, before, f.snapshot(t))
		})
	}
}
