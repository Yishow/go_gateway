## Context

This change is child phase 3 under umbrella `rule-centric-studio-refactor`, depends on `rule-driven-foundation` and `tag-review-first-flow`, and is a sibling of `local-modbus-review-apply` rather than its prerequisite.

The prior phases establish source-rule revisions, shared candidate contracts, and effective tag review/apply. This phase uses those contracts to make Database output a first-class rule-scoped review/apply target without depending on Local Modbus delivery.

Database output is operationally distinct because it depends on shared connector context rather than a register allocator. The design therefore focuses on context validity, revalidation, and target-isolated apply behavior.

## Goals / Non-Goals

**Goals:**

- Make database output candidates consume the effective tag state for a source-rule revision.
- Keep connector, schema, table, and column context explicit and reviewable.
- Revalidate blocked or `out_of_sync` candidates when connector context becomes valid again.
- Add database-specific apply APIs with per-item results.

**Non-Goals:**

- Implement Local Modbus apply behavior.
- Redefine shared candidate identity, readiness, or recompute semantics.
- Collapse database output into generic unmanaged CRUD.
- Perform final route cleanup for `/studio`.

## Decisions

### 1. Database output remains a sibling phase to Local Modbus

This change depends on foundation plus tag review, but it does not depend on Local Modbus. The two output phases consume the same shared candidate contracts and effective tag state while keeping target-specific apply logic isolated.

### 2. Connector context is explicit review scope, not implicit fallback

Database apply always operates within the selected connector, schema, table, and column context. If that context becomes invalid, the system marks the candidate or applied state `blocked` or `out_of_sync` rather than silently rebinding to another destination.

### 3. Revalidation is automatic once context becomes valid again

Operators should not have to recreate database candidates after a connector, schema, or table issue is fixed. The phase therefore restores candidates to apply-eligible state when the original referenced context becomes valid again.

### 4. Database apply is explicit and per-target

Applying database candidates does not apply Local Modbus candidates. Responses must be per-item so failures are isolated and later repair does not require replaying successful rows.

## Risks / Trade-offs

- **[Database context can drift outside the workbench]** External schema or connector changes can invalidate candidates unexpectedly.  
  **Mitigation:** keep invalidation explicit and support automatic revalidation when the same context becomes valid again.

- **[Database target state is richer than tag state]** Connector, schema, table, and column context add more failure modes.  
  **Mitigation:** keep the target review surface layered and preserve blocking reasons with the candidate state.

- **[Output siblings can diverge in UX detail]** Database and Local Modbus have different target rules.  
  **Mitigation:** keep shared revision/candidate contracts inherited from earlier phases and isolate only target-specific apply behavior here.

## Migration Plan

1. Bind database candidates to the effective post-tag-review state for the active rule revision.
2. Add connector-aware invalidation and automatic revalidation behavior.
3. Add database-specific apply APIs with per-item result reporting.
4. Validate sibling independence by confirming Local Modbus state is preserved during database apply.

Rollback for this phase removes database apply behavior and returns database output to candidate-only state while preserving the shared revision and tag contracts from earlier phases.

## Open Questions

- No blocking questions remain for this phase.
- Detailed payload expansion and UI layout choices remain local implementation work as long as the target stays revision-scoped, connector-aware, and independently applyable.
