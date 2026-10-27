package dbtarget

// Loader and TLS connection helper for the disposable gw-circuit
// MySQL/MariaDB fixtures. The fixture file is a 0600 JSON outside the
// repository; its contents, the password and any DSN are never logged.

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

const (
	mysqlFixtureEnv      = "GW_CIRCUIT_MYSQL_FIXTURE_FILE"
	mariaDBFixtureEnv    = "GW_CIRCUIT_MARIADB_FIXTURE_FILE"
	mysqlFixtureFamily   = "mysql"
	mariaDBFixtureFamily = "mariadb"
	longTestTablePrefix  = "gw_ltest_"
	longTestReceiptTable = "gw_effect_receipts"
	// recordingplan managedOwnerDefault, pinned here as a fixed expectation.
	longTestOwnerDefault = "gateway"
	fixtureHost          = "127.0.0.1"
	fixtureDatabase      = "gw_circuit_test"
)

var ownedTableNameRE = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

// mysqlLongFixtureConfig mirrors the private fixture file keys. Every field is
// required: a zero value means a malformed fixture packet, never a default.
type mysqlLongFixtureConfig struct {
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
	Database string `json:"Database"`
	User     string `json:"User"`
	Password string `json:"Password"`
	CAPath   string `json:"CAPath"`
}

func mysqlFixtureEnvFor(family string) string {
	if family == mariaDBFixtureFamily {
		return mariaDBFixtureEnv
	}
	return mysqlFixtureEnv
}

// loadMySQLLongFixtureConfig skips (NOT RUN) when the packet env is absent and
// otherwise enforces the fixture contract: loopback host, the dedicated
// disposable database and a 0600 private file.
func loadMySQLLongFixtureConfig(t *testing.T, family string) mysqlLongFixtureConfig {
	t.Helper()
	switch family {
	case mysqlFixtureFamily, mariaDBFixtureFamily:
	default:
		t.Fatalf("unsupported fixture family %q", family)
	}
	path := os.Getenv(mysqlFixtureEnvFor(family))
	if path == "" {
		t.Skipf("NOT RUN: %s is not set; the disposable %s fixture is unavailable", mysqlFixtureEnvFor(family), family)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat fixture file: %v", err)
	}
	if !info.Mode().IsRegular() {
		t.Fatalf("fixture file %s is not a regular file", path)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("fixture file %s must not be readable by group or other", path)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture file: %v", err)
	}
	var cfg mysqlLongFixtureConfig
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
		t.Fatalf("fixture User, Password and CAPath are all required")
	}
	if _, err := os.Stat(cfg.CAPath); err != nil {
		t.Fatalf("fixture CA: %v", err)
	}
	return cfg
}

// openMySQLLongFixtureDB opens the fixture over certificate-verified TLS only:
// RootCAs from the packet CA, ServerName pinned to the loopback host, TLS
// 1.2 minimum, no skip-verify and no plaintext fallback. The per-test TLS
// config is registered under a unique name and deregistered on cleanup.
func openMySQLLongFixtureDB(t *testing.T, cfg mysqlLongFixtureConfig) *sql.DB {
	t.Helper()
	pem, err := os.ReadFile(cfg.CAPath)
	if err != nil {
		t.Fatalf("read fixture CA: %v", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pem) {
		t.Fatalf("fixture CA PEM is unusable")
	}
	entropy := make([]byte, 6)
	if _, err := rand.Read(entropy); err != nil {
		t.Fatalf("tls config name entropy: %v", err)
	}
	tlsName := longTestTablePrefix + "tls_" + fmt.Sprintf("%x", entropy)
	if err := mysql.RegisterTLSConfig(tlsName, &tls.Config{
		RootCAs:    pool,
		ServerName: fixtureHost,
		MinVersion: tls.VersionTLS12,
	}); err != nil {
		t.Fatalf("register tls config: %v", err)
	}
	t.Cleanup(func() { mysql.DeregisterTLSConfig(tlsName) })

	dsn := mysql.NewConfig()
	dsn.Net = "tcp"
	dsn.Addr = net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))
	dsn.User = cfg.User
	dsn.Passwd = cfg.Password
	dsn.DBName = cfg.Database
	dsn.TLSConfig = tlsName
	dsn.ParseTime = true
	dsn.Loc = time.UTC
	dsn.Timeout = 10 * time.Second
	dsn.ReadTimeout = 30 * time.Second
	dsn.WriteTimeout = 30 * time.Second
	db, err := sql.Open("mysql", dsn.FormatDSN()) // DSN embeds the secret; never logged.
	if err != nil {
		t.Fatalf("open %s fixture: %v", cfg.Database, err)
	}
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping fixture over verified TLS: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

// requireMySQLFixtureFamily cross-checks the packet-declared family against
// the server identity reported by VERSION() on the verified connection.
func requireMySQLFixtureFamily(t *testing.T, db *sql.DB, family string) string {
	t.Helper()
	var version string
	if err := db.QueryRowContext(t.Context(), "SELECT VERSION()").Scan(&version); err != nil {
		t.Fatalf("select version(): %v", err)
	}
	isMariaDB := strings.Contains(version, "MariaDB")
	if family == mariaDBFixtureFamily && (!isMariaDB || !strings.HasPrefix(strings.TrimPrefix(version, "5.5.5-"), "11.4.")) {
		t.Fatalf("fixture packet says %s but VERSION() %q is not MariaDB 11.4", family, version)
	}
	if family == mysqlFixtureFamily && (isMariaDB || !strings.HasPrefix(version, "8.4.")) {
		t.Fatalf("fixture packet says %s but VERSION() %q is not MySQL 8.4", family, version)
	}
	return version
}

func randomHex(t *testing.T, bytes int) string {
	t.Helper()
	entropy := make([]byte, bytes)
	if _, err := rand.Read(entropy); err != nil {
		t.Fatalf("random hex: %v", err)
	}
	return fmt.Sprintf("%x", entropy)
}

func uniqueLongTestTable(t *testing.T) string {
	t.Helper()
	return longTestTablePrefix + randomHex(t, 8)
}

// uniqueLongTestOwnerColumn follows the production canonical owner-column
// shape: an underscore-prefixed _gw_owner_ name with a unique suffix.
func uniqueLongTestOwnerColumn(t *testing.T) string {
	t.Helper()
	return "_gw_owner_" + randomHex(t, 12)
}

// fixtureTableExists reports whether the named base table already exists in
// the fixture database, so pre-existing objects are never claimed or dropped.
func fixtureTableExists(t *testing.T, db *sql.DB, name string) bool {
	t.Helper()
	var one int
	err := db.QueryRowContext(t.Context(),
		"SELECT 1 FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", name).
		Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	}
	require.NoError(t, err)
	return one == 1
}

// dropFixtureTablesOwned drops only the tables this test created (verified
// by successful CREATE statements). Pre-existing or foreign tables, including a
// pre-existing receipt table, are never touched.
func dropFixtureTablesOwned(t *testing.T, db *sql.DB, created []string) {
	t.Helper()
	// t.Context() is canceled before Cleanup-registered functions run, so the
	// drops carry their own lifetime.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	drops := slices.Clone(created)
	slices.Reverse(drops)
	for _, name := range drops {
		if !ownedTableNameRE.MatchString(name) ||
			(name != longTestReceiptTable && !strings.HasPrefix(name, longTestTablePrefix)) {
			t.Errorf("refusing to drop unexpected table %q", name)
			continue
		}
		if _, err := db.ExecContext(ctx, "DROP TABLE `"+name+"`"); err != nil {
			t.Errorf("drop created table %q: %v", name, err)
		}
	}
}
