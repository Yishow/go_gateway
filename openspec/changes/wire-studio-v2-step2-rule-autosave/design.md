## Context

Step 1 只把設備變成 persisted data 還不夠；V2 的 Step 2 才是後續 point、tag、mapping、runtime 行為的上游。如果 Step 2 仍然是本地 state，使用者就看不出哪些規劃已真正進入後端，也無法把規則層與 workspace / device 邊界固定下來。

## Goals / Non-Goals

**Goals:**

- 讓 Step 2 規則內容合法即自動存
- 不合法規則不覆蓋最後成功版本
- 多條規則逐條處理，不因單條失敗阻塞其它條目
- 規則與所屬設備的 workspace 關係可重讀

**Non-Goals:**

- 不在這個 change 內完成 Step 3 mapping autosave
- 不處理 Step 4 第一次啟動
- 不在這個 change 內重做 candidate apply flows

## Decisions

### Valid-only rule autosave

Rule editor 只在規則內容滿足後端最小合法條件時才送出保存。暫時不合法的編輯值保留在畫面上，但後端仍以最後一次成功版本為準。

### Per-rule save isolation

多條規則以逐條 mutation 處理，單條失敗只影響該條規則的 save state。這對應使用者要求的 partial success。

### Workspace-scoped rule ownership

每條 Step 2 規則都必須能回到「屬於哪個 V2 workspace、哪台設備」的真實關係。否則多設備規劃一旦重整或重啟，前端無法判斷哪條規則屬於哪台設備。

## Implementation Contract

- Behavior:
  - 合法規則建立/修改後立刻寫入後端
  - 不合法規則保留本地值並標示未存成功
  - 其它合法規則不因單條失敗而被阻塞
  - 重整後可重建 workspace 內的規則與 device 關聯
- Interface / data shape:
  - `GET /api/v1/datalink/studio-v2/workspace/source-rules`
  - `POST /api/v1/datalink/studio-v2/workspace/source-rules`
  - `PUT /api/v1/datalink/studio-v2/workspace/source-rules/:id`
  - `DELETE /api/v1/datalink/studio-v2/workspace/source-rules/:id`
  - response 至少包含 `id`, `device_id`, `workspace_id`, `revision_id`, `save_state`
- Failure modes:
  - validation failure 不覆蓋最後成功版本
  - workspace/device 關聯錯誤時，不偷偷把 rule 掛到別的 device
- Acceptance criteria:
  - frontend tests 覆蓋合法 autosave、不合法不存、逐條錯誤隔離
  - backend tests 覆蓋 workspace-scoped rule CRUD 與 device ownership
- Scope boundaries:
  - In scope: Step 2 規則 autosave 與 workspace/device ownership
  - Out of scope: mapping autosave、activation、runtime 切換
- Not complete if:
  - 不合法規則仍會覆蓋後端最後成功版本
  - 重整後規則失去所屬設備
  - 單條保存失敗會讓其它合法規則一起停住

## Risks / Trade-offs

- [Risk] 規則合法性若與既有 `/studio` 條件不一致，可能造成雙軌驗證 → Mitigation: 直接重用既有 source rule validation/service contract
- [Risk] 規則 save 與 candidate snapshot 更新時序不清 → Mitigation: 先在這個 change 固定 CRUD 與 revision 保存，candidate/apply 仍由既有 source-rule runtime contract 管理
