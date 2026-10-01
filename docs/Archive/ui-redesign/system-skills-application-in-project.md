# 系統技能在 go_gateway（C 方案）專案落地方式

更新時間：2026-03-02  
目標：把「frontend-patterns / e2e-testing / security-review / api-design / verification-loop」轉成可直接執行的專案作法。  
限制：**不修改後端 API 契約**。

## 1) frontend-patterns（前端模式）

**套用重點**
- Container/Presentational 分離：沿用 `SmartDashboardPage.tsx`（編排）+ `smart-dashboard/*`（展示/子邏輯）模式。
- 可重用邏輯以 hook 抽離：比照 `useSmartDashboardCommitFlow.ts`、`useSmartDashboardGridOverlays.ts`。
- C 版 Gateway 頁面拆為：路由守門、流程狀態、提交面板三層。

**對應檔案**
- `frontend/src/pages/datalink/SmartDashboardPage.tsx`
- `frontend/src/pages/datalink/smart-dashboard/useSmartDashboardCommitFlow.ts`
- `frontend/src/router/gateway.tsx`
- `frontend/src/pages/gateway/GatewayQuickSetupPage.tsx`
- `frontend/src/pages/gateway/GatewayExpertWorkbenchPage.tsx`

**可執行命令**
```bash
cd frontend
npm run lint
```

---

## 2) e2e-testing（端到端測試）

**套用重點**
- 以入口導流、Quick 流程、Expert 流程為最小三條主線。
- 先 mock `ENABLE_GATEWAY_DUAL_ENTRY` 再走流程，避免環境不穩。
- 補 header 與錯誤路徑斷言（例如 invalid manifest disable submit）。

**對應檔案**
- `frontend/tests/e2e/gateway/dual-entry.spec.ts`
- `frontend/playwright.config.ts`
- `frontend/package.json`（`test:gateway:e2e`, `test:gateway:gate`）

**可執行命令**
```bash
cd frontend
npm run test:gateway:e2e
# 或只跑單檔
npx playwright test tests/e2e/gateway/dual-entry.spec.ts
```

---

## 3) security-review（安全檢視）

**套用重點**
- 檢查是否新增危險渲染（`dangerouslySetInnerHTML`）、未受控輸入、未預期 header 覆蓋。
- C 版僅允許使用既有 API client interceptor 注入 `X-UI-Version`，不可自行拼接敏感 header。
- Quick/Expert 的原始輸入（raw manifest）必須先 JSON parse/驗證再送出。

**對應檔案**
- `frontend/src/pages/gateway/GatewayQuickSetupPage.tsx`
- `frontend/src/pages/gateway/GatewayExpertWorkbenchPage.tsx`
- `frontend/src/features/gateway/gatewayAdapter.ts`
- `frontend/src/services/api.ts`
- `frontend/src/services/datalink.ts`

**可執行命令**
```bash
cd frontend
rg "dangerouslySetInnerHTML|eval\\(|new Function\\(" src/pages/gateway src/features/gateway src/services
npm run test:gateway:unit
```

---

## 4) api-design（API 契約設計）

**套用重點**
- 前端只做 Adapter，不改後端 contract。
- Quick/Expert 一律輸出 `ConnectRequest` 形狀（`protocol` + `config`）。
- 路由切換、旗標、埋點都在前端層處理，不動 `/api/v1/*` 端點。

**對應檔案**
- `frontend/src/features/gateway/gatewayAdapter.ts`
- `frontend/src/features/gateway/dualEntryFlag.ts`
- `frontend/src/features/gateway/uiVersion.ts`
- `frontend/src/router/gateway.tsx`

**可執行命令**
```bash
cd frontend
npm run test:gateway:unit
```

---

## 5) verification-loop（驗證閉環）

**套用重點**
- 每個小任務都走固定迴圈：`實作 -> unit -> e2e -> lint -> build`。
- 以 `test:gateway:gate` 當 release 前守門。
- 若失敗，先回到 Adapter/流程 state 修正，不做後端繞路。

**對應檔案**
- `frontend/package.json`（`test:gateway:gate`）
- `frontend/src/features/gateway/__tests__/*`
- `frontend/src/pages/gateway/__tests__/*`

**可執行命令**
```bash
cd frontend
npm run test:gateway:unit
npm run test:gateway:e2e
npm run lint
npm run build
npm run test:gateway:gate
```

---

## 6) 建議實作順序（技能對應）
1. 先用 **frontend-patterns** 拆層（避免後續難測）。  
2. 再用 **api-design** 固定 adapter 邊界（確保不破壞 API）。  
3. 接著補 **e2e-testing** 主線腳本。  
4. 每輪都跑 **verification-loop**。  
5. 最後做 **security-review** 全掃描與修補。

