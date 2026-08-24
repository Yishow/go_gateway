## Context

frontend-test-suite-stability 目前把完整 Vitest 的 release acceptance 定義為 repository default worker configuration。現有同 inventory 證據顯示，187 個 test files、1000 個 tests 的 default no-flags NormalFull 只有 991 pass，並出現 9 個動態 5 秒 test timeouts；explicit maxWorkers=6 在相同條件下兩次通過，耗時分別為 109286ms 與 134108ms。maxWorkers=2 反而超過 300 秒 external watchdog，因此問題不能以任意降低並行度或延長 timeout 處理。

這是一個小型、證據限定的前端測試穩定性 Bug Fix。現有 unstaged source patch 已將目標縮小到 Vite config 與其 contract test；本 design 只定義該 patch 應滿足的 durable handoff。diagnostic runner 的 NormalFull、DefaultGroup、SingleWorkerGroup 與 OneVariable 行為仍是診斷資料，不在本 change 修改。

## Goals / Non-Goals

**Goals:**

- 將 repository default 的 Vitest maxWorkers 明確固定為平行上限 6，讓無 CLI worker flags 的 release command 使用該設定。
- 以 config contract test 鎖定 test config 的 observable value，並保留 TDD RED/GREEN 證據。
- 把 default baseline failure、兩次 explicit cap full pass、contract RED/GREEN、三次 no-CLI full pass 組成不可省略的 acceptance gate。
- 保持 187-file inventory、既有 assertions、test and hook timeouts、file parallelism 與診斷分支語意不變。
- 將低核心主機上的 6 解釋為最大 worker 上限，而不是要求建立六個 worker。

**Non-Goals:**

- 不修改 diagnostic scripts、NormalFull evidence schema、DefaultGroup、SingleWorkerGroup、OneVariable 或其他 diagnosis branch。
- 不設定 minWorkers、pool、isolate、fileParallelism，不提高 timeout，不新增 skip、retry、coverage、exclude 或 dependency。
- 不修正 SourceStep mocks、QueryClient、product frontend code、legacy /studio、Go code 或其他 capability。
- 不把 maxWorkers=2、single-worker pass 或任何低並行診斷結果當作 release stability 證據。
- 不宣稱本 proposal 的 source patch 已完成 live/field acceptance；apply 時仍須重跑本 design 的 validation gate。

## Decisions

### Preserve the no-CLI release gate while bounding repository default parallelism

The release command remains npm --prefix frontend run test -- --run without worker flags. frontend/vite.config.ts supplies test.maxWorkers = 6, so the command remains a default-settings run while limiting concurrency. The setting is an upper bound: Vitest retains parallel file execution and uses no more than six workers, subject to host resources. No minWorkers or single-worker mode is introduced.

This choice is preferred over CLI-only flags because release and local default executions must exercise the same committed configuration. It is preferred over maxWorkers=2 because the same-inventory two-pass evidence for 6 is successful while the 2-worker branch exceeded the 300-second watchdog.

### Gate the configuration change on evidence and a TDD contract

The implementation must retain the evidence chain: a same-inventory default baseline failure; two explicit maxWorkers=6 full-suite passes; a config contract that is RED before the setting exists and GREEN with value 6 after it is added; then three no-CLI full-suite passes. The contract test calls the existing Vite config factory through the existing ambient-environment isolation helper and asserts config.test.maxWorkers equals 6.

This gate is preferred over inferring a fix from a single-worker or reduced-worker pass. A missing test result, mixed assertion/timeout result, watchdog event, orphaned process, or changed inventory blocks acceptance rather than being converted into a pass.

### Keep diagnostics and low-core behavior fail-closed

The diagnosis runner remains unchanged and continues to label worker branches as diagnosis-only. Its SingleWorkerGroup arguments may continue to use no-file-parallelism, maxWorkers=1, and minWorkers=1 for hypothesis testing; those arguments must not be copied into the repository config. The repository config must not set minWorkers, pool, isolate, or fileParallelism. On hosts with fewer than six effective cores, Vitest treats 6 as a ceiling and schedules only available workers.

This preserves the distinction between proving a test-local defect, comparing a hypothesis, and authorizing a repository default. It also prevents a future operator from interpreting the low-core ceiling as a required worker count.

## Implementation Contract

- Behavior: With no CLI worker options, npm --prefix frontend run test -- --run loads the repository Vite config and runs the complete discovered frontend inventory in parallel with maxWorkers capped at 6. The command must not need a new environment variable or script flag.
- Interface / data shape: frontend/vite.config.ts exports test.maxWorkers as the number 6. The test config must not add minWorkers, pool, isolate, or fileParallelism. frontend/tests/unit/utils/viteConfig.test.ts contains a named contract assertion for the value 6 and uses the existing buildConfigWithoutAmbientEnv helper.
- Failure modes: If the contract returns undefined or another value, the focused config test fails. If a full run has an assertion, test-timeout, worker-hang, resource, watchdog, orphan, environment, or unknown result, that run is not a stability pass and the gate remains blocked. A maxWorkers=2 run that exceeds the 300-second watchdog is retained as rejected evidence, not a fallback.
- Acceptance criteria: The pre-change baseline is recorded as 187 files / 1000 tests with 991 pass and 9 dynamic five-second timeouts; explicit maxWorkers=6 passes twice with 1000/1000 at 109286ms and 134108ms; the config contract is RED with undefined before the setting and GREEN with 6 after it; the no-CLI command passes three times on 187 files / 1001 tests at 104860ms, 135633ms, and 141675ms with zero timeout or watchdog events; frontend lint, frontend build, line-limit validation, and git diff check pass.
- Scope boundaries: Source changes are limited to frontend/vite.config.ts and frontend/tests/unit/utils/viteConfig.test.ts. The capability delta is limited to frontend-test-suite-stability. No diagnostic script, test fixture, product component, route, dependency manifest, or unrelated staged file is in scope.

## Risks / Trade-offs

- [Risk] A fixed cap can be conservative on high-core hosts. → Mitigation: preserve parallel execution and document 6 as an upper bound; later tuning requires a new evidence-backed change.
- [Risk] Dynamic test inventory or failure identity can drift between runs. → Mitigation: require the same 187-file inventory for baseline and explicit-cap comparison, and reject missing, extra, mixed, or watchdog evidence.
- [Risk] The config contract could pass while a full run is still unstable. → Mitigation: require three no-CLI full-suite passes plus lint, build, line-limit, and diff checks after the contract turns GREEN.
- [Risk] Diagnosis-only settings could leak into the repository config. → Mitigation: review the final config for absence of minWorkers, pool, isolate, and fileParallelism and keep the runner files outside the allowed source scope.

## Migration Plan

No data or dependency migration is required. Apply consists of the two scoped source edits and focused/full validation. Rollback removes only the repository maxWorkers=6 setting and its contract assertion if the evidence gate is invalidated; the canonical capability delta must then be re-evaluated before another worker setting is proposed.

## Open Questions

None for this bounded change. Any different worker cap, timeout policy, pool mode, or diagnostic-script modification requires a separate evidence-backed proposal.
