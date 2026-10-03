package recordingplan

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func groupPreviewScope() SchemaPreviewScope {
	return SchemaPreviewScope{
		WorkspaceID:       "ws-group",
		WorkspaceRevision: "setup-group-1",
		PlanID:            "group-1",
		PlanRevision:      "group-rev-1",
		ConnectorID:       "db-group",
		ConnectorRevision: "connector-rev-1",
		Dialect:           "sqlite",
		Database:          "/tmp/group.db",
		Schema:            "main",
		SourceDigest:      "sources-1",
		SchemaRevision:    "",
		SchemaDigest:      "",
		GroupLayout: &GroupSchemaLayout{
			TableName:   "readings",
			OwnerColumn: "_gw_owner_123456789012345678901234",
			Columns: []GroupSchemaColumn{
				{Name: "record_id", SQLType: "TEXT", Nullable: false, PrimaryKey: true},
				{Name: "bucket_start", SQLType: "TIMESTAMPTZ", Nullable: false},
				{Name: "temperature", SQLType: "REAL", Nullable: true},
				{Name: "_gw_owner_123456789012345678901234", SQLType: "TEXT", Nullable: false},
			},
		},
	}
}

func groupCompatibleState(layout GroupSchemaLayout) TargetTableInspection {
	state := TargetTableInspection{
		Status:           inspectionExists,
		Columns:          make([]string, 0, len(layout.Columns)),
		ColumnTypes:      make(map[string]string, len(layout.Columns)),
		ColumnNullable:   make(map[string]bool, len(layout.Columns)),
		ColumnPrimaryKey: make(map[string]bool, len(layout.Columns)),
	}
	for _, column := range layout.Columns {
		state.Columns = append(state.Columns, column.Name)
		state.ColumnTypes[column.Name] = column.SQLType
		state.ColumnNullable[column.Name] = column.Nullable
		state.ColumnPrimaryKey[column.Name] = column.PrimaryKey
	}
	return state
}

func groupReceiptCompatibleState() TargetTableInspection {
	state := TargetTableInspection{
		Status:           inspectionExists,
		Columns:          []string{"effect_key", "payload_digest", "committed_at", managedEffectReceiptMarker},
		ColumnTypes:      map[string]string{"effect_key": "TEXT", "payload_digest": "TEXT", "committed_at": "TEXT", managedEffectReceiptMarker: "TEXT"},
		ColumnNullable:   map[string]bool{"effect_key": true, "payload_digest": false, "committed_at": false, managedEffectReceiptMarker: false},
		ColumnPrimaryKey: map[string]bool{"effect_key": true, "payload_digest": false, "committed_at": false, managedEffectReceiptMarker: false},
	}
	return state
}

func TestPrepareSchemaPreview_CanonicalGroupCreatesOwnedTableAndReceipt(t *testing.T) {
	svc, _ := newPreviewService()
	scope := groupPreviewScope()
	var calls []string

	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		calls = append(calls, table)
		return TargetTableInspection{Status: inspectionMissing}, nil
	})

	if err != nil {
		t.Fatalf("canonical group preview: %v", err)
	}
	if token.GroupLayout == nil || token.PlanID != scope.PlanID || token.PlanRevision != scope.PlanRevision {
		t.Fatalf("group scope must bind to the existing token identity: %+v", token)
	}
	if strings.Join(calls, ",") != "readings,"+managedEffectReceiptTable {
		t.Fatalf("group preview must inspect the group table and shared receipt: %v", calls)
	}
	if len(token.Statements) != 2 || strings.Contains(token.Statements[0], "IF NOT EXISTS") || !strings.Contains(token.Statements[0], `CREATE TABLE "readings"`) {
		t.Fatalf("group preview must contain an explicit owned-table CREATE: %v", token.Statements)
	}
	if !strings.Contains(token.Statements[1], managedEffectReceiptMarker) || strings.Contains(token.Statements[1], "IF NOT EXISTS") {
		t.Fatalf("group preview must contain the fixed receipt marker and no IF NOT EXISTS: %v", token.Statements)
	}
	if !strings.Contains(token.Statements[1], `"effect_key" TEXT PRIMARY KEY`) || strings.Contains(token.Statements[1], "record_id") {
		t.Fatalf("group preview must preserve the fixed three-column receipt contract: %v", token.Statements[1])
	}
}

func TestPrepareSchemaPreview_GroupAcceptsCanonicalTableNameBeyondLegacyPrefixLimit(t *testing.T) {
	svc, _ := newPreviewService()
	scope := groupPreviewScope()
	layout := *scope.GroupLayout
	layout.TableName = "gw_group_" + strings.Repeat("a", 32)
	scope.GroupLayout = &layout
	// The workspace scope currently carries the exact group table name in the
	// legacy field; group normalization must validate the layout's 63-byte
	// identifier instead of applying the legacy 40-byte prefix limit.
	scope.TablePrefix = layout.TableName

	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == layout.TableName || table == managedEffectReceiptTable {
			return TargetTableInspection{Status: inspectionMissing}, nil
		}
		return TargetTableInspection{Status: inspectionForbidden}, nil
	})
	if err != nil {
		t.Fatalf("canonical group table name should bypass the legacy prefix limit: %v", err)
	}
	if token.GroupLayout == nil || token.GroupLayout.TableName != layout.TableName || !strings.Contains(token.Statements[0], `CREATE TABLE "`+layout.TableName+`"`) {
		t.Fatalf("group preview did not retain the canonical table name: token=%+v", token)
	}
}

func TestPrepareSchemaPreview_CanonicalGroupCompatibleOwnedTablesAreNoOp(t *testing.T) {
	svc, _ := newPreviewService()
	scope := groupPreviewScope()
	groupState := groupCompatibleState(*scope.GroupLayout)
	receiptState := groupReceiptCompatibleState()

	token, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == scope.GroupLayout.TableName {
			return groupState, nil
		}
		return receiptState, nil
	})

	if err != nil || len(token.Statements) != 0 || token.NoChangeReason != NoChangeSchemaCompatible {
		t.Fatalf("compatible owned group schema must be a verified no-op: token=%+v err=%v", token, err)
	}
}

func TestPrepareSchemaPreview_CanonicalGroupRequiresCompleteInspectionMetadata(t *testing.T) {
	svc, _ := newPreviewService()
	scope := groupPreviewScope()
	state := groupCompatibleState(*scope.GroupLayout)
	state.ColumnPrimaryKey = nil

	_, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
		if table == scope.GroupLayout.TableName {
			return state, nil
		}
		return groupReceiptCompatibleState(), nil
	})

	if !errors.Is(err, ErrTargetInspectionUnconfirmed) {
		t.Fatalf("missing group metadata must fail closed, got %v", err)
	}
}

func TestSQLRepositoryPersistsGroupSchemaFieldsInExistingTokenRow(t *testing.T) {
	repo := NewSQLRepository(recordingRepositoryDB(t))
	scope := groupPreviewScope()
	token := &SchemaPreviewToken{
		Token: "tok-group", OperationID: "op-group", Action: SchemaApplyAction,
		WorkspaceID: scope.WorkspaceID, WorkspaceRevision: scope.WorkspaceRevision,
		PlanID: scope.PlanID, PlanRevision: scope.PlanRevision,
		ConnectorID: scope.ConnectorID, ConnectorRevision: scope.ConnectorRevision,
		Dialect: scope.Dialect, Database: scope.Database, Schema: scope.Schema,
		GroupLayout: scope.GroupLayout, SourceDigest: scope.SourceDigest,
		SchemaRevision: scope.SchemaRevision, SchemaDigest: scope.SchemaDigest,
		Tables: []SchemaPreviewTable{{Name: scope.GroupLayout.TableName, Action: SchemaTableActionCreate}},
	}
	if err := repo.SavePreviewToken(t.Context(), token); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetPreviewToken(t.Context(), token.Token)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.GroupLayout, token.GroupLayout) || got.SourceDigest != token.SourceDigest ||
		got.SchemaRevision != token.SchemaRevision || got.SchemaDigest != token.SchemaDigest || !reflect.DeepEqual(got.Tables, token.Tables) {
		t.Fatalf("group token fields were not persisted in the existing token row: got=%+v want=%+v", got, token)
	}
}

type groupApplyStub struct {
	layout     GroupSchemaLayout
	states     map[string]TargetTableInspection
	executions int
}

func (s *groupApplyStub) inspect(_ context.Context, table string) (TargetTableInspection, error) {
	if state, ok := s.states[table]; ok {
		return state, nil
	}
	return TargetTableInspection{Status: inspectionMissing}, nil
}

func (s *groupApplyStub) execute(_ context.Context, statements []string) (SchemaExecution, error) {
	s.executions++
	for _, table := range []string{s.layout.TableName, managedEffectReceiptTable} {
		if table == s.layout.TableName {
			s.states[table] = groupCompatibleState(s.layout)
		} else {
			s.states[table] = groupReceiptCompatibleState()
		}
	}
	return SchemaExecution{Committed: len(statements)}, nil
}

func (s *groupApplyStub) target() SchemaApplyTarget {
	return SchemaApplyTarget{Inspect: s.inspect, Execute: s.execute}
}

func TestApplySchemaPreview_GroupCreatesBothTablesAndDuplicateDoesNotExecute(t *testing.T) {
	svc, _ := newPreviewService()
	target := &groupApplyStub{layout: *groupPreviewScope().GroupLayout, states: map[string]TargetTableInspection{}}
	token, err := svc.PrepareSchemaPreview(t.Context(), groupPreviewScope(), target.inspect)
	if err != nil {
		t.Fatal(err)
	}

	op, outcome, err := svc.ApplySchemaPreview(t.Context(), token, target.target())
	if err != nil || outcome != ClaimAcquired || op.Status != SchemaOperationSucceeded || target.executions != 1 {
		t.Fatalf("group apply must create and verify once: op=%+v outcome=%q executions=%d err=%v", op, outcome, target.executions, err)
	}
	again, outcome, err := svc.ApplySchemaPreview(t.Context(), token, target.target())
	if err != nil || outcome != ClaimCompleted || again.Status != SchemaOperationSucceeded || target.executions != 1 {
		t.Fatalf("duplicate group confirmation must return retained operation: op=%+v outcome=%q executions=%d err=%v", again, outcome, target.executions, err)
	}
}

func TestApplySchemaPreview_GroupRejectsCallerMutationBeforeClaim(t *testing.T) {
	svc, _ := newPreviewService()
	target := &groupApplyStub{layout: *groupPreviewScope().GroupLayout, states: map[string]TargetTableInspection{}}
	token, err := svc.PrepareSchemaPreview(t.Context(), groupPreviewScope(), target.inspect)
	if err != nil {
		t.Fatal(err)
	}

	mutated := *token.GroupLayout
	mutated.TableName = "caller_mutated"
	token.GroupLayout = &mutated
	token.Statements, err = GenerateGroupSchemaDDL(token.Dialect, mutated)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = svc.ApplySchemaPreview(t.Context(), token, target.target())
	if !errors.Is(err, ErrPreviewTokenStale) || target.executions != 0 {
		t.Fatalf("caller-mutated group token must be rejected before execution: err=%v executions=%d", err, target.executions)
	}
	if _, err := svc.GetSchemaOperation(t.Context(), token.WorkspaceID, token.OperationID); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("caller-mutated group token must not claim an operation: %v", err)
	}
}

func TestApplySchemaPreview_GroupTokenCannotFallbackToLegacyAfterLayoutRemoval(t *testing.T) {
	svc, _ := newPreviewService()
	target := &groupApplyStub{layout: *groupPreviewScope().GroupLayout, states: map[string]TargetTableInspection{}}
	token, err := svc.PrepareSchemaPreview(t.Context(), groupPreviewScope(), target.inspect)
	if err != nil {
		t.Fatal(err)
	}
	token.GroupLayout = nil
	token.SourceDigest = ""
	token.SchemaRevision = ""
	token.SchemaDigest = ""
	token.TablePrefix = "gw_record_"
	token.Statements, err = GenerateManagedSchemaDDL(token.Dialect, token.TablePrefix)
	if err != nil {
		t.Fatal(err)
	}
	token.Tables = nil
	if token.Digest, err = previewDigest(token); err != nil {
		t.Fatal(err)
	}

	_, _, err = svc.ApplySchemaPreview(t.Context(), token, target.target())
	if !errors.Is(err, ErrPreviewTokenStale) || target.executions != 0 {
		t.Fatalf("saved group token must not fall back to legacy execution: err=%v executions=%d", err, target.executions)
	}
	if _, err := svc.GetSchemaOperation(t.Context(), token.WorkspaceID, token.OperationID); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("saved group token fallback must not claim an operation: %v", err)
	}
}

func TestApplySchemaPreview_GroupRejectsNonPersistedTokenBeforeClaim(t *testing.T) {
	svc, repo := newPreviewService()
	target := &groupApplyStub{layout: *groupPreviewScope().GroupLayout, states: map[string]TargetTableInspection{}}
	token, err := svc.PrepareSchemaPreview(t.Context(), groupPreviewScope(), target.inspect)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.DeletePreviewToken(t.Context(), token.Token); err != nil {
		t.Fatal(err)
	}

	_, _, err = svc.ApplySchemaPreview(t.Context(), token, target.target())
	if !errors.Is(err, ErrPreviewTokenNotFound) || target.executions != 0 {
		t.Fatalf("non-persisted group token must be rejected before execution: err=%v executions=%d", err, target.executions)
	}
	if _, err := svc.GetSchemaOperation(t.Context(), token.WorkspaceID, token.OperationID); !errors.Is(err, ErrSchemaOperationNotFound) {
		t.Fatalf("non-persisted group token must not claim an operation: %v", err)
	}
}

func TestMemoryRepositoryPreviewTokenDeepClonesTableColumns(t *testing.T) {
	repo := NewMemoryRepository()
	token := &SchemaPreviewToken{
		Token:  "tok-clone",
		Tables: []SchemaPreviewTable{{Name: "readings", Columns: []string{"record_id"}}},
	}
	if err := repo.SavePreviewToken(t.Context(), token); err != nil {
		t.Fatal(err)
	}
	token.Tables[0].Columns[0] = "mutated-input"
	stored, err := repo.GetPreviewToken(t.Context(), token.Token)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Tables[0].Columns[0] != "record_id" {
		t.Fatalf("mutating the saved input must not alter the repository: %v", stored.Tables[0].Columns)
	}

	stored.Tables[0].Columns[0] = "mutated-output"
	again, err := repo.GetPreviewToken(t.Context(), token.Token)
	if err != nil {
		t.Fatal(err)
	}
	if again.Tables[0].Columns[0] != "record_id" {
		t.Fatalf("mutating a loaded token must not alter the repository: %v", again.Tables[0].Columns)
	}
}

func TestApplySchemaPreview_GroupLostAckRecoversFromEvidenceWithoutReplay(t *testing.T) {
	repo := &failingFinishRepository{MemoryRepository: NewMemoryRepository()}
	svc := NewService(repo)
	target := &groupApplyStub{layout: *groupPreviewScope().GroupLayout, states: map[string]TargetTableInspection{}}
	token, err := svc.PrepareSchemaPreview(t.Context(), groupPreviewScope(), target.inspect)
	if err != nil {
		t.Fatal(err)
	}

	op, outcome, err := svc.ApplySchemaPreview(t.Context(), token, target.target())
	if !errors.Is(err, ErrSchemaOperationUnacknowledged) || outcome != ClaimAcquired || op == nil || target.executions != 1 {
		t.Fatalf("group lost acknowledgement must retain one unresolved execution: op=%+v outcome=%q executions=%d err=%v", op, outcome, target.executions, err)
	}
	ageLastSchemaOperation(t, repo.MemoryRepository, token.OperationID)

	recovered, outcome, err := svc.ApplySchemaPreview(t.Context(), token, target.target())
	if err != nil || outcome != ClaimCompleted || recovered.Status != SchemaOperationSucceeded || recovered.Reason != SchemaReasonRecovered || target.executions != 1 {
		t.Fatalf("group stale operation must recover from existing tables without replay: op=%+v outcome=%q executions=%d err=%v", recovered, outcome, target.executions, err)
	}
}

func TestPrepareSchemaPreview_GroupRejectsUnownedConflictAndPermission(t *testing.T) {
	tests := map[string]struct {
		state TargetTableInspection
		want  error
	}{
		"unowned": {
			state: TargetTableInspection{
				Status:           inspectionExists,
				Columns:          []string{"record_id"},
				ColumnTypes:      map[string]string{"record_id": "TEXT"},
				ColumnNullable:   map[string]bool{"record_id": false},
				ColumnPrimaryKey: map[string]bool{"record_id": true},
			},
			want: ErrIncompatibleExistingTable,
		},
		"type conflict": {
			state: func() TargetTableInspection {
				scope := groupPreviewScope()
				state := groupCompatibleState(*scope.GroupLayout)
				state.ColumnTypes["temperature"] = "TEXT"
				return state
			}(),
			want: ErrIncompatibleExistingTable,
		},
		"permission": {state: TargetTableInspection{Status: inspectionForbidden}, want: ErrTargetPermissionDenied},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			svc, _ := newPreviewService()
			scope := groupPreviewScope()
			_, err := svc.PrepareSchemaPreview(t.Context(), scope, func(_ context.Context, table string) (TargetTableInspection, error) {
				if table == scope.GroupLayout.TableName {
					return test.state, nil
				}
				return groupReceiptCompatibleState(), nil
			})
			if !errors.Is(err, test.want) {
				t.Fatalf("group %s must fail closed with %v, got %v", name, test.want, err)
			}
		})
	}
}

func TestValidatePreviewTokenForApply_GroupBindsSourceSchemaLayoutAndDigest(t *testing.T) {
	svc, repo := newPreviewService()
	scope := groupPreviewScope()
	token, err := svc.PrepareSchemaPreview(t.Context(), scope, (&inspectorStub{}).inspect)
	if err != nil {
		t.Fatal(err)
	}

	mutations := map[string]func(*SchemaPreviewScope){
		"source digest":   func(s *SchemaPreviewScope) { s.SourceDigest = "sources-2" },
		"schema revision": func(s *SchemaPreviewScope) { s.SchemaRevision = "schema-2" },
		"schema digest":   func(s *SchemaPreviewScope) { s.SchemaDigest = "schema-digest-2" },
		"layout": func(s *SchemaPreviewScope) {
			s.GroupLayout.Columns[0].SQLType = "BIGINT"
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			current := groupPreviewScope()
			mutate(&current)
			if _, err := svc.ValidatePreviewTokenForApply(t.Context(), token.Token, current); !errors.Is(err, ErrPreviewTokenStale) {
				t.Fatalf("group %s mutation must stale the token, got %v", name, err)
			}
		})
	}

	tampered := *token
	tampered.SourceDigest = "tampered"
	if err := repo.MemoryRepository.SavePreviewToken(t.Context(), &tampered); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ValidatePreviewTokenForApply(t.Context(), token.Token, scope); !errors.Is(err, ErrPreviewTokenStale) {
		t.Fatalf("group digest mutation must be rejected, got %v", err)
	}
}

func TestValidateAndApplyExpiredGroupTokenOnlyReconcilesStaleOperation(t *testing.T) {
	svc, repo := newPreviewService()
	target := &groupApplyStub{layout: *groupPreviewScope().GroupLayout, states: map[string]TargetTableInspection{}}
	token, err := svc.PrepareSchemaPreview(t.Context(), groupPreviewScope(), target.inspect)
	if err != nil {
		t.Fatal(err)
	}
	op, outcome, err := svc.ClaimSchemaApply(t.Context(), token)
	if err != nil || outcome != ClaimAcquired {
		t.Fatalf("claim group operation: op=%+v outcome=%q err=%v", op, outcome, err)
	}
	repo.mu.Lock()
	stored := repo.tokens[token.Token]
	stored.ExpiresAt = time.Now().UTC().Add(-time.Minute)
	repo.tokens[token.Token] = stored
	ledger := repo.operations[token.OperationID]
	ledger.UpdatedAt = time.Now().UTC().Add(-SchemaOperationLease - time.Minute)
	repo.operations[token.OperationID] = ledger
	repo.mu.Unlock()

	validated, err := svc.ValidatePreviewTokenForApply(t.Context(), token.Token, groupPreviewScope())
	if err != nil {
		t.Fatalf("expired token with an existing operation must remain reconcilable: %v", err)
	}
	recovered, outcome, err := svc.ApplySchemaPreview(t.Context(), validated, target.target())
	if err != nil || outcome != ClaimCompleted || recovered.Status != SchemaOperationUnknown || target.executions != 0 {
		t.Fatalf("stale group operation must inspect without replaying DDL: op=%+v outcome=%q executions=%d err=%v", recovered, outcome, target.executions, err)
	}
}
