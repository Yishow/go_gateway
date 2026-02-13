## Context

此變更聚焦在工業資料採集操作員的核心任務：
1) 選擇來源設備與位址；2) 將來源資料映射到 Tag；3) 驗證轉換與品質；4) 讓 Tag 寫入時序資料庫。

現況已有多種 UI 入口，但缺乏「單一主流程畫面」的強約束與狀態一致性。

## Goals / Non-Goals

- Goals
  - 將 Source → Memory Grid → Tag → DB 寫入建立成單一可視化主流程。
  - 讓每個操作狀態可被診斷（哪一段成功、哪一段失敗）。
  - 建立一致的 flow state machine：`draft`、`validated`、`active`、`error`。
  - 確保 a11y（鍵盤可操作、focus 可見、icon-only aria-label）。
- Non-Goals
  - 不更換後端儲存引擎與資料表 schema。
  - 不新增新協定驅動。
  - 不重寫 test-ui 測試工具頁。

## Decisions

- Decision 1: 採用 Flow-first Workspace
  - 以流程段落作為主要資訊架構，而非純 CRUD 分頁。
  - 四段固定：Source / Grid / Tag / Write Target。

- Decision 2: 使用「持續可視化上下文」
  - 選取 Memory Grid 位址後，右側/下方同步展示已連結 Tag 與寫入狀態。
  - 切換設備時保留可回溯的最近流程快照（不自動丟失草稿）。

- Decision 3: 引入 flow state machine
  - `draft`: 僅建立來源或映射，未驗證。
  - `validated`: 通過 pipeline 驗證與預覽。
  - `active`: 已啟用並持續寫入。
  - `error`: 任一段異常（連線、解碼、轉換、寫入）。

- Decision 4: 觀測資料標準化
  - 每段顯示 `latest_value`, `quality`, `timestamp`, `last_error`。
  - 錯誤必須定位到段落（source / transform / sink）。

## Risks / Trade-offs

- 風險：UI 複雜度提高，初期開發工作量較大。
  - 緩解：分階段導入，先做主流程可視化，再做進階交互優化。
- 風險：與既有頁面並存期間狀態一致性問題。
  - 緩解：以 Query key 與 shared selectors 做單一狀態來源。

## Migration Plan

1. Phase A：新增 flow frame 與 state model（不移除舊入口）。
2. Phase B：將 Memory Grid 選取、Tag 對應、Preview 診斷串接同屏。
3. Phase C：將 active/error 狀態回寫到主流程視圖與列表。
4. Phase D：完成驗收後，逐步弱化舊分散式入口。

## Open Questions

- 是否要在第一版即加入「批次多 Tag 對應」可視化，或先以單點位主流程為主？
- 是否需要在主流程上直接顯示 DB partition/寫入延遲指標？
