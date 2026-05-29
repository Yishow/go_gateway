## Context

在目前的 v2 shell 中，Step 4 commit 成功後確實會出現 emerald success card，也有 `前往 Runtime Dashboard` 按鈕；`new_prototype/docs/step4-database.md` 也明確把這個按鈕定義成 `onCommit` callback 的出口。但目前 shell 傳入的 `onCommit` callback 仍只是 `console.log('Database committed successfully!')`。也就是說，產品語意已經暗示有下一步，實際互動卻沒有真的離開 setup flow。

前兩個 changes 已經把 runtime dashboard backend contract 與 route/page contract 拆出來了，這個 change 只負責把最後一段入口補上。它不應重新定義 dashboard route，也不應偷渡新的 Step 4 UI；它只需要把成功卡的 CTA 變成可依賴的 handoff。

## Goals / Non-Goals

**Goals:**

- 讓 Step 4 成功卡的 `前往 Runtime Dashboard` 按鈕成為真正的導航入口。
- 定義從現有 workbench state 解析 handoff `device_id` 的 deterministic 規則。
- 讓無法解析單一 device context 的情況也有明確 fallback，而不是 noop 或留在原頁面。
- 用 Step 4 與 shell routing tests 把 handoff 行為鎖定。

**Non-Goals:**

- 不建立新的 runtime dashboard route 或改動 backend runtime contract。
- 不把 success card 改成多按鈕 action center。
- 不擴張到 TopBar 新入口、fleet dashboard、或 Step 4 成功後自動跳轉。

## Decisions

### Decision: Keep navigation orchestration in WorkbenchV2Shell

- **Rationale**: `CommitSuccessCard` 是展示元件，`Step4Database` 是 Step 4 容器；真正知道 route 與整體 app 導航語意的是 `WorkbenchV2Shell`。導航應由 shell 注入 callback，而不是把 router 硬塞進 success card。
- **Alternatives considered**:
  - **Navigate directly inside `CommitSuccessCard`**: 會把 presentational component 綁死到 routing layer，降低可測性。
  - **Navigate directly inside `Step4Database`**: 較 `CommitSuccessCard` 好，但仍把 Step 元件和 app-level route contract 綁在一起。

### Decision: Resolve the handoff device from existing workbench state with a deterministic order

- **Rationale**: `WorkbenchV2State` 沒有單獨的 `selectedDeviceId`，但仍有足夠訊號可推導當前最合理的 handoff device。這個規則必須固定，否則多 device 狀態下每次實作都會猜不同答案。
- **Alternatives considered**:
  - **Always use `state.devices[0]`**: 在單 device happy path 可行，但多 device 時太脆弱。
  - **Block navigation unless an explicit selected device field exists**: 會把一個可恢復的情況變成無入口。

### Decision: Prefer explicit runtime route fallback over silent no-op

- **Rationale**: 如果無法解析單一 device，最糟不是看到 `missing-device-context`，而是按了按鈕什麼都沒發生。既然 route change 已經定義 clear empty state，這個 change 應主動利用它。
- **Alternatives considered**:
  - **Do nothing when device resolution fails**: 使用者會以為按鈕壞了。
  - **Invent a fallback fleet dashboard destination**: 會違反 focused monitoring 的前提。

## Implementation Contract

- **Behavior**:
  - 當 Step 4 success card 顯示且操作者點擊 `前往 Runtime Dashboard` 時，系統 SHALL 透過 shell-level callback 導向 runtime dashboard route。
  - 若可從當前 workbench state 解析出 handoff device，導向目標 SHALL 為 `/studio/runtime?device_id=<resolved-id>`。
  - 若無法解析出 handoff device，導向目標 SHALL 為 `/studio/runtime`，由 route change 定義的 `missing-device-context` 狀態接手。
- **Interface / data shape**:
  - `CommitSuccessCard` 仍只接收 `onCommit` callback，不直接依賴 router。
  - `WorkbenchV2Shell` SHALL 注入 route-aware `onCommit` callback。
  - device resolution helper SHALL 使用固定順序：
    1. `selectedRuleId` 對應 rule 的 `device_id`
    2. 若所有 active rules 屬於同一個 `device_id`，使用該 device
    3. 若 `state.devices.length === 1`，使用唯一 device
    4. 否則回傳無法解析
- **Failure modes**:
  - 若 helper 無法解析單一 device，點擊按鈕仍 SHALL 導航到 `/studio/runtime`，不得停留在原頁或只輸出 console log。
  - 若 runtime route 尚未載入 live data，這不屬於本 change 的責任；本 change 只保證 handoff destination 正確。
- **Acceptance criteria**:
  - Step 4 commit tests prove clicking the success-card button still invokes `onCommit` exactly once.
  - Shell routing tests prove a resolved device navigates to `/studio/runtime?device_id=<id>`.
  - Device-resolution unit tests prove the deterministic order and unresolved fallback case.
  - `spectra analyze wire-step4-runtime-dashboard-handoff --json` reports no Critical/Warning findings and `spectra validate wire-step4-runtime-dashboard-handoff` passes.
- **Scope boundaries**:
  - In scope: callback wiring, device resolution helper, route target construction, tests.
  - Out of scope: backend runtime loading, route page layout, auto-redirect after commit completes.

## Risks / Trade-offs

- `[Risk]` 以 rule/device 清單推導 handoff device 仍可能在極端 multi-device draft 狀態下失敗 -> `Mitigation`: failure path 明確導向 `/studio/runtime` 而非 noop，並用 helper tests 固定規則。
- `[Risk]` 若之後新增真正的 `selectedDeviceId`，這個 helper 可能需要調整 -> `Mitigation`: 將解析邏輯集中到單一 helper，而非散落在 shell 與 Step 元件中。
- `[Risk]` 使用者可能期待 commit 成功後自動跳轉 -> `Mitigation`: 本 change 保持顯式按鈕 handoff，避免在 setup flow 結束時突然切頁。

## Migration Plan

1. 先新增 device resolution helper 與 unit tests。
2. 在 shell 中將 `console.log` callback 替換為 route-aware navigation。
3. 更新 Step 4 / shell 測試，確認 resolved 與 unresolved 兩條 handoff 路徑。
4. 完成後由使用者或後續 change 決定是否要在其他入口重用相同 helper。
