package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/schema"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type reorderedAcquisitionRuntime struct {
	mu      sync.Mutex
	calls   int
	current *schema.Device
	first   chan struct{}
	release chan struct{}
}

func (r *reorderedAcquisitionRuntime) UpsertDevice(ctx context.Context, dev *schema.Device) error {
	r.mu.Lock()
	r.calls++
	first := r.calls == 1
	r.mu.Unlock()
	if first {
		close(r.first)
		select {
		case <-r.release:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	r.mu.Lock()
	deviceCopy := *dev
	r.current = &deviceCopy
	r.mu.Unlock()
	return nil
}

func (*reorderedAcquisitionRuntime) RemoveDevice(string) {}

func TestStudioConcurrentAcquisitionUpdatesApplyLatestSavedSettings(t *testing.T) {
	devSvc, _, workspaceSvc, _ := setupDeviceAcquisitionHandlerRig(t)
	runtime := &reorderedAcquisitionRuntime{first: make(chan struct{}), release: make(chan struct{})}
	handler := NewStudioV2WorkspaceDevicesHandler(workspaceSvc, devSvc, runtime)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	send := func(payload string) int {
		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)
		c.Params = gin.Params{{Key: "id", Value: "dev-a"}}
		c.Request = httptest.NewRequestWithContext(ctx, http.MethodPut, "/", strings.NewReader(payload))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.Update(c) //nolint:contextcheck // The request carries the test context into the handler.
		return response.Code
	}
	firstDone, secondDone := make(chan int, 1), make(chan int, 1)
	go func() { firstDone <- send(`{"acquisition_auto_interval":true}`) }()
	select {
	case <-runtime.first:
	case <-ctx.Done():
		t.Fatal("first update did not reach runtime")
	}
	go func() { secondDone <- send(`{"frozen_threshold_ms":600000}`) }()
	var secondCode int
	select {
	case secondCode = <-secondDone:
	case <-time.After(500 * time.Millisecond):
	}
	close(runtime.release)
	require.Equal(t, http.StatusOK, <-firstDone)
	if secondCode == 0 {
		secondCode = <-secondDone
	}
	require.Equal(t, http.StatusOK, secondCode)
	saved, err := devSvc.GetByID(t.Context(), "dev-a")
	require.NoError(t, err)
	require.True(t, saved.AcquisitionAutoInterval)
	require.Equal(t, 600000, saved.FrozenThresholdMs)
	runtime.mu.Lock()
	defer runtime.mu.Unlock()
	require.NotNil(t, runtime.current)
	require.True(t, runtime.current.AcquisitionAutoInterval)
	require.Equal(t, 600000, runtime.current.FrozenThresholdMs, "runtime must not end with the older successful image")
}

func TestStudioAndDeviceAPIsApplyLatestAcquisitionSettings(t *testing.T) {
	for _, firstStudio := range []bool{true, false} {
		name := "device first"
		if firstStudio {
			name = "Studio first"
		}
		t.Run(name, func(t *testing.T) {
			testCrossAPIAcquisitionSettings(t, firstStudio)
		})
	}
}

func testCrossAPIAcquisitionSettings(t *testing.T, firstStudio bool) {
	t.Helper()
	devSvc, _, workspaceSvc, _ := setupDeviceAcquisitionHandlerRig(t)
	runtime := &reorderedAcquisitionRuntime{first: make(chan struct{}), release: make(chan struct{})}
	studio := NewStudioV2WorkspaceDevicesHandler(workspaceSvc, devSvc, runtime)
	legacy := NewDeviceHandler(devSvc, runtime)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	send := func(payload string, update func(*gin.Context)) int {
		response := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(response)
		c.Params = gin.Params{{Key: "id", Value: "dev-a"}}
		c.Request = httptest.NewRequestWithContext(ctx, http.MethodPut, "/", strings.NewReader(payload))
		c.Request.Header.Set("Content-Type", "application/json")
		update(c)
		return response.Code
	}
	firstUpdate, secondUpdate := legacy.Update, studio.Update
	if firstStudio {
		firstUpdate, secondUpdate = studio.Update, legacy.Update
	}
	firstDone, secondDone := make(chan int, 1), make(chan int, 1)
	go func() { firstDone <- send(`{"acquisition_auto_interval":true}`, firstUpdate) }()
	select {
	case <-runtime.first:
	case <-ctx.Done():
		t.Fatal("first update did not reach runtime")
	}
	go func() { secondDone <- send(`{"frozen_threshold_ms":600000}`, secondUpdate) }()
	var secondCode int
	select {
	case secondCode = <-secondDone:
	case <-time.After(500 * time.Millisecond):
	}
	close(runtime.release)
	require.Equal(t, http.StatusOK, <-firstDone)
	if secondCode == 0 {
		secondCode = <-secondDone
	}
	require.Equal(t, http.StatusOK, secondCode)
	saved, err := devSvc.GetByID(t.Context(), "dev-a")
	require.NoError(t, err)
	require.True(t, saved.AcquisitionAutoInterval)
	require.Equal(t, 600000, saved.FrozenThresholdMs)
	runtime.mu.Lock()
	defer runtime.mu.Unlock()
	require.True(t, runtime.current.AcquisitionAutoInterval)
	require.Equal(t, 600000, runtime.current.FrozenThresholdMs, "runtime must retain the later API update")
}

func TestDeviceAPIAcquisitionRuntimeFailureIsSafeAndReportsSaved(t *testing.T) {
	devSvc, _, _, syncer := setupDeviceAcquisitionHandlerRig(t)
	syncer.upsertErr = errors.New("secret-runtime-detail")
	handler := NewDeviceHandler(devSvc, syncer)
	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Params = gin.Params{{Key: "id", Value: "dev-a"}}
	c.Request = httptest.NewRequestWithContext(t.Context(), http.MethodPut, "/", strings.NewReader(`{"acquisition_auto_interval":true}`))
	c.Request.Header.Set("Content-Type", "application/json")
	handler.Update(c) //nolint:contextcheck // The request carries the test context into the handler.
	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.Contains(t, response.Body.String(), "設定已儲存")
	require.NotContains(t, response.Body.String(), "secret-runtime-detail")
	saved, err := devSvc.GetByID(t.Context(), "dev-a")
	require.NoError(t, err)
	require.True(t, saved.AcquisitionAutoInterval)
}
