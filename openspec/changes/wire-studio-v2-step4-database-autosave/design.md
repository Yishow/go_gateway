## Context

Step 4 之後仍會保留資料庫 connector 與 target 編輯，但使用者已決定最後按鈕只做「第一次啟動 + 前往 runtime」，不是再次保存。因此 Step 4 的資料寫入必須先獨立成 autosave contract，才能避免 activation change 混進 data save。

## Goals / Non-Goals

**Goals:**

- 合法 connector / database target 編輯立即 autosave
- 不合法內容不覆蓋最後成功版本
- 多個 target 逐項保存、逐項回報錯誤
- 明確切開「資料 autosave」與「第一次啟動」

**Non-Goals:**

- 不在這個 change 內啟動 runtime
- 不改 Step 4 成功卡與 activation UX
- 不在這個 change 內處理 runtime 多設備頁

## Decisions

### Valid-only database autosave

connector 設定與 target row 只有在當前資料通過 validation 時才送出保存。invalid 編輯值保留在畫面上，不寫入後端。

### Per-target save isolation

database target row 逐項保存，單項失敗不阻塞其它 target 或 connector metadata 的成功保存。

### Database autosave remains pre-activation

Step 4 autosave 成功不代表設備開始跑。第一次啟動仍由後續 activation change 單獨處理，避免 data save 和 activation 再次耦合。

## Implementation Contract

- Behavior:
  - 合法 connector / target 編輯立即寫入後端
  - 不合法編輯不覆蓋最後成功版本
  - 單一 target save 失敗不阻塞其它 target
  - autosave 成功本身不會啟動 runtime
- Interface / data shape:
  - `GET /api/v1/datalink/studio-v2/workspace/database-config`
  - `PUT /api/v1/datalink/studio-v2/workspace/database-config`
  - `GET /api/v1/datalink/studio-v2/workspace/database-targets`
  - `PUT /api/v1/datalink/studio-v2/workspace/database-targets/:point_id`
  - database payload 至少區分 connector metadata、target rows、`save_state`
- Failure modes:
  - target validation failure 不清掉本地 row 值
  - autosave 成功時不得偷偷啟動 scheduler / runtime
- Acceptance criteria:
  - frontend Step 4 tests 覆蓋合法 autosave、不合法不存、逐項錯誤隔離
  - backend tests 覆蓋 workspace-scoped database config / target persistence
- Scope boundaries:
  - In scope: Step 4 connector/target autosave
  - Out of scope: first activation、runtime device switching
- Not complete if:
  - Step 4 還要等最後按鈕才真正寫入資料
  - autosave 會直接啟動 runtime
  - 單一 target 失敗會讓全部 target 一起失敗

## Risks / Trade-offs

- [Risk] target save 與 connector save 分開後，使用者可能短暫看到不一致狀態 → Mitigation: 前端清楚標示 row-level save state，並以 validation contract 阻擋不完整組合
- [Risk] Step 4 舊 commit 動畫仍暗示「此處才真正寫資料」 → Mitigation: 後續 activation change 會重寫文案與流程語意
