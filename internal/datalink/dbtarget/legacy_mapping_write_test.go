package dbtarget

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/tag"

	"github.com/stretchr/testify/require"
)

type legacyWriteCoordinatorProbe struct {
	db           *sql.DB
	preflightErr error
	checkErr     error
	events       []string
	mu           sync.Mutex
}

func (p *legacyWriteCoordinatorProbe) PreflightLegacyTargetWrite(context.Context, string, string, string) error {
	p.mu.Lock()
	p.events = append(p.events, "preflight")
	p.mu.Unlock()
	return p.preflightErr
}

func (p *legacyWriteCoordinatorProbe) RunLegacyTargetWrite(ctx context.Context, _ string, _ *schema.DatabaseTargetMapping, persist func(context.Context, *sql.Tx) error) error {
	p.mu.Lock()
	p.events = append(p.events, "run")
	p.mu.Unlock()
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := persist(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

func (p *legacyWriteCoordinatorProbe) CheckLegacyTargetWriteInTx(context.Context, *sql.Tx, string, string, string) error {
	p.mu.Lock()
	p.events = append(p.events, "check")
	p.mu.Unlock()
	return p.checkErr
}

func newLegacyMappingServiceFixture(t *testing.T, targetDSN string) (db *sql.DB, service *MappingService, tagID, connectorID string) {
	t.Helper()
	db = openMigratedTestDB(t)
	tagSvc := tag.NewService(tag.NewSQLRepository(db))
	tagEntity, err := tagSvc.Create(t.Context(), tag.CreateTagRequest{
		Key:         "temperature",
		DisplayName: "Temperature",
		DataType:    schema.DataTypeFloat64,
	})
	require.NoError(t, err)

	connectorID = "connector-legacy-write"
	config, err := json.Marshal(ConnectionConfig{"dsn": targetDSN})
	require.NoError(t, err)
	connectorRepo := NewSQLConnectorRepository(db)
	now := time.Now().UTC()
	require.NoError(t, connectorRepo.Create(t.Context(), &schema.DatabaseConnector{
		ID:                          connectorID,
		Name:                        "legacy-write-target",
		Kind:                        schema.DatabaseConnectorKindSQLite,
		ConnectionConfig:            string(config),
		IdentityRevision:            "connector-revision-1",
		Status:                      schema.DatabaseConnectorStatusReady,
		Enabled:                     true,
		DefaultWriteIntervalSeconds: 15,
		CreatedAt:                   now,
		UpdatedAt:                   now,
	}))

	mappingRepo := NewSQLTargetMappingRepository(db)
	service = NewMappingService(mappingRepo, connectorRepo, tagSvc)
	return db, service, tagEntity.ID, connectorID
}

func legacyMappingCreateRequest(tagID, connectorID string) CreateTargetMappingRequest {
	timestampColumn := "ts"
	return CreateTargetMappingRequest{
		TagID:           tagID,
		ConnectorID:     connectorID,
		TableName:       "sensor_values",
		ColumnName:      "value",
		WriteMode:       schema.DatabaseWriteModeUpsert,
		TimestampColumn: &timestampColumn,
	}
}

func TestMappingServicePreflightsBeforeConnectorInspection(t *testing.T) {
	missingTarget := t.TempDir() + "/missing.db"
	db, service, tagID, connectorID := newLegacyMappingServiceFixture(t, missingTarget)
	probe := &legacyWriteCoordinatorProbe{db: db, preflightErr: ErrConnectorRevisionConflict}
	WithLegacyWriteCoordinator(service, probe)

	_, err := service.Create(t.Context(), legacyMappingCreateRequest(tagID, connectorID))
	require.ErrorIs(t, err, ErrConnectorRevisionConflict)
	require.Equal(t, []string{"preflight"}, probe.events)
}

func TestMappingServicePersistsThroughCoordinatorTransaction(t *testing.T) {
	targetDSN := createTargetSQLite(t)
	db, service, tagID, connectorID := newLegacyMappingServiceFixture(t, targetDSN)
	probe := &legacyWriteCoordinatorProbe{db: db}
	WithLegacyWriteCoordinator(service, probe)

	mapping, err := service.Create(t.Context(), legacyMappingCreateRequest(tagID, connectorID))
	require.NoError(t, err)
	require.NotEmpty(t, mapping.ID)
	require.Equal(t, []string{"preflight", "preflight", "run", "check"}, probe.events)

	reloaded, err := service.GetByID(t.Context(), mapping.ID)
	require.NoError(t, err)
	require.Equal(t, mapping.ID, reloaded.ID)
}

func TestMappingServiceDeleteUsesCoordinatorAndEarlyOwnershipCheck(t *testing.T) {
	targetDSN := createTargetSQLite(t)
	db, service, tagID, connectorID := newLegacyMappingServiceFixture(t, targetDSN)
	mapping, err := service.Create(t.Context(), legacyMappingCreateRequest(tagID, connectorID))
	require.NoError(t, err)

	probe := &legacyWriteCoordinatorProbe{db: db}
	WithLegacyWriteCoordinator(service, probe)
	require.NoError(t, service.Delete(t.Context(), mapping.ID))
	_, err = service.GetByID(t.Context(), mapping.ID)
	require.Error(t, err)
	require.Equal(t, []string{"preflight", "run"}, probe.events)
}

func TestMappingServiceCoordinatorRollbackPreservesMapping(t *testing.T) {
	targetDSN := createTargetSQLite(t)
	db, service, tagID, connectorID := newLegacyMappingServiceFixture(t, targetDSN)
	mapping, err := service.Create(t.Context(), legacyMappingCreateRequest(tagID, connectorID))
	require.NoError(t, err)
	probe := &legacyWriteCoordinatorProbe{db: db, checkErr: errors.New("ownership changed during transaction")}
	WithLegacyWriteCoordinator(service, probe)

	newColumn := "value"
	_, err = service.Update(t.Context(), mapping.ID, UpdateTargetMappingRequest{ColumnName: &newColumn})
	require.Error(t, err)
	require.Contains(t, strings.ToLower(err.Error()), "ownership changed")
	require.Equal(t, []string{"preflight", "preflight", "run", "check"}, probe.events)

	reloaded, err := NewSQLTargetMappingRepository(db).GetByID(t.Context(), mapping.ID)
	require.NoError(t, err)
	require.Equal(t, mapping.ColumnName, reloaded.ColumnName)
}
