package dbtarget

import (
	"path/filepath"
	"testing"
	"time"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

func TestConnectorService_IdentityRevisionLifecycle(t *testing.T) {
	probeCalls := installIdentityProbeStub(t)
	ctx := t.Context()
	db := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(db)
	service := NewConnectorService(repo)

	connector, err := service.Create(ctx, CreateConnectorRequest{
		Name: "saved-postgres",
		Kind: schema.DatabaseConnectorKindPostgres,
		ConnectionConfig: ConnectionConfig{
			"host": "db-a.internal", "port": 5432, "user": "writer_a",
			"password": "  old\tsecret\n", "database": "metrics", "table": "samples",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, connector.IdentityRevision)
	initialRevision := connector.IdentityRevision

	name := "renamed"
	connector, err = service.Update(ctx, connector.ID, UpdateConnectorRequest{Name: &name})
	require.NoError(t, err)
	require.Equal(t, initialRevision, connector.IdentityRevision)

	connector, err = service.Update(ctx, connector.ID, UpdateConnectorRequest{
		ConnectionConfig: &ConnectionConfig{
			"host": "db-a.internal", "port": 5432, "user": "writer_a",
			"password": "", "database": "metrics", "table": "next_samples",
		},
	})
	require.NoError(t, err)
	require.Equal(t, initialRevision, connector.IdentityRevision)
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	require.NoError(t, err)
	require.Equal(t, "  old\tsecret\n", config["password"])

	interval := 30
	connector, err = service.Update(ctx, connector.ID, UpdateConnectorRequest{DefaultWriteIntervalSeconds: &interval})
	require.NoError(t, err)
	require.Equal(t, initialRevision, connector.IdentityRevision)

	connector, err = service.TestConnection(ctx, connector.ID)
	require.NoError(t, err)
	require.Equal(t, initialRevision, connector.IdentityRevision)

	checkAt := time.Now().UTC()
	require.NoError(t, repo.UpdateWriteOutcome(ctx, connector.ID, checkAt, "success", ""))
	require.NoError(t, repo.UpdateFlushOutcome(ctx, connector.ID, checkAt, "success", ""))
	require.NoError(t, repo.UpdateSchemaEnsureOutcome(ctx, connector.ID, checkAt, "success", ""))
	connector, err = service.GetByID(ctx, connector.ID)
	require.NoError(t, err)
	require.Equal(t, initialRevision, connector.IdentityRevision)

	kind := schema.DatabaseConnectorKindMySQL
	connector, err = service.Update(ctx, connector.ID, UpdateConnectorRequest{Kind: &kind})
	require.NoError(t, err)
	require.NotEqual(t, initialRevision, connector.IdentityRevision)
	kindRevision := connector.IdentityRevision
	config, err = parseConnectionConfig(connector.ConnectionConfig)
	require.NoError(t, err)
	_, hasPassword := config["password"]
	require.False(t, hasPassword, "a kind change must not reuse the previous password")

	connector, err = service.Update(ctx, connector.ID, UpdateConnectorRequest{
		ConnectionConfig: &ConnectionConfig{
			"host": "db-b.internal", "port": 5432, "user": "writer_a",
			"password": "", "database": "metrics", "table": "next_samples",
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, kindRevision, connector.IdentityRevision)
	endpointRevision := connector.IdentityRevision
	config, err = parseConnectionConfig(connector.ConnectionConfig)
	require.NoError(t, err)
	_, hasPassword = config["password"]
	require.False(t, hasPassword, "an endpoint change must not reuse the previous password")

	newPassword := "  new\tpassword\n"
	connector, err = service.Update(ctx, connector.ID, UpdateConnectorRequest{
		ConnectionConfig: &ConnectionConfig{
			"host": "db-b.internal", "port": 5432, "user": "writer_a",
			"password": newPassword, "database": "metrics", "table": "next_samples",
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, endpointRevision, connector.IdentityRevision)
	passwordRevision := connector.IdentityRevision
	config, err = parseConnectionConfig(connector.ConnectionConfig)
	require.NoError(t, err)
	require.Equal(t, newPassword, config["password"])

	clearPassword := true
	connector, err = service.Update(ctx, connector.ID, UpdateConnectorRequest{
		ConnectionConfig: &ConnectionConfig{
			"host": "db-b.internal", "port": 5432, "user": "writer_a",
			"password": "", "database": "metrics", "table": "next_samples",
		},
		ClearPassword: &clearPassword,
	})
	require.NoError(t, err)
	require.NotEqual(t, passwordRevision, connector.IdentityRevision)
	config, err = parseConnectionConfig(connector.ConnectionConfig)
	require.NoError(t, err)
	_, hasPassword = config["password"]
	require.False(t, hasPassword)
	require.NotZero(t, *probeCalls)
}

func TestConnectorService_ResolveSavedTargetRequiresEnabledMatchingRevision(t *testing.T) {
	probeCalls := installIdentityProbeStub(t)
	ctx := t.Context()
	db := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(db)
	service := NewConnectorService(repo)

	connector := &schema.DatabaseConnector{
		ID:               "saved-postgres",
		Name:             "saved postgres",
		Kind:             schema.DatabaseConnectorKindPostgres,
		ConnectionConfig: `{"host":"db.internal","user":"writer","database":"metrics","password":"secret"}`,
		Status:           schema.DatabaseConnectorStatusUnreachable,
		Enabled:          true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	require.NoError(t, repo.Create(ctx, connector))
	require.NotEmpty(t, connector.IdentityRevision)

	resolved, err := service.ResolveSavedTarget(ctx, connector.ID, connector.IdentityRevision)
	require.NoError(t, err)
	require.Equal(t, connector.ID, resolved.ID)
	require.Equal(t, connector.ConnectionConfig, resolved.ConnectionConfig)
	require.Zero(t, *probeCalls, "resolving a saved target must not probe or connect")

	_, err = service.ResolveSavedTarget(ctx, "missing", connector.IdentityRevision)
	require.ErrorIs(t, err, ErrConnectorNotFound)
	_, err = service.ResolveSavedTarget(ctx, connector.ID, "")
	require.ErrorIs(t, err, ErrConnectorRevisionConflict)
	_, err = service.ResolveSavedTarget(ctx, connector.ID, "stale-revision")
	require.ErrorIs(t, err, ErrConnectorRevisionConflict)

	disabled := *connector
	disabled.Enabled = false
	require.NoError(t, repo.UpdateWithExpectedIdentityRevision(ctx, &disabled, connector.IdentityRevision))
	_, err = service.ResolveSavedTarget(ctx, connector.ID, connector.IdentityRevision)
	require.ErrorIs(t, err, ErrConnectorDisabled)
	require.NotErrorIs(t, err, ErrConnectorNotFound, "a disabled connector exists and must not look missing")
	require.Zero(t, *probeCalls)
}

func TestConnectorService_EndpointChangeAcceptsExplicitPassword(t *testing.T) {
	installIdentityProbeStub(t)
	ctx := t.Context()
	db := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(db)
	service := NewConnectorService(repo)

	connector, err := service.Create(ctx, CreateConnectorRequest{
		Name: "explicit-password",
		Kind: schema.DatabaseConnectorKindPostgres,
		ConnectionConfig: ConnectionConfig{
			"host": "db-a.internal", "port": 5432, "user": "writer",
			"password": "same-secret", "database": "metrics",
		},
	})
	require.NoError(t, err)
	previousRevision := connector.IdentityRevision

	connector, err = service.Update(ctx, connector.ID, UpdateConnectorRequest{
		ConnectionConfig: &ConnectionConfig{
			"host": "db-b.internal", "port": 5432, "user": "writer",
			"password": "same-secret", "database": "metrics",
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, previousRevision, connector.IdentityRevision)
	config, err := parseConnectionConfig(connector.ConnectionConfig)
	require.NoError(t, err)
	require.Equal(t, "same-secret", config["password"])
}

func TestSQLConnectorRepository_RejectsStaleIdentityRevision(t *testing.T) {
	ctx := t.Context()
	db := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(db)
	connector := &schema.DatabaseConnector{
		ID:               "cas-connector",
		Name:             "original",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: `{"dsn":"/tmp/cas.db"}`,
		Status:           schema.DatabaseConnectorStatusReady,
		Enabled:          true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	require.NoError(t, repo.Create(ctx, connector))
	stale := *connector

	fresh := *connector
	fresh.Name = "fresh"
	fresh.IdentityRevision = "fresh-revision"
	require.NoError(t, repo.UpdateWithExpectedIdentityRevision(ctx, &fresh, connector.IdentityRevision))

	stale.Name = "stale"
	err := repo.UpdateWithExpectedIdentityRevision(ctx, &stale, connector.IdentityRevision)
	require.ErrorIs(t, err, ErrConnectorRevisionConflict)

	reloaded, err := repo.GetByID(ctx, connector.ID)
	require.NoError(t, err)
	require.Equal(t, "fresh", reloaded.Name)
	require.Equal(t, "fresh-revision", reloaded.IdentityRevision)
}

func TestConnectorService_UpdateRejectsStaleExpectedRevisionBeforeProbe(t *testing.T) {
	probeCalls := installIdentityProbeStub(t)
	ctx := t.Context()
	db := openMigratedTestDB(t)
	repo := NewSQLConnectorRepository(db)
	service := NewConnectorService(repo)

	connector, err := service.Create(ctx, CreateConnectorRequest{
		Name:             "stale-update",
		Kind:             schema.DatabaseConnectorKindSQLite,
		ConnectionConfig: ConnectionConfig{"dsn": "/tmp/stale-update.db"},
	})
	require.NoError(t, err)
	staleRevision := connector.IdentityRevision

	password := "new-secret"
	updated, err := service.Update(ctx, connector.ID, UpdateConnectorRequest{
		ConnectionConfig: &ConnectionConfig{"dsn": "/tmp/stale-update.db", "password": password},
	})
	require.NoError(t, err)
	require.NotEqual(t, staleRevision, updated.IdentityRevision)
	callsBeforeStaleUpdate := *probeCalls

	name := "must-not-win"
	_, err = service.Update(ctx, connector.ID, UpdateConnectorRequest{
		Name:                     &name,
		ExpectedIdentityRevision: &staleRevision,
	})
	require.ErrorIs(t, err, ErrConnectorRevisionConflict)
	require.Equal(t, callsBeforeStaleUpdate, *probeCalls, "stale revisions must fail before probing")

	reloaded, err := service.GetByID(ctx, connector.ID)
	require.NoError(t, err)
	require.NotEqual(t, name, reloaded.Name)
}

func installIdentityProbeStub(t *testing.T) *int {
	t.Helper()
	originalOpen := openExternalDBManagerFunc
	calls := 0
	probePath := filepath.Join(t.TempDir(), "probe.db")
	openExternalDBManagerFunc = func(kind schema.DatabaseConnectorKind, config ConnectionConfig) (*datalinkbase.DBManager, error) {
		calls++
		manager := datalinkbase.NewDBManager(datalinkbase.DBConfig{
			Type:         datalinkbase.DBTypeSQLite,
			DSN:          probePath,
			MaxOpenConns: 1,
			MaxIdleConns: 1,
		})
		if err := manager.Connect(); err != nil {
			return nil, err
		}
		return manager, nil
	}
	t.Cleanup(func() { openExternalDBManagerFunc = originalOpen })
	return &calls
}

func TestSavedConnectorIdentityKindChangeDoesNotReuseOmittedPassword(t *testing.T) {
	installIdentityProbeStub(t)
	repo := NewSQLConnectorRepository(openMigratedTestDB(t))
	service := NewConnectorService(repo)
	config := ConnectionConfig{"host": "127.0.0.1", "user": "writer", "database": "metrics", "password": "fixture password"}
	connector, err := service.Create(t.Context(), CreateConnectorRequest{Name: "kind change", Kind: schema.DatabaseConnectorKindPostgres, ConnectionConfig: config})
	require.NoError(t, err)
	delete(config, "password")
	nextKind := schema.DatabaseConnectorKindMySQL
	saved, err := service.Update(t.Context(), connector.ID, UpdateConnectorRequest{Kind: &nextKind, ConnectionConfig: &config})
	require.NoError(t, err)
	actual, err := parseConnectionConfig(saved.ConnectionConfig)
	require.NoError(t, err)
	require.NotContains(t, actual, "password")
}
