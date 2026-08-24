# cross-platform-test-contracts Specification Delta

## ADDED Requirements

### Requirement: Staged and unstaged release baselines are checked independently

The release audit SHALL record one fixed point and SHALL inspect staged, unstaged, and untracked path inventories as separate datasets. A combined diff or a single clean working-tree result SHALL NOT substitute for the three inventories. The audit SHALL identify ownership conflicts and forbidden paths before any gate result is accepted.

#### Scenario: Independent inventories expose a staged-only path

- **GIVEN** a file exists only in `git diff --cached --name-status` and not in `git diff --name-status`
- **WHEN** the fixed-point release audit runs
- **THEN** the staged dataset records that file independently
- **AND** the unstaged dataset remains separately recorded
- **AND** the audit does not silently treat the file as an unstaged change or as absent

#### Scenario: Inventory mismatch blocks acceptance

- **GIVEN** the active predecessor or another worker changes a path after the fixed point
- **WHEN** the final staged/unstaged comparison no longer matches the recorded ownership inventory
- **THEN** the audit reports the path and mismatch category
- **AND** the release-hardening acceptance remains blocked until the owning agent reconciles it

### Requirement: Cross-platform and environment failures are classified fail closed

The release evidence SHALL distinguish source defects, assertion failures, process or worker timeouts, missing toolchains, permission failures, and unavailable external systems. A mixed or insufficient signal SHALL be classified as `blocked` or `unknown`; it SHALL NOT be converted into a pass by changing worker count, rerunning without evidence, or relying on a different platform result.

#### Scenario: Toolchain unavailability blocks the gate

- **GIVEN** a required Windows, Go, Node, Playwright, or Spectra tool is unavailable or cannot start
- **WHEN** the gate runner records the result
- **THEN** the affected gate is marked incomplete with the command and environment error
- **AND** a pass is not inferred from a focused or different-platform run

#### Scenario: Mixed assertion and timeout evidence remains unknown

- **GIVEN** one run reports assertion failures and an external watchdog reports a timeout or process cleanup uncertainty
- **WHEN** the stability matrix classifies the run
- **THEN** the result is `blocked` or `unknown`
- **AND** a single-worker pass or a later partial rerun does not replace the original evidence
