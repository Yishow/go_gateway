## Context

你要求的終態不只是「流程跑得動」，而是「可觀察、可記錄」。在現況裡，即便 collector、runtime、writer 真的有在跑，操作員也很難回答三個基本問題：
1. 現在卡在哪一段？
2. 這是剛剛才壞，還是之前就壞？
3. 最近哪一次設定變更或 activation 導致這個狀態？

這表示系統目前缺兩種不同層級的資訊：一種是診斷導向的 flow diagnostics，另一種是時間序列導向的 audit history。這兩者不能互相取代。

## Goals / Non-Goals

**Goals:**

- 提供 collector -> mapping -> runtime -> database delivery 的 stage-oriented diagnostics。
- 提供 workspace setup、activation、runtime transition、delivery failure 的結構化 audit history。
- 讓 runtime dashboard 與 workspace surfaces 可直接使用 diagnostics/audit summary。

**Non-Goals:**

- 不把所有原始 server logs 都搬進前端。
- 不在這個 change 內重新定義 runtime dashboard truthfulness 或 database delivery 契約，只消費它們的結果。
- 不追蹤每一次前端鍵盤輸入；只追蹤 persisted、activation、runtime、delivery 相關的重要事件。

## Decisions

### Diagnostics use a stage-oriented model

runtime diagnostics 不應只是「某個 error count +1」，而要能表示資料在 collector、mapping、runtime projection、database delivery 哪一段停住。這樣操作員才能真正定位故障位置。

### Audit history records state transitions and delivery-impacting events

audit history 聚焦 persisted change、activation、runtime transition、delivery failure 這些會改變系統狀態或操作結果的事件，而不是把所有 ephemeral UI 互動都寫進去。這能保持 audit trail 可用、可查、可理解。

### Operator surfaces consume summaries, not raw logs

前端需要的是 latest failure context、recent activation history、recent successful/failed delivery，而不是整份 zap log。raw logs 可以留給後端維運，但 operator-facing surface 必須用 summary + drill-down model 呈現。

## Implementation Contract

- Behavior:
  - runtime diagnostics 能說出 selected scope 最近成功/失敗時間與失敗階段。
  - workspace audit history 能查到最近 activation、重要 persisted change、runtime transition、delivery failure。
  - runtime dashboard 與 Studio V2 workspace 可顯示這些 diagnostics/audit summaries。
- Interface / data shape:
  - diagnostics 至少包含 scope、last_success_at、last_failure_at、failure_stage、failure_reason。
  - audit history 至少包含 workspace_id、event_type、result、scope、occurred_at，以及能關聯 activation 或 persisted change 的欄位。
- Failure modes:
  - 若 diagnostics 來源暫時 unavailable，UI 顯示 diagnostics unavailable，而不是空白當正常。
  - 若 audit write 失敗，不得阻斷主要功能，但要有後端記錄與 degraded indicator。
- Acceptance criteria:
  - runtime dashboard 看得到最新失敗原因與時間。
  - Studio V2 看得到最近 activation 與重要 persisted change history。
  - 端到端失敗時，operator 能回答失敗發生在 collector、mapping、runtime、或 database delivery 哪一段。
- Scope boundaries:
  - In scope: diagnostics stage model、audit history persistence、operator-facing summaries。
  - Out of scope: full raw log viewer、SIEM integration、external notification routing。

## Risks / Trade-offs

- [新增 audit/diagnostics 狀態會增加寫入量] → 以重要事件與摘要為主，避免無界 log copy。
- [若 diagnostics 與 audit 模型設計過粗，仍然無法定位] → 先保證 stages 與 timestamps 足夠回答主要問題，再逐步增補細節。
- [UI 顯示更多錯誤狀態會增加壓力感] → 這是必要成本，因為現在真正的問題是失敗不可見。

## Migration Plan

1. 先定義 diagnostics stage model 與 audit event model。
2. 再接 runtime、activation、db delivery 的事件來源。
3. 最後把 summaries 接到 runtime dashboard 與 Studio V2 workspace。

## Open Questions

- audit history 是否需要 actor/session identity，還是先只記錄 workspace/system scope？
- diagnostics summary 要支援 point scope drill-down 到什麼程度才算足夠？
