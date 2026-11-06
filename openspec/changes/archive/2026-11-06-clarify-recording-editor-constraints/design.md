## Context

簡易記錄編輯頁（`GroupEditor` 的 simple presentation，經 `RecordingEditorLayout` 分成來源、記錄方式、目的地三段）同時承載三種記錄：沒有 `recording_profile` 的舊式寬表記錄、瞬時值／累積值設定檔記錄（`managed_long_v1` 或 `managed_wide_profile_v1`），以及各點位計算記錄（`recording_profile.kind = 'mixed'`）。各類記錄的限制來自既有規格：各點位記錄只用受管資料表（`recordingMixedDraft` 的 `mixedDraftIssues`）、長表在同一 Schema 共用一張受管表（`long-table-group-recording`）、簡易頁無法表達的設定須保留摘要與進階入口且不得自動轉換（`simple-recording-editor` 的 Basic and advanced recording semantics remain faithful）、瞬時值記錄只能經明確的新版本改為平均（`window-average-recording`）。

2026-11-04 實測：新建各點位記錄在儲存前 Schema 與資料表即被鎖住（`GroupEditor` 以 `profile || layout === managed_long_v1` 鎖定，各點位記錄帶有 profile）；此項已改為 `destinationLock`，並以 `tests/unit/workbench-v2/recording-mixed-destination-lock.test.tsx` 驗證。其餘問題：已儲存長表記錄鎖定欄位沒有說明；「記錄方式」區段在各點位模式只有 `RecordingIntentHelp` 說明；`GroupPolicyFields` 在 profile 記錄隱藏多個欄位但沒有說明；`GroupEditorPreparation` 只在 `managed && group` 時渲染；`RecordingDestinationFields` 不讀取型錄時 `GroupIdentityFields` 仍可能顯示讀取中；`GroupReadinessPanel` 的診斷與 `RuntimeSetupContextPanel` 直接顯示 `issue.code` 與英文 `issue.message`；記錄點位搜尋框曾在瀏覽器顯示 `null` 而篩選狀態為空（瀏覽器自動填入或表單還原，未在乾淨瀏覽器重現）。

## Goals / Non-Goals

**Goals:**

- 每個停用或隱藏的欄位都在旁邊用一句白話說明原因與替代做法。
- 各點位記錄的記錄方式在「記錄方式」區段可一眼看到，並能直接跳到修改位置。
- 瞬時值設定檔記錄可在簡易頁明確修改取值方式與週期，走既有保存與結構確認流程。
- 就緒問題對操作者只顯示中文說明；代碼留給診斷。

**Non-Goals:**

- 不改任何鎖定規則本身（儲存方式、長表共用、各點位只用受管表、保存後固定），只補說明與入口。
- 不把舊式寬表記錄原地轉成設定檔或各點位記錄；只提供複製成新記錄的捷徑。
- 不改後端 API 或結構確認流程；不新增演算法。
- 不在本次修改連接器狀態的後端更新（資料庫輸出失敗記錄屬 `decouple-acquisition-from-delivery`）。

## Decisions

### 鎖定判斷集中於 destinationLock

目的地欄位的鎖定由 `GroupEditorCapabilities.tsx` 的 `destinationLock(draft, saved)` 決定：各點位記錄在保存前只鎖儲存方式、保存後全部鎖；設定檔與長表記錄維持原規則。鎖定原因文字也依同一個結果選擇，避免畫面說明與實際鎖定不一致。替代方案：在每個欄位各自判斷——規則分散，已造成目前的過度鎖定。

### 鎖定說明附複製入口

已儲存且 `destinationLock` 回傳 locked 的記錄，在目的地欄位上方顯示一句原因（長表：「長表記錄與同一 Schema 的其他長表記錄共用資料表，儲存後固定」；設定檔寬表：「這份記錄的資料表已綁定，儲存後固定」）與既有「複製成新記錄」動作。替代方案：開放修改並由後端拒絕——操作者要到保存才知道失敗。

### 記錄方式摘要沿用 RecordingEditorSummary 的分組

「記錄方式」區段在各點位模式顯示與右側摘要相同的分組（算法 · 間隔 · 點數），每組提供「到來源表格」動作，捲動並聚焦該組第一個點位的算法欄位。替代方案：在記錄方式區段再放一套逐點控制——與來源表格形成兩個編輯入口。

### 瞬時值記錄的取值方式以既有 draft 欄位明確編輯

簡易頁對 `recording_profile.kind = 'instantaneous'` 的記錄提供「取值方式」（最新值／算術平均）與「記錄週期」控制，寫入同一個 draft 的 `recording_profile.sample_reduction` 與 `interval`，經 `GroupEditor` 既有保存（CAS）成為新版本；平均需既有的量測用途確認（gauge consent）；保存後準備區塊要求重新檢查並確認結構。累積值記錄不提供平均（規格拒絕）。替代方案：導回基本模式——操作者需要切換工作流程，且新頁已宣稱是主要編輯入口。

### 舊式寬表記錄只顯示與複製

沒有 `recording_profile` 的記錄在「記錄方式」顯示「最新值（舊式寬表）· 每 N 秒」，N 取 `interval_seconds`；提供「以這些點位建立新記錄」捷徑，開啟新記錄並帶入相同點位與目的地連接器，由操作者選擇模式後保存。替代方案：原地轉換——規格禁止。

### 就緒問題在地化對照表

新增以 `issue.code` 為鍵的中文說明對照（至少 `schema-confirmation-required`、`write-group-apply-required`、`database-target-missing`、`source-revision-stale`、`schema-check-blocked`），編輯頁與執行儀表板都以對照文字為主文；未知代碼顯示「需要處理的設定問題」並保留診斷。診斷內顯示「代碼：中文說明（英文原文）」。替代方案：後端回傳中文——後端訊息屬 API 契約，不在本次範圍。

## Implementation Contract

**Behavior**
- 新建各點位記錄：保存前 Schema 與資料表可從清單選擇或自行輸入；儲存方式停用並顯示「各點位計算只使用受管資料表」；已被其他記錄擁有的資料表在清單中停用。
- 已儲存的長表或設定檔記錄：目的地欄位停用，上方顯示原因與「複製成新記錄」。
- 各點位記錄的「記錄方式」區段：顯示「最新值 · 15 秒 · 4 點」這類分組列，每列有「到來源表格」。
- 各點位記錄的進階區塊：顯示「間隔在各點位設定；不完整的列固定為部分結果；受管資料表的重複保護、實體鍵與來源紀錄欄位由系統決定」。
- 未保存的受管記錄：準備區塊位置顯示「保存設定後才能檢查資料表結構」。
- 瞬時值記錄：可選最新值或算術平均及週期；選平均時要求量測用途確認；保存後顯示需重新檢查結構。
- 舊式寬表記錄：顯示「最新值（舊式寬表）· 每 60 秒」及「以這些點位建立新記錄」。
- 目的地型錄不讀取時不顯示讀取中；讀取失敗顯示「無法連到目的地資料庫」與重試。
- 就緒問題：主文為中文說明；「診斷」內為代碼、中文與英文原文。
- 記錄點位搜尋框帶 `autoComplete="off"`。

**Interface / data shape**
- `destinationLock(draft: EditorDraft, saved: boolean): { locked: boolean; storageLocked: boolean }`（已存在）。
- 就緒說明以 i18n 鍵 `readinessIssue.<code>` 提供 en／zh-TW 文字。
- 不新增 API；瞬時值記錄更新沿用既有寫入群組更新請求的 `recording_profile`。

**Failure modes**
- 型錄讀取失敗：顯示錯誤與重試，不顯示為空清單。
- 保存瞬時值記錄的新取值方式被後端拒絕（例如非數值來源選平均）：顯示後端欄位錯誤，草稿保留，不假裝已套用。

**Acceptance criteria**
- vitest：`destinationLock` 各情境（已有）；鎖定說明與複製入口只在 locked 時出現；記錄方式分組與「到來源表格」聚焦；進階區塊說明與未保存準備提示；瞬時值記錄改平均時要求確認、payload 帶 `sample_reduction: 'mean'` 與新週期；舊式記錄顯示實際週期與捷徑；型錄不讀取時無讀取中文字；就緒對照表涵蓋列出的代碼且未知代碼有預設文字；搜尋框屬性。
- 以 1440／768／390 寬度在 dev server 實際截圖新建各點位記錄、已儲存長表記錄、舊式記錄三種畫面，附於交接；不保存、不連真 PLC。

**Scope boundaries**
- 範圍內：簡易記錄編輯頁與執行儀表板就緒清單的顯示、說明與既有保存流程的入口。
- 範圍外：鎖定規則、後端 API、連接器狀態更新、基本模式頁面的移除。

## Risks / Trade-offs

- [GroupEditor.tsx 499 行、recording-editor.json 500 行已達上限] → 新增內容前先拆出元件與 locale 命名空間（獨立任務），每批維持可建置。
- [瞬時值記錄改平均會改變後續寫入內容] → 只在明確操作與量測用途確認後保存；保存後必須重新確認結構才會生效，歷史資料不重算。
- [就緒對照表與後端新增代碼不同步] → 未知代碼有預設文字並保留診斷，測試鎖定已知代碼清單。

## Migration Plan

純前端變更，無資料遷移；回滾為還原前端版本。

## Open Questions

- 無。
