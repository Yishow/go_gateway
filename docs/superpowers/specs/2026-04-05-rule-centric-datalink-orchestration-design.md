# Rule-centric Datalink Orchestration Design

- Status: Brainstorming-approved, spec-review passed; pending user review
- Date: 2026-04-05
- Scope: `/studio` end-to-end workflow across device readiness, source-rule lifecycle, derived tag review, database target output, and local Modbus output

## 1. Summary

This design converts the datalink workbench into a rule-centric orchestration flow.

`Device` and `SourceRule` are the only primary workflow objects the operator creates to define the data path. `Point` stays in the system as a runtime/internal read model for collection and scheduler integration, but it is no longer a primary operator concept. `DatabaseConnector` remains a separately persisted environment object that provides output context, but it is not part of the rule-owned workflow chain. Rule saves persist device/rule state and synchronize runtime points. Tags and output bindings are generated as reviewable candidates and are applied explicitly in Step 3 and Step 4.

Database and Local Modbus are equal first-class outputs. The first formal launch of this workflow does not need to preserve a manual point/tag/mapping-first mental model because that legacy flow is not yet a committed public product path.

## 2. Problem Statement

The current system is rule-aware but not rule-primary.

1. Source rule creation already auto-generates downstream `Point`, `Tag`, and `Mapping` state, but that generation happens immediately during rule persistence.
2. The UI still exposes manual `Point`, `Tag`, and `Mapping` CRUD as first-class workflow objects.
3. Database output remains a manual per-tag binding flow with no upstream rule context.
4. Local Modbus mappings remain manual and runtime-oriented instead of rule-driven and persistently reviewable.
5. Rule changes do not produce a formal candidate diff for downstream review.
6. Operators still reconstruct relationships step-by-step instead of approving system-derived candidates.

The result is a split-brain workflow:

- the backend partially behaves as if `SourceRule` is the primary aggregate
- the UI and API still encourage direct manipulation of lower-level derived objects
- output targets remain detached from the rule lifecycle

## 3. Goals

### Product goals

- Make `/studio` the one continuous operator workflow for source-to-output setup.
- Treat `SourceRule` as the primary acquisition object after device setup.
- Reframe Step 3 as review-first tag confirmation rather than manual tag construction.
- Reframe Step 4 as dual-target output review/apply for both Database and Local Modbus.

### Architecture goals

- Separate operator-owned objects from runtime-owned objects and review/apply objects.
- Add explicit ownership and revision tracking so downstream state can be regenerated safely.
- Keep Database and Local Modbus output state isolated from each other while sharing one rule-centric origin.
- Avoid introducing a new top-level `Pipeline` aggregate when a rule-centric design is sufficient.

### Operational goals

- Separate `connect` and `probe` clearly in the device step.
- Make conflict handling local to the affected target instead of blocking the whole workflow.
- Support batch-safe apply and re-apply after rule changes without silent overwrites.

## 4. Non-Goals

- Introducing a new top-level `Pipeline` aggregate or replacing the whole domain model with a second orchestration model.
- Keeping manual point/tag/mapping-first flows as first-class `/studio` workflow paths.
- Supporting database connector kinds beyond `SQLite` and `PostgreSQL` in the primary output flow.
- Reworking unrelated scheduler internals beyond what is required for rule-centric runtime synchronization.
- Preserving non-persistent Local Modbus mapping behavior.

## 5. Approved Decisions

1. The scope is a full domain/persistence/runtime/UI redesign of the workflow, not only a UI cleanup.
2. Database and Local Modbus are equal first-class outputs.
3. The design should be an umbrella blueprint with phased delivery, not one oversized implementation slice.
4. `SourceRule` automatically generates tag candidates, and Step 3 becomes review/rename/exception handling.
5. Rule revisions automatically generate both Database and Local Modbus output candidates, and Step 4 becomes review/apply; output apply remains blocked until the related tag candidate is approved or applied.
6. The activation model is hybrid: devices and rules persist earlier, while tags and outputs are reviewed and applied later.
7. Rule edits regenerate downstream diffs without silently overwriting already applied state.
8. The new workflow can be rule-driven only for the first formal launch; no legacy compatibility requirement is needed for manual-first flows.

## 6. Chosen Approach

Three approaches were considered:

1. **Conservative flow cleanup**
   - Keep the current CRUD-centric backend model and mostly improve UI sequencing.
   - Lowest risk, but it does not solve model drift between rules, points, tags, mappings, and outputs.

2. **Rule-centric orchestration**
   - Keep `Device` and `SourceRule` as the operator-facing persisted core.
   - Push `Point` into runtime/internal responsibilities.
   - Generate `Tag` and output candidates from rule revisions and apply them explicitly.

3. **Pipeline aggregate**
   - Introduce a higher-level flow object above device, rule, tags, and outputs.
   - Most complete on paper, but too large and too disconnected from current domain foundations.

The chosen approach is **Rule-centric orchestration** because it matches the approved user direction, aligns with the current backend trajectory, and avoids inventing an entirely new top-level aggregate.

## 7. Target Architecture

### 7.1 Layer model

The workflow is organized into five layers:

1. **Primary workflow layer**
   - `Device`
   - `SourceRule`

2. **Shared environment layer**
   - `DatabaseConnector`

3. **Runtime layer**
   - `Point`
   - `PollingGroup`

4. **Review/Apply layer**
   - `TagCandidate`
   - `DatabaseOutputCandidate`
   - `LocalModbusOutputCandidate`

5. **Applied layer**
   - `Tag`
   - `Point <-> Tag Mapping`
   - `DatabaseTargetMapping`
   - `LocalModbusMapping`

### 7.2 Responsibility boundaries

`Device` owns protocol capability, connect/probe readiness, address semantics, and default data-format behavior.

`SourceRule` owns source planning intent: start address, span count, data type, naming rules, scaling/casting intent, and enabled/disabled state.

`DatabaseConnector` is a persisted environment object. It is configured independently from a rule, but rule-owned database candidates reference it as context. If connector readiness, schema state, or selected table state changes, affected database candidates and applied mappings become `blocked` or `out_of_sync`; they are never silently rebound to a different connector/table.

`Point` remains necessary for scheduler/runtime behavior, but the operator should not manually create it in the primary `/studio` workflow.

Candidates are not a second domain root. They are projections derived from `SourceRule` revisions and used for review/apply decisions.

### 7.3 Ownership and revision tracking

All applied downstream objects must be traceable back to the rule that produced them. The exact storage shape may be normalized tables or persisted snapshots, but the contract must support:

- `managed_by_source_rule_id`
- `derived_from_rule_address`
- `last_applied_rule_revision`

Each candidate item must carry enough information to support safe diff calculation:

- `source_rule_id`
- `rule_revision`
- `proposed_signature`
- `effective_signature`
- `last_applied_signature`
- `blocking_reason`
- `override_state`
- `suppression_state`

This keeps `SourceRule` as the primary aggregate while still enabling review/apply workflows and revision-aware diffs.

## 8. Workflow Design

### Step 1: Device

- Show `connect` and `probe` as separate results.
- Allow saving a device when `connect` succeeds even if `probe` fails.
- Block `rule enable`, runtime collection, and downstream apply actions until `probe` succeeds.
- Keep device capability context visible to all later steps: address model, word order, unit/station identity, protocol traits, and supported data-format defaults.

### Step 2: Source Rule

- The operator creates and edits `SourceRule`, not `Point`.
- Saving a rule persists the rule immediately and synchronizes derived runtime points.
- The source canvas shows planning state, readiness state, conflicts, and the presence of downstream candidates.
- Every rule save increments a revision and triggers candidate recomputation.

### Step 3: Tag review

- Step 3 is no longer a manual tag-construction page.
- It presents `TagCandidate` items with:
  - proposed key
  - display name
  - data type
  - conflict state
  - diff against the last applied version
- The operator can batch approve, rename, skip, or resolve exceptions.
- Rename, skip, and field overrides are stored as rule-scoped review decisions so they survive later revisions until the operator clears them or the candidate identity changes materially.

### Step 4: Output review/apply

- Step 4 presents two parallel target surfaces in one workbench step:
  - `Database`
  - `Local Modbus`
- Both surfaces consume the same rule-derived candidate origin, but each target maintains isolated selection and apply state.
- Output candidates are generated at rule-recompute time and can be previewed immediately, but apply remains blocked until the related tag candidate is approved or applied.
- Output apply is target-specific:
  - Database failures do not block Local Modbus apply.
  - Local Modbus conflicts do not block Database apply.

### Step dependency rule

Output candidates may be visible before final tag apply, but output apply is blocked until the relevant tag is confirmed as usable by the rule-driven flow. Output candidates always derive from the effective tag state for the current revision, not from a stale pre-rename tag proposal.

## 9. Candidate / Diff / Apply Model

All rule-derived items share one status vocabulary:

| Status | Meaning |
| --- | --- |
| `proposed` | Latest candidate generated from the current rule revision |
| `applied` | Candidate has been formally applied |
| `out_of_sync` | Applied item exists, but the current rule revision proposes a different result |
| `conflict` | Candidate cannot be safely applied |
| `blocked` | Candidate is valid in principle but a prerequisite is missing |
| `overridden` | Operator has made a manual adjustment that the system must not silently overwrite |

### 9.1 Diff logic

The UI must not infer diffs by comparing ad-hoc visible fields. Diff state must be derived from persisted signatures:

- `proposed_signature`
- `effective_signature`
- `last_applied_signature`

`proposed_signature` represents the pure system derivation from the current rule revision. `effective_signature` represents the same candidate after persisted operator decisions such as rename or field override. `last_applied_signature` records what was last accepted into applied state.

This is required to support safe regeneration after rule edits and to avoid hidden partial overwrites.

### 9.2 Apply logic

- **Tag apply** creates or updates the formal `Tag` and `Point <-> Tag Mapping`.
- **Database apply** creates or updates the formal `DatabaseTargetMapping`.
- **Local Modbus apply** creates or updates the formal `LocalModbusMapping`.

Apply can happen in batches or per-item, but every apply operation must be explicit and visible to the operator.

### 9.3 Rule edit behavior

When a rule changes:

- the system regenerates candidates for the new revision
- already applied downstream objects are not silently overwritten
- affected items become `out_of_sync`
- manual overrides remain preserved and marked for review

### 9.4 Rename, skip, and override rules

- **Rename / field override**
  - stored as a rule-scoped operator override on the candidate
  - updates `effective_signature`
  - survives later revisions until the operator clears the override or the candidate identity changes materially

- **Skip**
  - stored as a suppression decision on the candidate
  - prevents accidental auto-apply for that item
  - survives later revisions while the candidate still represents the same rule-owned item
  - becomes `out_of_sync` and requires re-review if the candidate identity changes materially

- **Output regeneration**
  - output candidates derive from the effective tag state, including approved renames and persisted overrides
  - changing tag review decisions therefore recomputes downstream output candidates for the same rule revision

This preserves the approved design principle: auto-derive first, then show diff, then let the operator decide whether to apply.

## 10. Backend Services and Contracts

The backend should be organized around five capability groups.

### 10.1 Rule Orchestrator

Responsibilities:

- persist `SourceRule`
- synchronize derived runtime points
- trigger candidate recomputation
- enforce probe gating for enable/apply eligibility

Consistency contract:

- rule persistence, derived point synchronization, and revision increment are one orchestration boundary; if that boundary fails, the rule save fails
- candidate recomputation may complete immediately after the orchestration commit, but it must write one complete snapshot for the new revision before Step 3 or Step 4 apply actions become eligible
- mixed-revision candidate exposure is not allowed; while recomputation is pending, the rule is visible as `candidate_recompute_pending` and downstream apply remains blocked

### 10.2 Candidate Projector

Responsibilities:

- generate `TagCandidate`, `DatabaseOutputCandidate`, and `LocalModbusOutputCandidate`
- attach revision metadata and signatures
- calculate diff and conflict state

### 10.3 Apply services

Responsibilities:

- apply approved tag candidates
- apply approved database output candidates
- apply approved local Modbus output candidates
- return explicit success/failure summaries

### 10.4 Output-specific engines

- `DatabaseProposalEngine`
  - uses connector/schema/table context
  - proposes column bindings
  - validates connector readiness and table state

- `LocalModbusAllocator`
  - proposes register allocations
  - detects collisions
  - supports dry-run and auto-map strategies

### 10.5 Revision and ownership tracking

Applied downstream records must be queryable and explainable through rule ownership metadata and last-applied revision metadata.

### 10.6 API direction

The `/studio` workflow must stop treating lower-level CRUD as its orchestration backbone. The rule-centric API surface should add:

- `GET /source-rules/:id/candidates`
- `POST /source-rules/:id/candidates/recompute`
- `POST /source-rules/:id/tag-candidates/apply`
- `POST /source-rules/:id/database-candidates/apply`
- `POST /source-rules/:id/modbus-candidates/apply`

Existing lower-level CRUD may remain for lower-level tools or transition support, but the primary `/studio` workflow should not depend on `points`, `tags`, and `mappings` CRUD as if they were the main operator journey.

## 11. Output Target Design

### 11.1 Database

Database output remains a three-layer workflow:

1. connector
2. schema/table context
3. tag-to-column mapping

The primary connector scope in this design is limited to:

- `SQLite`
- `PostgreSQL`

`DatabaseConnector` persists independently from `SourceRule`. Rule-owned database candidates reference connector and table context explicitly. For the MVP:

- connector identity is persisted independently
- selected schema/table/column context is persisted as part of the rule-owned database candidate and applied mapping state
- candidate generation must respect the current connector and table context and must not reuse stale state from previously selected connectors or tables
- if a connector becomes disabled, unreadable, or loses the selected table/column, affected candidates become `blocked` and applied mappings become `out_of_sync`

### 11.2 Local Modbus

Local Modbus output requires a persistent mapping model. In-memory-only register mapping is not sufficient for a review/apply workflow.

This design therefore requires:

- a persistent `LocalModbusMapping` repository
- rule-aware proposal generation
- register conflict detection
- dry-run validation
- allocator strategies such as sequential, gap-aware, or aligned

The MVP allocator strategy is `sequential`. Alternative strategies such as `gap-aware` and `aligned` can remain behind the same allocator interface but are not required for the first implementation slice.

The Local Modbus server lifecycle remains separate from mapping proposal state. Starting or stopping the server must not be the source of truth for whether a mapping exists.

### 11.3 Target independence

Database and Local Modbus must remain equal first-class outputs, but their readiness and conflicts are isolated:

- a database connector problem blocks only database apply
- a register conflict blocks only Local Modbus apply

The workbench still shows an overall output summary, but target-local failures must remain target-local.

## 12. Lifecycle and Failure Rules

### 12.1 Connect and probe

- `connect fail`: device cannot be saved and downstream flow is blocked
- `connect success + probe fail`: device can be saved, but rule enable, collection, and downstream apply are blocked

### 12.2 Rule enable and disable

- `enable` requires successful probe readiness
- `disable` stops collection and runtime ingest
- `disable` does not delete already applied tags or outputs; those items remain visible and should be marked as sourced from an inactive rule

### 12.3 Rule delete

Rule delete is destructive and must show impact before confirmation. The confirm flow must summarize:

- derived runtime points to remove
- tag candidates to discard
- applied downstream assets still managed by the rule

On confirmed delete, the system removes:

- the rule record
- derived runtime points
- persisted candidate snapshots for that rule
- applied downstream assets that are still marked as managed by that rule

Delete must fail fast if the system cannot safely remove a still-managed downstream asset. Partial silent cleanup is not allowed.

### 12.4 Conflict handling

- tag conflicts block only the affected tag apply
- database readiness or schema conflicts block only database apply
- local Modbus register collisions block only Local Modbus apply

### 12.5 No silent overwrite rule

No rule change may silently replace applied downstream behavior. The system must surface `out_of_sync`, `conflict`, or `overridden` states explicitly and require an operator decision.

### 12.6 Lifecycle matrix

| Event | Immediate system action | Downstream result |
| --- | --- | --- |
| Create rule | Persist rule, sync runtime points, compute candidate snapshot | Tag/output candidates appear as `proposed` or `blocked` |
| Edit rule | Persist new revision, sync runtime points, recompute candidates | Applied downstream items become `out_of_sync` if signatures differ |
| Disable rule | Stop collection and runtime ingest | Applied tags/outputs remain visible and marked as source inactive |
| Reprobe success | Clear probe block for the device | Previously blocked rule enable/apply actions become eligible |
| Connector invalidation | Revalidate database candidate scope | Only database candidates/mappings become `blocked` or `out_of_sync` |

## 13. Testing Strategy

### 13.1 Unit tests

- candidate projector
- diff/signature comparison
- rule-driven tag naming and collision handling
- database proposal logic
- local Modbus allocation logic

### 13.2 Integration tests

- rule create/update/enable/disable
- runtime point synchronization
- candidate regeneration on rule revision change
- apply behavior for tag/database/modbus targets
- ownership and revision metadata persistence

### 13.3 End-to-end tests

1. Device connect/probe gating
2. Source rule creation and revision update
3. Tag review-first flow
4. Database apply
5. Local Modbus apply
6. Rule edit causing `out_of_sync` downstream state
7. Target-local conflicts leaving the other target available

## 14. Phased Delivery

### Phase 1: Ownership, revision, and persistent output foundations

- add ownership and revision metadata for rule-managed downstream assets
- introduce persistent Local Modbus mapping storage
- align existing schemas with rule-managed traceability

### Phase 2: Candidate projector and candidate APIs

- add candidate projection generation
- add rule-scoped candidate query/recompute endpoints
- keep lower-level CRUD available but no longer primary for `/studio`

### Phase 3: Tag review-first step

- convert Step 3 into `TagReviewBoard`
- remove manual-first tag creation from the primary workflow
- surface tag diff, collision, rename, and apply results

### Phase 4: Dual-target output review/apply

- convert Step 4 into one `OutputWorkspace` with isolated target state
- add database proposal/apply flow
- add local Modbus proposal/apply flow

### Phase 5: Workflow cleanup

- remove old point/tag/mapping manual-first workflow assumptions from `/studio`
- keep only lower-level tooling where it is still intentionally needed outside the main product flow

## 15. Why This Design Fits the Current Codebase

The repository already contains important pieces of the target direction:

- source rules are persisted
- device probe gating already exists
- rule save already knows how to derive runtime points and tag/mapping relationships

What is missing is not the idea of rule-driven orchestration. What is missing is the explicit candidate/review/apply layer and the output integration that make the workflow safe, understandable, and operator-first.

This design keeps the healthy foundations, removes the current split-brain workflow, and provides a path to implement the product as one coherent `/studio` experience.
