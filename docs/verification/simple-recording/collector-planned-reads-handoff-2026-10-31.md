# Scheduled planned reads：task 3.1 交接（2026-10-31）

## 狀態與歸屬

- 實作與本節列明的相關驗證已完成，source/test 已凍結；task 是否完成由 root 核對後透過 CLI 決定。本批未自行 done、commit 或 push。
- Change：`integrate-protocol-aware-acquisition-batching`，僅 task **3.1**；[tasks](../../../openspec/changes/integrate-protocol-aware-acquisition-batching/tasks.md)、[design](../../../openspec/changes/integrate-protocol-aware-acquisition-batching/design.md)。
- 比較基準：`f6916e9cfbc2eb02ba777cba5f1e0e8270a3fe96`。共享工作樹的 API/member/frontend 等其他 owner 變更不歸本批。
- **141 個 source/test 路徑**：87 production（48 修改、39 新增）、54 tests（1 修改、53 新增）；加本交接文件為 **142 路徑**。root 先前已 ingest 的 design/spec 不混入這份 source manifest。
- 最終 141 行 `sha256sum` manifest digest：`70275a89f81f65ee2bbe582f61d7357486ac31c7dfd601efc3b1449686036491`。
- 最後變更僅 membership fixture：`4275ef783bae77e4223108af2aa74c9ffb9741e43947ea396411f8495fbdf231`；不用先前 `d7df817f…`。

## 實際接線與保留契約

既有 `cmd/test_ui` services → SQLite/source-rule 安裝 → 原 scheduler ticker → 純 optimizer acquisition-plan-v5 → ManagedConnection/PlannedReader → 真 adapter wire peer → 原 ingestor → SampleSink。Fixture decorator **仍委派真 groupPipe**；SampleSink 被呼叫不等於 SQL commit，loopback peer 不是 physical PLC。

- Scheduled 才遷移；保留 device lock、原 group cadence、既有各設備並行與 group wait，不另造 scheduler/reducer。
- 同一短鎖捕捉 installed device/source/group、ReadProfile、PointRevision、InstallationLease/Scope。Plan FP 不含捕捉時間、AcquisitionID、RequestID、observer enable 或 ephemeral token。
- 原 CfgFP、EffectivePoint.SourceIdentity/SourceRevision、StableSampleID/SampleEnvelope、frozen raw→mapped binding、revision CAS/publication guards 不變；不增加 SourceEpoch。
- 每 Point 每 scan 最多一次 final emission。同 request 用實際 final gateway completion；不同 request 保留各自 completion，不冒稱 PLC atomic snapshot 或 SQL receipt time。
- 每 request 原 timeout、共享 `1+retry` actual-send budget；scan 固定 checked deadline 預留最多一次安全 fallback 的 request 上界，尊重較早 caller deadline，不延長。
- 每 actual send 在原 wirebudget 邊界先驗 context/captured install/config/stop，成功 admission 是線性化點；後發 disable 不撤回已 admitted send，但阻止後續 retry/fallback。
- 只有完整 matching unsupported response 可一次 replan unresolved Points；FATEK Cmd48/status **9**，不是 4。Timeout、CRC/短包/錯站/錯命令不授 fallback。
- 降級 cache 綁 final-wire 成功 connection token/capability revision；failed/late connect 不清較新版 downgrade。Token 不進 sample identity。
- Health/statistics 保持 final Point 語意，不把一包當一 Point；某 Point decode Bad 不污染其他 Good。Stale/canceled request eligibility 與新 device breaker failure 分開。
- Device config/protocol 改動使舊 installation stamp sticky-invalid，A→B→A 不復活舊 lease；合法新 install 才有效。Same-config refresh 不假造新 epoch。
- Native `RejectedObservationCount()` 只計既有 acceptance guard 拒收的 observation；不稱特定 source-quality/SQL error，未擴 runtime.Stats/JSON/API。

### 能力宣告

只沿既有 `Device.ConnectionConfig.acquisition` optional v1 物件：`version/profile`，可選 `limits` 與 `coverage`。同純 resolver 供 validation/install/schema；不新增 table/config source。

- 缺省 unknown-safe；`safe-v1` 禁 random。拒 null/unknown fields/version/profile/無效 limits/coverage。
- 顯式 family catalog：`fatek-fbs-ascii-v1`、`mitsubishi-fx5-cpu-3e-binary-v1`；後者只涵官方 FX5 CPU built-in Ethernet 情境，不泛稱 MC/Q/L/R/module。
- Client 存在、subtype、成功 connect 都不是設備 capability 證據；不 probe 猜能力。
- `native_units` 是每 request native component 更嚴上限；request/response bytes 各自取最嚴正值。
- `RandomProtocolLimit` 僅額外限制 random（FX5 192），不放寬既有 ProtocolLimit，也不套成 contiguous/bit 上限。
- Canonical area、實體半開 range/bounds、coverage content FP 均重用既有 optimizer authority，不能只信 opaque revision。

## 為何 source proof/hooks 是 3.1 必要依賴

純 planner 的 readable-gap 候選接到真 wire 後，active slots **不能證明** disabled/skipped cells 不存在。不得把 source-rule 範圍或 recording selection 當可讀授權。新增的是既有 source snapshot 的 derived immutable metadata，不是持久 source registry。

| 範圍 | 本批必要作用 |
|---|---|
| workspace/sourcerule | 同 canonical SQL source inventory/witness，包含 disabled/skipped/zero/layout_mode/mixed_items，拒 missing links、截斷、unlinked standalone、foreign/unknown。 |
| optimizer | 重用同 canonical span/parser/bounds 產 exclusion，active overlap 沿原規則；不建立第二地址解析 authority。 |
| point/sourcerule config writes | Durable write 前 quarantine；nested 同 owner context 不新增/提前釋 pending；direct write 只 invalidate，不自行 grant。Telemetry 完全不接。 |
| pollinggroup shared SQL repository | 同 writer reservation 先列精確 FK affected devices→nonblocking TryBegin all→DELETE→SQL cleanup→Finish，覆兩條正式 Delete path，防 concurrent assignment 漏 device。 |
| collector/runtime | 原 slot/config/mapping/scope receipts 與 SQL witness 一起接受；TryLock 全驗後一次 assign，cleanup 只 revoke exact own receipt。 |

未知完整性只禁 declared-readable gap，保留 forbidden/excluded，exact contiguous/active-only random 仍合法。只有真的宣告 Readable 卻因 proof 不足被停用才附 closed `PlanReason=source_coverage_unverified`；不是 QualityReason，也不改 Good 品質。

多操作 pending 不會因一個先完成就解除另一個 quarantine。Source-scoped B 變更可以使跨 B 地址的 A 舊 gap plan 失效，但不重裝或啟用 A 未選點。純 TagID/MappingID link 值與 recording selectors 不新增 acquisition authority；**Share/Scale/Naming 原 SourceRule.Update 的 bump revision/syncPoints/reinstall 行為保留**，不宣稱所有 output-looking 欄位都零 reinstall。

### 四個安全閉環及相容性

| 真失敗證據與修正 | 最小 guard |
|---|---|
| Full apply 選到 late foreign device/removal targets | 原 installed target set、full slots/config CAS、own removal receipt，不從晚 getter 採納 foreign。 |
| Provider 期間原 SQL projection payload/scope 被換 | Provider 前 FrozenCopy；原 scope/state receipt 比較，不把原 witness 拼到新 input。 |
| Profile failure cleanup 撤掉 foreign 新 lease | 只使用原 captured/owned expected lease；不借剛讀到的新 lease 清除。 |
| Compensation 後 fresh lease 或零點 skip-only rule 回滾失敗 | 原 SQL witness＋原 frozen profile/full slots＋成功 compensation lease；失敗只撤自己的 forward receipt。成功 restored-old/foreign pointer 不撤。 |

最後一項涵 observer=nil；native revoke 在 pending release 前，不靠 diagnostic enable 才安全。原 unknown 不會因舊 SQL graph 存在就復活 complete。

Final selection 另抓到三個既有 observation compatibility tests：ordinary duplicate device install、optional observer 同步 reentry 新版 full install。原 test 預期未改；只 journal 同 builder 自己成功的 device config，仍核 own marker/config/full slots；full apply 在 SQL cleanup/native publication 完成後、optional observation 前釋自己 pending。其後只餘原 exact-lease observation binding，沒有舊 geometry/state regrant。

## 已執行證據（選測，不是全產品）

環境：Linux/amd64，Go 1.25.5；`GOMAXPROCS=1`、`GOGC=50`、Go `-p 1`；全程單 heavy slot，使用既有 caches。Official golangci-lint 2.6.2、原 repo config，沒有 blanket suppression。

工作區證據目錄 `collector31-evidence` 位於 repo 旁，**不包含在這 142 個 Git 路徑內**；本文件記錄結果，完整 command/JSON/log/hash 由 root 本次驗收保留。

| Gate | 實際結果與邊界 |
|---|---|
| 真 production wiring | MC/FATEK/Modbus contiguous 3→1；顯式 sparse random、unknown-safe、strict unsupported/fallback、reconnect、late/disable/stop、partial Bad、source membership/coverage fixtures PASS。 |
| Complete coverage | MC/FATEK/Modbus readable coverage 3→1；unlinked/unknown 3→3；public disabled/skip-only 及 scoped fresh acceptance/compensation PASS。 |
| 四閉環 | `closure-first-red.jsonl` 保留 pre-fix failure；`closure-snapshot-green.jsonl`、`closure-rollback-green1.jsonl` PASS。 |
| 零點補償 | `zero-rollback-red.jsonl`→`zero-rollback-green.jsonl` PASS；6 top-level/13 含子測試，含失敗安全3 exact、成功原 exclusion2 requests、普通補償3→1。 |
| 相容性 | 首次 `final-normal.jsonl` FAIL 原三 tests；`compatibility-green.jsonl` PASS 原預期與 foreign/pending/fullslot/rollback negatives。 |
| Final normal selection | `final-normal2.*`：13/13 packages、**359 top-level／1,421 含 subtests PASS，0 FAIL/0 SKIP**。之後只改下列一個 fixture，其他 bytes 不變；最新該 fixture normal 再 PASS 3.483s。 |
| 首次完整 related race selection | `final-race.*` **aggregate FAIL**：359 top-level 中358 PASS/1 fixture FAIL，1,420 含子測試 PASS/1 FAIL；12 packages PASS，cmd 其餘22 top-level PASS。沒有 DATA RACE report。不能稱這一輪全 PASS。 |
| 修後 affected race | 只改 `acquisition_scheduled_membership_test.go`，`membership-race-diagnostics.*` PASS 29.286s；reuse 原 byte-identical 12 packages 及其餘 cmd cases，不再重跑全359。 |
| Official scoped lint | `final-lint2.log` 133 findings FAIL 保留→`final-lint3.log` **0 issues PASS**；最後單 fixture `membership-lint-final.*` **0 issues PASS**。 |
| Vet | `final-vet.*`：同13 packages `go vet -p 1` PASS。 |
| Static/source freeze | 最終 `git diff --check`、官方 line gate PASS；原超長 sourcerule/service.go 1,087→1,051，未增行。141 source/test SHA 僅上述 fixture 相對 candidate 改變。 |

上述 top-level 與「含 subtests」不是兩組不同測試，不能相加。Normal/race 是不同執行，不以加總製造測試數。

### Race fixture 修正不是 production timeout 修正

舊 membership fixture 以「整段 SQL update 前後恰好兩次 wall-clock tick」斷言 wire 數。Race 下真 SQL update 跨7輪，每個實際 request 都仍是完整21 addresses，故舊總次數 assertion 失敗。

修正保留真 ticker/SQL/ingestor：先觀測 update 返回後的一個**完整**21-point round，再取下一個 distinct AcquisitionID。既有 ticker 同步完成 pollGroup/wait 後才 capture 下一輪，故選到的是 capture-after-update，而不只是 completion-after-update。Before/after 每 scan 的21 unique Points、單 RequestID、原 PlanFP及所有額外 wire frame 的完整21 addresses仍硬驗。

首次新 barrier 的7s觀測等待在 race 仍不足；保留該 FAIL。實測丟棄63個舊-completion samples，barrier＋下一scan抵達 SampleSink 花15.6305s，所選 sample gateway-completion→sink 延遲13.9657s。最終只把 **test observation bound** 設30s；production wire timeout、scan deadline、retry、cadence、source pipeline 均未改。這是 race-instrumented SQLite FIFO backlog 證據，不是4.2 throughput benchmark或真設備性能結論。

### 命令與選擇

完整 literal regex/命令保存在 `final-normal2.command`、`final-race.command`、`membership-normal-final.command`、`membership-race-diagnostics.command`、`final-vet.command`、lint command files。`closure-test-selection.json` 明列每 package 實際359個 top-level names。以下短 regex 在本次凍結檔案已逐名核對，與該359-name selection **零 extra／零 missing**，不是執行零 tests 的泛稱：

```sh
RUN='Acquisition|Scheduled|ReadProfile|PointRevision|Observation|SourceProjection|SourceRuleShareFailure|WorkspaceProjection|Planned|FormatRecovery|Admission|SendObserver|Budget|^TestCommandRejectionRequiresValidatedRandomStatusNine$'
PKGS=(./cmd/test_ui ./internal/datalink/collector ./internal/datalink/connector
  ./internal/datalink/connector/adapters ./internal/datalink/device
  ./internal/datalink/optimizer ./internal/datalink/point ./internal/datalink/pollinggroup
  ./internal/datalink/runtime ./internal/datalink/sourcerule ./internal/datalink/workspace
  ./internal/protocol/fatek ./internal/protocol/wirebudget)
go test -p 1 -count=1 -json -run "$RUN" "${PKGS[@]}"
go test -race -p 1 -count=1 -json -run "$RUN" "${PKGS[@]}"
go test -p 1 -count=1 -json ./cmd/test_ui -run '^TestProductionScheduledAcquisitionRecordingSubsetNeverNarrowsSource$'
go test -race -p 1 -count=1 -json ./cmd/test_ui -run '^TestProductionScheduledAcquisitionRecordingSubsetNeverNarrowsSource$'
go vet -p 1 "${PKGS[@]}"
golangci-lint run --concurrency=1 --timeout=8m --new-from-patch="$OWNED_PATCH" "${PKGS[@]}"
```

`OWNED_PATCH` 僅本批141 paths；最後 fixture lint 使用單檔 patch及 `./cmd/test_ui`。未對全 dirty tree 套 auto-fix。5處 sqltx contextcheck 精確行註記保留原 independent bounded cleanup、其 terminal API 無 ctx，不修改 cleanup engine。

## 獨立 review 的實際覆蓋

原 A/B 分工審查與中途四閉環 review 是歷史證據，不用「A看兩處」替代 final141。最後另兩名 report-only reviewer 核對全部 attributable final production hunks/new files及對應 tests，產出 exact path/hash ledger：

| Final review | 覆蓋 | 結論 |
|---|---|---|
| Core | cmd/collector/connector/adapters/device/optimizer/protocol：41 production＋37 tests＝78 | 0 blocking findings；最後 membership `4275ef78…` 已重看。 |
| Source | runtime/point/pollinggroup/sourcerule/workspace：46 production＋17 tests＝63 | 0 blocking findings；含該分區的 witness/quarantine/atomic publication/rollback 與 lint/相容修正。 |

兩 ledger 聯集與 final manifest 實際比對 **141/141 相符、無缺漏/重複**（`final-review-coverage.json`）。Core report SHA：`6ebca47daddac760ee678a587f6914029100e89f1addefcec275355219c9dd22`；Source report SHA：`353474e5ffcb158a700255115ca28a2676130d18db6c5e815f4bb9f4af66ae84`。

Reviewers 沒有跑 Go/test/build/lint，不能把靜態 review 當獨立 race PASS；未逐檔覆核全部既有未修改 regression tests，也不替其他 owner 的 dirty paths 背書。Final source差異中沒有未被這兩分區覆核的 production path。

## 未做／下一步邊界

- 不宣稱全 repo/全產品/API/frontend 全套通過；未做 physical PLC、正式 DB/SCADA、SQL commit durability或速度提升驗收。
- 3.2 fair due queue/coalescing/backpressure/drop counters 未實作；原 group wait 仍在。
- 3.3 PollNow planner migration 未實作；現有 scalar/manual ingestion guards 保留。
- 4.1 preview API/UI 與 public Stats exposure 未實作；4.2 100-scan/load/100–10000 planner/HSL benchmark 未做。
- 任意 external SQL 即時同步不是本批正式服務 config-write hook 契約。
- user 要求全部本地整合完成後才同步 main；本交接不是 push/部署授權。

## 精確 source/test 路徑（141）

本文件自身為第142個 attributable path；下列87 production/54 tests分類依 `_test.go`，包別依上述權威說明。完整 bytes manifest及patch為 `final-sha256.txt`、`final-owned.patch`；不包含 unrelated owner paths。

- `cmd/test_ui/acquisition_batch_wiring_test.go`
- `cmd/test_ui/acquisition_scheduled_capability_test.go`
- `cmd/test_ui/acquisition_scheduled_coverage_test.go`
- `cmd/test_ui/acquisition_scheduled_fallback_safety_test.go`
- `cmd/test_ui/acquisition_scheduled_fixture_test.go`
- `cmd/test_ui/acquisition_scheduled_geometry_exclusion_test.go`
- `cmd/test_ui/acquisition_scheduled_geometry_failure_test.go`
- `cmd/test_ui/acquisition_scheduled_geometry_mutation_test.go`
- `cmd/test_ui/acquisition_scheduled_membership_test.go`
- `cmd/test_ui/acquisition_scheduled_partial_test.go`
- `cmd/test_ui/acquisition_scheduled_reconnect_test.go`
- `cmd/test_ui/acquisition_scheduled_safety_test.go`
- `cmd/test_ui/acquisition_wire_peer_test.go`
- `cmd/test_ui/service_wiring.go`
- `internal/datalink/collector/collected_value.go`
- `internal/datalink/collector/scheduler.go`
- `internal/datalink/collector/scheduler_acquisition_dispatch.go`
- `internal/datalink/collector/scheduler_acquisition_execute.go`
- `internal/datalink/collector/scheduler_acquisition_execution_cache_test.go`
- `internal/datalink/collector/scheduler_acquisition_execution_test.go`
- `internal/datalink/collector/scheduler_acquisition_full_mutation.go`
- `internal/datalink/collector/scheduler_acquisition_geometry.go`
- `internal/datalink/collector/scheduler_acquisition_geometry_admission.go`
- `internal/datalink/collector/scheduler_acquisition_geometry_test.go`
- `internal/datalink/collector/scheduler_acquisition_health_test.go`
- `internal/datalink/collector/scheduler_acquisition_installation_test.go`
- `internal/datalink/collector/scheduler_acquisition_mutation_batch.go`
- `internal/datalink/collector/scheduler_acquisition_mutation_batch_test.go`
- `internal/datalink/collector/scheduler_acquisition_partial_test.go`
- `internal/datalink/collector/scheduler_acquisition_plan_test.go`
- `internal/datalink/collector/scheduler_acquisition_publication.go`
- `internal/datalink/collector/scheduler_acquisition_publication_batch_test.go`
- `internal/datalink/collector/scheduler_acquisition_publication_test.go`
- `internal/datalink/collector/scheduler_acquisition_read.go`
- `internal/datalink/collector/scheduler_acquisition_reader_test.go`
- `internal/datalink/collector/scheduler_acquisition_rollback.go`
- `internal/datalink/collector/scheduler_acquisition_rollback_test.go`
- `internal/datalink/collector/scheduler_acquisition_safety_test.go`
- `internal/datalink/collector/scheduler_acquisition_snapshot.go`
- `internal/datalink/collector/scheduler_config_methods.go`
- `internal/datalink/collector/scheduler_dispatch.go`
- `internal/datalink/collector/scheduler_observation.go`
- `internal/datalink/collector/scheduler_observation_test.go`
- `internal/datalink/collector/scheduler_polling.go`
- `internal/datalink/collector/scheduler_profile_identity.go`
- `internal/datalink/connector/acquisition_connection_token.go`
- `internal/datalink/connector/acquisition_connection_token_test.go`
- `internal/datalink/connector/acquisition_profile.go`
- `internal/datalink/connector/acquisition_profile_decode.go`
- `internal/datalink/connector/acquisition_profile_schema.go`
- `internal/datalink/connector/acquisition_profile_test.go`
- `internal/datalink/connector/acquisition_profile_validation_test.go`
- `internal/datalink/connector/adapters/acquisition_connection_generation.go`
- `internal/datalink/connector/adapters/acquisition_connection_generation_test.go`
- `internal/datalink/connector/adapters/fatek.go`
- `internal/datalink/connector/adapters/fatek_context.go`
- `internal/datalink/connector/adapters/fatek_persistent.go`
- `internal/datalink/connector/adapters/mc3e.go`
- `internal/datalink/connector/managed_planned_read.go`
- `internal/datalink/connector/manager.go`
- `internal/datalink/connector/manager_acquire.go`
- `internal/datalink/connector/planned_read.go`
- `internal/datalink/connector/registry.go`
- `internal/datalink/device/acquisition_validation.go`
- `internal/datalink/device/acquisition_validation_test.go`
- `internal/datalink/device/service_validation.go`
- `internal/datalink/optimizer/acquisition_builder.go`
- `internal/datalink/optimizer/acquisition_coverage_validation.go`
- `internal/datalink/optimizer/acquisition_coverage_validation_test.go`
- `internal/datalink/optimizer/acquisition_plan.go`
- `internal/datalink/optimizer/acquisition_random.go`
- `internal/datalink/optimizer/acquisition_random_protocol_limit_test.go`
- `internal/datalink/optimizer/acquisition_source_span_test.go`
- `internal/datalink/point/acquisition_group_mutation_test.go`
- `internal/datalink/point/acquisition_mutation.go`
- `internal/datalink/point/acquisition_mutation_test.go`
- `internal/datalink/point/format_mutation.go`
- `internal/datalink/point/service.go`
- `internal/datalink/point/service_delete_revision.go`
- `internal/datalink/point/service_materialize_revision.go`
- `internal/datalink/point/service_ordinary_revision.go`
- `internal/datalink/point/service_point_crud.go`
- `internal/datalink/point/service_polling_group.go`
- `internal/datalink/point/service_restore_revision.go`
- `internal/datalink/point/service_update_revision.go`
- `internal/datalink/pollinggroup/acquisition_mutation.go`
- `internal/datalink/pollinggroup/acquisition_mutation_failure_test.go`
- `internal/datalink/pollinggroup/acquisition_mutation_test.go`
- `internal/datalink/pollinggroup/sql_repo.go`
- `internal/datalink/runtime/ingestor.go`
- `internal/datalink/runtime/rejected_observations.go`
- `internal/datalink/runtime/rejected_observations_test.go`
- `internal/datalink/runtime/service.go`
- `internal/datalink/runtime/service_projection_state.go`
- `internal/datalink/runtime/service_read_profile.go`
- `internal/datalink/runtime/service_workspace_projection.go`
- `internal/datalink/runtime/service_workspace_projection_scope.go`
- `internal/datalink/runtime/source_acquisition_finalization_test.go`
- `internal/datalink/runtime/source_acquisition_full_ownership_test.go`
- `internal/datalink/runtime/source_acquisition_geometry.go`
- `internal/datalink/runtime/source_acquisition_group_mutation.go`
- `internal/datalink/runtime/source_acquisition_input_ownership_test.go`
- `internal/datalink/runtime/source_acquisition_profile_ownership_test.go`
- `internal/datalink/runtime/source_acquisition_projection.go`
- `internal/datalink/runtime/source_acquisition_publication.go`
- `internal/datalink/runtime/source_acquisition_removal.go`
- `internal/datalink/runtime/source_acquisition_rollback.go`
- `internal/datalink/runtime/source_acquisition_scoped.go`
- `internal/datalink/runtime/source_projection_lease_guard.go`
- `internal/datalink/runtime/source_projection_rollback.go`
- `internal/datalink/runtime/source_rule_observation.go`
- `internal/datalink/runtime/source_rule_projection_scope.go`
- `internal/datalink/sourcerule/acquisition_link_mutation.go`
- `internal/datalink/sourcerule/acquisition_source_exclusions.go`
- `internal/datalink/sourcerule/acquisition_source_exclusions_test.go`
- `internal/datalink/sourcerule/acquisition_source_identity.go`
- `internal/datalink/sourcerule/acquisition_source_snapshot.go`
- `internal/datalink/sourcerule/acquisition_source_snapshot_test.go`
- `internal/datalink/sourcerule/runtime_acquisition_mutation.go`
- `internal/datalink/sourcerule/runtime_acquisition_mutation_test.go`
- `internal/datalink/sourcerule/runtime_projection_rollback.go`
- `internal/datalink/sourcerule/runtime_reconcile.go`
- `internal/datalink/sourcerule/service.go`
- `internal/datalink/sourcerule/service_links.go`
- `internal/datalink/sourcerule/service_rule_transform.go`
- `internal/datalink/sourcerule/share_recovery_snapshot.go`
- `internal/datalink/workspace/runtime_acquisition_acceptance.go`
- `internal/datalink/workspace/runtime_acquisition_acceptance_test.go`
- `internal/datalink/workspace/runtime_acquisition_batch_acceptance_test.go`
- `internal/datalink/workspace/runtime_acquisition_completeness_test.go`
- `internal/datalink/workspace/runtime_acquisition_fixture_test.go`
- `internal/datalink/workspace/runtime_acquisition_sql_snapshot.go`
- `internal/datalink/workspace/runtime_acquisition_sql_snapshot_test.go`
- `internal/datalink/workspace/runtime_projection_frozen.go`
- `internal/datalink/workspace/service_runtime_projection.go`
- `internal/datalink/workspace/service_runtime_projection_acquisition.go`
- `internal/protocol/fatek/command_rejection.go`
- `internal/protocol/fatek/command_rejection_test.go`
- `internal/protocol/wirebudget/admission.go`
- `internal/protocol/wirebudget/admission_test.go`
- `internal/protocol/wirebudget/budget.go`
