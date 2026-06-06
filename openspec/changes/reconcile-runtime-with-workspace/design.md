## Context

Studio V2 目前已經有 activation、runtime service、live config apply、runtime dashboard 等多個入口，但缺少一個共同的 runtime workspace projection 契約。結果是保存成功不代表 runtime 已更新，activation 成功不代表 restart 之後還會以同樣投影恢復，dashboard 也無法區分顯示的是最新投影還是舊投影。

這個 change 的焦點不是 readiness，也不是 observability，而是把 runtime 對 workspace 的關係講清楚：runtime 只能跑 persisted workspace projection；任何影響 runtime 的 persisted change 都要有 reconcile outcome；任何 drift 都要被顯式暴露。

## Goals / Non-Goals

**Goals:**

- 建立 persisted workspace -> runtime projection 的單一來源契約。
- 讓 activation、restart、live apply 共用同一套 projection/reconcile 流程。
- 讓 runtime API 與 workspace runtime view 能暴露 projection alignment/drift。

**Non-Goals:**

- 不在這個 change 內補 database write diagnostics 或 audit history。
- 不在這個 change 內改造 Step 1 到 Step 4 的 reload recovery；那屬於 setup-state change。
- 不保證所有 persisted change 都能立即 hot apply；允許明確的 deferred 或 restart-required outcome。

## Decisions

### Workspace projection is the only runtime source of truth

runtime 只能從 persisted workspace projection 建立或更新，不可從前端 session state、臨時 reducer state、或單一 handler 的局部 patch 推導真相。這保證 activation、restart、runtime restore 的語意一致。

### Reconciliation is targeted and idempotent

reconcile 不該每次都全量重建所有 runtime state。設計上允許 device、rule、mapping、database target 各自 targeted reconcile，但要求重複請求仍收斂到同一結果，不重複註冊或殘留舊 binding。

### Live apply returns explicit reconcile outcome

保存成功不代表 runtime 已更新，因此 live apply 的回應必須攜帶 reconcile outcome。這可以讓前端明白顯示「已套用」、「已延後」、「需要 restart」等狀態，避免成功訊息誤導操作員。

### Runtime surfaces projection drift instead of hiding it

runtime dashboard 與 workspace runtime view 不再假設 snapshot 就代表最新配置，而是直接顯示 projection alignment/drift。這會讓 UI 在 drift 存在時比較「刺眼」，但這比繼續展示假同步更安全。

## Implementation Contract

- Behavior:
  - activation 與 restart 會從 persisted workspace projection 建立相同 runtime scope。
  - 任何影響 runtime 的 persisted change 都會回傳 reconcile outcome。
  - runtime snapshot/stream 與 workspace runtime view 都能表示 projection 是否對齊最新 persisted state。
- Interface / data shape:
  - runtime projection 需有穩定的 version/alignment/deferred state 表達。
  - live apply、activation、runtime snapshot、runtime stream 至少要能表達 aligned、deferred、restart-required、stale 或等效狀態。
- Failure modes:
  - 若 targeted reconcile 失敗，runtime 保留舊 projection，但 API 必須回報 stale/deferred 狀態。
  - 若 restart 後無法重建 projection，runtime 必須回報 degraded，而不是假裝沿用上次記憶體狀態。
- Acceptance criteria:
  - activation 與 restart 對同一 persisted workspace 產生一致 projection。
  - live apply response 能區分 immediate reconcile 與 deferred/restart-required。
  - dashboard/runtime view 能看到 projection drift。
- Scope boundaries:
  - In scope: projection model、reconcile outcome、runtime alignment surfacing。
  - Out of scope: write failure observability、audit history、reload recovery UI。

## Risks / Trade-offs

- [更多 runtime 狀態類型] → 前後端要一起理解這些狀態，但換來較少的假成功。
- [targeted reconcile 複雜度提高] → 必須靠 idempotent tests 控制回歸風險。
- [顯示 stale/deferred 可能讓 UI 更常出現警示] → 這是設計目標，因為 drift 本來就不該被隱藏。

## Migration Plan

1. 定義 persisted workspace projection 與 reconcile outcome。
2. 讓 activation、restart、live apply 都改走 projection/reconcile 契約。
3. 最後讓 runtime API 與前端 surface 顯示 alignment/drift。

## Open Questions

- projection version 要採遞增 revision、hash，還是 workspace updated_at-based fingerprint？
- 某些 mapping/database target 變更是否允許 background reconcile，而不必中斷 collector？
