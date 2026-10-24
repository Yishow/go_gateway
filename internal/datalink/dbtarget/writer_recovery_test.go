package dbtarget

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

type writerRecoveryFixture struct {
	writer     *Writer
	connector  *schema.DatabaseConnector
	mappings   []*schema.DatabaseTargetMapping
	repo       *SQLTargetMappingRepository
	connectors *SQLConnectorRepository
}

func newWriterRecoveryFixture(t *testing.T, kind schema.DatabaseConnectorKind, config ConnectionConfig, schemaName string, grouped bool) writerRecoveryFixture {
	t.Helper()
	db := openMigratedTestDB(t)
	db.SetMaxOpenConns(1)
	tags := tag.NewService(tag.NewSQLRepository(db))
	connectors := NewSQLConnectorRepository(db)
	repo := NewSQLTargetMappingRepository(db)
	connector, err := NewConnectorService(connectors, repo).Create(t.Context(), CreateConnectorRequest{
		Name: "missing-table-recovery", Kind: kind, ConnectionConfig: config,
	})
	require.NoError(t, err)
	writer := NewWriterWithConfig(connectors, repo, WriterConfig{TagReader: tags, FlushTick: make(chan time.Time)})
	t.Cleanup(func() { _ = writer.Close(t.Context()) })
	f := writerRecoveryFixture{writer: writer, connector: connector, repo: repo, connectors: connectors}
	for i, column := range []string{"temperature", "pressure"} {
		tagRecord, err := tags.Create(t.Context(), tag.CreateTagRequest{
			Key: column, DisplayName: column, DataType: schema.DataTypeFloat64,
		})
		require.NoError(t, err)
		mapping := &schema.DatabaseTargetMapping{
			ID: fmt.Sprintf("recovery-%d", i), TagID: tagRecord.ID, ConnectorID: connector.ID,
			TableSchema: schemaName, TableName: "sensor_readings", ColumnName: column,
			WriteMode: schema.DatabaseWriteModeUpsert, TimestampColumn: common.Ptr("ts"), Enabled: true,
			CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		if grouped {
			mapping.GroupKey = common.Ptr("meter")
		}
		require.NoError(t, repo.Create(t.Context(), mapping))
		f.mappings = append(f.mappings, mapping)
	}
	return f
}

func sqliteWriterRecoveryFixture(t *testing.T, grouped bool) (writerRecoveryFixture, *sql.DB, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "target.db")
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	require.NoError(t, db.PingContext(t.Context()))
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return newWriterRecoveryFixture(t, schema.DatabaseConnectorKindSQLite, ConnectionConfig{"dsn": path}, "main", grouped), db, path
}

func TestWriter_RecreatesMissingTableAndRetriesSameSample(t *testing.T) {
	for _, grouped := range []bool{false, true} {
		t.Run(fmt.Sprintf("grouped=%t", grouped), func(t *testing.T) {
			f, db, _ := sqliteWriterRecoveryFixture(t, grouped)
			at := time.Date(2026, 10, 24, 0, 55, 0, 0, time.UTC)
			for cycle := range 2 {
				if cycle > 0 {
					_, err := db.ExecContext(t.Context(), `DROP TABLE sensor_readings`)
					require.NoError(t, err)
				}
				for i, mapping := range f.mappings {
					require.NoError(t, f.writer.WriteTagValue(t.Context(), mapping.TagID, 12.5+float64(i), at))
				}
				if grouped {
					require.NoError(t, f.writer.flushBuckets(t.Context(), at.Add(time.Minute), false))
				}
				var temperature, pressure float64
				require.NoError(t, db.QueryRowContext(t.Context(), `SELECT temperature, pressure FROM sensor_readings`).Scan(&temperature, &pressure))
				require.Equal(t, 12.5, temperature)
				require.Equal(t, 13.5, pressure)
				var count int
				require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM sensor_readings`).Scan(&count))
				require.Equal(t, 1, count)
			}
			updated, err := f.connectors.GetByID(t.Context(), f.connector.ID)
			require.NoError(t, err)
			require.Equal(t, "success", updated.LastSchemaEnsureStatus)
			if grouped {
				require.Equal(t, "success", updated.LastFlushStatus)
			} else {
				require.Equal(t, "success", updated.LastWriteStatus)
			}
		})
	}
}

func TestWriter_MissingTableConcurrentRecovery(t *testing.T) {
	f, db, _ := sqliteWriterRecoveryFixture(t, false)
	var wg sync.WaitGroup
	errors := make(chan error, 8)
	for i := range 8 {
		wg.Go(func() {
			errors <- f.writer.WriteTagValue(t.Context(), f.mappings[0].TagID, float64(i), time.Date(2026, 10, 24, 0, i, 0, 0, time.UTC))
		})
	}
	wg.Wait()
	close(errors)
	for err := range errors {
		require.NoError(t, err)
	}
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM sensor_readings`).Scan(&count))
	require.Equal(t, 8, count)
}

func TestWriter_MissingTableInsertRetriesWithoutDuplicate(t *testing.T) {
	f, db, _ := sqliteWriterRecoveryFixture(t, false)
	for _, mapping := range f.mappings {
		mapping.WriteMode = schema.DatabaseWriteModeInsert
		mapping.TimestampColumn = nil
		require.NoError(t, f.repo.Update(t.Context(), mapping))
	}
	require.NoError(t, f.writer.WriteTagValue(t.Context(), f.mappings[0].TagID, 12.5, time.Now()))
	var count int
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM sensor_readings WHERE temperature = 12.5`).Scan(&count))
	require.Equal(t, 1, count)
}

func TestWriter_RecoveryExcludesOutputsOwnedByWriteGroups(t *testing.T) {
	f, db, _ := sqliteWriterRecoveryFixture(t, false)
	f.writer.suppress = func(_, tagID string) bool { return tagID == f.mappings[1].TagID }
	require.NoError(t, f.writer.WriteTagValue(t.Context(), f.mappings[0].TagID, 12.5, time.Now()))
	var value float64
	require.NoError(t, db.QueryRowContext(t.Context(), `SELECT temperature FROM sensor_readings`).Scan(&value))
	_, err := db.ExecContext(t.Context(), `SELECT pressure FROM sensor_readings`)
	require.Error(t, err, "canonical group members must not be provisioned by the legacy writer")
}

func TestWriter_RecoveryDoesNotAlterExistingTableOrCreateMissingDatabase(t *testing.T) {
	t.Run("existing-table", func(t *testing.T) {
		f, db, _ := sqliteWriterRecoveryFixture(t, false)
		_, err := db.ExecContext(t.Context(), `CREATE TABLE sensor_readings (ts DATETIME PRIMARY KEY, sentinel TEXT)`)
		require.NoError(t, err)
		_, err = db.ExecContext(t.Context(), `INSERT INTO sensor_readings VALUES ('preserved', 'keep')`)
		require.NoError(t, err)
		require.Error(t, f.writer.WriteTagValue(t.Context(), f.mappings[0].TagID, 1.0, time.Now()))
		var sentinel string
		require.NoError(t, db.QueryRowContext(t.Context(), `SELECT sentinel FROM sensor_readings`).Scan(&sentinel))
		require.Equal(t, "keep", sentinel)
		_, err = db.ExecContext(t.Context(), `SELECT temperature FROM sensor_readings`)
		require.Error(t, err, "a missing column must not trigger automatic ALTER")
	})
	t.Run("missing-database", func(t *testing.T) {
		f, db, path := sqliteWriterRecoveryFixture(t, false)
		require.NoError(t, db.Close())
		require.NoError(t, os.Remove(path))
		require.Error(t, f.writer.WriteTagValue(t.Context(), f.mappings[0].TagID, 1.0, time.Now()))
		require.NoFileExists(t, path)
	})
}
