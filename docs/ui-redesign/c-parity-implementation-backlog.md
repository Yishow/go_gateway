# C 方案功能對齊實作 Backlog（Phase 0）

更新時間：2026-03-02  
原則：僅前端重構與整合，**不修改後端 API 契約**。  
估時單位：工程日（1d = 8h）

## P0（先達到可提交閉環）

| ID | 任務 | 主要檔案 | 估時 | 測試方式 | 完成定義 |
|---|---|---|---:|---|---|
| P0-01 | Quick Setup 五步 state machine（設備/點位/Tag/驗證/Commit） | `frontend/src/pages/gateway/GatewayQuickSetupPage.tsx`、`frontend/src/features/gateway/*` | 2.5d | unit + e2e + manual | 可從 Step1 走到 Step5，失敗可回上一步且草稿不丟。 |
| P0-02 | Quick/Expert 共用 draft store，支援入口切換無損 | `frontend/src/features/gateway/*store*`（新增）、`GatewayQuickSetupPage.tsx`、`GatewayExpertWorkbenchPage.tsx` | 1.5d | unit + e2e | 入口切換後欄位保持一致；Expert→Quick 降級有警告。 |
| P0-03 | 擴充 `gatewayAdapter` 對齊 datalink 常見欄位與保留策略 | `frontend/src/features/gateway/gatewayAdapter.ts`、`__tests__/gatewayAdapter.test.ts` | 1.5d | unit | 複雜 payload round-trip 後不遺失未知欄位。 |
| P0-04 | 接入流程驗證 gate（結構驗證 + pipeline 驗證） | `frontend/src/pages/gateway/*`、復用 `frontend/src/features/datalink/validationFlow.ts` | 1.0d | unit + manual | Commit 前必跑驗證，錯誤訊息一致可重試。 |
| P0-05 | 接入 commit lifecycle（chunk/retry/rollback）最小版 | `frontend/src/pages/gateway/*`、復用 `useSmartDashboardCommitFlow.ts` | 2.0d | unit + e2e | 能重現 partial fail、retry 成功、rollback。 |
| P0-06 | 補齊 gateway e2e 覆蓋（5-step + switch + commit） | `frontend/tests/e2e/gateway/dual-entry.spec.ts` | 1.0d | e2e | 新增 3 條核心腳本並穩定通過。 |
| P0-07 | i18n 化 Gateway 三頁（移除硬編碼） | `frontend/src/pages/gateway/*.tsx`、`frontend/src/i18n/locales/*` | 1.0d | unit + manual | zh-TW/en 切換正常，無硬編碼文案。 |

## P1（補齊作業深度）

| ID | 任務 | 主要檔案 | 估時 | 測試方式 | 完成定義 |
|---|---|---|---:|---|---|
| P1-01 | Quick Step1 補設備搜尋/篩選/連線測試 | `GatewayQuickSetupPage.tsx`、`hooks/datalink/useDevices.ts` | 1.5d | unit + manual | 可篩選 active/disabled/draft 並測連線。 |
| P1-02 | Expert 佈局接入 `MemoryGrid` 與衝突視圖 | `GatewayExpertWorkbenchPage.tsx`、`components/datalink/MemoryGrid.tsx` | 2.0d | integration + e2e | 衝突格可視化、可選取、可回寫草稿。 |
| P1-03 | 接入 Tag 連結/新建/guardrail | `GatewayQuickSetupPage.tsx`、`GatewayExpertWorkbenchPage.tsx`、`useSmartDashboardTagLinking.ts` | 2.0d | unit + manual | 多映射 Tag 編輯會觸發 guardrail。 |
| P1-04 | 接入批次建立點位與命名預覽 | `SmartDashboardPanels.tsx`（復用）、`GatewayExpertWorkbenchPage.tsx` | 1.5d | integration | 批次建立後可自動關聯到 Tag。 |
| P1-05 | 接入 Import/Export + Undo/Redo 工具列 | `GatewayExpertWorkbenchPage.tsx`、`usePointHistory.ts` | 1.0d | unit + manual | 匯入可撤銷、重做可回復。 |
| P1-06 | Local Modbus 快捷操作與 deep-link 工作台 | `GatewayExpertWorkbenchPage.tsx`、`SmartDashboardModbusPanel.tsx` | 1.0d | manual | 3 個快捷操作可用，且可跳轉完整工作台。 |

## P2（品質與上線收斂）

| ID | 任務 | 主要檔案 | 估時 | 測試方式 | 完成定義 |
|---|---|---|---:|---|---|
| P2-01 | Commit audit trace/UI 觀測資訊補齊 | `GatewayExpertWorkbenchPage.tsx`、`services/api.ts` | 0.8d | unit + manual | UI 可看 traceId，請求 header 完整。 |
| P2-02 | 效能優化（大資料量 route/grid） | `GatewayExpertWorkbenchPage.tsx` | 1.0d | manual | 100+ 條 route 編輯不卡頓（互動延遲可接受）。 |
| P2-03 | 可及性與快捷鍵一致化 | `frontend/src/pages/gateway/*.tsx` | 0.8d | manual | Tab/Enter/快捷鍵可完整操作核心流程。 |
| P2-04 | 文件與守門命令整併（CI gate） | `docs/ui-redesign/*`、`frontend/package.json`（若需） | 0.5d | manual | 文件與 `test:gateway:gate` 對齊。 |

## 建議執行節奏（可直接跑）

```bash
cd frontend

# 1) 先跑 gateway 基線守門
npm run test:gateway:gate

# 2) 日常開發迴圈（每完成一個 backlog item）
npm run test:gateway:unit
npm run test:gateway:e2e
npm run lint
npm run build
```

## 依賴關係（先後順序）
1. **P0-01 ~ P0-03** 完成後，才進行 P0-04/P0-05。  
2. **P0 全部完成** 才能進 P1（否則會在錯誤基底上擴充）。  
3. P2 為收斂階段，建議在 P1 穩定後一次處理。

