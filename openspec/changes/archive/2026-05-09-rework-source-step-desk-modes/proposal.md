## Why

目前 Source step 雖然已經導入 rule-first summary、Tag preview 與 database-aware hints，但 `巡檢 / 建造 / 分診` 三個 mode 的視覺語意仍然不夠清楚：巡檢容易只被理解成看 point 畫布、建造模式讓規劃規則層吃掉主要視窗、分診則和巡檢只剩排版差異。這讓 operator 很難一眼理解「現在是在檢查規劃、編輯規劃，還是處理衝突」，也削弱了 Source -> Tag -> Output 的主線敘事。

現在需要把 Source step 重整成共享同一條規劃主線、但明確區分三種工作姿勢的工作桌，讓規劃摘要、handoff 與 diagnostics 各自回到正確層級。

## What Changes

- 重構 Source step 的 desk mode 語意，讓 `巡檢` 明確代表「先讀 active rule summary，再用畫布驗證規劃」。
- 重構 `建造` 版面，讓 active rule editor 成為主要工作區，但不再壓縮到幾乎看不到下方畫布。
- 重構 `分診` 模式，讓它聚焦 conflict / blocked / unmanaged 項目與修復動作，而不是巡檢的排版變體。
- 固定 Source step 的共同骨架：active rule summary、handoff strip、database-aware hints 與 primary CTA 在 mode 切換時不丟失。
- 抽離 mode-specific panel 元件，降低 `SourceCanvasSection.tsx` 的版面責任，並補齊對應的 mode-contract 測試與文案。

## Capabilities

### New Capabilities

- `source-step-desk-modes`: 定義 Source step 中巡檢、建造、分診三種工作模式的目的、共享骨架與可辨識行為。

### Modified Capabilities

- `datalink-ui`: 調整 guided workflow 對 Source step 的敘事，要求 mode 切換不得丟失 active rule、Tag handoff 與 diagnostics 上下文。
- `datalink-workbench-desktop`: 調整 desktop workbench shell 中 Source workspace 的區域優先順序與 mode-specific layout 行為。

## Impact

- Affected specs: `source-step-desk-modes`, `datalink-ui`, `datalink-workbench-desktop`
- Affected code:
  - New: `frontend/src/pages/datalink/workbench/SourceRuleLayerPanel.tsx`, `frontend/src/pages/datalink/workbench/SourceTriagePanel.tsx`
  - Modified: `frontend/src/pages/datalink/workbench/MuiSourceCommandDeck.tsx`, `frontend/src/pages/datalink/workbench/MuiWorkbenchSourceStyles.tsx`, `frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx`, `frontend/src/pages/datalink/workbench/SourceStepRuleSummary.tsx`, `frontend/src/i18n/locales/zh-TW/common.json`, `frontend/src/i18n/locales/en/common.json`, `frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx`, `frontend/src/pages/datalink/workbench/__tests__/MuiSourceIncidentDesk.reopen.test.tsx`, `frontend/tests/unit/pages/datalink/workbench-source-step.test.tsx`
  - Removed: none
