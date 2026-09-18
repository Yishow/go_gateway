package api_test

import (
	"encoding/json"
	"strings"
	"testing"

	swagger "go-gateway/docs/swagger"
)

const (
	studioV2RecordingSchemaApplySwaggerPath = "/v1/datalink/studio-v2/workspace/recording-plans/schema-apply"
	studioV2RecordingTestWriteSwaggerPath   = "/v1/datalink/studio-v2/workspace/recording-plans/test-write"
	studioV2RecordingSwaggerBasePath        = "/api"
)

type recordingSwaggerDocument struct {
	BasePath    string                                          `json:"basePath"`
	Definitions map[string]recordingSwaggerSchema               `json:"definitions"`
	Paths       map[string]map[string]recordingSwaggerOperation `json:"paths"`
}

type recordingSwaggerOperation struct {
	Parameters []recordingSwaggerParameter         `json:"parameters"`
	Responses  map[string]recordingSwaggerResponse `json:"responses"`
}

type recordingSwaggerParameter struct {
	In       string                          `json:"in"`
	Name     string                          `json:"name"`
	Required bool                            `json:"required"`
	Schema   recordingSwaggerSchemaReference `json:"schema"`
}

type recordingSwaggerResponse struct {
	Description string                          `json:"description"`
	Schema      recordingSwaggerSchemaReference `json:"schema"`
}

type recordingSwaggerSchemaReference struct {
	Ref     string          `json:"$ref"`
	Type    string          `json:"type"`
	Example json.RawMessage `json:"example"`
}

type recordingSwaggerSchema struct {
	Properties map[string]recordingSwaggerSchemaReference `json:"properties"`
	Required   []string                                   `json:"required"`
}

func TestStudioV2RecordingSwagger(t *testing.T) {
	var doc recordingSwaggerDocument
	if err := json.Unmarshal([]byte(swagger.SwaggerInfo.ReadDoc()), &doc); err != nil {
		t.Fatalf("parse generated swagger: %v", err)
	}

	cases := []struct {
		name         string
		path         string
		expectedURL  string
		implemented  bool
		code         string
		message      string
		required     string
		moreRequired []string
		optional     []string
	}{
		{
			name:         "schema apply",
			path:         studioV2RecordingSchemaApplySwaggerPath,
			expectedURL:  "/api/v1/datalink/studio-v2/workspace/recording-plans/schema-apply",
			implemented:  true,
			required:     "token",
			moreRequired: []string{"operation_id", "expected_workspace_revision", "expected_plan_revision", "expected_connector_revision"},
		},
		{
			name:        "test write",
			path:        studioV2RecordingTestWriteSwaggerPath,
			expectedURL: "/api/v1/datalink/studio-v2/workspace/recording-plans/test-write",
			implemented: false,
			code:        "RECORDING_TEST_WRITE_NOT_IMPLEMENTED",
			message:     "recording test write is not implemented",
			required:    "plan_id",
			optional:    []string{"stream_id", "table_prefix"},
		},
	}

	if doc.BasePath != studioV2RecordingSwaggerBasePath {
		t.Fatalf("swagger base path = %q, want %q", doc.BasePath, studioV2RecordingSwaggerBasePath)
	}
	errorDefinition := doc.Definitions["handlers.APIErrorResponse"]
	success := errorDefinition.Properties["success"]
	if success.Type != "boolean" || string(success.Example) != "false" {
		t.Fatal("error response must document a boolean success=false example")
	}
	if errorDefinition.Properties["error"].Ref != "#/definitions/handlers.TypedAPIErrorEnvelope" {
		t.Fatal("error response must reference the typed error envelope")
	}
	envelope := doc.Definitions["handlers.TypedAPIErrorEnvelope"]
	for field, expectedType := range map[string]string{
		"code": "string", "message": "string", "action": "string", "request_id": "string", "retryable": "boolean",
	} {
		if envelope.Properties[field].Type != expectedType {
			t.Fatalf("error field %s type = %q, want %q", field, envelope.Properties[field].Type, expectedType)
		}
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := doc.BasePath + tc.path; got != tc.expectedURL {
				t.Fatalf("documented URL = %q, want %q", got, tc.expectedURL)
			}

			methods, ok := doc.Paths[tc.path]
			if !ok {
				t.Fatalf("swagger path %q is missing", tc.path)
			}
			op, ok := methods["post"]
			if !ok {
				t.Fatalf("swagger path %q has no POST operation", tc.path)
			}
			if tc.implemented {
				if _, ok := op.Responses["200"]; !ok {
					t.Fatal("implemented schema apply must document a 200 response")
				}
				if _, ok := op.Responses["202"]; !ok {
					t.Fatal("implemented schema apply must document a 202 running response")
				}
			} else {
				if _, ok := op.Responses["200"]; ok {
					t.Fatal("unimplemented recording mutation must not document a successful 200 response")
				}
				for _, status := range []string{"400", "501"} {
					response, ok := op.Responses[status]
					if !ok {
						t.Fatalf("response %s is missing", status)
					}
					if response.Schema.Ref != "#/definitions/handlers.APIErrorResponse" {
						t.Fatalf("response %s has no APIErrorResponse schema", status)
					}
				}

				response501 := op.Responses["501"]
				for _, fragment := range []string{
					"success=false",
					tc.code,
					tc.message,
					"retryable=false",
					"action=wait_for_supported_operation",
					"request_id",
				} {
					if !strings.Contains(response501.Description, fragment) {
						t.Fatalf("501 description %q does not contain %q", response501.Description, fragment)
					}
				}
			}

			body := recordingSwaggerBodyParameter(op.Parameters)
			if body == nil {
				t.Fatal("request body parameter is missing")
			}
			if !body.Required {
				t.Fatal("request body parameter must be required")
			}
			definitionName := strings.TrimPrefix(body.Schema.Ref, "#/definitions/")
			definition, ok := doc.Definitions[definitionName]
			if !ok {
				t.Fatalf("request schema %q is missing", body.Schema.Ref)
			}
			allRequired := append([]string{tc.required}, tc.moreRequired...)
			for _, field := range allRequired {
				if !containsRecordingSwaggerField(definition.Required, field) {
					t.Fatalf("request schema required fields = %v, missing %q", definition.Required, field)
				}
				if _, ok := definition.Properties[field]; !ok {
					t.Fatalf("request schema property %q is missing", field)
				}
			}
			for _, field := range tc.optional {
				if _, ok := definition.Properties[field]; !ok {
					t.Fatalf("request schema property %q is missing", field)
				}
				if containsRecordingSwaggerField(definition.Required, field) {
					t.Fatalf("optional request field %q is marked required", field)
				}
			}
		})
	}
}

func recordingSwaggerBodyParameter(parameters []recordingSwaggerParameter) *recordingSwaggerParameter {
	for index := range parameters {
		if parameters[index].In == "body" {
			return &parameters[index]
		}
	}
	return nil
}

func containsRecordingSwaggerField(fields []string, wanted string) bool {
	for _, field := range fields {
		if field == wanted {
			return true
		}
	}
	return false
}
