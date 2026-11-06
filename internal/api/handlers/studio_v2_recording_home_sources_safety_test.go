package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"
	"go-gateway/internal/datalink/workspace"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecordingHomeSourcesRealGraphProjectionAndNoMutation(t *testing.T) {
	f := newHomeRealFixture(t)
	before := f.snapshot(t)
	out := homeSourcesRequest(f.handler, "?workspace_id="+f.id)
	require.Equal(t, 200, out.Code, out.Body.String())
	var body RecordingHomeSourcesEnvelope
	require.NoError(t, json.Unmarshal(out.Body.Bytes(), &body))
	require.Equal(t, "available", body.Data.State)
	require.Len(t, body.Data.Candidates, 1)
	require.Equal(t, workspace.RecordingHomeCandidate{Key: f.mapping.PointID + "|" + f.tag.ID, DeviceID: "dev-A", PointID: f.mapping.PointID, TagID: f.tag.ID, TagKey: f.tag.Key, Label: "Temperature", DeviceName: "dev-A", Address: "40001", TargetType: "float64"}, body.Data.Candidates[0])
	require.Empty(t, body.Data.Excluded)
	require.Len(t, body.Data.SourceDigest, 64)
	require.Equal(t, before, f.snapshot(t))
	for _, secret := range []string{"connection_config", "127.0.0.1", "last_value", "last_error", "transform_pipeline", "password="} {
		require.NotContains(t, out.Body.String(), secret)
	}
}
func TestRecordingHomeSourcesScopeValidation(t *testing.T) {
	for _, query := range []string{"", "?workspace_id=", "?workspace_id=x&workspace_id=x", "?workspace_id=%20x", "?workspace_id=x%00", "?workspace_id=%zz", "?workspace_id=x&extra=1", "?workspace_id=x;extra=1"} {
		t.Run(query, func(t *testing.T) {
			var h *StudioV2RecordingHomeSourcesHandler
			out := homeSourcesRequest(h, query)
			require.Equal(t, 400, out.Code)
			require.Contains(t, out.Body.String(), "invalid-workspace-scope")
		})
	}
	f := newHomeRealFixture(t)
	require.Equal(t, 404, homeSourcesRequest(f.handler, "?workspace_id=other").Code)
	f.handler.workspace = homeWorkspaceFunc(func(context.Context) (*workspace.Record, error) { return nil, nil })
	require.Equal(t, 404, homeSourcesRequest(f.handler, "?workspace_id="+f.id).Code)
}
func TestRecordingHomeSourcesNilDependencies(t *testing.T) {
	cases := []struct {
		name    string
		replace func(*StudioV2RecordingHomeSourcesHandler)
	}{
		{"workspace", func(h *StudioV2RecordingHomeSourcesHandler) { h.workspace = nil }}, {"typed workspace", func(h *StudioV2RecordingHomeSourcesHandler) { h.workspace = (*workspace.Service)(nil) }},
		{"device", func(h *StudioV2RecordingHomeSourcesHandler) { h.devices = nil }}, {"typed device", func(h *StudioV2RecordingHomeSourcesHandler) { h.devices = (*device.Service)(nil) }},
		{"rules", func(h *StudioV2RecordingHomeSourcesHandler) { h.rules = nil }}, {"typed rules", func(h *StudioV2RecordingHomeSourcesHandler) { h.rules = (*homeRuleStub)(nil) }},
		{"points", func(h *StudioV2RecordingHomeSourcesHandler) { h.points = nil }}, {"typed points", func(h *StudioV2RecordingHomeSourcesHandler) { h.points = (*point.Service)(nil) }},
		{"mappings", func(h *StudioV2RecordingHomeSourcesHandler) { h.mappings = nil }}, {"typed mappings", func(h *StudioV2RecordingHomeSourcesHandler) { h.mappings = (*mapping.Service)(nil) }},
		{"tags", func(h *StudioV2RecordingHomeSourcesHandler) { h.tags = nil }}, {"typed tags", func(h *StudioV2RecordingHomeSourcesHandler) { h.tags = (*tag.Service)(nil) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newHomeRealFixture(t)
			tc.replace(f.handler)
			require.Equal(t, 503, homeSourcesRequest(f.handler, "?workspace_id="+f.id).Code)
		})
	}
	var h *StudioV2RecordingHomeSourcesHandler
	require.Equal(t, 503, homeSourcesRequest(h, "?workspace_id=W1").Code)
}
func TestRecordingHomeSourcesReadFailuresAreSafe(t *testing.T) {
	bad := errors.New("password=secret private.host")
	cases := []struct {
		name    string
		replace func(*StudioV2RecordingHomeSourcesHandler)
	}{
		{"workspace", func(h *StudioV2RecordingHomeSourcesHandler) {
			h.workspace = homeWorkspaceFunc(func(context.Context) (*workspace.Record, error) { return nil, bad })
		}},
		{"device", func(h *StudioV2RecordingHomeSourcesHandler) {
			h.devices = homeDeviceFunc(func(context.Context, string) (*schema.Device, error) { return nil, bad })
		}},
		{"rules", func(h *StudioV2RecordingHomeSourcesHandler) {
			h.rules = &homeRuleStub{list: func(context.Context, []string) ([]*schema.SourceRule, error) { return nil, bad }}
		}},
		{"links", func(h *StudioV2RecordingHomeSourcesHandler) {
			h.rules = &homeRuleStub{RecordingHomeRuleReader: h.rules, links: func(context.Context, string) ([]*schema.SourceRuleLink, error) { return nil, bad }}
		}},
		{"points", func(h *StudioV2RecordingHomeSourcesHandler) {
			h.points = homePointFunc(func(context.Context, point.ListFilter) ([]*schema.Point, error) { return nil, bad })
		}},
		{"mapping list", func(h *StudioV2RecordingHomeSourcesHandler) {
			h.mappings = &homeMappingStub{list: func(context.Context, mapping.ListFilter) ([]*schema.Mapping, error) { return nil, bad }}
		}},
		{"mapping get", func(h *StudioV2RecordingHomeSourcesHandler) {
			h.mappings = &homeMappingStub{RecordingHomeMappingReader: h.mappings, get: func(context.Context, string) (*schema.Mapping, error) { return nil, bad }}
		}},
		{"tag", func(h *StudioV2RecordingHomeSourcesHandler) {
			h.tags = homeTagFunc(func(context.Context, string) (*schema.Tag, error) { return nil, bad })
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newHomeRealFixture(t)
			before := f.snapshot(t)
			tc.replace(f.handler)
			out := homeSourcesRequest(f.handler, "?workspace_id="+f.id)
			require.Equal(t, http.StatusServiceUnavailable, out.Code)
			require.JSONEq(t, `{"success":false,"code":"source-unavailable"}`, out.Body.String())
			require.Equal(t, before, f.snapshot(t))
		})
	}
}
func TestRecordingHomeSourcesDanglingRecordsAreUnknown(t *testing.T) {
	f := newHomeRealFixture(t)
	f.handler.devices = homeDeviceFunc(func(context.Context, string) (*schema.Device, error) { return nil, sql.ErrNoRows })
	requireHomeUnknown(t, f)
	f = newHomeRealFixture(t)
	f.handler.tags = homeTagFunc(func(context.Context, string) (*schema.Tag, error) { return nil, tag.ErrTagNotFound })
	requireHomeUnknown(t, f)
	f = newHomeRealFixture(t)
	f.handler.mappings = &homeMappingStub{RecordingHomeMappingReader: f.handler.mappings, get: func(context.Context, string) (*schema.Mapping, error) { return nil, mapping.ErrMappingNotFound }}
	requireHomeUnknown(t, f)
}
