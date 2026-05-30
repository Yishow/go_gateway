## Why

目前 `studio/v2` 已大致對齊 `new_prototype` 的 setup shell，但舊的 `runtime-dashboard` change 把 route、TopBar 導覽、fleet-ish 監控面板與尚未穩定的資料來源綁在一起，導致 scope 過寬且與 prototype 不符。這個 change 需要被收斂成「post-setup runtime page surface」本身，讓它能建立在獨立 route 與 backend contract 之上。

## What Changes

- 將 `runtime-dashboard` 收斂為 post-setup runtime dashboard 的頁面表面與監看元件組合，不再擁有 route 或 TopBar 導覽。
- 以單一 device 為主體，新增 focused device header、runtime summary、collector health、live points table 與 degraded state banner。
- 讓頁面只消費 backend-supported runtime contract 與 route-level dashboard state，不自行發明 queue backlog、diagnostic logs 或 fleet KPI。
- 明確定義 live 連線中斷時的 page-level 降級呈現，保留最後成功資料，避免首版 dashboard 因 SSE 抖動而變成空白頁。

## Non-Goals

- 不新增 `/studio/dashboard` 或任何新的 dashboard-first 產品入口。
- 不修改 `studio/v2` shell 的 TopBar、StepRail、SummaryRail 或首頁 IA。
- 不在首版 page surface 內加入 queue backlog、diagnostic logs、歷史圖表、fleet-wide card wall。

## Capabilities

### New Capabilities

- `runtime-dashboard`: 定義 post-setup runtime dashboard 的 focused page surface、monitoring panels 與 degraded monitoring 呈現。

### Modified Capabilities

(none)

## Impact

- Affected specs: `runtime-dashboard`
- Affected code:
  - New:
    - `frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardPage.tsx`
    - `frontend/src/features/datalink/runtime-dashboard/components/FocusedDeviceHeader.tsx`
    - `frontend/src/features/datalink/runtime-dashboard/components/RuntimeSummaryPanel.tsx`
    - `frontend/src/features/datalink/runtime-dashboard/components/CollectorHealthPanel.tsx`
    - `frontend/src/features/datalink/runtime-dashboard/components/LivePointsTable.tsx`
    - `frontend/src/features/datalink/runtime-dashboard/components/LiveStateBanner.tsx`
    - `frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx`
    - `frontend/src/i18n/locales/en/runtime-dashboard.json`
    - `frontend/src/i18n/locales/zh-TW/runtime-dashboard.json`
  - Modified:
    - `frontend/src/features/datalink/runtime-dashboard/RuntimeDashboardRoute.tsx`
    - `frontend/src/i18n/config.ts`
