## 1. 群組導向的受控 schema

- [x] 1.1 補 RED（Canonical group managed preparation）：只有 canonical group、無 RecordingPlan／manual targets 時，可對空白 SQLite／PostgreSQL 提案並確認建表；preview 不產生目的副作用。
- [x] 1.2 用既有 schema service／ledger 接 group-scoped adapter 與薄 API，讓 1.1 GREEN；同步 DTO、generated Swagger 與 operation 文件，不另建 token 或 ledger。
- [x] 1.3 驗證 Managed preparation preserves operation and table ownership、Revision-bound schema preview and explicit creation：補失去回覆、duplicate confirm、restart、stale/foreign scope、無 CREATE 權限與同名未知 ownership 的回歸；修到拒絕前無副作用，未知結果可查。

## 2. 可用且可恢復的 SQL 資料

- [x] 2.1 補 RED（Basic managed SQL preserves record time and origin identity）：managed 真 SQL 必須有凍結 record/group/device identity、UTC bucket 時間與每點 observed_at／品質，延遲補送不改時間。
- [x] 2.2 落實 Stable exact managed column generation：接上系統 metadata 欄位與 exact typed layout，保存到 A 的恢復描述，讓 2.1 GREEN；驗證七個既有型別、uint64 2^53+1／2^63／最大值精確、UI 不提案 unsafe BIGINT、dialect compatibility 與欄名衝突。
- [x] 2.3 共用既有 receipt 策略，真 SQL 驗證 record/receipt 同交易及 commit acknowledgement loss 不重複；不新增 unique_key 產品模式。

## 3. UI 準備與相容保護

- [x] 3.1 補 saved group B 與 connector table A 的 read-only metadata、stale/foreign scope 及未保存表名回歸，再把 proposal-only UI 接成預覽→明確確認→真實結果查詢，SQLite 亦可用；先補元件回歸，再同步 en/zh-TW 與安全修復文案。
- [x] 3.2 驗證目的 SQLite 不可等於內部 DB、existing custom 表不 ALTER、不改 payload、相容 owned 表才 no-op；更新 rollout／rollback 限制。

## 4. 整合驗證

- [x] 4.1 執行真 SQLite／PostgreSQL schema/codec 整合及受影響前後端最低檢查；對照全部 scenarios，review 後重跑修正範圍。
