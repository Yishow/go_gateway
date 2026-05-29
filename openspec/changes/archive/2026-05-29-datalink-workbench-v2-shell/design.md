## Context

`new_prototype/` 已交付完整 4 步驟工作台原型（Dark Industrial Telemetry 風格、8 個 jsx 檔、3,515 行），但原型是 CDN-Tailwind + Babel-in-browser 的展示型 React，無法直接搬進 production。同時，現行 `/studio` 由 `datalink-workbench-desktop` spec 守護，是 datalink 系統當前唯一的正式主入口；任何 shell 級重構直接覆蓋會引入回歸風險。

本 change 是 `datalink-workbench-v2` phase 1（共 6 個並列 change）的入口，先把 shell + design tokens + 共用元件 + 4 個 step 空骨架 + 空 settings 落地到 `/studio/v2`，使後續 5 個 change（step1/step2/step3/step4/settings）能各自獨立替換對應 placeholder。

## Goals / Non-Goals

**Goals:**

- 在不破壞 `/studio` 任何既有行為的前提下，提供 `/studio/v2` 作為新版漸進落地容器。
- 落地 Dark Industrial Telemetry 設計 tokens（顏色、字型、格線底紋、漸層光暈），所有後續 step change 都從同一份 tokens 取值。
- 移植原型中 `shared.jsx` 與 `tweaks-panel.jsx` 的共用元件到 TypeScript + 真實檔案結構，符合 repo 既有「200–400 行 / 500 行硬上限」規範。
- 建立可由後續 step change 各自替換的 step placeholder 與 settings placeholder。

**Non-Goals:**

- 不接後端 API、不使用 React Query；所有 state 以 in-memory 形式承載，後續 backend-wiring change 才處理。
- 不刪除 `/studio` 路由或 `DatalinkWorkbenchPage`；舊主線保留。
- 不引入新 UI library（shadcn/ui、Radix 等），維持「Tailwind utility + 自寫 component」一致性。
- 不在本 change 完成 light mode；tokens 僅交付 dark mode。

## Decisions

### 路由策略：/studio/v2 並存，/studio 保留為 fallback

採用「並存路由 + 漸進切換」模式：

- `App.tsx` 新增 `<Route path="/studio/v2" element={<DatalinkWorkbenchV2Page />} />`，與既有 `<Route path="/studio" element={<DatalinkWorkbenchPage />} />` 並存。
- 舊 `/studio` 與所有 `/datalink/*` legacy redirect 行為不變。
- 後續所有 v2 change 都對齊 `/studio/v2`；直到 phase 3 結束、e2e 全綠後，由獨立 change 決定是否把 `/studio` 指向 v2 或保留並存。

**Alternatives considered**：(A) 直接覆蓋 `/studio` — 風險過高，任何中間 phase 的 bug 都會影響主線；(B) 用 feature flag 切換同一路由 — 增加 flag 系統複雜度，且 e2e 測試難以同時驗證新舊兩版。

### Feature 目錄結構：`frontend/src/features/datalink/workbench-v2/`

採用 feature-first 目錄而非 `frontend/src/pages/datalink/workbench-v2/`：page-level component 仍放在 `pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`（路由綁定），但 shell、共用元件、tokens、step、settings、state 全部集中在 `features/datalink/workbench-v2/` 子目錄。

理由：repo 既有 `frontend/src/features/datalink/legacyRoutes.ts` 已是 feature-first 慣例；shell + step + settings 是高度耦合的整體，應對齊在一個 feature 目錄裡。

### 共用元件移植：1:1 對應原型 + TypeScript 化

- `Icon`、`StatusChip`、`SectionCard`、`Field`、`Input`、`Select`、`Textarea`、`Button`、`Toggle` 全部對應原型 `shared.jsx` 內的同名元件。
- 全部以 TypeScript 重寫，props 型別嚴格定義；避免 `any`。
- 圖示集（device / rule / map / db / check / chevron / plus / play / refresh / cable / sliders / tag / table / alert / info / spark / arrow / close / save / eye / bolt / flow / grid）以單一 `Icon name="..."` 元件承載，內部用 SVG `paths` 字典；禁用 emoji 作為結構性圖示。
- `Toggle` 沿用原型的精確像素位移實作，避免 ON 狀態 dot 視覺溢出。

**Alternatives considered**：用 Radix Primitives 或 shadcn/ui — 會引入新依賴，且原型風格特殊（工業儀表板）難以套用既有 component theme。

### Design Tokens 落地：CSS variables + Tailwind theme.extend

- `frontend/src/features/datalink/workbench-v2/tokens.ts` export TypeScript 常數：`COLORS`、`FONTS`、`RADII`、`SHADOWS`、`ANIMATIONS`。
- `frontend/src/features/datalink/workbench-v2/styles/workbench-v2.css` 定義 CSS variables (`--wbv2-bg`、`--wbv2-surface`、`--wbv2-grid-color`) 以及原型的 keyframes（`pulseDot`、`sweepIn`、`shimmer`）與 utility classes（`.bg-canvas`、`.bg-grid`、`.pulse-dot`、`.sweep-in`、`.shimmer`、`.mem-cell`、`.kbd`、`.chip`、`.field`、`.label`）。
- `DatalinkWorkbenchV2Page` 為 v2 內容掛上 `data-workbench-v2` data attribute，CSS 變數作用域限縮於該 subtree，不污染既有頁面。
- 不修改全域 `tailwind.config`；維持與既有 `/studio` 完全隔離。

**Alternatives considered**：把 tokens 寫入 `tailwind.config.theme.extend.colors.wbv2` — 全域 Tailwind 設定容易讓 token 外洩到既有頁面，且 v2 仍在實驗階段不宜污染共用設定。

### 字型：`@fontsource/inter` + `@fontsource/jetbrains-mono` 本地化

- 原型用 Google Fonts CDN（`Inter` + `JetBrains Mono`），生產環境不可走 CDN（離線部署 + CSP）。
- 在 `frontend/package.json` 新增 `@fontsource/inter` 與 `@fontsource/jetbrains-mono`；於 `main.tsx`（或 v2 專屬 entry）import `@fontsource/inter/400.css` `/500.css` `/600.css` `/700.css` 與 `@fontsource/jetbrains-mono/400.css` `/500.css` `/600.css`。
- 在 v2 subtree CSS 設定 `font-family: 'Inter', system-ui, sans-serif`，mono 區塊用 `'JetBrains Mono', ui-monospace, monospace`。

**Alternatives considered**：用 `next/font` 風格的 build-time 子集化 — 目前 repo 是 Vite，沒有 `next/font`；用 `vite-plugin-fonts` 會額外引入 build 步驟。

### State 管理：本地 React state + 可序列化形狀

- `state/types.ts` 定義 `WorkbenchV2State` interface，對應原型的 `state` shape（`devices` / `rules` / `points` / `mappings` / `db` / `settings` / `commit` / `committed`）。
- `state/useWorkbenchV2State.ts` 提供 `useWorkbenchV2State()` hook，封裝 `useReducer` 與 immutable update。
- 不使用 Zustand / Redux / Context；後續 backend-wiring change 才會引入 React Query，並把 state 切割為 server cache + UI state。

### i18n：新增 `workbench-v2` namespace

- `frontend/src/i18n/locales/zh-TW/workbench-v2.json` 與 `frontend/src/i18n/locales/en/workbench-v2.json`。
- 所有 microcopy 走 `useTranslation('workbench-v2')`，遵守原型 `_overview.md` 的繁中規範（省略號用 `…`、loading「載入中…」、按鈕「儲存」「取消」「建立」「下一步」「繼續」「提交並啟動排程器」、狀態 draft=草稿/active=啟用/disabled=停用/tested=已測試/committed=已部署）。

### Tweaks Panel：簡化原型版本

- 原型 `tweaks-panel.jsx` 有 host postMessage 通訊（用於 Omelette / Anthropic Edit Mode），在 production app 不需要。
- 本 change 提供精簡版：`TweaksPanel` 元件僅在 dev mode 或 `localStorage.WBV2_TWEAKS === '1'` 時顯示，控制版面切換（收合側欄、顯示摘要欄）與示範流程入口（重置、加入第二台設備、加入第二條規則、跳到 Step 4）。
- 不保留 `useTweaks` 的 postMessage 行為；用 `useState` + `localStorage` 持久化。

## Implementation Contract

#### Behavior

- 開啟 `/studio/v2` 載入 `DatalinkWorkbenchV2Page`：顯示頂列（含 collapse 按鈕、品牌標識 "Datalink Workbench"、breadcrumb `Datalink › Workbench › Step 1: 新增裝置`、儲存草稿/取消、scheduler 狀態 `scheduler idle`）+ 左側 StepRail（4 步驟可見、`新增裝置` highlighted、`設定` 入口在分隔線下）+ 中央 Step 內容（顯示 STEP 01 / 04 + 進度橫條 + Step1 placeholder「Step 1 內容即將上線」訊息）+ 右側 SummaryRail（≥1280px 顯示，呈現 0 個設備 / 0 條規則 / 0 Tags / 尚未設定 DB）。
- 按 ⌘B（macOS）或 Ctrl+B：左側 StepRail 在 232px ↔ 64px 間切換；收合時僅顯示 step icon，hover 浮現 tooltip（step 名稱）。
- 點擊左側 StepRail 任一已 reachable 的 step：切換 `current` state，中央內容區換成對應 placeholder，breadcrumb 同步更新。
- 點擊 StepRail 內「設定」：`view` state 從 `'flow'` 切到 `'settings'`，中央內容區換成 Settings placeholder，右側 SummaryRail 隱藏。
- 開啟 Tweaks panel（dev mode 或 localStorage flag）：可切換「收合側邊欄」「顯示右側摘要欄」，狀態持久化到 localStorage。
- 開啟 `/studio`：行為與本 change 前完全一致；不發生 redirect。

#### Interface / Data Shape

- 路由：`/studio/v2` 對應 `DatalinkWorkbenchV2Page`，無 query/path params。
- `WorkbenchV2State`（TypeScript interface）：
  ```ts
  type WorkbenchV2State = {
    view: 'flow' | 'settings';
    current: 1 | 2 | 3 | 4;
    completed: Set<number>;
    sidebarCollapsed: boolean;
    showSummaryRail: boolean;
    devices: Device[];
    rules: Rule[];
    points: Point[];
    mappings: Record<string, Mapping>;
    db: { connector?: DbConnector; targets?: Record<string, DbTarget> };
    settings: Settings;
    commit?: CommitState;
    committed: boolean;
  };
  ```
  完整型別於 `state/types.ts` 對齊 `new_prototype/docs/_overview.md` 的 State shape。本 change 僅實作 view / current / completed / sidebarCollapsed / showSummaryRail 對應的行為；其餘欄位保留 default value 供後續 change 使用。
- 共用元件 props 介面：每個元件導出 named TypeScript interface（`ButtonProps`、`FieldProps` …），不使用 `any`。

#### Failure Modes

- 載入 `@fontsource` 失敗時，font-family 退回 `system-ui, sans-serif`（fallback chain）；不阻塞 UI。
- `localStorage` 不可用（隱私模式）：Tweaks 狀態回退到記憶體 state，重整後重置；console.warn 一次。
- 進入 `/studio/v2` 時若使用者瀏覽器 viewport < 1024px：SummaryRail 不顯示，StepRail 強制收合，仍可使用主流程；不顯示行動裝置警告（v2 暫不支援 mobile）。

#### Acceptance Criteria

- `npm run lint` 通過。
- `npm run test` 中 `frontend/tests/workbench-v2/shell.test.tsx` 全綠，覆蓋：
  - 進入 `/studio/v2` 後 shell 結構正確（頂列、StepRail、SummaryRail、Step1 placeholder 都 in DOM）。
  - 點擊 StepRail step 2：current 切到 2，中央顯示 Step2 placeholder，breadcrumb 顯示 Step 2 標題。
  - 按 ⌘B（mock `metaKey + b` keydown）：`sidebarCollapsed` 切換，aria-expanded / data-collapsed 屬性同步。
  - 點擊「設定」入口：view 切到 settings，中央顯示 Settings placeholder。
- `npm run build` 通過。
- 手動驗證：開啟 `/studio` 確認舊版完全不受影響；開啟 `/studio/v2` 確認 shell 視覺對齊原型（dark canvas、漸層光暈、24px 格線、Inter + JetBrains Mono、StepRail highlight 行為）。

#### Scope Boundaries

**In scope:**

- `/studio/v2` 路由註冊與 page 元件。
- Shell 五大區塊（TopBar / StepRail / 中央內容 / SummaryRail / TweaksPanel）。
- 設計 tokens、CSS variables、keyframes、utility classes。
- 共用元件 9 個（Icon / Button / Field / Input / Select / Textarea / StatusChip / SectionCard / Toggle）。
- 4 個 step 與 settings 的空骨架（placeholder）。
- i18n namespace `workbench-v2`（zh-TW + en）。
- `@fontsource/inter` 與 `@fontsource/jetbrains-mono` 引入。
- 上述項目的 unit test。

**Out of scope:**

- Step1/Step2/Step3/Step4/Settings 的實際內容（交給後續 5 個 change）。
- 任何後端 API 呼叫、React Query hooks（交給 phase 2 的 backend-wiring change）。
- E2E 測試（交給 phase 3 的 e2e change）。
- Light mode tokens。
- Mobile responsive（< 1024px viewport）。
- 把 `/studio` 切換到 v2（屬於未來獨立 change）。

## Risks / Trade-offs

- [字型 self-host 增加 bundle 體積] → 用 `@fontsource` 只 import 用到的字重（400/500/600/700），並讓 `font-display: swap`，影響可接受。
- [兩條主路由共存增加 UX 認知負擔] → `/studio/v2` 是工程內部用詞，不直接從 `/studio` 連結；對外暫不宣傳，僅由工程團隊使用。
- [tokens 與既有 `/studio` 視覺分歧] → 透過 `data-workbench-v2` subtree 隔離 CSS variables，避免互相污染；review checklist 包含「在 `/studio` 確認無視覺改動」。
- [Tweaks panel 暴露在生產] → 用 `import.meta.env.DEV` 與 localStorage flag 雙重 gate，預設關閉。
- [⌘B 衝突 macOS 既有快捷鍵] → 多數瀏覽器 ⌘B 並非預設行為；同時提供 collapse 按鈕作為替代入口，且 keydown listener 用 `preventDefault()` 僅在 v2 page 內生效。

## Migration Plan

1. 安裝 `@fontsource/inter` 與 `@fontsource/jetbrains-mono`（`npm install`）。
2. 建立 `frontend/src/features/datalink/workbench-v2/` 目錄與所有檔案；先寫 tokens / styles / 共用元件，再寫 shell，再寫 placeholder。
3. 註冊 `/studio/v2` 路由於 `App.tsx`，確認舊 `/studio` 路由與 redirect 規則完全沒動。
4. 寫 `frontend/tests/workbench-v2/shell.test.tsx`，TDD 順序：先測 shell 結構 → 測 step 切換 → 測 ⌘B → 測 settings 切換。
5. `npm run lint`、`npm run test`、`npm run build` 全部通過。
6. 手動 smoke test：`/studio` 與 `/studio/v2` 並存且都能載入。
7. PR 合併後，不變更任何使用者文件；後續 phase 2 完成且 e2e 通過後，再決定 `/studio` 切換策略。

**Rollback strategy**：本 change 不刪除任何既有檔案；rollback 只需移除 `/studio/v2` route 註冊與 v2 feature 目錄；package.json 兩個 `@fontsource` 依賴可選擇保留或移除（保留無害）。

## Open Questions

- v2 是否要在 dev mode 自動 redirect `/studio` → `/studio/v2` 以方便工程內部 dogfooding？目前傾向 **不**，避免污染 staging/preview 環境；最終方案在 backend-wiring change 階段再決定。
- Tweaks panel 的「重置流程」「加入第二台設備」「加入第二條規則」「跳到 Step 4」這四個示範按鈕，是否仍需在本 change 提供？傾向 **是**（純 client-side state mutation，不依賴後端），方便 step change 開發時快速產生資料。本 change 落地骨架，按鈕 onClick handler 留 TODO 註解，由後續 step change 補完。
