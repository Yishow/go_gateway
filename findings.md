# Findings

## 2026-03-08

### 規範來源盤查
- `AGENTS.md` 已作為專案共通規範入口，並已納入 UI/UX 主流程目標、TDD 與 `planning-with-files`。
- `CLAUDE.md` 已更新為：先讀 `AGENTS.md`，再依檔案類型讀 `.github/instructions/`，最後讀 `CLAUDE.md`。
- `GEMINI.md` 已改寫成與目前 repo 現況一致的 Gemini Agent 補充文件，並納入同樣的規範分工。
- `.github/instructions/` 目前共有三份正式規範：
  - `go.instructions.md`
  - `reactjs.instructions.md`
  - `typescript-5-es2022.instructions.md`

### 建議的規範分工
- `AGENTS.md` 管專案怎麼做。
- `CLAUDE.md` / `GEMINI.md` 管 Agent 怎麼工作。
- `.github/instructions/` 管程式怎麼寫。

### 已完成的文件對齊
- `AGENTS.md`
  新增「規範來源分層」，明確宣告 `.github/instructions/` 為正式實作規範來源，並定義優先順序。
- `CLAUDE.md`
  補上文件分工、閱讀順序與規範優先順序。
- `GEMINI.md`
  從舊版原型描述整理為目前專案可直接採用的 Agent 補充文件。

### 與目前 UI/UX 計畫的關聯
- `SmartDashboard` 應收斂成單人主工作流。
- `LocalModbusWorkbenchPage` 應作為主流程後段工作台。
- `TestPage` 保持工程工具定位，只做風格一致化。
- 後續前端實作需要同時遵守 React 與 TypeScript instructions。

### TDD 基線盤查
- `SmartDashboard` 已有多組互動與 hook 測試，包含：
  - `frontend/src/pages/datalink/__tests__/SmartDashboard.interaction.test.tsx`
  - `frontend/src/pages/datalink/__tests__/SmartDashboardGridOverlaysSection.test.tsx`
  - `frontend/src/pages/datalink/__tests__/useSmartDashboardCommitFlow.test.ts`
  - `frontend/src/pages/datalink/__tests__/useSmartDashboardWorkspaceContentState.test.tsx`
  - `frontend/src/pages/datalink/__tests__/useSmartDashboardWorkspaceState.test.ts`
  - `frontend/src/pages/datalink/__tests__/useSmartDashboardPanelsState.test.ts`
- `LocalModbusWorkbenchPage` 已有基礎頁面測試：
  - `frontend/src/pages/datalink/__tests__/LocalModbusWorkbenchPage.test.tsx`
- 目前未找到 `TestPage` 對應測試檔，代表後續若要做樣式一致化，必須先建立 `TestPage` 測試基線。

### 已補上的 TestPage 測試基線
- 已新增 `frontend/src/pages/__tests__/TestPage.test.tsx`
- 目前覆蓋的穩定頁面行為：
  - 初始化時會批量最小化非 config 卡片
  - 連線成功後會自動收折 config 卡片，且可手動展開
  - 協議切換時至少能穩定驗證模式同步到 `serial`
- 測試刻意聚焦頁面殼層與使用者可觀察行為，避免綁定過多內部實作細節，方便後續做 UI 風格一致化。

### TestPage 第一輪 UI 收斂
- `frontend/src/pages/TestPage.tsx` 已先完成低風險的外殼收斂：
  - 頁面背景改為與 datalink 主介面相近的深色漸層
  - 頂部加入工程工作台定位說明
  - 主要卡片表面統一為深色面板語言
  - 將本次碰到的 `transition-all` 改為更明確的 transition 屬性
- 此輪不改變 `TestPage` 的工具定位與核心流程，只收斂視覺語言與資訊層級。

### 前端測試入口重整
- 目前前端測試原本散落在：
  - `frontend/src/**/__tests__/*`
  - `frontend/src/**/*.test.ts(x)`
  - `frontend/tests/e2e/*`
- 已建立新的正式測試入口與分類：
  - `frontend/tests/unit/`
  - `frontend/tests/integration/`
  - `frontend/tests/e2e/`
- 本輪採用「root 測試入口收斂 + 舊測試內容逐步遷移」策略：
  - 先讓 Vitest 只收斂執行 `frontend/tests/unit` 與 `frontend/tests/integration`
  - 既有 source-adjacent 測試先由 root wrapper 匯入
  - `TestPage` 測試則已實際移到 `frontend/tests/unit/pages/test-page.test.tsx`
- `Playwright` 原本就使用 `frontend/tests/e2e/`，此次維持不變。
- `Go` 測試不適合集中到 root `tests/`，需維持 `*_test.go` 與實作檔相鄰，這一點已同步寫進 `AGENTS.md`。

### Root 測試入口的實務發現
- 若把多個有 `vi.mock` 隔離需求的 page test 直接聚合到同一個 wrapper，會因模組先後載入順序讓 mock 失效。
- 因此 root `tests/` 的分類不應只看資料夾層級，還要考慮測試隔離邊界：
  - `Gateway` page 測試拆成 `expert-workbench.test.ts` 與 `quick-setup.test.ts`
  - `SmartDashboard` page 測試拆成 interaction / grid-overlays / state 三組 wrapper
- 這個粒度可同時滿足：
  - 測試入口集中在 root `tests/`
  - 類型與領域有清楚分類
  - `vi.mock` 與模組快取不互相污染
