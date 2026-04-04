# UI/UX 改善執行計畫

**開始時間**：2026-03-22  
**目標**：依 findings.md 改善專案 UI/UX 缺失  
**策略**：優先處理 P0，再逐步處理 P1、P2

---

## 臨時任務：2026-03-23 規範文件對齊 ✅

- [x] 深度盤點 `README.md`、`Makefile`、`frontend/package.json`、`scripts/build.ps1`、`frontend/tsconfig.json`、`frontend/eslint.config.js`、`.golangci.yml`
- [x] 對齊 `AGENTS.md` 的專案總覽、建置與測試指令、程式碼樣式、測試規範、安全考量與 repo 特定規則
- [x] 對齊 `CLAUDE.md` 的專案總覽、常用命令、樣式 / 測試 / 安全摘要與工作流程
- [x] 修正兩份文件對規範優先順序的描述，使其一致為：`AGENTS.md -> Agent 專屬文件 -> .github/instructions/`
- [x] 以 `git diff --check -- AGENTS.md CLAUDE.md` 驗證文件 patch 無 whitespace / 格式錯誤

**備註**：
- 本次僅更新文件與 planning records，未碰其他進行中的程式碼變更。
- 由於屬文件對齊任務，未額外執行 build / test；驗證以來源盤點與 diff 檢查為主。

---

## Phase 1: P0 高優先級改善（影響可用性）

### 1.1 統一載入狀態與錯誤處理 ✅
- [x] 建立統一 Spinner 元件（`frontend/src/components/ui/spinner.tsx`）
- [x] 建立統一 Skeleton 元件（`frontend/src/components/ui/skeleton.tsx`）
- [x] 建立 ErrorBoundary 元件（`frontend/src/components/ErrorBoundary.tsx`）
- [x] 更新 App.tsx 加入 ErrorBoundary
- [x] 更新 WorkbenchDeviceStep 使用 Spinner
- [x] 補充測試

### 1.2 補齊 ARIA 屬性（螢幕閱讀器支援）
- [ ] WorkbenchStepRail 補充 `aria-current`、`aria-label`
- [ ] MemoryGrid 補充 `role="grid"`、`aria-label`
- [ ] TagBindingStudio 批量操作按鈕補充 `aria-describedby`
- [ ] SourceCanvasSection 地址選擇器補充 `aria-selected`
- [ ] 補充焦點管理（步驟切換後自動聚焦）
- [ ] 補充測試

### 1.3 補充 Studio 主流程 E2E 測試
- [ ] 建立 `frontend/tests/e2e/studio/device-to-output.spec.ts`
- [ ] 建立 `frontend/tests/e2e/studio/cross-step-traceability.spec.ts`
- [ ] 執行測試驗證

---

## Phase 2: P1 中優先級改善（影響體驗）

### 2.1 補充跨步驟追蹤視覺反饋
- [ ] 為 focused 項目補充視覺高亮樣式
- [ ] 在 Inspector Panel 顯示追蹤路徑
- [ ] 補充測試

### 2.2 補充移動端響應式設計
- [ ] WorkbenchFrame 補充響應式斷點
- [ ] 小螢幕下 Inspector Panel 改用 Drawer
- [ ] 移動端下 StepRail 改為頂部 Tab
- [ ] 補充移動端 E2E 測試

### 2.3 補充 Inspector Panel 與 BottomSummaryBar 互動設計
- [ ] BottomSummaryBar 就緒指示器補充點擊跳轉
- [ ] Inspector Panel 補充 "返回列表" 按鈕
- [ ] 補充互動設計文件

---

## Phase 3: P2 低優先級改善（改進體驗）

### 3.1 補充鍵盤導航支援
- [ ] 補充鍵盤快捷鍵文件
- [ ] 為主流程補充快捷鍵（Ctrl+1/2/3/4）
- [ ] 為 Grid/Canvas 補充方向鍵導航

### 3.2 補充無障礙性測試
- [ ] 補充 a11y 單元測試（jest-axe）
- [ ] 補充 a11y E2E 測試（@axe-core/playwright）

### 3.3 補充元件測試覆蓋
- [ ] 補齊核心元件測試（DeviceForm、MemoryGrid、MappingCanvas）
- [ ] 補充 Wizard 子元件測試

---

## 當前狀態

**Phase**: 1.2 補齊 ARIA 屬性  
**進度**: 0/6 完成  
**下一步**: WorkbenchStepRail 補充 aria-current、aria-label
