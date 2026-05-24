## Why

目前系統已經有 `/studio`、`/gateway/quick-setup`、`/gateway/expert-workbench` 等可用功能，但首次進入產品時缺少一個清楚的前置入口，讓操作員能快速理解「我要快速完成設定」與「我要進入既有完整功能」之間的差別。這導致新使用者需要自己理解多個路徑，而不是先從一個入口頁做選擇。

本變更的目標是新增一個入口頁，讓操作員能快速進入完整設定流程，同時明確保留現有功能、現有深連結與現有工作台，不以任何收斂或替換為前提。

## What Changes

- 新增或正式化一個產品入口頁，作為快速完整設定的前置入口。
- 入口頁提供清楚的 `快速完整設定` 主行動，導向現有 quick setup 流程完成設定，不重寫其既有能力。
- 入口頁同時保留前往 `/studio` 與 expert workbench 的明確入口，讓現有功能全部可達。
- 現有 `/studio`、`/gateway/quick-setup`、`/gateway/expert-workbench`、`/test` 路由與能力維持有效，不因新增入口頁而退役或降級為相容模式。
- 補齊入口頁、路由、文案與測試，確保「新增入口」不造成既有功能遺失。

## Non-Goals

- 不移除、合併或改寫既有 `/studio`、quick setup、expert workbench 的核心功能。
- 不把所有產品流程重構成單一路由或單一 shell。
- 不修改 Step 2 到 Step 4 的核心業務邏輯、protocol adapter 或 backend contract。
- 不將 `/test` 工程工具併入產品入口頁。

## Capabilities

### New Capabilities

- `gateway-setup-entry`: 提供一個前置入口頁，讓操作員快速選擇完整設定流程或進入既有進階功能，且不隱藏任何既有能力。

### Modified Capabilities

- `datalink-ui`: 主產品流程新增入口頁起點，但既有 `/studio` 與其他功能路徑仍保持可直接進入。
- `datalink-workbench-desktop`: `/studio` 作為主要工作台的行為不變，但系統需允許入口頁在工作台之前導向既有流程而不取代工作台。

## Impact

- Affected specs: `gateway-setup-entry`, `datalink-ui`, `datalink-workbench-desktop`
- Affected code:
  - Modified: `frontend/src/App.tsx`, `frontend/src/router/gateway.tsx`, `frontend/src/pages/gateway/GatewayEntryPage.tsx`, `frontend/src/features/gateway/dualEntryFlag.ts`, `frontend/src/features/gateway/useGatewayDualEntryFlag.ts`, `frontend/src/i18n/locales/en/common.json`, `frontend/src/i18n/locales/zh-TW/common.json`
  - Possibly adjusted for entry handoff copy or links: `frontend/src/pages/gateway/GatewayQuickSetupPage.tsx`, `frontend/src/pages/gateway/GatewayExpertWorkbenchPage.tsx`
  - Modified tests: `frontend/src/pages/gateway/__tests__/GatewayQuickSetupPage.test.tsx`, `frontend/src/pages/gateway/__tests__/GatewayExpertWorkbenchPage.test.tsx`, route-level tests for `frontend/src/App.tsx` or `frontend/src/router/gateway.tsx`
