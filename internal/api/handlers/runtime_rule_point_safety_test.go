package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRuntimeRulePointSafetyRemoteDiagnostics(t *testing.T) {
	repo := point.NewMemoryRepository()
	raw := `123.45`
	require.NoError(t, repo.Create(t.Context(), &schema.Point{ID: "point", DeviceID: "device", Address: "40001", DataType: schema.DataTypeFloat64, LastValue: &raw, LastError: "driver postgres://user:fixture-password@private.internal/db SELECT secret"}))
	h := NewPointHandler(point.NewService(repo, nil))
	for _, method := range []string{"list", "get"} {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		c.Request = httptest.NewRequest(http.MethodGet, "/points?device_id=device", nil)
		c.Request.RemoteAddr = "203.0.113.9:1234"
		c.Params = gin.Params{{Key: "id", Value: "point"}}
		if method == "list" {
			h.List(c)
		} else {
			h.Get(c)
		}
		require.Equal(t, 200, rec.Code)
		require.NotContains(t, rec.Body.String(), "fixture-password")
		require.NotContains(t, rec.Body.String(), "private.internal")
		require.Contains(t, rec.Body.String(), "123.45")
		require.Contains(t, rec.Body.String(), "40001")
	}
	saved, err := repo.GetByID(t.Context(), "point")
	require.NoError(t, err)
	require.Contains(t, saved.LastError, "fixture-password")
	require.Equal(t, raw, *saved.LastValue)
}

func TestRuntimeRulePointSafetyScalarCompatibility(t *testing.T) {
	for _, tc := range []struct {
		value  string
		unsafe bool
	}{
		{`18446744073709551615`, false}, {`123.45`, false}, {`true`, false}, {`"OK/NG"`, false}, {`"2026-10-28"`, false}, {`"SELECT mode"`, false}, {`"°C"`, false}, {`"40001"`, false}, {`plain text`, false},
		{`{"nested":{"password":"fixture-password"}}`, true}, {`["fixture-password"]`, true}, {`"postgres://user:fixture-password@internal/db"`, true}, {`"password=fixture-password"`, true}, {`"/private/fixture.sqlite"`, true}, {`"SELECT secret FROM credentials"`, true},
	} {
		t.Run(tc.value, func(t *testing.T) {
			repo := point.NewMemoryRepository()
			raw := tc.value
			require.NoError(t, repo.Create(t.Context(), &schema.Point{ID: "point", DataType: schema.DataTypeString, LastValue: &raw}))
			h := NewPointHandler(point.NewService(repo, nil))
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodGet, "/points", nil)
			h.List(c)
			var body struct {
				Data []struct {
					LastValue  *string `json:"last_value"`
					ValueState string  `json:"value_state"`
				} `json:"data"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			require.Len(t, body.Data, 1)
			if tc.unsafe {
				require.Nil(t, body.Data[0].LastValue)
				require.Equal(t, "unsupported", body.Data[0].ValueState)
			} else {
				require.NotNil(t, body.Data[0].LastValue)
				require.Equal(t, tc.value, *body.Data[0].LastValue)
			}
		})
	}
}

type runtimePointFailedRepository struct{ point.Repository }

func (runtimePointFailedRepository) List(context.Context, point.ListFilter) ([]*schema.Point, error) {
	return nil, errors.New("driver postgres://fixture-password@private.internal/db SELECT secret")
}
func TestRuntimeRulePointSafetyReadFailure(t *testing.T) {
	h := NewPointHandler(point.NewService(runtimePointFailedRepository{point.NewMemoryRepository()}, nil))
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/points", nil)
	h.List(c)
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	require.NotContains(t, rec.Body.String(), "fixture-password")
	require.Contains(t, rec.Body.String(), ErrCodeRuntimeSnapshotUnavailable)
}

func TestRuntimeRulePointSafetyMetadata(t *testing.T) {
	for _, address := range []string{"40001", "D0100", "0xFFFF", "400001", "R100.2"} {
		p := &schema.Point{ID: "point", DeviceID: "device", Name: "Pump A (°C)", Description: "SELECT mode OK/NG", Address: address, Function: "03", DataType: schema.DataTypeFloat64, DataFormat: "CDAB", Mode: schema.PointModeReadOnly}
		result := safeRuntimePointRead(p)
		require.Equal(t, address, result.Address)
		require.Equal(t, p.Name, result.Name)
		require.Equal(t, p.Description, result.Description)
		require.Equal(t, "03", result.Function)
		require.Equal(t, "CDAB", result.DataFormat)
	}
	for _, field := range []string{"name", "description", "address", "function", "format", "mode", "id", "device", "group"} {
		p := &schema.Point{ID: "point", DeviceID: "device", Name: "Point", Address: "40001", DataType: schema.DataTypeFloat64}
		unsafe := "password=fixture-password /private/fixture.sqlite"
		switch field {
		case "name":
			p.Name = unsafe
		case "description":
			p.Description = unsafe
		case "address":
			p.Address = unsafe
		case "function":
			p.Function = unsafe
		case "format":
			p.DataFormat = unsafe
		case "mode":
			p.Mode = schema.PointMode(unsafe)
		case "id":
			p.ID = unsafe
		case "device":
			p.DeviceID = unsafe
		case "group":
			p.PollingGroupID = &unsafe
		}
		result := safeRuntimePointRead(p)
		raw, err := json.Marshal(result)
		require.NoError(t, err)
		require.NotContains(t, string(raw), "fixture", field)
		require.Contains(t, string(raw), `"metadata_state":"unsupported"`)
	}
}
