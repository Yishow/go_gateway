## Why

執行日誌目前主要是程序啟停、HTTP 與沒有群組識別的泛用錯誤，無法回答哪個設備、規則或寫入群組出問題及如何修復。需要保留安全且有限的程序日誌設計，同時加入來自真正執行路徑的領域事件，而不是逐筆 sample 刷屏或把 log 當成 SQL receipt。

## What Changes

- 加入設備／來源規則／mapping／寫入群組的安全事件，攜帶 opaque identity、實際 revision、階段、typed code、可行動建議及證據參照。
- 事件涵蓋設定套用結果、採集／轉換故障及恢復、群組 intake／交付故障與恢復、容量拒收、storage repair／cutover 結果；正常 SQL 成功採有界摘要，不逐 sample 記錄。
- 執行日誌提供 entity／stage／需處理篩選，展開事件能辨識影響、時間、重複次數及精準設定／runtime 規則連結，保留原有 level/source/search/pause/resume。
- 相同故障有界合併，保留首次／最近時間與 occurrence；恢復事件獨立，logging 不阻塞採集或 SQL。
- 保持 snapshot／replay／live、loss counters、字串 cursor、容量上限、遮蔽先於 sink 以及 loopback-only 門禁。
- 依賴 expose-runtime-recording-rules 的共用讀取與定位契約；此案只擁有事件 producer、白名單投影、合併與 viewer，不新增交付權威。
- 本次只建立提案；不實作、連線私人主機、放寬遠端權限或宣稱產品測試已過。

## Capabilities

### New Capabilities

- `actionable-runtime-events`: 可定位、可修復、可界定歷史與當前故障的安全領域事件，以及有界事件產生與篩選呈現。

### Modified Capabilities

無。既有 active change 中的 runtime-log-viewer 契約保持；本案增補領域內容而非改寫 transport／local access 規則。

## Impact

- Producer 接線：cmd/test_ui/service_wiring.go、internal/datalink/runtime/、internal/datalink/grouppipeline/、internal/datalink/groupdelivery/、internal/api/handlers/；沿 durable transition 發出安全 identity，不傳 raw error。
- 事件／API：internal/diagnostics/、internal/api/runtime_log_query.go、internal/api/runtime_log_handler.go；保留 runtime_log_guard.go 的門禁。
- UI：frontend/src/features/runtime-logs/、frontend/src/types/runtimeLogs.ts、frontend/src/services/runtimeLogs.ts、en／zh-TW locales；使用 expose-runtime-recording-rules 的連結／safe identity。
- 硬依賴 expose-runtime-recording-rules；manage-recording-lifecycle-safely 與 fix-step4-recording-authoring 擁有 lifecycle／authoring 結果，guide-step4-recording-workflow 接收設定定位；本案不得改變這些寫入行為。
- add-windows-tray-runtime-logs 保留 broker／file／SSE／guard／Windows lifecycle 的唯一責任及歷史驗證；此案只新增領域投影，不能重置既有完成項目或 archive 該案。

## 執行順序調整

**Supersedes:** 原先整個 change 必須等待前案全部完成的執行安排，僅改為可隔離準備先行、正式啟用與整合驗收維持原依賴順序。runtime 共用契約固定後，獨立 domain projection／coalescer 可先行；production emitters、sink 接線、viewer locator 與最終驗收仍待 runtime 與 backend/lifecycle 交接。 七個 change 與本案需求／驗收不變；本次只更新 artifacts，不代表產品已實作。詳見 design.md「執行分期與交接」。
