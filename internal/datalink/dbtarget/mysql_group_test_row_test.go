package dbtarget

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"go-gateway/internal/datalink/schema"

	_ "modernc.org/sqlite"
)

const groupLongTestOwner = "_gw_owner_0123456789abcdef01234567"

func openGroupLongTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	for _, statement := range []string{
		"CREATE TABLE `gw_long_test` (`record_id` TEXT PRIMARY KEY, `workspace_id` TEXT, `group_id` TEXT, `group_revision` TEXT, `provenance` TEXT, `" + groupLongTestOwner + "` TEXT, `value_text` TEXT, `value_decimal` TEXT)",
		"CREATE TABLE `gw_effect_receipts` (`effect_key` TEXT PRIMARY KEY, `payload_digest` TEXT NOT NULL, `committed_at` TEXT NOT NULL)",
	} {
		if _, err := db.ExecContext(t.Context(), statement); err != nil {
			t.Fatal(err)
		}
	}
	return db
}

func longTestID(char string) string { return strings.Repeat(char, 64) }

func insertLongTestRow(t *testing.T, db *sql.DB, id, sampleID, value string) {
	t.Helper()
	_, err := db.ExecContext(t.Context(), "INSERT INTO `gw_long_test` (`record_id`,`workspace_id`,`group_id`,`group_revision`,`provenance`,`"+groupLongTestOwner+"`,`value_text`) VALUES (?,?,?,?,?,?,?)",
		id, "ws-1", "group-1", "rev-1", fmt.Sprintf(`{"sample_id":%q,"member":"m"}`, sampleID), "gateway", value)
	if err != nil {
		t.Fatal(err)
	}
}

func longTestRef(ids ...string) GroupLongTestRowRef {
	samples := make([]string, len(ids))
	keys := make([]string, len(ids))
	digests := make(map[string]string, len(ids))
	for i := range ids {
		samples[i] = fmt.Sprintf("gw-test-op1-%d", i)
		keys[i] = fmt.Sprintf("effect-%d", i+1)
		digests[keys[i]] = "digest-" + keys[i]
	}
	return GroupLongTestRowRef{
		Kind: schema.DatabaseConnectorKindMySQL, SchemaName: "main", TableName: "gw_long_test",
		OwnerColumn: groupLongTestOwner, OwnerValue: "gateway", ProvenanceColumn: "provenance",
		WorkspaceID: "ws-1", GroupID: "group-1", GroupRevision: "rev-1", OperationOwner: "gw-test-op1",
		RecordIDs: ids, SampleIDs: samples, EffectKeys: keys, EffectDigests: digests,
	}
}

func TestGroupLongTestRowsUseExactRecordIDs(t *testing.T) {
	db := openGroupLongTestDB(t)
	ref := longTestRef(longTestID("a"), longTestID("b"))
	insertLongTestRow(t, db, ref.RecordIDs[0], ref.SampleIDs[0], "owned-1")
	insertLongTestRow(t, db, ref.RecordIDs[1], ref.SampleIDs[1], "owned-2")
	insertLongTestRow(t, db, longTestID("c"), "gw-test-other-0", "production")
	for _, key := range ref.EffectKeys {
		if _, err := db.ExecContext(t.Context(), "INSERT INTO `gw_effect_receipts` (`effect_key`,`payload_digest`,`committed_at`) VALUES (?,?,?)", key, ref.EffectDigests[key], "2026-01-01T00:00:00Z"); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := ReadGroupLongTestRows(t.Context(), db, ref, []string{"record_id", "value_text"})
	if err != nil || len(rows) != 2 {
		t.Fatalf("read exact rows: rows=%d err=%v", len(rows), err)
	}
	removed, err := RemoveGroupLongTestRows(t.Context(), db, ref)
	if err != nil || removed != 2 {
		t.Fatalf("remove exact rows: removed=%d err=%v", removed, err)
	}
	var count int
	if err := db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `gw_long_test` WHERE `record_id`=?", longTestID("c")).Scan(&count); err != nil || count != 1 {
		t.Fatalf("production row was changed: count=%d err=%v", count, err)
	}
	if err := db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `gw_effect_receipts` WHERE `effect_key` IN (?,?)", ref.EffectKeys[0], ref.EffectKeys[1]).Scan(&count); err != nil || count != 0 {
		t.Fatalf("test receipts remain: count=%d err=%v", count, err)
	}
}

func TestGroupLongTestRowsRefuseMissingProofAndForeignRows(t *testing.T) {
	db := openGroupLongTestDB(t)
	ref := longTestRef(longTestID("a"))
	insertLongTestRow(t, db, ref.RecordIDs[0], ref.SampleIDs[0], "owned")
	noKey := ref
	noKey.EffectKeys = nil
	if _, err := RemoveGroupLongTestRows(t.Context(), db, noKey); err == nil {
		t.Fatal("missing receipt proof must abort cleanup")
	}
	noMarker := ref
	noMarker.OwnerValue = "gw-test-op1"
	if _, err := RemoveGroupLongTestRows(t.Context(), db, noMarker); err == nil {
		t.Fatal("an operation owner in the table marker must be refused")
	}
	noSample := ref
	noSample.SampleIDs = nil
	if _, err := ReadGroupLongTestRows(t.Context(), db, noSample, []string{"record_id"}); err == nil {
		t.Fatal("missing sample proof must abort readback")
	}
	foreign := ref
	foreign.SampleIDs = []string{"other-operation-0"}
	if _, err := ReadGroupLongTestRows(t.Context(), db, foreign, []string{"record_id"}); err == nil {
		t.Fatal("foreign row proof must abort readback")
	}
	if _, err := RemoveGroupLongTestRows(t.Context(), db, foreign); err == nil {
		t.Fatal("foreign row proof must abort cleanup")
	}
	var count int
	if err := db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `gw_long_test` WHERE `record_id`=?", ref.RecordIDs[0]).Scan(&count); err != nil || count != 1 {
		t.Fatalf("row changed after refused proof: count=%d err=%v", count, err)
	}
}

func TestGroupLongTestRowsRefuseForeignReceipt(t *testing.T) {
	db := openGroupLongTestDB(t)
	ref := longTestRef(longTestID("a"))
	insertLongTestRow(t, db, ref.RecordIDs[0], ref.SampleIDs[0], "owned")
	if _, err := db.ExecContext(t.Context(), "INSERT INTO `gw_effect_receipts` (`effect_key`,`payload_digest`,`committed_at`) VALUES (?,?,?)", ref.EffectKeys[0], "foreign-digest", "2026-01-01T00:00:00Z"); err != nil {
		t.Fatal(err)
	}
	if _, err := RemoveGroupLongTestRows(t.Context(), db, ref); err == nil {
		t.Fatal("a foreign receipt digest must abort cleanup")
	}
	var count int
	if err := db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM `gw_long_test` WHERE `record_id`=?", ref.RecordIDs[0]).Scan(&count); err != nil || count != 1 {
		t.Fatalf("row changed after foreign receipt: count=%d err=%v", count, err)
	}
}

func TestGroupLongTestCleanupAcceptsCanonicalJSON(t *testing.T) {
	db := openGroupLongTestDB(t)
	ref := longTestRef(longTestID("a"))
	insertLongTestRow(t, db, ref.RecordIDs[0], ref.SampleIDs[0], "owned")
	_, err := db.ExecContext(t.Context(), "UPDATE gw_long_test SET provenance=? WHERE record_id=?",
		`{"member": "m", "sample_id": "gw-test-op1-0"}`, ref.RecordIDs[0])
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(t.Context(), "INSERT INTO gw_effect_receipts VALUES (?,?,?)", ref.EffectKeys[0], ref.EffectDigests[ref.EffectKeys[0]], "2026-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	removed, err := RemoveGroupLongTestRows(t.Context(), db, ref)
	if err != nil || removed != 1 {
		t.Fatalf("canonical JSON cleanup: removed=%d err=%v", removed, err)
	}
}

func TestGroupLongTestRowsRejectProductionSampleNamespace(t *testing.T) {
	db := openGroupLongTestDB(t)
	ref := longTestRef(longTestID("a"))
	ref.OperationOwner, ref.SampleIDs = "production", []string{"production-0"}
	insertLongTestRow(t, db, ref.RecordIDs[0], ref.SampleIDs[0], "production")
	if _, err := RemoveGroupLongTestRows(t.Context(), db, ref); err == nil {
		t.Fatal("a production sample namespace must never authorize test cleanup")
	}
	var count int
	if err := db.QueryRowContext(t.Context(), "SELECT COUNT(*) FROM gw_long_test").Scan(&count); err != nil || count != 1 {
		t.Fatalf("production row must remain: count=%d err=%v", count, err)
	}
}
