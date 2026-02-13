# Change: Add Device Onboarding and Readiness Check

## Why

根據深度架構分析，當前系統存在以下關鍵問題：

1. **缺少使用者引導**：新增設備後，系統沒有明確指引下一步該做什麼。使用者不知道需要完成哪些步驟（建立點位、標籤、映射）才能開始收集資料。

2. **缺少狀態驗證**：設備啟用時只檢查連線，不檢查配置完整性（點位、群組、映射）。導致設備可能處於 `active` 狀態但實際上無法收集資料。

3. **缺少即時反饋**：無法即時查看設備是否正在收集資料，缺少「最後收集時間」和「收集統計」資訊。

這些問題導致：
- 使用者體驗差：配置流程不清晰，容易遺漏步驟
- 系統可靠性低：配置錯誤無法及時發現
- 運維效率低：無法快速診斷問題

## What Changes

### 1. 設備配置完成度檢查（Readiness Check）

新增 `DeviceReadiness` 功能，檢查設備是否已完整配置並可開始收集資料：

- **API 端點**：`GET /api/v1/datalink/devices/:id/readiness`
- **檢查項目**：
  - ✅ 設備狀態 = `active`
  - ✅ 至少一個 `enabled` 點位
  - ✅ 點位指派到 `enabled` 輪詢群組
  - ✅ 點位有對應的 `enabled` 映射
  - ✅ Scheduler 正在運行
- **回傳資訊**：完成度狀態、缺少的配置項目、建議的下一步操作

### 2. 新增設備引導流程（Onboarding Wizard）

建立 `DeviceOnboardingWizard` 前端組件，引導使用者完成所有必要配置：

- **觸發時機**：新增設備成功後自動開啟
- **步驟流程**：
  1. 建立設備 ✅
  2. 測試連線 → 啟用設備
  3. 建立點位（支援批次）
  4. 指派輪詢群組
  5. 建立標籤（支援批次）
  6. 建立映射（整合現有 MappingWizard）
  7. 確認開始收集
- **完成度顯示**：每個步驟顯示完成狀態和驗證結果

### 3. 設備狀態儀表板（Status Dashboard）

新增 `DeviceStatusDashboard` 頁面，顯示設備收集狀態和配置完成度：

- **統計資訊**：
  - 設備總數、啟用數、收集中數
  - 配置完成度統計
- **設備狀態**：
  - 最後收集時間
  - 收集次數、錯誤次數
  - 資料品質指標
- **即時監控**：使用 WebSocket/SSE 推送即時狀態更新

### 4. 資料模型擴充

擴充 `Device` 模型，新增收集統計欄位：

- `last_collected_at` - 最後收集時間
- `collection_count` - 收集次數
- `error_count` - 錯誤次數
- `readiness_status` - 配置完成度狀態

### 5. Sidebar 導航改進與 Points 頁面

改善導航結構和使用者體驗：

- **新增 Points 導航項目**：在 Sidebar 中新增 Points 頁面入口
- **實作 PointsPage 組件**：完整的點位管理頁面（CRUD + 篩選）
- **優化導航排序**：按照資料管線流程重新排序（Device → Point → Tag → Mapping）
- **視覺分組**：使用分隔線或標題分組相關功能
- **響應式優化**：行動端抽屜式 Sidebar，觸控優化

## Impact

- **Affected specs**: `device-registry`, `datalink-api`, `datalink-ui`
- **Affected code**:
  - `internal/datalink/device/service.go` - 新增 `CheckReadiness()` 方法
  - `internal/datalink/api/device_handler.go` - 新增 readiness 端點
  - `internal/datalink/schema/models.go` - 擴充 Device 模型
  - `internal/datalink/schema/migrations/` - 新增資料庫遷移
  - `frontend/src/pages/datalink/DevicesPage.tsx` - 整合引導流程
  - `frontend/src/components/datalink/DeviceOnboardingWizard.tsx` - 新增組件
  - `frontend/src/pages/datalink/DeviceStatusDashboard.tsx` - 新增頁面
  - `frontend/src/layouts/DatalinkLayout.tsx` - 更新 Sidebar 導航
  - `frontend/src/pages/datalink/PointsPage.tsx` - 新增點位管理頁面
  - `frontend/src/App.tsx` - 更新路由配置

## Risks

- **資料庫遷移**：需要新增欄位，需確保向下相容
- **效能影響**：Readiness Check 需要查詢多個表，需優化查詢效能
- **前端複雜度**：Onboarding Wizard 增加前端複雜度，需確保可維護性

## Migration Plan

1. **階段一**：實作 Readiness Check API 和資料模型擴充
2. **階段二**：實作 Onboarding Wizard 前端組件
3. **階段三**：實作 Status Dashboard 和即時監控
4. **階段四**：整合測試和文件更新
