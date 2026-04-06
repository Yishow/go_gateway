## Context

This change is child phase 1 under umbrella `rule-centric-studio-refactor` and has no earlier child prerequisite.

The umbrella change locked the final-state workflow and clarified that `rule-driven-foundation` owns the shared contracts that every later child change inherits. This phase must deliver the runtime and API baseline for revision-aware orchestration before any tag or output apply workflow is implemented.

The current system already persists source rules and already separates connect and probe semantics in parts of the codebase, but those behaviors are not yet normalized into one contract that later phases can depend on. This change therefore focuses on persistence boundaries, readiness semantics, and cross-target candidate metadata instead of user-facing apply flows.

## Goals / Non-Goals

**Goals:**

- Persist source-rule revisions as the authoritative orchestration boundary.
- Persist candidate snapshots with one canonical identity/signature model for later tag and output flows.
- Define readiness semantics that allow planning before probe success but block activation/apply until probe success.
- Expose shared candidate query and recompute APIs that later child changes can reuse unchanged.
- Extend mapping lifecycle semantics with `out_of_sync` handling for rule-derived assets.

**Non-Goals:**

- Implement tag apply UX or output apply UX.
- Finalize database-specific or Local Modbus-specific apply payloads.
- Remove legacy routes or cleanup `/studio` shell behavior.
- Replace lower-level CRUD APIs in this phase.

## Decisions

### 1. Source-rule save creates a new revision and a persisted candidate snapshot

This phase treats source-rule save as one orchestration boundary:

1. persist the updated rule definition
2. assign a new revision id
3. synchronize runtime rule-owned point state
4. persist a complete candidate snapshot for the revision

This keeps later review/apply flows stateless relative to rule editing and gives every child change one shared revision anchor.

### 2. Readiness is multi-dimensional, not one final boolean

This change defines separate `planning_allowed`, `activation_allowed`, and `apply_allowed` flags plus explicit `connect_status` and `probe_status`. Planning can proceed after a successful connect even if probe fails, while runtime activation and downstream apply remain blocked until probe succeeds.

This decision moves readiness away from the old point/polling-group/mapping-first assumptions and aligns it with the approved rule-centric workflow.

### 3. Canonical identity and signature belong in foundation

Rule-owned identity and proposed-signature comparison are shared semantics, not target-specific behavior. They must be defined here so tag review, database output, and Local Modbus output all inherit the same diff model instead of implementing incompatible matching logic later.

### 4. Candidate APIs stop at query and recompute in this phase

This phase adds shared query/recompute contracts but deliberately does not define target-specific apply endpoints yet. That split keeps the foundation complete while letting later child changes own target-specific apply contracts without redefining the candidate model itself.

### 5. Snapshot completeness allows blocked and deferred targets

Candidate snapshot completeness does not mean every target is apply-ready. It means the snapshot includes the current tag candidate set plus an explicit state entry for each supported downstream target, including `blocked` or `deferred` entries with reasons when context is missing.

This prevents later phases from inferring meaning from omitted data.

## Risks / Trade-offs

- **[Foundation touches multiple services]** Cross-cutting changes can drift across modules.  
  **Mitigation:** keep the phase strictly limited to shared runtime, readiness, mapping diff, and candidate API contracts.

- **[Legacy readiness consumers may expect one boolean]** Existing callers may assume readiness is a single pass/fail signal.  
  **Mitigation:** update readiness consumers in this phase to use the explicit contract and treat any summary boolean as derived data only.

- **[Candidate persistence increases storage complexity]** Persisted revision snapshots add new records and cleanup needs.  
  **Mitigation:** keep snapshot structure minimal and revision-scoped, with explicit ownership by source-rule runtime.

## Migration Plan

1. Add source-rule revision and candidate snapshot persistence.
2. Update readiness services and connector diagnostics to expose the shared planning/activation/apply contract.
3. Update shared mapping lifecycle semantics to support `out_of_sync`.
4. Add candidate query and recompute APIs.
5. Validate restart restore, readiness gating, and candidate recompute behavior before starting `tag-review-first-flow`.

Rollback for this phase restores the previous source-rule persistence path and removes newly created candidate snapshots if the foundation contract is withdrawn before downstream phases depend on it.

## Open Questions

- No blocking questions remain for this phase.
- Later child changes still need to define target-specific apply payload expansion, but they must not redefine readiness, identity, signature, or snapshot completeness.
