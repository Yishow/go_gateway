# frontend-test-suite-stability Specification Delta

## ADDED Requirements

### Requirement: Default-worker full Vitest is the release acceptance mode

The frontend release gate SHALL run the complete Vitest inventory with the repository's default worker configuration in non-watch mode. A single-worker or reduced-worker run SHALL be recorded only as a diagnostic branch and SHALL NOT establish full-suite stability. Each normal run SHALL retain its command, worker mode, discovered test inventory, exit result, timeout/assertion classification, and stdout/stderr evidence.

#### Scenario: Normal-worker full suite passes

- **GIVEN** the complete frontend test inventory is discovered
- **WHEN** `npm run test -- --run` executes with the default worker configuration
- **THEN** every discovered test passes without an external timeout or orphaned worker
- **AND** the release evidence records the normal-worker mode and complete inventory

#### Scenario: Single-worker pass does not override normal instability

- **GIVEN** a single-worker diagnostic run passes
- **AND** a normal-worker run has a timeout, assertion failure, worker hang, or mixed process evidence
- **WHEN** the stability result is classified
- **THEN** the normal-worker gate is `blocked` or `unknown`
- **AND** the single-worker result is labeled diagnosis-only
- **AND** no full-suite stability pass is reported

#### Scenario: Dynamic failure identity is preserved

- **GIVEN** two normal-worker runs fail with different files or test titles
- **WHEN** the evidence is compared
- **THEN** each run retains its own failure files and titles, discovered totals, exit code, duration, and classification
- **AND** the runner does not replace the observed inventory with a fixed expected failure list
