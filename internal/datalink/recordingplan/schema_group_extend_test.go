package recordingplan

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
)

const newManagedValue = "v_aaaaaaaaaaaaaaaaaaaaaaaa"

func extensionFixture() (SchemaPreviewScope, *groupApplyStub) {
	scope := groupPreviewScope()
	original := *scope.GroupLayout
	scope.GroupLayout = &GroupSchemaLayout{TableName: original.TableName, OwnerColumn: original.OwnerColumn, Columns: append(slices.Clone(original.Columns), GroupSchemaColumn{Name: newManagedValue, SQLType: "BIGINT", Nullable: true})}
	target := &groupApplyStub{layout: *scope.GroupLayout, states: map[string]TargetTableInspection{original.TableName: groupCompatibleState(original), managedEffectReceiptTable: groupReceiptCompatibleState()}}
	return scope, target
}

func TestManagedGroupExtensionPreviewConfirmAndReplay(t *testing.T) {
	svc, _ := newPreviewService()
	scope, target := extensionFixture()
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, target.inspect)
	if err != nil {
		t.Fatal(err)
	}
	if target.executions != 0 || len(token.Statements) != 1 || token.Statements[0] != `ALTER TABLE "readings" ADD COLUMN "v_aaaaaaaaaaaaaaaaaaaaaaaa" BIGINT;` {
		t.Fatalf("preview must only add the nullable managed value column: %+v", token)
	}
	if token.Tables[0].Action != "add_columns" || len(token.Tables[0].Columns) != 1 {
		t.Fatalf("missing exact addition disclosure: %+v", token.Tables)
	}
	op, _, err := svc.ApplySchemaPreview(t.Context(), token, target.target())
	if err != nil || op.Status != SchemaOperationSucceeded {
		t.Fatalf("apply: %+v %v", op, err)
	}
	_, claim, err := svc.ApplySchemaPreview(t.Context(), token, target.target())
	if err != nil || claim != ClaimCompleted || target.executions != 1 {
		t.Fatalf("replay executed twice: %v %v", claim, err)
	}
}

func TestManagedGroupExtensionRejectsDriftAndMissingOwnership(t *testing.T) {
	for _, kind := range []string{"owner", "required", "foreign", "type", "readonly-history"} {
		t.Run(kind, func(t *testing.T) {
			svc, _ := newPreviewService()
			scope, target := extensionFixture()
			state := target.states[scope.GroupLayout.TableName]
			switch kind {
			case "owner":
				delete(state.ColumnTypes, scope.GroupLayout.OwnerColumn)
			case "required":
				scope.GroupLayout.Columns[len(scope.GroupLayout.Columns)-1].Nullable = false
			case "foreign":
				state.Columns = append(state.Columns, "foreign_column")
				state.ColumnTypes["foreign_column"] = "TEXT"
				state.ColumnNullable["foreign_column"] = true
				state.ColumnPrimaryKey["foreign_column"] = false
			case "type":
				state.ColumnTypes["temperature"] = "TEXT"
			case "readonly-history":
				state.Columns = append(state.Columns, "v_bbbbbbbbbbbbbbbbbbbbbbbb")
				state.ColumnTypes["v_bbbbbbbbbbbbbbbbbbbbbbbb"] = "TEXT"
				state.ColumnNullable["v_bbbbbbbbbbbbbbbbbbbbbbbb"] = false
				state.ColumnPrimaryKey["v_bbbbbbbbbbbbbbbbbbbbbbbb"] = false
			}
			target.states[scope.GroupLayout.TableName] = state
			if _, err := svc.PrepareSchemaPreview(t.Context(), scope, target.inspect); err == nil {
				t.Fatal("unsafe extension accepted")
			}
			if target.executions != 0 {
				t.Fatal("preview wrote DDL")
			}
		})
	}
}

func TestManagedGroupExtensionRechecksExactPendingColumns(t *testing.T) {
	svc, _ := newPreviewService()
	scope, target := extensionFixture()
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, target.inspect)
	if err != nil {
		t.Fatal(err)
	}
	target.states[scope.GroupLayout.TableName] = groupCompatibleState(*scope.GroupLayout)
	_, _, err = svc.ApplySchemaPreview(t.Context(), token, target.target())
	if !errors.Is(err, ErrPreviewTokenStale) || target.executions != 0 {
		t.Fatalf("changed target must fail before execution: %v", err)
	}
}

func TestManagedGroupRemovalPreservesHistoricalColumns(t *testing.T) {
	scope, target := extensionFixture()
	// Current layout has fewer members; the physical historical value is retained.
	state := groupCompatibleState(*scope.GroupLayout)
	scope.GroupLayout.Columns = scope.GroupLayout.Columns[:len(scope.GroupLayout.Columns)-1]
	if err := ValidateManagedGroupTableInspection(*scope.GroupLayout, state); err != nil {
		t.Fatal(err)
	}
	svc, _ := newPreviewService()
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == scope.GroupLayout.TableName {
			return state, nil
		}
		return target.states[managedEffectReceiptTable], nil
	})
	if err != nil || len(token.Statements) != 0 {
		t.Fatalf("removal must not drop history: %+v %v", token, err)
	}
}

func TestManagedGroupPostgresExtensionKeepsExactNumericType(t *testing.T) {
	svc, _ := newPreviewService()
	scope, target := extensionFixture()
	scope.Dialect = "postgres"
	scope.GroupLayout.Columns[len(scope.GroupLayout.Columns)-1].SQLType = "NUMERIC(20,0)"
	receipt := target.states[managedEffectReceiptTable]
	receipt.ColumnNullable["effect_key"] = false
	target.states[managedEffectReceiptTable] = receipt
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, target.inspect)
	if err != nil || !strings.Contains(token.Statements[0], "NUMERIC(20,0)") {
		t.Fatalf("numeric extension: %+v %v", token, err)
	}
}

func TestManagedGroupExtensionLostAckRecoversOriginalOperationWithoutReplay(t *testing.T) {
	repo := &failingFinishRepository{MemoryRepository: NewMemoryRepository()}
	svc := NewService(repo)
	scope, target := extensionFixture()
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, target.inspect)
	if err != nil {
		t.Fatal(err)
	}
	if len(token.Tables) == 0 || token.Tables[0].Action != SchemaTableActionAddColumns {
		t.Fatalf("fixture must exercise additive DDL: %+v", token.Tables)
	}
	op, outcome, err := svc.ApplySchemaPreview(t.Context(), token, target.target())
	if !errors.Is(err, ErrSchemaOperationUnacknowledged) || outcome != ClaimAcquired || op == nil || op.OperationID != token.OperationID || target.executions != 1 {
		t.Fatalf("lost additive acknowledgement must keep original operation: op=%+v outcome=%q executions=%d err=%v", op, outcome, target.executions, err)
	}
	pending, outcome, err := svc.ApplySchemaPreview(t.Context(), token, target.target())
	if err != nil || outcome != ClaimInProgress || pending == nil || pending.OperationID != token.OperationID || target.executions != 1 {
		t.Fatalf("active additive retry must not execute again: op=%+v outcome=%q executions=%d err=%v", pending, outcome, target.executions, err)
	}
	ageLastSchemaOperation(t, repo.MemoryRepository, token.OperationID)
	recovered, outcome, err := svc.ApplySchemaPreview(t.Context(), token, target.target())
	if err != nil || outcome != ClaimCompleted || recovered == nil || recovered.OperationID != token.OperationID || recovered.Status != SchemaOperationSucceeded || recovered.Reason != SchemaReasonRecovered || target.executions != 1 {
		t.Fatalf("additive recovery must inspect original operation without replay: op=%+v outcome=%q executions=%d err=%v", recovered, outcome, target.executions, err)
	}
}
