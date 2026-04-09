# V2 Calm Summary Harvest Design

- Status: Brainstorming-approved; pending spec-review
- Date: 2026-04-09
- Scope: v2-only follow-up after Phase 6 compare, limited to Shell and Output summary surfaces

## 1. Summary

Phase 6 compare confirmed that `v2 / Sentry Incident Desk` should remain the canonical `/studio` owner, but it also identified a clear follow-up opportunity: some `v2` summary surfaces are louder than they need to be.

The approved follow-up is therefore a **small calming pass**, not a redesign:

- keep `v2` unmistakably `v2`
- harvest selected calm-framing strengths from `v1`
- touch only `Shell + Output summaries`
- avoid route, API, contract, or step-ownership changes

This work is intentionally narrow. It is meant to reduce noise in the highest-visibility summaries without weakening blocker clarity, action hierarchy, or the incident-desk workflow model.

## 2. Problem Statement

The final compare locked three conclusions that matter here:

1. `v2` is still the best-balanced product surface for operator efficiency.
2. `v1` contains reusable calm-framing strengths, especially in blocker wording and summary pacing.
3. `v3` is too dense to reuse wholesale.

The problem is not that `v2` has the wrong architecture. The problem is that some of its Shell and Output summaries push too much incident intensity into surfaces that should support the workflow rather than compete with it.

That creates two follow-up costs:

1. **Summary noise cost**
   - status, blocker, revision, and attention signals all fight too hard for first-glance priority.

2. **Readability cost**
   - `v2` remains fast, but some summary layers feel more escalated than necessary during steady-state review.

3. **Harvest-delay cost**
   - if the compare-approved `v1` strengths are not folded back in, the project keeps known quality improvements outside the canonical track.

## 3. Goals

### Product goals

- Keep `v2` clearly recognizable as the incident-desk canonical owner.
- Reduce summary noise in Shell and Output without reducing surfaced state.
- Preserve the shortest blocker-triage and return path.
- Make the summary language feel calmer without becoming vague or passive.

### Architecture goals

- Keep the change isolated to `v2` summary surfaces.
- Avoid touching shared API, shared domain state, or route behavior.
- Avoid changing the ownership boundary of Shell or Output.

### Delivery goals

- Land a low-risk follow-up that can be validated quickly.
- Preserve room for a future Source/Tag summary pass if the user wants it later.

## 4. Non-Goals

- Reopening the Phase 6 winner decision.
- Making `v2` behave like `v1`.
- Reworking Source or Tag in this same follow-up.
- Reordering command-dock actions or changing Output workboard behavior.
- Fixing Output reload recovery semantics in this pass.
- Introducing a new shared tone framework before implementation.

## 5. Approved User Decisions

The approved follow-up direction is:

1. harvest `v1` calm framing into `v2`
2. keep `v2` obviously incident-desk in tone and structure
3. use **copy plus small summary layout adjustments**
4. prioritize **Shell + Output summaries**
5. avoid larger rail or shell restructuring

## 6. Approaches Considered

### Approach A — Small calming pass

Update only `v2` Shell and Output summary surfaces, keeping the existing incident-desk structure and action order intact.

- **Pros:** lowest risk, shortest path, matches the approved user direction
- **Cons:** leaves Source and Tag unchanged for now

### Approach B — Medium sweep across all four phases

Apply the same calming pass to Shell, Source, Tag, and Output together.

- **Pros:** more unified tone
- **Cons:** larger regression surface and more difficult validation

### Approach C — Shared summary contract first

Define a new shared tone/summary layer before applying changes to `v2`.

- **Pros:** clean long-term structure
- **Cons:** turns a focused follow-up into architecture work

The chosen approach is **Approach A — Small calming pass**.

## 7. Design

### 7.1 Boundaries

This follow-up is intentionally limited to these units:

1. `frontend/src/pages/datalink/workbench/MuiWorkbenchIncidentStrip.tsx`
2. `frontend/src/pages/datalink/workbench/MuiOutputIncidentDesk.tsx`
3. related `v2` i18n strings
4. targeted tests that lock summary behavior

This follow-up does **not** modify:

- shared board logic in `LocalModbusBoard` or `DatabaseTargetBoard`
- step navigation or active-step ownership
- Shell action placement
- Output command-dock action order
- shared tokens or backend contracts

### 7.2 Shell summary changes

`MuiWorkbenchIncidentStrip` keeps the same four shell-owned surfaces:

- readiness summary
- active blocker summary
- diagnostics refresh status
- shortest return-to-mainline action

The change is presentational and copy-focused:

1. **Readiness line**
   - keep the same information
   - reduce “incident bulletin” intensity
   - make it read more like stable operational context

2. **Blocker copy**
   - keep the next action explicit
   - soften wording toward calm operational guidance
   - do not hide severity when the tone is warning or critical

3. **Refresh status**
   - keep diagnostic identity and state changes
   - lower its visual competition with the blocker copy

4. **Actions**
   - remain in the same order and position
   - no ownership changes

### 7.3 Output summary changes

`MuiOutputIncidentDesk` keeps the existing four-part structure:

1. priority card
2. summary strip
3. command dock
4. handoff panel

The calming pass works as follows:

1. **Priority card stays primary**
   - it remains the first-glance area for blocker/readiness interpretation
   - no structural demotion

2. **Summary strip becomes more contextual**
   - keep linked/scoped/attention/revision information
   - lower the visual priority of revision and attention when they are not the main decision driver
   - make the strip feel like supporting context rather than a second priority panel

3. **Handoff panel becomes calmer**
   - preserve the `Go to Tag` repair path
   - rewrite the copy toward calmer repair guidance
   - do not weaken the requirement to return when prerequisites are missing

4. **Command dock remains intact**
   - no button removal
   - no action reorder
   - only reduce surrounding summary noise so the dock does not compete with the priority card

### 7.4 Data flow and state invariants

No domain or data-flow meaning changes in this follow-up.

The following must remain identical before and after the change:

- active-step transitions
- blocker detection
- review-set status handling (`ready`, `blocked`, `deferred`)
- retry/refetch behavior
- Local Modbus and Database target switching
- return-to-Tag behavior

If the calmer copy causes a state to become less actionable or less explicit, the change is incorrect.

### 7.5 Error handling and edge cases

These states must stay explicit and actionable:

- no selected device
- no focused rule
- stale review data
- blocked output review
- deferred output review
- candidates query failure
- refresh in progress

This follow-up may soften phrasing, but it may not:

- remove the reason text
- hide the required next step
- demote a true blocker to decorative secondary copy

### 7.6 Deferred work

The compare found that Output recovery after reload still depends on restoring context through `Tag -> Output` handoff.

That issue is explicitly **out of scope** for this follow-up. This spec only records it so the implementation does not quietly expand into recovery redesign.

## 8. Testing and Validation

### 8.1 Targeted test updates

Update or add targeted tests that prove:

1. Shell actions and ownership remain unchanged.
2. Output command-dock order remains unchanged.
3. Summary and handoff copy changes do not break Local Modbus review.
4. Summary and handoff copy changes do not break Database review.
5. Existing blocker, stale, deferred, and retryable states still render clearly.

### 8.2 Validation workflow

The implementation validation should run:

1. targeted `v2` tests for Shell and Output surfaces
2. `tsc --noEmit`
3. targeted ESLint for changed files
4. `npm run build`
5. browser validation on real `/studio` data using `UI 4.3 Modbus TCP`

### 8.3 Browser evidence expectations

The acceptance evidence should focus on:

- Shell summary after calming pass
- Output summary strip and handoff panel for Local Modbus
- Output summary strip and handoff panel for Database
- proof that `v2` still reads as `v2`, not as a relabeled `v1`

## 9. Risks and Mitigations

### Risk 1 — Over-softening `v2`

If the copy becomes too calm, `v2` loses its incident-desk identity.

**Mitigation:** keep structure, action order, and primary-card hierarchy unchanged.

### Risk 2 — Summary drift into redesign

A small calming pass could accidentally turn into broader phase refactoring.

**Mitigation:** confine the implementation to `MuiWorkbenchIncidentStrip`, `MuiOutputIncidentDesk`, related i18n strings, and targeted tests.

### Risk 3 — Hidden behavior regressions

Even copy-level changes can break test selectors or state assumptions.

**Mitigation:** update targeted tests first and validate with the same real browser path used in Phase 6.

## 10. Success Criteria

This follow-up is successful when all of the following are true:

1. `v2` still reads immediately as the canonical incident-desk product.
2. Shell and Output summaries feel calmer and less noisy than before.
3. Blockers, retries, and return actions remain explicit and fast to interpret.
4. No API, route, or ownership semantics change.
5. Targeted tests and real-browser validation both pass.
