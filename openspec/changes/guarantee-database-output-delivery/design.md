## Context

database output 現在橫跨 Step 4 connector/target 設定、schema generate/ensure、runtime dbtarget writer、以及操作員最終在外部資料庫看到的結果。你的症狀「有在採集但完全沒寫入」說明目前這幾段之間缺乏共同的 delivery contract，而且 stale hidden mappings、connector readiness、writer flush outcome、schema ensure 成敗都可能在不同位置把資料吃掉。

這個 change 的重點是把 database output 當成一條 delivery pipeline 來規格化，而不是把 Step 4 設定畫面、schema ensure、writer 當成各自獨立的 feature。

## Goals / Non-Goals

**Goals:**

- 讓 schema ensure 與 writer 只處理 live、enabled、可解析的 mapping 集合。
- 讓 connector readiness、schema ensure outcome、write outcome 可被 Step 4 與 runtime 看見。
- 讓 collected value 到 DB write result 的失敗路徑可診斷。

**Non-Goals:**

- 不在這個 change 內重做整個 runtime dashboard truthfulness；那是另一個 change。
- 不把所有 diagnostics 都做成完整 audit history；那是 observability/audit 的範圍。
- 不新增新的 connector kinds；仍以既有 Studio V2 支援的 connector 為主。

## Decisions

### Delivery operates only on the live visible mapping set

schema ensure 與 writer 都以 live、enabled、可解析 tag 的 mapping set 為準，任何 stale hidden mapping 或 orphaned target 都要被隔離，而不是讓它們阻斷整條管線。

### Schema ensure and write path share the same filtered target projection

如果 schema ensure 用一套 target 集合、writer 又用另一套 target 集合，操作員仍會遇到「建表成功但寫不進去」或相反的情況。因此兩者必須共享同一套 filtered target projection。

### Writer persists last-write status and flush outcome

只靠 runtime counters 不足以解釋外部資料庫為什麼沒資料。writer 需要保留至少最近一次成功/失敗寫入與 flush outcome，讓上層 UI 和診斷 API 可以引用。

### Step 4 surfaces delivery truth, not only config truth

Step 4 不能只顯示 connector fields 和 target config。它還必須顯示這些設定是否真的 delivery-ready、schema-ready、以及最近一次寫入是否成功。這樣操作員才能在 setup 階段就看見 delivery truth。

## Implementation Contract

- Behavior:
  - hidden stale targets 不得阻擋 live schema ensure 或 live writer delivery。
  - schema ensure 與 writer 共用同一套 live target projection。
  - Step 4 能顯示 connector readiness、schema outcome、write outcome。
  - runtime/diagnostics 能指出 collected value 是在哪一段 delivery path 失敗。
- Interface / data shape:
  - database delivery 狀態至少要表達 connector readiness、last schema ensure status、last write status、last flush status、last failure reason。
  - target projection 需能辨識 live visible target 與 orphaned or stale target。
- Failure modes:
  - connector 不可用時，Step 4 顯示 readiness/problem，而不是只在 writer 時晚爆。
  - writer flush 失敗時，delivery truth 更新為 failed，不可只有 log。
- Acceptance criteria:
  - hidden stale target 不再阻斷 live activation/write。
  - write failure 後 Step 4 或 runtime surface 可見 failed delivery truth。
  - collected-but-not-written 的情境可被診斷到 DB write stage。
- Scope boundaries:
  - In scope: filtered mapping set、schema ensure/write consistency、delivery truth surfaces。
  - Out of scope: full audit history、runtime dashboard visual redesign、connector kind expansion。

## Risks / Trade-offs

- [隔離 stale hidden mappings 可能暴露舊資料髒污] → 這是預期效果，先隔離再由後續 change 處理清理策略。
- [writer status 落盤或持久化會增加狀態維護成本] → 但沒有這些狀態就無法回答「為什麼沒寫進去」。
- [Step 4 介面資訊變多] → 需要控制層級，避免把診斷全部塞進主要編輯表單。

## Migration Plan

1. 先統一 live target projection 與 stale target isolation 規則。
2. 再讓 schema ensure 與 writer 共用該 projection。
3. 最後把 delivery truth 帶到 Step 4 與 runtime-facing surfaces。

## Open Questions

- last write status 要存 connector scope、target scope、還是兩者都保留？
- orphaned stale mappings 是直接清除、標記、還是只在 projection 中隔離？
