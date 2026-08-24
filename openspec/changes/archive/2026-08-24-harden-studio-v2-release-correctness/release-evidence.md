# Release evidence ledger

Date: 2026-08-24
Change: `harden-studio-v2-release-correctness`
Repository: `C:\AIProject\go_gateway`
Fixed point: `7ad76e54c4ca6488463b878c3f396613b6c83502`
Overall classification: **BLOCKED (fail closed)**

This ledger records the action-time command, working directory, execution mode, exit/result, and classification for each release gate. A focused or single-worker result never overrides a blocked normal-worker gate. Generated Playwright reports, `.last-run.json`, and `.antigravitycli` state are inventory observations only; they are not source-of-truth trace paths.

## Evidence provenance and phase boundary

The baseline, focused-contract, source/delivery, SQLite, and external tables below retain the 2026-08-24 first-round pre-repair / partial-pre-repair review evidence, including `42/42`, `97`, `242`, `3`, and `975`. Those rows are historical review results, not post-repair acceptance results. The independent post-repair verifier results are recorded in a separate section below and do not overwrite the first-round rows. The overall classification remains **BLOCKED**.

## Baseline and ownership inventory

| Check | Command | CWD | Mode | Exit/result | Classification |
| --- | --- | --- | --- | --- | --- |
| Fixed point | `git rev-parse HEAD` | `C:\AIProject\go_gateway` | read-only baseline | `7ad76e54c4ca6488463b878c3f396613b6c83502` | PASS |
| Active change status | `spectra list --json` | `C:\AIProject\go_gateway` | read-only Spectra inventory | `[]`; the change is already archived | PASS (status evidence) |
| Predecessor status | `rg --files openspec/changes/archive | rg fix-protocol-address-adaptation-v2` | `C:\AIProject\go_gateway` | read-only archive inventory | predecessor archive is present under `2026-08-24-fix-protocol-address-adaptation-v2` | PASS (artifact status) |
| Staged inventory | `git diff --cached --name-status` | `C:\AIProject\go_gateway` | read-only Git inventory | only `D .antigravitycli/fd0ca231-1a9a-4e65-8569-14c49e7cfa1d.json` | PASS (foreign path excluded from A) |
| Unstaged inventory | `git diff --name-status` | `C:\AIProject\go_gateway` | read-only Git inventory | source, test, and canonical-spec paths owned by the concurrent implementation work | PASS (ownership recorded) |
| Untracked inventory | `git ls-files --others --exclude-standard` | `C:\AIProject\go_gateway` | read-only Git inventory | archive/canonical docs, source tests, and generated Playwright artifacts | PASS (generated artifacts retained) |

The previous acceptance run reported identical source state before and after its gates. This documentation correction does not stage, restore, delete, or otherwise alter the staged `.antigravitycli` path, source files, concurrent tests, or generated artifacts.

## Focused contract evidence

| Gate | Command | CWD | Mode | Exit/result | Classification |
| --- | --- | --- | --- | --- | --- |
| Frontend address/source-rule/Step 2 tests | `npm --prefix frontend run test -- --run tests/unit/utils/addressParser.test.ts tests/unit/workbench-v2/sourceRule.test.ts tests/unit/workbench-v2/step2-rule.test.tsx` | `C:\AIProject\go_gateway` | focused Vitest | `42/42` | PASS |
| Source-rule package | `go test ./internal/datalink/sourcerule/...` | `C:\AIProject\go_gateway` | focused Go package | `97 tests` | PASS |
| Studio V2 source-rule handlers | `go test ./internal/api/handlers/...` | `C:\AIProject\go_gateway` | focused Go package | `242 tests` | PASS |
| SQLite migrator | `go test ./internal/datalink -run "TestMigrator_SQLiteSourceRule(Migration_PartialSharePresenceTable|RollbackRequiresBackupRestore|RevisionPartialUpgradeBackfillsBlankValues)" -count=1` | `C:\AIProject\go_gateway` | focused Go migration tests | exit `0`; `3` test functions selected | PASS (pre-repair review evidence) |
| Embedded SQLite foreign-key contract | `go test ./internal/datalink -run TestDefaultEmbeddedSQLiteDSNEnablesForeignKeysOnReusedConnection -count=1` | `C:\AIProject\go_gateway` | focused Go test with task-specific writable `GOCACHE` | exit `0`; `PRAGMA foreign_keys=1` on reused connections and child-row cascade pass | PASS (latest verified result) |
| Studio V2 device-owned source-rule cascade | `go test ./internal/api/handlers -run TestStudioV2WorkspaceDevicesHandler_DeleteCascadesOwnedSourceRulesInSQLite -count=1` | `C:\AIProject\go_gateway` | focused handler test with task-specific writable `GOCACHE` | exit `0`; deleted device rule removed, peer device/rule retained, workspace order isolated | PASS (latest verified result) |

These results prove only the focused contracts. They do not prove a normal-worker frontend release, a formal line-script run, or field acceptance.

## Source and delivery gates

| Gate | Command | CWD | Mode | Exit/result | Classification |
| --- | --- | --- | --- | --- | --- |
| Frontend lint | `npm --prefix frontend run lint` | `C:\AIProject\go_gateway` | normal local | exit `0` | PASS |
| Frontend default-worker full suite | `npm --prefix frontend run test -- --run` | `C:\AIProject\go_gateway` | default workers, non-watch | exit `1`; `185` files / `987` tests, `973` pass, `14` timeout failures across `11` files | **BLOCKED** |
| Frontend build | `npm --prefix frontend run build` | `C:\AIProject\go_gateway` | normal local | exit `0`; `3310` modules | PASS |
| Go full test | `go test ./...` | `C:\AIProject\go_gateway` | normal local | `975` tests / `45` packages | PASS |
| Go vet | `go vet ./...` | `C:\AIProject\go_gateway` | normal local | exit `0` | PASS |
| Go lint | `golangci-lint run ./...` | `C:\AIProject\go_gateway` | local `golangci-lint` v1.64.8 | exit `3`; repository config is v2 | **BLOCKED** |
| Formal line gate | `bash ./scripts/check_file_lines.sh` | `C:\AIProject\go_gateway` | Windows Bash | launch failed with `E_ACCESSDENIED` | **BLOCKED** |
| Line-gate equivalent | PowerShell equivalent read-only audit | `C:\AIProject\go_gateway` | environment fallback | `9` files, no hard failure, `3` warnings; `sourceCanvasModel.ts` `715 -> 715` | PASS (partial only) |
| Diff whitespace gate | `git diff --check -- <owned canonical docs>` plus `git diff --no-index --check -- NUL <owned untracked docs>` | `C:\AIProject\go_gateway` | read-only | exit `0`; no whitespace diagnostics | PASS |
| Named Spectra analysis | `spectra analyze harden-studio-v2-release-correctness --json` | `C:\AIProject\go_gateway` | normal Spectra | active change not found because it is archived | **BLOCKED** |
| Named strict validation | `spectra validate harden-studio-v2-release-correctness --strict` | `C:\AIProject\go_gateway` | normal Spectra | active change not found because it is archived | **BLOCKED** |
| All-spec strict validation | `spectra validate --all --strict --json` | `C:\AIProject\go_gateway` | normal Spectra | exit `0`; output `[]` | PASS |
| Normal Spectra validation | `spectra validate` | `C:\AIProject\go_gateway` | normal Spectra | exit `0` | PASS |
| Embedded build (sandbox attempt) | `pwsh -NoProfile -File scripts/build.ps1` | `C:\AIProject\go_gateway` | restricted local | Node/esbuild child process failed with `spawn EPERM` | BLOCKED (environment attempt) |
| Embedded build (release attempt) | `pwsh -NoProfile -File scripts/build.ps1` | `C:\AIProject\go_gateway` | elevated Windows | frontend sync completed and `bin/test-ui.exe` produced | PASS |
| Embedded Playwright smoke | `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts` | `C:\AIProject\go_gateway\frontend` | clean embedded-server smoke | `6/6`: `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, `/gateway/quick-setup`, and missing-asset `404` | PASS |

The elevated embedded build and Playwright smoke are source/delivery evidence only. They are not proof of a PLC connection, actual database target, GUI acceptance, network reachability, or factory operation.

## SQLite rollback boundary

The focused migration evidence covers row preservation, partial/repeated upgrade behavior, and safe defaults. The current SQLite path has no executable, validated `down` capability. The formal rollback is an operator-confirmed pre-upgrade database backup followed by restore when the target SQLite version cannot safely drop the added columns. No SQLite `down` operation is claimed as executed or accepted, and no rollback step may rewrite or delete legacy rows.

## External acceptance gates

| Gate | Command/action | CWD | Mode | Exit/result | Classification |
| --- | --- | --- | --- | --- | --- |
| Physical PLC communication | no action-time run recorded | N/A | external field check | not run | **BLOCKED / UNVERIFIED** |
| Actual database target | no action-time run recorded | N/A | external deployment check | not run | **BLOCKED / UNVERIFIED** |
| GUI interaction | no action-time run recorded | N/A | external operator check | not run | **BLOCKED / UNVERIFIED** |
| Network reachability | no action-time run recorded | N/A | external environment check | not run | **BLOCKED / UNVERIFIED** |
| Factory-field acceptance | no action-time run recorded | N/A | external owner acceptance | not run | **BLOCKED / UNVERIFIED** |

No aggregate release pass is claimed while any required source gate or external acceptance gate remains blocked or unverified.

## Post-repair verifier evidence (2026-08-24)

This section records the second action-time verifier phase after the documentation and source repairs. It is independent of the pre-repair rows above. The verifier reported identical before/after staged, unstaged, and untracked inventories; no source, documentation, staged, or generated artifact was changed by verification. A focused pass, an isolated single-test pass, or an elevated build pass does not override a blocked full-suite gate.

| Gate | Command/action | CWD | Mode | Exit/result | Classification |
| --- | --- | --- | --- | --- | --- |
| Before/after Git inventories | `git status --short`; `git diff --cached --name-status`; `git diff --name-status`; `git ls-files --others --exclude-standard` | `C:\AIProject\go_gateway` | read-only verifier inventory | before and after identical; staged inventory still only `.antigravitycli` deletion | PASS (inventory evidence) |
| Focused frontend contracts | `npm --prefix frontend run test -- --run tests/unit/utils/addressParser.test.ts tests/unit/workbench-v2/selectors.test.tsx tests/unit/workbench-v2/sourceRule.test.ts tests/unit/workbench-v2/step2-rule.test.tsx tests/unit/workbench-v2/step3-mapping.test.tsx tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx tests/unit/workbench-v2/mapping-autosave-page.hydration.test.tsx tests/unit/workbench-v2/mapping-autosave-page.validation.test.tsx` | `C:\AIProject\go_gateway` | focused Vitest | exit `0`; `8` files, `64/64` | PASS |
| Frontend lint | `npm --prefix frontend run lint` | `C:\AIProject\go_gateway` | normal local | exit `0`; lint clean | PASS |
| Frontend default-worker full suite | `npm --prefix frontend run test -- --run` | `C:\AIProject\go_gateway` | default workers, non-watch | exit `1`; `185` files (`151` pass / `34` fail), `991` tests (`917` pass / `74` timeout), plus `1` unhandled `[vitest-worker]: Timeout calling onTaskUpdate` | **BLOCKED** |
| Frontend build | `npm --prefix frontend run build` | `C:\AIProject\go_gateway` | normal local | exit `0`; Vite `3310` modules | PASS |
| Focused Go packages, first attempt | `go test ./internal/datalink/sourcerule ./internal/api/handlers ./internal/datalink -count=1` | `C:\AIProject\go_gateway` | local Go test with default cache | exit `1`; Go cache `Access is denied` | **BLOCKED** (environment attempt) |
| Focused Go packages, rerun | `go test ./internal/datalink/sourcerule ./internal/api/handlers ./internal/datalink -count=1` | `C:\AIProject\go_gateway` | rerun with temporary task cache | exit `0`; all `3` packages pass | PASS |
| Go full test | `go test ./...` | `C:\AIProject\go_gateway` | normal local | exit `1`; `990` pass, `1` timeout failure, `4` skipped across `45` packages; `dbtarget/TestWriter_InsertModeRowGroupsEmitSeparateRowsForSharedColumn` failed in the full run | **BLOCKED** (full-suite FAIL) |
| Isolated full-test failure rerun | `go test ./internal/datalink/dbtarget -run "TestWriter_InsertModeRowGroupsEmitSeparateRowsForSharedColumn" -count=1` | `C:\AIProject\go_gateway` | diagnostic single-test rerun | exit `0`; isolated test passes, which does not override the full-suite failure | PASS (diagnostic only) |
| Go vet | `go vet ./...` | `C:\AIProject\go_gateway` | normal local | exit `0`; clean | PASS |
| Go lint | `golangci-lint run ./...` | `C:\AIProject\go_gateway` | local `golangci-lint` v1.64.8 | exit `3`; installed v1 cannot read repository v2 config | **BLOCKED** |
| Formal line gate | `bash ./scripts/check_file_lines.sh` | `C:\AIProject\go_gateway` | Git Bash | exit `1`; `23` files checked; `frontend/tests/unit/workbench-v2/sourceRule.test.ts` is `523` lines and exceeds hard `500` | **BLOCKED** (formal FAIL) |
| PowerShell line-gate equivalent | PowerShell equivalent line audit | `C:\AIProject\go_gateway` | environment fallback | exit `1`; same `523`-line failure plus `6` warnings; `internal/datalink/sourcerule/service.go` remains legacy `1242 -> 1242` | **BLOCKED** |
| Historical scope line boundary | `sourceCanvasModel.ts` line audit | `C:\AIProject\go_gateway` | read-only verifier audit | fixed point `715 -> 715`, delta `0` | PASS |
| Diff whitespace gate | `git diff --check HEAD -- .` (reports/artifacts excluded from the verifier comparison) | `C:\AIProject\go_gateway` | read-only | exit `0`; source/docs whitespace clean | PASS |
| All-spec strict validation | `spectra validate --all --strict --json` | `C:\AIProject\go_gateway` | normal Spectra | exit `0`; output `[]` | PASS |
| Normal Spectra validation | `spectra validate` | `C:\AIProject\go_gateway` | normal Spectra | exit `0`; clean | PASS |
| Named Spectra analysis | `spectra analyze harden-studio-v2-release-correctness --json` | `C:\AIProject\go_gateway` | normal Spectra | exit `1`; active change is unavailable because it is archived | **BLOCKED** |
| Embedded build, sandbox attempt | `pwsh -NoProfile -File scripts/build.ps1` | `C:\AIProject\go_gateway` | restricted local | exit `1`; Node/esbuild child process `spawn EPERM` | **BLOCKED** (environment attempt) |
| Embedded build, elevated attempt | `pwsh -NoProfile -File scripts/build.ps1` | `C:\AIProject\go_gateway` | elevated Windows | exit `0`; frontend sync and Go embedded `bin/test-ui.exe` completed | PASS |
| Embedded Playwright smoke and process check | `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts`; post-smoke `test-ui` process check | `C:\AIProject\go_gateway\frontend` | clean embedded-server smoke | exit `0`; `6/6` routes/assets pass, missing asset `404`, post-smoke `test-ui` process count `0` | PASS |

At the time of that second-round verifier, the post-repair full Vitest, full Go test, formal line gate, and golangci-lint gates remained incomplete or failed; this paragraph is historical and is not overwritten by the final validation below. The six line-gate warnings are `en/workbench-v2.json` (`485`), `zh-TW/workbench-v2.json` (`485`), `useStudioV2AutosaveState.ts` (`499`), `step2-rule.test.tsx` (`433`), `studio_v2_workspace_source_rules_handler_test.go` (`381`), and `migrator.go` (`451`). If `sourceRule.test.ts` is split, the line gate must be rerun and this ledger amended with a new action-time result; no resolution is claimed here.

The post-repair verifier did not establish physical PLC communication, actual database-target behavior, GUI acceptance, network reachability, or factory-field acceptance. Those external gates remain **BLOCKED / UNVERIFIED**, and no aggregate release pass is claimed.

## Historical repair validation and diagnosis before P1 (2026-08-24)

This historical third section records the primary validation and read-only dynamic Vitest diagnosis before the later P1 hydration repair. It does not overwrite the prior post-repair Go timeout row: the primary `go test ./... -count=1` rerun passed, while the earlier full-suite timeout remains preserved as historical evidence. The diagnostic artifacts were written under `C:\Users\Yishow\AppData\Local\Temp\frontend-vitest-hardening-85681beed6c54555a830fa8abb5a432c`; repository status, tracked-file, and frontend-test inventories remained unchanged. The final P1 results are recorded separately below.

| Gate | Command/action | CWD | Mode | Exit/result | Classification |
| --- | --- | --- | --- | --- | --- |
| Focused Vitest sandbox attempt | `npm --prefix frontend run test -- --run tests/unit/utils/addressParser.test.ts tests/unit/workbench-v2/selectors.test.tsx tests/unit/workbench-v2/sourceRule.test.ts tests/unit/workbench-v2/sourceRule-readiness.test.ts tests/unit/workbench-v2/step2-rule.test.tsx tests/unit/workbench-v2/step3-mapping.test.tsx tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx tests/unit/workbench-v2/mapping-autosave-page.hydration.test.tsx tests/unit/workbench-v2/mapping-autosave-page.validation.test.tsx` | `C:\AIProject\go_gateway` | restricted local | child process failed with `spawn EPERM` before test execution | **BLOCKED** (environment attempt; excluded from test conclusion) |
| Primary focused Vitest contracts | `npm --prefix frontend run test -- --run tests/unit/utils/addressParser.test.ts tests/unit/workbench-v2/selectors.test.tsx tests/unit/workbench-v2/sourceRule.test.ts tests/unit/workbench-v2/sourceRule-readiness.test.ts tests/unit/workbench-v2/step2-rule.test.tsx tests/unit/workbench-v2/step3-mapping.test.tsx tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx tests/unit/workbench-v2/mapping-autosave-page.hydration.test.tsx tests/unit/workbench-v2/mapping-autosave-page.validation.test.tsx` | `C:\AIProject\go_gateway` | elevated Windows | exit `0`; `9` files, `64/64` tests pass | PASS |
| Primary focused Go packages | `go test ./internal/datalink/sourcerule ./internal/api/handlers ./internal/datalink -count=1` | `C:\AIProject\go_gateway` | primary local rerun | exit `0`; all `3` packages pass | PASS |
| Primary Go full test | `go test ./... -count=1` | `C:\AIProject\go_gateway` | primary local rerun | exit `0`; all packages pass | PASS |
| Primary Go vet | `go vet ./...` | `C:\AIProject\go_gateway` | primary local rerun | exit `0`; clean | PASS |
| Primary frontend lint | `npm --prefix frontend run lint` | `C:\AIProject\go_gateway` | primary local rerun | exit `0`; clean | PASS |
| Primary formal line gate | `bash ./scripts/check_file_lines.sh` | `C:\AIProject\go_gateway` | elevated Git Bash | exit `0`; `24` files checked; six warnings; `sourceCanvasModel.ts` `715 -> 715`; legacy `service.go` `1242 -> 1242` allowed | PASS |
| NormalFull Repeat 2 diagnosis | `pwsh -NoProfile -Command '& { & "scripts\diagnose-frontend-vitest-matrix.ps1" -Mode NormalFull -Repeat 2 -WatchdogSeconds 600 -OutputPath "C:\Users\Yishow\AppData\Local\Temp\frontend-vitest-hardening-85681beed6c54555a830fa8abb5a432c\normal-full-default-repeat2" }'` | `C:\AIProject\go_gateway` | diagnostic matrix; temp artifact root recorded above | exit `1` overall; run 1 `991` tests / `985` pass / `6` timeouts in `71.3s`; run 2 `991` / `982` pass / `9` timeout identities across `8` unique files in `94.4s`; `4` exact identities overlapped; watchdog false | **BLOCKED** |
| Dynamic grouped diagnosis | matrix runner `DefaultGroup`, `SingleWorkerGroup`, and `FreshIsolated` using the NormalFull summary as the dynamic `10`-file union | `C:\AIProject\go_gateway` | diagnosis-only branches | exit `0`; DefaultGroup `45/45`, SingleWorker `45/45`, FreshIsolated `10` files / `45` tests pass | PASS (diagnosis only; not release acceptance) |
| Diagnostic process evidence | matrix artifacts and process inventory under the recorded temp root | `C:\AIProject\go_gateway` | read-only evidence audit | no `[vitest-worker]: Timeout calling onTaskUpdate` in these runs; no orphan Vitest process; sandbox direct-run `spawn EPERM` classified as environment | PASS (diagnosis evidence) |
| Go lint | `golangci-lint run ./...` | `C:\AIProject\go_gateway` | local `golangci-lint` v1.64.8 | exit `3`; installed v1 cannot read repository v2 config | **BLOCKED** |
| Named Spectra analysis | `spectra analyze harden-studio-v2-release-correctness --json` | `C:\AIProject\go_gateway` | normal Spectra | exit `1`; active change is unavailable because it is archived | **BLOCKED** |
| Named strict validation | `spectra validate harden-studio-v2-release-correctness --strict` | `C:\AIProject\go_gateway` | normal Spectra | exit `1`; active change is unavailable because it is archived | **BLOCKED** |

The final elevated line gate warnings were the two locale files at `485` lines, `useStudioV2AutosaveState.ts` at `499`, `step2-rule.test.tsx` at `433`, the handler test at `381`, and `migrator.go` at `451`; no hard failure remained. The matrix shows dynamic timeout identities and four exact overlaps, but the available runner telemetry is insufficient to attribute a deterministic source-A defect. The strongest diagnosis is shared resource/state/load interaction in the full default-worker context, so the default full release gate remains **BLOCKED**. Single-worker, grouped, and fresh-isolated passes do not change that classification.

## Previous frontend repair validation before P1 (2026-08-24)

This historical section records the primary post-repair frontend and focused-Go results before the later P1 hydration repair. It remains historical evidence and does not replace the final P1 validation below.

| Gate | Command/action | CWD | Mode | Exit/result | Classification |
| --- | --- | --- | --- | --- | --- |
| Primary focused frontend contracts | `npm --prefix frontend run test -- --run tests/unit/utils/addressParser.test.ts tests/unit/workbench-v2/selectors.test.tsx tests/unit/workbench-v2/sourceRule.test.ts tests/unit/workbench-v2/sourceRule-readiness.test.ts tests/unit/workbench-v2/step2-rule.test.tsx tests/unit/workbench-v2/step3-mapping.test.tsx tests/unit/workbench-v2/rule-tab-rail.test.tsx tests/unit/workbench-v2/hydratedProgress.test.ts tests/unit/workbench-v2/shell-readiness.test.tsx tests/unit/workbench-v2/shell-redesign.test.tsx tests/unit/workbench-v2/workspace-boot.test.tsx tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx tests/unit/workbench-v2/mapping-autosave-page.hydration.test.tsx tests/unit/workbench-v2/mapping-autosave-page.validation.test.tsx` | `C:\AIProject\go_gateway` | primary elevated Windows | exit `0`; `14` files, `95/95` tests pass | PASS |
| Primary frontend lint | `npm --prefix frontend run lint` | `C:\AIProject\go_gateway` | primary elevated Windows | exit `0`; clean | PASS |
| Primary formal line gate before P1 | `bash ./scripts/check_file_lines.sh` | `C:\AIProject\go_gateway` | primary elevated Git Bash | exit `0`; `34` files checked; nine warnings; `sourceCanvasModel.ts` `715 -> 715`; legacy `service.go` `1242 -> 1242` allowed | PASS (historical before P1) |
| Focused Go packages with default cache, attempt 1 | `go test ./internal/datalink/sourcerule ./internal/api/handlers ./internal/datalink -count=1` | `C:\AIProject\go_gateway` | primary local, default Go cache | exit `1`; `Access is denied` while using the existing cache | **BLOCKED** (environment attempt) |
| Focused Go packages with default cache, attempt 2 | same command as attempt 1 | `C:\AIProject\go_gateway` | primary local, default Go cache | exit `1`; repeated `Access is denied` | **BLOCKED** (environment attempt) |
| Focused Go packages with task-specific cache | `pwsh -NoProfile -Command '$env:GOCACHE="C:\Users\Yishow\AppData\Local\Temp\go_gateway_luna_gocache"; go test ./internal/datalink/sourcerule ./internal/api/handlers ./internal/datalink -count=1'` | `C:\AIProject\go_gateway` | primary elevated Windows, writable task-specific TEMP cache | exit `0`; all `3` packages pass | PASS |

The previously recorded primary `go test ./... -count=1` PASS and the NormalFull Repeat 2 default-worker Vitest blocker remain unchanged. The focused frontend and task-specific-cache Go results do not close the default-worker release gate.

## Final P1 and post-A validation (2026-08-24)

This final section records the P1 hydration repair and post-A source evidence. During hydration, string-valued blank `device_id` and `start_address` are retained rather than replaced by defaults; the focused tests prove `unknown_device` and `invalid_address` remain blocked and that a valid sibling rule cannot make the workspace ready. These results establish the A-class focused contract, but do not close the unrelated NormalFull, toolchain, historical TDD RED, or external/live gates.

| Gate | Command | CWD | Mode | Exit/result | Classification |
| --- | --- | --- | --- | --- | --- |
| Final P1 focused frontend contracts | `npm --prefix frontend run test -- --run tests/unit/workbench-v2/studioV2RuleAutosave.test.ts tests/unit/workbench-v2/sourceRule-readiness.test.ts tests/unit/workbench-v2/workspace-boot.test.tsx` | `C:\AIProject\go_gateway` | primary elevated Windows | exit `0`; `3` files / `18` tests pass, covering blank-string hydration, `unknown_device`, `invalid_address`, and valid-sibling readiness isolation | PASS |
| Final frontend lint | `npm --prefix frontend run lint` | `C:\AIProject\go_gateway` | primary elevated Windows | exit `0`; clean | PASS |
| Final frontend build | `npm --prefix frontend run build` | `C:\AIProject\go_gateway` | primary elevated Windows | exit `0`; Vite build pass | PASS |
| Final post-A Go full test | `go test ./... -count=1` | `C:\AIProject\go_gateway` | primary Windows with task-specific writable `GOCACHE` | exit `0`; all packages pass | PASS |
| Final post-A Go vet | `go vet ./...` | `C:\AIProject\go_gateway` | primary Windows with task-specific writable `GOCACHE` | exit `0`; clean | PASS |
| Final formal line gate | `bash ./scripts/check_file_lines.sh` | `C:\AIProject\go_gateway` | final elevated Git Bash | exit `0`; `39` files checked; same nine warnings; `sourceCanvasModel.ts` `715 -> 715`; legacy `service.go` `1242 -> 1242` | PASS |

The final nine warnings are `WorkbenchV2Shell.tsx` (`465`), `en/workbench-v2.json` (`492`), `zh-TW/workbench-v2.json` (`492`), `useStudioV2AutosaveState.ts` (`487`), `shell-readiness.test.tsx` (`346`), `step2-rule.test.tsx` (`437`), `workspace-boot.test.tsx` (`469`), `studio_v2_workspace_source_rules_handler_test.go` (`381`), and `migrator.go` (`451`). No new file exceeds the hard `500`-line limit. A-class source evidence is **GREEN**; the overall release remains **BLOCKED** only for the unrelated NormalFull default-worker Vitest gate, golangci-lint v1/v2 mismatch, archived named analysis where applicable, unchecked historical TDD RED tasks, and external/live acceptance.

## Resolved architecture decision: workspace ownership closure (Decision A)

Decision A is resolved: keep source-rule ownership indirect through the single owning device, and treat a normal workspace-device deletion as cascade removal of that device's source rules. The persisted SQLite `source_rules` table has `device_id` but intentionally has no `workspace_id` or tombstone. Workspace projection remains isolated by `OrderedDeviceIDs`, and `ListByDeviceIDs` is the only source-rule listing boundary for the workspace; the system must not search globally or enumerate rules outside those device IDs.

When normal device deletion succeeds, the existing `ON DELETE CASCADE` removes the device-owned source rules (and their dependent source-rule records), so they become unavailable rather than persisted orphans for later repair. This is the normal deleted-device lifecycle and satisfies the ownership boundary without adding cross-workspace state or global leakage. Cleanup or repair of pre-existing orphan rows is explicitly outside this change.

Frontend unknown/deleted readiness remains a defensive fail-closed check only for a stale or inconsistent source-rule record that the handler/API actually returns or that stale workspace state retains. Such a record must preserve rule/device/address identity, block Step 2, and never fall back to `devices[0]` or the first Modbus device. The check does not require a List API to bypass `OrderedDeviceIDs`/`ListByDeviceIDs` or to return rules already removed by cascade.

Rationale: this matches the existing schema and FK lifecycle, preserves workspace isolation, avoids cross-workspace/global leakage, and does not create a new `workspace_id` or tombstone contract. The backend cascade contract path is `internal/api/handlers/studio_v2_workspace_devices_handler_cascade_test.go`, and `internal/datalink/db_sqlite_fk_test.go` proves the per-connection `foreign_keys=1` DSN contract. Latest action-time runs with a task-specific writable `GOCACHE` both exit `0`: `go test ./internal/datalink -run TestDefaultEmbeddedSQLiteDSNEnablesForeignKeysOnReusedConnection -count=1` proves PRAGMA and generic child-row cascade, while `go test ./internal/api/handlers -run TestStudioV2WorkspaceDevicesHandler_DeleteCascadesOwnedSourceRulesInSQLite -count=1` proves the device-owned rule is removed and the peer workspace rule remains isolated. A3 cascade evidence is therefore **GREEN / PASS**. Cleanup or repair of pre-existing orphan rows remains explicitly outside this change. Focused package tests and `go vet` also pass, but the overall classification stays **BLOCKED (fail closed)** while the unrelated NormalFull Vitest, toolchain, and external gates remain blocked/unverified.

## Line warning rationale and split plan

The final formal line gate passed with `39` files and the same nine warnings; no newly added file exceeded `500` lines. This release-hardening round does not split them broadly because doing so would expand the A-class scope; each split is assigned to the next maintenance change, prior to further line growth:

| Warning | Rationale and owner/action |
| --- | --- |
| `frontend/src/features/datalink/workbench-v2/shell/WorkbenchV2Shell.tsx` (`465`) | The frontend shell maintainer will extract shell layout/readiness composition helpers. |
| `frontend/src/i18n/locales/en/workbench-v2.json` (`492`) and `frontend/src/i18n/locales/zh-TW/workbench-v2.json` (`492`) | Paired locale catalogs will be split by namespace together by the frontend i18n maintainer to avoid key drift. |
| `frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts` (`487`) | This change only removed the fallback caller; the Studio V2 state maintainer will extract hydration/reconciliation helpers. |
| `frontend/tests/unit/workbench-v2/shell-readiness.test.tsx` (`346`) | The frontend test maintainer will extract readiness banner/summary-rail scenarios into a dedicated test file. |
| `frontend/tests/unit/workbench-v2/step2-rule.test.tsx` (`437`) | The frontend test maintainer will extract Step 2 readiness UI scenarios into a dedicated test file. |
| `frontend/tests/unit/workbench-v2/workspace-boot.test.tsx` (`469`) | The frontend page-test maintainer will extract workspace bootstrap API/error fixtures and helpers. |
| `internal/api/handlers/studio_v2_workspace_source_rules_handler_test.go` (`381`) | The Go API maintainer will extract validation-context tests into a focused test file. |
| `internal/datalink/migrator.go` (`451`) | This change only added revision-idempotent backfill; the Go data-layer maintainer will extract a SQLite schema helper. |

The legacy `internal/datalink/service.go` fixed point remains `1242 -> 1242`; it is unchanged legacy context, not a new warning.

Physical PLC communication, actual database-target behavior, GUI acceptance, network reachability, and factory-field acceptance remain **BLOCKED / UNVERIFIED**. No aggregate release pass is claimed.
