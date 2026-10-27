package groupdelivery

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/snapshot"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

const (
	mysqlFixtureEnv      = "GW_CIRCUIT_MYSQL_FIXTURE_FILE"
	mariaDBFixtureEnv    = "GW_CIRCUIT_MARIADB_FIXTURE_FILE"
	mysqlFixtureFamily   = "mysql"
	mariaDBFixtureFamily = "mariadb"
	longTestTablePrefix  = "gw_ltest_"
	longTestReceiptTable = dbtarget.EffectReceiptTable
	fixtureHost          = "127.0.0.1"
	fixtureDatabase      = "gw_circuit_test"
)

// mysqlSenderFixtureConfig mirrors the private 0600 JSON fixture configuration.
// Secrets, passwords, and connection paths must never be logged or echoed.
type mysqlSenderFixtureConfig struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	Database string `json:"Database"`
	User     string `json:"User"`
	Password string `json:"Password"`
	CAPath   string `json:"CAPath"`
}

func mysqlSenderFixtureEnvFor(family string) string {
	if family == mariaDBFixtureFamily {
		return mariaDBFixtureEnv
	}
	return mysqlFixtureEnv
}

// loadMySQLSenderFixtureConfig reads and validates the private 0600 fixture packet,
// skipping (NOT RUN) if the environment variable is not configured.
func loadMySQLSenderFixtureConfig(t *testing.T, family string) mysqlSenderFixtureConfig {
	t.Helper()
	switch family {
	case mysqlFixtureFamily, mariaDBFixtureFamily:
	default:
		t.Fatalf("unsupported fixture family %q", family)
	}

	path := os.Getenv(mysqlSenderFixtureEnvFor(family))
	if path == "" {
		t.Skipf("NOT RUN: %s is not set; the disposable %s fixture is unavailable", mysqlSenderFixtureEnvFor(family), family)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat fixture file: %v", err)
	}
	if !info.Mode().IsRegular() {
		t.Fatalf("fixture file %s is not a regular file", path)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("fixture file %s must have 0600 permissions", path)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture file: %v", err)
	}

	var cfg mysqlSenderFixtureConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("decode fixture file: %v", err)
	}
	if cfg.Host != fixtureHost {
		t.Fatalf("fixture Host must be %s", fixtureHost)
	}
	if cfg.Database != fixtureDatabase {
		t.Fatalf("fixture Database must be %s", fixtureDatabase)
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		t.Fatalf("fixture Port out of range")
	}
	if cfg.User == "" || cfg.Password == "" || cfg.CAPath == "" {
		t.Fatalf("fixture User, Password, and CAPath are required")
	}
	if _, err := os.Stat(cfg.CAPath); err != nil {
		t.Fatalf("fixture CA not found: %v", err)
	}
	return cfg
}

type mysqlSenderConnectorRepo struct {
	dbtarget.ConnectorRepository
	connector *schema.DatabaseConnector
}

func (r *mysqlSenderConnectorRepo) GetByID(context.Context, string) (*schema.DatabaseConnector, error) {
	return r.connector, nil
}

type mysqlSenderResolver struct {
	fixture *mysqlSenderFixture
	proof   *dbtarget.MySQLManagedRowProof
}

func (r mysqlSenderResolver) Resolve(ctx context.Context, item OutboxItem) (Target, error) {
	opened, err := r.fixture.service.OpenDestination(ctx, item.Destination.ConnectorID, item.Destination.ConnectorRevision)
	if err != nil {
		return Target{}, err
	}
	return Target{DB: opened.DB, Kind: opened.Kind, Close: opened.Close, MySQL: r.proof}, nil
}

type mysqlSenderFixture struct {
	cfg        mysqlSenderFixtureConfig
	db         *sql.DB
	dsn        string
	service    *dbtarget.ConnectorService
	connector  *schema.DatabaseConnector
	inspector  *dbtarget.ManagedTableInspector
	capability *dbtarget.MySQLRecordingCapability
}

// newMySQLSenderFixture initializes the verified TLS connection and probing capability.
func newMySQLSenderFixture(t *testing.T, family string) *mysqlSenderFixture {
	t.Helper()
	cfg := loadMySQLSenderFixtureConfig(t, family)

	pem, err := os.ReadFile(cfg.CAPath)
	require.NoError(t, err)
	pool := x509.NewCertPool()
	require.True(t, pool.AppendCertsFromPEM(pem))

	entropy := make([]byte, 6)
	_, err = rand.Read(entropy)
	require.NoError(t, err)
	tlsKey := longTestTablePrefix + "tls_" + hex.EncodeToString(entropy)
	require.NoError(t, mysqldriver.RegisterTLSConfig(tlsKey, &tls.Config{
		RootCAs:    pool,
		ServerName: fixtureHost,
		MinVersion: tls.VersionTLS12,
	}))
	t.Cleanup(func() { mysqldriver.DeregisterTLSConfig(tlsKey) })

	dsn := mysqldriver.NewConfig()
	dsn.Net = "tcp"
	dsn.Addr = net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	dsn.User = cfg.User
	dsn.Passwd = cfg.Password
	dsn.DBName = cfg.Database
	dsn.TLSConfig = tlsKey
	dsn.ParseTime = true
	dsn.Loc = time.UTC
	dsn.Timeout = 10 * time.Second
	dsn.ReadTimeout = 30 * time.Second
	dsn.WriteTimeout = 30 * time.Second

	dsnStr := dsn.FormatDSN()
	adminDB, err := sql.Open("mysql", dsnStr)
	require.NoError(t, err)
	t.Cleanup(func() { _ = adminDB.Close() })
	require.NoError(t, adminDB.PingContext(t.Context()))

	rawConfig, err := json.Marshal(map[string]string{"dsn": dsnStr})
	require.NoError(t, err)

	connectorID := "sender-live-" + hex.EncodeToString(entropy)
	connector := &schema.DatabaseConnector{
		ID:               connectorID,
		IdentityRevision: "crev-1",
		Kind:             schema.DatabaseConnectorKindMySQL,
		Enabled:          true,
		ConnectionConfig: string(rawConfig),
	}

	service := dbtarget.NewConnectorService(&mysqlSenderConnectorRepo{connector: connector})
	inspector := dbtarget.NewManagedTableInspector(service, nil)
	capability, err := inspector.ProbeMySQLRecordingAtRevision(t.Context(), connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	require.True(t, capability.Connected)
	require.Equal(t, "verified", capability.MetadataStatus)
	require.True(t, capability.TransactionalEngine)
	require.Equal(t, family, capability.ServerFamily)

	return &mysqlSenderFixture{
		cfg:        cfg,
		db:         adminDB,
		dsn:        dsnStr,
		service:    service,
		connector:  connector,
		inspector:  inspector,
		capability: capability,
	}
}

// createTables generates canonical long layout DDL, creates the group table,
// only creates the receipt table if absent, and tracks exact cleanup so preexisting
// tables are never dropped.
func (f *mysqlSenderFixture) createTables(t *testing.T) (string, *dbtarget.MySQLManagedRowProof) {
	t.Helper()
	entropy := make([]byte, 12)
	_, err := rand.Read(entropy)
	require.NoError(t, err)
	table := longTestTablePrefix + hex.EncodeToString(entropy)
	ownerMarker := "_gw_owner_" + hex.EncodeToString(entropy)

	collation := "utf8mb4_0900_bin"
	if f.capability.ServerFamily == mariaDBFixtureFamily {
		collation = "utf8mb4_nopad_bin"
	}

	layout, err := recordingplan.ManagedMySQLLongGroupLayout(table, ownerMarker, recordingplan.MySQLGroupSchemaOptions{
		ServerFamily:       f.capability.ServerFamily,
		ServerVersion:      f.capability.ServerVersion,
		CapabilityRevision: f.capability.CapabilityRevision,
		Engine:             "InnoDB",
		Collation:          collation,
	})
	require.NoError(t, err)

	statements, err := recordingplan.GenerateGroupSchemaDDL("mysql", layout)
	require.NoError(t, err)
	require.Len(t, statements, 3)

	var receiptCount int
	err = f.db.QueryRowContext(t.Context(), `
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = ? AND table_name = ?
	`, f.cfg.Database, longTestReceiptTable).Scan(&receiptCount)
	require.NoError(t, err)
	receiptExists := receiptCount > 0

	var owned []string
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for _, tbl := range owned {
			if _, err := f.db.ExecContext(cleanupCtx, "DROP TABLE IF EXISTS `"+f.cfg.Database+"`.`"+tbl+"`"); err != nil {
				t.Errorf("cleanup drop table %s failed: %v", tbl, err)
			}
		}
	})

	for i, stmt := range statements {
		if i == 2 && receiptExists {
			continue
		}
		_, err := f.db.ExecContext(t.Context(), stmt)
		require.NoError(t, err)
		switch i {
		case 0:
			owned = append(owned, table)
		case 2:
			owned = append(owned, longTestReceiptTable)
		}
	}

	groupInspection, err := f.inspector.InspectMySQLManagedTableAtRevision(t.Context(),
		f.capability.ConnectorID, f.capability.ConnectorRevision, f.cfg.Database, f.capability.CapabilityRevision, table)
	require.NoError(t, err)
	receiptInspection, err := f.inspector.InspectMySQLManagedTableAtRevision(t.Context(),
		f.capability.ConnectorID, f.capability.ConnectorRevision, f.cfg.Database, f.capability.CapabilityRevision, longTestReceiptTable)
	require.NoError(t, err)

	proof := &dbtarget.MySQLManagedRowProof{
		ConnectorID:       f.capability.ConnectorID,
		ConnectorRevision: f.capability.ConnectorRevision,
		Group:             groupInspection,
		Receipt:           receiptInspection,
	}
	return table, proof
}

// enqueueCanonicalLong commits a single outbox item into the local SQLite store
// with all required canonical long columns, setting value_int to the specified value.
func (f *mysqlSenderFixture) enqueueCanonicalLong(t *testing.T, store *Store, table, effectKey string, valueInt int64) {
	t.Helper()
	bucketTime := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	recordID := strings.Repeat("e", 64)

	row := dbtarget.EncodedRow{
		RecordID:    recordID,
		EffectKey:   effectKey,
		BucketStart: bucketTime,
		Cells: []dbtarget.EncodedCell{
			{Column: "bucket_start", Value: bucketTime},
			{Column: "bucket_end", Value: bucketTime.Add(15 * time.Second)},
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
			{Column: "value_int", Value: valueInt},
			{Column: "provenance", Value: `{"source":"sender-mysql-test"}`},
		},
	}

	destination := FrozenDestination{
		ConnectorID:       f.connector.ID,
		ConnectorRevision: f.connector.IdentityRevision,
		Database:          f.cfg.Database,
		TableSchema:       f.cfg.Database,
		TableName:         table,
		DedupeCapability:  DedupeReceipt,
	}
	var err error
	destination.Scope, err = snapshot.DestinationScope(destination.ConnectorID, destination.ConnectorRevision,
		destination.Database, destination.TableSchema, destination.TableName)
	require.NoError(t, err)

	bucket := ClosedBucket{
		Outcome: snapshot.Outcome{
			Kind:        snapshot.OutcomeRow,
			RecordID:    recordID,
			EffectKey:   effectKey,
			BucketStart: bucketTime,
			BucketEnd:   bucketTime.Add(15 * time.Second),
			Members: []snapshot.MemberResult{{
				MemberKey: "power",
				Status:    snapshot.MemberOK,
				Sample: &snapshot.Sample{
					SampleID:   "s-1",
					MemberKey:  "power",
					ObservedAt: bucketTime,
					Quality:    schema.QualityGood,
					Value:      measurement.NewInt64(valueInt),
				},
			}},
		},
		Row: &row,
	}

	closure := Closure{
		Destination: destination,
		NextClose:   bucketTime.Add(15 * time.Second),
		Buckets:     []ClosedBucket{bucket},
	}
	groupKey := GroupKey{
		WorkspaceID:   "ws-live",
		GroupID:       "group-live",
		GroupRevision: "rev-1",
	}

	require.NoError(t, store.CommitClosure(t.Context(), groupKey, closure))
}

func (f *mysqlSenderFixture) countRows(t *testing.T, table, where string, args ...any) int {
	t.Helper()
	query := "SELECT COUNT(*) FROM `" + f.cfg.Database + "`.`" + table + "`"
	if where != "" {
		query += " WHERE " + where
	}
	var count int
	require.NoError(t, f.db.QueryRowContext(t.Context(), query, args...).Scan(&count))
	return count
}

func (f *mysqlSenderFixture) readValueInt(t *testing.T, table, recordID string) int64 {
	t.Helper()
	var val int64
	query := "SELECT value_int FROM `" + f.cfg.Database + "`.`" + table + "` WHERE record_id = ?"
	require.NoError(t, f.db.QueryRowContext(t.Context(), query, recordID).Scan(&val))
	return val
}
