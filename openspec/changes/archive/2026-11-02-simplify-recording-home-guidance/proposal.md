## Summary

在既有 Studio V2 增加可選用的記錄首頁，以目前工作區證據引導首次使用、已有來源、未完成草稿、正常記錄與異常處理。本批不切換預設入口，不改採集或資料庫引擎。

## Motivation

既有四步驟把使用者帶進設定，但日常開啟時仍需自行找出下一步。新版入口需讓同時存在的草稿、成功與異常記錄都可找到，不能以空畫面或綠燈掩蓋載入失敗與未知結果。

## Proposed Solution

- 新增 simple-recording-home 能力，首頁投影既有 workspace、source、group、delivery 與 lifecycle 證據，不新增設定 authority。
- 以同一 resolver 決定首頁提示優先序；提供空白、已有設備、未完成、正常、異常情境及混合／載入失敗回歸。
- 在 /studio/v2 以明確 view 選擇加入新版入口，保留原流程、entity locator、runtime return focus 與進階路徑。
- 依系統亮暗主題、單一主要下一步與按需詳情呈現，不暴露原始錯誤、憑證或內部識別碼。

## Impact

- Affected specs: simple-recording-home (new); datalink-workbench-v2-shell (modified); datalink-workbench-v2-settings (modified)
- Affected code:
  - Modified: `frontend/src/services/studioEntityLocator.ts`; `frontend/src/App.tsx`; `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`; `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx`
  - New: frontend/src/features/datalink/workbench-v2/recording-home/; frontend/tests/unit/workbench-v2/recording-home.test.tsx
  - Removed: none
- 依賴已完成 guide-step4-recording-workflow、expose-runtime-recording-rules、make-runtime-logs-actionable 的現行契約。後續 streamline-recording-task-editor、simplify-recording-lifecycle-recovery、adopt-simple-recording-entry 接續；不改寫舊完成證據。

本次更新納入 `DESIGN.md` 語意 tokens、system/light/dark 主題、RWD、簡潔文案與字型階層；新彈出設定為 viewport 寬高各75%的置中 modal。已選手機近滿版／小確認框較小例外，分類依 DESIGN.md；舊 demo 不構成此次確認。

- Modified design scope: DESIGN.md; frontend/src/contexts/ThemeContext.tsx; frontend/src/components/ThemeToggle.tsx; frontend/src/index.css; frontend/src/styles/tokens.ts; frontend/src/features/datalink/workbench-v2/tokens.ts; settings/UiSection.tsx 與 backendMappings.ts 的 theme intent adapter。

首頁改用可搜尋／篩選／分頁的緊湊記錄管理列表，設備僅短摘要與同頁管理，不無限展開卡片；目前以既有已載入資料投影，不假稱後端跨設備搜尋索引已存在。
