## Why

操作者在簡易記錄編輯頁反覆遇到「欄位不能選、區塊不見、代碼看不懂」：新建的各點位記錄在儲存前就鎖住 Schema 與資料表；已儲存的長表記錄鎖住儲存方式、資料表、Schema 與佈局卻沒有說明；「記錄方式」區段只有說明文字，實際的算法與間隔藏在來源表格；各點位記錄的進階區塊少了成列間隔等欄位，「受管資料表準備」要保存後才出現，畫面都沒有交代原因；瞬時值記錄的取值方式只能回舊的基本模式修改；就緒問題直接顯示 `schema-confirmation-required` 等英文代碼。這些限制多半是規格刻意的，但沒有說明時操作者會判斷為故障。

## What Changes

- 新建的各點位記錄在第一次儲存前可從清單選擇或自行輸入 Schema 與資料表，儲存方式固定為受管資料表並說明原因（已實作並通過測試，本提案補上規格）。
- 已儲存而鎖定目的地欄位的記錄（長表、設定檔記錄）在欄位旁說明為何固定，並提供「複製成新記錄」入口。
- 各點位記錄的「記錄方式」區段列出各點位的算法與間隔摘要（依算法與間隔分組計數），並可跳到來源表格對應列修改。
- 各點位記錄的進階區塊說明為何沒有成列間隔、不完整的列、重複保護、實體鍵與來源紀錄欄位；尚未保存時在準備區塊顯示「保存後才能檢查資料表結構」。
- 瞬時值設定檔記錄可在簡易編輯頁明確修改取值方式（最新值／算術平均）與記錄週期，儲存為新版本並要求重新確認結構；沒有設定檔的舊式寬表記錄顯示「最新值（舊式寬表）」與實際週期，並提供以目前點位複製成新記錄的捷徑，不做原地轉換。
- 目的地資料表清單在不會讀取時不顯示「正在讀取」；讀取失敗時顯示連不上目的地資料庫與重試。
- 就緒問題以中文說明呈現，原始代碼與英文訊息只保留在「診斷」內，且附上對應中文；涵蓋記錄編輯頁與執行儀表板。
- 記錄點位搜尋框關閉瀏覽器自動填入，避免出現與篩選狀態不一致的殘留文字（例如 `null`）。

## Non-Goals (optional)

詳見 design.md 的 Goals / Non-Goals。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `simple-recording-editor`: 新增目的地欄位鎖定與說明、記錄方式摘要、瞬時值記錄的明確取值方式編輯、就緒問題在地化與目的地讀取狀態等要求。

## Impact

- 前端：frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor.tsx、GroupEditorCapabilities.tsx、GroupIdentityFields.tsx、GroupPolicyFields.tsx、RecordingDestinationFields.tsx、GroupReadinessPanel.tsx；frontend/src/features/datalink/workbench-v2/recording-editor/RecordingMixedEditor.tsx、GroupEditorPreparation.tsx、RecordingCompactMemberTable.tsx、RecordingEditorSummary.tsx；frontend/src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel.tsx；en／zh-TW locales。
- 後端：無新 API；瞬時值記錄的取值方式與週期沿用既有寫入群組更新與結構確認流程。
- 行數：GroupEditor.tsx 為 499 行、兩個 recording-editor.json 為 500 行，新增內容前需先拆檔。
