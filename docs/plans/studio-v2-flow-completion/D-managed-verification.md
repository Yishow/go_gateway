# 受管寫入群組的 schema 與 SQL 驗證

2026-10-04；change：`complete-write-group-managed-storage`；起點：`b560926c88ce4b477f53ec73e20e5d721964ba81`。

## 實作邊界

- 保存的 canonical WriteGroup 直接接既有 recordingplan schema service、preview token、operation ledger 與查詢入口；沒有假的 RecordingPlan、measurement、第二套 token 或 ledger。
- 新 managed 群組保存穩定的 `gw_group_<identity digest>` 表名、`v_<member identity digest>` 欄名及系統 metadata bindings；Create／Update 的留白表名均在保存時解析，不執行 DDL。
- 基本單設備列保存 record/group/device identity、UTC bucket 及逐 member provenance。跨設備群組保留各 member 來源，沒有猜測的單一 device_id。
- production boundary 使用 applied version 的凍結 layout；ACK journal、payload、record/effect key、digest 與 receipts 仍由既有 delivery pipeline 處理。
- SQLite uint64 使用 TEXT；PostgreSQL 使用 NUMERIC(20,0)。PG inspection 保留 precision/scale，拒絕 NUMERIC(20,2)、NUMERIC(18,0)，保留 codec 已證明安全的 unbounded NUMERIC。
- schema preview／metadata 只讀；missing SQLite file 不會先建立。managed 檢查 actual internal SQLite main/attached 路徑、symlink/hardlink 身分及可驗證的目的地；只接受 main schema。
- 薄 API 只接受保存的 revisions 與 server token/operation identity，拒絕 client SQL、表名或 dialect。執行前後再核對 group/source/connector/layout scope。
- 只建立明確確認的缺少 owned structures；既有表須完整型別、nullable、PK 及 owner marker 相容才 no-op。沒有 ALTER、DROP、rename 或資料搬移流程。
- managed readiness 同時檢查 canonical data table、receipt marker 與已驗證 operation proof；確認建表、套用及 SQL committed 分別呈現。
- 舊 custom/all-good 表不增加 metadata／receipt，不改既有 accepted payload。舊 managed contract 沒有新 metadata bindings 時維持原路徑。
- rollout／rollback 邊界記於 `docs/technical/studio-v2-write-groups.md`；回退保留外部資料、receipts、journal/outbox 與 operation 證據。

## 需求、測試與執行

| Requirement | 實作位置 | Scenario 測試 | Example 測試 | 執行證據 |
| --- | --- | --- | --- | --- |
| Canonical group managed preparation | workspace `write_group_schema.go`、API group schema adapter、`GroupSchemaPanel` | EmptyNoPlan、internal destination、readonly metadata／preview | 無 | 真 SQLite／PG API，正常 embedded SQLite UI→SQL；passed-current |
| Managed preparation preserves operation and table ownership | recordingplan `schema_group.go`、既有 ledger、group metadata handler | restart/replay/scope/conflict/permission/metadata；custom pipeline | 無 | 真兩 DB、scope 反例、Partial PG 故障注入；passed-current |
| Revision-bound schema preview and explicit creation | `schema_preview.go`、`schema_apply.go`、group schema adapter | expired/legacy/tampered/source/layout、public guards、activation gate | 無 | 全 Go suite、兩 DB group API、既有 schema gate；passed-current |
| Basic managed SQL preserves record time and origin identity | runtime `group_boundary_layout.go`、`GroupRowLayout`、saved row policy | delayed restart、partial/bad/no_data、cross-device | 無 | 真兩 DB SQL/readback/receipt；passed-current |
| Stable exact managed column generation | managed save、exact layout、PG column inspection、frontend proposal/parser | same labels、uint64、rename、ack loss | Unsigned counter boundaries 全三列 | 真兩 DB exact SQL tests、前端 proposal tests；passed-current |

以下 22 個 scenarios 全分類為有測試；沒有 test-scope exclusions。

| Scenario | 對應斷言 |
| --- | --- |
| Empty managed destination | `TestManagedGroupPreparationEmptyNoPlan`：兩 DB 只有 canonical group、0 RecordingPlans，確認後 data/receipt 表及 readiness；UI 不手寫 SQL／target columns |
| SQLite preview remains read-only | EmptyNoPlan、metadata missing-file router/handler tests；實際 UI reload 與 preview 後 destination file 均不存在 |
| Internal database is not a recording destination | `TestManagedGroupSchemaRejectsInternalDestinationBeforePreview`、metadata internal router test、managed destination alias tests；實際 UI 422、安全提示、0 recording tables |
| Duplicate or lost confirmation response | `TestManagedGroupSchemaRestartReconcilesRealLostDDLReply`、completed replay tests；0/1/2 個實際 tables，重開 config DB、過期 token 僅 reconcile；前端同 operation lookup 回歸 |
| Scope changes before confirmation | `TestManagedGroupSchemaRejectsChangedScopeBeforeMutation`、GroupBindsSourceSchemaLayoutAndDigest；兩 DB 各 12 反例，0 target effects |
| Existing table conflict or insufficient privilege | conflict／PG no-CREATE tests；既有 sentinel 保留，permission denied 有 durable operation，無 schema proof |
| Existing custom table remains unchanged | `TestProductionEntityRowsSQLAndFrozenRestart`、既有 custom SQL tests；新 metadata bindings 為空時保留 payload/layout，未執行 ALTER |
| Metadata belongs to the saved group table | UsesSavedGroupTable、RejectsUnprovenGroupScope、frontend saved/dirty table tests：group B 的 persisted table，不使用 connector default A 或 draft table |
| Target changed after preview | changed-scope connector/destination/group cases；既有 SchemaApplyRejectsTargetChangedAfterPreview |
| Preview token expired or belongs elsewhere | RejectsExpiredForeignUnknownLegacyAndTampered；expired group token 沒有既有 operation 時不授權 DDL，有 operation 時僅 evidence reconcile |
| Missing unknown or legacy confirmation | changed-scope missing-confirmation/unknown-token、既有 UnsafeConfirmations、legacy token unit；拒絕前無外部 effects |
| Another session changes the setup | changed-scope workspace case：409，0 target tables／missing SQLite file 仍不存在 |
| Client alters statements | strict DTO 的 client-statements/table/dialect cases：400，stored server statements 才可執行 |
| Compatible existing schema | EmptyNoPlan owned no-op、CanonicalGroupCompatibleOwnedTablesAreNoOp；實際 inspection，executed statements=0 |
| Existing incompatible table | real conflict tests、schema drift：owner/receipt marker 被移除或 PG NUMERIC scale 變更時 readiness/preview 拒絕 |
| Legacy endpoint bypass attempt | GenerateSchemaNeedsConfirmation、GenerateSchema_DryRunAndExecute、SchemaApplyRejectsUnsafeConfirmationsWithoutMutation；共用既有 public safeguards |
| Activation without explicit schema confirmation | ReportsSchemaPreparationWithoutActivating、LocalModbusOnlyNeedsNoDatabase；managed proof 未確認時 UI Apply 停用 |
| Canonical group or layout changes after preview | GroupBindsSourceSchemaLayoutAndDigest、changed-scope tests、completed replay source/connector 反例 |
| Delayed delivery retains original time | `TestProductionManagedRowsKeepTimeIdentityAndExactTypesAfterRestart`：accepted row 後改 display、停機/reopen/offline destination、SQL 保留原 bucket/observed_at/identity/digest |
| Missing and bad values are not fabricated | `TestProductionManagedPartialRowsKeepNullAndQualityReasons`：NULL+missing/bad reason；required bad=skipped、silent=no_data，沒有零／假 good |
| Same labels and exact uint64 values | 同名、中文、保留字 labels 仍生成不同欄；三個 uint64 各自在兩 DB 逐字 readback；UI 不提案 unsafe BIGINT |
| Replay and display rename | frozen mapping/payload/record/effect/digest、commit acknowledgement loss，remote receipt reconcile、duplicate delivery 沒有第二筆 effect |

`Unsigned counter boundaries`：`9007199254740993`、`9223372036854775808`、`18446744073709551615` 在同一真 SQL matrix 對 SQLite TEXT／PG NUMERIC(20,0) 每列逐字斷言。

design 的「七型別」依實際既有 type contracts 驗證：10 種來源 Tag types 全覆蓋；既有 exact value codec 的 text/bool/int64/uint64/decimal/float64 及 NULL 路徑保留。沒有新增 decimal Tag 產品能力。

## 先失敗、再修補

- 新 group-scoped schema API：起初 404；新增 adapter/route 後 EmptyNoPlan GREEN。
- managed SQL metadata bindings 缺少時，移除 GroupID binding 的 Go overlay 使真兩 DB runtime SQL 失敗；修正後原時間、身份與 exact readback 通過。
- completed replay 在 source/connector 改變後錯誤回 200；修正後 409，歷史 operation 仍可查。
- owner/receipt marker drift 與 PG NUMERIC(20,2) 先錯誤 Ready；修正後 readiness/preview fail closed。
- PG inspection 原本把四種不同 NUMERIC 全讀成 numeric；保留 precision/scale 後 bounded type 回歸通過。
- embedded UI 200 preview token 的空 `table_prefix` 原本解析失敗；完整 group layout token 回歸先失敗，修正後可顯示並明確確認。
- 真 UI 找到 metadata GET 先建立 SQLite file：HEAD handler overlay 的 generic/group GET 回歸實際失敗，readonly inspector 修正後不建檔，internal destination 回 422。
- managed Update 留白表名先回 422；兩 DB API 回歸先失敗，保存時解析穩定群組表名後 200，0 DDL。
- succeeded 無 reason 原本顯示「需要處理」；真 response shape 回歸修正成功結果文案。
- connector 在驗證與 DDL 之間改變時，原本會打開替換的 SQLite；回歸先證實副作用，再改為執行已驗證且 revision-bound 的 immutable connector snapshot。
- inspection 原本重新讀 connector、metadata 原本混用前後 group scope；非 main／切換目的地／metadata revision 變動回歸先失敗，修正後安全拒絕且不建檔。
- caller 修改 persisted token layout／statements、移除 group layout 或污染 nested columns，原本可繞過保存內容；回歸先失敗，再使用 repository 的完整 token 與深層複製。
- managed readiness 現在使用同一 destination guard；internal alias／opaque file URI、缺少 runtime frozen descriptor／digest／owner proof 的回歸通過，既有 accepted journal 保留而不重建未知 layout。
- confirmation 回覆遺失後再改同群組 scope，原本會丟掉 operation lookup；before/after reply 回歸修正後保留同 operation 查詢。permission_denied 亦以安全、在地化文案呈現。
- 日誌位於 `/private/tmp/gw-D-*-red.log` 與對應 green/final logs。overlay 僅在暫存目錄，不覆寫 worktree。

## 真實 embedded UI 與獨立 SQL

- 正常 `cmd/test_ui` build（無 fixture build tag）提供 `/studio/v2`；production frontend build/static sync，agent-browser／Chrome for Testing 149。沒有 network mocks 或直接插入目的地資料列。
- 自有 loopback Modbus A/B（15030/15031），沿用 C 已驗證的 source configuration 複本；只重置複本內群組與 connector，保留 sources，重新走實際 group API/UI。沒有正式 PLC、DB、LAN 或部署操作。
- SQLite UI 保存四個 members、留白表名及 server-generated target columns；preview 後仍沒有 destination file。明確確認後獨立 SQL 看到恰好 data table＋receipt table。
- 群組套用後收集真 simulator samples。獨立 readback：A=`215/1013`、B=`187/777`，group/record identity、UTC bucket、四點 observed_at/good provenance 與 receipt 數量一致。
- UI 區分 saved/applied/collecting/local pending/SQL confirmed；正常 binary 冷啟動後持續交付，獨立 SQLite snapshot 有 30 筆 rows／30 筆 receipts，四個來源值及 provenance 一致，沒有固定成功資料。
- 最新 bundle 對同一 compatible owned tables 再 preview／明確確認為 no-op succeeded；operation `op-01260210-7799-4a0d-a0c2-c08cd127e32a`，applied frozen group revision 未變。
- 錯誤情境：可丟棄 configuration DB 同時被指定為 managed destination。UI preview 得到實際 422，顯示「請選擇獨立的記錄儲存檔案供受管儲存使用」，data/receipt tables 未新增。
- 部分成功情境：自有 PG namespace 的限定 event trigger 在 CREATE TABLE transaction 中移除新 receipt table。最新 operation `op-a3645904-20d8-4cc5-850f-eef7d3b01a8f` 實際回 partial/verification_mismatch、executed statements=1；只有 owned data table 存在、沒有 proof、Apply 停用。查詢維持同 operation ID，table 數未改；不是 mock response。
- 鍵盤 Enter 實測 Save、preview、confirm、Apply、operation lookup；390px schema table region 可 focus，ArrowRight 的 scrollLeft=40，tabIndex=0。
- 最新正常 binary／bundle 的三寬度主流程／錯誤／部分完成截圖在 `evidence-d/`，九張均已目視比對。390px 使用既有側欄及 connector 收合按鈕，main document width=390，表格維持自身 scroll container；768/1440px partial document width 與 viewport 一致。
- 既有 connector summary 的長 schema/table 膠囊在 partial fixture 390px 造成 document width=435；新 schema/member tables 留在各自 scroll containers。這是未修改的既有摘要限制，獨立列為範圍外觀察，未據此增加本案修復範圍。
- 沒有額外 Design Source；依 design 的 canonical group、明確確認、真 operation states、逐 member facts、custom 相容邊界比對，沿用現有 editor/card/table/keyboard patterns。

## 執行檢查

darwin/arm64、Go 1.27.1（module 1.25.5）、golangci-lint 2.14.0；自有 PG 16.14、loopback 55432／gwtest。既有 5432 服務未修改。

| 命令／範圍 | 結果與證據 |
| --- | --- |
| `POSTGRES_DSN= go test -p 1 ./... -count=1` | passed-current；`gw-D-go-latest-test.log` |
| `go vet ./...` | passed-current；`gw-D-go-latest-vet.log` |
| `golangci-lint run ./...` | passed-current，0 issues；`gw-D-go-latest-lint.log` |
| `go test -race ./internal/datalink/dbtarget ./internal/datalink/recordingplan ./internal/datalink/workspace ./internal/datalink/grouppipeline ./internal/api ./cmd/test_ui -run 'TestExact\|TestAtomicRow\|TestManaged\|TestApplySchemaPreview_Group\|TestMemoryRepositoryPreviewToken\|TestProductionManaged\|TestProductionEntityRowsSQLAndFrozenRestart\|TestProductionGroupOutageRecoveryPostgresDestinationReceiptsAndRestart\|TestCanonicalManaged' -count=1`，自有 POSTGRES_DSN | passed-current；六個 packages 均有實際 matching tests，真 SQLite／PG；`gw-D-sql-latest-race.log` |
| `cd frontend && npm run lint` | passed-current；`gw-D-frontend-latest-lint.log` |
| `cd frontend && npm test -- --run` | passed-current；185 files／1133 tests；`gw-D-frontend-latest-test.log` |
| `cd frontend && npm run build`、static sync、正常 Go binary build | passed-current；`gw-D-frontend-latest-build.log`、`gw-D-binary-reviewed-build.log`；此產物用於九張截圖 |
| `git diff --check`、`make check-lines`、Spectra scope/validate | 文件與封存後另核對；不以舊檢查代替最後差異 |

本機 simulator＋可丟棄 SQL／畫面證據不代替現場或部署驗收。未重跑整套帶 POSTGRES_DSN 的 repo suite；A 已記錄的歷史 recording-plan managed/partition migration 失敗未納入成功聲明。未執行完整 Playwright e2e suite；實際主線/error/partial 由正常 embedded UI 見證。

## 最後驗證與 review

已完成上述需求、design、production wiring、真 SQL、embedded UI 與修補後的最低檢查。最後 scope 及 task attribution 由 CLI 保存；archive 前再檢查 snapshot 穩定性與全部任務完成。

Review 涵蓋正確性、效率、重用、慣例及安全邊界；已修正本案發現的缺陷。local operation ledger 與外部 DDL 無法形成跨資料庫原子交易，這是 design 已明列的限制。確認被接受後另一次設定變更不會重新導向新 connector；若執行後 scope 已變，回傳 Unknown、不保存 schema proof。跨服務的全域設定鎖不在本案契約內，未據此新增 migration／locking engine。
