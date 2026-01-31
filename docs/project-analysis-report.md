# Go Gateway 專案深度分析報告

> **分析日期**: 2026-01-31
> **專案版本**: go 1.25.5 + React 19 + TypeScript
> **分析範圍**: 後端架構、前端架構、規格文檔、建置部署

---

## 📊 專案總覽

**Go Gateway** 是一個工業級數據採集閘道系統,用於從 PLC 採集數據並映射至時序資料庫。

| 維度 | 評分 | 說明 |
|------|------|------|
| 後端架構 | 6.4/10 | 良好的模組化，但存在安全與效能問題 |
| 前端架構 | 6.5/10 | 現代化堆疊，但測試覆蓋率極低 |
| 規格文檔 | 7.5/10 | OpenSpec 結構完善，但部分過時 |
| 建置部署 | 4.0/10 | 嚴重滯後，無 CI/CD 與容器化 |

---

## 🎯 40 項改善建議

### 一、安全性改善 (Security) - 5 項

#### 1. 🔴 [高優先] 新增 API 認證中間件
**問題**: 所有 API 端點無權限檢查，任何人可存取  
**位置**: `internal/api/router.go`  
**建議**:
```go
router.Use(authMiddleware())  // JWT 或 API Key
apiV1.DELETE("/devices/:id", requireRole("admin"), deviceHandler.Delete)
```
**工作量**: 2-3 天

---

#### 2. 🔴 [高優先] 敏感資訊保護
**問題**: 連線配置（含密碼）直接序列化暴露  
**位置**: `cmd/test_ui/main.go`  
**建議**:
```go
type DeviceConnectionConfig struct {
    Password string `json:"-"` // 不序列化
}
```
**工作量**: 1 天

---

#### 3. 🟡 [中優先] SQL 注入防護審計
**問題**: 大部分使用參數化查詢，但需全面審計  
**位置**: `internal/datalink/*/sql_repo.go`  
**建議**: 執行 `go install github.com/securego/gosec/v2/cmd/gosec@latest && gosec ./...`  
**工作量**: 0.5 天

---

#### 4. 🟡 [中優先] 輸入驗證強化
**問題**: 缺少統一的輸入驗證框架  
**位置**: API handlers  
**建議**: 引入 `go-playground/validator` 進行結構化驗證  
**工作量**: 2 天

---

#### 5. 🟡 [中優先] CORS 配置收緊
**問題**: 開發環境允許 `*` 來源  
**位置**: `internal/api/middleware.go`  
**建議**: 生產環境限制為特定域名列表  
**工作量**: 0.5 天

---

### 二、效能優化 (Performance) - 6 項

#### 6. 🔴 [高優先] 啟動連線清理例程
**問題**: `CleanupRoutine` 未自動啟動，可能導致連線洩漏  
**位置**: `cmd/test_ui/main.go:86-87`  
**建議**:
```go
connMgr := connector.GetConnectionManager()
connMgr.StartCleanupRoutine(context.Background())
defer connMgr.CloseAll()
```
**工作量**: 0.5 天

---

#### 7. 🔴 [高優先] 新增資料庫索引
**問題**: 關鍵查詢無索引，每次 O(n) 掃描  
**位置**: `internal/datalink/schema/migrations/`  
**建議**:
```sql
CREATE INDEX idx_points_polling_group ON points(polling_group_id);
CREATE INDEX idx_mappings_point_tag ON mappings(point_id, tag_id);
CREATE INDEX idx_timeseries_tag_time ON timeseries(tag_id, timestamp DESC);
```
**工作量**: 1 天

---

#### 8. 🟡 [中優先] 排程器群組索引優化
**問題**: `pollGroup()` 每次輪詢線性掃描所有點位  
**位置**: `internal/datalink/collector/scheduler.go:340`  
**建議**:
```go
type Scheduler struct {
    groupPointIndex map[string][]string // groupID -> [pointID] 預建索引
}
```
**工作量**: 2 天

---

#### 9. 🟡 [中優先] 時序資料批次寫入
**問題**: 每筆資料單獨 INSERT，效能低落  
**位置**: `internal/datalink/storage/timeseries.go`  
**建議**: 實作批次插入 (Batch Insert)，每 100 筆或每秒寫入一次  
**工作量**: 2 天

---

#### 10. 🟡 [中優先] 資料庫連線池調整
**問題**: SQLite 預設 `MaxOpenConns=1`，PostgreSQL 預設 `10` 過低  
**位置**: `internal/datalink/db.go`  
**建議**: 根據硬體調整，PostgreSQL 生產環境建議 50-100  
**工作量**: 0.5 天

---

#### 11. 🟡 [中優先] 前端虛擬列表化
**問題**: `MemoryGrid` 大地址空間時所有 cell 都會渲染  
**位置**: `frontend/src/components/datalink/MemoryGrid.tsx`  
**建議**: 引入 `react-window` 實現虛擬化  
**工作量**: 2 天

---

### 三、錯誤處理 (Error Handling) - 5 項

#### 12. 🔴 [高優先] 統一錯誤類型
**問題**: 錯誤訊息不一致，SQL 錯誤未區分類型  
**位置**: 全專案  
**建議**:
```go
type APIError struct {
    Code    string `json:"code"`    // "entity_not_found", "conflict"
    Message string `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}
```
**工作量**: 3 天

---

#### 13. 🔴 [高優先] 前端全局 Error Boundary
**問題**: 無全局錯誤邊界，未捕獲錯誤導致白屏  
**位置**: `frontend/src/App.tsx`  
**建議**: 包裝 `ErrorBoundary` 組件  
**工作量**: 1 天

---

#### 14. 🟡 [中優先] 排程器重試指數退避
**問題**: 重試使用線性延遲，無指數退避  
**位置**: `internal/datalink/collector/scheduler.go`  
**建議**:
```go
delay := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
```
**工作量**: 0.5 天

---

#### 15. 🟡 [中優先] 背壓處理機制
**問題**: 值通道溢出時舊值被丟棄，可能遺漏關鍵數據  
**位置**: `internal/datalink/collector/scheduler.go`  
**建議**: 實現帶持久化的背壓隊列  
**工作量**: 2 天

---

#### 16. 🟡 [中優先] 前端統一錯誤提示
**問題**: 部分 catch 只 console.error，未通知用戶  
**位置**: `frontend/src/components/datalink/*.tsx`  
**建議**: 所有 catch 加入 `showError()` 呼叫  
**工作量**: 1 天

---

### 四、測試覆蓋 (Testing) - 5 項

#### 17. 🔴 [高優先] 後端整合測試
**問題**: 只有單元測試，無端對端整合測試  
**位置**: 新增 `internal/integration_test/`  
**建議**: 測試完整流程 Device → Point → Mapping → Storage  
**工作量**: 1 週

---

#### 18. 🔴 [高優先] 前端關鍵路徑測試
**問題**: 測試覆蓋率約 3-5%，核心 hooks/services 完全無測試  
**位置**: `frontend/src/hooks/datalink/*.ts`  
**建議**: 優先測試 `useDevices`, `useTags`, `useMappings`  
**工作量**: 1 週

---

#### 19. 🟡 [中優先] 協議契約測試
**問題**: 協議實現無合規性驗證  
**位置**: `internal/protocol/*`  
**建議**: 建立 Modbus/FATEK 協議契約測試  
**工作量**: 3 天

---

#### 20. 🟡 [中優先] 效能基準測試擴展
**問題**: 只有 `lib/hsllogic` 有基準測試  
**位置**: `internal/datalink/collector/`  
**建議**: 新增 `BenchmarkSchedulerPollGroup` 測試 10000+ 點位  
**工作量**: 2 天

---

#### 21. 🟡 [中優先] 前端 E2E 測試
**問題**: 無端對端用戶流程測試  
**位置**: 新增 `frontend/e2e/`  
**建議**: 使用 Playwright 測試關鍵工作流程  
**工作量**: 1 週

---

### 五、建置部署 (Build & Deploy) - 6 項

#### 22. 🔴 [高優先] 建立 CI/CD 流程
**問題**: 無 `.github/workflows`，無自動化測試與部署  
**位置**: 新增 `.github/workflows/ci.yml`  
**建議**: 實現 lint → test → build → artifact 流程  
**工作量**: 2 天

---

#### 23. 🔴 [高優先] 容器化支援
**問題**: 無 Dockerfile，無法容器化部署  
**位置**: 新增 `Dockerfile`, `docker-compose.yml`  
**建議**: 多階段構建，最小化映像體積  
**工作量**: 1 天

---

#### 24. 🔴 [高優先] 修復 Makefile 跨平台問題
**問題**: 使用 `rm` 命令，Windows 不原生支援  
**位置**: `Makefile`  
**建議**: 使用條件判斷實現跨平台命令  
**工作量**: 0.5 天

---

#### 25. 🟡 [中優先] 依賴安全掃描
**問題**: 無漏洞掃描工具集成  
**位置**: CI 流程  
**建議**: 
```bash
go install github.com/google/osv-scanner/cmd/osv-scanner@latest
npm audit --audit-level=moderate
```
**工作量**: 0.5 天

---

#### 26. 🟡 [中優先] Go 版本回退
**問題**: `go 1.25.5` 為 RC 版本，不穩定  
**位置**: `go.mod`  
**建議**: 回退至 `go 1.23.x` 穩定版  
**工作量**: 0.5 天

---

#### 27. 🟡 [中優先] React 版本評估
**問題**: React 19 剛發佈，存在未知風險  
**位置**: `frontend/package.json`  
**建議**: 評估是否回退至 React 18 LTS  
**工作量**: 0.5 天

---

### 六、架構設計 (Architecture) - 5 項

#### 28. 🔴 [高優先] 解決單例連線管理器問題
**問題**: `GetConnectionManager()` 使用 `sync.Once` 單例，難以測試  
**位置**: `internal/datalink/connector/manager.go`  
**建議**: 改為依賴注入，測試時可替換為 Mock  
**工作量**: 3 天

---

#### 29. 🟡 [中優先] 實現熔斷器模式
**問題**: 協議連線失敗無熔斷機制，持續重試浪費資源  
**位置**: 新增 `internal/datalink/connector/circuit_breaker.go`  
**建議**: 參考 `p1-implement-circuit-breaker` 提案  
**工作量**: 1 週

---

#### 30. 🟡 [中優先] 外掛式協議載入
**問題**: 協議硬編碼在程式碼中，新增協議需重新編譯  
**位置**: `internal/datalink/connector/registry.go`  
**建議**:
```go
type ProtocolPlugin interface {
    Version() string
    Init() error
    NewClient(config string) (Protocol, error)
}
```
**工作量**: 2 週

---

#### 31. 🟡 [中優先] API 限流機制
**問題**: 無 API 限流，可能被濫用  
**位置**: `internal/api/middleware.go`  
**建議**: 實現 Token Bucket 限流  
**工作量**: 2 天

---

#### 32. 🟢 [低優先] 分散式支援準備
**問題**: 單機架構，無法水平擴展  
**位置**: 架構設計  
**建議**: 設計集群支援，分片點位索引  
**工作量**: 2-4 週

---

### 七、程式碼品質 (Code Quality) - 4 項

#### 33. 🔴 [高優先] 修復技術債務 TODO
**問題**: 多處 TODO 未實現，包含關鍵邏輯  
**位置**:
- `tag/service.go:45` - 檢查映射使用此標籤
- `device/service.go:88` - Check Points
- `point_handler.go:156` - 實際呼叫協議連接器  
**工作量**: 3 天

---

#### 34. 🟡 [中優先] 消除 TypeScript any 類型
**問題**: 部分 any 類型讓類型檢查失效  
**位置**: `frontend/src/services/api.ts`, `DeviceForm.tsx`  
**建議**: 定義精確的聯合類型  
**工作量**: 1 天

---

#### 35. 🟡 [中優先] 前端 React.memo 優化
**問題**: 關鍵組件未使用 React.memo，可能重複渲染  
**位置**: `frontend/src/components/datalink/*.tsx`  
**建議**: 對純展示組件加上 React.memo  
**工作量**: 1 天

---

#### 36. 🟡 [中優先] 使用 react-hook-form
**問題**: 表單使用多個 useState，容易失同步  
**位置**: `PointForm.tsx`, `DeviceForm.tsx`  
**建議**: 引入 react-hook-form 統一表單管理  
**工作量**: 2 天

---

### 八、規格文檔 (Documentation) - 4 項

#### 37. 🔴 [高優先] 生成 OpenAPI 規範
**問題**: 無正式 API 文檔，前後端契約不明確  
**位置**: 新增 `docs/openapi.yaml`  
**建議**: 使用 swag 從 Go 註解自動生成  
**工作量**: 3 天

---

#### 38. 🟡 [中優先] 更新規格 Purpose 欄位
**問題**: 多個規格 Purpose 為 TBD  
**位置**: `openspec/specs/*/spec.md`  
**建議**: 填寫明確的規格目的說明  
**工作量**: 0.5 天

---

#### 39. 🟡 [中優先] 補齊缺失的變更提案
**問題**: `add-ops-tooling` 缺少 proposal.md  
**位置**: `openspec/changes/add-ops-tooling/`  
**建議**: 補寫 proposal.md 說明變更原因與影響  
**工作量**: 0.5 天

---

#### 40. 🟡 [中優先] 文檔結構重整
**問題**: `docs/` 下 80+ 個 HSL 文檔，信噪比低  
**位置**: `docs/`  
**建議**:
```
docs/
├── GETTING_STARTED.md
├── API_REFERENCE.md  
├── ARCHITECTURE.md
├── TROUBLESHOOTING.md
└── archive/HslDocs_*/
```
**工作量**: 1 天

---

## 📋 優先級總覽

| 優先級 | 數量 | 估計工時 |
|--------|------|---------|
| 🔴 高優先 | 15 項 | 3-4 週 |
| 🟡 中優先 | 23 項 | 6-8 週 |
| 🟢 低優先 | 2 項 | 2-4 週 |

---

## 🚀 建議實施路線圖

### 第 1-2 週：關鍵修復
- [ ] #1 新增 API 認證中間件
- [ ] #6 啟動連線清理例程
- [ ] #7 新增資料庫索引
- [ ] #12 統一錯誤類型
- [ ] #13 前端全局 Error Boundary
- [ ] #22 建立 CI/CD 流程
- [ ] #23 容器化支援
- [ ] #24 修復 Makefile 跨平台問題

### 第 3-4 週：基礎設施
- [ ] #2 敏感資訊保護
- [ ] #17 後端整合測試
- [ ] #18 前端關鍵路徑測試
- [ ] #28 解決單例連線管理器問題
- [ ] #33 修復技術債務 TODO
- [ ] #37 生成 OpenAPI 規範

### 第 5-8 週：優化完善
- [ ] #8 排程器群組索引優化
- [ ] #9 時序資料批次寫入
- [ ] #11 前端虛擬列表化
- [ ] #19 協議契約測試
- [ ] #21 前端 E2E 測試
- [ ] #29 實現熔斷器模式
- [ ] #31 API 限流機制

---

## 📈 預期效益

| 改善領域 | 改善前 | 改善後 | 提升幅度 |
|---------|--------|--------|---------|
| 安全性 | 5/10 | 8/10 | +60% |
| 效能 | 6/10 | 8.5/10 | +42% |
| 測試覆蓋 | 10% | 60%+ | +500% |
| 部署效率 | 2/10 | 9/10 | +350% |
| 文檔品質 | 7/10 | 9/10 | +29% |

---

**報告產生時間**: 2026-01-31 17:59 UTC
**分析工具**: Claude Code + OpenSpec
