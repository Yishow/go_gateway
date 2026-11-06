## 1. 目的地欄位鎖定

- [x] 1.1 依「鎖定判斷集中於 destinationLock」實作 Destination locks are scoped and explained 的鎖定部分：新建各點位記錄保存前可選擇或輸入 Schema 與資料表、儲存方式固定為受管並顯示原因、他人擁有的表停用；已儲存記錄維持鎖定。證據：`tests/unit/workbench-v2/recording-mixed-destination-lock.test.tsx` 3 項通過，記錄相關 231 個測試檔 2964 項通過（2026-11-04），dev server 實測新記錄 Schema／資料表可操作且未送出寫入。
- [x] 1.2 先拆檔以容納新內容：把 `GroupEditor.tsx` 的目的地與記錄方式區段組裝抽成獨立元件、把 recording-editor 新文字放入新的 locale 命名空間；行為不變，以 `bash scripts/check_file_lines.sh` 通過且既有記錄測試全數通過驗證。
- [x] 1.3 依「鎖定說明附複製入口」完成 Destination locks are scoped and explained 的說明部分：已儲存且鎖定的長表與設定檔記錄，在目的地欄位上方顯示對應原因句與「複製成新記錄」；vitest 驗證只在 locked 時出現、長表與設定檔原因文字不同、按鈕呼叫既有複製動作。 [after: 1.2]

## 2. 記錄方式

- [x] 2.1 依「記錄方式摘要沿用 RecordingEditorSummary 的分組」實作 Recording method is visible where it is named 的各點位部分：記錄方式區段列出「算法 · 間隔 · 點數」分組，「到來源表格」捲動並聚焦該組第一個點位的算法欄位；vitest 以 3 個 latest 15 秒、1 個 mean 60 秒驗證分組文字與聚焦目標。 [after: 1.2]
- [x] 2.2 依「舊式寬表記錄只顯示與複製」實作 Recording method is visible where it is named 的舊式部分：沒有設定檔的記錄顯示「最新值（舊式寬表）· 每 N 秒」與「以這些點位建立新記錄」，捷徑開啟新記錄並帶入相同點位與連接器、原記錄不變；vitest 驗證文字、帶入內容與未送出原記錄更新。 [after: 1.2]
- [x] 2.3 依「瞬時值記錄的取值方式以既有 draft 欄位明確編輯」實作 Instantaneous recordings edit their reduction explicitly：瞬時值記錄提供最新值／算術平均與週期控制，選平均需量測用途確認，保存後準備區塊要求重新檢查結構，累積值不提供平均；vitest 驗證保存 payload 的 `sample_reduction` 與 `interval`、確認流程與累積值無平均選項。 [after: 1.2]

## 3. 說明與狀態

- [x] 3.1 實作 Hidden and pending areas state why：各點位記錄的進階區塊說明被隱藏欄位的原因；未保存受管記錄在準備區塊位置顯示「保存設定後才能檢查資料表結構」；目的地型錄不讀取時不顯示讀取中、讀取失敗顯示無法連到目的地資料庫與重試；vitest 驗證三種狀態。 [after: 1.2]
- [x] 3.2 依「就緒問題在地化對照表」實作 Readiness issues are shown in plain language：`GroupReadinessPanel` 與 `RuntimeSetupContextPanel` 以 `readinessIssue.<code>` 中文說明為主文，診斷內顯示代碼、中文與英文原文，未知代碼用通用說明；記錄點位搜尋框設定 `autoComplete="off"`；vitest 驗證五個已知代碼、未知代碼與搜尋框屬性。 [after: 1.2]

## 4. 驗收

- [x] 4.1 在 dev server 以 1440／768／390 寬度截圖新建各點位記錄、已儲存長表記錄、舊式記錄與瞬時值記錄（不保存、不連真 PLC），並執行 `npm test -- --run`、`npx tsc --noEmit`、相關檔案 eslint 與 `bash scripts/check_file_lines.sh`，結果與未執行項目記入交接。 [after: 1.3, 2.1, 2.2, 2.3, 3.1, 3.2]

2026-11-06 前次覆核（歷史）：8／9 項完成。已修復「到來源表格」被搜尋／問題篩選遮住及非平均錯誤誤導為平均確認的問題，相關測試、TypeScript 型別及選定 ESLint 通過。前端全套兩輪皆 632 檔通過、1 檔失敗，失敗的 `recording-shared-configuration-navigation.test.tsx` 在修復前隔離快照亦可重現；修復後瀏覽器因沙箱 Mach port 權限無法啟動，1440／768／390 截圖未重驗，既有 evidence 圖片僅為歷史版本。當時 4.1 維持未完成，不因局部通過而 archive。共同測試結果與限制見 [審查修復後覆核](../../decouple-acquisition-from-delivery/verification-evidence.md#2026-11-06-審查修復後覆核)。

2026-11-06 本輪驗收：重新執行前端全套 633 檔／7277 項全數通過，TypeScript、47 檔 ESLint、行數與 diff 檢查通過；四種記錄於 1440／768／390 的 12 張截圖及 3 張操作畫面已補齊。僅用既有模擬頁，不保存、不連真 PLC；未修改產品程式，歷史失敗此次未重現且未判定根因。詳見 [本輪驗收紀錄](verification-evidence.md)；本輪證據位於 `evidence/2026-11-06/`，未 archive 或提交。
