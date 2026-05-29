## 1. 基礎架構與路由設定 (Setup)

- [ ] [P] 1.1 註冊多語言資源。新增 `frontend/src/i18n/locales/zh-TW/runtime-dashboard.json` 與 `frontend/src/i18n/locales/en/runtime-dashboard.json` 包含 Dashboard 所有繁中與英文文案，並在 `frontend/src/i18n/config.ts` 內註冊 `runtime-dashboard` 命名空間。驗證：執行 `npm run lint` 且載入 config 時無 namespace warnings。
- [ ] [P] 1.2 註冊路由與導覽列入口。在 `frontend/src/App.tsx` 中新增 `/studio/dashboard` 路由指向 `RuntimeDashboard` 頁面，並在 `frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx` 元件中新增對應的導覽按鈕，對齊 "Navigation and Layout" 規格。驗證：撰寫 `runtime-dashboard.test.tsx` 中 "User navigates to dashboard" 測試，模擬點擊 TopBar 中的連結，斷言 React Router 成功轉址至 `/studio/dashboard` 且渲染主容器。

## 2. 實時 SSE 數據源對接 (Core Data Feed)

- [ ] [P] 2.1 實作 SSE 資料流訂閱。在 `frontend/src/features/datalink/workbench-v2/dashboard/RuntimeDashboard.tsx` 內實作對 `/api/v1/datalink/runtime/stream` 端點的訂閱，落實設計決策 "Decision: 使用 SSE (Server-Sent Events) 進行實時數據推送"。驗證：單元測試中使用 `vi.spyOn(window, 'EventSource')`，斷言 SSE 成功發起連線，且當連線中斷時頂部會顯示黃色 Warning Banner，並在 5 秒後自動嘗試重連。

## 3. 模組化監控元件開發 (Modular Component Integration)

- [ ] [P] 3.1 實作指標卡片。依據設計決策 "Decision: 前端監控元件模組化與狀態隔離" 建立 `MetricsOverview.tsx` 元件，展示 Points/Sec、Success Rate、Avg Latency 採集指標，滿足規格 "Real-time Metrics Processing"。驗證：在 `runtime-dashboard.test.tsx` 中對 "Dashboard metrics update" 進行 TDD 驗證，提供 mock SSE `metrics` 事件資料，斷言畫面上正確顯示計算後的 95.0% 成功率與 15.0ms 平均延遲。
- [ ] [P] 3.2 實作設備健康狀態網格。建立 `DeviceStatusGrid.tsx` 元件，以格狀展示所有 PLC 設備，對齊規格 "Device Status Grid Visualization"。驗證：在 `runtime-dashboard.test.tsx` 中撰寫 TDD 測試 "Device connectivity changes"，模擬 SSE 推送設備離線狀態，斷言對應設備項目背景色轉紅，並顯示 "Offline" 與最新時間戳記。
- [ ] [P] 3.3 實作點位即時數值表。建立 `LiveValuesTable.tsx` 元件，以表格展示點位即時轉換值，對齊規格 "Live Value Streaming"。驗證：在 `runtime-dashboard.test.tsx` 中撰寫 TDD 測試 "Value transform and display"，給定 point (multiplier=0.1, offset=5)，模擬 SSE 推送 raw=150 事件，斷言畫面上即時值更新為 "20.0" 並附帶 Throttle 限制（100ms 內防抖）。
- [ ] [P] 3.4 實作輸出佇列積壓與錯誤日誌。建立 `OutputQueueStatus.tsx` 與 `RealtimeLogs.tsx` 元件，滿足規格 "Output Queue Backlog Indicator" 與 "Real-time Log Stream"。驗證：在 `runtime-dashboard.test.tsx` 撰寫 TDD 測試 "Queue backlog alert" 與 "Error log output"，模擬 SSE 推送 backlog=501 與 TCP 逾時 error 事件，斷言畫面呈現琥珀色警告 Banner，且日誌控制台以紅字追加 timestamped 錯誤訊息。

## 4. 全量驗證與生產打包 (Final Integration)

- [ ] 4.1 執行完整前端建置與測試。在本機環境執行 `npm run lint && npm run test && npm run build`。驗證：前端編譯與 201+ 個測試全綠通過，生產 JS/CSS 資產打包成功無 warnings。
