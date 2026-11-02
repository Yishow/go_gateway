# 本輪狹義 lint 例外（Go：19 處）

完整官方 lint 0 issues 包含下列明列例外，不是無 waiver。未禁用全域規則、整檔或整包。原基線 e42f70ef 到 c88 的新增比較；後續 FE-only 修改不變更這些 Go bytes。

| 位置 | 規則 | 具體理由及直接測試 |
| --- | --- | --- |
| `internal/datalink/device/service_crud.go:123` | `contextcheck` | Single-device retirement intentionally completes independently of request cancellation; see TestDeviceMutationCleanupDoesNotInheritRequestCancellation. |
| `internal/datalink/device/service_crud.go:154` | `contextcheck` | Single-device retirement intentionally completes independently of request cancellation; see TestDeviceMutationCleanupDoesNotInheritRequestCancellation. |
| `internal/datalink/device/service_status.go:218` | `contextcheck` | Single-device retirement intentionally completes independently of request cancellation; see TestDeviceMutationCleanupDoesNotInheritRequestCancellation. |
| `internal/datalink/modbusshare/sql_revision_dirty.go:26` | `errorlint` | Ignore only bare TxDone; wrapped/joined foreign-driver errors can carry real cleanup failures (TestWorkspaceRevisionRollbackIgnoresOnlyBareTxDone). |
| `internal/datalink/modbusshare/sql_revision_store.go:70` | `errorlint` | Ignore only bare TxDone; wrapped/joined foreign-driver errors can carry real cleanup failures (TestWorkspaceRevisionRollbackIgnoresOnlyBareTxDone). |
| `internal/datalink/modbusshare/sql_revision_store.go:191` | `errorlint` | Ignore only bare TxDone; wrapped/joined foreign-driver errors can carry real cleanup failures (TestWorkspaceRevisionRollbackIgnoresOnlyBareTxDone). |
| `internal/datalink/point/format_mutation_invalidation.go:52` | `errorlint` | Only bare ErrTxDone is benign; wrapped/joined cleanup failures must remain fatal (`TestFormatInspectForeignRollbackOnlyBareTerminalIsBenign`). |
| `internal/datalink/point/format_mutation_invalidation.go:103` | `errorlint` | Only bare ErrTxDone is benign; wrapped/joined cleanup failures must remain fatal (`TestFormatInspectForeignRollbackOnlyBareTerminalIsBenign`). |
| `internal/datalink/point/format_mutation_invalidation.go:113` | `errorlint` | Only bare no-rows proves absence; wrapped/joined driver errors remain unavailable (`TestFormatInspectOnlyBareMissingReceiptIsNotFound`). |
| `internal/datalink/point/format_mutation_sql.go:52` | `errorlint` | Only bare ErrTxDone is benign; wrapped/joined cleanup failures must remain fatal (`TestFormatInspectForeignRollbackOnlyBareTerminalIsBenign`). |
| `internal/datalink/point/format_recovery_snapshot.go:61` | `errorlint` | Only bare ErrTxDone is benign; wrapped/joined cleanup failures must remain fatal (`TestFormatInspectForeignRollbackOnlyBareTerminalIsBenign`). |
| `internal/datalink/point/sql_repo_revision.go:29` | `errorlint` | Only bare ErrTxDone is benign; wrapped/joined cleanup failures must remain fatal (`TestFormatInspectForeignRollbackOnlyBareTerminalIsBenign`). |
| `internal/datalink/runtime/service.go:304` | `contextcheck` | Workspace-less fallback: modern profiles fail closed; legacy/disabled profiles do not inherit request evidence. |
| `internal/datalink/runtime/service.go:331` | `contextcheck` | Independent fresh-profile install; TestSnapshotInstallKeepsIndependentFreshProfileBoundary proves cancellation/cache isolation. |
| `internal/datalink/runtime/service_device_sync.go:46` | `contextcheck` | Admitted install is independent; TestDeviceUpsertInstallKeepsIndependentFreshProfileBoundary proves cancellation/cache isolation. |
| `internal/datalink/sourcerule/candidate_snapshot_revision.go:28` | `errorlint` | Ignore only bare TxDone; foreign-driver wrapped/joined cleanup causes must survive (TestRuleRepositoryRollbackIgnoresOnlyBareTxDone). |
| `internal/datalink/sourcerule/repository_links_atomic.go:40` | `errorlint` | Ignore only bare TxDone; foreign-driver wrapped/joined cleanup causes must survive (TestRuleRepositoryRollbackIgnoresOnlyBareTxDone). |
| `internal/datalink/sourcerule/sql_repo.go:222` | `errorlint` | Ignore only bare TxDone; foreign-driver wrapped/joined cleanup causes must survive (TestRuleRepositoryRollbackIgnoresOnlyBareTxDone). |
| `internal/datalink/sourcerule/sql_repo.go:284` | `errorlint` | Ignore only bare TxDone; foreign-driver wrapped/joined cleanup causes must survive (TestRuleRepositoryRollbackIgnoresOnlyBareTxDone). |

另本輪 Step3 FE 身分等價 memo 新增一處逐行 `react-hooks/exhaustive-deps` 例外，依 persistedSubscriptionKey 的 local/stream/persistedID 語意穩定依賴而非 snapshot object identity；兩個舊 typed recovery/retry 回歸先 FAIL 後修復，全 13 個相關案例 PASS。位置為 `frontend/src/features/datalink/workbench-v2/steps/step3/useStep3LiveValues.ts:143`，與上述 19 處新增 Go 例外分列；該檔其他既有 semantic-key memo 註解未算成本輪新增，並非宣稱全庫只有一處 FE 例外。
