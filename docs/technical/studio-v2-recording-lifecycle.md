# Studio V2 已保存記錄：停止、移除、重建與切換

本文說明 `/studio/v2` Step 4 的已保存記錄管理。Basic（基本設定）與 Advanced（進階）共用同一份伺服器操作狀態；切換模式不會建立另一份記錄，也不會重新送出變更。本文對照目前程式與畫面文字，不代表已完成正式瀏覽器驗收。

一般設定、來源及計量語意見 [迴路記錄說明](studio-v2-circuit-recording.md)；既有 partial-source／custom 的已審核切換仍使用該頁原流程，沒有被這份「完整來源的新世代」流程取代。

## 從哪裡開始

1. 在 Step 4 選取設備，找到已保存記錄的狀態與操作區；Advanced 也可選取同一群組。
2. 依畫面顯示的可用操作選擇下表動作。按鈕是否可用由伺服器的 `actions[]` 決定，不只看 draft／applied／disabled／deleted 字樣。
3. 若顯示狀態尚未確認，按「重新讀取可用操作」。過期版本、無法確認的來源或目的地、缺少資料歸屬證據都不能跳過。

| 畫面操作 | 用途與界線 |
| --- | --- |
| 編輯已保存設定 | 編輯同一群組的設定；保存草稿不等於套用。已有生效版本時，舊版本繼續依原規則處理，不會因草稿改動而換用新規則。 |
| 停止新記錄 | 停止接收新資料，不清除已接受的待送資料，也不刪 SQL 表、歷史或交付收據。 |
| 移除設定 | 先檢視影響，再按「移除設定並保留歷史」；保留可查閱的已移除項目，不是刪除外部 SQL 資料。 |
| 重新建立記錄 | 對可重建的已移除項目開啟審核，使用新群組 ID 與下一世代；不復活原群組，也不走一般新增。 |
| 變更格式或目的地 | 對符合條件、已有套用版本的受管 profile 開啟完整來源的變更審核；不是任意轉換或直接套用。 |

已停用的項目不會再提供可執行的「停止新記錄」。停止後若只要恢復相同設定，使用原群組既有的審核開始流程，不必先移除再重建。受管 profile 的開始仍回到既有的指定群組審核，不能用一般套用繞過。

### 找回已移除的項目

Basic 在所選設備下提供「已移除的記錄（數量）」展開區，保留原名稱及管理操作；「檢視已移除記錄」可開啟 Advanced 的該群組。Advanced 可選取已移除項目，原設定編輯器唯讀，但符合條件時仍可按「重新建立記錄」。這些項目不會重新出現在普通新增記錄候選中。

已有保留中的變更時，伺服器會阻擋競爭操作；請使用「查看原操作」，不要藉切換模式、設備或新增群組繞過。

## 重建採用哪份設定

伺服器提供的 `successor_seed` 是審核起點，不是套用權限：

- 曾經套用的記錄以原不可變的已套用版本為基礎，不採用後來編輯但未套用的草稿。
- 只有確實從未套用、且伺服器證明沒有歷史的已移除記錄，才以原保存定義為基礎；空的 `applied_revision` 本身不是證明。
- 保留完整來源、profile kind、記錄週期、reduction 與已確認的 counter 語意；這個對話框不提供來源或計量語意編輯。
- 缺失、過期或不相容的原始證據會阻擋審核。畫面上的候選與按鈕可用性都不能取代後端再次驗證。

### 格式與目的地限制

| 情境 | 支援範圍 |
| --- | --- |
| 受管 profile 的 plain 寬表 ↔ plain 長表 | 僅同引擎 SQLite／PostgreSQL，完整同型態來源，且每個來源都有等價表示。 |
| 同布局、同儲存模式換目的地 | 同引擎且重新查證目標能力、結構與擁有權。MySQL、MariaDB、SQL Server 各自僅限受管長表，不能彼此視為同引擎。 |
| 已移除的 Basic 原基本寬表 | 僅 SQLite／PostgreSQL plain，原布局、目的地與完整語意不變的重建；新世代仍需自己的資料表準備證據。 |
| 不支援的組合 | 跨引擎、plain／Timescale／native_monthly 互轉、部分來源的新 successor、改 profile kind、非等價欄位、MySQL／SQL Server 寬表，以及任意 custom 重建。 |

Timescale 只保留既有長表同模式契約，不提供寬表 Timescale。選單列出布局或連線不表示該組合已獲准；最終以伺服器審核與阻擋理由為準。SQL Server 本次 successor 實機驗收限制見文末。

## 三步完成重建或變更

從「重新建立記錄」或「變更格式或目的地」進入「變更已保存的記錄」對話框。一次只進行畫面當下的主要操作，不把所有底層 API 階段當成可任意排列的步驟。

### 1. 選擇變更

- 檢視「記錄名稱」、「資料表布局」及「目的地」。Basic 原基本寬表的布局與目的地固定；非 plain 的 profile 也不能在此切換布局。
- 按「檢視變更影響」。這是伺服器本機唯讀審核，不會保留新世代、不會產生資料表預覽 token、不會探測外部目的地或執行 DDL。
- 尚未送出保留設定時，可按「取消」、Escape 關閉，或在影響畫面按「返回」重新選擇；沒有因此建立保留或執行資料庫動作。

### 2. 檢視影響與準備資料表

1. 檢查新世代、保留來源數、原目的地、新目的地及「歷史與待交付明細」。有阻擋理由時先處理，不可繼續。
2. 明確按「保留設定並預覽資料表」。先保留同一條記錄歷程的新世代，再取得資料表預覽；必要能力檢查及外部唯讀探測由後端執行。此時新世代仍不接收資料，也沒有因預覽而執行 DDL。
3. 檢視「伺服器回傳的資料表準備計畫」，包括表名、欄位及「建立資料表／沿用相容資料表／新增欄位」。只有明確按「確認此資料表計畫」才授權該計畫的結構準備；沒有 DDL 的相容計畫也要確認。

若只完成保留設定，畫面可提供「預覽資料表準備」繼續原操作。資料表確認成功只表示準備證據已保存，不會啟動新記錄、停止舊版本或移動已完成的歷史資料。

### 3. 確認開始時間

1. 在「未來開始時間（您的當地時間）」選擇明確的未來時間 B。第一步的審核不是開始時間的最終授權。
2. 按「審核開始時間與就緒狀態」。伺服器以目前版本重新審核，並查證就緒狀態；此步不執行 DDL，也不是最終啟用。
3. 檢查「已審核邊界」及阻擋訊息，再按伺服器結果對應的「確認首次開始記錄」或「確認排程切換」。時間、設定或連線版本變更後須重新審核，不能拿舊審核摘要配新時間。

| 最後確認 | 實際意義 |
| --- | --- |
| 確認首次開始記錄（`first_start`） | 僅限伺服器證明整條歷程從未真正套用或接收資料的已移除來源；不虛構舊啟用版本、結束事件或切換歷史。是否符合仍由來源與引擎契約決定。 |
| 確認排程切換（`cutover`） | 依已審核未來邊界把符合條件的新觀測交給新世代；舊表、歷史及已定稿效果保留。 |

歸屬以來源觀測時間 `observed_at` 判定：切換時早於 B 的觀測屬舊世代，等於或晚於 B 的觀測屬新世代，不用到達伺服器的時間替代。首次開始沒有舊 owner，B 前的新觀測會拒收。切換中尚未定稿的觀測可能依已保存的邊界安全交接，保留原接受來源與交接紀錄；已定稿的 payload、outbox 與 receipt 不搬家、不重算、不雙寫。

累積計數器在新世代重新建立基準。例如舊世代最後是 100，新世代第一筆 107 只建立 baseline，差值為 NULL，不是跨世代補算的 7；下一筆可靠值 110 才得到差值 3。NULL 也不是零。

## 每個成功訊息只代表自己的階段

| 畫面或證據 | 不能據此推論 |
| --- | --- |
| 「設定已保留，尚未開始新記錄。」 | 不代表已準備資料表、已套用或正在收集。 |
| 「資料表準備已確認，尚未開始新記錄。」 | 不代表新世代已開始，也不代表任何記錄列已提交。 |
| 「已記錄確認時間的開始／切換排程。」 | 不代表設備正在執行、採集成功或 SQL 已提交。 |
| Runtime running／本機已接受 | 不代表外部 SQL 已提交。 |
| 匹配版本的 SQL receipt／SQL committed | 不等於另外完成外部 SQL 讀回驗證；試寫與清理也須各看自己的結果。 |

開始結果須與 runtime、同一已套用版本的交付狀態及實際 SQL 證據分開確認。HTTP 200 或某一 operation 的 `succeeded` 不可代替上述證據。

## 停止、移除與歷史查閱

- 「停止新記錄」停止新 intake；已接受的待送資料仍依原凍結目的地交付。資料庫離線時可能持續受阻，停止不是清空、撤回或保證已送完。
- 「移除設定」需先查證已接受資料的歸屬；連線身分已變更或歸屬不明時可能拒絕。若畫面仍允許停止，可先停止新記錄調查，不可略過移除保護。
- 移除確認及變更第二步都可檢視歷史。每個世代與已套用版本保留原連線資訊；尚未套用的世代會明示沒有已套用版本。

| 明細欄位 | 單位與解讀 |
| --- | --- |
| 待處理樣本數 | 尚待處理的來源樣本。 |
| 待交付資料列數 | 未完成交付的 outbox 資料列。 |
| 受阻資料列數／結果未知資料列數／隔離資料列數 | 各類交付列的狀態，不是另一批樣本，也不可與待交付總數再次相加。 |
| 已交接樣本數 | 由原群組／版本交給目標群組／版本的樣本；明細保留兩邊身分。 |

樣本、資料列與收據不是同一單位；不要相加成單一「待送筆數」。此歷史面板提供身分及待送摘要，不是 SQL 原始歷史值查詢器。

## 回應不明、重整與切換模式

### 先回到原操作

送出前，瀏覽器會保存原請求的復原資料。它以 workspace、原群組、lineage 及最初 request operation ID 固定識別，不改用新世代 ID、schema operation ID、first-start result ID 或最新 preview 取代。

1. 關閉對話框、切 Basic／Advanced 或在同一瀏覽器工作階段重整後，從提示按「查看原操作」。尚未完成的操作先保留未知狀態，不會自動重送。
2. 按「查詢原操作狀態」。此查詢只讀本機已保存證據，不探測目標、不執行 DDL、不取消操作，也不修復外部表。
3. 查到原結果後，按畫面允許的下一步處理。原 reservation 回傳的是不可變歷史結果，不是目前版本；後續步驟會另外讀取 current CAS，不能修改原 envelope。

復原使用 `sessionStorage` 與記憶體，不保證跨新分頁、新瀏覽器或清除儲存後自動找回。儲存失敗會阻止新請求；復原資料損壞但仍可辨識原操作時，只能查詢，不能據此重送或建立替代操作。不要以清除瀏覽器資料解除保留。

### 何時能重送

網路中斷、逾時、HTML／代理回應、5xx 或無法驗證的成功內容均可能是 unknown，不能當成確定失敗。先查原操作；只有畫面依已驗證證據提供「重送完全相同請求」時才能明確重送，不會自動重試 DDL。

| 原送出階段 | 可能提供 exact retry 的條件 |
| --- | --- |
| reserve | 原 scope 查詢回傳可辨識的 `WRITE_GROUP_NOT_FOUND`，且保留了完整原請求；一般查詢錯誤不是此證明。 |
| schema-confirm | 原 schema operation 為 `preview_not_submitted`，或為 `succeeded` 但尚不能確認群組已綁定該準備證據。重送仍使用原 token、digest、schema operation ID 與完整原請求；不是另造 DDL。 |
| 最終首次開始／切換確認 | 對應原 root 的 `first_start`／`cutover` 為 `not_submitted`；重送原確認，不偷偷改 B 或操作 ID。 |

唯讀模式、缺少完整請求、明確拒絕或不符合當下證據時，不一定提供重送。若查到資料表成功且群組綁定吻合，直接回到確認開始時間；若查到最終成功，保留已排程結果，後續查詢失敗也不能把該成功改成未知。`pending`、`running`、`failed`、`partial`、`unknown` 的 schema 結果不能被新計畫遮蓋。

### 遺失預覽回覆的窄例外

只有遺失的是 `schema-preview` 回覆、沒有可用本地 preview，且從未送出本地 schema-confirm（包括結果未知的確認），才可能顯示「重新取得預覽」。還必須由成功、完整且吻合原 root 的狀態查詢證明：

- 原 reserve 身分完全一致；
- `cutover` 與 `first_start` 都是 `not_submitted`；
- 所有 schema 項目都是 `preview_not_submitted`，或完全沒有項目。

按下後會再次查原狀態並更新 current CAS，只為同一保留設定取得預覽；不重新 reserve、不換 intent、不重建 root、不自動確認 DDL。先前 preview 的操作編號與狀態仍保留可見。任何 pending／running／failed／partial／unknown／succeeded schema、任何 activation 或身分查證失敗都排除此例外。

已取得且仍可驗證的原預覽，則依原計畫與目前版本檢查繼續；它與「遺失預覽」例外不同。新預覽本身永遠不是證明先前 DDL 沒有發生的依據。

### 關閉不等於取消，並非所有阻擋都能在畫面復原

一旦送出保留請求，對話框的「關閉」只關閉畫面；Escape 也不取消已保留、準備或排程的操作。「返回」只在尚未保留的影響審核提供。雙擊、重整或切模式不應用來嘗試撤回。

原來源、目標連線身分或定義漂移後，舊證據可能永久不再吻合；保留中的操作會安全阻擋。現在沒有一般「取消保留」、「重新綁定」或 rebase 按鈕，也不能保證重新讀取就能修好。請保留原操作識別與狀態，交由維護者診斷；不要改操作 ID、直接改設定庫、刪表或用一般新增繞過。此文件不授權人工清除保留或歷史。

## API 與實作對照

以下是實際路由，不是新增七、八條階段 URL：

- `GET /api/v1/datalink/studio-v2/workspace/write-groups/{id}/lifecycle`：本機唯讀的可用操作、原始 seed 與歷史摘要。
- `POST /api/v1/datalink/studio-v2/workspace/write-groups/{id}/successor-review`：初次本機唯讀審核。
- `POST /api/v1/datalink/studio-v2/workspace/write-groups/{id}/successor`：以 `action` 分派 `reserve`、`schema-preview`、`schema-confirm`、`activation-review`、`activation-readiness`、`activation-confirm`、`first-start-confirm`、`operation-status`。

`{id}` 固定為原 predecessor。`operation-status` 只帶原 `workspace_id`、`lineage_id`、`operation_id`；schema operation ID 與 current stage revisions 是不同資料，不可互換。合法舊 v1 activation envelope 沒有 `activation_kind` 時僅走原 cutover；不可猜成 first-start。

- [路由註冊](../../internal/api/router_studio_v2_write_group_routes.go)、[階段 dispatcher](../../internal/api/handlers/studio_v2_workspace_write_group_successor_handler.go)、[首次開始與原操作查詢](../../internal/api/handlers/studio_v2_workspace_write_group_successor_recovery.go)。
- [可用操作](../../internal/datalink/workspace/write_group_lifecycle_actions.go)、[原始設定 seed](../../internal/datalink/workspace/write_group_lifecycle_successor.go)、[原操作狀態](../../internal/datalink/workspace/write_group_successor_operation_status.go)。
- [共同管理入口](../../frontend/src/features/datalink/workbench-v2/steps/step4/SavedRecordingManagement.tsx)、[三步對話框](../../frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeDialog.tsx)、[操作與復原控制器](../../frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/RecordingChangeController.ts)、[復原條件](../../frontend/src/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeHelpers.ts)。
- [繁中操作文字](../../frontend/src/i18n/locales/zh-TW/recording-lifecycle.json)、[繁中三步文字](../../frontend/src/i18n/locales/zh-TW/recording-lifecycle-dialog.json)、[英文操作文字](../../frontend/src/i18n/locales/en/recording-lifecycle.json)、[英文三步文字](../../frontend/src/i18n/locales/en/recording-lifecycle-dialog.json)。
- [本次規格](../../openspec/changes/manage-recording-lifecycle-safely/specs/recording-lifecycle-management/spec.md)、[設計與驗收例外](../../openspec/changes/manage-recording-lifecycle-safely/design.md)、[任務狀態](../../openspec/changes/manage-recording-lifecycle-safely/tasks.md)。

## 驗證範圍

本頁與 task 4.2 文件檢查只核對操作文字、程式契約、連結、行數及空白；不把程式存在或文件完成當作功能驗收。正式 browser tests、截圖及逐 scenario 驗收屬 task 4.3，截至本頁此次更新為 NOT RUN，未提供假圖或替代截圖。

[SQL Server 文件](studio-v2-sqlserver-circuit-recording.md) 的 2026-10-27 PASS 屬先前長表功能。本次 successor lifecycle 的實際 SQL 寫入／讀回／cleanup 為 NOT RUN：依 2026-10-28 使用者「同意，若無法就跳過」的同機驗收例外，同機原生 private-prefix preflight 未找到官方文件支持的 workspace-only 完整 config/system root，因而安全停止；沒有下載、安裝或啟動 SQL Server engine。這不是已證明 runtime 啟動失敗，也不是本次 SQL 驗收通過。產品能力、TLS、結構、擁有權保護均未放寬。
