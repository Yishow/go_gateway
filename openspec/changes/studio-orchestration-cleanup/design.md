## Context

This change is child phase 5 under umbrella `rule-centric-studio-refactor` and depends on `rule-driven-foundation`, `tag-review-first-flow`, `database-output-review-apply`, and `local-modbus-review-apply`.

The earlier child changes define the runtime foundation and the rule-centric tag/output flows, but they intentionally do not finish route consolidation or remove all legacy primary-path assumptions. This final cleanup phase converts those delivered flows into one coherent product experience centered on `/studio`.

The key risk is doing cleanup too early. This phase therefore assumes the earlier child changes have already delivered the replacement rule-centric path and focuses only on consolidating the product shell and de-emphasizing manual-first tooling.

## Goals / Non-Goals

**Goals:**

- Make `/studio` the only primary route for the product workflow.
- Redirect legacy `/datalink/workbench` entry points into `/studio`.
- Demote manual-first tooling from the normal product path while keeping it available as secondary engineering tooling where needed.
- Remove legacy shell assumptions that suggest multiple equal primary workbench routes.

**Non-Goals:**

- Redefine foundation, tag review, or output apply contracts.
- Reintroduce separate primary routes for database or Local Modbus.
- Remove lower-level engineering tools entirely if they are still intentionally retained as secondary utilities.

## Decisions

### 1. `/studio` becomes the sole primary route

This phase finishes the product-level route decision made by the umbrella change. Any legacy `/datalink/workbench` entry must redirect into `/studio` instead of behaving like another primary workflow.

### 2. Manual-first tooling becomes secondary, not deleted by default

The cleanup phase does not need to destroy every lower-level tool. It needs to stop presenting them as equal primary workflow choices. Engineering utilities may remain, but they are no longer the normal operator path.

### 3. Shell cleanup happens only after replacement flows exist

Cleanup is sequenced last because it depends on the earlier phases already providing a usable rule-centric replacement. This prevents the product from losing workflow coverage during the transition.

## Risks / Trade-offs

- **[Cleanup can break discoverability]** Redirecting routes and demoting manual tooling can confuse existing operators.  
  **Mitigation:** preserve context on redirect where possible and keep secondary engineering surfaces explicitly labeled.

- **[Legacy UI assumptions can be scattered]** Manual-first behavior may still appear in multiple shell components.  
  **Mitigation:** scope this phase to guided workflow copy, route ownership, and shell-level entry behavior rather than reopening lower-level logic.

- **[Secondary tooling may still leak primary language]** Old labels and navigation can continue to imply multiple primary routes.  
  **Mitigation:** update shell, navigation, and guided workflow text together in this phase.

## Migration Plan

1. Redirect legacy `/datalink/workbench` entry points into `/studio`.
2. Update guided workflow and shell requirements to make `/studio` primary.
3. Demote manual-first tooling from normal-path navigation and copy.
4. Validate route redirects and product-path consistency after earlier phases are in place.

Rollback for this phase restores the prior routing and presentation layer only; it does not remove the earlier rule-centric behavior already delivered by prerequisite phases.

## Open Questions

- No blocking questions remain for this phase.
- Any retained secondary engineering route can vary in implementation as long as `/studio` remains the sole primary product route.
