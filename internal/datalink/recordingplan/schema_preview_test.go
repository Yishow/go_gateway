package recordingplan

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"
)

// managedColumns lists the columns of every managed table as defined by the
// managed recording schema; a table is compatible only when all are present.
var managedColumns = map[string][]string{
	"samples": {"workspace_id", "plan_id", "stream_id", "record_id", "measurement_id", "series_epoch", "observed_at",
		"received_at", "quality", "quality_reason", "value_type", "val_num", "val_dec", "val_str", "val_bool", "is_test"},
	"intervals": {"workspace_id", "plan_id", "stream_id", "record_id", "measurement_id", "series_epoch", "interval_start",
		"interval_end", "calculation_revision", "sample_count", "mean_val", "min_val", "max_val", "quantity_delta",
		"known_subtotal", "is_estimated", "is_complete", "is_test"},
	"events": {"workspace_id", "plan_id", "stream_id", "record_id", "measurement_id", "series_epoch", "observed_at",
		"event_type", "state_from", "state_to", "duration_ms", "message", "is_test"},
	"snapshots": {"workspace_id", "plan_id", "stream_id", "record_id", "batch_id", "trigger_id", "observed_at",
		"completeness", "payload_json", "is_test"},
	"definitions": {"workspace_id", "plan_id", "measurement_id", "definition_revision", "series_epoch", "quantity",
		"unit", "semantic_kind", "schema_json", "created_at"},
	"receipts": {"workspace_id", "plan_id", "destination_id", "batch_id", "last_record_id", "delivered_at", "record_count"},
}

var managedTableOrder = []string{"samples", "intervals", "events", "snapshots", "definitions", "receipts"}

func previewScope() SchemaPreviewScope {
	return SchemaPreviewScope{
		WorkspaceID: "ws-1", WorkspaceRevision: "setup-1", PlanID: "plan-1", PlanRevision: "rev-1",
		ConnectorID: "db-1", ConnectorRevision: "identity-1", Dialect: "sqlite", Database: "/tmp/line-a.db",
		Schema: "main", TablePrefix: "gw_record_",
	}
}

type inspectorStub struct {
	states map[string]TargetTableInspection
	err    error
	calls  []string
}

func (s *inspectorStub) inspect(_ context.Context, table string) (TargetTableInspection, error) {
	s.calls = append(s.calls, table)
	if s.err != nil {
		return TargetTableInspection{}, s.err
	}
	if state, ok := s.states[table]; ok {
		return state, nil
	}
	return TargetTableInspection{Status: "missing"}, nil
}

func compatibleState(table string) TargetTableInspection {
	return TargetTableInspection{Status: "exists", Columns: append([]string{}, managedColumns[table]...)}
}

type countingPreviewRepo struct {
	*MemoryRepository
	saves int
}

func (r *countingPreviewRepo) SavePreviewToken(ctx context.Context, token *SchemaPreviewToken) error {
	r.saves++
	return r.MemoryRepository.SavePreviewToken(ctx, token)
}

func newPreviewService() (*Service, *countingPreviewRepo) {
	repo := &countingPreviewRepo{MemoryRepository: NewMemoryRepository()}
	return NewService(repo), repo
}

var hexDigest = regexp.MustCompile(`^[0-9a-f]{64}$`)

func TestPrepareSchemaPreview_BindsScopeAndPersistsProtectedToken(t *testing.T) {
	svc, repo := newPreviewService()
	stub := &inspectorStub{}
	before := time.Now().UTC()

	token, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), stub.inspect)

	if err != nil {
		t.Fatalf("prepare preview: %v", err)
	}
	scope := previewScope()
	if token.WorkspaceID != scope.WorkspaceID || token.WorkspaceRevision != scope.WorkspaceRevision || token.PlanID != scope.PlanID ||
		token.PlanRevision != scope.PlanRevision || token.ConnectorID != scope.ConnectorID || token.ConnectorRevision != scope.ConnectorRevision ||
		token.Dialect != "sqlite" || token.Database != scope.Database || token.Schema != scope.Schema || token.TablePrefix != "gw_record_" {
		t.Fatalf("preview lost its scope: %+v", token)
	}
	if token.Action != SchemaApplyAction || token.Token == "" || token.OperationID == "" || token.OperationID == token.Token {
		t.Fatalf("preview needs a distinct token and operation identity: %+v", token)
	}
	if !hexDigest.MatchString(token.Digest) {
		t.Fatalf("preview digest must be a sha256 hex value, got %q", token.Digest)
	}
	if token.ExpiresAt.Before(before.Add(9*time.Minute)) || token.ExpiresAt.After(before.Add(11*time.Minute)) {
		t.Fatalf("preview expiry must be about ten minutes, got %v", token.ExpiresAt)
	}
	if len(token.Statements) == 0 || !strings.Contains(token.Statements[0], "CREATE TABLE IF NOT EXISTS gw_record_samples") || token.NoChangeReason != "" {
		t.Fatalf("missing tables must produce creation statements: %+v", token)
	}
	if len(token.Tables) != len(managedTableOrder) || token.Tables[0].Name != "gw_record_samples" ||
		token.Tables[0].Action != SchemaTableActionCreate || strings.Join(token.Tables[0].Columns, ",") != strings.Join(managedColumns["samples"], ",") {
		t.Fatalf("preview must summarize each table to create with its columns: %+v", token.Tables)
	}
	wantCalls := []string{"gw_record_samples", "gw_record_intervals", "gw_record_events", "gw_record_snapshots", "gw_record_definitions", "gw_record_receipts"}
	if strings.Join(stub.calls, ",") != strings.Join(wantCalls, ",") {
		t.Fatalf("preview must inspect each managed table once: %v", stub.calls)
	}
	stored, err := repo.GetPreviewToken(t.Context(), token.Token)
	if err != nil || stored.Digest != token.Digest || stored.OperationID != token.OperationID || stored.WorkspaceRevision != "setup-1" {
		t.Fatalf("preview must be persisted with its protection: stored=%+v err=%v", stored, err)
	}
}

func TestPrepareSchemaPreview_DefaultsAnEmptyPrefix(t *testing.T) {
	svc, _ := newPreviewService()
	scope := previewScope()
	scope.TablePrefix = ""

	token, err := svc.PrepareSchemaPreview(t.Context(), scope, (&inspectorStub{}).inspect)

	if err != nil || token.TablePrefix != "gw_record_" {
		t.Fatalf("empty prefix must use the default: token=%+v err=%v", token, err)
	}
}

func TestPrepareSchemaPreview_CompatibleExistingSchemaIsNoChange(t *testing.T) {
	svc, _ := newPreviewService()
	stub := &inspectorStub{states: map[string]TargetTableInspection{}}
	for _, table := range managedTableOrder {
		stub.states["gw_record_"+table] = compatibleState(table)
	}

	token, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), stub.inspect)

	if err != nil {
		t.Fatalf("prepare preview: %v", err)
	}
	if len(token.Statements) != 0 || token.NoChangeReason != NoChangeSchemaCompatible || !hexDigest.MatchString(token.Digest) {
		t.Fatalf("compatible schema must be an explained no-op: %+v", token)
	}
	if len(token.Tables) != len(managedTableOrder) {
		t.Fatalf("compatible schema must still list every managed table: %+v", token.Tables)
	}
	for _, table := range token.Tables {
		if table.Action != SchemaTableActionUnchanged || len(table.Columns) != 0 {
			t.Fatalf("compatible table must be reported unchanged: %+v", table)
		}
	}
}

func TestPrepareSchemaPreview_SkipsOnlyCompatibleExistingTables(t *testing.T) {
	svc, _ := newPreviewService()
	stub := &inspectorStub{states: map[string]TargetTableInspection{"gw_record_samples": compatibleState("samples")}}

	token, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), stub.inspect)

	if err != nil {
		t.Fatalf("prepare preview: %v", err)
	}
	joined := strings.Join(token.Statements, "\n")
	if strings.Contains(joined, "gw_record_samples") || !strings.Contains(joined, "CREATE TABLE IF NOT EXISTS gw_record_intervals") {
		t.Fatalf("only missing tables may be created: %v", token.Statements)
	}
}

func TestPrepareSchemaPreview_IncompatibleExistingTableBlocks(t *testing.T) {
	svc, repo := newPreviewService()
	columns := compatibleState("samples").Columns
	incompatible := TargetTableInspection{Status: "exists"}
	for _, column := range columns {
		if column != "val_num" {
			incompatible.Columns = append(incompatible.Columns, column)
		}
	}
	stub := &inspectorStub{states: map[string]TargetTableInspection{"gw_record_samples": incompatible}}

	_, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), stub.inspect)

	if !errors.Is(err, ErrIncompatibleExistingTable) || repo.saves != 0 {
		t.Fatalf("incompatible table must block without a token: err=%v saves=%d", err, repo.saves)
	}
}

func TestPrepareSchemaPreview_UnconfirmedInspectionCreatesNothing(t *testing.T) {
	for name, stub := range map[string]*inspectorStub{
		"forbidden": {states: map[string]TargetTableInspection{"gw_record_events": {Status: "forbidden"}}},
		"failed":    {states: map[string]TargetTableInspection{"gw_record_samples": {Status: "failed"}}},
		"error":     {err: errors.New("target unreachable")},
	} {
		t.Run(name, func(t *testing.T) {
			svc, repo := newPreviewService()
			_, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), stub.inspect)
			if !errors.Is(err, ErrTargetInspectionUnconfirmed) || repo.saves != 0 {
				t.Fatalf("unconfirmed target must not produce a token: err=%v saves=%d", err, repo.saves)
			}
		})
	}
}

func TestPrepareSchemaPreview_RejectsUnsafePrefix(t *testing.T) {
	for _, prefix := range []string{"gw;drop table x;", "1record_", "gw record_", `gw"_`, strings.Repeat("a", 41)} {
		svc, repo := newPreviewService()
		stub := &inspectorStub{}
		scope := previewScope()
		scope.TablePrefix = prefix
		_, err := svc.PrepareSchemaPreview(t.Context(), scope, stub.inspect)
		if !errors.Is(err, ErrInvalidTablePrefix) || len(stub.calls) != 0 || repo.saves != 0 {
			t.Fatalf("prefix %q must be rejected before inspection: err=%v calls=%v", prefix, err, stub.calls)
		}
	}
}

func TestValidatePreviewTokenForApply_AcceptsTheExactPreviewedScope(t *testing.T) {
	svc, _ := newPreviewService()
	token, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), (&inspectorStub{}).inspect)
	if err != nil {
		t.Fatalf("prepare preview: %v", err)
	}

	validated, err := svc.ValidatePreviewTokenForApply(t.Context(), token.Token, previewScope())

	if err != nil || validated.OperationID != token.OperationID || len(validated.Statements) != len(token.Statements) {
		t.Fatalf("exact scope must validate: validated=%+v err=%v", validated, err)
	}
}

func TestValidatePreviewTokenForApply_RejectsStaleScope(t *testing.T) {
	mutations := map[string]func(*SchemaPreviewScope){
		"workspace revision": func(s *SchemaPreviewScope) { s.WorkspaceRevision = "setup-2" },
		"plan revision":      func(s *SchemaPreviewScope) { s.PlanRevision = "rev-2" },
		"connector":          func(s *SchemaPreviewScope) { s.ConnectorID = "db-2" },
		"connector revision": func(s *SchemaPreviewScope) { s.ConnectorRevision = "identity-2" },
		"dialect":            func(s *SchemaPreviewScope) { s.Dialect = "postgres" },
		"database":           func(s *SchemaPreviewScope) { s.Database = "/tmp/other.db" },
		"schema":             func(s *SchemaPreviewScope) { s.Schema = "aux" },
		"table prefix":       func(s *SchemaPreviewScope) { s.TablePrefix = "gw_other_" },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			svc, _ := newPreviewService()
			token, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), (&inspectorStub{}).inspect)
			if err != nil {
				t.Fatalf("prepare preview: %v", err)
			}
			current := previewScope()
			mutate(&current)
			if _, err := svc.ValidatePreviewTokenForApply(t.Context(), token.Token, current); !errors.Is(err, ErrPreviewTokenStale) {
				t.Fatalf("changed %s must make the preview stale, got %v", name, err)
			}
		})
	}
}

func TestValidatePreviewTokenForApply_RejectsExpiredForeignUnknownLegacyAndTampered(t *testing.T) {
	svc, repo := newPreviewService()
	ctx := t.Context()
	token, err := svc.PrepareSchemaPreview(ctx, previewScope(), (&inspectorStub{}).inspect)
	if err != nil {
		t.Fatalf("prepare preview: %v", err)
	}

	foreign := previewScope()
	foreign.WorkspaceID = "ws-2"
	if _, err := svc.ValidatePreviewTokenForApply(ctx, token.Token, foreign); !errors.Is(err, ErrPreviewTokenNotFound) {
		t.Fatalf("foreign workspace must look like an unknown token, got %v", err)
	}
	if _, err := svc.ValidatePreviewTokenForApply(ctx, "tok-unknown", previewScope()); !errors.Is(err, ErrPreviewTokenNotFound) {
		t.Fatalf("unknown token must be rejected, got %v", err)
	}

	tampered := *token
	tampered.Statements = []string{"DROP TABLE gw_record_samples"}
	if err := repo.MemoryRepository.SavePreviewToken(ctx, &tampered); err != nil {
		t.Fatalf("store tampered token: %v", err)
	}
	if _, err := svc.ValidatePreviewTokenForApply(ctx, token.Token, previewScope()); !errors.Is(err, ErrPreviewTokenStale) {
		t.Fatalf("statements that no longer match the digest must be rejected, got %v", err)
	}

	expired := *token
	expired.ExpiresAt = time.Now().UTC().Add(-time.Minute)
	if err := repo.MemoryRepository.SavePreviewToken(ctx, &expired); err != nil {
		t.Fatalf("store expired token: %v", err)
	}
	if _, err := svc.ValidatePreviewTokenForApply(ctx, token.Token, previewScope()); !errors.Is(err, ErrPreviewTokenExpired) {
		t.Fatalf("expired token must be rejected, got %v", err)
	}

	legacy := &SchemaPreviewToken{
		Token: "tok-legacy", WorkspaceID: "ws-1", PlanID: "plan-1", PlanRevision: "rev-1", ConnectorID: "db-1",
		TablePrefix: "gw_record_", Statements: []string{"CREATE TABLE gw_record_samples (id TEXT)"}, ExpiresAt: time.Now().UTC().Add(time.Hour),
	}
	if err := repo.MemoryRepository.SavePreviewToken(ctx, legacy); err != nil {
		t.Fatalf("store legacy token: %v", err)
	}
	if _, err := svc.ValidatePreviewTokenForApply(ctx, "tok-legacy", previewScope()); !errors.Is(err, ErrPreviewTokenLegacy) {
		t.Fatalf("legacy token without protection must require a new preview, got %v", err)
	}
}

func TestPrepareSchemaPreview_SeparatesPermissionDeniedFromUnavailableTarget(t *testing.T) {
	svc, _ := newPreviewService()
	forbidden := &inspectorStub{states: map[string]TargetTableInspection{"gw_record_intervals": {Status: "forbidden"}}}
	_, err := svc.PrepareSchemaPreview(t.Context(), previewScope(), forbidden.inspect)
	if !errors.Is(err, ErrTargetPermissionDenied) || !errors.Is(err, ErrTargetInspectionUnconfirmed) {
		t.Fatalf("forbidden table must be reported as permission denied, got %v", err)
	}

	failed := &inspectorStub{states: map[string]TargetTableInspection{"gw_record_samples": {Status: "failed"}}}
	_, err = svc.PrepareSchemaPreview(t.Context(), previewScope(), failed.inspect)
	if errors.Is(err, ErrTargetPermissionDenied) || !errors.Is(err, ErrTargetInspectionUnconfirmed) {
		t.Fatalf("failed inspection must stay an unconfirmed target, got %v", err)
	}
}
