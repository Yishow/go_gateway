package recordingplan

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func recordingRepositoryDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "recording.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	applyRecordingDDL(t, db)
	return db
}

var recordingMigrationFiles = []string{
	"019_recording_plans_sqlite.up.sql",
	"022_schema_preview_scope_sqlite.up.sql",
	"023_schema_operations_sqlite.up.sql",
}

func applyRecordingDDL(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, name := range recordingMigrationFiles {
		ddl, err := os.ReadFile(filepath.Join("..", "schema", "migrations", name))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.ExecContext(t.Context(), string(ddl)); err != nil {
			t.Fatalf("apply %s: %v", name, err)
		}
	}
}

func TestSQLRepositoryPersistsPreviewScope(t *testing.T) {
	repo := NewSQLRepository(recordingRepositoryDB(t))
	created := time.Date(2026, 9, 16, 1, 2, 3, 0, time.UTC)
	want := SchemaPreviewToken{
		Token: "tok-scope", OperationID: "op-scope", Action: SchemaApplyAction, WorkspaceID: "ws-1", WorkspaceRevision: "setup-1",
		PlanID: "plan-1", PlanRevision: "rev-1", ConnectorID: "db-1", ConnectorRevision: "identity-1", Dialect: "sqlite",
		Database: "/tmp/line-a.db", Schema: "main", TablePrefix: "gw_record_", Statements: []string{"CREATE TABLE gw_record_x (id TEXT)"},
		Tables:         []SchemaPreviewTable{{Name: "gw_record_x", Action: SchemaTableActionCreate, Columns: []string{"id"}}},
		NoChangeReason: "", Digest: strings.Repeat("a", 64), ExpiresAt: created.Add(10 * time.Minute), CreatedAt: created,
	}
	saved := want
	if err := repo.SavePreviewToken(t.Context(), &saved); err != nil {
		t.Fatal(err)
	}

	got, err := repo.GetPreviewToken(t.Context(), "tok-scope")

	if err != nil {
		t.Fatal(err)
	}
	if !got.ExpiresAt.Equal(want.ExpiresAt) || !got.CreatedAt.Equal(want.CreatedAt) {
		t.Fatalf("preview times changed: got %v/%v want %v/%v", got.ExpiresAt, got.CreatedAt, want.ExpiresAt, want.CreatedAt)
	}
	got.ExpiresAt, got.CreatedAt = want.ExpiresAt, want.CreatedAt
	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("preview scope was not persisted:\n got %+v\nwant %+v", *got, want)
	}
}

func TestSQLRepositoryRejectsInvalidPlanJSON(t *testing.T) {
	for _, column := range []string{"members", "streams", "destinations", "retention", "limits"} {
		t.Run(column, func(t *testing.T) {
			db := recordingRepositoryDB(t)
			repo := NewSQLRepository(db)
			plan := &RecordingPlan{ID: "plan-json", WorkspaceID: "ws-json", Name: "JSON fixture"}
			if err := repo.CreatePlan(t.Context(), plan); err != nil {
				t.Fatal(err)
			}
			query := fmt.Sprintf("UPDATE recording_plans SET %s = ? WHERE id = ?", column)
			if _, err := db.ExecContext(t.Context(), query, "{", plan.ID); err != nil {
				t.Fatal(err)
			}
			if got, err := repo.GetPlanByID(t.Context(), plan.ID); err == nil || got != nil || !strings.Contains(err.Error(), column) {
				t.Fatalf("corrupt %s must fail without a partial plan: got=%v err=%v", column, got, err)
			}
			if got, err := repo.ListPlansByWorkspace(t.Context(), plan.WorkspaceID); err == nil || got != nil {
				t.Fatalf("list must reject corrupt plan: got=%v err=%v", got, err)
			}
		})
	}
}

func TestSQLRepositoryRejectsUnencodablePlanBeforeWrite(t *testing.T) {
	db := recordingRepositoryDB(t)
	repo := NewSQLRepository(db)
	plan := &RecordingPlan{ID: "plan-finite", WorkspaceID: "ws-json", Name: "Original"}
	if err := repo.CreatePlan(t.Context(), plan); err != nil {
		t.Fatal(err)
	}
	nan := math.NaN()
	plan.Streams = []PlanStream{{OnChangeDeadband: &nan}}
	plan.Name = "Must not persist"
	if err := repo.UpdatePlan(t.Context(), plan); err == nil {
		t.Fatal("non-finite stream must reject update")
	}
	got, err := repo.GetPlanByID(t.Context(), plan.ID)
	if err != nil || got.Name != "Original" {
		t.Fatalf("failed serialization changed stored plan: got=%v err=%v", got, err)
	}
	plan.ID = "invalid-new-plan"
	if err := repo.CreatePlan(t.Context(), plan); err == nil {
		t.Fatal("non-finite stream must reject create")
	}
	var count int
	if err := db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM recording_plans").Scan(&count); err != nil || count != 1 {
		t.Fatalf("invalid create changed row count: count=%d err=%v", count, err)
	}
}

func TestSQLRepositoryRejectsInvalidPreviewStatements(t *testing.T) {
	db := recordingRepositoryDB(t)
	repo := NewSQLRepository(db)
	preview := &SchemaPreviewToken{Token: "preview-json", Statements: []string{"SELECT 1"}, ExpiresAt: time.Now().Add(time.Minute)}
	if err := repo.SavePreviewToken(t.Context(), preview); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(t.Context(), "UPDATE managed_schema_preview_tokens SET statements = ?", "{"); err != nil {
		t.Fatal(err)
	}
	if got, err := repo.GetPreviewToken(t.Context(), preview.Token); err == nil || got != nil {
		t.Fatalf("corrupt preview must fail: got=%v err=%v", got, err)
	}
}
