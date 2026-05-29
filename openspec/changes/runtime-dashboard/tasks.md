## 1. Focused page surface 先行

- [ ] [P] 1.1 以 Red-Green 方式鎖定 `Focused runtime dashboard surface` 與 `Decision: Keep runtime dashboard as a post-setup page surface, not a workflow shell replacement`：先新增 `frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx`，完成後可觀察到 route-level state 提供有效 device context 時，頁面會渲染 focused device header、runtime summary、collector health 與 live points 四個主要區塊；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx`。
- [ ] [P] 1.2 以 Red-Green 方式鎖定 `Runtime summary presents backend-supported metrics`、`Collector health panel` 與 `Decision: Limit the first page surface to backend-supported monitoring panels`：在同一測試檔先寫出只接受正式 contract 欄位的斷言，完成後可觀察到頁面不需要 queue backlog 或 diagnostic logs 也能成立；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx`。

## 2. Page composition 與 panel 實作

- [ ] 2.1 落實 `Decision: Compose the page around the route-level dashboard state contract`：建立 `RuntimeDashboardPage.tsx` 與 page props contract，完成後 route container 只需傳入 `selectedDeviceId`、`selectedDevice`、`snapshot`、`streamState`、`liveValues`、`onSelectDevice` 就能渲染整頁；驗證：重新執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx`。
- [ ] 2.2 落實 `Focused runtime dashboard surface`：建立 `FocusedDeviceHeader.tsx`、`RuntimeSummaryPanel.tsx`、`CollectorHealthPanel.tsx` 與 `LivePointsTable.tsx`，完成後單一 device 的 runtime 狀態可在一頁中被讀懂；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx`，並人工檢查不會出現 shell-level 導覽改動。

## 3. Degraded monitoring 呈現

- [ ] 3.1 落實 `Degraded runtime monitoring banner` 與 `Decision: Preserve last-known values during degraded runtime monitoring`：建立 `LiveStateBanner.tsx` 並在 `streamState=degraded` 時保留最後成功 summary 與 live values，完成後使用者可以分辨是 live feed 降級而不是整頁失效；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx`，確認 degraded state 仍保留最後資料。
- [ ] 3.2 落實 `Live point values table`：讓 `LivePointsTable` 在 snapshot 已就緒但 live values 尚未到達時顯示 placeholder，當 live values 到達後顯示 raw/transformed/timestamp，完成後不會因為首筆 SSE 尚未抵達而產生空白錯覺；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx`。

## 4. 收斂與驗證

- [ ] [P] 4.1 註冊 `runtime-dashboard` i18n namespace 並補齊 focused page surface 文案，完成後頁面所有新字串都來自 `frontend/src/i18n/locales/zh-TW/runtime-dashboard.json` 與 `frontend/src/i18n/locales/en/runtime-dashboard.json`；驗證：執行 `cd frontend && npm run lint` 並人工檢查 i18n config 無遺漏 namespace。
- [ ] 4.2 完成 `runtime-dashboard` change 的 artifact 與頁面收斂：執行 `spectra analyze runtime-dashboard --json`、`spectra validate runtime-dashboard`、`cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx`，確認無 Critical/Warning，且 page surface 可直接建立在 route/backend contract 之上。
