## Context

runtime dashboard 目前已具備 snapshot-first、stream-overlay、degraded fallback 等結構，但你反映的核心問題是「畫面上的數據不是你在 Step 1 到 Step 4 設的那份真相」。這通常不是單一 endpoint 壞掉，而是 UI 在資料缺口時仍用推測值、fallback 值、或其他裝置上下文補畫面，導致使用者看見的不是 backend-backed truth。

這個 change 的工作是把 runtime dashboard 拉回「寧可明確顯示沒有資料，也不可以造資料」。它依賴 runtime reconciliation change 提供 alignment/drift state，但專注在 UI 與 backend contract 的 truthfulness，而不是 runtime refresh 本身。

## Goals / Non-Goals

**Goals:**

- 移除 dashboard 對 fake/synthetic runtime data 的依賴。
- 讓 selected device context、snapshot、stream、degraded fallback 一致對齊 backend contract。
- 在沒有真實資料時顯示 truthful empty/degraded/error states。

**Non-Goals:**

- 不在這個 change 內解決 runtime refresh 本身；那是 reconcile-runtime-with-workspace 的範圍。
- 不在這個 change 內加入完整 diagnostics timeline；那是 observability/audit 的範圍。
- 不保證每個 panel 一定有資料；只保證沒有資料時不會造假。

## Decisions

### Dashboard truthfulness beats visual continuity

若 backend 沒有真實資料，dashboard 必須顯示 empty/degraded/error，而不是為了視覺完整性補 synthetic values。這會讓某些畫面看起來更空，但可避免 operator 被假健康狀態誤導。

### Selected device context never silently falls back

selected device 是 runtime dashboard 的主要上下文，不能因為它缺資料就默默切到其他 device 或 fleet-wide summary。這個決策會讓 missing-context 問題更明顯，但可避免跨 device 誤讀。

### Snapshot and stream failures stay distinct

snapshot 可用但 stream 壞掉，與 snapshot 根本沒有資料，是兩種不同失效模式。前者代表「目前只剩靜態最後真相」，後者代表「連最後真相都沒有」，UI 必須區分。

## Implementation Contract

- Behavior:
  - dashboard 只顯示 runtime backend 正式提供的欄位與狀態。
  - selected device 缺資料時顯示該 device 的 empty/degraded/error，不切到別的 device。
  - stream unavailable 時保留最後 truthful snapshot，但顯示 degraded live state。
- Interface / data shape:
  - runtime snapshot/stream 需能表示 empty、degraded、unavailable、stale 等狀態。
  - 前端 dashboard state 需保留 selected device truth、snapshot truth、stream truth、degraded fallback truth 的分層。
- Failure modes:
  - backend 缺欄位時，panel 顯示 unsupported/empty，不可補推測值。
  - stream 中斷時，UI 只能展示最後可驗證 snapshot，而不可捏造持續更新。
- Acceptance criteria:
  - selected device 沒資料時，不會顯示其他 device 的 summary。
  - stream failure 只顯示 degraded，而不是假 live data。
  - runtime dashboard tests 覆蓋 empty/degraded/stale cases。
- Scope boundaries:
  - In scope: truthful rendering、selected device truth、snapshot/stream/degraded distinction。
  - Out of scope: runtime refresh algorithm、audit history、write diagnostics。

## Risks / Trade-offs

- [畫面可能更常顯示 empty/degraded] → 這是故意讓真實缺口可見，而不是 regression。
- [前端需要移除既有 fallback convenience] → 初期會讓一些 panel 看起來「比較不漂亮」，但可信度更高。
- [backend contract 需要補狀態欄位] → 必須和 reconcile/observability changes 協調，不然前端無法分辨不同失效模式。

## Migration Plan

1. 先盤點 dashboard 現有 synthetic/fallback sources。
2. 補齊 backend empty/degraded state 表達。
3. 再調整前端 panels 只吃 backend-backed truth。

## Open Questions

- 某些 panel 在 missing field 時要顯示隱藏、空白、還是 diagnostic chip？
- degraded state 是否要分 stream unavailable 與 projection stale 兩種視覺語意？
