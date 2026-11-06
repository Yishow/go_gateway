package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRecordingHomeSourcesReadCostFor2000Links(t *testing.T) {
	const count = 2000
	points := make([]*schema.Point, 0, count)
	links := make([]*schema.SourceRuleLink, 0, count)
	maps := map[string]*schema.Mapping{}
	tags := map[string]*schema.Tag{}
	for i := 0; i < count; i++ {
		id := fmt.Sprintf("P%d", i)
		points = append(points, &schema.Point{ID: id, DeviceID: "D1", Address: id, Enabled: true, DataType: schema.DataTypeFloat32})
		links = append(links, &schema.SourceRuleLink{ID: "L" + id, RuleID: "R1", PointID: id, Address: id, MappingID: common.Ptr(id), TagID: common.Ptr(id)})
		maps[id] = &schema.Mapping{ID: id, PointID: id, TagID: id, Enabled: true}
		tags[id] = &schema.Tag{ID: id, Key: id, DisplayName: id}
	}
	calls := map[string]int{}
	h := NewStudioV2RecordingHomeSourcesHandler(
		homeWorkspaceFunc(func(context.Context) (*workspace.Record, error) {
			calls["workspace"]++
			return &workspace.Record{ID: "W1", OrderedDeviceIDs: []string{"D1"}}, nil
		}),
		homeDeviceFunc(func(context.Context, string) (*schema.Device, error) {
			calls["devices"]++
			return &schema.Device{ID: "D1", Name: "Device", Protocol: schema.ProtocolModbusTCP}, nil
		}),
		&homeRuleStub{list: func(context.Context, []string) ([]*schema.SourceRule, error) {
			calls["rules"]++
			return []*schema.SourceRule{{ID: "R1", DeviceID: "D1", Enabled: true, SkippedAddresses: "[]"}}, nil
		}, links: func(context.Context, string) ([]*schema.SourceRuleLink, error) { calls["links"]++; return links, nil }},
		homePointFunc(func(context.Context, point.ListFilter) ([]*schema.Point, error) {
			calls["points"]++
			return points, nil
		}),
		&homeMappingStub{list: func(_ context.Context, f mapping.ListFilter) ([]*schema.Mapping, error) {
			calls["mapping_list"]++
			return []*schema.Mapping{maps[*f.PointID]}, nil
		}, get: func(_ context.Context, id string) (*schema.Mapping, error) {
			calls["mapping_get"]++
			return maps[id], nil
		}},
		homeTagFunc(func(_ context.Context, id string) (*schema.Tag, error) { calls["tags"]++; return tags[id], nil }),
	)
	start := time.Now()
	out := homeSourcesRequest(h, "?workspace_id=W1")
	elapsed := time.Since(start)
	require.Equal(t, 200, out.Code)
	var body RecordingHomeSourcesEnvelope
	require.NoError(t, json.Unmarshal(out.Body.Bytes(), &body))
	require.Equal(t, "available", body.Data.State)
	require.Len(t, body.Data.Candidates, count)
	require.Equal(t, map[string]int{"workspace": 2, "devices": 2, "rules": 2, "links": 2, "points": 2, "mapping_list": 4000, "mapping_get": 4000, "tags": 4000}, calls)
	t.Logf("in-memory adapter only: 2000 links, reads=%v, elapsed=%s, response_bytes=%d; not database or latency-improvement evidence", calls, elapsed, out.Body.Len())
}
