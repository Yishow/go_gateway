package recordingplan

import (
	"errors"
	"testing"
)

func TestManagedDecimalColumnAcceptsExactStorageTypes(t *testing.T) {
	for _, kind := range []string{"text", "numeric"} {
		t.Run(kind, func(t *testing.T) {
			service, _ := newPreviewService()
			state := compatibleState("samples")
			state.ColumnTypes = map[string]string{"val_dec": kind}
			stub := &inspectorStub{states: map[string]TargetTableInspection{"gw_record_samples": state}}
			if _, err := service.PrepareSchemaPreview(t.Context(), previewScope(), stub.inspect); err != nil {
				t.Fatalf("generated exact decimal representation %q rejected: %v", kind, err)
			}
		})
	}
}

func TestManagedDecimalColumnRejectsLossyStorageTypes(t *testing.T) {
	for _, kind := range []string{"real", "double precision", "integer", "boolean", "numeric(38,12)", "decimal(20,6)"} {
		t.Run(kind, func(t *testing.T) {
			service, repo := newPreviewService()
			state := compatibleState("samples")
			state.ColumnTypes = map[string]string{"val_dec": kind}
			stub := &inspectorStub{states: map[string]TargetTableInspection{"gw_record_samples": state}}
			if _, err := service.PrepareSchemaPreview(t.Context(), previewScope(), stub.inspect); !errors.Is(err, ErrIncompatibleExistingTable) {
				t.Fatalf("lossy decimal representation %q accepted: %v", kind, err)
			}
			if repo.saves != 0 {
				t.Fatal("incompatible decimal representation saved a token")
			}
		})
	}
}
