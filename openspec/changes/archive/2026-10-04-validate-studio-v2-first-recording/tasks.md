## 1. 空白目的地正式驗收

- [x] 1.1 補 harness 前置與負測試：缺檔 SQLite 用 stat、已有檔只用 read-only URI；每 run 用 fresh file/schema，證明目的地無採集表，setup 不可用直接 API／fixture SQL 偷建表或寫值；沿用既有 loopback、命名空間與 cleanup 保護。
- [x] 1.2 驗證 Production device-to-SQL witness：先以 D/E 的單設備空白 SQLite 真 UI→首列取得產品回饋，再真 UI 完成 SQLite／PostgreSQL 空白設定、schema 確認與開始記錄，獨立 SQL 比對 identities/時間/品質及至少三個正式 bucket；先記失敗，再隨 A-E 修正轉綠。
- [x] 1.3 執行七型別含精確 uint64 2^53+1／2^63／最大值 的單/雙設備路徑；同群組多 entity 另外驗證，不用兩群組結果代替；更新 witness 說明。

## 2. 恢復與互動矩陣

- [x] 2.1 驗證 Review-regression witnesses cover actual recovery combinations：執行草稿跨桶/重啟、慢 sender、離線冷啟動新 ACK、ACK 後停用/刪除/切版恢復及並行 cutoff；SQL 與本地 durable state 均符合 A-C。
- [x] 2.2 執行補送時間、schema/start 的 duplicate/lost response/stale/partial、custom/Share-only 保留，並重跑既有品質/容量/poison/worker/未知 commit/cleanup 保護。
- [x] 2.3 對兩個 DB 各量測三次預設 60 秒區間的空白 setup 到首筆 SQL，用獨立 fresh namespace 保留每次、最慢值及 setup/system wait/delivery 的 monotonic 分段與 polling 限制；操作者流程觀察另記，對照 300 秒目標；不縮減 correctness 驗收以過關。

## 3. 完成與證據

- [x] 3.1 保存 source/build、命令、環境、真 UI/SQL、durable identities 及 owned cleanup 的 sanitized witness；未執行平台、現場與工具限制明列 NOT RUN，不改寫歷史證據。
- [x] 3.2 核對 A-E 所有 scenarios、focused tests、repository 全套最低檢查與本矩陣皆完成；早期 SQLite 首列結果不能代替本項完成；review 只回修唯一 owner 的本範圍問題，不加新能力湊驗收。
