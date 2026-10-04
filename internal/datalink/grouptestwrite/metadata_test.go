package grouptestwrite

import (
	"errors"
	"testing"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

func configureManagedMetadataTable(t *testing.T, h *harness) {
	t.Helper()
	for _, ddl := range []string{
		`DROP TABLE readings`,
		`CREATE TABLE readings (
			record_id TEXT PRIMARY KEY NOT NULL,
			group_id TEXT NOT NULL,
			device_id TEXT NOT NULL,
			bucket_start TEXT NOT NULL,
			provenance TEXT NOT NULL,
			entity TEXT NOT NULL,
			temperature REAL NOT NULL,
			running INTEGER NOT NULL,
			batch TEXT NOT NULL,
			big INTEGER NOT NULL,
			note TEXT
		)`,
		`INSERT INTO readings (record_id, group_id, device_id, bucket_start, provenance, entity, temperature, running, batch, big, note)
			VALUES ('production-record', 'production-group', 'd1', '2026-01-01T00:00:00Z', '[]', 'line-a', 20.5, 0, 'B-1', 7, 'production')`,
	} {
		if _, err := h.probe.ExecContext(t.Context(), ddl); err != nil {
			t.Fatal(err)
		}
	}
	h.groups.group.RowPolicy = workspace.WriteGroupRowPolicy{
		IntervalSeconds: 10, EntityKeyColumn: "entity", ProvenanceColumn: "provenance",
		RecordKeyColumn: "record_id", BucketStartColumn: "bucket_start",
		GroupIDColumn: "group_id", DeviceIDColumn: "device_id",
	}
}

func TestWriteManagedMetadataColumnsAreIncludedInTheConfirmedRow(t *testing.T) {
	h := newHarness(t, harnessOptions{})
	configureManagedMetadataTable(t, h)

	preview := h.preview()
	outcome, err := h.confirm(preview)
	if err != nil {
		t.Fatalf("confirm: %v", err)
	}
	if outcome.Operation.WriteOutcome != WriteVerified || outcome.Operation.CleanupStatus != CleanupCleaned {
		t.Fatalf("a managed metadata table must accept and clean its test row: %+v", outcome.Operation)
	}
	if note, temperature := h.productionRow(); note != "production" || temperature != 20.5 {
		t.Fatalf("the production neighbor changed: note=%q temperature=%v", note, temperature)
	}
}

func TestWriteManagedMetadataMutationStalesThePreview(t *testing.T) {
	h := newHarness(t, harnessOptions{})
	configureManagedMetadataTable(t, h)
	preview := h.preview()
	h.groups.group.RowPolicy.DeviceIDColumn = ""

	if _, err := h.confirm(preview); !errors.Is(err, recordingplan.ErrPreviewTokenStale) {
		t.Fatalf("metadata layout mutation must stale the preview, got %v", err)
	}
	if h.dest.opened.Load() != 0 {
		t.Fatal("a stale metadata preview must not open the destination")
	}
}

func TestPostgresTestWriteIncludesManagedMetadataColumns(t *testing.T) {
	h := newPGHarness(t, false)
	for _, ddl := range []string{
		`DROP TABLE "` + h.schemaName + `".readings`,
		`CREATE TABLE "` + h.schemaName + `".readings (
			record_id TEXT PRIMARY KEY NOT NULL,
			group_id TEXT NOT NULL,
			device_id TEXT NOT NULL,
			bucket_start TEXT NOT NULL,
			provenance JSONB NOT NULL,
			entity TEXT NOT NULL,
			temperature DOUBLE PRECISION NOT NULL,
			running BOOLEAN NOT NULL,
			batch TEXT NOT NULL,
			big BIGINT NOT NULL,
			note TEXT
		)`,
		`INSERT INTO "` + h.schemaName + `".readings (record_id, group_id, device_id, bucket_start, provenance, entity, temperature, running, batch, big, note)
			VALUES ('production-record', 'production-group', 'd1', '2026-01-01T00:00:00Z', '[]', 'line-a', 20.5, false, 'B-1', 7, 'production')`,
	} {
		if _, err := h.admin.ExecContext(t.Context(), ddl); err != nil {
			t.Fatal(err)
		}
	}
	h.groups.group.RowPolicy.RecordKeyColumn = "record_id"
	h.groups.group.RowPolicy.BucketStartColumn = "bucket_start"
	h.groups.group.RowPolicy.GroupIDColumn = "group_id"
	h.groups.group.RowPolicy.DeviceIDColumn = "device_id"
	h.groups.group.RowPolicy.ProvenanceColumn = "provenance"

	_, outcome := h.confirmPreview()
	if outcome.Operation.WriteOutcome != WriteVerified || outcome.Operation.CleanupStatus != CleanupCleaned {
		t.Fatalf("a PostgreSQL managed metadata table must accept and clean its test row: %+v", outcome.Operation)
	}
	if h.count(`note = 'production'`) != 1 {
		t.Fatal("the production neighbor must remain after test-write cleanup")
	}
}
