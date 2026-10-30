package handlers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/schema"
	"net/http/httptest"
	"testing"
)

func TestPointLegacyFormatPUTRequiresExplicitGuards(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := point.NewMemoryRepository()
	svc := point.NewService(repo, nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: "dev", Name: "point", Address: "HR0", DataType: schema.DataTypeUint32})
	require.NoError(t, err)
	h := NewPointHandler(svc)
	router := gin.New()
	router.PUT("/points/:id", h.Update)
	for _, payload := range []string{`{"data_format":"CDAB"}`, `{"data_format":""}`, `{"format_mode":"override","configured_format":"ABCD"}`} {
		r := httptest.NewRequest("PUT", "/points/"+p.ID, bytes.NewBufferString(payload))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)
		require.Equal(t, 422, w.Code, payload)
		got, err := svc.GetByID(t.Context(), p.ID)
		require.NoError(t, err)
		require.Equal(t, p.RevisionID, got.RevisionID)
		require.Equal(t, p.FormatMode, got.FormatMode)
	}
}

func TestPointOrdinaryTypePUTRejectsInapplicableOverride(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := point.NewMemoryRepository()
	svc := point.NewService(repo, nil)
	p, err := svc.Create(t.Context(), point.CreatePointRequest{DeviceID: "d", Name: "p", Address: "40001", DataType: schema.DataTypeInt32})
	require.NoError(t, err)
	p.FormatMode = schema.PointFormatOverride
	p.FormatOverride = "ABCD"
	p.DataFormat = "ABCD"
	require.NoError(t, repo.Update(t.Context(), p))
	router := gin.New()
	router.PUT("/points/:id", NewPointHandler(svc).Update)
	request := httptest.NewRequest("PUT", "/points/"+p.ID, bytes.NewBufferString(`{"data_type":"bool"}`))
	request.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	require.Equal(t, 422, w.Code)
	stored, err := svc.GetByID(t.Context(), p.ID)
	require.NoError(t, err)
	require.Equal(t, p.RevisionID, stored.RevisionID)
}
