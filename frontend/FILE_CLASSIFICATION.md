# Frontend 檔案分類

更新時間：2026-03-20

## 1) 路由與頁面分類

### 核心產品頁（現行）
- `src/pages/datalink/workbench/DatalinkWorkbenchPage.tsx`（對外主入口對應 `/studio`）

### 相容/過渡路由（Legacy Redirect）
- `/datalink`：redirect 到 `/studio`
- `/datalink/workbench`：redirect 到 `/studio`
- `/datalink/dashboard-legacy`：redirect 到 `/studio`
- `/datalink/devices-legacy`：redirect 到 legacy intent 對應的新主線

### 協議測試工具頁（獨立於 Datalink 主流程）
- `src/pages/TestPage.tsx`
- `src/pages/TestPageShell.tsx`

### 頁面測試
- `src/pages/datalink/workbench/__tests__/DatalinkWorkbenchFoundation.test.tsx`
- `src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx`
- `src/pages/datalink/workbench/__tests__/DatalinkWorkbenchTagStep.test.tsx`
- `src/pages/datalink/workbench/__tests__/DatalinkWorkbenchOutputStep.test.tsx`

## 2) UI 與元件分類

### Datalink 領域元件
- 目錄：`src/components/datalink/`
- 用途：MemoryGrid、QuickActions、DeviceForm、Import/Export、Wizard 等核心互動元件。

### 協議測試工具元件
- 檔案：`src/components/ProtocolSelector.tsx`、`src/components/ConfigForm.tsx`、`src/components/TestOperations.tsx`、`src/components/DeviceScanner.tsx`、`src/components/MonitorControl.tsx`、`src/components/RTUPollingCard.tsx` 等。

### 共用基礎元件
- 檔案：`src/components/Toast.tsx`、`src/components/ThemeToggle.tsx`、`src/components/ui/*`

## 3) 狀態與領域邏輯分類

### Feature 層（純邏輯/策略）
- 目錄：`src/features/datalink/`、`src/features/flow/`
- 用途：commit lifecycle、route intent、allocation strategy、motion guidance、state machine 等。

### Hooks 層
- 目錄：`src/hooks/datalink/`：API Query/Mutation hooks
- 目錄：`src/hooks/`：快捷鍵、歷史、profiles、preview stream

### Service 層
- 目錄：`src/services/`
- 用途：HTTP API 呼叫封裝（`api.ts`、`datalink.ts`）

### 型別與工具
- 型別：`src/types/`
- 工具：`src/utils/`
- 共用函式：`src/lib/`

## 4) 全域基礎設施分類

- 入口與路由：`src/main.tsx`、`src/App.tsx`
- i18n：`src/i18n/`
- Context：`src/contexts/`
- 樣式：`src/index.css`、`src/styles/tokens.ts`
- 測試初始化：`src/setupTests.ts`
- 建置設定：`vite.config.ts`、`tsconfig.json`、`tailwind.config.cjs`、`eslint.config.js`

## 5) 測試檔案分類

- 頁面互動測試：`src/pages/**/__tests__`
- 元件測試：`src/components/**/__tests__`
- 邏輯測試：`src/features/**/__tests__`
- Hook 測試：`src/hooks/**/__tests__`
- Utils 測試：`src/utils/**/__tests__` 與 `*.test.ts`

## 6) 本次去重已刪除檔案

以下檔案已確認無路由與無引用，已移除：

- `src/pages/datalink/DeviceOnboardingPage.tsx`
- `src/pages/datalink/PointsPage.tsx`
- `src/pages/datalink/PollingGroupsPage.tsx`
- `src/pages/datalink/TagsPage.tsx`
- `src/pages/datalink/MappingsPage.tsx`
- `src/pages/datalink/MappingWizardPage.tsx`
- `src/pages/datalink/SettingsPage.tsx`
- `src/pages/datalink/Dashboard.tsx`
- `src/pages/datalink/DevicesPage.tsx`
- `src/pages/datalink/SmartDashboard.tsx`
- `src/pages/datalink/SmartDashboardPage.tsx`
- `src/pages/datalink/LocalModbusWorkbenchPage.tsx`
- `src/pages/TemplatesPage.tsx`
- `src/pages/HistoryPage.tsx`
- `src/pages/ComparePage.tsx`
- `src/pages/AnalyzerPage.tsx`
- `src/components/Layout.tsx`
