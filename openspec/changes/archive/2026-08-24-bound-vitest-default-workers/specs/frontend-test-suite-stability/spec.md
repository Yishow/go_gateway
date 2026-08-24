## MODIFIED Requirements

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
