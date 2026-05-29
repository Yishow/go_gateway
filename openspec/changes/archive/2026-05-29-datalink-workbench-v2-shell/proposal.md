## Why

現有 `/studio` 主介面（`datalink-workbench-desktop`）已歷經多輪疊代，但 UI 結構與互動模式仍偏向過去的 datalink legacy 心智模型，導致以下痛點：

- **流程入口分散**：device / source-rule / mapping / output 等步驟視覺對齊度不足，新手難以建立完整的「device → source rule → tag → output」線性認知。
- **設計缺乏統一風格**：缺乏明確的 design tokens，工業儀表板必備的「點位記憶體」隱喻、深色高密度資訊呈現未被一致採用。
- **缺少安全的版型實驗管道**：任何 shell 層調整都會直接影響生產主入口，無法以低風險方式驗證新版型。

`new_prototype/` 已完成全新 4 步驟工作台原型（`app.jsx` + `shared.jsx` + `tweaks-panel.jsx` 等 8 個 jsx 檔，共 3,515 行），確立了 Dark Industrial Telemetry 風格。本變更先把 shell 與共用元件以 `/studio/v2` 並存路由落地，作為後續 step1/step2/step3/step4/settings 五個 change 的容器，舊 `/studio` 保留作為 fallback 直到 v2 全部成熟。

## What Changes

- 新增 `/studio/v2` 路由，掛載新的 `WorkbenchV2Shell` 元件；舊 `/studio` 與 `/datalink/workbench*` 行為與既有 `datalink-workbench-desktop` spec 一致，不動。
- 在 `frontend/src/features/datalink/workbench-v2/` 建立新版 feature 目錄，含設計 tokens、共用元件、shell layout、StepRail、SummaryRail、Tweaks panel、4 個空 step 骨架、空 Settings 骨架。
- 落地設計 tokens：dark mode 主背景 `#0b1220`、surface `slate-900/60` + `backdrop-blur`、primary `blue-500/600`、success `emerald-500`、warning `amber-500`、error `red-500`、UI 字體 `Inter`、mono 字體 `JetBrains Mono`、24px 點位格線底紋、漸層光暈背景。字型以 `@fontsource` 本地化（不走 CDN）。
- 提供共用元件：`Icon`（內聯 SVG 圖示集，禁用 emoji 作為結構性圖示）、`Button`、`Field`、`Input`、`Select`、`Textarea`、`Toggle`、`StatusChip`、`SectionCard`，全部具備鍵盤可達性與 dark/light 雙態（v2 只交付 dark）。
- Shell 結構：頂列（collapse 按鈕、品牌標識、breadcrumb、儲存草稿/取消按鈕、scheduler 狀態指示）+ 左側可收合 StepRail（4 步驟 + 設定入口、⌘B/Ctrl+B 切換、收合僅顯示 icon + tooltip）+ 中央 step 內容區（含 STEP 0X/04 標頭與進度橫條）+ 右側 SummaryRail（≥xl 顯示，可在 Tweaks 開關）+ 浮動 Tweaks panel（版面切換、示範流程的入口）。
- 4 個 step 與 settings 都以空骨架佔位，後續 5 個 change 會逐一替換為完整內容；本 change 不引入任何後端 API 呼叫，所有資料以 React in-memory state 撐起。
- i18n：所有 microcopy 預先放入 `frontend/src/i18n/locales/zh-TW/workbench-v2.json` 與 `en/workbench-v2.json`，符合 repo 既有 i18n 工作流。

## Non-Goals (optional)

- 不在本 change 接後端 API；資料模型僅以 mock state 撐起，與 `datalink-react-query` spec 的整合留給 phase 2 的 backend-wiring change。
- 不刪除舊 `/studio` 的 `DatalinkWorkbenchPage` 程式碼與相關 spec；舊主線維持運作。
- 不在本 change 改動 light mode；v2 僅落地 dark mode tokens。
- 不導入新的 UI library；僅使用 Tailwind utility classes，不引入 shadcn/ui 等第三方元件庫。

## Capabilities

### New Capabilities

- `datalink-workbench-v2-shell`：定義 `/studio/v2` 主路由、shell 版面、設計 tokens、共用元件契約、StepRail/SummaryRail/Tweaks panel 行為，以及與 `/studio` legacy 並存的 fallback 規則。

### Modified Capabilities

- `datalink-workbench-desktop`：新增「`/studio/v2` 並存主路由」要求，明訂 `/studio` 仍為當前正式入口、`/studio/v2` 為 v2 漸進落地容器；不變更 `/studio` 任何既有 Scenario。

## Impact

- Affected specs:
  - 新增 `openspec/specs/datalink-workbench-v2-shell/spec.md`
  - 修改 `openspec/specs/datalink-workbench-desktop/spec.md`
- Affected code:
  - New:
    - `frontend/src/features/datalink/workbench-v2/index.ts`
    - `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx`
    - `frontend/src/features/datalink/workbench-v2/shell/TopBar.tsx`
    - `frontend/src/features/datalink/workbench-v2/shell/StepRail.tsx`
    - `frontend/src/features/datalink/workbench-v2/shell/SummaryRail.tsx`
    - `frontend/src/features/datalink/workbench-v2/shell/TweaksPanel.tsx`
    - `frontend/src/features/datalink/workbench-v2/components/Icon.tsx`
    - `frontend/src/features/datalink/workbench-v2/components/Button.tsx`
    - `frontend/src/features/datalink/workbench-v2/components/Field.tsx`
    - `frontend/src/features/datalink/workbench-v2/components/StatusChip.tsx`
    - `frontend/src/features/datalink/workbench-v2/components/SectionCard.tsx`
    - `frontend/src/features/datalink/workbench-v2/components/Toggle.tsx`
    - `frontend/src/features/datalink/workbench-v2/components/inputs.tsx`
    - `frontend/src/features/datalink/workbench-v2/tokens.ts`
    - `frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css`
    - `frontend/src/features/datalink/workbench-v2/steps/Step1DevicePlaceholder.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/Step2RulePlaceholder.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/Step3MappingPlaceholder.tsx`
    - `frontend/src/features/datalink/workbench-v2/steps/Step4DatabasePlaceholder.tsx`
    - `frontend/src/features/datalink/workbench-v2/settings/SettingsPlaceholder.tsx`
    - `frontend/src/features/datalink/workbench-v2/state/useWorkbenchV2State.ts`
    - `frontend/src/features/datalink/workbench-v2/state/types.ts`
    - `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`
    - `frontend/src/i18n/locales/zh-TW/workbench-v2.json`
    - `frontend/src/i18n/locales/en/workbench-v2.json`
    - `frontend/tests/workbench-v2/shell.test.tsx`
  - Modified:
    - `frontend/src/App.tsx`（新增 `/studio/v2` 路由）
    - `frontend/src/main.tsx`（若需要載入 `@fontsource/inter` 與 `@fontsource/jetbrains-mono`）
    - `frontend/package.json`（新增 `@fontsource/inter` 與 `@fontsource/jetbrains-mono`）
- 不變更後端任何檔案；不變更 `cmd/test_ui/` 嵌入式靜態資源以外的 Go 程式碼。
