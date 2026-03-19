## Context

The workbench currently delivers a coherent desktop shell, but its core operational model is still split across frontend-only planning state, runtime-facing points, manually maintained tag bindings, and output pages that share unstable state. This causes real operational drift:

- Step 1 cannot distinguish between `transport connect` and `protocol probe`, so network and protocol failures collapse into one ambiguous result.
- Step 2 rules are still planning-only constructs in frontend state, while runtime collection only understands persisted devices, points, mappings, and polling groups.
- Step 3 still behaves like a manual binding step even though the selected product direction is that rule creation should establish the primary Point→Tag relationship automatically.
- Step 4 and Database output surfaces still carry shared selection state and tightly coupled UI state that can drift away from the actual mapping target.

The selected direction keeps the existing domain layers (`device -> point -> tag -> output`) but changes their roles. `Rule / Source span` becomes the primary operator-created object, `Point` becomes more explicitly runtime/internal, and `Tag` remains the external semantic/output object. The workbench must orchestrate these layers automatically instead of forcing the operator to rebuild the relationship by hand in every step.

## Goals / Non-Goals

**Goals:**
- Separate Step 1 device testing into `connect` and `probe` phases with explicit diagnostics and gating behavior.
- Introduce persisted source-rule lifecycle behavior that supports create/update, enable/disable, restart restoration, and rule-driven live value display.
- Align Step 2 planning UI with device capability context and make `planned / used / conflict / unmanaged` states explicit.
- Move the primary Point→Tag creation path from manual binding to automatic rule-derived creation with strict `1 Point : 1 Tag`.
- Reframe Step 3 as a review/verification surface plus exception handling path.
- Isolate Step 4 output state per target and move binding state to one authoritative source.
- Rebuild the Database workbench flow around explicit `Connector / Schema / Mapping` layering for SQLite and PostgreSQL.

**Non-Goals:**
- Replacing the workbench shell or returning to the old SmartDashboard orchestration model.
- Removing the `Point` domain object entirely in this change.
- Expanding Database connector scope beyond SQLite and PostgreSQL in this round.
- Solving every advanced diagnostic/reporting feature in the first implementation slice (for example, full batch diagnostics or full mapping history UI).
- Reworking runtime scheduler internals beyond what is necessary to support persisted rule lifecycle and restart restoration.

## Decisions

### 1. Introduce a persisted `SourceRule` aggregate instead of keeping rules in frontend-only planning state

The system needs a first-class rule model in the backend and database because the selected behavior requires enable/disable semantics, runtime collection gating, restart restoration, and stable live-value display. Keeping rules only in `WorkbenchProvider` state cannot satisfy any of those behaviors.

**Chosen approach**
- Add a persisted rule model keyed by device and rule identity.
- Store enabled/disabled state, address model inputs, merge/type inputs, and metadata required to derive runtime-facing points and output-facing tags.
- Treat existing `SourceRule` frontend data as the seed for this persisted aggregate rather than the source of truth.

**Alternatives considered**
- Keep rules frontend-only and persist only points: rejected because restart restoration and lifecycle semantics would remain implicit and drift-prone.
- Make points themselves the only persisted source of truth: rejected because it loses the operator’s rule-level intent and makes enable/disable semantics too low-level.

### 2. Treat Step 1 device testing as a two-phase result (`connect`, `probe`) with partial-success gating

The user explicitly wants both results shown and judged separately. The current single success/failure shape cannot represent “host reachable enough to connect, but protocol probe not ready”.

**Chosen approach**
- Return and render two phase results:
  - `connect`: transport/socket-level reachability
  - `probe`: protocol-specific validation/readiness
- Permit saving a device when `connect = success` and `probe = failed`.
- Block rule activation and runtime collection until `probe = success`.

**Alternatives considered**
- Keep a single boolean result with better messages: rejected because the gating semantics still remain ambiguous.
- Require both phases before save: rejected because it blocks valid “save first, finish protocol later” field workflows.

### 3. Keep `Point` and `Tag` as separate domain objects, but automate the primary linkage path

The current domain still benefits from separating runtime acquisition (`Point`) from external semantic/output identity (`Tag`). However, the operator should not be forced to manually rebuild that relationship when the rule already describes the source.

**Chosen approach**
- Preserve `Point` as the runtime/internal acquisition asset.
- Preserve `Tag` as the output-facing semantic asset.
- Automatically create the primary Tag and Mapping when a rule becomes active in the main flow.
- Make Step 3 a review/verification and exception-handling surface rather than a mandatory manual first-pass binding step.

**Alternatives considered**
- Collapse to a single Rule/Tag model: rejected for this round because runtime internals and existing APIs still depend heavily on points.
- Keep manual Point→Tag binding as the primary path: rejected because it conflicts with the selected user workflow and causes repeated ceremony.

### 4. Use target-isolated output state plus a single authoritative binding model

Selection drift in Step 4 comes from shared cross-target state and partially duplicated local UI state. The solution is not more UI chrome; it is cleaner state ownership.

**Chosen approach**
- Maintain separate selection state per output target.
- Centralize `selected tag`, `selected register/column target`, and bind/unbind readiness in one authoritative state model.
- Let the UI reflect state instead of owning hidden competing state.

**Alternatives considered**
- Keep shared selection state and add more guards: rejected because it treats symptoms, not root cause.
- Split Local Modbus and Database into unrelated routes again: rejected because it breaks the chosen one-step Output workflow.

### 5. Model Database output as three layers: Connector, Schema/Table context, Mapping

The current Database page couples connector choice, schema fetch, table selection, validation, and binding actions so tightly that stale state and confusing actions are hard to avoid.

**Chosen approach**
- Treat connector selection as the outer scope.
- Treat schema/table context as a separate middle scope.
- Treat tag-to-column mappings as the innermost scope.
- Limit this change to SQLite and PostgreSQL so the layered model can be stabilized before expanding connector coverage.

**Alternatives considered**
- Keep a single coupled editor and just improve validation messaging: rejected because it does not solve stale state and action ambiguity.
- Fully generic connector abstraction first: rejected because the user selected a narrower connector scope for this round.

## Risks / Trade-offs

- **[Rule migration ambiguity]** Existing points may exist without persisted rules → **Mitigation:** explicitly classify them as unmanaged/legacy in Step 2 and avoid silently inferring rules during the first rollout.
- **[Auto-generated tag identity collisions]** Rule-name-plus-address generation can still collide in edge cases → **Mitigation:** keep case-insensitive uniqueness enforcement and surface deterministic collision diagnostics in Step 3 review.
- **[Runtime refactor scope]** Persisted rule lifecycle touches frontend, backend, schema, and bootstrap logic → **Mitigation:** stage implementation so persistence/model work lands before UI semantics depend on it.
- **[Output state refactor regressions]** Consolidating Local Modbus and Database state can break existing direct-binding flows → **Mitigation:** isolate target state by contract first, then refactor one target at a time behind regression tests.
- **[Operator expectation mismatch]** Users may assume disabling a rule deletes downstream state → **Mitigation:** specify clearly that disable stops collection only and preserves Point/Tag/Output relationships.

## Migration Plan

1. Add persisted rule storage and APIs without removing the current workbench flow.
2. Introduce runtime bootstrap for persisted rules and explicit enable/disable semantics.
3. Update Step 2 UI to read persisted rule state, device capability context, and unmanaged-point explanations.
4. Add automatic Tag + Mapping creation for the primary rule flow while preserving Step 3 as a review/exception surface.
5. Refactor Step 4 Local Modbus and Database state into target-isolated, authoritative models.
6. Treat legacy unmanaged points as visible but distinct until operators migrate them into persisted rules.
7. Roll forward by enabling the new rule-driven flow in workbench; roll back by disabling rule activation while keeping underlying points/tags intact.

## Open Questions

- Should disabling a rule freeze the last displayed live value, or clearly mark it as stale/inactive in the grid?
- How much manual editing should remain available for auto-created Point/Tag relationships before the system considers them “exception handling” instead of “normal flow”?
- Should unmanaged legacy points be migratable into rules inside the workbench, or remain read-only until a later change?
