# 迴路記錄規格情境覆蓋與實跑證據索引 (Circuit Recording Scenario Coverage)

> 對應變更：`openspec/changes/add-circuit-recording-profiles/`（Task 5.2 證據索引結算）
> 可攜證據：[`circuit-recording-test-results-20261024.json`](./circuit-recording-test-results-20261024.json)（11套件、2,476個Go PASS路徑，另有policy-only與無擴充focused實跑；260檔1,812個前端測試PASS）。
> 收尾 review 修正與增量重測：[`studio-v2-circuit-recording-final-review.md`](../studio-v2-circuit-recording-final-review.md)（修正後 1,245 個 Go PASS 路徑；與前述完整結算分開記錄）。
> 資料庫矩陣源：[`docs/technical/studio-v2-circuit-recording-db-verification.md`](../studio-v2-circuit-recording-db-verification.md)、[`docs/technical/evidence/circuit-recording-db-matrix-20261024.json`](./circuit-recording-db-matrix-20261024.json)
> 真瀏覽器實測：Root 實測雙 profile、Space 鍵、批次只覆寫迴路 B、reload、越界、confirm/start、自動重建 generation 2 警告、23:26:07 cutover、23:28:11 rollback 及 SQL partial（截圖見 [`docs/technical/evidence/circuit-recording-20261024/`](./circuit-recording-20261024/)）。

---

## 一、Spec 1: `circuit-recording-profiles` (14 項情境)

| # | 需求 (Requirement) | 情境 (Scenario) | 對應驗證實作 / 測試完整路徑與函式 | 本輪實跑證據狀態 |
|---|---|---|---|---|
| 1 | Stable circuit membership and independent profiles | Power and energy on one circuit | `internal/datalink/workspace/write_group_circuit_profile_test.go` (`TestCircuitProfileCanonicalIdentity`), `internal/api/router_studio_v2_write_group_profile_test.go` (`TestStudioV2WriteGroupsEnsureCircuitProfileRoundTrip`) | **Go 本輪實跑 PASS** |
| 2 | Stable circuit membership and independent profiles | One device contains several circuits | `internal/datalink/workspace/write_group_circuit_profile_test.go` (`TestCircuitProfileCanonicalIdentity`), `frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx` | **Go 實跑 PASS** + **前端 final PASS** |
| 3 | Stable circuit membership and independent profiles | Unknown semantics and duplicate requests | `internal/datalink/workspace/write_group_circuit_profile_rules_test.go` (`TestCircuitProfileRejectsInvalidIdentity`, `TestCircuitProfileDuplicateOwnership`), `internal/api/router_studio_v2_write_group_profile_test.go` (`TestStudioV2WriteGroupsProfileConflictAndLegacyCompatibility`) | **Go 本輪實跑 PASS** |
| 4 | Bounded recording intervals independent of acquisition | Interval boundaries | `internal/datalink/workspace/write_group_recording_profile_test.go` (`TestRecordingProfileIntervalBounds`), `internal/api/router_studio_v2_write_group_profile_test.go` (`TestStudioV2WriteGroupsProfileIntervalBoundaries`), `frontend/tests/unit/workbench-v2/profile-boundaries.test.ts` | **Go 實跑 PASS** + **前端 final PASS**（真瀏覽器越界校驗） |
| 5 | Bounded recording intervals independent of acquisition | Slow acquisition | `internal/datalink/runtime/group_boundary_long_test.go` (`TestInstantaneousProfileWindowEdges`, `TestInstantaneousProfileEligibility`) | **Go 本輪實跑 PASS** |
| 6 | Deterministic per-source instantaneous records | Boundary and missing member | `internal/datalink/runtime/group_boundary_long_test.go` (`TestInstantaneousProfileWindowEdges`, `TestLongProfileIndependentMissing`) | **Go 本輪實跑 PASS** |
| 7 | Deterministic per-source instantaneous records | Reordering and late arrivals | `internal/datalink/runtime/group_boundary_long_test.go` (`TestProfileLateSample`) | **Go 本輪實跑 PASS** |
| 8 | Cumulative raw values and exact differences | First reading and decimal progression | `internal/datalink/measurement/exact_arith_test.go`, `internal/datalink/dbtarget/circuit_recording_integration_test.go` (`TestCircuitRecordingIntegration`) | **5.1 矩陣實跑 PASS** (12000.0→12000.4→12000.7 跨 3 模式等價) |
| 9 | Cumulative raw values and exact differences | Exact large integer | `internal/datalink/aggregation/profile_counter_test.go` (`TestProfileCounterExactDelta`), `internal/datalink/grouppipeline/circuit_recording_integration_test.go` (`TestCircuitRecordingIntegration`) | **5.1 矩陣實跑 PASS** (uint64 9007199254740992→993, delta 1) |
| 10 | Counter discontinuities and gaps remain explicit | Reset inside a recording window | `internal/datalink/aggregation/profile_counter_discontinuity_test.go` (`TestProfileCounterResetWithinWindow`), `internal/datalink/grouppipeline/circuit_recording_integration_test.go` (`TestCircuitRecordingIntegration/counter_reset_within_window/*`) | **5.1 矩陣實跑 PASS** (reset_or_uncertain, delta NULL) |
| 11 | Counter discontinuities and gaps remain explicit | Signed and rollover evidence | `internal/datalink/aggregation/profile_counter_discontinuity_test.go` (`TestProfileSignedAndRollover`) | **Go 本輪實跑 PASS** |
| 12 | Counter discontinuities and gaps remain explicit | Missing window and changed meter | `internal/datalink/aggregation/profile_counter_gap_test.go` (`TestProfileCounterGapAndEpoch`), `internal/datalink/grouppipeline/circuit_recording_gap_integration_test.go` (`TestCircuitRecordingIntegration/gap_span_total/*`) | **5.1 矩陣實跑 PASS** (gap, span_total=true) |
| 13 | Atomic baseline and delivery progress | Crash at closure and uncertain SQL acknowledgement | `internal/datalink/grouppipeline/long_recovery_test.go` (`TestProfileClosureCrashRecovery`, `TestProfileCommitAckLost`), `internal/datalink/grouppipeline/circuit_recording_integration_test.go` (`TestCircuitRecordingIntegration/sql_ack_lost_and_replay_dedup/*`, `closure_crash_and_recovery/*`) | **5.1 矩陣實跑 PASS** (重送不重複前移基準) |
| 14 | Atomic baseline and delivery progress | Offline restart | `internal/datalink/grouppipeline/long_recovery_test.go` (`TestProfileOfflineRestart`) | **Go 本輪實跑 PASS** |

---

## 二、Spec 2: `guided-recording-workflow` (22 項情境)

| # | 需求 (Requirement) | 情境 (Scenario) | 對應驗證實作 / 測試完整路徑與函式 | 本輪實跑證據狀態 |
|---|---|---|---|---|
| 15 | Four-step intent-led setup | Mixed sensor configured without SQL | `frontend/tests/unit/workbench-v2/circuit-profile-mixed-sources.test.tsx`, `frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx` | **前端 final PASS**（真瀏覽器混合感測實測） |
| 16 | Four-step intent-led setup | Reload and edited connection | `frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx` | **前端 final PASS**（真瀏覽器 reload 恢復實測） |
| 17 | Four-step intent-led setup | Unavailable protocol | `frontend/tests/unit/workbench-v2/step1-capability-and-drafts.test.tsx` (`protocol availability`) | **前端 final PASS** |
| 18 | Four-step intent-led setup | Same address on different devices | `frontend/tests/unit/workbench-v2/address-conflicts.test.ts` (`two devices may both use 40001`) | **前端 final PASS** |
| 19 | Four-step intent-led setup | Missing advanced measurement | `frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx` | **前端 final PASS** |
| 20 | Four-step intent-led setup | Offline draft navigation | `frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx` | **前端 final PASS** |
| 21 | Four-step intent-led setup | Routine mapping does not require duplicate entry | `frontend/tests/unit/workbench-v2/step3-mapping.test.tsx` | **前端 final PASS** |
| 22 | Four-step intent-led setup | Existing advanced configuration is reopened | `frontend/tests/unit/workbench-v2/circuit-profile-compatibility.test.tsx` | **前端 final PASS** |
| 23 | Four-step intent-led setup | New basic setup is resumed | `frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx` | **前端 final PASS** |
| 24 | Complete basic preparation with one effective configuration | Fresh basic recording | `frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx` (佐證: `docs/technical/evidence/circuit-recording-20261024/final-prepare-only.png`) | **前端 final PASS**（真瀏覽器全新錄製 + 截圖） |
| 25 | Complete basic preparation with one effective configuration | Inspect effective settings | `frontend/tests/unit/workbench-v2/circuit-profile-status.test.tsx` | **前端 final PASS** |
| 26 | Complete basic preparation with one effective configuration | Destination before recording choices | `frontend/tests/unit/workbench-v2/basic-recording-panel.test.tsx` (`blocks a missing destination without calling ensure or start`), `frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx` | **前端 final PASS** |
| 27 | Complete basic preparation with one effective configuration | Prepare action and focused repair | `frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx` (佐證: `docs/technical/evidence/circuit-recording-20261024/final-prepare-only.png`) | **前端 final PASS**（真瀏覽器 confirm/start + 截圖） |
| 28 | Basic missing-value policy reveals required gating | Required gap remains blocked | `frontend/tests/unit/workbench-v2/basic-required-policy.test.tsx` + `internal/datalink/snapshot/partial_test.go` (`TestExplicitPartialPolicyRequiredMemberStillGatesRow`) | **前端與Go實跑 PASS**（提示與引擎分層驗證） |
| 29 | Basic missing-value policy reveals required gating | Optional gap and recovery | `internal/datalink/snapshot/partial_test.go` (`TestExplicitPartialPolicyMissingMemberIsNullWithReason`, `TestExplicitPartialPolicyCompleteRowIsNotMarkedPartial`) + `internal/datalink/dbtarget/group_row_layout_test.go` | **Go實跑 PASS**（NULL原因與正常值分層驗證） |
| 30 | Basic missing-value policy reveals required gating | Edit existing exact managed members | `frontend/tests/unit/workbench-v2/writeGroupColumns.test.ts` + `frontend/tests/unit/workbench-v2/basic-required-policy.test.tsx` + `frontend/tests/unit/workbench-v2/circuit-profile-advanced.test.tsx` | **前端 final PASS**（exact欄位相容、同ID入口與wide required控制分層驗證） |
| 31 | Basic missing-value policy reveals required gating | Confirm unchanged schema after member policy edit | `internal/api/router_studio_v2_write_group_schema_test.go` + `frontend/tests/unit/workbench-v2/groupSchemaPanel.test.tsx` + `internal/api/router_studio_v2_circuit_profile_partial_test.go` (`TestManagedWidePolicyEditRequiresNoChangeConfirmation`) | **Go與前端實跑 PASS**（另以SQLite＋PG focused實跑policy-only空DDL與新revision明確確認） |
| 32 | Basic missing-value policy reveals required gating | Long profile has a missing source | `frontend/tests/unit/workbench-v2/circuit-profile-status.test.tsx` | **前端 final PASS** |
| 33 | Simple circuit choices and scoped batch edits | Batch edit with override | `frontend/tests/unit/workbench-v2/circuit-profile-batch.test.tsx`, `frontend/tests/unit/workbench-v2/circuit-profile-batch-scope.test.tsx` (佐證: `docs/technical/evidence/circuit-recording-20261024/final-batch-preview.png`) | **前端 final PASS**（真瀏覽器批次只套用迴路 B + 截圖） |
| 34 | Simple circuit choices and scoped batch edits | Cumulative first record | `internal/datalink/grouppipeline/circuit_recording_integration_test.go`, `frontend/tests/unit/workbench-v2/circuit-profile-batch.test.tsx` (佐證: `docs/technical/evidence/circuit-recording-20261024/final-reconstruction-warning.png`) | **5.1 矩陣實跑 PASS** + **前端 final PASS**（真瀏覽器首讀數 + 截圖） |
| 35 | Explain advanced groups and workspace responsibilities | Advanced and workspace review | `frontend/tests/unit/workbench-v2/circuit-profile-advanced.test.tsx` | **前端 final PASS** |
| 36 | Explain advanced groups and workspace responsibilities | Database preparation failed | `frontend/tests/unit/workbench-v2/circuit-profile-basic.test.tsx` | **前端 final PASS** |

---

## 三、Spec 3: `long-table-group-recording` (13 項情境)

| # | 需求 (Requirement) | 情境 (Scenario) | 對應驗證實作 / 測試完整路徑與函式 | 本輪實跑證據狀態 |
|---|---|---|---|---|
| 37 | Versioned shared typed long storage | New point without column migration | `internal/api/router_studio_v2_write_group_long_test.go` (`TestSharedStorageConfirmationScope`), `frontend/tests/unit/workbench-v2/long-layout-contract.test.ts` | **Go 實跑 PASS** + **前端 final PASS** |
| 38 | Versioned shared typed long storage | Same Tag label on different circuits | `internal/datalink/recordingplan/schema_group_long_test.go` (`TestLongLayoutExactTypes`), `internal/datalink/recordingplan/schema_group_long_postgres_test.go` (`TestLongLayoutExactTypesPostgres`) | **Go 本輪實跑 PASS** |
| 39 | Optional explicit TimescaleDB mode | Plain PostgreSQL without extension | `internal/datalink/dbtarget/timescale_test.go` (`TestTimescaleCapabilityGate`), `internal/datalink/dbtarget/circuit_recording_integration_test.go` (`TestCircuitRecordingIntegration/postgres` on PG18 Alpine) | **5.1 矩陣實跑 PASS** (無 extension 環境明選 plain 成功) |
| 40 | Optional explicit TimescaleDB mode | Hypertable and unique effect | `internal/datalink/grouppipeline/long_timescale_test.go` (`TestTimescaleSpansChunksAndReplayAddsNoRow`), `internal/datalink/grouppipeline/circuit_recording_integration_test.go` (`TestCircuitRecordingIntegration/timescale_catalog_and_hypertable/*`) | **5.1 矩陣實跑 PASS** (7日 chunk, 唯一效果) |
| 41 | Optional explicit TimescaleDB mode | Unavailable extension | `internal/datalink/dbtarget/timescale_test.go` (`TestTimescaleCapabilityGate`), `internal/api/router_studio_v2_write_group_timescale_test.go` (`TestTimescaleSelectionRejectedOnSQLite`, `TestTimescaleUnavailableBlocksWithoutChange`) | **Go 本輪實跑 PASS** (SQLite 拒絕 PASS；API 端因 TS 環境 skip，已由 task5.2 PG18 無 extension focused 實跑補足) |
| 42 | Scoped confirmation and truthful delivery evidence | Shared table with stale group confirmation | `internal/api/router_studio_v2_write_group_long_test.go` (`TestSharedStorageConfirmationScope`) | **Go 本輪實跑 PASS** |
| 43 | Scoped confirmation and truthful delivery evidence | Confirmation response lost | `internal/api/router_studio_v2_write_group_long_recovery_test.go` (`TestSchemaOperationLostResponse`), `internal/api/circuit_recording_integration_test.go` (`TestCircuitRecordingIntegration/schema_operation_lost_recovery/*`) | **5.1 矩陣實跑 PASS** (verified_digest 重放冪等) |
| 44 | Restore only previously confirmed missing managed storage | Confirmed table removed during delivery | `internal/api/router_studio_v2_write_group_recovery_test.go` (`TestLongStorageMissingRecovery`), `internal/api/circuit_recording_integration_test.go` (`TestCircuitRecordingIntegration/missing_table_reconstruction/*`) (佐證: `docs/technical/evidence/circuit-recording-20261024/final-reconstruction-warning.png`, `docs/technical/evidence/circuit-recording-20261024/final-reconstruction-warning.png`) | **5.1 矩陣實跑 PASS**（真瀏覽器自動重建 generation 2 歷史警告 + 截圖） |
| 45 | Restore only previously confirmed missing managed storage | Unsafe or disabled recovery | `internal/api/router_studio_v2_write_group_recovery_denied_test.go` (`TestLongStorageRecoveryDenied`), `internal/api/router_studio_v2_write_group_recovery_denied_pg_test.go` (`TestLongStorageRecoveryDeniedPostgres`) | **Go 本輪實跑 PASS** |
| 46 | Restore only previously confirmed missing managed storage | Permissions or SQLite file missing | `internal/api/router_studio_v2_write_group_recovery_denied_test.go` (`TestLongStorageRecoveryDenied`) | **Go 本輪實跑 PASS** |
| 47 | Explicit forward cutover preserves wide history | Successful switch with pending old rows | `internal/api/router_studio_v2_write_group_cutover_test.go` (`TestWideToLongCutover`), `internal/api/circuit_recording_integration_test.go` (`TestCircuitRecordingIntegration/cutover_drain/*`) (佐證: `docs/technical/evidence/circuit-recording-20261024/final-cutover-history.png`) | **5.1 矩陣實跑 PASS**（真瀏覽器 23:26:07 cutover + SQL partial + 截圖） |
| 48 | Explicit forward cutover preserves wide history | Failed switch and rollback | `internal/api/router_studio_v2_write_group_cutover_test.go` (`TestProfileCutoverFailure`, `TestProfileRollbackDrain`), `internal/api/router_studio_v2_write_group_cutover_guard_test.go` (佐證: `docs/technical/evidence/circuit-recording-20261024/final-cutover-history.png`) | **Go 本輪實跑 PASS**（真瀏覽器 23:28:11 rollback + 截圖） |
| 49 | Stable history access without automatic deletion | Equivalent query across storage modes | `internal/datalink/grouppipeline/circuit_recording_integration_test.go` (`TestCircuitRecordingIntegration/logical_equivalence/timescale`) | **5.1 矩陣實跑 PASS** (SQLite/PG/Timescale 等價) |

---

## 四、範圍邊界、跳過案例與排除宣告 (Scope Boundaries & Disclosures)

1. **Go 4 個 Skip 案例明列**：
   - `internal/datalink/groupdelivery::TestCrashChildProcess`：子程序崩潰輔助測試 (helper)。
   - `cmd/test_ui::TestProductionLifecycleCrashChild`：生產生命週期崩潰輔助測試 (helper)。
   - `cmd/test_ui::TestProductionEntityRowsSQLAndFrozenRestart/sqlite/quoted-distinct`：未執行 (明確標註 NOT_RUN)。
   - `internal/api::TestTimescaleUnavailableBlocksWithoutChange`：因當前環境具備 TimescaleDB extension 而跳過；同一測試另於task5.2 PG18無擴充fixture實跑PASS；其5個需要擴充的子案例則在有Timescale的全套回歸PASS。
2. **唯一排除之舊 API Baseline**：
   - `TestStudioV2RecordingRoutes_PostgresApplyCreatesAndVerifiesTheManagedSchema`：非本輪新增之情境，不列入 PASS。
3. **資料庫引擎邊界**：
   - 支援範圍：SQLite、plain PostgreSQL（17.11 與 18.6）、TimescaleDB（2.31.0）。
   - 排除範圍：MySQL 與 SQL Server 屬前期 parked 提案，本輪不提供或宣稱支援。
4. **非現場驗收聲明**：
   - 本覆蓋表係依容器化環境、暫存檔與模擬器之自動化測試及本地真瀏覽器驗收結果，**不代表真實 PLC、正式生產 SCADA、實體 LAN 或 Windows 桌面現場驗收**。
5. **無自動資料清理/壓縮作業**：
   - 長表儲存明確排除 retention policy 與 compression job 自動執行，歷史存取保留原貌。
