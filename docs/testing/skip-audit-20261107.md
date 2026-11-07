# Go 測試跳過項目核對（2026-11-07）

## 原始 524 個測試事件

以下取自先前完整 go test JSON 的 test-level skip；不計沒有測試檔的 package，也不把 skip 當成 pass。

| 分類 | 數量 | 判斷 |
| --- | ---: | --- |
| PostgreSQL fixture 未配置 | 294 | 尚未驗證資料庫整合，不是測試成功 |
| MySQL/MariaDB fixture 或驗證 TLS 未配置 | 192 | 尚未驗證資料庫整合，不是測試成功 |
| SQL Server fixture 或明確操作 gate 未設定 | 26 | 尚未驗證資料庫整合，不應繞過 gate |
| 無條件停用的 handler 測試 | 3 | 本次修復並恢復執行 |
| 明確選入的耗時／模擬器測量 | 3 | 一般測試預設跳過；本次補測其中兩項 |
| 僅供父測試啟動的子程序入口 | 3 | 正常；不能在主測試程序直接執行 helper |
| 實體 PLC 不可連線 | 1 | 不屬模擬器驗證；本次改為明確 opt-in |
| 嵌入式 SPA 未建置 | 1 | 建置並同步 static 後補測 |
| SQLite 不支援只差大小寫的名稱 | 1 | 正常 dialect 限制，不應假裝可測 |
| 合計 | 524 | skip 數不是失敗數，也不是成功數 |

## 本次修復

- 設備 activate 與 activate-then-disable 測試改用動態埠的 loopback Modbus peer，驗證 API 回應、active/disabled 持久狀態與成功的 probe。
- 輪詢群組重複名稱空測試改為實際建立兩次並驗證拒絕。記憶體儲存庫現在遵守 SQL 已有的 UNIQUE 契約，拒絕衝突更新且保留原欄位；自身名稱可正常保留。
- 真 PLC 測試只有 GATEWAY_REAL_DEVICE_INTEGRATION=1 才會進入網路操作。執行前仍需設備目標個別授權，一般完整測試不會自動探測硬體。
- 資料庫 proof fixture 必須區分 stale revision（409）與缺少 generation 的 lineage corruption（500），不能降低斷言來消除失敗。

## 已完成的本機針對性驗證

- handlers 與 pollinggroup 兩個 package 全套通過。
- 三個恢復測試與 memory uniqueness 回歸通過；相同選測的 race 通過。
- 硬體 opt-in 條件回歸通過；真設備項目在預設執行時仍明列 SKIP。
- GATEWAY_MIXED_SIMULATOR=1 下的 TestProductionScheduledAcquisitionMixedRuleWireRequests 通過：38 個地址以一個 MC random request 讀取。
- GATEWAY_TASK13_THROUGHPUT=1 下的 TestProductionSnapshotWorker_MCUnreachableThroughput 通過：22 秒觀測 704 點，約 32 點/秒；不可連線的目的端確實失敗，改回可丟棄 SQLite 後恢復寫入。
- PostgreSQL successor proof 最終選測 11 個測試事件通過、0 fail/skip；MySQL/SQL Server 同源 fixture 修正已靜態審查，實際 dialect 測試不由此結果推論通過。
- TestAcquisitionAcceptanceScheduledLoad 的完整 18 組耗時矩陣本輪未執行，不以兩個其他 opt-in 測試取代。

資料庫補測與最後完整批次結果以本次交付報告為準，不能沿用原始 524 的數字當成修復後結論。

## 後續：依要求移除剩餘八個測試身份

此節取代上文對測試入口「目前仍保留」的描述，不改寫先前測試結果。

- 刪除真 PLC 整合測試檔及其 opt-in helper，測試套件不再含該硬編碼目標。
- 刪除 TestAcquisitionAcceptanceScheduledLoad 的 18 組長負載矩陣入口；兩個短排程 pilot 與共用測量 helper 保留。
- SQLite quoted-distinct 組合不再註冊；PostgreSQL 對應組合及 SQLite 其餘案例保留。
- 三個 crash/kill 子程序改為普通 helper，由既有父測試根據專用環境參數分派，保留終止程序及恢復驗證，沒有把空執行改報 PASS。
- SQL Server gate 的拒絕分支改用 recorder 核對拒絕與原因，不再刻意建立兩個 SKIP 子測試；整合測試的能力檢查不變。

移除不代表真 PLC 或長負載已通過。本次未修改產品程式，也未取消其他尚缺環境的 skip。Mac 手測仍須依一般建置流程建立前端與內嵌資源。
