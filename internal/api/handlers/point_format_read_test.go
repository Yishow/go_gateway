package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
)

func TestPointFormatReadHTTPScopedWhitelist(t *testing.T) {
	db := openStudioV2CascadeTestDB(t)
	t.Cleanup(func() { _ = db.Close() })
	now := time.Now().UTC()
	require.NoError(t, device.NewSQLRepository(db).Create(t.Context(), &schema.Device{ID: "format-read-device", Name: "d", Protocol: schema.ProtocolModbusTCP, Status: schema.DeviceStatusDraft, ConnectionConfig: `{"data_format":"ABCD","password":"private-secret"}`, CreatedAt: now, UpdatedAt: now}))
	ws, err := workspace.NewService(workspace.NewSQLRepository(db)).AttachDevice(t.Context(), "format-read-device")
	require.NoError(t, err)
	svc := point.NewService(point.NewSQLRepository(db), nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: "format-read-device", Name: "p", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	router := gin.New()
	router.GET("/points/:id/format", NewPointHandler(svc).ReadFormat)
	path := "/points/" + p.ID + "/format?workspace_id=" + ws.ID + "&device_id=" + p.DeviceID
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
	require.Equal(t, 200, w.Code, w.Body.String())
	require.NotContains(t, w.Body.String(), "private-secret")
	require.NotContains(t, w.Body.String(), "connection_config")
	require.NotContains(t, w.Body.String(), "last_error")
	var response PointFormatReadResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	require.Equal(t, p.RevisionID, response.Data.PointRevision)
	require.Equal(t, schema.PointFormatInherit, response.Data.Resolution.Mode)
	require.Equal(t, "ABCD", response.Data.Resolution.EffectiveFormat)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", path+"&rule_id=foreign", nil))
	require.Equal(t, 403, w.Code)
}
