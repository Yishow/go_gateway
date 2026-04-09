# Design: Workbench UX Operator Efficiency Experiment

## Goal

This correction keeps the approved Phase 1R winner decision, but reopens Phase 2–5 after a live review found that the post-Device surfaces still read too similarly outside Device.

The remaining rollout must keep one real `/studio` workflow over one real API while preserving three archetype-level versions under a winner-led delivery model.

## Approved Phase 1R Outcome

The accepted downstream decision is:

1. `v2 / Sentry Incident Desk` is the canonical functional direction.
2. `v1 / Linear Control Room` remains a full-flow high-polish alternate.
3. `v3 / ClickHouse Data Cockpit` remains a limited comparison/control track.

The correction therefore changes only the **remaining rollout model**, not the completed Device result.

## Reopen Trigger

The current Source/Tag previews exposed a governance gap:

1. the shared logic was mostly preserved
2. the branches did touch different Source/Tag files
3. but the operator-facing result still felt too similar outside Device

The issue is therefore not backend logic drift. The issue is that the remaining phases were allowed to fall back toward shared surfaces with lighter wrappers, instead of carrying the Device archetypes through the rest of the product.

Phase 2–5 are therefore reopened under the same shared logic and version roles, but with a stronger surface-differentiation rule.

## Hard Constraints

These rules remain absolute:

1. same real API
2. same `/studio` route
3. same shared semantic token system
4. no mock-only flow
5. no version-specific backend contract
6. no additional product route

Existing worktree identity remains unchanged:

| Surface | Branch / worktree identity | Role |
| --- | --- | --- |
| `baseline` | current frozen mainline snapshot | control entry for compare |
| `v1` | `woe-v1-radix` | high-polish alternate |
| `v2` | `woe-v2-mui` | canonical functional track |
| `v3` | `woe-v3-antd` | minimum-obligation comparison track |

## Logic Invariant

Across reopened Phase 2–5 work:

- the domain flow stays the same
- the shared API contract stays the same
- the state ownership / blocker / retry semantics stay the same

What changes is the operator surface:

- action placement and ordering
- primary working surface
- preview / summary scaffolding
- visual hierarchy and interaction pacing

If a review can still reasonably conclude that only Device looks different, the phase is not complete.

## Version Roles

### `baseline`

For Phase 2–5, baseline is the `main` branch `/studio` surface:

- after the shared acceptance update for that phase is committed
- before any phase-specific variant UI work starts

It exists only to provide a frozen control snapshot for compare.

### `v2`

`v2` is the primary functional track. Later workflow behavior is considered established only after it is correct in `v2`:

- main task flow
- blocker handling
- diagnostics surfaced state
- step handoff behavior
- shared acceptance contract

### `v1`

`v1` must implement the same accepted domain flow as `v2`, but can express it through:

- calmer hierarchy
- cleaner pacing
- better long-session readability
- more refined operator-facing detail

`v1` may not drop required surfaced state, steps, or blockers in the name of polish.

### `v3`

`v3` is the necessary-consistency comparison track.

It is complete for a phase only when:

1. the same phase entry and exit conditions can be exercised
2. the phase mainline action can complete through the same handoff
3. blocker, error, and retry visibility remain present
4. compare evidence exists for overview, focused, and handoff states
5. any extra cockpit surface is additive and does not replace the shared workflow

`v3` may omit secondary refinements that are not required to satisfy those conditions.

## Archetype Continuity Rule

The Device 1R identities must remain visible in every later phase:

- `v1 / Linear Control Room`: calm task zoning, deliberate pacing, long-session readability, and visibly different work/readout framing
- `v2 / Sentry Incident Desk`: command-deck structure, surfaced blocker/retry narrative, and task-centered incident-style operation
- `v3 / ClickHouse Data Cockpit`: dense banner / KPI / board composition, cockpit-like summary density, and clearly different primary readout surfaces

The following do **not** satisfy this rule on their own:

- kit-only substitution
- CSS-only reskinning
- identical action choreography with slightly different component chrome

### Canonical ownership rule

Compare may recommend harvesting ideas from `v1` or `v3`, or stopping for a spec correction, but it may not silently move canonical ownership away from `v2`. Any owner change requires an explicit OpenSpec amendment.

## Fixed Execution Rhythm

Every remaining phase follows this sequence:

1. shared contract / acceptance update
2. baseline check and evidence refresh
3. `v2`
4. `v1`
5. `v3`
6. compare gate

The next phase does not begin until the compare gate is complete.

## Reopen Rule

Previously accepted Phase 2 work and partial Phase 3 work are downgraded to reference evidence only. They may inform the rerun, but they do not count as completed acceptance under this amendment.

### Baseline check checklist

Each baseline check must confirm:

1. `/studio` loads on the real API without route or contract drift
2. the phase entry state is reachable from the current mainline flow
3. blocker and handoff states still render from the shared contract
4. baseline evidence is captured before any variant branch begins phase-specific UI work

## Shared Domain Contract

### Shared flow

The flow remains:

- Device -> Source -> Tag -> Output

### Shared phase state expectations

Each phase must expose the same core state shape across baseline, `v2`, `v1`, and `v3`:

- focused working context
- actionable draft/list/candidate state
- loading / empty / blocked / failed / retryable state
- next-step handoff state

Variants may change presentation and grouping, but not state ownership or domain meaning.

### Source -> Tag boundary

Source rules and derived candidate state determine:

- tag review candidates
- diff preview scope
- batch action eligibility

`diff preview scope` is fixed as the delta set produced by the currently selected device context plus the currently selected source rule revision. It must:

- include create / update / unbind outcomes for affected candidates
- exclude unrelated devices or manual-only tags
- invalidate stale previews when the rule, live scope, or candidate snapshot changes

### Tag -> Output boundary

Tag binding/apply state determines:

- output candidate readiness
- blocker explanation
- dry-run/apply eligibility

Phase 4 Output must always cover both target families:

- Local Modbus register binding
- Database schema/column binding

### Cross-step shell ownership

The shell owns only four cross-step surfaces:

- readiness summary
- active blocker summary
- diagnostics refresh status
- shortest return-to-mainline action

Detailed editing, mutation, and step-local validation remain owned by the corresponding phase surfaces.

## Testing and Evidence

### Shared tests first

The shared contract is owned on `main` before variant work starts. It must be branch-neutral and define, at minimum:

1. mainline happy path
2. selected / focused actionable state
3. empty or not-yet-ready state
4. loading or in-progress state
5. blocker or validation-failure state
6. request failure plus retry path
7. successful recovery and handoff to the next step

Variant branches may add UI-specific tests, but they may not redefine shared domain expectations.

### Minimum per-phase scenario focus

#### Phase 2 — Source

- create or edit a rule and hand off to Tag
- template apply and plan / live / link switching
- invalid rule or missing device context blocker
- stale preview invalidation, failure, retry, and recovery

#### Phase 3 — Tag

- review candidates, inspect diff, apply decision, and hand off to Output
- unresolved source issues or no candidates blocker
- apply failure, refresh failure, retry, and recovery

#### Phase 4 — Output

- inspect readiness and blockers for both Local Modbus and Database targets
- run dry-run when required
- apply output actions and recover from blocked or failed runs

#### Phase 5 — Shell / diagnostics

- show readiness summary and correct next action
- surface global blocker and shortest valid return path
- refresh diagnostics, handle failure, retry, and recovery

#### Phase 6 — Final compare

- re-run the end-to-end path across baseline, `v2`, `v1`, and `v3`
- produce the final recommendation and retained ideas

### Browser evidence minimum

Each version must keep at least:

- overview state
- selected / focused state
- next-step handoff state

`v2` evidence must prove the full mainline flow.

`v1` evidence must prove the same flow remains intact while the language becomes more refined and consistent.

`v3` evidence must prove cockpit strengths remain present without breaking the shared flow.

## Phase Model

### Phase 2 — Source

`v2` defines the canonical rule create/edit, template apply, mode switching, inspector, and stale-preview recovery flow through an incident-desk command surface. `v1` keeps the same flow, but must express it through a clearly different control-room editing skeleton with different action zoning, preview staging, and handoff scaffolding. `v3` keeps cockpit density through clearly different banner / KPI / board composition while preserving the same Source -> Tag handoff.

### Phase 3 — Tag

`v2` defines the canonical review queue, diff preview, batch decision, apply, retry, and recovery flow through an incident-desk review command deck. `v1` keeps those same actions while moving them into a clearly different control-room review skeleton optimized for long batch sessions. `v3` keeps board/summary density only if the primary work surface still reads as a cockpit and the same apply / recovery contract remains visible.

### Phase 4 — Output

`v2` defines the canonical flow for both Local Modbus and Database targets, and both must cover readiness, dry-run, apply, and blocker diagnosis through an incident-desk output command surface. `v1` keeps the same targets and actions through a clearly different control-room operator console. `v3` may keep denser mapping/state views, but it must still read as a cockpit surface while surfacing both targets and the same apply/blocker contract.

### Phase 5 — Cross-step shell / diagnostics

`v2` is the canonical shell implementation. `v1` re-expresses the same shell abilities with a clearly different control-room summary skeleton and lower interference. `v3` keeps only the cockpit summary/alert surfaces needed to satisfy the same shell ownership model, but those surfaces must remain recognizably cockpit-oriented rather than a light shell reskin.

### Phase 6 — Final compare

Phase 6 re-runs the end-to-end validation across baseline, `v2`, `v1`, and `v3`, then produces the final recommendation and retained follow-up ideas.

## Compare Gate

Each remaining phase must output:

- 操作順暢度
- 邏輯清晰度
- 對系統的完整性
- 首屏資訊密度
- 關鍵操作時間
- 實作 / 維護風險
- 推薦版本與理由
- 三版本在非 Device phase 是否仍然明確可辨

The compared entries remain:

- baseline
- `v2`
- `v1`
- `v3`

## Stop Conditions

Pause and amend artifacts before continuing if:

1. a version needs route-specific or backend-specific semantics
2. baseline cannot be frozen for the current phase
3. shared token semantics are no longer sufficient for all compared surfaces
4. compare output is not sufficient to support a recommendation
5. a new variant is needed before proposal / design / tasks are extended

## Implementation Boundary

This correction updates only:

- `proposal.md`
- `design.md`
- `tasks.md`

It does not change the shared backend API spec. After the correction is committed, execution resumes with Phase 2 shared acceptance work.
