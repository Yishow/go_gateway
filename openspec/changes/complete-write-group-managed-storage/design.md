## Context

主線已有受確認保護的 schema operation ledger，但掛在 recording-plan 路徑；GroupColumnProposal 只展示提案，SchemaSetupSection 對 SQLite 不顯示。production boundary 未把 BucketStartColumn／RecordKeyColumn 接入基本群組。詳見 proposal.md 與總覽來源。

## Goals / Non-Goals

讓空白目標透過正式 UI 完成準備，保留既有 schema token、operation identity、CAS、readback 與 ownership 保護。不建立第二套 plan、token、ledger、SQL editor 或 schema migration engine。

## Decisions

1. 在現有 schema service 加 canonical group 輸入 adapter，以保存的 group/member/connector revisions 生成 layout；薄 group-scoped preview/confirm 入口共用既有 ledger、operation lookup 及錯誤契約。舊 recording-plan API 保留，不偽造 measurement/plan。
2. 新 managed 群組採現有寬列模型，一個群組一個受控表；命名由穩定群組與 member identity 產生並保存，不由畫面列順序或 display name 反覆生成。每台設備一個預設群組的建立與重用由後續操作案負責。
3. 新基本單設備 managed 表具有系統控制的 record_id、group_id、device_id、bucket_start 與 provenance 欄位；值來自該筆凍結 row/member，不取 flush 時間或目前 UI 選擇。每點 observed_at、quality、reason 沿用既有 provenance。既有跨設備 advanced 表不能被填上猜測的單一 device_id，仍遵守原本逐 member/列身分契約。欄名避免與使用者 Tag 衝突，具體 dialect 型別以 exact codec 的安全策略決定，uint64 全範圍不得經浮點。
4. 管理表用 record identity 約束加既有 receipt 策略，row 與 receipt 同一遠端交易；本案不提供新的 unique_key 選項。把必要欄位與 metadata bindings 保存到 A 的恢復描述，送出與重啟均用同一版本。
5. preview 只讀，不先建立 SQLite 檔案。確認後可在 backend 管理的資料目錄建立獨立目的檔；已存在檔案必須驗證身分，禁止指向 gateway 設定/journal DB。PostgreSQL 僅使用已存在資料庫及允許的 schema，不替使用者建帳號或升權。
6. 群組是 canonical table configuration：既有 token 的 PlanID 映射 group.ID、PlanRevision 映射 group.Revision，不建立假的 RecordingPlan。群組 adapter 在同一 token/ledger 保存並檢查 destination.schema_revision/schema_digest（未 Apply 時可為空）、凍結 members/source revisions 及 layout digest；不能把空值當作已確認 schema。每次確認綁定 workspace/group/source/connector revisions、目標範圍、statement/layout digest、expiry 及 operation_id；不執行 client SQL。相同 operation 重送／重啟先查既有證據；失去回覆保持 unknown，不能猜成功或盲目 DDL。
7. 同名未知 ownership 或不相容表必須阻擋並保留原資料；既有 owned 且相容的表可真實 no-op。需要不同結構時明確重新提案建立新受控表或選進階路徑，不能自動 ALTER 舊表。既有 custom/all-good 表不被強迫增加 metadata 或 receipt。

## Implementation Contract

- 行為：保存的 managed group 即可預覽→明確確認→真 schema inspection，不要求 legacy targets 或 RecordingPlan。SQLite 缺檔 preview 無副作用，PostgreSQL 不建資料庫/帳號。
- 介面：薄 group-scoped schema API 共用 recordingplan schema service、preview token、operation lookup。PlanID/PlanRevision 的 group 映射如決策 6；群組/schema/source/connector 任一版本變更，確認前回 stale 並要求重新預覽。
- 欄位查詢：database-metadata 可帶 group_id、expected_group_revision、expected_connector_revision；server 以目前 workspace 的持久 group 解析 table/schema。未知、deleted、foreign 或 stale 在 inspection 前拒絕，不接受任意 client table/DSN。未保存的 table 變更不冒充已查證欄位。
- 型別：uint64 完整範圍採 SQLite TEXT／PostgreSQL NUMERIC(20,0) 或 codec 已證明等價的精確表示；BIGINT 不足以容納完整範圍。9007199254740993、9223372036854775808、18446744073709551615 不經 JS Number/浮點；proposal、metadata compatibility 與 production codec 一致。
- 失敗：unknown ownership、不相容或無權限阻擋且保留原表；same operation 重送先查證，不盲目執行 DDL。custom 表不 ALTER。
- 驗收：真空白 SQLite UI→schema confirmation→SQL 垂直流程先取得回饋；七型別、兩 DB、scope/replay/ownership/receipt/恢復測試仍是完整交付條件。
- 範圍：群組 schema、必要 metadata/exact layout 與恢復接線；不新增 token 家族、ledger、sampling mode 或 migration engine。

## Risks / Trade-offs

- 外部 DDL 與本地 ledger 不是同一交易：保留 succeeded/partial/failed/unknown 與實際 inspection，不宣稱跨 DB rollback。
- 自動欄名與 SQL 型別：必須保存生成結果、quote identifiers 並測試同名、中文、保留字及大整數；不新增型別猜測。
- 舊 managed 表少欄位：不偷偷升級，不修改 accepted payload；透過確認的新表設定處理。

## Migration Plan

只 additive 擴充已有 group/version/operation 所需資料；不重建 journal/outbox。回退保留已建立表、外部 committed effects 與操作證據。能力開放以 SQLite／PostgreSQL 各自實際驗證為準。

## Verification

先 RED：空目標無人工 SQL、SQLite preview 無檔案副作用、confirmed create、重送／失去回覆、同名衝突、權限不足、來源與連線 revision 過期。再驗證真 metadata 與延遲補送的原時間／身分；custom 既有表不變。
