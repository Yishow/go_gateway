package storage_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupSQLiteDB(t *testing.T) *sql.DB {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "writer_test.db")
	dsn := "file:" + dbPath + "?cache=shared&mode=rwc"
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)

	mig := datalink.NewMigrator()
	require.NoError(t, mig.Migrate(db))
	return db
}

func TestSQLiteWriterWrite(t *testing.T) {
	db := setupSQLiteDB(t)
	defer db.Close()

	w := storage.NewSQLiteWriter(db)
	v := 12.34
	r := storage.TimeSeriesRecord{
		TagID:     "tag-1",
		Timestamp: time.Now(),
		ValueNum:  &v,
		Quality:   schema.QualityGood,
		RawValue:  map[string]interface{}{"src": "test"},
	}
	require.NoError(t, w.Write(context.Background(), r))

	var count int
	err := db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM timeseries WHERE tag_id = ?", "tag-1").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestSQLiteWriterWriteBatch(t *testing.T) {
	db := setupSQLiteDB(t)
	defer db.Close()

	w := storage.NewSQLiteWriter(db)
	n1 := 1.0
	n2 := 2.0
	recs := []storage.TimeSeriesRecord{
		{TagID: "tag-b", Timestamp: time.Now(), ValueNum: &n1, Quality: schema.QualityGood},
		{TagID: "tag-b", Timestamp: time.Now().Add(time.Second), ValueNum: &n2, Quality: schema.QualityGood},
	}
	require.NoError(t, w.WriteBatch(context.Background(), recs))

	var count int
	err := db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM timeseries WHERE tag_id = ?", "tag-b").Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
