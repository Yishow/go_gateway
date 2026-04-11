# Proposal: Workbench UX Operator Efficiency Experiment

## Why

Phase 1R Device compare is complete and the downstream direction is now approved:

1. `v2 / Sentry Incident Desk` is the primary functional direction.
2. `v1 / Linear Control Room` remains the full-flow high-polish alternate.
3. `v3 / ClickHouse Data Cockpit` remains a limited comparison/control track.

The current OpenSpec still models Phase 2–5 as three equal-cost parallel experiments. That is no longer accurate and now creates four problems:

- **spec drift**: approved direction and written tasks no longer match
- **delivery inefficiency**: `v2` has already won the functional role, but the spec still spreads equal effort
- **quality dilution**: `v1` risks becoming a half-maintained skin instead of a deliberate full-flow language
- **over-investment in `v3`**: the cockpit track should remain valuable, but not consume equal depth in every later phase

This correction keeps the original experiment constraints intact while changing the remaining rollout model to match the approved Phase 1R outcome.

## What Changes

The change now treats Phase 2–5 as a **winner-led rollout** instead of a three-equal-track race.

### Fixed version roles

| Version | New role | Meaning |
| --- | --- | --- |
| `baseline` | phase control snapshot | frozen `main` `/studio` surface used for comparison |
| `v2` | primary functional track | canonical workflow behavior for later phases |
| `v1` | full-flow high-polish track | same workflow as `v2`, expressed with a more refined operator language |
| `v3` | necessary-consistency comparison track | compare-capable cockpit surface with explicit minimum obligations |

### Fixed phase rhythm

Every remaining phase now follows:

1. shared contract / acceptance update
2. baseline check and evidence refresh
3. `v2`
4. `v1`
5. `v3`
6. compare gate

The compare gate remains mandatory before the next phase can begin.

### Baseline rule

For Phase 2–5, baseline means the `main` branch `/studio` surface:

- after the shared acceptance update for that phase is committed
- before any phase-specific variant UI work begins

### Compare authority rule

Compare may recommend:

- keeping `v2` as canonical and continuing
- harvesting ideas from `v1` or `v3`
- reopening the spec if the canonical direction is no longer acceptable

Compare may **not** silently transfer canonical ownership away from `v2`. Any owner change requires an explicit OpenSpec amendment.

## Shared Constraints

These constraints remain absolute:

1. same real API
2. same `/studio` route
3. same shared semantic token system
4. no mock-only flow
5. no version-specific backend contract
6. no additional product route

The remaining rollout also keeps one shared domain flow:

- Device -> Source -> Tag -> Output

Additional boundary rules:

- Phase 4 Output must cover both `Local Modbus register binding` and `Database schema/column binding`.
- Phase 5 shell may summarize readiness/blockers, but detailed editing and step-local validation remain owned by the underlying step surfaces.
- Shared acceptance tests must be defined on `main` before variant work begins.

## `v3` Minimum Obligations

`v3` no longer receives equal investment depth, but it still must:

1. run on the same `/studio` route and real API
2. complete the same phase entry and exit conditions as the shared contract
3. keep required blocker, error, and retry visibility
4. preserve compare evidence for overview, focused, and handoff states
5. keep any cockpit-specific additions additive, not workflow-replacing

## Compare Output Contract

Each remaining compare gate still must output:

- 操作順暢度
- 邏輯清晰度
- 對系統的完整性
- 首屏資訊密度
- 關鍵操作時間
- 實作 / 維護風險
- 推薦版本與理由

The compared entries remain:

- baseline
- `v2`
- `v1`
- `v3`

## Future Variant Expansion

If `v4+` is ever added later:

1. update proposal / design / tasks first
2. keep the same `/studio` route and shared API contract
3. keep the same shared semantic token system
4. allocate a new branch / worktree / port before implementation
5. extend the same compare matrix before the new variant can participate

## Non-Goals

This correction does **not**:

1. reopen the Phase 1R Device winner decision
2. turn `v1` into a skin-only layer
3. keep `v3` as an equal-cost delivery path
4. create new product routes or variant-only backend semantics
5. let compare silently replace the canonical owner

## Impact

### OpenSpec artifacts

- update `openspec/changes/workbench-ux-operator-efficiency/proposal.md`
- update `openspec/changes/workbench-ux-operator-efficiency/design.md`
- update `openspec/changes/workbench-ux-operator-efficiency/tasks.md`
- keep `specs/datalink-api/spec.md` unchanged as the shared backend contract

### Next execution step

After this correction, execution resumes with:

1. Phase 2 shared Source contract
2. Phase 2 baseline evidence
3. Phase 2 `v2`
4. Phase 2 `v1`
5. Phase 2 `v3`
6. Phase 2 compare
