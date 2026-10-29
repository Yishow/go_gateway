package api_test

import (
	"slices"
	"strings"
	"testing"
)

func TestStudioV2SuccessorSwaggerPublicModels(t *testing.T) {
	doc := readSuccessorSwagger(t)
	for name, fields := range map[string][]string{
		"workspace.WriteGroupLifecycle":                    {"workspace_id", "workspace_revision", "group_id", "group_revision", "applied_revision", "actions", "history", "successor_seed"},
		"workspace.WriteGroupSuccessorSeed":                {"basis_revision", "basis_digest", "group"},
		"workspace.WriteGroupLifecycleHistory":             {"canonical", "lineage_id", "lineage_revision", "current_group_id", "reserved_group_id", "generations"},
		"workspace.WriteGroupLifecycleGeneration":          {"group_id", "generation", "reserved", "saved_destination", "revisions"},
		"workspace.WriteGroupLifecycleRevision":            {"group_revision", "effective_at", "ends_at", "sources", "destination", "pending"},
		"workspace.WriteGroupSuccessorEnvelope":            {"version", "predecessor_group_id", "generation", "source_applied_revision", "source_definition_digest", "request"},
		"workspace.WriteGroupSuccessorPreparation":         {"workspace_revision", "lineage_id", "lineage_revision", "generation", "predecessor_group_id", "operation_id", "intent_digest", "group"},
		"workspace.WriteGroupSuccessorActivationEnvelope":  {"version", "activation_kind", "request", "schema_revision", "schema_digest", "counter_proof_digest"},
		"workspace.WriteGroupSuccessorActivationReview":    {"envelope", "review_digest", "history", "blockers"},
		"workspace.WriteGroupSuccessorActivationReadiness": {"review", "readiness", "readiness_digest"},
		"workspace.WriteGroupSuccessorFirstStartResult":    {"operation_id", "root_operation_id", "workspace_revision", "lineage_id", "lineage_revision", "predecessor_group_id", "group_id", "applied_revision", "boundary", "generation", "status"},
		"workspace.WriteGroupSuccessorOperationStatus":     {"workspace_id", "lineage_id", "predecessor_group_id", "operation_id", "reserve", "schema_operations", "cutover", "first_start"},
		"workspace.WriteGroupSuccessorSchemaStatus":        {"operation_id", "status", "reason_code", "next_action"},
		"recordingplan.SchemaPreviewToken":                 {"token", "operation_id", "digest", "statements", "group_layout", "expires_at"},
		"recordingplan.SchemaOperation":                    {"operation_id", "action", "status", "executed_statements", "verified_digest"},
	} {
		t.Run(name, func(t *testing.T) {
			definition, ok := doc.Definitions[name]
			if !ok {
				t.Fatalf("missing public DTO %s", name)
			}
			for _, field := range fields {
				if _, ok := definition.Properties[field]; !ok {
					t.Errorf("missing public field %s", field)
				}
			}
		})
	}
	for name, fields := range map[string][]string{
		"handlers.WriteGroupSuccessorReviewRequest": {"workspace_id", "expected_workspace_revision", "lineage_id", "expected_lineage_revision", "expected_predecessor_revision", "expected_predecessor_connector_revision", "expected_target_connector_revision", "operation_id", "boundary", "group"},
		"handlers.WriteGroupSuccessorCurrentCAS":    {"workspace_revision", "lineage_revision", "predecessor_revision", "predecessor_connector_revision", "group_id", "group_revision", "target_connector_revision"},
		"handlers.WriteGroupSuccessorStageRequest":  {"envelope", "review_digest", "current"},
	} {
		definition := doc.Definitions[name]
		for _, field := range fields {
			if !slices.Contains(definition.Required, field) {
				t.Errorf("%s must require %s", name, field)
			}
		}
	}
	for _, name := range []string{"workspace.WriteGroupLifecycle", "workspace.WriteGroupSuccessorOperationStatus", "workspace.WriteGroupSuccessorActivationReadiness", "recordingplan.SchemaOperation"} {
		assertSuccessorSwaggerSafe(t, doc, successorSwaggerSchema{Ref: "#/definitions/" + name}, map[string]bool{})
	}
	first := doc.Definitions["workspace.WriteGroupSuccessorFirstStartResult"]
	for _, invented := range []string{"ended", "cutover_id", "event", "fence", "handoff"} {
		if _, ok := first.Properties[invented]; ok {
			t.Errorf("first-start must not fabricate %s", invented)
		}
	}
}

func assertSuccessorSwaggerSafe(t *testing.T, doc successorSwaggerDocument, schema successorSwaggerSchema, seen map[string]bool) {
	t.Helper()
	if schema.Ref != "" {
		if seen[schema.Ref] {
			return
		}
		seen[schema.Ref] = true
		definition, exists := doc.Definitions[strings.TrimPrefix(schema.Ref, "#/definitions/")]
		if !exists {
			t.Errorf("unresolved schema reference %s", schema.Ref)
			return
		}
		assertSuccessorSwaggerSafe(t, doc, definition, seen)
	}
	for field, property := range schema.Properties {
		if slices.Contains([]string{"token", "detail", "snapshot", "password", "dsn", "connection_config", "scope_key", "owner", "source_capability_proof"}, strings.ToLower(field)) {
			t.Errorf("operator-safe DTO graph exposes private field %q", field)
		}
		assertSuccessorSwaggerSafe(t, doc, property, seen)
	}
	if schema.Items != nil {
		assertSuccessorSwaggerSafe(t, doc, *schema.Items, seen)
	}
	for _, part := range schema.AllOf {
		assertSuccessorSwaggerSafe(t, doc, part, seen)
	}
}
