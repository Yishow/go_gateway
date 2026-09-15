package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
)

func reserveTCPPort(t *testing.T) int {
	t.Helper()
	var listenConfig net.ListenConfig
	ln, err := listenConfig.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve port failed: %v", err)
	}
	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port
}

func setupModbusShareHandlerForLifecycle(t *testing.T) (*ModbusShareHandler, *modbusshare.Service) {
	t.Helper()
	tagSvc := tag.NewService(tag.NewMemoryRepository())
	svc := modbusshare.NewService(tagSvc, 65536)
	t.Cleanup(func() {
		_ = svc.Stop()
	})
	return NewModbusShareHandler(svc), svc
}

func TestModbusShareHandler_StartAndStopLifecycle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, svc := setupModbusShareHandlerForLifecycle(t)
	router := gin.New()
	router.POST("/start", handler.Start)
	router.POST("/stop", handler.Stop)
	router.GET("/status", handler.Status)

	startPayload := map[string]int{"port": reserveTCPPort(t)}
	body, _ := json.Marshal(startPayload)
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on start, got %d: %s", w.Code, w.Body.String())
	}

	if !svc.Status().Enabled {
		t.Fatal("expected service enabled after start")
	}

	stopReq := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/stop", http.NoBody)
	stopW := httptest.NewRecorder()
	router.ServeHTTP(stopW, stopReq)
	if stopW.Code != http.StatusOK {
		t.Fatalf("expected 200 on stop, got %d: %s", stopW.Code, stopW.Body.String())
	}
	if svc.Status().Enabled {
		t.Fatal("expected service disabled after stop")
	}
}

func TestModbusShareHandler_StartPortConflictReturns409(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, svc := setupModbusShareHandlerForLifecycle(t)
	router := gin.New()
	router.POST("/start", handler.Start)

	port := reserveTCPPort(t)
	if err := svc.Start(port); err != nil {
		t.Fatalf("pre-start failed: %v", err)
	}

	body, _ := json.Marshal(map[string]int{"port": port})
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 on conflict, got %d: %s", w.Code, w.Body.String())
	}
}

func TestModbusShareHandler_StatusExposesSafeBindFailureDiagnostic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, svc := setupModbusShareHandlerForLifecycle(t)
	svc.SetHydrationState(modbusshare.HydrationState{State: modbusshare.HydrationStateReady, Readiness: true})
	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("occupy port: %v", err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	err = svc.ApplySettings(context.Background(), modbusshare.Settings{
		Enabled: true, BindAddress: "127.0.0.1", Port: port, SlaveID: 1, CapacityRegisters: 32768,
	})
	if err == nil {
		t.Fatal("expected bind conflict")
	}

	router := gin.New()
	router.GET("/status", handler.Status)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/status", http.NoBody))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Data struct {
			LastFailure *modbusshare.Diagnostic `json:"last_failure"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode status: %v", err)
	}
	if body.Data.LastFailure == nil {
		t.Fatal("expected bind failure diagnostic")
	}
	if body.Data.LastFailure.Code != modbusshare.ErrCodeListenerBindFailed || !body.Data.LastFailure.Retryable || body.Data.LastFailure.Action == "" {
		t.Fatalf("unexpected bind diagnostic: %+v", body.Data.LastFailure)
	}
	if strings.Contains(body.Data.LastFailure.Message, "address") || strings.Contains(body.Data.LastFailure.Message, "listen tcp") {
		t.Fatalf("raw bind error leaked: %q", body.Data.LastFailure.Message)
	}
}

func TestModbusShareHandler_StartRequiresExplicitPortWhenNoDurableGate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupModbusShareHandlerForLifecycle(t)
	router := gin.New()
	router.POST("/start", handler.Start)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/start", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 without explicit or persisted port, got %d: %s", w.Code, w.Body.String())
	}
}

func TestModbusShareHandler_StopIdempotent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, svc := setupModbusShareHandlerForLifecycle(t)
	router := gin.New()
	router.POST("/stop", handler.Stop)

	_ = svc.Stop()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/stop", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for stop idempotent call, got %d: %s", w.Code, w.Body.String())
	}
}

func TestModbusShareHandler_StartRejectsInvalidPort(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupModbusShareHandlerForLifecycle(t)
	router := gin.New()
	router.POST("/start", handler.Start)

	body, _ := json.Marshal(map[string]int{"port": -1})
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for invalid port, got %d: %s", w.Code, w.Body.String())
	}
}

func TestModbusShareHandler_StatusContainsBindState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, svc := setupModbusShareHandlerForLifecycle(t)
	router := gin.New()
	router.GET("/status", handler.Status)

	port := reserveTCPPort(t)
	if err := svc.Start(port); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/status", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var envelope struct {
		Success bool `json:"success"`
		Data    struct {
			BindState string `json:"bind_state"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !envelope.Success {
		t.Fatal("expected success true")
	}
	if envelope.Data.BindState != "pass" {
		t.Fatalf("expected bind_state=pass, got %q", envelope.Data.BindState)
	}

	if err := svc.Stop(); err != nil {
		t.Fatalf("stop failed: %v", err)
	}
	if svc.Status().BindState != "disabled" && svc.Status().BindState != "fail" {
		t.Fatalf("expected bind_state=disabled or fail after stop, got %q", svc.Status().BindState)
	}
}

func TestModbusShareHandler_StatusAvailableWithoutMappings(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupModbusShareHandlerForLifecycle(t)
	router := gin.New()
	router.GET("/status", handler.Status)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/status", http.NoBody)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var envelope struct {
		Success bool `json:"success"`
		Data    struct {
			MappingCount int    `json:"mapping_count"`
			BindState    string `json:"bind_state"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !envelope.Success {
		t.Fatal("expected success true")
	}
	if envelope.Data.MappingCount != 0 {
		t.Fatalf("expected mapping_count=0, got %d", envelope.Data.MappingCount)
	}
	if envelope.Data.BindState != "disabled" && envelope.Data.BindState != "fail" {
		t.Fatalf("expected bind_state=disabled or fail, got %q", envelope.Data.BindState)
	}
}

func TestModbusShareHandler_StartAndStopWithContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, svc := setupModbusShareHandlerForLifecycle(t)
	ctx := context.Background()
	if err := svc.Start(reserveTCPPort(t)); err != nil {
		t.Fatalf("setup start failed: %v", err)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequestWithContext(ctx, http.MethodPost, "/stop", http.NoBody)
	handler.Stop(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
