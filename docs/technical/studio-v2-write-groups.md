# Studio V2 WriteGroup：設定、相容移轉與回復

本文件對應 A `unify-studio-v2-write-group-contract`，並在「Durable 交付（C）」「確認試寫（D）」說明 C `wire-durable-write-group-delivery`、D `implement-confirmed-write-group-test-write` 的運行行為。A 建立本地設定 authority、相容 API 及交易邊界；B 的 sample 規則已實作並歸檔；C 的 durable production delivery 實作與驗證見 C 的 validation（任務尚未全部完成）；D 的真實試寫、E 的 UI 及 F 的 UI→SQL 驗收仍按順序交付。完成設定或 domain Apply 不表示 SQL committed／readback verified。

## 本地資料結構

`internal/datalink/schema/migrations/024_write_groups_sqlite.up.sql` 擴充 gateway 的 SQLite 設定資料庫，與 destination SQLite／PostgreSQL 是不同資料庫。

| 表 | 保存內容 |
| --- | --- |
| `write_groups` | workspace/group ID、draft/applied revision、狀態、destination identity、row/write policy、migration provenance |
| `write_group_members` | 真實 device/point/tag、source/mapping revision、optional measurement、column、required/freshness 與 optional entity key |
| `write_group_versions` | group/revision 的不可變 payload 及 UTC effective_at，保留舊 bucket 的設定 |
| `write_group_migration_maps` | workspace/source kind/source ID→stable group ID、原 intent/source revision、review digest/adapter version |

foreign key 使用 `ON DELETE RESTRICT`。logical delete 留 tombstone；不 cascade、purge 或重用 group ID。`entity_key` 僅是 row partition，不代表 SQL business-key 值。group/unique-key columns 是已保存的規劃 metadata，不表示 destination 已有 verified SQL constraint。

Migrator 可重跑建立缺少的 024 表，並以 additive `ALTER TABLE` 補既有 member table 的 `entity_key`；不重建或刪除 group/member。部分 schema 修補不能還原已經遺失的資料，仍須使用有效備份。

## API 與 revision

共同前綴為 `/api/v1/datalink/studio-v2/workspace/write-groups`。正式 DTO／responses 見 `docs/swagger/swagger.json`；frontend SDK 在 `frontend/src/services/studioV2WorkspaceWriteGroups.ts` 與 `studioV2WriteGroupMigration.ts`。

| 方法／尾碼 | 行為 |
| --- | --- |
| `GET`（無尾碼）、`GET /:id` | workspace-scoped list/get；已移轉資料由 canonical authority 重載 |
| `POST`（無尾碼）、`PUT /:id` | 保存 draft，group/member/projection/workspace revision 同一 local transaction |
| `GET /:id/readiness` | 使用 persisted group 驗來源與 read-only destination schema，不建表、不啟用設備 |
| `POST /:id/disable`、`DELETE /:id` | 停新 intake／logical delete；保留 immutable snapshots 與 accepted backlog 歸屬 |
| `POST /migrations/{single-mappings,row-groups,recording-plans}/preview` | 本地只讀預覽，帶 revisions/source intent/digest/issues |
| `POST /migrations/{single-mappings,row-groups,recording-plans}/review` | 綁定同 revision/digest 的明確審閱，支援來源只保存 draft；blocked plan 不轉換 |

mutation 帶 `workspace_id`、`expected_workspace_revision`、`expected_connector_revision`；create 以外另帶 `expected_group_revision`。`expected_workspace_revision` 對應 `Record.DatabaseSetupRevision`，不是 Share 的獨立 workspace revision。版本過期409；unknown／foreign安全404；invalid422；operational storage failure安全500，不顯示 SQL／DSN／credential。

A 沒有 HTTP Apply endpoint 或 production owner consumer。domain `WriteGroupService.Apply` 在 readiness preflight 後重查 local CAS/live connector/source，於下一 UTC bucket 排 immutable applied version。可注入 `WriterOwnershipActivationBarrier.ActivateInTx`，在同一 local transaction 寫 owner projection；version/applied/owner/workspace 同 commit 或 rollback。metadata-only rename 不重啟 owner transition。C 接線時必須注入真實 consumer barrier，並保留現有 Share hydration/settings/readiness-token/revision gate；不得把 domain applied metadata 當成 production writer 成功。

## 三批相容行為

1. single mapping：原 writer 是逐筆；preview 明列改為週期 snapshot 的選值、late、missing/freshness 與 timestamp 差異。`confirm_snapshot_conversion=true` 與 current revisions/digest 才保存 draft，原 target 及 enabled 狀態不變。
2. row group：原 row ID、targets、GroupKey→entity_key、member 與 key metadata 保留。共享 column 只有 distinct entity 可表達；同 row/column 碰撞、缺身份或目標混合等情況 blocked。review 不啟用 writer。
3. recording plan：保留完整原 plan/source intent，沒有猜測 candidate。current review422、stale409、foreign404；原 plan revision/status/applied identity 不變，基本新寫入使用 canonical Create。

single/row reviewed migration 保存 durable ID map。相同 source revision 重跑回同一 group，不再增 workspace revision、不覆寫後續 canonical edit。來源變動需重新 preview。legacy reads 對 owned mapping 使用 canonical projection；無法表達時回 actionable conflict。

已移轉 scope 的 global／Studio target CRUD 與 owned row-group replacement 回409 `WRITE_GROUP_LEGACY_WRITE_CONFLICT`、`open_write_groups`。guard 同時檢查 mapping ID、原始 Tag/connector 及 canonical current Tag/connector；在 mutex／同 transaction 內先於 legacy persistence 重查。未移轉 CRUD 保留。nil row_groups 是不替換；explicit `[]` 不能抹掉 owned layout。

## Durable 交付（C）

`025_write_group_delivery_sqlite.up.sql` 在同一個 gateway SQLite 設定資料庫加入 `wg_delivery_samples`（sample journal）、`wg_delivery_checkpoints`、`wg_delivery_buckets`（每個 entity／bucket 的 row／skipped／no_data 結果）、`wg_delivery_outbox`（凍結 group／destination revision 與 payload digest 的待交付 row）、`wg_delivery_receipts`。與 destination 資料庫無關。

流程：樣本在 journal transaction commit 後才 ACK → 時間驅動關閉 bucket，把 row、outbox、consumed 與 checkpoint 在**同一個本地交易**提交 → delivery worker 以 claim（owner＋遞增 fencing epoch＋lease）逐 partition（group＋entity）依序交付，同一 partition 前面的 row 未完成（retrying／blocked／quarantined／unknown／sending）就不會越過，其他 partition 不受影響。`Pipeline`（`internal/datalink/grouppipeline`）在 production 啟動時 reconcile 已 Apply 的群組、驅動 tick 與 worker，並告訴 legacy writer 哪些輸出已被群組擁有（一個輸出只有一個 writer）。

交付狀態只有 destination 確認過才是 `sql_committed`：

| 狀態 | 意義 |
| --- | --- |
| collecting（journal） | 已 ACK，bucket 尚未關閉 |
| queued／retrying | 已關閉成 row 並持久保存，等待或重試；重試有上限，用盡轉 blocked，資料不刪 |
| blocked | 需修復（destination revision 改變、停用、憑證被拒、schema 不符、重試用盡） |
| quarantined | destination 永遠拒絕這個 row（資料或完整性錯誤）；payload 保留，只擋同 partition 後面的 row，需 operator 明確 `retry`／`skip` |
| unknown | 無法確認 destination 是否已 commit（無 dedupe 能力時的 commit 回應遺失、中途崩潰）；不會自動重送，需核對 |
| sql_committed | destination 已確認；`last_sql_committed_at` 只來自這個狀態 |

讀取：`GET .../write-groups/:id/delivery` 回傳各階段計數、最舊未完成年齡、按 group revision × destination 分組的 backlog（含該 backlog 被接受時的 connector revision 與安全 error code）、intake 狀態與配額。**這是 C 的後端來源；UI 顯示屬 E。**

**dedupe 能力與限制（沒有 universal exactly-once）**：群組的 `write_policy.dedupe_capability` 決定 commit 結果不確定時能否安全重試——`receipt`：destination 同交易寫入 `gw_effect_receipts`（需存在，由 managed schema 流程或操作者建立；sender 不執行 DDL），重試先查 receipt，相同 digest 視為已 commit、不同 digest 阻擋且不覆蓋；`unique_key`：需要 inspection 確認為 UNIQUE／PK 的 record key 欄位（A 的 row policy 目前沒有此欄位設定，需 E 補）；`none`（custom 表預設）：只有單純 insert，commit 後回應遺失即 `unknown`，**不會自動重插**。MySQL 目前沒有群組交付策略（回明確不支援）。PostgreSQL 16 與 SQLite 已有真實驗證。

運行界限（都是設定值而非效能保證）：每個 partition 每輪最多 50 筆、最多 4 個 partition 並行；單次 destination 嘗試 30 秒、寫回本地狀態 5 秒；lease 至少涵蓋一次嘗試；關閉時先停新 row、等在途交付到 deadline，逾時才取消，被中斷的 row 由下次啟動恢復；累積未交付資料上限預設 500 MiB（量測 payload 位元組），達上限只拒絕新的 ACK、不刪已接受資料，並顯示範圍與 loss-risk。

**回復保護**：025 只新增表，重跑 migration 不改動已接受資料（有測試）。停用群組、編輯 destination endpoint、binary 回退都不會改變舊 backlog 的去向——舊 row 保持凍結的 group／destination revision，endpoint 被編輯後成為 `blocked`，不會被送到新 endpoint。停用群組只停止新 intake，且在 endpoint 被編輯之後仍可執行。若需回退 binary，先停新 intake，保留 journal／outbox／receipt 與 destination 已提交的效果；不得用舊備份覆蓋新 accepted data。

## 確認試寫（D）

`026_operation_test_write_sqlite.up.sql` 為既有 `managed_schema_operations`（建表 operation 帳本）新增 `payload_digest`、`write_outcome`、`cleanup_status`、`cleanup_reason`、`detail` 五欄；舊 row 與建表流程不變。試寫沒有另一套 token 或 operation repository，只是同一帳本中 `action = test_write` 的 operation。

| 路由 | 行為 |
| --- | --- |
| `POST …/write-groups/:id/test-write-preview` | 讀已保存 group 與真實 destination 表，回傳要寫入的每個欄位值、擁有者欄位與值、dedupe、清理方式，並保存綁定 workspace／group／connector revision 與內容 digest 的 token。**不寫 target。** |
| `POST …/write-groups/:id/test-write` | 帶 `token`、`operation_id`；原子 claim 一次後寫入、讀回、清理。 |
| `POST …/recording-plans/test-write-preview`、`…/test-write` | 舊 plan 路由，只靠 group 的 migration provenance 解析到**唯一**群組，解析不到 422 `RECORDING_TEST_WRITE_PLAN_UNRESOLVED`，不猜第一筆；確認一樣必須帶 token 與 operation_id，裸 `plan_id` 回 400。目前沒有任何 production 流程寫入這種 provenance（recording plan 的 migration review 仍是 blocked），所以今天這兩條實際上都會回 422。 |
| `GET …/database-operations/:operation_id` | 沿用的共用狀態端點，也服務 test_write。 |

狀態碼：running 重送 202；已保存結果 200（token 之後過期也一樣，不再執行）；同表另一個 operation 佔用、stale、首次 claim 已過期 409；缺欄 400；未知／他人的 token 404；他種 action 的 token（例如建表 token）或不支援的目標 422；結果無法保存 503 `WRITE_GROUP_TEST_WRITE_RESULT_UNKNOWN`（附 operation_id，重送會接手並從 target 證據收尾，不是重寫）。

**擁有權（為什麼有些表不能試寫）**：試寫 row 的 entity key 欄位寫入 `gw-test-<operation uuid>`，讀回與清理都只用這個完全相同的值：`DELETE … WHERE <entity 欄位> = <該值>`，值不在 `gw-test-` 命名空間內一律拒絕。因此 group 的 row policy 必須設定 `entity_key_column`；沒有的表（無法證明哪一列是這次試寫的）preview 回 422 `WRITE_GROUP_TEST_WRITE_UNSUPPORTED`，不做任何寫入、也不做廣泛 DELETE。`receipt` dedupe 另需 `gw_effect_receipts` 存在；`unique_key` 目前沒有 record key 欄位設定，試寫改用無 dedupe 的單純 insert。

結果欄位彼此獨立：

- `write_outcome`：`written_verified`（讀回的型別化值、擁有者標記、provenance 與 receipt 都與寫入內容相符）、`written_unverified`（已 commit 但讀不到、被拒讀取、或內容不符；`reason` 為安全代碼）、`failed`（確定沒寫入）、`unknown`（commit 結果無法確認，且 target 證據也無法判定）。
- `cleanup_status`：`cleaned`、`failed`（例如缺 DELETE 權限，row 仍在）、`unknown`（清理 commit 回應遺失）、`not_attempted`（沒有寫入）。
- operation `status`：`succeeded`↔verified、`partial`↔unverified、`failed`、`unknown`。

寫入使用 production row layout（`GroupRowLayout.EncodeRow`＋`InsertGroupRow`）與 C 的 dedupe 策略，但**不經 delivery outbox**：試寫不進 backlog／quota／partition 順序，且必須等到 target commit 才有結果，佇列中不會被當成 written。

重啟：每一步的副作用前先保存進度（`claimed`／`writing`／`cleaning`）。lease（3 分鐘）過期後，重送同一 operation 會由新 owner 接手：`writing` 階段若 target 已有擁有者 row 就直接讀回、清理；沒有時只有 `receipt` dedupe 會用同一 effect key 重試，其他一律 `unknown`（可能已 commit，不重插）；`cleaning` 階段只完成清理，不會因為 row 已消失而重寫。group 在試寫後被編輯則無法重建預期 row，擁有者 row 只回報 `written_unverified`（`group-changed-after-write`）並清理。

**目前只是後端／API／前端 service 與 hook。** Step 4 的試寫按鈕仍由 `supports_test_writes`（刻意維持 false）把關，直到 E 整合；新增的 `supports_group_test_writes`（SQLite／PostgreSQL 為 true）只表示後端能力，個別 group 仍可能在 preview 被拒。MySQL 沒有群組試寫。

## 備份與回復邊界

任何真實設定資料庫 migration／部署需另獲授權；本輪只驗可丟棄 fixtures。下列是 deployment-owner 的操作契約，不表示本次已在正式資料執行。

1. 確认實際設定 DB 路徑、binary/version、schema version、workspace/connector/group revisions 與 pending/accepted backlog；不得由範例路徑推定 target。
2. 停止該設定 DB 的 writer，使用 SQLite backup API／`.backup` 建立一致備份，或所有 DB connections 關閉後一併保存 DB/WAL/SHM；不要僅複製仍在寫入的 `.db` 檔。備份含 secrets，沿用既有權限保護，不提交到 Git。
3. 在備份副本執行 migration，核對 `integrity_check`、原 workspace/target/recording-plan payload、group/member/version/map 與 revisions，再允許正式啟動。重跑 migration 不等於證明資料已轉換；review 仍要求當前 digest 與明確確認。
4. C 前尚無新的 production intake，本地 schema 可保留 additive tables。回復程式碼／UI不刪群組表或 map，也不解除 owned scope 的 legacy write 保護。舊版若不理解 canonical ownership，不能直接開啟已移轉 scope 的 legacy writer；先停該 scope intake，交由相容 reader／owner-aware版本處理。
5. C 以後只可先停新 intake。accepted payload、immutable revision、destination provenance、outbox／receipt 與外部已提交效果必須保留；不得用舊備份覆蓋新 accepted records，不得清空表或重新插入「補資料」。實際 owner rollback/redeploy 由 C/F 的 verified procedure 與 deployment-owner驗收決定。

024 沒有 destructive down migration。回復既有資料副本只限確認沒有新增 accepted data／external effect 的停機情境；否則保留新 journal／outbox／receipt，不能把 binary rollback 與 data rollback 混為一件事。

## 驗證證據

本輪命令、RED/GREEN、before/after payload、平台與未執行項目記於 `openspec/changes/archive/2026-10-02-unify-studio-v2-write-group-contract/validation.md`。本地 tests／readiness／domain lifecycle 不能代替真實 PostgreSQL、production UI→SQL、Windows／ARM、LAN／PLC、現場或部署回復驗收。
