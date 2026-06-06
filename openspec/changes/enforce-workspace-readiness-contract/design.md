## Context

目前 Studio V2 對 readiness 的判斷分散在 device validity、source-rule activation readiness、mapping/database handler validation、以及 activation handler 自己的保險步驟。這些判斷點缺少共同的 issue model，導致有些問題能保存卻不能啟動，有些問題直到 runtime 或 write path 才出現，前端也無法用一致方式告知操作員「現在缺什麼」。

這個 change 的角色是把 readiness 從零散條件拉成正式契約：用 persisted state 算出一份 normalized issue set，供 shell、summary、activation、live apply 共用。它不直接處理 runtime refresh 或 database write 成功與否，而是先保證系統在啟動前能準確說出是否 ready。

## Goals / Non-Goals

**Goals:**

- 以 persisted Step 1 到 Step 4 資料建立同一套 workspace readiness issue model。
- 讓 activation 與 live apply 使用同一套 blocking/warning 規則。
- 讓前端 shell 與 Step 4 activation surface 顯示同一份 readiness summary。
- 讓操作員在 activation 前就知道缺的是哪個 step、哪個 device、哪種 issue code。

**Non-Goals:**

- 不在這個 change 內解決 runtime refresh 細節或 database writer flush 細節。
- 不把所有 warning 都升成 blocker；只處理 readiness 契約與 issue 分級。
- 不在這個 change 內引入完整 audit history；那是 observability/audit change 的範圍。

## Decisions

### Backend readiness snapshot owns activation gating

activation 與 live apply 的 gating 以 backend readiness snapshot 為準，而不是前端自行拼湊 validity。這樣才能保證重載後、不同入口、不同前端 surface 的結論一致。

### Blocking and warning issues are first-class

不是所有問題都該阻止 activation。設計上將 issue 分成 blocking 與 warning，兩者都要顯示，但只有 blocking 會阻止 activation 或 live apply。這能避免把所有診斷都擠成 fatal，又能保留操作決策空間。

### Readiness evaluates step coverage and downstream integrity

readiness 不只看 Step 1 到 Step 4 欄位是否填滿，還要看 rule-derived downstream integrity，例如 point/tag/mapping/database target 是否存在、是否對應當前 persisted identity。否則只做表單層 validity 仍會把問題留到 runtime 晚爆。

### Frontend consumes the same readiness contract everywhere

shell、summary rail、Step 4 activation button、live apply 結果都消費同一份 readiness summary。前端不再各自發明 validity rules，而是只負責顯示與引導。

## Implementation Contract

- Behavior:
  - workspace 的 persisted state 會被評估成 readiness summary，包含 issue code、severity、scope、step。
  - activation 遇到 blocking readiness issues 時直接拒絕，並回傳 normalized issues。
  - live apply 對會引入 blocking readiness issue 的 persisted change 也必須拒絕或 defer。
  - shell 與 Step 4 都能顯示同一份 readiness 結果。
- Interface / data shape:
  - backend 提供 workspace readiness summary 與 issue list。
  - issue 至少包含 code、severity、step、scope、message。
  - activation/live apply 回應能帶回 blocking issues，不再只回 generic failure。
- Failure modes:
  - 若 readiness 計算失敗，前端只可顯示 readiness unavailable/degraded，不可假設 workspace ready。
  - 若 issue 是 warning，前端不得把它誤標成 blocking；若 issue 是 blocking，前端不得允許 activation 假通過。
- Acceptance criteria:
  - incomplete workspace 會在 activation 前看見 blocking issue。
  - live apply 遇到 blocking issue 會被拒絕或 defer，且帶出同一個 issue code。
  - shell/summary/Step 4 顯示相同 readiness summary。
- Scope boundaries:
  - In scope: readiness issue model、activation/live apply gating、front-end readiness display。
  - Out of scope: runtime refresh algorithm、database write execution、audit history。

## Risks / Trade-offs

- [Issue model 一旦公開就需要穩定] → 以 normalized code 為主、message 可調整，降低前後端耦合。
- [更多 blocking 可能讓 activation 看起來更難成功] → 這是刻意把晚爆轉成早爆，降低 silent failure。
- [不同 subsystem 對「warning」理解可能不同] → 先由 workspace readiness 統一定義 severity，再逐步細化。

## Migration Plan

1. 先定義 readiness summary 與 issue model。
2. 再讓 activation 與 live apply 都改走這個 gating。
3. 最後讓前端 shell、summary、Step 4 消費 readiness summary。

## Open Questions

- 某些 database-only 問題在沒有啟用 database sink 的 workspace 中是否應降級為 warning？
- readiness summary 是否要保存最近一次 successful evaluation timestamp？
