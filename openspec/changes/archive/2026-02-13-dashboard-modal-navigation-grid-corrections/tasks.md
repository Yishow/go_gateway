## 1. Sidebar to Modal Migration (except /test)

- [x] 1.1 建立 sidebar 路由到 modal 的映射表（devices/settings/points/mappings/wizard）
- [x] 1.2 實作 `/datalink` 內 modal state orchestration（query + state）
- [x] 1.3 確認 `/test` 維持獨立頁且不受 modal 導向影響
- [x] 1.4 補齊 legacy route redirect 測試（含 modal 開啟意圖）

## 2. Local Modbus Entry Connectivity

- [x] 2.1 在 context bar 新增固定 `Server Memory Grid` 入口
- [x] 2.2 在 Quick Actions 新增固定 `Server Memory Grid` CTA
- [x] 2.3 在 dashboard Grid 卡片新增 `完整工作台` 連結
- [x] 2.4 local-modbus 頁新增回跳 dashboard 入口並保留 section

## 3. Dashboard Grid Enhancement

- [x] 3.1 下修 dashboard memory grid 視覺高度（desktop/tablet/mobile）
- [x] 3.2 新增 grid 摘要資訊（bind_state/mapping_count/conflict_count）
- [x] 3.3 保持 preflight gate 一致視覺（PASS/FAIL）

## 4. Typed Occupancy Fix

- [x] 4.1 修正 `float32` 佔 2 格連續規則（渲染 + 計算）
- [x] 4.2 修正 `int64` 佔 4 格連續規則（渲染 + 計算）
- [x] 4.3 補齊 typed occupancy 單元測試（含 pair/group 標示）
- [x] 4.4 補齊 SmartDashboard 互動測試驗證連續格視覺與提交前檢核

## 5. Explicit Exclusions

- [x] 5.1 不實作原建議第 10（conflict deep-link query）
- [x] 5.2 不實作原建議第 15（recent write events）
- [x] 5.3 不實作原建議第 16（Top3 conflict registers）
- [x] 5.4 在設計文件與 PR 說明明確標註排除項目

## 6. Verification

- [x] 6.1 執行 `cd frontend; npm run test`
- [x] 6.2 執行 `cd frontend; npm run lint`
- [x] 6.3 補上 UI 截圖證據（dashboard + modal + local-modbus 導流）
