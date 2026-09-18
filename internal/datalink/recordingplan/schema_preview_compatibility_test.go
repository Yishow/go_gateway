package recordingplan

import (
	"errors"
	"strings"
	"testing"
)

func TestPrepareSchemaPreview_IncompatibleColumnTypeBlocks(t *testing.T) {
	svc, repo := newPreviewService()
	state := compatibleState("samples")
	state.ColumnTypes = make(map[string]string)
	for _, col := range rangeStateColumns(state) {
		state.ColumnTypes[col] = compatibleColumnType(col)
	}
	// Make val_num the single incompatible column, so the reported column does
	// not depend on map iteration order inside incompatibleColumnType.
	state.ColumnTypes["val_num"] = "varchar(100)"

	stub := &inspectorStub{states: map[string]TargetTableInspection{"gw_record_samples": state}}

	_, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), stub.inspect)

	if !errors.Is(err, ErrIncompatibleExistingTable) || repo.saves != 0 {
		t.Fatalf("table with incompatible column type must block without a token: err=%v saves=%d", err, repo.saves)
	}
	if !strings.Contains(err.Error(), "val_num has incompatible type varchar(100)") {
		t.Fatalf("error must mention incompatible column and type: %v", err)
	}
}

func rangeStateColumns(state TargetTableInspection) []string {
	return state.Columns
}

func compatibleColumnType(col string) string {
	switch col {
	case "observed_at", "received_at":
		return "datetime"
	case "val_num":
		return "real"
	case "val_bool", "is_test":
		return "integer"
	default:
		return "text"
	}
}

func TestPrepareSchemaPreview_CompatibleColumnTypesPass(t *testing.T) {
	svc, _ := newPreviewService()
	state := compatibleState("samples")
	state.ColumnTypes = make(map[string]string)
	for _, col := range state.Columns {
		state.ColumnTypes[col] = compatibleColumnType(col)
	}

	stub := &inspectorStub{states: map[string]TargetTableInspection{"gw_record_samples": state}}

	token, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), stub.inspect)
	if err != nil {
		t.Fatalf("compatible column types must pass preview: %v", err)
	}
	if len(token.Statements) == 0 {
		t.Fatalf("remaining tables should still produce statements: %v", token.Statements)
	}
}
