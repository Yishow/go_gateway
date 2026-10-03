# 實作待辦與移交

2026-10-01：依新一輪V2基本群組寫入計畫，五項未完成工作移交給唯一owner；不是完成、取消或archive。原基準為main 4c00af81fb51b3239b4e0d0d9c6b83680c808b36，原5/10任務全文可由该版讀取。詳見 [handoff.md](handoff.md)。

目前本案保留五項既有完成紀錄，另有一項未完成的移交closeout gate。數量改變是工作歸屬整理，不能當成產品完成率提高。移交的2.3、3.3、3.4、4.1、4.3不得再於本案平行實作，後續以新owner的tasks及證據為準。

## 1. 保留已完成實作與原始描述

- [x] 1.2 確認 `fix-studio-v2-database-result-truthfulness` 1.1–1.8 已驗證後，在 `TargetMappingTable` 補 `PreserveEditingDraft`、`SingleEnterCommit`、`CompositionDoesNotSubmit`，實作每列草稿／同列衝突／刪除清理及單一送出途徑。驗證 A 列回覆不抹掉 B 列，中文組字與 Enter 加 blur 不重複送出；單元測試與型別檢查通過。 對照：Protected row editing and single submission；設計「讓輸入與預覽不互相搗亂」。
- [x] 1.3 [after: 1.2] 在 `SchemaSetupSection` 與 `Step4Database` 補 `StalePreviewIgnored`、`ReadonlyBlocksSchemaActions`，統一範圍／版本檢查與進行中鎖定。驗證切表、切設備、返回頁面和啟用中不能使用舊結果或重複建表；不將忽略回覆當成取消後端。 對照：Fresh schema results and comprehensive readonly controls。
- [x] 2.1 確認 `fix-studio-v2-database-result-truthfulness` 1.1–1.8 已驗證後，以 `SavedConnectorIdentity` 測試讓記錄方案從已存連線編號取得真正種類與憑證，前後端驗證所屬範圍；由 `schema_dbtarget_models.go`、dbtarget repository 及必要本地 migration 保存 connector identity revision，使用 server target resolver，身份／憑證變更遞增 revision，健康檢查不代替 identity revision；移除寫死目標。測試兩個不同種類的連線、未儲存／foreign／unknown connector、身分改變、dialect spoof、重用遮蔽密碼、原密碼含前後空白、明確清除及無權使用連線。 對照：Credential and path safety；Connector configuration form；設計「遵循已完成的安全封鎖」「目標只能來自已儲存資料」。
- [x] 2.2 [after: 2.1] 以 `MultiDeviceMembershipAndPlanSelection` 測試讓記錄方案使用已儲存的測量項目與各自設備，選取正確方案而非清單第一筆。測試兩台設備、缺少測量項目、停用／移除項目、跨工作區項目、空清單與讀取失敗；不得補造編號。List／Get／Create／Update／Delete 都核對歸屬，缺少或跨工作區資源以同一安全 404 拒絕且不讀改其他範圍；測試工作區 A 對工作區 B 的查詢、更新與刪除。 對照：Persisted recording membership and explicit plan selection。
- [x] 2.4 [after: 2.2, 1.3] 在本地連線／分組及對應／工作區參照儲存加入 `AtomicSetupSave` 故障注入測試，重現第二段失敗後，改為同一交易與版本檢查，將 connector identity revision 與設定 revision 同交易保存。驗證失敗保留原狀、舊版本拒絕、重載一致及未存成功不取得可執行資格；不得把外部建表混入本地原子性承諾。 對照：Consistent database setup persistence；設計「一次儲存要有一致的結果」。

## 5. 移交後核對

- [x] 5.1 [after: 2.4] 確認 handoff.md 指定的D/E/F及其A/B/C前置已完成實作與驗收，逐項對照原2.3/3.3/3.4/4.1/4.3需求及保留的schema/readiness/Share回歸；核對移轉delta與既有完成證據一致後再進行本案verify/review及archive評估，不能只因本次文件合併勾選。

5.1 移交核對對照（不改寫以上五項歷史完成描述）：design「真實欄位與建表建議分開」與「第 4 步用四個問題引導」由 E 完成；「消費已驗證的建表結果」「試寫、讀回與清理各自有證據」由 D 實作、F 以真 UI／SQL 與故障矩陣驗收。逐項證據與移轉／回復限制見 [handoff-closeout.md](../../../../docs/plans/studio-v2-write-groups/handoff-closeout.md)。2026-10-03 最後 review／verify、完整基準與受影響實跑齊全，依實際證據完成 closeout。
