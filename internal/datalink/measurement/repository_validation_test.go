package measurement

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func measurementRepositoryDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "measurement.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	ddl, err := os.ReadFile("../schema/migrations/018_measurement_semantics_sqlite.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(t.Context(), string(ddl)); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestSQLRepositoryRejectsInvalidMeasurementJSON(t *testing.T) {
	for _, column := range []string{"counter_policy", "state_map", "bitmask_labels"} {
		t.Run(column, func(t *testing.T) {
			db := measurementRepositoryDB(t)
			repo := NewSQLRepository(db)
			def := &MeasurementDefinition{ID: "meas-json", WorkspaceID: "ws-json", Name: "JSON fixture"}
			if err := repo.Create(t.Context(), def); err != nil {
				t.Fatal(err)
			}
			query := fmt.Sprintf("UPDATE measurement_definitions SET %s = ? WHERE id = ?", column)
			if _, err := db.ExecContext(t.Context(), query, "{", def.ID); err != nil {
				t.Fatal(err)
			}
			if got, err := repo.GetByID(t.Context(), def.ID); err == nil || got != nil || !strings.Contains(err.Error(), column) {
				t.Fatalf("corrupt %s must fail without a partial definition: got=%v err=%v", column, got, err)
			}
			if got, err := repo.ListByWorkspace(t.Context(), def.WorkspaceID); err == nil || got != nil {
				t.Fatalf("list must reject corrupt definition: got=%v err=%v", got, err)
			}
		})
	}
}

func TestSQLRepositoryRejectsUnencodableMeasurementBeforeWrite(t *testing.T) {
	db := measurementRepositoryDB(t)
	repo := NewSQLRepository(db)
	def := &MeasurementDefinition{ID: "meas-finite", WorkspaceID: "ws-json", Name: "Original"}
	if err := repo.Create(t.Context(), def); err != nil {
		t.Fatal(err)
	}
	nan := math.NaN()
	def.CounterPolicy = &CounterPolicy{ResetThreshold: &nan}
	def.Name = "Must not persist"
	if err := repo.Update(t.Context(), def); err == nil {
		t.Fatal("non-finite policy must reject update")
	}
	got, err := repo.GetByID(t.Context(), def.ID)
	if err != nil || got.Name != "Original" {
		t.Fatalf("failed serialization changed stored definition: got=%v err=%v", got, err)
	}
	def.ID = "invalid-new-definition"
	if err := repo.Create(t.Context(), def); err == nil {
		t.Fatal("non-finite policy must reject create")
	}
	var count int
	if err := db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM measurement_definitions").Scan(&count); err != nil || count != 1 {
		t.Fatalf("invalid create changed row count: count=%d err=%v", count, err)
	}
}
