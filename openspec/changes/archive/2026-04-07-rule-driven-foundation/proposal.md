## Why

The umbrella contract `rule-centric-studio-refactor` is now stable, but later child changes still need one shared foundation for rule revisions, candidate persistence, readiness gating, and candidate APIs. Without this change, Step 3, Database output, and Local Modbus output would each invent their own state model and reopen the architecture debate the umbrella was meant to close.

## What Changes

- Execute this change as child phase 1 under umbrella `rule-centric-studio-refactor`.
- Treat this change as depending only on the umbrella contract, with no earlier child prerequisite.
- Establish source-rule revision persistence as the authoritative runtime boundary for candidate generation.
- Add persisted candidate snapshots, canonical rule-owned identity, deterministic proposed signatures, and explicit snapshot completeness rules.
- Replace one-dimensional device readiness with a planning-versus-activation/apply contract tied to separate connect/probe outcomes.
- Extend connector diagnostics with planning capability hints used by rule-driven planning.
- Extend rule-derived mapping lifecycle semantics with `out_of_sync` support and shared identity/signature comparison rules.
- Add rule-scoped candidate query and recompute APIs as the shared contract for later child changes.
- **BREAKING**: saving a source rule is no longer treated as immediate downstream persistence; it becomes a revision-and-candidate boundary.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `source-rule-runtime`: Persist source-rule revisions, candidate snapshots, and canonical diff semantics.
- `device-registry`: Change readiness from one final boolean into explicit planning, activation, and apply eligibility.
- `protocol-connectors`: Add capability hints to connector diagnostics used by rule-driven planning.
- `mapping-pipeline`: Add `out_of_sync` lifecycle semantics and shared identity/signature comparison for rule-derived mappings.
- `datalink-api`: Add rule-scoped candidate query/recompute contracts and update readiness response semantics.

## Impact

- **Backend services**: `sourcerule`, `device`, connector diagnostics, candidate orchestration, and shared mapping lifecycle handling.
- **Persistence**: source-rule revision records, candidate snapshot storage, and identity/signature metadata.
- **APIs**: readiness response shape plus shared candidate query/recompute endpoints.
- **Tests**: source-rule revision persistence, snapshot restore, readiness gating, and candidate API coverage.
