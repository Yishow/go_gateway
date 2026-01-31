package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go-gateway/internal/virtual/memory"
	modbusserver "go-gateway/internal/virtual/server/modbus"
	"go-gateway/internal/virtual/simulation"
)

// =============================================================================
// 任務 4.1-4.2: API Handler 測試
// =============================================================================

func setupTestRouter() (*gin.Engine, *VirtualDeviceHandler) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	bank := memory.NewMemoryBank(1024)
	server := modbusserver.NewServer(bank)
	sim := simulation.NewSimulationEngine(bank)

	handler := NewVirtualDeviceHandler(bank, server, sim)
	handler.RegisterRoutes(r.Group("/api"))

	return r, handler
}

func TestAPI_GetMemoryDump(t *testing.T) {
	r, _ := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/virtual/memory/dump", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("預期狀態碼 200，實際 %d", w.Code)
	}

	var resp MemoryDumpResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析回應失敗: %v", err)
	}

	if resp.Size != 1024 {
		t.Errorf("預期大小 1024，實際 %d", resp.Size)
	}
}

func TestAPI_WriteAndReadMemory(t *testing.T) {
	r, handler := setupTestRouter()

	// 寫入數據
	writeReq := WriteMemoryRequest{
		Offset: 0,
		Value:  intPtr(12345),
	}
	body, _ := json.Marshal(writeReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/virtual/memory/write", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("寫入失敗，狀態碼: %d, 回應: %s", w.Code, w.Body.String())
	}

	// 驗證記憶體值
	val, _ := handler.bank.ReadWord(0)
	if val != 12345 {
		t.Errorf("預期值 12345，實際 %d", val)
	}

	// 讀取範圍
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/virtual/memory/range?offset=0&length=4", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("讀取失敗，狀態碼: %d", w.Code)
	}

	var rangeResp MemoryRangeResponse
	json.Unmarshal(w.Body.Bytes(), &rangeResp)

	if len(rangeResp.Words) == 0 || rangeResp.Words[0] != 12345 {
		t.Errorf("讀取的 Words 不正確: %v", rangeResp.Words)
	}
}

func TestAPI_ClearMemory(t *testing.T) {
	r, handler := setupTestRouter()

	// 先寫入數據
	handler.bank.WriteWord(0, 9999)

	// 清空
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/virtual/memory/clear", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("清空失敗，狀態碼: %d", w.Code)
	}

	// 驗證
	val, _ := handler.bank.ReadWord(0)
	if val != 0 {
		t.Errorf("清空後值應為 0，實際 %d", val)
	}
}

func TestAPI_ServerStartStop(t *testing.T) {
	r, _ := setupTestRouter()

	// 啟動伺服器
	startReq := StartServerRequest{Port: 0}
	body, _ := json.Marshal(startReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/virtual/server/start", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("啟動伺服器失敗，狀態碼: %d, 回應: %s", w.Code, w.Body.String())
	}

	// 取得狀態
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/virtual/server/status", nil)
	r.ServeHTTP(w, req)

	var status map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &status)

	if status["running"] != true {
		t.Errorf("伺服器應該在運行中")
	}

	// 停止伺服器
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/virtual/server/stop", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("停止伺服器失敗，狀態碼: %d", w.Code)
	}
}

func TestAPI_AddSimulationRule(t *testing.T) {
	r, _ := setupTestRouter()

	ruleReq := AddSimulationRuleRequest{
		ID:         "test-rule",
		Type:       "increment",
		Offset:     100,
		IntervalMs: 1000,
		Increment:  5,
	}
	body, _ := json.Marshal(ruleReq)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/virtual/simulation/rules", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("新增規則失敗，狀態碼: %d, 回應: %s", w.Code, w.Body.String())
	}
}

// intPtr 建立 int 指標
func intPtr(i int) *int {
	return &i
}
