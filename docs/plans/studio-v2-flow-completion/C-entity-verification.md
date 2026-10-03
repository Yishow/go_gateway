# 同一寫入群組的 entity 列驗證

2026-10-04；change：`fix-write-group-entity-row-layout`；起點：`c26c8a52a1712d922015ba2d3822263c76a5dd7d`。

## 實作與相容性

- 沿用既有 entity-to-members 分割：同一群組 A/B 可共用欄位，各自成列，不拆成兩個群組。
- 前端以 trim／小寫判斷同列欄位衝突，保留原欄名供提示；不同 entity 的同欄位仍合法。
- readiness 的同列碰撞會提供安全且可操作的欄位修正說明，保留既有 issue code；兩資料庫的反例先失敗、修補後通過。
- SQL 使用 inspected 欄名，先採精確匹配再解析大小寫別名，避免 PostgreSQL quoted identifiers 寫錯欄位。
- durable closure 的結構性編碼失敗回傳安全 `layout-blocked`，保留 ACK journal、checkpoint 與未完成 closure；修復並重開 store 後可重播。
- 既有品質 incomplete／partial／no_data 與 `sql-value-blocked` 拒絕不變。舊 group/member key、record ID、effect key、digest、entity 與 pending payload 均未 migration。
- production 仍由 `service_wiring.go` 接既有 GroupRowLayout、boundary、SQLite journal/outbox 與實際 SQL sender；沒有測試專用成功路徑。
- custom SQL 列身分來自設定的 `entity` 欄與 effect receipt；不宣稱 production 自動新增 SQL business key／record_id 欄。

## 規格對照

| Requirement／Scenario | 實作與測試 | 此次證據 |
| --- | --- | --- |
| Entity-scoped column layout agreement：Shared column in distinct rows | `findColumnConflicts`、`NewGroupRowLayout`；`TestProductionEntityLayoutReadinessAndApplyAgree`、前端 columns/draft tests | 前端、真 SQLite／PostgreSQL readiness 與 Apply 一致接受；4 members 保留 |
| 同上：Competing values in one row | 同 entity、大小寫變體的前端反例、Readiness／Apply 反例 | UI 提示並停用儲存；後端拒絕，無 applied revision |
| Entity-local encoding and recoverable structural failure：Independent values in one group | `membersForEntity`、`EncodeRow`；`TestProductionEntityRowsSQLAndFrozenRestart` | 共欄、不同欄、大小寫別名及 PostgreSQL 精確 quoted 欄名的真 SQL readback；payload／record／effect／receipt digest 與兩次重啟一致 |
| 同上：One entity lacks data | SQL matrix、`TestEntityRowsKeepMissingPartialAndSilentOutcomesLocal` | A required missing 時 B 正常；partial 的 A NULL／provenance 只含 A；silent 每 entity 一個 no_data，沒有 SQL 假列 |
| 同上：Structural mismatch after durable acceptance | `structuralEncodingFailure`、durable Tick；`TestEntityStructuralFailureKeepsAcknowledgedJournalAndCheckpoint` | 5 個 ACK 樣本保持 unconsumed；checkpoint／bucket／outbox 均 0；修復並重開 store 後成功重播 |

2 requirements、5 scenarios 全有測試；0 Examples、0 scope exclusions。analyze 的 COV-1／COV-2 是 tasks 未逐字重複 requirement 名稱，上表完成語意對照。

## 失敗先於修補的證據

- 前端同 entity 的 `Value/value`：修補前 2 項行為失敗，修補後受影響 3 檔 33 項通過。
- structural mismatch：修補前 Tick 回 nil、closure 被當 skipped 消耗；修補後保留 journal、安全錯誤及實際重開重播通過。
- SQLite／PostgreSQL 欄名別名：修補前輸出原 `Temperature/ENTITY`，與 inspected 名稱不一致；修補後通過。
- PostgreSQL 精確 quoted case-distinct 欄名：新回歸測試先失敗（被改到小寫欄），精確匹配修補後 unit 與實際 PostgreSQL SQL 通過。
- 既有 entity 分割未重做；Go overlay 恢復 global column claim 後，兩個資料庫的合法共欄 readiness 反例均失敗。遍歷全群組 members 的 overlay 亦用獨立列 SQL 反例核對。
- 日誌：`/private/tmp/go-gateway-C-column-red.log`、`go-gateway-C-structural-red.log`、`go-gateway-C-identifier-red.log`、`go-gateway-C-quoted-red.log`、`go-gateway-C-action-red.log`；overlay 在 `/private/tmp/gw-C-mutations/`，未覆寫 worktree。

## 真實畫面與 SQL

- 正常 `cmd/test_ui` binary 提供 embedded `/studio/v2`；實際 frontend build／static sync，沒有 mock API、fixture build tag 或直接插入目的地資料列。
- agent-browser／Chrome for Testing；本次 loopback Modbus simulator A/B（15030／15031）及可丟棄 SQLite。UI 完成兩設備、各兩點、映射、單群組四 members、A/A/B/B、10 秒 bucket、Save／Apply／啟動。
- 獨立 SQL：A=`215/1013`、B=`187/777`。A simulator 停止後 A 維持 11 列，B 由 11 增到 15、後續到 67；local buckets 顯示 A skipped／品質不良，B 持續 row／SQL committed。
- UI 區分已儲存、已套用、採集、本地待交付、資料庫已確認；silent bucket 仍呈現 no_data。
- 鍵盤：entity input Tab 到下一列 checkbox；Enter 完成 Save／Apply。390px 表格 region 可 focus，ArrowRight 實測 scrollLeft=513、document width=390。
- 390px 用既有按鈕收合導覽，保留表格局部橫向捲動；錯誤摘要與儲存／交付操作可見。768／1440px 亦實測。沒有 Design Source；依 design 的同群組、成列與品質語意比對，未新增視覺模式。
- `evidence-c/{main,error,partial}-{390,768,1440}.png` 共 9 張已視覺檢查；error 三張含實際衝突、停用儲存及最後已存版本提示。
- frontend build 內容在畫面驗證期間未再變動。最後 backend 欄名修補後重建正常 binary、重啟同一 cfg DB；applied revision 不變、readiness 通過、B SQL 持續增加。最終 binary SHA256：`529e4bd51d20f51a5b58d496d549845c63b82952ac1d159d1b98844b6825a0f0`。
- browser daemon 等待異常經 doctor 查證後重開本次工作階段；不是產品通過證據。未使用正式 PLC／DB、LAN 或部署。

## 執行檢查

darwin/arm64、Go 1.27.1（module 1.25.5）、golangci-lint 2.14.0；自有 PostgreSQL 16.14，loopback 55432／`gwtest`。未碰既有 5432 服務。

| 命令與範圍 | 結果 |
| --- | --- |
| `go test -race ./cmd/test_ui ./internal/datalink/workspace ./internal/datalink/dbtarget ./internal/datalink/runtime -run 'TestProductionEntity\|TestEntity\|TestAtomicRowOutboxCheckpointBoundaryBlockedEncoding\|TestWriteGroupReadiness' -count=1`，自有 POSTGRES_DSN | 通過；SQL matrix 13 情境；SQLite case-only quoted 欄位情境不適用，明確 skip |
| `POSTGRES_DSN= go test -p 1 ./... -count=1` | 通過；本案最後變動後另補受影響檢查 |
| `go vet ./...` | 通過 |
| `golangci-lint run ./...` | 通過，0 issues；測試 if/switch gate 修正後通過 |
| `cd frontend && npm run lint` | 通過 |
| `cd frontend && npm test -- --run` | 182 檔、1107 tests 通過 |
| `make sync-frontend-static`（含 frontend production build）及正常 Go binary build | 通過 |
| `git diff --check`、`make check-lines` | 通過；最後文件／封存後另核對 |

未重跑整套帶 POSTGRES_DSN 的 repository suite；A 已記錄的 recording operation／歷史 partition migration 既有失敗不納入本案成功聲明。本案是本機 simulator、可丟棄 SQL 與畫面驗證，並非現場設備或部署驗收。

