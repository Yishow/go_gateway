## Context

目前 `/studio` 的 Source step 已有完整的 rule-scoped state、candidate snapshot 與 downstream continuity，但 UI 結構仍偏向以 canvas 與模式切換為主，而不是以 active source rule 為主。使用者回饋的核心不是缺功能，而是缺少一個清楚、可掃描、能預告下一步的規劃主視角。

現有系統已具備幾個重要約束：
- `crossStepContext.focusedRuleId` 與 `sourcePlanningState.selectedRuleId` 是跨步驟 continuity 的核心。
- Tag review 與 Database output 已依賴 rule-scoped candidate snapshot。
- database grouping semantics 已存在於 `databaseGroupingSuggestions.ts` 與 downstream row planner model。
- Source step 已可設定 target data type、scale、naming，但這些預設尚未依 output target context 調整。

因此本次設計重點不是重做資料模型，而是利用既有 state / hook / review model，把 Source step 重新排列成 rule-first、database-aware 的規劃工作台。

## Goals / Non-Goals

**Goals:**
- 讓 Source step 一眼看出目前 active rule 的規劃內容與成形狀態。
- 將 primary CTA 收斂為「套用規劃到畫布」，並提供明確完成回饋。
- 將 Tag review 的下一步結果在 Source step 內提前預告。
- 讓 `target=database` 影響 Source step 的預設值、文案與 grouping hints。
- 保持既有 Source → Tag → Database 的 rule-scoped contract 不變。

**Non-Goals:**
- 不新增新的跨步驟資料模型或第二套 route。
- 不把 Database connector / schema / table setup 移到 Source step。
- 不重寫 Tag review 或 Output planner 的核心 domain logic。
- 不變更 active rule continuity 的擁有者（仍由既有 Provider state 管理）。

## Decisions

### 1. 維持既有 active-rule contract，只重組 Source step UI 結構
- 決定：保留 `focusedRuleId` / `selectedRuleId` / candidate snapshot 架構，不引入新的 source-step 專屬 selection model。
- 原因：這些欄位已被 Tag / Output 使用，若改動 contract，影響面會從 UI 調整擴大成整條 flow 的 state migration。
- 替代方案：建立新的 `activePlanningRuleId` 或 source-step preview state。放棄原因是會與既有 cross-step rule continuity 重疊，增加同步成本。

### 2. 將 rule summary + preview 提升到 canvas 之前
- 決定：把 active-rule summary 與 Tag preview 放到 Source step 上半部／前景，canvas 降為 rule explanation surface。
- 原因：使用者主要痛點是看不到「目前規劃了什麼」，不是看不到格子本身。
- 替代方案：維持 canvas 為主，只在旁邊補更多提示。放棄原因是仍無法改變首要閱讀順序。

### 3. Tag preview 採用 rule-scoped candidate snapshot，而不是複製 Tag step state
- 決定：重用 `useSourceRuleCandidatesQuery`、`resolveActiveRuleId` 與現有 Tag review surface 的資料邏輯，做最小預覽版 preview。
- 原因：這樣能保證 Source preview 與 Tag step 實際看到的 scope 一致。
- 替代方案：在 Source step 內自行重算 preview。放棄原因是容易與 Tag step 的既有 contract 漂移。

### 4. `target=database` 僅影響 planning defaults 與 grouping hints
- 決定：當 output target 為 database 時，只調整 Source step 的 target data type / scale / naming defaults 與 grouping-oriented copy，不提前帶入 connector setup UI。
- 原因：使用者要的是 database-aware planning，不是把 Source step 變成 output config。
- 替代方案：把部分 connector/schema/table 設定前移。放棄原因是會污染 Step 2 主語意，並違反現有 workbench 分層。

### 5. mode tabs 改為次級控制，不再佔主導
- 決定：保留 `plan/live/link` 能力，但弱化其視覺階層，避免 mode tabs 成為最醒目的控制元件。
- 原因：使用者明確指出 mode tabs 搶戲，且目前最大問題是規劃摘要不足而非模式能力缺失。
- 替代方案：完全移除 mode tabs。放棄原因是現有 live/link 仍有價值，直接移除風險過高。

## Risks / Trade-offs

- [Rule summary 過重，造成首屏資訊堆疊] → 只顯示 active rule 必要摘要，細節仍可由 canvas / expanded surfaces 補充。
- [Source preview 與 Tag step 實際結果不同步] → 只重用既有 candidate snapshot hook 與 selection resolver，不做第二套推導。
- [Database-aware defaults 造成既有 protocol-driven 規劃習慣混亂] → 將 database context 限制為預設值與 hint，不強制覆寫使用者手動選擇。
- [弱化 mode tabs 影響既有 power users] → 保留功能、只調整層級與預設閱讀順序。
- [修改 Source step UI 時連帶破壞既有 tests] → 以現有 Source / Tag / Database tests 為主線補齊，避免先改 contract 再補測。

## Migration Plan

1. 先更新 OpenSpec requirements，明確規定 Source step 的 rule-first summary、Tag preview、database-aware defaults。
2. 先在 Source step 補 active-rule summary 與 primary CTA 調整，保持舊 contract 不變。
3. 接上 rule-scoped Tag preview。
4. 最後接入 database-aware defaults 與 grouping hints。
5. 以既有 Source / Tag / Database tests 驗證 continuity，再跑 frontend lint/test/build。

Rollback 策略：若新 UI 組織導致 Source → Tag continuity 混亂，可先保留 underlying state 與 hooks，僅回退 summary / preview surface 與 CTA 文案層，不需資料 migration。

## Open Questions

- rule hover preview 是否需要在 keyboard-only flow 中提供等價焦點行為，而不只依賴 pointer hover。
- active rule summary 要顯示到哪個粒度才足夠，而不會把首屏再度做成另一種資訊過載。
- database-aware naming / scale defaults 是否需要依 protocol 類型再做細分映射。