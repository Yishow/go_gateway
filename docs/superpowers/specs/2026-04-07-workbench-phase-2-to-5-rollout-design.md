# Workbench Phase 2–5 Winner-led Rollout Design

- Status: Brainstorming-approved; pending spec-review
- Date: 2026-04-07
- Scope: OpenSpec correction and delivery model for `/studio` Phase 2–5 after Phase 1R Device compare

## 1. Summary

Phase 1R established a new Device round with three visibly distinct archetypes and finished with a compare outcome:

1. `v2 / Sentry Incident Desk` is the best-balanced version for workflow correctness, diagnostics prominence, and operator efficiency.
2. `v1 / Linear Control Room` is the strongest low-risk fallback and the best candidate for a higher-polish full-flow treatment.
3. `v3 / ClickHouse Data Cockpit` contains valuable high-density cockpit ideas, but its maintenance cost is too high to keep as an equal-investment track for every later phase.

The current OpenSpec still assumes Phase 2–5 are three-equal variant races. That no longer matches the approved direction.

This design changes the rollout model for the remaining workbench phases:

- `v2` becomes the **primary functional delivery track**
- `v1` becomes the **full-flow high-polish parallel track**
- `v3` becomes the **necessary-consistency comparison track**

The compare gate remains mandatory at the end of every phase, but the investment model changes from “three equal implementation efforts” to “one canonical workflow track, one full-quality alternate expression, one limited comparison/control track.”

## 2. Problem Statement

The approved user direction is now:

- accept the Phase 1R compare result
- continue through the entire flow, not only Device
- carry `v1` forward as a full-flow high-quality language, not a cosmetic skin
- stop treating `v3` as an equal-cost delivery path unless its cockpit-specific strengths justify the extra effort

If Phase 2–5 continue under the previous three-equal-track assumption, the project will pay three costs at once:

1. **Spec drift cost**
   - The operational decision has changed, but the OpenSpec task matrix still encodes the old shape.

2. **Delivery inefficiency**
   - All three versions would continue receiving equal effort even though the approved compare result already selected `v2` as the primary direction.

3. **Quality dilution**
   - `v1` would risk becoming a half-maintained alternate branch instead of the deliberate, refined end-to-end language the user asked for.

The remaining work must therefore be re-scoped before implementation continues.

## 3. Goals

### Product goals

- Keep `/studio` as one real workflow over the same real API and same route contract.
- Carry the winning `v2` interaction model through `Source`, `Tag`, `Output`, and cross-step diagnostics.
- Carry `v1` through the same end-to-end workflow with a cleaner, more refined, more consistent product language.
- Retain `v3` only as a comparison-capable cockpit track with explicit limits.

### Architecture goals

- Keep one shared domain/data-flow contract across all remaining phases.
- Prevent branch-local logic forks, hidden behavior divergence, or version-only backend semantics.
- Ensure each phase still ends with a comparable evidence-based decision.

### Delivery goals

- Re-center development effort on the chosen winner without losing the value of the other two tracks.
- Preserve the compare discipline while avoiding unnecessary equal-cost implementation.

## 4. Non-Goals

- Reopening the Phase 1R Device winner selection.
- Reintroducing three equal-investment tracks for every remaining phase.
- Allowing `v1` to diverge from the canonical domain flow in the name of polish.
- Allowing `v3` to become a second primary product direction.
- Creating new routes, mock backends, or variant-specific API contracts.

## 5. Approved Decisions

1. The Phase 1R compare result is accepted for downstream planning.
2. The remaining work is a multi-phase rollout across `Source`, `Tag`, `Output`, and cross-step shell/diagnostics.
3. `v2` is the primary implementation authority for later workflow behavior.
4. `v1` must remain a full-flow implementation, not a skin-only layer.
5. `v3` is retained for necessary consistency and compare value only.
6. Every later phase still requires a compare gate before moving forward.

## 6. Approaches Considered

### Approach A — Direct winner-led rollout

Make `v2` the only serious track for later phases, keep `v1` mostly visual, and reduce `v3` to minimal parity.

- **Pros:** fastest
- **Cons:** turns `v1` into a skin track and conflicts with the approved requirement that `v1` stay a full-flow polished implementation

### Approach B — Shared-contract-first rollout

First update the spec, then implement the remaining phases with explicit role separation:

- `v2` = primary functional track
- `v1` = full-flow refined alternate track
- `v3` = necessary-consistency comparison track

- **Pros:** matches the approved user direction, minimizes spec drift, preserves compare discipline
- **Cons:** requires spec work before implementation resumes

### Approach C — v1 aesthetic-first rollout

Let `v1` become the lead track for the remaining phases, with `v2` following later.

- **Pros:** maximizes the polish-first goal
- **Cons:** contradicts the accepted compare outcome and raises delivery risk

The chosen approach is **Approach B — Shared-contract-first rollout**.

## 7. Execution Architecture

### 7.1 Hard constraints that remain unchanged

These rules stay absolute for all remaining phases:

- same real API
- same `/studio` route
- same shared semantic token system
- no mock data
- no version-specific backend contracts
- no additional product routes

### 7.2 Version roles

#### `v2 / Sentry Incident Desk`

`v2` is the **primary functional delivery track**.

For each later phase, the new workflow behavior is considered “established” only after it is correct in `v2`:

- primary task flow
- blocker handling
- diagnostics surfaced state
- step handoff behavior
- acceptance test contract

#### `v1 / Linear Control Room`

`v1` is the **full-flow high-polish parallel track**.

It must implement the same accepted domain flow as `v2`, but it is allowed to express that flow differently through:

- calmer hierarchy
- cleaner pacing
- better long-session readability
- more refined operator-facing details

`v1` is explicitly **not** allowed to drop steps, reduce surfaced state, or hide blockers in order to look cleaner.

#### `v3 / ClickHouse Data Cockpit`

`v3` is the **necessary-consistency comparison track**.

Its job is to:

- remain runnable on the same route and API
- execute the compare scenarios for each phase
- preserve any cockpit-specific strengths worth reusing later

It is not required to receive equal feature-investment depth as `v1` and `v2`.

### 7.3 Fixed phase rhythm

For every remaining phase, the sequence is:

1. shared contract / acceptance update
2. baseline check and evidence refresh
3. `v2` full implementation and validation
4. `v1` full-flow polished implementation and validation
5. `v3` necessary consistency implementation and validation
6. compare gate

The next phase does not start until the compare gate completes.

## 8. Shared Data-Flow Contract

The data-flow contract stays shared across all versions.

### 8.1 Device → Source

Selected `Device` context determines:

- source planner defaults
- plan/live/link scope
- runtime readiness and diagnostics context

### 8.2 Source → Tag

Source rules and derived point/candidate state determine:

- tag review candidates
- diff preview scope
- batch action eligibility

### 8.3 Tag → Output

Tag binding/apply state determines:

- output candidate readiness
- blocker explanation
- dry-run/apply eligibility

### 8.4 Allowed implementation surface

Versions may only change:

- presentation
- grouping
- operator guidance hierarchy
- interaction arrangement

Versions may not introduce:

- branch-local data transforms
- version-only query logic
- alternate state ownership
- route-specific workflow semantics

## 9. Phase-by-Phase Surface Model

### 9.1 Phase 2 — Source

#### `v2`

`v2` defines the canonical operator flow for:

- rule creation/editing
- template application
- plan/live/link mode switching
- rule layer and inspector behavior

#### `v1`

`v1` keeps the exact same rule workflow but re-expresses it as a cleaner, lower-noise control surface optimized for prolonged editing sessions.

#### `v3`

`v3` preserves the high-density canvas and summary strengths needed for compare, without requiring equal expansion work.

### 9.2 Phase 3 — Tag

#### `v2`

`v2` defines the canonical review queue, batch decision, diff preview, apply, and failure-recovery flow.

#### `v1`

`v1` keeps the same workflow and actions, but emphasizes composure, legibility, and refined feedback during large review batches.

#### `v3`

`v3` keeps the board and summary density required for compare, but only within the shared flow contract.

### 9.3 Phase 4 — Output

#### `v2`

`v2` defines the canonical output flow for:

- state overview
- dry-run
- apply
- blocker diagnosis

#### `v1`

`v1` keeps the same flow, but presents blockers and state transitions in a calmer operator-console language instead of a dense dashboard idiom.

#### `v3`

`v3` preserves dense mapping/state views only where they add compare value.

### 9.4 Phase 5 — Cross-step shell / diagnostics

#### `v2`

`v2` becomes the canonical global shell for:

- readiness surfaced state
- global blockers
- shortest return-to-mainline actions

#### `v1`

`v1` expresses the same shell abilities with lower interference and stronger visual consistency across the entire flow.

#### `v3`

`v3` keeps only the most valuable cockpit summary/alert ideas needed for compare and future extraction.

## 10. Error Handling Rules

### 10.1 Canonical behavior source

`v2` is the behavior source for later phases.

This means:

- pending states
- error surfaced states
- blocker copy expectations
- retry paths
- empty/loading states

must be made correct in `v2` first.

### 10.2 `v1` error-handling obligations

`v1` must fully preserve the same error/blocker scenarios as `v2`.

It may only change:

- emphasis
- hierarchy
- pacing
- visual framing

It may not:

- hide required blockers
- swallow failures
- reduce surfaced state
- change what is considered actionable

### 10.3 `v3` error-handling obligations

`v3` must preserve compare-required error and blocker visibility, but does not need a second full innovation track.

## 11. Testing and Evidence Rules

### 11.1 Shared tests first

Every phase begins by updating or adding shared acceptance tests that define the phase contract.

### 11.2 Delivery order

1. `v2` goes red → green first
2. `v1` aligns to the same shared tests, plus any needed v1-specific UI contract tests
3. `v3` implements only the compare-required shared surface

No compare result is valid if a branch depends on branch-local expectation drift.

### 11.3 Browser evidence minimum

Each version must keep at least:

- overview state
- selected/focused state
- next-step handoff state

`v2` evidence must prove the full mainline flow.

`v1` evidence must prove the same flow remains intact while the language becomes more refined and consistent.

`v3` evidence must prove cockpit strengths remain present without breaking the shared flow.

## 12. Compare Gate Output Contract

Each remaining phase must output:

- 操作順暢度
- 邏輯清晰度
- 對系統的完整性
- 首屏資訊密度
- 關鍵操作時間
- 實作 / 維護風險
- 推薦版本與理由

The compare is complete only when those items are written against:

- baseline
- `v2`
- `v1`
- `v3`

using the role model defined in this design.

## 13. OpenSpec Amendment Plan

Before implementation resumes, the change set must update:

1. `openspec/changes/workbench-ux-operator-efficiency/proposal.md`
2. `openspec/changes/workbench-ux-operator-efficiency/design.md`
3. `openspec/changes/workbench-ux-operator-efficiency/tasks.md`

The amendment must explicitly redefine Phase 2–5 roles as:

- `v2` primary functional track
- `v1` full-flow high-polish track
- `v3` necessary-consistency comparison track

## 14. Implementation Entry Criteria

Implementation may start only after:

1. this spec passes spec review
2. the user reviews the written spec
3. the OpenSpec correction plan is accepted
4. a concrete implementation plan is written for Phase 2

That implementation starts with:

- Phase 2 shared acceptance update
- baseline check
- `v2`
- `v1`
- `v3`
- compare

## 15. Risks and Mitigations

### Risk: `v1` becomes cosmetic-only

**Mitigation:** require `v1` to pass the same shared flow and error-handling tests as `v2`.

### Risk: `v3` still consumes too much effort

**Mitigation:** explicitly limit `v3` to necessary consistency plus compare value.

### Risk: compare language becomes vague

**Mitigation:** keep one fixed compare output template for every remaining phase.

### Risk: new spec role split is not reflected in implementation behavior

**Mitigation:** amend OpenSpec before any Phase 2 code work resumes.
