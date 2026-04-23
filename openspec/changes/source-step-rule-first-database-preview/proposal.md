## Why

`/studio?target=database&step=source` 目前已具備 source-rule 規劃、跨步驟 handoff 與 database downstream 語意，但實際操作時仍無法讓使用者快速看出「目前規劃了什麼」。Source step 的主敘事被大型 canvas、mode tabs 與過早的 handoff 語氣稀釋，導致 source rule 不夠像主角，也無法在這一步就預見 Tag review 下一步會看到的內容。

這個調整現在要做，因為 `/studio` 已明確是唯一主線，且 source → tag → output 的 rule-scoped contract 已存在；目前缺的是把既有能力重組成更清楚的 rule-first、database-aware 規劃體驗，而不是再增加另一套流程。

## What Changes

- 將 Source step 重構為以 active source rule 為中心的規劃工作台，讓使用者先看懂目前正在編哪一條 rule，再看地址覆蓋與下一步 handoff。
- 在畫布上方建立 active rule summary 前景區塊，左側 rule list 保留高亮與切換職責；summary 至少顯示起始位址、count、data type 與 planned / used / conflict 覆蓋摘要。
- 將 Source step 的 primary CTA 收斂為「套用規劃到畫布」，並在套用後直接回饋 rule 已成形與下一步 Tag review 將看到的摘要。
- 降低大型 canvas 與 mode tabs 的主導性，讓 canvas 轉為驗證與診斷 surface，而不是主要規劃入口；使用者應先從 summary 理解當前狀態，再進入格位檢查 conflict、unmanaged 與 live value。
- 在 Source step 內加入 rule-scoped Tag preview，支援以 hover / click rule 預覽下一步 Tag review 的候選結果。
- 讓 `target=database` 影響 Source step 的預設值與規劃語言，例如 data type、scale、naming 與 grouping hints，而不是把 output connector setup 細節提前帶入 Source step。
- 保持既有 cross-step contract：active rule、rule-scoped candidate snapshot、Tag / Database downstream continuity 不變。

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `datalink-ui`: 調整 `/studio` Source step 的資訊階層、主操作與跨步驟預覽，讓 source rule 的當前規劃與 Tag handoff 在 UI 上更可見。
- `datalink-workbench-desktop`: 調整 Source workspace 的 rule-first 結構、canvas 說明方式與 Source → Tag handoff 可預見性。
- `database-target-workbench`: 讓 database downstream 語意提前影響 Source step 的規劃預設與 grouping hints，但不提前引入 connector/schema/table setup。
- `source-rule-target-datatype`: 調整 target data type / scale / naming 在 database context 下的 Source-step 規劃預設與呈現方式。

## Impact

- 前端 workbench Source step 相關檔案：`MuiSourceCommandDeck.tsx`、`SourceCanvasSection.tsx`、`WorkbenchProvider.tsx`、`sourceRuleSelection.ts`。
- 前端 database grouping 與 downstream preview 相關檔案：`databaseGroupingSuggestions.ts`、相關 Tag / Database review surfaces。
- i18n 文案：`frontend/src/i18n/locales/zh-TW/common.json`、`frontend/src/i18n/locales/en/common.json`。
- 測試：Source step、Tag review continuity、database grouping / preview 相關 unit / integration tests。
