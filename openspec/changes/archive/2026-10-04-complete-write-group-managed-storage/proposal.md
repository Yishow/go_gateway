## Why

新的基本群組 UI 能提出欄位，但尚未把空白 SQLite／PostgreSQL 的受控準備接成完整流程；目前端到端工具先用 SQL 建表。基本落表也沒有保證記錄時間與來源資訊。本案只補齊 canonical WriteGroup 的安全建表及必要欄位。

## What Changes

- 以保存的群組產生受控資料表與既有 receipt 表提案，明確確認後才建立，再查真 schema 驗證。
- 本機 SQLite 與既有 PostgreSQL 連線皆可從無採集表開始，不必先建立 RecordingPlan 或手動 targets。
- 新 managed 表自動包含穩定 record/group/device 身分、UTC bucket 時間及每點 observed_at／品質資訊。
- 沿用 exact codec 與已驗證 receipt 交付，不新增 UI 的 unique_key 模式或通用 schema migration。

## Capabilities

### New Capabilities

無；沿用現有 capability，不建立第二套採集或交付模型。

### Modified Capabilities

- `recording-database-setup`
- `write-group-row-semantics`

## Impact

限定既有 workspace WriteGroup、dbtarget layout/codec、recordingplan 的 schema operation ledger、對應 group handler/route 及 Step 4 schema 準備元件。不得讓基本群組必須建立 advanced RecordingPlan。必要 generated Swagger／SDK／locale 與技術文件隨實作同步。

## Scope and Dependencies

只建立本產品受控的新表，或驗證已有且具相同 ownership 的相容表；不 ALTER、DROP、改名或搬移使用者既有表，不建立 PostgreSQL server/database/account。目的 SQLite 與內部設定/journal DB 必須不同。自訂表維持進階入口及原有合約。

單設備基本路徑前置：`fix-write-group-runtime-lifecycle`, `fix-write-group-delivery-deadlines`。`fix-write-group-entity-row-layout` 是多 entity 的正式驗收前置，不阻擋基本 SQLite UI→首列的早期流程回饋；共用 layout 的修改仍按 ownership 整合。
本案只起草；產品 tasks 全部未完成。共同邊界與來源見 [總覽](../../../docs/plans/studio-v2-flow-completion/README.md)。
