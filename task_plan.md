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

### 階段 10：建立 studio surface 與 API 維護文件
- [x] 盤點 `/studio`、`/studio/v2`、`/studio/runtime`、`/gateway/*`、`/test` 的定位與成熟度
- [x] 盤點前端實際使用的 query / mutation / SSE 與對應後端 API
- [x] 建立多份 Markdown 文件，分開記錄主線頁面、V2/runtime、gateway/test、API registry 與 gap roadmap
- [x] 建立一份可直接開啟的 HTML 總覽頁
- **狀態：** complete

### 階段 11：依產品優先順序重構 surface inventory 文件
- [x] 將 `/studio` 標示為完整版主線但先暫停
- [x] 將 `/studio/v2` 標示為預設入口與重點施作
- [x] 將 `/test` 拆成獨立文件並標示先暫停
- [x] 將 `/gateway/*` 拆成獨立 experimental 記錄
- [x] 更新 HTML 總覽以反映上述優先順序
- **狀態：** complete

### 階段 12：修正 HTML 文件台無資料並拆分為 html/js/css
- [x] 找出 HTML 無資料的 root cause
- [x] 將文件台拆為 `index.html`、`inventory.css`、`inventory.js`
- [x] 驗證 JS 語法與基本檔案關聯
- **狀態：** complete

### 階段 13：為 studio surface inventory 建立 changelog 機制
- [x] 在 `AGENTS.md` / `CLAUDE.md` 寫明 inventory 更新必須記錄
- [x] 新增 SQLite changelog CLI
- [x] 初始化 changelog database 並寫入本次變更
- [x] 更新 inventory README 使用說明
- [x] 驗證 CLI 與資料庫內容
- **狀態：** complete

### 階段 14：補齊 inventory onboarding 入口
- [x] 新增 `START_HERE.md` 作為第一入口
- [x] 新增 `context.json` 作為 machine-readable 摘要
- [x] 更新 `AGENTS.md` / `CLAUDE.md` 指向 onboarding 入口
- [x] 更新 `README.md` 閱讀順序與 onboarding 規則
- [x] 將本次文件更新寫入 SQLite changelog
- [x] 驗證文件與 changelog
- **狀態：** complete

### 階段 15：強化新對話自動接手入口
- [x] 新增 `CURRENT_STATE.md` 作為最新狀態快照
- [x] 更新 onboarding 文件與 `context.json`
- [x] 強化 `AGENTS.md` / `CLAUDE.md`，把 `studio` surface 任務必讀順序寫死
- [x] 補一筆 changelog 記錄這次接手機制升級
- [x] 補一筆 Codex memory note，讓新對話更容易直接接上
- [x] 驗證文件、JSON 與 changelog
- **狀態：** complete

## 關鍵問題
1. 每一步驟的變更有無互相依賴？（需依序實作，後面的變更可能基於前面已實作的程式碼）
2. 測試與 Lint 是否在每個階段實作完都要通過？（是，確保每個階段提交時皆為 clean 狀態）

## 已做決策
| 決策 | 理由 |
|------|------|
| 新建分支 `feature/datalink-workbench-v2` | 使用者要求另開分支實作，不直接污染 main |
| 日誌改為時間遞增排序 (Append to array end) | 與 auto-scroll-to-bottom 對齊，提供流暢的控制台滾動 UX。 |
| 將 surface inventory 拆成多份 md + 一份 html 總覽 | 避免單一文件過大，讓後續維護與 change 規畫更容易定位 |

## 遇到的錯誤
| 錯誤 | 嘗試次數 | 解決方案 |
|------|---------|---------|

## 備註
- 隨著進度更新階段狀態：pending → in_progress → complete
- 做重大決策前重新讀取此計畫（注意力操縱）
- 記錄所有錯誤，避免重複
