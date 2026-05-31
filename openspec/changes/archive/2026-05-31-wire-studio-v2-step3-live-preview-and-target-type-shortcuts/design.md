## Context

Step 3 已經有 workspace-scoped mapping autosave，但右側 `轉換管線預覽` 仍使用 deterministic seed + 本地 transform 計算。使用者要求 `/studio/v2` 全站對接，因此 preview 不能再停留在 mock。

目前 backend 已經提供 `POST /api/v1/datalink/mappings/preview`，而 Step 3 reducer / autosave 也已有穩定資料流。這次 change 的正確邊界是「重用現有 preview API 與既有 reducer」，而不是新增另一套 V2 preview 契約。

## Goals / Non-Goals

**Goals:**

- Step 3 右側 preview 使用真實 backend preview API
- preview 正確處理 loading、error 與 stale response
- 新增 `目標型態` 快捷設定，支援當前列與全部列
- payload preview 保持 display-only

**Non-Goals:**

- 不新增新的 workspace preview endpoint
- 不改 Step 3 autosave 主體
- 不改 payload preview JSON shape
- 不重做 footer 既有 `scale/offset/target_type` bulk apply

## Decisions

### Reuse existing mapping preview API

Step 3 直接重用 `POST /api/v1/datalink/mappings/preview`。這是最小改動，也符合 `/studio/v2` 接真後端的主線。

### Deterministic raw sample stays client-chosen

目前 Step 3 已有 deterministic raw seeds `[243, 251, 1024, 985, 67, 542, 18, 1450]`。這次不改 sample 策略，而是把該 raw value 送進 backend preview API，讓 transform 真相來自 backend。

### Preview container owns async state

`TransformPreview` 需要持有 debounce、loading、error 與 stale response guard。這些 async concerns 不應塞進 mapping reducer。

### Payload preview remains local-only

即使 transform preview 改為 backend-backed，payload preview 仍然只顯示當前 draft payload。它不得自己再發 request。

### Target type shortcuts live in the preview card

快捷操作放在 preview 卡最合理，因為：

- 與 preview 結果放在同一個決策區塊
- 不會讓表格列再膨脹
- 可直接重用 selected row context

### Apply-all reuses existing bulkApplyTransform

`套用到全部列` 不新增新 action，而是沿用既有 `bulkApplyTransform`，但只傳 `['target_type']`。

## Implementation Contract

- Behavior:
  - 當 selected row 存在且其 draft transform 改變時，系統會在短 debounce 後發送 preview request
  - preview success 時顯示 backend `step_results` 與 `final_value`
  - preview failure 時顯示 backend error，不回退成本地假成功
  - 點 target type chip 只更新當前列 `target_type`
  - 點 `套用到全部列` 會將目前 selected row 的 `target_type` 複製到 `state.mappings`
- Interface / data shape:
  - request: `POST /api/v1/datalink/mappings/preview`
  - payload 至少包含 `raw_value` 與 `transform_pipeline`
  - response 使用既有 `raw_value`、`final_value`、`step_results`、`error`
- Failure modes:
  - 沒有 selected row 時不發 request
  - 較舊 request 回來時不得覆蓋較新 preview
  - payload preview 不得因 transform preview request 而觸發新的 request
- Acceptance criteria:
  - 前端測試覆蓋 loading/success/error/stale response
  - 前端測試覆蓋 target type 當前列快捷設定與 apply-all
  - `cd frontend && npm run build`
  - `git diff --check`
- Scope boundaries:
  - In scope: Step 3 preview wiring、快捷操作、spec/task 更新
  - Out of scope: Step 1/2/4、settings、runtime、new API creation
- Not complete if:
  - preview 仍靠本地 `runScale/castValue/formatFinal` 直接算出 operator-visible 結果
  - preview error 仍被掩蓋成 success
  - `目標型態` 仍只能逐列 select，沒有當前列/全部列快捷操作

## Risks / Trade-offs

- [Risk] preview request 頻率過高 → Mitigation: 固定短 debounce
- [Risk] async preview 邏輯讓元件責任變重 → Mitigation: 把 request state 限制在 preview 容器，不滲透到 reducer
- [Risk] footer bulk apply 與 preview 卡 apply-all 看起來重複 → Mitigation: preview 卡只處理 `target_type`，footer 保留完整 transform 套用
