package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gateway/internal/datalink/modbusshare"
	"go-gateway/internal/datalink/tag"

	"github.com/gin-gonic/gin"
)

func reserveTCPPort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
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
	return NewModbusShareHandler(svc, nil, nil), svc
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
	req := httptest.NewRequest(http.MethodPost, "/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on start, got %d: %s", w.Code, w.Body.String())
	}

	if !svc.Status().Enabled {
		t.Fatal("expected service enabled after start")
	}

	stopReq := httptest.NewRequest(http.MethodPost, "/stop", nil)
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
	req := httptest.NewRequest(http.MethodPost, "/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409 on conflict, got %d: %s", w.Code, w.Body.String())
	}
}

func TestModbusShareHandler_StartDefaultsTo5020WhenPayloadEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupModbusShareHandlerForLifecycle(t)
	router := gin.New()
	router.POST("/start", handler.Start)

	req := httptest.NewRequest(http.MethodPost, "/start", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusConflict {
		t.Fatalf("expected 200 or 409 with empty payload, got %d: %s", w.Code, w.Body.String())
	}

	if w.Code == http.StatusOK {
		var envelope struct {
			Data struct {
				Port int `json:"port"`
			} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
			t.Fatalf("unmarshal failed: %v", err)
		}
		if envelope.Data.Port != 5020 {
			t.Fatalf("expected default port 5020, got %d", envelope.Data.Port)
		}
	}
}

func TestModbusShareHandler_StopIdempotent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, svc := setupModbusShareHandlerForLifecycle(t)
	router := gin.New()
	router.POST("/stop", handler.Stop)

	_ = svc.Stop()

	req := httptest.NewRequest(http.MethodPost, "/stop", nil)
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
	req := httptest.NewRequest(http.MethodPost, "/start", bytes.NewReader(body))
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

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
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
	if svc.Status().BindState != "fail" {
		t.Fatalf("expected bind_state=fail after stop, got %q", svc.Status().BindState)
	}
}

func TestModbusShareHandler_StatusAvailableWithoutMappings(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _ := setupModbusShareHandlerForLifecycle(t)
	router := gin.New()
	router.GET("/status", handler.Status)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
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
	if envelope.Data.BindState != "fail" {
		t.Fatalf("expected bind_state=fail, got %q", envelope.Data.BindState)
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
	c.Request = httptest.NewRequest(http.MethodPost, "/stop", nil).WithContext(ctx)
	handler.Stop(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
