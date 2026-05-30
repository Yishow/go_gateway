## Context

Step 3 是把 rule-derived points 轉成 tags / mappings 的關鍵層。如果這一層還只是本地 state，Step 4 的 database targets、runtime 套用與可用性判斷都沒有可信的後端基礎。這個 change 把 mapping 列自動存獨立出來，不和 database target 或 activation 混做。

## Goals / Non-Goals

**Goals:**

- 合法 mapping 列立即 autosave
- 不合法 mapping 列不覆蓋最後成功版本
- 多列 mapping 逐列處理，保留 partial success
- 明確保存每列 mapping 的 save state

**Non-Goals:**

- 不在這個 change 內處理 Step 4 database target 保存
- 不處理 activation 或 runtime device switching
- 不改既有 mapping preview / validate API 的主體邏輯

## Decisions

### Valid-only mapping autosave

每一列 mapping 只有在 tag key、target type、transform 組合通過當前 validation contract 時才送出保存。

### Per-point mapping save isolation

每個 point 對應的 mapping row 分開保存與回報錯誤。使用者指定的 partial success 必須在這一層被固定下來，不能等到 Step 4 才補。

### Save-state visibility

前端每列都要顯示最後一次保存是否成功，否則使用者無法區分哪列只是本地暫存、哪列已經成為後端 truth。

## Implementation Contract

- Behavior:
  - 合法 mapping 列會立即寫入後端
  - 不合法 mapping 列只保留在本地畫面，不覆蓋最後成功版本
  - 單列失敗不會阻塞其它列保存
- Interface / data shape:
  - `GET /api/v1/datalink/studio-v2/workspace/mappings`
  - `POST /api/v1/datalink/studio-v2/workspace/mappings`
  - `PUT /api/v1/datalink/studio-v2/workspace/mappings/:id`
  - `DELETE /api/v1/datalink/studio-v2/workspace/mappings/:id`
  - mapping row view model 至少區分 `local_value`, `persisted_value`, `save_state`
- Failure modes:
  - validation failure 不清掉本地輸入值
  - save failure 不回滾其它已成功列
- Acceptance criteria:
  - frontend tests 覆蓋合法 autosave、不合法不存、逐列 save isolation
  - backend tests 覆蓋 workspace-scoped mapping CRUD 與 validation error
- Scope boundaries:
  - In scope: Step 3 mapping autosave 與 row save state
  - Out of scope: database target autosave、first activation、runtime page
- Not complete if:
  - mapping row save 仍要靠整批 commit 才落地
  - 不合法 row 仍覆蓋最後成功版本
  - 使用者看不出哪列沒存成功

## Risks / Trade-offs

- [Risk] mapping validation 若過重，可能導致頻繁 autosave 抖動 → Mitigation: 先用 debounce + existing validation API，避免每個 key stroke 都直接 hit backend
- [Risk] local/persisted row 雙狀態增加前端複雜度 → Mitigation: 用統一 row view model 固定欄位與 save_state
