## 1. 後端 API 完善

### 1.1 新增 Polling Groups 模組

- [ ] 1.1.1 建立 `internal/datalink/pollinggroup/service.go`
- [ ] 1.1.2 建立 `internal/datalink/pollinggroup/repository.go`（介面定義）
- [ ] 1.1.3 建立 `internal/datalink/pollinggroup/memory_repo.go`
- [ ] 1.1.4 建立 `internal/api/handlers/polling_group_handler.go`
- [ ] 1.1.5 在 `router.go` 註冊路由

### 1.2 新增 Protocols API

- [ ] 1.2.1 建立 `internal/api/handlers/protocol_handler.go`
- [ ] 1.2.2 實現協議列表與 JSON Schema 返回
- [ ] 1.2.3 在 `router.go` 註冊路由

### 1.3 新增 Health API

- [ ] 1.3.1 建立 `internal/api/handlers/health_handler.go`
- [ ] 1.3.2 在 `router.go` 註冊路由

### 1.4 擴展 Tag Handler

- [ ] 1.4.1 實現 `POST /tags/:id/activate`
- [ ] 1.4.2 實現 `POST /tags/:id/retire`
- [ ] 1.4.3 實現 `POST /tags/batch`
- [ ] 1.4.4 實現 `POST /tags/validate-key`
- [ ] 1.4.5 更新 `router.go` 註冊新路由

### 1.5 擴展 Device Handler

- [ ] 1.5.1 實現 `POST /devices/test-batch`
- [ ] 1.5.2 修復 `Activate` 方法（目前空殼）
- [ ] 1.5.3 修復 `Disable` 方法（目前空殼）
- [ ] 1.5.4 更新 `router.go` 註冊新路由

### 1.6 擴展 Point Handler

- [ ] 1.6.1 實現 `POST /points/:id/poll`
- [ ] 1.6.2 實現 `POST /points/poll`（批量）
- [ ] 1.6.3 更新 `router.go` 註冊新路由

### 1.7 擴展 Mapping Handler

- [ ] 1.7.1 實現 `POST /mappings/validate-pipeline`
- [ ] 1.7.2 更新 `router.go` 註冊新路由

### 1.8 重構 Settings Handler

- [ ] 1.8.1 建立 `internal/datalink/settings/service.go`
- [ ] 1.8.2 建立 `internal/datalink/settings/repository.go`
- [ ] 1.8.3 重構 `settings_handler.go` 連接 Service
- [ ] 1.8.4 移除硬編碼 Mock

## 2. 儲存層持久化

### 2.1 資料庫工廠

- [ ] 2.1.1 建立 `internal/datalink/db.go`（連接管理）
- [ ] 2.1.2 實現 Repository 工廠模式
- [ ] 2.1.3 整合現有 migrations 至啟動流程

### 2.2 Device SQL Repository

- [ ] 2.2.1 建立 `internal/datalink/device/sql_repo.go`
- [ ] 2.2.2 實現 CRUD 操作（SQLite + Postgres 相容）
- [ ] 2.2.3 撰寫單元測試

### 2.3 Point SQL Repository

- [ ] 2.3.1 建立 `internal/datalink/point/sql_repo.go`
- [ ] 2.3.2 實現 CRUD 操作
- [ ] 2.3.3 撰寫單元測試

### 2.4 Tag SQL Repository

- [ ] 2.4.1 建立 `internal/datalink/tag/sql_repo.go`
- [ ] 2.4.2 實現 CRUD 操作
- [ ] 2.4.3 撰寫單元測試

### 2.5 Mapping SQL Repository

- [ ] 2.5.1 建立 `internal/datalink/mapping/sql_repo.go`
- [ ] 2.5.2 實現 CRUD 操作
- [ ] 2.5.3 撰寫單元測試

### 2.6 Polling Group SQL Repository

- [ ] 2.6.1 建立 `internal/datalink/pollinggroup/sql_repo.go`
- [ ] 2.6.2 實現 CRUD 操作
- [ ] 2.6.3 撰寫單元測試

### 2.7 Settings SQL Repository

- [ ] 2.7.1 建立 `internal/datalink/settings/sql_repo.go`
- [ ] 2.7.2 實現 CRUD 操作
- [ ] 2.7.3 撰寫單元測試

## 3. SSE 即時預覽

### 3.1 後端 SSE 端點

- [ ] 3.1.1 建立 `internal/api/handlers/datalink_sse.go`
- [ ] 3.1.2 實現 `GET /datalink/preview/stream`
- [ ] 3.1.3 整合 Mapping Preview 邏輯
- [ ] 3.1.4 在 `router.go` 註冊路由

### 3.2 前端 SSE 客戶端

- [ ] 3.2.1 建立 `web/test-ui/src/hooks/usePreviewStream.ts`
- [ ] 3.2.2 實現 SSE 連接與自動重連
- [ ] 3.2.3 整合至 LivePreviewPanel 組件

## 4. 導引式工作流程 UI

### 4.1 Stepperize 整合

- [ ] 4.1.1 安裝 `@stepperize/react`
- [ ] 4.1.2 建立 `web/test-ui/src/components/datalink/wizard/MappingWizard.tsx`
- [ ] 4.1.3 定義 6 步驟結構與驗證 schema

### 4.2 步驟組件實現

- [ ] 4.2.1 建立 `DeviceStep.tsx`（步驟 1）
- [ ] 4.2.2 建立 `PointStep.tsx`（步驟 2）
- [ ] 4.2.3 建立 `TagStep.tsx`（步驟 3）
- [ ] 4.2.4 建立 `TransformStep.tsx`（步驟 4，重用 TransformBuilder）
- [ ] 4.2.5 建立 `PreviewStep.tsx`（步驟 5，整合 SSE）
- [ ] 4.2.6 建立 `CompleteStep.tsx`（步驟 6）

### 4.3 共用組件

- [ ] 4.3.1 建立 `StepIndicator.tsx`
- [ ] 4.3.2 建立 `LivePreviewPanel.tsx`
- [ ] 4.3.3 整合至 MappingWizard 佈局

### 4.4 路由與導航

- [ ] 4.4.1 新增 `/datalink/wizard` 路由
- [ ] 4.4.2 從現有頁面新增入口連結

## 5. 驗證與測試

### 5.1 後端測試

- [ ] 5.1.1 執行現有 Datalink API 測試
- [ ] 5.1.2 新增 Polling Groups API 測試
- [ ] 5.1.3 新增 Protocols API 測試
- [ ] 5.1.4 新增 Health API 測試
- [ ] 5.1.5 新增擴展 API 測試

### 5.2 整合測試

- [ ] 5.2.1 SQLite Repository 整合測試
- [ ] 5.2.2 Postgres Repository 整合測試（如有環境）

### 5.3 前端驗證

- [ ] 5.3.1 手動驗證 Wizard 完整流程
- [ ] 5.3.2 驗證 SSE 即時預覽功能
- [ ] 5.3.3 驗證草稿儲存與恢復
