package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	swagger "go-gateway/docs/swagger"
	"go-gateway/internal/api"
	"go-gateway/internal/datalink/workspace"
)

const successorSwaggerRoot = "/v1/datalink/studio-v2/workspace/write-groups/{id}"

type successorSwaggerSchema struct {
	Ref        string                            `json:"$ref"`
	Type       string                            `json:"type"`
	Properties map[string]successorSwaggerSchema `json:"properties"`
	Required   []string                          `json:"required"`
	Enum       []string                          `json:"enum"`
	Items      *successorSwaggerSchema           `json:"items"`
	AllOf      []successorSwaggerSchema          `json:"allOf"`
}

type successorSwaggerAction struct {
	SuccessStatuses []int  `json:"success_statuses"`
	Error           string `json:"error"`
	Request         string `json:"request"`
	ResponseData    string `json:"response_data"`
}

type successorSwaggerOperation struct {
	recordingSwaggerOperation
	Description string                            `json:"description"`
	Actions     map[string]successorSwaggerAction `json:"x-action-contracts"`
}

type successorSwaggerDocument struct {
	BasePath    string                                          `json:"basePath"`
	Definitions map[string]successorSwaggerSchema               `json:"definitions"`
	Paths       map[string]map[string]successorSwaggerOperation `json:"paths"`
}

func readSuccessorSwagger(t *testing.T) successorSwaggerDocument {
	t.Helper()
	var doc successorSwaggerDocument
	if err := json.Unmarshal([]byte(swagger.SwaggerInfo.ReadDoc()), &doc); err != nil {
		t.Fatal(err)
	}
	return doc
}

func TestStudioV2LifecycleSuccessorSwaggerRoutes(t *testing.T) {
	doc := readSuccessorSwagger(t)
	router := api.NewRouter(&api.DatalinkServices{Workspace: workspace.NewService(nil)})
	for suffix, method := range map[string]string{"lifecycle": "get", "successor-review": "post", "successor": "post"} {
		t.Run(suffix, func(t *testing.T) {
			path := successorSwaggerRoot + "/" + suffix
			fullURL := "/api/v1/datalink/studio-v2/workspace/write-groups/:id/" + suffix
			if got := strings.ReplaceAll(doc.BasePath+path, "{id}", ":id"); got != fullURL {
				t.Fatalf("composed Swagger URL = %q, want %q", got, fullURL)
			}
			if !slices.ContainsFunc(router.Routes(), func(route gin.RouteInfo) bool {
				return route.Path == fullURL && route.Method == strings.ToUpper(method)
			}) {
				t.Fatalf("actual router missing %s %s", method, fullURL)
			}
			op, ok := doc.Paths[path][method]
			if !ok {
				t.Fatalf("missing %s %s (actual router URL %s)", method, path, fullURL)
			}
			for _, status := range []string{"404", "409", "422", "500", "503"} {
				if op.Responses[status].Schema.Ref != "#/definitions/handlers.APIErrorResponse" {
					t.Errorf("%s must document typed %s failure", suffix, status)
				}
			}
			for _, parameter := range op.Parameters {
				if parameter.In == "path" && parameter.Name == "id" && parameter.Required {
					return
				}
			}
			t.Error("missing required predecessor/group path ID")
		})
	}
	for path := range doc.Paths {
		if strings.HasPrefix(path, successorSwaggerRoot+"/successor/") || path == strings.TrimPrefix(successorSwaggerRoot, "/v1")+"/lifecycle" {
			t.Errorf("fabricated stage URL or incomplete lifecycle URL: %s", path)
		}
	}
}

func TestStudioV2SuccessorSwaggerActionContracts(t *testing.T) {
	doc := readSuccessorSwagger(t)
	op := doc.Paths[successorSwaggerRoot+"/successor"]["post"]
	cases := []struct {
		action, request, result string
		required                []string
	}{
		{"reserve", "WriteGroupSuccessorReserveRequest", "workspace.WriteGroupSuccessorPreparation", []string{"action", "envelope", "review_digest"}},
		{"schema-preview", "WriteGroupSuccessorPreviewRequest", "recordingplan.SchemaPreviewToken", []string{"action", "envelope", "review_digest", "current"}},
		{"schema-confirm", "WriteGroupSuccessorConfirmRequest", "recordingplan.SchemaOperation", []string{"action", "envelope", "review_digest", "current", "token", "digest", "schema_operation_id"}},
		{"activation-review", "WriteGroupSuccessorActivationReviewRequest", "workspace.WriteGroupSuccessorActivationReview", []string{"action", "stage", "boundary"}},
		{"activation-readiness", "WriteGroupSuccessorReadinessRequest", "workspace.WriteGroupSuccessorActivationReadiness", []string{"action", "envelope", "review_digest"}},
		{"activation-confirm", "WriteGroupSuccessorCutoverConfirmRequest", "workspace.WriteGroupCutover", []string{"action", "envelope", "review_digest", "readiness_digest"}},
		{"first-start-confirm", "WriteGroupSuccessorFirstStartConfirmRequest", "workspace.WriteGroupSuccessorFirstStartResult", []string{"action", "envelope", "review_digest", "readiness_digest"}},
		{"operation-status", "WriteGroupSuccessorStatusRequest", "workspace.WriteGroupSuccessorOperationStatus", []string{"action", "workspace_id", "lineage_id", "operation_id"}},
	}
	if len(op.Actions) != len(cases) {
		t.Errorf("got %d action contracts, want %d", len(op.Actions), len(cases))
	}
	for _, tc := range cases {
		t.Run(tc.action, func(t *testing.T) {
			action := op.Actions[tc.action]
			if action.Request != "#/definitions/handlers."+tc.request || action.ResponseData != "#/definitions/"+tc.result {
				t.Fatalf("wrong action-specific request/result: %+v", action)
			}
			expectedStatuses := []int{200}
			if tc.action == "schema-confirm" {
				expectedStatuses = append(expectedStatuses, 202)
			}
			if !slices.Equal(action.SuccessStatuses, expectedStatuses) || action.Error != "#/definitions/handlers.APIErrorResponse" {
				t.Errorf("wrong action-specific success/error contract: %+v", action)
			}
			request := doc.Definitions["handlers."+tc.request]
			if len(request.Properties) != len(tc.required) || len(request.Required) != len(tc.required) {
				t.Fatalf("request must have exactly its action's fields and requirements: %+v", request)
			}
			for _, field := range tc.required {
				if _, exists := request.Properties[field]; !exists || !slices.Contains(request.Required, field) {
					t.Errorf("missing required %s", field)
				}
			}
			if !slices.Equal(request.Properties["action"].Enum, []string{tc.action}) {
				t.Errorf("request action discriminator is not fixed to %s", tc.action)
			}
			if _, exists := doc.Definitions[tc.result]; !exists {
				t.Errorf("missing response data model %s", tc.result)
			}
		})
	}
	body := recordingSwaggerBodyParameter(op.Parameters)
	if body == nil || !body.Required {
		t.Fatal("required dispatcher body missing")
	}
	discriminator := doc.Definitions[strings.TrimPrefix(body.Schema.Ref, "#/definitions/")]
	if len(discriminator.Properties) != 1 || len(discriminator.Properties["action"].Enum) != 8 {
		t.Fatal("dispatcher must expose only a discriminator, not a mix-and-match mutation object")
	}
	for _, fragment := range []string{"Swagger 2", "x-action-contracts", "immutable", "current", "schema_operation_id", "first_start", "local read", "not runtime running", "not SQL committed"} {
		if !strings.Contains(op.Description, fragment) {
			t.Errorf("dispatcher description missing %q", fragment)
		}
	}
	if op.Responses["202"].Schema.Ref == "" || !strings.Contains(op.Responses["202"].Description, "schema-confirm") {
		t.Error("schema-confirm must document typed in-progress 202")
	}
}

func TestStudioV2SuccessorSwaggerDispatchMatchesActualMissingScope(t *testing.T) {
	router := api.NewRouter(&api.DatalinkServices{Workspace: workspace.NewService(nil)})
	for _, action := range []string{"reserve", "schema-preview", "schema-confirm", "activation-review", "activation-readiness", "activation-confirm", "first-start-confirm", "operation-status"} {
		t.Run(action, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/datalink/studio-v2/workspace/write-groups/g1/successor", strings.NewReader(`{"action":"`+action+`"}`))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(rr, req)
			if rr.Code != http.StatusBadRequest || !strings.Contains(rr.Body.String(), "WRITE_GROUP_INVALID_REQUEST") {
				t.Fatalf("missing action scope must fail typed 400, got %d %s", rr.Code, rr.Body.String())
			}
		})
	}
}

func TestStudioV2LifecycleSwaggerAnnotationFullURL(t *testing.T) {
	source, err := os.ReadFile("handlers/studio_v2_workspace_write_group_lifecycle.go")
	if err != nil {
		t.Fatal(err)
	}
	match := regexp.MustCompile(`@Router\s+(\S+)\s+\[get\]`).FindSubmatch(source)
	if len(match) != 2 {
		t.Fatal("lifecycle route annotation missing")
	}
	fullURL := readSuccessorSwagger(t).BasePath + string(match[1])
	if fullURL != "/api"+successorSwaggerRoot+"/lifecycle" {
		t.Fatalf("lifecycle annotation composes to %s, want /api%s/lifecycle", fullURL, successorSwaggerRoot)
	}
}
