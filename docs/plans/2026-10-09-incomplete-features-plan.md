# go_gateway 未完成功能規劃（2026-10-09）

依據：README 能力邊界、docs/plans/studio-v2-write-groups（A→F 已歸檔）、
openspec/changes/add-windows-tray-runtime-logs（規格已提交、部分實作）、
2026-10-09 現場 5 台 MC 真機 E2E 實測。事實優先，猜測處標注。

## 現況定位

- 核心「設備→採集→分組寫 DB」已端到端驗證（真機：MC×5、Modbus TCP、FATEK probe×2 → PG 有列有 provenance）。
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
先證 GET 不補 mapping_id）→ GREEN（141 案 sourcerule 全過＋golangci-lint 乾淨）。

**剩餘（縮小後）**：API 文件把 candidates query scope（workspace_id＋
expected_workspace_revision＋revision_id）與空字串語意寫進 README/技術文件。
驗收不變：新手依 UI 完成 2 設備→新 PG 表，不需看 API。

### P2 — test-write 契約統一（D 案遺留）
- `recording-plans/test-write`（legacy plan 路徑）與 `write-groups/:id/test-write`
  並存，兩者都要 token＋operation_id；README 明講 plan 路徑仍是 501
  `RECORDING_TEST_WRITE_NOT_IMPLEMENTED`。group 路徑已實作（grouptestwrite）。
- 動作：確認 UI 只用 group 路徑後，把 plan 路徑退役或補齊，README 對應更新。
  驗收：README 能力邊界不再列 501；一次 confirm→write→readback→cleanup 實跑紀錄。

### P3 — Modbus UDP / RTU / FATEK 採集 E2E（本次只驗到通訊層）
- UDP：原始封包對 ET-7017 已證真回應，但未跑 gateway 採集→DB。
- RTU：現場無串列設備；FATEK：probe 過（.162/.165）未長跑。
- 動作：UDP 用現有 ET-7017 補一條採集 E2E（一小時批次即可）；FATEK 找一台
  可連續採集的機台跑 24h。RTU 列「無現場條件，持續不保證」。
- 附帶：test-draft 的必填欄位錯誤訊息一次列齊（目前 mode→station_no 逐次擠牙膏）。

### P4 — in-memory grouped buffer 耐重啟（README 明講限制）
- grouped buffer 在記憶體，重啟即失；C 案的 durable delivery 已涵蓋
  journal/outbox，需確認 grouped 寫入路徑是否已接 C 案機制（事實待核）。
- 動作：核對 grouppipeline→dbtarget writer 是否已用 journal；未接則把
  buffer 落 journal，接了就更新 README 移除這條限制。
- 已知相關行為：gateway 重啟後 recording intent 自動續寫（本次實測），
  文件要寫明這個語意與停止方式。

### P5 — Windows tray/無 console 發行（進行中，部分完成）
- add-windows-tray-runtime-logs：2.1/2.3/3.1/3.2/6.1–6.5 已完成（9 項），
  剩 1.1/1.2（lifecycle 模式）、2.2（file sink 輪替）、3.3（/studio/logs UI）、
  4.1–4.3（tray/單實例/shutdown）、5.1–5.4（打包/實機/文件）。
- 現場部署目前用 ScheduledTask 跑 console exe；此案完成前維持現狀，
  完成後才換 GUI 發行。

### P6 — 觀測性小項（低）
- MC EOF 類錯誤已帶 port 提示（1456c8d）；再加裝置級 metric
  （per-device reconnect 計數、最後成功讀值時間）讓 UI 不用翻 log。

## 明確不做
- 不建 V3、不恢復 /studio、不新增第二套 write ledger。
- aggregation/measurement/delivery 進階能力維持「明確選擇才啟用」。
