package api_test

import (
	"net/http"
	"slices"
	"testing"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/api"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

// These read contracts must document the actual mounted URLs, not a second API.
func TestRecordingRecoveryAndRuleEvidenceSwagger(t *testing.T) {
	doc := readSuccessorSwagger(t)
	router := api.NewRouter(&api.DatalinkServices{Workspace: workspace.NewService(nil), RecordingPlan: recordingplan.NewService(nil)})
	root := "/v1/datalink/studio-v2/workspace/"
	for _, suffix := range []string{"database-operations/test-writes", "runtime-context"} {
		t.Run(suffix, func(t *testing.T) {
			path := root + suffix
			op, ok := doc.Paths[path]["get"]
			if !ok {
				t.Fatalf("missing GET %s; actual mounted URL requires /api/v1", path)
			}
			if !slices.ContainsFunc(router.Routes(), func(route gin.RouteInfo) bool {
				return route.Method == http.MethodGet && route.Path == doc.BasePath+path
			}) {
				t.Fatalf("documented URL has no mounted route: %s", doc.BasePath+path)
			}
			if op.Responses["503"].Schema.Ref != "#/definitions/handlers.APIErrorResponse" {
				t.Error("read failure must document typed unavailable response")
			}
			if suffix != "runtime-context" {
				if op.Responses["200"].Schema.Ref != "#/definitions/handlers.TestWriteRecoveryResponse" {
					t.Error("recovery response contract missing")
				}
				return
			}
			for _, name := range []string{"workspace_id", "workspace_revision", "device_id", "rule_id", "rule_revision", "point_id", "mapping_id", "source_list", "q", "limit", "cursor"} {
				if !slices.ContainsFunc(op.Parameters, func(p recordingSwaggerParameter) bool {
					return p.Name == name && p.In == "query" && !p.Required
				}) {
					t.Errorf("missing optional scoped query field %s", name)
				}
			}
			for _, status := range []string{"400", "403", "409"} {
				if op.Responses[status].Schema.Ref != "#/definitions/handlers.APIErrorResponse" {
					t.Errorf("missing typed scope error %s", status)
				}
			}
			response := doc.Definitions["handlers.runtimeWorkspaceContextResponse"]
			if response.Properties["rule_evidence"].Ref != "#/definitions/runtimeview.RuleEvidenceView" {
				t.Error("scoped context must reuse the existing RuleEvidenceView contract")
			}
		})
	}
	if _, exists := doc.Paths["/datalink/studio-v2/workspace/runtime-context"]; exists {
		t.Error("obsolete incomplete runtime-context URL remains documented")
	}
}

// Locator entry and its native read follow-ups share the mounted workspace API.
func TestLocatorEntryReadSwagger(t *testing.T) {
	doc := readSuccessorSwagger(t)
	root := "/v1/datalink/studio-v2/workspace"
	entry := doc.Paths[root]["get"]
	for _, name := range []string{"entry", "workspace_id", "step", "device_id", "rule_id", "mapping_id", "group_id", "group_revision", "editor", "expected_workspace_revision", "expected_group_revision"} {
		if !slices.ContainsFunc(entry.Parameters, func(p recordingSwaggerParameter) bool {
			return p.Name == name && p.In == "query" && !p.Required
		}) {
			t.Errorf("missing optional locator entry query %s", name)
		}
	}
	for _, suffix := range []string{"", "/database-metadata", "/database-operations/test-writes", "/database-operations/{operation_id}", "/measurements", "/recording-plans/capabilities"} {
		t.Run(suffix, func(t *testing.T) {
			op, ok := doc.Paths[root+suffix]["get"]
			if !ok {
				t.Fatalf("missing scoped GET %s", root+suffix)
			}
			if suffix != "" {
				for _, name := range []string{"locator_workspace_id", "locator_workspace_revision"} {
					if !slices.ContainsFunc(op.Parameters, func(p recordingSwaggerParameter) bool {
						return p.Name == name && p.In == "query" && !p.Required
					}) {
						t.Errorf("missing optional locator read query %s", name)
					}
				}
			}
			for _, status := range []string{"400", "403", "409", "503"} {
				if op.Responses[status].Schema.Ref != "#/definitions/handlers.APIErrorResponse" {
					t.Errorf("missing typed locator response %s", status)
				}
			}
			if _, exists := doc.Paths["/datalink/studio-v2/workspace"+suffix]; exists {
				t.Errorf("obsolete URL without /v1 for %s", suffix)
			}
		})
	}
}
