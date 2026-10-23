package recordingplan

import (
	"context"
	"encoding/json"
	"testing"
)

func TestCompatibleGroupPreviewSerializesEmptyStatementArray(t *testing.T) {
	svc, _ := newPreviewService()
	scope := groupPreviewScope()
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == scope.GroupLayout.TableName {
			return groupCompatibleState(*scope.GroupLayout), nil
		}
		return groupReceiptCompatibleState(), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(token)
	if err != nil {
		t.Fatal(err)
	}
	var response struct {
		Statements json.RawMessage `json:"statements"`
	}
	if err := json.Unmarshal(encoded, &response); err != nil {
		t.Fatal(err)
	}
	if string(response.Statements) != "[]" {
		t.Fatalf("no-op preview requires an empty statement array, got %s", response.Statements)
	}
	if token.NoChangeReason != NoChangeSchemaCompatible || token.Token == "" || token.OperationID == "" {
		t.Fatal("no-op still requires revision-bound explicit confirmation")
	}
}
