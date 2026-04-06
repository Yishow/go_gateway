## Context

This change is child phase 4 under umbrella `rule-centric-studio-refactor`, depends on `rule-driven-foundation` and `tag-review-first-flow`, and is a sibling of `database-output-review-apply` rather than its prerequisite.

The earlier phases provide revision-aware candidates and effective tag state. Local Modbus must now consume those contracts while solving a target-specific problem Database does not have: register allocation conflict governance inside shared local memory space.

This phase therefore focuses on persistent rule-owned mapping state, conflict detection, and verification-first apply without redefining the shared candidate model or depending on the Database phase.

## Goals / Non-Goals

**Goals:**

- Persist Local Modbus mapping ownership and revision state for rule-derived output.
- Make Local Modbus candidate review/apply consume the active source-rule revision.
- Detect overlapping register allocations and block only the conflicting subset.
- Add Local Modbus-specific apply APIs with per-item results.

**Non-Goals:**

- Implement database output behavior.
- Redefine shared candidate identity, signature, readiness, or recompute semantics.
- Complete final `/studio` route consolidation.
- Reintroduce Local Modbus as a manual-first primary workflow.

## Decisions

### 1. Local Modbus state becomes persistent and rule-owned

This phase treats Local Modbus output state as persistent workflow state instead of in-memory runtime convenience state. Restart safety and review/apply semantics both require persisted ownership and revision metadata.

### 2. Conflict governance blocks only the conflicting subset

Multiple rules may target the same Local Modbus memory grid. The system therefore detects overlapping register ranges and marks only the conflicting candidates as `blocked_conflict`, allowing non-conflicting candidates to continue through review/apply.

### 3. Verification remains part of apply, not a separate legacy tool path

Local Modbus keeps verification-first behavior, but verification now operates on the revision-scoped candidate set and the pending apply selection rather than on ad hoc manual register edits.

### 4. Local Modbus apply is explicit and per-target

Applying Local Modbus candidates does not apply database candidates. Responses must be per-item so operators can resolve conflicts or write failures without replaying unaffected candidates.

## Risks / Trade-offs

- **[Persistent Local Modbus state adds data-model complexity]** This target historically behaved like a runtime helper.  
  **Mitigation:** keep persistence scoped to rule-owned candidate and applied mapping state only.

- **[Allocator conflicts can be noisy]** Shared local memory space can surface many conflicts at once.  
  **Mitigation:** block only the conflicting subset and keep conflict reasons explicit at candidate level.

- **[Verification can look like legacy tooling]** The workflow must not regress to manual-first editing.  
  **Mitigation:** keep verification tied to the selected candidate set and explicit apply flow.

## Migration Plan

1. Persist Local Modbus ownership and revision metadata.
2. Update Local Modbus review state to consume the active source-rule revision.
3. Add conflict detection and `blocked_conflict` handling for overlapping ranges.
4. Add Local Modbus-specific apply APIs with per-item results.
5. Validate restart restore, conflict blocking, and sibling independence from Database output.

Rollback for this phase removes Local Modbus apply behavior and returns the target to candidate-only output state while preserving earlier shared contracts.

## Open Questions

- No blocking questions remain for this phase.
- Detailed allocator heuristics remain implementation work as long as the workflow preserves explicit conflict signaling and rule-owned persistence.
