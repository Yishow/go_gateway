## Purpose

Provide an ownership-preserving lifecycle for saved recording configurations. Operators can distinguish ordinary edits, stopping intake, removing configuration, and reviewed successor transitions without deleting recorded SQL history.

## ADDED Requirements

### Requirement: Server classified recording actions

The system SHALL expose revision-scoped lifecycle actions with `action`, `allowed`, `reason_code`, `next_action`, and `requires_review`. Name, interval, and compatible membership edits SHALL create a draft revision of the same group and preserve the prior applied snapshot. Layout, fixed storage identity, and recreation of a removed recording SHALL use the successor contract rather than ordinary update. Changed semantics SHALL invalidate dependent evidence before new intake uses them.

#### Scenario: Edit while the old revision is recording

- **WHEN** an operator changes the saved interval while a group has an applied revision
- **THEN** the draft and applied rules remain distinguishable, old accepted data retains its original revision, and the edited rules do not become effective without current readiness and explicit review.

##### Example: Edit while the old revision is recording

- **GIVEN** group G1 runs revision R1 every 15 seconds
- **WHEN** the operator saves R2 with a 30-second interval
- **THEN** R1 remains effective until reviewed apply; R1 effects retain their original destination

#### Scenario: An unavailable action has a repair route

- **WHEN** stale connector identity or unproven ownership prevents removal
- **THEN** the system returns a safe typed reason and the available stop-intake action without bypassing the failed guard.

##### Example: An unavailable action has a repair route

- **GIVEN** G1 retains connector revision C1 while the live connector is C2
- **WHEN** the operator inspects removal actions
- **THEN** remove reports stale connector identity and stop using G1/C1 remains available


### Requirement: Stop and remove preserve accepted data

Stopping SHALL prevent new intake while preserving immutable versions, accepted samples, outbox rows, destination identity, and receipts. Removing SHALL retain a readable tombstone and operation history and SHALL NOT delete destination tables or recorded SQL rows. Removal of an applied group MUST prove backlog ownership within the configuration transaction. Stop SHALL remain available using the group's saved connector revision after a live connector is edited or removed.

#### Scenario: Stop with queued deliveries

- **WHEN** an operator stops a group with queued or retrying rows
- **THEN** no new intake is accepted and existing workers continue delivery to each row's frozen destination without relabeling it as new-generation data.

##### Example: Stop with queued deliveries

- **GIVEN** G1 has queued row Q1 bound to R1 and destination D1
- **WHEN** stop commits for G1
- **THEN** new intake is rejected and Q1 remains eligible for delivery to D1

#### Scenario: Ownership proof fails during removal

- **WHEN** an accepted sample or undelivered row cannot be resolved to its immutable owner
- **THEN** removal fails without changing the group, backlog, or historical SQL data.

##### Example: Ownership proof fails during removal

- **GIVEN** pending Q1 names a missing immutable revision R0
- **WHEN** remove G1 is requested
- **THEN** removal fails and G1 plus Q1 remain unchanged


### Requirement: Stable recording lineage with distinct successor generations

The system SHALL preserve the existing canonical workspace, device/circuit, profile-kind, and role identity across a persisted lineage. Each successor SHALL receive a new group ID and increasing generation; predecessor group IDs, versions, and receipt identities MUST NOT be reused. The canonical identity registry SHALL reference current and reserved successor groups without duplicating their configuration. Only a reviewed transition SHALL transfer the canonical head. Ordinary create retries SHALL NOT revive tombstones or create successors.

Same-shape recreation SHALL also support a tombstoned, canonical Basic-only `managed_wide_v1` recording with no `RecordingProfile` on its existing SQLite or PostgreSQL plain-wide contract. It SHALL preserve the canonical device/role, complete device/point/tag source set, source and mapping evidence, member and measurement-binding meaning, interval, reduction, row/write policy, layout, and plain mode. Where an applied revision exists, its scoped immutable definition SHALL remain the semantic source basis; an unapplied draft MUST NOT disguise changed source or recording meaning, and existing exact-source guards SHALL remain enforced. Missing or incompatible source evidence SHALL cause a typed refusal. This bounded recreation SHALL NOT introduce Basic layout/profile/mode conversion, partial-source successors, custom or unkeyed recreation, MySQL or SQL Server wide, cross-engine transitions, destination replacement capabilities, or generic ownership exemptions; the Explicit supported transition matrix and existing legacy wide/custom conversion contracts SHALL remain unchanged.

#### Scenario: Recreate a removed recording

- **WHEN** an operator explicitly reviews and confirms recreation for a tombstoned canonical identity
- **THEN** a distinct generation is created under the same logical identity, the tombstone remains readable, and existing accepted data keeps the predecessor owner.

##### Example: Recreate a removed recording

- **GIVEN** lineage L1 generation 1 points to tombstone G1
- **WHEN** explicit recreate operation O2 is confirmed
- **THEN** generation 2 uses new group G2 under L1 while G1 stays readable

#### Scenario: Recreate a Basic recording with the same complete shape

- **WHEN** an operator requests explicit recreation of a tombstoned canonical Basic-only `managed_wide_v1` recording on SQLite or PostgreSQL with unchanged complete device/source identity and recording semantics
- **THEN** the system SHALL reserve a distinct non-accepting draft group and next generation under the same lineage, preserve the original key, current head, tombstone, history, and accepted data, and assign a new owned table/storage identity without copying prior schema proof; intake SHALL remain blocked until fresh review, schema preparation and confirmation, and all existing readiness and live-transition gates pass.

##### Example: Recreate a Basic recording with the same complete shape

- **GIVEN** SQLite lineage L1 has tombstoned Basic-only G1 at generation 1 with `managed_wide_v1`, no `RecordingProfile`, device D1, complete sources P1/T1 and P2/T2, and immutable applied rules R1
- **WHEN** operation O2 requests the same format, both original sources, and R1 recording semantics
- **THEN** reserved G2 uses generation 2 and a new owned wide-table identity with no applied revision or intake; G1, its canonical key, original table, accepted data, and receipts remain unchanged, and fresh review/schema proof is still required before recording

#### Scenario: Concurrent recreation and retry

- **WHEN** the same recreate operation is retried or a competing operation targets the same canonical identity
- **THEN** identical intent resolves to the same successor and different concurrent intent receives a typed conflict without a duplicate writer.

##### Example: Concurrent recreation and retry

- **GIVEN** operation O2 reserves generation 2 for L1
- **WHEN** O2 is repeated and conflicting operation O3 targets L1
- **THEN** O2 returns the same G2 and O3 receives a conflict


### Requirement: Explicit supported transition matrix

New successor transitions SHALL support full-group `managed_wide_profile_v1/plain` to `managed_long_v1/plain` and the reverse on SQLite and PostgreSQL only, preserving profile kind, source identity, reduction, and confirmed counter meaning. Reverse conversion SHALL require an equivalent wide representation for every member. Same-layout, same-storage-mode destination replacement SHALL support verified same-engine destinations, with MySQL, MariaDB, and SQL Server independently verified and restricted to managed long. Existing legacy wide/custom-to-long contracts SHALL remain available under their existing gates. All other new transition combinations SHALL be rejected with typed reasons.

The supported matrix SHALL be enforced as follows; each allowed row still requires current proof and review.

For this delivery only, the owner authorized skipping actual SQL Server acceptance if a disposable instance cannot be established on the existing computer without another host/container. In that case, passing local adapter and proof tests SHALL satisfy the implementation-only completion gate and unblock subsequent work, but actual SQL Server write/readback and cleanup MUST remain explicitly NOT RUN. This delivery exception SHALL NOT weaken the product matrix, runtime capability/TLS/schema/ownership verification, or other database acceptance requirements; release handoff MUST disclose the unverified SQL Server scope.

| Source | Target | Allowed destination engines | Condition |
| --- | --- | --- | --- |
| Managed wide profile/plain | Managed long/plain | SQLite, PostgreSQL | Complete same-kind equivalent sources |
| Managed long/plain | Managed wide profile/plain | SQLite, PostgreSQL | Every source has a verified equivalent wide representation |
| Managed layout/mode | Same layout/mode on another connector | Same engine only; MySQL, MariaDB, SQL Server are long-only | Independent destination capability and schema proof |
| Removed managed recording | New generation using one of the preceding targets | Same supported matrix | Explicit recreate, never ordinary ensure |
| Legacy wide/custom | Long | Existing supported contract only | Existing review and ownership gates retained |
| Any other combination | Any | None in this change | Typed rejection before effects |

#### Scenario: Bidirectional managed layout transition

- **WHEN** a SQLite or PostgreSQL managed plain group with compatible members is reviewed for the other managed layout
- **THEN** the system prepares a new table generation, displays the complete transition impact, and keeps the original SQL table and historical rows unchanged.

##### Example: Bidirectional managed layout transition

- **GIVEN** SQLite G1 uses managed_wide_profile_v1/plain with instantaneous state role
- **WHEN** a full-group transition to managed_long_v1/plain is confirmed
- **THEN** new generation G2 preserves latest state meaning and G1 SQL history is not copied or deleted

#### Scenario: Unsupported layout or shape

- **WHEN** a request selects MySQL or SQL Server wide, a cross-engine target, a storage-mode change, a partial-source successor, a profile-kind change, or a non-equivalent wide member shape
- **THEN** the request is blocked before DDL or intake mutation and the response identifies the unsupported condition without silently converting values or modes.

##### Example: Unsupported layout or shape

- **GIVEN** a SQL Server group uses managed_long_v1/plain
- **WHEN** a transition to managed_wide_profile_v1/plain is requested
- **THEN** unsupported-layout is returned before table creation or intake mutation

#### Scenario: Same-mode long destination replacement

- **WHEN** a verified same-engine managed-long destination on MySQL, MariaDB, SQL Server, SQLite, or PostgreSQL is selected with unchanged storage mode
- **THEN** its connector revision, capability, schema, and ownership are independently confirmed and no predecessor proof is reused for the new target.

##### Example: Same-mode long destination replacement

- **GIVEN** a MariaDB long/plain group writes to connector A revision A1
- **WHEN** connector B revision B1 of the same engine is independently verified and confirmed
- **THEN** the successor writes only newly owned effects to B; finalized A effects remain at A


### Requirement: Reserved preparation and atomic effective ownership

A successor draft SHALL reserve only the reviewed lineage operation and SHALL NOT accept source data before the confirmed effective boundary. The predecessor and reserved successor overlap exception SHALL be limited to that lineage; unrelated same-source same-profile owners MUST still be rejected. Review and confirmation SHALL bind workspace, canonical head, source and target revisions, connector revisions, exact boundary, and review digest. New live transitions MUST remain unavailable until the frozen-counter-replay and existing-cutover-boundary correctness dependencies pass their relevant acceptance checks.

Preparation SHALL use explicit stages on the existing successor endpoints. The initial `successor-review` SHALL be a coherent local read with no reservation, persisted schema token, external inspection, or DDL. The `successor` endpoint SHALL distinguish `reserve`, `schema-preview`, and `schema-confirm`; unknown actions and activation before the later live-transition task SHALL be rejected. A post-reservation schema preview SHALL inspect the target and, on success, persist the existing revision-bound preview token but SHALL NOT execute DDL.

The immutable reservation envelope and its versioned digest SHALL bind the original reviewed intent, operation, scope, exact boundary, and original expected revisions. Later preparation stages SHALL separately validate current revisions without rewriting that replay identity. Schema preparation SHALL reuse the existing schema-operation ledger, bind its token to the exact reserved group and source identity, and distinguish the schema operation ID from the lineage transition operation ID. Lost acknowledgement or an outcome made uncertain by a concurrent revision change SHALL retain the original ledger result, not imply that no DDL occurred or authorize a new operation.

Active-predecessor preparation SHALL use a distinct server-authorized, same-transaction, same-lineage reservation path limited to the existing managed-profile matrix and complete source meaning. It SHALL NOT weaken the deleted-only recreate contract, ordinary reserved guards, or unrelated owner conflicts, and SHALL NOT add active Basic conversion capability. Schema success SHALL only record verified preparation proof; it SHALL NOT move the canonical head, create an applied version or ownership fence, stop the predecessor, or enable successor intake. Activation SHALL require a fresh review and explicit confirmation of current evidence and the effective boundary in the later live-transition task.

#### Scenario: Cancel an initial successor review

- **WHEN** an operator reviews a successor intent and leaves before confirming reservation
- **THEN** the local database is unchanged, no reserved generation or schema token exists, and no target inspection or DDL has occurred.

##### Example: Cancel an initial successor review

- **GIVEN** active G1 has no reserved successor and the initial review names boundary B1
- **WHEN** the operator closes that review without sending `reserve`
- **THEN** G1, its lineage, and all schema operations remain unchanged, and no G2 is created

#### Scenario: Prepare exactly one reserved successor

- **WHEN** a valid reserve confirmation is repeated with the identical frozen envelope and then its schema preparation succeeds
- **THEN** the same reserved generation is returned, current-stage revisions are checked independently, and only its verified schema proof is saved; ordinary routes cannot activate it or reuse its authority for another owner.

##### Example: Prepare exactly one reserved successor

- **GIVEN** operation T1 reserves G2 from active G1 and schema operation S1 belongs to the exact G2 preview
- **WHEN** S1 succeeds and the original T1 reserve request is repeated after the workspace revision advances
- **THEN** T1 still resolves to G2, S1 remains a distinct schema operation, G1 keeps its lifecycle, and G2 remains reserved with no applied revision or intake

#### Scenario: Reject foreign schema authority

- **WHEN** a schema token or operation belongs to another reservation, generation, target, source, or bound revision
- **THEN** the successor adapter rejects it before new target work and does not change the reservation, canonical head, or intake state.

##### Example: Reject foreign schema authority

- **GIVEN** T1 reserves G2 and T2 reserves a different lineage's G9
- **WHEN** the G2 schema confirmation supplies G9's preview token
- **THEN** it is refused without executing new DDL or modifying either lineage

#### Scenario: Preparation fails

- **WHEN** target schema preparation or capability verification fails
- **THEN** the predecessor continues according to its existing lifecycle and no successor intake becomes active.

##### Example: Preparation fails

- **GIVEN** G1 is active and reserved G2 has no confirmed schema
- **WHEN** G2 schema creation is denied
- **THEN** G1 remains active and G2 accepts no new samples

#### Scenario: Confirmed transition crosses the boundary

- **WHEN** a valid confirmed transition reaches its effective boundary
- **THEN** source `observed_at < boundary` selects the old owner and `observed_at >= boundary` selects the new owner regardless of received, accepted, acknowledgement, or reconcile time; delayed reconciliation and late unfinalized observations follow that boundary and lateness contract through the existing persisted ownership fence; unfinalized observations accepted before reconciliation SHALL be atomically handed off when required, preserving original acceptance provenance, frozen binding proof, and audited disposition without deletion or double counting, while finalized effects, outbox rows, and receipts remain frozen.

##### Example: Confirmed transition crosses the boundary

- **GIVEN** boundary is 12:00:12.500Z and an unfinalized sample observed at 12:00:13Z is buffered before delayed reconciliation
- **WHEN** reconciliation executes at 12:00:18Z
- **THEN** the sample is audited into the new consuming owner using observed_at; a finalized old receipt is unchanged

#### Scenario: Successor counter establishes its own baseline

- **WHEN** a new cumulative successor receives its first valid post-boundary sample
- **THEN** that sample establishes the successor baseline with NULL delta, identical binding or epoch evidence does not silently carry an old baseline, and the system does not fabricate usage across a gap.

##### Example: Successor counter establishes its own baseline

- **GIVEN** the predecessor last value is 100 and the successor first valid value is 107 after the boundary
- **WHEN** the successor forms its initial counter result
- **THEN** 107 becomes its baseline and delta is NULL, not a fabricated cross-generation value of 7

#### Scenario: A reviewed revision becomes stale

- **WHEN** any bound revision or canonical head changes before confirmation
- **THEN** confirmation fails with a typed conflict and performs no new transition.

##### Example: A reviewed revision becomes stale

- **GIVEN** review binds source R1 and target T1
- **WHEN** the source becomes R2 before confirmation
- **THEN** the R1 review is rejected without scheduling a transition


### Requirement: Historical destinations and pending units remain explicit

The lifecycle review SHALL list predecessor and successor destination identities, affected sources, generations, boundary, and pending evidence grouped by immutable revision. Samples, outbox rows, and blocked or unknown delivery states SHALL be separately labeled and SHALL NOT be summed as a single record count. Review, draft saved, schema confirmed, intake enabled, SQL committed, readback verified, and cleanup completed SHALL remain distinct facts.

#### Scenario: Remove or transition with mixed pending evidence

- **WHEN** a group contains collecting samples and queued or unknown outbox rows from older revisions
- **THEN** the confirmation displays their distinct units and original destinations and does not imply that removal clears them or that transition retargets finalized rows; any unfinalized observation handoff is separately identified with its original provenance and consuming owner.

##### Example: Remove or transition with mixed pending evidence

- **GIVEN** R1 has 3 collecting samples and 2 queued outbox rows including 1 unknown outcome
- **WHEN** the operator opens impact review
- **THEN** the counts remain separately labeled and the UI does not claim 5 homogeneous records


### Requirement: Successor operations survive interrupted navigation

The system SHALL retain the same transition and operation identity after lost responses, reloads, scope changes, or restarts. An unknown outcome SHALL be reconciled from the server ledger before any replay; a permitted retry SHALL reuse the exact frozen request. Canceling an unsubmitted review SHALL have no effects, while reversing a confirmed transition SHALL require a separately reviewed forward boundary and SHALL NOT rewrite history.

#### Scenario: Lost confirmation response

- **WHEN** confirmation commits but its response is lost
- **THEN** status recovery finds the original transition by its bound identity and does not create a new generation or reapply an unrelated recent transition.

##### Example: Lost confirmation response

- **GIVEN** operation O2 commits transition T2 but its HTTP response is lost
- **WHEN** the client queries O2 and repeats its identical frozen request if allowed
- **THEN** the ledger returns T2 and generation 2, never generation 3

#### Scenario: Operator leaves a pending transition

- **WHEN** an operator switches groups or reloads while a transition is unresolved
- **THEN** the transition remains accessible with its original target and cannot be replaced by a fresh preview that hides the unresolved outcome.

##### Example: Operator leaves a pending transition

- **GIVEN** T2 is unresolved for G1
- **WHEN** the operator views G9 and reloads
- **THEN** T2 remains reachable with G1 and its original target identity


### Requirement: Delivery acceptance evidence remains explicit

For this delivery only, the owner's explicit browser acceptance exception on 2026-10-28 SHALL allow completion of the local implementation and dependent change handoff without actual browser execution or screenshots. The formal tests, per-scenario status, reproducible commands, and environment blocker evidence MUST remain available. Actual browser interaction, layout, keyboard/focus, and screenshots MUST remain SKIP or NOT RUN, never PASS based on discovery, jsdom tests, type checks, lint, or builds. The exception SHALL NOT weaken any product capability, revision, ownership, confirmation, or backend correctness prerequisite. The separate SQL Server exception SHALL retain its original bounded scope.

#### Scenario: Browser execution is skipped under the explicit delivery exception

- **WHEN** the local implementation checks pass and the approved existing-computer browser environment remains unavailable
- **THEN** the delivery SHALL retain all 42 prepared formal browser cases as NOT RUN, disclose that screenshots and actual browser behavior are unverified, and allow the adjusted implementation handoff without representing skipped checks as passed.

##### Example: Browser execution is skipped under the explicit delivery exception

- **GIVEN** discovery finds 42 formal cases and local regression reports 315 files and 2,942 tests passed
- **WHEN** the authorized browser exception is applied for this delivery
- **THEN** browser and screenshot results remain NOT RUN; the local pass counts do not become browser acceptance evidence, and any separately reproduced unresolved pre-existing ordinary-start race remains disclosed
