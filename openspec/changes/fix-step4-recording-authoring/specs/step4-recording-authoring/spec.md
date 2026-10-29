## Purpose

Make Step 4 recording choices match the selected source meaning and verified destination capabilities. Preserve unresolved test-write evidence independently of the currently viewed editor.

## ADDED Requirements

### Requirement: Numeric state is an explicit latest-value role

The guided source editor SHALL offer `state` with latest-value semantics for numeric machine states, `gauge` with mean semantics for compatible numeric measurements, and `counter_raw` only with confirmed counter meaning. Non-numeric sources SHALL expose only compatible roles. Existing saved roles SHALL remain unchanged on navigation, reload, or introduction of new options; changing a saved role SHALL require the lifecycle review and frozen-binding protections.

#### Scenario: Integer machine state is not averaged

- **WHEN** an operator explicitly chooses state for integer samples 1 and 2 in a window
- **THEN** the authored rule uses latest state 2 rather than mean 1.5 and the review states that choice.

##### Example: Integer machine state is not averaged

- **GIVEN** source machine_state has type int16 and samples 1 and 2
- **WHEN** the operator selects state
- **THEN** the draft says latest and the reviewed result is 2

#### Scenario: Saved role stays locked without a role-changing review capability

- **WHEN** an operator reopens a saved numeric role or navigates between Basic, Advanced, and the supported lifecycle layout or destination review without an available role-changing review capability
- **THEN** the editor SHALL retain and lock the saved role, explain that role changes are unsupported, and expose only supported lifecycle navigation without implying role-editing authority
- **AND** navigation, reload, and review dismissal SHALL NOT change the saved role or frozen binding, and SHALL NOT submit an ensure or update fallback to change that role

##### Example: Saved gauge cannot become state through layout review

- **GIVEN** saved group G1 has numeric source S1 with role gauge and an unchanged frozen binding
- **WHEN** the operator opens the supported layout review and returns to Basic without a role-changing review capability
- **THEN** S1 remains gauge, its role control stays locked with an unsupported explanation, and no role-change mutation is sent

#### Scenario: Gauge and counter remain distinct

- **WHEN** an operator chooses gauge or counter_raw for a numeric source
- **THEN** gauge keeps its supported mean behavior and counter_raw cannot be prepared until its real counter definition and binding are confirmed.

##### Example: Gauge and counter remain distinct

- **GIVEN** numeric temperature has no counter definition
- **WHEN** gauge is selected and then counter_raw is selected
- **THEN** gauge uses mean; counter_raw requires a real confirmed definition before preparation


### Requirement: Counter-dependent names have visible prerequisites

Whenever counter confirmation requires a circuit name, the editor SHALL display an inline required-name reason before enabling confirmation, including instantaneous-only counter_raw recording. Ordinary non-counter instantaneous recording SHALL retain an optional circuit name. Names SHALL NOT determine counter meaning or replace stable circuit or source identities.

#### Scenario: Instantaneous-only counter needs a name

- **WHEN** counter_raw is selected without a required circuit name while cumulative recording is disabled
- **THEN** the name field and blocked confirmation show the same actionable requirement and filling it satisfies only that prerequisite without starting recording.

##### Example: Instantaneous-only counter needs a name

- **GIVEN** counter_raw is selected, cumulative is disabled, and circuit_name is empty
- **WHEN** the operator opens confirmation
- **THEN** the name is marked required with a reason and no confirmation is sent

#### Scenario: Ordinary instantaneous recording has no forced name

- **WHEN** only non-counter instantaneous sources are selected
- **THEN** an empty circuit name does not introduce a counter-confirmation blocker.

##### Example: Ordinary instantaneous recording has no forced name

- **GIVEN** a numeric gauge is selected and circuit_name is empty
- **WHEN** the operator prepares instantaneous recording
- **THEN** no counter-name requirement is introduced


### Requirement: New cumulative selection is explicit

A new cumulative profile SHALL start with no sources selected and SHALL direct the operator to choose and confirm counters. Enabling cumulative recording SHALL NOT implicitly select temperature, pressure, text, or other source candidates. Existing persisted explicit cumulative selections SHALL survive reload and toggle changes. Source selection controls SHALL NOT silently select sources for a disabled profile or a different device scope.

#### Scenario: Enable cumulative with mixed candidates

- **WHEN** a new device draft contains numeric temperature, numeric pressure, and text candidates and cumulative recording is enabled
- **THEN** none becomes a cumulative member automatically, the selected count is zero, and an actionable source-selection prompt is visible.

##### Example: Enable cumulative with mixed candidates

- **GIVEN** temperature, pressure, and batch_text are candidates in a new draft
- **WHEN** cumulative is enabled
- **THEN** cumulative_keys stays empty and the source picker explains that counters must be selected

#### Scenario: Reopen a saved cumulative group

- **WHEN** a saved cumulative selection is reloaded or its enabled toggle is changed and restored
- **THEN** the explicit saved membership is retained rather than replaced by either an empty default or all device sources.

##### Example: Reopen a saved cumulative group

- **GIVEN** saved group G1 explicitly contains counter tag K1 only
- **WHEN** G1 is reloaded and cumulative is toggled off and on
- **THEN** K1 stays selected and temperature T1 stays unselected


### Requirement: Every selected counter exposes measurement-read recovery

A measurement-definition lookup failure SHALL show a safe error and scoped retry whenever any selected source requires counter definitions, including instantaneous-only counter_raw. Retry SHALL retain the draft and refetch the failed scope only; it SHALL NOT save groups, apply schema, or start recording. Only actions dependent on the missing evidence SHALL remain blocked.

#### Scenario: Instantaneous-only lookup fails

- **WHEN** counter_raw is selected and measurement lookup fails with cumulative disabled
- **THEN** the editor exposes the failure and retry, retains the role and name, and explains why prepare or confirm is blocked.

##### Example: Instantaneous-only lookup fails

- **GIVEN** instantaneous counter_raw K1 is selected with cumulative disabled
- **WHEN** the measurement query returns an error
- **THEN** retry and a safe error are visible while K1 selection and its name remain

#### Scenario: Retry returns after navigation

- **WHEN** a lookup retry completes after the operator changes device scope
- **THEN** the response cannot overwrite the newly viewed scope or start any recording and the original scope remains recoverable.

##### Example: Retry returns after navigation

- **GIVEN** a definition retry is issued for device A
- **WHEN** the operator selects B before A replies
- **THEN** A evidence is associated only with A and no recording start occurs


### Requirement: Destination layout choices are capability bound

Both visible layout choices and outbound draft construction SHALL use capabilities bound to the saved connector ID and revision. MySQL and SQL Server SHALL NOT offer or submit new managed-wide recording. Unknown, loading, failed, stale, or foreign capability evidence SHALL NOT authorize a write. An existing unsupported configuration SHALL remain readable with an explanation and SHALL NOT be silently changed to another layout.

#### Scenario: Long-only destination

- **WHEN** verified MySQL or SQL Server capability supports managed long
- **THEN** the operator sees long as the supported recording layout and wide is unavailable with a reason before reaching an API rejection.

##### Example: Long-only destination

- **GIVEN** SQL Server connector C1 revision R1 proves managed-long support
- **WHEN** the operator chooses a new layout
- **THEN** long is selectable and wide cannot be submitted

#### Scenario: Connector revision changes

- **WHEN** the saved connector revision changes while capability lookup is in flight
- **THEN** the old capability response cannot enable a layout or schema action for the new destination.

##### Example: Connector revision changes

- **GIVEN** capability query is issued for C1/R1
- **WHEN** C1 is changed to R2 before the R1 response arrives
- **THEN** the R1 result does not authorize any R2 schema operation


### Requirement: Unresolved test writes remain independently recoverable

Test-write operation identity SHALL be retained separately from editor preview state and bound to workspace, group, revision, and target. Pending, running, unknown-write, and failed or unknown-cleanup outcomes SHALL remain queryable through the server ledger after reload, scope change, or local-reference loss. New preview for the same group SHALL NOT replace unresolved evidence. Write, readback, and cleanup results SHALL be displayed separately; cleanup repair SHALL reuse a server-authorized existing operation and SHALL NOT insert another test row.

#### Scenario: Lost response followed by new preview

- **WHEN** a test-write confirmation loses its response and the operator attempts another preview
- **THEN** the original operation and target remain visible, the new preview is blocked for that unresolved group, and the operator can query the same operation.

##### Example: Lost response followed by new preview

- **GIVEN** test-write O1 targets G1/R1 and its confirmation response is lost
- **WHEN** the operator clicks new preview for G1
- **THEN** O1 remains queryable and no replacement preview is created

#### Scenario: Scope change or local storage loss

- **WHEN** the operator changes group revision, reloads, or loses browser operation references
- **THEN** a workspace-scoped recovery query can locate the original unresolved operation without displaying foreign-workspace evidence or assuming that an absent local reference proves no write occurred.

##### Example: Scope change or local storage loss

- **GIVEN** O1 is unresolved in workspace W1 and browser references are cleared
- **WHEN** the W1 recovery view loads
- **THEN** O1 is recovered from the server ledger without exposing operations from W2

#### Scenario: Write succeeded but cleanup did not

- **WHEN** write and readback are confirmed but cleanup is failed or unknown
- **THEN** the UI reports both facts, keeps the operation unresolved for recovery, and offers only supported same-operation cleanup repair or safe diagnostic guidance rather than another insert.

##### Example: Write succeeded but cleanup did not

- **GIVEN** O1 reports written_verified and cleanup failed
- **WHEN** the operator requests the supported cleanup repair
- **THEN** the same O1 is repaired and no additional test insert is sent

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
