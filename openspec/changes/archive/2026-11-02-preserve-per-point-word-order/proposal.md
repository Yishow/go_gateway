## Why

同一 PLC 的不同點位可能來自不同設備，需要各自解碼 ABCD／CDAB。現有點位雖保存 data_format，來源規則更新仍會覆寫衍生點位，部分協議 runtime 又僅取連線預設，不能靠畫面欄位宣稱需求完成。

## What Changes

- 新增有來源可追溯的逐點格式意圖：繼承、明確覆寫、舊版相容；提供有效格式及來源讀取契約。
- 規則更新／重新衍生保留同一點位覆寫，重設才回復繼承；以範圍、版本與冪等控制避免跨設備或舊畫面覆蓋。
- 將現有 Modbus TCP／RTU／UDP、MC3E、FATEK 生產解碼接線與適用型別矩陣納入驗收，不新增協議。
- 舊資料缺少意圖標記時保存原協議有效解碼語意，不猜測使用者意圖、不靜默改值。
- 作為 streamline-recording-task-editor 的獨立前置能力；本案提供 API／runtime 契約，後者負責簡易點位編輯畫面。

## Capabilities

### New Capabilities

- `per-point-word-order`: 點位格式意圖、繼承、覆寫保存、協議實際解碼與安全相容。

### Modified Capabilities

無；新增契約補充 point-catalog 與 source-rule-runtime，保留既有 lifecycle／candidate 安全要求。

## Impact

- internal/datalink/schema、internal/datalink/point、internal/datalink/sourcerule：內部設定持久化、穩定 identity、規則同步及 rollback。
- internal/datalink/collector、internal/datalink/connector/adapters：排程／手動讀取的格式傳遞與協議 decoder。
- internal/api/handlers、cmd/test_ui：既有 scoped API 與 production wiring 整合。
- frontend/src/types、frontend/src/services：供既有與簡易編輯器共用的契約；本案不改頁面布局。
- 不新增環境或套件；產品實作與測試尚未執行。
