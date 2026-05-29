## Why

當前閘道器缺乏一個直觀、實時的運作監控儀表板。用戶無法在單一介面中掌握各 PLC 設備的即時連線健康度、輪詢調度績效、資料庫/訊息佇列的寫入進度與即時錯誤日誌。引入 Runtime Dashboard 將提供實時數據監控介面，顯著提升維運與故障排障的效率。

## What Changes

- 在前端新增 `/studio/dashboard` 路由，提供實時運行儀表板。
- 在 Workbench V2 的 Shell TopBar 新增「監控面板 (Dashboard)」導覽入口。
- 新增即時指標元件，顯示採集吞吐量、輪詢延遲、丟包率與資料庫佇列積壓狀態。
- 新增 PLC 設備狀態格格，直觀展示所有連線中的設備狀態。
- 新增即時點位數值表，顯示各 Tag 的原始讀取值、轉換值與更新時間。
- 新增錯誤診斷與日誌面板，顯示最近的採集錯誤（如 TCP 逾時、Modbus 異常碼）。

## Capabilities

### New Capabilities

- `runtime-dashboard`: 定義 Runtime Dashboard 的版面結構、與 SSE/API 即時數據源的對接協議，以及各元件（指標、設備狀態、點位值、錯誤日誌）的互動規格。

### Modified Capabilities

(none)

## Impact

- Affected specs: `runtime-dashboard`
- Affected code:
  - New:
    - `frontend/src/features/datalink/workbench-v2/dashboard/RuntimeDashboard.tsx`
    - `frontend/src/features/datalink/workbench-v2/dashboard/components/MetricsOverview.tsx`
    - `frontend/src/features/datalink/workbench-v2/dashboard/components/DeviceStatusGrid.tsx`
    - `frontend/src/features/datalink/workbench-v2/dashboard/components/LiveValuesTable.tsx`
    - `frontend/src/features/datalink/workbench-v2/dashboard/components/OutputQueueStatus.tsx`
    - `frontend/src/features/datalink/workbench-v2/dashboard/components/RealtimeLogs.tsx`
    - `frontend/tests/unit/workbench-v2/runtime-dashboard.test.tsx`
    - `frontend/src/i18n/locales/en/runtime-dashboard.json`
    - `frontend/src/i18n/locales/zh-TW/runtime-dashboard.json`
  - Modified:
    - `frontend/src/App.tsx`
    - `frontend/src/i18n/config.ts`
    - `frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx`
