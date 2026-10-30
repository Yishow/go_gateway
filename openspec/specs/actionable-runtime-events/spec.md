# actionable-runtime-events Specification

## Purpose

Make local runtime logs identify affected device and recording rules and provide safe corrective navigation. Keep event capture bounded and separate from the authoritative configuration and delivery evidence.

## Requirements

### Requirement: Production events carry attributable evidence

The system SHALL emit typed device, source-rule, mapping, group-intake, delivery, quota, storage-repair and cutover events from their actual execution boundaries. Each event SHALL carry only verified applicable identities, revisions, stages, actions and opaque evidence references from the shared runtime-rule contract. Unproven identity SHALL be represented as unavailable rather than inferred.

#### Scenario: Recording group fails and recovers

- **WHEN** an applied group's delivery transitions from retrying to destination-confirmed commitment
- **THEN** the failure and recovery events identify the same group and applied revision with typed stage and action
- **AND** the recovery references matching durable receipt evidence

##### Example: Recording group fails and recovers

- **GIVEN** G revision R1 has a retrying effect E and later a matching committed receipt
- **WHEN** the producer observes both durable transitions
- **THEN** failure and recovery identify G/R1 and recovery references E

#### Scenario: No proof of SQL commitment

- **WHEN** input is accepted, a row is queued, an acknowledgement is uncertain or a receipt mismatches the revision
- **THEN** no SQL-committed event is emitted for that input
- **AND** a log event never substitutes for readback verification or cleanup evidence

##### Example: No proof of SQL commitment

- **GIVEN** effect E is queued and its only receipt belongs to revision R0 instead of R1
- **WHEN** the event producer checks R1 delivery
- **THEN** no R1 SQL-committed event is emitted

#### Scenario: Event capture fails

- **WHEN** the logging admission queue rejects an event after a successful delivery transition
- **THEN** delivery retains its actual outcome and the capture loss is accounted for
- **AND** the writer is not blocked or failed by logging

##### Example: Event capture fails

- **GIVEN** effect E is durably committed and the admission queue is at its existing byte cap
- **WHEN** its summary event is rejected
- **THEN** E stays committed and capture loss increases without a writer error




---
### Requirement: Event projection preserves redaction and local access

The system SHALL validate each event code and its scalar field allowlist before memory, file or SSE admission. It SHALL exclude credentials, endpoints, DSNs, private paths, raw SQL, driver exceptions, nested arbitrary payloads and sample values. The existing process-log loopback and same-origin restrictions SHALL remain unchanged.

#### Scenario: Hostile and secret-bearing fields

- **WHEN** a producer supplies an unknown code, raw exception, secret-bearing field, nested value or oversized text
- **THEN** only the safe suppressed or bounded projected event reaches any sink
- **AND** the viewer renders allowed text without executing markup

##### Example: Hostile and secret-bearing fields

- **GIVEN** a candidate contains raw_error=password_fixture and a nested credential object
- **WHEN** the projector receives it
- **THEN** neither candidate value is retained in memory, file, SSE or rendered DOM

#### Scenario: Remote viewing is denied without a proxy bypass

- **WHEN** a LAN, forwarded or foreign-origin caller requests process logs
- **THEN** the existing guard rejects the request and a denied viewer stops automatic retries
- **AND** the runtime monitor does not proxy the restricted feed and still uses its safe workspace rule and delivery APIs

##### Example: Remote viewing is denied without a proxy bypass

- **GIVEN** the log request has a non-loopback peer or an X-Forwarded-For header
- **WHEN** the guarded endpoint receives it
- **THEN** it returns the existing logs_local_only denial and the client does not retry automatically




---
### Requirement: Repeated events remain bounded and nonblocking

The system SHALL emit a first fault immediately and coalesce repeated faults by workspace, entity, applied revision, stage and code. Summaries SHALL use the existing runtime-log heartbeat interval, include first and last observation times and an exact occurrence count, and distinguish intentional coalescing from capture loss. A recovery SHALL be a separate event. Coalescer memory and output SHALL respect the existing admission, broker, replay, subscriber and browser count and byte limits.

#### Scenario: Repeated failure and revision change

- **WHEN** one fault repeats within a summary window and the same entity later fails under a different applied revision
- **THEN** the first revision has an accurate bounded summary and the new revision starts a distinct fault
- **AND** a recovery for one revision does not clear the other revision's fault

##### Example: Repeated failure and revision change

- **GIVEN** G/R1 faults three times in one heartbeat window and G/R2 faults once
- **WHEN** the coalescer emits its summary
- **THEN** R1 reports occurrences=3 as a decimal string and R2 remains an independent first fault

#### Scenario: High-cardinality burst

- **WHEN** distinct event scopes exceed the existing admission count or byte capacity
- **THEN** capture remains nonblocking, bounded and loss-accounted
- **AND** it does not allocate an unbounded coalescing map or emit per-sample success logs

##### Example: High-cardinality burst

- **GIVEN** the coalescer has reached its existing AdmissionRecords limit
- **WHEN** another distinct entity key arrives
- **THEN** memory remains capped and any rejected capture is reflected in loss accounting

#### Scenario: Restart after a previous commit

- **WHEN** the process restarts with historical receipts
- **THEN** it does not replay them as newly committed events
- **AND** the viewer identifies the process reset and links to current durable evidence

##### Example: Restart after a previous commit

- **GIVEN** process I1 committed receipt E before shutdown
- **WHEN** process I2 restores delivery state
- **THEN** the viewer receives a process reset but E is not emitted as a new commitment




---
### Requirement: Actionable filters retain event and transport identity

The system SHALL support applicable workspace, device, rule, group, stage and actionable filters alongside existing level, source and literal search. Snapshot, replay and live delivery SHALL apply the same filters and preserve string cursor precision, bounded retention and gap semantics. Missing recovery evidence SHALL NOT establish a healthy state.

#### Scenario: Filters change with an old response in flight

- **WHEN** the operator switches from device A to device B while A's response or stream callback arrives
- **THEN** B's view ignores A's late data and uses B's filter-bound replay progress
- **AND** pause, resume and local clear retain their existing server-retention semantics

##### Example: Filters change with an old response in flight

- **GIVEN** device A filter has cursor sequence 9007199254740993 and device B is selected
- **WHEN** A returns after B connects
- **THEN** A rows are ignored and its cursor is not rounded or reused for B progress

#### Scenario: Retention gap prevents complete fault history

- **WHEN** the viewer has a retention, capture or process-reset gap
- **THEN** the actionable view identifies the uncertainty and offers current rule evidence
- **AND** absence of a retained fault or recovery event is not presented as healthy

##### Example: Retention gap prevents complete fault history

- **GIVEN** a fault is retained but its recovery lies outside the retained history
- **WHEN** the actionable filter is selected
- **THEN** the event is unresolved or unknown with a gap notice rather than asserted healthy




---
### Requirement: Events navigate to exact settings without mutations

The system SHALL render a safe problem description, affected identity, stage, first and last occurrence, count and typed next action. It SHALL use the EntityLocator and runtime-rule resolver owned by expose-runtime-recording-rules. Historical event revisions SHALL remain immutable and event actions SHALL only navigate, not execute retry, skip, activation or DDL.

#### Scenario: Operator opens the exact failed group

- **WHEN** the operator follows a delivery event action
- **THEN** the runtime rule or Studio owning step selects the matching workspace and group identity
- **AND** the original event revision and current configuration are visibly distinguished

##### Example: Operator opens the exact failed group

- **GIVEN** an event identifies workspace W, group G and applied revision R1 while draft R2 exists
- **WHEN** the operator opens its action
- **THEN** the resolver selects W/G and retains R1 as event context while identifying R2 separately

#### Scenario: Configuration was deleted or renamed

- **WHEN** an event references an entity that has been deleted, renamed or moved outside the available scope
- **THEN** the event retains its verified historical identity and offers a safe unavailable explanation or workspace list
- **AND** it does not relabel history with a new draft or choose a same-name replacement

##### Example: Configuration was deleted or renamed

- **GIVEN** event E references deleted G1 and active G2 has G1 old name
- **WHEN** E is expanded and its action followed
- **THEN** E retains G1 identity and shows a safe unavailable result rather than opening G2




---
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
