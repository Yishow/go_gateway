package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubPointManualPoller struct {
	results map[string]collector.CollectedValue
}

func (s *stubPointManualPoller) PollNow(pointIDs []string) []collector.CollectedValue {
	items := make([]collector.CollectedValue, 0, len(pointIDs))
	for _, pointID := range pointIDs {
		if result, ok := s.results[pointID]; ok {
			items = append(items, result)
		}
	}
	return items
}

type stubPointDirectReader struct {
	results map[string]collector.CollectedValue
}

func (s *stubPointDirectReader) PollDirect(
	_ context.Context,
	pointIDs []string,
) []collector.CollectedValue {
	items := make([]collector.CollectedValue, 0, len(pointIDs))
	for _, pointID := range pointIDs {
		if result, ok := s.results[pointID]; ok {
			items = append(items, result)
		}
	}
	return items
}

type failingPointReadResultRepo struct {
	*point.MemoryRepository
}

func (r *failingPointReadResultRepo) UpdateReadResult(
	ctx context.Context,
	id string,
	value interface{},
	errMsg string,
) error {
	return assert.AnError
}

func setupPointPollContractRouter(t *testing.T, poller *stubPointManualPoller) (*gin.Engine, *point.Service, *point.MemoryRepository) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	pointRepo := point.NewMemoryRepository()
	pointSvc := point.NewService(pointRepo, nil)

	mappingRepo := mapping.NewMemoryRepository()
	mappingSvc := mapping.NewService(mappingRepo)

	groupRepo := pollinggroup.NewMemoryRepository()
	groupSvc := pollinggroup.NewService(groupRepo)

	ctx := context.Background()
	groupID := "group-1"
	require.NoError(t, groupRepo.Create(ctx, &schema.PollingGroup{
		ID:         groupID,
		Name:       "group-1",
		IntervalMs: 1000,
		Priority:   100,
		Enabled:    true,
		CreatedAt:  time.Now().Add(-time.Hour),
		UpdatedAt:  time.Now().Add(-time.Hour),
	}))

	require.NoError(t, pointRepo.Create(ctx, &schema.Point{
		ID:             "point-1",
		DeviceID:       "device-1",
		Name:           "Point 1",
		Address:        "40001",
		DataType:       schema.DataTypeInt16,
		Mode:           schema.PointModeReadOnly,
		PollingGroupID: &groupID,
		Enabled:        true,
		CreatedAt:      time.Now().Add(-time.Hour),
		UpdatedAt:      time.Now().Add(-time.Hour),
	}))
	require.NoError(t, pointRepo.Create(ctx, &schema.Point{
		ID:             "point-2",
		DeviceID:       "device-1",
		Name:           "Point 2",
		Address:        "40002",
		DataType:       schema.DataTypeInt16,
		Mode:           schema.PointModeReadOnly,
		PollingGroupID: &groupID,
		Enabled:        true,
		CreatedAt:      time.Now().Add(-time.Hour),
		UpdatedAt:      time.Now().Add(-time.Hour),
	}))

	olderPipeline, err := json.Marshal([]schema.TransformStep{{
		Type:  schema.TransformScale,
		Order: 1,
		Params: map[string]interface{}{
			"multiplier": 0.1,
		},
	}})
	require.NoError(t, err)

	newerPipeline, err := json.Marshal([]schema.TransformStep{{
		Type:  schema.TransformScale,
		Order: 1,
		Params: map[string]interface{}{
			"multiplier": 0.2,
		},
	}})
	require.NoError(t, err)

	require.NoError(t, mappingRepo.Create(ctx, &schema.Mapping{
		ID:                "mapping-older",
		PointID:           "point-1",
		TagID:             "tag-1",
		TransformPipeline: string(olderPipeline),
		Enabled:           true,
		CreatedAt:         time.Date(2026, 3, 16, 6, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2026, 3, 16, 6, 0, 0, 0, time.UTC),
	}))
	require.NoError(t, mappingRepo.Create(ctx, &schema.Mapping{
		ID:                "mapping-newer",
		PointID:           "point-1",
		TagID:             "tag-2",
		TransformPipeline: string(newerPipeline),
		Enabled:           true,
		CreatedAt:         time.Date(2026, 3, 16, 6, 1, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2026, 3, 16, 6, 1, 0, 0, time.UTC),
	}))

	handler := NewPointHandler(pointSvc).WithPolling(poller, mappingSvc, groupSvc)

	router := gin.Default()
	router.POST("/datalink/points/:id/poll", handler.Poll)
	router.POST("/datalink/points/poll", handler.PollBatch)

	return router, pointSvc, pointRepo
}

func setupPointPollContractRouterWithDirectReader(
	t *testing.T,
	poller *stubPointManualPoller,
	directReader *stubPointDirectReader,
) (*gin.Engine, *point.Service, *point.MemoryRepository) {
	t.Helper()

	router, pointSvc, pointRepo := setupPointPollContractRouter(t, poller)
	gin.SetMode(gin.TestMode)

	mappingRepo := mapping.NewMemoryRepository()
	mappingSvc := mapping.NewService(mappingRepo)

	groupRepo := pollinggroup.NewMemoryRepository()
	groupSvc := pollinggroup.NewService(groupRepo)

	ctx := context.Background()
	groupID := "group-1"
	require.NoError(t, groupRepo.Create(ctx, &schema.PollingGroup{
		ID:         groupID,
		Name:       "group-1",
		IntervalMs: 1000,
		Priority:   100,
		Enabled:    true,
		CreatedAt:  time.Now().Add(-time.Hour),
		UpdatedAt:  time.Now().Add(-time.Hour),
	}))

	handler := NewPointHandler(pointSvc).
		WithPolling(poller, mappingSvc, groupSvc).
		WithDirectReader(directReader)

	router = gin.Default()
	router.POST("/datalink/points/:id/poll", handler.Poll)
	router.POST("/datalink/points/poll", handler.PollBatch)
	return router, pointSvc, pointRepo
}

func TestPointHandler_Poll_UsesManualPollAndLatestEnabledMapping(t *testing.T) {
	polledAt := time.Now().UTC().Truncate(time.Second)
	poller := &stubPointManualPoller{
		results: map[string]collector.CollectedValue{
			"point-1": {
				PointID:   "point-1",
				DeviceID:  "device-1",
				Value:     250,
				Timestamp: polledAt,
				Quality:   schema.QualityGood,
			},
		},
	}

	router, pointSvc, _ := setupPointPollContractRouter(t, poller)

	req := httptest.NewRequest(http.MethodPost, "/datalink/points/point-1/poll", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		Success bool       `json:"success"`
		Data    PollResult `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.True(t, body.Success)
	assert.Equal(t, "point-1", body.Data.PointID)
	assert.Equal(t, float64(250), body.Data.Value)
	assert.Equal(t, float64(50), body.Data.TransformedValue)
	assert.Equal(t, polledAt.Format(time.RFC3339), body.Data.Timestamp)
	assert.Equal(t, 192, body.Data.Quality)
	assert.False(t, body.Data.Stale)
	assert.Empty(t, body.Data.Error)

	storedPoint, err := pointSvc.GetByID(context.Background(), "point-1")
	require.NoError(t, err)
	require.NotNil(t, storedPoint.LastValue)
	assert.Equal(t, "250", *storedPoint.LastValue)
	assert.NotNil(t, storedPoint.LastReadAt)
}

func TestPointHandler_Poll_FallsBackToCachedStateAndComputesStale(t *testing.T) {
	router, _, pointRepo := setupPointPollContractRouter(t, nil)

	valueJSON, err := json.Marshal(123)
	require.NoError(t, err)
	readAt := time.Now().Add(-2 * time.Second)

	pt, err := pointRepo.GetByID(context.Background(), "point-2")
	require.NoError(t, err)
	pt.LastValue = ptrString(string(valueJSON))
	pt.LastReadAt = &readAt
	require.NoError(t, pointRepo.Update(context.Background(), pt))

	pollReq := httptest.NewRequest(http.MethodPost, "/datalink/points/point-2/poll", bytes.NewBuffer(nil))
	pollResp := httptest.NewRecorder()
	router.ServeHTTP(pollResp, pollReq)

	require.Equal(t, http.StatusOK, pollResp.Code)

	var body struct {
		Success bool       `json:"success"`
		Data    PollResult `json:"data"`
	}
	require.NoError(t, json.Unmarshal(pollResp.Body.Bytes(), &body))
	assert.True(t, body.Success)
	assert.Equal(t, "point-2", body.Data.PointID)
	assert.Equal(t, float64(123), body.Data.Value)
	assert.Equal(t, float64(123), body.Data.TransformedValue)
	assert.True(t, body.Data.Stale)
}

func TestPointHandler_Poll_FallbackWithoutLastReadAtIsStale(t *testing.T) {
	router, _, pointRepo := setupPointPollContractRouter(t, nil)

	valueJSON, err := json.Marshal(55)
	require.NoError(t, err)

	pt, err := pointRepo.GetByID(context.Background(), "point-2")
	require.NoError(t, err)
	pt.LastValue = ptrString(string(valueJSON))
	pt.LastReadAt = nil
	require.NoError(t, pointRepo.Update(context.Background(), pt))

	req := httptest.NewRequest(http.MethodPost, "/datalink/points/point-2/poll", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var body struct {
		Success bool       `json:"success"`
		Data    PollResult `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.True(t, body.Success)
	assert.True(t, body.Data.Stale)
	assert.Empty(t, body.Data.Timestamp)
}

func TestPointHandler_PollBatch_FallsBackToDirectReaderWhenManualPollIsUnavailable(t *testing.T) {
	polledAt := time.Now().UTC().Truncate(time.Second)
	directReader := &stubPointDirectReader{
		results: map[string]collector.CollectedValue{
			"point-1": {
				PointID:   "point-1",
				DeviceID:  "device-1",
				Value:     432,
				Timestamp: polledAt,
				Quality:   schema.QualityGood,
			},
		},
	}

	router, pointSvc, _ := setupPointPollContractRouterWithDirectReader(t, &stubPointManualPoller{}, directReader)

	body, err := json.Marshal(PollBatchRequest{PointIDs: &[]string{"point-1"}})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/datalink/points/poll", bytes.NewBuffer(body))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var payload struct {
		Success bool         `json:"success"`
		Data    []PollResult `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	require.True(t, payload.Success)
	require.Len(t, payload.Data, 1)
	assert.Equal(t, "point-1", payload.Data[0].PointID)
	assert.Equal(t, float64(432), payload.Data[0].Value)
	assert.Equal(t, polledAt.Format(time.RFC3339), payload.Data[0].Timestamp)
	assert.Empty(t, payload.Data[0].Error)

	storedPoint, err := pointSvc.GetByID(context.Background(), "point-1")
	require.NoError(t, err)
	require.NotNil(t, storedPoint.LastValue)
	assert.Equal(t, "432", *storedPoint.LastValue)
}

func TestPointHandler_PollBatch_UsesManualPollResultsAndKeepsMissingPoints(t *testing.T) {
	firstPolledAt := time.Now().UTC().Truncate(time.Second)
	secondPolledAt := firstPolledAt.Add(time.Second)
	poller := &stubPointManualPoller{
		results: map[string]collector.CollectedValue{
			"point-1": {
				PointID:   "point-1",
				DeviceID:  "device-1",
				Value:     250,
				Timestamp: firstPolledAt,
				Quality:   schema.QualityGood,
			},
			"point-2": {
				PointID:   "point-2",
				DeviceID:  "device-1",
				Value:     99,
				Timestamp: secondPolledAt,
				Quality:   schema.QualityGood,
			},
		},
	}

	router, _, _ := setupPointPollContractRouter(t, poller)

	body, err := json.Marshal(map[string]interface{}{
		"point_ids": []string{"point-1", "missing-point", "point-2"},
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/datalink/points/poll", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var payload struct {
		Success bool         `json:"success"`
		Data    []PollResult `json:"data"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.True(t, payload.Success)
	require.Len(t, payload.Data, 3)

	assert.Equal(t, "point-1", payload.Data[0].PointID)
	assert.Equal(t, float64(50), payload.Data[0].TransformedValue)
	assert.Empty(t, payload.Data[0].Error)

	assert.Equal(t, "missing-point", payload.Data[1].PointID)
	assert.Equal(t, "Point not found", payload.Data[1].Error)

	assert.Equal(t, "point-2", payload.Data[2].PointID)
	assert.Equal(t, float64(99), payload.Data[2].Value)
	assert.Equal(t, float64(99), payload.Data[2].TransformedValue)
}

func TestPointHandler_Poll_ReturnsInternalServerErrorForWritebackFailure(t *testing.T) {
	polledAt := time.Now().UTC().Truncate(time.Second)
	poller := &stubPointManualPoller{
		results: map[string]collector.CollectedValue{
			"point-1": {
				PointID:   "point-1",
				DeviceID:  "device-1",
				Value:     250,
				Timestamp: polledAt,
				Quality:   schema.QualityGood,
			},
		},
	}

	gin.SetMode(gin.TestMode)

	pointRepo := &failingPointReadResultRepo{MemoryRepository: point.NewMemoryRepository()}
	pointSvc := point.NewService(pointRepo, nil)
	mappingRepo := mapping.NewMemoryRepository()
	mappingSvc := mapping.NewService(mappingRepo)
	groupRepo := pollinggroup.NewMemoryRepository()
	groupSvc := pollinggroup.NewService(groupRepo)

	ctx := context.Background()
	groupID := "group-1"
	require.NoError(t, groupRepo.Create(ctx, &schema.PollingGroup{
		ID:         groupID,
		Name:       "group-1",
		IntervalMs: 1000,
		Priority:   100,
		Enabled:    true,
		CreatedAt:  time.Now().Add(-time.Hour),
		UpdatedAt:  time.Now().Add(-time.Hour),
	}))
	require.NoError(t, pointRepo.Create(ctx, &schema.Point{
		ID:             "point-1",
		DeviceID:       "device-1",
		Name:           "Point 1",
		Address:        "40001",
		DataType:       schema.DataTypeInt16,
		Mode:           schema.PointModeReadOnly,
		PollingGroupID: &groupID,
		Enabled:        true,
		CreatedAt:      time.Now().Add(-time.Hour),
		UpdatedAt:      time.Now().Add(-time.Hour),
	}))

	handler := NewPointHandler(pointSvc).WithPolling(poller, mappingSvc, groupSvc)
	router := gin.Default()
	router.POST("/datalink/points/:id/poll", handler.Poll)

	req := httptest.NewRequest(http.MethodPost, "/datalink/points/point-1/poll", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusInternalServerError, resp.Code)
}

func ptrString(value string) *string {
	return &value
}
