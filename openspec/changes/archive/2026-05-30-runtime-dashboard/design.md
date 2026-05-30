## Context

`new_prototype/docs/_overview.md` 把 `studio/v2` 定位成四步驟 setup/workbench shell，加上一個獨立 settings 頁；`new_prototype/docs/step4-database.md` 只把 Runtime Dashboard 當成 commit 成功後的 `onCommit` 目的地，而不是 shell 首頁或 TopBar 常駐入口。現有 `WorkbenchV2Shell`、step rail 與 state shape 基本上已符合這個方向。

真正失焦的是舊的 `runtime-dashboard` change：它假設 `/studio/dashboard` route、TopBar 導覽、queue backlog、realtime logs 與五個 widget 會一次到位，這和 prototype 的 post-setup 角色、以及目前 backend contract 已知可支援的資料範圍都不對齊。這個 design 要做的是把 change 收斂成 focused page surface，建立在其他 changes 已拆出的 route 與 backend contract 之上。

## Goals / Non-Goals

**Goals:**

- 實作 post-setup runtime dashboard 的 focused page surface，承接獨立 route 與 backend runtime contract。
- 在單一 device 脈絡下呈現 device header、runtime summary、collector health、live points table 與 degraded monitoring banner。
- 讓頁面在 live stream 中斷時仍保留最後成功資料並明確顯示 degraded 狀態。
- 以可測試、可拆分的 page/panel 元件完成首版 dashboard 表面，不重做 v2 shell。

**Non-Goals:**

- 不建立或重新命名 route；`/studio/runtime` 與 query-state 由 `add-post-setup-runtime-dashboard-route` change 擁有。
- 不在 shell TopBar 新增 dashboard 常駐導覽，也不把 `/studio/v2` 改成 dashboard 首頁。
- 不在首版頁面加入 queue backlog、diagnostic logs、歷史趨勢圖、fleet-wide overview 或控制型操作。
- 不修改 backend event shapes；這由 `add-runtime-dashboard-backend-contract` change 擁有。

## Decisions

### Decision: Keep runtime dashboard as a post-setup page surface, not a workflow shell replacement

- **Rationale**: `studio/v2` 目前已經是 setup shell，不應再被 dashboard-first 需求反向擴張。把 runtime dashboard 視為 post-setup page surface，才能保持 prototype 的角色分工，也能讓 Step 4 handoff 成為主要入口。
- **Alternatives considered**:
  - **Make dashboard a fifth step inside `/studio/v2`**: 會把 setup 與運轉監看混成同一條任務流，違反 prototype。
  - **Add a TopBar dashboard destination for all times**: 會讓 dashboard 重新被誤解成新的首頁或平行產品入口。

### Decision: Compose the page around the route-level dashboard state contract

- **Rationale**: route change 已經擁有 `device_id` query-state、snapshot/live/fallback 狀態機與 data hooks。這個 change 不該重做這些邏輯，而應由 `RuntimeDashboardPage` 消費 route-level contract，例如 `selectedDeviceId`、`selectedDevice`、`snapshot`、`streamState`、`liveValues`、`onSelectDevice`。
- **Alternatives considered**:
  - **Have the page create its own status/stream hooks**: 會讓 route/state ownership 回流到 page surface，破壞拆分邊界。
  - **Render everything inside the route container**: 雖然少一層元件，但 route 與 page composition 會再度糾纏。

### Decision: Limit the first page surface to backend-supported monitoring panels

- **Rationale**: 首版 dashboard 成功標準是「看懂剛上線那台 device 現在有沒有在跑」，不是把所有觀測構想都塞上去。既然 backend contract 明確支援的是 runtime summary、collector health 與 live points，頁面就應只呈現這些可依賴資訊。
- **Alternatives considered**:
  - **Include queue backlog and diagnostic logs in the first page**: 目前沒有同等穩定的正式 contract，會逼前端回頭依賴舊端點或臨時欄位。
  - **Build a fleet KPI card wall first**: 會把 focused monitoring 變成另一個產品問題。

### Decision: Preserve last-known values during degraded runtime monitoring

- **Rationale**: 對 post-setup 驗證來說，最糟不是資料稍舊，而是頁面一失去 stream 就清空成空白。保留最後一次成功 snapshot/live values，並把 stream badge/bannner 轉為 degraded，能讓操作者維持連續判讀。
- **Alternatives considered**:
  - **Clear all panels on stream disconnect**: 會讓操作者失去比較基準，也難以分辨是設備停了還是瀏覽器斷線。
  - **Hide degraded state and silently keep polling**: 使用者會誤以為畫面仍是完全即時。

## Implementation Contract

- **Behavior**:
  - 當 route-level dashboard state 提供有效 `device_id` 與 snapshot 時，頁面 SHALL 顯示 focused device header、runtime summary、collector health 與 live points table。
  - 當 `streamState` 為 `live` 時，頁面 SHALL 以 live badge / copy 呈現正在接收即時資料。
  - 當 `streamState` 轉為 `degraded` 時，頁面 SHALL 顯示 degraded banner，並保留最後成功的 summary 與 point values。
  - 當 route-level state 為 `missing-device-context` 或 `loading` 時，頁面 SHALL 顯示對應的 focused empty/loading surface，而不是自行改成 fleet dashboard。
- **Interface / Data Shape**:
  - 頁面 SHALL 消費 route-level state contract，而不是自行決定資料來源。至少包含：
    - `selectedDeviceId: string | null`
    - `selectedDevice: { id: string; name: string; protocol?: string } | null`
    - `snapshot: RuntimeStatusSnapshot | null`
    - `streamState: 'loading' | 'live' | 'degraded' | 'missing-device-context' | 'error'`
    - `liveValues: RuntimePointValue[]`
    - `onSelectDevice(deviceId: string): void`
  - `RuntimeSummaryPanel` SHALL 只呈現 backend contract 已保證的 metrics 與 collector summary 欄位。
  - `CollectorHealthPanel` SHALL 呈現至少 `status`、`points_total`、`points_healthy`、`points_stale`、`points_error`、`last_read_at`、`breaker_state`。
  - `LivePointsTable` SHALL 呈現 point address、raw value、transformed value、quality/stale 與 timestamp。
- **Failure modes**:
  - 缺少 `selectedDeviceId` 時，頁面 SHALL 呈現 focused empty state，不自行嘗試列出全 fleet。
  - route-level state 若回傳 error/degraded，頁面 SHALL 保留最後成功資料，並只把 live 狀態呈現為 degraded/error。
  - 若 `liveValues` 尚未到達但 snapshot 已成功，頁面 SHALL 先呈現 snapshot summary 與空的 live points placeholder。
- **Acceptance criteria**:
  - `frontend/tests/unit/runtime-dashboard/runtime-dashboard-page.test.tsx` 驗證 Focused runtime dashboard surface 會在有效 state 下渲染 header、summary、health、points 四個主要區塊。
  - 同一測試檔驗證 Runtime summary presents backend-supported metrics 與 Collector health panel 只顯示正式 contract 欄位，不依賴 queue/log/backlog 資料。
  - 同一測試檔驗證 Live point values table 會顯示 raw/transformed/timestamp，並在 degraded 狀態下保留最後資料。
  - `spectra analyze runtime-dashboard --json` reports no Critical/Warning findings and `spectra validate runtime-dashboard` passes.
- **Scope boundaries**:
  - In scope: page composition, monitoring panels, i18n strings, degraded banner, page-level tests.
  - Out of scope: route registration, Step 4 navigation wiring, backend event shapes, fleet dashboard, queue/log/history panels.

## Risks / Trade-offs

- `[Risk]` route/state contract 之後若再變動，page surface 可能需要跟著調整 prop shape -> `[Mitigation]`: 將頁面依賴集中在單一 route-level interface，避免各 panel 直接呼叫多個 hooks。
- `[Risk]` live points 數量一多時表格可能變重 -> `[Mitigation]`: 首版先保持 focused device scope，必要時僅在 table 元件內做輕量 rendering control，不擴張頁面範圍。
- `[Risk]` degraded banner 若過於搶眼會讓使用者誤判設備故障 -> `[Mitigation]`: 文案明確指出是 live stream 連線降級，而不是直接宣告 device offline。

## Migration Plan

1. 先由 route change 提供穩定的 route-level dashboard state contract。
2. 再由本 change 將 focused page surface 與 panels 接上該 contract。
3. 完成 page tests 後，再讓 Step 4 handoff change 把成功卡入口接到新 route。

## Open Questions

- `selectedDevice` selector 是否由 route state 直接提供完整 display model，或只提供 `device_id` 讓 page 再查名稱；本 change 先以「route state 直接提供最小 display model」為預設，減少 page 重複查找。
