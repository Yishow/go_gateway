## Why

Under umbrella `rule-centric-studio-refactor`, Local Modbus is the other first-class output target, but it still behaves like a runtime/manual engineering tool instead of a rule-owned review/apply destination. This phase makes Local Modbus consume the shared revision contract, persist rule-owned mapping state, and enforce allocator conflicts explicitly so it can stand beside Database output as a true sibling workflow.

## What Changes

- Execute this change as child phase 4 under umbrella `rule-centric-studio-refactor`.
- Treat `rule-driven-foundation` and `tag-review-first-flow` as prerequisites, and treat `database-output-review-apply` as a sibling phase rather than a prerequisite.
- Convert Local Modbus output into rule-scoped candidate review/apply backed by persisted ownership and revision state.
- Make conflict governance explicit for overlapping register allocations and block only the conflicting subset.
- Keep verification and health feedback inside the same review/apply surface.
- Add explicit rule-scoped Local Modbus apply APIs with per-item results.
- **BREAKING**: Local Modbus output for the primary workflow is no longer treated as an in-memory-only or manual-only binding path.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `local-modbus-memory-workbench`: Change Local Modbus from manual/runtime-only behavior to persistent rule-scoped candidate review/apply with explicit conflict governance.
- `datalink-api`: Add Local Modbus-specific rule-scoped apply contracts.

## Impact

- **Frontend**: Local Modbus candidate review board, conflict surfacing, verification-first apply, and health/result feedback.
- **Backend services**: persistent Local Modbus ownership state, conflict detection, and apply execution.
- **APIs**: Local Modbus apply endpoint and response contract.
- **Tests**: persistence across restart, conflict blocking, and partial apply result coverage.
