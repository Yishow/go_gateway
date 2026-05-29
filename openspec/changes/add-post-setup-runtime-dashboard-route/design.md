## Context

`/studio/v2` 現在是一個明確的四步驟 setup/workbench shell，這和 `new_prototype/docs/_overview.md` 的 shell 定位一致；`new_prototype/docs/step4-database.md` 也已經把 Runtime Dashboard 定義成 commit 成功後 `onCommit` 的目的地，而不是 shell 首頁。使用者已經決定 runtime dashboard 不應取代 setup flow，而是作為「這台剛上線的 device 現在是否正常運轉」的確認面板。

這代表 route 設計必須同時滿足兩件事：第一，它要和 `/studio/v2` 的 setup flow 分開，不把 v2 shell 再長成 generic dashboard；第二，它要能直接吃剛 commit 的 `device_id`，並在 stream 不可用時維持最低限度可觀測性，讓操作者不用退回 legacy dashboard。這個 change 只負責 route/query-state/state contract，不擁有 page surface。

## Goals / Non-Goals

**Goals:**

- 新增一個明確的 post-setup runtime dashboard route，供後續 Step 4 handoff 與直接 deep-link 使用。
- 讓 dashboard 以單一 `device_id` 為主上下文，首屏聚焦剛完成設定的 device。
- 用 backend runtime contract 載入 snapshot 與 live stream，並提供 stream 失敗時的 polling fallback。
- 提供在頁面內切換其他 device 的能力，但維持同一個 focused monitoring mental model。

**Non-Goals:**

- 不把 `/studio/v2` 改成 dashboard 首頁。
- 不在這個 change 內接上 Step 4 commit 成功後的導航 wiring。
- 不做 fleet-wide KPI 首頁、歷史圖表、diagnostic logs、queue backlog、或全域 alert center。
- 不重寫 TopBar、StepRail、SummaryRail 等 v2 setup shell 元件。

## Decisions

### Decision: Place the dashboard on a dedicated route outside the setup steps

- **Rationale**: runtime dashboard 是 setup 之後的操作面，不是 Step 5，也不是 `/studio/v2` 的首頁。獨立 route 最能保留 workbench shell 的任務導向定位，同時也方便直接 deep-link。
- **Alternatives considered**:
  - **Embed as a fifth step inside `/studio/v2`**: 會把 setup flow 與運轉監看混成一頁，違反你剛剛定下的角色分工。
  - **Replace `/studio/v2` landing with dashboard**: 這會改變 v2 IA，超出本 change 範圍。

### Decision: Keep route ownership separate from page-surface ownership

- **Rationale**: route 的責任是決定 URL、query-state、snapshot/live/fallback 狀態機與 page props contract；page surface 的責任是把這些狀態渲染成 header/panels/banner。分開 ownership 才能讓 `runtime-dashboard` change 專心做 UI surface，不重複實作 route 邏輯。
- **Alternatives considered**:
  - **Let this change also own the page panels**: 會重新把 route 與 UI surface 綁回同一個 change，失去拆分價值。
  - **Let the page surface own route parsing**: 會讓 `device_id` query-state 與 fallback 行為散落在 UI 元件裡。

### Decision: Use device_id query state as the canonical dashboard context

- **Rationale**: post-setup dashboard 的核心是「先看剛 setup 完的那台」，而 query state 既能被 Step 4 handoff 帶入，也能被手動切換或直接貼連結重現。它比只放在 memory state 或 reducer 裡更穩定。
- **Alternatives considered**:
  - **Store current device only in in-memory state**: 重新整理頁面後會丟失上下文，不利於分享與除錯。
  - **Use path params only**: 也可行，但 query param 更容易與未來額外 filter 共存，且不需要為切換動作重建 route tree。

### Decision: Render snapshot first, then layer live stream with polling fallback

- **Rationale**: dashboard 不能等到 SSE 建立完成才有畫面。先拿 `runtime/status` 可以讓頁面立即顯示最新狀態，之後再用 `runtime/stream` 補 live update；若 stream 中斷，退回 snapshot polling 仍能維持可觀測性。
- **Alternatives considered**:
  - **Stream-only page**: 首屏空白時間過長，且一旦 EventSource 失敗就只剩錯誤畫面。
  - **Polling-only page**: 可做，但會失去 live point updates，無法滿足 post-setup「立即確認寫入是否起跑」的需求。

### Decision: Keep the page focused on one device instead of becoming a fleet dashboard

- **Rationale**: 這個 route 的成功標準不是「把所有指標塞滿」，而是讓操作者在 setup 完之後可以最快看懂一台設備的 runtime 狀態。fleet overview 是另一個產品問題，不該偷偷在這裡長出來。
- **Alternatives considered**:
  - **Show all devices by default with cards and filters**: 容易變成 generic dashboard，且與剛 commit 後的 immediate verification 情境不符。

## Implementation Contract

- **Behavior**:
  - 系統 SHALL 新增一個 post-setup runtime dashboard route，例如 `/studio/runtime`，可直接以 URL 開啟。
  - 當 URL 帶有 `device_id=<id>` 時，頁面 SHALL 以該 device 作為主上下文，先顯示 snapshot，再接上 live stream。
  - 當使用者在頁面內切換 device 時，URL 中的 `device_id` SHALL 同步更新，並重新建立 snapshot/stream 訂閱。
  - 當 live stream 不可用或中斷時，頁面 SHALL 顯示 degraded/live disconnected 狀態，並以固定節奏重新輪詢 snapshot，直到 stream 恢復。
- **Interface / data shape**:
  - Route: `/studio/runtime` with query `device_id` as canonical selected device
  - Data sources:
    - `runtimeAPI.getStatus(deviceId)` for initial and fallback snapshot
    - `runtimeAPI.getStreamUrl(deviceId, pointIds?)` for EventSource connection
  - Route-level state SHALL distinguish at least `loading`, `live`, `degraded`, and `missing-device-context`
  - Route-level state SHALL expose at least `selectedDeviceId`, `selectedDevice`, `snapshot`, `streamState`, `liveValues`, and `onSelectDevice` for the page-surface change to consume
- **Failure modes**:
  - Missing `device_id` SHALL produce a clear empty or choose-device state instead of silently opening a fleet dashboard.
  - Snapshot load failure SHALL surface a recoverable error state with retry.
  - Stream failure SHALL preserve the last successful snapshot while fallback polling runs.
- **Acceptance criteria**:
  - Route tests prove `/studio/runtime?device_id=<id>` renders the runtime dashboard page.
  - State tests prove the page enters `missing-device-context` without `device_id`, `loading` before snapshot, `live` after stream attach, and `degraded` after stream disconnect.
  - Integration-level UI tests prove switching devices updates the query string and rebinds data loading.
  - `spectra analyze add-post-setup-runtime-dashboard-route --json` reports no Critical/Warning findings and `spectra validate add-post-setup-runtime-dashboard-route` passes.
- **Scope boundaries**:
  - In scope: new route, query-state, route-level state contract, device context handling, snapshot/stream/fallback wiring.
  - Out of scope: page panels and layout composition, Step 4 button wiring, fleet homepage, history charts, queue/log panels, backend contract changes.

## Risks / Trade-offs

- `[Risk]` query-driven context 讓頁面在沒有 `device_id` 時看起來比較空 -> `Mitigation`: 明確設計 `missing-device-context` 狀態，而不是偷偷降級成全局 dashboard。
- `[Risk]` snapshot 與 stream 切換時可能出現短暫狀態閃爍 -> `Mitigation`: 保留最後成功 snapshot，只有 stream 狀態 badge 改為 degraded，不清空整頁。
- `[Risk]` 後續若真的需要 fleet dashboard，這個 route 名稱與定位可能不夠通用 -> `Mitigation`: 現階段刻意接受 focused naming，未來 fleet overview 另開 change。

## Migration Plan

1. 先新增 route、page state 與 focused layout skeleton。
2. 接上 `runtime/status` snapshot 與 `runtime/stream` EventSource。
3. 完成 stream disconnect → polling fallback 後，再交給 Step 4 handoff change 連接入口。
4. 若 rollout 過程發現 live stream 不穩，頁面可暫時只保留 snapshot polling 而不需要移除 route。
