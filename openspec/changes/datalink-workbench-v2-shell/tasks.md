## 1. 依賴與字型

- [ ] 1.1 安裝 `@fontsource/inter`（400/500/600/700）與 `@fontsource/jetbrains-mono`（400/500/600），在 `frontend/package.json` 寫入 dependency，落地設計決策「字型：`@fontsource/inter` + `@fontsource/jetbrains-mono` 本地化」。**行為**：執行 `npm install` 後 `node_modules/@fontsource/inter/400.css` 與 `node_modules/@fontsource/jetbrains-mono/400.css` 可被解析。**驗證**：`cd frontend && npm install && node -e "require.resolve('@fontsource/inter/400.css'); require.resolve('@fontsource/jetbrains-mono/400.css')"` 結束碼為 0。
- [ ] 1.2 在 `frontend/src/main.tsx` import Inter 與 JetBrains Mono 的所需字重 CSS，確保 v2 page 載入時字型可用且 fallback 為 `system-ui, sans-serif`。**行為**：v2 page 文字 computed font-family 第一個解析到 `Inter`，mono 區塊解析到 `JetBrains Mono`。**驗證**：`shell.test.tsx` 中新增測試 `font-loading` 斷言 `document.fonts` 包含 `Inter` 與 `JetBrains Mono` family（jsdom 用 mock fonts；若 jsdom 不支援則在 `frontend/tests/workbench-v2/shell.test.tsx` 改以斷言 `<link>` / inline style 存在替代）。

## 2. Design tokens scope 與基礎樣式

- [ ] 2.1 建立 `frontend/src/features/datalink/workbench-v2/tokens.ts`，匯出 `COLORS`、`FONTS`、`RADII`、`SHADOWS`、`ANIMATIONS` 常數對齊原型，落地設計決策「Feature 目錄結構：`frontend/src/features/datalink/workbench-v2/`」的入口檔案。**行為**：其他 v2 模組可以 `import { COLORS } from '../../tokens'` 取得 token；`COLORS.bg === '#0b1220'`、`COLORS.success === '#10b981'`。**驗證**：新增 `frontend/tests/workbench-v2/tokens.test.ts`，斷言關鍵 token 值；`npm run test workbench-v2/tokens` 通過。
- [ ] 2.2 建立 `frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css`，定義 CSS variables、keyframes（`pulseDot`、`sweepIn`、`shimmer`）與 utility classes（`.bg-canvas`、`.bg-grid`、`.pulse-dot`、`.sweep-in`、`.shimmer`、`.mem-cell`、`.kbd`、`.chip`、`.field`、`.label`），所有規則用 `[data-workbench-v2] ...` 範圍，落地設計決策「Design Tokens 落地：CSS variables + Tailwind theme.extend」。**行為**：CSS variables 只在 `[data-workbench-v2]` subtree 內生效（**Design tokens scope** 對齊）。**驗證**：在 `shell.test.tsx` 中渲染 v2 page，使用 `getComputedStyle` 確認 `.bg-canvas` 在 v2 內回傳 `#0b1220` 背景；另在純 wrapper 外的測試確認 token 不生效。
- [ ] 2.3 在 `DatalinkWorkbenchV2Page` root element 設定 `data-workbench-v2="true"`，並 import `workbench-v2.css`。**行為**：v2 page DOM root 帶 `data-workbench-v2` 屬性，全頁套用 canvas 背景。**驗證**：`shell.test.tsx` 斷言 `screen.getByTestId('workbench-v2-root')` 帶有 `data-workbench-v2` 屬性。

## 3. 共用元件移植（Shared component library for v2）

- [ ] 3.1 建立 `components/Icon.tsx`，匯出 `Icon` 元件與 `IconName` 型別，內建 23 個 SVG paths（device/rule/map/db/check/chevron/plus/play/refresh/cable/sliders/tag/table/alert/info/spark/arrow/close/save/eye/bolt/flow/grid），落地設計決策「共用元件移植：1:1 對應原型 + TypeScript 化」。**行為**：`<Icon name="device" />` 渲染為 SVG 元素並帶 `aria-hidden="true"`。**驗證**：`frontend/tests/workbench-v2/components.test.tsx` 中 `describe('Icon')` 斷言 SVG 結構與 ARIA 屬性，並斷言所有 23 個 name 都不會渲染 emoji。
- [ ] 3.2 建立 `components/Button.tsx`，支援 `variant`（primary/secondary/ghost/danger/success）+ `size`（sm/md/lg）+ `icon` + `disabled`。**行為**：disabled 時不觸發 `onClick`；keyboard Enter 觸發 click。**驗證**：`components.test.tsx` 中 `describe('Button')` 斷言 disabled props、keyboard event、icon render。
- [ ] 3.3 建立 `components/inputs.tsx`，匯出 `Input`、`Select`、`Textarea` 元件，全部 `forwardRef`。**行為**：value/onChange controlled；keyboard focus ring 出現（`focus:ring-blue-500/20`）。**驗證**：`components.test.tsx` 中對 `Input` 渲染控制 value 並觸發 onChange。
- [ ] 3.4 建立 `components/Field.tsx`，支援 `label` + `hint` + `error` + `required` + `children`；required 時 label 出現 `*` 紅色標記。**行為**：error 存在時隱藏 hint 並顯示 `alert` icon + 紅字。**驗證**：`components.test.tsx` 斷言 required asterisk、error 樣式優先於 hint。
- [ ] 3.5 建立 `components/Toggle.tsx`，使用原型精確像素位移實作（避免 dot 視覺溢出），帶 `role="switch"`、`aria-checked`、可選 `label`、`size`（sm/md）。**行為**：`<Toggle checked={false} onChange={fn} />` 點擊呼叫 `onChange(true)`；ARIA 狀態同步切換（**Shared component library for v2** 的 `Toggle exposes ARIA switch role` Scenario）。**驗證**：`components.test.tsx` 斷言 ARIA 屬性與 onChange 行為。
- [ ] 3.6 建立 `components/StatusChip.tsx`，支援 `tone`（neutral/info/success/warning/error/draft） + `dot`（true/false） + `children`。**行為**：info/success tone 帶 `pulse-dot` 動效。**驗證**：`components.test.tsx` 斷言每個 tone 對應 className 與 dot 顯示。
- [ ] 3.7 建立 `components/SectionCard.tsx`，支援 `title` + `subtitle` + `icon` + `aside` + `children` + `className` + `contentClassName`。**行為**：title 或 aside 任一存在時 header 區塊顯示；children 渲染在 content 區塊（**Shared component library for v2** 的 `SectionCard renders header and content` Scenario）。**驗證**：`components.test.tsx` 斷言 header/content 結構。
- [ ] 3.8 匯出 barrel `components/index.ts`，把所有共用元件與 props 型別 re-export。**行為**：外部可 `import { Button, Field, StatusChip } from '../../components'`。**驗證**：`tsc --noEmit` 在 `cd frontend && npm run typecheck` 中通過。

## 4. State 層

- [ ] 4.1 建立 `state/types.ts`，定義 `WorkbenchV2State` 與所有子型別（`Device`、`Rule`、`Point`、`Mapping`、`DbConnector`、`DbTarget`、`Settings`、`CommitState`），對齊 `new_prototype/docs/_overview.md` 的 State shape。**行為**：型別嚴格，禁止 `any`。**驗證**：`tsc --noEmit` 通過；新增 `state/types.test-d.ts` 以 type-level 斷言 `WorkbenchV2State['view']` 為 `'flow' | 'settings'`。
- [ ] 4.2 建立 `state/useWorkbenchV2State.ts`，封裝 `useReducer` 與 action（`setView`、`setCurrent`、`completeStep`、`toggleSidebar`、`toggleSummaryRail`），預設值帶一筆 `PLC-生產線-01` device 與一筆 `Holding Registers` rule，落地設計決策「State 管理：本地 React state + 可序列化形狀」。**行為**：state 變更全部 immutable；reducer 預設值對齊 `new_prototype/app.jsx` 的 `useState` 初值。**驗證**：`frontend/tests/workbench-v2/state.test.ts` 對 reducer 寫單元測試（toggleSidebar、completeStep、setCurrent）。

## 5. Shell 五大區塊（Shell layout regions）

- [ ] 5.1 建立 `shell/TopBar.tsx`：collapse button、品牌標識、breadcrumb、儲存草稿/取消按鈕、scheduler 狀態指示。**行為**：點 collapse button 呼叫 `onToggleSidebar`；breadcrumb 顯示 `Datalink › Workbench › {step.title | "設定"}`（**Shell layout regions** 的 `Default render of all regions` Scenario）。**驗證**：`shell.test.tsx` 斷言 sticky 行為、breadcrumb 文字、collapse button click handler。
- [ ] 5.2 建立 `shell/StepRail.tsx`：4 個 step + 分隔線 + 設定入口；支援 `collapsed`、`view`、`current`、`completed`、`onJump`、`onSwitchView` props。**行為**：reachability rule（**Step rail navigation** 的 `Initial reachability` + `Step completion unlocks the next step`）；當 collapsed=true 只顯示 icon + tooltip（**Shell layout regions** 的 `Step rail collapsed state`）；點設定切換 view（**Step rail navigation** 的 `Settings entry switches view mode` + `Returning from settings to flow`）。**驗證**：`shell.test.tsx` 多個 it block 對應每個 Scenario。
- [ ] 5.3 建立 `shell/SummaryRail.tsx`：顯示 devices / rules / mappings / db 摘要卡片；viewports < 1280px 時不渲染（**Shell layout regions** 的 `Summary rail hidden at narrow viewports`）。**行為**：根據 state 推導摘要數字（devices count、tested count、rule count、enabled point count、mapping count、db connector 摘要）。**驗證**：`shell.test.tsx` 用 jsdom 模擬 viewport width 1024px 斷言 SummaryRail 不在 DOM。
- [ ] 5.4 建立 `shell/TweaksPanel.tsx`：浮動面板，dev mode 或 `localStorage.WBV2_TWEAKS === '1'` 顯示；含「收合側邊欄」「顯示右側摘要欄」toggle 與「重置流程」「加入第二台設備」「加入第二條規則」「跳到 Step 4」四個示範按鈕（後三者 onClick 留 TODO 註解 + console.warn `'pending follow-up change'`），落地設計決策「Tweaks Panel：簡化原型版本」。**行為**：`Tweaks panel availability` 三個 Scenario 全綠；持久化 `sidebarCollapsed` 與 `showSummaryRail` 到 localStorage。**驗證**：`shell.test.tsx` 中模擬 `import.meta.env.DEV` true/false、`localStorage` 可用/不可用四種組合。
- [ ] 5.5 建立 `shell/WorkbenchV2Shell.tsx`，組合 TopBar + StepRail + 中央 step 內容 + SummaryRail + TweaksPanel；中央內容含 `STEP 0X / 04` 標頭與進度橫條（active=blue/done=emerald/pending=slate-700）。**行為**：`Default render of all regions` Scenario 對齊；切換 view 時中央內容換 placeholder。**驗證**：`shell.test.tsx` 主測試 case `renders all regions with default state`。

## 6. Keyboard shortcut（Cmd+B / Ctrl+B）

- [ ] 6.1 在 `WorkbenchV2Shell` 內加 keydown listener，按 `Cmd+B`（macOS）或 `Ctrl+B`（其他）切換 `sidebarCollapsed`，並對 keydown event 呼叫 `preventDefault`。**行為**：**Keyboard shortcut for step rail collapse** 的 `Shortcut toggles step rail` Scenario。**驗證**：`shell.test.tsx` fire `metaKey + b` keydown 斷言 sidebar 狀態翻轉。
- [ ] 6.2 對 keydown listener 加 `target.tagName in {INPUT, TEXTAREA, SELECT}` 或 `target.isContentEditable === true` 的早退判斷。**行為**：**Keyboard shortcut for step rail collapse** 的 `Shortcut suppressed inside text input` Scenario。**驗證**：`shell.test.tsx` 在 `<input>` 內 fire keydown 斷言 sidebar 狀態不變。

## 7. Placeholder step 與 settings

- [ ] 7.1 建立 `steps/Step1DevicePlaceholder.tsx`、`steps/Step2RulePlaceholder.tsx`、`steps/Step3MappingPlaceholder.tsx`、`steps/Step4DatabasePlaceholder.tsx`，每個都顯示「Step N · {title}」標題與「完整內容由後續 change 交付」訊息卡。**行為**：**Placeholder step and settings surfaces** 的 `Step 1 placeholder renders identifying content` Scenario；不發任何 network request（`fetch` / `axios` 均不被呼叫）。**驗證**：`shell.test.tsx` 對每個 step 切換後斷言文字與「未發 network」（用 jest.spyOn(window, 'fetch') 斷言 not called）。
- [ ] 7.2 建立 `settings/SettingsPlaceholder.tsx`，顯示「設定」標題與佔位訊息。**行為**：**Placeholder step and settings surfaces** 的 `Settings placeholder renders identifying content` Scenario。**驗證**：`shell.test.tsx` 切換到 settings 後斷言文字。

## 8. i18n namespace `workbench-v2`

- [ ] 8.1 建立 `frontend/src/i18n/locales/zh-TW/workbench-v2.json`，含 step 標題/副標、breadcrumb 片段、placeholder 訊息、按鈕、摘要欄 label，遵守原型繁中規範（`…`、`儲存草稿`、`取消`、`下一步` 等）。**行為**：`zh-TW locale shows traditional Chinese microcopy` Scenario。**驗證**：`shell.test.tsx` 在 zh-TW locale 斷言 `新增裝置` 出現。
- [ ] 8.2 建立 `frontend/src/i18n/locales/en/workbench-v2.json` 對應 key，提供英文翻譯（`Add Device`、`Source Rule`、`Point Mapping`、`Save to Database`、`Settings`、`Save Draft` 等）。**行為**：`en locale shows English microcopy` Scenario。**驗證**：`shell.test.tsx` 在 en locale 斷言 `Add Device` 出現。
- [ ] 8.3 在 i18n config 註冊新 namespace `workbench-v2`（檢視 `frontend/src/i18n/config.ts` 既有 namespace 列表並擴充），落地設計決策「i18n：新增 `workbench-v2` namespace」並實作需求 **i18n namespace for v2**。**行為**：`useTranslation('workbench-v2')` 不噴 missing namespace warning。**驗證**：`shell.test.tsx` console.error spy 斷言無相關 warning。

## 9. 路由註冊（Coexisting v2 workbench route + Coexisting v2 workbench at `/studio/v2`）

- [ ] 9.1 建立 `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`，內部用 `useWorkbenchV2State()` + `<WorkbenchV2Shell />`。**行為**：page 元件可獨立 mount。**驗證**：`shell.test.tsx` `render(<DatalinkWorkbenchV2Page />)` 不噴 error。
- [ ] 9.2 在 `frontend/src/App.tsx` 新增 `<Route path="/studio/v2" element={<DatalinkWorkbenchV2Page />} />`，置於 `/studio` route 之前。確認 `/studio`、`/datalink/*` legacy route 完全不動，落地設計決策「路由策略：/studio/v2 並存，/studio 保留為 fallback」並實作 datalink-workbench-desktop 的新需求 **Coexisting v2 workbench at `/studio/v2`** 與 datalink-workbench-v2-shell 的 **Coexisting v2 workbench route**。**行為**：兩條 spec 的所有 Scenario 全綠。**驗證**：`frontend/tests/workbench-v2/routing.test.tsx` 用 `MemoryRouter` 分別測試 `/studio`（loads `DatalinkWorkbenchPage`）、`/studio/v2`（loads `DatalinkWorkbenchV2Page`）、`/datalink/workbench`（redirect to `/studio`）。

## 10. 驗證

- [ ] 10.1 跑 `cd frontend && npm run lint`，零錯誤。
- [ ] 10.2 跑 `cd frontend && npm run test -- --run workbench-v2`，所有 `workbench-v2` 測試綠燈。
- [ ] 10.3 跑 `cd frontend && npm run build`，產出乾淨。
- [ ] 10.4 跑 `make check-lines`（或 `bash scripts/check_file_lines.sh`），確認所有 v2 新檔均 ≤ 300 行（硬上限 500 行）。
- [ ] 10.5 手動 smoke：開 `/studio` 確認 legacy 未受影響；開 `/studio/v2` 確認 shell 視覺對齊原型（dark canvas、漸層、24px 格線、Inter + JetBrains Mono、StepRail 行為、SummaryRail 行為、⌘B、Tweaks panel dev gate）。
