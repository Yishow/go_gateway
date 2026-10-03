# 基本記錄設定與範圍啟動驗證

2026-10-04；change：`streamline-studio-v2-recording-setup`；起點：`38c023e403efefc7601dbf8628d2fd8e8627b86f`。

## 實作邊界

- 保留 `/studio/v2` 四步、既有 Point-to-Tag 保存與真實讀值。新基本群組使用 backend 持久 `(workspace_id, device_id, managed role)` 唯一身分；既有 advanced 群組不自動轉換。
- 基本畫面顯示設備、實際選取資料、目的地、60 秒預設快照間隔與缺值政策。canonical group 不依賴 legacy targets/row_groups；自訂表、欄位、revisions、schema 準備及試寫保留可達的進階入口。
- `recording-start` 是既有 save/readiness/Apply/activation services 的狹窄協調，使用既有 operation ledger 的獨立 action、request identity、原意圖 digest、scope facts、owner fencing 與進度。
- 保存／Apply 的本地設定與 ledger checkpoint 同一 transaction；外部 schema confirmation 仍獨立明確執行。Start 不執行 DDL，不宣稱跨資源 rollback。
- 原請求重送與失去回覆使用同一 operation；初始無 effects 的中斷可恢復。已 Apply 的部分結果不重做 Apply；外部 intent/source/schema/connector 改變要求重新驗證。
- server 回傳 operation 的目前 setup revision，使 UI 能區分該 operation 自己產生的新版本與外部 schema preparation。完成準備後的新意圖保留同一基本群組；未知回覆不丟掉原 request identity。
- restart 只在先前 barrier 已持久驗證、原 workspace/settings revision 未變及 current hydration 仍 ready 時取得目前 server token，並重跑原 barrier。原 request/digest 不變；缺 proof、版本漂移或 failed hydration 仍拒絕。
- delivery 按目前 applied revision 查詢採集／接受階段與真正 receipt-backed SQL effect。歷史 revision、running、queued、unknown 或成功試寫都不冒充目前群組落庫；SQL committed 不等於另外執行過 readback。

## 規格情境與證據

兩份 delta 共 4 個 requirements、18 個 scenarios，沒有 Examples；不新增 test-scope exclusions。Spectra validate 通過；analyze 的 Coverage／Consistency／Gaps／Localization clean，15 個 Suggestion 僅提醒抽象 scenario 沒有 concrete Examples；本案不新增規格數據或驗收條件。

| Scenario | 對應實作／測試 |
| --- | --- |
| Mixed sensor configured without SQL | `basic-point-tag-to-group.test.ts`、`basic-recording-state.test.ts`；沿用 typed candidates，不新增衍生量前置 |
| Reload and edited connection | `step1-capability-and-drafts.test.tsx`、既有 device persistence/readiness、autosave/hydration；Start 的 source/tag/connector/schema drift regression |
| Unavailable protocol | `step1-capability-and-drafts.test.tsx`、`protocols.test.ts`；MQTT unavailable reason 與 selection guard，本案沒有新增協議或 parser |
| Same address on different devices | `address-conflicts.test.ts`；persistent basic key 使用 device identity，非名稱／位址 |
| Missing advanced measurement | `basic-point-tag-to-group.test.ts`；basic raw member 沒有 fake measurement identity |
| Offline draft navigation | `step1-capability-and-drafts.test.tsx`、既有 save barrier suites；實際 UI 修改 A port 後 probe 失效，失敗 probe、保存並繼續、reload 保留 15039 與未測試；B 的成功狀態保留 |
| Routine mapping does not require duplicate entry | 既有 mapping autosave/reconciliation suites、`basic-point-tag-to-group.test.ts`；Step 3 真值 A=215/1013、B=187/777，四個 persisted mapping 保留 |
| Existing advanced configuration is reopened | Basic panel/state tests；既有 advanced group identity/layout 與手動欄位／試寫入口保留 |
| New basic setup is resumed | `write_group_basic_managed_test.go`、Basic panel regressions；lost reply/reload/concurrent/restart、同名設備、目的地 CAS conflict、tombstone、advanced separation |
| Ready managed group without legacy targets | Basic panel 與 coordinator fixture；只有 canonical group 與 confirmed proof，legacy collections 為空 |
| First bucket is not yet due | Basic evidence tests；60 秒及 UTC 完整桶語意，不倒數假成功或偷寫半桶 |
| Missing schema confirmation | `TestRecordingStartMissingPreparationHasNoDDL`；保存可查、prepare_schema、安全停止，0 DDL |
| Duplicate start or response loss | ConcurrentSameRequest、CompletedReplayAfterSQLConfigRestart、StaleLease、RecoversClaimBeforeFirstProgressCheckpoint；前端 unresolved request／token rotation regressions |
| Partial start and changed intent | PartialRetryReusesAppliedRevision、ApplyCheckpointRollback、RejectsChangedDigest、RejectsSourceTagConnectorAndSchemaDrift、RevalidatesAfterExplicitPreparation |
| Selected device scope is independent | scoped readiness/activation、runtime scoped projection、Share exact-scope tests；共享 polling group 及同設備未選點位的 ticker 回歸；實際 B start 前後 A running/projection version 不變 |
| Share-only setup | WithoutDatabaseRequiresEnabledShareOutput、Basic/SDK 空 revision、state reload、fresh rule Create reconcile 回歸；全新空白 configuration 的真 UI 成功、0 DB resources、FC3=215/1013；同 operation token rotation reload 通過 |
| Destination offline after acquisition | committed-effect 的 queued/unknown/mismatched receipt 反例、revision stage tests；Basic evidence 不顯示 SQL success |
| Production first row committed | receipt-backed effect/query/SDK tests；正常 UI A/B 生產 delivery effect 與獨立目的地 row/receipt/provenance 讀回一致 |

| Requirement | 實作位置 | 有測試的情境／Examples | 本次執行 |
| --- | --- | --- | --- |
| Four-step intent-led setup | 既有 Step 1–3、`state/basicRecording.ts`、workspace `write_group_basic_managed.go`、migration 028 | 9/9；無 Examples | full frontend／Go；離線草稿與真映射 UI |
| Basic recording controls hide incidental complexity | `BasicRecordingPanel.tsx`、`Step4Database.tsx`、`BasicRecordingEvidence.tsx` | 2/2；無 Examples | Basic/SDK tests；正常 UI 等首桶與三寬度 |
| Recoverable scoped recording start | workspace `recording_start*.go`、`service_activation_scope.go`、runtime scoped projection、既有 recordingplan ledger、HTTP handler/router | 5/5；無 Examples | full Go／race；缺 proof 0 DDL、partial/restart 與健康鄰近設備 UI；fresh Share-only 真啟動／FC3 通過；同 operation reload 通過 |
| First recording result is evidence-backed | groupdelivery `committed_effect.go`、grouppipeline `delivery_view.go`、delivery DTO/SDK、Basic evidence | 2/2；無 Examples | full Go/frontend；queued/receipt 反例與實際 SQL 讀回 |

## 先失敗、再修補

- Start route、基本 create-once identity 與 safe SDK 的缺口由新回歸先重現，再接實際 service／ledger；不增加另一套 job、queue 或 token repository。
- completed replay 在後來刪除群組後錯誤阻擋；修正後回既有歷史成功，不再執行 live effects。
- hydration token 在 disabled Share／DB-only 路徑被漏掉；修正後保留原全域 readiness/ownership/revision gate，不使用 full-workspace restore 影響 B。
- 舊 revision 的 collecting/committed counts 原本可能被當成新 applied revision 證據；改為 current-revision stage counts 及 matching receipt identity。
- 真 UI 準備 schema 後，舊 start request 被重用且 GET 舊狀態覆蓋新 POST；frontend 回歸修正為已知外部 setup 變更建立新意圖，原群組 ID 不變；POST/GET 依 updated_at 保留最新結果。
- 部分 Apply 後重啟旋轉 token 原本無法接續；`gw-E-restart-barrier-red.log` 重現，修正後同 operation／同 applied revision 接續。workspace/settings/hydration 三個反例仍拒絕。
- Claim 後第一個 progress checkpoint 前中斷，原本永久 running/stale_intent；`gw-E-initial-checkpoint-red.log` 重現。新 Claim 先捕捉完整 source/schema facts；僅無任何效果的初始空 checkpoint 可重新初始化。
- fresh DB 查詢不存在的 operation 原本建立 workspace；`gw-E-readonly-get-red.log` 重現。GET 改為 repository 唯讀，缺 workspace 回 not found。
- 成功回覆遺失、reload 看見 own Apply 的新版本／重啟 token 時，原本建立第二 request identity；前端有效回歸先失敗，修正後未知結果重送原 request，backend 決定 historical result 或 stale/revalidate。
- 已知 running 操作若程序中斷，原本沒有 UI 重送入口；回歸先找不到 Retry。修正後可手動重送原 request，仍保留 external readonly、scope mismatch 與 actionBusy 保護，由 backend fencing 決定是否可接續。
- shared polling group 原本可能刪除／重建未選點位，包含同設備其他點位；ticker 回歸先失敗，修正後沿用未選點位與原 scheduler tick。
- fresh Share-only 無 DB setup revision 原本被 Basic/SDK 擋下；空字串只在已載入正確 workspace snapshot、沒有 DB groups 時接受，missing／其他格式與 DB 群組仍拒絕。
- fresh source rule Create 原本只保存、不 reconcile，形成 desired=2/raw=0；回歸先失敗，改用既有 `CreateWithRuntimeReconcile`，inactive/not-running 也能建立真正 Share projection。其 production running 路徑的舊全量 polling-group 更新會重建健康 B ticker，改走既有 scoped projection，Create 只處理該規則點位；更新／刪除保留完整 device 重新核對，避免留下已移除 scheduler 點位。running A／inactive B 草稿回歸先失敗，修補後 known inactive 回 deferred、兩個 missing/inconsistent scope 仍 stale；健康 A ticker 不變。
- 真 backend `dev-...` ID 原本被 UUID-only handoff 丟掉；回歸先得到 null，修正後只接受 server persisted/saved ID，草稿與未知 ID 仍拒絕；實際 B 導航 query 已核對。
- 已知 operation 的 token rotation 原本會建立新 request；反例先失敗，改為比較 settings／Share workspace revisions，保留原 request/token 交給 backend replay；外部 revision 漂移仍撤銷目前成功／handoff。
- Share-only 空 DB revision 原本被 session intent loader 拒絕，reload 丟失 request identity；新 state regression 先失敗，最小修補後 55 個相關測試通過，原 DB 格式保護保留。
- 最後 runtime 邊界回歸先重現零點位 Create 重建健康同設備 ticker，以及未知／空 status 被誤當成 deferred；修補後零點位 runtime no-op，只有合法 draft/disabled deferred，其餘 missing/inconsistent/unknown/empty status stale，Share authority 仍獨立驗證。
- 英文 SQL 狀態改為 commit confirmed；confirmed SQL receipt 與另外執行過的 readback verification 不混用。

## 本次執行

darwin/arm64、Go 1.27.1（module 1.25.5）、golangci-lint 2.14；Chrome for Testing 149／agent-browser 0.38.1。只使用自有 loopback simulators、可丟棄 configuration／destination DB；未操作正式 PLC、LAN、部署或正式資料。

| 命令／範圍 | 結果 |
| --- | --- |
| `cd frontend && npm run lint` | passed-current；`gw-E-frontend-fixture-final-lint.log` |
| `cd frontend && npm test -- --run` | passed-current；189 files／1210 passed；`gw-E-frontend-fixture-final-test.log` |
| `cd frontend && npm run build` | passed-current；`gw-E-frontend-match-final-build.log`；保留既有 chunk size warning |
| `go test -p 1 ./... -count=1` | passed-current；`gw-E-backend-scope-final-test.log`；最後 runtime 修補後全部 Go packages，避免固定 listener 平行衝突 |
| `go vet ./...` | passed-current；`gw-E-backend-scope-final-vet.log` |
| `golangci-lint run ./...` | passed-current；0 issues；`gw-E-backend-scope-final-lint.log` |
| 7 packages 的 focused `go test -race` | passed-current；handlers/API/runtime/recordingplan/groupdelivery 有 matching tests；workspace/modbusshare 此 filter 沒有 matching tests；`gw-E-final-race.log` |
| source-rule runtime `go test -race` | passed-current；same-device／other-device ticker、zero-point Create、inactive draft、unknown/empty/inconsistent scope、Update/Delete 9 functions；`gw-E-source-reconcile-scope-final-race.log` |
| `make gen-docs`，使用已安裝 swag 的 PATH | passed-current；三份 generated Swagger；`gw-E-swag-final.log` |
| normal embedded binary、主流程／錯誤／部分完成、鍵盤及三寬度 PNG | passed-current：實際主流程、缺 proof、partial 與原 operation restart；Share-only 另驗原 operation token rotation reload 通過 |
| `git diff --check`、`make check-lines`、Spectra validate | passed-current；line gate 保留兩份原有 511 行 locale，皆未增行；`gw-E-lines-scope-final.log`、`gw-E-validate-scope-final.json` |

## 實際 UI 與獨立 SQL

E 使用既有來源的可丟棄 configuration 複本驗證協調、production SQL 與真 states；Share-only 另用全新空白 configuration。完整新 workspace／空白目的地、七型別、三桶與三次計時矩陣仍由 F 負責，不能用 E 或 D 的證據代替。

- 只有正常 frontend build/static sync/`go build ./cmd/test_ui`，由真 embedded `/studio/v2` 操作；未使用 UI network mock、固定成功動畫或替代 handler。
- 原主流程 Basic A 群組 `a6087a28-2128-4443-9399-648d49e34aa5`、revision `83c3d6c2-3da9-439d-acd6-39d7cc2bd058`；第一次完整桶 `22:14:00Z`、record `d9fabf85b8a9efa8090ec7193eacc55007c02c9a7df68c05bb3eb6ca99f713b7`；SQL 215/1013 與兩個 good provenance 讀回一致。
- 最新正常 build 再驗 A 群組 `1f68a355-0f8f-4e2c-9c39-0667b1aa3a1f`、applied revision `5c2fdb3c-aab6-47fe-ba81-ac33ac142e78`；B 群組 `69a66511-05a2-41e3-8b81-506654b3c84c`、applied revision `a102977e-37e3-4730-99c1-0efbc8070f87`。
- 獨立 SQLite readback 的 `22:37:00Z` 桶：A=215/1013、B=187/777；record IDs 分別 `e52973b3c8550072d556200c71b28032926c059a060e81e162712dcac96ddbf3`／`4e03128737859fab05c80b0ac78368170359f1e810421e3e7449d8742c22a21b`；group/device/revision、receipt digest 與 provenance 一致。
- A 的目的地 commit `22:38:01.638518Z`、本地 ack `22:38:01.639637Z`；B 分別 `22:38:04.638171Z`／`22:38:04.639536Z`。兩個時間是不同 checkpoint，不要求完全相等；UI SQL committed 只表示 receipt-backed delivery，本報告另外完成真正 readback。
- 首桶前曾實際查到 target 0 rows，UI 保留 60 秒完整 UTC 桶等待與 SQL 未確認。畫面的 local accepted 是 queued/retrying 當下筆數；交付完成後 0 不代表從未接受。
- 缺 schema proof 的 error A 群組 `a25417c2-1d54-46bb-8eda-e9483492f3b8` 已保存，applied revision 空；獨立查證其 target table 不存在，Start 沒有 DDL。
- partial A 先因本次自有 SQLite trigger 阻擋 activation，顯示 saved/ready/applied=true、activated=false；僅移除該 trigger，再重啟／reload／鍵盤重送原 request。
- 原 operation `start-642b8e5fcd4e2e836d0a4edd21cf9a514ad0ff5548cc0f5ba315c5cd2ae2781e` 恢復成功；operation count、原 request、group/applied revision 均相同。server hydration token 已旋轉、原 workspace/settings revision 未變，原 barrier 重新驗證。
- 最新 DB token-only restart 後真 UI 鍵盤重送 `start-3f7ebd84960225ef3f744791fa5dda9c4751a544cac638288b0f49f59f5d5f11`：原 request、digest、private progress 完全相同，ledger 5→5；current token 旋轉且 settings／Share workspace revision 不變。較早 B 設定造成的 DB setup revision 變更會建立新意圖，與此條件分開核對。
- B start 前後 A 都 running/aligned，runtime projection 保留 `sha256:105956cd243f4173`；B start 後 `/studio/runtime?device_id=dev-1791051561497-6`，真正選 B，沒有落回預設 A。
- 全新 Share-only workspace `bc0bf69a-27c3-4243-a5d2-11e2d702d8cb` 經真 UI 新增設備、probe、rule、映射與明確 Share geometry，首次成功前沒有重啟或複製既有 configuration。
- Share-only operation `start-68c87067da17483e02269d47f0daf75efa3a20a313f4e218e81ae84df2c81a81` succeeded、groups 為空、設備 activated；獨立檢查 connectors/groups/outbox/receipts 均 0，真正 FC3 40001=215／40002=1013。
- Share-only restart／reload 後原 request `bf91b1d2-60b6-435a-83d1-0edf0409cdcc`、同 operation/digest 保留；server token 已旋轉但 settings／Share workspace revisions 未變，真 UI Enter POST 200／GET 同一 operation 200。獨立 ledger 仍 3 筆（兩次修復前失敗與一次成功），DB resources 仍 0、FC3 仍 215/1013。
- 同一 fresh Share-only setup 再以真 UI 新增 B 離線未測試草稿及 disabled／Share-disabled rule；A 仍可 start，B 保留 draft、未 activation，A 的兩個 Share mappings 與 FC3 215/1013 不變。新的 A 意圖 `start-9f79ef1f5965efc3d888fbe5dd49b8dc9943f23cdcafd6628fabddcd8e03ad6f` succeeded、groups 為空，connectors／groups／recording plans／outbox／receipts 都沒有建立。
- 最後完整 binary 再重啟此 A-ready/B-draft setup，真 UI reload／Enter POST 200、GET 同一成功 operation。原 request `fd4092d6-66ed-43bb-807e-be97ade52adc`、digest 與所有 ledger detail 完全相同，ledger 4→4；獨立只讀 SQL 核對 A active、B draft/未測試、B rule disabled，七種 DB resource counts 都是 0，FC3 仍 215/1013。
- 既有 advanced 群組 `8c8869ac-6b82-46b3-b545-79260df9b8a1` 與四個 members 保留；同一已知 start identity 的連續 Enter 重送沒有新增 operation/group/Apply。

18 張 PNG 在 `evidence-e/`，主代理已逐張檢視 main/error/partial 的 390/768/1440，及 waiting/mapping/offline/runtime、Share-only 三寬度與 offline-neighbor 兩張截圖。主流程三寬度已用最後 frontend identity 修補後的正常 build 更新，含真正 token-only replay；最後 source-rule scoped reconcile 修補不改 controls/layout。Share-only 另用最後 frontend build 驗原 operation reload；手機圖使用現成鍵盤側欄／connector 收合，展開的舊側欄仍有窄版佈局限制，不以此宣稱整個既有頁面已重設 responsive。

| 畫面／互動 | 設計比對與結果 |
| --- | --- |
| main 三寬度 | 必要 controls、selected members、60 秒語意、真 receipt effect；390 控制項直排、ID 換行，768 兩欄，1440 主區與摘要，無本案控制項水平裁切 |
| error 三寬度 | 清楚保存／未就緒／未 Apply／未啟用，安全 repair action；沒有 SQL 成功 |
| partial 三寬度 | 保留已 Apply 成果、未 activation、原請求 retry，不顯示 rollback 或整體成功 |
| Share-only 三寬度 | 沒有 DB writer／schema 前置；設備 activated 與 Runtime handoff 依真 operation 顯示，390 基本 controls 完整；既有進階 connector 區仍保留 |
| offline-neighbor | Step 1 保留 A probe 成功及 B 未測試草稿；Step 4 整體摘要仍提示 B 待修正，A 的 scoped start 成功，沒有冒稱全部設備已就緒 |
| 鍵盤 | Start Enter、原請求 Retry Enter、連接器與側欄收合 Enter，interval→Tab completeness 有 focus；沿用既有進階 schema preview/confirm |
| handoff | B 真 backend ID 保留；只在原 selected operation succeeded/activated 後顯示入口 |

Runtime 既有 renderer 在 focused A/B 時混入另一設備同位址標籤，值卻跟隨 focused device；該 renderer 不在 E diff。本案保留此範圍外限制，不宣稱 Runtime 表格已修復；設備 identity/採集以 runtime context、Step 3 真值、scoped point API 與獨立 SQL 證據查證。

沒有額外 Design Source；畫面依 E design 的四步、必要 operator controls、可達 advanced、明確 schema confirmation、逐 scope 進度與三段交付事實比對。完整 Playwright e2e、Windows／現場與 human acceptance 不由本機測試替代。

既有 source-rule Update/Delete 的 mapping cache 殘留在此次 review 被發現；原 baseline 已存在。本案只保證新限定投影不留下被移除的 scheduler／point metadata，未擴張為 mapping cache 全面修復。

最後獨立 review 對 correctness、efficiency、simplification/reuse、convention 四個面向沒有剩餘 findings；原 effect/count 候選經 receipt-backed UI 的反例核對後撤回，沒有為此增加範圍。兩個 runtime 修補 hash 為 `b98ee9bd3fc12d5a745e44071142f6c7ee991b027f53a100981d84013feeb667`／`eb8219338beb65dd144924e98c7b46847a04240f7ed6f42a849e985149a6f6f0`，主代理在最後全套 Go gates 後再核對一致。前端 relevant identities 亦與 1210 項全套通過時完全相同。

Verify 結論：10/10 tasks 的工作均有本報告證據，4/4 requirements／18/18 scenarios 有相應測試，無 Examples 或新增 exclusions；production wiring 與 design 的 scope、原 operation 恢復、獨立 DDL 確認及三段事實一致。本案 inspected scope 沒有剩餘 Critical、Warning 或 Suggestion；完整 e2e、現場與人員驗收維持 not-run。最終範圍以 touched tracking 的完整 source／測試／文件／18 張已檢視 PNG 捕捉並於 archive 前核對，snapshot 身分保存在本機交付紀錄，不在本文內自我引用。
