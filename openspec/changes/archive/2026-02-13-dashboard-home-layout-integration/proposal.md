## Why

目前 `/datalink` 由外層 `DatalinkLayout` 提供 sidebar 與 top wrapper，造成視覺層級分散、垂直空間被占用，且與「設備通常只選一次、僅切換時才需要目錄」的實際操作習慣不一致。需要將 Dashboard 升級為首頁主容器，並把設備切換流程收斂為 context bar + 按需展開的 drawer。

## What Changes

- 移除 `/datalink` 外層 `DatalinkLayout` 的常駐 sidebar/top wrapper，改為 `SmartDashboard` 自主承載首頁框架。
- 將「功能分頁列」與「已選設備狀態列」改為同列水平布局，降低垂直占用。
- 保留設備樹能力，但從常駐左欄改為「按下切換設備才展開」的 drawer。
- 將設備切換異常處理納入統一互動：離線、切換失敗、未儲存變更、權限不足。
- 調整 legacy 路由導向策略，確保 `/datalink/devices`、`/datalink/settings` 能導向首頁內對應區段。
- 除測試頁外，設備建立/編輯流程改為 dashboard 內 modal（含 stepper）。
- 新增本機 `5020` memory grid 寫入的獨立工作頁，提供完整衝突治理與操作能力。

## Capabilities

### New Capabilities

- `dashboard-device-context-switching`: 定義首頁設備 context bar、切換設備 drawer 與異常切換流程。
- `local-modbus-memory-workbench`: 定義本機 `5020` 寫入映射的獨立頁與完整操作功能。

### Modified Capabilities

- `datalink-ui`: 導覽模型由外層 sidebar 轉為 dashboard 內整合導覽，並更新 legacy 導向規則。
- `datalink-smart-dashboard`: 設備樹由常駐區改為按需切換模式，並補齊切換設備異常狀態要求。
- `protocol-servers`: 補充本機 `5020` sink 的衝突檢核與工作台整合行為。

## Impact

- Frontend routing: `frontend/src/App.tsx`、`frontend/src/layouts/DatalinkLayout.tsx`
- Dashboard layout and interactions: `frontend/src/pages/datalink/SmartDashboard.tsx`
- Device navigation component behavior: `frontend/src/components/datalink/DeviceTreeNav.tsx`
- Local Modbus sink workspace: `frontend/src/pages/datalink/LocalModbusWorkbenchPage.tsx`（new）
- 可能影響既有測試：路由、layout 快照、Dashboard 互動測試
