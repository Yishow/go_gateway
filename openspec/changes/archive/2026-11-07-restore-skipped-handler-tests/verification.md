# 修復驗證與審查

## 範圍與證據

基準為 9e5fa7cb。此 change 包含 proposal Impact 的十個 Go 檔案、skip 核對文件與本 change 工件。同期 lint change 另行審查；Spectra scope 的 attribution 為 approximated_base，不視為逐 hunk 擁有權證明。

## 規格與實作核對

「輪詢群組名稱唯一性保持一致」由 MemoryRepository.Create/Update 的 mutex 內名稱檢查實現。三個 scenario 與具體資料例子皆由 TestMemoryRepositoryRejectsDuplicateNames 覆蓋：同名新增拒絕、衝突更新保持原名稱和 interval、保留自身名稱可改 interval。SQL 的 UNIQUE 行為未改變。

- Completeness：7/7 任務已完成，1 個 requirement、3 個 scenarios。
- Correctness：名稱回歸先失敗（重複建立回 201），修復後通過。fixture 型態與實際 HTTP/持久資料均有斷言。
- Coherence：無假成功、未刪環境 guard、未改 production successor CAS 優先順序、未放寬真 corruption 的 500 斷言。

## 實際執行

以下皆為本輪 passed-current，硬體項目例外明列：

- go test ./internal/api/handlers ./internal/datalink/pollinggroup -count=1：兩個 package 通過。
- go test -race 的四個 handler/name 回歸：通過；補強 interval 斷言後再跑 memory 回歸通過。
- go test ./internal/datalink -run TestRealDevice -count=1 -v：opt-in 回歸通過；真設備測試預設 SKIP，未執行硬體連線。
- PostgreSQL TestPostgresSuccessorActivationRejectsForeignOrMissingProof：最終 11 個事件通過，0 fail/skip。
- focused golangci-lint：0 issues、exit 0；有預設 goimports cache 目錄唯讀警告，整體最終 lint 由整合流程另跑可寫 cache。
- focused go vet：exit 0。
- 混合規則 MC wire 與不可連線 PostgreSQL 目的端吞吐測量：皆通過，詳細數據見 skip 核對文件。
- Spectra validate/analyze：通過，analyze 0 findings；diff check、行數檢查通過。既有超長設備測試由 725 減為 631 行。

## Code review

依 correctness、efficiency、simplification/reuse、convention 四個方向檢查十個 Go diff，未找到阻擋缺陷。獨立 reviewer 另覆核七個本機修復檔案，確認 loopback 綁定、cleanup 順序、opt-in 位於 Dial 前、同鎖更新與 SQL 約束一致。

三個 dialect proof fixture 保留操作前後內部/外部狀態不變斷言，以真正缺失 generation 驗證 corruption，stale revision 獨立要求 409，沒有放寬為任意錯誤。

## 未驗證與限制

- MySQL/SQL Server 的 proof fixture 在本 change 驗證時尚未完成真實 dialect 執行；其結果由資料庫補測獨立記錄，不由 PostgreSQL 通過推論。
- PostgreSQL 其他整合測試、TimescaleDB 與所有其他 DB 的完整結果屬本次獨立補測範圍。
- 18 組長時間負載矩陣、實體 PLC、其他平台驗收未執行。
- skip 不是 pass；本 change 不宣稱 524 個缺漏全部已補齊。最終完整檢查與發布由整合流程執行。
