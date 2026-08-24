package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
	"go-gateway/internal/datalink/workspace"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestStudioV2WorkspaceDevicesHandler_DeleteCascadesOwnedSourceRulesInSQLite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	db := openStudioV2CascadeTestDB(t)
	defer db.Close()

	deviceRepo := device.NewSQLRepository(db)
	deviceSvc := device.NewService(deviceRepo, nil)
	workspaceSvc := workspace.NewService(workspace.NewSQLRepository(db))
	sourceRuleRepo := sourcerule.NewSQLRepository(db)
	sourceRuleSvc := sourcerule.NewService(sourceRuleRepo, deviceSvc, nil, nil)

	seedSQLiteCascadeDevice(t, ctx, deviceRepo, "device-target")
	seedSQLiteCascadeDevice(t, ctx, deviceRepo, "device-peer")
	_, err := workspaceSvc.AttachDevice(ctx, "device-target")
	require.NoError(t, err)
	_, err = workspaceSvc.AttachDevice(ctx, "device-peer")
	require.NoError(t, err)

	seedSQLiteCascadeRule(t, ctx, sourceRuleRepo, "rule-target", "device-target")
	seedSQLiteCascadeRule(t, ctx, sourceRuleRepo, "rule-peer", "device-peer")

	handler := NewStudioV2WorkspaceDevicesHandler(workspaceSvc, deviceSvc)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/datalink/studio-v2/workspace/devices/device-target", nil)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "device-target"}}

	handler.Delete(c)

	require.Equal(t, http.StatusOK, resp.Code, resp.Body.String())
	_, err = deviceRepo.GetByID(ctx, "device-target")
	require.Error(t, err)
	_, err = sourceRuleRepo.GetByID(ctx, "rule-target")
	require.Error(t, err, "device-owned source rule must follow device FK cascade")

	peerDevice, err := deviceRepo.GetByID(ctx, "device-peer")
	require.NoError(t, err)
	require.Equal(t, "device-peer", peerDevice.ID)
	peerRules, err := sourceRuleSvc.ListByDeviceIDs(ctx, []string{"device-peer"})
	require.NoError(t, err)
	require.Equal(t, []string{"rule-peer"}, sourceRuleIDs(peerRules))
	targetRules, err := sourceRuleSvc.ListByDeviceIDs(ctx, []string{"device-target"})
	require.NoError(t, err)
	require.Empty(t, targetRules)

	record, err := workspaceSvc.GetOrCreate(ctx)
	require.NoError(t, err)
	require.Equal(t, []string{"device-peer"}, record.OrderedDeviceIDs)
}

func openStudioV2CascadeTestDB(t *testing.T) *sql.DB {
	t.Helper()
	databasePath := filepath.ToSlash(filepath.Join(t.TempDir(), "studio-v2-cascade.db"))
	dsn := strings.Replace(datalink.DefaultEmbeddedSQLiteDSN, "datalink.db", databasePath, 1)
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	datalink.ApplySQLitePoolDefaults(db)
	require.NoError(t, datalink.NewMigrator().Migrate(db))

	var foreignKeys int
	require.NoError(t, db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys))
	t.Logf("SQLite foreign_keys=%d", foreignKeys)
	return db
}

func seedSQLiteCascadeDevice(t *testing.T, ctx context.Context, repo *device.SQLRepository, id string) {
	t.Helper()
	now := time.Now().UTC()
	require.NoError(t, repo.Create(ctx, &schema.Device{
		ID:               id,
		Name:             id,
		Protocol:         schema.ProtocolModbusTCP,
		Status:           schema.DeviceStatusActive,
		ConnectionConfig: `{"host":"127.0.0.1","port":502,"slave_id":1}`,
		CreatedAt:        now,
		UpdatedAt:        now,
	}))
}

func seedSQLiteCascadeRule(t *testing.T, ctx context.Context, repo *sourcerule.SQLRepository, id, deviceID string) {
	t.Helper()
	now := time.Now().UTC()
	require.NoError(t, repo.Create(ctx, &schema.SourceRule{
		ID:               id,
		DeviceID:         deviceID,
		StartAddress:     "40001",
		Count:            1,
		DataType:         schema.DataTypeInt16,
		NamingPrefix:     id + "_",
		Enabled:          true,
		Origin:           "manual",
		SkippedAddresses: "[]",
		RevisionID:       id + ":contract",
		CreatedAt:        now,
		UpdatedAt:        now,
	}))
}

func sourceRuleIDs(rules []*schema.SourceRule) []string {
	ids := make([]string, 0, len(rules))
	for _, rule := range rules {
		ids = append(ids, rule.ID)
	}
	return ids
}
