package dbtarget

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/testutil/lostcommit"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

type mysqlGroupFixtureRepo struct {
	ConnectorRepository
	connector *schema.DatabaseConnector
}

func (r *mysqlGroupFixtureRepo) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return r.connector, nil
}

type mysqlGroupFixture struct {
	cfg          mysqlLongFixtureConfig
	db           *sql.DB
	dsn          string
	service      *ConnectorService
	connector    *schema.DatabaseConnector
	inspector    *ManagedTableInspector
	capability   *MySQLRecordingCapability
	receiptOwned bool
}

func newMySQLGroupFixture(t *testing.T, family string) *mysqlGroupFixture {
	t.Helper()
	cfg := loadMySQLLongFixtureConfig(t, family)
	db := openMySQLLongFixtureDB(t, cfg)
	requireMySQLFixtureFamily(t, db, family)

	pem, err := os.ReadFile(cfg.CAPath)
	require.NoError(t, err)
	pool := x509.NewCertPool()
	require.True(t, pool.AppendCertsFromPEM(pem))
	key := "group_fixture_" + randomHex(t, 6)
	require.NoError(t, mysqldriver.RegisterTLSConfig(key, &tls.Config{
		RootCAs: pool, ServerName: cfg.Host, MinVersion: tls.VersionTLS12,
	}))
	t.Cleanup(func() { mysqldriver.DeregisterTLSConfig(key) })

	dsn := mysqldriver.NewConfig()
	dsn.User, dsn.Passwd, dsn.DBName = cfg.User, cfg.Password, cfg.Database
	dsn.Net, dsn.Addr, dsn.TLSConfig = "tcp", net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port)), key
	dsn.ParseTime, dsn.Loc = true, time.UTC
	dsn.Timeout, dsn.ReadTimeout, dsn.WriteTimeout = 10*time.Second, 30*time.Second, 30*time.Second
	raw, err := json.Marshal(map[string]string{"dsn": dsn.FormatDSN()})
	require.NoError(t, err)
	connector := &schema.DatabaseConnector{
		ID: "group-live-" + randomHex(t, 6), IdentityRevision: "rev-1",
		Kind: schema.DatabaseConnectorKindMySQL, Enabled: true, ConnectionConfig: string(raw),
	}
	service := NewConnectorService(&mysqlGroupFixtureRepo{connector: connector})
	inspector := NewManagedTableInspector(service, nil)
	capability, err := inspector.ProbeMySQLRecordingAtRevision(t.Context(), connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	require.True(t, capability.Connected)
	require.Equal(t, "verified", capability.MetadataStatus)
	require.True(t, capability.TransactionalEngine)
	require.Equal(t, family, capability.ServerFamily)
	return &mysqlGroupFixture{
		cfg: cfg, db: db, dsn: dsn.FormatDSN(), service: service, connector: connector,
		inspector: inspector, capability: capability,
	}
}

func (f *mysqlGroupFixture) open(t *testing.T) *OpenedDestination {
	t.Helper()
	opened, err := f.service.OpenDestination(t.Context(), f.connector.ID, f.connector.IdentityRevision)
	require.NoError(t, err)
	t.Cleanup(func() { _ = opened.Close() })
	return opened
}

func (f *mysqlGroupFixture) create(t *testing.T) (string, *MySQLManagedRowProof) {
	t.Helper()
	table, owner := uniqueLongTestTable(t), uniqueLongTestOwnerColumn(t)
	layout, err := recordingplan.ManagedMySQLLongGroupLayout(table, owner, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily: f.capability.ServerFamily, ServerVersion: f.capability.ServerVersion,
		CapabilityRevision: f.capability.CapabilityRevision, Engine: "InnoDB",
		Collation: mysqlFamilyCollationFor(f.capability.ServerFamily),
	})
	require.NoError(t, err)
	statements, err := recordingplan.GenerateGroupSchemaDDL("mysql", layout)
	require.NoError(t, err)
	require.Len(t, statements, 3)
	receiptExists := fixtureTableExists(t, f.db, longTestReceiptTable)
	f.receiptOwned = !receiptExists
	var owned []string
	t.Cleanup(func() { dropFixtureTablesOwned(t, f.db, owned) })
	for i, statement := range statements {
		if i == 2 && receiptExists {
			continue
		}
		_, err := f.db.ExecContext(t.Context(), statement)
		require.NoError(t, err)
		switch i {
		case 2:
			owned = append(owned, longTestReceiptTable)
		case 0:
			owned = append(owned, table)
		}
	}
	group, err := f.inspector.InspectMySQLManagedTableAtRevision(t.Context(), f.capability.ConnectorID,
		f.capability.ConnectorRevision, f.cfg.Database, f.capability.CapabilityRevision, table)
	require.NoError(t, err)
	receipt, err := f.inspector.InspectMySQLManagedTableAtRevision(t.Context(), f.capability.ConnectorID,
		f.capability.ConnectorRevision, f.cfg.Database, f.capability.CapabilityRevision, longTestReceiptTable)
	require.NoError(t, err)
	return table, &MySQLManagedRowProof{
		ConnectorID: f.capability.ConnectorID, ConnectorRevision: f.capability.ConnectorRevision,
		Group: group, Receipt: receipt,
	}
}

func (f *mysqlGroupFixture) request(table, effect, digest string, proof *MySQLManagedRowProof) GroupInsertRequest {
	return f.requestWithValue(table, effect, digest, proof, 9007199254740993)
}

func (f *mysqlGroupFixture) requestWithValue(table, effect, digest string, proof *MySQLManagedRowProof, value int64) GroupInsertRequest {
	bucket := time.Date(2026, 1, 2, 3, 4, 5, 123456000, time.UTC)
	recordID := strings.Repeat("a", 64)
	return GroupInsertRequest{
		Kind: schema.DatabaseConnectorKindMySQL, SchemaName: f.cfg.Database, TableName: table,
		Strategy: GroupEffectReceipt, PayloadDigest: digest, CommittedAt: "2026-01-02T03:04:05.123456Z",
		MySQL: proof,
		Row: EncodedRow{
			RecordID: recordID, EffectKey: effect, BucketStart: bucket,
			Cells: []EncodedCell{
				{Column: "bucket_start", Value: bucket},
				{Column: "bucket_end", Value: bucket.Add(15 * time.Second)},
				{Column: "record_id", Value: recordID},
				{Column: "workspace_id", Value: "ws-live"},
				{Column: "group_id", Value: "group-live"},
				{Column: "group_revision", Value: "rev-1"},
				{Column: "device_id", Value: "device-live"},
				{Column: "point_id", Value: "point-live"},
				{Column: "tag_id", Value: "tag-live"},
				{Column: "profile_kind", Value: "instantaneous"},
				{Column: "value_type", Value: "int64"},
				{Column: "status", Value: "ok"},
				{Column: "value_int", Value: value},
				{Column: "provenance", Value: `{"source":"fixture"}`},
			},
		},
	}
}

func (f *mysqlGroupFixture) count(t *testing.T, table, where string, args ...any) int {
	t.Helper()
	query := "SELECT COUNT(*) FROM " + qualifiedTableName(schema.DatabaseConnectorKindMySQL, f.cfg.Database, table)
	if where != "" {
		query += " WHERE " + where
	}
	var count int
	require.NoError(t, f.db.QueryRowContext(t.Context(), query, args...).Scan(&count))
	return count
}

func (f *mysqlGroupFixture) insertRowDirect(t *testing.T, table string, req GroupInsertRequest) {
	t.Helper()
	columns := make([]string, 0, len(req.Row.Cells))
	args := make([]any, 0, len(req.Row.Cells))
	for _, cell := range req.Row.Cells {
		columns = append(columns, quoteIdentifier(schema.DatabaseConnectorKindMySQL, cell.Column))
		args = append(args, normalizeDBValue(cell.Value))
	}
	_, err := f.db.ExecContext(t.Context(), groupInsertStatement(req,
		qualifiedTableName(schema.DatabaseConnectorKindMySQL, f.cfg.Database, table), columns), args...)
	require.NoError(t, err)
}

func TestMySQLLongReceiptReplay(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			fixture := newMySQLGroupFixture(t, family)
			table, proof := fixture.create(t)
			opened := fixture.open(t)
			req := fixture.request(table, "effect-replay", "digest-replay", proof)
			_, err := InsertGroupRow(t.Context(), opened.DB, req)
			require.NoError(t, err)
			result, err := InsertGroupRow(t.Context(), opened.DB, req)
			require.NoError(t, err)
			require.True(t, result.AlreadyCommitted)
			require.Equal(t, 1, fixture.count(t, table, ""))
			require.Equal(t, 1, fixture.count(t, longTestReceiptTable, "effect_key = ?", req.Row.EffectKey))
		})
	}
}

func TestMySQLLongCommitAckLost(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			fixture := newMySQLGroupFixture(t, family)
			table, proof := fixture.create(t)
			req := fixture.request(table, "effect-ack", "digest-ack", proof)
			lose := &atomic.Bool{}
			lose.Store(true)
			lossy := lostcommit.Wrap(mysqldriver.MySQLDriver{}, fixture.dsn, lose)
			t.Cleanup(func() { _ = lossy.Close() })
			_, err := InsertGroupRow(t.Context(), lossy, req)
			var insertErr *GroupInsertError
			require.ErrorAs(t, err, &insertErr)
			require.True(t, insertErr.Ambiguous())
			require.ErrorIs(t, err, lostcommit.ErrResponseLost)
			lose.Store(false)
			require.NoError(t, lossy.Close())

			opened := fixture.open(t)
			result, err := InsertGroupRow(t.Context(), opened.DB, req)
			require.NoError(t, err)
			require.True(t, result.AlreadyCommitted)
			require.Equal(t, 1, fixture.count(t, table, ""))
			require.Equal(t, 1, fixture.count(t, longTestReceiptTable, "effect_key = ?", req.Row.EffectKey))
		})
	}
}

func TestMySQLLongConcurrentConflict(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			t.Run("same digest", func(t *testing.T) {
				fixture := newMySQLGroupFixture(t, family)
				table, proof := fixture.create(t)
				opened := fixture.open(t)
				req := fixture.request(table, "effect-concurrent", "digest-concurrent", proof)
				results, errs := runConcurrentMySQLInserts(t, opened, req, req)
				require.Len(t, errs, 0)
				require.Len(t, results, 2)
				require.Equal(t, 1, countAlready(results))
				require.Equal(t, 1, fixture.count(t, table, ""))
				require.Equal(t, 1, fixture.count(t, longTestReceiptTable, "effect_key = ?", req.Row.EffectKey))
			})
			t.Run("different digest", func(t *testing.T) {
				fixture := newMySQLGroupFixture(t, family)
				table, proof := fixture.create(t)
				opened := fixture.open(t)
				const firstValue = int64(9007199254740993)
				const secondValue = int64(9007199254740994)
				first := fixture.requestWithValue(table, "effect-conflict", "digest-1", proof, firstValue)
				second := fixture.requestWithValue(table, "effect-conflict", "digest-2", proof, secondValue)
				results, errs := runConcurrentMySQLInserts(t, opened, first, second)
				require.Len(t, results, 1)
				require.Len(t, errs, 1)
				require.ErrorIs(t, errs[0], ErrReceiptDigestMismatch)
				require.Equal(t, 1, fixture.count(t, table, ""))
				require.Equal(t, 1, fixture.count(t, longTestReceiptTable, "effect_key = ?", first.Row.EffectKey))

				var winningValue int64
				require.NoError(t, fixture.db.QueryRowContext(t.Context(),
					"SELECT value_int FROM "+qualifiedTableName(schema.DatabaseConnectorKindMySQL, fixture.cfg.Database, table)).Scan(&winningValue))
				var winningDigest string
				require.NoError(t, fixture.db.QueryRowContext(t.Context(),
					"SELECT payload_digest FROM "+qualifiedTableName(schema.DatabaseConnectorKindMySQL, fixture.cfg.Database, longTestReceiptTable)+" WHERE effect_key = ?", first.Row.EffectKey).Scan(&winningDigest))
				require.True(t,
					(winningValue == firstValue && winningDigest == "digest-1") ||
						(winningValue == secondValue && winningDigest == "digest-2"),
					"winning row and receipt must describe the same frozen request")
				losingValue := firstValue
				if winningValue == firstValue {
					losingValue = secondValue
				}
				require.Equal(t, 0, fixture.count(t, table, "value_int = ?", losingValue), "the losing row must never be written")
			})
		})
	}
}

func runConcurrentMySQLInserts(t *testing.T, opened *OpenedDestination, requests ...GroupInsertRequest) ([]GroupInsertResult, []error) {
	t.Helper()
	var wg sync.WaitGroup
	results := make(chan GroupInsertResult, len(requests))
	errs := make(chan error, len(requests))
	for _, req := range requests {
		wg.Go(func() {
			result, err := InsertGroupRow(t.Context(), opened.DB, req)
			if err != nil {
				errs <- err
				return
			}
			results <- result
		})
	}
	wg.Wait()
	close(results)
	close(errs)
	var gotResults []GroupInsertResult
	for result := range results {
		gotResults = append(gotResults, result)
	}
	var gotErrors []error
	for err := range errs {
		gotErrors = append(gotErrors, err)
	}
	return gotResults, gotErrors
}

func countAlready(results []GroupInsertResult) int {
	count := 0
	for _, result := range results {
		if result.AlreadyCommitted {
			count++
		}
	}
	return count
}

func TestMySQLLongAtomicRollback(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			fixture := newMySQLGroupFixture(t, family)
			table, proof := fixture.create(t)
			opened := fixture.open(t)
			req := fixture.request(table, "effect-rollback", "digest-rollback", proof)
			fixture.insertRowDirect(t, table, req)
			_, err := InsertGroupRow(t.Context(), opened.DB, req)
			require.ErrorIs(t, err, ErrExistingRowDiffers)
			require.Equal(t, 1, fixture.count(t, table, ""))
			require.Equal(t, 0, fixture.count(t, longTestReceiptTable, "effect_key = ?", req.Row.EffectKey))
		})
	}
}

func TestMySQLLongPhysicalDriftRejected(t *testing.T) {
	for _, family := range []string{mysqlFixtureFamily, mariaDBFixtureFamily} {
		t.Run(family, func(t *testing.T) {
			for name, mutate := range map[string]func(*testing.T, *mysqlGroupFixture, string){
				"owner": func(t *testing.T, f *mysqlGroupFixture, table string) {
					owner := mysqlGroupOwnerColumn(t, f, table)
					_, err := f.db.ExecContext(t.Context(), "ALTER TABLE "+qualifiedTableName(schema.DatabaseConnectorKindMySQL, f.cfg.Database, table)+" MODIFY COLUMN "+quoteIdentifier(schema.DatabaseConnectorKindMySQL, owner)+" VARCHAR(32) NOT NULL DEFAULT 'gateway'")
					require.NoError(t, err)
				},
				"engine": func(t *testing.T, f *mysqlGroupFixture, table string) {
					require.True(t, f.receiptOwned, "engine drift requires a receipt created by this fixture")
					_, err := f.db.ExecContext(t.Context(), "ALTER TABLE "+qualifiedTableName(schema.DatabaseConnectorKindMySQL, f.cfg.Database, longTestReceiptTable)+" ENGINE=MyISAM")
					require.NoError(t, err)
				},
				"physical": func(t *testing.T, f *mysqlGroupFixture, table string) {
					_, err := f.db.ExecContext(t.Context(), "ALTER TABLE "+qualifiedTableName(schema.DatabaseConnectorKindMySQL, f.cfg.Database, table)+" ADD COLUMN `_foreign` INT")
					require.NoError(t, err)
				},
			} {
				t.Run(name, func(t *testing.T) {
					fixture := newMySQLGroupFixture(t, family)
					table, proof := fixture.create(t)
					opened := fixture.open(t)
					mutate(t, fixture, table)
					_, err := InsertGroupRow(t.Context(), opened.DB, fixture.request(table, "effect-drift-"+name, "digest-drift-"+name, proof))
					require.ErrorIs(t, err, ErrGroupInsertUnsupported)
				})
			}
		})
	}
}

func mysqlGroupOwnerColumn(t *testing.T, f *mysqlGroupFixture, table string) string {
	t.Helper()
	var owner string
	require.NoError(t, f.db.QueryRowContext(t.Context(), `
		SELECT COLUMN_NAME FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME LIKE '_gw_owner_%' LIMIT 1`, f.cfg.Database, table).Scan(&owner))
	return owner
}
