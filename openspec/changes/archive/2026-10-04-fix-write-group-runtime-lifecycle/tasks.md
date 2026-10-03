## 1. 草稿與責任界線

- [x] 1.1 補 RED：執行中改名／改 members 只存草稿、不 Apply，跨兩桶及重啟仍使用舊版本；同語意在退休完成前重新啟用（10 秒桶，t=12 停用、t=15 Apply、t=21 新樣本須入 journal）、有效切換與 legacy writer 排他性都有斷言。
- [x] 1.2 修正執行資格及可恢復的截止界線，讓 1.1 GREEN；同步說明 draft、applied、disabled 與 drain-only 的差異。
- [x] 1.3 補 RED 的 AcceptSample 與 SetUntil 並行測試，再以同一鎖域修正；實際 package 的 `go test -race` 通過，不以隔離範例代替。

## 2. 離線與歷史版本恢復

- [x] 2.1 補 RED：目的 DB 離線時冷啟動仍可本地 ACK；ACK 後停用／刪除／切版、封桶前 kill，重啟可排空原 revision；逐一檢查原 destination 與 effect identity。
- [x] 2.2 保存並使用已驗證的凍結恢復描述，缺描述或版本不符明確 blocked；讓 2.1 的離線新 intake 測試 GREEN，且沒有遠端 metadata 前置。
- [x] 2.3 從 durable journal/checkpoint 發現歷史版本，恢復只排空 boundary；讓 2.1 的三種歷史狀態 GREEN，無新 intake、重複 row 或未完成 journal 遺留。
- [x] 2.4 補 migration 重跑、既有缺描述版本、connector 身分改變的回歸；同步回復文件，不修改已接受 payload 或移轉目的地。

## 3. 整合驗證

- [x] 3.1 執行受影響套件、真 SQLite／PostgreSQL 恢復整合與 race 測試；再跑 repository 後端最低檢查，保存命令、版本和未執行範圍。
- [x] 3.2 對照 Safe basic group lifecycle 與 Production durable intake and recoverable closure 的兩份 delta 的全部 scenarios 與 production wiring，確認沒有只修 domain 測試；review 修復後重跑受影響測試並交付同一範圍。
