package datalink

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
)

func TestHandoffConvergenceMigrationFreshUpgradeAndReopen(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		name := "fresh"
		if legacy {
			name = "upgrade-045"
		}
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "handoff.db")
			db := openConvergenceMigrationDB(t, path)
			require.NoError(t, NewMigrator().Migrate(db))
			if legacy {
				setLegacyHandoffIndex(t, db)
			}
			seedConvergenceMigrationRows(t, db)
			before := convergenceMigrationSnapshot(t, db)
			require.NoError(t, NewMigrator().Migrate(db))
			require.Equal(t, 0, convergenceTargetIndexUnique(t, db))
			require.Equal(t, before, convergenceMigrationSnapshot(t, db))
			insertConvergenceAlias(t, db)
			withAlias := convergenceMigrationSnapshot(t, db)
			require.NoError(t, db.Close())
			db = openConvergenceMigrationDB(t, path)
			require.NoError(t, NewMigrator().Migrate(db))
			require.NoError(t, ensureSQLiteSampleHandoffs(db))
			require.Equal(t, 0, convergenceTargetIndexUnique(t, db))
			require.Equal(t, withAlias, convergenceMigrationSnapshot(t, db))
			_, err := db.Exec(`DROP INDEX idx_wg_handoff_target_seq`)
			require.NoError(t, err)
			require.NoError(t, ensureSQLiteSampleHandoffs(db))
			require.Equal(t, 0, convergenceTargetIndexUnique(t, db))
			require.Equal(t, withAlias, convergenceMigrationSnapshot(t, db))
			var count int
			require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_samples`).Scan(&count))
			require.Equal(t, 4, count)
			require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM wg_delivery_sample_handoffs WHERE target_seq=4`).Scan(&count))
			require.Equal(t, 2, count)
			_, err = db.Exec(`INSERT INTO wg_delivery_samples SELECT NULL,workspace_id,group_id,group_revision,'another-address',member_key,observed_at,bucket_start,payload,payload_digest,consumed,created_at,acceptance_kind,origin_workspace_id,origin_group_id,origin_revision,origin_sample_id,'another-transfer' FROM wg_delivery_samples WHERE seq=4`)
			require.Error(t, err, "target origin uniqueness must remain enforced")
			_, err = db.Exec(`INSERT INTO wg_delivery_samples SELECT NULL,workspace_id,group_id,group_revision,'another-native-address',member_key,observed_at,bucket_start,payload,payload_digest,consumed,created_at,acceptance_kind,origin_workspace_id,origin_group_id,origin_revision,origin_sample_id,handoff_id FROM wg_delivery_samples WHERE seq=1`)
			require.Error(t, err, "native logical uniqueness must remain enforced")
		})
	}
}

func TestHandoffConvergenceMigrationAtomicFailureAndLostResponse(t *testing.T) {
	for _, mode := range []string{"ddl-failure", "lost-commit-response"} {
		t.Run(mode, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "fault.db")
			db := openConvergenceMigrationDB(t, path)
			require.NoError(t, NewMigrator().Migrate(db))
			setLegacyHandoffIndex(t, db)
			seedConvergenceMigrationRows(t, db)
			before := convergenceMigrationSnapshot(t, db)
			require.NoError(t, db.Close())
			fault := &convergenceMigrationFault{mode: mode}
			wrapped := sql.OpenDB(convergenceMigrationConnector{path: path, fault: fault})
			t.Cleanup(func() { require.NoError(t, wrapped.Close()) })
			require.ErrorIs(t, ensureSQLiteSampleHandoffs(wrapped), errConvergenceMigrationInjected)
			require.True(t, fault.fired, "fault must intercept actual migration execution")
			require.Equal(t, 1, fault.drops, "failure must follow the transactional DROP")
			wantUnique := 1
			if mode == "lost-commit-response" {
				wantUnique = 0
			}
			require.Equal(t, wantUnique, convergenceTargetIndexUnique(t, wrapped))
			require.Equal(t, before, convergenceMigrationSnapshot(t, wrapped))
			require.NoError(t, ensureSQLiteSampleHandoffs(wrapped))
			require.Equal(t, 0, convergenceTargetIndexUnique(t, wrapped))
			insertConvergenceAlias(t, wrapped)
			before = convergenceMigrationSnapshot(t, wrapped)
			drops := fault.drops
			require.NoError(t, ensureSQLiteSampleHandoffs(wrapped))
			require.Equal(t, drops, fault.drops, "reentry must not rebuild the nonunique index")
			require.Equal(t, before, convergenceMigrationSnapshot(t, wrapped))
			require.NoError(t, wrapped.Close())
			db = openConvergenceMigrationDB(t, path)
			require.NoError(t, NewMigrator().Migrate(db))
			require.Equal(t, before, convergenceMigrationSnapshot(t, db))
		})
	}
}

func openConvergenceMigrationDB(t *testing.T, path string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", path)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return db
}

func setLegacyHandoffIndex(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`DROP INDEX idx_wg_handoff_target_seq;
		CREATE UNIQUE INDEX idx_wg_handoff_target_seq ON wg_delivery_sample_handoffs(target_seq) WHERE target_seq IS NOT NULL`)
	require.NoError(t, err)
}

func convergenceTargetIndexUnique(t *testing.T, db *sql.DB) int {
	t.Helper()
	var unique int
	require.NoError(t, db.QueryRow(`SELECT "unique" FROM pragma_index_list('wg_delivery_sample_handoffs') WHERE name='idx_wg_handoff_target_seq'`).Scan(&unique))
	return unique
}

func seedConvergenceMigrationRows(t *testing.T, db *sql.DB) {
	t.Helper()
	const payload = ` {"sample_id":"s-25","proof":{"epoch":"e-1"},"value":25} `
	for i, group := range []string{"A", "I", "C", "A"} {
		kind, origin, transfer := "handoff", "s-25", "transfer-"+string(rune('1'+i))
		revision := "R1"
		if i == 0 {
			kind, transfer = "native", ""
		}
		if i == 3 {
			revision = "R3"
		}
		_, err := db.Exec(`INSERT INTO wg_delivery_samples(seq,workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,payload,payload_digest,consumed,created_at,acceptance_kind,origin_workspace_id,origin_group_id,origin_revision,origin_sample_id,handoff_id)
			VALUES(?,'workspace',?,?,'physical-address','member','25','20',?,'unchanged-digest',0,'8',?,'workspace','A','R1',?,?)`, i+1, group, revision, payload, kind, origin, transfer)
		require.NoError(t, err)
	}
	for _, edge := range []struct {
		id, source, target, revision string
		sourceSeq, targetSeq         int
	}{{"transfer-2", "A", "I", "R1", 1, 2}, {"transfer-3", "A", "C", "R1", 1, 3}, {"transfer-4", "I", "A", "R3", 2, 4}} {
		_, err := db.Exec(`INSERT INTO wg_delivery_sample_handoffs(transfer_id,workspace_id,source_group_id,source_revision,source_seq,target_group_id,target_revision,target_seq,origin_seq,state,record)
			VALUES(?,'workspace',?,'R1',?,?,?,?,1,'accepted',?)`, edge.id, edge.source, edge.sourceSeq, edge.target, edge.revision, edge.targetSeq, ` {"version":1,"proof":"unchanged","id":"`+edge.id+`"} `)
		require.NoError(t, err)
	}
}

func insertConvergenceAlias(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec(`INSERT INTO wg_delivery_sample_handoffs(transfer_id,workspace_id,source_group_id,source_revision,source_seq,target_group_id,target_revision,target_seq,origin_seq,state,record)
		VALUES('alias','workspace','C','R1',3,'A','R3',4,1,'accepted',' {"version":1,"CanonicalTransferID":"transfer-4"} ')`)
	require.NoError(t, err, "independently verified edges may reference one canonical target")
}

func convergenceMigrationSnapshot(t *testing.T, db *sql.DB) []string {
	t.Helper()
	var snapshot []string
	for _, query := range []string{
		`SELECT 'table:'||name||':'||sql FROM sqlite_master WHERE type='table' AND name IN ('wg_delivery_samples','wg_delivery_sample_handoffs') ORDER BY name`,
		`SELECT 'index:'||name||':'||sql FROM sqlite_master WHERE type='index' AND name LIKE 'idx_wg_delivery_%' ORDER BY name`,
		`SELECT json_array(seq,workspace_id,group_id,group_revision,sample_id,member_key,observed_at,bucket_start,hex(payload),payload_digest,consumed,created_at,acceptance_kind,origin_workspace_id,origin_group_id,origin_revision,origin_sample_id,handoff_id) FROM wg_delivery_samples ORDER BY seq`,
		`SELECT json_array(transfer_id,workspace_id,source_group_id,source_revision,source_seq,target_group_id,target_revision,target_seq,origin_seq,state,hex(record)) FROM wg_delivery_sample_handoffs ORDER BY transfer_id`,
	} {
		rows, err := db.Query(query)
		require.NoError(t, err)
		for rows.Next() {
			var raw string
			require.NoError(t, rows.Scan(&raw))
			snapshot = append(snapshot, raw)
		}
		require.NoError(t, rows.Err())
		require.NoError(t, rows.Close())
	}
	return snapshot
}

var errConvergenceMigrationInjected = errors.New("injected convergence migration failure")

type convergenceMigrationFault struct {
	mode    string
	fired   bool
	changed bool
	drops   int
}

type convergenceMigrationConnector struct {
	path  string
	fault *convergenceMigrationFault
}

func (c convergenceMigrationConnector) Connect(context.Context) (driver.Conn, error) {
	conn, err := (&sqlite.Driver{}).Open(c.path)
	if err != nil {
		return nil, err
	}
	return &convergenceMigrationConn{Conn: conn, fault: c.fault}, nil
}

func (convergenceMigrationConnector) Driver() driver.Driver { return &sqlite.Driver{} }

type convergenceMigrationConn struct {
	driver.Conn
	fault *convergenceMigrationFault
}

func (c *convergenceMigrationConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if strings.Contains(query, "DROP INDEX") && strings.Contains(query, "idx_wg_handoff_target_seq") {
		c.fault.changed = true
		c.fault.drops++
	}
	if c.fault.mode == "ddl-failure" && !c.fault.fired && strings.Contains(query, "CREATE INDEX") && strings.Contains(query, "idx_wg_handoff_target_seq") {
		c.fault.fired = true
		return nil, errConvergenceMigrationInjected
	}
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}

func (c *convergenceMigrationConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
}

func (c *convergenceMigrationConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	tx, err := c.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return convergenceMigrationTx{Tx: tx, fault: c.fault}, nil
}

type convergenceMigrationTx struct {
	driver.Tx
	fault *convergenceMigrationFault
}

func (tx convergenceMigrationTx) Commit() error {
	if err := tx.Tx.Commit(); err != nil {
		return err
	}
	if tx.fault.mode == "lost-commit-response" && tx.fault.changed && !tx.fault.fired {
		tx.fault.fired = true
		return errConvergenceMigrationInjected
	}
	return nil
}
