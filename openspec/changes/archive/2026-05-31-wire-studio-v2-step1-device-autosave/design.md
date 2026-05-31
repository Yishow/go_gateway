## Context

使用者已決定 `/studio/v2` 需要每一步都直接接後端，且「合法的先存，不合法的留在畫面上」。Step 1 是整條主線的入口，因此要先把設備層的 create / update / delete / ordering 契約從本地 reducer 拉成 workspace-scoped 真實資料。

## Goals / Non-Goals

**Goals:**

- Step 1 合法設備內容會自動存到後端
- 不合法內容不會覆蓋後端最後一次成功版本
- 多台設備逐台保存、逐台顯示錯誤
- 刪除已保存設備時同步刪除 V2 關聯資料

**Non-Goals:**

- 不在這個 change 內完成 Step 2 規則保存
- 不處理第一次啟動或 runtime 切換 UI
- 不把 `/studio` legacy device flows 合併進 V2

## Decisions

### Valid-only device autosave

前端保留每台設備的本地編輯值，但只有在設備表單合法時才觸發 create / update。這樣能同時滿足「自動存」與「不合法不寫入」兩個決策。

### Per-device save isolation

多台設備以每台獨立 mutation 與錯誤狀態處理，不使用整批 transactional save。這和使用者指定的「成功的先成立，失敗逐台回報」一致。

### Persisted device deletion cascade

在 V2 刪除已保存設備時，直接刪掉該設備與其 V2 關聯資料，而不是只從工作區隱藏。這個語意必須在 Step 1 階段就定清楚，避免後續 rule / mapping changes 各自猜測刪除結果。

### Workspace-owned order

device tab 排序需寫回 workspace，因為 runtime 後續要以 V2 排列順序決定預設打開哪台成功設備。

## Implementation Contract

- Behavior:
  - 合法設備建立後會立刻變成 workspace 的 persisted device
  - 合法設備修改後會立刻更新後端版本
  - 不合法設備修改不會寫入後端，但畫面保留使用者輸入並標示未存成功
  - 刪除已保存設備時，該設備與其 V2 關聯資料同步刪除
- Interface / data shape:
  - `GET /api/v1/datalink/studio-v2/workspace/devices`
  - `POST /api/v1/datalink/studio-v2/workspace/devices`
  - `PUT /api/v1/datalink/studio-v2/workspace/devices/:id`
  - `DELETE /api/v1/datalink/studio-v2/workspace/devices/:id`
  - `PUT /api/v1/datalink/studio-v2/workspace/device-order`
  - 前端每台設備至少區分 `draft-invalid`、`saving`、`saved`、`save-error`
- Failure modes:
  - 單台 save 失敗時，不清掉該台本地輸入值
  - 其它已合法設備的 save 不因一台失敗而回滾
  - delete cascade 失敗時，前端不得假裝設備已被完全移除
- Acceptance criteria:
  - frontend unit / integration tests 覆蓋合法自動存、不合法不存、逐台失敗顯示、刪除 cascade
  - backend handler / service tests 覆蓋 workspace-scoped create / update / delete / order
- Scope boundaries:
  - In scope: Step 1 device autosave、排序、刪除 cascade
  - Out of scope: Step 2 rules、Step 4 activation、runtime switch UI
- Not complete if:
  - 不合法內容仍會覆蓋後端最後一次成功版本
  - 一台 save 失敗會把其它合法設備一起擋住
  - 刪除設備後，V2 關聯資料還留在後端

## Risks / Trade-offs

- [Risk] 前端需要同時持有本地未存值與後端最後成功值 → Mitigation: 顯式保存 per-device save state，避免 reducer 混用單一欄位
- [Risk] delete cascade 若 scope 不清，後續 Step 2/3 可能重複實作 → Mitigation: 先把 cascade contract 寫在這個 change 的 design 與 spec
