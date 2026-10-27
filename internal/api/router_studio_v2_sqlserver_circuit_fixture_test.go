package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/workspace"

	"github.com/stretchr/testify/require"
)

const sqlServerCircuitFixtureEnv = "GW_SQLSERVER_TEST_FIXTURE"

// ErrSQLServerCircuitFixtureUnavailable marks a missing or unsafe real SQL
// fixture. Callers must treat it as NOT_RUN, never as a passing engine case.
var ErrSQLServerCircuitFixtureUnavailable = errors.New("SQL Server circuit fixture unavailable")

// SQLServerCircuitFixtureManifest is the opt-in real-fixture manifest. The
// manifest contains no credentials; connection_config_file points at an
// existing saved SQL Server connection-config JSON document.
type SQLServerCircuitFixtureManifest struct {
	SchemaVersion        int    `json:"schema_version"`
	Disposable           bool   `json:"disposable"`
	NativeAMD64          bool   `json:"native_amd64"`
	HostOS               string `json:"host_os"`
	HostArch             string `json:"host_arch"`
	Database             string `json:"database"`
	Schema               string `json:"schema"`
	ServerVersionPrefix  string `json:"server_version_prefix"`
	ServerEdition        string `json:"server_edition,omitempty"`
	EngineEdition        int64  `json:"engine_edition"`
	CompatibilityLevel   int64  `json:"compatibility_level"`
	ConnectionConfigFile string `json:"connection_config_file"`
	ImageDigest          string `json:"image_digest,omitempty"`
	Provenance           string `json:"provenance,omitempty"`
}

// SQLServerCircuitFixture is the validated manifest plus the existing saved
// connection-config payload. It intentionally exposes no credential fields.
type SQLServerCircuitFixture struct {
	Manifest       SQLServerCircuitFixtureManifest
	connectionJSON string
}

func (f *SQLServerCircuitFixture) ConnectionConfigJSON() string {
	if f == nil {
		return ""
	}
	return f.connectionJSON
}

// SQLServerCircuitRuntimeOptions is the small consumer contract for GLM/root:
// validated manifest scope plus the existing saved connection-config JSON.
type SQLServerCircuitRuntimeOptions struct {
	Database             string
	Schema               string
	ConnectionConfigJSON string
	Disposable           bool
	NativeAMD64          bool
	ServerVersionPrefix  string
	EngineEdition        int64
	CompatibilityLevel   int64
}

// SQLServerCircuitFixtureOptions are the small reusable setup choices for
// consumer tests. The fixture still owns the saved connector, source, and
// confirmation flow; consumers may select the storage mode and tag type.
type SQLServerCircuitFixtureOptions struct {
	StorageMode workspace.WriteGroupStorageMode
	TagType     string
}

func (f *SQLServerCircuitFixture) RuntimeOptions() SQLServerCircuitRuntimeOptions {
	if f == nil {
		return SQLServerCircuitRuntimeOptions{}
	}
	return SQLServerCircuitRuntimeOptions{
		Database: f.Manifest.Database, Schema: f.Manifest.Schema, ConnectionConfigJSON: f.connectionJSON,
		Disposable: f.Manifest.Disposable, NativeAMD64: f.Manifest.NativeAMD64,
		ServerVersionPrefix: f.Manifest.ServerVersionPrefix, EngineEdition: f.Manifest.EngineEdition,
		CompatibilityLevel: f.Manifest.CompatibilityLevel,
	}
}

func (f *SQLServerCircuitFixture) ValidateCapability(capability *dbtarget.SQLServerRecordingCapability) error {
	if f == nil || capability == nil {
		return fmt.Errorf("%w: capability evidence is missing", ErrSQLServerCircuitFixtureUnavailable)
	}
	if !capability.Connected || capability.MetadataStatus != "verified" {
		return fmt.Errorf("%w: SQL Server metadata is not verified", ErrSQLServerCircuitFixtureUnavailable)
	}
	if capability.Database != f.Manifest.Database || capability.Schema != f.Manifest.Schema {
		return fmt.Errorf("%w: SQL Server database/schema scope differs from manifest", ErrSQLServerCircuitFixtureUnavailable)
	}
	if !strings.HasPrefix(capability.ServerVersion, f.Manifest.ServerVersionPrefix) {
		return fmt.Errorf("%w: SQL Server version is outside the declared baseline", ErrSQLServerCircuitFixtureUnavailable)
	}
	if f.Manifest.ServerEdition != "" && capability.ServerEdition != f.Manifest.ServerEdition {
		return fmt.Errorf("%w: SQL Server edition differs from manifest", ErrSQLServerCircuitFixtureUnavailable)
	}
	if capability.EngineEdition != f.Manifest.EngineEdition {
		return fmt.Errorf("%w: SQL Server engine edition differs from manifest", ErrSQLServerCircuitFixtureUnavailable)
	}
	if capability.CompatibilityLevel != f.Manifest.CompatibilityLevel {
		return fmt.Errorf("%w: SQL Server compatibility level differs from manifest", ErrSQLServerCircuitFixtureUnavailable)
	}
	permissions := capability.Permissions
	if !permissions.Verified || !permissions.DatabaseViewDefinition || !permissions.SchemaViewDefinition ||
		!permissions.CreateTable || !permissions.AlterSchema || !permissions.SelectSchema ||
		!permissions.InsertSchema || !permissions.DeleteSchema {
		return fmt.Errorf("%w: required SQL Server catalog/recording permissions are missing", ErrSQLServerCircuitFixtureUnavailable)
	}
	return nil
}

// LoadSQLServerCircuitFixture reads the opt-in manifest and validates local
// file/scope guards. It does not connect to SQL Server or reveal credentials.
func LoadSQLServerCircuitFixture() (*SQLServerCircuitFixture, error) {
	path := strings.TrimSpace(os.Getenv(sqlServerCircuitFixtureEnv))
	if path == "" {
		return nil, fmt.Errorf("%w: %s is not set", ErrSQLServerCircuitFixtureUnavailable, sqlServerCircuitFixtureEnv)
	}
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%w: manifest path is not a regular file", ErrSQLServerCircuitFixtureUnavailable)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: manifest cannot be read", ErrSQLServerCircuitFixtureUnavailable)
	}
	var manifest SQLServerCircuitFixtureManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, fmt.Errorf("%w: manifest is not valid JSON", ErrSQLServerCircuitFixtureUnavailable)
	}
	if err := validateSQLServerCircuitManifest(manifest); err != nil {
		return nil, err
	}
	configPath := strings.TrimSpace(manifest.ConnectionConfigFile)
	configInfo, err := os.Stat(configPath)
	if err != nil || !configInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("%w: connection_config_file is not a regular file", ErrSQLServerCircuitFixtureUnavailable)
	}
	if configInfo.Mode().Perm()&0o077 != 0 {
		return nil, fmt.Errorf("%w: connection_config_file permissions are too broad", ErrSQLServerCircuitFixtureUnavailable)
	}
	configRaw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("%w: connection_config_file cannot be read", ErrSQLServerCircuitFixtureUnavailable)
	}
	var config map[string]any
	if err := json.Unmarshal(configRaw, &config); err != nil {
		return nil, fmt.Errorf("%w: connection_config_file is not a JSON object", ErrSQLServerCircuitFixtureUnavailable)
	}
	if config["database"] != manifest.Database || config["schema"] != manifest.Schema {
		return nil, fmt.Errorf("%w: connection_config_file scope differs from manifest", ErrSQLServerCircuitFixtureUnavailable)
	}
	password, passwordOK := config["password"].(string)
	if strings.TrimSpace(stringConfigValue(config, "host")) == "" || strings.TrimSpace(stringConfigValue(config, "user")) == "" ||
		!passwordOK || password == "" {
		return nil, fmt.Errorf("%w: connection_config_file lacks required SQL authentication fields", ErrSQLServerCircuitFixtureUnavailable)
	}
	return &SQLServerCircuitFixture{Manifest: manifest, connectionJSON: string(configRaw)}, nil
}

func validateSQLServerCircuitManifest(manifest SQLServerCircuitFixtureManifest) error {
	version16 := strings.HasPrefix(manifest.ServerVersionPrefix, "16.")
	version17 := strings.HasPrefix(manifest.ServerVersionPrefix, "17.")
	switch {
	case manifest.SchemaVersion != 1:
		return fmt.Errorf("%w: manifest schema_version must be 1", ErrSQLServerCircuitFixtureUnavailable)
	case !manifest.Disposable:
		return fmt.Errorf("%w: manifest must attest a disposable scope", ErrSQLServerCircuitFixtureUnavailable)
	case manifest.HostOS != "linux" && manifest.HostOS != "windows" && manifest.HostOS != "darwin":
		return fmt.Errorf("%w: manifest host_os must be linux, windows, or darwin", ErrSQLServerCircuitFixtureUnavailable)
	case manifest.HostArch != "amd64" && manifest.HostArch != "arm64":
		return fmt.Errorf("%w: manifest host_arch must be amd64 or arm64", ErrSQLServerCircuitFixtureUnavailable)
	case manifest.NativeAMD64 && manifest.HostArch != "amd64":
		return fmt.Errorf("%w: native SQL Server host must attest amd64 architecture", ErrSQLServerCircuitFixtureUnavailable)
	case !manifest.NativeAMD64 && manifest.HostArch != "arm64":
		return fmt.Errorf("%w: non-native SQL Server host must attest arm64 emulation architecture", ErrSQLServerCircuitFixtureUnavailable)
	case strings.TrimSpace(manifest.Database) == "" || strings.TrimSpace(manifest.Schema) == "":
		return fmt.Errorf("%w: manifest database and schema are required", ErrSQLServerCircuitFixtureUnavailable)
	case !version16 && !version17:
		return fmt.Errorf("%w: manifest server baseline must be SQL Server 2022 (16.*) or 2025 (17.*)", ErrSQLServerCircuitFixtureUnavailable)
	case manifest.EngineEdition != 2 && manifest.EngineEdition != 3 && manifest.EngineEdition != 4:
		return fmt.Errorf("%w: manifest engine_edition must be 2, 3, or 4", ErrSQLServerCircuitFixtureUnavailable)
	case (version16 && manifest.CompatibilityLevel != 160) || (version17 && manifest.CompatibilityLevel != 170):
		return fmt.Errorf("%w: manifest compatibility_level must pair 160 with 16.* or 170 with 17.*", ErrSQLServerCircuitFixtureUnavailable)
	case strings.TrimSpace(manifest.ConnectionConfigFile) == "":
		return fmt.Errorf("%w: manifest connection_config_file is required", ErrSQLServerCircuitFixtureUnavailable)
	case !filepath.IsAbs(manifest.ConnectionConfigFile):
		return fmt.Errorf("%w: manifest connection_config_file must be absolute", ErrSQLServerCircuitFixtureUnavailable)
	}
	return nil
}

// SQLServerCircuitFixtureOrSkip enforces the real-fixture opt-in. A missing or
// unsafe fixture is a truthful NOT_RUN, not a passing engine case.
func SQLServerCircuitFixtureOrSkip(t *testing.T) *SQLServerCircuitFixture {
	t.Helper()
	fixture, err := LoadSQLServerCircuitFixture()
	if errors.Is(err, ErrSQLServerCircuitFixtureUnavailable) {
		t.Skipf("NOT_RUN: %v", err)
	}
	require.NoError(t, err)
	return fixture
}

func stringConfigValue(config map[string]any, key string) string {
	value, _ := config[key].(string)
	return value
}

func TestSQLServerCircuitFixtureManifestGuards(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "connection.json")
	manifestPath := filepath.Join(dir, "manifest.json")
	config := map[string]any{"host": "127.0.0.1", "user": "fixture", "password": " ", "database": "recordings", "schema": "dbo"}
	rawConfig, err := json.Marshal(config)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(configPath, rawConfig, 0o600))
	manifest := SQLServerCircuitFixtureManifest{
		SchemaVersion: 1, Disposable: true, NativeAMD64: true, HostOS: "windows", HostArch: "amd64",
		Database: "recordings", Schema: "dbo", ServerVersionPrefix: "16.", EngineEdition: 3,
		CompatibilityLevel: 160, ConnectionConfigFile: configPath, ImageDigest: "sha256:fixture",
	}
	rawManifest, err := json.Marshal(manifest)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(manifestPath, rawManifest, 0o600))
	t.Setenv(sqlServerCircuitFixtureEnv, manifestPath)
	fixture, err := LoadSQLServerCircuitFixture()
	require.NoError(t, err)
	require.Equal(t, "recordings", fixture.Manifest.Database)
	require.NotEmpty(t, fixture.ConnectionConfigJSON())

	manifest.HostOS = "linux"
	rawManifest, err = json.Marshal(manifest)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(manifestPath, rawManifest, 0o600))
	_, err = LoadSQLServerCircuitFixture()
	require.NoError(t, err, "linux/amd64 is also a supported native SQL Server fixture host")

	manifest.HostOS = "darwin"
	manifest.HostArch = "arm64"
	manifest.NativeAMD64 = false
	manifest.ServerVersionPrefix = "17."
	manifest.CompatibilityLevel = 170
	rawManifest, err = json.Marshal(manifest)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(manifestPath, rawManifest, 0o600))
	_, err = LoadSQLServerCircuitFixture()
	require.NoError(t, err, "darwin/arm64 emulation fixture is supported for SQL Server 2025")

	require.NoError(t, fixture.ValidateCapability(&dbtarget.SQLServerRecordingCapability{
		Database: "recordings", Schema: "dbo", ServerVersion: "16.0.4180.1", EngineEdition: 3,
		CompatibilityLevel: 160, Connected: true, MetadataStatus: "verified",
		Permissions: dbtarget.SQLServerRecordingPermissions{
			Verified: true, DatabaseViewDefinition: true, SchemaViewDefinition: true, CreateTable: true,
			AlterSchema: true, SelectSchema: true, InsertSchema: true, DeleteSchema: true,
		},
	}))
	require.ErrorIs(t, fixture.ValidateCapability(&dbtarget.SQLServerRecordingCapability{
		Database: "foreign", Schema: "dbo", ServerVersion: "16.0.4180.1", EngineEdition: 3,
		CompatibilityLevel: 160, Connected: true, MetadataStatus: "verified",
	}), ErrSQLServerCircuitFixtureUnavailable)

	manifest.Disposable = false
	rawManifest, err = json.Marshal(manifest)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(manifestPath, rawManifest, 0o600))
	_, err = LoadSQLServerCircuitFixture()
	require.ErrorIs(t, err, ErrSQLServerCircuitFixtureUnavailable)

	require.NoError(t, os.Chmod(configPath, 0o644))
	manifest.Disposable = true
	manifest.HostOS = "windows"
	rawManifest, err = json.Marshal(manifest)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(manifestPath, rawManifest, 0o600))
	_, err = LoadSQLServerCircuitFixture()
	require.ErrorIs(t, err, ErrSQLServerCircuitFixtureUnavailable)
}

func TestSQLServerCircuitFixtureMissingEnvIsNotRun(t *testing.T) {
	t.Setenv(sqlServerCircuitFixtureEnv, "")
	_, err := LoadSQLServerCircuitFixture()
	require.ErrorIs(t, err, ErrSQLServerCircuitFixtureUnavailable)
}

func TestSQLServerCircuitFixtureManifestStrictGuards(t *testing.T) {
	base := SQLServerCircuitFixtureManifest{
		SchemaVersion: 1, Disposable: true, NativeAMD64: true, HostOS: "linux", HostArch: "amd64",
		Database: "recordings", Schema: "dbo", ServerVersionPrefix: "16.", EngineEdition: 3,
		CompatibilityLevel: 160, ConnectionConfigFile: "/abs/path/connection.json",
	}
	require.NoError(t, validateSQLServerCircuitManifest(base))

	t.Run("native true with arm64 must fail", func(t *testing.T) {
		m := base
		m.HostArch = "arm64"
		require.ErrorIs(t, validateSQLServerCircuitManifest(m), ErrSQLServerCircuitFixtureUnavailable)
	})

	t.Run("native false with amd64 must fail", func(t *testing.T) {
		m := base
		m.NativeAMD64 = false
		m.HostArch = "amd64"
		require.ErrorIs(t, validateSQLServerCircuitManifest(m), ErrSQLServerCircuitFixtureUnavailable)
	})

	t.Run("unsupported host_os must fail", func(t *testing.T) {
		m := base
		m.HostOS = "freebsd"
		require.ErrorIs(t, validateSQLServerCircuitManifest(m), ErrSQLServerCircuitFixtureUnavailable)
	})

	t.Run("version 16 paired with 170 must fail", func(t *testing.T) {
		m := base
		m.ServerVersionPrefix = "16."
		m.CompatibilityLevel = 170
		require.ErrorIs(t, validateSQLServerCircuitManifest(m), ErrSQLServerCircuitFixtureUnavailable)
	})

	t.Run("version 17 paired with 160 must fail", func(t *testing.T) {
		m := base
		m.NativeAMD64 = false
		m.HostArch = "arm64"
		m.HostOS = "darwin"
		m.ServerVersionPrefix = "17."
		m.CompatibilityLevel = 160
		require.ErrorIs(t, validateSQLServerCircuitManifest(m), ErrSQLServerCircuitFixtureUnavailable)
	})
}
