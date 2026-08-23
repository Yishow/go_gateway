## Context

目前前端 Vitest 完整 suite 的 failure identity 會在執行間漂移。拆分後觀測為 181/938、928 pass、10 個 5 秒 timeout、wall 78.9 秒；拆分前為 913/938、25 failures、wall 149.5 秒。focused 201/201、standalone 22/22、targeted Modbus surface 5/5、canvas overlay 8/8、SourcePreview 2/2、SourceRuleTargetDatatype 3/3，以及 default 和 single-worker 四檔 18/18 均通過。這些結果支持 diagnosis-first，但尚未證明 worker/resource saturation、test-local lifecycle、mock cleanup、order/isolation 或 environment 哪一項是 root cause。

本 change 只建立可重跑的診斷矩陣與 evidence contract。任何 test-only repair 都必須先通過單變因 root-cause gate；不能把某次出現的十個 timeout suite 當成固定 identity，也不能以 timeout、worker、config 或 coverage workaround 取代證據。

## Goals / Non-Goals

**Goals:**

- 以當次完整 suite 實際觀察到的 failure set 為輸入，不預先硬編碼 suite identity。
- 建立 normal-settings full baseline、fresh isolated per-file repeat、default group、single-worker group，以及必要的 isolate、order、pool 單變因矩陣。
- 每次執行保存可比較的 command、settings、failure、phase、process、heap、worker 與 timeout evidence。
- 將結果明確分類為 assertion、test-timeout、worker-hang、resource leak 或 environment，未知結果進入 blocked 分支。
- 只有證實 test-local lifecycle 或 cleanup root cause 時，才允許 bounded test-only repair；修復後須以 focused、group 與三次 normal full-suite evidence 驗證。
- 保持既有 test title、assertion、mock、fixture、route、cleanup 與 discovery contract。

**Non-Goals:**

- 不修改 testTimeout、hookTimeout、Vitest/Vite config、package manifest、lockfile、dependency、skip、only、retry、exclude 或 coverage。
- 不把降低 worker 數量、改 pool、改 order 或提高外部 watchdog 當成產品修復。
- 不修改 production frontend/src 或 Go code；frontend/src/setupTests.ts 只有在證據證明 test-setup cleanup 是 root cause 時，才可作為 test-only conditional candidate。
- 行數 gate 使用 make check-lines 或 repo-supported equivalent（Git Bash scripts/check_file_lines.sh）；不新增 line-limit ignore。
- 不執行與本 diagnosis 無關的 full suite；不宣稱這個 change 已完成現場、瀏覽器或部署 acceptance。

## Decisions

### Decision: 動態 failure-set 與 normal baseline

每次完整 suite run 先收集當次 discovered test count、pass/fail/timeout count、實際 failure file/title 與 runner exit outcome，再將 failure set 以穩定排序寫入 evidence。normal baseline 使用現有 frontend package script 與既有設定，不調整任何 Vitest timeout 或 worker 設定。後續矩陣只引用上一階段當次觀察到的檔案集合；當 identity 漂移時保留所有觀測，不以固定十檔替代動態結果。

Alternative rejected: 直接重跑目前列出的十個 suite。這會把一次性 failure identity 固定化，無法回答 failure 漂移是否來自 worker、order 或環境。

### Decision: 單變因矩陣與 sequential execution

矩陣依序執行，避免同時啟動多個 jsdom/Vitest worker 污染資源與 log。階段順序如下：

1. Environment capture：記錄 OS、Node/npm、Vitest 版本、CPU/memory 摘要、既有 config/package hash，以及 clean child-process baseline。
2. Normal full baseline：使用既有 normal settings 執行三次，保存每次 dynamic failure set 與 elapsed time；診斷外部 watchdog 只負責防止永久 hang，不改變 test timeout。
3. Fresh isolated per-file：由每次 normal run 的 observed failure files 去重後逐檔建立 fresh process，重複十次；每檔 run 不與另一檔共用 process、temp output 或 mutable mock state。
4. Default group：以 observed failure files 的當次順序建立 group，保持預設 worker/isolation settings，重複至 evidence 足以判斷 group effect。
5. Single-worker group：只將 worker 變因設為 single-worker，其餘保持 default group 相同；若此階段改善，仍須記錄這只是 diagnosis signal，不是修復。
6. One-variable branches：只有前述差異支持假說時，依序測試 isolate、order、pool 或 process lifetime，一次只改一個變因；若沒有可否證的假說，停止矩陣並標記 blocked。

Alternative rejected: 平行執行各組以縮短 wall time。並行本身會改變 resource pressure 與 process ordering，破壞可比較性。

### Decision: 證據 schema 與 classification gate

每筆 run evidence 至少包含下列欄位，並以 JSON 保存；schemaVersion 固定為 1：

    {
      schemaVersion,
      runId,
      mode,
      repeatIndex,
      command,
      settings: { isolate, pool, workers, order, groupFiles },
      startedAtUtc,
      durationMs,
      discovered: { files, tests, passed, failed, timedOut },
      failures: [
        { file, title, kind, message, timeoutMs, phase }
      ],
      resources: {
        runnerPid,
        childPids,
        workerCount,
        peakWorkingSetBytes,
        heapUsedBytes,
        processExitCode
      },
      result,
      artifactPaths
    }

failure file 必須是 project-root-relative path；title、kind、phase 與 runner output 保留原值。artifactPaths 只指向本次 run 的 evidence location，不得把 secrets、DSN、token 或 credential 寫入 artifact。

Classification gate 依以下順序判斷：

- assertion：process 正常結束但有 assertion failure，message 與 title 可重現，且沒有 timeout 或 worker-hang signal。
- test-timeout：runner 回報明確 test timeout，process 仍能完成或產生完整 timeout output，且不是 child process 無回應。
- worker-hang：runner 或 child process 在外部 watchdog 前後沒有 completion marker，或 process tree 持續存在且需由本次 run 的 owner cleanup。
- resource leak：fresh isolated per-file 通過，但 group/default 失敗，且 process/heap/worker 指標在同一單變因下呈現可重現增長或 saturation correlation。
- environment：同一 local settings 與 process lifecycle 無法重現，且 evidence 顯示 OS、background process、CPU/memory 或外部環境差異。
- unknown：證據同時符合多類或不足以分辨；不得進入 repair branch，結果為 blocked。

### Decision: Hypothesis falsification table

| Hypothesis | Falsification test | Supporting evidence | Consequence |
| --- | --- | --- | --- |
| Worker/resource saturation | Compare normal/default group with same-file single-worker group while recording process and heap | Group-only failures, worker-sensitive timing, and correlated resource peak repeat | Keep as diagnosis signal; no worker/config repair |
| Test-local lifecycle or cleanup | Fresh isolated per-file repeat ten times and compare same title/phase with grouped run | Same file/title fails in isolation and a minimal cleanup change removes it without config change | Permit bounded test-only repair branch |
| Order or isolation coupling | Change only isolate or order across repeated group runs | Failure follows one order/isolation setting while files and resources remain comparable | Permit test-local ordering/cleanup investigation only |
| Assertion defect | Repeat the same file in fresh process with deterministic assertion output | Same assertion/title/message repeats without timeout or worker-hang | Permit focused test repair preserving behavior contract |
| Environment effect | Repeat unchanged command after environment capture and compare process/resource state | Failure correlates with external process or host state and disappears in clean run | Record environment evidence; blocked for source repair |

A hypothesis is not considered proven from one passing or failing run. At least two comparable observations are required for a repair branch, and the branch must identify the one changed variable.

### Decision: Conditional bounded repair

When the classification gate proves test-local lifecycle, mock cleanup, fixture ownership, or deterministic assertion root cause, the implementer may modify only the smallest evidence-selected path from the proposal Impact union below. The final path set must come from dynamic evidence, not from a fixed failure identity. frontend/src/setupTests.ts remains conditional and is allowed only for demonstrable shared test setup cleanup, never for production behavior.

The proposal Impact union is the only pre-authorized conditional test path set:

- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchRuntimePhase.test.tsx
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourcePreview.test.tsx
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceRuleTargetDatatype.test.tsx
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-surface.test.tsx
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.batch-selection.test.tsx
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.canvas-overlays.test.tsx
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.conflict-queue.test.tsx
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.inline-editing.test.tsx
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.planning.test.tsx
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.rule-actions.test.tsx
- frontend/tests/unit/pages/datalink/workbench/MuiSourceIncidentDesk.reopen.test.tsx

frontend/src/setupTests.ts is permitted only when shared test-setup cleanup is proven as the root cause. The union is an observed cross-run union, not a fixed failure identity. If dynamic evidence selects any other path, stop before editing and run spectra-ingest to update proposal Impact and this gate; do not add a path locally.

Any repair must first add a characterization or mutation RED test for the proven defect, then implement the smallest GREEN change, and finally preserve title, assertion, mock, fixture, route, cleanup, discovery and line-limit contracts. If the only proposed repair changes timeout, worker, pool, config, skip/retry/exclude, dependency, coverage, or production code, the branch is rejected and recorded as blocked.

### Decision: Process/resource safety and rollback

The matrix runner creates a unique run directory outside the repository unless the caller supplies an explicit output path. Each run owns its child process tree and can terminate only that tree after the external watchdog; it must never call a global process kill or target unrelated PIDs. Each phase is sequential, records stdout/stderr and exit code, and cleans only its own temp directory after evidence persistence.

The runner never edits frontend/package configuration, Vitest/Vite settings, lockfiles, line-limit ignore rules, or production sources. If a process cannot be terminated safely, the run is blocked and the PID evidence is retained for operator review. A source repair rolls back by reverting only the exact evidence-selected test/harness/setup paths; diagnosis-only script rollback removes scripts/diagnose-frontend-vitest-matrix.ps1 and its own evidence output, leaving unrelated work untouched.

## Implementation Contract

### Observable behavior

An operator can invoke the diagnosis script with an explicit mode and output path, receive a bounded result, and inspect machine-readable evidence for the exact run. Repeated normal runs report their own dynamic failure set. Per-file isolation reports each observed file separately. Group modes preserve their declared settings. No phase runs concurrently with another phase.

### Interface and data shape

The new script is scripts/diagnose-frontend-vitest-matrix.ps1. Its interface accepts Mode, Repeat, OutputPath, GroupFiles, and WatchdogSeconds. When omitted, OutputPath is a unique temporary directory and WatchdogSeconds defaults to 600 seconds. WatchdogSeconds is recorded as diagnostic metadata and never written into Vitest configuration; it only prevents an indefinitely hung owned process and does not alter the existing Vitest test or hook timeout. Supported modes are Environment, NormalFull, FreshIsolated, DefaultGroup, SingleWorkerGroup, and OneVariable. Each JSON record follows schemaVersion 1 from the evidence schema decision and includes command, settings, discovered counts, dynamic failures, resource metrics, result, and artifact paths.

The script exits zero only when the requested phase completes and writes a result record. It exits nonzero when the phase is classified as assertion, test-timeout, worker-hang, resource leak, environment, or unknown; the JSON result remains available. Process termination errors are surfaced as worker-hang or unknown with cleanup evidence.

### Failure modes

- A deterministic assertion is surfaced with file, title, message, phase, and exit code.
- A test timeout is surfaced separately from an external watchdog event.
- A worker hang includes runner/child PID evidence and cleanup outcome.
- Resource or environment classifications include the compared runs and metrics that caused the classification.
- Unclassifiable evidence is explicitly blocked; the runner never silently treats it as a pass.
- Missing dependencies, malformed output, unsafe process ownership, or missing output path fail before any source repair is attempted.

### Acceptance criteria

- Dynamic failure-set records contain no hard-coded ten-suite assumption and can be compared by runId and repeatIndex.
- Normal full baseline, fresh isolated per-file repeat ten, default group, and single-worker group outputs are complete or explicitly blocked with preserved stdout/stderr.
- One-variable branches change exactly one setting and record the comparison pair.
- Focused evidence remains separate from full-suite classification.
- Any bounded repair passes affected focused repeat, group validation, three normal full-suite runs without classified timeout, frontend lint, frontend build, and make check-lines or repo-supported equivalent (Git Bash scripts/check_file_lines.sh).
- If no root cause is proven, proposal tasks end in a blocked/no-code outcome with evidence and no timeout or configuration workaround.
- No unrelated path, production source, dependency, timeout, skip, retry, exclusion, coverage, or line-limit ignore change is accepted.

### Scope boundaries

In scope are the diagnosis script, its machine-readable evidence contract, exact evidence-selected test-only cleanup or assertion changes after proof, and the proposal/design/spec/tasks artifacts. Out of scope are production code, frontend configuration, package or lockfile changes, worker reduction as a fix, test discovery changes, broad refactors, and any change based only on a single unstable run.

## Risks / Trade-offs

- [Risk] Repeated full suites consume substantial time and host resources. → Matrix phases are sequential, bounded, and stop when a hypothesis is falsified or evidence is blocked.
- [Risk] External watchdog output could be mistaken for Vitest timeout evidence. → Store separate runner timeout and test timeout fields with distinct classification values.
- [Risk] Dynamic failure files change between runs. → Store every observed set and derive later phases from the current run rather than a permanent list.
- [Risk] Resource metrics differ by host. → Record OS, runtime, process and heap context and classify environment when comparability is insufficient.
- [Risk] A test-only cleanup change can alter semantics. → Require RED/GREEN characterization, contract parity checks, focused/group tests and three normal full-suite runs before acceptance.

## Migration Plan

1. Add the diagnosis script without changing existing frontend settings.
2. Capture environment and execute the sequential matrix, preserving JSON, stdout and stderr.
3. Classify results and choose either the conditional bounded repair branch or the blocked/no-code branch.
4. If a repair is proven, run characterization RED, minimal GREEN change, focused/group checks, three normal full-suite runs, lint, build and make check-lines or repo-supported equivalent (Git Bash scripts/check_file_lines.sh).
5. Review evidence and exact diff; rollback only evidence-selected paths if the contract fails.

## Execution Decisions

- The first matrix host is the apply host. Environment capture records OS, Node/npm, Vitest, CPU/memory, background-process snapshot, config/package hashes, and host identity before any matrix phase.
- Reporter capability is detected during the environment phase. Supported phase-timing and worker metrics are recorded; unsupported metrics are recorded as unavailable and are never inferred.
- The external watchdog defaults to 600 seconds and is parameterizable per run. It only prevents an indefinitely hung owned process and is recorded separately; it does not change the existing Vitest 5-second test timeout or hook timeout.

