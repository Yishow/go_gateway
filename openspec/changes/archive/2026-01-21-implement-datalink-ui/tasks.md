## 1. 後端 API 完善

### 1.1 新增 Polling Groups 模組

- [x] 1.1.1 建立 `internal/datalink/pollinggroup/service.go`
- [x] 1.1.2 建立 `internal/datalink/pollinggroup/repository.go`（介面定義）
- [x] 1.1.3 建立 `internal/datalink/pollinggroup/memory_repo.go`
- [x] 1.1.4 建立 `internal/api/handlers/polling_group_handler.go`
- [x] 1.1.5 在 `router.go` 註冊路由

### 1.2 新增 Protocols API

- [x] 1.2.1 建立 `internal/api/handlers/protocol_handler.go`
- [x] 1.2.2 實現協議列表與 JSON Schema 返回
- [x] 1.2.3 在 `router.go` 註冊路由

### 1.3 新增 Health API

- [x] 1.3.1 建立 `internal/api/handlers/datalink_health_handler.go`
- [x] 1.3.2 在 `router.go` 註冊路由

### 1.4 擴展 Tag Handler

- [x] 1.4.1 實現 `POST /tags/:id/activate`
- [x] 1.4.2 實現 `POST /tags/:id/retire`
- [x] 1.4.3 實現 `POST /tags/batch`
- [x] 1.4.4 實現 `POST /tags/validate-key`
- [x] 1.4.5 更新 `router.go` 註冊新路由

### 1.5 擴展 Device Handler

- [x] 1.5.1 實現 `POST /devices/test-batch`
- [x] 1.5.2 修復 `Activate` 方法（目前空殼）
- [x] 1.5.3 修復 `Disable` 方法（目前空殼）
- [x] 1.5.4 更新 `router.go` 註冊新路由

### 1.6 擴展 Point Handler

- [x] 1.6.1 實現 `POST /points/:id/poll`
- [x] 1.6.2 實現 `POST /points/poll`（批量）
- [x] 1.6.3 更新 `router.go` 註冊新路由

### 1.7 擴展 Mapping Handler

- [x] 1.7.1 實現 `POST /mappings/validate-pipeline`
- [x] 1.7.2 更新 `router.go` 註冊新路由

### 1.8 重構 Settings Handler

- [x] 1.8.1 建立 `internal/datalink/settings/service.go`
- [x] 1.8.2 建立 `internal/datalink/settings/repository.go`（整合於 service.go）
- [x] 1.8.3 重構 `settings_handler.go` 連接 Service
- [x] 1.8.4 移除硬編碼 Mock

## 2. 儲存層持久化

### 2.1 資料庫工廠

- [x] 2.1.1 建立 `internal/datalink/db.go`（連接管理）
- [x] 2.1.2 實現 Repository 工廠模式
- [x] 2.1.3 整合現有 migrations 至啟動流程

### 2.2 Device SQL Repository

- [x] 2.2.1 建立 `internal/datalink/device/sql_repo.go`
- [x] 2.2.2 實現 CRUD 操作（SQLite + Postgres 相容）
- [x] 2.2.3 撰寫單元測試

### 2.3 Point SQL Repository

- [x] 2.3.1 建立 `internal/datalink/point/sql_repo.go`
- [x] 2.3.2 實現 CRUD 操作
- [x] 2.3.3 撰寫單元測試

### 2.4 Tag SQL Repository

- [x] 2.4.1 建立 `internal/datalink/tag/sql_repo.go`
- [x] 2.4.2 實現 CRUD 操作
- [x] 2.4.3 撰寫單元測試

### 2.5 Mapping SQL Repository

- [x] 2.5.1 建立 `internal/datalink/mapping/sql_repo.go`
- [x] 2.5.2 實現 CRUD 操作
- [x] 2.5.3 撰寫單元測試

### 2.6 Polling Group SQL Repository

- [x] 2.6.1 建立 `internal/datalink/pollinggroup/sql_repo.go`
- [x] 2.6.2 實現 CRUD 操作
- [x] 2.6.3 撰寫單元測試

### 2.7 Settings SQL Repository

- [x] 2.7.1 建立 `internal/datalink/settings/sql_repo.go`
- [x] 2.7.2 實現 CRUD 操作
- [x] 2.7.3 撰寫單元測試

## 3. SSE 即時預覽

### 3.1 後端 SSE 端點

- [x] 3.1.1 建立 `internal/api/handlers/datalink_sse_handler.go`
- [x] 3.1.2 實現 `GET /datalink/preview/stream`
- [x] 3.1.3 整合 Mapping Preview 邏輯
- [x] 3.1.4 在 `router.go` 註冊路由

### 3.2 前端 SSE 客戶端

- [x] 3.2.1 建立 `web/test-ui/src/hooks/usePreviewStream.ts`
- [x] 3.2.2 實現 SSE 連接與自動重連
- [x] 3.2.3 整合至 LivePreviewPanel 組件

## 4. 導引式工作流程 UI

### 4.1 Stepperize 整合

- [x] 4.1.1 安裝 `@stepperize/react`（改用自訂實現）
- [x] 4.1.2 建立 `web/test-ui/src/components/datalink/wizard/MappingWizard.tsx`
- [x] 4.1.3 定義 6 步驟結構與驗證 schema

### 4.2 步驟組件實現

- [x] 4.2.1 建立 `DeviceStep.tsx`（步驟 1）
- [x] 4.2.2 建立 `PointStep.tsx`（步驟 2）
- [x] 4.2.3 建立 `TagStep.tsx`（步驟 3）
- [x] 4.2.4 建立 `TransformStep.tsx`（步驟 4，重用 TransformBuilder）
- [x] 4.2.5 建立 `PreviewStep.tsx`（步驟 5，整合 SSE）
- [x] 4.2.6 建立 `CompleteStep.tsx`（步驟 6）

### 4.3 共用組件

- [x] 4.3.1 建立 `StepIndicator.tsx`
- [x] 4.3.2 建立 `LivePreviewPanel.tsx`
- [x] 4.3.3 整合至 MappingWizard 佈局

### 4.4 路由與導航

- [x] 4.4.1 新增 `/datalink/wizard` 路由
- [x] 4.4.2 從現有頁面新增入口連結（Dashboard 卡片、MappingsPage 按鈕）

## 5. 驗證與測試

### 5.1 後端測試

### 5.1 後端測試

- [x] 5.1.1 執行現有 Datalink API 測試
- [x] 5.1.2 新增 Polling Groups API 測試
- [x] 5.1.3 新增 Protocols API 測試
- [x] 5.1.4 新增 Health API 測試
- [x] 5.1.5 新增擴展 API 測試

### 5.2 整合測試

- [x] 5.2.1 SQLite Repository 整合測試
- [x] 5.2.2 Postgres Repository 整合測試（如有環境）

### 5.3 前端驗證

- [x] 5.3.1 手動驗證 Wizard 完整流程
- [x] 5.3.2 驗證 SSE 即時預覽功能
- [x] 5.3.3 驗證草稿儲存與恢復
