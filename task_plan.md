# 任務計畫：實作 Datalink Workbench V2 變更

## 目標
依序完成 6 個 `spectra-apply` 變更的實作，確保前端與後端功能整合完畢，並通過所有測試與 lint 檢查。

## 目前階段
階段 1：建立分支與規劃

## 各階段

### 階段 1：建立分支與規劃
- [x] 另開 git branch `feature/datalink-workbench-v2`
- [x] 建立並初始化規劃檔案（`task_plan.md`, `findings.md`, `progress.md`）
- **狀態：** complete

### 階段 2：實作 datalink-workbench-v2-shell
- [x] 執行 `spectra status --change "datalink-workbench-v2-shell" --json` 取得任務
- [x] 依序完成 shell 的任務
- [x] 執行單元測試與驗證
- [x] 將變更標記為 done 並封存或進行下一階段
- **狀態：** complete

### 階段 3：實作 datalink-workbench-v2-step1-device
- [x] 執行 `spectra status --change "datalink-workbench-v2-step1-device" --json` 取得任務
- [x] 依序完成 device 的任務
- [x] 執行單元測試與驗證
- **狀態：** complete

### 階段 4：實作 datalink-workbench-v2-step2-rule
- [x] 執行 `spectra status --change "datalink-workbench-v2-step2-rule" --json` 取得任務
- [x] 依序完成 rule 的任務
- [x] 執行單元測試與驗證
- **狀態：** complete

### 階段 5：實作 datalink-workbench-v2-step3-mapping
- [x] 執行 `spectra status --change "datalink-workbench-v2-step3-mapping" --json` 取得任務
- [x] 依序完成 mapping 的任務
- [x] 執行單元測試與驗證
- **狀態：** complete

### 階段 6：實作 datalink-workbench-v2-step4-database
- [x] 執行 `spectra status --change "datalink-workbench-v2-step4-database" --json` 取得任務
- [x] 依序完成 database 的任務
- [x] 執行單元測試與驗證
- **狀態：** complete

### 階段 7：實作 datalink-workbench-v2-settings
- [x] 執行 `spectra status --change "datalink-workbench-v2-settings" --json` 取得任務
- [x] 依序完成 settings 的任務
- [x] 執行單元測試與驗證
- **狀態：** complete

### 階段 8：最終驗證與整合
- [x] 執行完整建置：前端 build、複製 embed static、後端 build
- [x] 執行 Go 和 React 的所有單元與整合測試
- [x] 檢查單檔行數限制
- **狀態：** complete

### 階段 9：優化項目補齊與日誌系統升級
- [x] 實作診斷日誌面板 (RealtimeLogsPanel.tsx) 與 `stale` / `recovered` 狀態去重日誌。
- [x] 將日誌改為 append 末尾以匹配 auto-scroll 效果。
- [x] 撰寫整合單元測試並通過 `npm run test` 與 `go test`。
- **狀態：** complete

## 關鍵問題
1. 每一步驟的變更有無互相依賴？（需依序實作，後面的變更可能基於前面已實作的程式碼）
2. 測試與 Lint 是否在每個階段實作完都要通過？（是，確保每個階段提交時皆為 clean 狀態）

## 已做決策
| 決策 | 理由 |
|------|------|
| 新建分支 `feature/datalink-workbench-v2` | 使用者要求另開分支實作，不直接污染 main |
| 日誌改為時間遞增排序 (Append to array end) | 與 auto-scroll-to-bottom 對齊，提供流暢的控制台滾動 UX。 |

## 遇到的錯誤
| 錯誤 | 嘗試次數 | 解決方案 |
|------|---------|---------|

## 備註
- 隨著進度更新階段狀態：pending → in_progress → complete
- 做重大決策前重新讀取此計畫（注意力操縱）
- 記錄所有錯誤，避免重複
