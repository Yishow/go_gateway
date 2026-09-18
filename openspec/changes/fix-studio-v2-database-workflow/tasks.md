# 實作待辦

所有產品實作待辦維持未完成。先讀 evidence.md、validation.md 及 design.md 的 Implementation Contract。1.2／2.1 開始前須完成 `fix-studio-v2-database-result-truthfulness` 1.1–1.8；先執行本案 1.2、1.3、2.1、2.2、2.4，再切到 `implement-studio-v2-verified-schema-setup` 3.1–3.10，完成後才回來做 2.3 及後續。跨案 gate 必須由主代理核對待辦與測試證據，不以 CLI 只驗證本地依賴當成已解鎖。

每項行為修正都先加入能重現問題的測試，確認修正前失敗，再修改並確認通過。不得直接把上一輪下載包整包當成已驗證程式套入；只能比較後逐項採用。測試名稱以下是預定驗收標籤，不宣稱目前已存在。

## 1. 先停止錯誤成功與操作互相干擾

- [x] 1.2 確認 `fix-studio-v2-database-result-truthfulness` 1.1–1.8 已驗證後，在 `TargetMappingTable` 補 `PreserveEditingDraft`、`SingleEnterCommit`、`CompositionDoesNotSubmit`，實作每列草稿／同列衝突／刪除清理及單一送出途徑。驗證 A 列回覆不抹掉 B 列，中文組字與 Enter 加 blur 不重複送出；單元測試與型別檢查通過。 對照：Protected row editing and single submission；設計「讓輸入與預覽不互相搗亂」。
- [x] 1.3 [after: 1.2] 在 `SchemaSetupSection` 與 `Step4Database` 補 `StalePreviewIgnored`、`ReadonlyBlocksSchemaActions`，統一範圍／版本檢查與進行中鎖定。驗證切表、切設備、返回頁面和啟用中不能使用舊結果或重複建表；不將忽略回覆當成取消後端。 對照：Fresh schema results and comprehensive readonly controls。

## 2. 把畫面接回真正的設定

- [x] 2.1 確認 `fix-studio-v2-database-result-truthfulness` 1.1–1.8 已驗證後，以 `SavedConnectorIdentity` 測試讓記錄方案從已存連線編號取得真正種類與憑證，前後端驗證所屬範圍；由 `schema_dbtarget_models.go`、dbtarget repository 及必要本地 migration 保存 connector identity revision，使用 server target resolver，身份／憑證變更遞增 revision，健康檢查不代替 identity revision；移除寫死目標。測試兩個不同種類的連線、未儲存／foreign／unknown connector、身分改變、dialect spoof、重用遮蔽密碼、原密碼含前後空白、明確清除及無權使用連線。 對照：Credential and path safety；Connector configuration form；設計「遵循已完成的安全封鎖」「目標只能來自已儲存資料」。
- [x] 2.2 [after: 2.1] 以 `MultiDeviceMembershipAndPlanSelection` 測試讓記錄方案使用已儲存的測量項目與各自設備，選取正確方案而非清單第一筆。測試兩台設備、缺少測量項目、停用／移除項目、跨工作區項目、空清單與讀取失敗；不得補造編號。List／Get／Create／Update／Delete 都核對歸屬，缺少或跨工作區資源以同一安全 404 拒絕且不讀改其他範圍；測試工作區 A 對工作區 B 的查詢、更新與刪除。 對照：Persisted recording membership and explicit plan selection。
- [x] 2.4 [after: 2.2, 1.3] 在本地連線／分組及對應／工作區參照儲存加入 `AtomicSetupSave` 故障注入測試，重現第二段失敗後，改為同一交易與版本檢查，將 connector identity revision 與設定 revision 同交易保存。驗證失敗保留原狀、舊版本拒絕、重載一致及未存成功不取得可執行資格；不得把外部建表混入本地原子性承諾。 對照：Consistent database setup persistence；設計「一次儲存要有一致的結果」。
- [ ] 2.3 [after: 2.4] 確認 `implement-studio-v2-verified-schema-setup` 3.1–3.10 已驗證後，消費其 metadata 契約接入 `RealMetadataAndReviewableAssignment`：真正欄位、無表、無權限、失敗與新表建議分開，切換身分後重新查詢；保留相容配對，歧義或欄位不足時待確認而非循環分配。測試合法同組共用、跨組衝突及缺少唯一識別的更新策略仍受既有規則控制。 對照：Tag-to-column auto-assignment；Column conflict detection；設計「真實欄位與建表建議分開」「消費已驗證的建表結果」。
  - 狀態註記（2026-09-16）：metadata 消費的基礎已隨 3.1–3.10 提前落地——`useStep4TargetColumns`／`TargetMetadataStatus` 依 scope 讀取真實欄位並區分 not_checked／checking／exists／missing／forbidden／failed，`Step4Database` 的 `autoAssignTargets` 保留既有配對並以衝突標記欄位不足。本項的完整驗收（`RealMetadataAndReviewableAssignment` 場景與合法同組共用／跨組衝突／缺少唯一識別的測試證據）仍待完成，故維持未勾選。

## 3. 接上可核對的建表與試寫

- [ ] 3.3 [after: 2.3] 以 `IdempotentRealTestWrite` 沿用 C3 的 preview／operation repository，擴充 `test_write` 種類、test-write-preview route、既有試寫 handler、operation 查詢與前端服務／hook，實作 design 第 6 節的 token／operation_id／版本及 202／200／409 契約；測預覽零 target mutation、建表 token 不可試寫、foreign status 安全 404。真實試寫只到所選連線並能辨識自己的測試列；驗證同要求重試／程序重啟不重複、舊版僅送 plan_id 被拒、不確定結果只查原操作；對外試寫開放仍須等待 3.4 完成。 對照：Truthful explicit test writes；設計「試寫、讀回與清理各自有證據」。
- [ ] 3.4 [after: 3.3] 以 `ReadbackAndCleanupEvidence` 比較回讀內容並只清理本次測試列，成套更新 Go／TypeScript／UI 的已寫入已驗證、已寫入未驗證、失敗、不確定及獨立清理狀態。驗證缺讀權、缺刪權、逾時、清理後重試及正式資料未受影響；成功判斷不再使用另一套字串或自行填時間。 對照：Truthful explicit test writes；設計「試寫、讀回與清理各自有證據」。

## 4. 收斂畫面與完成驗收

- [ ] 4.1 [after: 3.4, 2.3, 1.2, 1.3] 以 `GuidedDatabaseSetup` 將第 4 步排成四個問題，系統安排記錄表與手動配表分開；加入搜尋／問題篩選、明確批量範圍、逐列儲存狀態、讀取失敗重試與 en／zh-TW 白話文案。測試無設備、無連線、無表、無方案及唯讀；390／768／1440 CSS px 下鍵盤可操作、表格獨立捲動且主要按鈕可到達。 對照：Guided database setup and accessible controls；設計「第 4 步用四個問題引導」。
- [ ] 4.3 [after: 4.1] 將上述場景串成整合與瀏覽器驗收，回歸 `fix-studio-v2-database-result-truthfulness` 的 BackendBackedActivationResults（部分成功、全部失敗、原有設備、重載），於可丟棄資料庫確認建表、試寫、讀回、清理及回覆遺失；執行 AGENTS.md 規定的前後端測試／檢查／建置、`make check-lines` 與 `git diff --check`。同步實際 API 文件、必要的 Studio inventory 及其 SQLite changelog，保存畫面與資料庫證據；未執行項目明列，重新 analyze／validate，不部署。

## 移交紀錄

初稿 1.1、3.1、3.2、4.2 已移交其他 change，映射與理由見 design.md 的 Supersedes；不是已完成工作。保留其餘編號以利追溯。
