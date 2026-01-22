# Tasks: Add Device Onboarding and Readiness Check

## 1. 資料模型擴充（TDD）

- [x] **RED**: 撰寫 Device 模型擴充測試
  - [x] 測試 `last_collected_at` 欄位
  - [x] 測試 `collection_count` 欄位
  - [x] 測試 `error_count` 欄位
  - [x] 測試 `readiness_status` 欄位
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 Device 模型擴充
  - [x] 擴充 `internal/datalink/schema/models.go` 中的 Device 結構
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化程式碼
  - [x] 檢查命名和結構
  - [x] 確保測試仍通過

## 2. 資料庫遷移（TDD）

- [x] **RED**: 撰寫遷移測試
  - [x] 測試新增欄位的 SQL 語句
  - [x] 測試向下相容性（現有資料不受影響）
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作資料庫遷移
  - [x] 建立遷移檔案 `002_add_device_collection_stats.up.sql`
  - [x] 建立回滾檔案 `002_add_device_collection_stats.down.sql`
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化遷移腳本
  - [x] 檢查索引和約束
  - [x] 確保遷移可重複執行

## 3. Readiness Check Service（TDD）

- [x] **RED**: 撰寫 Readiness Check 單元測試
  - [x] 測試設備狀態檢查（draft/active/disabled）
  - [x] 測試點位存在性檢查
  - [x] 測試輪詢群組指派檢查
  - [x] 測試映射存在性檢查
  - [x] 測試 Scheduler 運行狀態檢查
  - [x] 測試完整配置場景（所有檢查通過）
  - [x] 測試部分配置場景（部分檢查失敗）
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 Readiness Check 服務
  - [x] 在 `internal/datalink/device/service.go` 新增 `CheckReadiness()` 方法
  - [x] 實作各項檢查邏輯
  - [x] 實作建議生成邏輯
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化程式碼
  - [x] 提取重複邏輯
  - [x] 優化查詢效能（避免 N+1 查詢）
  - [x] 確保測試仍通過

## 4. Readiness Check API（TDD）

- [x] **RED**: 撰寫 API 端點測試
  - [x] 測試 `GET /devices/:id/readiness` 成功場景
  - [x] 測試設備不存在場景（404）
  - [x] 測試回應格式驗證
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 API 端點
  - [x] 在 `internal/datalink/api/device_handler.go` 新增 `GetReadiness()` 處理器
  - [x] 在 `internal/datalink/api/router.go` 註冊路由
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化 API 回應
  - [x] 統一錯誤處理格式
  - [x] 確保測試仍通過

## 5. 收集統計更新機制（TDD）

- [x] **RED**: 撰寫統計更新測試
  - [x] 測試 `UpdateCollectionStats()` 方法
  - [x] 測試 `last_collected_at` 更新
  - [x] 測試 `collection_count` 遞增
  - [x] 測試 `error_count` 遞增
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作統計更新機制
  - [x] 在 `internal/datalink/device/service.go` 新增 `UpdateCollectionStats()` 方法
  - [x] 在 Scheduler 中整合統計更新
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化更新邏輯
  - [x] 使用批次更新減少資料庫操作
  - [x] 確保測試仍通過

## 6. Onboarding Wizard 前端組件（TDD）

- [x] **RED**: 撰寫組件測試
  - [x] 測試組件渲染
  - [x] 測試步驟導航（下一步/上一步）
  - [x] 測試表單驗證
  - [x] 測試 API 整合
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 Onboarding Wizard
  - [x] 建立 `frontend/src/components/datalink/DeviceOnboardingWizard.tsx`
  - [x] 實作步驟流程（7 步驟）
  - [x] 整合現有 API 和組件
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化組件結構
  - [x] 提取可重用邏輯
  - [x] 優化狀態管理
  - [x] 確保測試仍通過

## 7. 設備狀態儀表板（TDD）

- [x] **RED**: 撰寫儀表板測試
  - [x] 測試統計資訊顯示
  - [x] 測試設備列表渲染
  - [x] 測試即時更新機制
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 Status Dashboard
  - [x] 建立 `frontend/src/pages/datalink/DeviceStatusDashboard.tsx`
  - [x] 實作統計資訊顯示
  - [x] 實作設備狀態列表
  - [x] 整合 WebSocket/SSE 即時更新
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化儀表板效能
  - [x] 優化資料查詢
  - [x] 實作虛擬滾動（如需要）
  - [x] 確保測試仍通過

## 8. 整合測試（TDD）

- [x] **RED**: 撰寫端到端測試
  - [x] 測試完整引導流程
  - [x] 測試 Readiness Check 整合
  - [x] 測試狀態儀表板資料流
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作整合測試
  - [x] 建立整合測試檔案
  - [x] 實作測試場景
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化測試結構
  - [x] 提取測試輔助函數
  - [x] 確保測試仍通過

## 9. 文件更新

- [x] 更新 API 文件
- [x] 更新使用者文件
- [x] 更新開發者文件

## 10. 驗證與部署

- [x] 執行完整測試套件
- [x] 執行效能測試
- [x] 執行安全性檢查
- [x] 部署到測試環境
- [x] 驗證功能正常運作

## 11. Sidebar 導航改進（TDD）

- [x] **RED**: 撰寫導航測試
  - [x] 測試 Points 導航項目顯示
  - [x] 測試導航項目排序
  - [x] 測試響應式 Sidebar（行動端）
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 Sidebar 改進
  - [x] 在 `frontend/src/layouts/DatalinkLayout.tsx` 新增 Points 導航項目
  - [x] 重新排序導航項目（按照資料管線流程）
  - [x] 新增視覺分組（分隔線或標題）
  - [x] 實作響應式 Sidebar（行動端抽屜式）
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化導航結構
  - [x] 統一圖標風格
  - [x] 優化 hover 和 active 狀態
  - [x] 確保測試仍通過

## 12. Points 頁面實作（TDD）

- [x] **RED**: 撰寫 PointsPage 組件測試
  - [x] 測試頁面渲染
  - [x] 測試點位列表顯示
  - [x] 測試篩選功能（設備、狀態、資料型別）
  - [x] 測試排序功能
  - [x] 測試 CRUD 操作（建立、編輯、刪除）
  - [x] 測試批量操作
  - [x] 執行測試確認失敗（RED）

- [x] **GREEN**: 實作 PointsPage 組件
  - [x] 建立 `frontend/src/pages/datalink/PointsPage.tsx`
  - [x] 參考 `TagsPage.tsx` 的設計模式
  - [x] 實作表格列表顯示
  - [x] 實作篩選和排序功能
  - [x] 實作 CRUD 操作（使用現有 API）
  - [x] 實作批量操作（enable/disable, assign group, delete）
  - [x] 整合 `PointForm` 組件（如需要）
  - [x] 執行測試確認通過（GREEN）

- [x] **REFACTOR**: 優化頁面效能
  - [x] 實作虛擬滾動（如點位數量多）
  - [x] 優化查詢效能
  - [x] 確保測試仍通過

## 13. 路由配置更新

- [x] 更新 `frontend/src/App.tsx` 路由配置
  - [x] 確保 `/datalink/points` 路由正確配置
  - [x] 移除 PointsPage 佔位組件
  - [x] 驗證路由導航正常運作
