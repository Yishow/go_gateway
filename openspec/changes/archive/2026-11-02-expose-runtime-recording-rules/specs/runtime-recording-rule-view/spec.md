## Purpose

Explain the effective device and database recording rules in the focused runtime monitor. Connect revision-bound configuration to durable delivery evidence without changing runtime behavior or revealing connection secrets.

## ADDED Requirements

### Requirement: Scoped device rule inspection

The system SHALL provide a workspace-owned, selected-device rule view containing source identities, register addresses, protocol functions, read and result types, effective polling interval provenance, enabled state, safe decoding and transformation details, and Tag bindings. Saved facts and effective runtime facts SHALL be explicitly distinguished. Lists SHALL support bounded search and pagination rather than an inaccessible remainder count.

#### Scenario: More rules than the summary limit

- **WHEN** a device has six rules and ten mappings
- **THEN** the operator can inspect every rule and mapping through the bounded list
- **AND** selecting another device does not include the previous device's rule details or late responses

##### Example: More rules than the summary limit

- **GIVEN** Device A owns rules A1 through A6 and device B owns B1
- **WHEN** the operator opens A6 and then switches to B
- **THEN** B1 is shown and a delayed A6 response does not replace it

#### Scenario: Polling and writing use different schedules

- **WHEN** a point polls every second and its group writes a mean every fifteen seconds
- **THEN** the view identifies the polling interval and the recording interval separately with their authoritative sources
- **AND** it does not infer one from the other

##### Example: Polling and writing use different schedules

- **GIVEN** Point P has polling interval 1000 ms and group G has mean interval 15 seconds
- **WHEN** the rule details load
- **THEN** the acquisition field reads 1000 ms from its polling group and recording reads mean / 15 seconds

#### Scenario: Effective transformation cannot be proven

- **WHEN** a saved source rule differs from the runtime projection and the effective transform is unavailable
- **THEN** the saved transform remains labeled saved-only and the effective transform is unavailable
- **AND** no draft transform is presented as currently executing

##### Example: Effective transformation cannot be proven

- **GIVEN** saved rule revision R2 has scale 0.1 and the runtime only identifies projection R1
- **WHEN** the effective transform is requested without an R1 snapshot
- **THEN** R2 scale is saved-only and the active scale is unavailable

### Requirement: Revision-bound recording rule inspection

The system SHALL expose group membership, safe destination identity, table and column mapping, layout, storage mode, recording interval, latest or mean reduction or cumulative-delta semantics, write policy, draft and applied revisions, and effective boundary. Effective recording facts SHALL come from immutable applied snapshots and verified runtime layouts.

#### Scenario: Draft destination changed while old rules run

- **WHEN** draft revision B changes table and interval while applied revision A remains active
- **THEN** the active rule shows revision A's destination and interval
- **AND** the draft changes and activation boundary remain separately identified

##### Example: Draft destination changed while old rules run

- **GIVEN** applied A writes table meter_old every 15 seconds and draft B targets meter_new every 30 seconds
- **WHEN** runtime is still executing A
- **THEN** active shows meter_old / 15 seconds and draft shows meter_new / 30 seconds

#### Scenario: Snapshot changed during a read

- **WHEN** a read cannot establish one consistent workspace and applied revision
- **THEN** the API returns a typed stale-reference result and the UI offers a read-only refresh
- **AND** it does not combine facts from incompatible revisions

##### Example: Snapshot changed during a read

- **GIVEN** rule metadata is read at workspace revision W1 and group evidence arrives at W2
- **WHEN** the API assembles a rule view
- **THEN** it returns typed stale-reference instead of a W1/W2 combined result

### Requirement: Delivery evidence remains authoritative

The system SHALL show current-revision stages, historical backlog, source windows, receipt identities, pending age, quota and storage events from the durable delivery authorities. Accepted samples, queues, healthy collectors and logs SHALL NOT establish SQL commitment. SQL commitment SHALL NOT establish independent readback verification or cleanup completion.

#### Scenario: Disabled group still drains old backlog

- **WHEN** intake is disabled and records remain for an earlier connector revision and table
- **THEN** the view distinguishes stopped intake from outstanding delivery and retains the original backlog destination
- **AND** a draft destination change does not relabel those records

##### Example: Disabled group still drains old backlog

- **GIVEN** disabled G has three pending records frozen to connector C1 revision 7 and table old_table
- **WHEN** draft G targets connector C2
- **THEN** backlog remains C1 / revision 7 / old_table / three records

#### Scenario: Missing or uncertain evidence

- **WHEN** a source has no data, unknown delta, partial mean coverage, an unmatched receipt, or an unsupported per-source evidence layout
- **THEN** the view presents the corresponding typed missing, unknown, partial or unavailable state
- **AND** it does not substitute zero, complete coverage or a successful write

##### Example: Missing or uncertain evidence

- **GIVEN** source S reports delta_known=false and source T reports mean coverage=partial
- **WHEN** the source detail is rendered
- **THEN** S is unknown rather than zero and T remains partial rather than complete

#### Scenario: Shared group totals are not device totals

- **WHEN** a group contains sources from multiple devices
- **THEN** the view marks group-level totals as shared
- **AND** device or Tag success is asserted only from matching source evidence

##### Example: Shared group totals are not device totals

- **GIVEN** G contains device A and B and has ten committed group rows but no per-source evidence
- **WHEN** the operator selects A
- **THEN** ten is labeled a shared group total and A has no inferred source commit count

#### Scenario: Timing is unavailable

- **WHEN** the scheduler or outbox does not supply next-due, next-retry or attempt-duration evidence
- **THEN** those fields are unavailable
- **AND** the browser does not manufacture an execution countdown from the configured interval

##### Example: Timing is unavailable

- **GIVEN** G has a 15-second interval and next_retry_at is absent
- **WHEN** the UI opens at 10:00:01Z
- **THEN** next retry is unavailable rather than an invented 10:00:15Z deadline

### Requirement: Exact safe return to settings

The system SHALL use the shared EntityLocator contract with applicable step, workspace_id, device_id, rule_id, mapping_id and group_id fields. Ownership SHALL be checked before selection. A group_revision selection SHALL be informational and SHALL NOT replace mutation revision checks.

#### Scenario: Return to one rule or write group

- **WHEN** the operator follows a rule, mapping or group action from the runtime view
- **THEN** Studio selects the exact workspace-owned entity in its owning step
- **AND** Back and Forward preserve the selected runtime context

##### Example: Return to one rule or write group

- **GIVEN** workspace W owns device D, rule R and group G
- **WHEN** the operator follows the group setting action
- **THEN** Studio receives workspace_id=W and group_id=G for Step 4 without an activation request

#### Scenario: Entity was removed or belongs elsewhere

- **WHEN** an entity locator is deleted, stale or foreign to the workspace
- **THEN** the UI explains the safe failure and offers the matching workspace list
- **AND** it does not select another entity by name or apply configuration

##### Example: Entity was removed or belongs elsewhere

- **GIVEN** group G1 is deleted and G2 uses its old display name
- **WHEN** a historical G1 locator is opened
- **THEN** the deleted G1 explanation is shown and G2 is not selected

### Requirement: Safe focused monitoring and truthful freshness

The system SHALL prioritize device rules, database recording rules and recent delivery evidence in the focused monitor, while keeping reports secondary. Operator API and DOM projections SHALL exclude credentials, network endpoints, private file paths, DSNs, raw SQL and driver exceptions. Source register addresses and validated safe destination labels SHALL remain inspectable. Process logs SHALL retain their local-only access boundary.

#### Scenario: Stream disconnects

- **WHEN** the runtime stream disconnects after a successful snapshot
- **THEN** the UI retains that snapshot with its observation time and stale state, uses existing polling recovery and stops claiming LIVE
- **AND** browser transport events remain labeled connection events rather than database execution logs

##### Example: Stream disconnects

- **GIVEN** the last successful snapshot was observed at 10:00:00Z
- **WHEN** the stream disconnects at 10:00:05Z
- **THEN** the 10:00:00Z data remains stale, LIVE is absent and polling recovery remains available

#### Scenario: Remote monitor and hostile connector data

- **WHEN** a remote operator opens the monitor and a connector contains credentials, an endpoint or a private database path
- **THEN** the rule view exposes only the safe destination projection
- **AND** it neither requests a privileged process-log proxy nor weakens the local-only log guard

##### Example: Remote monitor and hostile connector data

- **GIVEN** connector C includes a password fixture and a private SQLite path
- **WHEN** a remote monitor requests C rule evidence
- **THEN** the response includes only safe C identity and table labels and contains neither fixture string

#### Scenario: Bounded details

- **WHEN** the caller requests a rule list without a limit or with a limit above one hundred
- **THEN** the server uses a default of fifty or clamps to one hundred respectively and returns a revision-bound cursor
- **AND** expanding one record does not trigger an unbounded per-group request fan-out

##### Example: Bounded details

- **GIVEN** a workspace contains 120 rules
- **WHEN** the caller requests limit=500
- **THEN** the page has at most 100 rules and a revision-bound cursor for the remainder


### Requirement: Browser verification exceptions remain explicit

Under the user-authorized exception for the existing cloud environment, delivery SHALL retain formal browser test cases, test discovery, E2E-specific typechecking, lint, and scenario coverage. Only actual browser execution and screenshots SHALL be eligible for a conditional delivery-gate exception when the verified environment blockers remain unchanged and execution remains unavailable. The exception SHALL NOT weaken product behavior, dependency order, ownership, backend database or frozen-proof checks, other tests or reviews, or platform and field verification boundaries. It SHALL NOT authorize a new environment or a security workaround.

#### Scenario: Blocked browser execution is reported without a false pass

- **WHEN** unchanged verified cloud blockers prevent actual browser execution or screenshot capture during this change's acceptance
- **THEN** the delivery record SHALL identify the affected cases and screenshots as NOT RUN, SKIP, or unverified, retain the blocker evidence and rerun commands, and report preparation checks separately from browser execution
- **AND** the delivery SHALL NOT report browser or screenshot PASS from discovery, typechecking, lint, simulated DOM results, or synthetic images

##### Example: Discovered case without a browser result

- **GIVEN** browser case B1 is discoverable and its E2E-specific typecheck and lint pass, but the verified launch blocker remains
- **WHEN** acceptance evidence is recorded under this authorized exception
- **THEN** B1 and its screenshot remain NOT RUN with the blocker reason; only the executed preparation checks are PASS
