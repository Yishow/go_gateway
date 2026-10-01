# V2 Step 3 Live Preview And Target Type Shortcuts Design

- Status: Brainstorming-approved; pending written-spec review
- Date: 2026-05-30
- Scope: `/studio/v2` Step 3 only

## 1. Summary

`/studio/v2` 的 Point → Tag 頁面目前仍保留兩個 operator-visible 缺口：

1. `轉換管線預覽` 仍是前端本地 seed + 本地 transform 計算，不是 backend 真實 preview。
2. `目標型態` 只有表格內逐列 select，缺少快速設定當前列與一鍵套用到全部列的操作。

這次 follow-up 的目標很窄：

- 把 Step 3 preview 改成 backend-backed preview
- 保留 payload preview 為 display-only
- 新增 `目標型態` 快捷設定，但不重做 Step 3 autosave 主資料流

## 2. Problem Statement

前一輪 `/studio/v2` 全站對接已經完成 Step 1 live diagnostics、workspace autosave 與 settings backend wiring，但 Step 3 的 preview 仍停在本地假值。這造成兩個問題：

1. preview 顯示的結果不是 backend 真相
2. 使用者在大量 mapping 列上調整 `target_type` 的效率偏低

Step 3 已經有 workspace mapping autosave，也已經有既有 `POST /datalink/mappings/preview` 能力，因此這裡缺的是接線與操作收斂，不是再造第二套契約。

## 3. Goals

### Product goals

- 讓 `轉換管線預覽` 顯示 backend 真實 preview 結果
- 讓 `目標型態` 支援當前列快速設定與一鍵套用到全部列
- 保持 Step 3 畫面仍然是 review-first flow，而不是新增一塊獨立工作台

### Architecture goals

- 重用既有 `POST /datalink/mappings/preview`
- 不新增 V2-only preview endpoint
- 不重做 Step 3 autosave、mapping reducer、workspace CRUD 契約

### Delivery goals

- 只觸碰 Step 3 相關前端檔案與對應 spec/tasks
- 以最小測試面覆蓋 loading/success/error/stale response 與 target type 快捷操作

## 4. Non-Goals

- 不改 Step 2 / Step 4
- 不改 payload preview JSON 結構
- 不新增新的 workspace preview API
- 不把 `scale/offset` 快捷操作搬進 preview 卡
- 不重做既有 footer 的 bulk apply `scale/offset/target_type` 契約

## 5. Approved User Decisions

已確認的設計決策如下：

1. `目標型態` 快捷設定必須同時支援「當前列」與「一鍵套用到全部列」
2. `套用到全部列` 指的是目前 Step 3 畫面中的全部 `state.mappings`
3. preview 要改接真 API，不再接受本地 mock preview
4. payload preview 繼續保持 display-only

## 6. Approaches Considered

### Approach A — Reuse existing mapping preview API

右側 preview 直接重用既有 `POST /datalink/mappings/preview`，以前端 deterministic raw seed + draft transform pipeline 發 request。

- Pros: 最小改動、最快落地、無須新 backend surface
- Cons: 前端要補 debounce、loading/error、stale response guard

### Approach B — Add V2 workspace preview endpoint

新增 `/studio-v2/workspace/mappings/preview`，讓 Step 3 只走 workspace API。

- Pros: 邊界最乾淨
- Cons: 為同一 preview 能力再開一層，成本高且現在沒有必要

### Approach C — Keep local preview, add shortcuts only

- Pros: 成本最低
- Cons: 不符合「mock 要接 API」的核心要求

選擇 **Approach A**。

## 7. Design

### 7.1 Boundaries

這次只改動以下責任邊界：

- `Step3Mapping`
  - 組裝 selected row 與 preview container
- `TransformPreview`
  - 從純展示元件升級為帶 async preview state 的容器
- `PipelineSteps`
  - 由 backend preview response 驅動 steps 顯示
- `MappingRow`
  - 維持既有 inline select，不新增額外快捷操作
- Step 3 tests / i18n / spec artifacts

這次不改：

- `mappingReducer` 的主體資料結構
- `useStudioV2MappingAutosave`
- Step 4 database target flow

### 7.2 Live preview data flow

當有 selected row 時，preview 容器會：

1. 取得目前列對應的 deterministic raw seed
2. 以當前 draft `scale/offset/target_type` 組出 preview pipeline
3. debounce 後呼叫 `POST /datalink/mappings/preview`
4. 用 backend `step_results` / `final_value` 渲染管線顯示

這保留了目前 Step 3 的 deterministic sample 操作感，但把實際 transform truth 移回 backend。

### 7.3 Preview states

preview 卡需要三種明確狀態：

- `loading`
  - 顯示預覽計算中
- `success`
  - 顯示 backend 回來的 decode / scale / cast / final 結果
- `error`
  - 顯示 backend 回傳的錯誤訊息，不回退成本地假成功

沒有 selected row 時，維持現有 empty state。

### 7.4 Stale response guard

因為使用者可能連續修改 `scale/offset/target_type`，preview request 必須丟棄過期回應：

- 只採用最後一次 request 的結果
- 較早送出的 request 即使較晚回來，也不得覆蓋較新的 preview

這個邏輯應封裝在 preview 容器層，而不是 spread 到 reducer。

### 7.5 Payload preview boundary

`PayloadPreview` 仍然只顯示 JSON payload 結構。

這次 live preview request 不會改變：

- payload preview 的資料形狀
- payload preview display-only 的定位
- payload preview 不主動發 request 的邊界

### 7.6 Target type quick actions

快捷操作放在右側 preview 卡，而不是再往表格新增一列 UI。

操作如下：

- 常用 chips：`bool`、`int16`、`float64`、`string`
- 點 chip：
  - dispatch `updateMapping`
  - 只改當前 selected row 的 `target_type`
- 點 `套用到全部列`：
  - 沿用既有 `bulkApplyTransform`
  - 但 `fields` 只帶 `['target_type']`

這樣可以重用既有 autosave 與 reducer flow，不引入第二套批次更新機制。

## 8. Error Handling

- preview request 失敗時顯示 backend error
- preview error 不得偷偷保留上一個 success 畫面作為目前 truth
- payload preview 不因 preview error 而消失；它仍可顯示當前 draft payload 結構

## 9. Testing Strategy

最小必要驗證如下：

1. `step3-mapping.test.tsx`
   - preview loading / success / error
   - stale response guard
   - target type chip 更新當前列
   - `套用到全部列` dispatch `bulkApplyTransform` with `['target_type']`
2. 如 async 邏輯過重，再補一個小型 preview hook test
3. `cd frontend && npm run build`
4. `git diff --check`

若目前 backend preview contract 已足夠，這次不需要新增 backend product code。

## 10. Risks And Mitigations

- Risk: debounce 太短會產生多餘 preview request
  - Mitigation: 使用固定短 debounce，並以 stale response guard 補上 correctness
- Risk: preview 與 payload preview 職責混淆
  - Mitigation: 明確限制只有 transform preview 會發 request，payload preview 永遠 display-only
- Risk: 表格內既有 bulk apply 與 preview 卡 apply-all 語意看起來重複
  - Mitigation: preview 卡只提供 `target_type` apply-all；footer 仍保留 `scale/offset/target_type` 整體套用
