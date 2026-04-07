## Context

The repository already contains a partially rule-aware datalink system:

- persisted devices
- persisted source rules
- runtime points derived from rule or manual flows
- tag and mapping services
- database target mappings
- Local Modbus output behavior

The problem is not a lack of domain foundations. The problem is that the product workflow, API surface, and output targets still behave as if operators should manually assemble `point -> tag -> output` relationships, while the backend increasingly behaves as if `SourceRule` should be the primary orchestration object.

This umbrella change exists to stop that split-brain evolution. It defines the final target-state contract for a rule-centric `/studio` workflow and gives the five executable child changes one shared source of truth.

The umbrella scope intentionally spans:

- product workflow
- workbench UI shell and step semantics
- source-rule runtime orchestration
- tag ownership and review/apply behavior
- database target review/apply behavior
- Local Modbus review/apply behavior
- device and connector readiness semantics
- API contracts that support candidate/recompute/apply flows

## Goals / Non-Goals

**Goals:**

- Define the final target-state workflow as `Device -> SourceRule -> Tag review -> Output review/apply`.
- Establish `Device` and `SourceRule` as the primary workflow objects.
- Establish a candidate/review/apply model that prevents silent downstream drift after rule edits.
- Make Database and Local Modbus equal first-class output targets.
- Provide a phase map that lets the five child changes be implemented in sequence without re-opening core architecture debates.
- Provide one explicit spec map across the affected capability surfaces.

**Non-Goals:**

- Deliver all implementation details inside this umbrella change.
- Introduce a new top-level `Pipeline` aggregate above `SourceRule`.
- Preserve manual point/tag/mapping-first workflow assumptions in the primary `/studio` flow.
- Expand database connector scope beyond SQLite and PostgreSQL in the primary output flow.
- Solve every low-level payload, storage, or screen-level detail here; child changes will own those details within this contract.

## Decisions

### 1. Use one umbrella change plus five executable child changes

The approved structure is:

1. `rule-centric-studio-refactor` (umbrella)
2. `rule-driven-foundation`
3. `tag-review-first-flow`
4. `database-output-review-apply`
5. `local-modbus-review-apply`
6. `studio-orchestration-cleanup`

**Why this over one mega change**

- One mega change would make tasks, review, and implementation boundaries too large.
- The child changes need clear execution order and independent validation.
- The user explicitly wants both a large refactor mother change and implementable sequential child changes.

**Why this over child changes only**

- Child changes without a mother change would duplicate architecture arguments and phase sequencing.
- Shared rules like ownership metadata, revision semantics, and target independence need one parent contract.

### 2. Make `SourceRule` the center of the operator workflow without inventing a new top-level aggregate

`SourceRule` becomes the primary acquisition object after device selection. `Point` remains necessary, but it becomes runtime/internal rather than the thing operators assemble manually.

**Why this over a new pipeline aggregate**

- The current codebase already knows how to derive runtime points and mappings from source rules.
- A new top-level aggregate would multiply concepts and delay delivery.
- The approved direction is rule-centric, not pipeline-centric.

### 3. Introduce a candidate/review/apply layer as the missing orchestration boundary

The final workflow requires three distinct states:

- persisted workflow inputs (`Device`, `SourceRule`)
- generated candidates (`TagCandidate`, `DatabaseOutputCandidate`, `LocalModbusOutputCandidate`)
- applied downstream assets (`Tag`, mappings, output bindings)

This layer is the key architectural addition because it lets rule revisions regenerate downstream intent without silently overwriting applied state.

**Why this over immediate cascade**

- Immediate downstream persistence on rule save makes review-first UX impossible.
- The user approved diff-driven regeneration with explicit apply decisions.
- Both output targets need the same orchestration principle.

### 4. Treat Database connector context as shared environment state, not a rule-owned primary input

The operator configures database connectors, but connectors are not part of the `Device -> SourceRule` chain. They are shared environment configuration that database candidates reference.

**Why this matters**

- It keeps the workflow model clean: rules define data intent, connectors define destination context.
- It explains why connector invalidation blocks database apply without forcing the whole workflow to restart.

### 5. Make Local Modbus output persistent and rule-owned

The final target state cannot rely on in-memory-only Local Modbus mapping state. Review/apply, out-of-sync, and sequential implementation all require persistent ownership and revision tracking.

**Why this matters**

- The user wants Local Modbus to be a first-class output, not an auxiliary runtime tool.
- Persistent mapping state is required for restart safety and review/apply semantics.

### 6. Use the umbrella change to define the final-state spec map directly

This mother change is not only a narrative program document. It explicitly maps the final-state requirement changes across these capability surfaces:

- `datalink-ui`
- `datalink-workbench-desktop`
- `source-rule-runtime`
- `tag-dictionary`
- `database-target-workbench`
- `local-modbus-memory-workbench`
- `device-registry`
- `protocol-connectors`
- `mapping-pipeline`
- `datalink-api`

The child changes can then implement subsets of this target state in sequence.

### 7. Enforce child-change gates through explicit contract inheritance

Later child changes are not allowed to invent their own workflow, ownership, revision, or candidate contracts.

Each child change MUST:

- name its prerequisite child change(s)
- reference the inherited umbrella contract in its proposal and design
- avoid contradicting any delivered ownership, revision, gating, or apply semantics from an earlier child change

### 8. Canonical rule-owned identity, signature, and snapshot completeness are umbrella-level contracts

The umbrella change fixes these shared terms so child changes do not redefine them:

- **Rule-owned identity**: the stable identity for a rule-derived downstream object across revisions. At minimum it includes `source_rule_id`, `candidate_type`, `derived_from_rule_address`, and target binding scope.
- **Proposed signature**: the deterministic fingerprint of the effective candidate payload for one rule-owned identity.
- **Out of sync**: the same rule-owned identity still exists, but its proposed signature differs from the last applied signature.
- **Snapshot complete**: the current revision has a persisted tag candidate set plus an explicit `ready`, `blocked`, or `deferred` entry for each supported downstream target. Omission is invalid.

These definitions are fixed before `rule-driven-foundation` starts and are inherited by all later child changes.

## Risks / Trade-offs

- **[Umbrella and child spec overlap]** The umbrella and child changes may both mention the same capabilities.  
  **Mitigation:** The umbrella defines final-state contract and phase ownership; child changes implement bounded slices of that contract.

- **[Rule-centric design increases orchestration complexity]** Candidate generation, revision tracking, and apply semantics add a new layer.  
  **Mitigation:** Keep `SourceRule` as the top workflow object and avoid introducing a second top-level aggregate.

- **[Database and Local Modbus diverge operationally]** Equal first-class targets still have different readiness and conflict models.  
  **Mitigation:** Keep one shared candidate contract and isolate target-specific blocking/apply rules.

- **[Large refactor scope can blur execution order]** Cross-cutting work often gets re-opened mid-implementation.  
  **Mitigation:** Use the approved child change order and require later changes to inherit the earlier ownership/revision foundations.

- **[Existing lower-level CRUD still exists]** The codebase will temporarily contain both old and new operating models.  
  **Mitigation:** Child changes must progressively remove manual-first assumptions from `/studio` while leaving lower-level tooling only where intentionally needed.

## Migration Plan

This umbrella change is not deployed on its own. Its rollout happens through the child changes in this order:

1. `rule-driven-foundation`
   - establish ownership, revision, and consistency rules
   - align device probe gating and rule orchestration
2. `tag-review-first-flow`
   - convert Step 3 into candidate review/apply
3. `database-output-review-apply`
   - add rule-scoped database candidate/apply behavior
4. `local-modbus-review-apply`
   - add persistent Local Modbus mapping and candidate/apply behavior
5. `studio-orchestration-cleanup`
   - remove manual-first workflow assumptions from `/studio`

### Child-change ownership and dependency matrix

| Child change | Primary ownership surfaces | Depends on | Relationship | Phase exit criteria |
| --- | --- | --- | --- | --- |
| `rule-driven-foundation` | `source-rule-runtime`, `device-registry`, `protocol-connectors`, shared `mapping-pipeline` identity/diff semantics, shared `datalink-api` candidate/recompute contracts | umbrella only | foundation for all later phases | candidate store, readiness gating, canonical identity/signature, and snapshot completeness are implemented and accepted |
| `tag-review-first-flow` | `tag-dictionary`, Step 3 behavior in `datalink-workbench-desktop`, tag-specific `datalink-api` apply contracts | `rule-driven-foundation` | prerequisite for both output phases | Step 3 uses candidate review/apply and preserves compatible review decisions |
| `database-output-review-apply` | `database-target-workbench`, database-specific `datalink-api` apply contracts | `rule-driven-foundation`, `tag-review-first-flow` | sibling of Local Modbus phase, not its prerequisite | database candidates use inherited candidate contracts and support explicit revalidation/apply |
| `local-modbus-review-apply` | `local-modbus-memory-workbench`, Local Modbus-specific `datalink-api` apply contracts | `rule-driven-foundation`, `tag-review-first-flow` | sibling of Database phase, not its prerequisite | Local Modbus candidates use inherited candidate contracts and support conflict-governed apply |
| `studio-orchestration-cleanup` | `datalink-ui`, shell and route cleanup in `datalink-workbench-desktop`, remaining legacy path cleanup | all earlier child changes | final cleanup and consolidation phase | `/studio` is the only primary product workflow and legacy manual-first assumptions are removed or redirected |

Rollback strategy is phase-local, not umbrella-local:

- each child change must define its own deploy and rollback boundary
- the umbrella change remains the architecture contract even if execution pauses between phases

The umbrella rollback rule is:

- rollback restores the last successfully applied state per target
- unapplied candidate snapshots from the rolled-back revision are removed
- Database and Local Modbus rollback independently while still honoring shared rule ownership and revision metadata

## Open Questions

- The umbrella no longer leaves workflow order, rule-owned identity/signature, snapshot completeness, or phase dependency rules open for debate.
- Child changes still need to detail local implementation choices such as payload field expansion, persistence table layout, and UI component boundaries.
- Those implementation details remain non-blocking only while they stay within the workflow, ownership, revision, readiness, and rollback contracts defined here.
