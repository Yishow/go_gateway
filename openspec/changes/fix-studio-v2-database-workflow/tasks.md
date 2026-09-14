# 實作待辦

本次僅起草，所有產品實作待辦維持未完成。先閱讀 evidence.md 與 validation.md，重新核對 main；完成正式文件檢查後才開始。以下 [after: ...] 是必須先完成的工作；有共用檔案的項目不得同時修改。

每項行為修正都先加入能重現問題的測試，確認修正前失敗，再修改並確認通過。不得直接把上一輪下載包整包當成已驗證程式套入；只能比較後逐項採用。測試名稱以下是預定驗收標籤，不宣稱目前已存在。

## 1. 先停止錯誤成功與操作互相干擾

- [ ] 1.1 針對記錄方案 `SchemaApply`／`TestWrite` 補上 `NoExecutionNoSuccess` 路由測試，未接線路徑回安全 501 且無成功資料；前端辨識未開放且保留設定，一般 `GenerateSchema` 的既有可用路徑不受影響。對應 recording-database-setup 的能力與真實試寫需求。
- [ ] 1.2 在 `TargetMappingTable` 補 `PreserveEditingDraft`、`SingleEnterCommit`、`CompositionDoesNotSubmit`，實作每列草稿／同列衝突／刪除清理及單一送出途徑。驗證 A 列回覆不抹掉 B 列，中文組字與 Enter 加 blur 不重複送出；單元測試與型別檢查通過。
- [ ] 1.3 [after: 1.2] 在 `SchemaSetupSection` 與 `Step4Database` 補 `StalePreviewIgnored`、`ReadonlyBlocksSchemaActions`，統一範圍／版本檢查與進行中鎖定。驗證切表、切設備、返回頁面和啟用中不能使用舊結果或重複建表；不將忽略回覆當成取消後端。

## 2. 把畫面接回真正的設定

- [ ] 2.1 [after: 1.1] 以 `SavedConnectorIdentity` 測試讓記錄方案從已存連線編號取得真正種類與憑證，前後端驗證所屬範圍；移除寫死目標。測試兩個不同種類的連線、未儲存連線、身分改變、重用遮蔽密碼、原密碼含前後空白、明確清除及無權使用連線。
- [ ] 2.2 [after: 2.1] 以 `MultiDeviceMembershipAndPlanSelection` 測試讓記錄方案使用已儲存的測量項目與各自設備，選取正確方案而非清單第一筆。測試兩台設備、缺少測量項目、停用／移除項目、跨工作區項目、空清單與讀取失敗；不得補造編號。讀取、建立與更新都核對歸屬。
- [ ] 2.3 [after: 2.2, 1.3] 沿用資料表讀取服務接入 `RealMetadataAndReviewableAssignment`：真正欄位、無表、無權限、失敗與新表建議分開，切換身分後重新查詢；保留相容配對，歧義或欄位不足時待確認而非循環分配。測試合法同組共用、跨組衝突及缺少唯一識別的更新策略仍受既有規則控制。
- [ ] 2.4 [after: 2.3] 在本地連線／分組及對應／工作區參照儲存加入 `AtomicSetupSave` 故障注入測試，重現第二段失敗後，改為同一交易與版本檢查。驗證失敗保留原狀、舊版本拒絕、重載一致及未存成功不取得可執行資格；不得把外部建表混入本地原子性承諾。

## 3. 接上可核對的建表與試寫

- [ ] 3.1 [after: 2.4, 1.3] 以 `PreviewBoundToCurrentSetup` 建立後端持久保存的預覽範圍、版本、期限與操作狀態，來源取真實連線／方案；相容讀取舊設定，舊預覽缺少保護欄位則要求重建。測試變更方案、目標或設定後過期，預覽不建表、不試寫；補受影響 API 型別與文件。
- [ ] 3.2 [after: 3.1, 1.1] 以 `ApplyExactlyCurrentPreview` 將記錄方案建表接入真正目標；原子取得執行權，拒絕未知、跨範圍、過期、重複及舊版確認資料。測試成功後真表存在、雙擊／回覆遺失不重做、失敗可核對部分結果，並驗證一般建表路徑的版本保護；未通過能力驗證者維持安全未開放。
- [ ] 3.3 [after: 3.2] 以 `IdempotentRealTestWrite` 實作需確認的真實試寫與持久操作身分，寫到所選連線並能辨識自己產生的測試列。驗證指定目標才有該列、同要求重試／程序重啟不重複、不確定結果不盲目插入；對外試寫開放仍須等待 3.4 完成。
- [ ] 3.4 [after: 3.3] 以 `ReadbackAndCleanupEvidence` 比較回讀內容並只清理本次測試列，成套更新 Go／TypeScript／UI 的已寫入已驗證、已寫入未驗證、失敗、不確定及獨立清理狀態。驗證缺讀權、缺刪權、逾時、清理後重試及正式資料未受影響；成功判斷不再使用另一套字串或自行填時間。

## 4. 收斂畫面與完成驗收

- [ ] 4.1 [after: 3.4, 2.3, 1.2, 1.3] 以 `GuidedDatabaseSetup` 將第 4 步排成四個問題，系統安排記錄表與手動配表分開；加入搜尋／問題篩選、明確批量範圍、逐列儲存狀態、讀取失敗重試與 en／zh-TW 白話文案。測試無設備、無連線、無表、無方案及唯讀；390／768／1440 CSS px 下鍵盤可操作、表格獨立捲動且主要按鈕可到達。
- [ ] 4.2 [after: 4.1] 以 `BackendBackedActivationResults` 核對第 4 步進度與結果只依後端事實顯示，移除與固定動畫成功規格相關的依賴。測試部分設備成功、全部失敗、原有設備仍運作及重載；保留現有工作區／設定版本、readiness token 與 Modbus Share 阻擋條件，不重寫 runtime。
- [ ] 4.3 [after: 4.2] 將上述場景串成整合與瀏覽器驗收，於可丟棄資料庫確認建表、試寫、讀回、清理及回覆遺失；執行 AGENTS.md 規定的前後端測試／檢查／建置、`make check-lines` 與 `git diff --check`。同步實際 API 文件、必要的 Studio inventory 及其 SQLite changelog，保存畫面與資料庫證據；未執行項目明列，重新 analyze／validate，不部署。

## 驗收對照

| 需求 | 主要待辦 |
| --- | --- |
| Capability-backed database preparation | 1.1、2.3、3.2 |
| Revision-bound schema preview and explicit creation | 1.3、3.1、3.2 |
| Credential and path safety | 2.1、3.1、3.2 |
| Truthful explicit test writes | 1.1、3.3、3.4 |
| Persisted recording membership and explicit plan selection | 2.2 |
| Consistent database setup persistence | 2.4 |
| Connector configuration form / Tag-to-column auto-assignment | 2.1、2.3、4.1 |
| Column conflict detection | 2.3 |
| Commit sequence and animation | 4.2 |
| Protected row editing and single submission | 1.2 |
| Fresh schema results and comprehensive readonly controls | 1.3 |
| Guided database setup and accessible controls | 4.1、4.3 |
