package dbtarget

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"go-gateway/internal/datalink/schema"

	_ "modernc.org/sqlite"
)

func ownedRowFixture(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "owned.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	db.SetMaxOpenConns(1)
	ctx := context.Background()
	for _, ddl := range []string{
		`CREATE TABLE readings (entity TEXT, temperature REAL, note TEXT)`,
		`INSERT INTO readings VALUES ('line-a', 20.5, 'production'), ('gw-test-op1', 1.5, 'test'), ('gw-test-op2', 2.5, 'other op')`,
	} {
		if _, err := db.ExecContext(ctx, ddl); err != nil {
			t.Fatal(err)
		}
	}
	if err := CreateEffectReceiptTable(ctx, db, schema.DatabaseConnectorKindSQLite, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO gw_effect_receipts VALUES ('effect-op1', 'digest', 't'), ('effect-prod', 'digest', 't')`); err != nil {
		t.Fatal(err)
	}
	return db
}

func ownedRef(owner string) OwnedRowRef {
	return OwnedRowRef{Kind: schema.DatabaseConnectorKindSQLite, TableName: "readings", OwnerColumn: "entity", OwnerValue: owner}
}

func TestOwnedRowsReadAndRemoveTouchOnlyTheOperationsRow(t *testing.T) {
	db := ownedRowFixture(t)
	ctx := context.Background()

	rows, err := ReadOwnedRows(ctx, db, ownedRef("gw-test-op1"), []string{"temperature", "note"})
	if err != nil || len(rows) != 1 || rows[0][0] != 1.5 || rows[0][1] != "test" {
		t.Fatalf("read owned rows: rows=%v err=%v", rows, err)
	}
	removed, err := RemoveOwnedRows(ctx, db, ownedRef("gw-test-op1"), "effect-op1")
	if err != nil || removed != 1 {
		t.Fatalf("remove owned rows: removed=%d err=%v", removed, err)
	}
	var production, other, receipts int
	for query, dest := range map[string]*int{
		`SELECT COUNT(*) FROM readings WHERE entity = 'line-a'`:                    &production,
		`SELECT COUNT(*) FROM readings WHERE entity = 'gw-test-op2'`:               &other,
		`SELECT COUNT(*) FROM gw_effect_receipts WHERE effect_key = 'effect-prod'`: &receipts,
	} {
		if err := db.QueryRowContext(ctx, query).Scan(dest); err != nil {
			t.Fatal(err)
		}
	}
	if production != 1 || other != 1 || receipts != 1 {
		t.Fatalf("neighbors must be untouched: production=%d other-operation=%d production-receipt=%d", production, other, receipts)
	}
	var ownReceipt int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gw_effect_receipts WHERE effect_key = 'effect-op1'`).Scan(&ownReceipt); err != nil || ownReceipt != 0 {
		t.Fatalf("the operation's own receipt is removed with its row: count=%d err=%v", ownReceipt, err)
	}
}

func TestOwnedRowsRefuseValuesOutsideTheTestNamespace(t *testing.T) {
	db := ownedRowFixture(t)
	ctx := context.Background()
	for _, owner := range []string{"line-a", "", "gw-test-", "%", "gw-test"} {
		if _, err := RemoveOwnedRows(ctx, db, ownedRef(owner), ""); !errors.Is(err, ErrOwnedRowRefInvalid) {
			t.Fatalf("owner %q: cleanup must be refused, got %v", owner, err)
		}
		if _, err := ReadOwnedRows(ctx, db, ownedRef(owner), []string{"note"}); !errors.Is(err, ErrOwnedRowRefInvalid) {
			t.Fatalf("owner %q: readback must be refused, got %v", owner, err)
		}
	}
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM readings`).Scan(&count); err != nil || count != 3 {
		t.Fatalf("a refused cleanup deletes nothing: count=%d err=%v", count, err)
	}
}

func TestOwnedRowsCleanupFailureRollsBackTheRowDeletion(t *testing.T) {
	db := ownedRowFixture(t)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `DROP TABLE gw_effect_receipts`); err != nil {
		t.Fatal(err)
	}
	if _, err := RemoveOwnedRows(ctx, db, ownedRef("gw-test-op1"), "effect-op1"); err == nil {
		t.Fatal("a missing receipt table with a receipt to remove is a failure, not a silent skip")
	}
	rows, err := ReadOwnedRows(ctx, db, ownedRef("gw-test-op1"), []string{"note"})
	if err != nil || len(rows) != 1 {
		t.Fatalf("the half-done cleanup must roll back and keep the row: rows=%v err=%v", rows, err)
	}
}

func TestReadEffectReceiptDigest(t *testing.T) {
	db := ownedRowFixture(t)
	digest, found, err := ReadEffectReceiptDigest(context.Background(), db, schema.DatabaseConnectorKindSQLite, "", "effect-op1")
	if err != nil || !found || digest != "digest" {
		t.Fatalf("existing receipt: digest=%q found=%v err=%v", digest, found, err)
	}
	if _, found, err := ReadEffectReceiptDigest(context.Background(), db, schema.DatabaseConnectorKindSQLite, "", "missing"); err != nil || found {
		t.Fatalf("missing receipt: found=%v err=%v", found, err)
	}
}
