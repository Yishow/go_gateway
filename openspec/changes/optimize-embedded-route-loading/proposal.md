## Summary

縮小 Studio 記錄首頁首次需要的程式與傳輸量，保留穩定載入框架、既有路由及完整嵌入式交付。

## Motivation

2026-11-06 雲端 main 0f54740 的本機 HTTP／build 證據量得首次資源約 2.77 MB，其中 JS 約 2.19 MB，皆為未壓縮位元組，不是瀏覽器下載時間或 Mac 體感。首頁 eager graph 包含進階編輯程式；改進不能以空白頁、錯誤快取或缺 chunk 換取較小數字。

## Proposed Solution

- 保留立即可見的記錄頁 reading frame，將非首頁必需的 advanced/editor 模組延後載入，deep link 仍到原目的地。
- 對嵌入式 hashed JS/CSS 提供可協商 gzip 與 identity，明確區分 HTML 重驗與 immutable hashed assets。
- 以資源圖、HTTP 契約、cold/warm/deep-link 的瀏覽器紀錄分別驗證傳輸、解析與互動改善。
- 不修改 product route 意義；舊 embedded spec 的 /studio 已退休敘述依 legacy-studio-retirement 修正，保留完整 asset graph 契約。

## Impact

- Affected specs: embedded-frontend-delivery
- Affected code:
  - Modified: frontend/src/App.tsx, frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx, frontend/vite.config.ts, internal/web/embed.go, frontend/tests/unit/app-routing-lazy-load.test.tsx, frontend/tests/e2e/embedded-frontend-delivery.spec.ts
  - New: frontend/src/pages/datalink/workbench-v2/AdvancedWorkbenchV2Page.tsx, internal/web/embed_delivery_test.go
  - Removed: 無。
- 與 add-windows-tray-runtime-logs 共用 App 接點，runtime logs 的路由、初始化、shutdown 與 diagnostics 仍由原案負責。本案不調整其完成狀態。

## 驗收順序更新（2026-11-06）

使用者最新要求：「後續要做browser驗證的，等最後做完在一起上Mac驗」。四案先在雲端完成非瀏覽器實作與verify/review/fix，全部實作完成後才統一在Mac做真瀏覽器驗收。沒有取消browser條件，也沒有授權雲端擴權。未驗browser與最後結案維持待辦；未完整完成不archive。
