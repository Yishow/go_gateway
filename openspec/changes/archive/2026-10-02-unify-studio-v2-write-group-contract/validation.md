# 實作驗證紀錄

2026-10-02：本輪已授權依總覽的 A→F 順序實作與驗證。目標仍是完整六案及舊案移交驗收；本文件記錄實際執行，不以提案齊全表示產品完成。

## 基準與範圍

- 起始 HEAD：`3b09292afc858ba5b96668389ccbebe3aad8c363`；工作目錄乾淨，live remote main 同 SHA。
- 平台：macOS arm64；本機 Go 1.27.1、Node 24.15.0；專案 Go 語言版本以 `go.mod` 為準。
- 按 A 的任務順序執行。未完成 A 的必要驗證前，不推進 B；未授權 commit、push 或部署。
- 測試只使用可丟棄資料庫；沒有執行現有 gateway 設定資料庫的 migration 或外部試寫。

## 修改前前端基準

| 命令 | 實際結果 |
| --- | --- |
| `cd frontend && npm run lint` | PASS |
| `cd frontend && npm run build` | PASS |
| `cd frontend && npm test -- --run` | FAIL：172 個檔案，974 項通過、1 項失敗；失敗為 `database-autosave-page.row-groups.test.tsx` 的 delayed row-group save 案例 |
| `cd frontend && npm test -- --run tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx --reporter=dot` | PASS：3/3；完整執行的失敗尚不能視為已解決 |

以上前端檢查在修改前端來源之前執行；不得稱為最後實作版本的驗證。後續群組相容／UI 批次需留意此不穩定案例。

## 尚未完成的產品驗收

- A1–A3 已完成下列本地設定、API、移轉及交接介面驗證；A3.3 的最後規格／程式核對依本文件結果收尾。domain applied metadata 不是 production writer proof。
- B–F 尚待依序實作。C 接實際 intake／owner／journal／sender；D 接 confirmed test-write；E 接群組 UI；F 做 embedded production UI→SQL。
- 本機 Docker daemon 未啟動；live PostgreSQL、Linux production UI→SQL、Windows／ARM、LAN／真 PLC／SCADA、部署與現場驗收沒有通過證據。

## A1.1 保存與重載

- Worker 記錄行為 RED：可編譯的 repository stub 回傳 `write-group repository is not implemented`，新增保存測試因此失敗；不是缺少 symbol 或環境錯誤。
- 主代理對交回版本執行 `go test ./internal/datalink/workspace ./internal/datalink -run 'BasicGroupWithoutMeasurement|WriteGroup|WriteGroupsMigration' -count=1`：兩個 package PASS。
- 主代理執行 `go vet ./internal/datalink/workspace ./internal/datalink`：PASS；`git diff --check`：PASS。
- 已回讀 source，確認 server ID/revision、draft、不改 applied revision、optional measurement、來源與 connector 綁定，以及 caller transaction seam。
- 審查未發現已證實功能錯誤；已補過期 connector/source/mapping 版本、跨 workspace Create/Get/List、同 workspace 跨 device/point、disabled mapping、missing tag 及 partial migration 測試。
- 主代理對最後交回版本執行 `go test ./internal/datalink/workspace ./internal/datalink -count=1`：兩個 package 全部 PASS。
- 主代理執行 `go test -race ./internal/datalink/workspace -count=1`：PASS。
- 尚無 PostgreSQL live repository、HTTP、CAS、readiness、lifecycle 或新 delivery 的通過證據；後續任務各自驗證。

## A1.2 原子保存與 API

- Service RED：可編譯 stub 的 `TestWriteGroupServiceCreatePersistsDraftAndProjection` 回傳 unavailable，測試失敗。
- 正式 router RED：SQLite fixture 成功初始化後，`TestNewRouter_AtomicGroupSaveAndCAS` 的 POST 回 404，預期 201；新增 route、service injection 及保存行為後轉 GREEN。
- 整合補測 RED：設備已刪除但 point 尚存在時，POST 回 500，預期安全 404；補 source join 的 `sql.ErrNoRows` 分類後轉 GREEN。
- 主代理 API focused suite：25 個測試事件全部 PASS；包括第二段 workspace save trigger 故障的 group/member/projection/revision 全回滾、兩 client 同版本最多一方提交、三層 CAS、Create/Delete live connector stale、安全 resource 404、server readonly fields、tombstone 保存及 reload。
- 主代理對最後 source 執行 `go test -race ./internal/api ./internal/datalink/workspace -run 'WriteGroup|AtomicGroupSaveAndCAS' -count=1`：兩個 package PASS。
- 前端服務測試：9/9 PASS；初次 draft 可省略 source/mapping revisions，canonical response 必須含 server revisions；Create/Update 剝除 server fields/proofs，保留 caller draft，malformed 200 拒絕。
- 已修正並測試 legacy reference 保留、schema binding 變更清除舊 proof、deleted group 不可復活。尚未開放 applied group delete；A1.4 負責 lifecycle/backlog 保護。
- Production API 已注入主程式共用的 SQLite handle；只保存本地 draft，不新增外部 DDL、probe 或 delivery。
- Swagger 由專案鎖定的 `github.com/swaggo/swag/cmd/swag@v1.16.6` 重產，工具與 cache 使用 `/private/tmp`；未手改 generated Swagger，也未改 `go.mod`／`go.sum`。

### 完整檢查與基準差異

- `go test ./... -count=1`：FAIL，只有既有 `modbusshare/TestService_Status_AfterStartAndStop` 無法綁定固定 5020；單獨相同測試重跑 PASS。此結果與並行埠競爭相符，未修改 listener 測試。
- `go test -p 1 ./... -count=1`：PASS，44 個有測試 package；6 個 package 沒有測試。使用套件循序執行避開固定埠競爭。
- `go vet ./...`：PASS；最後 source 修正後，受影響 workspace/API/cmd 再跑 PASS。
- `golangci-lint run ./...`：FAIL，只剩未修改 `internal/datalink/dbtarget/service_probe.go:201` 的 gocritic 提示。另以起始 SHA `3b09292` 的完整 snapshot 執行同 lint，得到完全相同一項；本次新增程式問題已修正，未擴大修復範圍。
- 前端 `npm run lint`、`npm run build`：PASS；`npm test -- --run`：173 個檔案、984 項 PASS。最後僅增強原服務測試的 assertions，再 focused 9/9 PASS。
- `git diff --check`、`make check-lines`：PASS；有超過 300 行警告，沒有超過 500 行阻擋。
- Readiness、Apply/bucket lifecycle、相容移轉、writer owner、durable delivery 與 Linux/PostgreSQL UI→SQL 驗收仍待後續任務；A1.2 的 local/API 證據不能代替它們。

## A1.3 群組 readiness

- 主代理 workspace 行為 RED：managed/custom 都將 canonical member 誤列為 `database-target-missing`；預期只留下未遷移點位。接 canonical coverage 後，同一測試轉 GREEN，既有 `TestService_Readiness` cases 同時 PASS。
- 主代理 router 行為 RED：群組已保存後 GET `/:id/readiness` 回 404，預期 200；接 scoped route 後轉 GREEN。缺少 MeasurementDefinition、零額外 target mapping 的基本群組可得到 `config_ready=true`，未驗 schema 仍 `schema_ready=false`、`ready=false`。
- Workspace summary 保留 Step 2/3 ownership 檢查，canonical coverage 以 connector/point/tag tuple 比對；未 Apply 明列 `write-group-apply-required`，不將配置 readiness 說成 runtime running。
- 前端 worker RED/GREEN：先用 unavailable stub 驗證 GET service；另以矛盾 `ready=true` 加 blocking issue 的回應證明 parser 必須拒絕。focused 16/16 PASS。
- 主代理對最後前端來源執行 `npm run lint`、`npm run build`：PASS；`npm test -- --run`：173 個檔案、991 項 PASS。
- 最後前端窄幅邊界修正後，worker focused 18/18、lint、build PASS；前述全 suite 是該窄幅修正前的版本。
- 真實可丟棄 SQLite + `ReadOnlyTableInspector` 已驗證 schema/columns；未配置 measurement 或額外 target mapping 仍可 config/schema ready，group applied revision、schema proof 與 membership 不變。missing file 不建檔，DSN/path priority 與 saved scope 一致。
- 精確 scope/column 與 SQL 型別 allowlist fail closed；型別收窄 RED 曾證明 float→INTEGER、int64→REAL、任意 string→UUID/JSON 被誤認，相同 focused cases 修正後 PASS。未知 SQL 型別不再用 substring 猜測。
- Production wiring 使用唯讀 inspector 與 canonical readiness reader；PostgreSQL catalog inspection 使用同一 saved connector snapshot，未實際連線 live PostgreSQL。
- 主代理最後 `go test -p 1 ./... -count=1`：44 個測試 package PASS、6 個無測試；`go vet ./...`：PASS；readiness/readonly/CAS focused `go test -race`：workspace/API/dbtarget 三個 package PASS。
- 主代理最後全 repo lint：FAIL，只有前述相同 baseline `service_probe.go:201`；`git diff --check`、`make check-lines`：PASS。
- Swagger 已由鎖定工具重產 readiness route/model。Readiness 不套用版本、不建表、不表示 SQL committed；A1.4 與後續 B→F 仍未完成。

## A1.4 Lifecycle 與 backlog

- Router 行為 RED：POST `/:id/disable` 回 404，預期 200；接 route 後同測試 GREEN。四層身分／CAS、stale 409、缺欄 400、安全 404、不可復活 tombstone，以及 reload 保留 membership/applied revision 已驗證。
- Workspace 行為 RED：disabled canonical scope 被重新要求 legacy target mapping；修正 coverage 後，只剩真正未移轉的點位要求 mapping，disabled group 不要求 Apply。
- Worker 行為 RED/GREEN：60 秒 interval 降至 15 秒時，原實作提早切換；改用安全的共同 UTC 邊界，舊 bucket 保持原 snapshot 至 12:02:00。
- 主代理新增行為 RED：6→10→5 秒連續 Apply 可讓較早排程稍後覆蓋新 policy；修正後，同時只允許一個待生效的 semantic version，拒絕時不新增 snapshot、不改 applied identity。metadata-only Apply 仍保留 writer identity。
- 主代理審查新增行為 RED：`time.Truncate` 的 year-one 基準使 7 秒／7→11 秒共同邊界偏離契約的 Unix epoch bucket；改為 epoch 對齊後，含負時間、時區與次秒邊界的 lifecycle focused suite PASS，workspace vet PASS。此窄幅修正晚於下列完整 Go 測試。
- applied group delete 預設 fail closed；安裝 transaction-scoped backlog ownership guard 才能 tombstone。可丟棄 SQLite fixture 保留 immutable version、舊 revision payload 與 target identity，既有 delivery worker 完成 local handoff 並保存 local receipt。
- `write_group_versions` 使用 `ON DELETE RESTRICT`；migration repair 補缺表而保留既有群組及成員。
- 主代理 `go test ./internal/datalink/workspace ./internal/api ./internal/datalink -run 'GroupLifecycle|WriteGroup|WriteGroups' -count=1`：三個 package PASS；workspace/API 同範圍 `go test -race`：兩個 package PASS；受影響三個 package `go vet`：PASS。
- 主代理 lifecycle 最終版 `go test -p 1 ./... -count=1`：PASS；受影響後端 lint 仍只有既有 `service_probe.go:201` 的同一項 gocritic，未新增 lint 問題。
- 前端 disable service 可編譯 unavailable stub 行為 RED；worker focused 21/21、lint、build PASS。UI 尚未修改；此批沒有再次執行全 frontend suite。
- Swagger 已重產 disable route；`make check-lines` 通過。Apply domain seam 尚未接 production route/settings/readiness token/owner 切換；本批 local receipt 不等於外部 SQL commit 或 readback。

## A2.1 Single mapping：已確認 snapshot 轉換並完成

- 實際 production writer：`writer.go:123-130` 只有 `GroupKey` 進 grouped buffer；無 GroupKey 每筆立即交給 `writeMapping`。`writer_statements.go:74-120` 的 insert 不保存 timestamp，upsert 使用每筆 observed time，不因 interval 設定成為 snapshot。
- 新 basic 契約每個 UTC interval 只選一列 snapshot；把 legacy interval 直接填入新 row policy 會改變列數、時間及 late/retry 語意，不能稱為無損等價移轉。具 GroupKey 者仍屬 A2.2，不能挑第一筆當 single mapping。
- 使用者於 2026-10-02 回覆「依建議」，確認審閱逐筆→snapshot 差異並顯式確認才保存草稿；保留原設定及舊 writer 到有效 Apply。已以 spectra-ingest 更新 proposal/design/spec/tasks，保留完成項目與 provenance；analyze 沒有 Critical/Warning（仍有具體例子的 Suggestion），validate PASS。
- 已完成不依賴答案的唯讀差異預覽與服務契約。主代理 router 行為 RED：真實 SQLite legacy fixture 的 preview POST 回 404，預期 200；同測試在正式 route/domain 接線後 GREEN。
- Review 保存、stable identity map、canonical legacy reads 已完成；A2.2/A2.3/A3 仍待實作，B→F 尚未解鎖。此處 review 不等於 runtime Apply。

- 正式 review route 行為 RED：真實 SQLite fixture 的 preview 回 200，review POST 回 404、預期 200；接線後轉 GREEN（見下方實測）。

### 唯讀預覽實作與實際驗證

- 一次本地 read transaction 取得 workspace、connector、selected legacy mappings 與完整 source chain；不選第一筆 active source、不把 disabled unrelated target 當競爭者。歧義、shared-column、advanced metadata、upsert 或 unknown write mode 標示 blocked。
- before intent 保存原 timestamp/group key/interval，不將空白 trim 成新保存值；candidate 沒有 persisted ID/revision/applied identity，不建表、不切 runtime。Unknown/foreign resource 回安全 404；缺欄 400；錯誤不回傳 SQL 或 connector credential。
- 主代理新增行為 RED/GREEN：JSON findings 的 nil/array 不一致、connector default interval 未納 source revision、只檢查第一個 potential collision、unknown mode 被猜成 append；各原測試在最小修正後通過。
- SDK 行為 RED/GREEN：blocked 空 source revision／invalid legacy intent 原值被誤拒，以及 needs_review/blocked 矛盾狀態被誤接受；修正後 focused 32/32 PASS，caller payload 不變、只送 workspace/source IDs。
- 主代理完整 `go test -p 1 ./... -count=1`：44 個測試 package PASS、6 個無測試；`go vet ./...`：PASS；preview/lifecycle focused race：workspace/API 兩個 package PASS。
- 主代理全前端 `npm test -- --run`：175 個檔案、1007 項 PASS。隨後只移除測試中未使用的 type import，最終 `npm run lint`、`npm run build` PASS。
- 主代理完整 backend lint 仍只報起始版本既有 `service_probe.go:201` 的 gocritic；未將它算成本次通過。
- 最後審查補測 RED：missing Tag 的 before intent 丟失原綁定 ID；saved connector schema 與 legacy schema 不符時，preview 提議改 destination。修正為保留原 Tag ID／block scope mismatch；晚於前述全 repo Go gate，最終受影響兩個 package 檢查另列下方。
- 跨 workspace 多來源先使用了會 AttachDevice 的 helper，該失敗不算有效 RED；改成真正 foreign fixture，再對修正前程式確認行為 RED，修正後維持安全 not-found，不回傳 foreign 的歧義預覽。
- 修正 scope gate 時，原正向 fixture 的 connector schema 與 target 不一致，已改為一致的 SQLite `main` scope。最後版本 `go test ./internal/datalink/workspace ./internal/api -count=1`、兩個 package preview/lifecycle focused race、vet 與 scoped lint：PASS；scoped lint 0 issues。
- Swagger 由鎖定版本工具重產 preview route/DTO。所有 SQL fixture 都是本機可丟棄 SQLite；PostgreSQL scope 案例只驗設定解析，沒有 live PostgreSQL probe/write。

### 明確確認保存與相容讀取

- 同一 transaction 重算 preview、檢查 workspace/connector revisions 與 digest、建立 draft/projection/持久化來源 ID map，保存原始 before intent。缺欄 400、未確認或 blocked 422、unknown/foreign 404、stale 409；map 或 workspace 第二段保存故障全回滾。
- 相同來源 revision 取得新 preview 後重跑回同群組 ID，不增 workspace revision、不覆寫 canonical edit；舊 request 因 workspace stale 回 409。已 tombstone 的 ID 不重建；來源改變或 connector ID/revision 不符拒絕重跑。
- Preview 明列最大 observed time、穩定 sample ID、late、skip_row/freshness 與 timestamp 差異。原 interval=0 保存於 before intent，候選明示 connector fallback 或設計預設 15 秒，max age 同 interval；沒有保存後才默默補預設。
- Review 正式 route 的 404 RED 已轉 GREEN；真實 API fixture 完成 preview→review→canonical edit→global mapping GET／workspace target GET。保存不改 applied revision、不開啟外部 DB，原 legacy mapping 的 timestamp/column/enabled 保留。
- 舊讀取只在 HTTP handler 投影，raw MappingService/runtime/EnsureWorkspaceSchema 未改。回應附 canonical_group 與仍由 legacy writer 使用的 legacy_intent；後者不是 map 內凍結的 review snapshot。canonical disable 的 filter 在投影後正確生效。
- 讀取投影 RED/GREEN：過時 column、workspace target 遺漏已移轉 point、多 member 的錯誤修復動作、無關 connector filter 被卡住，以及 canonical connector 變更被錯當可相容。最終不能表達時回 `WRITE_GROUP_LEGACY_READ_CONFLICT`／`open_write_groups`；相關 canonical API 仍可讀，無關 scope 不受阻擋。
- 補 map-only upgrade fixture：既有 group/member/version 表均存在、map 表尚未有時，真實 migrator 建立 ID map 而保留群組；再次執行保留 map。此補測是 GREEN 回歸，不另宣稱有獨立 RED。
- SDK 首次 undefined seam 失敗不是有效 RED；worker 另用可編譯 unavailable stub 取得有效行為 RED，還原後 focused 20/20 PASS，tsc/scope lint PASS。SDK 僅傳明確確認與 revision/digest/來源 IDs，拒絕 unsaved/foreign/空或重複群組結果，保留安全錯誤資訊。
- 完整 Go 初次只失敗既有 Share production gate 的臨時埠占用，單獨重跑 PASS；`go test -p 1 -parallel 1 ./... -count=1` PASS：44 個有測試 package、6 個無測試。此完整結果在最後 connector read guard 前；最後版本 API/handlers/workspace/datalink 四個完整 package、migration focused race 兩 package、受影響 vet 與 scoped lint PASS（0 issues）。
- `go vet ./...` PASS；完整 backend lint 仍只剩起始版本相同 `service_probe.go:201` gocritic，未擴大修復。最終 frontend lint/build PASS；完整 frontend：176 檔案、1016 項 PASS。
- Swagger 已由鎖定版本重產 review route/DTO；`git diff --check`、`make check-lines` PASS。PostgreSQL live、production owner switch、新 delivery 及 UI→SQL 尚未驗收。

## A2.2 Row-group：已完成

- 使用者已授權 A→F；本批技術表示補入 A2.2 Implementation Contract 與具體 shared-column fixture。analyze 無 Critical/Warning，18 項 Suggestion 是其他抽象 scenarios 缺具體 examples；validate PASS。
- 真實 legacy writer 的 `TestWriter_InsertModeRowGroupsEmitSeparateRowsForSharedColumn` 與 `TestWriter_GroupedMappingsFlushOneRowPerBucketOnTimer` 本輪實跑 PASS。不同 persisted GroupKey 可在同 column 分列，同 GroupKey 的不同 column 會合列；GroupKey 只作 buffer partition，不是 SQL business value。
- 新 `member.entity_key` 保存 row partition；row policy 保留 group/unique-key 規劃 metadata，migration 保存原 row-group ID 及 target mapping ID→point ID。這些 metadata 不代表外部 SQL unique constraint 已驗證。原 target mapping、enabled 與 writer 保持原狀；review 只存 draft。
- API 有效 RED：正向 fixture preview 回 404、預期 200。preview／confirmed review 接線後，四個整合測試 PASS，涵蓋 shared column 分列、stable provenance、CAS／explicit confirmation、blocked／safe404、交易 rollback、重跑不覆寫 canonical edit，以及舊讀取反映 canonical 成員。
- 正式 handler 在舊讀取中逐 member 投影；removed／新增無對應 member、connector 或 SQL business binding 不可表示時，回 `WRITE_GROUP_LEGACY_READ_CONFLICT`／`open_write_groups`。原 MappingService/runtime/EnsureWorkspaceSchema 未接此 HTTP adapter。
- 審查新有效 RED：刪除原 `legacy-B` 後，兩個 List endpoint 都回 200 並少列成員、預期 409。補足全部來源存在性檢查後 GREEN；無關 connector／仍完整的 Tag filter 正常回覆，缺失 Tag scope 回 conflict。
- SDK 以可編譯 unavailable stub 取得有效行為 RED（5 項失敗）；還原後 focused 5 檔案、49 項 PASS。最後審查確認 blocked intent 原本已可保存空 DeviceID，另以有效 RED 證明 needs_review 也錯誤接受空來源；限定放寬只用於 blocked 後，focused 50/50、tsc／scoped lint PASS。主代理最後完整 frontend：177 檔案、1025 項 PASS；lint／build PASS。保留 false confirmation，不自動 Apply。
- Swagger 由鎖定工具重產；核對沒有移除原 API path，新欄位與 row-group preview/review paths 存在。重產後 API 移轉測試 PASS。
- domain 保留原 row-group/target scope、完整來源 metadata revision、跨批次來源 ownership 與舊 schema upgrade；同一 legacy target 不可同時由 single/row-group adapter 建立另一份 canonical authority。原 row-group scope、member/key/ref 與 connector default interval 納入 digest；重啟後 ID map 與 draft 保留。
- 主代理新增有效行為 RED：同一 Tag 已有另一個 point 的 enabled mapping 時，preview 回 needs_review、預期 blocked。改為檢查同 Tag 的全部 enabled source 後 GREEN，review 拒絕且 group 數為零。
- 最後 `go test ./internal/datalink/workspace ./internal/api ./internal/datalink -count=1`：三個完整 package PASS；API/workspace 移轉 focused `go test -race` PASS。`go vet ./...` PASS；最後完整 backend lint FAIL，僅剩已確認起始版本相同的 `service_probe.go:201` gocritic。本次新增 handler 的重複字串 lint 已修正，並在最後 race 與 lint 重驗。
- 最後另跑 `go test -p 1 -parallel 1 ./... -count=1 -json`：44 個有測試 package PASS、6 個無測試，沒有 fail event。`git diff --check`、`make check-lines` PASS（300 行以上警告，沒有 500 行阻擋）；task 以明確 source/test 路徑記錄完成。A2.3/A3/B→F、live PostgreSQL、owner switch、新 delivery 及 UI→SQL 尚未驗收。

## A2.3 已確認的保留方案：已完成

- 現有 recording plan 尚未注入 production runtime，且 destination/member 缺少 canonical column／row identity 設定。`raw_history` 的 every_sample/on_change、triggered snapshot 與其他 advanced streams 均不能證明與週期 snapshot 等價；`sampled` 常數沒有對應 sampling 實作。
- 使用者於2026-10-02確認最小方案：保留等價要求，提供完整來源預覽與 blocked／修復提示並保存原 plan；新基本寫入使用 canonical WriteGroup。已以 spectra-ingest 更新 proposal/design/spec/tasks，保留已完成1.1–2.2及 provenance；analyze 無 Critical/Warning、18項 Suggestion，validate PASS。2.3 已完成下列實作與本機驗證。
- 主代理以真實 SQLite plan/source/connector fixture 取得 API 有效 RED：preview POST 回404、預期200；原 plan 每筆記錄、AppliedRevision 與原 target 都已存在，沒有利用 compile/setup 錯誤當成行為 RED。domain 與 SDK 正依同一安全 DTO 契約實作；沒有正向等價候選。
- SDK worker 的可編譯 unavailable stub 取得行為 RED；主代理審查並完成正式 error code/action 的安全 allowlist。原 plan 的 null collections、未知 mode、空 revision 與 unresolved source 原值保留；needs_review/candidate、錯誤 identity、prototype key 及任何 review success 均拒絕。最後 migration focused 38/38 PASS。
- 前端完整檢查首次有一項本次整合缺漏：新 allowlisted error code 尚無 en/zh-TW 訊息；補兩個訊息後，相關 12/12 PASS。最後 `npm run lint`、`npm run build` PASS，`npm test -- --run`：178 檔案、1034 項全部 PASS。歷史 locale 檔維持 511 行未增行；拆分計畫是在 E 涉及群組 UI 文字時，按既有 i18n namespace 模式分出該部分；本批未執行拆分。
- 同一 local transaction 讀完整 plan、referenced measurement/source/mapping 與各 saved destination connector；沒有 Validate/default 改原 payload。來源歧義、disabled/missing、空 revision、unsupported stream／raw policy、多 target、缺 column／row identity 均有明確 issues；缺來源的 sources 回空陣列，原 plan 的 null collections 保留。
- Domain 可編譯 unavailable stub 取得有效 RED；完整原 intent、unresolved measurement 內容變更、同 Tag 不同 point 歧義、目的地設定變更及 stale review 的回歸 PASS。另補 all advanced/on_change/sampled/batch/latest/unknown mode 與 foreign chain 的 GREEN 契約覆蓋；這批補測未另宣稱獨立 RED。
- 正式 API 的四個 `BasicPlanMigration` 整合測試 PASS；原 plan 的 running/status/applied revision、workspace revision、target enabled 保持原值，group/map/version 皆無新增，外部 SQLite 檔沒有建立。初次整合僅因主代理測試誤用 conflict code 失敗，改成既有 `revision_mismatch` 後通過，產品行為無需修正。
- 本批 `go test -p 1 -parallel 1 ./... -count=1 -json`：44 個測試 package PASS、6 個無測試、零 fail event；`go vet ./...` PASS。全 repo 結果在最後補契約測試前，產品 source 未再變動；新增 cases 的 focused tests/race 隨後 PASS。最後 API/workspace migration focused race PASS；domain 三個完整 package、scoped lint PASS。完整 backend lint FAIL，僅有已確認起始版本相同的 `service_probe.go:201` gocritic，未擴大修復。
- Swagger 已由鎖定 `swag@v1.16.6` 重產，recording-plan preview/review paths 及 intent DTO 存在，沒有移除起始 API path；go.mod/go.sum 不變。`git diff --check`、`make check-lines` PASS（歷史 locale 511→511）；沒有 UI editor、live PostgreSQL、owner switch、新 delivery 或外部 write/readback 的驗收。

## A3.1 修改前唯讀盤點

- 已核對 `DatabaseTargetHandler.CreateMapping/UpdateMapping/DeleteMapping`：直接呼叫 legacy MappingService，沒有 migrated-source 寫入 guard 或 workspace CAS。POST 沒有原 target ID，必須依 output scope/provenance 判定衝突，不能只比對 SourceIDs。
- Studio 的 `commitWorkspaceTargetSetup` 與 `commitWorkspaceConnectorSetup` 已使用 `UpdateDatabaseSetup`；target SQL／ref 與 connector／row groups 均在同一 setup-revision 交易中。guard 須在該 transaction 中讀取 ownership 並於 legacy 寫入前拒絕，handler 外的 read snapshot 不足以防止同時移轉。
- `writeGroupMigrationReader.List` 目前只供 operator reads，不是 transaction-scoped ownership checker。canonical Save/Review 的合法 projection 更新須持續可用；既有 single/row migration API fixture 與 typed conflict envelope 可重用。
- 此盤點由唯讀 reader 提供，主代理已回讀以上 transaction seams；沒有新增 A3 source、測試或完成勾選。A3.1 仍依賴 A2.3 完成，尚未執行 A3 gates。

## A3.1 舊寫入入口收斂：已完成

- 採既有設計允許的 actionable409，未移轉 scope 保留原 CRUD。全域 Create/Update/Delete、Studio target 保存共用 ownership guard；同一 workspace mutex 與 SQL transaction 內重查，先於實際 legacy persistence。主程式服務初始化及 HTTP handlers 都注入同一 coordinator；runtime 仍讀原 mapping repository，沒有切 writer。
- 已移轉 mapping ID、canonical current Tag/connector 及 durable migration map 的原 Tag/connector 都受保護，包含 disabled/tombstone、原 legacy row 不存在、canonical destination 已變更及不同 table/column。row-group 原 intent 的來源須唯一且完整覆蓋；missing/malformed provenance 拒絕寫入。
- Studio row groups 省略代表不替換；explicit-empty、移除 member、改 table/key/connector 均409。等值保存及 canonical edit 後保留目前 projection 可成功；不以 frozen before intent 誤擋。early check 與 transaction recheck 都先依本次 config 正規化候選 layout。
- 主代理取得有效行為 RED：移轉後 DELETE 仍200、預期409；空 row-group table/schema 繼承新 config 可改 layout 並回200、預期409；服務初始化後尚未建 router 的 owned Create 進外部 table inspection；ownership SQL read error 被誤報409、預期500。最小修正後相同案例與 before/after 保存資料 assertions PASS。測試期望誤將既有安全 generic retry action 當成應省略，已依既有 renderer 修正，未另改產品 action。
- ownership DB 操作錯誤保留原始 cause 並回安全500/internal、retryable=true；不顯示 SQL/DSN，不要求改群組。已實測 global 三個 mutation 與 Studio target/config 在缺少 map table 時均不寫入。missing persisted map 與 malformed provenance 仍是409，兩者不混用。
- 主代理最後 `go test -p 1 -parallel 1 ./... -count=1 -json`：44 個有測試 package PASS、6 個無測試、零 fail events；`go vet ./...` PASS。API/handlers/workspace/dbtarget/cmd 的相關 `go test -race` 五個 package PASS。完整 backend lint FAIL，僅剩起始 SHA 已確認相同的 `service_probe.go:201` gocritic。
- 前端 safe parser/localized error 的有效 RED 為 conflict code 被丟棄；補 shared allowlist 及 en/zh-TW 訊息、開啟群組動作後 focused31/31 PASS。最後完整 frontend：178 檔案、1035 項 PASS；lint/build PASS。沒有新 UI editor 或視覺驗收。
- `git diff --check`、`make check-lines` PASS；本次曾因 service.go 增行被阻擋，將此次修改的 CRUD 移至窄責任檔後，歷史 service.go 923→895；兩個 locale 511→511。沒有新行數豁免。A3.2/A3.3 與 B→F 尚未完成，沒有 production owner switch、live PostgreSQL 或 UI→SQL 的通過證據。

## A3.2 Apply 交易內的 writer 交接介面：已完成

- 可注入 barrier 在同一 local transaction 內讀到新的 immutable version，收到 workspace/group/connector、前後 applied revisions 及下一 UTC bucket effective_at。未 ready、stale CAS/source 不呼叫；rename 不重啟。主程式未注入，也沒有新增 HTTP Apply route。
- 有效行為 RED/GREEN 涵蓋 callback 缺失及 unavailable 原因遺失；精確 group-A／group-revision-1／15 秒／12:00:07Z→12:00:15Z fixture 證明 owner 在交易內改變後，後段 workspace 保存失敗會將 owner/version/applied/workspace 全回滾；原 targets/outbox/receipt 不變。
- 主代理新增第二次 semantic Apply 的精確 request assertions，將 previous applied revision 暫改空白時測試確實失敗；立即還原後 activation/lifecycle focused tests PASS。
- 主代理完整 Go：44 有測試 package PASS、6 無測試、零 fail；全 repo vet PASS。workspace activation/lifecycle race PASS。此完整結果早於上述六項 test assertions，產品 source 未再變動；新增 assertions 隨後 focused PASS。
- Share API/handler 14 個 top-level regression PASS；modbusshare domain 74 個 pass events、零 fail。Share hydration/settings/token/revision 與 local setup CAS 分別驗證。
- 全 repo lint FAIL：僅起始 SHA 相同的 service_probe.go:201 gocritic；沒有擴大修復。diff/check-lines PASS。production owner switch、真正 runtime rollback、外部 SQL commit/readback 仍由 C/F 驗收。

## A3.3 收尾實作與最後 gates

- 主代理核對 sole persisted authority、CAS、lifecycle 及已確認的兩項移轉決策。新增 docs/technical/studio-v2-write-groups.md 說明 024 additive schema、13 個 API operations、revision／相容 adapter、backup 及無 destructive down／不丟 accepted backlog 的回復邊界。沒有操作正式資料庫、備份或部署。
- 審查發現並以有效行為 RED→GREEN 修正：同批 row groups 認領同一 legacy target；external schema inspection 期間 source/connector stale 仍 ConfigReady=true；single before intent 遺漏原 database；SDK 接受 needs_review 缺 before/candidate 與 canonical list 重複 ID。最後主代理回讀來源及測試；readiness/migration 保存資料不變 assertions PASS。
- BeforeIntent.Database fixture 在 review 後關閉、重開真實 disposable SQLite，完整 intent 一致；來源 destination 改變後 stale review 拒絕，已審閱 database/provenance 保留。初次測試 API 用法錯誤是 compile failure，未算 RED；修正測試後 original source 出現 empty database 的有效 RED。
- 主代理最後 go test -p 1 -parallel 1 ./... -count=1 -json：44 個有測試 package PASS、6 無測試、零 fail events；go vet ./... PASS；workspace/API group/migration/lifecycle focused race PASS。來源版本為本次 A3.3 最後產品檔案，結果 log 為 /private/tmp/go-gateway-a33-all-go.jsonl 及同前綴 vet/race logs。
- 最後 frontend lint、完整 Vitest（178 files／1037 tests）、build PASS；其後只將一個 error case 擴成三個 parameter cases，affected 40/40、tsc/scoped eslint PASS，未冒稱再次全測 1039 項。
- 全 repo golangci-lint FAIL，只有起始 SHA 同一 service_probe.go:201 gocritic（1 issue）；scoped domain/SDK lint PASS。此既有問題不阻擋本次 scoped 行為驗證，但不能稱全 repo lint 通過。
- Swagger 使用 swag@v1.16.6 重產，補可達 list/get 安全 500/503，原 API paths 無移除、新 WriteGroup operations 共 13；沒有 Apply endpoint。舊全域 Swagger basePath=/api 與 production /api/v1 差異屬起始版本既有設定，本文給實際 route，沒有擴大重寫所有 API。
- 真正 activation switches writer、binary/runtime rollback accepted backlog、live SQL commit/readback 與 UI→SQL scenarios 仍未執行，移交 C/F；A 不 archive 產品。B 可接續已通過的 domain/preparation gate；這不表示全六案完成。

## spectra-verify／spectra-review（2026-10-02）

- verify：無 Critical，10/10 任務完成；3 個 scenario（Activation switches writer、Rollback with accepted backlog、Delete with backlog 的 worker 完成部分）本來就依賴 production wiring，已從本 change 的 delta 移到 `wire-durable-write-group-delivery` 的 delta（新 requirement「Single production writer owner and lifecycle backlog」），本 change 的 spec 不再宣稱未實作的行為；兩個 change 重新通過 `spectra validate`。
- review 共 4 個 Warning、1 個 Suggestion，處置如下：
  1. **已修**：`ConnectorService.Delete` 沒有歸屬保護，可刪掉 canonical group 仍使用的 connector 並使之後所有 migrated read 失敗。新增 `WriteGroupService.PreflightConnectorDelete`（含 tombstone；缺 connector ID 回驗證錯誤）、`ConnectorService.SetDeleteGuard`，主程式與 handler 注入，`DeleteConnector` 以 409 `WRITE_GROUP_LEGACY_WRITE_CONFLICT`＋`open_write_groups` 回應；workspace 與 router 各有測試（router 測試確認 mapping 與 connector 都未被動）。限制：preflight 與實際刪除之間沒有同一交易（與 mapping 的 transaction 內重查不同），同時發生的新 migration 仍可能競爭；後續若要完全封閉需把 connector 刪除納入 owner transaction。
  2. **不改（設計）**：對 disabled 群組執行 Apply 是顯式的重新啟用，既有 `write_group_lifecycle_test.go` 明確釘住此行為（`Update` 保持 disabled、`Apply` 才回 ready 且受 CAS／readiness 把關）；沒有 Apply route，也無法被無意觸發。
  3. **接受並記錄**：workspace readiness 現在會對每個非 deleted／disabled 群組做 destination inspection（外部連線）。這是 A1.3 為了誠實顯示 schema readiness 的刻意選擇；風險是頁面載入時有外部連線與讀寫競爭。後續 E 改 UI 時應改為 summary 僅用 config-only readiness，詳細 inspection 保留在單群組 readiness route 與 Apply。
  4. **接受並記錄**：未 Apply 的 saved group 會讓 workspace readiness 出現阻擋性的 `write-group-apply-required`，而目前沒有 Apply route（待 C／E）。此為預期的過渡狀態（草稿尚未成為 writer），且有測試釘住；使用者可 disable／delete。
  5. **已修**：前端 allowlist 補上 handlers 實際回傳的 6 個 `WRITE_GROUP_*` code 與 `review_request`／`review_group`／`disable_group` 三個 action，並補 en／zh-TW 訊息（既有長行內增加，locale 仍 511 行）；擴充 `backendErrorCodes.test.ts`。
