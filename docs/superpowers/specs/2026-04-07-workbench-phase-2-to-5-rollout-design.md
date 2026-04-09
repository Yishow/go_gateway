# Workbench Phase 2–5 Archetype-preserving Winner-led Rollout Design

- Status: Brainstorming-approved; pending spec-review
- Date: 2026-04-07
- Scope: OpenSpec correction and delivery model for `/studio` Phase 2–5 after Phase 1R Device compare

## 1. Summary

Phase 1R established a new Device round with three visibly distinct archetypes and finished with a compare outcome:

1. `v2 / Sentry Incident Desk` is the best-balanced version for workflow correctness, diagnostics prominence, and operator efficiency.
2. `v1 / Linear Control Room` is the strongest low-risk fallback and the best candidate for a higher-polish full-flow treatment.
3. `v3 / ClickHouse Data Cockpit` contains valuable high-density cockpit ideas, but its maintenance cost is too high to keep as an equal-investment track for every later phase.

The current OpenSpec no longer has the old three-equal-track problem, but a new live-review problem emerged: outside Device, the current Source/Tag work still reads too similarly across versions.

This amendment keeps the winner-led owner model, but reopens Phase 2–5 so the remaining workbench phases preserve three genuinely different archetype-level surfaces:

- `v2` remains the **primary functional delivery track**
- `v1` remains the **full-flow high-polish parallel track**
- `v3` remains the **necessary-consistency comparison track**

The compare gate remains mandatory at the end of every phase, but completion now also requires the resulting phase to still read as three recognizably different versions beyond Device.

## 2. Problem Statement

The approved user direction is now:

- accept the Phase 1R compare result
- continue through the entire flow, not only Device
- carry `v1` forward as a full-flow high-quality language, not a cosmetic skin
- stop treating `v3` as an equal-cost delivery path unless its cockpit-specific strengths justify the extra effort

If Phase 2–5 continue from the current partial state without reopening the surface requirement, the project pays a new cost:

1. **Archetype drift cost**
   - Different branches touch different files, but the operator-facing result still collapses back toward one shared product.

2. **Review failure cost**
   - The user already reviewed the live previews and rejected the current outcome because only Device reads clearly different.

3. **Mixed-suite cost**
   - If only later phases are corrected, `/studio` will still contain a half-distinct product line where Source remains too close while later steps diverge.

The remaining work must therefore be re-scoped before implementation continues.

## 3. Goals

### Product goals

- Keep `/studio` as one real workflow over the same real API and same route contract.
- Keep the shared domain logic unchanged while reworking all operator-facing operations.
- Carry `v2` through `Source`, `Tag`, `Output`, and cross-step diagnostics as a clearly recognizable incident-desk product.
- Carry `v1` through the same end-to-end workflow as a clearly recognizable control-room product.
- Retain `v3` as a clearly recognizable cockpit product with explicit limits.

### Architecture goals

- Keep one shared domain/data-flow contract across all remaining phases.
- Prevent branch-local logic forks, hidden behavior divergence, or version-only backend semantics.
- Ensure each phase still ends with a comparable evidence-based decision.
- Ensure each reopened phase stays recognizably different beyond Device.

### Delivery goals

- Re-center development effort on the chosen winner without losing the value of the other two tracks.
- Preserve the compare discipline while avoiding unnecessary equal-cost implementation.

## 4. Non-Goals

- Reopening the Phase 1R Device winner selection.
- Reintroducing three equal-investment tracks for every remaining phase.
- Allowing `v1` to diverge from the canonical domain flow in the name of polish.
- Allowing `v3` to become a second primary product direction.
- Creating new routes, mock backends, or variant-specific API contracts.
- Accepting wrapper-only or skin-only variation as sufficient completion.

## 5. Approved Decisions

1. The Phase 1R compare result is accepted for downstream planning.
2. The remaining work is a multi-phase rollout across `Source`, `Tag`, `Output`, and cross-step shell/diagnostics.
3. `v2` is the primary implementation authority for later workflow behavior.
4. `v1` must remain a full-flow implementation, not a skin-only layer.
5. `v3` is retained for necessary consistency and compare value only.
6. Every later phase still requires a compare gate before moving forward.
7. Phase 2–5 are reopened because current non-Device surfaces are not distinct enough.

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

Each compare gate may still recommend harvesting ideas from `v1` or `v3`, or halting for a spec correction, but it does **not** silently transfer canonical ownership away from `v2`. Changing the primary track requires an explicit OpenSpec amendment.

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

`v3` is considered complete for a phase only when all of the following are true:

1. the same phase entry and exit conditions as the shared contract can be exercised on the same `/studio` route
2. the phase's mainline action can be completed through handoff to the next step
3. required blocker, error, and retry states remain visible to the operator
4. compare evidence exists for overview, focused, and next-step handoff states
5. any extra cockpit-only surface is additive and does not replace the shared workflow contract

`v3` may omit secondary refinements that are not needed to satisfy those five conditions.

### 7.2A Archetype continuity rule

Every reopened phase must preserve the Device 1R identities:

- `v1`: control-room zoning, calmer pacing, long-session readability, clearly different work/readout framing
- `v2`: incident-desk command surfaces, blocker-first narrative, task-centered operation
- `v3`: cockpit density, KPI/banner/board composition, clearly different primary readout surfaces

The following are insufficient on their own:

- kit-only substitution
- CSS-only overlays
- identical action choreography with different component chrome

### 7.3 Fixed phase rhythm

For every remaining phase, the sequence is:

1. shared contract / acceptance update
2. baseline check and evidence refresh
3. `v2` full implementation and validation
4. `v1` full-flow polished implementation and validation
5. `v3` necessary consistency implementation and validation
6. compare gate

The next phase does not start until the compare gate completes.

### 7.3A Reopen rule

Completed Phase 2 work and partial Phase 3 work become reference evidence only until rerun under this amendment.

### 7.4 Baseline definition

For Phase 2–5, **baseline** means:

- the `main` branch `/studio` surface at the start of the phase
- after the shared contract / acceptance update for that phase is committed
- before any phase-specific `v2`, `v1`, or `v3` UI implementation starts

Baseline is therefore the frozen control snapshot for that phase, not:

- the last recommended variant
- the production deployment at some unrelated point in time
- an evolving branch tip after variant work has started

Each phase must refresh baseline evidence from that exact starting point and use it as the control entry in the compare gate.

The baseline check for each phase must confirm only four things:

1. `/studio` loads on the real API without route or contract drift
2. the phase entry state is reachable from the current mainline flow
3. blocker and handoff states still render from the shared contract
4. baseline evidence is captured before any variant branch starts phase-specific UI work

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

`diff preview scope` is fixed as the delta set produced by the currently selected device context plus the currently selected source rule revision. It must include create/update/unbind outcomes for the affected candidates, exclude unrelated devices or manual-only tags, and invalidate stale previews whenever the rule, live scope, or candidate snapshot changes.

### 8.3 Tag → Output

Tag binding/apply state determines:

- output candidate readiness
- blocker explanation
- dry-run/apply eligibility

Phase 4 Output always covers two required target families under that same contract:

- Local Modbus register binding
- Database schema/column binding

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

`v2` defines the canonical flow for rule create/edit, template apply, plan/live/link switching, inspector behavior, and stale-preview recovery through an incident-desk command surface. `v1` keeps the same flow through a clearly different control-room editing skeleton. `v3` keeps the same handoff through a clearly different cockpit banner / KPI / board composition.

### 9.2 Phase 3 — Tag

`v2` defines the canonical review queue, diff preview, batch decision, apply, retry, and recovery flow through an incident-desk review command surface. `v1` keeps those same actions through a clearly different control-room review skeleton. `v3` keeps board/summary density only where the same apply and recovery contract stays visible and the phase still reads as a cockpit.

### 9.3 Phase 4 — Output

`v2` defines the canonical output flow for both Local Modbus and Database targets, and both must cover readiness, dry-run, apply, and blocker diagnosis through an incident-desk output surface. `v1` keeps those same target families and actions with a clearly different control-room console. `v3` may keep denser mapping/state views, but it must still surface both target families through a clearly different cockpit surface.

### 9.4 Phase 5 — Cross-step shell / diagnostics

The shared shell contract owns only four cross-step surfaces: readiness summary, active blocker summary, diagnostics refresh status, and shortest return-to-mainline action. Detailed editing, mutation, and step-local validation stay owned by the underlying phase surfaces. `v2` is the canonical shell implementation, `v1` re-expresses the same shell through a clearly different control-room summary skeleton, and `v3` keeps only the cockpit summary/alert surfaces needed to satisfy that same ownership model while remaining recognizably cockpit-like.

## 10. Error Handling Rules

### 10.1 Canonical behavior source

`v2` is the behavior source for later phases.

Pending, error, blocker copy, retry, stale/recovery, and empty/loading states must therefore be made correct in `v2` first.

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

The shared phase contract is owned on `main` by the primary agent before variant work starts.

That contract must be branch-neutral and must define, at minimum:

1. the mainline happy path for the phase
2. the selected/focused state needed to act
3. empty or not-yet-ready state
4. loading or in-progress state
5. blocker or validation-failure state
6. request failure plus retry path
7. successful recovery and handoff to the next step

Variant branches may add UI-specific tests, but they may not redefine those shared domain expectations.

### 11.1.1 Minimum per-phase scenario matrix

Every Phase 2–5 shared acceptance update must encode the following minimum scenario set.

#### Phase 2 — Source

- mainline: create or edit a source rule and hand off to Tag
- blocker: no valid device context or invalid rule configuration
- loading: source planner or live-state fetch in progress
- failure/retry: rule save, load, or diagnostics fetch fails and can be retried
- recovery: corrected rule returns the step to ready state

#### Phase 3 — Tag

- mainline: review candidates, inspect diff, apply decision, and hand off to Output
- blocker: no eligible source-driven candidates or unresolved rule issues
- loading: candidate/review state fetch in progress
- failure/retry: apply or refresh fails and can be retried
- recovery: successful re-run returns the queue to actionable state

#### Phase 4 — Output

- mainline: inspect readiness, run dry-run if required, apply binding/output action
- blocker: no eligible tags/bindings or unresolved upstream dependency
- loading: target/readiness state fetch in progress
- failure/retry: dry-run or apply fails and can be retried
- recovery: corrected state returns the output step to ready/applyable state

#### Phase 5 — Cross-step shell / diagnostics

- mainline: global shell shows readiness and returns the operator to the correct next action
- blocker: global blocker is surfaced with the shortest valid return-to-mainline action
- loading: global readiness or diagnostics refresh is in progress
- failure/retry: shell-level refresh/diagnostics action fails and can be retried
- recovery: refreshed shell returns to valid readiness state

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

### 12.1 Required compared surfaces

Each remaining phase must output:

- 操作順暢度 (flow efficiency)
- 邏輯清晰度 (workflow clarity)
- 對系統的完整性 (system completeness)
- 首屏資訊密度 (first-screen information density)
- 關鍵操作時間 (critical task time)
- 實作 / 維護風險 (implementation and maintenance risk)
- 推薦版本與理由 (recommended version and rationale)

The compare is complete only when those items are written against:

- baseline
- `v2`
- `v1`
- `v3`

using the role model defined in this design.

The compare is also incomplete unless each entry includes evidence for:

- baseline as defined in section 7.4
- `v2` as the canonical behavior implementation
- `v1` as the full-flow polished alternate
- `v3` as the minimum-obligation comparison cockpit

## 13. OpenSpec Amendment Plan

Before implementation resumes, the change set must update:

1. `openspec/changes/workbench-ux-operator-efficiency/proposal.md`
2. `openspec/changes/workbench-ux-operator-efficiency/design.md`
3. `openspec/changes/workbench-ux-operator-efficiency/tasks.md`

The amendment must explicitly redefine Phase 2–5 roles as:

- `v2` primary functional track
- `v1` full-flow high-polish track
- `v3` necessary-consistency comparison track

### 13.1 Required amendment outputs

The OpenSpec correction is not complete until those files also contain:

1. the fixed phase rhythm: shared -> baseline -> `v2` -> `v1` -> `v3` -> compare
2. the baseline definition from section 7.4
3. the `v3` minimum-obligation completion rule from section 7.2
4. the shared phase-contract ownership rule from section 11.1
5. the minimum per-phase scenario matrix from section 11.1.1
6. the fixed compare output contract from section 12

### 13.2 Planning handoff criteria

Phase 2 planning may begin only when all of the following are true:

1. `proposal.md` states the winner-led role split and unchanged hard constraints
2. `design.md` encodes the baseline definition, `v3` minimum obligations, and shared scenario matrix
3. `tasks.md` breaks each remaining phase into shared, baseline, `v2`, `v1`, `v3`, and compare tasks with explicit done conditions
4. those OpenSpec amendments are committed and accepted as the active source of truth

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
