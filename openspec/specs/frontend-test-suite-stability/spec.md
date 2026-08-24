# frontend-test-suite-stability Specification

## Purpose

TBD - created by archiving change 'diagnose-frontend-full-suite-timeouts'. Update Purpose after archive.

## Requirements

### Requirement: The diagnostic runner reproduces and classifies dynamic full-suite failures

The diagnosis runner SHALL execute the existing frontend full-suite command under normal settings and SHALL collect the failure set observed in each run instead of assuming a fixed suite identity. For every run it SHALL write schemaVersion 1 evidence containing runId, mode, repeatIndex, command, settings, timestamps, duration, discovered counts, failure file and title, classification, process result, and preserved stdout and stderr. The runner SHALL classify each non-passing result as assertion, test-timeout, worker-hang, resource leak, environment, or unknown.

#### Scenario: normal baseline preserves each run identity

- **WHEN** the operator runs the normal full-suite phase three times with the existing settings
- **THEN** the runner writes three distinct run records with dynamic failure files and titles, discovered totals, elapsed duration, exit code, and separate test-timeout and external-watchdog fields

#### Scenario: drifting failures remain visible

- **WHEN** two normal runs report different failure files or titles
- **THEN** each observed set is retained under its own runId and the runner does not replace the sets with a fixed ten-suite list

#### Scenario: classification is explicit

- **WHEN** a run exits with an assertion, test timeout, worker hang, resource correlation, environment correlation, or insufficient evidence
- **THEN** the evidence result contains exactly the corresponding classification and the raw output needed to review that decision


<!-- @trace
source: diagnose-frontend-full-suite-timeouts
updated: 2026-08-23
code:
  - scripts/lib/FrontendVitestMatrix.psm1
  - scripts/lib/FrontendVitestFreshIsolated.psm1
  - scripts/diagnose-frontend-vitest-matrix.ps1
  - scripts/lib/FrontendVitestGroups.psm1
tests:
  - scripts/tests/diagnose-frontend-final-audit.Tests.ps1
  - scripts/tests/diagnose-frontend-process-safety.Tests.ps1
  - scripts/tests/diagnose-frontend-classification-review.Tests.ps1
  - scripts/tests/diagnose-frontend-stability-handoff.Tests.ps1
  - scripts/tests/diagnose-frontend-fresh-isolated.Tests.ps1
  - scripts/tests/diagnose-frontend-groups.Tests.ps1
  - scripts/tests/diagnose-frontend-repair-gate.Tests.ps1
  - scripts/tests/diagnose-frontend-vitest-matrix.Tests.ps1
-->

---
### Requirement: Root-cause diagnosis uses a sequential single-variable matrix

The runner SHALL execute matrix phases sequentially in this order: environment capture, normal full baseline, fresh isolated per-file repeat ten, default group, single-worker group, and hypothesis-specific one-variable branches. The fresh isolated phase SHALL derive its file set from the current observed failure set. Every branch SHALL change exactly one setting while retaining comparable files, command, and other settings. The runner SHALL record isolate, pool, worker, order, group, process-lifetime, and watchdog settings for every run and SHALL NOT modify Vitest or Vite configuration.

#### Scenario: fresh isolation tests the current observed files

- **WHEN** the normal phase produces an observed failure file set
- **THEN** the runner starts a separate fresh process for each observed file and repeats that file ten times without sharing mutable process state or output directories

#### Scenario: worker saturation is tested without becoming a fix

- **WHEN** the default group is compared with a single-worker group
- **THEN** only the worker setting differs, resource and timing evidence is recorded for both groups, and the result is labeled diagnosis evidence rather than a worker-configuration repair

#### Scenario: one-variable branches stop on an unfalsifiable hypothesis

- **WHEN** prior evidence does not identify a falsifiable isolate, order, pool, or lifecycle hypothesis
- **THEN** the runner stops before opening another branch and records an unknown or blocked outcome without changing settings or source files


<!-- @trace
source: diagnose-frontend-full-suite-timeouts
updated: 2026-08-23
code:
  - scripts/lib/FrontendVitestMatrix.psm1
  - scripts/lib/FrontendVitestFreshIsolated.psm1
  - scripts/diagnose-frontend-vitest-matrix.ps1
  - scripts/lib/FrontendVitestGroups.psm1
tests:
  - scripts/tests/diagnose-frontend-final-audit.Tests.ps1
  - scripts/tests/diagnose-frontend-process-safety.Tests.ps1
  - scripts/tests/diagnose-frontend-classification-review.Tests.ps1
  - scripts/tests/diagnose-frontend-stability-handoff.Tests.ps1
  - scripts/tests/diagnose-frontend-fresh-isolated.Tests.ps1
  - scripts/tests/diagnose-frontend-groups.Tests.ps1
  - scripts/tests/diagnose-frontend-repair-gate.Tests.ps1
  - scripts/tests/diagnose-frontend-vitest-matrix.Tests.ps1
-->

---
### Requirement: Proven test-local defects receive bounded contract-preserving repairs

A repair SHALL start only after at least two comparable observations prove a test-local lifecycle, mock cleanup, fixture ownership, or deterministic assertion root cause. The repair SHALL use characterization RED evidence followed by the smallest GREEN test-only change in an evidence-selected test, test-only harness, or proven test setup path. The repair path SHALL be one of the exact conditional test paths listed in proposal Impact, or implementation SHALL stop and run spectra-ingest before modifying a new path. The repair SHALL preserve test titles, assertions, mock request shapes, fixture values, route transitions, cleanup semantics, discovery, and line-limit rules. A repair SHALL NOT change timeout values, worker or pool settings, test configuration, dependency files, coverage, skip or retry behavior, production frontend code, or Go code.

#### Scenario: lifecycle evidence unlocks a bounded repair

- **WHEN** a repeated fresh-isolated and grouped comparison proves one test-local cleanup root cause and the characterization test fails before the change
- **THEN** the implementer makes a test-only repair, turns the characterization test green, and verifies the preserved title, assertion, mock, fixture, route, cleanup, discovery, and line-limit contracts

#### Scenario: resource-only evidence does not unlock source repair

- **WHEN** only a worker, process, heap, or host-resource correlation explains the timeout and no test-local defect is proven
- **THEN** the change records the evidence as blocked and changes no worker setting, timeout, configuration, dependency, or source code


<!-- @trace
source: diagnose-frontend-full-suite-timeouts
updated: 2026-08-23
code:
  - scripts/lib/FrontendVitestMatrix.psm1
  - scripts/lib/FrontendVitestFreshIsolated.psm1
  - scripts/diagnose-frontend-vitest-matrix.ps1
  - scripts/lib/FrontendVitestGroups.psm1
tests:
  - scripts/tests/diagnose-frontend-final-audit.Tests.ps1
  - scripts/tests/diagnose-frontend-process-safety.Tests.ps1
  - scripts/tests/diagnose-frontend-classification-review.Tests.ps1
  - scripts/tests/diagnose-frontend-stability-handoff.Tests.ps1
  - scripts/tests/diagnose-frontend-fresh-isolated.Tests.ps1
  - scripts/tests/diagnose-frontend-groups.Tests.ps1
  - scripts/tests/diagnose-frontend-repair-gate.Tests.ps1
  - scripts/tests/diagnose-frontend-vitest-matrix.Tests.ps1
-->

---
### Requirement: Stability acceptance separates focused correctness from full-suite evidence

When a bounded repair is accepted, the implementer SHALL run the affected focused tests repeatedly, the relevant default group, three normal-settings full-suite runs, frontend lint, frontend build, and make check-lines or repo-supported equivalent (Git Bash scripts/check_file_lines.sh). Acceptance SHALL require no classified timeout in the three normal full-suite runs, preserved focused and group behavior, and no out-of-scope diff. When root cause is unproven or the matrix is not reproducible, the implementer SHALL record blocked evidence and SHALL NOT claim stability.

#### Scenario: repaired contract passes the stability gate

- **WHEN** a bounded repair has passed characterization RED/GREEN and affected focused checks
- **THEN** the affected focused repeat, default group, three normal full-suite runs, npm --prefix frontend run lint, npm --prefix frontend run build, and make check-lines or repo-supported equivalent (Git Bash scripts/check_file_lines.sh) all pass without a classified timeout or out-of-scope change

#### Scenario: non-reproducible diagnosis is reported honestly

- **WHEN** the same normal command produces non-comparable failure sets or the evidence cannot distinguish a root cause
- **THEN** the report contains the preserved matrix records, marks the outcome blocked, and leaves timeout, worker, config, skip, retry, exclusion, coverage, dependency, and production paths unchanged


<!-- @trace
source: diagnose-frontend-full-suite-timeouts
updated: 2026-08-23
code:
  - scripts/lib/FrontendVitestMatrix.psm1
  - scripts/lib/FrontendVitestFreshIsolated.psm1
  - scripts/diagnose-frontend-vitest-matrix.ps1
  - scripts/lib/FrontendVitestGroups.psm1
tests:
  - scripts/tests/diagnose-frontend-final-audit.Tests.ps1
  - scripts/tests/diagnose-frontend-process-safety.Tests.ps1
  - scripts/tests/diagnose-frontend-classification-review.Tests.ps1
  - scripts/tests/diagnose-frontend-stability-handoff.Tests.ps1
  - scripts/tests/diagnose-frontend-fresh-isolated.Tests.ps1
  - scripts/tests/diagnose-frontend-groups.Tests.ps1
  - scripts/tests/diagnose-frontend-repair-gate.Tests.ps1
  - scripts/tests/diagnose-frontend-vitest-matrix.Tests.ps1
-->

---
### Requirement: Matrix process handling is isolated and resource-safe

The runner SHALL create a unique output directory for each run, SHALL track only the child process tree it starts, and SHALL terminate only that owned tree after an external watchdog boundary. The external watchdog SHALL default to 600 seconds and SHALL be parameterizable per run; it only prevents an indefinitely hung owned process and SHALL NOT alter Vitest's existing five-second test or hook timeout. The runner SHALL keep external watchdog events distinct from Vitest test-timeout events, SHALL preserve process and heap metrics when available, and SHALL surface unsafe termination or missing output as a failure. The runner SHALL NOT invoke a global process kill or delete unrelated files.

#### Scenario: owned process cleanup is bounded

- **WHEN** a child runner exceeds the external watchdog boundary
- **THEN** the runner records the owned runner and child PIDs, termination result, stdout, stderr, and classification, then cleans only its unique output directory

#### Scenario: unsafe process ownership blocks the run

- **WHEN** the runner cannot prove that a process belongs to the current run
- **THEN** it refuses termination, records an unknown or worker-hang result with the available evidence, and leaves the process for operator review without touching unrelated processes

#### Scenario: scope gate blocks an unlisted repair path

- **WHEN** dynamic evidence selects a root-relative test path outside the exact conditional path union in proposal Impact
- **THEN** implementation stops, records the blocked evidence, and runs spectra-ingest before any new path is modified

<!-- @trace
source: diagnose-frontend-full-suite-timeouts
updated: 2026-08-23
code:
  - scripts/lib/FrontendVitestMatrix.psm1
  - scripts/lib/FrontendVitestFreshIsolated.psm1
  - scripts/diagnose-frontend-vitest-matrix.ps1
  - scripts/lib/FrontendVitestGroups.psm1
tests:
  - scripts/tests/diagnose-frontend-final-audit.Tests.ps1
  - scripts/tests/diagnose-frontend-process-safety.Tests.ps1
  - scripts/tests/diagnose-frontend-classification-review.Tests.ps1
  - scripts/tests/diagnose-frontend-stability-handoff.Tests.ps1
  - scripts/tests/diagnose-frontend-fresh-isolated.Tests.ps1
  - scripts/tests/diagnose-frontend-groups.Tests.ps1
  - scripts/tests/diagnose-frontend-repair-gate.Tests.ps1
  - scripts/tests/diagnose-frontend-vitest-matrix.Tests.ps1
-->

---
### Requirement: Default-worker full Vitest is the release acceptance mode

The frontend release gate SHALL run the complete Vitest inventory with the repository's committed default worker configuration in non-watch mode. The repository default SHALL remain parallel. A repository default maxWorkers cap of 6 SHALL be authorized only after the same-inventory default baseline fails, at least two comparable explicit maxWorkers=6 full-suite runs pass, a configuration contract is RED before the setting and GREEN with value 6 after the setting, and three no-CLI full-suite runs pass on the resulting inventory. The cap SHALL be an upper bound and SHALL NOT set minWorkers, pool, isolate, or fileParallelism, increase test or hook timeouts, or add skip or retry behavior. A single-worker or reduced-worker CLI run SHALL be recorded only as a diagnostic branch and SHALL NOT establish full-suite stability. Each normal run SHALL retain its command, worker mode, discovered test inventory, exit result, timeout/assertion classification, and stdout/stderr evidence.

#### Scenario: Evidence-qualified parallel cap is accepted

- **GIVEN** the same current inventory contains 187 files and 1000 tests
- **AND** the no-flags default baseline records 991 passed tests and 9 dynamic five-second timeouts
- **AND** two explicit maxWorkers=6 full-suite runs pass 1000/1000 in 109286ms and 134108ms
- **AND** the configuration contract is RED with an undefined value before the setting and GREEN with value 6 after the setting
- **WHEN** the no-flags default command runs three times after the setting
- **THEN** every run passes on 187 files and 1001 tests at 104860ms, 135633ms, and 141675ms
- **AND** no run records a test timeout, external watchdog event, or orphaned worker
- **AND** the repository default remains parallel with maxWorkers=6 as an upper bound

#### Scenario: Normal-worker full suite passes

- **GIVEN** the complete frontend test inventory is discovered
- **WHEN** npm --prefix frontend run test -- --run executes without CLI worker flags using the committed default configuration
- **THEN** every discovered test passes without an external timeout or orphaned worker
- **AND** the release evidence records normal-worker mode, the complete inventory, and the effective repository cap

#### Scenario: Single-worker pass does not override normal instability

- **GIVEN** a single-worker diagnostic run passes
- **AND** a normal-worker run has a timeout, assertion failure, worker hang, or mixed process evidence
- **WHEN** the stability result is classified
- **THEN** the normal-worker gate is blocked or unknown
- **AND** the single-worker result is labeled diagnosis-only
- **AND** no full-suite stability pass is reported

#### Scenario: Dynamic failure identity is preserved

- **GIVEN** two normal-worker runs fail with different files or test titles
- **WHEN** the evidence is compared
- **THEN** each run retains its own failure files and titles, discovered totals, exit code, duration, and classification
- **AND** the runner does not replace the observed inventory with a fixed expected failure list

#### Scenario: Unsupported worker workaround is rejected

- **GIVEN** maxWorkers=2 exceeds the 300-second external watchdog
- **WHEN** a worker-cap candidate is evaluated
- **THEN** the candidate is rejected as a solution
- **AND** no timeout increase, minWorkers setting, pool change, isolate change, fileParallelism change, skip, retry, or dependency change is authorized by this requirement
- **AND** diagnostic runner branches remain unchanged and diagnosis-only

<!-- @trace
source: bound-vitest-default-workers
updated: 2026-08-24
code:
  - frontend/vite.config.ts
tests:
  - frontend/tests/unit/utils/viteConfig.test.ts
-->
