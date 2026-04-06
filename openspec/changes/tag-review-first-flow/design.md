## Context

This change is child phase 2 under umbrella `rule-centric-studio-refactor` and depends on `rule-driven-foundation`.

The foundation phase defines revisions, candidate snapshots, and shared identity/signature rules, but it intentionally leaves target-specific apply flows for later phases. `tag-review-first-flow` is the first consumer of that foundation and establishes the operator-facing review/apply behavior that both output phases will inherit.

Step 3 is therefore not just a UI reshuffle. It is the phase that turns rule-generated tag intent into explicit operator-reviewed state, while preserving meaningful decisions across compatible revisions.

## Goals / Non-Goals

**Goals:**

- Make Step 3 operate on tag candidates and pending mapping intent from the active source-rule revision.
- Persist rename, skip, and override decisions across compatible revisions.
- Block tag apply on stale candidate revisions.
- Add explicit tag apply APIs with per-item results.

**Non-Goals:**

- Implement database output apply behavior.
- Implement Local Modbus output apply behavior.
- Consolidate final `/studio` routing and legacy redirects.
- Redefine candidate identity or readiness semantics from the foundation phase.

## Decisions

### 1. Tag review owns the first downstream apply boundary

This phase is where a rule-derived candidate first becomes active downstream state. Tag apply creates or updates the active tag record and the active mapping intent for the approved candidate set.

This keeps the workflow aligned with the approved model: save rule first, review tags second, review outputs third.

### 2. Review-decision persistence follows the canonical identity contract

Rename, skip, and override decisions are attached to the rule-owned tag identity from the foundation phase. If the identity remains materially the same, the decision carries forward. If the identity changes, the old decision becomes stale and the operator must review again.

### 3. Stale candidate revisions block tag apply

Tag review cannot apply a stale candidate set. Step 3 must compare the open review revision against the latest persisted source-rule revision and require refresh before apply when they diverge.

### 4. Tag apply is target-specific and does not auto-apply outputs

This phase adds tag apply only. Database and Local Modbus remain reviewable downstream targets that later phases will handle independently. Tag apply may recompute downstream candidate state, but it must not silently apply outputs.

## Risks / Trade-offs

- **[Decision carry-forward can surprise operators]** Preserving rename or skip decisions is useful but can hide underlying rule changes.  
  **Mitigation:** treat identity changes as stale and surface revision-aware review state clearly.

- **[Step 3 gets more stateful]** Review/apply introduces more explicit states than the old immediate-sync model.  
  **Mitigation:** keep the state model anchored on the persisted foundation revision and candidate contracts.

- **[Tag apply may trigger downstream candidate churn]** Output candidates depend on effective tag state.  
  **Mitigation:** update downstream candidates only as candidate state; later output phases remain the only apply owners for their targets.

## Migration Plan

1. Update rule-generated tag behavior to remain candidate-only before apply.
2. Persist review decisions keyed by the canonical rule-owned identity.
3. Update Step 3 to load candidate state, surface exception actions, and block stale apply.
4. Add explicit rule-scoped tag apply APIs.
5. Validate carry-forward and stale blocking before starting output-focused child changes.

Rollback for this phase removes tag review/apply behavior and returns tag handling to the foundation candidate-only state without changing readiness or candidate contracts.

## Open Questions

- No blocking questions remain for this phase.
- Later output phases still need to define how they consume the effective post-tag-review state, but they must reuse the same revision and identity contracts from the foundation phase.
