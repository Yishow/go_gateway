# go_gateway 未完成功能規劃（2026-10-09）

依據：README 能力邊界、docs/plans/studio-v2-write-groups（A→F 已歸檔）、
openspec/changes/add-windows-tray-runtime-logs（規格已提交、部分實作）、
2026-10-09 現場 MC 真機 E2E 實測（先 5 台，後擴展至 7 台＋ET7017）、
2026-10-11 現場 8 台設備真機 E2E 驗證、2026-10-13 現場 Modbus UDP 實體採集 E2E 結案。
事實優先，猜測處標注。

> 下列 10/09–10/13 現場結果來自既有未提交工作紀錄，保留為歷史來源陳述；原始現場見證未隨本文提交，這輪未連現場或正式 DB 重驗。私有位址、Tag/group/資料表執行識別已省略。「結案」僅指當時紀錄；2026-10-23 的可重跑本機驗證另見本輪交付說明。

## 現況定位

- 核心「設備→採集→分組寫 DB」已端到端驗證（真機：MC×7、Modbus TCP、FATEK probe×2 → PG 有列有 provenance）。
- **2026-10-13 現場 Modbus UDP 實體採集 E2E 結案（最新）**：
  - 設備：泓格 ET-7017（私有位址省略，`modbus_udp` 協議），點位 40001 AI0，Tag 專用 UDP 測試 Tag。
  - 入庫：Basic Write Group （識別省略）；Schema Preview／Apply 後
    在 PG 派表 指定受管測試表。
  - 根因排除（採集→入庫原先未通的卡點）：
    1. 建立對應的 `source_rules` 及 `source_rule_links`，納入 RuntimeProjection。
    2. 對齊 Tag labels 中的 source-rule 歸屬中繼標籤，並校準 Revision 雜湊。
  - 結果：重啟後每秒 sample 穩定流入 buffer，並連續每分鐘以 `kind: row` 穩定寫入 PostgreSQL
    （讀值 0，附完整 provenance 審計記錄）。Modbus UDP E2E 正式結案；P3 的 UDP 項目隨之關閉。
- **2026-10-11 現場真機 E2E 實測驗證**：
  - 設備：8 台在線設備全數建立連線。MC×7（設備／製程名稱與私有 port 省略）；Modbus TCP×1＝ET-7017。
  - 入庫：8 台皆納入 PG 測試表 指定混合設備測試表；26 個資料欄位中 25 個每分鐘持續穩定落庫
    真實讀值，一台現場斷電設備為 NULL。
  - 兩台 MC 設備落庫異常的根因與排除：
    1. `source_rule_links` 補齊 `mapping_id`。
    2. mappings 與 points 的 enabled／active 狀態對齊。
    3. `write_group_members` 與 snapshot 版本的 Revision SHA-256 雜湊一致性校準。
- **2026-10-09 實測（現場擴展）**：
  - 設備：8 台實體設備、共 25 個點位。MC×7（設備／製程名稱省略）；另 1 台 Modbus-ET7017。
  - 入庫：PG 表 指定多設備測試表 每分鐘連續聚合入庫；25 個點位欄位皆有真實讀值
    （含整數、負數與 0）；provenance 25/25 為 `status: ok`、`quality: good`。
  - 唯讀鐵律：實測全程無寫入 PLC 指令。
  - 品質閘：mc3e adapter 的 contextcheck 缺陷已修（`withConnectionRetry`／`retryExecuteRead`
    改傳入 ctx，退避由 `time.Sleep` 改為可被 ctx 取消、重連不再用 `context.Background()`）；
    全庫 golangci-lint 0 issues；`make check-lines` Passed（`shell.test.tsx` 為 656 行的 legacy
    超限檔，維持 656→656 未增加；尚未拆到 500 行以下）。
  - 狀態：上述「仍在工作樹」為當日歷史狀態，不能用來推論目前是否未提交；本輪未修改這些歷史修復。
- A→F 六案 54/54 驗收歸檔（2026-10-03），durable delivery（journal/outbox/sender/restart）已實作。
- 測試基線：go test 全綠（2522 tests / 59 pkgs）；`internal/diagnostics` 與 `cmd/test_ui`
  tasks.md 提到的 5＋1 個 file-sink 失敗**已不復存在**（52＋95 全過，HEAD 1456c8dc）。

## 缺口清單（依優先序）

### P1 — UI 流程收斂（2026-10-09 開工調查後大幅翻案縮小）
**調查結論（翻案）**：原列的「7 個非直覺步驟」幾乎全是**進階 legacy API 鏈**的行為；
UI 基本路徑（Step4 BasicRecordingPanel）本來就走全自動正道：
`buildBasicManagedDrafts`（members 只帶 device/point/tag，`target_column: ''`、`table_name: ''`）
→ `POST /write-groups/basic/:deviceId`（ensureBasicManaged，server 派表名 `gw_group_<hash>`）
→ `POST /recording-start`（帶 draft，server 內部 update→readiness→apply→activate 一條跑完，
欄位 `v_<hash>` 也是 server 派）。我現場手組的那套（config PUT 加成員、
database-target 掛 row group、write-group PUT 完整 members、managed 表加欄）
是不受支援的組合， 卡點屬**設計上的 fail-closed**，不是 UI 缺口。

**真正修掉的缺陷（已修，commit 見 git log）**：
`database_outputs` 候選快照在 tag 的 database target mapping **晚建立**時，
GET `/source-rules/:id/candidates` 不會自動重算（`shouldRecomputeDatabaseCandidateView`
只認 blocked），導致下一次 `database-outputs/apply` 一次就 `schema_missing`
（"database mapping scope is not configured"），必須手動 `POST .../candidates/recompute`。
修法：ready 快照若解碼出任何候選缺 `mapping_id` 也觸發自動重算
（`databaseCandidatesMissingMappings`）。TDD：RED（`candidate_api_late_mapping_test.go`
先證 GET 不補 mapping_id）→ GREEN（sourcerule 141 案＋golangci-lint 乾淨）。
全套迴歸：**go test 2542/59 packages 全綠**。commit `b517910`（2026-10-09）。

**剩餘（縮小後）**：API 文件把 candidates query scope（workspace_id＋
expected_workspace_revision＋revision_id）與空字串語意寫進 README/技術文件。
驗收不變：新手依 UI 完成 2 設備→新 PG 表，不需看 API。

### P2 — test-write 契約統一（D 案遺留）— 已清理更新
- 前端已移除 legacy plan testWrite 呼叫，Step 4 只走 `write-groups/:id/test-write`
  （grouptestwrite）；後端 `recording-plans/test-write(-preview)` 改為解析到同一
  group 與 service（雙軌對齊），不再有獨立 token／ledger；README 已移除 501 限制。
- 殘留：前端 `TestWriteResult`／`supports_test_writes` 型別與 501 錯誤碼的 i18n、
  測試仍在；`openspec/specs/recording-database-setup` 仍描述 501。
- 未含：confirm→write→readback→cleanup 實跑紀錄未在本次更新補上。

### P3 — Modbus UDP / RTU / FATEK 採集 E2E（UDP 已於 2026-10-13 結案；RTU／FATEK 待辦）
- UDP — **已結案（2026-10-13）**：原本只證明原始封包對 ET-7017 有真回應；現已跑完
  gateway 採集→PG 全鏈路。
  - 設備：泓格 ET-7017（私有位址省略，`modbus_udp`），點位 40001 AI0，Tag 專用 UDP 測試 Tag。
  - 入庫：Basic Write Group （識別省略），Schema Preview／Apply
    後派表 指定受管測試表。
  - 根因排除：建立對應 `source_rules`／`source_rule_links` 並納入 RuntimeProjection；
    對齊 Tag labels 的 source-rule 歸屬中繼標籤與 Revision 雜湊。
  - 驗證：重啟後每秒 sample 穩定流入 buffer，連續每分鐘以 `kind: row` 寫入 PostgreSQL
    （讀值 0，provenance 審計記錄完整）。此項不再需要原訂的「一小時批次」補測。
- RTU：現場無串列設備；FATEK：probe 過（私有位址省略）未長跑。
- 動作：FATEK 找一台可連續採集的機台跑 24h。RTU 列「無現場條件，持續不保證」。
- 2026-10-23：連線必填問題由本輪獨立 change 彙整 missing／empty，保留既有協定及 range／enum 規則；此項以本機測試證據驗收，不等於新增實體協定 E2E。

### P4 — grouped buffer 耐重啟 — 已驗證就緒（2026-10-11 結案）
- 結論：經代碼與實體重啟驗證，Grouped Delivery 已完整由 SQLite WAL 的
  `wg_delivery_samples`、`wg_delivery_buckets`、`wg_delivery_outbox` 機制落盤；
  已 journal commit 的資料可在重啟後依原 scope／revision／ownership 證據恢復；無法證明外部效果時保持 unknown／blocked，不能盲目重送。原「全部 grouped buffer 僅記憶體」描述不適用正式 WriteGroup 路徑。
- 代碼依據：`internal/datalink/schema/migrations/025_write_group_delivery_sqlite.up.sql`
  建表、`internal/datalink/groupdelivery` 為 store／dispatcher，內部 SQLite DSN 為
  `journal_mode(WAL)`；崩潰重啟測試見 `cmd/test_ui/group_pipeline_lifecycle_crash_test.go`。
- README 能力邊界「資料庫寫入」已同步更新，移除不耐重啟的描述。
- 剩餘：gateway 重啟後 recording intent 自動續寫（實測行為），文件仍需寫明此語意與停止方式。

### P5 — Windows tray/無 console 發行（進行中，部分完成）
- add-windows-tray-runtime-logs：2.1/2.3/3.1/3.2/3.3/6.1–6.5 有實作及自動化證據（原20項中10項），
  剩 1.1/1.2（lifecycle 模式）、2.2（file sink 輪替）、
  4.1–4.3（tray/單實例/shutdown）、5.1–5.4（打包/實機/文件）。
- 現場部署目前用 ScheduledTask 跑 console exe；此案完成前維持現狀，
  完成後才換 GUI 發行。

### P6 — 觀測性小項（低）
- MC EOF 類錯誤已帶 port 提示（1456c8d）；再加裝置級 metric
  （per-device reconnect 計數、最後成功讀值時間）讓 UI 不用翻 log。

## 明確不做
- 不建 V3、不恢復 /studio、不新增第二套 write ledger。
- aggregation/measurement/delivery 進階能力維持「明確選擇才啟用」。
