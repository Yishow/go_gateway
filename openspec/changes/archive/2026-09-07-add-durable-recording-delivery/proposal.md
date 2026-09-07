# 資料庫斷線時安全暫存，恢復後不重複補送

## Why

有曲線與用量就不能只靠記憶體 buffer。外部資料庫失聯、Gateway 重啟或回應遺失，都可能漏資料或重複計量。
需要把「讀到」「本機保存」「送到目的地」「查回確認」分開，才能讓使用者知道資料真正在哪裡。

## What Changes

- 新增本機 journal、transactional outbox、每目的地 receipt 與可重入 delivery worker。
- 將樣本、計量 checkpoint 與待送結果的提交邊界明確化；無法持久化時不得宣稱成功。
- 加入重送去重、逐目的地失敗隔離、退避、容量告警、恢復與保留策略。
- 對執行中的計畫版次與目標變更提供安全切換，不將舊佇列改送新目的地。
- 不承諾任意外部系統的 exactly-once；只在具備交易與去重能力的目標提供可驗證單次效果。

## Capabilities

### New Capabilities
- `durable-recording-delivery`: 耐久入口、送出、回執、容量與版次切換。

### Modified Capabilities
- `database-output-delivery`: 將「新資料只取啟用目標」與「已耐久接受的舊資料仍需送達」分清楚；其餘 truth 與 activation gate 不變。

## Impact

依賴 1、2 的 sample/storage 介面；整合驗收包括 3 的 stateful checkpoints。
擬新增本機 durable recording repository/outbox worker，整合現有 runtime ingestion 與 dbtarget adapters。
不重用 in-memory grouped bucket 作為已保存證據，不改現場資料庫內容或連線設定。
