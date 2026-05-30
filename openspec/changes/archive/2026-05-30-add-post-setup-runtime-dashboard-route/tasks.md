## 1. Route 與狀態契約先行

- [x] [P] 1.1 以 Red-Green 方式鎖定 `Dedicated post-setup runtime dashboard route` 與 `Decision: Place the dashboard on a dedicated route outside the setup steps`：先新增 route 測試，完成後可觀察到 `/studio/runtime?device_id=device-A` 直接渲染 runtime dashboard，且 `/studio/v2` 仍維持 setup flow；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx`。
- [x] [P] 1.2 以 Red-Green 方式鎖定 `Runtime dashboard uses device-focused context`、`Decision: Keep route ownership separate from page-surface ownership` 與 `Decision: Use device_id query state as the canonical dashboard context`：先新增 state 測試，完成後可觀察到沒有 `device_id` 時 route-level state 進入 `missing-device-context`，切換 device 時 query string 與狀態輸出同步改變；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx`。

## 2. Focused route 與 route-level state 實作

- [x] 2.1 落實 `Dedicated post-setup runtime dashboard route`：在 `frontend/src/App.tsx` 註冊 `/studio/runtime` 與 `RuntimeDashboardRoute` 容器，完成後 runtime dashboard 可以被直接 deep-link 開啟，且不需要經過 `/studio` 或 `/studio/v2` redirect；驗證：重新執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx`。
- [x] 2.2 落實 `Runtime dashboard uses device-focused context`、`Decision: Keep the page focused on one device instead of becoming a fleet dashboard` 與 `Decision: Use device_id query state as the canonical dashboard context`：在 `useRuntimeDashboardState` 中只維持單一 selected device query-state，完成後缺少 `device_id` 時輸出 focused empty state、有 `device_id` 時 route-level state 只載入該 device 的 dashboard data；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx`，並人工確認 route container 不會自行產生全局 dashboard。

## 3. Snapshot/live/fallback 載入實作

- [x] 3.1 落實 `Decision: Render snapshot first, then layer live stream with polling fallback` 與 `Runtime dashboard layers snapshot, live stream, and degraded fallback`：建立 `useRuntimeStatus`、`useRuntimeStream` 與 `useRuntimeDashboardState` 的組合流程，完成後 route-level state 會先取 snapshot、再接上 EventSource，且 stream attach 成功後狀態切換為 `live`；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx`，確認 `loading -> live` 遷移。
- [x] 3.2 落實 `Runtime dashboard layers snapshot, live stream, and degraded fallback`：在 stream 中斷時保留最後成功 snapshot、切為 `degraded` 並啟動 snapshot polling，完成後 page-surface change 可以直接消費穩定的 route-level fallback state；驗證：執行 `cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx`，確認 `live -> degraded` 後仍保留最後資料。

## 4. 收斂與驗證

- [x] 4.1 完成 `post-setup-runtime-dashboard` spec 的 route/state 收斂：執行 `spectra analyze add-post-setup-runtime-dashboard-route --json`、`spectra validate add-post-setup-runtime-dashboard-route`、`cd frontend && npm run test -- --run frontend/tests/unit/runtime-dashboard/runtime-dashboard-route.test.tsx frontend/tests/unit/runtime-dashboard/runtime-dashboard-state.test.tsx`，確認無 Critical/Warning，且 route、query-state、fallback 行為可被 `runtime-dashboard` page-surface change 與後續 Step 4 handoff change 直接依賴。
