package handlers

import (
	"encoding/json"
	"reflect"
	"testing"

	"go-gateway/internal/datalink/workspace"
)

// Documentation projections must match the strict decoder's actual JSON shape;
// required tags and action-specific enums add documentation, not new authority.
func TestSuccessorSwaggerRequestsMatchStrictDecoderModels(t *testing.T) {
	for _, tc := range []struct {
		name             string
		documented, wire any
	}{
		{"review", WriteGroupSuccessorReviewRequest{}, workspace.WriteGroupSuccessorReviewRequest{}},
		{"reserve", WriteGroupSuccessorReserveRequest{}, successorReserveRequest{}},
		{"schema-preview", WriteGroupSuccessorPreviewRequest{}, successorPreviewRequest{}},
		{"schema-confirm", WriteGroupSuccessorConfirmRequest{}, successorConfirmRequest{}},
		{"activation-review", WriteGroupSuccessorActivationReviewRequest{}, successorActivationReviewRequest{}},
		{"activation-readiness", WriteGroupSuccessorReadinessRequest{}, successorActivationReadinessRequest{}},
		{"activation-confirm", WriteGroupSuccessorCutoverConfirmRequest{}, successorActivationConfirmRequest{}},
		{"first-start-confirm", WriteGroupSuccessorFirstStartConfirmRequest{}, successorActivationConfirmRequest{}},
		{"operation-status", WriteGroupSuccessorStatusRequest{}, successorOperationStatusRequest{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			documented, err := json.Marshal(tc.documented)
			if err != nil {
				t.Fatal(err)
			}
			wire, err := json.Marshal(tc.wire)
			if err != nil {
				t.Fatal(err)
			}
			var documentedShape, wireShape any
			if err := json.Unmarshal(documented, &documentedShape); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal(wire, &wireShape); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(documentedShape, wireShape) {
				t.Fatalf("documentation/decoder JSON shape drift:\ndoc: %s\nwire: %s", documented, wire)
			}
			if err := decodeSuccessorRequest(documented, reflect.New(reflect.TypeOf(tc.wire)).Interface()); err != nil {
				t.Fatalf("documented shape is rejected by strict decoder: %v", err)
			}
		})
	}
}
