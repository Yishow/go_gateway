# 命名與換算映射 404 恢復

舊流程收到 `workspace_mapping_not_found` 後，明確重試仍以過期的映射 ID 呼叫更新，導致同一來源規則的待存列持續失敗。此次沒有追溯映射最初消失的原因。

## 修改範圍

- `frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts`：明確重試 404 時先查詢目前映射；依來源規則、設備與位址取得現在的 ID，缺少映射則走既有建立 API。查詢期間保留最新草稿、合併重複重試，查詢失敗保持暫停。後端 ownership、來源規則與候選衝突檢查維持原契約。
- `frontend/src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page.tsx`：初次載入映射失敗仍顯示錯誤；已有資料的重新查詢失敗保留編輯介面與列上的錯誤／重試操作。
- `frontend/tests/unit/workbench-v2/mapping-autosave-queue.missing.test.tsx`：涵蓋多列映射消失、替換 ID、查詢期間新草稿、重複重試，以及查詢失敗後再次重試。
- `frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx`：驗證重新查詢失敗後仍可看到草稿並再次重試。

相關契約：[Studio V2 workspace](../../openspec/specs/studio-v2-workspace/spec.md)。本次未修改規格或後端。

## 2026-10-24 本機驗證

- PASS：新增佇列測試在修正前 3/3 失敗，修正後通過；整頁測試套回原先錯誤 gate 時失敗，恢復修正後通過。
- PASS：`cd frontend && npm run lint`。
- PASS：`cd frontend && npm test -- --run`，216 個檔案／1,382 個測試。
- PASS：`cd frontend && npm run build`；保留 Browserslist 資料較舊與 bundle 大小警告，沒有擴大處理。
- PASS：`git diff --check`、`make check-lines`。
- PASS：agent-browser／Chromium 1440×900，真實 MappingTable 與 autosave hook、替代 API 資料；一次失敗更新後，按另一列的重試按鈕，3 列建立成功，保留修改過的名稱。當時的 404／恢復截圖未留存在儲存庫，無法作為此次可重現證據；臨時測試頁已移除。
- NOT RUN：真設備、現有工作區寫入、SQL 寫入驗收、Windows、正式部署。瀏覽器成功狀態只代表隔離替代 API 流程；未修補目前工作區資料，未重啟服務或同步嵌入式 assets，未 commit／push。

行數警告：autosave hook 412 行、reconciliation 測試 315 行，均低於 500 行阻擋值。後續若需擴充這兩處，先分離明確重試的恢復邏輯及獨立恢復測試情境；此次保留小範圍修正。原有已暫存的 launcher 變更未納入本次範圍。

## 2026-10-25 提交前覆核

- 前端完整測試 1,812／1,812 通過，lint 與 build 通過；404 恢復以既有佇列與整頁測試覆核，未重建原臨時瀏覽器頁。
- Playwright／Chromium 的 8 個頁面案例全部通過：Basic／Advanced 1280×900 捲動邊界、1600×1000 與 1024×720 即時值溢出、1600×900 與 1024×720 工作台版面。使用實際 React 元件與隔離 fixture／替代 API。
- 此次實測截圖：[Basic 1024×720](evidence/mapping-layout-20261025/basic-1024.png)、[Advanced 1024×720](evidence/mapping-layout-20261025/advanced-1024.png)。這些是版面驗證，不代表原工作區、真設備或 SQL 寫入驗收。
