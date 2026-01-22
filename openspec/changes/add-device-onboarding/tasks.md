# Tasks: Add Device Onboarding and Readiness Check

## 1. 資料模型擴充（TDD）

- [ ] **RED**: 撰寫 Device 模型擴充測試
  - [ ] 測試 `last_collected_at` 欄位
  - [ ] 測試 `collection_count` 欄位
  - [ ] 測試 `error_count` 欄位
  - [ ] 測試 `readiness_status` 欄位
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作 Device 模型擴充
  - [ ] 擴充 `internal/datalink/schema/models.go` 中的 Device 結構
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化程式碼
  - [ ] 檢查命名和結構
  - [ ] 確保測試仍通過

## 2. 資料庫遷移（TDD）

- [ ] **RED**: 撰寫遷移測試
  - [ ] 測試新增欄位的 SQL 語句
  - [ ] 測試向下相容性（現有資料不受影響）
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作資料庫遷移
  - [ ] 建立遷移檔案 `002_add_device_collection_stats.up.sql`
  - [ ] 建立回滾檔案 `002_add_device_collection_stats.down.sql`
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化遷移腳本
  - [ ] 檢查索引和約束
  - [ ] 確保遷移可重複執行

## 3. Readiness Check Service（TDD）

- [ ] **RED**: 撰寫 Readiness Check 單元測試
  - [ ] 測試設備狀態檢查（draft/active/disabled）
  - [ ] 測試點位存在性檢查
  - [ ] 測試輪詢群組指派檢查
  - [ ] 測試映射存在性檢查
  - [ ] 測試 Scheduler 運行狀態檢查
  - [ ] 測試完整配置場景（所有檢查通過）
  - [ ] 測試部分配置場景（部分檢查失敗）
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作 Readiness Check 服務
  - [ ] 在 `internal/datalink/device/service.go` 新增 `CheckReadiness()` 方法
  - [ ] 實作各項檢查邏輯
  - [ ] 實作建議生成邏輯
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化程式碼
  - [ ] 提取重複邏輯
  - [ ] 優化查詢效能（避免 N+1 查詢）
  - [ ] 確保測試仍通過

## 4. Readiness Check API（TDD）

- [ ] **RED**: 撰寫 API 端點測試
  - [ ] 測試 `GET /devices/:id/readiness` 成功場景
  - [ ] 測試設備不存在場景（404）
  - [ ] 測試回應格式驗證
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作 API 端點
  - [ ] 在 `internal/datalink/api/device_handler.go` 新增 `GetReadiness()` 處理器
  - [ ] 在 `internal/datalink/api/router.go` 註冊路由
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化 API 回應
  - [ ] 統一錯誤處理格式
  - [ ] 確保測試仍通過

## 5. 收集統計更新機制（TDD）

- [ ] **RED**: 撰寫統計更新測試
  - [ ] 測試 `UpdateCollectionStats()` 方法
  - [ ] 測試 `last_collected_at` 更新
  - [ ] 測試 `collection_count` 遞增
  - [ ] 測試 `error_count` 遞增
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作統計更新機制
  - [ ] 在 `internal/datalink/device/service.go` 新增 `UpdateCollectionStats()` 方法
  - [ ] 在 Scheduler 中整合統計更新
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化更新邏輯
  - [ ] 使用批次更新減少資料庫操作
  - [ ] 確保測試仍通過

## 6. Onboarding Wizard 前端組件（TDD）

- [ ] **RED**: 撰寫組件測試
  - [ ] 測試組件渲染
  - [ ] 測試步驟導航（下一步/上一步）
  - [ ] 測試表單驗證
  - [ ] 測試 API 整合
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作 Onboarding Wizard
  - [ ] 建立 `frontend/src/components/datalink/DeviceOnboardingWizard.tsx`
  - [ ] 實作步驟流程（7 步驟）
  - [ ] 整合現有 API 和組件
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化組件結構
  - [ ] 提取可重用邏輯
  - [ ] 優化狀態管理
  - [ ] 確保測試仍通過

## 7. 設備狀態儀表板（TDD）

- [ ] **RED**: 撰寫儀表板測試
  - [ ] 測試統計資訊顯示
  - [ ] 測試設備列表渲染
  - [ ] 測試即時更新機制
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作 Status Dashboard
  - [ ] 建立 `frontend/src/pages/datalink/DeviceStatusDashboard.tsx`
  - [ ] 實作統計資訊顯示
  - [ ] 實作設備狀態列表
  - [ ] 整合 WebSocket/SSE 即時更新
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化儀表板效能
  - [ ] 優化資料查詢
  - [ ] 實作虛擬滾動（如需要）
  - [ ] 確保測試仍通過

## 8. 整合測試（TDD）

- [ ] **RED**: 撰寫端到端測試
  - [ ] 測試完整引導流程
  - [ ] 測試 Readiness Check 整合
  - [ ] 測試狀態儀表板資料流
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作整合測試
  - [ ] 建立整合測試檔案
  - [ ] 實作測試場景
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化測試結構
  - [ ] 提取測試輔助函數
  - [ ] 確保測試仍通過

## 9. 文件更新

- [ ] 更新 API 文件
- [ ] 更新使用者文件
- [ ] 更新開發者文件

## 10. 驗證與部署

- [ ] 執行完整測試套件
- [ ] 執行效能測試
- [ ] 執行安全性檢查
- [ ] 部署到測試環境
- [ ] 驗證功能正常運作

## 11. Sidebar 導航改進（TDD）

- [ ] **RED**: 撰寫導航測試
  - [ ] 測試 Points 導航項目顯示
  - [ ] 測試導航項目排序
  - [ ] 測試響應式 Sidebar（行動端）
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作 Sidebar 改進
  - [ ] 在 `frontend/src/layouts/DatalinkLayout.tsx` 新增 Points 導航項目
  - [ ] 重新排序導航項目（按照資料管線流程）
  - [ ] 新增視覺分組（分隔線或標題）
  - [ ] 實作響應式 Sidebar（行動端抽屜式）
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化導航結構
  - [ ] 統一圖標風格
  - [ ] 優化 hover 和 active 狀態
  - [ ] 確保測試仍通過

## 12. Points 頁面實作（TDD）

- [ ] **RED**: 撰寫 PointsPage 組件測試
  - [ ] 測試頁面渲染
  - [ ] 測試點位列表顯示
  - [ ] 測試篩選功能（設備、狀態、資料型別）
  - [ ] 測試排序功能
  - [ ] 測試 CRUD 操作（建立、編輯、刪除）
  - [ ] 測試批量操作
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作 PointsPage 組件
  - [ ] 建立 `frontend/src/pages/datalink/PointsPage.tsx`
  - [ ] 參考 `TagsPage.tsx` 的設計模式
  - [ ] 實作表格列表顯示
  - [ ] 實作篩選和排序功能
  - [ ] 實作 CRUD 操作（使用現有 API）
  - [ ] 實作批量操作（enable/disable, assign group, delete）
  - [ ] 整合 `PointForm` 組件（如需要）
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化頁面效能
  - [ ] 實作虛擬滾動（如點位數量多）
  - [ ] 優化查詢效能
  - [ ] 確保測試仍通過

## 13. 路由配置更新

- [ ] 更新 `frontend/src/App.tsx` 路由配置
  - [ ] 確保 `/datalink/points` 路由正確配置
  - [ ] 移除 PointsPage 佔位組件
  - [ ] 驗證路由導航正常運作
